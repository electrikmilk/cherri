/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"encoding/json"
	"fmt"
	plists "howett.net/plist"
	"os"
	"reflect"
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

	var writeErr = os.WriteFile(basename+"_decompiled.cherri", []byte(code.String()), 0600)
	handle(writeErr)
}

// mapVariables creates a map of variables that are assigned throughout the Shortcut, so we know if an identifier is an assigned variable.
func mapVariables() {
	variables = make(map[string]variableValue)
	for _, action := range data.WFWorkflowActions {
		if action.WFWorkflowActionIdentifier != "is.workflow.actions.setvariable" && action.WFWorkflowActionIdentifier != "is.workflow.actions.appendvariable" {
			continue
		}
		var varName = action.WFWorkflowActionParameters["WFVariableName"].(string)
		if _, found := variables[varName]; !found {
			variables[varName] = variableValue{}
		}
	}
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
	printIdentifierMap()
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
		case "is.workflow.actions.gettext":
			decompTextValue(&action)
		case "is.workflow.actions.number":
			decompNumberValue(&action)
		case "is.workflow.actions.dictionary":
			currentVariableValue = decompDictionary(action.WFWorkflowActionParameters["WFItems"].(map[string]interface{}))

			var customOutputName = action.WFWorkflowActionParameters["CustomOutputName"].(string)
			if _, found := variables[customOutputName]; !found {
				newCodeLine(fmt.Sprintf("const %s = ", customOutputName))
				code.WriteString(currentVariableValue)
				code.WriteRune('\n')
				currentVariableValue = ""
			}
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

func decompTextValue(action *ShortcutAction) {
	currentVariableValue = decompValue(action.WFWorkflowActionParameters["WFTextActionText"])
	if reflect.TypeOf(action.WFWorkflowActionParameters["WFTextActionText"]).String() == "string" {
		return
	}

	var customOutputName = action.WFWorkflowActionParameters["CustomOutputName"].(string)
	if _, found := variables[customOutputName]; !found {
		newCodeLine(fmt.Sprintf("const %s = ", customOutputName))
		code.WriteString(decompValue(action.WFWorkflowActionParameters["WFTextActionText"]))
		code.WriteRune('\n')
	} else {
		currentVariableValue = decompValue(action.WFWorkflowActionParameters["WFTextActionText"])
	}
}

var macDefinition bool

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

func decompVariable(action *ShortcutAction) {
	newCodeLine("@%s", action.WFWorkflowActionParameters["WFVariableName"].(string))

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
		}

		if _, found := action.WFWorkflowActionParameters["WFNumberValue"]; found {
			var numberValue, convErr = strconv.Atoi(action.WFWorkflowActionParameters["WFNumberValue"].(string))
			handle(convErr)
			code.WriteString(decompValue(numberValue))
		} else if _, foundStr := action.WFWorkflowActionParameters["WFConditionalActionString"]; foundStr {
			code.WriteString(decompValue(action.WFWorkflowActionParameters["WFConditionalActionString"]))
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

func decompDictionary(value map[string]interface{}) string {
	var Value = value["Value"].(map[string]interface{})
	var dictionary = decompDictionaryItems(Value["WFDictionaryFieldValueItems"].([]interface{}))
	var jsonBytes, jsonErr = json.Marshal(dictionary)
	handle(jsonErr)

	return string(jsonBytes)
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
		return value["OutputName"].(string)
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

				var chars = strings.Split(attachmentString, "")
				chars[position] = fmt.Sprintf("{%s}", variableName)
				attachmentString = strings.Join(chars, "")
			}

			attachmentString = fmt.Sprintf("\"%s\"", attachmentString)
		}

		return attachmentString
	default:
		return fmt.Sprintf("%v", value)
	}
}

func matchAction(action *ShortcutAction) (identifier string, definition actionDefinition) {
	for ident, def := range actions {
		var shortcutsIdentifier = "is.workflow.actions."
		if def.identifier != "" {
			shortcutsIdentifier += def.identifier
		} else {
			shortcutsIdentifier += ident
		}
		if shortcutsIdentifier == action.WFWorkflowActionIdentifier || definition.appIdentifier == action.WFWorkflowActionIdentifier {
			identifier = ident
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
			}
			break
		}
	}
	return
}

func printIdentifierMap() {
	for identifier, actions := range identifierMap {
		fmt.Println(identifier)
		for _, action := range actions {
			fmt.Print("\t")
			setCurrentAction(action.identifier, action.definition)
			fmt.Println(generateActionDefinition(parameterDefinition{}, false, false))
		}
		fmt.Print("\n")
	}
}

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
			newCodeLine(fmt.Sprintf("const %s = ", customOutputName))
			isConstant = true
		} else {
			isVariableValue = true
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
