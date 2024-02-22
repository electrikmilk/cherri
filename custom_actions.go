/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	"github.com/electrikmilk/args-parser"
	"regexp"
	"strings"
)

// customAction contains the collected declaration of a custom action.
type customAction struct {
	definition actionDefinition
	body       string
}

// customActions is a map of all the custom actions that have been defined.
var customActions map[string]customAction

// parseCustomActions parses defined actions and collects them.
func parseCustomActions() {
	if !regexp.MustCompile(`action (.*?)\((.*?)\)`).MatchString(contents) {
		return
	}
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

	replaceCustomActionRefs()
	resetParse()
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

	var arguments []parameterDefinition
	if next(1) != ')' {
		advance()
		skipWhitespace()
		arguments = collectParameterDefinitions()
	} else {
		advanceTimes(2)
	}

	advanceUntilExpect('{', 3)
	advance()

	var body = strings.TrimSpace(collectObject())

	var endLine = lineIdx

	for i := 0; i <= endLine && i >= startLine; i++ {
		lines[i] = ""
	}

	customActions[identifier] = customAction{
		definition: actionDefinition{
			parameters: arguments,
		},
		body: body,
	}

	if args.Using("debug") {
		printCustomActionsDebug()
	}
}

func collectParameterDefinitions() (arguments []parameterDefinition) {
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

		if char == ',' {
			advance()
		}
		skipWhitespace()
	}
	advance()

	return
}

func replaceCustomActionRefs() {
	resetParse()

	for char != -1 {
		switch {
		case isToken(ForwardSlash):
			collectComment()
		case strings.Contains(lookAheadUntil('\n'), "("):
			customActionCall()
			continue
		}
		advance()
	}
}

func customActionCall() {
	var identifier = collectIdentifier()
	if _, found := actions[identifier]; found {
		advanceUntil('\n')
		return
	}
	if _, found := customActions[identifier]; !found {
		parserError(fmt.Sprintf("Undefined custom action '%s()'", identifier))
	}
	var action = customActions[identifier]

	setCurrentAction(identifier, &action.definition)

	advance()
	skipWhitespace()
	if char != ')' {
		var arguments = collectArguments()
		fmt.Println(arguments)
	}
	/*
		TODO:
			const myFunctionCall = {
			    "cherri_functions": 1,
			    "function": "add",
			    "arguments": ["{operandOne}", "{operandTwo}"]
			}
			const myFunctionReturn = runSelf(myFunctionCall)
	*/
	lines[lineIdx] = ""
}

func printCustomActionsDebug() {
	fmt.Println(ansi("### CUSTOM ACTIONS ###", bold) + "\n")
	for identifier, customAction := range customActions {
		fmt.Println("identifier:", identifier+"()")
		fmt.Println("parameters:")
		fmt.Println(customAction.definition.parameters)
		fmt.Println("body:")
		fmt.Println(customAction.body)
		fmt.Println("(end)")
		fmt.Print("\n")
	}
}
