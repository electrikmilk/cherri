/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/electrikmilk/args-parser"
)

var usingCustomActions bool

// customAction contains the collected declaration of a custom action.
type customAction struct {
	definition actionDefinition
	body       string
	used       bool
}

// customActions is a map of all the custom actions that have been defined.
var customActions map[string]*customAction

// handleCustomActions parses defined custom actions and checks their usage.
func handleCustomActions() {
	if !regexp.MustCompile(`action (.*?)\((.*?)\)`).MatchString(contents) {
		return
	}
	parseCustomActions()

	checkCustomActionUsage(contents)
	for _, action := range customActions {
		if strings.ContainsAny(action.body, "()") {
			checkCustomActionUsage(action.body)
		}
	}

	usingCustomActions = isUsingCustomActions()
	if usingCustomActions {
		var customActionsHeader = generateCustomActionHeader()
		lines = append([]string{customActionsHeader}, lines...)

		resetParse()
	}

	if args.Using("debug") {
		printCustomActionsDebug()
		fmt.Println(contents)
	}
}

func parseCustomActions() {
	customActions = make(map[string]*customAction)
	for char != -1 {
		switch {
		case isChar('/'):
			collectComment()
			continue
		case lineCharIdx == 0 && tokenAhead(Action):
			advance()
			collectCustomActionDefinition()
			continue
		}
		advance()
	}

	resetParse()
}

func isUsingCustomActions() bool {
	for _, action := range customActions {
		if action.used {
			hasShortcutInputVariables = true
			return true
		}
	}
	return false
}

func collectCustomActionDefinition() {
	var lineRef = newLineReference()
	var identifier, arguments, outputType = collectActionDefinition('{')

	advanceUntilExpect('{', 3)
	advance()

	var body = strings.TrimSpace(collectObject())

	lineRef.replaceLines()

	customActions[identifier] = &customAction{
		definition: actionDefinition{
			parameters: arguments,
			outputType: outputType,
		},
		body: body,
	}
}

var actionUsageRegex = regexp.MustCompile(`([a-zA-Z0-9]+)\(`)

func checkCustomActionUsage(content string) {
	var matches = actionUsageRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return
	}
	for _, match := range matches {
		var ref = strings.TrimSpace(match[1])
		if customAction, found := customActions[ref]; found {
			if !customAction.used {
				customActions[ref].used = true
			}
		}
	}
}

func generateCustomActionHeader() string {
	var outputActionRegex = regexp.MustCompile(`(?:must)?[o|O]utput(?:OrClipboard)?\((.*?)\)`)
	var customActionsHeader strings.Builder
	customActionsHeader.WriteString("if ShortcutInput {\n")
	customActionsHeader.WriteString("    @_cherri_empty_dictionary: dictionary\n")
	customActionsHeader.WriteString("    const _cherri_dictionary_type_name = typeOf(_cherri_empty_dictionary)\n")
	customActionsHeader.WriteString("    const _cherri_inputType = typeOf(ShortcutInput)\n")
	customActionsHeader.WriteString("    if _cherri_inputType == _cherri_dictionary_type_name {\n")
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
			case Integer, Float, Bool:
				customActionsHeader.WriteString(fmt.Sprintf("@%s = number(%s)\n", param.name, argumentReference))
			case Dict:
				customActionsHeader.WriteString(fmt.Sprintf("@%s = getDictionary(%s)\n", param.name, argumentReference))
			case Arr:
				customActionsHeader.WriteString(fmt.Sprintf("const %s_array_dictionary = getDictionary(%s)\n", argumentReference, argumentReference))
				customActionsHeader.WriteString(fmt.Sprintf("                const %s_array = getValue(%s,\"array\")\n", argumentReference, argumentReference))
				customActionsHeader.WriteString(fmt.Sprintf("                for %s_array_item in %s_array {\n", argumentReference, argumentReference))
				customActionsHeader.WriteString(fmt.Sprintf("                    @%s += %s_array_item\n                }", param.name, argumentReference))
			default:
				customActionsHeader.WriteString(fmt.Sprintf("@%s = %s\n", param.name, argumentReference))
			}

			if param.defaultValue != nil {
				var defaultValue = param.defaultValue
				if reflect.TypeOf(param.defaultValue).Kind() == reflect.String {
					defaultValue = fmt.Sprintf("\"%s\"", defaultValue)
				}

				customActionsHeader.WriteString(fmt.Sprintf("                if !%s {\n", argumentReference))
				customActionsHeader.WriteString(fmt.Sprintf("                   @%s = %v\n", param.name, defaultValue))
				customActionsHeader.WriteString("                }\n")
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

	return customActionsHeader.String()
}

func makeCustomActionRef(identifier *string) any {
	var customAction = customActions[*identifier]
	setCurrentAction(*identifier, &customAction.definition)

	var arguments []actionArgument
	var paramsSize = len(customAction.definition.parameters)
	if paramsSize != 0 {
		advance()
		arguments = collectArguments()

		currentArgumentsSize = len(arguments)
		checkAction()
	}

	var variableIdentifier = fmt.Sprintf("_%s_cherri_call", *identifier)
	var customActionCall = makeCustomActionCall(identifier, &arguments)
	insertReference(variableIdentifier, Dict, customActionCall, true)

	advanceUntil('\n')

	var runSelfAction = makeActionValue("runSelf", []actionArgument{
		{
			valueType: Variable,
			value: varValue{
				valueType: Variable,
				value:     variableIdentifier,
			},
		},
	})

	if customAction.definition.outputType != "" {
		var outputIdentifier = fmt.Sprintf("_%s_cherri_call_output", *identifier)
		insertReference(outputIdentifier, Action, runSelfAction, true)

		return coerceOutputValue(outputIdentifier, customAction.definition.outputType, runSelfAction)
	}

	return runSelfAction
}

func coerceOutputValue(value any, valueType tokenType, defaultValue any) any {
	switch valueType {
	case String:
		return makeActionValue("text", []actionArgument{
			{
				valueType: String,
				value:     fmt.Sprintf("{%s}", value),
			},
		})
	case Bool, Integer:
		return makeActionValue("number", []actionArgument{
			{
				valueType: Variable,
				value: varValue{
					valueType: Variable,
					value:     value,
				},
			},
		})
	case Dict:
		return makeActionValue("getDictionary", []actionArgument{
			{
				valueType: Variable,
				value: varValue{
					valueType: Variable,
					value:     value,
				},
			},
		})
	default:
		return defaultValue
	}
}

func makeCustomActionCall(identifier *string, arguments *[]actionArgument) map[string]any {
	var argumentValues []any
	if len(*arguments) != 0 {
		for _, argument := range *arguments {
			var argumentValue = fmt.Sprintf("%v", argument.value)
			switch argument.valueType {
			case String:
				argumentValues = append(argumentValues, fmt.Sprintf("\"%s\"", argumentValue))
			case Variable:
				var identifier = argument.value.(varValue).value.(string)
				var variableValue, found = getVariableValue(identifier)
				if !found {
					parserError(fmt.Sprintf("Undefined reference '%s'", argumentValue))
				}
				if variableValue.valueType == Arr {
					var wrappedArray = map[string]any{"array": fmt.Sprintf("{%s}", identifier)}
					argumentValues = append(argumentValues, wrappedArray)
				} else {
					argumentValues = append(argumentValues, fmt.Sprintf("\"{%s}\"", identifier))
				}
			case Arr:
				var wrappedArray = map[string]any{"array": argument.value}
				argumentValues = append(argumentValues, wrappedArray)
			default:
				argumentValues = append(argumentValues, argument.value)
			}
		}
	}

	var customActionCall = map[string]any{
		"cherri_functions": 1,
		"function":         *identifier,
		"arguments":        argumentValues,
	}

	return customActionCall
}

func printCustomActionsDebug() {
	fmt.Println(ansi("### CUSTOM ACTIONS ###", bold) + "\n")
	for identifier, customAction := range customActions {
		fmt.Println("identifier:", identifier+"()")
		fmt.Println("used:", customAction.used)
		fmt.Println("output type:", customAction.definition.outputType)
		fmt.Println("parameters:")
		fmt.Println(customAction.definition.parameters)
		fmt.Println("body:")
		fmt.Println(customAction.body)
		fmt.Println("(end)")
		fmt.Print("\n")
	}
}
