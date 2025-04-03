/*
 * Copyright (c) Cherri
 */

package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/electrikmilk/args-parser"
	plists "howett.net/plist"
)

const (
	UUID          = "UUID"
	ShortcutInput = "ShortcutInput"
)

var shortcut Shortcut
var code strings.Builder
var specialCharsRegex *regexp.Regexp

func decompile(b []byte) {
	var _, marshalIndexedErr = plists.Unmarshal(b, &shortcut)
	handle(marshalIndexedErr)

	mapIdentifiers()
	mapSplitActions()

	defineName()
	decompileIcon()

	decompileActions()

	if args.Using("debug") {
		printDecompDebug()
	}

	outputPath = getOutputPath(fmt.Sprintf("%s%s.cherri", relativePath, strings.ReplaceAll(basename, " ", "_")))

	var writeErr = os.WriteFile(outputPath, []byte(code.String()), 0600)
	handle(writeErr)
}

// mapIdentifiers creates a map of variable identifiers and UUIDs that are assigned throughout the Shortcut.
func mapIdentifiers() {
	specialCharsRegex = regexp.MustCompile("[^a-zA-Z0-9_]+")
	variables = make(map[string]variableValue)
	uuids = make(map[string]string)
	for _, action := range shortcut.WFWorkflowActions {
		currentActionIdentifier = action.WFWorkflowActionIdentifier
		var params = action.WFWorkflowActionParameters
		if action.WFWorkflowActionIdentifier == SetVariableIdentifier || action.WFWorkflowActionIdentifier == AppendVariableIdentifier {
			var varName = params["WFVariableName"].(string)
			sanitizeIdentifier(&varName)
			if _, found := variables[varName]; !found {
				variables[varName] = variableValue{}
			}

			if action.WFWorkflowActionParameters["WFInput"] != nil {
				var wfInput WFInput
				mapToStruct(action.WFWorkflowActionParameters["WFInput"], &wfInput)
				varUUIDs = append(varUUIDs, wfInput.Value.OutputUUID)
			}
		}

		if params["UUID"] != nil && params["CustomOutputName"] != nil {
			mapUUID(params["UUID"].(string), params["CustomOutputName"].(string))
		}

		checkParamIdentifiers(params)
	}
}

func peekActions(peek int) ShortcutAction {
	return shortcut.WFWorkflowActions[actionIndex+peek]
}

// insertCodeComment creates a comment in the decompiled code.
func insertCodeComment(comment string) {
	newCodeLine(fmt.Sprintf("// %s\n", comment))
}

func checkParamIdentifiers(params map[string]interface{}) {
	for _, value := range params {
		if value == nil || reflect.TypeOf(value).String() != dictType {
			continue
		}

		var paramValues = value.(map[string]interface{})
		checkParamValueAttachments(paramValues)

		if _, found := paramValues["Variable"]; found {
			var paramVariable = paramValues["Variable"].(map[string]interface{})
			checkParamValueAttachments(paramVariable)
		}
	}
}

// checkParamValueAttachments checks for attachments on values to map them.
func checkParamValueAttachments(param map[string]interface{}) {
	if param["Value"] == nil {
		return
	}

	var paramValue = param["Value"].(map[string]interface{})

	var inputValue Value
	mapToStruct(paramValue, &inputValue)
	mapValueReference(inputValue)

	if paramValue["attachmentsByRange"] != nil {
		mapAttachmentIdentifiers(paramValue["attachmentsByRange"].(map[string]interface{}))
	}
}

// mapAttachmentIdentifiers maps the UUID and output name of an attachment value.
func mapAttachmentIdentifiers(attachments map[string]interface{}) {
	for _, attachment := range attachments {
		var attachmentValue Value
		mapToStruct(attachment, &attachmentValue)
		mapValueReference(attachmentValue)
	}
}

func mapValueReference(value Value) {
	if value.OutputName == "" || value.OutputUUID == "" {
		return
	}
	mapUUID(value.OutputUUID, value.OutputName)
}

func mapUUID(uuid string, varName string) {
	var outputName string
	if _, found := uuids[uuid]; !found {
		outputName = varName
		sanitizeIdentifier(&outputName)
		uuids[uuid] = checkDuplicateOutputName(outputName)
	}
}

// sanitizeIdentifier strips special characters and replaces dashes with underscores.
func sanitizeIdentifier(identifier *string) {
	*identifier = strings.ReplaceAll(*identifier, "-", "_")
	*identifier = specialCharsRegex.ReplaceAllString(*identifier, "")
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

// newCodeLine writes s on a new line in the code.
func newCodeLine(s string) {
	if tabLevel > 0 {
		for i := 0; i < tabLevel; i++ {
			code.WriteRune('\t')
		}
	}
	code.WriteString(s)
}

// tabbedLine returns s prepended with tab characters at the current tabLevel.
func tabbedLine(s string) string {
	if tabLevel < 1 {
		return s
	}
	var str strings.Builder
	for i := 0; i < tabLevel; i++ {
		str.WriteRune('\t')
	}
	str.WriteString(s)

	return str.String()
}

// defineName writes the name of the imported Shortcut to code as a name definition if the basename contains a space.
func defineName() {
	if strings.Contains(basename, " ") {
		newCodeLine(fmt.Sprintf("#define name %s\n", basename))
	}
}

func decompileIcon() {
	var hasDefinitions bool
	var icon = shortcut.WFWorkflowIcon
	if icon.WFWorkflowIconStartColor != iconColor {
		for name, i := range colors {
			if icon.WFWorkflowIconStartColor != i {
				continue
			}

			newCodeLine(fmt.Sprintf("#define color %s\n", name))
			hasDefinitions = true
		}
	}

	if icon.WFWorkflowIconGlyphNumber != iconGlyph {
		for name, i := range glyphs {
			if icon.WFWorkflowIconGlyphNumber != int64(i) {
				continue
			}

			newCodeLine(fmt.Sprintf("#define glyph %s\n", name))
			hasDefinitions = true
		}
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
			var nonLiteral = decompNumberValue(&action)
			if nonLiteral {
				decompAction(&action)
			}
		case "is.workflow.actions.dictionary":
			decompDictionary(&action)
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

var varUUIDs []string

// checkConstantLiteral determines if action is a constant literal and if it should be
// written out on a new line as a constant and clear the current variable value.
func checkConstantLiteral(action *ShortcutAction) {
	if _, found := action.WFWorkflowActionParameters[UUID]; !found {
		return
	}
	var actionUUID = action.WFWorkflowActionParameters[UUID].(string)
	if slices.Contains(varUUIDs, actionUUID) {
		return
	}
	if outputName, found := uuids[actionUUID]; found {
		newCodeLine(fmt.Sprintf("const %s = ", outputName))
		code.WriteString(currentVariableValue)
		code.WriteRune('\n')
		currentVariableValue = ""
	}
}

func decompComment(action *ShortcutAction) {
	var commentText = action.WFWorkflowActionParameters["WFCommentActionText"].(string)
	if strings.Contains(commentText, "\n") {
		newCodeLine(fmt.Sprintf("/*\n%s\n*/\n\n", commentText))
	} else {
		newCodeLine(fmt.Sprintf("// %s\n\n", commentText))
	}
}

func decompTextValue(action *ShortcutAction) {
	currentVariableValue = decompValue(action.WFWorkflowActionParameters["WFTextActionText"])
	if currentVariableValue == "" {
		currentVariableValue = "\"\""
	}
	checkConstantLiteral(action)
}

func decompNumberValue(action *ShortcutAction) (nonLiteral bool) {
	var value = action.WFWorkflowActionParameters["WFNumberActionNumber"]
	if reflect.TypeOf(value).String() == dictType {
		var mapValue = value.(map[string]interface{})
		var Value = mapValue["Value"].(map[string]interface{})
		value = decompValue(value)

		match, _ := regexp.MatchString("[0-9.]", value.(string))
		if !match {
			nonLiteral = true
			return
		}

		if Value["Type"] == "ActionOutput" {
			currentVariableValue = value.(string)
		}
		return
	}

	var number int
	if value != "" {
		if reflect.TypeOf(value).String() == stringType {
			var convErr error
			number, convErr = strconv.Atoi(value.(string))
			handle(convErr)
		} else {
			number = int(value.(float64))
		}
	}
	currentVariableValue = decompValue(number)
	checkConstantLiteral(action)
	return
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
	sanitizeIdentifier(&variableName)
	newCodeLine(fmt.Sprintf("@%s", variableName))

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

var controlFlowUUIDs []string

func collectControlFlowUUID(action *ShortcutAction) {
	if action.WFWorkflowActionParameters["UUID"] != nil {
		var uuid = action.WFWorkflowActionParameters["UUID"].(string)
		controlFlowUUIDs = append(controlFlowUUIDs, uuid)
	}
}

func decompMenu(action *ShortcutAction) {
	if len(menus) == 0 {
		menus = make(map[string][]variableValue)
	}
	var controlFlowMode = action.WFWorkflowActionParameters["WFControlFlowMode"].(uint64)
	var groupingUUID = action.WFWorkflowActionParameters["GroupingIdentifier"].(string)
	switch controlFlowMode {
	case startStatement:
		menus[groupingUUID] = []variableValue{}
		var items = action.WFWorkflowActionParameters["WFMenuItems"]
		for _, item := range items.([]interface{}) {
			menus[groupingUUID] = append(menus[groupingUUID], variableValue{value: item})
		}
		newCodeLine("menu ")
		code.WriteString(decompValue(action.WFWorkflowActionParameters["WFMenuPrompt"]))
		code.WriteString(" {\n")
		tabLevel++
	case statementPart:
		tabLevel--
		newCodeLine("item ")
		if _, found := action.WFWorkflowActionParameters["WFMenuItemAttributedTitle"]; found {
			code.WriteString(decompValue(action.WFWorkflowActionParameters["WFMenuItemAttributedTitle"]))
		} else if menus[groupingUUID] != nil {
			var menuItem = menus[groupingUUID][0]
			code.WriteString(decompValue(menuItem.value))
			var menu = menus[groupingUUID]
			menus[groupingUUID] = append(menu[:0], menu[1:]...)
		}
		code.WriteString(":\n")
		tabLevel++
	case endStatement:
		collectControlFlowUUID(action)
		tabLevel--
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
		collectControlFlowUUID(action)
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
		collectControlFlowUUID(action)
		tabLevel--
		newCodeLine("}\n")
	}
}

func decompConditional(action *ShortcutAction) {
	var controlFlowMode = action.WFWorkflowActionParameters["WFControlFlowMode"].(uint64)
	switch controlFlowMode {
	case startStatement:
		var conditionInt = action.WFWorkflowActionParameters["WFCondition"].(uint64)
		var conditionString = strconv.FormatUint(conditionInt, 10)
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
					if reflect.TypeOf(action.WFWorkflowActionParameters["WFNumberValue"]).String() == stringType {
						var numberValue, convErr = strconv.Atoi(action.WFWorkflowActionParameters["WFNumberValue"].(string))
						handle(convErr)
						code.WriteString(decompValue(numberValue))
					} else {
						code.WriteString(decompValue(action.WFWorkflowActionParameters["WFNumberValue"]))
					}
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
		collectControlFlowUUID(action)
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

var decompilingDictinary = false

func decompDictionary(action *ShortcutAction) {
	decompilingDictinary = true

	var params DictionaryActionParameters
	mapToStruct(action.WFWorkflowActionParameters, &params)

	var dictionary = decompDictionaryItems(params.WFItems.Value.WFDictionaryFieldValueItems)
	var jsonBytes, jsonErr = json.MarshalIndent(dictionary, strings.Repeat("\t", tabLevel), "\t")
	handle(jsonErr)

	currentVariableValue = string(jsonBytes)

	checkConstantLiteral(action)

	decompilingDictinary = false
}

func decompDictionaryItems(items []WFDictionaryFieldValueItem) (dictionary map[string]interface{}) {
	dictionary = make(map[string]interface{})
	for _, item := range items {
		var itemKey = decompValue(item.WFKey)
		var itemStringValue = decompValue(item.WFValue.Value)
		var itemValueType = item.WFItemType
		var itemValue any
		switch dictDataType(itemValueType) {
		case itemTypeNumber:
			itemStringValue = item.WFValue.Value.(map[string]interface{})["string"].(string)
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
			itemValue = decompArray(item.WFValue.Value.([]interface{}))
		case itemTypeDict:
			var dictionaryItems []WFDictionaryFieldValueItem
			var value = item.WFValue.Value.(map[string]interface{})
			mapToStruct(value["Value"].(map[string]interface{})["WFDictionaryFieldValueItems"], &dictionaryItems)
			itemValue = decompDictionaryItems(dictionaryItems)
		default:
			itemValue = itemStringValue
		}
		dictionary[itemKey] = itemValue
	}
	return
}

func decompArray(items []interface{}) (array []interface{}) {
	for _, item := range items {
		var arrayItem = item.(map[string]interface{})
		var itemType = dictDataType(int(arrayItem["WFItemType"].(float64)))
		var itemValue interface{}
		switch itemType {
		case itemTypeText:
			itemValue = decompValue(arrayItem["WFValue"].(map[string]interface{}))
		case itemTypeNumber:
			var convErr error
			itemValue, convErr = strconv.ParseInt(decompValue(arrayItem["WFValue"].(map[string]interface{})), 10, 64)
			handle(convErr)
		case itemTypeArray:
			var wfValue = arrayItem["WFValue"].(map[string]interface{})
			array = append(array, decompArray(wfValue["Value"].([]interface{})))
		case itemTypeDict:
			var dictionaryItems []WFDictionaryFieldValueItem
			var wfValue = arrayItem["WFValue"].(map[string]interface{})
			var rootValue = wfValue["Value"].(map[string]interface{})
			var value = rootValue["Value"].(map[string]interface{})
			for _, item := range value["WFDictionaryFieldValueItems"].([]interface{}) {
				var dictionaryItem = item.(map[string]interface{})
				var itemValue = dictionaryItem["WFValue"].(map[string]interface{})
				var value = itemValue["Value"].(map[string]interface{})
				dictionaryItems = append(dictionaryItems, WFDictionaryFieldValueItem{
					WFKey:      dictionaryItem["WFKey"],
					WFItemType: int(dictionaryItem["WFItemType"].(float64)),
					WFValue: WFValue{
						String: value["string"].(string),
					},
				})
			}
			itemValue = decompDictionaryItems(dictionaryItems)
		case itemTypeBool:
			var wfValue = arrayItem["WFValue"].(map[string]interface{})
			itemValue = wfValue["Value"].(bool)
		}
		if itemValue != nil {
			array = append(array, itemValue)
		}
	}
	return
}

func decompValue(value any) string {
	if value == nil {
		return "nil"
	}
	var valueType = reflect.TypeOf(value).String()
	switch valueType {
	case dictType:
		return decompValueObject(value.(map[string]interface{}))
	case stringType:
		return fmt.Sprintf("\"%s\"", escapeString(value.(string)))
	default:
		return fmt.Sprintf("%v", value)
	}
}

func escapeString(value string) string {
	var escapes = map[string]string{
		"\n": "\\n",
		"\t": "\\t",
		"\r": "\\r",
		"\"": "\\\"",
	}
	for chr, e := range escapes {
		value = strings.ReplaceAll(value, chr, e)
	}

	return value
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
			sanitizeIdentifier(&variableName)

			return variableName
		}

		var variableValue = value["Variable"].(map[string]interface{})
		return decompValue(variableValue["Value"])
	case "ActionOutput":
		if _, found := value["OutputUUID"]; found {
			var outputName = uuids[value["OutputUUID"].(string)]
			sanitizeIdentifier(&outputName)

			if isControlFlowUUID(value["OutputUUID"].(string)) {
				decompWarning("Usage of control flow action output is not supported. This can be manually corrected by assigning a variable within the control flow branches and then using that variable instead.")
			}

			return outputName
		}
	case globals[ShortcutInput].variableType:
		return ShortcutInput
	}

	return decompObjectValue(value)
}

func isControlFlowUUID(uuid string) bool {
	return slices.Contains(controlFlowUUIDs, uuid)
}

func decompObjectValue(valueObj any) string {
	var valueType = reflect.TypeOf(valueObj).String()
	if valueType != dictType {
		return fmt.Sprintf("%v", valueObj)
	}

	var value = valueObj.(map[string]interface{})

	var attachmentString string
	if value["value"] != nil {
		if reflect.TypeOf(value["value"]).String() != dictType {
			return fmt.Sprintf("%v", valueObj)
		}
		value = value["value"].(map[string]interface{})
	}

	if _, found := value["string"]; found {
		attachmentString = value["string"].(string)
	}

	if attachments, found := value["attachmentsByRange"]; found {
		decompAttachmentString(&attachmentString, attachments.(map[string]interface{}))
	}

	return attachmentString
}

func decompAttachmentString(attachmentString *string, attachments map[string]interface{}) {
	var originalString = *attachmentString
	var attachmentChars = strings.Split(*attachmentString, "")

	for attachmentRange, a := range attachments {
		var attachmentRanges = strings.Split(attachmentRange, ",")
		var attachmentPosition = strings.TrimPrefix(attachmentRanges[0], "{")
		var position, convErr = strconv.Atoi(attachmentPosition)
		handle(convErr)

		var attachment Value
		mapToStruct(a, &attachment)

		var variableName = attachment.VariableName
		if _, uuid := uuids[attachment.OutputUUID]; uuid {
			variableName = uuids[attachment.OutputUUID]
		}
		sanitizeIdentifier(&variableName)
		if variableName == "" && attachment.Type == globals[ShortcutInput].variableType {
			variableName = ShortcutInput
		}

		if len(attachment.Aggrandizements) != 0 {
			decompAggrandizements(&variableName, attachment.Aggrandizements)
		}

		if isControlFlowUUID(attachment.OutputUUID) {
			decompWarning("Usage of control flow action output is not supported. This can be manually corrected by assigning a variable within the control flow branches and then using that variable instead.")
		}

		attachmentChars[position] = fmt.Sprintf("{%s}", variableName)
	}

	*attachmentString = escapeString(strings.Join(attachmentChars, ""))

	if !decompilingDictinary && originalString == ObjectReplaceCharStr && len(attachments) == 1 {
		*attachmentString = strings.Trim(*attachmentString, "{}")
	} else {
		*attachmentString = fmt.Sprintf("\"%s\"", *attachmentString)
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
	if skipDecompAction(action) {
		return
	}

	var actionCallCode strings.Builder
	var matchedIdentifier, matchedAction = matchAction(action)
	if matchedIdentifier == "" {
		actionCallCode.WriteString(makeRawAction(action))
	}

	if matchedAction.mac && !macDefinition {
		var saveCode = code.String()
		code.Reset()
		code.WriteString(fmt.Sprintf("#define mac true\n%s", saveCode))
		macDefinition = true
	}

	var isConstant, isVariableValue = checkOutputType(action)

	if matchedIdentifier != "" {
		var actionCallStart = fmt.Sprintf("%s(", matchedIdentifier)
		if !isConstant && !isVariableValue {
			newCodeLine(actionCallStart)
		} else {
			actionCallCode.WriteString(actionCallStart)
		}

		if matchedAction.make != nil || matchedAction.decomp != nil {
			decompActionCustom(&actionCallCode, &matchedAction, action)
		} else {
			var matchedParamsSize = len(matchedAction.parameters)
			if matchedParamsSize > 0 {
				decompActionArguments(&actionCallCode, &matchedAction, action)
			}
		}

		actionCallCode.WriteString(")")
	} else if !isConstant && !isVariableValue {
		var saveCode = tabbedLine(actionCallCode.String())
		actionCallCode.Reset()
		actionCallCode.WriteString(saveCode)
	}

	if isVariableValue {
		currentVariableValue = actionCallCode.String()
	} else {
		code.WriteString(actionCallCode.String())
		code.WriteRune('\n')
		currentVariableValue = ""
	}
	actionIndex++
}

// checkOutputType determines if action output is a constant or a variable.
// If it is a constant we will write a constant statement on a new line to prepend the action.
func checkOutputType(action *ShortcutAction) (isConstant bool, isVariableValue bool) {
	if action.WFWorkflowActionParameters["CustomOutputName"] != nil {
		var customOutputName = action.WFWorkflowActionParameters["CustomOutputName"].(string)
		sanitizeIdentifier(&customOutputName)
		if ref, foundVar := variables[customOutputName]; foundVar {
			isConstant = ref.constant
			isVariableValue = !ref.constant
		}
	}

	if action.WFWorkflowActionParameters[UUID] != nil && !isVariableValue {
		var uuid = action.WFWorkflowActionParameters[UUID].(string)
		if _, found := uuids[uuid]; found {
			isConstant = !slices.Contains(varUUIDs, uuid)
			isVariableValue = slices.Contains(varUUIDs, uuid)
			if isConstant {
				newCodeLine(fmt.Sprintf("const %s = ", uuids[uuid]))
			}
		}
	}
	return
}

// skipDecompAction skips actions we don't support or when necessary.
func skipDecompAction(action *ShortcutAction) bool {
	if action.WFWorkflowActionIdentifier == "is.workflow.actions.getvariable" {
		var varName = decompValue(action.WFWorkflowActionParameters["WFVariable"])

		insertCodeComment(fmt.Sprintf("TODO: Get Variable not supported: Assign variable here to '%s'.", varName))
		decompWarning(fmt.Sprintf("Get variable '%s' is not supported. Set a variable to that value instead if something was depending on it's output.", varName))

		actionIndex++
		return true
	}

	if action.WFWorkflowActionIdentifier == "is.workflow.actions.nothing" {
		var nextAction = peekActions(1)
		var controlflowActionIdentifiers = []string{"is.workflow.actions.conditional", "is.workflow.actions.repeat.each", "is.workflow.actions.repeat.count", "is.workflow.acitons.choosefrommenu"}
		if slices.Contains(controlflowActionIdentifiers, nextAction.WFWorkflowActionIdentifier) {
			var controlFlowMode = nextAction.WFWorkflowActionParameters["WFControlFlowMode"]
			if controlFlowMode == endStatement || controlFlowMode == statementPart {
				return true
			}
		}
	}

	return false
}

func decompActionArguments(actionCallCode *strings.Builder, matchedAction *actionDefinition, action *ShortcutAction) {
	for i, param := range matchedAction.parameters {
		if param.key == "" {
			continue
		}

		var argValue string
		if value, found := action.WFWorkflowActionParameters[param.key]; found {
			argValue = decompValue(value)
		} else if !param.optional {
			argValue = makeDefaultValue(param)
		}

		switch param.validType {
		case Integer:
			argValue = strings.Trim(argValue, "\"")
		}

		if argValue != "" {
			if i == 0 {
				actionCallCode.WriteString(argValue)
			} else {
				actionCallCode.WriteString(fmt.Sprintf(",%s", argValue))
			}
		}
	}
}

func makeDefaultValue(param parameterDefinition) string {
	if param.defaultValue != nil {
		if reflect.TypeOf(param.defaultValue).Kind() == reflect.String {
			return fmt.Sprintf("\"%s\"", param.defaultValue)
		}

		return fmt.Sprintf("%v", param.defaultValue)
	}

	switch param.validType {
	case Integer:
		return "0"
	case String:
		return "\"\""
	case Arr:
		return "[]"
	case Dict:
		return "{}"
	case Bool:
		return "false"
	}

	return "nil"
}

func makeRawAction(action *ShortcutAction) string {
	var rawActionCode strings.Builder
	rawActionCode.WriteString(fmt.Sprintf("rawAction(\"%s\"", action.WFWorkflowActionIdentifier))
	var paramsSize = len(action.WFWorkflowActionParameters)
	var onlyUUIDParam = false
	var onlyOutputNameParam = false
	if paramsSize > 1 {
		if _, found := action.WFWorkflowActionParameters[UUID]; found {
			onlyUUIDParam = found
		}
		if _, found := action.WFWorkflowActionParameters["CustomOutputName"]; found {
			onlyOutputNameParam = found
		}
	}
	if paramsSize > 1 && (!onlyUUIDParam && !onlyOutputNameParam) {
		tabLevel++
		rawActionCode.WriteString(", [\n")
		rawActionCode.WriteString(tabbedLine("{\n"))
		var index = 0
		for key, param := range action.WFWorkflowActionParameters {
			index++

			if key == UUID || key == "CustomOutputName" {
				continue
			}

			rawActionCode.WriteString(strings.Repeat("\t", tabLevel+1))
			rawActionCode.WriteString(fmt.Sprintf("\"%s\": ", key))

			var value = decompValue(param)
			if !strings.Contains(value, "\"") {
				value = fmt.Sprintf("\"{%s}\"", value)
			}
			rawActionCode.WriteString(value)

			if index < paramsSize {
				rawActionCode.WriteRune(',')
			}
			rawActionCode.WriteRune('\n')
		}

		rawActionCode.WriteString(tabbedLine("}\n"))
		tabLevel--
		rawActionCode.WriteString(tabbedLine("])"))
	} else {
		rawActionCode.WriteRune(')')
	}

	return rawActionCode.String()
}

func matchAction(action *ShortcutAction) (name string, definition actionDefinition) {
	for call, def := range actions {
		var identifier = strings.ToLower(call)
		if def.identifier != "" {
			identifier = def.identifier
		}
		var appIdentifier = "is.workflow.actions"
		if def.appIdentifier != "" {
			appIdentifier = def.appIdentifier
		}
		var actionIdentifier = fmt.Sprintf("%s.%s", appIdentifier, identifier)
		if actionIdentifier == action.WFWorkflowActionIdentifier || definition.overrideIdentifier == action.WFWorkflowActionIdentifier {
			definition = *def
			name = call

			if splitActions, found := identifierMap[identifier]; found {
				matchSplitAction(&splitActions, action.WFWorkflowActionParameters, &name, &definition)
				if name != "run" && name != "runSelf" {
					break
				}
				var workflow = action.WFWorkflowActionParameters["WFWorkflow"].(map[string]interface{})
				if _, isSelf := workflow["isSelf"]; !isSelf {
					break
				}

				if workflow["isSelf"].(bool) || action.WFWorkflowActionParameters["WFWorkflowName"] == basename {
					name = "runSelf"
					definition = *actions["runSelf"]
				} else {
					name = "run"
				}
			}
			break
		}
	}
	return
}

type actionMatch struct {
	params uint
	values uint
	action actionValue
}

var matches []actionMatch

func matchSplitAction(splitActions *[]actionValue, parameters map[string]any, identifier *string, definition *actionDefinition) {
	if args.Using("debug") {
		fmt.Println("## MATCHING SPLIT ACTIONS ##")
		fmt.Println("parameters", parameters)
	}
	matches = []actionMatch{}

	var defaultAction, hasDefaultAction = getDefaultAction(splitActions)
	if hasDefaultAction {
		*identifier = defaultAction.identifier
		*definition = *defaultAction.definition
		if args.Using("debug") {
			fmt.Println("has default action", defaultAction.identifier)
		}
	}

	for _, splitAction := range *splitActions {
		var matchedParams, matchedValues = scoreActionMatch(splitAction, splitAction.definition.parameters, parameters)

		if matchedParams == 0 {
			continue
		}

		matches = append(matches, actionMatch{
			params: matchedParams,
			values: matchedValues,
			action: splitAction,
		})
	}

	if len(matches) == 0 {
		return
	}

	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].values > matches[j].values
	})

	if args.Using("debug") {
		for _, match := range matches[1:] {
			fmt.Printf("%s()\n", match.action.identifier)
			fmt.Println("params:", match.params, ", values:", match.values)
			fmt.Println("---")
		}
		fmt.Print("\n\n")
	}

	if !competingMatches(matches) {
		return
	}

	var matchedAction = matches[0]
	*identifier = matchedAction.action.identifier
	*definition = *matchedAction.action.definition
}

// competingMatches determines if the matches for this identifier have more values than 1 matching this action.
func competingMatches(matches []actionMatch) bool {
	var matchedValues int
	for _, match := range matches {
		if match.values > 0 {
			matchedValues++
		}
	}

	return matchedValues > 0
}

func scoreActionMatch(splitAction actionValue, splitActionParams []parameterDefinition, parameters map[string]any) (matchedParams uint, matchedValues uint) {
	var paramMatches, valueMatches = scoreActionParams(&splitActionParams, parameters)
	matchedParams += paramMatches
	matchedValues += valueMatches

	var splitActionAddParams []parameterDefinition
	if splitAction.definition.addParams != nil {
		for _, addParam := range splitAction.definition.addParams([]actionArgument{}) {
			splitActionAddParams = append(splitActionAddParams, parameterDefinition{
				key:          addParam.key,
				defaultValue: addParam.value,
			})
		}
	}

	var addParamMatches, addValueMatches = scoreActionAddParams(&splitActionAddParams, parameters)
	matchedParams += addParamMatches
	matchedValues += addValueMatches

	return
}

func scoreActionParams(splitActionParams *[]parameterDefinition, parameters map[string]any) (matchedParams uint, matchedValues uint) {
	for _, param := range *splitActionParams {
		if param.key == "" {
			continue
		}
		if value, found := parameters[param.key]; found {
			matchedParams++
			if len(param.enum) > 0 && slices.Contains(param.enum, fmt.Sprintf("%s", value)) {
				matchedValues++
			}
			if param.defaultValue != nil {
				var defaultValue = fmt.Sprintf("%v", param.defaultValue)
				var rawValue = strings.Trim(decompValue(value), "\"")
				if defaultValue == rawValue && defaultValue != "" && rawValue != "" {
					matchedValues++
				}
			}
		}
	}
	return
}

func scoreActionAddParams(splitActionAddParams *[]parameterDefinition, parameters map[string]any) (matchedParams uint, matchedValues uint) {
	for _, param := range *splitActionAddParams {
		if param.key == "" {
			continue
		}
		if value, found := parameters[param.key]; found {
			matchedParams++

			var defaultValueType = reflect.TypeOf(param.defaultValue).String()
			var valueType = reflect.TypeOf(value).String()
			if defaultValueType == dictType && valueType == dictType {
				if maps.Equal(param.defaultValue.(map[string]interface{}), value.(map[string]interface{})) {
					matchedValues++
				}
			} else if param.defaultValue == value {
				matchedValues++
			}
		}
	}
	return
}

func decompActionCustom(actionCode *strings.Builder, matchedAction *actionDefinition, action *ShortcutAction) {
	var arguments = matchedAction.decomp(action)
	if len(arguments) == 0 {
		return
	}
	actionCode.WriteString(strings.Join(arguments, ","))
}

func glueToChar(glue string) string {
	switch glue {
	case "New Lines":
		return "\n"
	case "Spaces":
		return " "
	case "Every Character":
		return ""
	default:
		return glue
	}
}

// getDefaultAction gets the default action from a slice of split actions.
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

	fmt.Println("### VARIABLES ###")
	printVariables()
	fmt.Print("\n")

	fmt.Println("### VARIABLE UUIDs ###")
	fmt.Println(varUUIDs)
	fmt.Print("\n")

	fmt.Println("### UUIDS ###")
	for uuid, name := range uuids {
		fmt.Printf("%s | %s\n", uuid, name)
	}
	fmt.Print("\n")
}

func decompWarning(message string) {
	var linesLen = strings.Count(code.String(), "\n")

	var filePath string
	if args.Using("no-ansi") {
		filePath = basename + ".cherri"
	} else {
		filePath = getOutputPath(relativePath + basename + ".cherri")
	}

	fmt.Println(ansi("Warning:", yellow, bold), fmt.Sprintf("%s (%s:%d:0)\n", message, filePath, linesLen+1))
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
