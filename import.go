/*
 * Copyright (c) Cherri
 */

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

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
func importShortcut() (shortcutBytes []byte) {
	importPath = args.Value("import")

	var icloudURLRegex = regexp.MustCompile(`^https://(?:www.)?icloud\.com/shortcuts/.+$`)
	if icloudURLRegex.MatchString(importPath) {
		shortcutBytes = downloadShortcut()
	} else {
		shortcutBytes = readShortcutFile()
	}

	if hasSignedBytes(shortcutBytes) {
		exit("import: Signed Shortcuts are currently not supported :(\nYou can use an iCloud link instead by sharing the Shortcut and selecting \"Copy iCloud Link\".")
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

func hasSignedBytes(b []byte) bool {
	var rawUnsignedBytes = []byte{98, 112, 108, 105, 115, 116, 48, 48}
	var unsignedBytes = []byte{60, 63, 120, 109, 108, 32, 118, 101}
	for i, ub := range unsignedBytes {
		if ub == b[i] || b[i] == rawUnsignedBytes[i] {
			continue
		}

		return true
	}

	return false
}
