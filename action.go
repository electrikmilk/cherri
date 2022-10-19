/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"fmt"
	"strings"
)

var currentAction string

type argumentDefinition struct {
	field     string
	validType tokenType
}

type actionArgument struct {
	valueType tokenType
	value     any
}

type action struct {
	ident string
	args  []actionArgument
}

type actionCall func(args []actionArgument) []plistData
type actionCheck func(args []actionArgument)

type actionDefinition struct {
	ident string
	args  []argumentDefinition
	check actionCheck
	call  actionCall
}

var actions map[string]actionDefinition

var hashTypes = []string{"MD5", "SHA1", "SHA256", "SHA512"}

func callAction(arguments []actionArgument, outputName plistData, actionUUID plistData) {
	var ident string
	if actions[currentAction].ident != "" {
		ident = actions[currentAction].ident
	} else {
		ident = strings.ToLower(currentAction)
	}
	var params []plistData
	if actions[currentAction].call != nil {
		params = actions[currentAction].call(arguments)
	}
	checkIdentify(&params, outputName, actionUUID)
	shortcutActions = append(shortcutActions, makeAction(ident, params))
}

func checkAction(arguments []actionArgument) {
	if len(actions[currentAction].args) > 0 {
		enoughArgs(arguments, len(actions[currentAction].args))
		checkTypes(arguments, actions[currentAction].args)
	}
	if actions[currentAction].check != nil {
		actions[currentAction].check(arguments)
	}
}

func realVariableValue(varName string) (varValue variableValue) {
	if _, global := globals[varName]; global {
		varValue = globals[varName]
		hasInputVariables(varName)
	} else if _, found := variables[strings.ToLower(varName)]; found {
		varName = strings.ToLower(varName)
		var argValueType = variables[varName].valueType
		var value = variables[varName].value
		if argValueType == Variable {
			varValue = realVariableValue(value.(string))
		} else {
			varValue = variables[varName]
		}
	} else {
		parserError(fmt.Sprintf("Variable or Global '%s' does not exist", varName))
	}
	return
}

func checkTypes(arguments []actionArgument, checks []argumentDefinition) {
	for i, check := range checks {
		typeCheck(check.field, check.validType, arguments[i])
	}
}

func typeCheck(field string, validType tokenType, argument actionArgument) {
	var argValueType = argument.valueType
	var argVal = argument.value
	if argValueType == Action {
		// FIXME: Identify the output type of action
		return
	}
	if argValueType == Variable {
		var varName = argVal.(string)
		var getVar = realVariableValue(varName)
		argValueType = getVar.valueType
		argVal = getVar.value
		if argValueType == Action {
			// FIXME: Identify the output type of action
			return
		}
		if argValueType != validType && validType != Variable {
			lineIdx--
			parserError(fmt.Sprintf("Invalid variable value '%v' of type '%s' for argument '%s' of type '%s' in '%s()'",
				argVal,
				typeName(argValueType),
				field,
				typeName(validType),
				currentAction,
			))
		}
	} else if argument.valueType != validType {
		lineIdx -= 2
		parserError(fmt.Sprintf("Invalid value '%v' of type '%s' for argument '%s' of type '%s' in '%s()'",
			argVal,
			typeName(argValueType),
			field,
			typeName(validType),
			currentAction,
		))
	}
}

func getArgValue(variable actionArgument) any {
	if variable.valueType == Variable {
		var varName = variable.value.(string)
		if _, found := variables[varName]; found {
			if variables[varName].valueType == Variable {
				return getArgValue(actionArgument{
					valueType: variables[varName].valueType,
					value:     variables[varName].value,
				})
			}
			return variables[varName].value
		}
	}
	return variable.value
}

// FIXME: default values for arguments instead of requiring all of them?
func enoughArgs(arguments []actionArgument, min int) {
	if len(arguments) < min {
		lineIdx--
		parserError(fmt.Sprintf(
			"Not enough arguments to call '%s()'. Minimum is %d arg(s), but provided %d arg(s)",
			currentAction,
			min,
			len(arguments),
		))
	}
}

func checkIdentify(params *[]plistData, outputName plistData, actionUUID plistData) {
	if outputName.value != nil {
		*params = append(*params, outputName)
	}
	if actionUUID.value != nil {
		*params = append(*params, actionUUID)
	}
}

func changeCase(textCase string, args []actionArgument) []plistData {
	return []plistData{
		{
			key:      "Show-text",
			dataType: Boolean,
			value:    true,
		},
		{
			key:      "WFCaseType",
			dataType: Text,
			value:    textCase,
		},
		argumentValue("text", args, 0),
	}
}

func textParts(args []actionArgument) []plistData {
	return []plistData{
		{
			key:      "Show-text",
			dataType: Boolean,
			value:    true,
		},
		{
			key:      "WFTextSeparator",
			dataType: Text,
			value:    "Custom",
		},
		argumentValue("text", args, 0),
		argumentValue("WFTextCustomSeparator", args, 1),
	}
}

func replaceText(caseSensitive bool, regExp bool, args []actionArgument) []plistData {
	return []plistData{
		{
			key:      "WFReplaceTextCaseSensitive",
			dataType: Boolean,
			value:    caseSensitive,
		},
		{
			key:      "WFReplaceTextRegularExpression",
			dataType: Boolean,
			value:    regExp,
		},
		argumentValue("WFReplaceTextFind", args, 0),
		argumentValue("WFReplaceTextReplace", args, 1),
		argumentValue("WFInput", args, 2),
	}
}

func count(countType string, args []actionArgument) []plistData {
	return []plistData{
		{
			key:      "WFCountType",
			dataType: Text,
			value:    countType,
		},
		inputValue("Input", args[0].value.(string), ""),
	}
}
