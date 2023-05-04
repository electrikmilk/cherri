/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// include is a data structure to track include statements.
type include struct {
	file   string
	start  int
	end    int
	lines  []string
	length int
}

var includes []include

// parseIncludes() searches for include statements within the file and
// injects the contents of the file at the specified path.
func parseIncludes() {
	lines = strings.Split(contents, "\n")
	for l, line := range lines {
		var lineChars = strings.Split(line, "")
		if len(lineChars) == 0 {
			continue
		}
		if !startsWith(strings.Trim(line, " "), "#include") {
			collectUntil('\n')
			continue
		}

		// Prepare for possible error
		chars = strings.Split(line, "")
		lineIdx = l
		idx = len("#include") + 1
		lineCharIdx = idx

		r := regexp.MustCompile("\"(.*?)\"")
		var includePath = strings.Trim(r.FindString(line), "\"")
		if includePath == "" {
			parserError("Expected file path")
		}

		if contains(included, includePath) {
			parserError(fmt.Sprintf("File '%s' has already been included.", includePath))
		}

		var includeFileBytes []byte
		var includeReadErr error
		if includePath == "stdlib" {
			includeFileBytes, includeReadErr = stdLib.ReadFile("stdlib.cherri")
		} else {
			if !strings.Contains(includePath, "..") {
				includePath = relativePath + includePath
			}
			checkFile(includePath)
			includeFileBytes, includeReadErr = os.ReadFile(includePath)
		}
		handle(includeReadErr)

		var includeContents = string(includeFileBytes)
		var includeLines = strings.Split(includeContents, "\n")
		var includeLinesCount = len(includeLines)
		if strings.Contains(includeContents, "#include") {
			includeLinesCount--
		}

		updateIncludesMap(l, includeLinesCount)

		includes = append(includes, include{
			file:   includePath,
			start:  l,
			end:    l + includeLinesCount,
			lines:  includeLines,
			length: includeLinesCount,
		})

		lines[l] = includeContents
		included = append(included, includePath)
	}
	contents = strings.Join(lines, "\n")
	lines = strings.Split(contents, "\n")
	lineIdx = 0
	if strings.Contains(contents, "#include") {
		parseIncludes()
	}
}

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
