/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/electrikmilk/args-parser"
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

	resetParse()

	checkCustomActionUsage(contents)
	for _, action := range customActions {
		if strings.ContainsAny(action.body, "()") {
			checkCustomActionUsage(action.body)
		}
	}
	makeCustomActionsHeader()

	if args.Using("debug") {
		printCustomActionsDebug()
		fmt.Println(contents)
	}
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

func checkCustomActionUsage(content string) {
	var actionUsageRegex = regexp.MustCompile(`(action )?([a-zA-Z0-9]+)\(`)
	var matches = actionUsageRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return
	}
	for _, match := range matches {
		var ref = strings.TrimSpace(match[2])
		if _, found := customActions[ref]; found {
			customActions[ref].used = true
		}
	}
}

func makeCustomActionsHeader() {
	var outputActionRegex = regexp.MustCompile(`(?:must)?[o|O]utput(?:OrClipboard)?\((.*?)\)`)
	var customActionsHeader strings.Builder
	customActionsHeader.WriteString("if ShortcutInput {\n")
	customActionsHeader.WriteString("    const _cherri_inputType = typeOf(ShortcutInput)\n")
	customActionsHeader.WriteString("    if _cherri_inputType == \"Dictionary\" {\n")
	customActionsHeader.WriteString("        const _cherri_input = getDictionary(ShortcutInput)\n")
	customActionsHeader.WriteString("        const _cherri_identifier = getValue(_cherri_input, \"cherri_functions\")\n")
	customActionsHeader.WriteString("        const _cherri_valid = number(_cherri_identifier)\n")
	customActionsHeader.WriteString("        if _cherri_valid == true {\n")
	customActionsHeader.WriteString("            const _cherri_function = getValue(_cherri_input, \"function\")\n")
	customActionsHeader.WriteString("            const _cherri_function_name = \"{_cherri_function}\"\n")
	customActionsHeader.WriteString("            const _cherri_function_args = getValue(_cherri_input, \"arguments\")\n")

	for identifier, customAction := range customActions {
		if !customAction.used {
			continue
		}

		customActionsHeader.WriteString("            if _cherri_function_name == \"")
		customActionsHeader.WriteString(identifier)
		customActionsHeader.WriteString("\" {\n")

		for i, param := range customAction.definition.parameters {
			var idx = i + 1
			var argumentReference = fmt.Sprintf("_cherri_%s_arg_%d_%s", identifier, idx, param.name)

			customActionsHeader.WriteString(fmt.Sprintf("                const %s = getListItem(_cherri_function_args, %d)\n", argumentReference, idx))
			customActionsHeader.WriteString(fmt.Sprintf("                @%s: %s\n                ", param.name, param.validType))

			switch param.validType {
			case String:
				customActionsHeader.WriteString(fmt.Sprintf("@%s = \"{%s}\"\n", param.name, argumentReference))
			case Integer:
				customActionsHeader.WriteString(fmt.Sprintf("@%s = number(%s)\n", param.name, argumentReference))
			case Dict:
				customActionsHeader.WriteString(fmt.Sprintf("@%s = getDictionary(%s)\n", param.name, argumentReference))
			case Arr:
				customActionsHeader.WriteString(fmt.Sprintf("const %s_array_dictionary = getDictionary(%s)\n", argumentReference, argumentReference))
				customActionsHeader.WriteString(fmt.Sprintf("                const %s_array = getValue(%s,\"array\")\n", argumentReference, argumentReference))
				customActionsHeader.WriteString(fmt.Sprintf("                for %s_array_item in %s_array {\n", argumentReference, argumentReference))
				customActionsHeader.WriteString(fmt.Sprintf("                    @%s += %s_array_item\n                }", param.name, argumentReference))
			default:
				customActionsHeader.WriteString(argumentReference)
			}

			customActionsHeader.WriteRune('\n')
		}

		customActionsHeader.WriteString(customAction.body)

		if !outputActionRegex.MatchString(customAction.body) {
			customActionsHeader.WriteString("\noutput(nil)")
		}

		customActionsHeader.WriteRune('\n')

		customActionsHeader.WriteString("            }\n")
	}

	customActionsHeader.WriteString("            output(nil)\n")
	customActionsHeader.WriteString("        }\n    }\n}")

	lines = append([]string{customActionsHeader.String()}, lines...)

	resetParse()
}

func handleCustomActionRef(identifier *string) action {
	var customAction = customActions[*identifier]

	var arguments []actionArgument
	advance()
	if char != ')' {
		setCurrentAction(*identifier, &customAction.definition)
		arguments = collectArguments()
	}

	var customActionCall = makeCustomActionCall(identifier, &arguments)

	tokens = append(tokens, token{
		typeof:    Var,
		ident:     fmt.Sprintf("_%s_cherri_call", *identifier),
		valueType: Dict,
		value:     customActionCall,
	})

	variables[*identifier] = variableValue{
		variableType: "Variable",
		valueType:    Dict,
		value:        customActionCall,
		constant:     true,
	}

	advanceUntil('\n')

	return action{
		ident: "runSelf",
		args: []actionArgument{
			{
				valueType: Var,
				value:     fmt.Sprintf("_%s_cherri_call", *identifier),
			},
		},
	}
}

func makeCustomActionCall(identifier *string, arguments *[]actionArgument) (customActionCall interface{}) {
	var customActionCallJSON strings.Builder
	customActionCallJSON.WriteString("{\"cherri_functions\": 1,\"function\": \"")
	customActionCallJSON.WriteString(*identifier)
	customActionCallJSON.WriteString("\",\"arguments\": [")
	if len(*arguments) > 0 {
		for i, argument := range *arguments {
			var argumentValue = fmt.Sprintf("%v", argument.value)
			switch argument.valueType {
			case String:
				argumentValue = fmt.Sprintf("\"%s\"", argumentValue)
			case Variable:
				var variableValue, found = getVariableValue(argumentValue)
				if !found {
					parserError(fmt.Sprintf("Variable not found: %s", argumentValue))
				}
				if variableValue.valueType == Arr {
					var jsonBytes, jsonErr = json.Marshal(variableValue.value)
					handle(jsonErr)
					argumentValue = fmt.Sprintf("{\"array\":%s}", string(jsonBytes))
				} else {
					argumentValue = fmt.Sprintf("\"{%s}\"", argumentValue)
				}
			case Arr:
				var jsonBytes, jsonErr = json.Marshal(argument.value)
				handle(jsonErr)
				argumentValue = fmt.Sprintf("{\"array\":%s}", string(jsonBytes))
			case Dict:
				var jsonBytes, jsonErr = json.Marshal(argument.value)
				handle(jsonErr)
				argumentValue = string(jsonBytes)
			}
			customActionCallJSON.WriteString(argumentValue)

			if len(*arguments)-1 != i {
				customActionCallJSON.WriteRune(',')
			}
		}
	}
	customActionCallJSON.WriteString("]}")

	if err := json.Unmarshal([]byte(customActionCallJSON.String()), &customActionCall); err != nil {
		parserError(fmt.Sprintf("Custom action JSON error: %s", err))
	}

	return
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
