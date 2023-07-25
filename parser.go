/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/electrikmilk/args-parser"
)

var idx int
var lines []string
var chars []string
var char rune

var lineIdx int
var lineCharIdx int

var closureUUIDs map[int]string
var closureTypes map[int]tokenType
var closureIdx int
var currentGroupingUUID string

func parse() {
	tokenChars = make(map[tokenType][]string)
	if strings.Contains(contents, "action") {
		standardActions()
		parseCustomActions()
	}
	if args.Using("debug") {
		fmt.Printf("Parsing %s... ", filename)
	}
	variables = make(map[string]variableValue)
	questions = make(map[string]*question)
	menus = make(map[string][]variableValue)
	closureUUIDs = make(map[int]string)
	closureTypes = make(map[int]tokenType)
	makeGlobals()
	chars = strings.Split(contents, "")
	idx = -1
	advance()
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
			collectVariable()
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
			collectEndClosure()
		case strings.Contains(lookAheadUntil(' '), "("):
			reachable()
			var identifier, value = collectAction()
			tokens = append(tokens, token{
				typeof:    Action,
				ident:     identifier,
				valueType: Action,
				value:     value,
			})
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

func collectUntilIgnoreStrings(ch rune) (collected string) {
	var insideString = false
	for char != -1 {
		if char == ch && !insideString {
			break
		}
		if char == '"' {
			if insideString && prev(1) != '\\' {
				insideString = false
			} else {
				insideString = true
			}
		}
		collected += string(char)
		advance()
	}
	collected = strings.Trim(collected, " ")
	return
}

// collectUntil advances ahead until the current character is `ch`,
// This should be used in cases where we are unsure how many characters will occur before we reach this character.
// For instance a string collector would need to use this.
func collectUntil(ch rune) (collected string) {
	for char != ch && char != -1 {
		collected += string(char)
		advance()
	}
	collected = strings.Trim(collected, " ")
	return
}

// collectUntilExpect advances ahead until the current character is `ch`.
// If we advance more times than `maxAdvances` before finding `ch`, we throw
// an error that we expected `ch` and return the characters collected.
func collectUntilExpect(ch rune, maxAdvances int) (collected string) {
	var advances int
	for char != ch && char != -1 {
		if advances > maxAdvances {
			parserError(fmt.Sprintf("Expected %s, got: %s", string(ch), collected))
		}
		collected += string(char)
		advances++
		advance()
	}
	return
}

func lookAheadUntil(until rune) (ahead string) {
	var nextIdx = idx
	var nextChar rune
	for nextChar != until {
		if len(chars) > nextIdx {
			nextChar = []rune(chars[nextIdx])[0]
			ahead += chars[nextIdx]
			nextIdx++
		} else {
			break
		}
	}
	ahead = strings.Trim(strings.ToLower(ahead), " \t\n")
	return
}

func collectVariableValue(valueType *tokenType, value *any, varType *tokenType, coerce *string, getAs *string) {
	if tokenAhead(AddTo) {
		*varType = AddTo
	} else {
		tokensAhead(Set)
	}
	advance()
	collectValue(valueType, value, '\n')
	if *valueType != Variable {
		return
	}
	var stringValue = fmt.Sprintf("%v", *value)
	if strings.Contains(stringValue, ".") {
		var dotParts = strings.Split(stringValue, ".")
		*coerce = strings.Trim(dotParts[1], " ")
		if strings.Contains(stringValue, "[") {
			var varParts = strings.Split(dotParts[0], "[")
			*value = varParts[0]
			*getAs = strings.Trim(strings.TrimSuffix(varParts[1], "]"), " ")
		} else {
			*value = dotParts[0]
		}
	} else if strings.Contains(stringValue, "[") {
		var varParts = strings.Split(stringValue, "[")
		*value = varParts[0]
		*getAs = strings.TrimSuffix(varParts[1], "]")
	}
}

func collectValue(valueType *tokenType, value *any, until rune) {
	switch {
	case intChar():
		var integer = collectInteger()
		*valueType = Integer
		*value = integer
		advance()
		if !tokensAhead(Plus, Minus, Multiply, Divide, Modulus) {
			break
		}
		var expression string
		idx -= len(integer) + 4
		advance()
		*valueType = Expression
		expression = collectUntil(until)
		*value = expression
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
	case strings.Contains(lookAheadUntil(until), "("):
		*valueType = Action
		_, *value = collectAction()
	default:
		collectReference(valueType, value, &until)
	}
}

func collectReference(valueType *tokenType, value *any, until *rune) {
	if lookAheadUntil(*until) == "" {
		parserError("Value expected")
	}
	var identifier string
	var fullIdentifier string
	switch {
	case strings.Contains(lookAheadUntil(*until), "["):
		identifier = collectUntil('[')
		advance()
		fullIdentifier = identifier + "[" + collectUntil(*until)
		advance()
	case strings.Contains(lookAheadUntil(*until), "."):
		identifier = collectUntil('.')
		advance()
		fullIdentifier = identifier + "." + collectUntil(*until)
		advance()
	default:
		identifier = collectUntil(*until)
		fullIdentifier = identifier
		advance()
	}
	var lowerIdentifier = strings.ToLower(identifier)
	if _, global := globals[identifier]; global {
		*valueType = Variable
		*value = fullIdentifier
		isInputVariable(identifier)
	} else if _, found := variables[lowerIdentifier]; found {
		*valueType = Variable
		*value = fullIdentifier
	} else if _, found := questions[lowerIdentifier]; found {
		if questions[lowerIdentifier].used {
			parserError(fmt.Sprintf("Duplicate usage of '%s', import questions can only be referenced once.", fullIdentifier))
		}
		*valueType = Question
		*value = fullIdentifier
		questions[lowerIdentifier].used = true
	} else {
		if fullIdentifier == "" {
			parserError("Value expected")
		}
		if args.Using("debug") {
			fmt.Println("\nvariables", variables)
			fmt.Println("questions", questions)
		}
		parserError(fmt.Sprintf("Unknown value type: '%s'", fullIdentifier))
	}
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
		if argIndex == paramsSize && !param.infinite {
			parserError(fmt.Sprintf("Too many arguments for action %s() (max: %d)", currentAction, argIndex))
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
		var argument = actionArgument{
			valueType: valueType,
			value:     value,
		}
		if !param.infinite {
			checkArg(param, argument)
		}
		arguments = append(arguments, argument)
		argIndex++
	}
	return
}

func collectComment() {
	var comment string
	if isToken(ForwardSlash) {
		comment = collectUntil('\n')
	} else {
		advanceTimes(2)
		for {
			if char == '*' && next(1) == '/' {
				break
			}
			comment += string(char)
			advance()
		}
		comment = strings.Trim(comment, "\n")
		advanceTimes(3)
	}
	comment = strings.Trim(comment, " ")
	tokens = append(tokens, token{
		typeof:    Comment,
		ident:     "",
		valueType: String,
		value:     comment,
	})
}

func collectVariable() {
	reachable()
	var identifier string
	if strings.Contains(lookAheadUntil('\n'), "=") {
		identifier = collectUntil(' ')
		advance()
	} else {
		identifier = collectUntil('\n')
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
		collectVariableValue(&valueType, &value, &varType, &coerce, &getAs)
	}
	tokens = append(tokens, token{
		typeof:    varType,
		ident:     identifier,
		valueType: valueType,
		value:     value,
	})
	if varType == Var {
		variables[strings.ToLower(identifier)] = variableValue{
			variableType: "Variable",
			valueType:    valueType,
			value:        value,
			getAs:        getAs,
			coerce:       coerce,
		}
	}
}

func collectDefinition() {
	advance()
	switch {
	case tokenAhead(Name):
		advance()
		workflowName = collectUntil('\n')
	case tokenAhead(Color):
		advance()
		var collectColor = collectUntil('\n')
		makeColors()
		collectColor = strings.ToLower(collectColor)
		if _, found := colors[collectColor]; found {
			iconColor = colors[collectColor]
		} else {
			var list = makeKeyList("Available icon colors:", colors)
			parserError(fmt.Sprintf("Invalid icon color '%s'\n\n%s", collectColor, list))
		}
	case tokenAhead(Glyph):
		advance()
		var collectGlyph = collectUntil('\n')
		makeGlyphs()
		collectGlyph = strings.ToLower(collectGlyph)
		if _, found := glyphs[collectGlyph]; found {
			glyphInt, hexErr := strconv.ParseInt(fmt.Sprintf("%d", glyphs[collectGlyph]), 10, 64)
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
				if _, found := contentItems[input]; found {
					inputs = append(inputs, contentItems[input])
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
				if _, found := contentItems[output]; found {
					outputs = append(outputs, contentItems[output])
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
			for _, wtype := range definedWorkflowTypes {
				wtype = strings.Trim(wtype, " ")
				if _, found := workflowTypes[wtype]; found {
					types = append(types, workflowTypes[wtype])
				} else {
					var list = makeKeyList("Available workflow types:", workflowTypes)
					parserError(fmt.Sprintf("Invalid workflow type '%s'\n\n%s", wtype, list))
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
			var wtype = collectUntil('\n')
			makeContentItems()
			if _, found := contentItems[wtype]; found {
				noInput = noInputParams{
					name: "WFWorkflowNoInputBehaviorAskForInput",
					params: []plistData{
						{
							key:      "ItemClass",
							dataType: Text,
							value:    contentItems[wtype],
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
		if _, found := versions[collectVersion]; found {
			minVersion = versions[collectVersion]
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
	if _, found := libraries[collectedLibrary]; found {
		libraries[collectedLibrary].make(libraries[collectedLibrary].identifier)
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
	var identifier = collectUntilExpect(' ', 3)
	if _, found := questions[identifier]; found {
		parserError(fmt.Sprintf("Duplicate declaration of import question '%s'.", identifier))
	}
	if _, found := variables[identifier]; found {
		parserError(fmt.Sprintf("Import question conflicts with defined variable '%s'.", identifier))
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

func collectRepeat() {
	reachable()
	advance()
	currentGroupingUUID = shortcutsUUID()
	closureIdx++
	closureUUIDs[closureIdx] = currentGroupingUUID
	closureTypes[closureIdx] = Repeat
	var timesType tokenType
	var timesValue any
	collectValue(&timesType, &timesValue, '{')
	collectUntilExpect('{', 3)
	advance()
	tokens = append(tokens, token{
		typeof:    Repeat,
		ident:     currentGroupingUUID,
		valueType: timesType,
		value:     timesValue,
	})
}

func collectRepeatEach() {
	reachable()
	currentGroupingUUID = shortcutsUUID()
	closureIdx++
	closureUUIDs[closureIdx] = currentGroupingUUID
	closureTypes[closureIdx] = RepeatWithEach
	var iterableType tokenType
	var iterableValue any
	collectValue(&iterableType, &iterableValue, '{')
	advance()
	tokens = append(tokens, token{
		typeof:    RepeatWithEach,
		ident:     currentGroupingUUID,
		valueType: iterableType,
		value:     iterableValue,
	})
}

func collectConditional() {
	reachable()
	advance()
	makeConditions()
	currentGroupingUUID = shortcutsUUID()
	closureIdx++
	closureUUIDs[closureIdx] = currentGroupingUUID
	closureTypes[closureIdx] = Conditional
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
		if _, found := conditions[collectConditionalToken]; found {
			conditionType = conditions[collectConditionalToken]
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
		ident:     currentGroupingUUID,
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
	currentGroupingUUID = shortcutsUUID()
	closureIdx++
	closureUUIDs[closureIdx] = currentGroupingUUID
	closureTypes[closureIdx] = Menu
	var promptType tokenType
	var promptValue any
	collectValue(&promptType, &promptValue, '{')
	collectUntilExpect('{', 3)
	advance()
	menus[currentGroupingUUID] = []variableValue{}
	tokens = append(tokens, token{
		typeof:    Menu,
		ident:     currentGroupingUUID,
		valueType: promptType,
		value:     promptValue,
	})
}

func collectMenuItem() {
	advance()
	if currentGroupingUUID == "" {
		parserError("Case has no starting menu statement.")
	}
	var itemType tokenType
	var itemValue any
	collectValue(&itemType, &itemValue, ':')
	collectUntil(':')
	advance()
	menus[currentGroupingUUID] = append(menus[currentGroupingUUID], variableValue{
		valueType: itemType,
		value:     itemValue,
	})
	tokens = append(tokens, token{
		typeof:    Item,
		ident:     currentGroupingUUID,
		valueType: itemType,
		value:     itemValue,
	})
}

func collectEndClosure() {
	advance()
	if tokenAhead(Else) {
		advance()
		if currentGroupingUUID == "" {
			parserError("Else has no starting if statement.")
		}
		tokens = append(tokens, token{
			typeof:    Conditional,
			ident:     closureUUIDs[closureIdx],
			valueType: Else,
			value:     nil,
		})
		tokenAhead(LeftBrace)
	} else {
		if currentGroupingUUID == "" {
			parserError("Ending closure has no starting statement.")
		}
		var closureType = closureTypes[closureIdx]
		if closureType == Repeat || closureType == RepeatWithEach {
			reachable()
		}
		tokens = append(tokens, token{
			typeof:    closureType,
			ident:     closureUUIDs[closureIdx],
			valueType: EndClosure,
			value:     nil,
		})
		closureIdx--
	}
}

const intTypeString = string(Integer)

func intChar() bool {
	var charStr = string(char)
	return strings.Contains(intTypeString, charStr)
}

func collectInteger() (integer string) {
	for intChar() {
		integer += string(char)
		advance()
	}
	return
}

func collectString() (str string) {
	for char != -1 {
		if char == '"' && prev(1) != '\\' {
			break
		}
		if char == '\\' && next(1) == '"' {
			advance()
			continue
		}
		str += string(char)
		advance()
	}
	advance()
	str = strings.Trim(str, " ")
	return
}

func collectArray() (array interface{}) {
	var rawJSON = "{\"array\":[" + collectUntilIgnoreStrings(']') + "]}"
	if err := json.Unmarshal([]byte(rawJSON), &array); err != nil {
		if args.Using("debug") {
			fmt.Println(ansi("\n### COLLECTED ARRAY ###", bold))
			fmt.Println(rawJSON)
			fmt.Print("\n")
		}
		parserErr(err)
	}
	array = array.(map[string]interface{})["array"]
	advance()
	return
}

func collectDictionary() (dictionary interface{}) {
	var rawJSON = "{" + collectObject() + "}"
	if args.Using("debug") {
		fmt.Println(ansi("\n\n### COLLECTED DICTIONARY ###", bold))
		fmt.Println(rawJSON)
	}
	if err := json.Unmarshal([]byte(rawJSON), &dictionary); err != nil {
		parserErr(err)
	}
	advance()
	return
}

func collectObject() (jsonStr string) {
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
		jsonStr += string(char)
		advance()
	}
	return
}

func collectAction() (identifier string, value action) {
	standardActions()

	identifier = collectUntil('(')
	if _, found := actions[identifier]; !found {
		parserError(fmt.Sprintf("Unknown action '%s()'", identifier))
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

func advance() {
	idx++
	if len(chars) > idx {
		char = []rune(chars[idx])[0]
		if char == '\n' {
			lineCharIdx = 0
			lineIdx++
		} else {
			lineCharIdx++
		}
	} else {
		char = -1
	}
}

func advanceTimes(times int) {
	for i := 0; i < times; i++ {
		advance()
	}
}

func isToken(token tokenType) bool {
	if strings.ToLower(string(char)) != string(token) {
		return false
	}
	var tokenLength = len(string(token))
	advanceTimes(tokenLength)
	return true
}

func tokenAhead(token tokenType) bool {
	if len(token) == 1 && strings.ToLower(string(char)) == string(token) {
		advance()
		return true
	}
	var tChars []string
	if _, found := tokenChars[token]; found {
		tChars = tokenChars[token]
	} else {
		tChars = strings.Split(string(token), "")
		tokenChars[token] = tChars
	}
	for i, tokenChar := range tChars {
		if tokenChar == "\t" || tokenChar == "\n" {
			continue
		}
		if i == 0 {
			if strings.ToLower(string(char)) != tokenChar {
				return false
			}
		} else if next(i) != []rune(tokenChar)[0] {
			return false
		}
	}
	advanceTimes(len(token))
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
		return []rune(chars[atIndex])[0]
	}
	return -1
}

func firstChar() {
	lineIdx = 0
	lineCharIdx = 0
	idx = -1
	advance()
}

func parserErr(err error) {
	if err != nil {
		parserError(fmt.Sprintf("%s", err))
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
		printDebug()
		panic("debug")
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
