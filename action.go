/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"embed"
	"fmt"
	"math"
	"strings"
)

//go:embed stdlib.cherri
var stdLib embed.FS

// currentAction holds the current action identifier between functions.
var currentAction string
var currentArguments []actionArgument
var currentArgumentsSize int

// isMac is set based on if the mac definition is set.
var isMac = false

// parameterDefinition is used to define an actions parameters and to check against collected argument values.
type parameterDefinition struct {
	name         string
	validType    tokenType
	key          string
	defaultValue actionArgument
	enum         []string
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

// paramsFunc is a function that can be passed a collected actions arguments as a slice of actionArgument that must return a slice of plistData as a result.
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
var actions map[string]*actionDefinition

// libraryDefinition defines a 3rd-party actions library that can be imported using the `#import` syntax.
type libraryDefinition struct {
	identifier string
	// make is the function to call to add the actions in this library to the actions map.
	make func(identifier string)
}

// plistAction builds an action based on its actionDefinition and adds it to the shortcutActions map which makePlist will use to build the actions section of the Shortcut file format.
func plistAction(arguments []actionArgument, outputName plistData, actionUUID plistData) {
	// Check for question arguments
	questionArgs(arguments)
	// Determine identifier
	var ident = actionIdentifier()
	// Determine parameters
	var params = actionParameters(arguments)
	// Additionally add the output name and UUID of this action if provided
	if outputName.value != nil {
		params = append(params, outputName)
	}
	if actionUUID.value != nil {
		params = append(params, actionUUID)
	}
	shortcutActions = append(shortcutActions, makeAction(ident, params))
}

// actionIdentifier determines the identifier of currentAction.
func actionIdentifier() (ident string) {
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
	return
}

// actionParameters creates the actions' parameters by injecting the values of the arguments into the defined parameters.
func actionParameters(arguments []actionArgument) (params []plistData) {
	if actions[currentAction].make == nil && actions[currentAction].parameters != nil {
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
	if actions[currentAction].make != nil {
		params = actions[currentAction].make(arguments)
	}
	if actions[currentAction].addParams != nil {
		var addParams = actions[currentAction].addParams(arguments)
		params = append(params, addParams...)
	}
	return
}

// questionArgs updates questions to target the action parameter
// that it's identifier matches the arguments value.
func questionArgs(arguments []actionArgument) {
	for i, a := range arguments {
		if a.valueType != Question {
			continue
		}
		var lowerIdentifier = strings.ToLower(a.value.(string))
		if _, found := questions[lowerIdentifier]; found {
			var parameter = actions[currentAction].parameters[i]
			questions[lowerIdentifier].parameter = parameter.key
			questions[lowerIdentifier].actionIndex = len(shortcutActions)
			arguments[i].value = ""
		}
	}
}

// makeAction constructs the action for the plist using ident and params.
func makeAction(ident string, params []plistData) (action string) {
	action = plistValue(Dictionary, []plistData{
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
func checkAction() {
	var action = actions[currentAction]
	if len(action.parameters) > 0 {
		checkRequiredArgs(action.parameters)
		checkTypes(action.parameters)
	}
	if action.check != nil {
		action.check(currentArguments)
	}
	if action.minVersion != 0 {
		if action.minVersion > iosVersion {
			parserError(
				fmt.Sprintf("Action '%s()' is not available in set minimum version '%.1f'", currentAction, math.Ceil(iosVersion)),
			)
		}
	}
	if !isMac && action.mac {
		parserError(
			fmt.Sprintf("You've set your Shortcut as non-Mac. Action '%s()' is a Mac only action.", currentAction),
		)
	}
}

func checkRequiredArgs(params []parameterDefinition) {
	for i, param := range params {
		if param.infinite {
			return
		}
		if i+1 > currentArgumentsSize && !param.optional {
			var argIndex = idx + 1
			var suffix string
			switch argIndex {
			case 1:
				suffix = "st"
			case 2:
				suffix = "nd"
			case 3:
				suffix = "rd"
			default:
				suffix = "th"
			}
			parserError(fmt.Sprintf("Missing required %d%s argument '%s' for action '%s'", argIndex, suffix, param.name, currentAction))
		}
	}
}

// checkEnum checks an argument value against a string slice.
func checkEnum(name string, enum []string, argument actionArgument) {
	var value = getArgValue(argument).(string)
	if !contains(enum, value) {
		var enumList string
		for _, e := range enum {
			enumList += "- " + e + "\n"
		}
		parserError(
			fmt.Sprintf(
				"Invalid argument '%s' for %s.\n\nAvailable %ss:\n%s\nNote: Values must be in the exact case listed to work properly.",
				value,
				name,
				name,
				enumList,
			),
		)
	}
}

// realVariableValue recurses to get the real value of a variable given its name.
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
func checkTypes(checks []parameterDefinition) {
	for i, check := range checks {
		if currentArgumentsSize > i {
			typeCheck(check.name, check.validType, currentArguments[i])
		}
	}
}

// validActionOutput checks the output of an action in the case that the output has been assigned to a variable.
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

// typeCheck is used to check the types of arguments given for actions.
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
		if argValueType == String {
			argVal = "\"" + argVal.(string) + "\""
		}
		parserError(fmt.Sprintf("%s(): Invalid value '%v' (%s) for argument '%s' (%s).",
			currentAction,
			argVal,
			typeName(argValueType),
			field,
			typeName(validType),
		))
	}
}

// getArgValue recurses to find the actual value of an argument
// in the case that the argument is a variable.
func getArgValue(argument actionArgument) any {
	if argument.valueType != Variable {
		return argument.value
	}
	var variable = argument.value.(string)
	if _, found := variables[variable]; !found {
		return argument.value
	}
	if variables[variable].valueType == Variable {
		return getArgValue(actionArgument{
			valueType: variables[variable].valueType,
			value:     variables[variable].value,
		})
	}
	return variables[variable].value
}

// checkArg checks to ensure the collected argument for the current action is valid.
func checkArg(idx int, param parameterDefinition, argument actionArgument) {
	if param.infinite {
		return
	}
	if param.enum != nil {
		checkEnum(param.name, param.enum, argument)
	}
	if currentArgumentsSize < idx {
		return
	}
	var realValue = getArgValue(argument)
	if param.defaultValue.value != nil && param.defaultValue.value == realValue {
		var argumentPlacement = currentAction + "("
		for argIndex := 0; argIndex < idx+1; argIndex++ {
			if argIndex == idx {
				argumentPlacement += fmt.Sprintf("%s = %v", param.name, argument.value)
			} else {
				argumentPlacement += "..."
			}
			if argIndex < len(actions[currentAction].parameters)-1 {
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

// makeLibraries makes the library variable, this is where 3rd party action library definitions will start.
func makeLibraries() {
	libraries = make(map[string]libraryDefinition)
}
