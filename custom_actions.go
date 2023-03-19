/*
 * Copyright (c) 2023 Brandon Jordan
 */

package main

import (
	"fmt"
	"strings"
)

// customAction contains the collected declaration of a custom action.
type customAction struct {
	body      string
	arguments []string
}

// customActions is a map of all the custom actions that have been defined.
var customActions map[string]customAction

// parseCustomActions parses defined actions and collects them.
func parseCustomActions() {
	customActions = make(map[string]customAction)
	variables = make(map[string]variableValue)
	chars = strings.Split(contents, "")
	idx = -1
	advance()
	for char != -1 {
		if tokenAhead(At) {
			var identifier string
			if strings.Contains(lookAheadUntil('\n'), "=") {
				identifier = collectUntil(' ')
			} else {
				identifier = collectUntil('\n')
			}
			// Pre-define variables so the value checker doesn't freak out
			variables[identifier] = variableValue{}
		}
		if lineCharIdx != 1 || !tokenAhead(CustomAction) {
			advance()
			continue
		}
		var startActionLineIdx = lineIdx
		var identifier = collectUntil('(')
		if _, found := customActions[identifier]; found {
			parserError(fmt.Sprintf("Duplication declaration of action '%s()'", identifier))
		}
		if _, found := actions[identifier]; found {
			parserError(fmt.Sprintf("Duplication declaration of action '%s()'", identifier))
		}
		advance()
		var argumentDefinitions []string
		if !isToken(RightParen) && next(1) != ')' {
			argumentDefinitions = collectArgumentDefinitions()
		}
		collectUntilExpect('{', 1)
		advance()
		var body string
		var insideString = false
		for char != -1 {
			if char == '}' && !insideString {
				break
			}
			if char == '"' {
				if insideString {
					insideString = false
				} else {
					insideString = true
				}
			}
			body += string(char)
			advance()
		}
		var endActionLineIdx = lineIdx
		advance()
		customActions[identifier] = customAction{
			body:      body,
			arguments: argumentDefinitions,
		}
		lines = strings.Split(contents, "\n")
		for i := range lines {
			if i >= startActionLineIdx && i <= endActionLineIdx {
				lines[i] = ""
			}
		}
		contents = strings.Join(lines, "\n")
	}
	chars = strings.Split(contents, "")
	findCustomActionRefs()
	firstChar()
}

// findCustomActionRefs replaces references to defined actions with their collected body.
func findCustomActionRefs() {
	firstChar()
	for char != -1 {
		if !strings.Contains(lookAheadUntil('\n'), "(") {
			advance()
			continue
		}

		var returnVariable string
		if tokenAhead(At) {
			if !strings.Contains(lookAheadUntil('\n'), "=") {
				continue
			}
			returnVariable = collectUntil(' ')
			collectUntilExpect('=', 1)
			advance()
		}

		var identifier = strings.Trim(collectUntil('('), " \t\n")
		if _, found := customActions[identifier]; !found {
			collectUntil('\n')
			continue
		}
		advance()

		if char == ')' || next(1) == ')' {
			lines[lineIdx] = customActions[identifier].body
			contents = strings.Join(lines, "\n")
			lines = strings.Split(contents, "\n")
			chars = strings.Split(contents, "")
			firstChar()
			continue
		}

		var arguments = collectArguments()
		if len(arguments) < len(customActions[identifier].arguments) {
			parserError(fmt.Sprintf("Not enough arguments to call declared action '%s()'", identifier))
		}
		if len(arguments) > len(customActions[identifier].arguments) {
			parserError(fmt.Sprintf("Too many arguments to call declared action '%s()'", identifier))
		}

		var actionBody = customActions[identifier].body
		var argumentDefinitions = customActions[identifier].arguments
		for i, argName := range argumentDefinitions {
			if strings.Contains(actionBody, argName) {
				var replacementValue = fmt.Sprintf("%v", arguments[i].value)
				if arguments[i].valueType == String {
					replacementValue = "\"" + replacementValue + "\""
				}
				if strings.Contains(actionBody, "{"+argName+"}") && arguments[i].valueType != Variable {
					var uniqueVariable = identifier + "-" + argName
					actionBody = strings.ReplaceAll(actionBody, argName, uniqueVariable)
					actionBody = "@" + uniqueVariable + " = " + replacementValue + "\n" + actionBody
				} else {
					actionBody = strings.ReplaceAll(actionBody, argName, replacementValue)
				}
			}
		}

		var actionBodyLines = strings.Split(actionBody, "\n")
		for i, line := range actionBodyLines {
			if len(line) == 0 {
				continue
			}
			if startsWith(strings.Trim(line, " "), "return") && returnVariable != "" {
				actionBodyLines[i] = strings.Replace(line, "return ", "@"+returnVariable+" = ", 1)
			}
		}
		actionBody = strings.Join(actionBodyLines, "\n")
		lines[lineIdx] = actionBody

		contents = strings.Join(lines, "\n")
		lines = strings.Split(contents, "\n")
		chars = strings.Split(contents, "")
		firstChar()
	}
}

// collectArgumentDefinitions loosely collects argument names for an action definition into a string slice.
func collectArgumentDefinitions() (arguments []string) {
	for strings.Contains(lookAheadUntil(')'), ",") {
		arguments = append(arguments, strings.Trim(collectUntil(','), " \t\n"))
		advance()
	}
	if next(1) != ')' {
		arguments = append(arguments, strings.Trim(collectUntil(')'), " \t\n"))
		advance()
	}
	return
}
