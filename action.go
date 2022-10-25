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
	field        string
	validType    tokenType
	defaultValue actionArgument
	noMax        bool
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
		enoughArgs(&arguments)
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
		if len(arguments) > i {
			typeCheck(check.field, check.validType, arguments[i])
		} else if check.defaultValue.value == nil {
			parserError(fmt.Sprintf("Missing required argument '%s' to call action '%s'", check.field, currentAction))
		}
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

func enoughArgs(arguments *[]actionArgument) {
	var actionArgs = actions[currentAction].args
	for a, arg := range *arguments {
		if (len(actionArgs) - 1) < a {
			break
		}
		if actionArgs[a].noMax {
			break
		}
		if actionArgs[a].defaultValue.value == nil {
			if arg.value == nil {
				parserError(fmt.Sprintf("Missing required argument '%s' to call action '%s'", actionArgs[a].field, currentAction))
			}
		} else if arg.value == nil {
			arg = actionArgs[a].defaultValue
		}
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

var contactValues []string

func contactValue(class string, contentKit string, args []actionArgument) []plistData {
	contactValues = []string{}
	var entryType int
	switch contentKit {
	case "emailaddress":
		entryType = 2
	case "phonenumber":
		entryType = 1
	}
	for _, item := range args {
		contactValues = append(contactValues, plistDict("", []plistData{
			{
				key:      "EntryType",
				dataType: Number,
				value:    entryType,
			},
			{
				key:      "SerializedEntry",
				dataType: Dictionary,
				value: []plistData{
					{
						key:      "link.contentKit." + contentKit,
						dataType: Text,
						value:    item.value,
					},
				},
			},
		}))
	}
	return []plistData{
		{
			key:      class,
			dataType: Dictionary,
			value: []plistData{
				{
					key:      "Value",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "WFContactFieldValues",
							dataType: Array,
							value:    contactValues,
						},
					},
				},
				{
					key:      "WFSerializationType",
					dataType: Text,
					value:    "WFContactFieldValue",
				},
			},
		},
	}
}

func roundingValue(mode string, args []actionArgument) []plistData {
	switch args[2].value {
	case "1":
		args[2].value = "Ones Place"
	case "10":
		args[2].value = "Tens Place"
	case "100":
		args[2].value = "Hundreds Place"
	case "1000":
		args[2].value = "Thousands"
	case "10000":
		args[2].value = "Ten Thousands"
	case "100000":
		args[2].value = "Hundred Thousands"
	case "1000000":
		args[2].value = "Millions"
	}
	return []plistData{
		{
			key:      "WFRoundMode",
			dataType: Text,
			value:    mode,
		},
		argumentValue("WFInput", args, 0),
		{
			key:      "WFRoundTo",
			dataType: Text,
			value:    args[2].value,
		},
	}
}

func calculateStatistics(operation string, args []actionArgument) []plistData {
	return []plistData{
		{
			key:      "WFStatisticsOperation",
			dataType: Text,
			value:    operation,
		},
		variableInput("WFInput", args[1].value.(string)),
	}
}

var ipTypes = []string{"IPv4", "IPv6"}

func checkIPType(args []actionArgument) {
	var ipType = strings.ToUpper(getArgValue(args[0]).(string))
	if !contains(ipTypes, ipType) {
		parserError(fmt.Sprintf("Invalid IP address type of '%s'. Available IP types: %v", ipType, ipTypes))
	}
}

func adjustDate(operation string, unit string, args []actionArgument) []plistData {
	var adjustDateParams = []plistData{
		{
			key:      "WFAdjustOperation",
			dataType: Text,
			value:    operation,
		},
		argumentValue("WFDate", args, 0),
	}
	if unit != "" {
		adjustDateParams = append(adjustDateParams, plistData{
			key:      "WFDuration",
			dataType: Dictionary,
			value: []plistData{
				{
					key:      "Value",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "Unit",
							dataType: Text,
							value:    unit,
						},
						argumentValue("Magnitude", args, 1),
					},
				},
				{
					key:      "WFSerializationType",
					dataType: Text,
					value:    "WFQuantityFieldValue",
				},
			},
		})
	}
	return adjustDateParams
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
		variableInput("Input", args[0].value.(string)),
	}
}
