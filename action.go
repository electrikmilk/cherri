/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"fmt"
	"math"
	"strings"
)

var currentAction string
var isMac = false

type parameterDefinition struct {
	name         string
	validType    tokenType
	key          string
	defaultValue actionArgument
	optional     bool
	infinite     bool
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
	identifier    string
	appIdentifier string
	parameters    []parameterDefinition
	check         paramCheck
	make          makeParams
	outputType    tokenType
	mac           bool
	minVersion    float64
}

var actions map[string]actionDefinition

type libraryDefinition struct {
	identifier string
	make       func(identifier string)
}

var libraries map[string]libraryDefinition

func callAction(arguments []actionArgument, outputName plistData, actionUUID plistData) {
	var ident string
	if actions[currentAction].appIdentifier != "" {
		ident = actions[currentAction].appIdentifier
	} else {
		if actions[currentAction].identifier != "" {
			ident = actions[currentAction].identifier
		} else {
			ident = strings.ToLower(currentAction)
		}
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
	if outputName.value != nil {
		params = append(params, outputName)
	}
	if actionUUID.value != nil {
		params = append(params, actionUUID)
	}
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

func makeStdAction(ident string, params []plistData) string {
	ident = "is.workflow.actions." + ident
	return makeAction(ident, params)
}

func checkAction(arguments []actionArgument) {
	if len(actions[currentAction].parameters) > 0 {
		checkArgs(&arguments)
		checkTypes(arguments, actions[currentAction].parameters)
	}
	if actions[currentAction].check != nil {
		actions[currentAction].check(arguments)
	}
	if actions[currentAction].minVersion != 0 {
		if actions[currentAction].minVersion > iosVersion {
			parserWarning(
				fmt.Sprintf("Action '%s()' is not available in set minimum version '%.1f'", currentAction, math.Ceil(iosVersion)),
			)
		}
	}
	if !isMac && actions[currentAction].mac {
		parserWarning(
			fmt.Sprintf("You've set your Shortcut as non-Mac. Action '%s()' is a Mac only action.", currentAction),
		)
	}
}

func checkEnum(name string, enum []string, args []actionArgument, idx int) {
	if len(args) < idx {
		return
	}
	if args[idx].value == nil {
		return
	}
	var value = strings.ToLower(getArgValue(args[idx]).(string))
	if !contains(enum, value) {
		var enumList string
		for _, e := range enum {
			enumList += "- " + e + "\n"
		}
		parserError(fmt.Sprintf("Invalid %s of '%s'.\n\nAvailable %ss:\n%s", name, value, name, enumList))
	}
}

func realVariableValue(varName string, lastValueType tokenType) (varValue variableValue) {
	if _, global := globals[varName]; global {
		varValue = globals[varName]
		hasInputVariables(varName)
		return
	}
	if _, found := variables[strings.ToLower(varName)]; !found {
		parserError(fmt.Sprintf("Variable or Global '%s' does not exist", varName))
	}
	varName = strings.ToLower(varName)
	var argValueType = variables[varName].valueType
	var value = variables[varName].value
	if argValueType == Variable {
		if lastValueType == Variable && argValueType == Variable {
			parserError("Passed variable value that evaluates to variable")
		}
		varValue = realVariableValue(value.(string), argValueType)
	} else {
		varValue = variables[varName]
	}
	return
}

func checkTypes(arguments []actionArgument, checks []parameterDefinition) {
	for i, check := range checks {
		if len(arguments) > i {
			typeCheck(check.name, check.validType, arguments[i])
			if check.defaultValue.value == arguments[i].value {
				parserWarning(
					fmt.Sprintf(
						"Value for argument %d '%s' for action '%s()' of '%v', is the same as the default value.",
						i+1, check.name, currentAction, arguments[i].value,
					),
				)
			}
		} else if check.defaultValue.value == nil && !check.optional {
			parserError(
				fmt.Sprintf("Missing required argument %d '%s' to call action '%s()'", i+1, check.name, currentAction),
			)
		}
	}
}

func validActionOutput(field string, validType tokenType, value any) {
	var actionIdent = value.(action).ident
	if _, found := actions[actionIdent]; found {
		var actionOutputType = actions[actionIdent].outputType
		if actionOutputType != "" {
			if actionOutputType != validType {
				parserError(
					fmt.Sprintf(
						"Invalid variable value of action '%v' that outputs type '%s' for argument '%s' of type '%s' in '%s()'",
						actionIdent+"()",
						typeName(actionOutputType),
						field,
						typeName(validType),
						currentAction,
					),
				)
			}
		}
	}
}

func typeCheck(field string, validType tokenType, argument actionArgument) {
	var argValueType = argument.valueType
	var argVal = argument.value
	switch {
	case argValueType == Action:
		validActionOutput(field, validType, argVal)
		return
	case argValueType == Variable:
		var varName = argVal.(string)
		var getVar = realVariableValue(varName, String)
		argValueType = getVar.valueType
		argVal = getVar.value
		if argValueType == Action {
			validActionOutput(field, validType, argVal)
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
	case argument.valueType != validType:
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

func checkArgs(arguments *[]actionArgument) {
	var actionArgs = actions[currentAction].parameters
	for a, arg := range *arguments {
		if (len(actionArgs) - 1) < a {
			break
		}
		if actionArgs[a].infinite {
			break
		}
		if !actionArgs[a].optional {
			if arg.value == nil {
				parserError(fmt.Sprintf("Missing required argument '%s' to call action '%s'", actionArgs[a].name, currentAction))
			}
		}
	}
}

func makeLibraries() {
	libraries = make(map[string]libraryDefinition)
}
