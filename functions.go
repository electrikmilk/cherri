/*
 * Copyright (c) Cherri
 */

/*

Shortcut Functions

This implementation runs actions within the Shortcut isolated from the
rest of the Shortcut via input to a Run Shortcut action and interception
at the beginning of the Shortcut using a generated header.

By the mechanisms of input and output in a Shortcut and our internal references,
the assignment of a function call is assigned to the output of the Run Shortcut
action with defined input. Thus, any output the function produces using an output
action will be returned to the original call.

*/

package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/electrikmilk/args-parser"
)

var usingFunctions bool

// function contains the collected declaration of a function.
type function struct {
	definition actionDefinition
	body       string
	used       bool
}

// functions is a map of all the functions that have been defined.
var functions map[string]*function

// handleFunctions parses defined functions and checks their usage.
func handleFunctions() {
	if !regexp.MustCompile(`function (.*?)\((.*?)\)`).MatchString(contents) {
		return
	}
	parseFunctions()

	checkFunctionUsage(contents)
	for _, action := range functions {
		if strings.ContainsAny(action.body, "()") {
			checkFunctionUsage(action.body)
		}
	}

	usingFunctions = isUsingFunctions()
	if usingFunctions {
		var functionsHeader = generateFunctionsHeader()
		lines = append([]string{functionsHeader}, lines...)

		resetParse()
	}

	if args.Using("debug") {
		printFunctionsDebug()
		fmt.Println(contents)
	}
}

func parseFunctions() {
	functions = make(map[string]*function)
	for char != -1 {
		switch {
		case commentAhead():
			collectComment()
		case startOfLineTokenAhead(Function):
			advance()
			collectFunctionDefinition()
		}
		advance()
	}

	resetParse()
}

func isUsingFunctions() bool {
	for _, action := range functions {
		if action.used {
			hasShortcutInputVariables = true
			return true
		}
	}
	return false
}

func collectFunctionDefinition() {
	var lineRef = newLineReference()
	var identifier, arguments, outputType = collectActionDefinition('{')

	advanceUntilExpect('{', 3)
	advance()

	var body = strings.TrimSpace(collectObject())

	lineRef.replaceLines()

	functions[identifier] = &function{
		definition: actionDefinition{
			parameters: arguments,
			outputType: outputType,
		},
		body: body,
	}
}

var actionUsageRegex = regexp.MustCompile(`([a-zA-Z0-9]+)\(`)

func checkFunctionUsage(content string) {
	var matches = actionUsageRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return
	}
	for _, match := range matches {
		var ref = strings.TrimSpace(match[1])
		if function, found := functions[ref]; found {
			if !function.used {
				functions[ref].used = true
			}
		}
	}
}

func generateFunctionsHeader() string {
	var functionsHeader strings.Builder
	functionsHeader.WriteString("if ShortcutInput {\n")
	functionsHeader.WriteString("    @_cherri_empty_dictionary: dictionary\n")
	functionsHeader.WriteString("    const _cherri_dictionary_type_name = typeOf(@_cherri_empty_dictionary)\n")
	functionsHeader.WriteString("    const _cherri_inputType = typeOf(ShortcutInput)\n")
	functionsHeader.WriteString("    if _cherri_inputType == _cherri_dictionary_type_name {\n")
	functionsHeader.WriteString("        const _cherri_input = getDictionary(ShortcutInput)\n")
	functionsHeader.WriteString("        const _cherri_identifier = getValue(_cherri_input, \"cherri_functions\")\n")
	functionsHeader.WriteString("        const _cherri_valid = number(_cherri_identifier)\n")
	functionsHeader.WriteString("        if _cherri_valid == true {\n")
	functionsHeader.WriteString("            const _cherri_function = getValue(_cherri_input, \"function\")\n")
	functionsHeader.WriteString("            const _cherri_function_name = \"{_cherri_function}\"\n")
	functionsHeader.WriteString("            const _cherri_function_args = getValue(_cherri_input, \"arguments\")\n")

	generateFunctions(&functionsHeader)

	functionsHeader.WriteString("            output(nil)\n")
	functionsHeader.WriteString("        }\n    }\n}")

	return functionsHeader.String()
}

func generateFunctions(functionsHeader *strings.Builder) {
	var outputActionRegex = regexp.MustCompile(`(?:must)?[o|O]utput(?:OrClipboard)?\((.*?)\)`)
	for identifier, function := range functions {
		if !function.used {
			continue
		}

		functionsHeader.WriteString("            if _cherri_function_name == \"")
		functionsHeader.WriteString(identifier)
		functionsHeader.WriteString("\" {\n")

		handleFunctionArguments(functionsHeader, identifier, function)

		functionsHeader.WriteString(function.body)

		if !outputActionRegex.MatchString(function.body) {
			functionsHeader.WriteString("\noutput(nil)")
		}

		functionsHeader.WriteRune('\n')

		functionsHeader.WriteString("            }\n")
	}
}

func handleFunctionArguments(functionsHeader *strings.Builder, identifier string, function *function) {
	for i, param := range function.definition.parameters {
		var idx = i + 1
		var argumentReference = fmt.Sprintf("_cherri_%s_arg_%d_%s", identifier, idx, param.name)

		functionsHeader.WriteString(fmt.Sprintf("                const %s = getListItem(_cherri_function_args, %d)\n", argumentReference, idx))
		functionsHeader.WriteString(fmt.Sprintf("                @%s: %s\n                ", param.name, param.validType))

		switch param.validType {
		case String:
			functionsHeader.WriteString(fmt.Sprintf("@%s = \"{%s}\"\n", param.name, argumentReference))
		case Integer, Float, Bool:
			functionsHeader.WriteString(fmt.Sprintf("@%s = number(%s)\n", param.name, argumentReference))
		case Dict:
			functionsHeader.WriteString(fmt.Sprintf("@%s = getDictionary(%s)\n", param.name, argumentReference))
		case Arr:
			functionsHeader.WriteString(fmt.Sprintf("const %s_array_dictionary = getDictionary(%s)\n", argumentReference, argumentReference))
			functionsHeader.WriteString(fmt.Sprintf("                const %s_array = getValue(%s,\"array\")\n", argumentReference, argumentReference))
			functionsHeader.WriteString(fmt.Sprintf("                for %s_array_item in %s_array {\n", argumentReference, argumentReference))
			functionsHeader.WriteString(fmt.Sprintf("                    @%s += @%s_array_item\n                }", param.name, argumentReference))
		default:
			functionsHeader.WriteString(fmt.Sprintf("@%s = %s\n", param.name, argumentReference))
		}

		if param.defaultValue != nil {
			var defaultValue = param.defaultValue
			if reflect.TypeOf(param.defaultValue).Kind() == reflect.String {
				defaultValue = fmt.Sprintf("\"%s\"", defaultValue)
			}

			functionsHeader.WriteString(fmt.Sprintf("                if !%s {\n", argumentReference))
			functionsHeader.WriteString(fmt.Sprintf("                   @%s = %v\n", param.name, defaultValue))
			functionsHeader.WriteString("                }\n")
		}

		functionsHeader.WriteRune('\n')
	}
}

func makeFunctionRef(identifier *string) any {
	var function = functions[*identifier]
	setCurrentAction(*identifier, &function.definition)

	var arguments []actionArgument
	var paramsSize = len(function.definition.parameters)
	if paramsSize != 0 {
		advance()
		arguments = collectArguments()

		currentArgumentsSize = len(arguments)
		checkAction()
	}

	var variableIdentifier = fmt.Sprintf("_%s_cherri_call", *identifier)
	var functionCall = makeFunctionCall(identifier, &arguments)
	insertReference(variableIdentifier, Dict, functionCall, true)

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

	if function.definition.outputType != "" {
		var outputIdentifier = fmt.Sprintf("_%s_cherri_call_output", *identifier)
		insertReference(outputIdentifier, Action, runSelfAction, true)

		return coerceOutputValue(outputIdentifier, function.definition.outputType, runSelfAction)
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

func makeFunctionCall(identifier *string, arguments *[]actionArgument) map[string]any {
	var argumentValues []any
	if len(*arguments) != 0 {
		for _, argument := range *arguments {
			var argumentValue = fmt.Sprintf("%v", argument.value)
			switch argument.valueType {
			case String:
				argumentValues = append(argumentValues, fmt.Sprintf("%s", argumentValue))
			case Variable:
				var identifier = argument.value.(varValue).value.(string)
				var variableValue, found = getVariableValue(identifier)
				if !found {
					parserError(fmt.Sprintf("Undefined reference '%s'", identifier))
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

	var functionCall = map[string]any{
		"cherri_functions": 1,
		"function":         *identifier,
		"arguments":        argumentValues,
	}

	return functionCall
}

func printFunctionsDebug() {
	fmt.Println(ansi("### FUNCTIONS ###", bold) + "\n")
	for identifier, function := range functions {
		fmt.Println("identifier:", identifier+"()")
		fmt.Println("used:", function.used)
		fmt.Println("output type:", function.definition.outputType)
		fmt.Println("parameters:")
		fmt.Println(function.definition.parameters)
		fmt.Println("body:")
		fmt.Println(function.body)
		fmt.Println("(end)")
		fmt.Print("\n")
	}
}
