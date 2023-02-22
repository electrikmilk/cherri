/*
 * Copyright (c) 2023 Brandon Jordan
 */

package main

import (
	"fmt"
	"math"
	"strings"
)

// currentAction holds the current action identifier between functions.
var currentAction string

// isMac is set based on if the mac definition is set.
var isMac = false

// parameterDefinition is used to define an actions parameters and to check against collected argument values.
type parameterDefinition struct {
	name         string
	validType    tokenType
	key          string
	defaultValue actionArgument
	optional     bool
	infinite     bool
}

// actionArgument is a variableValue value used to define collected argument values by the parser.
type actionArgument struct {
	valueType tokenType
	value     any
}

// action is a variableValue value that represents a collected action and arguments.
type action struct {
	// ident is the identifier of the action collected (e.g. identifier(...)).
	ident string
	// args are each of the arguments collected between the actions' parenthesis.
	args []actionArgument
}

// argsFunc is a function that can be passed a collected actions arguments as a slice of actionArgument.
type argsFunc func(args []actionArgument)

// argsFunc is a function that can be passed a collected actions arguments as a slice of actionArgument that must return a slice of plistData as a result.
type paramsFunc func(args []actionArgument) []plistData

// actionDefinition defines an action, what it expects and has functions for checking the arguments and creating the parameters.
type actionDefinition struct {
	identifier    string
	appIdentifier string
	parameters    []parameterDefinition
	check         argsFunc
	make          paramsFunc
	addParams     paramsFunc
	outputType    tokenType
	mac           bool
	minVersion    float64
}

// actions is the data structure that determines every action the compiler knows about.
// The key determines the identifier of the identifier that must be used in the syntax, it's value defines its behavior, etc. using an actionDefinition.
var actions map[string]actionDefinition

// libraryDefinition defines a 3rd-party actions library that can be imported using the `#import` syntax.
type libraryDefinition struct {
	identifier string
	// make is the function to call to add the actions in this library to the actions map.
	make func(identifier string)
}

// libraries is a map of the 3rd party libraries defined in the compiler.
// The key determines the identifier of the identifier name that must be used in the syntax, it's value defines its behavior, etc. using an libraryDefinition.
var libraries map[string]libraryDefinition

// callAction builds an action based on its actionDefinition and adds it to the shortcutActions map which makePlist will use to build the actions section of the Shortcut file format.
func callAction(arguments []actionArgument, outputName plistData, actionUUID plistData) {
	for i, a := range arguments {
		if a.valueType == Question {
			var lowerIdentifier = strings.ToLower(a.value.(string))
			if _, found := questions[lowerIdentifier]; found {
				var q = questions[lowerIdentifier]
				var parameter = actions[currentAction].parameters[i]
				questions[lowerIdentifier] = question{
					parameter:    parameter.key,
					actionIndex:  len(shortcutActions),
					text:         q.text,
					defaultValue: q.defaultValue,
				}
				arguments[i].value = ""
			}
		}
	}
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
	if actions[currentAction].addParams != nil {
		var addParams = actions[currentAction].addParams(arguments)
		for _, param := range addParams {
			params = append(params, param)
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

// makeAction constructs the action for the plist using ident and params.
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

// makeStdAction is an alias of makeAction that simply prepends the shortcuts bundle identifier to ident.
func makeStdAction(ident string, params []plistData) string {
	ident = "is.workflow.actions." + ident
	return makeAction(ident, params)
}

// checkAction checks the parsed arguments provided for an action and if it can be used based on definitions set.
// If an action has a check function defined this will be called and provided the parsed arguments.
func checkAction(arguments []actionArgument) {
	if len(actions[currentAction].parameters) > 0 {
		checkArgs(arguments)
		checkTypes(arguments, actions[currentAction].parameters)
	}
	if actions[currentAction].check != nil {
		actions[currentAction].check(arguments)
	}
	if actions[currentAction].minVersion != 0 {
		if actions[currentAction].minVersion > iosVersion {
			parserError(
				fmt.Sprintf("Action '%s()' is not available in set minimum version '%.1f'", currentAction, math.Ceil(iosVersion)),
			)
		}
	}
	if !isMac && actions[currentAction].mac {
		parserError(
			fmt.Sprintf("You've set your Shortcut as non-Mac. Action '%s()' is a Mac only action.", currentAction),
		)
	}
}

// checkEnum checks an argument value against a string slice
func checkEnum(name string, enum []string, args []actionArgument, idx int) {
	if len(args) < idx {
		return
	}
	if args[idx].value == nil {
		return
	}
	var value = getArgValue(args[idx]).(string)
	if !contains(enum, value) {
		var enumList string
		for _, e := range enum {
			enumList += "- " + e + "\n"
		}
		parserError(fmt.Sprintf("Invalid %s of '%s'.\n\nAvailable %ss:\n%s\nValues must be in the exact case listed to work properly.", name, value, name, enumList))
	}
}

// realVariableValue recurses to get the real value of a variable given its name
func realVariableValue(varName string, lastValueType tokenType) (varValue variableValue) {
	if _, global := globals[varName]; global {
		varValue = globals[varName]
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

// checkTypes iterates through `arguments` against `checks` to determine if the valid type defined
// for an action argument is the same as the type of the argument that was parsed.
func checkTypes(arguments []actionArgument, checks []parameterDefinition) {
	for i, check := range checks {
		if len(arguments) > i {
			typeCheck(check.name, check.validType, arguments[i])
		}
	}
}

// validActionOutput checks the output of an action in the case that the output has been assigned to a variable
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

// typeCheck is used to check the types of arguments given for actions
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
			parserError(fmt.Sprintf("Invalid variable value '%v' of type '%s' for argument '%s' of type '%s' in '%s()'",
				argVal,
				typeName(argValueType),
				field,
				typeName(validType),
				currentAction,
			))
		}
	case argValueType == Question:
	case argument.valueType != validType:
		parserError(fmt.Sprintf("Invalid argument value '%v' of type '%s' for argument '%s' of type '%s' in '%s()'",
			argVal,
			typeName(argValueType),
			field,
			typeName(validType),
			currentAction,
		))
	}
}

// getArgValue recurses to find the actual value of an argument
// in the case that the argument is a variable
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

// checkArgs checks to ensure all the required arguments for an action were entered.
func checkArgs(arguments []actionArgument) {
	var actionParams = actions[currentAction].parameters
	for i, param := range actionParams {
		if param.infinite {
			break
		}
		if i+1 > len(arguments) && !param.optional {
			var argIndex = i + 1
			var suffix = "th"
			switch argIndex {
			case 1:
				suffix = "st"
			case 2:
				suffix = "nd"
			case 3:
				suffix = "rd"
			}
			parserError(fmt.Sprintf("Missing required %d%s argument '%s' for action '%s'", argIndex, suffix, param.name, currentAction))
		}
		var realValue = getArgValue(arguments[i])
		if param.defaultValue.value == realValue {
			var argumentPlacement = currentAction + "("
			for argIndex := 0; argIndex < i+1; argIndex++ {
				if argIndex == i {
					argumentPlacement += fmt.Sprintf("%s = %v", param.name, arguments[i].value)
				} else {
					argumentPlacement += "..."
				}
				if argIndex < len(actionParams)-1 {
					argumentPlacement += ","
				}
			}
			argumentPlacement += ")"
			parserWarning(
				fmt.Sprintf(
					"Value for action argument is the same as the default value\n%s.",
					argumentPlacement,
				),
			)
		}
	}
}

func makeLibraries() {
	libraries = make(map[string]libraryDefinition)
}
