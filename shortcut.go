/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"maps"
	"reflect"
	"regexp"
	"strings"
)

/*
 Shortcut File Format Data Structures
*/

type ShortcutIcon struct {
	WFWorkflowIconGlyphNumber int64
	WFWorkflowIconStartColor  int
}

type Shortcut struct {
	WFWorkflowIcon                       ShortcutIcon
	WFWorkflowActions                    []ShortcutAction
	WFQuickActionSurfaces                []string
	WFWorkflowInputContentItemClasses    []string
	WFWorkflowClientVersion              string
	WFWorkflowMinimumClientVersion       int
	WFWorkflowMinimumClientVersionString string
	WFWorkflowImportQuestions            interface{}
	WFWorkflowTypes                      []string
	WFWorkflowOutputContentItemClasses   []string
	WFWorkflowHasShortcutInputVariables  bool
	WFWorkflowHasOutputFallback          bool
	WFWorkflowNoInputBehavior            WFWorkflowNoInputBehavior
	WFWorkflowName                       string
}

var shortcut Shortcut

type ShortcutAction struct {
	WFWorkflowActionIdentifier string
	WFWorkflowActionParameters map[string]any
}

type WFWorkflowNoInputBehavior struct {
	Name       string
	Parameters map[string]string
}

type Value struct {
	Type                        string
	VariableName                string
	OutputUUID                  string
	OutputName                  string
	Value                       any
	Variable                    any
	WFDictionaryFieldValueItems []WFDictionaryFieldValueItem
	AttachmentsByRange          map[string]Value
	String                      string
	Aggrandizements             []Aggrandizement
}

type Aggrandizement struct {
	Type              string
	CoercionItemClass string
	DictionaryKey     string
	PropertyName      string
	PropertyUserInfo  any
}

type WFDictionaryFieldValueItem struct {
	WFKey      any
	WFItemType int
	WFValue    any
}

type WFValue struct {
	Value               any
	String              string
	WFSerializationType string
}

type WFInput struct {
	Value Value
}

type WFContactFieldValue struct {
	EntryType       int
	SerializedEntry map[string]interface{}
}

type ImageSize struct {
	Value SizeValue
}

type SizeValue struct {
	Unit      string
	Magnitude string
}

type WFMeasurementUnit struct {
	WFNSUnitSymbol any
	Value          SizeValue
}

var uuids map[string]string

type dataType string

const Text dataType = "string"
const Number dataType = "integer"
const Dictionary dataType = "dictionary"
const Array dataType = "array"
const Boolean dataType = "boolean"

type dictDataType int

const itemTypeText dictDataType = 0
const itemTypeNumber dictDataType = 3
const itemTypeArray dictDataType = 2
const itemTypeDict dictDataType = 1
const itemTypeBool dictDataType = 4

var noInput WFWorkflowNoInputBehavior

var hasShortcutInputVariables = false

// ObjectReplaceChar is a Shortcuts convention to mark the placement of inline variables in a string.
const ObjectReplaceChar = '\uFFFC'
const ObjectReplaceCharStr = "\uFFFC"

func conditionalParameter(key string, conditionalParams map[string]any, typeOf *tokenType, value any) {
	if key == "" {
		if *typeOf == String {
			key = "WFConditionalActionString"
		} else if *typeOf == Integer || *typeOf == Bool {
			key = "WFNumberValue"
		}
	}
	switch *typeOf {
	case String:
		conditionalParams[key] = paramValue(actionArgument{
			valueType: *typeOf,
			value:     value,
		}, String)
	case Integer:
		conditionalParams[key] = paramValue(actionArgument{
			valueType: *typeOf,
			value:     value,
		}, String)
	case Bool:
		var boolNumber = "0"
		if value == true {
			boolNumber = "1"
		}
		conditionalParams[key] = paramValue(actionArgument{
			valueType: Integer,
			value:     boolNumber,
		}, Integer)
	case Variable:
		conditionalParameterVariable(conditionalParams, value)
	}
}

func conditionalParameterVariable(conditionalParams map[string]any, value any) {
	var condVarValue = value.(varValue)
	var variable = variables[condVarValue.value.(string)]
	switch variable.valueType {
	case Integer:
		conditionalParams["WFNumberValue"] = variableValue(condVarValue)
	default:
		conditionalParams["WFConditionalActionString"] = attachmentValues(fmt.Sprintf("{%s}", makeVariableReferenceString(condVarValue)))
	}
}

func makeVariableValueAction(token *token, customOutputName *string, varUUID *string) {
	var reference = map[string]any{
		"CustomOutputName": *customOutputName,
		"UUID":             *varUUID,
	}

	if (token.typeof == AddTo || token.typeof == SubFrom || token.typeof == MultiplyBy || token.typeof == DivideBy) &&
		token.valueType != Arr &&
		variables[token.ident].valueType != Arr {
		variableValueModifier(token, &reference)
		return
	}

	makeVariableValue(&reference, token.valueType, &token.value)
}

func makeVariableValue(reference *map[string]any, valueType tokenType, value *any) {
	switch valueType {
	case Integer, Float:
		makeIntValue(reference, value)
	case Bool:
		makeBoolValue(reference, value)
	case String:
		makeStringValue(reference, value)
	case RawString:
		makeRawStringValue(reference, value)
	case Expression:
		makeExpressionValue(reference, value)
	case Action:
		var valuePtr = *value
		var action = valuePtr.(action)
		setCurrentAction(action.ident, actions[action.ident])
		makeAction(action.args, reference)
	case Dict:
		buildStdAction("dictionary", attachReferenceParams(&map[string]any{
			"WFItems": makeDictionaryValue(value),
		}, reference))
	}
}

func makeIntValue(reference *map[string]any, value *any) {
	buildStdAction("number", attachReferenceParams(&map[string]any{
		"WFNumberActionNumber": *value,
	}, reference))
}

func makeStringValue(reference *map[string]any, value *any) {
	buildStdAction("gettext", attachReferenceParams(&map[string]any{
		"WFTextActionText": attachmentValues(fmt.Sprintf("%s", *value)),
	}, reference))
}

func makeRawStringValue(reference *map[string]any, value *any) {
	buildStdAction("gettext", attachReferenceParams(&map[string]any{
		"WFTextActionText": fmt.Sprintf("%s", *value),
	}, reference))
}

func makeBoolValue(reference *map[string]any, value *any) {
	var boolValue = "0"
	if *value == true {
		boolValue = "1"
	}

	buildStdAction("number", attachReferenceParams(&map[string]any{
		"WFNumberActionNumber": boolValue,
	}, reference))
}

var formattedExpression []string

func makeExpressionValue(reference *map[string]any, value *any) {
	formattedExpression = []string{}
	var expression = fmt.Sprintf("%s", *value)
	var expressionParts = strings.Split(expression, " ")
	if len(expressionParts) == 3 && containsTokens(&expression, Plus, Minus, Multiply, Divide) {
		var operandOne string
		var operandTwo string

		var operation = expressionParts[1]
		expressionParts = strings.Split(expression, operation)
		operandOne = strings.Trim(expressionParts[0], " ")
		operandTwo = strings.Trim(expressionParts[1], " ")
		wrapVariableReference(&operandOne)
		wrapVariableReference(&operandTwo)

		buildStdAction("math", attachReferenceParams(&map[string]any{
			"WFScientificMathOperation": attachmentValues(operation),
			"WFMathOperation":           attachmentValues(operation),
			"WFInput":                   attachmentValues(operandOne),
			"WFMathOperand":             attachmentValues(operandTwo),
		}, reference))

		return
	}

	expressionParts = strings.Split(expression, " ")
	for _, part := range expressionParts {
		var p = part
		wrapVariableReference(&p)
		formattedExpression = append(formattedExpression, p)
	}

	expression = strings.Join(formattedExpression, " ")

	buildStdAction("calculateexpression", attachReferenceParams(&map[string]any{
		"Input": attachmentValues(expression),
	}, reference))
}

func makeDictionaryValue(value *any) map[string]any {
	return map[string]any{
		"Value": map[string]any{
			"WFDictionaryFieldValueItems": makeDictionary(*value),
		},
		"WFSerializationType": "WFDictionaryFieldValue",
	}
}

func variableValueModifier(token *token, reference *map[string]any) {
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
		buildStdAction("math", attachReferenceParams(&map[string]any{
			"WFMathOperand": paramValue(actionArgument{
				valueType: token.valueType,
				value:     token.value,
			}, token.valueType),
			"WFInput": paramValue(actionArgument{
				valueType: Variable,
				value:     varValue{value: token.ident},
			}, token.valueType),
			"WFMathOperation": operation,
		}, reference))
	case String, RawString:
		var varInput = token.value.(string)
		wrapVariableReference(&varInput)

		buildStdAction("gettext", attachReferenceParams(&map[string]any{
			"WFTextActionText": paramValue(actionArgument{
				valueType: String,
				value:     fmt.Sprintf("{%s}%s", token.ident, varInput),
			}, token.valueType),
		}, reference))
	}
}

func attachReferenceParams(params *map[string]any, reference *map[string]any) map[string]any {
	maps.Copy(*params, *reference)

	return *params
}

func inputValue(name string, varUUID string) map[string]any {
	var value = make(map[string]any)

	if varUUID != "" {
		value["OutputUUID"] = varUUID
	}

	if variable, found := variables[name]; found {
		if !variable.repeatItem && (variable.constant && variable.valueType != Variable) {
			value["OutputName"] = name
			value["Type"] = "ActionOutput"
		} else {
			value["VariableName"] = name
			value["Type"] = "Variable"
		}
	} else if global, found := globals[name]; found {
		isInputVariable(name)
		value["Type"] = global.variableType
	} else {
		value["OutputName"] = name
		value["Type"] = "ActionOutput"
	}

	return map[string]any{
		"Value":               value,
		"WFSerializationType": "WFTextTokenAttachment",
	}
}

func variableValue(variable varValue) map[string]any {
	var identifier = variable.value.(string)
	var variableReference varValue
	var aggrandizements []map[string]any
	if global, found := globals[identifier]; found {
		isInputVariable(identifier)
		variable.variableType = global.variableType
	} else if v, found := variables[identifier]; found {
		variableReference = v
		variable.constant = v.constant
		variable.repeatItem = v.repeatItem
	}
	if variable.getAs != "" {
		var refValueType = variable.valueType
		if variable.valueType == Variable && variableReference.valueType != "" {
			refValueType = variableReference.valueType
		}
		if refValueType == Dict {
			aggrandizements = append(aggrandizements, map[string]any{
				"Type":          "WFDictionaryValueVariableAggrandizement",
				"DictionaryKey": variable.getAs,
			})
		} else {
			aggrandizements = append(aggrandizements, map[string]any{
				"PropertyUserInfo": 0,
				"Type":             "WFPropertyVariableAggrandizement",
				"PropertyName":     variable.getAs,
			})
		}
	}
	if variable.coerce != "" {
		if contentItem, found := contentItems[variable.coerce]; found {
			aggrandizements = append(aggrandizements, map[string]any{
				"Type":              "WFCoercionVariableAggrandizement",
				"CoercionItemClass": contentItem,
			})
		}
	}
	var varType = "Variable"
	if variable.variableType != "" {
		varType = variable.variableType
	}
	var varValue map[string]any
	if variable.constant {
		var varUUID = uuids[identifier]
		varValue = map[string]any{
			"OutputName": identifier,
			"OutputUUID": varUUID,
			"Type":       "ActionOutput",
		}
	} else {
		varValue = map[string]any{
			"VariableName": identifier,
			"Type":         varType,
		}
	}
	if len(aggrandizements) > 0 {
		varValue["Aggrandizements"] = aggrandizements
	}

	return map[string]any{
		"Value":               varValue,
		"WFSerializationType": "WFTextTokenAttachment",
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

var varPositions map[string]any
var inlineVars []inlineVar
var varIndex []attachmentVariable

func attachmentValues(str string) any {
	if !strings.ContainsAny(str, "{}") {
		return str
	}

	varPositions = make(map[string]any)
	inlineVars = []inlineVar{}
	varIndex = []attachmentVariable{}

	var noVarString = collectInlineVariables(&str)
	makeAttachmentValues()

	return map[string]any{
		"Value": map[string]any{
			"attachmentsByRange": varPositions,
			"string":             noVarString,
		},
		"WFSerializationType": "WFTextTokenString",
	}
}

func makeAttachmentValues() {
	for _, stringVar := range inlineVars {
		var storedVar varValue
		if g, global := globals[stringVar.identifier]; global {
			isInputVariable(stringVar.identifier)
			storedVar = g
			stringVar.identifier = g.value.(string)
		} else if v, found := variables[stringVar.identifier]; found {
			storedVar = v
		} else {
			exit(fmt.Sprintf("Undefined reference '%s'", stringVar.identifier))
		}
		var variable = variables[stringVar.identifier]
		var varUUID = uuids[stringVar.identifier]
		var varValue = map[string]any{}
		var varType = "Variable"
		var aggr []map[string]string
		if storedVar.variableType != "" {
			varType = storedVar.variableType
		}
		if !variable.constant {
			varValue = map[string]any{
				"VariableName": stringVar.identifier,
				"Type":         varType,
			}
		} else {
			varValue = map[string]any{
				"OutputName": stringVar.identifier,
				"OutputUUID": varUUID,
				"Type":       "ActionOutput",
			}
		}
		if stringVar.getAs != "" {
			var refValueType = variable.valueType
			if refValueType == Dict {
				aggr = append(aggr, map[string]string{
					"Type":          "WFDictionaryValueVariableAggrandizement",
					"DictionaryKey": stringVar.getAs,
				})
			} else {
				aggr = append(aggr, map[string]string{
					"Type":         "WFPropertyVariableAggrandizement",
					"PropertyName": stringVar.getAs,
				})
			}
		}
		if stringVar.coerce != "" {
			if contentItem, found := contentItems[stringVar.coerce]; found {
				aggr = append(aggr, map[string]string{
					"Type":              "WFCoercionVariableAggrandizement",
					"CoercionItemClass": contentItem,
				})
			} else {
				var list = makeKeyList("Available content item types:", contentItems, stringVar.coerce)
				parserError(fmt.Sprintf("Invalid content item for type coerce '%s'\n\n%s\n", stringVar.coerce, list))
			}
		}
		if stringVar.getAs != "" || stringVar.coerce != "" {
			varValue["Aggrandizements"] = aggr
		}

		var positionsKey = fmt.Sprintf("{%d, 1}", stringVar.col)
		varPositions[positionsKey] = varValue
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

var replaceVarRegex = regexp.MustCompile(`(\{.*?})`)

// collectInlineVariables collects inline variables from `str` and adds them to a slice of attachmentVariable.
// It then replaces all instances of inline variables in `str` with ObjectReplaceChar.
func collectInlineVariables(str *string) (noVarString string) {
	var matches = collectVarRegex.FindAllStringSubmatch(*str, -1)
	if matches != nil {
		for _, match := range matches {
			var attachmentVar attachmentVariable
			if len(match) < 2 {
				continue
			}
			attachmentVar.identifier = match[1]
			if len(match[2]) > 0 {
				attachmentVar.getAs = match[2]
			}
			if len(match[3]) > 0 {
				attachmentVar.coerce = match[3]
			}
			varIndex = append(varIndex, attachmentVar)
		}

		noVarString = replaceVarRegex.ReplaceAllString(*str, ObjectReplaceCharStr)
	}

	mapInlineVars(&noVarString)
	return
}

func convertDataTypeToTokenType(dataType dataType) tokenType {
	switch dataType {
	case Text:
		return String
	case Number:
		return Integer
	case Array:
		return Arr
	case Dictionary:
		return Dict
	case Boolean:
		return Bool
	default:
		return ""
	}
}

func argumentValue(args []actionArgument, idx int) any {
	var actionParameter parameterDefinition
	if len(currentAction.parameters) <= idx {
		// First parameter is likely infinite
		actionParameter = currentAction.parameters[0]
	} else {
		actionParameter = currentAction.parameters[idx]
	}
	var arg actionArgument
	if len(args) <= idx {
		if actionParameter.optional || actionParameter.defaultValue != nil {
			return map[string]any{}
		}
	} else {
		arg = args[idx]
	}

	return paramValue(arg, actionParameter.validType)
}

func paramValue(arg actionArgument, handleAs tokenType) any {
	if arg.valueType == Nil || arg.value == nil {
		return map[string]any{}
	}
	switch arg.valueType {
	case Variable:
		if handleAs == String {
			var refStr = makeVariableReferenceString(arg.value.(varValue))
			return attachmentValues(fmt.Sprintf("{%s}", refStr))
		}

		return variableValue(arg.value.(varValue))
	case Dict:
		return makeDictionaryValue(&arg.value)
	case Integer:
		fallthrough
	case Bool:
		fallthrough
	case Float:
		return arg.value
	default:
		return attachmentValues(arg.value.(string))
	}
}

func wrapVariableReference(s *string) {
	if validReference(*s) {
		*s = fmt.Sprintf("{%s}", *s)
	}
}

// isInputVariable checks if identifier is the ShortcutInput global to set the global boolean in the final plist.
func isInputVariable(identifier string) {
	if hasShortcutInputVariables {
		return
	}
	hasShortcutInputVariables = identifier == ShortcutInput
}

const (
	startStatement uint64 = 0
	statementPart  uint64 = 1
	endStatement   uint64 = 2
)

// makeDictionary creates a Shortcut dictionary value.
func makeDictionary(value interface{}) (dictItems []map[string]any) {
	for key, item := range value.(map[string]interface{}) {
		dictItems = append(dictItems, dictionaryItemValue(key, item))
	}
	return
}

// dictionaryItemValue creates an inner dictionary value.
func dictionaryItemValue(key string, value any) map[string]any {
	if value == nil {
		value = ""
	}
	var itemType dictDataType
	var serializedType string
	var wfValue = map[string]any{
		"Value": map[string]any{
			"string": value,
		},
	}

	if value != "" {
		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			if strings.ContainsAny(value.(string), "{}") {
				wfValue["Value"] = paramValue(actionArgument{
					valueType: String,
					value:     value,
				}, String)
				if reflect.TypeOf(wfValue["Value"]).String() == "map[string]interface {}" {
					for _, val := range wfValue {
						wfValue = val.(map[string]any)
						break
					}
				}
			}
			itemType = itemTypeText
			serializedType = "WFTextTokenString"
		case reflect.Int, reflect.Float64:
			itemType = itemTypeNumber
			serializedType = "WFTextTokenString"
			wfValue = map[string]any{
				"Value": map[string]any{
					"string": fmt.Sprintf("%v", value),
				},
			}
		case reflect.Slice:
			itemType = itemTypeArray
			serializedType = "WFArrayParameterState"
			var arrayValue []map[string]interface{}
			for _, item := range value.([]interface{}) {
				arrayValue = append(arrayValue, dictionaryItemValue("", item))
			}
			wfValue = map[string]any{
				"Value": arrayValue,
			}
		case reflect.Map:
			itemType = itemTypeDict
			serializedType = "WFDictionaryFieldValue"
			wfValue = map[string]any{
				"Value": map[string]any{
					"Value": map[string]any{
						"WFDictionaryFieldValueItems": makeDictionary(value),
					},
					"WFSerializationType": "WFDictionaryFieldValue",
				},
			}
		case reflect.Bool:
			itemType = itemTypeBool
			serializedType = "WFNumberSubstitutableState"
			wfValue = map[string]any{
				"Value": value,
			}
		default:
			exit(fmt.Sprintf("Unsupported dictionary item value type %s", reflect.TypeOf(value)))
		}
	} else {
		itemType = itemTypeText
		serializedType = "WFTextTokenString"
		wfValue = map[string]any{}
	}

	return dictionaryValue(key, itemType, serializedType, wfValue)
}

func dictionaryValue(key string, itemType dictDataType, serializedType string, wfValue map[string]any) map[string]any {
	var wfValueParams = map[string]any{
		"WFSerializationType": serializedType,
	}
	maps.Copy(wfValueParams, wfValue)
	var valueData = map[string]any{
		"WFItemType": itemType,
		"WFValue":    wfValueParams,
	}

	if key != "" {
		var wfKey = map[string]any{
			"Value": map[string]string{
				"string": key,
			},
		}
		if strings.ContainsAny(key, "{}") {
			wfKey["Value"] = paramValue(actionArgument{
				valueType: String,
				value:     key,
			}, String)
			if reflect.TypeOf(wfKey["Value"]).String() == "map[string]interface {}" {
				for _, val := range wfKey["Value"].(map[string]any) {
					wfKey = val.(map[string]any)
					break
				}
			}
		}

		var wfKeyParams = map[string]any{
			"WFSerializationType": "WFTextTokenString",
		}
		maps.Copy(wfKeyParams, wfKey)
		valueData["WFKey"] = wfKeyParams
	}

	return valueData
}
