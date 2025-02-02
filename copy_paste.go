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
var copyPasteRegex = regexp.MustCompile(`copy .+ \{`)

func parseCopyPastes() {
	if !copyPasteRegex.MatchString(contents) {
		return
	}

	pasteables = make(map[string]string)
	for char != -1 {
		switch {
		case isChar('/'):
			collectComment()
		case tokenAhead(Copy):
			advance()
			collectCopy()
		case tokenAhead(Paste):
			advance()
			pasteCopy()
			continue
		}
		advance()
	}

	resetParse()

	if args.Using("debug") {
		printPasteablesDebug()
	}
}

func collectCopy() {
	var startLine = lineIdx
	var identifier = collectIdentifier()

	if _, found := pasteables[identifier]; found {
		parserError(fmt.Sprintf("Duplication declaration of copy/paste '%s'", identifier))
	}

	advanceUntil('{')
	advance()
	var contents = collectObject()

	for i := startLine; i <= lineIdx; i++ {
		lines[i] = ""
	}

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
