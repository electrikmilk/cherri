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
	"howett.net/plist"
)

const (
	UUID          = "UUID"
	ShortcutInput = "ShortcutInput"
	Ask           = "Ask"
)

var tabLevel int

var code strings.Builder
var specialCharsRegex *regexp.Regexp

func decompile(b []byte) {
	var _, marshalIndexedErr = plist.Unmarshal(b, &shortcut)
	handle(marshalIndexedErr)

	specialCharsRegex = regexp.MustCompile("[^a-zA-Z0-9_]+")
	variables = make(map[string]varValue)
	uuids = make(map[string]string)
	controlFlowGroups = make(map[int]controlFlowGroup)

	basename = strings.ReplaceAll(basename, " ", "_")
	outputPath = getOutputPath(basename + ".cherri")
	if args.Using("no-ansi") {
		filePath = basename + ".cherri"
	} else {
		filePath = outputPath
	}

	loadBasicStandardActions()
	resetParse()
	firstChar()

	mapVariables()
	mapSplitActions()
	mapIdentifiers()
	waitFor(
		mapControlFlowOutputs,
		defineToggleSetActions,
	)

	defineName()
	decompileIcon()

	decompileActions()

	if args.Using("debug") {
		printDecompDebug()
	}

	var writeErr = os.WriteFile(outputPath, []byte(code.String()), 0600)
	handle(writeErr)
}

// mapIdentifiers creates a map of variable identifiers and UUIDs that are assigned throughout the Shortcut.
func mapIdentifiers() {
	for _, action := range shortcut.WFWorkflowActions {
		currentActionIdentifier = action.WFWorkflowActionIdentifier
		var params = action.WFWorkflowActionParameters
		if action.WFWorkflowActionIdentifier == SetVariableIdentifier || action.WFWorkflowActionIdentifier == AppendVariableIdentifier {
			continue
		}

		if params["UUID"] != nil && params["CustomOutputName"] != nil {
			mapUUID(params["UUID"].(string), params["CustomOutputName"].(string))
		}

		checkParamIdentifiers(params)
	}
}

func mapControlFlowOutputs() {
	for _, action := range shortcut.WFWorkflowActions {
		switch action.WFWorkflowActionIdentifier {
		case "is.workflow.actions.conditional":
			fallthrough
		case "is.workflow.actions.choosefrommenu":
			fallthrough
		case "is.workflow.actions.repeat.count":
			fallthrough
		case "is.workflow.actions.repeat.each":
			var params = action.WFWorkflowActionParameters
			if params["WFControlFlowMode"].(uint64) != endStatement {
				continue
			}

			var groupName string
			if params["UUID"] != nil {
				var endingUUID = params["UUID"].(string)
				if uuids[endingUUID] != "" {
					groupName = uuids[endingUUID]
				}
			}

			controlFlowGroups[groupingIdx] = controlFlowGroup{
				groupType:  Conditional,
				identifier: groupName,
			}
			groupingIdx++
		}
	}
	groupingIdx = 0
}

// Map out variables in the Shortcut and their UUIDs for later checks.
func mapVariables() {
	for _, action := range shortcut.WFWorkflowActions {
		currentActionIdentifier = action.WFWorkflowActionIdentifier
		var params = action.WFWorkflowActionParameters
		if action.WFWorkflowActionIdentifier == SetVariableIdentifier || action.WFWorkflowActionIdentifier == AppendVariableIdentifier {
			var varName = params["WFVariableName"].(string)
			sanitizeIdentifier(&varName)
			if _, found := variables[varName]; !found {
				variables[varName] = varValue{}
			}

			if action.WFWorkflowActionParameters["WFInput"] != nil {
				var wfInput WFInput
				mapToStruct(action.WFWorkflowActionParameters["WFInput"], &wfInput)
				if _, found := uuids[wfInput.Value.OutputUUID]; !found {
					varUUIDs = append(varUUIDs, wfInput.Value.OutputUUID)
				}
			}
		}
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
		if value == nil || reflect.TypeOf(value).Kind() != reflect.Map {
			continue
		}

		var paramValues = value.(map[string]interface{})
		checkParamValueAttachments(paramValues)

		if _, found := paramValues["Variable"]; found {
			var paramVariable = paramValues["Variable"].(map[string]interface{})
			checkParamValueAttachments(paramVariable)
		}

		if _, found := paramValues["WFConditions"]; found {
			var wfConditions = paramValues["WFConditions"].(map[string]interface{})
			var value = wfConditions["Value"].(map[string]interface{})
			if value["WFActionParameterFilterTemplates"] != nil {
				for _, filtertemplate := range value["WFActionParameterFilterTemplates"].(map[string]interface{}) {
					var paramFilterTemplate = filtertemplate.(map[string]interface{})
					checkParamIdentifiers(paramFilterTemplate["WFInput"].(map[string]interface{}))
				}
			}
		}
	}
}

// checkParamValueAttachments checks for attachments on values to map them.
func checkParamValueAttachments(param map[string]interface{}) {
	if param["Value"] == nil || reflect.TypeOf(param["Value"]).Kind() != reflect.Map {
		return
	}

	var paramValue = param["Value"].(map[string]interface{})

	var inputValue Value
	mapToStruct(paramValue, &inputValue)
	mapValueReference(inputValue)
	checkConstantUUID(inputValue)

	if inputValue.AttachmentsByRange != nil {
		mapAttachmentIdentifiers(inputValue.AttachmentsByRange)
	}
	if inputValue.WFDictionaryFieldValueItems != nil {
		mapDictionaryValueIdentifiers(inputValue.WFDictionaryFieldValueItems)
	}
}

func mapDictionaryValueIdentifiers(items []WFDictionaryFieldValueItem) {
	for _, item := range items {
		if item.WFValue == nil || reflect.TypeOf(item.WFValue).Kind() != reflect.Map {
			continue
		}
		var wfValue = item.WFValue.(map[string]interface{})
		if wfValue["Value"] == nil {
			continue
		}

		checkParamValueAttachments(wfValue)

		var valueKind = reflect.TypeOf(wfValue["Value"]).Kind()
		if valueKind == reflect.Map {
			var valueTwo = wfValue["Value"].(map[string]interface{})
			if valueTwo != nil && valueTwo["Value"] != nil && reflect.TypeOf(valueTwo["Value"]).Kind() == reflect.Map {
				var itemsValue = valueTwo["Value"].(map[string]interface{})

				var dictionaryItems []WFDictionaryFieldValueItem
				mapToStruct(itemsValue["WFDictionaryFieldValueItems"], &dictionaryItems)
				mapDictionaryValueIdentifiers(dictionaryItems)
			}
		} else if valueKind == reflect.Slice {
			var dictionaryItems = mapDictionaryItems(wfValue["Value"].([]interface{}))
			mapDictionaryValueIdentifiers(dictionaryItems)
		}
	}
}

func mapDictionaryItems(items []interface{}) (dictionaryItems []WFDictionaryFieldValueItem) {
	for _, value := range items {
		if reflect.TypeOf(value).Kind() != reflect.Map {
			continue
		}

		var dictionaryItem WFDictionaryFieldValueItem
		mapToStruct(value, &dictionaryItem)
		dictionaryItems = append(dictionaryItems, dictionaryItem)
	}
	return
}

// mapAttachmentIdentifiers maps the UUID and output name of an attachment value.
func mapAttachmentIdentifiers(attachments map[string]Value) {
	for _, attachmentValue := range attachments {
		mapValueReference(attachmentValue)
		checkConstantUUID(attachmentValue)
	}
}

func checkConstantUUID(value Value) {
	if value.OutputUUID != "" && value.Type == "ActionOutput" {
		constUUIDs = append(constUUIDs, value.OutputUUID)
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
	if strings.Contains(*identifier, " ") {
		var words = strings.Split(strings.TrimSpace(*identifier), " ")
		for i, word := range words {
			words[i] = capitalize(word)
		}
		*identifier = strings.Join(words, "")
	}

	*identifier = specialCharsRegex.ReplaceAllString(*identifier, "")
	*identifier = strings.ReplaceAll(*identifier, "-", "_")
}

type actionValue struct {
	identifier string
	definition *actionDefinition
}

var identifierMap map[string][]actionValue

// mapSplitActions creates a map of actions that have been split into a few actions to reduce the number of arguments.
func mapSplitActions() {
	if identifierMap == nil {
		identifierMap = make(map[string][]actionValue)
	}
	for identifier, action := range actions {
		var ident = action.identifier
		if action.identifier == "" {
			ident = identifier
		}

		ident = strings.ToLower(ident)

		var splitActionValue = actionValue{
			identifier: identifier,
			definition: action,
		}
		if identifierMap[ident] != nil && slices.Contains(identifierMap[ident], splitActionValue) {
			continue
		}

		identifierMap[ident] = append(identifierMap[ident], splitActionValue)
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
	var icon = shortcut.WFWorkflowIcon
	var hasIcon = icon.WFWorkflowIconStartColor != iconColor || icon.WFWorkflowIconGlyphNumber != iconGlyph

	if icon.WFWorkflowIconStartColor != iconColor {
		defineColors(icon, colors)
		defineColors(icon, altColors)
	}

	if icon.WFWorkflowIconGlyphNumber != iconGlyph {
		for name, i := range glyphs {
			if icon.WFWorkflowIconGlyphNumber != int64(i) {
				continue
			}

			newCodeLine(fmt.Sprintf("#define glyph %s\n", name))
		}
	}

	if hasIcon {
		newCodeLine("\n")
	}
}

func defineColors(icon ShortcutIcon, colors map[string]int) {
	for name, i := range colors {
		if icon.WFWorkflowIconStartColor != i {
			continue
		}

		newCodeLine(fmt.Sprintf("#define color %s\n", name))
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
		case "is.workflow.actions.math":
			decompBasicExpression(&action)
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
		case "is.workflow.actions.getvalueforkey":
			var dictionaryKey = action.WFWorkflowActionParameters["WFDictionaryKey"]
			if dictionaryKey != nil && action.WFWorkflowActionParameters[UUID] != nil &&
				!slices.Contains(constUUIDs, action.WFWorkflowActionParameters[UUID].(string)) {
				decompDictionaryGetValue(&action)
				continue
			}
			fallthrough
		default:
			decompAction(&action)
		}
		actionIndex++
	}
}

var varUUIDs []string
var constUUIDs []string

// checkConstantLiteral determines if action should be written out on a new line as a constant and clear the current variable value.
func checkConstantLiteral(action *ShortcutAction) {
	if _, found := action.WFWorkflowActionParameters[UUID]; !found {
		return
	}
	var actionUUID = action.WFWorkflowActionParameters[UUID].(string)
	if slices.Contains(varUUIDs, actionUUID) && !slices.Contains(constUUIDs, actionUUID) {
		return
	}
	if outputName, found := uuids[actionUUID]; found {
		newCodeLine(fmt.Sprintf("const %s = ", outputName))
		code.WriteString(currentVariableValue)
		code.WriteRune('\n')
		currentVariableValue = ""
	}
}

func makeConstantLiteral(action *ShortcutAction) {
	if _, found := action.WFWorkflowActionParameters[UUID]; !found {
		return
	}
	var actionUUID = action.WFWorkflowActionParameters[UUID].(string)
	if outputName, found := uuids[actionUUID]; found {
		newCodeLine(fmt.Sprintf("const %s = ", outputName))
		code.WriteString(currentVariableValue)
		code.WriteRune('\n')
		currentVariableValue = ""
	}
}

func decompComment(action *ShortcutAction) {
	var commentText = action.WFWorkflowActionParameters["WFCommentActionText"].(string)
	if args.Using("comments") {
		if strings.Contains(commentText, "\n") {
			newCodeLine(fmt.Sprintf("comment('\n%s\n')\n\n", commentText))
		} else {
			newCodeLine(fmt.Sprintf("comment('%s')\n", commentText))
		}
	}

	if strings.Contains(commentText, "\n") {
		newCodeLine("/*\n")
		for _, line := range strings.Split(commentText, "\n") {
			newCodeLine(fmt.Sprintln(line))
		}
		newCodeLine("*/\n")
	} else {
		newCodeLine(fmt.Sprintf("// %s\n\n", commentText))
	}
}

var decompilingText = false

func decompTextValue(action *ShortcutAction) {
	decompilingText = true
	currentVariableValue = decompValue(action.WFWorkflowActionParameters["WFTextActionText"])
	if currentVariableValue == "" {
		currentVariableValue = "\"\""
	} else {
		currentVariableValue = fmt.Sprintf("\"%s\"", escapeString(strings.Trim(currentVariableValue, "\"")))
	}
	checkConstantLiteral(action)
	decompilingText = false
}

func decompNumberValue(action *ShortcutAction) (nonLiteral bool) {
	var value = action.WFWorkflowActionParameters["WFNumberActionNumber"]
	if reflect.TypeOf(value).Kind() == reflect.Map {
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

	var number any
	var convErr error
	if reflect.TypeOf(value).Kind() == reflect.String {
		if strings.Contains(value.(string), ".") {
			number, convErr = strconv.ParseFloat(value.(string), 64)
		} else if value != "" {
			number, convErr = strconv.Atoi(value.(string))
		}
		handle(convErr)
	} else {
		number = int(value.(uint64))
	}

	currentVariableValue = decompValue(number)
	checkConstantLiteral(action)

	return
}

func decompBasicExpression(action *ShortcutAction) {
	var input = action.WFWorkflowActionParameters["WFInput"]
	var operand = action.WFWorkflowActionParameters["WFMathOperand"]
	var operation = "+"
	if action.WFWorkflowActionParameters["WFMathOperation"] != nil {
		switch action.WFWorkflowActionParameters["WFMathOperation"] {
		case "×":
			operation = "*"
		case "÷":
			operation = "/"
		default:
			operation = action.WFWorkflowActionParameters["WFMathOperation"].(string)
		}
	}
	var expression = strings.Trim(decompValue(input), "\"") + " " + strings.Trim(operation, "\"") + " " + strings.Trim(decompValue(operand), "\"")
	var varRegex = regexp.MustCompile(`{(.*?)}`)
	currentVariableValue = varRegex.ReplaceAllString(expression, "$1")

	checkConstantLiteral(action)
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

func decompDictionaryGetValue(action *ShortcutAction) {
	var dictionaryValueRef strings.Builder
	dictionaryValueRef.WriteString(decompValue(action.WFWorkflowActionParameters["WFInput"]))

	if action.WFWorkflowActionParameters["WFDictionaryKey"] != nil {
		dictionaryValueRef.WriteRune('[')
		if reflect.TypeOf(action.WFWorkflowActionParameters["WFDictionaryKey"]).Kind() == reflect.String {
			dictionaryValueRef.WriteRune('\'')
			dictionaryValueRef.WriteString(action.WFWorkflowActionParameters["WFDictionaryKey"].(string))
			dictionaryValueRef.WriteRune('\'')
		} else {
			dictionaryValueRef.WriteString(decompValue(action.WFWorkflowActionParameters["WFDictionaryKey"]))
		}
		dictionaryValueRef.WriteRune(']')
		currentVariableValue = dictionaryValueRef.String()
		checkConstantLiteral(action)
	}
}

func checkControlFlowOutput(action *ShortcutAction) bool {
	switch action.WFWorkflowActionIdentifier {
	case "is.workflow.actions.conditional":
		fallthrough
	case "is.workflow.actions.choosefrommenu":
		fallthrough
	case "is.workflow.actions.repeat.count":
		fallthrough
	case "is.workflow.actions.repeat.each":
		if action.WFWorkflowActionParameters["WFControlFlowMode"].(uint64) == startStatement {
			var group = controlFlowGroups[groupingIdx]
			if group.identifier != "" {
				newCodeLine(fmt.Sprintf("const %s = ", group.identifier))
				groupingIdx++
				return true
			}
		}
	}
	return false
}

func beginStatement(t tokenType, action *ShortcutAction) {
	if tabLevel == 0 {
		newCodeLine("\n")
	}

	var statementKeyword = fmt.Sprintf("%s ", strings.TrimSpace(string(t)))
	if checkControlFlowOutput(action) {
		code.WriteString(statementKeyword)
	} else {
		newCodeLine(statementKeyword)
	}
}

func decompMenu(action *ShortcutAction) {
	if len(menus) == 0 {
		menus = make(map[string][]varValue)
	}
	var controlFlowMode = action.WFWorkflowActionParameters["WFControlFlowMode"].(uint64)
	var groupingUUID = action.WFWorkflowActionParameters["GroupingIdentifier"].(string)
	switch controlFlowMode {
	case startStatement:
		beginStatement(Menu, action)

		menus[groupingUUID] = []varValue{}
		var items = action.WFWorkflowActionParameters["WFMenuItems"]
		for _, item := range items.([]interface{}) {
			menus[groupingUUID] = append(menus[groupingUUID], varValue{value: item})
		}
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
		tabLevel--
		newCodeLine("}\n")
	}
}

func decompRepeat(action *ShortcutAction) {
	var controlFlowMode = action.WFWorkflowActionParameters["WFControlFlowMode"].(uint64)
	switch controlFlowMode {
	case startStatement:
		beginStatement(Repeat, action)

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
		beginStatement(RepeatWithEach, action)

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
		beginStatement(If, action)

		if action.WFWorkflowActionParameters["WFConditions"] != nil {
			var conditions = action.WFWorkflowActionParameters["WFConditions"].(map[string]interface{})
			var conditionValue = conditions["Value"].(map[string]interface{})
			var intFilterPrefix, convErr = strconv.Atoi(fmt.Sprintf("%d", conditionValue["WFActionParameterFilterPrefix"]))
			handle(convErr)
			var paramFilterPrefix = convertFilterPrefix(intFilterPrefix)
			var filterTemplates = conditionValue["WFActionParameterFilterTemplates"].([]interface{})
			var conditionsLen = len(filterTemplates)
			for i, condition := range filterTemplates {
				decompCondition(condition.(map[string]interface{}), action)

				if i != conditionsLen-1 {
					code.WriteString(fmt.Sprintf(" %s ", paramFilterPrefix))
				}
			}
		} else {
			decompCondition(action.WFWorkflowActionParameters, action)
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

func matchConditionOperator(number int) tokenType {
	for operator, cond := range conditions {
		if cond == number {
			return operator
		}
	}
	return ""
}

func decompCondition(condition map[string]interface{}, action *ShortcutAction) {
	var conditionInt = condition["WFCondition"].(uint64)
	var conditionalOperator = matchConditionOperator(int(conditionInt))
	if conditionalOperator == "" {
		decompError(fmt.Sprintf("Invalid conditional %v", conditionInt), action)
	}
	if conditionalOperator == Empty {
		code.WriteRune('!')
	}

	code.WriteString(decompValue(condition["WFInput"]))

	if conditionalOperator == Any || conditionalOperator == Empty {
		return
	}

	code.WriteString(fmt.Sprintf(" %s ", conditionalOperator))

	if condition["WFNumberValue"] != nil && condition["WFNumberValue"] != "" {
		if reflect.TypeOf(condition["WFNumberValue"]).Kind() == reflect.String {
			var numberValue, convErr = strconv.Atoi(condition["WFNumberValue"].(string))
			handle(convErr)
			code.WriteString(decompValue(numberValue))
		} else {
			code.WriteString(decompValue(condition["WFNumberValue"]))
		}
	} else if _, foundStr := condition["WFConditionalActionString"]; foundStr {
		code.WriteString(decompValue(condition["WFConditionalActionString"]))
	}
}

type DictionaryActionParameters struct {
	WFItems WFItems
}

type WFItems struct {
	Value Value
}

var decompilingDictionary = false

func decompDictionary(action *ShortcutAction) {
	var params DictionaryActionParameters
	mapToStruct(action.WFWorkflowActionParameters, &params)

	currentVariableValue = decompDictionaryValue(params.WFItems.Value.WFDictionaryFieldValueItems)

	checkConstantLiteral(action)
}

func decompDictionaryValue(items []WFDictionaryFieldValueItem) string {
	decompilingDictionary = true

	var dictionary = decompDictionaryItems(items)
	var jsonBytes, jsonErr = json.MarshalIndent(dictionary, strings.Repeat("\t", tabLevel), "\t")
	handle(jsonErr)

	decompilingDictionary = false

	return string(jsonBytes)
}

func isReferenceValue(value any) bool {
	if value == nil {
		return false
	}
	if reflect.TypeOf(value).Kind() == reflect.Map {
		var mapValue = value.(map[string]interface{})
		if mapValue["Value"] != nil && mapValue["WFSerializationType"] == "WFTextTokenAttachment" {
			return true
		}
	}

	return false
}

func decompDictionaryItems(items []WFDictionaryFieldValueItem) (dictionary map[string]interface{}) {
	dictionary = make(map[string]interface{})
	for _, item := range items {
		var itemKey = strings.Trim(decompValue(item.WFKey), "\"")
		dictionary[itemKey] = decompDictionaryItem(item)
	}
	return
}

var numericRegex = regexp.MustCompile("^[0-9]+$")

func decompDictionaryItem(item WFDictionaryFieldValueItem) any {
	var itemStringValue = decompValue(item.WFValue)
	var itemValueType = item.WFItemType
	var itemValueMap = item.WFValue.(map[string]interface{})
	var itemValue any

	if isReferenceValue(itemValueMap["Value"]) {
		return fmt.Sprintf("{%s}", itemStringValue)
	}

	switch dictDataType(itemValueType) {
	case itemTypeNumber:
		itemValue = decompValue(item.WFValue.(map[string]interface{}))
		if !numericRegex.MatchString(itemValue.(string)) {
			break
		}

		var convErr error
		itemValue, convErr = strconv.ParseInt(itemValue.(string), 10, 64)
		handle(convErr)
	case itemTypeBool:
		var wfValue = item.WFValue.(map[string]interface{})
		itemValue = wfValue["Value"].(bool)
	case itemTypeText:
		itemValue = strings.Trim(itemStringValue, "\"")
	case itemTypeArray:
		itemValue = decompArray(itemValueMap["Value"].([]interface{}))
	case itemTypeDict:
		var dictionaryItems []WFDictionaryFieldValueItem
		var value = itemValueMap["Value"].(map[string]interface{})
		mapToStruct(value["Value"].(map[string]interface{})["WFDictionaryFieldValueItems"], &dictionaryItems)
		itemValue = decompDictionaryItems(dictionaryItems)
	default:
		itemValue = itemStringValue
	}

	return itemValue
}

func decompArray(items []interface{}) (array []interface{}) {
	for _, item := range items {
		if reflect.TypeOf(item).Kind() == reflect.Map {
			var fieldValueItem WFDictionaryFieldValueItem
			mapToStruct(item, &fieldValueItem)
			array = append(array, decompDictionaryItem(fieldValueItem))
			continue
		}

		array = append(array, strings.Trim(decompValue(item), "\""))
	}
	return
}

func decompValue(value any) string {
	if value == nil {
		return "nil"
	}

	var valueType = reflect.TypeOf(value).Kind()
	switch valueType {
	case reflect.Map:
		return decompValueObject(value.(map[string]interface{}))
	case reflect.String:
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
		"\"": `\"`,
	}
	for chr, e := range escapes {
		value = strings.ReplaceAll(value, chr, e)
	}

	return value
}

func decompValueObject(value map[string]interface{}) string {
	if v, found := value["Value"]; found {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			value = v.(map[string]interface{})
		}
	}

	if value["WFDictionaryFieldValueItems"] != nil {
		var items []WFDictionaryFieldValueItem
		mapToStruct(value["WFDictionaryFieldValueItems"], &items)

		return decompDictionaryValue(items)
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

			if value["Aggrandizements"] == nil {
				return outputName
			}
		}
		break
	case "Ask":
		if value["Prompt"] == nil {
			return Ask
		}

		return fmt.Sprintf("Ask: \"%s\"", value["Prompt"])
	case globals[ShortcutInput].variableType:
		return ShortcutInput
	}

	return decompObjectValue(value)
}

func decompObjectValue(valueObj any) string {
	var valueType = reflect.TypeOf(valueObj).Kind()
	if valueType != reflect.Map {
		return fmt.Sprintf("%v", valueObj)
	}

	var value = valueObj.(map[string]interface{})

	var attachmentString string
	if value["Value"] != nil {
		if reflect.TypeOf(value["Value"]).Kind() != reflect.Map {
			return fmt.Sprintf("%v", valueObj)
		}
		value = value["Value"].(map[string]interface{})
	}

	if _, found := value["string"]; found {
		attachmentString = value["string"].(string)
	}

	if value["Aggrandizements"] != nil {
		attachmentString = ObjectReplaceCharStr
		decompAttachmentString(&attachmentString, map[string]any{
			"{0, 1}": value,
		})
		attachmentString = strings.Trim(attachmentString, "\"{}")
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
		if attachment.Type != "Variable" {
			for name, global := range globals {
				if global.variableType == attachment.Type {
					variableName = name
				}
			}
		}

		if len(attachment.Aggrandizements) != 0 {
			decompAggrandizements(&variableName, attachment.Aggrandizements)
		}

		attachmentChars[position] = fmt.Sprintf("{%s}", variableName)
	}

	*attachmentString = escapeString(strings.Join(attachmentChars, ""))

	if !decompilingDictionary && !decompilingText {
		if originalString == ObjectReplaceCharStr {
			*attachmentString = strings.Trim(*attachmentString, "{}")
		} else if !strings.Contains("\"", *attachmentString) {
			*attachmentString = fmt.Sprintf("\"%s\"", *attachmentString)
		}
	}
}

func decompAggrandizements(reference *string, aggrs []Aggrandizement) {
	var index string
	var coerce string
	var revContentItems = reversedContentItems()
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

// popLine adds line to the top of the generated Cherri code.
func popLine(line string) {
	var saveCode = code.String()
	code.Reset()
	code.WriteString(fmt.Sprintf("%s\n%s", line, saveCode))
}

var macDefinition bool
var setMacDefinition bool

func decompAction(action *ShortcutAction) {
	if skipDecompAction(action) {
		return
	}

	currentVariableValue = makeActionCallCode(action)

	var isConstant, isVariableValue = checkOutputType(action)
	if isConstant && isVariableValue {
		makeConstantLiteral(action)
	} else if isConstant {
		checkConstantLiteral(action)
	} else if !isVariableValue {
		code.WriteString(tabbedLine(currentVariableValue))
		code.WriteRune('\n')
		currentVariableValue = ""
	}
}

func makeActionCallCode(action *ShortcutAction) string {
	var actionCallCode strings.Builder
	var matchedIdentifier, matchedAction = matchAction(action)
	if matchedIdentifier == "" {
		checkMissingStandardInclude(&action.WFWorkflowActionIdentifier, false)

		matchedIdentifier, matchedAction = matchAction(action)
		if matchedIdentifier == "" {
			actionCallCode.WriteString(makeRawAction(action))
			return actionCallCode.String()
		}
	}

	if (matchedAction.macOnly || matchedAction.nonMacOnly) && !setMacDefinition {
		macDefinition = matchedAction.macOnly && !matchedAction.nonMacOnly
		popLine(fmt.Sprintf("#define mac %v", macDefinition))
		setMacDefinition = true
	}

	actionCallCode.WriteString(fmt.Sprintf("%s(", matchedIdentifier))

	if matchedAction.make != nil || matchedAction.decomp != nil {
		decompActionCustom(&actionCallCode, &matchedAction, action)
	} else {
		var matchedParamsSize = len(matchedAction.parameters)
		if matchedParamsSize > 0 {
			decompActionArguments(&actionCallCode, &matchedAction, action)
		}
	}

	actionCallCode.WriteString(")")

	return actionCallCode.String()
}

// checkOutputType determines if action output is a constant or a variable.
// If it is a constant we will write a constant statement on a new line to prepend the action.
func checkOutputType(action *ShortcutAction) (isConstant bool, isVariableValue bool) {
	if action.WFWorkflowActionParameters[UUID] == nil {
		return
	}
	var uuid = action.WFWorkflowActionParameters[UUID].(string)
	if _, found := uuids[uuid]; !found {
		return
	}

	isConstant = slices.Contains(constUUIDs, uuid) || !slices.Contains(varUUIDs, uuid)
	isVariableValue = slices.Contains(varUUIDs, uuid)

	return
}

// skipDecompAction skips actions we don't support or when necessary.
func skipDecompAction(action *ShortcutAction) bool {
	var identifier = actionIdentifierEnd(action.WFWorkflowActionIdentifier)
	if identifier == "getvariable" {
		var varName = decompValue(action.WFWorkflowActionParameters["WFVariable"])

		insertCodeComment(fmt.Sprintf("TODO: Get Variable not supported: Assign variable here to '%s'.", varName))
		decompWarning(fmt.Sprintf("Get variable '%s' is not supported. Set a variable to that value instead if something was depending on it's output.", varName))

		return true
	}

	if identifier == "nothing" {
		var nextAction = peekActions(1)
		var controlflowActionIdentifiers = []string{"conditional", "repeat.each", "repeat.count", "choosefrommenu"}
		var nextActionIdentifier = actionIdentifierEnd(nextAction.WFWorkflowActionIdentifier)
		if slices.Contains(controlflowActionIdentifiers, nextActionIdentifier) {
			var controlFlowMode = nextAction.WFWorkflowActionParameters["WFControlFlowMode"]
			if controlFlowMode == endStatement || controlFlowMode == statementPart {
				return true
			}
		}
	}

	return false
}

func actionIdentifierEnd(identifier string) string {
	return strings.Replace(identifier, "is.workflow.actions.", "", 1)
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
			if startsWith(Ask, argValue) {
				break
			}
			argValue = strings.Trim(argValue, "\"")
		}

		if argValue != "" {
			if i == 0 {
				actionCallCode.WriteString(argValue)
			} else {
				actionCallCode.WriteString(fmt.Sprintf(", %s", argValue))
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
	if paramsSize == 0 || (onlyUUIDParam && onlyOutputNameParam) {
		return fmt.Sprintf("rawAction(\"%s\")", action.WFWorkflowActionIdentifier)
	}

	var rawParams = processRawParameters(action.WFWorkflowActionParameters)
	var jb, jsonErr = json.MarshalIndent(rawParams, strings.Repeat("\t", tabLevel), "\t")
	handle(jsonErr)

	var arguments = strings.Join([]string{fmt.Sprintf("\"%s\"", action.WFWorkflowActionIdentifier), string(jb)}, ", ")

	return fmt.Sprintf("rawAction(%s)", arguments)
}

func processRawParameters(params map[string]any) map[string]any {
	for key, value := range params {
		if key == UUID || key == "CustomOutputName" {
			delete(params, key)
		}

		if reflect.TypeOf(value).Kind() == reflect.Map {
			decompilingDictionary = true
			params[key] = decompValueObject(value.(map[string]interface{}))
			decompilingDictionary = false
		}
	}

	return params
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
		if actionIdentifier != action.WFWorkflowActionIdentifier && def.overrideIdentifier != action.WFWorkflowActionIdentifier {
			continue
		}

		definition = *def
		name = call

		if splitActions, found := identifierMap[identifier]; found {
			matchSplitAction(&splitActions, action.WFWorkflowActionParameters, &name, &definition)
			if name != "run" && name != "runSelf" {
				return
			}
			var workflow = action.WFWorkflowActionParameters["WFWorkflow"].(map[string]interface{})
			if _, isSelf := workflow["isSelf"]; !isSelf {
				return
			}
			if workflow["isSelf"].(bool) || action.WFWorkflowActionParameters["WFWorkflowName"] == basename {
				name = "runSelf"
				definition = *actions["runSelf"]
			} else {
				name = "run"
			}
		}

		return
	}
	return
}

type actionMatch struct {
	params uint
	values uint
	action actionValue
}

func matchSplitAction(splitActions *[]actionValue, parameters map[string]any, identifier *string, definition *actionDefinition) {
	if args.Using("debug") {
		fmt.Println("## MATCHING SPLIT ACTIONS ##")
		fmt.Println("parameters", parameters)
	}

	var defaultAction, hasDefaultAction = getDefaultAction(splitActions)
	if hasDefaultAction {
		*identifier = defaultAction.identifier
		*definition = *defaultAction.definition
		if args.Using("debug") {
			fmt.Println("has default action", defaultAction.identifier)
		}
	}

	var matches = getSplitActionMatches(splitActions, parameters)

	if len(matches) == 0 {
		return
	}

	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].params > matches[j].params || matches[i].values > matches[j].values
	})

	if args.Using("debug") {
		for _, match := range matches[0:] {
			fmt.Printf("%s()\n", match.action.identifier)
			fmt.Println("params:", match.params, " values:", match.values)
			fmt.Println("---")
		}
		fmt.Print("\n")
	}

	if !competingMatches(matches) {
		if args.Using("debug") {
			fmt.Println("matches are not competing! using default")
		}
		return
	}

	var matchedAction = matches[0]
	*identifier = matchedAction.action.identifier
	*definition = *matchedAction.action.definition
}

func getSplitActionMatches(splitActions *[]actionValue, parameters map[string]any) (matches []actionMatch) {
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
	return
}

// competingMatches determines if the matches for this identifier have more values than 1 matching this action.
func competingMatches(matches []actionMatch) bool {
	var hasCompetingValues bool
	for _, match := range matches {
		if match.values > 1 {
			hasCompetingValues = true
			break
		}
	}
	if hasCompetingValues {
		return true
	}
	for _, match := range matches {
		if match.params > 1 {
			return true
		}
	}
	return false
}

func scoreActionMatch(splitAction actionValue, splitActionParams []parameterDefinition, parameters map[string]any) (matchedParams uint, matchedValues uint) {
	var paramMatches, valueMatches = scoreActionParams(&splitActionParams, parameters)
	matchedParams += paramMatches
	matchedValues += valueMatches

	var splitActionAddParams []parameterDefinition
	if splitAction.definition.addParams != nil {
		for key, value := range splitAction.definition.addParams([]actionArgument{}) {
			splitActionAddParams = append(splitActionAddParams, parameterDefinition{
				key:          key,
				defaultValue: value,
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
			if len(param.enum) > 0 && slices.Contains(getEnum(param.enum), fmt.Sprintf("%s", value)) {
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

			var defaultValueType = reflect.TypeOf(param.defaultValue).Kind()
			var valueType = reflect.TypeOf(value).Kind()
			if defaultValueType == reflect.Map && valueType == reflect.Map {
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
	actionCode.WriteString(strings.Join(arguments, ", "))
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

	fmt.Println("### CONSTANT UUIDs ###")
	fmt.Println(constUUIDs)
	fmt.Print("\n")

	fmt.Println("### UUIDS ###")
	for uuid, name := range uuids {
		fmt.Printf("%s | %s\n", uuid, name)
	}
	fmt.Print("\n")
}

func decompWarning(message string) {
	var linesLen = strings.Count(code.String(), "\n")
	fmt.Println(ansi("Warning:", orange, bold), fmt.Sprintf("%s (%s:%d:0)\n", message, filePath, linesLen+1))
}

func decompError(message string, action *ShortcutAction) {
	fmt.Println(ansi(fmt.Sprintf("Error: %s\n\n", message), red, bold))

	fmt.Println("Action identifier:", action.WFWorkflowActionIdentifier)
	lines = strings.Split(code.String(), "\n")
	var linesLen = len(lines)
	var lastWrittenLine = lines[linesLen-1]
	var prevWrittenLine = lines[linesLen-2]
	fmt.Printf("\nStopped while writing line %d:\n", linesLen)
	fmt.Printf("%s %s\n", ansi(fmt.Sprintf("%d |", linesLen-1), dim), ansi(prevWrittenLine, dim))
	fmt.Printf("%s %s\n", ansi(fmt.Sprintf("%d |", linesLen), dim), ansi(lastWrittenLine, red))

	os.Exit(1)
}
