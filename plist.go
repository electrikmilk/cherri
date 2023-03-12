/*
 * Copyright (c) 2023 Brandon Jordan
 */

package main

import (
	"fmt"
	"html"
	"reflect"
	"strings"
)

var shortcutActions []string
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

func plistKeyValue(key string, dataType plistDataType, value any) (pair string) {
	if key != "" {
		pair = "<key>" + key + "</key>\n"
	}
	switch dataType {
	case Boolean:
		if value == true {
			pair += "<true/>\n"
		} else {
			pair += "<false/>\n"
		}
	case Array:
		if len(value.([]string)) == 0 {
			pair += "<array/>\n"
			break
		}
		pair += "<array>\n"
		for _, val := range value.([]string) {
			pair += html.EscapeString(val)
		}
		pair += "</array>\n"
	case Dictionary:
		pair += "<dict>\n"
		var empty = plistData{}
		for _, data := range value.([]plistData) {
			if data != empty {
				pair += plistKeyValue(data.key, data.dataType, data.value)
			}
		}
		pair += "</dict>\n"
	default:
		if reflect.TypeOf(value).String() == "string" {
			value = html.EscapeString(value.(string))
		}
		pair += fmt.Sprintf("<%s>%v</%s>\n", dataType, value, dataType)
	}
	return
}

func plistValue(dataType plistDataType, value any) string {
	return plistKeyValue("", dataType, value)
}

func condParam(key string, conditionalParams *[]plistData, typeOf *tokenType, value any) {
	if key == "" {
		switch *typeOf {
		case String:
			key = "WFConditionalActionString"
		case Integer, Bool:
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
		var variable = variables[strings.ToLower(value.(string))]
		switch variable.valueType {
		case Integer:
			*conditionalParams = append(*conditionalParams, variablePlistValue("WFNumberValue", value.(string), uuids[value.(string)]))
		case String:
			*conditionalParams = append(*conditionalParams, attachmentValues("WFConditionalActionString", fmt.Sprintf("{%s}", value.(string)), uuids[value.(string)], Text))
		default:
			var realVar = realVariableValue(value.(string), String)
			condParam(key, conditionalParams, &realVar.valueType, realVar.value)
		}
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
	switch token.valueType {
	case Integer:
		shortcutActions = append(shortcutActions, makeStdAction("number", []plistData{
			UUID,
			outputName,
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
		shortcutActions = append(shortcutActions, makeStdAction("number", []plistData{
			UUID,
			outputName,
			{
				key:      "WFNumberActionNumber",
				dataType: Text,
				value:    boolValue,
			},
		}))
	case String:
		shortcutActions = append(shortcutActions, makeStdAction("gettext", []plistData{
			UUID,
			outputName,
			attachmentValues("WFTextActionText", token.value.(string), "", Text),
		}))
	case Expression:
		shortcutActions = append(shortcutActions, makeStdAction("calculateexpression", []plistData{
			UUID,
			outputName,
			{
				key:      "Input",
				dataType: Text,
				value:    token.value,
			},
		}))
	case Action:
		currentAction = token.value.(action).ident
		callAction(token.value.(action).args, outputName, UUID)
	case Dict:
		shortcutActions = append(shortcutActions, makeStdAction("dictionary", []plistData{
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

func variableInput(key string, name string) plistData {
	if _, found := uuids[name]; found {
		return inputValue(key, name, uuids[name])
	}
	return inputValue(key, name, "")
}

func inputValue(key string, name string, varUUID string) plistData {
	var value []plistData
	if varUUID != "" {
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
	} else if _, found := globals[name]; found {
		value = []plistData{
			{
				key:      "Type",
				dataType: Text,
				value:    globals[name].variableType,
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

func variablePlistValue(key string, varName string, ident string) plistData {
	var variable variableValue
	var getAs string
	var coerce string
	var aggrandizements []string
	var lowerVarName = strings.ToLower(varName)
	var lowerIdent = strings.ToLower(ident)
	if _, global := globals[varName]; global {
		variable = globals[varName]
		hasInputVariables(varName)
		varName = variable.value.(string)
	} else if _, found := variables[lowerVarName]; found {
		variable = variables[lowerVarName]
	}
	if _, found := variables[lowerIdent]; found {
		getAs = variables[lowerIdent].getAs
		coerce = variables[lowerIdent].coerce
		if getAs != "" {
			aggrandizements = append(aggrandizements, plistValue(Dictionary, []plistData{
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
			}))
		}
		if coerce != "" {
			makeContentItems()
			if _, found := contentItems[coerce]; found {
				aggrandizements = append(aggrandizements, plistValue(Dictionary, []plistData{
					{
						key:      "Type",
						dataType: Text,
						value:    "WFCoercionVariableAggrandizement",
					},
					{
						key:      "CoercionItemClass",
						dataType: Text,
						value:    contentItems[coerce],
					},
				}))
			}
		}
	}
	var varType = "Variable"
	if variable.variableType != "" {
		varType = variable.variableType
	}
	var varValue = []plistData{
		{
			key:      "VariableName",
			dataType: Text,
			value:    varName,
		},
		{
			key:      "Type",
			dataType: Text,
			value:    varType,
		},
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

type stringVar struct {
	varName string
	col     int
	getAs   string
	coerce  string
}

type attachmentVariable struct {
	varName string
	getAs   string
	coerce  string
}

var varPositions []plistData

func attachmentValues(key string, variable string, varUUID string, outputType plistDataType) plistData {
	if !strings.Contains(variable, "{") {
		return plistData{
			key:      key,
			dataType: outputType,
			value:    variable,
		}
	}
	varPositions = []plistData{}
	var stringVars []stringVar
	var varIndex = make(map[int]attachmentVariable)
	var variableChars = strings.Split(variable, "")
	var currentVariable string
	var collectingVariable bool
	var collectingGetAs bool
	var collectingCoerce bool
	var varNum int
	var noVarString = variable
	var getAs string
	var coerce string
	for _, chr := range variableChars {
		if collectingVariable {
			switch {
			case collectingGetAs:
				if chr == "]" {
					collectingGetAs = false
					continue
				}
				getAs += chr
			case collectingCoerce:
				if chr == ")" {
					collectingCoerce = false
					continue
				}
				coerce += chr
			default:
				if chr == "}" {
					varIndex[varNum] = attachmentVariable{
						varName: currentVariable,
						getAs:   getAs,
						coerce:  coerce,
					}
					var varName = currentVariable
					if getAs != "" {
						varName = currentVariable + "[" + getAs + "]"
					}
					if coerce != "" {
						varName = currentVariable + "(" + coerce + ")"
					}
					// Replace with OBJECT REPLACEMENT character
					noVarString = strings.Replace(
						noVarString,
						"{"+varName+"}",
						"\uFFFC",
						1)
					currentVariable = ""
					getAs = ""
					coerce = ""
					collectingVariable = false
					varNum++
					continue
				}
				if chr == "[" {
					collectingGetAs = true
					continue
				}
				if chr == "(" {
					collectingCoerce = true
					continue
				}
				currentVariable += chr
			}
		} else if chr == "{" {
			collectingVariable = true
		}
	}
	var variableIdx int
	var noVarChars = strings.Split(noVarString, "")
	for c, s := range noVarChars {
		if s == "\uFFFC" {
			stringVars = append(stringVars, stringVar{
				varName: varIndex[variableIdx].varName,
				col:     c,
				getAs:   varIndex[variableIdx].getAs,
				coerce:  varIndex[variableIdx].coerce,
			})
			variableIdx++
		}
	}
	for _, stringVar := range stringVars {
		var storedVar variableValue
		if _, global := globals[stringVar.varName]; global {
			storedVar = globals[stringVar.varName]
			hasInputVariables(stringVar.varName)
			stringVar.varName = storedVar.value.(string)
		} else if _, found := variables[strings.ToLower(stringVar.varName)]; found {
			storedVar = variables[stringVar.varName]
		} else {
			exit(fmt.Sprintf("Variable '%s' does not exist!", stringVar.varName))
		}
		var varValue []plistData
		var varType = "Variable"
		var aggr []string
		if storedVar.variableType != "" {
			varType = storedVar.variableType
		}
		if varUUID == "" {
			varValue = []plistData{
				{
					key:      "VariableName",
					dataType: Text,
					value:    stringVar.varName,
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
					key:      "OutputUUID",
					dataType: Text,
					value:    varUUID,
				},
				{
					key:      "Type",
					dataType: Text,
					value:    "ActionOutput",
				},
				{
					key:      "OutputName",
					dataType: Text,
					value:    stringVar.varName,
				},
			}
		}
		if stringVar.getAs != "" {
			aggr = append(aggr, plistValue(Dictionary, []plistData{
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
			}))
		}
		if stringVar.coerce != "" {
			makeContentItems()
			if _, found := contentItems[stringVar.coerce]; found {
				aggr = append(aggr, plistValue(Dictionary, []plistData{
					{
						key:      "Type",
						dataType: Text,
						value:    "WFCoercionVariableAggrandizement",
					},
					{
						key:      "CoercionItemClass",
						dataType: Text,
						value:    contentItems[stringVar.coerce],
					},
				}))
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
						key:      "string",
						dataType: Text,
						value:    noVarString,
					},
					{
						key:      "attachmentsByRange",
						dataType: Dictionary,
						value:    varPositions,
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

func argumentValue(key string, args []actionArgument, idx int) plistData {
	var actionArg = actions[currentAction].parameters[idx]
	var arg actionArgument
	if len(args) <= idx {
		if actionArg.optional {
			return plistData{}
		}
		if actionArg.defaultValue.value != nil {
			arg = actionArg.defaultValue
		}
	} else {
		arg = args[idx]
	}
	var plistType plistDataType
	switch actionArg.validType {
	case String:
		plistType = Text
	case Integer:
		plistType = Number
	case Arr:
		plistType = Array
	case Dict:
		plistType = Dictionary
	case Bool:
		plistType = Boolean
	}
	return paramValue(key, arg, actionArg.validType, plistType)
}

func paramValue(key string, arg actionArgument, handleAs tokenType, outputType plistDataType) plistData {
	switch arg.valueType {
	case Variable:
		if handleAs == String {
			return attachmentValues(key, fmt.Sprintf("{%s}", arg.value), "", Text)
		}

		return variablePlistValue(key, arg.value.(string), "")
	case Bool:
		return plistData{
			key:      key,
			dataType: Boolean,
			value:    arg.value,
		}
	default:
		return attachmentValues(key, arg.value.(string), "", outputType)
	}
}

func hasInputVariables(varName string) {
	if varName == "ShortcutInput" {
		hasShortcutInputVariables = true
	}
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

func makeDictionary(value interface{}) (dictItems []string) {
	for key, item := range value.(map[string]interface{}) {
		dictItems = append(dictItems, dictionaryValue(key, item))
	}
	return
}

func dictionaryValue(key string, value any) string {
	var itemType dictDataType
	var serializedType string
	var WFValue = plistData{
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
		itemType = itemTypeText
		serializedType = "WFTextTokenString"
	case intType:
		itemType = itemTypeNumber
		serializedType = "WFTextTokenString"
	case arrayType:
		itemType = itemTypeArray
		serializedType = "WFArrayParameterState"
		var arrayValue []string
		for _, item := range value.([]interface{}) {
			arrayValue = append(arrayValue, dictionaryValue("", item))
		}
		WFValue = plistData{
			key:      "Value",
			dataType: Array,
			value:    arrayValue,
		}
	case dictType:
		itemType = itemTypeDict
		serializedType = "WFDictionaryFieldValue"
		WFValue = plistData{
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
		WFValue = plistData{
			key:      "Value",
			dataType: Boolean,
			value:    value,
		}
	}
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
				WFValue,
				{
					key:      "WFSerializationType",
					dataType: Text,
					value:    serializedType,
				},
			},
		},
	}
	if key != "" {
		valueData = append(valueData, plistData{
			key:      "WFKey",
			dataType: Dictionary,
			value: []plistData{
				{
					key:      "Value",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "string",
							dataType: Text,
							value:    key,
						},
					},
				},
				{
					key:      "WFSerializationType",
					dataType: Text,
					value:    "WFTextTokenString",
				},
			},
		})
	}
	return plistValue(Dictionary, valueData)
}
