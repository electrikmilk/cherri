/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	"html"
	"reflect"
	"regexp"
	"strings"
)

var uuids map[string]string

type plistDataType string

const Text plistDataType = "string"
const Number plistDataType = "integer"
const Dictionary plistDataType = "dictionary"
const Array plistDataType = "array"
const Boolean plistDataType = "boolean"

type plistData struct {
	key      string
	dataType plistDataType
	value    any
}

type dictDataType string

const itemTypeText dictDataType = "0"
const itemTypeNumber dictDataType = "3"
const itemTypeArray dictDataType = "2"
const itemTypeDict dictDataType = "1"
const itemTypeBool dictDataType = "4"

type noInputParams struct {
	name   string
	params []plistData
}

var noInput noInputParams

var hasShortcutInputVariables = false

// ObjectReplaceChar is a Shortcuts convention to mark the placement of inline variables in a string.
const ObjectReplaceChar = '\uFFFC'
const ObjectReplaceCharStr = "\uFFFC"

var tabLevel = 0

func plistKeyValue(key string, dataType plistDataType, value any) string {
	var pair strings.Builder
	var tabs = strings.Repeat("\t", tabLevel)
	if key != "" {
		pair.WriteString(fmt.Sprintf("%s<key>%s</key>\n%s", tabs, key, tabs))
	} else {
		pair.WriteString(tabs)
	}
	switch dataType {
	case Boolean:
		pair.WriteString(fmt.Sprintf("<%t/>\n", value))
	case Array:
		if value == nil || len(value.([]plistData)) == 0 {
			pair.WriteString("<array/>\n")
			break
		}
		pair.WriteString(fmt.Sprintf("<array>\n%s%s</array>\n", plistDictValue(value.([]plistData)), strings.Repeat("\t", tabLevel)))
	case Dictionary:
		pair.WriteString(fmt.Sprintf("<dict>\n%s%s</dict>\n", plistDictValue(value.([]plistData)), strings.Repeat("\t", tabLevel)))
	default:
		if value != nil {
			if reflect.TypeOf(value).String() == "string" {
				value = html.EscapeString(value.(string))
			}
		}
		pair.WriteString(fmt.Sprintf("<%s>%v</%s>\n", dataType, value, dataType))
	}

	return pair.String()
}

var emptyPlistData = plistData{}

func plistDictValue(value []plistData) string {
	tabLevel++
	var pair strings.Builder
	for _, data := range value {
		if data == emptyPlistData {
			continue
		}

		pair.WriteString(plistKeyValue(data.key, data.dataType, data.value))
	}
	tabLevel--
	return pair.String()
}

// appendPlist grows and writes to the plist string builder.
func appendPlist(data []plistData) {
	var xmlStr = plistDictValue(data)
	plist.WriteString(xmlStr)
}

func conditionalParameter(key string, conditionalParams *[]plistData, typeOf *tokenType, value any) {
	if key == "" {
		if *typeOf == String {
			key = "WFConditionalActionString"
		} else if *typeOf == Integer || *typeOf == Bool {
			key = "WFNumberValue"
		}
	}
	switch *typeOf {
	case String:
		*conditionalParams = append(*conditionalParams, paramValue(key, actionArgument{
			valueType: *typeOf,
			value:     value,
		}, String, Text))
	case Integer:
		*conditionalParams = append(*conditionalParams, paramValue(key, actionArgument{
			valueType: *typeOf,
			value:     value,
		}, String, Text))
	case Bool:
		if value == true {
			*conditionalParams = append(*conditionalParams, paramValue(key, actionArgument{
				valueType: Integer,
				value:     "1",
			}, Integer, Text))
		} else {
			*conditionalParams = append(*conditionalParams, paramValue(key, actionArgument{
				valueType: Integer,
				value:     "0",
			}, Integer, Text))
		}
	case Variable:
		conditionalParameterVariable(conditionalParams, value)
	}
}

func conditionalParameterVariable(conditionalParams *[]plistData, value any) {
	var stringValue = value.(string)
	var variable = variables[stringValue]
	switch variable.valueType {
	case Integer:
		*conditionalParams = append(*conditionalParams, variablePlistValue("WFNumberValue", stringValue, uuids[stringValue]))
	default:
		*conditionalParams = append(*conditionalParams, attachmentValues("WFConditionalActionString", fmt.Sprintf("{%s}", stringValue), Text))
	}
}

func makeVariableValue(token *token, varUUID *string) {
	var UUID = plistData{
		key:      "UUID",
		dataType: Text,
		value:    *varUUID,
	}
	var outputName = plistData{
		key:      "CustomOutputName",
		dataType: Text,
		value:    token.ident,
	}

	if token.typeof != Var &&
		(token.typeof == AddTo || token.typeof == SubFrom || token.typeof == MultiplyBy || token.typeof == DivideBy) &&
		token.valueType != Arr &&
		variables[token.ident].valueType != Arr {
		variableValueModifier(token, &outputName, &UUID)
		return
	}

	switch token.valueType {
	case Integer:
		appendPlist(makeStdAction("number", []plistData{
			outputName,
			UUID,
			{
				key:      "WFNumberActionNumber",
				dataType: Text,
				value:    token.value,
			},
		}))
	case Bool:
		var boolValue string
		if token.value == true {
			boolValue = "1"
		} else {
			boolValue = "0"
		}
		appendPlist(makeStdAction("number", []plistData{
			outputName,
			UUID,
			{
				key:      "WFNumberActionNumber",
				dataType: Text,
				value:    boolValue,
			},
		}))
	case String:
		appendPlist(makeStdAction("gettext", []plistData{
			outputName,
			UUID,
			attachmentValues("WFTextActionText", token.value.(string), Text),
		}))
	case Expression:
		var expression = token.value.(string)
		var expressionParts []string
		if tokensOccur(&expression, Plus, Minus, Multiply, Divide, Modulus) {
			var operandOne string
			var operandTwo string
			var operation string
			var expressionParts = strings.Split(expression, " ")
			if len(expressionParts) == 3 {
				operation = ""
			}
			if operation != "" {
				expressionParts = strings.Split(expression, operation)
				operandOne = strings.Trim(expressionParts[0], " ")
				operandTwo = strings.Trim(expressionParts[1], " ")
				wrapVariableReference(&operandOne)
				wrapVariableReference(&operandTwo)

				appendPlist(makeStdAction("math", []plistData{
					outputName,
					UUID,
					attachmentValues("WFScientificMathOperation", operation, Text),
					attachmentValues("WFInput", operandOne, Number),
					attachmentValues("WFMathOperand", operandTwo, Number),
				}))
				break
			}
		}
		var formattedExpression []string
		expressionParts = strings.Split(expression, " ")
		for _, part := range expressionParts {
			var p = part
			wrapVariableReference(&p)
			formattedExpression = append(formattedExpression, p)
		}
		expression = strings.Join(formattedExpression, " ")
		appendPlist(makeStdAction("calculateexpression", []plistData{
			outputName,
			attachmentValues("Input", expression, Text),
			UUID,
		}))
	case Action:
		currentAction = token.value.(action).ident
		plistAction(token.value.(action).args, outputName, UUID)
	case Dict:
		appendPlist(makeStdAction("dictionary", []plistData{
			outputName,
			UUID,
			{
				key:      "WFItems",
				dataType: Dictionary,
				value: []plistData{
					{
						key:      "Value",
						dataType: Dictionary,
						value: []plistData{
							{
								key:      "WFDictionaryFieldValueItems",
								dataType: Array,
								value:    makeDictionary(token.value),
							},
						},
					},
					{
						key:      "WFSerializationType",
						dataType: Text,
						value:    "WFDictionaryFieldValue",
					},
				},
			},
		}))
	}
}

func variableValueModifier(token *token, outputName *plistData, UUID *plistData) {
	var valueType = token.valueType
	if valueType == Variable {
		valueType = variables[token.value.(string)].valueType
	}
	switch valueType {
	case Integer:
		var operation string
		switch token.typeof {
		case AddTo:
			operation = "+"
		case SubFrom:
			operation = "-"
		case MultiplyBy:
			operation = "×"
		case DivideBy:
			operation = "÷"
		}
		var tokenType = convertTypeToken(token.valueType)
		appendPlist(makeStdAction("math", []plistData{
			*outputName,
			*UUID,
			paramValue("WFMathOperand", actionArgument{
				valueType: token.valueType,
				value:     token.value,
			}, token.valueType, tokenType),
			paramValue("WFInput", actionArgument{
				valueType: Var,
				value:     token.ident,
			}, token.valueType, tokenType),
			{
				key:      "WFMathOperation",
				dataType: Text,
				value:    operation,
			},
		}))
	case String:
		var varInput = token.value.(string)
		wrapVariableReference(&varInput)
		appendPlist(makeStdAction("gettext", []plistData{
			*outputName,
			*UUID,
			paramValue("WFTextActionText", actionArgument{
				valueType: String,
				value:     fmt.Sprintf("{%s}%s", token.ident, varInput),
			}, token.valueType, convertTypeToken(valueType)),
		}))
	}
}

func variableInput(key string, name string) plistData {
	if uuid, found := uuids[name]; found {
		return inputValue(key, name, uuid)
	}
	return inputValue(key, name, "")
}

func inputValue(key string, name string, varUUID string) plistData {
	var value []plistData
	if varUUID != "" {
		var variable = variables[name]
		if variable.valueType != Variable {
			value = []plistData{
				{
					key:      "OutputName",
					dataType: Text,
					value:    name,
				},
				{
					key:      "OutputUUID",
					dataType: Text,
					value:    varUUID,
				},
				{
					key:      "Type",
					dataType: Text,
					value:    "ActionOutput",
				},
			}
		} else {
			value = []plistData{
				{
					key:      "VariableName",
					dataType: Text,
					value:    name,
				},
				{
					key:      "Type",
					dataType: Text,
					value:    "Variable",
				},
			}
		}
	} else if global, found := globals[name]; found {
		value = []plistData{
			{
				key:      "Type",
				dataType: Text,
				value:    global.variableType,
			},
		}
	} else {
		value = []plistData{
			{
				key:      "OutputName",
				dataType: Text,
				value:    name,
			},
			{
				key:      "Type",
				dataType: Text,
				value:    "ActionOutput",
			},
		}
	}
	return plistData{
		key:      key,
		dataType: Dictionary,
		value: []plistData{
			{
				key:      "Value",
				dataType: Dictionary,
				value:    value,
			},
			{
				key:      "WFSerializationType",
				dataType: Text,
				value:    "WFTextTokenAttachment",
			},
		},
	}
}

func variablePlistValue(key string, identifier string, ident string) plistData {
	var variable variableValue
	var getAs string
	var coerce string
	var aggrandizements []plistData
	if global, found := globals[identifier]; found {
		variable = global
		isInputVariable(identifier)
		identifier = variable.value.(string)
	} else if v, found := variables[identifier]; found {
		variable = v
	}
	if v, found := variables[ident]; found {
		getAs = v.getAs
		coerce = v.coerce
		if getAs != "" {
			aggrandizements = append(aggrandizements, plistData{
				dataType: Dictionary,
				value: []plistData{
					{
						key:      "PropertyUserInfo",
						dataType: Number,
						value:    0,
					},
					{
						key:      "Type",
						dataType: Text,
						value:    "WFPropertyVariableAggrandizement",
					},
					{
						key:      "PropertyName",
						dataType: Text,
						value:    getAs,
					},
				},
			})
		}
		if coerce != "" {
			makeContentItems()
			if contentItem, found := contentItems[coerce]; found {
				aggrandizements = append(aggrandizements, plistData{
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "Type",
							dataType: Text,
							value:    "WFCoercionVariableAggrandizement",
						},
						{
							key:      "CoercionItemClass",
							dataType: Text,
							value:    contentItem,
						},
					},
				})
			}
		}
	}
	var varType = "Variable"
	if variable.variableType != "" {
		varType = variable.variableType
	}
	var varValue []plistData
	if variable.constant {
		var varUUID = uuids[identifier]
		varValue = []plistData{
			{
				key:      "OutputName",
				dataType: Text,
				value:    identifier,
			},
			{
				key:      "OutputUUID",
				dataType: Text,
				value:    varUUID,
			},
			{
				key:      "Type",
				dataType: Text,
				value:    "ActionOutput",
			},
		}
	} else {
		varValue = []plistData{
			{
				key:      "VariableName",
				dataType: Text,
				value:    identifier,
			},
			{
				key:      "Type",
				dataType: Text,
				value:    varType,
			},
		}
	}
	if len(aggrandizements) > 0 {
		varValue = append(varValue, plistData{
			key:      "Aggrandizements",
			dataType: Array,
			value:    aggrandizements,
		})
	}
	return plistData{
		key:      key,
		dataType: Dictionary,
		value: []plistData{
			{
				key:      "Value",
				dataType: Dictionary,
				value:    varValue,
			},
			{
				key:      "WFSerializationType",
				dataType: Text,
				value:    "WFTextTokenAttachment",
			},
		},
	}
}

type inlineVar struct {
	identifier string
	col        int
	getAs      string
	coerce     string
}

type attachmentVariable struct {
	identifier string
	getAs      string
	coerce     string
}

var varPositions []plistData
var inlineVars []inlineVar
var varIndex []attachmentVariable

func attachmentValues(key string, str string, outputType plistDataType) plistData {
	if !strings.ContainsAny(str, "{}") {
		return plistData{
			key:      key,
			dataType: outputType,
			value:    str,
		}
	}
	varPositions = []plistData{}
	inlineVars = []inlineVar{}
	varIndex = []attachmentVariable{}
	var noVarString = collectInlineVariables(&str)
	for _, stringVar := range inlineVars {
		var storedVar variableValue
		if g, global := globals[stringVar.identifier]; global {
			storedVar = g
			isInputVariable(stringVar.identifier)
			stringVar.identifier = g.value.(string)
		} else if v, found := variables[stringVar.identifier]; found {
			storedVar = v
		} else {
			exit(fmt.Sprintf("Undefined reference '%s'", stringVar.identifier))
		}
		var variable = variables[stringVar.identifier]
		var varUUID = uuids[stringVar.identifier]
		var varValue []plistData
		var varType = "Variable"
		var aggr []plistData
		if storedVar.variableType != "" {
			varType = storedVar.variableType
		}
		if !variable.constant {
			varValue = []plistData{
				{
					key:      "VariableName",
					dataType: Text,
					value:    stringVar.identifier,
				},
				{
					key:      "Type",
					dataType: Text,
					value:    varType,
				},
			}
		} else {
			varValue = []plistData{
				{
					key:      "OutputName",
					dataType: Text,
					value:    stringVar.identifier,
				},
				{
					key:      "OutputUUID",
					dataType: Text,
					value:    varUUID,
				},
				{
					key:      "Type",
					dataType: Text,
					value:    "ActionOutput",
				},
			}
		}
		if stringVar.getAs != "" {
			aggr = append(aggr, plistData{
				dataType: Dictionary,
				value: []plistData{
					{
						key:      "Type",
						dataType: Text,
						value:    "WFPropertyVariableAggrandizement",
					},
					{
						key:      "PropertyName",
						dataType: Text,
						value:    stringVar.getAs,
					},
				},
			})
		}
		if stringVar.coerce != "" {
			makeContentItems()
			if contentItem, found := contentItems[stringVar.coerce]; found {
				aggr = append(aggr, plistData{
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "Type",
							dataType: Text,
							value:    "WFCoercionVariableAggrandizement",
						},
						{
							key:      "CoercionItemClass",
							dataType: Text,
							value:    contentItem,
						},
					},
				})
			} else {
				var list = makeKeyList("Available content item types:", contentItems)
				parserError(fmt.Sprintf("Invalid content item for type coerce '%s'\n\n%s\n", stringVar.coerce, list))
			}
		}
		if stringVar.getAs != "" || stringVar.coerce != "" {
			varValue = append(varValue, plistData{
				key:      "Aggrandizements",
				dataType: Array,
				value:    aggr,
			})
		}
		varPositions = append(varPositions, plistData{
			key:      fmt.Sprintf("{%d, 1}", stringVar.col),
			dataType: Dictionary,
			value:    varValue,
		})
	}
	return plistData{
		key:      key,
		dataType: Dictionary,
		value: []plistData{
			{
				key:      "Value",
				dataType: Dictionary,
				value: []plistData{
					{
						key:      "attachmentsByRange",
						dataType: Dictionary,
						value:    varPositions,
					},
					{
						key:      "string",
						dataType: Text,
						value:    noVarString,
					},
				},
			},
			{
				key:      "WFSerializationType",
				dataType: Text,
				value:    "WFTextTokenString",
			},
		},
	}
}

// mapInlineVars finds occurrences of ObjectReplaceChar and adds them to inlineVars to map the inline variables in noVarString.
func mapInlineVars(noVarString *string) {
	var variableIdx int
	var i = 0
	for _, c := range *noVarString {
		if c == ObjectReplaceChar {
			inlineVars = append(inlineVars, inlineVar{
				identifier: varIndex[variableIdx].identifier,
				col:        i,
				getAs:      varIndex[variableIdx].getAs,
				coerce:     varIndex[variableIdx].coerce,
			})
			variableIdx++
		}
		i++
	}
}

// collectInlineVariables collects inline variables from `str` and adds them to a slice of attachmentVariable.
// It then replaces all instances of inline variables in `str` with ObjectReplaceChar.
func collectInlineVariables(str *string) (noVarString string) {
	var collectVarRegex = regexp.MustCompile(`\{(.*?)(?:\[(.*?)])?(?:\.(.*?))?}`)
	var matches = collectVarRegex.FindAllStringSubmatch(*str, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			var attachmentVar attachmentVariable
			if len(match) < 2 {
				continue
			}
			attachmentVar.identifier = match[1]
			if !validReference(attachmentVar.identifier) {
				parserError(fmt.Sprintf("Inline var '%s' does not exist!", attachmentVar.identifier))
			}
			if len(match[2]) > 0 {
				attachmentVar.getAs = match[2]
			}
			if len(match[3]) > 0 {
				attachmentVar.coerce = match[3]
			}
			varIndex = append(varIndex, attachmentVar)
		}

		var replaceVarRegex = regexp.MustCompile(`(\{.*?})`)
		noVarString = replaceVarRegex.ReplaceAllString(*str, ObjectReplaceCharStr)
	}

	mapInlineVars(&noVarString)

	return
}

func convertTypeToken(tokenType tokenType) plistDataType {
	switch tokenType {
	case String:
		return Text
	case Integer:
		return Number
	case Arr:
		return Array
	case Dict:
		return Dictionary
	case Bool:
		return Boolean
	default:
		return ""
	}
}

func argumentValue(key string, args []actionArgument, idx int) plistData {
	var actionArg parameterDefinition
	if len(actions[currentAction].parameters) <= idx {
		// First parameter is likely infinite
		actionArg = actions[currentAction].parameters[0]
	} else {
		actionArg = actions[currentAction].parameters[idx]
	}
	var arg actionArgument
	if len(args) <= idx {
		if actionArg.optional || actionArg.defaultValue != nil {
			return plistData{}
		}
	} else {
		arg = args[idx]
	}
	return paramValue(key, arg, actionArg.validType, convertTypeToken(actionArg.validType))
}

func paramValue(key string, arg actionArgument, handleAs tokenType, outputType plistDataType) plistData {
	if arg.valueType == Nil || arg.value == nil {
		return plistData{}
	}
	switch arg.valueType {
	case Variable:
		if handleAs == String {
			return attachmentValues(key, fmt.Sprintf("{%s}", arg.value), Text)
		}
		return variablePlistValue(key, arg.value.(string), "")
	case Bool:
		return plistData{
			key:      key,
			dataType: Boolean,
			value:    arg.value,
		}
	default:
		return attachmentValues(key, arg.value.(string), outputType)
	}
}

func wrapVariableReference(s *string) {
	if validReference(*s) {
		*s = fmt.Sprintf("{%s}", *s)
	}
}

// isInputVariable checks if identifier is the ShortcutInput global to set the global boolean in the final plist.
func isInputVariable(identifier string) {
	hasShortcutInputVariables = identifier == "ShortcutInput"
}

const (
	startStatement = 0
	statementPart  = 1
	endStatement   = 2
)

const (
	stringType string = "string"
	intType    string = "float64"
	arrayType  string = "[]interface {}"
	dictType   string = "map[string]interface {}"
	boolType   string = "bool"
)

// makeDictionary creates a Shortcut dictionary value.
func makeDictionary(value interface{}) (dictItems []plistData) {
	for key, item := range value.(map[string]interface{}) {
		dictItems = append(dictItems, dictionaryValue(key, item))
	}
	return
}

// dictionaryValue creates an inner dictionary value.
func dictionaryValue(key string, value any) plistData {
	var itemType dictDataType
	var serializedType string
	var wfValue = plistData{
		key:      "Value",
		dataType: Dictionary,
		value: []plistData{
			{
				key:      "string",
				dataType: Text,
				value:    value,
			},
		},
	}
	switch reflect.TypeOf(value).String() {
	case stringType:
		if strings.ContainsAny(value.(string), "{}") {
			wfValue = paramValue("Value", actionArgument{
				valueType: String,
				value:     value,
			}, String, Text)
			if reflect.TypeOf(wfValue.value).String() == "[]main.plistData" {
				for _, val := range wfValue.value.([]plistData) {
					wfValue = val
					break
				}
			}
		}
		itemType = itemTypeText
		serializedType = "WFTextTokenString"
	case intType:
		itemType = itemTypeNumber
		serializedType = "WFTextTokenString"
	case arrayType:
		itemType = itemTypeArray
		serializedType = "WFArrayParameterState"
		var arrayValue []plistData
		for _, item := range value.([]interface{}) {
			arrayValue = append(arrayValue, dictionaryValue("", item))
		}
		wfValue = plistData{
			key:      "Value",
			dataType: Array,
			value:    arrayValue,
		}
	case dictType:
		itemType = itemTypeDict
		serializedType = "WFDictionaryFieldValue"
		wfValue = plistData{
			key:      "Value",
			dataType: Dictionary,
			value: []plistData{
				{
					key:      "Value",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "WFDictionaryFieldValueItems",
							dataType: Array,
							value:    makeDictionary(value),
						},
					},
				},
				{
					key:      "WFSerializationType",
					dataType: Text,
					value:    "WFDictionaryFieldValue",
				},
			},
		}
	case boolType:
		itemType = itemTypeBool
		serializedType = "WFNumberSubstitutableState"
		wfValue = plistData{
			key:      "Value",
			dataType: Boolean,
			value:    value,
		}
	}
	return dictionaryPlistValue(key, itemType, serializedType, wfValue)
}

func dictionaryPlistValue(key string, itemType dictDataType, serializedType string, wfValue plistData) plistData {
	var valueData = []plistData{
		{
			key:      "WFItemType",
			dataType: Number,
			value:    string(itemType),
		},
		{
			key:      "WFValue",
			dataType: Dictionary,
			value: []plistData{
				wfValue,
				{
					key:      "WFSerializationType",
					dataType: Text,
					value:    serializedType,
				},
			},
		},
	}
	if key != "" {
		var WFKey = plistData{
			key:      "Value",
			dataType: Dictionary,
			value: []plistData{
				{
					key:      "string",
					dataType: Text,
					value:    key,
				},
			},
		}
		if strings.ContainsAny(key, "{}") {
			WFKey = paramValue("Value", actionArgument{
				valueType: String,
				value:     key,
			}, String, Text)
			if reflect.TypeOf(WFKey.value).String() == "[]main.plistData" {
				for _, val := range WFKey.value.([]plistData) {
					WFKey = val
					break
				}
			}
		}
		valueData = append(valueData, plistData{
			key:      "WFKey",
			dataType: Dictionary,
			value: []plistData{
				WFKey,
				{
					key:      "WFSerializationType",
					dataType: Text,
					value:    "WFTextTokenString",
				},
			},
		})
	}
	return plistData{
		dataType: Dictionary,
		value:    valueData,
	}
}
