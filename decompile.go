/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"encoding/json"
	"fmt"
	plists "howett.net/plist"
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
	decompileIcon()
	decompileActions()

	fmt.Println(code.String())
}

func decompileIcon() {
	var icon = data.WFWorkflowIcon
	if icon.WFWorkflowIconStartColor != iconColor {
		makeColors()
		for name, i := range colors {
			if icon.WFWorkflowIconStartColor != i {
				continue
			}

			code.WriteString(fmt.Sprintf("#define color %s\n", name))
		}
	}

	if icon.WFWorkflowIconGlyphNumber != iconGlyph {
		for name, i := range glyphs {
			if icon.WFWorkflowIconGlyphNumber != i {
				continue
			}

			code.WriteString(fmt.Sprintf("#define glyph %s\n", name))
		}
	}

	code.WriteRune('\n')
}

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

func decompileActions() {
	var variableValue string
	for _, action := range data.WFWorkflowActions {
		switch action.WFWorkflowActionIdentifier {
		case "is.workflow.actions.gettext":
			variableValue = decompValue(action.WFWorkflowActionParameters["WFTextActionText"])
		case "is.workflow.actions.number":
			var value = action.WFWorkflowActionParameters["WFNumberActionNumber"]
			if reflect.TypeOf(value).String() == dictType {
				value = decompValue(value)
			}
			var number int
			if value != "" {
				var convErr error
				number, convErr = strconv.Atoi(value.(string))
				handle(convErr)
			}
			variableValue = decompValue(number)
		case "is.workflow.actions.dictionary":
			variableValue = decompDictionary(action.WFWorkflowActionParameters["WFItems"].(map[string]interface{}))
		case "is.workflow.actions.setvariable":
			code.WriteRune('@')
			code.WriteString(action.WFWorkflowActionParameters["WFVariableName"].(string))

			if variableValue != "" {
				code.WriteString(fmt.Sprintf(" = %s", variableValue))
			}

			variableValue = ""
			code.WriteRune('\n')
		case "is.workflow.actions.appendvariable":
			code.WriteRune('@')
			code.WriteString(action.WFWorkflowActionParameters["WFVariableName"].(string))

			if variableValue != "" {
				code.WriteString(fmt.Sprintf(" += %s", variableValue))
			}

			variableValue = ""
			code.WriteRune('\n')
		default:
			var matchedAction actionDefinition
			var matchedIdentifier string
			for identifier, definition := range actions {
				var shortcutsIdentifier = "is.workflow.actions." + definition.identifier
				if shortcutsIdentifier == action.WFWorkflowActionIdentifier || definition.appIdentifier == action.WFWorkflowActionIdentifier {
					matchedIdentifier = identifier
					if matchedIdentifier == "confirm" {
						if value, found := action.WFWorkflowActionParameters["WFAlertActionCancelButtonShown"]; found {
							if value == false {
								matchedIdentifier = "alert"
							}
						}
					}
					matchedAction = *definition
					break
				}
			}
			if matchedAction.identifier == "" {
				continue
			}

			var isVariableValue = false
			var actionCallCode strings.Builder
			if customOutputName, found := action.WFWorkflowActionParameters["CustomOutputName"]; found {
				if _, foundVar := variables[customOutputName.(string)]; !foundVar {
					code.WriteString(fmt.Sprintf("const %s = ", customOutputName))
				} else {
					isVariableValue = true
				}
			}
			actionCallCode.WriteString(fmt.Sprintf("%s(", matchedIdentifier))

			var matchedParamsSize = len(matchedAction.parameters)
			for i, param := range matchedAction.parameters {
				if param.key == "" {
					// TODO: Run make functions
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
				variableValue = actionCallCode.String()
			} else {
				code.WriteString(actionCallCode.String())
				code.WriteRune('\n')
			}
		}
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
			switch itemStringValue {
			case "true":
				itemValue = true
			case "false":
				itemValue = false
			}
		case itemTypeText:
			itemValue = strings.Trim(itemStringValue, "\"")
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

func decompValue(value any) string {
	var valueType = reflect.TypeOf(value).String()
	switch valueType {
	case "map[string]interface {}":
		return decompValueObject(value.(map[string]interface{}))
	case stringType:
		return fmt.Sprintf("\"%s\"", value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func decompValueObject(value map[string]interface{}) string {
	var valueType = reflect.TypeOf(value["Value"]).String()
	switch valueType {
	case "map[string]interface {}":
		fmt.Println("value", value)
		var attachmentString string
		var Value = value["Value"].(map[string]interface{})
		if _, found := Value["string"]; found {
			attachmentString = Value["string"].(string)
		}

		if _, found := Value["attachmentsByRange"]; found {
			for attachmentRange, a := range Value["attachmentsByRange"].(map[string]interface{}) {
				var position, convErr = strconv.Atoi(strings.TrimPrefix(strings.Split(attachmentRange, ",")[0], "{"))
				handle(convErr)

				var attachment = a.(map[string]interface{})
				var variableName = attachment["VariableName"]
				var chars = strings.Split(attachmentString, "")
				chars[position] = fmt.Sprintf("{%s}", variableName)
				attachmentString = strings.Join(chars, "")
			}

			attachmentString = fmt.Sprintf("\"%s\"", attachmentString)
		}

		return attachmentString
	default:
		return fmt.Sprintf("%v", value["Value"])
	}
}
