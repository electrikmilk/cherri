/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type iCloudRecord struct {
	Fields RecordFields
}

type RecordFields struct {
	Shortcut ShortcutRecord
	Name     ShortcutName
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

func downloadShortcut() []byte {
	importPath = strings.Replace(importPath, "/shortcuts/", "/shortcuts/api/records/", 1)

	fmt.Println("Retrieving record from iCloud...")

	var apiResponse, apiErr = http.Get(importPath)
	handle(apiErr)
	defer apiResponse.Body.Close()
	if apiResponse.StatusCode != http.StatusOK {
		exit(fmt.Sprintf("icloud: Failed to get Shortcut file URL from API."))
	}

	var record iCloudRecord
	var decodeErr = json.NewDecoder(apiResponse.Body).Decode(&record)
	handle(decodeErr)

	var downloadURL = record.Fields.Shortcut.Value.DownloadURL
	filename = record.Fields.Name.Value
	record = iCloudRecord{}

	fmt.Println("Downloading Shortcut...")

	var response, httpErr = http.Get(downloadURL)
	handle(httpErr)
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		exit(fmt.Sprintf("icloud: Failed to download Shortcut file."))
	}

	var b, readErr = io.ReadAll(response.Body)
	handle(readErr)

	return b
}

func importShortcut() []byte {
	var _, statErr = os.Stat(importPath)
	if os.IsNotExist(statErr) {
		exit("import: File does not exist!")
	}

	var segments = strings.Split(importPath, "/")
	filename = segments[len(segments)-1]
	var nameSegments = strings.Split(filename, ".")
	var extension = nameSegments[len(nameSegments)-1]
	if extension != "shortcut" {
		exit(fmt.Sprintf("import: File is not a Shortcut."))
	}

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
