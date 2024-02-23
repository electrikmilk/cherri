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
	used       bool
}

// customActions is a map of all the custom actions that have been defined.
var customActions map[string]*customAction

// parseCustomActions parses defined actions and collects them.
func parseCustomActions() {
	if !regexp.MustCompile(`action (.*?)\((.*?)\)`).MatchString(contents) {
		return
	}
	customActions = make(map[string]*customAction)

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

	if args.Using("debug") {
		printCustomActionsDebug()
	}

	replaceCustomActionRefs()
	makeCustomActionsHeader()
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

	customActions[identifier] = &customAction{
		definition: actionDefinition{
			parameters: arguments,
		},
		body: body,
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
	action.used = true

	advance()
	skipWhitespace()
	if char != ')' {
		setCurrentAction(identifier, &action.definition)
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

func makeCustomActionsHeader() {
	var customActionsHeader strings.Builder
	customActionsHeader.WriteString("if ShortcutInput {\n")
	customActionsHeader.WriteString("    const inputType = typeOf(ShortcutInput)\n")
	customActionsHeader.WriteString("    if inputType == \"Dictionary\" {\n")
	customActionsHeader.WriteString("        const input = getDictionary(ShortcutInput)\n")
	customActionsHeader.WriteString("        const identifier = getValue(input, \"cherri_functions\")\n")
	customActionsHeader.WriteString("        const valid = number(identifier)\n")
	customActionsHeader.WriteString("        if valid == true {\n")
	customActionsHeader.WriteString("            const function_name = getValue(input, \"function\")\n")
	customActionsHeader.WriteString("            const function = \"{function_name}\"\n")
	customActionsHeader.WriteString("            const args = getValue(input, \"arguments\")\n")

	for identifier, customAction := range customActions {
		if !customAction.used {
			continue
		}

		customActionsHeader.WriteString("            if function == \"")
		customActionsHeader.WriteString(identifier)
		customActionsHeader.WriteString("\" {\n")

		for i, param := range currentAction.parameters {
			var idx = i + 1
			customActionsHeader.WriteString(fmt.Sprintf("                const arg%d = ", idx))
			customActionsHeader.WriteString(fmt.Sprintf("getListItem(args, %d)\n", idx))
			customActionsHeader.WriteString(fmt.Sprintf("                const %s = ", param.name))

			switch param.validType {
			case String:
				customActionsHeader.WriteString(fmt.Sprintf("\"{arg%d}\"", idx))
			case Integer:
				customActionsHeader.WriteString(fmt.Sprintf("number(arg%d)", idx))
			}

			customActionsHeader.WriteRune('\n')
		}

		customActionsHeader.WriteString(customAction.body)
		customActionsHeader.WriteRune('\n')

		customActionsHeader.WriteString("            }\n")
		fmt.Println(identifier, customAction)
	}

	customActionsHeader.WriteString("        }\n    }\n}")

	lines = append([]string{customActionsHeader.String()}, lines...)
}

func printCustomActionsDebug() {
	fmt.Println(ansi("### CUSTOM ACTIONS ###", bold) + "\n")
	for identifier, customAction := range customActions {
		fmt.Println("identifier:", identifier+"()")
		fmt.Println("used:", customAction.used)
		fmt.Println("parameters:")
		fmt.Println(customAction.definition.parameters)
		fmt.Println("body:")
		fmt.Println(customAction.body)
		fmt.Println("(end)")
		fmt.Print("\n")
	}
}
