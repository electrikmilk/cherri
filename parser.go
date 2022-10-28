/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func parse() {
	variables = make(map[string]variableValue)
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
		case tokenAhead(Definition):
			switch {
			case tokenAhead(Name):
				workflowName = collectUntil('\n')
			case tokenAhead(Color):
				var collectColor = collectUntil('\n')
				makeColors()
				collectColor = strings.ToLower(collectColor)
				if _, found := colors[collectColor]; found {
					iconColor = colors[collectColor]
				} else {
					var list = makeKeyList("Available icon colors:", colors)
					parserError(fmt.Sprintf("Invalid icon color '%s'\n\n%s\n", collectColor, list))
				}
			case tokenAhead(Glyph):
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
					parserError(fmt.Sprintf("Invalid icon glyph '%s'\n\n%s\n", collectGlyph, list))
				}
			case tokenAhead(Inputs):
				if len(contentItems) == 0 {
					makeContentItems()
				}
				var collectInputs = collectUntil('\n')
				if collectInputs != "" {
					var definedInputs = strings.Split(collectInputs, ",")
					for _, input := range definedInputs {
						input = strings.Trim(input, " ")
						if _, found := contentItems[input]; found {
							inputs = append(inputs, contentItems[input])
						} else {
							var list = makeKeyList("Available content items:", contentItems)
							parserError(fmt.Sprintf("Invalid input content item '%s'\n\n%s\n", input, list))
						}
					}
				}
			case tokenAhead(Outputs):
				if len(contentItems) == 0 {
					makeContentItems()
				}
				var collectOutputs = collectUntil('\n')
				if collectOutputs != "" {
					var definedOutputs = strings.Split(collectOutputs, ",")
					for _, output := range definedOutputs {
						output = strings.Trim(output, " ")
						if _, found := contentItems[output]; found {
							outputs = append(outputs, contentItems[output])
						} else {
							var list = makeKeyList("Available content items:", contentItems)
							parserError(fmt.Sprintf("Invalid output content item '%s'\n\n%s\n", output, list))
						}
					}
				}
			case tokenAhead(From):
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
							parserError(fmt.Sprintf("Invalid workflow type '%s'\n\n%s\n", wtype, list))
						}
					}
				}
			case tokenAhead(NoInput):
				switch {
				case tokenAhead(StopWith):
					lineCharIdx -= 2
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
					makeContentItems()
					var wtype = collectUntil('\n')
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
						parserError(fmt.Sprintf("Invalid workflow type '%s'\n\n%s\n", wtype, list))
					}
				case tokenAhead(GetClipboard):
					noInput = noInputParams{
						name:   "WFWorkflowNoInputBehaviorGetClipboard",
						params: []plistData{},
					}
				}
			case tokenAhead(Version):
				var collectVersion = collectUntil('\n')
				makeVersions()
				if _, found := versions[collectVersion]; found {
					minimumVersion = versions[collectVersion]
				} else {
					parserError(fmt.Sprintf("Invalid minimum version '%s'", collectVersion))
				}
			}
		case tokenAhead(Import):
			makeLibraries()
			var collectedLibrary = collectUntil('\n')
			if _, found := libraries[collectedLibrary]; found {
				libraries[collectedLibrary].make()
			}
		case tokenAhead(Var):
			idx -= 2
			advance()
			var identifier string
			if strings.Contains(lookAheadUntil('\n'), "=") {
				identifier = collectUntil(' ')
			} else {
				identifier = collectUntil('\n')
				idx -= 2
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
				col:       idx,
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
					col:       idx,
					typeof:    Comment,
					ident:     "",
					valueType: String,
					value:     collectComment(singleLine),
				})
			} else {
				tokens = append(tokens, token{
					col:       idx,
					typeof:    Comment,
					ident:     "",
					valueType: String,
					value:     collectComment(multiLine),
				})
			}
		case tokenAhead(Repeat):
			groupingUUID = shortcutsUUID()
			closureIdx++
			closureUUIDs[closureIdx] = groupingUUID
			closureTypes[closureIdx] = Repeat
			var timesType tokenType
			var timesValue any
			collectValue(&timesType, &timesValue, '{')
			advance()
			tokens = append(tokens, token{
				col:       idx,
				typeof:    Repeat,
				ident:     groupingUUID,
				valueType: timesType,
				value:     timesValue,
			})
		case tokenAhead(RepeatWithEach):
			groupingUUID = shortcutsUUID()
			closureIdx++
			closureUUIDs[closureIdx] = groupingUUID
			closureTypes[closureIdx] = RepeatWithEach
			var iterableType tokenType
			var iterableValue any
			collectValue(&iterableType, &iterableValue, '{')
			advance()
			tokens = append(tokens, token{
				col:       idx,
				typeof:    RepeatWithEach,
				ident:     groupingUUID,
				valueType: iterableType,
				value:     iterableValue,
			})
		case tokenAhead(Menu):
			groupingUUID = shortcutsUUID()
			closureIdx++
			closureUUIDs[closureIdx] = groupingUUID
			closureTypes[closureIdx] = Menu
			var promptType tokenType
			var promptValue any
			collectValue(&promptType, &promptValue, '{')
			collectUntil('{')
			advance()
			menus[groupingUUID] = []variableValue{}
			tokens = append(tokens, token{
				col:       idx,
				typeof:    Menu,
				ident:     groupingUUID,
				valueType: promptType,
				value:     promptValue,
			})
		case tokenAhead(Case):
			if groupingUUID == "" {
				parserError("Case has no starting menu statement.")
			}
			var itemType tokenType
			var itemValue any
			collectValue(&itemType, &itemValue, ':')
			collectUntil(':')
			advance()
			menus[groupingUUID] = append(menus[groupingUUID], variableValue{
				valueType: itemType,
				value:     itemValue,
			})
			tokens = append(tokens, token{
				col:       idx,
				typeof:    Case,
				ident:     groupingUUID,
				valueType: itemType,
				value:     itemValue,
			})
		case tokenAhead(If):
			if len(conditions) == 0 {
				makeConditions()
			}
			groupingUUID = shortcutsUUID()
			closureIdx++
			closureUUIDs[closureIdx] = groupingUUID
			closureTypes[closureIdx] = Conditional
			var cond string
			if tokensAhead(Exclamation) {
				cond = conditions[Empty]
				idx -= 2
				advance()
			} else {
				cond = conditions[Any]
			}
			var variableOneType tokenType
			var variableOneValue any
			var variableTwoType tokenType
			var variableTwoValue any
			var variableThreeType tokenType
			var variableThreeValue any
			idx--
			advance()
			collectValue(&variableOneType, &variableOneValue, ' ')
			if !isToken(LeftBrace) {
				var collectConditional = collectUntil(' ')
				var collectConditionalToken = tokenType(collectConditional)
				if _, found := conditions[collectConditionalToken]; found {
					cond = conditions[collectConditionalToken]
				} else {
					parserError(fmt.Sprintf("Invalid conditional '%s'", collectConditional))
				}
				collectValue(&variableTwoType, &variableTwoValue, ' ')
				if !isToken(LeftBrace) {
					collectValue(&variableThreeType, &variableThreeValue, '{')
				}
			}
			advance()
			tokens = append(tokens, token{
				col:       idx,
				typeof:    Conditional,
				ident:     groupingUUID,
				valueType: If,
				value: condition{
					variableOneType:    variableOneType,
					variableOneValue:   variableOneValue,
					condition:          cond,
					variableTwoType:    variableTwoType,
					variableTwoValue:   variableTwoValue,
					variableThreeType:  variableThreeType,
					variableThreeValue: variableThreeValue,
				},
			})
		case tokenAhead(RightBrace):
			if tokenAhead(Else) {
				if groupingUUID == "" {
					parserError("Else has no starting if statement.")
				}
				tokens = append(tokens, token{
					col:       idx,
					typeof:    Conditional,
					ident:     closureUUIDs[closureIdx],
					valueType: Else,
					value:     nil,
				})
				tokenAhead(LeftBrace)
			} else {
				if groupingUUID == "" {
					parserError("Ending } has no starting statement.")
				}
				tokens = append(tokens, token{
					col:       idx,
					typeof:    closureTypes[closureIdx],
					ident:     closureUUIDs[closureIdx],
					valueType: EndIf,
					value:     nil,
				})
				closureIdx--
			}
		case strings.Contains(lookAhead(), "("):
			if len(actions) == 0 {
				standardActions()
			}
			var identifier = collectUntil('(')
			if _, found := actions[identifier]; found {
				var arguments = collectArguments()
				currentAction = identifier
				checkAction(arguments)
				tokens = append(tokens, token{
					col:       idx,
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
		default:
			parserError(fmt.Sprintf("Illegal character '%s'", string(char)))
		}
	}
}

func collectUntil(ch rune) (collected string) {
	for char != ch && char != -1 {
		collected += string(char)
		advance()
	}
	advance()
	collected = strings.Trim(collected, " ")
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
	ahead = strings.Trim(strings.Trim(strings.ToLower(ahead), " "), "\n")
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
		if char == '\n' {
			idx -= 1
			advance()
		} else {
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
	case strings.Contains(lookAheadUntil('\n'), "("):
		var identifier = collectUntil('(')
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
			parserError(fmt.Sprintf("Invalid action '%s()'", identifier))
		}
	default:
		var identifier string
		var fullIdentifier string
		switch {
		case strings.Contains(lookAheadUntil(until), "["):
			identifier = collectUntil('[')
			fullIdentifier = identifier + "[" + collectUntil(until)
		case strings.Contains(lookAheadUntil(until), "."):
			identifier = collectUntil('.')
			fullIdentifier = identifier + "." + collectUntil(until)
		default:
			identifier = collectUntil(until)
			fullIdentifier = identifier
		}
		var lowerIdentifier = strings.ToLower(identifier)
		if _, global := globals[identifier]; global {
			*valueType = Variable
			*value = fullIdentifier
			hasInputVariables(identifier)
		} else if _, found := variables[lowerIdentifier]; found {
			*valueType = Variable
			*value = fullIdentifier
		} else {
			lineIdx--
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
		if strings.Contains(lookAhead(), ",") {
			collectValue(&valueType, &value, ',')
		} else {
			collectValue(&valueType, &value, ')')
		}
		arguments = append(arguments, actionArgument{
			valueType: valueType,
			value:     value,
		})
	}
	if char == ')' {
		advance()
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
		advanceTimes(2)
		advance()
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
	return collectUntil('"')
}

func collectArray(until rune) (array interface{}) {
	var rawJson = "{\"array\":["
	for char != until && char != -1 {
		rawJson += string(char)
		advance()
	}
	rawJson += "]}"
	if err := json.Unmarshal([]byte(rawJson), &array); err != nil {
		lineIdx -= 2
		parserErr(err)
	}
	array = array.(map[string]interface{})["array"]
	advance()
	return
}

func collectDictionary() (dictionary interface{}) {
	var rawJson = "{"
	var insideInnerObject = false
	for {
		rawJson += string(char)
		if char == '{' {
			insideInnerObject = true
		} else if char == '}' {
			if !insideInnerObject {
				break
			}
			insideInnerObject = false
		}
		advance()
	}
	if err := json.Unmarshal([]byte(rawJson), &dictionary); err != nil {
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
	var i int
	for i = 0; i < times; i++ {
		advance()
	}
}

func isToken(token tokenType) bool {
	if strings.ToLower(string(char)) == string(token) {
		var tokenLength = len(string(token))
		advanceTimes(tokenLength)
		return true
	} else {
		return false
	}
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
		advanceTimes(len(tokenChars) + 1)
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

func getChar(atIndex int) (returnChar rune) {
	if atIndex == -1 {
		return []rune(chars[0])[0]
	}
	if len(chars) > atIndex {
		returnChar = []rune(chars[atIndex])[0]
	} else {
		returnChar = -1
	}
	return
}

func parserErr(err error) {
	parserError(fmt.Sprintf("%s", err))
}

func parserWarning(message string) {
	fmt.Printf("\n\033[33mWarning: %s (%d:%d)\033[0m\n", message, lineIdx+1, lineCharIdx+1)
}

func parserError(message string) {
	lines = strings.Split(contents, "\n")
	if char == '\n' || prev(1) == '\n' {
		lineIdx--
	}
	fmt.Print("\n\n")
	fmt.Printf("\033[31m\033[1m%s\033[0m\n", message)
	fmt.Printf("\n\033[2m--- \033[0m%s:%d:%d\n", filePath, lineIdx+1, lineCharIdx+1)
	if lineIdx != -1 {
		if len(lines) > (lineIdx - 1) {
			fmt.Println("\033[2m" + fmt.Sprintf("%d | ", lineIdx) + lines[lineIdx-1] + "\u001B[0m")
		}
		if len(lines) > lineIdx {
			fmt.Println("\033[31m\033[1m" + fmt.Sprintf("%d | ", lineIdx+1) + lines[lineIdx] + "\u001B[0m")
		}
		var spaces string
		for i := 0; i < (lineCharIdx + 3); i++ {
			spaces += " "
		}
		fmt.Println("\033[31m" + spaces + "^\u001B[0m")
		if len(lines) > (lineIdx + 1) {
			fmt.Println("\033[2m" + fmt.Sprintf("%d | ", lineIdx+2) + lines[lineIdx+1] + "\u001B[0m")
		}
	}
	if arg("debug") {
		panic("parser error")
	} else {
		os.Exit(1)
	}
}

func makeKeyList(title string, list map[string]string) (formattedList string) {
	formattedList = title + "\n"
	for key := range list {
		formattedList += "- " + key + "\n"
	}
	return
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
	fmt.Println(currentChar, fmt.Sprintf("%d:%d", lineIdx+1, idx))
}
*/
