/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
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
			if tokenAction.ident == "rawAction" && len(tokenAction.args) > 0 {
				currentAction.definition.overrideIdentifier = getArgValue(tokenAction.args[0]).(string)
			}
			makeAction(tokenAction.args, &WFActionReference{})
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
		addStdAction("comment", map[string]any{
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
			addStdAction("setvariable", setVariableParams)
			return
		}

		addStdAction("appendvariable", setVariableParams)
		return
	}

	if v, found := variables[t.ident]; found {
		if v.constant {
			return
		}
	}
	addStdAction("setvariable", setVariableParams)

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
	var reference = WFActionReference{
		CustomOutputName: *customOutputName,
		UUID:             *varUUID,
	}

	if (token.typeof == AddTo || token.typeof == SubFrom || token.typeof == MultiplyBy || token.typeof == DivideBy) &&
		token.valueType != Arr &&
		variables[token.ident].valueType != Arr {
		variableValueModifier(token, &reference)
		return
	}

	makeVariableValue(&reference, token.valueType, &token.value)
}

func makeVariableValue(reference *WFActionReference, valueType tokenType, value *any) {
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
		addStdAction("dictionary", attachReferenceToParams(map[string]any{
			"WFItems": makeDictionaryValue(value),
		}, reference))
	}
}

func makeIntValue(reference *WFActionReference, value *any) {
	addStdAction("number", attachReferenceToParams(map[string]any{
		"WFNumberActionNumber": *value,
	}, reference))
}

func makeStringValue(reference *WFActionReference, value *any) {
	addStdAction("gettext", attachReferenceToParams(map[string]any{
		"WFTextActionText": attachmentValues(fmt.Sprintf("%s", *value)),
	}, reference))
}

func makeRawStringValue(reference *WFActionReference, value *any) {
	addStdAction("gettext", attachReferenceToParams(map[string]any{
		"WFTextActionText": fmt.Sprintf("%s", *value),
	}, reference))
}

func makeBoolValue(reference *WFActionReference, value *any) {
	var boolValue = "0"
	if *value == true {
		boolValue = "1"
	}

	addStdAction("number", attachReferenceToParams(map[string]any{
		"WFNumberActionNumber": boolValue,
	}, reference))
}

func makeExpressionValue(reference *WFActionReference, value *any) {
	var expression = fmt.Sprintf("%s", *value)
	var expressionParts = strings.Split(expression, " ")
	if len(expressionParts) == 3 && containsTokens(&expression, Plus, Minus, Multiply, Divide) {
		makeMathValue(reference, expression, expressionParts)
		return
	}

	addStdAction("calculateexpression", attachReferenceToParams(map[string]any{
		"Input": attachmentValues(expression),
	}, reference))
}

func makeMathValue(reference *WFActionReference, expression string, expressionParts []string) {
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

	addStdAction("math", attachReferenceToParams(map[string]any{
		"WFMathOperation": operation,
		"WFInput":         attachmentValues(operandOne),
		"WFMathOperand":   attachmentValues(operandTwo),
	}, reference))

	return
}

func makeDictionaryValue(value *any) WFDictionaryFieldValue {
	return WFDictionaryFieldValue{
		Value: WFDictionaryFieldValueWrapper{
			WFDictionaryFieldValueItems: makeDictionary(*value),
		},
		WFSerializationType: "WFDictionaryFieldValue",
	}
}

func variableValueModifier(token *token, reference *WFActionReference) {
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
		addStdAction("math", attachReferenceToParams(map[string]any{
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

		addStdAction("gettext", attachReferenceToParams(map[string]any{
			"WFTextActionText": paramValue(actionArgument{
				valueType: String,
				value:     fmt.Sprintf("{%s}%s", token.ident, varInput),
			}, token.valueType),
		}, reference))
	}
}

func attachReferenceToParams(params map[string]any, reference *WFActionReference) map[string]any {
	if reference.CustomOutputName != "" {
		params["CustomOutputName"] = reference.CustomOutputName
	}
	if reference.UUID != "" {
		params["UUID"] = reference.UUID
	}
	return params
}

func inputValue(name string, varUUID string) WFTextTokenAttachment {
	var value = Value{}

	if varUUID != "" {
		value.OutputUUID = varUUID
	}

	if variable, found := variables[name]; found {
		if !variable.repeatItem && (variable.constant && variable.valueType != Variable) {
			value.OutputName = name
			value.Type = "ActionOutput"
		} else {
			value.VariableName = name
			value.Type = "Variable"
		}
	} else if global, found := globals[name]; found {
		isInputVariable(name)
		value.Type = global.variableType
	} else {
		value.OutputName = name
		value.Type = "ActionOutput"
	}

	return WFTextTokenAttachment{
		Value:               value,
		WFSerializationType: "WFTextTokenAttachment",
	}
}

func variableValue(variable varValue) any {
	return variableValueWithSerialization(variable, "WFTextTokenAttachment")
}

func variableValueWithSerialization(variable varValue, serializationType string) any {
	var identifier = variable.value.(string)
	var variableReference varValue
	var aggrandizements []Aggrandizement
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
			aggrandizements = append(aggrandizements, Aggrandizement{
				Type:          "WFDictionaryValueVariableAggrandizement",
				DictionaryKey: variable.getAs,
			})
		} else {
			aggrandizements = append(aggrandizements, Aggrandizement{
				PropertyUserInfo: 0,
				Type:             "WFPropertyVariableAggrandizement",
				PropertyName:     variable.getAs,
			})
		}
	}
	if variable.coerce != "" {
		if contentItem, found := contentItems[variable.coerce]; found {
			aggrandizements = append(aggrandizements, Aggrandizement{
				Type:              "WFCoercionVariableAggrandizement",
				CoercionItemClass: contentItem,
			})
		}
	}
	var varType = "Variable"
	if variable.variableType != "" {
		varType = variable.variableType
	}
	var varValue = Value{}
	if variable.constant {
		var varUUID = uuids[identifier]
		varValue.OutputName = identifier
		varValue.OutputUUID = varUUID
		varValue.Type = "ActionOutput"
	} else {
		varValue.VariableName = identifier
		varValue.Type = varType
		if varType == Ask && variable.prompt != "" {
			varValue.Prompt = variable.prompt
		}
	}
	if len(aggrandizements) > 0 {
		varValue.Aggrandizements = aggrandizements
	}

	if serializationType == "" {
		return varValue
	}

	return WFTextTokenAttachment{
		Value:               varValue,
		WFSerializationType: serializationType,
	}
}

type inlineVariable struct {
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

var varPositions map[string]Value
var inlineVariables []inlineVariable
var varIndex []attachmentVariable

func attachmentValues(str string) any {
	if !strings.ContainsAny(str, "{}") {
		return str
	}

	varPositions = make(map[string]Value)
	inlineVariables = []inlineVariable{}
	varIndex = []attachmentVariable{}

	var noVarString = collectInlineVariables(&str)
	makeAttachmentValues()

	return WFTextTokenString{
		Value: WFTextTokenStringValue{
			AttachmentsByRange: varPositions,
			String:             noVarString,
		},
		WFSerializationType: "WFTextTokenString",
	}
}

func makeAttachmentValues() {
	for _, inlineVar := range inlineVariables {
		var varValue, found = getVariableValue(inlineVar.identifier)
		if !found {
			exit(fmt.Sprintf("Undefined reference '%s'", inlineVar.identifier))
		}

		var attachmentValue Value
		var aggrandizements []Aggrandizement
		var varType = "Variable"
		if varValue.variableType != "" {
			varType = varValue.variableType
		}
		if !varValue.constant {
			attachmentValue = Value{
				VariableName: inlineVar.identifier,
				Type:         varType,
			}
		} else {
			var varUUID = uuids[inlineVar.identifier]
			attachmentValue = Value{
				OutputName: inlineVar.identifier,
				OutputUUID: varUUID,
				Type:       "ActionOutput",
			}
		}

		if inlineVar.getAs != "" {
			aggrandizements = append(aggrandizements, makeAggrandizement(&varValue.valueType, varValue, inlineVar.getAs))
		}
		if inlineVar.coerce != "" {
			if contentItem, found := contentItems[inlineVar.coerce]; found {
				aggrandizements = append(aggrandizements, Aggrandizement{
					Type:              "WFCoercionVariableAggrandizement",
					CoercionItemClass: contentItem,
				})
			} else {
				var list = makeKeyList("Available content item types:", contentItems, inlineVar.coerce)
				parserError(fmt.Sprintf("Invalid content item for type coerce '%s'\n\n%s\n", inlineVar.coerce, list))
			}
		}
		if inlineVar.getAs != "" || inlineVar.coerce != "" {
			attachmentValue.Aggrandizements = aggrandizements
		}

		var positionsKey = fmt.Sprintf("{%d, 1}", inlineVar.col)
		varPositions[positionsKey] = attachmentValue
	}
}

func makeAggrandizement(valueType *tokenType, variable *varValue, getAs string) (aggrandizement Aggrandizement) {
	switch *valueType {
	case Dict:
		aggrandizement.Type = "WFDictionaryValueVariableAggrandizement"
	case Action:
		var variableAction = *variable.value.(action).def
		if variableAction.outputType == Dict {
			aggrandizement.Type = "WFDictionaryValueVariableAggrandizement"
		} else {
			aggrandizement.Type = "WFPropertyVariableAggrandizement"
		}
	default:
		aggrandizement.Type = "WFPropertyVariableAggrandizement"
	}

	if aggrandizement.Type == "WFDictionaryValueVariableAggrandizement" {
		aggrandizement.DictionaryKey = getAs
	} else {
		aggrandizement.PropertyName = getAs
	}

	return aggrandizement
}

const utf16BMPThreshold = 0x10000

// mapInlineVariables finds occurrences of ObjectReplaceChar and adds them to inlineVariables to map the inline variables in noVarString.
// Accounts for UTF-16 characters with code units which require the inline variable position to be doubled (e.g. emoji, bold text).
func mapInlineVariables(noVarString *string) {
	var variableIdx int
	var charPos = 0

	for _, r := range *noVarString {
		if r == ObjectReplaceChar {
			inlineVariables = append(inlineVariables, inlineVariable{
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

	mapInlineVariables(&noVarString)
	return
}

func argumentValue(args []actionArgument, idx int) any {
	var actionParameter parameterDefinition
	if len(currentAction.definition.parameters) <= idx {
		// First parameter is likely infinite
		actionParameter = currentAction.definition.parameters[0]
	} else {
		actionParameter = currentAction.definition.parameters[idx]
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
		return WFColorValue{
			WFColorRepresentationType: "WFColorRepresentationTypeCGColor",
			RedComponent:              colorArgs[0].value.(float64),
			GreenComponent:            colorArgs[1].value.(float64),
			BlueComponent:             colorArgs[2].value.(float64),
			AlphaComponent:            colorArgs[3].value,
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
func makeDictionary(value interface{}) (dictItems []WFDictionaryFieldValueItem) {
	if value == nil {
		return []WFDictionaryFieldValueItem{}
	}
	for key, item := range value.(map[string]interface{}) {
		dictItems = append(dictItems, makeDictionaryItem(key, item))
	}
	return
}

func makeDictionaryItem(key string, value any) WFDictionaryFieldValueItem {
	if value == nil {
		value = ""
	}
	var itemType dictDataType
	var wfValue any
	switch v := value.(type) {
	case string:
		itemType = itemTypeText
		if strings.ContainsAny(v, "{}") {
			wfValue = attachmentValues(v)
		} else {
			wfValue = WFTextTokenString{
				WFSerializationType: "WFTextTokenString",
				Value:               WFTextTokenStringValue{String: v},
			}
		}
	case int, float64:
		itemType = itemTypeNumber
		wfValue = WFTextTokenString{
			WFSerializationType: "WFTextTokenString",
			Value:               WFTextTokenStringValue{String: fmt.Sprintf("%v", v)},
		}
	case []interface{}:
		itemType = itemTypeArray
		var items []WFDictionaryFieldValueItem
		for _, item := range v {
			items = append(items, makeDictionaryItem("", item))
		}
		wfValue = WFArrayValue{
			WFSerializationType: "WFArrayParameterState",
			Value:               items,
		}
	case map[string]interface{}:
		itemType = itemTypeDict
		// Shortcuts wraps nested dict items in two WFDictionaryFieldValue layers:
		// the outer marks the value type, the inner holds the items. Top-level
		// dict actions use only one layer (handled by makeDictionaryValue).
		wfValue = WFDictionaryFieldValue{
			WFSerializationType: "WFDictionaryFieldValue",
			Value: WFDictionaryFieldValue{
				WFSerializationType: "WFDictionaryFieldValue",
				Value: WFDictionaryFieldValueWrapper{
					WFDictionaryFieldValueItems: makeDictionary(v),
				},
			},
		}
	case bool:
		itemType = itemTypeBool
		wfValue = WFBoolValue{
			WFSerializationType: "WFNumberSubstitutableState",
			Value:               v,
		}
	default:
		exit(fmt.Sprintf("Unsupported dictionary item value type %T", value))
		return WFDictionaryFieldValueItem{}
	}
	return buildDictionaryItem(key, itemType, wfValue)
}

func buildDictionaryItem(key string, itemType dictDataType, wfValue any) WFDictionaryFieldValueItem {
	var item = WFDictionaryFieldValueItem{
		WFItemType: int(itemType),
		WFValue:    wfValue,
	}
	if key != "" {
		item.WFKey = buildDictionaryKey(key)
	}
	return item
}

func buildDictionaryKey(key string) any {
	if strings.ContainsAny(key, "{}") {
		return attachmentValues(key)
	}
	return WFTextTokenString{
		WFSerializationType: "WFTextTokenString",
		Value:               WFTextTokenStringValue{String: key},
	}
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

		addStdAction("appendvariable", map[string]any{
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
			// Apply the same duration-aware remap as the modern path: Shortcuts
			// uses 1000/1001 for date-relative qty comparisons, not 2/0.
			var conditionCode = firstCondition.condition
			if len(firstCondition.arguments) > 1 && firstCondition.arguments[1].valueType == Quantity {
				switch conditionCode {
				case conditions[GreaterThan]:
					conditionCode = 1000
				case conditions[LessThan]:
					conditionCode = 1001
				}
			}
			if len(firstCondition.arguments) > 1 {
				conditionalParameterLegacy(conditionalParams, firstCondition.arguments[1])
			}
			if len(firstCondition.arguments) > 2 {
				conditionalParameterLegacy(conditionalParams, firstCondition.arguments[2])
			}
			conditionalParams["WFCondition"] = conditionCode
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

	addStdAction("conditional", conditionalParams)
}

var filterTemplates []WFConditionParam

func makeConditions(wfConditions *WFConditions) WFContentPredicateTableTemplate {
	filterTemplates = []WFConditionParam{}
	for _, condition := range wfConditions.conditions {
		// Shortcuts uses 1000 ("is more than N time before now") and 1001
		// ("is within N time of now") when comparing dates by duration. The
		// generic GreaterThan (2) and LessThan (0) codes only apply to numbers.
		var conditionCode = condition.condition
		if len(condition.arguments) > 1 && condition.arguments[1].valueType == Quantity {
			switch conditionCode {
			case conditions[GreaterThan]:
				conditionCode = 1000
			case conditions[LessThan]:
				conditionCode = 1001
			}
		}

		var conditionParam = WFConditionParam{
			WFCondition: conditionCode,
			WFInput: WFInputVariable{
				Type:     "Variable",
				Variable: variableValue(condition.arguments[0].value.(varValue)).(WFTextTokenAttachment),
			},
		}

		if len(condition.arguments) > 1 {
			conditionalParameter(&conditionParam, condition.arguments[1])
		}
		if len(condition.arguments) > 2 {
			conditionalParameter(&conditionParam, condition.arguments[2])
		}

		filterTemplates = append(filterTemplates, conditionParam)
	}

	return WFContentPredicateTableTemplate{
		Value: WFConditionValue{
			WFActionParameterFilterPrefix:    wfConditions.WFActionParameterFilterPrefix,
			WFActionParameterFilterTemplates: filterTemplates,
		},
		WFSerializationType: "WFContentPredicateTableTemplate",
	}
}

func conditionalParameter(param *WFConditionParam, arg actionArgument) {
	switch arg.valueType {
	case String:
		param.WFConditionalActionString = paramValue(arg, String)
	case Integer, Float:
		var val = paramValue(arg, Integer)
		if param.WFNumberValue == nil {
			param.WFNumberValue = val
		} else {
			param.WFAnotherNumber = val
		}
	case Bool:
		var boolNumber = 0
		if arg.value == true {
			boolNumber = 1
		}
		var val = paramValue(actionArgument{valueType: Integer, value: boolNumber}, Integer)
		if param.WFNumberValue == nil {
			param.WFNumberValue = val
		} else {
			param.WFAnotherNumber = val
		}
	case Date:
		var val = paramValue(arg, String)
		if param.WFDate == nil {
			param.WFDate = val
		} else {
			param.WFAnotherDate = val
		}
	case Quantity:
		param.WFDuration = makeQuantityFieldValue(arg.value.([]actionArgument))
	case Variable:
		conditionalParameterVariable(param, arg)
	}
}

func conditionalParameterVariable(param *WFConditionParam, arg actionArgument) {
	var condVarValue = arg.value.(varValue)
	var variable = variables[condVarValue.value.(string)]
	var val = variableValue(condVarValue)
	// When a variable was assigned from an action, resolve the action's declared
	// output type so the comparison value routes to the right plist key.
	var effectiveType = variable.valueType
	if effectiveType == Action {
		if a, ok := variable.value.(action); ok && a.def != nil {
			effectiveType = a.def.outputType
		}
	}
	switch effectiveType {
	case Integer, Float:
		if param.WFNumberValue == nil {
			param.WFNumberValue = val
		} else {
			param.WFAnotherNumber = val
		}
	case Date:
		if param.WFDate == nil {
			param.WFDate = val
		} else {
			param.WFAnotherDate = val
		}
	default:
		param.WFConditionalActionString = attachmentValues(fmt.Sprintf("{%s}", makeVariableReferenceString(condVarValue)))
	}
}

// conditionalParameterLegacy is for iOS < 18 compatibility where we still use map[string]any
func conditionalParameterLegacy(params map[string]any, arg actionArgument) {
	switch arg.valueType {
	case String:
		params["WFConditionalActionString"] = paramValue(arg, String)
	case Integer, Float:
		var val = paramValue(arg, Integer)
		if _, exists := params["WFNumberValue"]; !exists {
			params["WFNumberValue"] = val
		} else {
			params["WFAnotherNumber"] = val
		}
	case Bool:
		var boolNumber = 0
		if arg.value == true {
			boolNumber = 1
		}
		var val = paramValue(actionArgument{valueType: Integer, value: boolNumber}, Integer)
		if _, exists := params["WFNumberValue"]; !exists {
			params["WFNumberValue"] = val
		} else {
			params["WFAnotherNumber"] = val
		}
	case Date:
		var val = paramValue(arg, String)
		if _, exists := params["WFDate"]; !exists {
			params["WFDate"] = val
		} else {
			params["WFAnotherDate"] = val
		}
	case Quantity:
		params["WFDuration"] = makeQuantityFieldValue(arg.value.([]actionArgument))
	case Variable:
		conditionalParameterVariableLegacy(params, arg)
	}
}

func conditionalParameterVariableLegacy(params map[string]any, arg actionArgument) {
	var condVarValue = arg.value.(varValue)
	var variable = variables[condVarValue.value.(string)]
	var val = variableValue(condVarValue)
	// When a variable was assigned from an action, resolve the action's declared
	// output type so the comparison value routes to the right plist key.
	var effectiveType = variable.valueType
	if effectiveType == Action {
		if a, ok := variable.value.(action); ok && a.def != nil {
			effectiveType = a.def.outputType
		}
	}
	switch effectiveType {
	case Integer, Float:
		if _, exists := params["WFNumberValue"]; !exists {
			params["WFNumberValue"] = val
		} else {
			params["WFAnotherNumber"] = val
		}
	case Date:
		if _, exists := params["WFDate"]; !exists {
			params["WFDate"] = val
		} else {
			params["WFAnotherDate"] = val
		}
	default:
		params["WFConditionalActionString"] = attachmentValues(fmt.Sprintf("{%s}", makeVariableReferenceString(condVarValue)))
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
		var menuItems = make([]WFMenuItem, 0, len(menuItemParams))
		for _, item := range menuItemParams {
			menuItems = append(menuItems, WFMenuItem{
				WFItemType: 0,
				WFValue: paramValue(actionArgument{
					valueType: item.valueType,
					value:     item.value,
				}, String),
			})
		}

		menuParams["WFMenuItems"] = menuItems
	}

	addStdAction("choosefrommenu", menuParams)
}

func makeMenuItemAction(t *token) {
	addStdAction("choosefrommenu", map[string]any{
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

	addStdAction("repeat.count", repeatParams)
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

	addStdAction("repeat.each", repeatEachParams)
}

type WFQuestion struct {
	ParameterKey string `plist:",omitempty"`
	Category     string `plist:",omitempty"`
	ActionIndex  int    `plist:",omitempty"`
	Text         string `plist:",omitempty"`
	DefaultValue any    `plist:",omitempty"`
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
