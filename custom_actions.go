/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	"regexp"
	"strings"
)

// customAction contains the collected declaration of a custom action.
type customAction struct {
	arguments []parameterDefinition
	body      string
}

// customActions is a map of all the custom actions that have been defined.
var customActions map[string]customAction

// parseCustomActions parses defined actions and collects them.
func parseCustomActions() {
	if !regexp.MustCompile(`action (.*?)\((.*?)\) \{`).MatchString(contents) {
		return
	}
	standardActions()

	customActions = make(map[string]customAction)
	for char != -1 {
		switch {
		case isToken(ForwardSlash):
			collectComment()
		case tokenAhead(Action):
			advance()
			collectActionDefinition()
			continue
		}
		advance()
	}

	fmt.Println(customActions)

	contents = strings.Join(lines, "\n")

	firstChar()
}

func collectActionDefinition() {
	var startLine = lineIdx

	var identifier = collectIdentifier()
	if _, found := customActions[identifier]; found {
		parserError(fmt.Sprintf("Duplication declaration of custom action '%s()'", identifier))
	}
	if _, found := actions[identifier]; found {
		parserError(fmt.Sprintf("Declaration conflicts with built-in action '%s()'", identifier))
	}

	var arguments = collectArgumentDefinitions()
	advance()
	var body = strings.TrimSpace(collectObject())

	var endLine = lineIdx

	for i := 0; i <= endLine && i >= startLine; i++ {
		lines[i] = ""
	}

	customActions[identifier] = customAction{
		arguments: arguments,
		body:      body,
	}
}

func collectArgumentDefinitions() (arguments []parameterDefinition) {
	for char != ')' {
		var valueType tokenType
		var value any
		collectType(&valueType, &value)
		value = nil
		advance()

		var identifier = collectIdentifier()
		arguments = append(arguments, parameterDefinition{
			name:      identifier,
			validType: valueType,
		})

		if char == ',' || char == ' ' {
			advance()
		}
	}
	advanceTimes(2)

	return
}

func printCustomActionsDebug() {
	fmt.Println(ansi("### CUSTOM ACTIONS ###", bold) + "\n")
	for identifier, customAction := range customActions {
		fmt.Println("identifier:", identifier+"()")
		fmt.Println("arguments:")
		fmt.Println(customAction.arguments)
		fmt.Println("body:")
		fmt.Println(customAction.body)
		fmt.Print("\n")
	}
}
