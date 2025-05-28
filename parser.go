/*
 * Copyright (c) Cherri
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/electrikmilk/args-parser"
	"github.com/google/uuid"
)

var idx int
var lines []string
var chars []rune
var char rune

var lineIdx int
var lineCharIdx int

var groupingUUIDs map[int]string
var groupingTypes map[int]tokenType
var groupingIdx int

var preParsing bool

// resetParse will take the current lines and merge them together to create new contents,
// then reset the chars and lines, then reset the parser cursor position.
// This is usually done when something modifies the contents of the file like custom actions or includes.
func resetParse() {
	contents = strings.Join(lines, "\n")
	chars = []rune(contents)
	lines = strings.Split(contents, "\n")
	firstChar()
}

func initParse() {
	if args.Using("debug") {
		fmt.Printf("Parsing %s...\n", filename)
	}
	variables = make(map[string]varValue)
	questions = make(map[string]*question)
	groupingUUIDs = make(map[int]string)
	groupingTypes = make(map[int]tokenType)
	chars = []rune(contents)
	lines = strings.Split(contents, "\n")
	idx = -1
	advance()

	preParse()

	for char != -1 {
		parse()
	}
	if args.Using("debug") {
		printParsingDebug()
	}

	contents = ""
	char = -1
	idx = -1
	lineIdx = 0
	lineCharIdx = -1
	actionIndex = 0
	chars = []rune{}
	lines = []string{}
	groupingUUIDs = map[int]string{}
	groupingTypes = map[int]tokenType{}
	groupingIdx = 0
	includes = []include{}

	if args.Using("debug") {
		fmt.Println(ansi("Done.", green) + "\n")
	}
}

func preParse() {
	preParsing = true

	defineRawAction()
	waitFor(
		defineToggleSetActions,
		handleIncludes,
	)

	handleCopyPastes()
	handleCustomActions()

	writeProcessed()

	preParsing = false
}

func writeProcessed() {
	if !args.Using("debug") {
		return
	}
	if len(included) == 0 && len(customActions) == 0 && len(pasteables) == 0 {
		return
	}

	var path = fmt.Sprintf("%s%s_processed.cherri", relativePath, workflowName)
	fmt.Printf("Writing processed version to %s...", path)

	var writeErr = os.WriteFile(path, []byte(contents), 0600)
	handle(writeErr)

	fmt.Println(ansi("Done.", green))
}

func printParsingDebug() {
	fmt.Println(ansi("### PARSING ###", bold) + "\n")

	lineReport("Current Line")

	fmt.Println(ansi("## TOKENS ##", bold))
	printTokens(tokens)
	fmt.Print("\n")

	fmt.Println(ansi("## DEFINITIONS ##", bold))
	fmt.Println("Name: " + workflowName)
	fmt.Println("Color:", iconColor)
	fmt.Printf("Glyph: %d\n", iconGlyph)
	fmt.Printf("Inputs: %v\n", inputs)
	fmt.Printf("Outputs: %v\n", outputs)
	fmt.Printf("Workflows: %v\n", definedWorkflowTypes)
	fmt.Printf("Quick Actions: %v\n", definedQuickActions)
	fmt.Printf("No Input: %v\n", noInput)
	if def, found := definitions["mac"]; found {
		fmt.Printf("macOS Only: %v\n", def)
	} else {
		fmt.Println("macOS Only: false")
	}
	fmt.Printf("Client Version: %s\n", clientVersion)
	fmt.Printf("iOS Version: %.1f\n", iosVersion)
	fmt.Print("\n")

	fmt.Println(ansi("## ENUMERATIONS ##", bold))
	fmt.Println(enumerations)
	fmt.Print("\n")

	fmt.Println(ansi("## VARIABLES ##", bold))
	printVariables()
	fmt.Print("\n")

	fmt.Println(ansi("## MENUS ##", bold))
	fmt.Println(menus)
	fmt.Print("\n")

	fmt.Println(ansi("## IMPORT QUESTIONS ##", bold))
	fmt.Println(questions)
	fmt.Print("\n")
}

func parse() {
	switch {
	case char == ' ' || char == '\t' || char == '\n':
		advance()
	case tokenAhead(Question):
		collectQuestion()
	case tokenAhead(Definition):
		collectDefinition()
	case tokenAhead(Import):
		collectImport()
	case isChar('@'):
		collectVariable(false)
	case tokenAhead(Constant):
		advance()
		collectVariable(true)
	case isChar('/'):
		collectComment()
	case tokenAhead(Repeat):
		collectRepeat()
	case tokenAhead(RepeatWithEach):
		collectRepeatEach()
	case tokenAhead(Menu):
		collectMenu()
	case tokenAhead(Item):
		collectMenuItem()
	case tokenAhead(If):
		collectConditionals()
	case tokenAhead(RightBrace):
		collectEndStatement()
	case tokenAhead(Enumeration):
		collectEnumeration()
	case strings.Contains(lookAheadUntil(' '), "("):
		collectActionCall()
	default:
		parserError(fmt.Sprintf("Illegal character '%s'", string(char)))
	}
}

var lastToken token

// reachable checks if the last token was a "stopper" and throws a warning if so,
// should only be run when we are about to parse a new statement.
func reachable() {
	if len(tokens) == 0 {
		return
	}
	lastToken = tokens[len(tokens)-1]
	if lastToken.valueType != Action {
		return
	}
	var lastActionIdentifier = lastToken.value.(action).ident
	var stoppers = []string{"stop", "output", "mustOutput", "outputOrClipboard"}
	if slices.Contains(stoppers, lastActionIdentifier) {
		parserWarning(fmt.Sprintf("Dead actions: Statement appears to be unreachable or does not loop as %s() was called outside of conditional.", lastActionIdentifier))
	}
}

func collectUntilIgnoreStrings(ch rune) string {
	var collected strings.Builder
	var insideString = false
	for char != -1 {
		if char == ch && !insideString {
			break
		}
		if char == '"' {
			insideString = prev(1) == '\\'
		}
		collected.WriteRune(char)
		advance()
	}

	return strings.Trim(collected.String(), " ")
}

// collectUntil advances ahead until the current character is `ch`,
// This should be used in cases where we are unsure how many characters
// will occur before we reach this character.
func collectUntil(ch rune) string {
	var collected strings.Builder
	for char != ch && char != -1 {
		collected.WriteRune(char)
		advance()
	}

	return strings.Trim(collected.String(), " ")
}

// lookAheadUntil does a pseudo string collection stopping when we reach `until` and returning the collected string.
func lookAheadUntil(until rune) string {
	var ahead strings.Builder
	var nextIdx = idx
	var nextChar rune
	for nextChar != until {
		if len(chars) <= nextIdx {
			break
		}

		nextChar = chars[nextIdx]
		ahead.WriteRune(chars[nextIdx])
		nextIdx++
	}

	return strings.Trim(ahead.String(), " \t\n")
}

func collectVariableValue(constant bool, valueType *tokenType, value *any) {
	collectValue(valueType, value, '\n')

	if constant && (*valueType == Arr || *valueType == Variable) {
		parserError(fmt.Sprintf("Type %v values cannot be constants.", *valueType))
	}
	if *valueType == Question {
		parserError(fmt.Sprintf("Illegal reference to import question '%s'. Shortcuts does not support import questions as variable values.", *value))
	}
}

func collectValue(valueType *tokenType, value *any, until rune) {
	var ahead = lookAheadUntil(until)
	if ahead == "" {
		parserError("Value expected")
	}
	switch {
	case intChar():
		*valueType = Integer
		if strings.Contains(ahead, ".") {
			*valueType = Float
		}
		collectIntegerValue(valueType, value, &until)
	case char == '"':
		collectStringValue(valueType, value)
	case char == '\'':
		advance()
		*valueType = RawString
		*value = collectRawString()
	case char == '[':
		advance()
		*valueType = Arr
		*value = collectArray()
	case char == '{':
		advance()
		*valueType = Dict
		*value = collectDictionary()
	case tokenAhead(True):
		*valueType = Bool
		*value = true
	case tokenAhead(False):
		*valueType = Bool
		*value = false
	case tokenAhead(Nil):
		*valueType = Nil
		advanceUntil(until)
	case strings.Contains(ahead, "("):
		collectActionValue(valueType, value)
	case containsTokens(&ahead, Plus, Minus, Multiply, Divide, Modulus):
		*valueType = Expression
		*value = collectUntil(until)
	default:
		collectReference(valueType, value, &until)
	}
}

func collectStringValue(valueType *tokenType, value *any) {
	advance()
	*valueType = String
	*value = collectString()

	var stringValue = fmt.Sprintf("%s", *value)
	if strings.ContainsAny(stringValue, "{}") {
		checkInlineVars(&stringValue)
	}
}

func collectActionValue(valueType *tokenType, value *any) {
	*valueType = Action
	var identifier = collectIdentifier()

	if _, found := customActions[identifier]; found {
		*value = makeCustomActionRef(&identifier)
		return
	}

	*value = collectAction(&identifier)
}

var collectVarRegex = regexp.MustCompile(`\{(.*?)(?:\['(.*?)'])?(?:\.(.*?))?}`)

func checkInlineVars(value *string) {
	var matches = collectVarRegex.FindAllStringSubmatch(*value, -1)
	if matches == nil {
		return
	}

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		var identifier = match[1]

		if startsWith(Ask, identifier) {
			identifier = Ask
		}

		if !validReference(identifier) {
			parserError(fmt.Sprintf("Undefined inline variable reference '%s'", identifier))
		}
	}
}

func collectReference(valueType *tokenType, value *any, until *rune) {
	var reference = collectIdentifier()
	var getAs string
	var coerce string
	var prompt string

	if q, found := questions[reference]; found {
		if q.used {
			parserError(fmt.Sprintf("Duplicate usage of import question '%s'. Import questions can only be used in one action argument.", reference))
		}

		*valueType = Question
		*value = reference
		q.used = true

		advance()
		return
	}

	if !validReference(reference) {
		parserError(fmt.Sprintf("Undefined reference '%s'", reference))
	}

	if reference == "Ask" && char == ':' {
		advance()
		skipWhitespace()
		if char != '"' {
			parserError(fmt.Sprintf("Expected prompt string, got: %c", char))
		}
		advance()
		prompt = collectString()
	}

	if char == '[' {
		advance()
		if char != '\'' {
			parserError("Expected raw string for key.")
		}
		advance()
		var key = collectRawString()
		advanceUntil(']')
		advance()
		getAs = key
	}
	if char == '.' {
		advance()
		coerce = collectUntil(*until)
	}

	isInputVariable(reference)

	*valueType = Variable
	*value = varValue{
		variableType: "Variable",
		valueType:    Variable,
		value:        reference,
		coerce:       coerce,
		getAs:        getAs,
		prompt:       prompt,
	}
}

func collectEnumeration() {
	advance()
	if enumerations == nil {
		enumerations = make(map[string][]string)
	}

	var identifier = collectIdentifier()

	advanceUntilExpect('{', 3)
	advance()
	skipWhitespace()

	if char != '\'' {
		parserError(fmt.Sprintf("Expected enumeration raw string ('), got: %c", char))
	}
	advance()

	var enumeration []string
	for char != '}' && char != -1 {
		var permutation = collectRawString()
		enumeration = append(enumeration, permutation)
		if char == ',' {
			advance()
		}
		skipWhitespace()
		if char == '\'' {
			advance()
		}
	}

	enumerations[identifier] = enumeration
	advance()
}

func collectArguments() (arguments []actionArgument) {
	var params = currentAction.parameters
	var paramsSize = len(params)
	var argIndex = 0
	var param parameterDefinition
	for {
		if char == ')' || char == -1 {
			break
		}
		if argIndex < paramsSize {
			param = params[argIndex]
		}
		arguments = append(arguments, collectArgument(&argIndex, &param, &paramsSize))
		argIndex++
	}
	return
}

func collectArgument(argIndex *int, param *parameterDefinition, paramsSize *int) (argument actionArgument) {
	if *argIndex == *paramsSize && !param.infinite {
		parserError(
			fmt.Sprintf("Too many arguments\n\n%s",
				generateActionDefinition(parameterDefinition{}, false, false),
			),
		)
	}
	if char == ',' {
		advance()
	}
	skipWhitespace()

	var valueType tokenType
	var value any
	collectValue(&valueType, &value, endOfNextArgument())

	skipWhitespace()

	argument = actionArgument{
		valueType: valueType,
		value:     value,
	}
	if !param.infinite && (valueType != Nil && value != nil) {
		checkArg(param, &argument)
	}
	return
}

func endOfNextArgument() (until rune) {
	until = ')'
	if strings.Contains(lookAheadUntil(')'), ",") {
		until = ','
	}
	return
}

func collectComment() {
	var collect = args.Using("comments")
	var comment strings.Builder
	if isChar('/') {
		if collect {
			comment.WriteString(collectUntil('\n'))
		} else {
			advanceUntil('\n')
		}
	} else if char == '*' {
		collectMultilineComment(&comment, &collect)
	}
	if collect && !preParsing {
		var commentStr = strings.Trim(comment.String(), " \n")
		tokens = append(tokens, token{
			typeof:    Comment,
			ident:     "",
			valueType: String,
			value:     commentStr,
		})
	}
}

func collectMultilineComment(comment *strings.Builder, collect *bool) {
	advanceTimes(2)
	for char != 1 {
		if char == '*' && next(1) == '/' {
			break
		}
		if *collect {
			comment.WriteRune(char)
		}
		advance()
	}
	advanceTimes(3)
}

func collectVariable(constant bool) {
	reachable()

	var identifier = collectIdentifier()
	availableIdentifier(&identifier)

	var valueType tokenType
	var value any
	var varType = Variable
	switch {
	case strings.Contains(lookAheadUntil('\n'), "="):
		advance()
		switch {
		case tokenAhead(AddTo):
			varType = AddTo
		case tokenAhead(SubFrom):
			varType = SubFrom
		case tokenAhead(MultiplyBy):
			varType = MultiplyBy
		case tokenAhead(DivideBy):
			varType = DivideBy
		case tokenAhead(Set):
		}
		if varType != Variable && constant {
			parserError("Constants cannot be added to.")
		}
		advance()

		collectVariableValue(constant, &valueType, &value)

		if valueType == Variable && value.(varValue).value == "Ask" {
			parserError("Ask global cannot be used as a variable value.")
		}
	case tokenAhead(Colon):
		if constant {
			parserError("Constants cannot be initialized without a value")
		}
		skipWhitespace()
		collectType(&valueType, &value, '\n')
	case constant:
		parserError("Constants must be initialized with a value.")
	}

	tokens = append(tokens, token{
		typeof:    varType,
		ident:     identifier,
		valueType: valueType,
		value:     value,
	})

	if varType != Variable {
		return
	}
	if _, found := variables[identifier]; !found {
		variables[identifier] = varValue{
			variableType: "Variable",
			valueType:    valueType,
			value:        value,
			constant:     constant,
		}
	}
}

func collectType(valueType *tokenType, value *any, until rune) {
	switch {
	case tokenAhead(String):
		*valueType = String
		*value = ""
	case tokenAhead(Integer):
		*valueType = Integer
		*value = "0"
	case tokenAhead(Float):
		*valueType = Float
		*value = "0.0"
	case tokenAhead(Bool):
		*valueType = Bool
		*value = false
	case tokenAhead(Arr):
		*valueType = Arr
	case tokenAhead(Dict):
		*valueType = Dict
		*value = make(map[string]interface{})
	case tokenAhead(Variable):
		*valueType = Variable
	default:
		parserError(fmt.Sprintf("Unknown type '%s'\n\nAvailable types: \n- text\n- number\n- bool\n- array\n- dictionary\n- var", collectUntil(until)))
	}
}

func collectIdentifier() string {
	var identifier strings.Builder
	for char != -1 {
		if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_' {
			identifier.WriteRune(char)
			advance()
			continue
		}

		break
	}

	return identifier.String()
}

func collectDefinition() {
	if len(definitions) == 0 {
		definitions = make(map[string]any)
	}
	advance()

	switch {
	case tokenAhead(Name):
		advance()
		workflowName = collectUntil('\n')
		if strings.Trim(workflowName, " \n\t") == "" {
			parserError("Expected name")
		}
		if !args.Using("output") {
			outputPath = relativePath + workflowName + ".shortcut"
		}
	case tokenAhead(Color):
		advance()
		collectColorDefinition()
	case tokenAhead(Glyph):
		advance()
		collectGlyphDefinition()
	case tokenAhead(Inputs):
		advance()
		inputs = collectTypeValues("content item", contentItems)
	case tokenAhead(Outputs):
		advance()
		outputs = collectTypeValues("content item", contentItems)
	case tokenAhead(From):
		advance()
		definedWorkflowTypes = collectTypeValues("workflow", workflowTypes)
	case tokenAhead(NoInput):
		advance()
		collectNoInputDefinition()
	case tokenAhead(Mac):
		advance()
		collectBooleanDefinition("mac")
	case tokenAhead(QuickActions):
		advance()
		definedQuickActions = collectTypeValues("quick action", quickActions)
	case tokenAhead(Version):
		var collectVersion = collectUntil('\n')
		if version, found := versions[collectVersion]; found {
			clientVersion = version
			iosVersion, _ = strconv.ParseFloat(collectVersion, 32)
		} else {
			var list = makeKeyList("Available versions:", versions, collectVersion)
			parserError(fmt.Sprintf("Invalid minimum version '%s'\n\n%s", collectVersion, list))
		}
	case tokenAhead(Action):
		advance()
		collectDefinedAction()
	}
}

func collectDefinedAction() {
	if char != '\'' {
		parserError("Expected workflow action identifier raw string (').")
	}
	advance()

	var workflowIdentifier = collectRawString()
	var shortIdentifier string
	var overrideIdentifier string
	if len(strings.Split(workflowIdentifier, ".")) < 4 {
		shortIdentifier = workflowIdentifier
	} else {
		overrideIdentifier = workflowIdentifier
	}
	advance()

	var identifier, arguments, outputType = collectActionDefinition('\n')
	if !usingAction(contents, identifier) {
		return
	}

	skipWhitespace()

	var addParams paramsFunc
	if char == '{' {
		advance()
		var dict = collectDictionary()
		addParams = func(args []actionArgument) map[string]any {
			handleRawParams(dict.(map[string]interface{}))
			return dict.(map[string]any)
		}
	}

	actions[identifier] = &actionDefinition{
		identifier:         shortIdentifier,
		overrideIdentifier: overrideIdentifier,
		parameters:         arguments,
		outputType:         outputType,
		addParams:          addParams,
	}

	if args.Using("debug") {
		setCurrentAction(identifier, actions[identifier])
		fmt.Println("\ndefined:", currentAction.appIdentifier, generateActionDefinition(parameterDefinition{}, true, true))
		fmt.Print("\n")
	}
}

func collectColorDefinition() {
	var collectColor = collectUntil('\n')
	collectColor = strings.ToLower(collectColor)
	if color, found := colors[collectColor]; found {
		iconColor = color
	} else {
		var list = "Available icon colors:\n"
		for c := range colors {
			list += fmt.Sprintf("- %s\n", c)
		}
		parserError(fmt.Sprintf("Invalid icon color '%s'\n\n%s", collectColor, list))
	}
}

func collectTypeValues(typeName string, valueTypes map[string]string) (types []string) {
	var collectedTypes = collectUntil('\n')
	if collectedTypes == "" {
		parserError("Expected type")
	}
	var definedTypes = strings.Split(collectedTypes, ",")
	for _, definedType := range definedTypes {
		definedType = strings.Trim(definedType, " ")
		if _, found := valueTypes[definedType]; !found {
			var list = makeKeyList(fmt.Sprintf("Available %s types:", typeName), workflowTypes, definedType)
			parserError(fmt.Sprintf("Invalid %s type '%s'\n\n%s", typeName, definedType, list))
			return
		}
		types = append(types, fmt.Sprintf("%v", valueTypes[definedType]))
	}
	return
}

func collectGlyphDefinition() {
	var collectGlyph = collectUntil('\n')
	if glyph, found := glyphs[collectGlyph]; found {
		glyphInt, hexErr := strconv.ParseInt(fmt.Sprintf("%d", glyph), 10, 64)
		handle(hexErr)
		iconGlyph = glyphInt
	} else {
		if !args.Using("no-ansi") {
			fmt.Println(ansi("Build an icon for your Shortcut at https://glyphs.cherrilang.org/.", cyan))
			args.Args["glyph"] = collectGlyph
			glyphsSearch()
		}
		parserError(fmt.Sprintf("Invalid icon glyph '%s'", collectGlyph))
	}
}

func collectNoInputDefinition() {
	switch {
	case tokenAhead(StopWith):
		advanceTimes(2)
		var stopWithError = collectString()
		noInput = WFWorkflowNoInputBehavior{
			Name: "WFWorkflowNoInputBehaviorShowError",
			Parameters: map[string]string{
				"Error": stopWithError,
			},
		}
	case tokenAhead(AskFor):
		advance()
		var contentItemType = collectUntil('\n')
		if itemClass, found := contentItems[contentItemType]; found {
			noInput = WFWorkflowNoInputBehavior{
				Name: "WFWorkflowNoInputBehaviorAskForInput",
				Parameters: map[string]string{
					"ItemClass": itemClass,
				},
			}
		} else {
			var list = makeKeyList("Available content item types:", contentItems, itemClass)
			parserError(fmt.Sprintf("Invalid content item type '%s'\n\n%s", itemClass, list))
		}
	case tokenAhead(GetClipboard):
		noInput = WFWorkflowNoInputBehavior{
			Name: "WFWorkflowNoInputBehaviorGetClipboard",
		}
	}
}

func collectBooleanDefinition(key string) {
	switch {
	case tokenAhead(True):
		definitions[key] = true
	case tokenAhead(False):
		definitions[key] = false
	default:
		parserError(fmt.Sprintf("Invalid value for boolean definition '%s'", key))
	}
}

// libraries is a map of the 3rd party libraries defined in the compiler.
// The key determines the identifier of the identifier name that must be used in the syntax, it's value defines its behavior, etc. using an libraryDefinition.
var libraries map[string]libraryDefinition

func collectImport() {
	makeLibraries()
	advanceTimes(2)
	var collectedLibrary = collectString()
	if lib, found := libraries[collectedLibrary]; found {
		lib.make(lib.identifier)
	} else {
		parserError(fmt.Sprintf("Import library '%s' does not exist!", collectedLibrary))
	}
}

var questions map[string]*question

type question struct {
	parameter    string
	actionIndex  int
	text         string
	defaultValue string
	used         bool
}

func collectQuestion() {
	advance()
	var identifier = collectIdentifier()
	if _, found := questions[identifier]; found {
		parserError(fmt.Sprintf("Duplicate declaration of import question '%s'.", identifier))
	}
	if validReference(identifier) {
		parserError(fmt.Sprintf("Import question conflicts with defined variable or global '%s'.", identifier))
	}
	advance()

	if char != '"' {
		parserError("Expected question prompt string.")
	}
	advance()

	var text = collectString()
	advance()

	if char != '"' {
		parserError("Expected question default string value.")
	}
	advance()

	var defaultValue = collectString()
	questions[identifier] = &question{
		text:         text,
		defaultValue: defaultValue,
	}
}

var repeatItemIndex = 1

func collectRepeat() {
	reachable()
	var groupingUUID = groupStatement(Repeat)

	var index string
	if repeatItemIndex > 1 {
		index = fmt.Sprintf(" %d", repeatItemIndex)
	}
	var repeatIndexIdentifier = collectIdentifier()

	advance()
	if !tokenAhead(RepeatWithEach) {
		parserError(fmt.Sprintf("Expected `for`, got: %c", char))
	}

	if char == '{' {
		parserError("Expected number of times to repeat")
	}

	var timesType tokenType
	var timesValue any
	collectValue(&timesType, &timesValue, '{')
	advanceTimes(2)

	tokens = append(tokens,
		token{
			typeof:    Repeat,
			ident:     groupingUUID,
			valueType: timesType,
			value:     timesValue,
		}, token{
			typeof:    Variable,
			ident:     repeatIndexIdentifier,
			valueType: Variable,
			value: varValue{
				valueType: Variable,
				value:     fmt.Sprintf("Repeat Index%s", index),
			},
		},
	)

	variables[repeatIndexIdentifier] = varValue{
		variableType: "Variable",
		valueType:    Integer,
		value:        repeatIndexIdentifier,
		repeatItem:   true,
	}
}

func collectRepeatEach() {
	reachable()
	var groupingUUID = groupStatement(RepeatWithEach)

	var index string
	if repeatItemIndex > 1 {
		index = fmt.Sprintf(" %d", repeatItemIndex)
	}
	var repeatItemIdentifier = collectIdentifier()

	advance()
	if !tokenAhead(In) {
		parserError(fmt.Sprintf("Expected `in`, got: %c", char))
	}
	advance()

	if char == '{' {
		parserError("Expected value")
	}

	var iterableType tokenType
	var iterableValue any
	collectValue(&iterableType, &iterableValue, '{')
	advanceTimes(2)

	tokens = append(tokens,
		token{
			typeof:    RepeatWithEach,
			ident:     groupingUUID,
			valueType: iterableType,
			value:     iterableValue,
		}, token{
			typeof:    Variable,
			ident:     repeatItemIdentifier,
			valueType: Variable,
			value: varValue{
				valueType: Variable,
				value:     fmt.Sprintf("Repeat Item%s", index),
			},
		},
	)

	variables[repeatItemIdentifier] = varValue{
		variableType: "Variable",
		valueType:    String,
		value:        repeatItemIdentifier,
		repeatItem:   true,
	}

	repeatItemIndex++
}

func collectConditionals() {
	reachable()
	advance()

	var groupingUUID = groupStatement(Conditional)

	var conditions WFConditions
	conditions.WFActionParameterFilterPrefix = -1
	for char != 1 && char != '{' {
		var conditional = collectConditional()

		collectFilterPrefix(&conditions)
		conditions.conditions = append(conditions.conditions, conditional)
	}

	advance()

	tokens = append(tokens, token{
		typeof:    Conditional,
		ident:     groupingUUID,
		valueType: If,
		value:     conditions,
	})
}

func collectFilterPrefix(wfConditions *WFConditions) {
	var ahead = lookAheadUntil(' ')
	if !containsTokens(&ahead, And, Or) {
		return
	}

	var collectedFilterPrefix tokenType
	switch {
	case tokenAhead(And):
		collectedFilterPrefix = And
	case tokenAhead(Or):
		collectedFilterPrefix = Or
	}

	if wfConditions.WFActionParameterFilterPrefix != -1 &&
		wfConditions.WFActionParameterFilterPrefix != conditionFilterPrefixes[collectedFilterPrefix] {
		var initLogicOperator = convertFilterPrefix(wfConditions.WFActionParameterFilterPrefix)
		parserError(fmt.Sprintf("Due to limitations in Shortcuts, only the initially set logical operator '%s' is allowed for all other comparisons.", initLogicOperator))
	}

	skipWhitespace()

	if wfConditions.WFActionParameterFilterPrefix == -1 {
		wfConditions.WFActionParameterFilterPrefix = conditionFilterPrefixes[collectedFilterPrefix]
	}
}

func convertFilterPrefix(filterPrefix int) string {
	switch filterPrefix {
	case conditionFilterPrefixes[And]:
		return "&&"
	case conditionFilterPrefixes[Or]:
		return "||"
	}
	return "||"
}

func collectConditional() (conditional condition) {
	if isChar('!') {
		conditional.condition = conditions[Empty]
	} else {
		conditional.condition = conditions[Any]
	}

	var variableOneType tokenType
	var variableOneValue any
	collectValue(&variableOneType, &variableOneValue, ' ')
	conditional.arguments = append(conditional.arguments, actionArgument{
		valueType: variableOneType,
		value:     variableOneValue,
	})

	skipWhitespace()

	var ahead = lookAheadUntil(' ')
	if containsTokens(&ahead, And, Or) {
		return
	}

	if char != '{' {
		var collectConditional = collectUntil(' ')
		var collectConditionalToken = tokenType(collectConditional)
		if condition, found := conditions[collectConditionalToken]; found {
			conditional.condition = condition
			checkConditionalTypes(&collectConditionalToken, variableOneType, variableOneValue)
		} else {
			parserError(fmt.Sprintf("Invalid conditional '%s'", collectConditional))
		}

		skipWhitespace()

		var variableTwoType tokenType
		var variableTwoValue any
		collectValue(&variableTwoType, &variableTwoValue, ' ')
		conditional.arguments = append(conditional.arguments, actionArgument{
			valueType: variableTwoType,
			value:     variableTwoValue,
		})

		skipWhitespace()

		if char != '{' && char != '&' && char != '|' {
			var variableThreeType tokenType
			var variableThreeValue any

			collectValue(&variableThreeType, &variableThreeValue, '{')
			conditional.arguments = append(conditional.arguments, actionArgument{
				valueType: variableThreeType,
				value:     variableThreeValue,
			})
			skipWhitespace()
		}
	}

	return
}

func checkConditionalTypes(conditional *tokenType, variableType tokenType, value any) {
	if variableType == Variable {
		var variable = value.(varValue)
		variableType = variable.valueType
		var variableValue, found = getVariableValue(variable.value.(string))
		if found && variableValue.valueType != Variable {
			variableType = variableValue.valueType
		}
		if variable.coerce != "" {
			variableType = tokenType(variable.coerce)
		}
	}

	if len(allowedConditionalTypes[*conditional]) != 0 && !slices.Contains(allowedConditionalTypes[*conditional], variableType) {
		parserError(
			fmt.Sprintf("Invalid type '%s' for conditional '%s'\nAllowed types: %s",
				variableType,
				*conditional,
				allowedConditionalTypes[*conditional],
			),
		)
	}
}

func collectMenu() {
	if len(menus) == 0 {
		menus = make(map[string][]varValue)
	}

	reachable()
	advance()
	var groupingUUID = groupStatement(Menu)

	var promptType tokenType
	var promptValue any
	collectValue(&promptType, &promptValue, '{')
	advanceUntil('{')
	advance()

	menus[groupingUUID] = []varValue{}
	tokens = append(tokens, token{
		typeof:    Menu,
		ident:     groupingUUID,
		valueType: promptType,
		value:     promptValue,
	})
}

func collectMenuItem() {
	advance()
	if _, ok := groupingUUIDs[groupingIdx]; !ok {
		parserError("Item has no starting menu statement.")
	}
	var groupingUUID = groupingUUIDs[groupingIdx]

	var itemType tokenType
	var itemValue any
	collectValue(&itemType, &itemValue, ':')
	advanceUntil(':')
	advance()

	if len(menus[groupingUUID]) > 0 {
		addNothing()
	}

	menus[groupingUUID] = append(menus[groupingUUID], varValue{
		valueType: itemType,
		value:     itemValue,
	})

	tokens = append(tokens,
		token{
			typeof:    Item,
			ident:     groupingUUID,
			valueType: itemType,
			value:     itemValue,
		},
	)
}

func collectEndStatement() {
	advance()
	if tokenAhead(Else) {
		addNothing()

		advance()
		if _, ok := groupingUUIDs[groupingIdx]; !ok {
			parserError("Else has no starting if statement.")
		}
		tokens = append(tokens, token{
			typeof:    Conditional,
			ident:     groupingUUIDs[groupingIdx],
			valueType: Else,
			value:     nil,
		})
		tokenAhead(LeftBrace)
	} else {
		if _, ok := groupingUUIDs[groupingIdx]; !ok {
			parserError("Ending has no starting statement.")
		}
		var groupType = groupingTypes[groupingIdx]
		if groupType == Repeat || groupType == RepeatWithEach {
			repeatItemIndex--
			reachable()
		}

		addNothing()

		tokens = append(tokens, token{
			typeof:    groupType,
			ident:     groupingUUIDs[groupingIdx],
			valueType: EndClosure,
			value:     nil,
		})
		groupingIdx--
	}
}

// groupStatement creates a grouping UUID for a statement and adds to the statement groupings.
func groupStatement(groupType tokenType) (groupingUUID string) {
	groupingUUID = uuid.New().String()
	groupingIdx++
	groupingUUIDs[groupingIdx] = groupingUUID
	groupingTypes[groupingIdx] = groupType

	return
}

// addNothing helps reduce memory usage by not passing anything to the next action.
func addNothing() {
	lastToken = tokens[len(tokens)-1]
	if lastToken.typeof == Action && lastToken.ident == "nothing" {
		return
	}

	tokens = append(tokens, token{
		typeof:    Action,
		ident:     "nothing",
		valueType: Action,
		value: action{
			ident: "nothing",
		},
	})
}

func intChar() bool {
	return (char >= '0' && char <= '9') || char == '-' || char == '.'
}

func collectInteger() string {
	var collection strings.Builder
	for intChar() {
		collection.WriteRune(char)
		advance()
	}

	return collection.String()
}

func collectIntegerValue(valueType *tokenType, value *any, until *rune) {
	var ahead = lookAheadUntil(*until)
	if !containsTokens(&ahead, Plus, Minus, Multiply, Divide, Modulus) {
		var integerString = collectInteger()

		if *valueType == Integer {
			var integer, convErr = strconv.Atoi(integerString)
			handle(convErr)

			*value = integer
		} else {
			var float, floatErr = strconv.ParseFloat(integerString, 64)
			handle(floatErr)

			*value = float
		}
		return
	}
	*valueType = Expression
	*value = collectUntil(*until)
}

func collectString() string {
	var collection strings.Builder
	for char != -1 {
		if char == '\\' {
			switch next(1) {
			case '"':
				collection.WriteRune('"')
			case 'n':
				collection.WriteRune('\n')
			case 't':
				collection.WriteRune('\t')
			case 'r':
				collection.WriteRune('\r')
			case '\\':
				collection.WriteRune('\\')
			default:
				collection.WriteRune(next(1))
			}
			advanceTimes(2)
			continue
		} else if char == '"' {
			break
		}

		collection.WriteRune(char)
		advance()
	}
	advance()

	return collection.String()
}

func collectRawString() string {
	var collection strings.Builder
	for char != -1 {
		if char == '\\' && next(1) == '\'' {
			collection.WriteRune('\'')
			advanceTimes(2)
			continue
		}
		if char == '\'' {
			break
		}

		collection.WriteRune(char)
		advance()
	}
	advance()

	return collection.String()
}

func collectArray() (array interface{}) {
	var rawJSON = "{\"array\":[" + collectUntilIgnoreStrings(']') + "]}"
	if err := json.Unmarshal([]byte(rawJSON), &array); err != nil {
		if args.Using("debug") {
			fmt.Println(ansi("\n### COLLECTED ARRAY ###", bold))
			fmt.Println(rawJSON)
			fmt.Print("\n")
		}
		parserError(fmt.Sprintf("JSON error: %s", err))
	}
	array = array.(map[string]interface{})["array"]
	advance()
	return
}

func collectDictionary() (dictionary interface{}) {
	if char == '}' {
		advance()
		return
	}
	var rawJSON = "{" + collectObject() + "}"
	if args.Using("debug") {
		fmt.Println(ansi("\n\n### COLLECTED DICTIONARY ###", bold))
		fmt.Println(rawJSON)
		fmt.Print("\n")
	}
	if err := json.Unmarshal([]byte(rawJSON), &dictionary); err != nil {
		parserError(fmt.Sprintf("JSON error: %s", err))
	}
	return
}

func collectObject() string {
	var jsonStr strings.Builder
	var innerObjectDepth = 0
	var insideString = false
	for char != -1 {
		if char == '"' {
			if insideString {
				if prev(1) != '\\' {
					insideString = false
				}
			} else {
				insideString = true
			}
		}
		if !insideString {
			if char == '{' {
				innerObjectDepth++
			} else if char == '}' {
				if innerObjectDepth == 0 {
					break
				}
				innerObjectDepth--
			}
		}
		jsonStr.WriteRune(char)
		advance()
	}
	advance()

	return jsonStr.String()
}

func collectActionCall() {
	reachable()
	var identifier = collectIdentifier()
	if _, found := customActions[identifier]; found {
		tokens = append(tokens, token{
			typeof:    Action,
			ident:     "runSelf",
			valueType: Action,
			value:     makeCustomActionRef(&identifier),
		})
		actionIndex++
		return
	}

	var value = collectAction(&identifier)

	if identifier == "comment" && usingCustomActions && isFirstCommentAction {
		isFirstCommentAction = false
		tokens = append([]token{{
			typeof:    Action,
			ident:     identifier,
			valueType: Action,
			value:     value,
		}}, tokens...)
		return
	}

	tokens = append(tokens, token{
		typeof:    Action,
		ident:     identifier,
		valueType: Action,
		value:     value,
	})
	actionIndex++
}

func collectAction(identifier *string) (value action) {
	if _, found := actions[*identifier]; !found {
		parserError(fmt.Sprintf("Undefined action '%s()'", *identifier))
	}
	advance()
	setCurrentAction(*identifier, actions[*identifier])

	var arguments = collectArguments()
	currentArguments = arguments
	currentArgumentsSize = len(currentArguments)

	checkAction()

	value = action{
		ident: *identifier,
		args:  arguments,
	}

	advance()
	return
}

// advance advances the character cursor.
func advance() {
	if char == '\n' {
		lineCharIdx = 0
		lineIdx++
	} else {
		lineCharIdx++
	}

	idx++
	if len(chars) <= idx {
		char = -1
		return
	}
	char = chars[idx]
}

// advanceTimes advances the character cursor by `times`.
func advanceTimes(times int) {
	for i := 0; i < times; i++ {
		advance()
	}
}

// advanceUntil advances the character cursor until we reach `ch`.
func advanceUntil(ch rune) {
	for char != ch && char != -1 {
		advance()
	}
}

// advanceUntilExpect advances the character cursor until we reach `ch`.
// However, it expects to reach this character by no more than `maxAdvances` advances and throws a parser error if it doesn't.
func advanceUntilExpect(ch rune, maxAdvances int) {
	var advances int
	for char != ch && char != -1 {
		if advances > maxAdvances {
			parserError(fmt.Sprintf("Expected %c", ch))
		}
		advances++
		advance()
	}
}

// If char is tokenChar, advance.
func isChar(tokenChar rune) bool {
	if char != tokenChar {
		return false
	}
	advance()
	return true
}

func tokenAhead(token tokenType) bool {
	var tokenLen = len(token)
	for i, tokenChar := range token {
		if (i == 0 && unicode.ToLower(char) != tokenChar) || next(i) != tokenChar {
			return false
		}
	}

	advanceTimes(tokenLen)
	return true
}

func containsTokens(str *string, v ...tokenType) bool {
	for _, aheadToken := range v {
		if strings.Contains(*str, string(aheadToken)) {
			return true
		}
	}
	return false
}

func next(mov int) rune {
	return seek(&mov, false)
}

func prev(mov int) rune {
	return seek(&mov, true)
}

func seek(mov *int, reverse bool) rune {
	var nextChar = idx
	if reverse {
		nextChar -= *mov
	} else {
		nextChar += *mov
	}

	return getChar(nextChar)
}

func getChar(atIndex int) rune {
	if atIndex < 0 {
		return -1
	}
	if len(chars) > atIndex {
		return chars[atIndex]
	}
	return -1
}

func firstChar() {
	lineIdx = 0
	lineCharIdx = -1
	idx = -1
	advance()
}

func skipWhitespace() {
	for char == ' ' || char == '\t' || char == '\n' {
		advance()
	}
}

func printVariables() {
	for identifier, v := range variables {
		if v.constant {
			fmt.Print("const ")
		} else {
			fmt.Print("@")
		}
		fmt.Print(identifier)

		if v.getAs != "" {
			fmt.Printf("[%s]", v.getAs)
		}
		if v.coerce != "" {
			fmt.Printf(".%s", v.coerce)
		}
		if v.variableType != "Variable" {
			fmt.Printf(" (%s)", v.variableType)
		}
		if v.value != nil {
			fmt.Printf(" = %s", v.value)
		}
		if string(v.valueType) != "" {
			fmt.Printf(" (%s)", v.valueType)
		}
		if v.repeatItem {
			fmt.Print(" (repeat item var)")
		}
		fmt.Print("\n")
	}
}

func printTokens(tokens []token) {
	var size = len(tokens)
	var pad = len(fmt.Sprintf("%d", size))
	for i, token := range tokens {
		var idx = i + 1
		var spaces = pad - len(fmt.Sprintf("%d", idx))
		fmt.Printf("%s%d | %s\n", strings.Repeat(" ", spaces), idx, token)
	}
}

func parserWarning(message string) {
	var errorFilename, errorLine, errorCol = delinquentFile()
	var warning = fmt.Sprintf("%s %s", ansi("\nWarning:", yellow, bold), message)

	if !args.Using("no-ansi") {
		warning += fmt.Sprintf(" %s:%d:%d", errorFilename, errorLine, errorCol)
	} else {
		warning += fmt.Sprintf(" (%d:%d)", errorLine, errorCol)
	}

	fmt.Println(warning + "\n")
}

func makeKeyList(title string, list map[string]string, value string) string {
	var formattedList strings.Builder
	formattedList.WriteString("\033[0m")
	formattedList.WriteString(fmt.Sprintf("%s\n", title))
	for key := range list {
		var matchedKey = key
		var matched, result = matchString(&key, &value)
		if matched {
			matchedKey = result
		}
		formattedList.WriteString(fmt.Sprintf("- %s\n", matchedKey))
	}
	formattedList.WriteString("\033[0m")

	return formattedList.String()
}

func parserError(message string) {
	lines = strings.Split(contents, "\n")
	var errorFilename, errorLine, errorCol = delinquentFile()

	if args.Using("no-ansi") {
		fmt.Printf("Error: %s (%d:%d)\n", message, errorLine, errorCol)
		os.Exit(1)
	}

	excerptError(message, errorFilename, errorLine, errorCol)

	if args.Using("debug") {
		panicDebug(nil)
	} else {
		os.Exit(1)
	}
}

func excerptError(message string, errorFilename string, errorLine int, errorCol int) {
	fmt.Print("\033[31m")
	fmt.Println("\n" + ansi(message, bold))
	fmt.Printf("\n\033[2m----- \033[0m%s:%d:%d\n", errorFilename, errorLine, errorCol)
	if len(lines) > (lineIdx-1) && lineIdx != 0 {
		fmt.Printf("\033[2m%d | %s\033[0m\n", errorLine-1, lines[lineIdx-1])
	}
	if len(lines) > lineIdx {
		fmt.Printf("\033[31m\033[1m%d | ", errorLine)
		for c, chr := range strings.Split(lines[lineIdx], "") {
			if c == idx {
				fmt.Print(ansi(chr, underline))
			} else {
				fmt.Print(chr)
			}
		}
		fmt.Print("\033[0m\n")
	}
	var spaces string
	for i := 0; i < (lineCharIdx + 4); i++ {
		spaces += " "
	}
	fmt.Println("\033[31m" + spaces + "^\033[0m")
	if len(lines) > (lineIdx + 1) {
		fmt.Printf("\033[2m%d | %s\n-----\033[0m\n\n", errorLine+1, lines[lineIdx+1])
	}
}
