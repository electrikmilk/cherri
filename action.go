/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"fmt"
	"strings"
)

var currentAction string

type parameterDefinition struct {
	name         string
	validType    tokenType
	key          string
	defaultValue actionArgument
	optional     bool
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

type makeParams func(args []actionArgument) []plistData
type paramCheck func(args []actionArgument)

type actionDefinition struct {
	stdIdentifier string
	identifier    string
	parameters    []parameterDefinition
	check         paramCheck
	make          makeParams
}

var actions map[string]actionDefinition

func callAction(arguments []actionArgument, outputName plistData, actionUUID plistData) {
	var ident string
	if actions[currentAction].stdIdentifier != "" {
		ident = actions[currentAction].stdIdentifier
		ident = "is.workflow.actions." + ident
	} else if actions[currentAction].identifier != "" {
		ident = actions[currentAction].identifier
	} else {
		ident = strings.ToLower(currentAction)
		ident = "is.workflow.actions." + ident
	}
	var params []plistData
	if actions[currentAction].make != nil {
		params = actions[currentAction].make(arguments)
	} else if actions[currentAction].parameters != nil {
		for i, a := range actions[currentAction].parameters {
			if len(arguments) > i {
				if a.validType == Variable {
					params = append(params, variableInput(a.key, arguments[i].value.(string)))
				} else {
					params = append(params, argumentValue(a.key, arguments, i))
				}
			}
		}
	}
	checkIdentify(&params, outputName, actionUUID)
	shortcutActions = append(shortcutActions, makeAction(ident, params))
}

func makeAction(ident string, params []plistData) (action string) {
	action = plistDict("", []plistData{
		{
			key:      "WFWorkflowActionIdentifier",
			dataType: Text,
			value:    ident,
		},
		{
			key:      "WFWorkflowActionParameters",
			dataType: Dictionary,
			value:    params,
		},
	})
	return
}

func standardAction(ident string, params []plistData) string {
	ident = "is.workflow.actions." + ident
	return makeAction(ident, params)
}

func checkAction(arguments []actionArgument) {
	if len(actions[currentAction].parameters) > 0 {
		enoughArgs(&arguments)
		checkTypes(arguments, actions[currentAction].parameters)
	}
	if actions[currentAction].check != nil {
		actions[currentAction].check(arguments)
	}
}

func checkEnum(name string, arg actionArgument, enum []string) {
	var value = strings.ToLower(getArgValue(arg).(string))
	if !contains(enum, value) {
		var enumList string
		for _, e := range enum {
			enumList += "- " + e + "\n"
		}
		parserError(fmt.Sprintf("Invalid %s of '%s'.\n\nAvailable %ss:\n%s", name, value, name, enumList))
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

func checkTypes(arguments []actionArgument, checks []parameterDefinition) {
	for i, check := range checks {
		if len(arguments) > i {
			typeCheck(check.name, check.validType, arguments[i])
			if check.defaultValue.value == arguments[i].value {
				parserWarning(fmt.Sprintf("Value for argument %d '%s' for action '%s()' of '%v', is the same as the default value.", i+1, check.name, currentAction, arguments[i].value))
			}
		} else if check.defaultValue.value == nil && check.optional != true {
			parserError(fmt.Sprintf("Missing required argument %d '%s' to call action '%s()'", i+1, check.name, currentAction))
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
		lineIdx--
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
	var actionArgs = actions[currentAction].parameters
	for a, arg := range *arguments {
		if (len(actionArgs) - 1) < a {
			break
		}
		if actionArgs[a].noMax {
			break
		}
		if actionArgs[a].defaultValue.value == nil && actionArgs[a].optional != true {
			if arg.value == nil {
				parserError(fmt.Sprintf("Missing required argument '%s' to call action '%s'", actionArgs[a].name, currentAction))
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

func contactValue(key string, contentKit string, args []actionArgument) []plistData {
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
						key:      "link.contentkit." + contentKit,
						dataType: Text,
						value:    item.value,
					},
				},
			},
		}))
	}
	return []plistData{
		{
			key:      key,
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
	switch args[1].value {
	case "1":
		args[1].value = "Ones Place"
	case "10":
		args[1].value = "Tens Place"
	case "100":
		args[1].value = "Hundreds Place"
	case "1000":
		args[1].value = "Thousands"
	case "10000":
		args[1].value = "Ten Thousands"
	case "100000":
		args[1].value = "Hundred Thousands"
	case "1000000":
		args[1].value = "Millions"
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
			value:    args[1].value,
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
