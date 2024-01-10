/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/electrikmilk/args-parser"
	"io"
	"net/http"
	"os"
	"regexp"
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

func downloadShortcut() []byte {
	var icloudURL = args.Value("icloud")

	var icloudURLRegex = regexp.MustCompile(`^https://(?:www.)?icloud\.com/shortcuts/.+$`)
	if !icloudURLRegex.MatchString(icloudURL) {
		exit(fmt.Sprintf("import: iCloud URL `%s` does not match format.", icloudURL))
	}

	icloudURL = strings.Replace(icloudURL, "/shortcuts/", "/shortcuts/api/records/", 1)

	fmt.Println("Retrieving record from iCloud...")

	var apiResponse, apiErr = http.Get(icloudURL)
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
	var path = args.Value("import")
	var _, statErr = os.Stat(path)
	if os.IsNotExist(statErr) {
		exit(fmt.Sprintf("import: File '%s' does not exist!", path))
	}

	var segments = strings.Split(path, "/")
	filename = segments[len(segments)-1]
	var nameSegments = strings.Split(filename, ".")
	var extension = nameSegments[len(nameSegments)-1]
	if extension != "shortcut" {
		exit(fmt.Sprintf("import: File is not a Shortcut."))
	}

	var b, readErr = os.ReadFile(path)
	handle(readErr)

	return b
}

func hasSignedBytes(b []byte) bool {
	var unsignedBytes = []byte{98, 112, 108, 105, 115, 116, 48, 48}
	for i, ub := range unsignedBytes {
		if ub == b[i] {
			continue
		}

		return true
	}

	return false
}
