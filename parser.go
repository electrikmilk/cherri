/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/electrikmilk/args-parser"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
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

func initParse() {
	if strings.Contains(contents, "action") {
		standardActions()
		parseCustomActions()
	}
	if args.Using("debug") {
		fmt.Printf("Parsing %s...\n", filename)
	}
	variables = make(map[string]variableValue)
	questions = make(map[string]*question)
	menus = make(map[string][]variableValue)
	groupingUUIDs = make(map[int]string)
	groupingTypes = make(map[int]tokenType)
	makeGlobals()
	chars = []rune(contents)
	lines = strings.Split(contents, "\n")
	idx = -1
	advance()

	parse()
	if args.Using("debug") {
		printParsingDebug()
	}

	contents = ""
	char = -1
	idx = -1
	lineIdx = 0
	lineCharIdx = 0
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

func printParsingDebug() {
	fmt.Println(ansi("### PARSING ###", bold) + "\n")

	if idx != 0 {
		fmt.Println("Previous Character:")
		printChar(prev(1))
	}

	fmt.Println("\nCurrent Character:")
	printChar(char)

	if len(contents) > idx+1 {
		fmt.Println("\nNext Character:")
		printChar(next(1))
	}

	if len(lines) > lineIdx {
		fmt.Println("\nCurrent Line: \n" + lines[lineIdx] + "\n")
	}

	fmt.Println(ansi("## TOKENS ##", bold))
	printTokens(tokens)
	fmt.Print("\n")

	fmt.Println(ansi("## DEFINITIONS ##", bold))
	fmt.Println("Name: " + workflowName)
	fmt.Println("Color: " + iconColor)
	fmt.Printf("Glyph: %d\n", iconGlyph)
	fmt.Printf("Inputs: %v\n", inputs)
	fmt.Printf("Outputs: %v\n", outputs)
	fmt.Printf("Workflows: %v\n", types)
	fmt.Printf("No Input: %v\n", noInput)
	fmt.Printf("macOS Only: %v\n", isMac)
	fmt.Printf("Mininum Version: %s\n", minVersion)
	fmt.Printf("iOS Version: %.1f\n", iosVersion)
	fmt.Print("\n")

	fmt.Println(ansi("## VARIABLES ##", bold))
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
			fmt.Printf(" (%s)", typeName(v.valueType))
		}
		if v.repeatItem {
			fmt.Print(" (repeat item var)")
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")

	fmt.Println(ansi("## MENUS ##", bold))
	fmt.Println(menus)
	fmt.Print("\n")

	fmt.Println(ansi("## IMPORT QUESTIONS ##", bold))
	fmt.Println(questions)
	fmt.Print("\n")
}

func parse() {
	for char != -1 {
		switch {
		case char == ' ' || char == '\t' || char == '\n':
			advance()
		case tokenAhead(Question):
			collectQuestion()
		case tokenAhead(Definition):
			collectDefinition()
		case tokenAhead(Import):
			collectImport()
		case isToken(At):
			collectVariable(false)
		case tokenAhead(Constant):
			advance()
			collectVariable(true)
		case isToken(ForwardSlash):
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
			collectConditional()
		case tokenAhead(RightBrace):
			collectEndStatement()
		case strings.Contains(lookAheadUntil(' '), "("):
			collectActionCall()
		default:
			parserError(fmt.Sprintf("Illegal character '%s'", string(char)))
		}
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
	if contains(stoppers, lastActionIdentifier) {
		parserWarning(fmt.Sprintf("Statement appears to be unreachable or does not loop as %s() was called outside of conditional.", lastActionIdentifier))
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
			insideString = insideString && prev(1) == '\\'
		}
		collected.WriteRune(char)
		advance()
	}

	return strings.Trim(collected.String(), " ")
}

// collectUntil advances ahead until the current character is `ch`,
// This should be used in cases where we are unsure how many characters will occur before we reach this character.
// For instance a string collector would need to use this.
func collectUntil(ch rune) string {
	var collected strings.Builder
	for char != ch && char != -1 {
		collected.WriteRune(char)
		advance()
	}

	return strings.Trim(collected.String(), " ")
}

// collectUntilExpect advances ahead until the current character is `ch`.
// If we advance more times than `maxAdvances` before finding `ch`, we throw
// an error that we expected `ch` and return the characters collected.
func collectUntilExpect(ch rune, maxAdvances int) string {
	var collected strings.Builder
	var advances int
	for char != ch && char != -1 {
		if advances > maxAdvances {
			parserError(fmt.Sprintf("Expected %c, got: %s", ch, collected.String()))
		}
		collected.WriteRune(char)
		advances++
		advance()
	}
	return collected.String()
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

	return strings.Trim(strings.ToLower(ahead.String()), " \t\n")
}

func collectVariableValue(constant bool, valueType *tokenType, value *any, varType *tokenType, coerce *string, getAs *string) {
	if tokenAhead(AddTo) {
		*varType = AddTo
		if constant {
			parserError("Constants cannot be added to.")
		}
	} else if !tokensAhead(Set) && constant {
		parserError("Constants must be initialized with a value.")
	}

	advance()
	collectValue(valueType, value, '\n')

	if constant && (*valueType == Arr || *valueType == Variable) {
		lineIdx--
		var valueTypeName = capitalize(typeName(*valueType))
		parserError(fmt.Sprintf("%v values cannot be constants.", valueTypeName))
	}
	if *valueType == Question {
		parserError(fmt.Sprintf("Illegal reference to import question '%s'. Shortcuts does not support import questions as variable values.", *value))
	}
	if *valueType != Variable {
		return
	}

	var stringValue = fmt.Sprintf("%v", *value)
	if !strings.ContainsAny(stringValue, "[]") && !strings.Contains(stringValue, ".") {
		return
	}

	var regex = regexp.MustCompile(`^(.*?)(?:\[(.*?)])?(?:\.(.*?))?$`)
	var matches = regex.FindAllStringSubmatch(stringValue, -1)
	for _, m := range matches {
		*value = m[1]
		if m[2] != "" {
			*getAs = m[2]
		}
		if m[3] != "" {
			*coerce = m[3]
		}
	}
}

func collectValue(valueType *tokenType, value *any, until rune) {
	var ahead = lookAheadUntil(until)
	if ahead == "" {
		parserError("Value expected")
	}
	switch {
	case intChar():
		collectIntegerValue(valueType, value, &until)
	case isToken(String):
		*valueType = String
		*value = collectString()
	case isToken(Arr):
		*valueType = Arr
		*value = collectArray()
	case isToken(Dict):
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
	case strings.Contains(lookAheadUntil(until), "("):
		*valueType = Action
		_, *value = collectAction()
	case containsTokens(lookAheadUntil(until), Plus, Minus, Multiply, Divide, Modulus):
		*valueType = Expression
		*value = collectUntil(until)
	default:
		collectReference(valueType, value, &until)
	}
}

func collectReference(valueType *tokenType, value *any, until *rune) {
	var identifier strings.Builder
	identifier.WriteString(collectIdentifier())

	if q, found := questions[identifier.String()]; found {
		if q.used {
			parserError(fmt.Sprintf("Duplicate usage of import question reference '%s', can only be used once.", identifier.String()))
		}

		*valueType = Question
		*value = identifier.String()
		q.used = true

		advance()
		return
	}

	if !validReference(identifier.String()) {
		parserError(fmt.Sprintf("Undefined reference '%s'", identifier.String()))
	}

	if char == '[' {
		identifier.WriteString(collectUntil(']') + `]`)
		advance()
	}
	if char == '.' {
		identifier.WriteString(collectUntil(*until))
	}

	*valueType = Variable
	*value = identifier.String()
	advance()
	return
}

func collectArguments() (arguments []actionArgument) {
	var params = actions[currentAction].parameters
	var paramsSize = len(params)
	var argIndex = 0
	var param parameterDefinition
	for {
		if char == ')' || char == '\n' || char == -1 {
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
			fmt.Sprintf("Too many arguments for action %s()\n\n%s",
				currentAction,
				generateActionDefinition(parameterDefinition{}, false, false),
			),
		)
	}
	if char == ',' {
		advance()
	}
	if char == ' ' {
		advance()
	}
	var valueType tokenType
	var value any
	if strings.Contains(lookAheadUntil('\n'), ",") {
		collectValue(&valueType, &value, ',')
	} else {
		collectValue(&valueType, &value, ')')
	}
	argument = actionArgument{
		valueType: valueType,
		value:     value,
	}
	if !param.infinite {
		checkArg(param, &argument)
	}
	return
}

func collectComment() {
	var collect = args.Using("comments")
	var comment strings.Builder
	if isToken(ForwardSlash) {
		if collect {
			comment.WriteString(collectUntil('\n'))
		} else {
			advanceUntil('\n')
		}
	} else {
		advanceTimes(2)
		for {
			if char == '*' && next(1) == '/' {
				break
			}
			if collect {
				comment.WriteRune(char)
			}
			advance()
		}
		advanceTimes(3)
	}
	if collect {
		var commentStr = strings.Trim(comment.String(), " \n")
		tokens = append(tokens, token{
			typeof:    Comment,
			ident:     "",
			valueType: String,
			value:     commentStr,
		})
	}
}

func collectVariable(constant bool) {
	reachable()

	var identifier = collectIdentifier()

	if v, found := variables[identifier]; found {
		if v.constant {
			parserError(fmt.Sprintf("Cannot redefine constant '%s'.", identifier))
		}
		if v.repeatItem {
			parserError(fmt.Sprintf("Cannot redefine repeat item '%s'.", identifier))
		}
	}
	if _, found := globals[identifier]; found {
		parserError(fmt.Sprintf("Cannot redefine global variable '%s'.", identifier))
	}
	if _, found := questions[identifier]; found {
		parserError(fmt.Sprintf("Variable conflicts with defined import question '%s'.", identifier))
	}
	var valueType tokenType
	var value any
	var getAs string
	var coerce string
	var varType = Var
	if strings.Contains(lookAheadUntil('\n'), "=") {
		advance()
		collectVariableValue(constant, &valueType, &value, &varType, &coerce, &getAs)
	} else if tokenAhead(Colon) {
		advance()
		switch {
		case tokenAhead(VarTextType):
			valueType = String
		case tokenAhead(VarNumberType):
			valueType = Integer
		case tokenAhead(VarBoolType):
			valueType = Bool
		case tokenAhead(VarArrayType):
			valueType = Arr
		case tokenAhead(VarDictType):
			valueType = Dict
		case tokenAhead(VarVariableType):
			valueType = VarVariableType
		default:
			parserError(fmt.Sprintf("Unknown type '%s'", lookAheadUntil('\n')))
		}
	}

	tokens = append(tokens, token{
		typeof:    varType,
		ident:     identifier,
		valueType: valueType,
		value:     value,
	})
	if varType != Var {
		return
	}
	variables[identifier] = variableValue{
		variableType: "Variable",
		valueType:    valueType,
		value:        value,
		getAs:        getAs,
		coerce:       coerce,
		constant:     constant,
	}
}

func collectIdentifier() string {
	var identifier strings.Builder
	for char != -1 {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
			break
		}
		identifier.WriteRune(char)
		advance()
	}

	return identifier.String()
}

func collectDefinition() {
	advance()
	switch {
	case tokenAhead(Name):
		advance()
		workflowName = collectUntil('\n')
		outputPath = relativePath + workflowName + ".shortcut"
	case tokenAhead(Color):
		advance()
		var collectColor = collectUntil('\n')
		makeColors()
		collectColor = strings.ToLower(collectColor)
		if color, found := colors[collectColor]; found {
			iconColor = color
		} else {
			var list = makeKeyList("Available icon colors:", colors)
			parserError(fmt.Sprintf("Invalid icon color '%s'\n\n%s", collectColor, list))
		}
	case tokenAhead(Glyph):
		advance()
		var collectGlyph = collectUntil('\n')
		makeGlyphs()
		collectGlyph = strings.ToLower(collectGlyph)
		if glyph, found := glyphs[collectGlyph]; found {
			glyphInt, hexErr := strconv.ParseInt(fmt.Sprintf("%d", glyph), 10, 64)
			handle(hexErr)
			iconGlyph = glyphInt
		} else {
			var list = "Available icon glyphs:\n"
			for key := range glyphs {
				list += "- " + key + "\n"
			}
			parserError(fmt.Sprintf("Invalid icon glyph '%s'\n\n%s", collectGlyph, list))
		}
	case tokenAhead(Inputs):
		advance()
		var collectInputs = collectUntil('\n')
		if collectInputs != "" {
			var inputTypes = strings.Split(collectInputs, ",")
			for _, input := range inputTypes {
				input = strings.Trim(input, " ")
				makeContentItems()
				if contentItem, found := contentItems[input]; found {
					inputs = append(inputs, contentItem)
				} else {
					var list = makeKeyList("Available content item types:", contentItems)
					parserError(fmt.Sprintf("Invalid input type '%s'\n\n%s", input, list))
				}
			}
		}
	case tokenAhead(Outputs):
		advance()
		var collectOutputs = collectUntil('\n')
		if collectOutputs != "" {
			var outputTypes = strings.Split(collectOutputs, ",")
			for _, output := range outputTypes {
				output = strings.Trim(output, " ")
				makeContentItems()
				if contentItem, found := contentItems[output]; found {
					outputs = append(outputs, contentItem)
				} else {
					var list = makeKeyList("Available content item types:", contentItems)
					parserError(fmt.Sprintf("Invalid output type '%s'\n\n%s", output, list))
				}
			}
		}
	case tokenAhead(From):
		advance()
		makeWorkflowTypes()
		var collectWorkflowTypes = collectUntil('\n')
		if collectWorkflowTypes != "" {
			var definedWorkflowTypes = strings.Split(collectWorkflowTypes, ",")
			for _, wt := range definedWorkflowTypes {
				wt = strings.Trim(wt, " ")
				if wtype, found := workflowTypes[wt]; found {
					types = append(types, wtype)
				} else {
					var list = makeKeyList("Available workflow types:", workflowTypes)
					parserError(fmt.Sprintf("Invalid workflow type '%s'\n\n%s", wt, list))
				}
			}
		}
	case tokenAhead(NoInput):
		advance()
		switch {
		case tokenAhead(StopWith):
			advance()
			var stopWithError = collectString()
			noInput = noInputParams{
				name: "WFWorkflowNoInputBehaviorShowError",
				params: []plistData{
					{
						key:      "Error",
						dataType: Text,
						value:    stopWithError,
					},
				},
			}
		case tokenAhead(AskFor):
			advance()
			var workflowType = collectUntil('\n')
			makeContentItems()
			if wtype, found := contentItems[workflowType]; found {
				noInput = noInputParams{
					name: "WFWorkflowNoInputBehaviorAskForInput",
					params: []plistData{
						{
							key:      "ItemClass",
							dataType: Text,
							value:    wtype,
						},
					},
				}
			} else {
				var list = makeKeyList("Available workflow types:", workflowTypes)
				parserError(fmt.Sprintf("Invalid workflow type '%s'\n\n%s", wtype, list))
			}
		case tokenAhead(GetClipboard):
			noInput = noInputParams{
				name:   "WFWorkflowNoInputBehaviorGetClipboard",
				params: []plistData{},
			}
		}
	case tokenAhead(Mac):
		var defValue = collectUntil('\n')
		switch defValue {
		case "true":
			isMac = true
		case "false":
			isMac = false
		default:
			parserError(fmt.Sprintf("Invalid value of '%s' for boolean definition 'mac'", defValue))
		}
	case tokenAhead(Version):
		var collectVersion = collectUntil('\n')
		makeVersions()
		if version, found := versions[collectVersion]; found {
			minVersion = version
			iosVersion, _ = strconv.ParseFloat(collectVersion, 8)
		} else {
			var list = makeKeyList("Available versions:", versions)
			parserError(fmt.Sprintf("Invalid minimum version '%s'\n\n%s", collectVersion, list))
		}
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
	if !isToken("\"") {
		parserError("Expected question prompt string.")
	}
	var text = collectString()
	advance()
	if !isToken("\"") {
		parserError("Expected question default string value.")
	}
	var defaultValue = collectString()
	questions[identifier] = &question{
		text:         text,
		defaultValue: defaultValue,
	}
}

var repeatIndexIndex = 1
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
	tokenAhead(RepeatWithEach)

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
			typeof:    Var,
			ident:     repeatIndexIdentifier,
			valueType: Variable,
			value:     fmt.Sprintf("Repeat Index%s", index),
		},
	)

	variables[repeatIndexIdentifier] = variableValue{
		variableType: "Variable",
		valueType:    String,
		value:        repeatIndexIdentifier,
		repeatItem:   true,
	}

	repeatIndexIndex++
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
	tokenAhead(In)
	advance()

	var iterableType tokenType
	var iterableValue any
	collectValue(&iterableType, &iterableValue, '{')

	advance()
	tokens = append(tokens,
		token{
			typeof:    RepeatWithEach,
			ident:     groupingUUID,
			valueType: iterableType,
			value:     iterableValue,
		}, token{
			typeof:    Var,
			ident:     repeatItemIdentifier,
			valueType: Variable,
			value:     fmt.Sprintf("Repeat Item%s", index),
		},
	)

	variables[repeatItemIdentifier] = variableValue{
		variableType: "Variable",
		valueType:    String,
		value:        repeatItemIdentifier,
		repeatItem:   true,
	}

	repeatItemIndex++
	repeatIndexIndex++
}

func collectConditional() {
	reachable()
	advance()
	makeConditions()

	var groupingUUID = groupStatement(Conditional)

	var conditionType string
	if isToken(Exclamation) {
		conditionType = conditions[Empty]
	} else {
		conditionType = conditions[Any]
	}

	var variableOneType tokenType
	var variableOneValue any
	var variableTwoType tokenType
	var variableTwoValue any
	var variableThreeType tokenType
	var variableThreeValue any
	collectValue(&variableOneType, &variableOneValue, ' ')

	if !isToken(LeftBrace) {
		var collectConditional = collectUntil(' ')
		var collectConditionalToken = tokenType(collectConditional)
		if condition, found := conditions[collectConditionalToken]; found {
			conditionType = condition
		} else {
			parserError(fmt.Sprintf("Invalid conditional '%s'", collectConditional))
		}
		advance()
		collectValue(&variableTwoType, &variableTwoValue, ' ')
		if char == ' ' {
			advance()
		}
		if !isToken(LeftBrace) {
			collectValue(&variableThreeType, &variableThreeValue, '{')
			advance()
		}
	}
	isToken(LeftBrace)

	tokens = append(tokens, token{
		typeof:    Conditional,
		ident:     groupingUUID,
		valueType: If,
		value: condition{
			variableOneType:    variableOneType,
			variableOneValue:   variableOneValue,
			condition:          conditionType,
			variableTwoType:    variableTwoType,
			variableTwoValue:   variableTwoValue,
			variableThreeType:  variableThreeType,
			variableThreeValue: variableThreeValue,
		},
	})
}

func collectMenu() {
	reachable()
	advance()
	var groupingUUID = groupStatement(Menu)

	var promptType tokenType
	var promptValue any
	collectValue(&promptType, &promptValue, '{')
	advanceUntil('{')
	advance()

	menus[groupingUUID] = []variableValue{}
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

	menus[groupingUUID] = append(menus[groupingUUID], variableValue{
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
			reachable()
			repeatItemIndex--
			repeatIndexIndex--
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
	groupingUUID = shortcutsUUID()
	groupingIdx++
	groupingUUIDs[groupingIdx] = groupingUUID
	groupingTypes[groupingIdx] = groupType

	return
}

// addNothing helps reduce memory usage by not passing anything to the next action.
func addNothing() {
	standardActions()
	tokens = append(tokens, token{
		typeof:    Action,
		ident:     "nothing",
		valueType: Action,
		value: action{
			ident: "nothing",
		},
	})
}

const intTypeString = string(Integer)

func intChar() bool {
	var charStr = string(char)
	return strings.Contains(intTypeString, charStr)
}

func collectInteger() string {
	var integer strings.Builder
	for intChar() {
		integer.WriteRune(char)
		advance()
	}
	return integer.String()
}

func collectIntegerValue(valueType *tokenType, value *any, until *rune) {
	if !containsTokens(lookAheadUntil(*until), Plus, Minus, Multiply, Divide, Modulus) {
		var integer = collectInteger()
		*valueType = Integer
		*value = integer
		advance()
		return
	}
	*valueType = Expression
	*value = collectUntil(*until)
}

func collectString() string {
	var collection strings.Builder
	escaped := false
	for char != -1 {
		if escaped {
			switch char {
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
				// handle other cases or throw an error if needed
			}
			escaped = false
		} else {
			if char == '\\' {
				escaped = true
			} else if char == '"' {
				break
			} else {
				collection.WriteRune(char)
			}
		}
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
	advance()
	return
}

func collectObject() string {
	var jsonStr strings.Builder
	var insideInnerObject = false
	var insideString = false
	for {
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
				insideInnerObject = true
			} else if char == '}' {
				if !insideInnerObject {
					break
				}
				insideInnerObject = false
			}
		}
		jsonStr.WriteRune(char)
		advance()
	}
	return jsonStr.String()
}

func collectActionCall() {
	reachable()
	var identifier, value = collectAction()
	tokens = append(tokens, token{
		typeof:    Action,
		ident:     identifier,
		valueType: Action,
		value:     value,
	})
}

func collectAction() (identifier string, value action) {
	standardActions()

	identifier = collectIdentifier()
	if _, found := actions[identifier]; !found {
		lineIdx--
		parserError(fmt.Sprintf("Undefined action '%s()'", identifier))
	}
	advance()
	currentAction = identifier

	var arguments = collectArguments()
	currentArguments = arguments
	currentArgumentsSize = len(currentArguments)

	checkAction()

	value = action{
		ident: identifier,
		args:  arguments,
	}

	if char == ')' {
		advance()
	}
	return
}

// advance advances the character cursor.
func advance() {
	idx++
	if len(chars) <= idx {
		char = -1
		return
	}

	char = chars[idx]
	if char == '\n' {
		lineCharIdx = 0
		lineIdx++
	} else {
		lineCharIdx++
	}
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

func isToken(token tokenType) bool {
	var tokenChar = []rune(token)[0]
	if unicode.ToLower(char) != tokenChar {
		return false
	}
	advance()
	return true
}

func tokenAhead(token tokenType) bool {
	var tokenLen = len(token)
	if tokenLen == 1 && unicode.ToLower(char) == []rune(token)[0] {
		advance()
		return true
	}
	for i, tokenChar := range token {
		if tokenChar == '\t' || tokenChar == '\n' {
			continue
		}

		if (i == 0 && unicode.ToLower(char) != tokenChar) || next(i) != tokenChar {
			return false
		}
	}

	if tokenLen == 1 {
		advance()
		return true
	}

	advanceTimes(tokenLen)
	return true
}

func tokensAhead(v ...tokenType) bool {
	for _, aheadToken := range v {
		if tokenAhead(aheadToken) {
			return true
		}
	}
	return false
}

func containsTokens(str string, v ...tokenType) bool {
	for _, aheadToken := range v {
		if strings.Contains(str, string(aheadToken)) {
			return true
		}
	}
	return false
}

func tokensOccur(str *string, v ...tokenType) bool {
	for _, aheadToken := range v {
		if strings.Count(*str, string(aheadToken)) != 0 {
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

func seek(mov *int, reverse bool) (requestedChar rune) {
	var nextChar = idx
	if reverse {
		nextChar -= *mov
	} else {
		nextChar += *mov
	}
	requestedChar = getChar(nextChar)
	for requestedChar == '\t' || requestedChar == '\n' {
		if reverse {
			nextChar--
		} else {
			nextChar++
		}
		requestedChar = getChar(nextChar)
	}
	return
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
	lineCharIdx = 0
	idx = -1
	advance()
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
	fmt.Println(ansi("\nWarning: ", yellow, bold) + fmt.Sprintf("%s %s:%d:%d", message, errorFilename, errorLine, errorCol))
}

func makeKeyList(title string, list map[string]string) (formattedList string) {
	formattedList = title + "\n"
	for key := range list {
		formattedList += "- " + key + "\n"
	}
	return
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
