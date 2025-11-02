/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/electrikmilk/args-parser"
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
var included []string

func handleIncludes() {
	if args.Using("debug") {
		fmt.Println("Parsing includes...")
	}

	parseIncludes()

	if args.Using("debug") {
		fmt.Println(ansi("Done.", green) + "\n")
	}
}

var includedFile bool

// parseIncludes() searches for include statements within the file and
// injects the contents of the file at the specified path.
func parseIncludes() {
	includedFile = false
	for char != -1 {
		switch {
		case char == '"':
			collectString()
			advanceUntil('\n')
		case commentAhead():
			collectComment()
		case tokenAhead(Include):
			advance()
			parseInclude()
			resetParse()
			includedFile = true
		}
		advance()
	}

	resetParse()

	if includedFile {
		parseIncludes()
	}
}

func collectIncludePath() string {
	var collection strings.Builder
	for char != -1 {
		if char == '\'' {
			break
		}

		collection.WriteRune(char)
		advance()
	}

	return collection.String()
}

func parseInclude() {
	if char == '"' {
		parserError("Use raw string (') for include file paths")
	}
	if char != '\'' {
		parserError("Expected file path")
	}
	advance()

	var includePath = collectIncludePath()

	if slices.Contains(included, includePath) {
		parserError(fmt.Sprintf("Path '%s' has already been included.", includePath))
	}

	var includeFileBytes []byte
	var includeReadErr error
	if startsWith("actions", includePath) && strings.Contains(includePath, "/") {
		var actionCat = end(strings.Split(includePath, "/"))
		includeFileBytes, includeReadErr = stdActions.ReadFile(fmt.Sprintf("actions/%s.cherri", actionCat))
		if includeReadErr != nil {
			parserError(fmt.Sprintf("Undefined actions include '%s'.", actionCat))
		}
	} else if includePath == "stdlib" {
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

	updateIncludesMap(lineIdx, includeLinesCount)

	includes = append(includes, include{
		file:   includePath,
		start:  lineIdx,
		end:    lineIdx + includeLinesCount,
		lines:  includeLines,
		length: includeLinesCount,
	})

	lines[lineIdx] = includeContents
	included = append(included, includePath)
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

// delinquentFile determines what file the current cursor exists within in relation to any included files.
func delinquentFile() (errorFilename string, errorLine int, errorCol int) {
	errorFilename = workflowName + ".cherri"
	errorLine = lineIdx + 1
	errorCol = lineCharIdx + 1
	if len(includes) == 0 {
		return
	}

	var currentLine = lines[lineIdx]
	var found bool
	for _, inc := range includes {
		if errorLine <= inc.start || errorLine >= inc.end {
			continue
		}
		errorFilename = inc.file

		for l, line := range inc.lines {
			if line == currentLine {
				errorLine = l + 1
				found = true
				break
			}
		}
	}

	if !found {
		findOriginalLine(&errorLine)
	}

	return
}

func findOriginalLine(errorLine *int) {
	for l, line := range strings.Split(originalContents, "\n") {
		if line == lines[lineIdx] {
			*errorLine = l
		}
	}
}

// insideInclude returns a boolean based on if we are within an included file with a name that contains needle.
func insideInclude(needle string) bool {
	for _, inc := range includes {
		if !strings.Contains(inc.file, needle) || lineIdx+1 < inc.start || lineIdx+1 > inc.end {
			continue
		}

		return true
	}

	return false
}

func printIncludesDebug() {
	fmt.Println(ansi("### INCLUDES ###", bold) + "\n")

	fmt.Println(ansi("## INCLUDED ##", bold))
	fmt.Println(included)

	fmt.Print("\n")

	fmt.Println(ansi("## INCLUDES MAP ##", bold))
	for i := range lines {
		for _, inc := range includes {
			if inc.start == i {
				fmt.Println(ansi("#include", magenta), inc.file)
				fmt.Println(ansi(fmt.Sprintf("%d |", inc.start+1), cyan))
				fmt.Print(ansi("...", dim))
				fmt.Print(ansi(fmt.Sprintf("%d lines", inc.length), magenta))
				fmt.Println(ansi("...", dim))
				fmt.Println(ansi(fmt.Sprintf("%d |", inc.end+1), cyan))
				fmt.Print("\n")
			}
		}
	}

	fmt.Print("\n")
}
