/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type plistDataType string

const Text plistDataType = "string"
const Number plistDataType = "integer"
const Dictionary plistDataType = "dictionary"
const Array plistDataType = "array"
const Boolean plistDataType = "boolean"

type dictDataType string

const itemTypeText dictDataType = "0"
const itemTypeNumber dictDataType = "3"
const itemTypeArray dictDataType = "2"
const itemTypeDict dictDataType = "1"
const itemTypeBool dictDataType = "4"

type plistData struct {
	key      string
	dataType plistDataType
	value    any
}

var shortcutActions []string

var uuids map[string]string

var hasShortcutInputVariables = false

type noInputParams struct {
	name   string
	params []plistData
}

var noInput noInputParams

func makeAction(ident string, paramsDict []plistData) (action string) {
	action = plistDict("", []plistData{
		{
			key:      "WFWorkflowActionIdentifier",
			dataType: Text,
			value:    "is.workflow.actions." + ident,
		},
		{
			key:      "WFWorkflowActionParameters",
			dataType: Dictionary,
			value:    paramsDict,
		},
	})
	return
}

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
	default:
		pair += fmt.Sprintf("<%s>%v</%s>\n", dataType, value, dataType)
	}
	return
}

func plistValue(dataType plistDataType, value any) string {
	return plistKeyValue("", dataType, value)
}

func plistArray(key string, values []string) (pair string) {
	if key != "" {
		pair = "<key>" + key + "</key>\n"
	}
	if len(values) == 0 {
		pair += "<array/>\n"
		return
	} else {
		pair += "<array>\n"
		for _, val := range values {
			pair += val
		}
		pair += "</array>\n"
		return
	}
}

func plistDict(key string, values []plistData) (pair string) {
	if key != "" {
		pair = "<key>" + key + "</key>\n"
	}
	pair += "<dict>\n"
	for _, data := range values {
		switch data.dataType {
		case Dictionary:
			pair += plistDict(data.key, data.value.([]plistData))
		case Array:
			pair += plistArray(data.key, data.value.([]string))
		default:
			pair += plistKeyValue(data.key, data.dataType, data.value)
		}
	}
	pair += "</dict>\n"
	return
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
		var realVar = realVariableValue(value.(string))
		condParam(key, conditionalParams, &realVar.valueType, realVar.value)
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
	switch {
	case token.valueType == Integer:
		shortcutActions = append(shortcutActions, makeAction("number", []plistData{
			UUID,
			outputName,
			{
				key:      "WFNumberActionNumber",
				dataType: Text,
				value:    token.value,
			},
		}))
	case token.valueType == Bool:
		var boolValue string
		if token.value == true {
			boolValue = "1"
		} else {
			boolValue = "0"
		}
		shortcutActions = append(shortcutActions, makeAction("number", []plistData{
			UUID,
			outputName,
			{
				key:      "WFNumberActionNumber",
				dataType: Text,
				value:    boolValue,
			},
		}))
	case token.valueType == String:
		shortcutActions = append(shortcutActions, makeAction("gettext", []plistData{
			UUID,
			outputName,
			attachmentValues("WFTextActionText", token.value.(string), "", Text),
		}))
	case token.valueType == Expression:
		shortcutActions = append(shortcutActions, makeAction("calculateexpression", []plistData{
			UUID,
			outputName,
			{
				key:      "Input",
				dataType: Text,
				value:    token.value,
			},
		}))
	case token.valueType == Action:
		currentAction = token.value.(action).ident
		callAction(token.value.(action).args, outputName, UUID)
	case token.valueType == Dict:
		shortcutActions = append(shortcutActions, makeAction("dictionary", []plistData{
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
								value:    makeDict(token.value),
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
	return plistData{
		key:      key,
		dataType: Dictionary,
		value: []plistData{
			{
				key:      "Value",
				dataType: Dictionary,
				value: []plistData{
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
				},
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
			aggrandizements = append(aggrandizements, plistDict("", []plistData{
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
			if len(contentItems) == 0 {
				makeContentItems()
			}
			if _, found := contentItems[coerce]; found {
				aggrandizements = append(aggrandizements, plistDict("", []plistData{
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
			if collectingGetAs {
				if chr == "]" {
					collectingGetAs = false
					continue
				}
				getAs += chr
			} else if collectingCoerce {
				if chr == ")" {
					collectingCoerce = false
					continue
				}
				coerce += chr
			} else {
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
					// Replace with OBJECT REPLACEMENT CHARACTER
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
			fmt.Printf("\n\n\033[31mVariable '%s' does not exist!\033[0m\n", stringVar.varName)
			os.Exit(1)
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
			aggr = append(aggr, plistDict("", []plistData{
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
			if _, found := contentItems[stringVar.coerce]; found {
				aggr = append(aggr, plistDict("", []plistData{
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
				var list = makeKeyList("Available content items:", contentItems)
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
	var actionArg = actions[currentAction].args[idx]
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
	return paramValue(key, args[idx], actionArg.validType, plistType)
}

func paramValue(key string, arg actionArgument, handleAs tokenType, outputType plistDataType) plistData {
	switch arg.valueType {
	case Variable:
		if handleAs == String {
			return attachmentValues(key, fmt.Sprintf("{%s}", arg.value), "", Text)
		} else {
			return variablePlistValue(key, arg.value.(string), "")
		}
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
	stringType string = "string"
	intType    string = "float64"
	arrayType  string = "[]interface {}"
	dictType   string = "map[string]interface {}"
	boolType   string = "bool"
)

func makeDict(value interface{}) (dictItems []string) {
	for key, item := range value.(map[string]interface{}) {
		var itemType dictDataType
		var serializedType string
		var WFValue = plistData{
			key:      "Value",
			dataType: Dictionary,
			value: []plistData{
				{
					key:      "string",
					dataType: Text,
					value:    item,
				},
			},
		}
		switch reflect.TypeOf(item).String() {
		case stringType:
			itemType = itemTypeText
			serializedType = "WFTextTokenString"
		case intType:
			itemType = itemTypeNumber
			serializedType = "WFTextTokenString"
		case arrayType:
			itemType = itemTypeArray
			serializedType = "WFArrayParameterState"
			WFValue = plistData{
				key:      "Value",
				dataType: Array,
				value:    dictArray(item),
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
								value:    makeDict(item),
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
				value:    item,
			}
		}
		dictItems = append(dictItems, plistDict("", []plistData{
			{
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
			},
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
		}))
	}
	return
}

func dictArray(value interface{}) (array []string) {
	for _, item := range value.([]interface{}) {
		switch reflect.TypeOf(item).String() {
		case stringType:
			array = append(array, plistValue(Text, item))
		case intType:
			array = append(array, plistDict("", []plistData{
				{
					key:      "WFItemType",
					dataType: Number,
					value:    itemTypeNumber,
				},
				{
					key:      "WFValue",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "Value",
							dataType: Dictionary,
							value: []plistData{
								{
									key:      "string",
									dataType: Text,
									value:    item,
								},
							},
						},
						{
							key:      "WFSerializationType",
							dataType: Text,
							value:    "WFTextTokenString",
						},
					},
				},
			}))
		case boolType:
			array = append(array, plistDict("", []plistData{
				{
					key:      "WFItemType",
					dataType: Number,
					value:    itemTypeBool,
				},
				{
					key:      "WFValue",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "Value",
							dataType: Boolean,
							value:    item,
						},
						{
							key:      "WFSerializationType",
							dataType: Text,
							value:    "WFNumberSubstitutableState",
						},
					},
				},
			}))
		case arrayType:
			array = append(array, plistDict("", []plistData{
				{
					key:      "WFItemType",
					dataType: Number,
					value:    itemTypeArray,
				},
				{
					key:      "WFValue",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "Value",
							dataType: Array,
							value:    dictArray(item),
						},
						{
							key:      "WFSerializationType",
							dataType: Text,
							value:    "WFArrayParameterState",
						},
					},
				},
			}))
		case dictType:
			array = append(array, plistDict("", []plistData{
				{
					key:      "WFItemType",
					dataType: Number,
					value:    itemTypeDict,
				},
				{
					key:      "WFValue",
					dataType: Dictionary,
					value: []plistData{
						{
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
											value:    makeDict(item),
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
	return
}
