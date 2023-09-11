/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"embed"
	"fmt"
	"math"
	"reflect"
	"strconv"
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
	defaultValue any
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
			if len(arguments) <= i || len(arguments) == 0 {
				break
			}
			if arguments[i].valueType == Nil {
				continue
			}

			if a.validType == Variable {
				params = append(params, variableInput(a.key, arguments[i].value.(string)))
			} else {
				params = append(params, argumentValue(a.key, arguments, i))
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
func makeAction(ident string, params []plistData) plistData {
	return plistData{
		dataType: Dictionary,
		value: []plistData{
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
		},
	}
}

// makeStdAction is an alias of makeAction that simply prepends the shortcuts bundle identifier to ident.
func makeStdAction(ident string, params []plistData) plistData {
	ident = "is.workflow.actions." + ident
	return makeAction(ident, params)
}

// checkAction checks the parsed arguments provided for an action and if it can be used based on definitions set.
// If an action has a check function defined this will be called and provided the parsed arguments.
func checkAction() {
	var action = actions[currentAction]
	if len(action.parameters) > 0 {
		checkRequiredArgs()

		for i, param := range actions[currentAction].parameters {
			if !param.infinite {
				continue
			}
			checkInfiniteArgs(i)
			break
		}
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

func checkInfiniteArgs(startIdx int) {
	for i, arg := range currentArguments {
		if i < startIdx {
			continue
		}
		checkArg(&actions[currentAction].parameters[startIdx], &arg)
	}
}

// checkRequiredArgs checks if all required arguments for an action have a value.
func checkRequiredArgs() {
	for i, param := range actions[currentAction].parameters {
		if param.infinite {
			return
		}
		if i+1 > currentArgumentsSize && !param.optional && param.defaultValue == nil {
			var argIndex = i + 1
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
			parserError(fmt.Sprintf("Missing required %d%s argument '%s' for action '%s'.\n%s", argIndex, suffix, param.name, currentAction, generateActionDefinition(param, false, true)))
		}
	}
}

// checkEnum checks an argument value against a string slice.
func checkEnum(param parameterDefinition, argument actionArgument) {
	var value = getArgValue(argument).(string)
	if !contains(param.enum, value) {
		parserError(
			fmt.Sprintf(
				"Invalid argument '%s' for %s.\n\n%s",
				value,
				param.name,
				generateActionDefinition(param, false, true),
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

// typeCheck is used to check the types of arguments given for actions.
func typeCheck(param *parameterDefinition, argument *actionArgument) {
	var argValueType = argument.valueType
	var argVal = argument.value
	switch {
	case argValueType == Action:
		validActionOutput(param.name, param.validType, argVal)
		return
	case argValueType == Variable:
		var varName = argVal.(string)
		var getVar = realVariableValue(varName, String)
		argValueType = getVar.valueType
		argVal = getVar.value
		if argValueType == Action {
			validActionOutput(param.name, param.validType, argVal)
			return
		}
		if argValueType != param.validType && param.validType != Variable {
			parserError(fmt.Sprintf("Invalid variable value %v (%s) for argument '%s' (%s).\n%s",
				argVal,
				typeName(argValueType),
				param.name,
				typeName(param.validType),
				generateActionDefinition(*param, false, false),
			))
		}
	case argValueType == Question:
	case argValueType == Nil:
	case argument.valueType != param.validType:
		if argValueType == String {
			argVal = "\"" + argVal.(string) + "\""
		}
		parserError(fmt.Sprintf("Invalid value %v (%s) for argument '%s' (%s).\n%s",
			argVal,
			typeName(argValueType),
			param.name,
			typeName(param.validType),
			generateActionDefinition(*param, false, false),
		))
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

// incrementValue increments a string integer value.
func incrementValue(value any) string {
	var intValue, intValueErr = strconv.ParseInt(value.(string), 10, 64)
	handle(intValueErr)
	intValue++
	return fmt.Sprintf("%d", intValue)
}

// checkArg checks to ensure the collected argument for the current action is valid.
func checkArg(param *parameterDefinition, argument *actionArgument) {
	if param.enum != nil {
		checkEnum(*param, *argument)
	}
	typeCheck(param, argument)
	var realValue = getArgValue(*argument)
	var stringDefaultValue = fmt.Sprintf("%v", param.defaultValue)
	if param.defaultValue != nil && stringDefaultValue == realValue {
		parserWarning(
			fmt.Sprintf(
				"Value for action argument '%s' is the same as the default value.\n%s",
				param.name,
				generateActionDefinition(*param, false, false),
			),
		)
	}
}

func makeMeasurementUnits() {
	if len(units) != 0 {
		return
	}
	units = make(map[string][]string)
	units["Acceleration"] = []string{"m/s²", "g-force"}
	units["Angle"] = []string{"degrees", "arcminutes", "arcseconds", "radians", "grad", "revolutions"}
	units["Area"] = []string{"Mm²", "square kilometers", "square meters", "square centimeters", "mm²", "um²", "nm²", "square inches", "square feet", "square yards", "square miles", "acres", "a", "hectares"}
	units["Concentration Mass"] = []string{"g/L", "mg/dL", "µg/m³"}
	units["Dispersion"] = []string{"ppm"}
	units["Duration"] = []string{"milliseconds", "microseconds", "nanoseconds", "ps", "seconds", "minutes", "hours"}
	units["Electric Charge"] = []string{"C", "MAh", "kAh", "Ah", "mAh", "µAh"}
	units["Electric Current"] = []string{"MA", "kA", "amp", "mA", "µA"}
	units["Electric Potential Difference"] = []string{"MV", "kV", "volt", "mV", "µV"}
	units["Electric Resistance"] = []string{"MΩ", "kΩ", "ohm", "mΩ", "µΩ"}
	units["Energy"] = []string{"kJ", "joule", "kcal", "cal", "kWh"}
	units["Frequency"] = []string{"tHz", "GHz", "MHz", "kHz", "Hz", "mHz", "µHz", "nHz", "fps"}
	units["Fuel Efficiency"] = []string{"L/100km", "mpg"}
	units["Illuminance"] = []string{"lux"}
	units["Information Storage"] = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	units["Length"] = []string{"Mm", "km", "hm", "dam", "meter", "dm", "cm", "mm", "µm", "nm", "pm", "in", "ft", "yd", "mi", "smi", "ly", "nmi", "fathom", "furlong", "au", "parsec"}
	units["Mass"] = []string{"kg", "gram", "dg", "cg", "mg", "µg", "ng", "pg", "oz", "lb", "stone", "t", "ton", "carat", "oz t", "slug"}
	units["Power"] = []string{"TW", "GW", "MW", "kW", "watt", "mW", "µW", "nW", "pw", "fw", "hp"}
	units["Pressure"] = []string{"N/m²", "GPa", "MPa", "kPa", "hPa", "\" Hg", "bar", "mbar", "mm Hg", "psi"}
	units["Speed"] = []string{"m/s", "km/hr", "mi/hr", "kn"}
	units["Temperature"] = []string{"K", "ºC", "ºF"}
	units["Volume"] = []string{"ML", "kL", "liter", "dL", "cL", "mL", "km³", "m³", "dm³", "cm³", "mm³", "in³", "ft³", "yd³", "mi³", "acre ft", "bushel", "tsp", "tbsp", "fl oz", "pt", "qt", "Imp gal", "mcup"}
}

func generateActionDefinition(focus parameterDefinition, restrictions bool, showEnums bool) (definition string) {
	var action = actions[currentAction]
	definition += currentAction + "("
	for i, param := range action.parameters {
		if i != 0 && i < len(action.parameters) {
			definition += ", "
		}
		if focus.name != "" {
			if param.name == focus.name {
				definition += generateActionParamDefinition(param)
			} else {
				definition += "..."
			}
			continue
		}
		definition += generateActionParamDefinition(param)
	}
	definition += ")"
	if restrictions && (action.minVersion != 0 || action.mac) {
		definition += generateActionRestrictions()
	}
	if showEnums {
		definition += generateActionParamEnums(focus)
	}
	return definition
}

func generateActionRestrictions() (definition string) {
	definition += "\nRestrictions: "
	if actions[currentAction].minVersion != 0 {
		definition += fmt.Sprintf("iOS %1.f+", actions[currentAction].minVersion)
	}
	if actions[currentAction].minVersion != 0 && actions[currentAction].mac {
		definition += ", "
	}
	if actions[currentAction].mac {
		definition += "macOS only"
	}
	return
}

func generateActionParamEnums(focus parameterDefinition) (definition string) {
	var hasEnum = false
	for _, param := range actions[currentAction].parameters {
		if param.enum == nil {
			continue
		}
		if focus.name != "" && focus.name != param.name {
			continue
		}
		hasEnum = true
		definition += "\n\nAvailable " + param.name + "s:\n"
		for _, e := range param.enum {
			definition += "- " + e + "\n"
		}
	}
	if hasEnum {
		definition += "\nNote: Enum values are case-sensitive."
	}
	return
}

func generateActionParamDefinition(param parameterDefinition) (definition string) {
	if param.enum == nil {
		definition += typeName(param.validType) + " "
	} else {
		definition += "enum "
	}
	if param.infinite {
		definition += "..."
	}
	if param.optional || param.defaultValue != nil {
		definition += "?"
	}
	definition += param.name
	if param.defaultValue != nil {
		if reflect.TypeOf(param.defaultValue).String() == stringType {
			definition += fmt.Sprintf(" = \"%v\"", param.defaultValue)
		} else {
			definition += fmt.Sprintf(" = %v", param.defaultValue)
		}
	}
	return
}

// makeLibraries makes the library variable, this is where 3rd party action library definitions will start.
func makeLibraries() {
	libraries = make(map[string]libraryDefinition)
}
