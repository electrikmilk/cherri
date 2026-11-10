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
	"os/exec"
	"slices"
	"time"

	"github.com/aixiansheng/lzfse"
	"github.com/electrikmilk/args-parser"
	"howett.net/plist"
)

var signFailed = false

// sign runs the shortcuts sign command on the unsigned shortcut file.
func sign() {
	if !darwin {
		useHubSign()
		return
	}

	var signingMode = "people-who-know-me"
	if args.Using("share") && args.Value("share") == "anyone" {
		signingMode = "anyone"
	}

	if args.Using("debug") {
		fmt.Printf("Signing %s to %s...", inputPath, outputPath)
	}
	var sign = exec.Command(
		"shortcuts",
		"sign",
		"-i", inputPath,
		"-o", outputPath,
		"-m", signingMode,
	)
	var stdErr bytes.Buffer
	sign.Stderr = &stdErr
	var signErr = sign.Run()
	if signErr != nil {
		signFailed = true
		if args.Using("debug") {
			fmt.Print(ansi("Failed!\n", red))
		}

		fmt.Printf("%s\n%s\n", ansi("Failed to sign Shortcut using macOS :(", orange, bold), ansi(stdErr.String(), orange))

		useHubSign()
	}

	if args.Using("debug") {
		fmt.Println(ansi("Done.", green))
	}
}

func useHubSign() {
	var hubSignService = hubSign()
	useSigningService(&hubSignService)
}

type SigningService struct {
	name string
	url  string
	info func() string
}

// Sign the Shortcut using a signing service.
func useSigningService(service *SigningService) {
	if !args.Using("no-ansi") {
		fmt.Println(ansi(fmt.Sprintf("Signing using %s service...", service.name), green))
		if service.info != nil {
			fmt.Println(service.info())
		}
	}

	var signedShortcut = requestSignedShortcut(service)

	if !looksLikeSignedShortcut(signedShortcut) {
		exit("Signing server response does not look like a Shortcut file.")
	}

	var writeErr = os.WriteFile(outputPath, signedShortcut, 0600)
	handle(writeErr)

	if args.Using("debug") {
		fmt.Println(ansi("Done.", green))
	}
}

func requestSignedShortcut(service *SigningService) []byte {
	var marshaledPlist, plistErr = plist.Marshal(shortcut, plist.XMLFormat)
	handle(plistErr)

	var payload = map[string]string{
		"shortcutName": basename,
		"shortcut":     string(marshaledPlist),
	}
	var jsonPayload, jsonErr = json.Marshal(payload)
	handle(jsonErr)

	var request, httpErr = http.NewRequest("POST", service.url, bytes.NewReader(jsonPayload))
	handle(httpErr)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", fmt.Sprintf("cherri/%s", version))

	var client = &http.Client{
		Timeout: time.Second * 30,
	}
	var response, resErr = client.Do(request)
	handle(resErr)
	defer response.Body.Close()

	// Checked before the content type because error responses can have any content type.
	if response.StatusCode != http.StatusOK {
		exit(fmt.Sprintf("Failed to sign Shortcut (%s)", response.Status))
	}

	var responseContentType = response.Header.Get("Content-Type")
	var allowedContentTypes = []string{"application/octet-stream", "application/x-plist", "application/x-apple-shortcut"}
	if !slices.Contains(allowedContentTypes, responseContentType) {
		exit(fmt.Sprintf("Unsupported response content type: %s", responseContentType))
	}

	var body, readErr = io.ReadAll(response.Body)
	handle(readErr)

	return body
}

// looksLikeSignedShortcut performs quick checks to make sure response is a signed Shortcut.
func looksLikeSignedShortcut(buffer []byte) bool {
	if len(buffer) >= 4 && string(buffer[:4]) == "AEA1" {
		return true
	}
	return false
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

func removeUnsigned() {
	var _, signedStatErr = os.Stat(fmt.Sprintf("%s%s.shortcut", relativePath, workflowName))
	if os.IsNotExist(signedStatErr) {
		return
	}
	var _, unsignedStatErr = os.Stat(fmt.Sprintf("%s%s%s", relativePath, workflowName, unsignedEnd))
	if os.IsNotExist(unsignedStatErr) {
		return
	}

	if args.Using("debug") {
		fmt.Printf("Removing %s%s...", workflowName, unsignedEnd)
	}

	var removeErr = os.Remove(inputPath)
	handle(removeErr)

	if args.Using("debug") {
		fmt.Println(ansi("Done.", green))
	}
}

// Sign the Shortcut using RoutineHub's signing service.
func hubSign() SigningService {
	return SigningService{
		name: "HubSign",
		url:  "https://hubsign.routinehub.services/sign",
		info: func() string {
			return fmt.Sprintf("Shortcut Signing Powered By %s", ansi("RoutineHub", red))
		},
	}
}
