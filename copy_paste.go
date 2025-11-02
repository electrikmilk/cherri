/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/electrikmilk/args-parser"
)

var pasteables map[string]string
var copyPasteRegex = regexp.MustCompile(`(copy )?([\W_]+)\{`)

func handleCopyPastes() {
	var matches = copyPasteRegex.FindAllStringSubmatch(contents, -1)
	if len(matches) == 0 {
		return
	}

	parseCopyPastes()

	if args.Using("debug") {
		printPasteablesDebug()
	}
}

func parseCopyPastes() {
	pasteables = make(map[string]string)
	for char != -1 {
		switch {
		case char == '"':
			collectString()
			advanceUntil('\n')
		case commentAhead():
			collectComment()
		case startOfLineTokenAhead(Copy):
			advance()
			collectCopy()
		case startOfLineTokenAhead(Paste):
			advance()
			pasteCopy()
		}
		advance()
	}

	resetParse()
}

func collectCopy() {
	var lineRef = newLineReference()
	var identifier = collectIdentifier()

	if _, found := pasteables[identifier]; found {
		parserError(fmt.Sprintf("Duplication declaration of copy/paste '%s'", identifier))
	}

	advanceUntil('{')
	advance()
	var contents = collectObject()

	lineRef.replaceLines()

	pasteables[identifier] = strings.TrimSpace(contents)
}

func pasteCopy() {
	var identifier = collectIdentifier()
	if contents, found := pasteables[identifier]; found {
		lines[lineIdx] = contents
	} else {
		parserError(fmt.Sprintf("Unable to paste undefined copy '%s'", identifier))
	}
}

func printPasteablesDebug() {
	fmt.Println(ansi("### COPY/PASTE ###", bold) + "\n")
	for identifier, contents := range pasteables {
		fmt.Println("identifier:", identifier+"()")
		fmt.Println("contents:")
		fmt.Println(contents)
		fmt.Println("(end)")
		fmt.Print("\n")
	}
}
