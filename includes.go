/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
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

// directory is where the file the line came from lives, which for the main file and generated lines is the main file's directory.
func (origin lineOrigin) directory() string {
	if origin.file == "" {
		return relativePath
	}

	return filepath.Dir(origin.file)
}

var included []string

func handleIncludes() {
	if args.Using("debug") {
		fmt.Println("Parsing includes...")
	}

	for expandIncludeLevel() {
		// Files that were just included can include more files, which the next level expands.
	}

	if args.Using("debug") {
		fmt.Println(ansi("Done.", green) + "\n")
	}
}

// expandIncludeLevel replaces each line holding an include with the lines of the included file, and reports whether it found any.
// The included files are not scanned yet, so the files they include are found by the next level.
func expandIncludeLevel() bool {
	var expandedLines = make([]string, 0, len(lines))
	var expandedOrigins = make([]lineOrigin, 0, len(lineOrigins))
	var nextUnexpandedLine int
	var foundInclude bool
	for char != -1 {
		switch {
		case char == '"':
			collectString()
			advanceUntil('\n')
		case commentAhead():
			collectComment()
		case tokenAhead(Include):
			advance()
			var includeLineIdx = lineIdx
			var includeLines, includeOrigins = expandInclude()

			expandedLines = append(expandedLines, lines[nextUnexpandedLine:includeLineIdx]...)
			expandedOrigins = append(expandedOrigins, lineOrigins[nextUnexpandedLine:includeLineIdx]...)
			expandedLines = append(expandedLines, includeLines...)
			expandedOrigins = append(expandedOrigins, includeOrigins...)
			nextUnexpandedLine = includeLineIdx + 1
			foundInclude = true

			// The whole line is replaced by the included file, so the rest of it is never scanned.
			advanceUntil('\n')
		}
		advance()
	}

	if foundInclude {
		lines = append(expandedLines, lines[nextUnexpandedLine:]...)
		lineOrigins = append(expandedOrigins, lineOrigins[nextUnexpandedLine:]...)
	}
	resetParse()

	return foundInclude
}

// expandInclude collects the include path at the cursor, checks it and returns the lines of the included file.
func expandInclude() ([]string, []lineOrigin) {
	if char == '"' {
		parserError("Use raw string (') for include file paths")
	}
	if char != '\'' {
		parserError("Expected file path")
	}
	advance()

	var includePath = collectIncludePath()

	var includer = lineOrigins[lineIdx]
	var includerDirectory = includer.directory()
	var packageDirectory, packageSignature, includerInPackage = packageDirectoryOf(includerDirectory)

	// A package has no standard library of its own, so inside one an actions path is an ordinary path.
	var standardLibraryInclude = !includerInPackage && isStandardLibraryInclude(includePath)
	var resolvedPath = includePath
	if !standardLibraryInclude {
		resolvedPath = resolveIncludePath(includePath, includerDirectory)
	}

	if includerInPackage && !isInsideDirectory(resolvedPath, packageDirectory) {
		parserError(fmt.Sprintf("Include '%s' resolves to '%s', outside package '%s'. Packages are isolated and cannot include files outside their own directory.", includePath, resolvedPath, packageSignature))
	}

	if alreadyIncluded(resolvedPath) {
		parserError(fmt.Sprintf("Path '%s' has already been included.", resolvedPath))
	}

	var includeSource []byte
	if standardLibraryInclude {
		includeSource = readStandardLibrarySource(includePath)
	} else {
		includeSource = readIncludeSource(resolvedPath)
	}
	included = append(included, resolvedPath)

	var includeLines = strings.Split(string(includeSource), "\n")

	return includeLines, originsFor(resolvedPath, len(includeLines))
}

// alreadyIncluded is also true for the main file, so a file that includes it back is a repeated include like any cycle.
func alreadyIncluded(resolvedPath string) bool {
	if filePath != "" && resolvedPath == filepath.Clean(filePath) {
		return true
	}

	return slices.Contains(included, resolvedPath)
}

// packageDirectoryOf finds the installed package that contains directory.
func packageDirectoryOf(directory string) (packageDirectory string, packageSignature string, inPackage bool) {
	var packagesDirectory = filepath.Join(relativePath, "packages")
	var pathInPackages, relErr = filepath.Rel(packagesDirectory, filepath.Clean(directory))
	if relErr != nil {
		return
	}

	// A directory in a package starts with the user and the package name.
	var pathParts = strings.Split(pathInPackages, string(filepath.Separator))
	if len(pathParts) < 2 || !strings.HasPrefix(pathParts[0], "@") {
		return
	}

	packageSignature = pathParts[0] + "/" + pathParts[1]
	packageDirectory = filepath.Join(packagesDirectory, pathParts[0], pathParts[1])
	inPackage = true

	return
}

// isInsideDirectory compares by path elements so a sibling that shares a name prefix is outside.
func isInsideDirectory(path string, directory string) bool {
	var pathInDirectory, relErr = filepath.Rel(directory, path)
	if relErr != nil {
		return false
	}

	return pathInDirectory != ".." && !strings.HasPrefix(pathInDirectory, ".."+string(filepath.Separator))
}

func isStandardLibraryInclude(includePath string) bool {
	return isStandardActionsInclude(includePath) || includePath == "stdfunc"
}

func isStandardActionsInclude(includePath string) bool {
	// The slash is part of the prefix so a project folder like actionsKit/ is not mistaken for the standard library.
	return startsWith("actions/", includePath)
}

// resolveIncludePath resolves includePath against the directory of the file that includes it.
func resolveIncludePath(includePath string, includerDirectory string) string {
	if filepath.IsAbs(includePath) {
		return filepath.Clean(includePath)
	}

	return filepath.Clean(filepath.Join(includerDirectory, includePath))
}

func readStandardLibrarySource(includePath string) []byte {
	if includePath == "stdfunc" {
		var stdfuncSource, readErr = stdFuncs.ReadFile("stdfunc.cherri")
		handle(readErr)

		return stdfuncSource
	}

	var actionCategory = end(strings.Split(includePath, "/"))
	var actionsSource, readErr = stdActions.ReadFile(fmt.Sprintf("actions/%s.cherri", actionCategory))
	if readErr != nil {
		parserError(fmt.Sprintf("Undefined actions include '%s'.", actionCategory))
	}

	return actionsSource
}

func readIncludeSource(resolvedPath string) []byte {
	if _, statErr := os.Stat(resolvedPath); os.IsNotExist(statErr) {
		parserError(fmt.Sprintf("File '%s' does not exist!", resolvedPath))
	}
	checkFile(resolvedPath)

	var fileSource, readErr = os.ReadFile(resolvedPath)
	handle(readErr)

	return fileSource
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
