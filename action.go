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
	maxVersion    float64
}

// actions is the data structure that determines every action the compiler knows about.
// The key determines the identifier of the identifier that must be used in the syntax, it's value defines its behavior, etc. using an actionDefinition.
var actions map[string]*actionDefinition

var usedActions []string

// libraryDefinition defines a 3rd-party actions library that can be imported using the `#import` syntax.
type libraryDefinition struct {
	identifier string
	// make is the function to call to add the actions in this library to the actions map.
	make func(identifier string)
}

var actionIndex int

// plistAction builds an action based on its actionDefinition and adds it to the plist.
func plistAction(arguments []actionArgument, outputName *plistData, actionUUID *plistData) {
	actionIndex++
	// Check for question arguments
	questionArgs(arguments)
	// Determine identifier
	var ident = actionIdentifier()
	// Determine parameters
	var params = actionParameters(arguments)
	// Additionally add the output name and UUID of this action if provided
	if outputName.value != nil {
		params = append(params, *outputName)
	}
	if actionUUID.value != nil {
		params = append(params, *actionUUID)
	}
	appendPlist(makeAction(ident, params))
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
	if actions[currentAction].addParams != nil {
		params = append(params, actions[currentAction].addParams(arguments)...)
	}
	if actions[currentAction].make != nil {
		params = actions[currentAction].make(arguments)
		return
	}
	if actions[currentAction].parameters != nil {
		var argumentsSize = len(arguments)
		if argumentsSize == 0 {
			return
		}
		for i, a := range actions[currentAction].parameters {
			if argumentsSize <= i {
				return
			}
			if arguments[i].valueType == Nil || a.key == "" {
				continue
			}
			if a.validType == Variable {
				params = append(params, variableInput(a.key, arguments[i].value.(string)))
				continue
			}

			params = append(params, argumentValue(a.key, arguments, i))
		}
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
		if question, found := questions[lowerIdentifier]; found {
			var parameter = actions[currentAction].parameters[i]
			question.parameter = parameter.key
			question.actionIndex = actionIndex
			arguments[i].value = ""
		}
	}
}

// makeAction constructs the action for the plist using ident and params.
func makeAction(ident string, params []plistData) []plistData {
	return []plistData{
		{
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
		},
	}
}

// makeStdAction is an alias of makeAction that simply prepends the shortcuts bundle identifier to ident.
func makeStdAction(ident string, params []plistData) []plistData {
	ident = "is.workflow.actions." + ident
	return makeAction(ident, params)
}

// checkAction checks the parsed arguments provided for an action and if it can be used based on definitions set.
// If an action has a check function defined this will be called and provided the parsed arguments.
func checkAction() {
	var action = actions[currentAction]
	if len(action.parameters) > 0 {
		checkRequiredArgs()
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
	if action.maxVersion != 0 {
		parserWarning(fmt.Sprintf("Action '%s()' has been deprecated as it was removed or significantly modified.", currentAction))
		if action.maxVersion < iosVersion {
			parserError(
				fmt.Sprintf("Action '%s()' is not available in set minimum version '%.1f'", currentAction, math.Ceil(iosVersion)),
			)
		}
	}
	if isMac, found := definitions["mac"]; found {
		if !isMac.(bool) && action.mac {
			parserError(
				fmt.Sprintf("You've set your Shortcut as non-Mac. Action '%s()' is a Mac only action", currentAction),
			)
		}
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
			checkInfiniteArgs(i)
			continue
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
	var value = getArgValue(argument)
	if value == nil {
		return
	}
	if !contains(param.enum, value.(string)) {
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
func realVariableValue(identifier string, lastValueType tokenType) (varValue variableValue) {
	if _, global := globals[identifier]; global {
		varValue = globals[identifier]
		return
	}
	if _, found := variables[identifier]; !found {
		parserError(fmt.Sprintf("Variable or Global '%s' does not exist", identifier))
	}
	var argValueType = variables[identifier].valueType
	var value = variables[identifier].value
	if argValueType == Variable {
		if lastValueType == Variable && argValueType == Variable {
			parserError("Passed variable value that evaluates to variable")
		}
		varValue = realVariableValue(value.(string), argValueType)
	} else {
		varValue = variables[identifier]
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
		var identifier = argVal.(string)
		var getVar = realVariableValue(identifier, String)
		argValueType = getVar.valueType
		argVal = getVar.value
		if argValueType == Action {
			validActionOutput(param.name, param.validType, argVal)
			return
		}
		if argValueType != param.validType && param.validType != Variable {
			parserError(fmt.Sprintf("Invalid variable value %v (%s) for argument '%s' (%s).\n%s",
				argVal,
				argValueType,
				param.name,
				param.validType,
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
			argValueType,
			param.name,
			param.validType,
			generateActionDefinition(*param, false, false),
		))
	}
}

// validActionOutput checks the output of an action in the case that the output has been assigned to a variable.
func validActionOutput(field string, validType tokenType, value any) {
	var actionIdent = value.(action).ident
	if action, found := actions[actionIdent]; found {
		var actionOutputType = action.outputType
		if actionOutputType != "" {
			if actionOutputType != validType {
				parserError(
					fmt.Sprintf(
						"Invalid variable value of action '%v' that outputs type '%s' for argument '%s' of type '%s' in '%s()'",
						actionIdent+"()",
						actionOutputType,
						field,
						validType,
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
	var identifier = argument.value.(string)
	if _, found := variables[identifier]; !found {
		return argument.value
	}
	if variables[identifier].valueType == Variable {
		return getArgValue(actionArgument{
			valueType: variables[identifier].valueType,
			value:     variables[identifier].value,
		})
	}
	return variables[identifier].value
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
	var stringDefaultValue = fmt.Sprintf("%s", param.defaultValue)
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
	units = map[string][]string{
		"Acceleration":                  {"m/s²", "g-force"},
		"Angle":                         {"degrees", "arcminutes", "arcseconds", "radians", "grad", "revolutions"},
		"Area":                          {"Mm²", "square kilometers", "square meters", "square centimeters", "mm²", "um²", "nm²", "square inches", "square feet", "square yards", "square miles", "acres", "a", "hectares"},
		"Concentration Mass":            {"g/L", "mg/dL", "µg/m³"},
		"Dispersion":                    {"ppm"},
		"Duration":                      {"milliseconds", "microseconds", "nanoseconds", "ps", "seconds", "minutes", "hours"},
		"Electric Charge":               {"C", "MAh", "kAh", "Ah", "mAh", "µAh"},
		"Electric Current":              {"MA", "kA", "amp", "mA", "µA"},
		"Electric Potential Difference": {"MV", "kV", "volt", "mV", "µV"},
		"Electric Resistance":           {"MΩ", "kΩ", "ohm", "mΩ", "µΩ"},
		"Energy":                        {"kJ", "joule", "kcal", "cal", "kWh"},
		"Frequency":                     {"tHz", "GHz", "MHz", "kHz", "Hz", "mHz", "µHz", "nHz", "fps"},
		"Fuel Efficiency":               {"L/100km", "mpg"},
		"Illuminance":                   {"lux"},
		"Information Storage":           {"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"},
		"Length":                        {"Mm", "km", "hm", "dam", "meter", "dm", "cm", "mm", "µm", "nm", "pm", "in", "ft", "yd", "mi", "smi", "ly", "nmi", "fathom", "furlong", "au", "parsec"},
		"Mass":                          {"kg", "gram", "dg", "cg", "mg", "µg", "ng", "pg", "oz", "lb", "stone", "t", "ton", "carat", "oz t", "slug"},
		"Power":                         {"TW", "GW", "MW", "kW", "watt", "mW", "µW", "nW", "pw", "fw", "hp"},
		"Pressure":                      {"N/m²", "GPa", "MPa", "kPa", "hPa", "\" Hg", "bar", "mbar", "mm Hg", "psi"},
		"Speed":                         {"m/s", "km/hr", "mi/hr", "kn"},
		"Temperature":                   {"K", "ºC", "ºF"},
		"Volume":                        {"ML", "kL", "liter", "dL", "cL", "mL", "km³", "m³", "dm³", "cm³", "mm³", "in³", "ft³", "yd³", "mi³", "acre ft", "bushel", "tsp", "tbsp", "fl oz", "pt", "qt", "Imp gal", "mcup"},
	}
}

func generateActionDefinition(focus parameterDefinition, restrictions bool, showEnums bool) string {
	var action = actions[currentAction]
	var definition strings.Builder
	definition.WriteString(fmt.Sprintf("%s(", currentAction))
	var arguments []string
	for _, param := range action.parameters {
		if param.name == focus.name || focus.name == "" {
			arguments = append(arguments, generateActionParamDefinition(param))
		} else {
			arguments = append(arguments, "...")
		}
	}
	definition.WriteString(strings.Join(arguments, ", "))
	definition.WriteRune(')')
	if restrictions && (action.minVersion != 0 || action.maxVersion != 0 || action.mac) {
		definition.WriteString(generateActionRestrictions())
	}
	if showEnums {
		definition.WriteString(generateActionParamEnums(focus))
	}

	return definition.String()
}

func generateActionRestrictions() string {
	var definition strings.Builder
	definition.WriteString("\nRestrictions: ")
	var restrictions []string
	if actions[currentAction].minVersion != 0 {
		restrictions = append(restrictions, fmt.Sprintf("iOS %1.f+", actions[currentAction].minVersion))
	}
	if actions[currentAction].maxVersion != 0 {
		restrictions = append(restrictions, fmt.Sprintf("Removed or significantly changed after iOS %1.f+", actions[currentAction].maxVersion))
	}
	if actions[currentAction].mac {
		restrictions = append(restrictions, "macOS only")
	}
	definition.WriteString(strings.Join(restrictions, ", "))

	return ansi(definition.String(), red, bold)
}

func generateActionParamEnums(focus parameterDefinition) string {
	var definition strings.Builder
	if len(actions[currentAction].parameters) != 0 {
		definition.WriteRune('\n')
	}
	var hasEnum = false
	for _, param := range actions[currentAction].parameters {
		if param.enum == nil {
			continue
		}
		if focus.name != "" && focus.name != param.name {
			continue
		}
		hasEnum = true
		definition.WriteString(ansi(fmt.Sprintf("\nAvailable %ss:\n", param.name), yellow))
		for _, e := range param.enum {
			definition.WriteString(fmt.Sprintf("- %s\n", e))
		}
	}
	if hasEnum {
		definition.WriteString(ansi("\nNote: Enum values are case-sensitive.", bold))
	}

	return definition.String()
}

func generateActionParamDefinition(param parameterDefinition) string {
	var definition strings.Builder
	if param.enum == nil {
		definition.WriteString(fmt.Sprintf("%s ", param.validType))
	} else {
		definition.WriteString("enum ")
	}
	if param.infinite {
		definition.WriteString("...")
	}
	if param.optional || param.defaultValue != nil {
		definition.WriteRune('?')
	}
	definition.WriteString(param.name)
	if param.defaultValue != nil {
		if reflect.TypeOf(param.defaultValue).String() == stringType {
			definition.WriteString(fmt.Sprintf(" = \"%v\"", param.defaultValue))
		} else {
			definition.WriteString(fmt.Sprintf(" = %v", param.defaultValue))
		}
	}

	return definition.String()
}

// makeLibraries makes the library variable, this is where 3rd party action library definitions will start.
func makeLibraries() {
	libraries = make(map[string]libraryDefinition)
}
