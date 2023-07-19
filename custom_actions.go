/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	"strings"

	"github.com/electrikmilk/args-parser"
)

// customAction contains the collected declaration of a custom action.
type customAction struct {
	body string
}

// customActions is a map of all the custom actions that have been defined.
var customActions map[string]customAction

// parseCustomActions parses defined actions and collects them.
func parseCustomActions() {
	if args.Using("debug") {
		fmt.Print(ansi("done!", green) + "\n")
		fmt.Print("Parsing custom actions... ")
	}
	customActions = make(map[string]customAction)
	chars = strings.Split(contents, "")
	idx = -1
	advance()
	for char != -1 {
		if lineCharIdx != 0 || !tokenAhead(CustomAction) {
			advance()
			continue
		}
		parseCustomAction()
	}
	splitContents()
	findCustomActionRefs()
}

func parseCustomAction() {
	var startActionLineIdx = lineIdx
	var identifier = collectUntil('(')
	if _, found := customActions[identifier]; found {
		parserError(fmt.Sprintf("Duplication definition of custom action '%s()'", identifier))
	}
	if _, found := actions[identifier]; found {
		parserError(fmt.Sprintf("Duplication definition of built-in action '%s()'", identifier))
	}

	advance()
	collectUntilExpect('{', 1)
	advance()

	var body = collectUntilIgnoreStrings('}')
	customActions[identifier] = customAction{
		body: body,
	}

	var endActionLineIdx = lineIdx
	advance()

	lines = strings.Split(contents, "\n")
	for i := range lines {
		if i >= startActionLineIdx && i <= endActionLineIdx {
			lines[i] = ""
		}
	}
	contents = strings.Join(lines, "\n")
}

// findCustomActionRefs replaces references to defined actions with their collected body.
func findCustomActionRefs() {
	firstChar()
	for char != -1 {
		if !strings.Contains(lookAheadUntil('\n'), "(") {
			advance()
			continue
		}

		var identifier = strings.Trim(collectUntil('('), " \t\n")
		if _, found := customActions[identifier]; !found {
			collectUntil('\n')
			continue
		}
		advance()

		if char == ')' || (char == ' ' && next(1) == ')') {
			lines[lineIdx] = customActions[identifier].body
			splitContents()
			firstChar()
			continue
		}

		collectUntilExpect(')', 1)

		var actionBody = customActions[identifier].body
		lines[lineIdx] = actionBody
	}
	splitContents()
	firstChar()
}
