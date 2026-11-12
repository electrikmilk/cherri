/*
 * Copyright (c) Cherri
 */

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/aixiansheng/lzfse"
	"github.com/electrikmilk/args-parser"
)

type iCloudRecord struct {
	Fields RecordFields
}

type RecordFields struct {
	Name     ShortcutName
	Shortcut ShortcutRecord
}

type ShortcutName struct {
	Value string
}

type ShortcutRecord struct {
	Value ShortcutRecordValue
}

type ShortcutRecordValue struct {
	DownloadURL string
}

var importPath string

// Imports a Shortcut for decompilation based on path given.
func importShortcut(path string) (shortcutBytes []byte) {
	importPath = path
	var icloudURLRegex = regexp.MustCompile(`^https://(?:www.)?icloud\.com/shortcuts/.+$`)
	if icloudURLRegex.MatchString(importPath) {
		shortcutBytes = downloadShortcut()
	} else {
		shortcutBytes = readShortcutFile()
	}

	if looksLikeSignedShortcut(shortcutBytes) {
		var extractErr error
		shortcutBytes, extractErr = extractSignedShortcut(shortcutBytes)
		if extractErr != nil {
			exit(fmt.Sprintf("import: %s", extractErr))
		}
	}

	return
}

func downloadShortcut() []byte {
	importPath = strings.Replace(importPath, "/shortcuts/", "/shortcuts/api/records/", 1)

	var apiResponse, apiErr = http.Get(importPath)
	handle(apiErr)
	defer apiResponse.Body.Close()
	if apiResponse.StatusCode != http.StatusOK {
		exit("import: Failed to get Shortcut file URL from iCloud Shortcuts API.")
	}

	var record iCloudRecord
	var decodeErr = json.NewDecoder(apiResponse.Body).Decode(&record)
	handle(decodeErr)

	var downloadURL = record.Fields.Shortcut.Value.DownloadURL
	basename = record.Fields.Name.Value
	record = iCloudRecord{}

	var response, httpErr = http.Get(downloadURL)
	handle(httpErr)
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		exit("icloud: Failed to download Shortcut file :(")
	}

	var b, readErr = io.ReadAll(response.Body)
	handle(readErr)

	if args.Using("debug") {
		writeErr := os.WriteFile(basename+"_decompile.plist", b, 0600)
		handle(writeErr)
	}

	return b
}

func readShortcutFile() []byte {
	var _, statErr = os.Stat(importPath)
	if os.IsNotExist(statErr) {
		exit("import: File does not exist!")
	}

	var segments = strings.Split(importPath, "/")
	filename = segments[len(segments)-1]

	var nameSegments = strings.Split(filename, ".")
	basename = nameSegments[0]

	var extension = nameSegments[len(nameSegments)-1]
	if extension != "shortcut" && extension != "plist" {
		exit("import: File is not a Shortcut or property list (plist) file.")
	}

	relativePath = strings.Replace(importPath, filename, "", 1)

	var b, readErr = os.ReadFile(importPath)
	handle(readErr)

	return b
}

const (
	authDataSizeOffset = 8
	// Signed Shortcuts are Apple Encrypted Archives without encryption. The signature, root header and cluster
	// headers between the authentication data and the compressed data have a constant size for Shortcuts.
	// https://theapplewiki.com/wiki/Apple_Encrypted_Archive
	compressedDataOffsetFromAuthDataSize = 0x495c
	lzfseMagicPrefix                     = "bvx"
	appleArchiveEntryMagic               = "AA01"
	appleArchiveEntryPrefixSize          = 6
	shortcutArchivePath                  = "Shortcut.wflow"
)

// archiveFieldLengths are the sizes of the Apple Archive field types a Shortcut archive is made of.
var archiveFieldLengths = map[byte]int{'1': 1, '2': 2, '4': 4, '8': 8, 'T': 12}

type archiveEntry struct {
	path string
	data []byte
	size int
}

// extractSignedShortcut returns the property list inside of a signed Shortcut. No key is needed because Shortcuts
// are only signed, not encrypted.
func extractSignedShortcut(signed []byte) ([]byte, error) {
	var archive, decompressErr = decompressSignedShortcut(signed)
	if decompressErr != nil {
		return nil, decompressErr
	}

	return findArchiveFile(archive, shortcutArchivePath)
}

func decompressSignedShortcut(signed []byte) ([]byte, error) {
	if len(signed) < authDataSizeOffset+4 {
		return nil, errors.New("signed Shortcut is too short")
	}

	var authDataSize = binary.LittleEndian.Uint32(signed[authDataSizeOffset:])
	var compressedStart = uint64(authDataSize) + compressedDataOffsetFromAuthDataSize
	if compressedStart+uint64(len(lzfseMagicPrefix)) > uint64(len(signed)) ||
		string(signed[compressedStart:compressedStart+uint64(len(lzfseMagicPrefix))]) != lzfseMagicPrefix {
		return nil, errors.New("unsupported signed Shortcut layout, compressed data not found where expected")
	}

	// Each block of up to 1 MiB of the archive is its own LZFSE stream, and the decompressor stops at the end of one.
	var compressed = bytes.NewReader(signed[compressedStart:])
	var archive bytes.Buffer
	for compressed.Len() > 0 {
		var remainingBefore = compressed.Len()
		var block, decompressErr = io.ReadAll(lzfse.NewReader(compressed))
		if decompressErr != nil {
			return nil, fmt.Errorf("failed to decompress signed Shortcut: %w", decompressErr)
		}
		if compressed.Len() == remainingBefore {
			return nil, errors.New("failed to decompress signed Shortcut, no data consumed")
		}
		archive.Write(block)
	}

	return archive.Bytes(), nil
}

func findArchiveFile(archive []byte, path string) ([]byte, error) {
	for offset := 0; offset < len(archive); {
		var entry, entryErr = parseArchiveEntry(archive[offset:])
		if entryErr != nil {
			return nil, entryErr
		}
		if entry.path == path {
			return entry.data, nil
		}
		offset += entry.size
	}

	return nil, fmt.Errorf("%s not found in signed Shortcut", path)
}

// parseArchiveEntry reads an Apple Archive entry, which is a header of typed fields followed by the entry's data.
func parseArchiveEntry(archive []byte) (entry archiveEntry, err error) {
	if len(archive) < appleArchiveEntryPrefixSize || string(archive[:4]) != appleArchiveEntryMagic {
		return entry, errors.New("signed Shortcut contains an invalid archive entry")
	}

	var headerSize = int(binary.LittleEndian.Uint16(archive[4:]))
	if headerSize < appleArchiveEntryPrefixSize || headerSize > len(archive) {
		return entry, errors.New("signed Shortcut contains an archive entry with an invalid header size")
	}

	var dataSize int
	var fields = archive[appleArchiveEntryPrefixSize:headerSize]
	for len(fields) > 0 {
		if len(fields) < 4 {
			return entry, errors.New("signed Shortcut contains a truncated archive field")
		}

		var key = string(fields[:4])
		fields = fields[4:]

		var fieldLength, lengthErr = archiveFieldLength(key[3], fields)
		if lengthErr != nil {
			return entry, fmt.Errorf("%w (%s)", lengthErr, key)
		}
		var value = fields[:fieldLength]
		fields = fields[fieldLength:]

		switch key {
		case "PATP":
			entry.path = string(value[2:])
		case "DATA", "DATB":
			var paddedSize [8]byte
			copy(paddedSize[:], value)
			dataSize = int(binary.LittleEndian.Uint64(paddedSize[:]))
		}
	}

	if headerSize+dataSize > len(archive) {
		return entry, errors.New("signed Shortcut contains an archive entry with truncated data")
	}

	entry.data = archive[headerSize : headerSize+dataSize]
	entry.size = headerSize + dataSize

	return entry, nil
}

// archiveFieldLength returns the length of a field's value, which comes after its 4 character key.
func archiveFieldLength(fieldType byte, value []byte) (int, error) {
	var length int
	switch fieldType {
	case 'P':
		if len(value) < 2 {
			return 0, errors.New("truncated archive field")
		}
		length = 2 + int(binary.LittleEndian.Uint16(value))
	case 'A':
		length = 2
	case 'B':
		length = 4
	default:
		var fixedLength, known = archiveFieldLengths[fieldType]
		if !known {
			return 0, errors.New("unsupported archive field type")
		}
		length = fixedLength
	}

	if length > len(value) {
		return 0, errors.New("truncated archive field")
	}

	return length, nil
}
