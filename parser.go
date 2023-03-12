/*
 * Copyright (c) 2023 Brandon Jordan
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
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
	variables = make(map[string]variableValue)
	questions = make(map[string]question)
	actions = make(map[string]actionDefinition)
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
			advance()
			var identifier = collectUntilExpect(' ', 3)
			identifier = strings.ToLower(identifier)
			advance()
			if !isToken("\"") {
				parserError("Expected string")
			}
			var text = collectString()
			advance()
			if !isToken("\"") {
				parserError("Expected string")
			}
			var defaultValue = collectString()
			questions[identifier] = question{
				text:         text,
				defaultValue: defaultValue,
			}
		case tokenAhead(Definition):
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
		case tokenAhead(Import):
			advance()
			makeLibraries()
			var collectedLibrary = strings.ToLower(collectUntil('\n'))
			if _, found := libraries[collectedLibrary]; found {
				libraries[collectedLibrary].make(libraries[collectedLibrary].identifier)
			} else {
				parserError(fmt.Sprintf("Import library '%s' does not exist!", collectedLibrary))
			}
		case isToken(At):
			var identifier string
			if strings.Contains(lookAheadUntil('\n'), "=") {
				identifier = collectUntil(' ')
				advance()
			} else {
				identifier = collectUntil('\n')
				advance()
			}
			var valueType tokenType
			var value any
			var getAs string
			var coerce string
			var varType = Var
			if strings.Contains(
				lookAheadUntil('\n'),
				"=",
			) {
				if tokenAhead(AddTo) {
					varType = AddTo
				} else {
					tokensAhead(Set)
				}
				advance()
				collectValue(&valueType, &value, '\n')
				if valueType == Variable {
					var stringValue = value.(string)
					if strings.Contains(stringValue, ".") {
						var dotParts = strings.Split(stringValue, ".")
						coerce = strings.Trim(dotParts[1], " ")
						if strings.Contains(stringValue, "[") {
							var varParts = strings.Split(dotParts[0], "[")
							value = varParts[0]
							getAs = strings.Trim(strings.TrimSuffix(varParts[1], "]"), " ")
						} else {
							value = dotParts[0]
						}
					} else if strings.Contains(stringValue, "[") {
						var varParts = strings.Split(stringValue, "[")
						value = varParts[0]
						getAs = strings.TrimSuffix(varParts[1], "]")
					}
				}
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
		case isToken(ForwardSlash):
			if isToken(ForwardSlash) {
				tokens = append(tokens, token{
					typeof:    Comment,
					ident:     "",
					valueType: String,
					value:     collectComment(singleLine),
				})
			} else {
				tokens = append(tokens, token{
					typeof:    Comment,
					ident:     "",
					valueType: String,
					value:     collectComment(multiLine),
				})
			}
		case tokenAhead(Repeat):
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
		case tokenAhead(RepeatWithEach):
			advance()
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
		case tokenAhead(Menu):
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
		case tokenAhead(Case):
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
				typeof:    Case,
				ident:     currentGroupingUUID,
				valueType: itemType,
				value:     itemValue,
			})
		case tokenAhead(If):
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
		case tokenAhead(RightBrace):
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
				tokens = append(tokens, token{
					typeof:    closureTypes[closureIdx],
					ident:     closureUUIDs[closureIdx],
					valueType: EndClosure,
					value:     nil,
				})
				closureIdx--
			}
		case strings.Contains(lookAhead(), "("):
			standardActions()
			var identifier = collectUntil('(')
			advance()
			if _, found := actions[identifier]; found {
				var arguments = collectArguments()
				currentAction = identifier
				checkAction(arguments)
				tokens = append(tokens, token{
					typeof:    Action,
					ident:     identifier,
					valueType: Action,
					value: action{
						ident: identifier,
						args:  arguments,
					},
				})
			} else {
				parserError(fmt.Sprintf("Unknown action '%s()'", identifier))
			}
			if char == ')' {
				advance()
			}
		default:
			parserError(fmt.Sprintf("Illegal character '%s'", string(char)))
		}
	}
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

func lookAhead() (ahead string) {
	return lookAheadUntil(' ')
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

func collectValue(valueType *tokenType, value *any, until rune) {
	switch {
	case strings.Contains(
		string(Integer),
		string(char),
	):
		var integer = collectInteger()
		*valueType = Integer
		*value = integer
		if char != '\n' {
			advance()
		}
		if tokensAhead(Plus, Minus, Multiply, Divide, Modulus) {
			var expression string
			idx -= len(integer) + 4
			advance()
			*valueType = Expression
			expression = collectUntil(until)
			*value = expression
		}
	case isToken(String):
		*valueType = String
		*value = collectString()
	case isToken(Arr):
		*valueType = Arr
		*value = collectArray(']')
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
		standardActions()
		var identifier = collectUntil('(')
		advance()
		if _, found := actions[identifier]; found {
			var arguments = collectArguments()
			currentAction = identifier
			checkAction(arguments)
			*valueType = Action
			*value = action{
				ident: identifier,
				args:  arguments,
			}
		} else {
			parserError(fmt.Sprintf("Unknown action as value '%s()'", identifier))
		}
		if char == ')' {
			advance()
		}
	default:
		if lookAheadUntil(until) == "" {
			parserError("Value expected")
		}
		var identifier string
		var fullIdentifier string
		switch {
		case strings.Contains(lookAheadUntil(until), "["):
			identifier = collectUntil('[')
			advance()
			fullIdentifier = identifier + "[" + collectUntil(until)
			advance()
		case strings.Contains(lookAheadUntil(until), "."):
			identifier = collectUntil('.')
			advance()
			fullIdentifier = identifier + "." + collectUntil(until)
			advance()
		default:
			identifier = collectUntil(until)
			fullIdentifier = identifier
			advance()
		}
		var lowerIdentifier = strings.ToLower(identifier)
		if _, global := globals[identifier]; global {
			*valueType = Variable
			*value = fullIdentifier
			hasInputVariables(identifier)
		} else if _, found := variables[lowerIdentifier]; found {
			*valueType = Variable
			*value = fullIdentifier
		} else if _, found := questions[lowerIdentifier]; found {
			*valueType = Question
			*value = fullIdentifier
		} else {
			if fullIdentifier == "" {
				parserError("Value expected")
			}
			parserError(fmt.Sprintf("Unknown value type: '%s'", fullIdentifier))
		}
	}
}

func collectArguments() (arguments []actionArgument) {
	for {
		if char == ')' || char == '\n' || char == -1 {
			break
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
		arguments = append(arguments, actionArgument{
			valueType: valueType,
			value:     value,
		})
	}
	return
}

type commentType int

const (
	singleLine commentType = iota
	multiLine
)

func collectComment(collect commentType) (comment string) {
	if collect == singleLine {
		comment = collectUntil('\n')
	} else if collect == multiLine {
		advance()
		for char != '*' && next(1) != '/' && char != -1 {
			comment += string(char)
			advance()
		}
		comment = strings.Trim(comment, "\n")
		advanceTimes(3)
	}
	comment = strings.Trim(comment, " ")
	return
}

func collectInteger() (integer string) {
	for strings.Contains(
		string(Integer),
		string(char),
	) {
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

func collectArray(until rune) (array interface{}) {
	var rawJSON = "{\"array\":["
	for char != until && char != -1 {
		rawJSON += string(char)
		advance()
	}
	rawJSON += "]}"
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
	var rawJSON = "{"
	var insideInnerObject = false
	var insideString = false
	for {
		rawJSON += string(char)
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
		advance()
	}
	if err := json.Unmarshal([]byte(rawJSON), &dictionary); err != nil {
		if args.Using("debug") {
			fmt.Println(ansi("\n\n### COLLECTED DICTIONARY ###", bold))
			fmt.Println(rawJSON)
		}
		parserErr(err)
	}
	advance()
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

func tokenAhead(token tokenType) (isAhead bool) {
	var tokenChars = strings.Split(string(token), "")
	isAhead = true
	for i, tchar := range tokenChars {
		if tchar == " " || tchar == "\t" || tchar == "\n" {
			continue
		}
		if i == 0 {
			if strings.ToLower(string(char)) != tchar {
				isAhead = false
				break
			}
		} else if next(i) != []rune(tchar)[0] {
			isAhead = false
			break
		}
	}
	if isAhead {
		advanceTimes(len(tokenChars))
	}
	return
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
	for requestedChar == ' ' || requestedChar == '\t' || requestedChar == '\n' {
		if reverse {
			nextChar -= 1
		} else {
			nextChar += 1
		}
		requestedChar = getChar(nextChar)
	}
	return
}

func getChar(atIndex int) rune {
	if atIndex == -1 {
		return []rune(chars[0])[0]
	}
	if len(chars) > atIndex {
		return []rune(chars[atIndex])[0]
	}
	return -1
}

type include struct {
	file   string
	start  int
	end    int
	lines  []string
	length int
}

var includes []include

func parseIncludes() {
	lines = strings.Split(contents, "\n")
	for l, line := range lines {
		if !strings.Contains(line, "#include") {
			continue
		}

		// Prepare for possible error
		chars = strings.Split(line, "")
		lineIdx = l
		idx = len("#include") + 1
		lineCharIdx = idx

		r := regexp.MustCompile("\"(.*?)\"")
		var includePath = strings.Trim(r.FindString(line), "\"")
		if includePath == "" {
			parserError("Expected file path")
		}

		if !strings.Contains(includePath, "..") {
			includePath = relativePath + includePath
		}

		if contains(included, includePath) {
			parserError(fmt.Sprintf("File '%s' has already been included.", includePath))
		}

		checkFile(includePath)
		var includeFileBytes, readErr = os.ReadFile(includePath)
		handle(readErr)

		var includeContents = string(includeFileBytes)
		var includeLines = strings.Split(includeContents, "\n")
		var includeLinesCount = len(includeLines)
		if strings.Contains(includeContents, "#include") {
			includeLinesCount--
		}

		for i, inc := range includes {
			if inc.start == l {
				includes[i].start = inc.start + includeLinesCount
				includes[i].end = (inc.start + includeLinesCount) + inc.length
			}
		}

		includes = append(includes, include{
			file:   includePath,
			start:  l,
			end:    l + includeLinesCount,
			lines:  includeLines,
			length: includeLinesCount,
		})

		lines[l] = includeContents
		included = append(included, includePath)
	}
	contents = strings.Join(lines, "\n")
	lineIdx = 0
	if strings.Contains(contents, "#include") {
		parseIncludes()
	}
}

func parserErr(err error) {
	if err != nil {
		parserError(fmt.Sprintf("%s", err))
	}
}

func delinquentFile() (errorFilename string, errorLine int, errorCol int) {
	errorFilename = filePath
	errorLine = lineIdx + 1
	errorCol = lineCharIdx + 1
	for _, inc := range includes {
		if lineIdx+1 >= inc.start && lineIdx+1 <= inc.end {
			errorFilename = inc.file
			for l, line := range inc.lines {
				if lineIdx < len(lines) {
					if line == lines[lineIdx] {
						errorLine = l + 1
					}
				}
			}
		}
	}
	return
}

func parserWarning(message string) {
	var errorFilename, errorLine, errorCol = delinquentFile()
	fmt.Print(ansi(fmt.Sprintf("\nWarning: %s %s:%d:%d\n", message, errorFilename, errorLine, errorCol), yellow))
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
	if lineIdx != -1 && !args.Using("no-ansi") {
		fmt.Print("\033[31m")
		fmt.Println("\n" + ansi(message, bold))
		var errorFilename, errorLine, errorCol = delinquentFile()
		fmt.Printf("\n\033[2m----- \033[0m%s:%d:%d\n", errorFilename, errorLine, errorCol)
		if len(lines) > (lineIdx-1) && lineIdx != 0 {
			fmt.Printf("\033[2m%d | %s\033[0m\n", lineIdx, lines[lineIdx-1])
		}
		if len(lines) > lineIdx {
			fmt.Printf("\033[31m\033[1m%d | ", lineIdx+1)
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
			fmt.Printf("\033[2m%d | %s\n-----\033[0m\n\n", lineIdx+2, lines[lineIdx+1])
		}
	} else {
		fmt.Printf("Error: %s (%d:%d)\n", message, lineIdx+1, lineCharIdx+1)
	}
	if args.Using("debug") {
		printDebug()
		panic("debug")
	} else {
		os.Exit(1)
	}
}

/*
func printCurrentChar() {
	var currentChar string
	switch char {
	case ' ':
		currentChar = "SPACE"
	case '\t':
		currentChar = "TAB"
	case '\n':
		currentChar = "LF"
	case -1:
		currentChar = "EMPTY"
	default:
		currentChar = string(char)
	}
	fmt.Printf("%s %d:%d\n", currentChar, lineIdx+1, lineCharIdx)
}
*/
