/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
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
		WFQuickActionSurfaces:                definedQuickActions,
		WFWorkflowNoInputBehavior:            noInput,
		WFWorkflowTriggers:                   automationTriggers,
	}

	waitFor(
		func() {
			shortcut.WFWorkflowTypes = generateWorkflowTypes(definedWorkflowTypes)
		},
		func() {
			shortcut.WFWorkflowInputContentItemClasses = generateInputContentItems()
		},
		func() {
			shortcut.WFWorkflowOutputContentItemClasses = generateOutputContentItems()
		},
	)

	generateActions()

	// Must run after generateActions(): question.actionIndex is only populated
	// once each action referencing it has actually been generated.
	shortcut.WFWorkflowImportQuestions = generateImportQuestions()

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
			setCurrentActionFromCall(tokenAction)
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

// rawAction has no identifier of its own, it takes it from its first argument.
func setCurrentActionFromCall(call action) {
	setCurrentAction(call.ident, actions[call.ident])
	if call.ident == "rawAction" {
		currentAction.definition.overrideIdentifier = call.args[0].value.(string)
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

	if _, hasInput := setVariableParams["WFInput"]; !hasInput {
		addStdAction("nothing", map[string]any{})
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
		setCurrentActionFromCall(action)
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
	var value Value
	var global, isGlobal = globals[name]
	switch {
	case isGlobal:
		// A global here points at the action output made for it, so it keeps its type without a variable name.
		isInputVariable(name)
		value = Value{Type: global.variableType}
	case validReference(name):
		value = referenceValue(varValue{value: name})
	default:
		value = Value{OutputName: name, Type: "ActionOutput"}
	}

	if varUUID != "" {
		value.OutputUUID = varUUID
	}

	return WFTextTokenAttachment{
		Value:               value,
		WFSerializationType: "WFTextTokenAttachment",
	}
}

// globalDisplayName translates a Cherri global identifier to the spaced
// Shortcuts display name stored in its varValue (e.g. RepeatIndex → "Repeat Index").
// The globals check is required because mutable user variables also store
// variableType "Variable" — the two maps are the only discriminator.
// Returns identifier unchanged for non-globals or globals without a "Variable" type.
func globalDisplayName(identifier string, v varValue) string {
	if _, isGlobal := globals[identifier]; isGlobal && v.variableType == "Variable" {
		if name, ok := v.value.(string); ok {
			return name
		}
	}
	return identifier
}

// referenceAggrandizements builds the key or property read and the coercion written after a
// variable reference. Which kind of read it is depends on the type of the definition behind
// the reference, because Shortcuts reads dictionary keys and content item properties differently.
func referenceAggrandizements(reference varValue) []Aggrandizement {
	var aggrandizements []Aggrandizement
	if reference.getAs != "" {
		var identifier, _ = reference.value.(string)
		var definition, _ = resolveDefinition(identifier)
		if effectiveType(definition) == Dict {
			aggrandizements = append(aggrandizements, Aggrandizement{
				Type:          "WFDictionaryValueVariableAggrandizement",
				DictionaryKey: reference.getAs,
			})
		} else {
			aggrandizements = append(aggrandizements, Aggrandizement{
				Type:         "WFPropertyVariableAggrandizement",
				PropertyName: reference.getAs,
			})
		}
	}
	if reference.coerce != "" {
		checkCoerce(reference.coerce)
		aggrandizements = append(aggrandizements, Aggrandizement{
			Type:              "WFCoercionVariableAggrandizement",
			CoercionItemClass: contentItems[reference.coerce],
		})
	}

	return aggrandizements
}

// referenceValue is the only place a variable reference becomes a plist Value.
func referenceValue(reference varValue) Value {
	var identifier = reference.value.(string)
	var value = Value{
		VariableName: identifier,
		Type:         "Variable",
	}
	if global, isGlobal := globals[identifier]; isGlobal {
		isInputVariable(identifier)
		value.VariableName = globalDisplayName(identifier, global)
		value.Type = global.variableType
		if global.variableType == Ask && reference.prompt != "" {
			value.Prompt = reference.prompt
		}
	} else if variable, found := variables[identifier]; found && variable.constant {
		value = Value{
			OutputName: identifier,
			OutputUUID: uuids[identifier],
			Type:       "ActionOutput",
		}
	}

	var aggrandizements = referenceAggrandizements(reference)
	if len(aggrandizements) > 0 {
		value.Aggrandizements = aggrandizements
	}

	return value
}

// attachmentReferenceValue is a reference written as a whole value instead of inside text.
// Shortcuts stores a property read there with PropertyUserInfo, but not in the text form of referenceValue.
func attachmentReferenceValue(reference varValue) Value {
	var value = referenceValue(reference)
	for i, aggrandizement := range value.Aggrandizements {
		if aggrandizement.Type == "WFPropertyVariableAggrandizement" {
			value.Aggrandizements[i].PropertyUserInfo = 0
		}
	}

	return value
}

func variableValue(reference varValue) WFTextTokenAttachment {
	return WFTextTokenAttachment{
		Value:               attachmentReferenceValue(reference),
		WFSerializationType: "WFTextTokenAttachment",
	}
}

// textTokenStringValue wraps a reference as a whole text value without writing it as `{name}`
// text and parsing it back, which would lose what the text form cannot carry, like an Ask prompt.
func textTokenStringValue(reference varValue) WFTextTokenString {
	return WFTextTokenString{
		Value: WFTextTokenStringValue{
			String: ObjectReplaceCharStr,
			AttachmentsByRange: map[string]Value{
				"{0, 1}": referenceValue(reference),
			},
		},
		WFSerializationType: "WFTextTokenString",
	}
}

type inlineVariable struct {
	reference varValue
	col       int
}

var varPositions map[string]Value
var inlineVariables []inlineVariable
var varIndex []varValue

func attachmentValues(str string) any {
	if !strings.ContainsAny(str, "{}") {
		return str
	}

	varPositions = make(map[string]Value)
	inlineVariables = []inlineVariable{}
	varIndex = []varValue{}

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
		var identifier = inlineVar.reference.value.(string)
		if !validReference(identifier) {
			exit(fmt.Sprintf("Undefined reference '%s'", identifier))
		}

		var positionsKey = fmt.Sprintf("{%d, 1}", inlineVar.col)
		varPositions[positionsKey] = referenceValue(inlineVar.reference)
	}
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
				reference: varIndex[variableIdx],
				col:       charPos,
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
			if len(match) < 2 {
				continue
			}
			varIndex = append(varIndex, varValue{
				value:  match[1],
				getAs:  match[2],
				coerce: match[3],
			})
		}

		noVarString = replaceVarRegex.ReplaceAllString(*str, ObjectReplaceCharStr)
	}

	mapInlineVariables(&noVarString)
	return
}

func argumentValue(args []actionArgument, idx int) any {
	var actionParameter parameterDefinition
	if len(currentAction.definition.parameters) <= idx {
		// Arguments past the last parameter belong to its variadic parameter
		actionParameter = currentAction.definition.parameters[len(currentAction.definition.parameters)-1]
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
			return textTokenStringValue(arg.value.(varValue))
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
	case RawString:
		return arg.value.(string)
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
		wfValue = makeArrayValue(items)
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
	case WFTextTokenAttachment, WFTextTokenString:
		itemType = itemTypeText
		wfValue = v
	case functionArgReference:
		itemType = itemTypeText
		wfValue = textTokenStringValue(varValue(v))
	case WFDictionaryFieldValue:
		itemType = itemTypeDict
		wfValue = v
	case WFArrayValue:
		itemType = itemTypeArray
		wfValue = v
	case WFBoolValue:
		itemType = itemTypeBool
		wfValue = v
	case varValue:
		itemType = itemTypeFile
		wfValue = WFValue{
			Value:               variableValue(v),
			WFSerializationType: "WFTokenAttachmentParameterState",
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

func makeArrayValue(items []WFDictionaryFieldValueItem) WFArrayValue {
	return WFArrayValue{
		WFSerializationType: "WFArrayParameterState",
		Value:               items,
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

		var cond = t.value.(WFConditions)
		if len(cond.conditions) == 1 || iosVersion < 18 {
			var firstCondition = cond.conditions[0]
			var firstArg = firstCondition.arguments[0]
			conditionalParams["WFInput"] = map[string]any{
				"Type":     "Variable",
				"Variable": variableValue(firstArg.value.(varValue)),
			}
			var conditionCode = firstCondition.condition
			if len(firstCondition.arguments) > 1 && firstCondition.arguments[1].valueType == Quantity {
				switch conditionCode {
				case conditions[GreaterThan]:
					conditionCode = 1000
				case conditions[LessThan]:
					conditionCode = 1001
				}
			}
			var inputType = referenceType(firstArg.value.(varValue))
			if len(firstCondition.arguments) > 1 {
				conditionalParameterLegacy(conditionalParams, firstCondition.arguments[1], inputType)
			}
			if len(firstCondition.arguments) > 2 {
				conditionalParameterLegacy(conditionalParams, firstCondition.arguments[2], inputType)
			}
			conditionalParams["WFCondition"] = conditionCode
		} else {
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

func routeConditionalValue(param *WFConditionParam, val any, inputType tokenType) {
	switch inputType {
	case Date:
		if param.WFDate == nil {
			param.WFDate = val
		} else {
			param.WFAnotherDate = val
		}
	case Integer, Float, Bool:
		if param.WFNumberValue == nil {
			param.WFNumberValue = val
		} else {
			param.WFAnotherNumber = val
		}
	default:
		param.WFConditionalActionString = val
	}
}

func routeConditionalValueLegacy(params map[string]any, val any, inputType tokenType) {
	switch inputType {
	case Date:
		if _, exists := params["WFDate"]; !exists {
			params["WFDate"] = val
		} else {
			params["WFAnotherDate"] = val
		}
	case Integer, Float, Bool:
		if _, exists := params["WFNumberValue"]; !exists {
			params["WFNumberValue"] = val
		} else {
			params["WFAnotherNumber"] = val
		}
	default:
		params["WFConditionalActionString"] = val
	}
}

func makeConditions(wfConditions *WFConditions) WFContentPredicateTableTemplate {
	filterTemplates = []WFConditionParam{}
	for _, condition := range wfConditions.conditions {
		var conditionCode = condition.condition
		if len(condition.arguments) > 1 && condition.arguments[1].valueType == Quantity {
			switch conditionCode {
			case conditions[GreaterThan]:
				conditionCode = 1000
			case conditions[LessThan]:
				conditionCode = 1001
			}
		}

		var inputType = referenceType(condition.arguments[0].value.(varValue))
		var conditionParam = WFConditionParam{
			WFCondition: conditionCode,
			WFInput: WFInputVariable{
				Type:     "Variable",
				Variable: variableValue(condition.arguments[0].value.(varValue)),
			},
		}

		if len(condition.arguments) > 1 {
			conditionalParameter(&conditionParam, condition.arguments[1], inputType)
		}
		if len(condition.arguments) > 2 {
			conditionalParameter(&conditionParam, condition.arguments[2], inputType)
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

func conditionalParameter(param *WFConditionParam, arg actionArgument, inputType tokenType) {
	switch arg.valueType {
	case Quantity:
		param.WFDuration = makeQuantityFieldValue(arg.value.([]actionArgument))
	case Variable:
		conditionalParameterVariable(param, arg, inputType)
	case Bool:
		var boolNumber = 0
		if arg.value == true {
			boolNumber = 1
		}
		routeConditionalValue(param, paramValue(actionArgument{valueType: Integer, value: boolNumber}, Integer), inputType)
	case Integer, Float:
		routeConditionalValue(param, paramValue(arg, arg.valueType), Integer)
	default:
		routeConditionalValue(param, paramValue(arg, arg.valueType), inputType)
	}
}

func conditionalParameterVariable(param *WFConditionParam, arg actionArgument, inputType tokenType) {
	var condVarValue = arg.value.(varValue)
	switch inputType {
	case Date, Integer, Float, Bool:
		routeConditionalValue(param, variableValue(condVarValue), inputType)
	default:
		routeConditionalValue(param, textTokenStringValue(condVarValue), inputType)
	}
}

// conditionalParameterLegacy is for iOS < 18 compatibility where we still use map[string]any
func conditionalParameterLegacy(params map[string]any, arg actionArgument, inputType tokenType) {
	switch arg.valueType {
	case Quantity:
		params["WFDuration"] = makeQuantityFieldValue(arg.value.([]actionArgument))
	case Variable:
		conditionalParameterVariableLegacy(params, arg, inputType)
	case Bool:
		var boolNumber = 0
		if arg.value == true {
			boolNumber = 1
		}
		routeConditionalValueLegacy(params, paramValue(actionArgument{valueType: Integer, value: boolNumber}, Integer), inputType)
	case Integer, Float:
		routeConditionalValueLegacy(params, paramValue(arg, arg.valueType), Integer)
	default:
		routeConditionalValueLegacy(params, paramValue(arg, arg.valueType), inputType)
	}
}

func conditionalParameterVariableLegacy(params map[string]any, arg actionArgument, inputType tokenType) {
	var condVarValue = arg.value.(varValue)
	switch inputType {
	case Date, Integer, Float, Bool:
		routeConditionalValueLegacy(params, variableValue(condVarValue), inputType)
	default:
		routeConditionalValueLegacy(params, textTokenStringValue(condVarValue), inputType)
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
	ActionIndex  int
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
	} else {
		for _, input := range inputs {
			inputContentItems = append(inputContentItems, input)
		}
	}
	sort.Strings(inputContentItems)
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

func generateWorkflowTypes(types []string) []string {
	if len(types) == 0 {
		return []string{"Watch", "WFWorkflowTypeShowInSearch"}
	}
	return types
}
