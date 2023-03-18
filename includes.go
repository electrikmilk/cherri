/*
 * Copyright (c) 2023 Brandon Jordan
 */

package main

import (
	"fmt"
	"os"
	"strings"
)

type include struct {
	file   string
	start  int
	end    int
	lines  []string
	length int
}

var includes []include

// updateIncludesMap checks if an included file starts on `line`.
// If so, it updates its start and end lines to account for the included file it overlaps with.
func updateIncludesMap(line int, includeLines int) {
	for i, inc := range includes {
		if inc.start == line {
			includes[i].start = inc.start + includeLines
			includes[i].end = (inc.start + includeLines) + inc.length
		}
	}
}

// delinquentFile checks if the current parsing error occurs on an included file.
func delinquentFile() (errorFilename string, errorLine int, errorCol int) {
	errorFilename = filePath
	errorLine = lineIdx + 1
	errorCol = lineCharIdx + 1
	for _, inc := range includes {
		if lineIdx+1 >= inc.start && lineIdx+1 <= inc.end {
			errorFilename = inc.file
			for l, line := range inc.lines {
				if lineIdx < len(lines) {
					if line == lines[lineIdx] {
						errorLine = l + 1
					}
				}
			}
		}
	}
	return
}

// autoInclude looks for files in the same directory and automatically includes them.
func autoInclude() {
	var files, readDirErr = os.ReadDir(".")
	handle(readDirErr)
	if len(files) == 0 {
		return
	}
	var cherriFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == os.Args[1] {
			continue
		}
		var ext = end(strings.Split(file.Name(), "."))
		if ext == "cherri" {
			cherriFiles = append(cherriFiles, file.Name())
		}
	}
	if len(cherriFiles) == 0 {
		return
	}
	lines = strings.Split(contents, "\n")
	for _, cherriFile := range cherriFiles {
		var includeLine = fmt.Sprintf("#include \"%s\"", cherriFile)
		lines = append([]string{includeLine}, lines...)
	}
	contents = strings.Join(lines, "\n")
}
