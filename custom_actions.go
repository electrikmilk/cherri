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
		fmt.Print("Parsing custom actions... ")
	}
	customActions = make(map[string]customAction)
	chars = []rune(contents)
	idx = -1
	advance()
	for char != -1 {
		if lineCharIdx != 1 || !tokenAhead(CustomAction) {
			advance()
			continue
		}
		parseCustomAction()
	}
	lines = strings.Split(contents, "\n")
	chars = []rune(contents)

	firstChar()
	findCustomActionRefs()
	firstChar()

	if args.Using("debug") {
		printCustomActionsDebug()
	}

	customActions = map[string]customAction{}

	if args.Using("debug") {
		fmt.Print(ansi("done!", green) + "\n")
	}
}

func printCustomActionsDebug() {
	fmt.Println(ansi("### CUSTOM ACTIONS ###", bold) + "\n")

	for identifier, customAction := range customActions {
		fmt.Println("identifier: " + identifier)
		fmt.Println("body:")
		fmt.Println(customAction.body)
		fmt.Print("\n")
	}

	fmt.Print("\n")
}

func parseCustomAction() {
	advance()
	var startActionLineIdx = lineIdx
	var identifier = collectIdentifier()
	if _, found := customActions[identifier]; found {
		parserError(fmt.Sprintf("Duplication definition of custom action '%s()'", identifier))
	}
	if _, found := actions[identifier]; found {
		parserError(fmt.Sprintf("Duplication definition of built-in action '%s()'", identifier))
	}

	advance()
	advanceUntilExpect('{', 1)
	advance()

	var body = collectObject()
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
	for char != -1 {
		if !strings.Contains(lines[lineIdx], "(") {
			advance()
			continue
		}

		var identifier = strings.Trim(collectUntil('('), " \t\n")
		if _, found := customActions[identifier]; !found {
			advanceUntil('\n')
			continue
		}
		advance()

		if char == ')' || (char == ' ' && next(1) == ')') {
			lines[lineIdx] = customActions[identifier].body
			splitContents()
			firstChar()
			continue
		}

		advanceUntilExpect(')', 1)

		var actionBody = customActions[identifier].body
		lines[lineIdx] = actionBody
	}
	splitContents()
	firstChar()
}
