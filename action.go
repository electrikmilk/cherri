/*
 * Copyright (c) Cherri
 */

package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/electrikmilk/args-parser"
)

//go:embed stdlib.cherri
var stdLib embed.FS

//go:embed actions
var stdActions embed.FS

// currentAction holds the current action definition between functions.
var currentAction actionDefinition
var currentActionIdentifier string
var currentArguments []actionArgument
var currentArgumentsSize int

// parameterDefinition is used to define an actions parameters and to check against collected argument values.
type parameterDefinition struct {
	name         string
	validType    tokenType
	key          string
	defaultValue any
	enum         string
	optional     bool
	infinite     bool
	literal      bool
}

// actionArgument is a varValue value used to define collected argument values by the parser.
type actionArgument struct {
	valueType tokenType
	value     any
}

// action is a varValue value that represents a collected action and arguments.
type action struct {
	// ident is the identifier of the action collected (e.g. identifier(...)).
	ident string
	// args are each of the arguments collected between the actions' parenthesis.
	args []actionArgument
}

// checkFunc is a function that can be passed a collected actions arguments as a slice of actionArgument and the current action's definition.
type checkFunc func(args []actionArgument, definition *actionDefinition)

// paramsFunc is a function that can be passed a collected actions arguments as a slice of actionArgument that must return action params as a result.
type paramsFunc func(args []actionArgument) map[string]any

type appIntent struct {
	name                string
	bundleIdentifier    string
	appIntentIdentifier string
}

// actionDefinition defines an action, what it expects and has functions for checking the arguments and creating the parameters.
type actionDefinition struct {
	identifier         string
	appIdentifier      string
	overrideIdentifier string
	parameters         []parameterDefinition
	check              checkFunc
	make               paramsFunc
	decomp             func(action *ShortcutAction) (arguments []string)
	addParams          paramsFunc
	appIntent          appIntent
	outputType         tokenType
	defaultAction      bool // Default action for this identifier during decompilation.
	mac                bool
	minVersion         float64
	maxVersion         float64
	setKey             string
	builtin            bool // builtin is based on if the action was in the actions map when it was first initialized.
	doc                selfDoc
}

type selfDoc struct {
	title       string
	description string
}

// libraryDefinition defines a 3rd-party actions library that can be imported using the `#import` syntax.
type libraryDefinition struct {
	identifier string
	// make is the function to call to add the actions in this library to the actions map.
	make func(identifier string)
}

var enumerations = map[string][]string{
	"measurementUnitType":          {"Acceleration", "Angle", "Area", "Concentration Mass", "Dispersion", "Duration", "Electric Charge", "Electric Current", "Electric Potential Difference", "V Electric Resistance", "Energy", "Frequency", "Fuel Efficiency", "Illuminance", "Information Storage", "Length", "Mass", "Power", "Pressure", "Speed", "Temperature", "Volume"},
	"storageUnit":                  {"bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"},
	"inputType":                    {"Text", "Number", "URL", "Date", "Time", "Date and Time"},
	"appSplitRatio":                {"half", "thirdByTwo"},
	"httpMethod":                   {"POST", "PUT", "PATCH", "DELETE"},
	"sortOrder":                    {"asc", "desc"},
	"windowSorting":                {"Title", "App Name", "Width", "Height", "X Position", "Y Position", "Window Index", "Name", "Random"},
	"timerDurations":                {"hr", "min", "sec"},
	"fileLabel":                    {"red", "orange", "yellow", "green", "blue", "purple", "gray"},
	"filesSortBy":                   {"File Size", "File Extension", "Creation Date", "File Path", "Last Modified Date", "Name", "Random"},
	"seekBehavior":                  {"To Time", "Forward By", "Backward By"},
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
	"Length":                        {"Mm", "km", "hm", "dam", "m", "dm", "cm", "mm", "µm", "nm", "pm", "in", "ft", "yd", "mi", "smi", "ly", "nmi", "fathom", "furlong", "au", "parsec"},
	"Mass":                          {"kg", "gram", "dg", "cg", "mg", "µg", "ng", "pg", "oz", "lb", "stone", "t", "ton", "carat", "oz t", "slug"},
	"Power":                         {"TW", "GW", "MW", "kW", "watt", "mW", "µW", "nW", "pw", "fw", "hp"},
	"Pressure":                      {"N/m²", "GPa", "MPa", "kPa", "hPa", "\" Hg", "bar", "mbar", "mm Hg", "psi"},
	"Speed":                         {"m/s", "km/hr", "mi/hr", "kn"},
	"Temperature":                   {"K", "ºC", "ºF"},
	"Volume":                        {"ML", "kL", "liter", "dL", "cL", "mL", "km³", "m³", "dm³", "cm³", "mm³", "in³", "ft³", "yd³", "mi³", "acre ft", "bushel", "tsp", "tbsp", "fl oz", "pt", "qt", "Imp gal", "mcup"},
	"fileOrderBy":                   {"Smallest First", "Biggest First", "Latest First", "Oldest First", "A to Z", "Z to A"},
}

var actionIndex int

// setCurrentAction sets the current action identifier and definition for use between functions.
func setCurrentAction(identifier string, definition *actionDefinition) {
	currentActionIdentifier = identifier
	currentAction = *definition
}

// undefinable checks if the current action cannot be defined using only Cherri because of the way it is defined.
func undefinable() bool {
	if currentAction.addParams != nil {
		var addedParams = currentAction.addParams([]actionArgument{})
		if len(addedParams) == 0 {
			return true
		}
	}

	return currentAction.builtin || currentAction.make != nil || currentAction.check != nil || currentAction.decomp != nil || currentAction.appIntent != emptyAppIntent
}

// makeAction builds an action based on its actionDefinition and adds it to the shortcut.
func makeAction(arguments []actionArgument, reference *map[string]any) {
	actionIndex++
	// Determine identifier
	var ident = getFullActionIdentifier()
	// Determine parameters
	var params = getActionParameters(arguments)
	// Additionally add the output name and UUID of this action if provided
	addAction(ident, attachReferenceToParams(&params, reference))
}

// getFullActionIdentifier determines the full identifier of currentAction.
func getFullActionIdentifier() (ident string) {
	if currentAction.overrideIdentifier != "" {
		return currentAction.overrideIdentifier
	}

	ident = "is.workflow.actions"
	if currentAction.appIdentifier != "" {
		ident = currentAction.appIdentifier
	}
	if currentAction.identifier != "" {
		ident = fmt.Sprintf("%s.%s", ident, currentAction.identifier)
	} else {
		ident = fmt.Sprintf("%s.%s", ident, strings.ToLower(currentActionIdentifier))
	}
	return
}

// getActionIdentifier determines the identifier of currentAction.
func getActionIdentifier() (ident string) {
	if currentAction.appIdentifier != "" {
		ident = fmt.Sprintf("%s.", currentAction.appIdentifier)
	}
	if currentAction.identifier != "" {
		ident = fmt.Sprintf("%s%s", ident, currentAction.identifier)
	} else {
		ident = fmt.Sprintf("%s%s", ident, strings.ToLower(currentActionIdentifier))
	}
	return
}

var emptyAppIntent = appIntent{}

// getActionParameters creates the actions' parameters by injecting the values of the arguments into the defined parameters.
func getActionParameters(arguments []actionArgument) map[string]any {
	var params = make(map[string]any)
	if currentAction.addParams != nil {
		maps.Copy(params, currentAction.addParams(arguments))
	}
	if currentAction.appIntent != emptyAppIntent {
		maps.Copy(params, appIntentDescriptor(currentAction.appIntent))
	}
	if currentAction.make != nil {
		return currentAction.make(arguments)
	}
	if currentAction.parameters != nil {
		var argumentsSize = len(arguments)
		if argumentsSize == 0 {
			return params
		}
		for i, param := range currentAction.parameters {
			if argumentsSize <= i {
				return params
			}
			if arguments[i].valueType == Nil || param.key == "" {
				continue
			}
			if param.validType == Variable {
				params[param.key] = variableValue(arguments[i].value.(varValue))
				continue
			}

			params[param.key] = argumentValue(arguments, i)
		}
	}

	return params
}

// addStdAction is an alias of addAction that simply prepends the shortcuts bundle identifier to ident.
func addStdAction(ident string, params *map[string]any) {
	addAction(fmt.Sprintf("is.workflow.actions.%s", ident), params)
}

// addAction adds an action to the shortcut.
func addAction(identifier string, params *map[string]any) {
	shortcut.WFWorkflowActions = append(shortcut.WFWorkflowActions,
		ShortcutAction{
			WFWorkflowActionIdentifier: identifier,
			WFWorkflowActionParameters: *params,
		},
	)
}

// checkAction checks the parsed arguments provided for an action and if it can be used based on definitions set.
// If an action has a check function defined this will be called and provided the parsed arguments.
func checkAction() {
	if len(currentAction.parameters) > 0 {
		checkRequiredArgs()
	}
	if currentAction.check != nil {
		currentAction.check(currentArguments, &currentAction)
	}
	if currentAction.minVersion != 0 {
		if currentAction.minVersion > iosVersion {
			parserError(
				fmt.Sprintf("Action '%s()' is not available in set minimum version '%.1f'", currentActionIdentifier, math.Ceil(iosVersion)),
			)
		}
	}
	if currentAction.maxVersion != 0 {
		parserWarning(fmt.Sprintf("Action '%s()' has been deprecated as it was removed or significantly modified.", currentActionIdentifier))
		if currentAction.maxVersion < iosVersion {
			parserError(
				fmt.Sprintf("Action '%s()' is not available in set minimum version '%.1f'", currentActionIdentifier, math.Ceil(iosVersion)),
			)
		}
	}
	if isMac, found := definitions["mac"]; found {
		if !isMac.(bool) && currentAction.mac {
			parserError(
				fmt.Sprintf("macOS action '%s()' in non-macOS Shortcut.", currentActionIdentifier),
			)
		} else if isMac.(bool) && !currentAction.mac {
			parserError(
				fmt.Sprintf("Non-macOS action '%s()' in macOS-only Shortcut.", currentActionIdentifier),
			)
		}
	}
}

func checkInfiniteArgs(startIdx int) {
	for i, arg := range currentArguments {
		if i < startIdx {
			continue
		}
		checkArg(&currentAction.parameters[startIdx], &arg)
	}
}

// checkRequiredArgs checks if all required arguments for an action have a value.
func checkRequiredArgs() {
	for i, param := range currentAction.parameters {
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
			parserError(fmt.Sprintf("Missing required %d%s argument '%s'.\n%s", argIndex, suffix, param.name, generateActionDefinition(param, true)))
		}
	}
}

func getEnum(identifier string) []string {
	if enumerations[identifier] == nil {
		return []string{}
	}

	return enumerations[identifier]
}

// checkEnum checks an argument value against a string slice.
func checkEnum(param *parameterDefinition, argument *actionArgument) {
	var value = getArgValue(*argument)
	if value == nil || reflect.TypeOf(value).Kind() != reflect.String || argument.valueType == Question {
		return
	}
	if !slices.Contains(getEnum(param.enum), value.(string)) {
		parserError(
			fmt.Sprintf(
				"Invalid argument '%s' for %s.\n\n%s",
				value,
				param.name,
				generateActionDefinition(*param, true),
			),
		)
	}
}

// realVariableValue recurses to get the real value of a variable given its name.
func realVariableValue(identifier string, lastValueType tokenType) (variableValue varValue) {
	if _, global := globals[identifier]; global {
		variableValue = globals[identifier]
		return
	}
	var front = strings.Split(identifier, "[")[0]
	if _, found := variables[front]; !found {
		parserError(fmt.Sprintf("Variable or Global '%s' does not exist", identifier))
	}
	var argValueType = variables[front].valueType
	var value = variables[front].value
	if argValueType == Variable {
		if lastValueType == Variable {
			parserError("Passed variable value that evaluates to variable")
		}
		if value != nil {
			variableValue = realVariableValue(value.(varValue).value.(string), argValueType)
		}
	} else {
		variableValue = variables[front]
	}
	return
}

func checkTypeTransform(valueType tokenType) tokenType {
	if valueType == Expression {
		valueType = Integer
	}

	return valueType
}

// typeCheck is used to check the types of arguments given for actions.
func typeCheck(param *parameterDefinition, argument *actionArgument) {
	var argValueType = checkTypeTransform(argument.valueType)
	var argVal = argument.value
	switch {
	case argValueType == Action:
		validActionOutput(param, argVal)
		return
	case argValueType == Variable:
		var identifier = argVal.(varValue).value.(string)
		var getVar = realVariableValue(identifier, String)
		argValueType = checkTypeTransform(getVar.valueType)
		argVal = getVar.value
		if argValueType == Action {
			validActionOutput(param, argVal)
			return
		}
		if argValueType != param.validType && param.validType != Variable && getVar.variableType != Ask {
			parserError(fmt.Sprintf("Invalid variable value %v (%s) for argument '%s' (%s).\n%s",
				argVal,
				argValueType,
				param.name,
				param.validType,
				generateActionDefinition(*param, false),
			))
		}
	case argValueType == Question:
	case argValueType == Nil:
	case param.validType == String && argument.valueType == RawString:
	case argValueType != param.validType:
		if argValueType == String {
			argVal = "\"" + argVal.(string) + "\""
		}
		parserError(fmt.Sprintf("Invalid value %v (%s) for argument '%s' (%s).\n%s",
			argVal,
			argValueType,
			param.name,
			param.validType,
			generateActionDefinition(*param, false),
		))
	}
}

// validActionOutput checks the output of an action in the case that the output has been assigned to a variable.
func validActionOutput(param *parameterDefinition, value any) {
	var actionIdent = value.(action).ident
	if action, found := actions[actionIdent]; found {
		var actionOutputType = action.outputType
		if actionOutputType == "" {
			return
		}
		if actionOutputType != param.validType && param.validType != Variable {
			parserError(fmt.Sprintf("Invalid variable value of action '%v' (%s) for argument '%s' (%s).\n%s",
				actionIdent+"()",
				actionOutputType,
				param.name,
				param.validType,
				generateActionDefinition(*param, false),
			))
		}
	}
}

// getArgValue recurses to find the actual value of an argument
// in the case that the argument is a variable.
func getArgValue(argument actionArgument) any {
	if argument.valueType != Variable {
		return argument.value
	}
	if argument.value == nil {
		return nil
	}

	var identifier = argument.value.(varValue).value.(string)
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

// checkArg checks to ensure the collected argument for the current action is valid.
func checkArg(param *parameterDefinition, argument *actionArgument) {
	if argument.valueType == Variable && argument.value == Ask {
		return
	}

	if param.enum != "" {
		checkEnum(param, argument)
	}

	typeCheck(param, argument)

	questionArg(param, argument)

	if param.literal {
		checkLiteralValue(param, argument)
	}

	var realValue = getArgValue(*argument)
	var stringDefaultValue = fmt.Sprintf("%s", param.defaultValue)
	if param.defaultValue != nil && stringDefaultValue == realValue {
		parserWarning(
			fmt.Sprintf(
				"Value for action argument '%s' is the same as the default value.\n%s",
				param.name,
				generateActionDefinition(*param, false),
			),
		)
	}
}

// questionArg checks if the argument references a question so that it can update the question to point to the current action's argument.
func questionArg(param *parameterDefinition, argument *actionArgument) {
	if argument.valueType != Question {
		return
	}
	var identifier = argument.value.(string)
	if question, found := questions[identifier]; found {
		question.parameter = param.key
		question.actionIndex = actionIndex
		argument.value = ""
	}
}

func checkLiteralValue(param *parameterDefinition, argument *actionArgument) {
	if argument.valueType != param.validType {
		parserError(fmt.Sprintf(
			"Shortcuts does not allow variables for this argument, use a literal for the argument value.\n\n%s",
			generateActionDefinition(*param, false),
		))
	}
}

func generateActionDefinition(focus parameterDefinition, showEnums bool) string {
	var definition strings.Builder
	definition.WriteRune('\n')

	var docTitle = currentAction.doc.title
	if currentAction.doc.title == "" {
		if args.Using("no-ansi") {
			docTitle = fmt.Sprintf("`%s()`", currentActionIdentifier)
		} else {
			docTitle = fmt.Sprintf("%s()", currentActionIdentifier)
		}
	}
	if args.Using("no-ansi") {
		definition.WriteString("### ")
	}
	definition.WriteString(fmt.Sprintf("%s\n", ansi(docTitle, bold, underline)))
	definition.WriteRune('\n')
	if currentAction.doc.description != "" {
		definition.WriteString(fmt.Sprintf("%s\n\n", ansi(currentAction.doc.description, italic)))
	}

	if args.Using("no-ansi") {
		definition.WriteString("```\n")
	}

	if showEnums {
		definition.WriteString(generateActionParamEnums(focus))
	}

	var definitionType string
	if currentAction.builtin {
		definitionType = "#builtin"
	} else {
		definitionType = "#define"
	}

	definition.WriteString(ansi(fmt.Sprintf("%s action ", definitionType), orange))

	if currentAction.defaultAction {
		definition.WriteString(ansi("default ", yellow))
	}
	if currentAction.mac {
		definition.WriteString(ansi("mac ", orange))
	}
	if currentAction.minVersion != 0 {
		definition.WriteString(ansi(fmt.Sprintf("v%1.f>", currentAction.minVersion), underline))
		definition.WriteRune(' ')
	}
	if currentAction.maxVersion != 0 {
		definition.WriteString(ansi(fmt.Sprintf("v%1.f<", currentAction.maxVersion), red, bold))
		definition.WriteRune(' ')
	}

	if currentAction.identifier != "" || currentAction.appIdentifier != "" {
		setCurrentAction(currentActionIdentifier, &currentAction)
		definition.WriteString(ansi(fmt.Sprintf("'%s' ", getActionIdentifier()), red))
	}

	definition.WriteString(fmt.Sprintf("%s(", ansi(currentActionIdentifier, blue, bold)))
	var arguments []string
	for _, param := range currentAction.parameters {
		if param.name == focus.name || focus.name == "" {
			arguments = append(arguments, generateActionParamDefinition(param))
		} else {
			arguments = append(arguments, ansi("...", dim))
		}
	}
	definition.WriteString(strings.Join(arguments, ", "))
	definition.WriteRune(')')

	if currentAction.outputType != "" {
		definition.WriteString(fmt.Sprintf(": %s", ansi(string(currentAction.outputType), magenta)))
	}

	if currentAction.addParams != nil {
		var addParams = currentAction.addParams([]actionArgument{})
		if len(addParams) != 0 {
			var jsonBytes, jsonErr = json.MarshalIndent(addParams, strings.Repeat("\t", tabLevel), "\t")
			handle(jsonErr)
			definition.WriteString(fmt.Sprintf(" %s", string(jsonBytes)))
		}
	}

	if args.Using("no-ansi") {
		definition.WriteString("\n```")
	}

	return definition.String()
}

func generateActionParamEnums(focus parameterDefinition) string {
	var definition strings.Builder
	for _, param := range currentAction.parameters {
		if param.enum == "" {
			continue
		}
		if focus.name != "" && focus.name != param.name {
			continue
		}
		definition.WriteString(ansi("enum ", orange))
		definition.WriteString(param.enum)
		definition.WriteString(ansi(" {\n", dim))
		for i, enum := range getEnum(param.enum) {
			var enumSize = len(param.enum)
			definition.WriteString(ansi(fmt.Sprintf("\t'%s'", enum), orange))
			if i < enumSize+1 {
				definition.WriteString(",\n")
			}
		}
		definition.WriteString(ansi("}\n\n", dim))
	}

	return definition.String()
}

func generateActionParamDefinition(param parameterDefinition) string {
	var definition strings.Builder
	var argType string
	if param.enum == "" {
		argType = fmt.Sprintf("%s ", param.validType)
	} else {
		argType = fmt.Sprintf("%s ", param.enum)
	}
	definition.WriteString(ansi(argType, magenta))

	if param.infinite {
		definition.WriteString("...")
	}
	if param.optional || param.defaultValue != nil {
		definition.WriteRune('?')
	}
	definition.WriteString(param.name)

	if param.key != "" && param.key != param.name {
		definition.WriteString(ansi(fmt.Sprintf(": '%s'", param.key), orange))
	}

	if param.defaultValue != nil {
		definition.WriteString(ansi(" = ", dim))
		var defaultValue string
		if reflect.TypeOf(param.defaultValue).Kind() == reflect.String {
			defaultValue = fmt.Sprintf("\"%v\"", strings.Replace(param.defaultValue.(string), "\n", "\\n", 1))
		} else {
			defaultValue = fmt.Sprintf("%v", param.defaultValue)
		}
		definition.WriteString(ansi(defaultValue, green))
	}

	return definition.String()
}

// makeLibraries makes the library variable, this is where 3rd party action library definitions will start.
func makeLibraries() {
	libraries = make(map[string]libraryDefinition)
}

func appIntentDescriptor(intent appIntent) map[string]any {
	return map[string]any{
		"AppIntentDescriptor": map[string]string{
			"TeamIdentifier":      "0000000000",
			"BundleIdentifier":    intent.bundleIdentifier,
			"Name":                intent.name,
			"AppIntentIdentifier": intent.appIntentIdentifier,
		},
	}
}

// handleActionDefinitions parses defined actions in the current file and collects them into the actions map.
func handleActionDefinitions() {
	if !regexp.MustCompile(`#define action (?:'(.+)')?(.*?)\(`).MatchString(contents) && !regexp.MustCompile(`enum (.*?) \{`).MatchString(contents) {
		return
	}
	parseActionDefinitions()

	resetParse()
}

func parseActionDefinitions() {
	for char != -1 {
		switch {
		case isChar('/'):
			args.Args["comments"] = ""
			preParsing = false
			collectComment()
			preParsing = true
			delete(args.Args, "comments")
		case tokenAhead(Enumeration):
			collectEnumeration()
		case tokenAhead(Definition):
			advance()
			if tokenAhead(Action) {
				advance()
				collectDefinedAction()
				continue
			}
		}
		advance()
	}
	tokens = []token{}
}

func collectDefinedAction() {
	var lineRef = newLineReference()

	var doc selfDoc
	var lastToken = getLastAddedToken()
	if lastToken.typeof == Comment {
		var comment = lastToken.value.(string)
		if !strings.Contains(comment, "\n") && strings.Contains(comment, ":") {
			var parts = strings.Split(comment, ":")
			if strings.TrimSpace(parts[0]) == "[Doc]" {
				if len(parts) > 2 {
					doc = selfDoc{title: strings.TrimSpace(parts[1]), description: strings.TrimSpace(parts[2])}
				} else {
					doc = selfDoc{description: strings.TrimSpace(parts[1])}
				}
			}
		}
	}

	var defaultAction bool
	if tokenAhead(Default) {
		defaultAction = true
		advance()
	}

	var macOnlyAction bool
	if tokenAhead(Mac) {
		macOnlyAction = true
		advance()
	} else if tokenAhead(NonMac) {
		macOnlyAction = false
		advance()
	}

	var minVersion, maxVersion = collectVersionDefinition()

	var shortIdentifier string
	var overrideIdentifier string
	if char == '\'' {
		advance()

		var workflowIdentifier = collectRawString()
		if len(strings.Split(workflowIdentifier, ".")) < 4 {
			shortIdentifier = workflowIdentifier
		} else {
			overrideIdentifier = workflowIdentifier
		}
		advance()
	}

	var identifier, arguments, outputType = collectActionDefinition('\n')
	if shortIdentifier == "" {
		shortIdentifier = strings.ToLower(identifier)
	}

	advance()

	var addParams paramsFunc
	if char == '{' {
		advance()
		var dict = collectDictionary()
		addParams = func(args []actionArgument) map[string]any {
			handleRawParams(dict.(map[string]interface{}))
			return dict.(map[string]any)
		}
	}

	lineRef.replaceLines()

	actions[identifier] = &actionDefinition{
		identifier:         shortIdentifier,
		overrideIdentifier: overrideIdentifier,
		parameters:         arguments,
		outputType:         outputType,
		addParams:          addParams,
		defaultAction:      defaultAction,
		mac:                macOnlyAction,
		minVersion:         minVersion,
		maxVersion:         maxVersion,
		doc:                doc,
	}

	if args.Using("debug") {
		setCurrentAction(identifier, actions[identifier])
		fmt.Println("\ndefined:", currentAction.appIdentifier, generateActionDefinition(parameterDefinition{}, true))
		fmt.Print("\n")
	}
}

func collectVersionDefinition() (minVersion float64, maxVersion float64) {
	for char != -1 && char != ' ' {
		if char != 'v' || char == 'v' && !intChar(next(1)) {
			break
		}
		advance()
		var valueType tokenType
		var version any
		collectIntegerValue(&valueType, &version, ' ')

		switch char {
		case '>':
			minVersion = version.(float64)
			advance()
		case '<':
			maxVersion = version.(float64)
			advance()
		default:
			minVersion = version.(float64)
		}
		skipWhitespace()
	}

	return minVersion, maxVersion
}

func collectActionDefinition(until rune) (identifier string, arguments []parameterDefinition, outputType tokenType) {
	identifier = collectIdentifier()
	if _, found := customActions[identifier]; found {
		parserError(fmt.Sprintf("Duplication declaration of custom action '%s()'", identifier))
	}
	if _, found := actions[identifier]; found {
		parserError(fmt.Sprintf("Duplication declaration of action '%s()'", identifier))
	}

	if next(1) != ')' {
		advance()
		skipWhitespace()
		arguments = collectParameterDefinitions()
	} else {
		advanceTimes(2)
	}

	if tokenAhead(Colon) {
		skipWhitespace()
		var value any
		collectType(&outputType, &value, until)
	}

	return
}

func collectParameterDefinitions() (arguments []parameterDefinition) {
	for char != ')' && char != -1 {
		var valueType tokenType
		var value any

		var enumeration string
		var ahead = lookAheadUntil(' ')
		if enumerations[ahead] != nil {
			enumeration = collectUntil(' ')
			valueType = String
		} else {
			collectType(&valueType, &value, ' ')
		}

		var literal bool
		if char == '!' {
			advance()
			literal = true
		}

		value = nil

		skipWhitespace()

		var optional bool
		if char == '?' {
			optional = true
			advance()
		}

		var identifier = collectIdentifier()

		var parameterKey string
		if char == ':' {
			advanceTimes(2)

			if char != '\'' {
				parserError("Expected parameter key raw string (').")
			}
			advance()
			parameterKey = collectRawString()
		} else {
			parameterKey = identifier
		}

		skipWhitespace()

		var defaultValue any
		switch char {
		case '=':
			advance()
			skipWhitespace()

			var defaultValueType tokenType
			collectValue(&defaultValueType, &defaultValue, endOfNextArgument())
			if defaultValueType != valueType {
				parserError(fmt.Sprintf("Invalid default value of type '%s' for '%s' type argument '%s'", defaultValueType, valueType, identifier))
			}
		}
		if char == ',' {
			advance()
		}

		arguments = append(arguments, parameterDefinition{
			name:         identifier,
			key:          parameterKey,
			validType:    valueType,
			optional:     optional,
			defaultValue: defaultValue,
			enum:         enumeration,
			literal:      literal,
		})

		skipWhitespace()
	}
	advance()

	return
}
