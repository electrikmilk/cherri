/*
 * Copyright (c) Cherri
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/electrikmilk/args-parser"
	plists "howett.net/plist"
	"math"
	"os"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const UUID = "UUID"

var shortcut Shortcut
var genericShortcut GenericShortcut
var code strings.Builder

func decompile(b []byte) {
	var _, marshalIndexedErr = plists.Unmarshal(b, &shortcut)
	handle(marshalIndexedErr)

	var _, marshalErr = plists.Unmarshal(b, &genericShortcut)
	handle(marshalErr)

	mapIdentifiers()
	mapSplitActions()
	decompileIcon()
	decompileActions()

	if args.Using("debug") {
		printDecompDebug()
	}

	var writeErr = os.WriteFile(relativePath+basename+"_decompiled.cherri", []byte(code.String()), 0600)
	handle(writeErr)
}

// mapIdentifiers creates a map of variable identifiers and UUIDs that are assigned throughout the Shortcut.
func mapIdentifiers() {
	variables = make(map[string]variableValue)
	uuids = make(map[string]string)
	for _, action := range shortcut.WFWorkflowActions {
		var params = action.WFWorkflowActionParameters
		if action.WFWorkflowActionIdentifier == SetVariableIdentifier || action.WFWorkflowActionIdentifier == AppendVariableIdentifier {
			var varName = strings.ReplaceAll(params["WFVariableName"].(string), " ", "")
			if _, found := variables[varName]; !found {
				variables[varName] = variableValue{}
			}
		}

		if params["UUID"] != nil && params["CustomOutputName"] != nil {
			mapUUID(params["UUID"].(string), params["CustomOutputName"].(string))
		}

		var input WFInput
		mapToStruct(params["WFInput"], &input)

		mapValueReference(input.Value)
		mapValueReference(input.Variable.Value)

		if params["WFInput"] != nil {
			var WFInput = params["WFInput"].(map[string]interface{})
			if WFInput["Value"] != nil {
				var value = WFInput["Value"].(map[string]interface{})
				if value["attachmentsByRange"] != nil {
					var attachments = value["attachmentsByRange"].(map[string]interface{})

					for _, attachment := range attachments {
						var attachmentValue Value
						mapToStruct(attachment, &attachmentValue)
						mapValueReference(attachmentValue)
					}
				}
			}
		}
	}
}

func mapValueReference(value Value) {
	if value.OutputName == "" || value.OutputUUID == "" {
		return
	}
	mapUUID(value.OutputUUID, value.OutputName)
}

func mapUUID(uuid string, varName string) {
	if _, found := uuids[uuid]; !found {
		var outputName = strings.ReplaceAll(varName, " ", "")
		uuids[uuid] = checkDuplicateOutputName(outputName)
		variables[outputName] = variableValue{}
	}
}

var currentOutputName string
var duplicateDelta int

func checkDuplicateOutputName(name string) string {
	if name != currentOutputName {
		currentOutputName = name
		duplicateDelta = 0
	}
	for _, outputName := range uuids {
		if outputName == currentOutputName {
			return checkDuplicateOutputName(duplicateOutputName())
		}
	}
	return currentOutputName
}

func duplicateOutputName() string {
	duplicateDelta++

	var revChars = []rune(currentOutputName)
	slices.Reverse(revChars)
	var numChars []rune
	for _, rc := range revChars {
		if rc >= '0' && rc <= '9' {
			numChars = append(numChars, rc)
		}
	}

	if len(numChars) != 0 {
		slices.Reverse(numChars)
		var num = string(numChars)
		var endingDelta, numErr = strconv.Atoi(num)
		handle(numErr)
		endingDelta++

		currentOutputName, _ = strings.CutSuffix(currentOutputName, num)
		duplicateDelta = endingDelta
	}

	currentOutputName = fmt.Sprintf("%s%d", currentOutputName, duplicateDelta)

	return currentOutputName
}

type actionValue struct {
	identifier string
	definition *actionDefinition
}

var identifierMap map[string][]actionValue

// mapSplitActions creates a map of actions that have been split into a few actions to reduce the number of arguments.
func mapSplitActions() {
	identifierMap = make(map[string][]actionValue)
	for identifier, action := range actions {
		var ident = action.identifier
		if action.identifier == "" {
			ident = identifier
		}

		ident = strings.ToLower(ident)

		identifierMap[ident] = append(identifierMap[ident], actionValue{
			identifier: identifier,
			definition: action,
		})
	}
	for identifier, actions := range identifierMap {
		if len(actions) < 2 {
			delete(identifierMap, identifier)
			continue
		}
	}
}

func newCodeLine(s string, v ...any) {
	if tabLevel > 0 {
		for i := 0; i < tabLevel; i++ {
			code.WriteRune('\t')
		}
	}
	code.WriteString(fmt.Sprintf(s, v...))
}

func decompileIcon() {
	var hasDefinitions bool
	var icon = shortcut.WFWorkflowIcon
	if icon.WFWorkflowIconStartColor != iconColor {
		makeColors()
		for name, i := range colors {
			if icon.WFWorkflowIconStartColor != i {
				continue
			}

			newCodeLine(fmt.Sprintf("#define color %s\n", name))
		}
		hasDefinitions = true
	}

	if icon.WFWorkflowIconGlyphNumber != iconGlyph {
		for name, i := range glyphs {
			if icon.WFWorkflowIconGlyphNumber != int64(i) {
				continue
			}

			newCodeLine(fmt.Sprintf("#define glyph %s\n", name))
		}
		hasDefinitions = true
	}

	if hasDefinitions {
		code.WriteRune('\n')
	}
}

var currentVariableValue string

func decompileActions() {
	for _, action := range shortcut.WFWorkflowActions {
		switch action.WFWorkflowActionIdentifier {
		case "is.workflow.actions.comment":
			decompComment(&action)
		case "is.workflow.actions.gettext":
			decompTextValue(&action)
		case "is.workflow.actions.number":
			decompNumberValue(&action)
		case "is.workflow.actions.dictionary":
			decompDictionary(&action)
		case "is.workflow.actions.list":
			decompList(&action)
		case "is.workflow.actions.url":
			decompURL(&action)
		case "is.workflow.actions.calculateexpression":
			decompExpression(&action)
		case SetVariableIdentifier, AppendVariableIdentifier:
			decompVariable(&action)
		case "is.workflow.actions.conditional":
			decompConditional(&action)
		case "is.workflow.actions.repeat.count":
			decompRepeat(&action)
		case "is.workflow.actions.repeat.each":
			decompFor(&action)
		case "is.workflow.actions.choosefrommenu":
			decompMenu(&action)
		default:
			decompAction(&action)
		}
	}
}

func checkConstantLiteral(action *ShortcutAction) {
	if _, found := action.WFWorkflowActionParameters["CustomOutputName"]; found {
		var customOutputName = strings.ReplaceAll(action.WFWorkflowActionParameters["CustomOutputName"].(string), " ", "")
		if _, found := variables[customOutputName]; !found {
			newCodeLine(fmt.Sprintf("const %s = ", customOutputName))
			code.WriteString(currentVariableValue)
			code.WriteRune('\n')
			currentVariableValue = ""
			return
		}
	}
	if _, found := action.WFWorkflowActionParameters[UUID]; found {
		var uuid = action.WFWorkflowActionParameters[UUID].(string)
		if _, found := uuids[uuid]; found {
			newCodeLine(fmt.Sprintf("const %s = ", uuids[uuid]))
			code.WriteString(currentVariableValue)
			code.WriteRune('\n')
			currentVariableValue = ""
		}
	}
}

func decompComment(action *ShortcutAction) {
	var commentText = action.WFWorkflowActionParameters["WFCommentActionText"].(string)
	if strings.Contains(commentText, "\n") {
		code.WriteString(fmt.Sprintf("/*\n%s\n*/\n\n", commentText))
	} else {
		code.WriteString(fmt.Sprintf("// %s\n\n", commentText))
	}
}

func decompTextValue(action *ShortcutAction) {
	currentVariableValue = decompValue(action.WFWorkflowActionParameters["WFTextActionText"])
	if reflect.TypeOf(action.WFWorkflowActionParameters["WFTextActionText"]).String() == "string" {
		return
	}

	checkConstantLiteral(action)
}

func decompNumberValue(action *ShortcutAction) {
	var customOutputName = strings.ReplaceAll(action.WFWorkflowActionParameters["CustomOutputName"].(string), " ", "")
	if _, found := variables[customOutputName]; !found {
		decompAction(action)
		return
	}

	var value = action.WFWorkflowActionParameters["WFNumberActionNumber"]
	if reflect.TypeOf(value).String() == dictType {
		var mapValue = value.(map[string]interface{})
		var Value = mapValue["Value"].(map[string]interface{})
		value = decompValue(value)

		if Value["Type"] == "ActionOutput" {
			currentVariableValue = value.(string)
		}
		return
	}

	var number int
	if value != "" {
		var convErr error
		number, convErr = strconv.Atoi(value.(string))
		handle(convErr)
	}
	currentVariableValue = decompValue(number)
}

func decompExpression(action *ShortcutAction) {
	var input = action.WFWorkflowActionParameters["Input"].(map[string]interface{})
	var expression = strings.Trim(decompValue(input["Value"]), "\"")
	var varRegex = regexp.MustCompile(`{(.*?)}`)
	currentVariableValue = varRegex.ReplaceAllString(expression, "$1")

	checkConstantLiteral(action)
}

func decompVariable(action *ShortcutAction) {
	var variableName = action.WFWorkflowActionParameters["WFVariableName"].(string)
	newCodeLine("@%s", strings.ReplaceAll(variableName, " ", ""))

	if currentVariableValue != "" {
		code.WriteRune(' ')
		if action.WFWorkflowActionIdentifier == AppendVariableIdentifier {
			code.WriteString("+= ")
		} else {
			code.WriteString("= ")
		}

		code.WriteString(currentVariableValue)
	} else {
		var decompInput = decompValue(action.WFWorkflowActionParameters["WFInput"])
		if decompInput != "" {
			code.WriteString(fmt.Sprintf(" = %s", decompInput))
		}
	}

	currentVariableValue = ""
	code.WriteRune('\n')
}

func decompList(action *ShortcutAction) {
	var list strings.Builder
	var listItems = action.WFWorkflowActionParameters["WFItems"].([]interface{})
	var listSize = len(listItems)
	list.WriteString("list(")
	for i, item := range listItems {
		var itemValue = item
		if reflect.TypeOf(item).String() != "string" {
			itemValue = item.(map[string]interface{})["WFValue"]
		}
		list.WriteString(decompValue(itemValue))

		if i < listSize-1 {
			list.WriteRune(',')
		}
	}

	list.WriteRune(')')
	currentVariableValue = list.String()
	list.Reset()

	checkConstantLiteral(action)
}

func decompURL(action *ShortcutAction) {
	var urlValueType = reflect.TypeOf(action.WFWorkflowActionParameters["WFURLActionURL"]).String()
	if urlValueType == dictType || urlValueType == "string" {
		currentVariableValue = fmt.Sprintf("url(%s)", decompValue(action.WFWorkflowActionParameters["WFURLActionURL"]))
		return
	}

	var urlAction strings.Builder
	var urls = action.WFWorkflowActionParameters["WFURLActionURL"].([]interface{})
	var urlsSize = len(urls)
	urlAction.WriteString("url(")
	for i, url := range urls {
		urlAction.WriteString(decompValue(url))

		if i < urlsSize-1 {
			urlAction.WriteRune(',')
		}
	}

	urlAction.WriteRune(')')
	currentVariableValue = urlAction.String()
	urlAction.Reset()

	checkConstantLiteral(action)
}

func decompMenu(action *ShortcutAction) {
	var controlFlowMode = action.WFWorkflowActionParameters["WFControlFlowMode"].(uint64)
	switch controlFlowMode {
	case startStatement:
		newCodeLine("menu ")
		code.WriteString(decompValue(action.WFWorkflowActionParameters["WFMenuPrompt"]))
		code.WriteString(" {\n")
		tabLevel++
	case statementPart:
		newCodeLine("item ")
		code.WriteString(decompValue(action.WFWorkflowActionParameters["WFMenuItemAttributedTitle"]))
		code.WriteString(":\n")
		tabLevel++
	case endStatement:
		tabLevel -= 2
		newCodeLine("}\n")
	}
}

func decompRepeat(action *ShortcutAction) {
	var controlFlowMode = action.WFWorkflowActionParameters["WFControlFlowMode"].(uint64)
	switch controlFlowMode {
	case startStatement:
		if tabLevel == 0 {
			newCodeLine("\nrepeat ")
		} else {
			newCodeLine("repeat ")
		}

		code.WriteString("_ for ")

		code.WriteString(decompValue(action.WFWorkflowActionParameters["WFRepeatCount"]))

		code.WriteString(" {\n")
		tabLevel++
	case endStatement:
		tabLevel--
		newCodeLine("}\n")
	}
}

func decompFor(action *ShortcutAction) {
	var controlFlowMode = action.WFWorkflowActionParameters["WFControlFlowMode"].(uint64)
	switch controlFlowMode {
	case startStatement:
		if tabLevel == 0 {
			newCodeLine("\nfor ")
		} else {
			newCodeLine("for ")
		}

		code.WriteString("_ in ")

		code.WriteString(decompValue(action.WFWorkflowActionParameters["WFInput"]))

		code.WriteString(" {\n")
		tabLevel++
	case endStatement:
		tabLevel--
		newCodeLine("}\n")
	}
}

func decompConditional(action *ShortcutAction) {
	var controlFlowMode = action.WFWorkflowActionParameters["WFControlFlowMode"].(uint64)
	switch controlFlowMode {
	case startStatement:
		makeConditions()

		var conditionInt = int(action.WFWorkflowActionParameters["WFCondition"].(uint64))
		var conditionString = strconv.Itoa(conditionInt)
		var conditionalOperator string
		for operator, cond := range conditions {
			if cond == conditionString {
				conditionalOperator = string(operator)
			}
		}
		if conditionalOperator == "" {
			decompError(fmt.Sprintf("Invalid conditional %s", conditionString), action)
		}

		newCodeLine("if ")

		if conditionalOperator == "!value" {
			code.WriteRune('!')
		}

		code.WriteString(decompValue(action.WFWorkflowActionParameters["WFInput"]))

		if conditionalOperator != "value" && conditionalOperator != "!value" {
			code.WriteRune(' ')
			code.WriteString(conditionalOperator)
			code.WriteRune(' ')

			if _, found := action.WFWorkflowActionParameters["WFNumberValue"]; found {
				var numberType = reflect.TypeOf(action.WFWorkflowActionParameters["WFNumberValue"]).Kind()
				if numberType == reflect.Uint64 {
					code.WriteString(decompValue(action.WFWorkflowActionParameters["WFNumberValue"]))
				} else {
					var numberValue, convErr = strconv.Atoi(action.WFWorkflowActionParameters["WFNumberValue"].(string))
					handle(convErr)
					code.WriteString(decompValue(numberValue))
				}
			} else if _, foundStr := action.WFWorkflowActionParameters["WFConditionalActionString"]; foundStr {
				code.WriteString(decompValue(action.WFWorkflowActionParameters["WFConditionalActionString"]))
			}
		}

		code.WriteString(" {\n")
		tabLevel++
	case statementPart:
		tabLevel--
		newCodeLine("} else {\n")
		tabLevel++
	case endStatement:
		tabLevel--
		newCodeLine("}\n")
	}
}

type DictionaryActionParameters struct {
	WFItems WFItems
}

type WFItems struct {
	Value Value
}

func decompDictionary(action *ShortcutAction) {
	var params DictionaryActionParameters
	mapToStruct(action.WFWorkflowActionParameters, &params)

	var dictionary = decompDictionaryItems(params.WFItems.Value.WFDictionaryFieldValueItems)
	var jsonBytes, jsonErr = json.MarshalIndent(dictionary, strings.Repeat("\t", tabLevel), "\t")
	handle(jsonErr)

	currentVariableValue = string(jsonBytes)

	checkConstantLiteral(action)
}

func decompDictionaryItems(items []WFDictionaryFieldValueItem) (dictionary map[string]interface{}) {
	dictionary = make(map[string]interface{})
	for _, item := range items {
		var itemKey = decompValue(item.WFKey)
		var itemStringValue = decompValue(item.WFValue.Value)
		var itemValueType = fmt.Sprintf("%d", item.WFItemType)
		var itemValue any
		switch dictDataType(itemValueType) {
		case itemTypeNumber:
			itemStringValue = item.WFValue.String
			if itemStringValue != "" {
				var convErr error
				itemValue, convErr = strconv.Atoi(itemStringValue)
				handle(convErr)
			}
		case itemTypeBool:
			itemValue = item.WFValue.Value
		case itemTypeText:
			itemValue = strings.Trim(itemStringValue, "\"")
		case itemTypeArray:
			var arrayItems []ArrayValue
			mapToStruct(item.WFValue.Value.([]interface{}), arrayItems)
			itemValue = decompArray(arrayItems)
		case itemTypeDict:
			var dictionaryItems []WFDictionaryFieldValueItem
			mapToStruct(item.WFValue.Value.(map[string]interface{}), dictionaryItems)
			itemValue = decompDictionaryItems(dictionaryItems)
		default:
			itemValue = itemStringValue
		}
		dictionary[itemKey] = itemValue
	}
	return
}

func decompArray(items []ArrayValue) (array []interface{}) {
	for _, item := range items {
		var itemStringValue = decompValue(item.WFValue)
		var itemValue any
		var itemValueType = fmt.Sprintf("%d", item.WFItemType)
		switch dictDataType(itemValueType) {
		case itemTypeText:
			itemValue = strings.Trim(itemStringValue, "\"")
		case itemTypeNumber:
			if itemStringValue != "" {
				var convErr error
				itemValue, convErr = strconv.Atoi(itemStringValue)
				handle(convErr)
			}
		case itemTypeBool:
			if itemStringValue == "true" {
				itemValue = true
			} else if itemStringValue == "false" {
				itemValue = false
			}
		default:
			itemValue = itemStringValue
		}
		array = append(array, itemValue)
	}
	return
}

func decompValue(value any) string {
	if value == nil {
		return ""
	}
	var valueType = reflect.TypeOf(value).String()
	switch valueType {
	case dictType:
		return decompValueObject(value.(map[string]interface{}))
	case stringType:
		return fmt.Sprintf("\"%s\"", value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func decompValueObject(value map[string]interface{}) string {
	if v, found := value["Value"]; found {
		if reflect.TypeOf(v).String() == dictType {
			value = v.(map[string]interface{})
		}
	}

	switch value["Type"] {
	case "Variable":
		if _, found := value["VariableName"]; found {
			var variableName = value["VariableName"].(string)

			return strings.ReplaceAll(variableName, " ", "")
		}

		var variableValue = value["Variable"].(map[string]interface{})
		return decompValue(variableValue["Value"])
	case "ActionOutput":
		if _, found := value["OutputUUID"]; found {
			return strings.ReplaceAll(uuids[value["OutputUUID"].(string)], " ", "")
		}
	case "ExtensionInput":
		return "ShortcutInput"
	}

	return decompObjectValue(value)
}

func decompObjectValue(valueObj any) string {
	var valueType = reflect.TypeOf(valueObj).String()
	switch valueType {
	case "map[string]interface {}":
		var value = valueObj.(map[string]interface{})

		var attachmentString string
		if value["value"] != nil {
			if reflect.TypeOf(value["value"]).String() != "map[string]interface {}" {
				return fmt.Sprintf("%v", valueObj)
			}
			value = value["value"].(map[string]interface{})
		}

		if _, found := value["string"]; found {
			attachmentString = value["string"].(string)
		}

		var attachmentChars = strings.Split(attachmentString, "")
		if attachments, found := value["attachmentsByRange"]; found {
			for attachmentRange, a := range attachments.(map[string]interface{}) {
				var attachmentRanges = strings.Split(attachmentRange, ",")
				var attachmentPosition = strings.TrimPrefix(attachmentRanges[0], "{")
				var position, convErr = strconv.Atoi(attachmentPosition)
				handle(convErr)

				var attachment Value
				mapToStruct(a, &attachment)

				var variableName = attachment.VariableName
				if attachment.OutputName != "" {
					variableName = attachment.OutputName
				}

				if len(attachment.Aggrandizements) != 0 {
					decompAggrandizements(&variableName, attachment.Aggrandizements)
				}

				attachmentChars[position] = fmt.Sprintf("{%s}", strings.ReplaceAll(variableName, " ", ""))
			}

			attachmentString = fmt.Sprintf("\"%s\"", strings.Join(attachmentChars, ""))
		}

		return attachmentString
	default:
		return fmt.Sprintf("%v", valueObj)
	}
}

var revContentItems map[string]string

func decompAggrandizements(reference *string, aggrs []Aggrandizement) {
	if len(revContentItems) == 0 {
		revContentItems = make(map[string]string)
		for key, item := range contentItems {
			revContentItems[item] = key
		}
	}

	var index string
	var coerce string
	for _, aggr := range aggrs {
		switch aggr.Type {
		case "WFCoercionVariableAggrandizement":
			if _, found := revContentItems[aggr.CoercionItemClass]; found {
				coerce = revContentItems[coerce]
			}
		case "WFDictionaryValueVariableAggrandizement":
			index = aggr.DictionaryKey
		case "WFPropertyVariableAggrandizement":
			index = aggr.PropertyName
		}
	}

	if index != "" {
		*reference = fmt.Sprintf("%s['%s']", *reference, index)
	}
	if index != "" && coerce != "" {
		*reference = fmt.Sprintf("%s.%s", *reference, coerce)
	}
}

var macDefinition bool

func decompAction(action *ShortcutAction) {
	var matchedIdentifier, matchedAction = matchAction(action)
	if matchedIdentifier == "" {
		makeRawAction(action)
		return
	}

	if matchedAction.mac && !macDefinition {
		var saveCode = code.String()
		code.Reset()
		code.WriteString(fmt.Sprintf("#define mac true\n%s", saveCode))
		macDefinition = true
	}

	var isVariableValue = false
	var isConstant = false
	var actionCallCode strings.Builder
	if action.WFWorkflowActionParameters["CustomOutputName"] != nil {
		var customOutputName = strings.ReplaceAll(action.WFWorkflowActionParameters["CustomOutputName"].(string), " ", "")
		if _, foundVar := variables[customOutputName]; !foundVar {
			newCodeLine(fmt.Sprintf("const %s = ", customOutputName))
			isConstant = true
		} else {
			isVariableValue = true
		}
	}
	if action.WFWorkflowActionParameters[UUID] != nil {
		var uuid = action.WFWorkflowActionParameters[UUID].(string)
		if _, found := uuids[uuid]; found {
			newCodeLine(fmt.Sprintf("const %s = ", uuids[uuid]))
			isConstant = true
		}
	}

	var actionCallStart = fmt.Sprintf("%s(", matchedIdentifier)
	if !isConstant && !isVariableValue {
		newCodeLine(actionCallStart)
	} else {
		actionCallCode.WriteString(actionCallStart)
	}

	var matchedParamsSize = len(matchedAction.parameters)
	for i, param := range matchedAction.parameters {
		if param.key == "" {
			continue
		}
		if value, found := action.WFWorkflowActionParameters[param.key]; found {
			if i != 0 && matchedParamsSize != 1 && matchedParamsSize > i {
				actionCallCode.WriteString(", ")
			}

			var dValue = decompValue(value)
			actionCallCode.WriteString(dValue)
		}
	}
	actionCallCode.WriteString(")")

	if isVariableValue {
		currentVariableValue = actionCallCode.String()
	} else {
		code.WriteString(actionCallCode.String())
		code.WriteRune('\n')
		currentVariableValue = ""
	}
}

func makeRawAction(action *ShortcutAction) {
	newCodeLine(fmt.Sprintf("rawAction(\"%s\", [\n", action.WFWorkflowActionIdentifier))
	tabLevel++
	newCodeLine("{\n")

	for key, param := range action.WFWorkflowActionParameters {
		if key == UUID {
			continue
		}

		code.WriteString(strings.Repeat("\t", tabLevel+1))
		code.WriteString(fmt.Sprintf("\"%s\": ", key))

		var value = decompValue(param)
		if !strings.Contains(value, "\"") {
			value = fmt.Sprintf("\"{%s}\"", value)
		}

		code.WriteString(value)
		code.WriteRune('\n')
	}

	newCodeLine("}\n")
	tabLevel--
	newCodeLine("])\n")
}

func matchAction(action *ShortcutAction) (name string, definition actionDefinition) {
	for call, def := range actions {
		var identifier = strings.ToLower(call)
		if def.identifier != "" {
			identifier = def.identifier
		}
		var shortcutsIdentifier = fmt.Sprintf("is.workflow.actions.%s", identifier)
		if shortcutsIdentifier == action.WFWorkflowActionIdentifier || definition.appIdentifier == action.WFWorkflowActionIdentifier {
			definition = *def
			name = call

			if splitActions, found := identifierMap[identifier]; found {
				matchSplitAction(&splitActions, action.WFWorkflowActionParameters, &name, &definition)
				break
			}
			if name == "run" {
				if _, isSelf := action.WFWorkflowActionParameters["isSelf"]; isSelf {
					name = "runSelf"
				} else if wfName, foundName := action.WFWorkflowActionParameters["workflowName"]; foundName {
					if wfName == basename {
						name = "runSelf"
					}
				}
			}
			break
		}
	}
	return
}

type actionMatch struct {
	params float64
	values float64
	action actionValue
}

var matches []actionMatch

func matchSplitAction(splitActions *[]actionValue, parameters map[string]any, identifier *string, definition *actionDefinition) {
	matches = []actionMatch{}

	var defaultAction, hasDefaultAction = getDefaultAction(splitActions)
	if hasDefaultAction {
		*identifier = defaultAction.identifier
		*definition = *defaultAction.definition
	}

	for _, splitAction := range *splitActions {
		var splitActionParams = splitAction.definition.parameters
		if splitAction.definition.addParams != nil {
			for _, addParam := range splitAction.definition.addParams([]actionArgument{}) {
				if addParam.key == "CustomOutputName" || addParam.key == UUID {
					continue
				}
				splitActionParams = append(splitActionParams, parameterDefinition{
					key:          addParam.key,
					defaultValue: addParam.value,
				})
			}
		}

		if !hasRequiredParams(parameters, &splitActionParams) {
			continue
		}

		var matchedParams float64
		var matchedValues float64
		for _, param := range splitActionParams {
			if param.key == "" || param.defaultValue == nil {
				continue
			}
			if value, found := parameters[param.key]; found {
				matchedParams++

				var defaultValue = fmt.Sprintf("%v", param.defaultValue)
				var rawValue = strings.Trim(decompValue(value), "\"")
				if defaultValue == rawValue && defaultValue != "" && rawValue != "" {
					matchedValues++
				}
			}
		}

		if matchedParams == 0 {
			continue
		}
		var splitActionParamsSize = float64(len(splitActionParams))
		if matchedParams > 0 {
			matchedParams = math.Max(splitActionParamsSize-matchedParams, 0)
		}
		if matchedValues > 0 {
			matchedValues = math.Max(splitActionParamsSize-matchedValues, 0)
		}

		matches = append(matches, actionMatch{
			params: matchedParams,
			values: matchedValues,
			action: splitAction,
		})
	}
	if len(matches) < 1 {
		return
	}

	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].params > matches[j].params || matches[i].values > matches[j].values
	})

	var matchedAction = matches[0]
	*identifier = matchedAction.action.identifier
	*definition = *matchedAction.action.definition
}

func hasRequiredParams(parameters map[string]any, definitions *[]parameterDefinition) bool {
	for _, def := range *definitions {
		if def.optional || def.key == "" {
			continue
		}
		if _, found := parameters[def.key]; !found {
			return false
		}
	}

	for key := range parameters {
		if key == "CustomOutputName" || key == UUID {
			continue
		}
		if isKeyDefined(definitions, &key) {
			continue
		}

		return false
	}

	return true
}

func isKeyDefined(definitions *[]parameterDefinition, key *string) bool {
	for _, def := range *definitions {
		if def.key == *key {
			return true
		}
	}

	return false
}

// Get default action from a slice of split actions for an identifier.
func getDefaultAction(splitActions *[]actionValue) (action actionValue, found bool) {
	for _, splitAction := range *splitActions {
		if splitAction.definition.defaultAction {
			action = splitAction
			found = true
			return
		}
	}

	return
}

func printDecompDebug() {
	fmt.Println(ansi("##### DEBUG #####\n", red))

	fmt.Println("### ACTIONS ###")
	for _, action := range shortcut.WFWorkflowActions {
		fmt.Println(action.WFWorkflowActionIdentifier)
		var maxKeySize int
		for key := range action.WFWorkflowActionParameters {
			var keySize = len(key)
			if keySize > maxKeySize {
				maxKeySize = keySize
			}
		}
		for key, value := range action.WFWorkflowActionParameters {
			fmt.Println("\t", key, strings.Repeat(" ", maxKeySize-len(key)), value)
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")

	fmt.Println("### VARIABLES ###")
	printVariables()
	fmt.Print("\n")

	fmt.Println("### UUIDS ###")
	for uuid, name := range uuids {
		fmt.Printf("%s | %s\n", uuid, name)
	}
	fmt.Print("\n")
}

func decompError(message string, action *ShortcutAction) {
	fmt.Print(ansi("[Decompilation Error]\n\n", red, bold))

	fmt.Println(ansi(fmt.Sprintf("%s\n", message), red))

	var identifier = strings.Replace(action.WFWorkflowActionIdentifier, "is.workflow.actions.", "", 1)

	fmt.Println("Action identifier:", identifier)
	fmt.Println("Full action identifier:", action.WFWorkflowActionIdentifier)

	lines = strings.Split(code.String(), "\n")
	var linesLen = len(lines)
	var lastWrittenLine = lines[linesLen-1]
	var prevWrittenLine = lines[linesLen-2]
	fmt.Printf("\nStopped while writing line %d:\n", linesLen)
	fmt.Printf("%s %s\n", ansi(fmt.Sprintf("%d |", linesLen-1), dim), ansi(prevWrittenLine, dim))
	fmt.Printf("%s %s\n", ansi(fmt.Sprintf("%d |", linesLen), dim), ansi(lastWrittenLine, red))

	os.Exit(1)
}
