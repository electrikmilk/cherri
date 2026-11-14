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

const (
	lzvnBlockMagic        = "bvxn"
	lzfseEndOfStreamMagic = "bvx$"
	lzvnBlockHeaderSize   = 12
	lzvnEndOfStreamOp     = 0x06
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
		var stream = signed[len(signed)-remainingBefore:]
		// Apple compresses small inputs as a single LZVN block. lzfse v0.2.4 (the latest) stores LZVN literal
		// and match lengths in a byte, so a length above 255 wraps, the decoder loses its place in the block and fails
		// with "Bad magic". Those streams are decoded here instead.
		if string(stream[:min(len(lzvnBlockMagic), len(stream))]) == lzvnBlockMagic {
			var block, consumed, lzvnErr = decodeLZVNStream(stream)
			if lzvnErr != nil {
				return nil, fmt.Errorf("failed to decompress signed Shortcut: %w", lzvnErr)
			}
			if _, seekErr := compressed.Seek(int64(consumed), io.SeekCurrent); seekErr != nil {
				return nil, fmt.Errorf("failed to decompress signed Shortcut: %w", seekErr)
			}
			archive.Write(block)
			continue
		}

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

// decodeLZVNStream decodes an LZFSE stream made of one LZVN block and the end of stream marker, which is how Apple's
// encoder writes small inputs, and returns the decoded data and the number of stream bytes it used.
func decodeLZVNStream(stream []byte) ([]byte, int, error) {
	if len(stream) < lzvnBlockHeaderSize {
		return nil, 0, errors.New("LZVN block header is truncated")
	}

	var rawSize = binary.LittleEndian.Uint32(stream[4:])
	var payloadEnd = uint64(lzvnBlockHeaderSize) + uint64(binary.LittleEndian.Uint32(stream[8:]))
	var streamEnd = payloadEnd + uint64(len(lzfseEndOfStreamMagic))
	if streamEnd > uint64(len(stream)) {
		return nil, 0, fmt.Errorf("LZVN block of %d bytes is truncated, %d bytes left", streamEnd, len(stream))
	}
	if string(stream[payloadEnd:streamEnd]) != lzfseEndOfStreamMagic {
		return nil, 0, fmt.Errorf("expected end of LZFSE stream after LZVN block, found %q", stream[payloadEnd:streamEnd])
	}

	var decoded, decodeErr = decodeLZVN(stream[lzvnBlockHeaderSize:payloadEnd], int(rawSize))
	if decodeErr != nil {
		return nil, 0, decodeErr
	}

	return decoded, int(streamEnd), nil
}

// decodeLZVN decodes an LZVN payload, following the opcode layout of Apple's reference lzvn_decode_base.c.
func decodeLZVN(payload []byte, rawSize int) ([]byte, error) {
	var decoded = make([]byte, 0, rawSize)
	var distance int
	for pos := 0; ; {
		if pos >= len(payload) {
			return nil, errors.New("LZVN payload ended without an end of stream opcode")
		}

		var op = payload[pos]
		var literalLength, matchLength, opSize int
		switch {
		case op == lzvnEndOfStreamOp:
			if len(decoded) != rawSize {
				return nil, fmt.Errorf("LZVN block decoded to %d bytes, expected %d", len(decoded), rawSize)
			}
			return decoded, nil
		case op == 0x0e || op == 0x16:
			pos++
			continue
		case op&0xf0 == 0x70 || op&0xf0 == 0xd0 || (op < 0x40 && op&7 == 6):
			return nil, fmt.Errorf("invalid LZVN opcode 0x%02x at offset %d", op, pos)
		case op == 0xe0:
			opSize = 2
		case op&0xf0 == 0xe0:
			opSize = 1
			literalLength = int(op & 0xf)
		case op == 0xf0:
			opSize = 2
		case op&0xf0 == 0xf0:
			opSize = 1
			matchLength = int(op & 0xf)
		case op&0xe0 == 0xa0, op&7 == 7:
			opSize = 3
		case op&7 == 6:
			opSize = 1
		default:
			opSize = 2
		}
		if pos+opSize > len(payload) {
			return nil, fmt.Errorf("LZVN opcode 0x%02x at offset %d is truncated", op, pos)
		}

		var operands = payload[pos+1 : pos+opSize]
		switch {
		case op == 0xe0:
			literalLength = int(operands[0]) + 16
		case op == 0xf0:
			matchLength = int(operands[0]) + 16
		case op >= 0xe0:
		case op&0xe0 == 0xa0:
			var packed = binary.LittleEndian.Uint16(operands)
			literalLength = int(op>>3) & 3
			matchLength = (int(op&7)<<2 | int(packed&3)) + 3
			distance = int(packed >> 2)
		default:
			literalLength = int(op >> 6)
			matchLength = int(op>>3&7) + 3
			if op&7 == 7 {
				distance = int(binary.LittleEndian.Uint16(operands))
			} else if op&7 != 6 {
				distance = int(op&7)<<8 | int(operands[0])
			}
		}
		pos += opSize

		if pos+literalLength > len(payload) {
			return nil, fmt.Errorf("LZVN literal of %d bytes at offset %d is truncated", literalLength, pos)
		}
		decoded = append(decoded, payload[pos:pos+literalLength]...)
		pos += literalLength

		if matchLength == 0 {
			continue
		}
		if distance == 0 || distance > len(decoded) {
			return nil, fmt.Errorf("LZVN match distance %d at offset %d is outside the %d decoded bytes", distance, pos, len(decoded))
		}
		// Byte by byte because a match may overlap the bytes it is copying.
		var matchStart = len(decoded) - distance
		for i := 0; i < matchLength; i++ {
			decoded = append(decoded, decoded[matchStart+i])
		}
	}
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
