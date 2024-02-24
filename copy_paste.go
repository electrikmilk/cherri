/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	"regexp"
	"strings"
)

var pastables map[string]string
var copyPasteRegex = regexp.MustCompile(`copy .+ \{`)

func parseCopyPastes() {
	if !copyPasteRegex.MatchString(contents) {
		return
	}

	pastables = make(map[string]string)
	for char != -1 {
		switch {
		case isToken(ForwardSlash):
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
	//fmt.Println(contents)
	//fmt.Println(pastables)
}

func collectCopy() {
	var startLine = lineIdx
	var identifier = collectIdentifier()
	advanceUntil('{')
	advance()
	var contents = collectObject()

	for i := startLine; i <= lineIdx; i++ {
		lines[i] = ""
	}

	pastables[identifier] = strings.TrimSpace(contents)
}

func pasteCopy() {
	var identifier = collectIdentifier()
	if char == '\n' {
		idx--
		lineIdx--
		lineCharIdx = len(lines[lineIdx])
	}
	if contents, found := pastables[identifier]; found {
		lines[lineIdx] = contents
	} else {
		parserError(fmt.Sprintf("Unable to paste undefined copy '%s'", identifier))
	}
}
