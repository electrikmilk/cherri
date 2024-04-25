/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/electrikmilk/args-parser"
	plists "howett.net/plist"
	"os"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Shortcut struct {
	WFWorkflowIcon                      ShortcutIcon
	WFWorkflowActions                   []ShortcutAction
	WFQuickActionSurfaces               []string
	WFWorkflowInputContentItemClasses   []string
	WFWorkflowClientVersion             string
	WFWorkflowMinimumClientVersion      int
	WFWorkflowImportQuestions           interface{}
	WFWorkflowTypes                     []string
	WFWorkflowOutputContentItemClasses  []string
	WFWorkflowHasShortcutInputVariables bool
	WFWorkflowHasOutputFallback         bool
}

type ShortcutIcon struct {
	WFWorkflowIconGlyphNumber int64
	WFWorkflowIconStartColor  int
}

type ShortcutAction struct {
	WFWorkflowActionIdentifier string
	WFWorkflowActionParameters map[string]any
}

var data Shortcut
var code strings.Builder

func decompile(b []byte) {
	var _, marshalErr = plists.Unmarshal(b, &data)
	handle(marshalErr)

	mapVariables()
	mapSplitActions()
	decompileIcon()
	decompileActions()

	if args.Using("debug") {
		printDecompDebug()
	}

	var writeErr = os.WriteFile(basename+"_decompiled.cherri", []byte(code.String()), 0600)
	handle(writeErr)
}

// mapVariables creates a map of variables that are assigned throughout the Shortcut, so we know if an identifier is an assigned variable.
func mapVariables() {
	variables = make(map[string]variableValue)
	uuids = make(map[string]string)
	for _, action := range data.WFWorkflowActions {
		if action.WFWorkflowActionIdentifier == "is.workflow.actions.setvariable" || action.WFWorkflowActionIdentifier == "is.workflow.actions.appendvariable" {
			var varName = strings.ReplaceAll(action.WFWorkflowActionParameters["WFVariableName"].(string), " ", "")
			if _, found := variables[varName]; !found {
				variables[varName] = variableValue{}
			}
			continue
		}

		if action.WFWorkflowActionParameters["WFInput"] != nil {
			var wfInput = action.WFWorkflowActionParameters["WFInput"].(map[string]interface{})
			if wfInput["Value"] != nil {
				var Value = wfInput["Value"].(map[string]interface{})
				if _, found := Value["OutputName"]; !found {
					continue
				}
				if _, found := Value["OutputUUID"]; found {
					if Value["OutputUUID"] == nil {
						continue
					}
					var outputUUID = Value["OutputUUID"].(string)
					if _, found := uuids[outputUUID]; !found {
						var outputName = strings.ReplaceAll(Value["OutputName"].(string), " ", "")
						uuids[outputUUID] = checkDuplicateOutputName(outputName)
						variables[outputName] = variableValue{}
					}
				}
			}
		}
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
		if action.identifier == "" {
			continue
		}

		identifierMap[action.identifier] = append(identifierMap[action.identifier], actionValue{
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
	var icon = data.WFWorkflowIcon
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
			if icon.WFWorkflowIconGlyphNumber != i {
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
	for _, action := range data.WFWorkflowActions {
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
		case "is.workflow.actions.setvariable", "is.workflow.actions.appendvariable":
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
		var customOutputName = action.WFWorkflowActionParameters["CustomOutputName"].(string)
		if _, found := variables[customOutputName]; !found {
			newCodeLine(fmt.Sprintf("const %s = ", customOutputName))
			code.WriteString(currentVariableValue)
			code.WriteRune('\n')
			currentVariableValue = ""
			return
		}
	}
	if _, found := action.WFWorkflowActionParameters["UUID"]; found {
		var uuid = action.WFWorkflowActionParameters["UUID"].(string)
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
	var customOutputName = action.WFWorkflowActionParameters["CustomOutputName"].(string)
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
		if action.WFWorkflowActionIdentifier == "is.workflow.actions.appendvariable" {
			code.WriteString("+= ")
		} else {
			code.WriteString("= ")
		}

		code.WriteString(currentVariableValue)
	} else {
		var decompInput = decompValue(action.WFWorkflowActionParameters["WFInput"])
		var wfInput = action.WFWorkflowActionParameters["WFInput"]
		if reflect.TypeOf(wfInput).String() == dictType {
			var Value = wfInput.(map[string]interface{})["Value"].(map[string]interface{})
			if _, found := Value["OutputName"]; found {
				decompInput = ""
			}
		}
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

func decompDictionary(action *ShortcutAction) {
	var value = action.WFWorkflowActionParameters["WFItems"].(map[string]interface{})
	var Value = value["Value"].(map[string]interface{})
	var dictionary = decompDictionaryItems(Value["WFDictionaryFieldValueItems"].([]interface{}))
	var jsonBytes, jsonErr = json.Marshal(dictionary)
	handle(jsonErr)

	currentVariableValue = string(jsonBytes)

	checkConstantLiteral(action)
}

func decompDictionaryItems(items []interface{}) (dictionary map[string]interface{}) {
	dictionary = make(map[string]interface{})
	for _, item := range items {
		var dictionaryItem = item.(map[string]interface{})
		var itemKey = decompValue(dictionaryItem["WFKey"])
		var itemStringValue = decompValue(dictionaryItem["WFValue"])
		var itemValueType = fmt.Sprintf("%d", dictionaryItem["WFItemType"])
		var itemValue any
		switch dictDataType(itemValueType) {
		case itemTypeNumber:
			if itemStringValue != "" {
				var convErr error
				itemValue, convErr = strconv.Atoi(itemStringValue)
				handle(convErr)
			}
		case itemTypeBool:
			var wfValue = dictionaryItem["WFValue"].(map[string]interface{})
			itemValue = wfValue["Value"]
		case itemTypeText:
			itemValue = strings.Trim(itemStringValue, "\"")
		case itemTypeArray:
			var wfValue = dictionaryItem["WFValue"].(map[string]interface{})
			var Value = wfValue["Value"].([]interface{})
			itemValue = decompArray(Value)
		case itemTypeDict:
			var wfValue = dictionaryItem["WFValue"].(map[string]interface{})
			var Value = wfValue["Value"].(map[string]interface{})
			var dictionaryValue = Value["Value"].(map[string]interface{})
			itemValue = decompDictionaryItems(dictionaryValue["WFDictionaryFieldValueItems"].([]interface{}))
		default:
			itemValue = itemStringValue
		}
		dictionary[itemKey] = itemValue
	}
	return
}

func decompArray(items []interface{}) (array []interface{}) {
	for _, item := range items {
		var itemInterface = item.(map[string]interface{})
		var itemStringValue = decompValue(itemInterface["WFValue"])
		var itemValue any
		var itemValueType = fmt.Sprintf("%d", itemInterface["WFItemType"])
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
		return strings.ReplaceAll(value["OutputName"].(string), " ", "")
	case "ExtensionInput":
		return "ShortcutInput"
	default:
		return decompObjectValue(value)
	}
}

func decompObjectValue(value any) string {
	var valueType = reflect.TypeOf(value).String()
	switch valueType {
	case "map[string]interface {}":
		var Value = value.(map[string]interface{})

		var attachmentString string
		if Value["Value"] != nil {
			if reflect.TypeOf(Value["Value"]).String() != "map[string]interface {}" {
				return fmt.Sprintf("%v", value)
			}
			Value = Value["Value"].(map[string]interface{})
		}

		if _, found := Value["string"]; found {
			attachmentString = Value["string"].(string)
		}

		var attachmentChars = strings.Split(attachmentString, "")
		if attachments, found := Value["attachmentsByRange"]; found {
			for attachmentRange, a := range attachments.(map[string]interface{}) {
				var attachmentRanges = strings.Split(attachmentRange, ",")
				var attachmentPosition = strings.TrimPrefix(attachmentRanges[0], "{")
				var position, convErr = strconv.Atoi(attachmentPosition)
				handle(convErr)

				var attachment = a.(map[string]interface{})
				var variableName string
				if _, found := attachment["OutputName"]; found {
					variableName = attachment["OutputName"].(string)
				} else {
					variableName = attachment["VariableName"].(string)
				}

				attachmentChars[position] = fmt.Sprintf("{%s}", strings.ReplaceAll(variableName, " ", ""))
			}

			attachmentString = fmt.Sprintf("\"%s\"", strings.Join(attachmentChars, ""))
		}

		return attachmentString
	default:
		return fmt.Sprintf("%v", value)
	}
}

var macDefinition bool

func decompAction(action *ShortcutAction) {
	var matchedIdentifier, matchedAction = matchAction(action)
	if matchedIdentifier == "" {
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
	if customOutputName, found := action.WFWorkflowActionParameters["CustomOutputName"]; found {
		if _, foundVar := variables[customOutputName.(string)]; !foundVar {
			newCodeLine(fmt.Sprintf("const %s = ", strings.ReplaceAll(customOutputName.(string), " ", "")))
			isConstant = true
		} else {
			isVariableValue = true
		}
	}
	if _, found := action.WFWorkflowActionParameters["UUID"]; found {
		var uuid = action.WFWorkflowActionParameters["UUID"].(string)
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
				actionCallCode.WriteRune(',')
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

func matchAction(action *ShortcutAction) (identifier string, definition actionDefinition) {
	for call, def := range actions {
		var ident string
		if def.identifier != "" {
			ident = def.identifier
		} else {
			ident = call
		}
		var shortcutsIdentifier = fmt.Sprintf("is.workflow.actions.%s", ident)
		if shortcutsIdentifier == action.WFWorkflowActionIdentifier || definition.appIdentifier == action.WFWorkflowActionIdentifier {
			identifier = call
			definition = *def

			switch identifier {
			case "confirm":
				if value, found := action.WFWorkflowActionParameters["WFAlertActionCancelButtonShown"]; found {
					if value == false {
						identifier = "alert"
					}
				}
			case "run":
				var runSelfIdentifier = "runSelf"
				if _, isSelf := action.WFWorkflowActionParameters["isSelf"]; isSelf {
					identifier = runSelfIdentifier
					break
				}
				if name, foundName := action.WFWorkflowActionParameters["workflowName"]; foundName {
					if name == basename {
						identifier = runSelfIdentifier
						break
					}
				} else {
					identifier = runSelfIdentifier
				}
			case "outputOrClipboard", "mustOutput":
				identifier = "output"
			}

			if splitActions, found := identifierMap[ident]; found {
				matchSplitAction(&splitActions, action.WFWorkflowActionParameters, &identifier, &definition)
			}

			break
		}
	}
	return
}

type actionMatch struct {
	matchedParams int
	action        actionValue
}

func matchSplitAction(splitAction *[]actionValue, parameters map[string]any, identifier *string, definition *actionDefinition) {
	var matches []actionMatch
	for _, splitAction := range *splitAction {
		if splitAction.definition.identifier == "getitemfromlist" {
			matchListAction(parameters, identifier, definition)
			return
		}

		var matchedParams int
		var paramsSize = len(splitAction.definition.parameters)
		for _, param := range splitAction.definition.parameters {
			if param.key == "" {
				continue
			}
			if _, found := parameters[param.key]; found {
				matchedParams++
			}
		}

		if paramsSize == matchedParams {
			matches = append(matches, actionMatch{
				matchedParams: matchedParams,
				action:        splitAction,
			})
		}
	}
	if len(matches) < 2 {
		return
	}
	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].matchedParams > matches[j].matchedParams
	})
	var matchedAction = matches[0]
	*identifier = matchedAction.action.identifier
	*definition = *matchedAction.action.definition
}

func matchListAction(parameters map[string]any, identifier *string, definition *actionDefinition) {
	switch parameters["WFItemSpecifier"] {
	case "First Item":
		*identifier = "getFirstItem"
		definition = actions["getFirstItem"]
	case "Last Item":
		*identifier = "getLastItem"
		definition = actions["getLastItem"]
	case "Random Item":
		*identifier = "getRandomItem"
		definition = actions["getRandomItem"]
	case "Item At Index":
		*identifier = "getListItem"
		definition = actions["getListItem"]
	case "Items in Range":
		*identifier = "getListItems"
		definition = actions["getListItems"]
	}
}

func printDecompDebug() {
	fmt.Println(ansi("##### DEBUG #####\n", red))

	fmt.Println("### ACTIONS ###")
	for _, action := range data.WFWorkflowActions {
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
