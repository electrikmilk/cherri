/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/electrikmilk/args-parser"
)

// lineOrigin is where a line of the combined source came from.
// The main file has an empty file, lines the compiler generated have a line of 0.
type lineOrigin struct {
	file string
	line int
}

// lineOrigins has an origin for each of the lines.
var lineOrigins []lineOrigin

// originsFor numbers count lines of file from 1.
func originsFor(file string, count int) []lineOrigin {
	var origins = make([]lineOrigin, count)
	for i := range origins {
		origins[i] = lineOrigin{file: file, line: i + 1}
	}

	return origins
}

func (origin lineOrigin) source() string {
	if origin.line == 0 {
		return "(generated)"
	}
	if origin.file == "" {
		return "(main file)"
	}

	return origin.file
}

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

	handlePackageIncludePath(&includePath)

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
	} else if includePath == "stdfunc" {
		includeFileBytes, includeReadErr = stdFuncs.ReadFile("stdfunc.cherri")
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

	lines = slices.Replace(lines, lineIdx, lineIdx+1, includeLines...)
	lineOrigins = slices.Replace(lineOrigins, lineIdx, lineIdx+1, originsFor(includePath, includeLinesCount)...)
	included = append(included, includePath)
}

var packageIncludeRegex = regexp.MustCompile("packages/@(.*?)/(.*?)/")

// handlePackageIncludePath resolves internal package includes.
func handlePackageIncludePath(includePath *string) {
	var matches = packageIncludeRegex.FindStringSubmatch(lineOrigins[lineIdx].file)
	if matches == nil {
		return
	}

	*includePath = fmt.Sprintf("packages/@%s/%s/%s", matches[1], matches[2], *includePath)
}

// delinquentFile determines what file and line the current cursor came from.
func delinquentFile() (errorFilename string, errorLine int, errorCol int) {
	errorFilename = filePath
	errorLine = lineIdx + 1
	errorCol = lineCharIdx + 1

	// The source lines are discarded once parsing finishes, so errors while generating have no origin.
	if lineIdx >= len(lineOrigins) {
		return
	}

	var origin = lineOrigins[lineIdx]
	if origin.line == 0 {
		return
	}

	errorLine = origin.line
	if origin.file != "" {
		errorFilename = origin.file
	}

	return
}

func printIncludesDebug() {
	fmt.Println(ansi("### INCLUDES ###", bold) + "\n")

	fmt.Println(ansi("## INCLUDED ##", bold))
	fmt.Println(included)

	fmt.Print("\n")

	fmt.Println(ansi("## LINE ORIGINS ##", bold))
	var runStart int
	for i := 1; i <= len(lineOrigins); i++ {
		if i < len(lineOrigins) && lineOrigins[i].source() == lineOrigins[runStart].source() {
			continue
		}

		fmt.Println(ansi(fmt.Sprintf("%d-%d |", runStart+1, i), cyan), lineOrigins[runStart].source())
		runStart = i
	}

	fmt.Print("\n")
}
