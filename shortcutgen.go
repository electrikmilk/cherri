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

	"github.com/electrikmilk/args-parser"
)

func generateShortcut() {
	if args.Using("debug") {
		fmt.Print("Generating Shortcut data...")
	}

	shortcut = Shortcut{
		WFWorkflowIcon: ShortcutIcon{
			iconGlyph,
			iconColor,
		},
		WFWorkflowClientVersion:              clientVersion,
		WFWorkflowHasShortcutInputVariables:  hasShortcutInputVariables,
		WFWorkflowMinimumClientVersion:       900,
		WFWorkflowMinimumClientVersionString: "900",
		WFWorkflowTypes:                      definedWorkflowTypes,
		WFQuickActionSurfaces:                definedQuickActions,
		WFWorkflowNoInputBehavior:            noInput,
	}

	waitFor(
		func() {
			shortcut.WFWorkflowInputContentItemClasses = generateInputContentItems()
		},
		func() {
			shortcut.WFWorkflowOutputContentItemClasses = generateOutputContentItems()
		},
		func() {
			shortcut.WFWorkflowImportQuestions = generateImportQuestions()
		},
	)

	generateActions()

	if args.Using("debug") {
		printShortcutGenDebug()
		fmt.Println(ansi("Done.\n", green))
	}

	resetShortcutGen()
}

func resetShortcutGen() {
	tokens = []token{}
	menus = map[string][]varValue{}
	uuids = map[string]string{}
	variables = map[string]varValue{}
	questions = map[string]*question{}
	noInput = map[string]any{}
	definedWorkflowTypes = []string{}
	definedQuickActions = []string{}
	inputs = []string{}
	outputs = []string{}
}

func printShortcutGenDebug() {
	fmt.Println(ansi("### SHORTCUT GEN ###", bold) + "\n")

	fmt.Println(ansi("## UUIDS ##", bold))
	fmt.Println(uuids)

	fmt.Print("\n")
}

func generateActions() {
	uuids = make(map[string]string)
	for _, t := range tokens {
		switch t.typeof {
		case Variable, AddTo, SubFrom, MultiplyBy, DivideBy:
			makeVariableAction(&t)
		case Comment:
			makeCommentAction(t.value.(string))
		case Action:
			var tokenAction = t.value.(action)
			if actions[tokenAction.ident] == nil {
				exit(fmt.Sprintf("Undefined action '%s'", tokenAction.ident))
			}
			setCurrentAction(tokenAction.ident, actions[tokenAction.ident])
			makeAction(tokenAction.args, &map[string]any{})
		case Repeat:
			makeRepeatAction(&t)
		case RepeatWithEach:
			makeRepeatEachAction(&t)
		case Menu:
			makeMenuAction(&t)
		case Item:
			makeMenuItemAction(&t)
		case Conditional:
			makeConditionalAction(&t)
		}
	}
}

func makeCommentAction(comment string) {
	if args.Using("comments") {
		addStdAction("comment", &map[string]any{
			"WFCommentActionText": comment,
		})
	}
}

func makeVariableAction(t *token) {
	var setVariableParams = map[string]any{
		"WFVariableName": t.ident,
	}

	makeVariableInput(t, setVariableParams)

	if t.typeof != Variable {
		if variables[t.ident].valueType != Arr {
			addStdAction("setvariable", &setVariableParams)
			return
		}

		addStdAction("appendvariable", &setVariableParams)
		return
	}

	if v, found := variables[t.ident]; found {
		if v.constant {
			return
		}
	}
	addStdAction("setvariable", &setVariableParams)

	if t.valueType == Arr {
		makeArrayVariable(t)
	}
}

func makeVariableInput(t *token, params map[string]any) {
	var outputName = makeOutputName(t)
	var varUUID string
	if uuids[outputName] == "" {
		varUUID = createUUID(&outputName)
		uuids[outputName] = varUUID
	} else {
		varUUID = uuids[outputName]
	}

	makeVariableValueAction(t, &outputName, &varUUID)
	if t.valueType != Arr && t.value != nil {
		if t.typeof == Variable && t.valueType == Variable {
			params["WFInput"] = variableValue(t.value.(varValue))
		} else {
			params["WFInput"] = inputValue(outputName, varUUID)
		}

		params["WFSerializationType"] = "WFTextTokenAttachment"
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
		addStdAction("dictionary", attachReferenceToParams(&map[string]any{
			"WFItems": makeDictionaryValue(value),
		}, reference))
	}
}

func makeIntValue(reference *map[string]any, value *any) {
	addStdAction("number", attachReferenceToParams(&map[string]any{
		"WFNumberActionNumber": *value,
	}, reference))
}

func makeStringValue(reference *map[string]any, value *any) {
	addStdAction("gettext", attachReferenceToParams(&map[string]any{
		"WFTextActionText": attachmentValues(fmt.Sprintf("%s", *value)),
	}, reference))
}

func makeRawStringValue(reference *map[string]any, value *any) {
	addStdAction("gettext", attachReferenceToParams(&map[string]any{
		"WFTextActionText": fmt.Sprintf("%s", *value),
	}, reference))
}

func makeBoolValue(reference *map[string]any, value *any) {
	var boolValue = "0"
	if *value == true {
		boolValue = "1"
	}

	addStdAction("number", attachReferenceToParams(&map[string]any{
		"WFNumberActionNumber": boolValue,
	}, reference))
}

func makeExpressionValue(reference *map[string]any, value *any) {
	var expression = fmt.Sprintf("%s", *value)
	var expressionParts = strings.Split(expression, " ")
	if len(expressionParts) == 3 && containsTokens(&expression, Plus, Minus, Multiply, Divide) {
		makeMathValue(reference, expression, expressionParts)
		return
	}

	addStdAction("calculateexpression", attachReferenceToParams(&map[string]any{
		"Input": attachmentValues(expression),
	}, reference))
}

func makeMathValue(reference *map[string]any, expression string, expressionParts []string) {
	var operandOne string
	var operandTwo string

	var operation = expressionParts[1]
	expressionParts = strings.Split(expression, operation)

	operandOne = strings.Trim(expressionParts[0], " ")
	operandTwo = strings.Trim(expressionParts[1], " ")

	switch operation {
	case "*":
		operation = "×"
	case "/":
		operation = "÷"
	}

	addStdAction("math", attachReferenceToParams(&map[string]any{
		"WFMathOperation": operation,
		"WFInput":         attachmentValues(operandOne),
		"WFMathOperand":   attachmentValues(operandTwo),
	}, reference))

	return
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
		var variable = token.value.(varValue)
		valueType = variable.valueType
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
		addStdAction("math", attachReferenceToParams(&map[string]any{
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

		addStdAction("gettext", attachReferenceToParams(&map[string]any{
			"WFTextActionText": paramValue(actionArgument{
				valueType: String,
				value:     fmt.Sprintf("{%s}%s", token.ident, varInput),
			}, token.valueType),
		}, reference))
	}
}

func attachReferenceToParams(params *map[string]any, reference *map[string]any) *map[string]any {
	maps.Copy(*params, *reference)

	return params
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
	var varValue = make(map[string]any)
	if variable.constant {
		var varUUID = uuids[identifier]
		varValue["OutputName"] = identifier
		varValue["OutputUUID"] = varUUID
		varValue["Type"] = "ActionOutput"
	} else {
		varValue["VariableName"] = identifier
		varValue["Type"] = varType
		if varType == Ask && variable.prompt != "" {
			varValue["Prompt"] = variable.prompt
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
		var varValue map[string]any
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
			aggr = append(aggr, makeAggrandizement(&variable.valueType, &variable, stringVar.getAs))
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

func makeAggrandizement(valueType *tokenType, variable *varValue, getAs string) map[string]string {
	var aggrandizement = make(map[string]string)
	switch *valueType {
	case Dict:
		aggrandizement["Type"] = "WFDictionaryValueVariableAggrandizement"
	case Action:
		var variableAction = *variable.value.(action).def
		if variableAction.outputType == Dict {
			aggrandizement["Type"] = "WFDictionaryValueVariableAggrandizement"
		} else {
			aggrandizement["Type"] = "WFPropertyVariableAggrandizement"
		}
	default:
		aggrandizement["Type"] = "WFPropertyVariableAggrandizement"
	}

	if aggrandizement["Type"] == "WFDictionaryValueVariableAggrandizement" {
		aggrandizement["DictionaryKey"] = getAs
	} else {
		aggrandizement["PropertyName"] = getAs
	}

	return aggrandizement
}

const utf16BMPThreshold = 0x10000

// mapInlineVars finds occurrences of ObjectReplaceChar and adds them to inlineVars to map the inline variables in noVarString.
// Accounts for UTF-16 characters with code units which require the inline variable position to be doubled (e.g. emoji, bold text).
func mapInlineVars(noVarString *string) {
	var variableIdx int
	var charPos = 0

	for _, r := range *noVarString {
		if r == ObjectReplaceChar {
			inlineVars = append(inlineVars, inlineVar{
				identifier: varIndex[variableIdx].identifier,
				col:        charPos,
				getAs:      varIndex[variableIdx].getAs,
				coerce:     varIndex[variableIdx].coerce,
			})
			variableIdx++
		}

		if r >= utf16BMPThreshold {
			charPos += 2
		} else {
			charPos += 1
		}
	}
}

var collectInlineVarRegex = regexp.MustCompile(`\{@?(.*?)(?:\['(.*?)'])?(?:\.(.*?))?}`)
var replaceVarRegex = regexp.MustCompile(`(\{@?.*?})`)

// collectInlineVariables collects inline variables from `str` and adds them to a slice of attachmentVariable.
// It then replaces all instances of inline variables in `str` with ObjectReplaceChar.
func collectInlineVariables(str *string) (noVarString string) {
	var matches = collectInlineVarRegex.FindAllStringSubmatch(*str, -1)
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
	case Color:
		var colorArgs = arg.value.([]actionArgument)
		return map[string]any{
			"WFColorRepresentationType": "WFColorRepresentationTypeCGColor",
			"redComponent":              colorArgs[0].value.(float64),
			"greenComponent":            colorArgs[1].value.(float64),
			"blueComponent":             colorArgs[2].value.(float64),
			"alphaComponent":            colorArgs[3].value,
		}
	case Quantity:
		return makeQuantityFieldValue(arg.value.([]actionArgument))
	default:
		return attachmentValues(arg.value.(string))
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
	if value == nil {
		return []map[string]any{}
	}
	for key, item := range value.(map[string]interface{}) {
		dictItems = append(dictItems, makeDictionaryItem(key, item))
	}
	return
}

// makeDictionaryItem creates an inner dictionary value.
func makeDictionaryItem(key string, value any) map[string]any {
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
				arrayValue = append(arrayValue, makeDictionaryItem("", item))
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

	return makeDictionaryItemValue(key, itemType, serializedType, wfValue)
}

func makeDictionaryItemValue(key string, itemType dictDataType, serializedType string, wfValue map[string]any) map[string]any {
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

func makeOutputName(token *token) string {
	if variable, found := variables[token.ident]; found {
		if variable.constant {
			return token.ident
		}
	}
	if token.valueType == Variable && token.value != nil {
		var identifier = token.value.(varValue).value.(string)
		if validReference(identifier) {
			return identifier
		}
	}
	var typeOfToken = string(token.valueType)
	if typeOfToken == "action" {
		typeOfToken = token.value.(action).ident
	}

	if token.valueType == "" {
		typeOfToken = token.ident
	}

	var customOutputName = fmt.Sprintf("%s%s", strings.ToTitle(string(typeOfToken[0])), typeOfToken[1:])

	return checkDuplicateOutputName(customOutputName)
}

func makeArrayVariable(t *token) {
	if t.value == nil {
		return
	}
	for _, value := range t.value.([]interface{}) {
		if value == nil {
			continue
		}
		var UUID = createUUID(&t.ident)
		var valueType tokenType
		var itemIdent string
		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			valueType = String
			itemIdent = "Text"
		case reflect.Float64:
			valueType = Integer
			itemIdent = "Number"
		case reflect.Map:
			valueType = Dict
			itemIdent = "Dictionary"
		default:
			exit(fmt.Sprintf("Invalid array value '%v' (%s)", value, reflect.TypeOf(value)))
		}
		makeVariableValueAction(&token{
			typeof:    valueType,
			ident:     itemIdent,
			valueType: valueType,
			value:     value,
		}, &itemIdent, &UUID)

		addStdAction("appendvariable", &map[string]any{
			"WFInput":        inputValue(itemIdent, UUID),
			"WFVariableName": t.ident,
		})
	}
}

func makeConditionalAction(t *token) {
	var conditionalParams = map[string]any{
		"GroupingIdentifier": t.ident,
	}
	switch t.valueType {
	case If:
		conditionalParams["WFControlFlowMode"] = startStatement

		if iosVersion < 18 {
			var cond = t.value.(WFConditions)
			var firstCondition = cond.conditions[0]
			var firstArg = firstCondition.arguments[0]
			conditionalParams["WFInput"] = map[string]any{
				"Type":     "Variable",
				"Variable": variableValue(firstArg.value.(varValue)),
			}
			if len(firstCondition.arguments) > 1 {
				var secondArg = firstCondition.arguments[1]
				conditionalParameter("", conditionalParams, &secondArg.valueType, secondArg.value)
			}
			if len(firstCondition.arguments) > 2 {
				var thirdArg = firstCondition.arguments[2]
				conditionalParameter("WFAnotherNumber", conditionalParams, &thirdArg.valueType, thirdArg.value)
			}
			conditionalParams["WFCondition"] = firstCondition.condition
			conditionalParams["WFControlFlowMode"] = startStatement
		} else {
			var cond = t.value.(WFConditions)
			conditionalParams["WFConditions"] = makeConditions(&cond)
		}
	case Else:
		conditionalParams["WFControlFlowMode"] = statementPart
	case EndClosure:
		conditionalParams["UUID"] = createUUIDReference(t.value.(string))
		conditionalParams["WFControlFlowMode"] = endStatement
	}

	addStdAction("conditional", &conditionalParams)
}

var filterTemplates []map[string]any

func makeConditions(wfConditions *WFConditions) map[string]any {
	filterTemplates = []map[string]any{}
	for _, condition := range wfConditions.conditions {
		var conditionParams = map[string]any{
			"WFCondition": condition.condition,
			"WFInput": map[string]any{
				"Type":     "Variable",
				"Variable": variableValue(condition.arguments[0].value.(varValue)),
			},
		}

		if len(condition.arguments) > 1 {
			var argumentTwo = condition.arguments[1]
			conditionalParameter("", conditionParams, &argumentTwo.valueType, argumentTwo.value)
		}
		if len(condition.arguments) > 2 {
			var argumentThree = condition.arguments[2]
			conditionalParameter("WFAnotherNumber", conditionParams, &argumentThree.valueType, argumentThree.value)
		}

		filterTemplates = append(filterTemplates, conditionParams)
	}

	return map[string]any{
		"Value": map[string]any{
			"WFActionParameterFilterPrefix":    wfConditions.WFActionParameterFilterPrefix,
			"WFActionParameterFilterTemplates": filterTemplates,
		},
		"WFSerializationType": "WFContentPredicateTableTemplate",
	}
}

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

func makeMenuAction(t *token) {
	var menuParams = map[string]any{
		"GroupingIdentifier": t.ident,
		"WFControlFlowMode":  startStatement,
	}
	if t.valueType == EndClosure {
		menuParams["WFControlFlowMode"] = endStatement
		menuParams["UUID"] = createUUIDReference(t.value.(string))
	}
	if t.valueType != EndClosure {
		if t.valueType != Nil {
			menuParams["WFMenuPrompt"] = paramValue(actionArgument{
				valueType: t.valueType,
				value:     t.value,
			}, String)
		}
		var menuItemParams = menus[t.ident]
		var menuItems []map[string]any
		for _, item := range menuItemParams {
			menuItems = append(menuItems, map[string]any{
				"WFItemType": 0,
				"WFValue": paramValue(actionArgument{
					valueType: item.valueType,
					value:     item.value,
				}, String),
			})
		}

		menuParams["WFMenuItems"] = menuItems
	}

	addStdAction("choosefrommenu", &menuParams)
}

func makeMenuItemAction(t *token) {
	addStdAction("choosefrommenu", &map[string]any{
		"GroupingIdentifier": t.ident,
		"WFControlFlowMode":  statementPart,
		"WFMenuItemAttributedTitle": paramValue(actionArgument{
			valueType: t.valueType,
			value:     t.value,
		}, String),
		"WFMenuItemTitle": paramValue(actionArgument{
			valueType: t.valueType,
			value:     t.value,
		}, String),
	})
}

func makeRepeatAction(t *token) {
	var controlFlowMode = startStatement
	if t.valueType == EndClosure {
		controlFlowMode = endStatement
	}
	var repeatParams = map[string]any{
		"WFControlFlowMode":  controlFlowMode,
		"GroupingIdentifier": t.ident,
	}
	if controlFlowMode == endStatement {
		repeatParams["UUID"] = createUUIDReference(t.value.(string))
	}
	if controlFlowMode == startStatement {
		repeatParams["WFRepeatCount"] = paramValue(actionArgument{
			valueType: t.valueType,
			value:     t.value,
		}, Integer)
	}

	addStdAction("repeat.count", &repeatParams)
}

func makeRepeatEachAction(t *token) {
	var controlFlowMode = startStatement
	if t.valueType == EndClosure {
		controlFlowMode = endStatement
	}
	var repeatEachParams = map[string]any{
		"WFControlFlowMode":  controlFlowMode,
		"GroupingIdentifier": t.ident,
	}
	if controlFlowMode == endStatement {
		repeatEachParams["UUID"] = createUUIDReference(t.value.(string))
	}
	if controlFlowMode == startStatement {
		repeatEachParams["WFInput"] = paramValue(actionArgument{
			valueType: t.valueType,
			value:     t.value,
		}, Integer)
	}

	addStdAction("repeat.each", &repeatEachParams)
}

type WFQuestion struct {
	ParameterKey string
	Category     string
	ActionIndex  int
	Text         string
	DefaultValue any
}

func generateImportQuestions() (importQuestions []WFQuestion) {
	if len(questions) == 0 {
		return
	}

	for _, q := range questions {
		if !q.used {
			continue
		}
		importQuestions = append(importQuestions, WFQuestion{
			ParameterKey: q.parameter,
			Category:     "Parameter",
			ActionIndex:  q.actionIndex,
			Text:         q.text,
			DefaultValue: q.defaultValue,
		})
	}
	return
}

func generateInputContentItems() (inputContentItems []string) {
	if len(inputs) == 0 {
		for _, input := range contentItems {
			inputContentItems = append(inputContentItems, input)
		}
		return
	}

	for _, input := range inputs {
		inputContentItems = append(inputContentItems, input)
	}
	return
}

func generateOutputContentItems() (outputContentItems []string) {
	if len(outputs) == 0 {
		return
	}
	for _, output := range outputs {
		outputContentItems = append(outputContentItems, output)
	}

	return
}
