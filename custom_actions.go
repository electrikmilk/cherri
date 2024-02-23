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

	for i := startLine; i <= lineIdx; i++ {
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
			var identifier = collectIdentifier()
			if _, found := actions[identifier]; found {
				advanceUntil('\n')
				return
			}
			makeCustomActionCall(&identifier)
			continue
		}
		advance()
	}
}

func makeCustomActionCall(identifier *string) {
	if _, found := customActions[*identifier]; !found {
		parserError(fmt.Sprintf("Undefined custom action '%s()'", *identifier))
	}
	var action = customActions[*identifier]
	action.used = true

	advance()
	skipWhitespace()
	var arguments []actionArgument
	if char != ')' {
		setCurrentAction(*identifier, &action.definition)
		arguments = collectArguments()
	}

	var collectSpaces strings.Builder
	for _, char := range lines[lineIdx] {
		if char != ' ' {
			break
		}
		collectSpaces.WriteRune(' ')
	}
	var spaces = collectSpaces.String()
	collectSpaces.Reset()

	var customActionCall strings.Builder
	customActionCall.WriteString(fmt.Sprintf("%sconst %sCherriCall = {\n", spaces, *identifier))
	customActionCall.WriteString(fmt.Sprintf("%s\"cherri_functions\": 1,\n\"function\": \"", spaces))
	customActionCall.WriteString(*identifier)
	customActionCall.WriteString("\",\n")
	customActionCall.WriteString(fmt.Sprintf("%s\"arguments\": [", spaces))
	if len(arguments) > 0 {
		for i, argument := range arguments {
			var argumentValue = fmt.Sprintf("%v", argument.value)
			if argument.valueType == String {
				argumentValue = fmt.Sprintf("\"%s\"", argumentValue)
			}
			customActionCall.WriteString(argumentValue)

			if len(arguments)-1 != i {
				customActionCall.WriteRune(',')
			}
		}
	}
	customActionCall.WriteString("]\n")
	customActionCall.WriteString(fmt.Sprintf("%s}\n", spaces))

	customActionCall.WriteString(fmt.Sprintf("%sconst myFunctionReturn = runSelf(", spaces))
	customActionCall.WriteString(fmt.Sprintf("%sCherriCall)\n", *identifier))

	lines[lineIdx] = customActionCall.String()
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
