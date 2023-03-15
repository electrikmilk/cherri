/*
 * Copyright (c) 2023 Brandon Jordan
 */

package main

import (
	"fmt"
	"os"
	"regexp"
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

func parseIncludes() {
	lines = strings.Split(contents, "\n")
	for l, line := range lines {
		var lineChars = strings.Split(line, "")
		if len(lineChars) == 0 {
			continue
		}
		if !startsWith(line, "#include") {
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

		if !strings.Contains(includePath, "..") {
			includePath = relativePath + includePath
		}

		if contains(included, includePath) {
			parserError(fmt.Sprintf("File '%s' has already been included.", includePath))
		}

		checkFile(includePath)
		var includeFileBytes, readErr = os.ReadFile(includePath)
		handle(readErr)

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
