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

	decompileIcon()
	decompileActions()
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

func decompileActions() {
	var variableValue any
	for _, action := range data.WFWorkflowActions {
		var identifier = matchAction(action.WFWorkflowActionIdentifier)
		if identifier == "" {
			switch action.WFWorkflowActionIdentifier {
			case "is.workflow.actions.gettext":
				variableValue = decompValue(action.WFWorkflowActionParameters["WFTextActionText"])
			case "is.workflow.actions.number":
				var number, convErr = strconv.Atoi(action.WFWorkflowActionParameters["WFNumberActionNumber"].(string))
				handle(convErr)
				variableValue = decompValue(number)
			case "is.workflow.actions.dictionary":
				variableValue = decompDictionary(action.WFWorkflowActionParameters["WFItems"].(map[string]interface{}))
			case "is.workflow.actions.setvariable":
				code.WriteRune('@')
				code.WriteString(action.WFWorkflowActionParameters["WFVariableName"].(string))

				if variableValue != nil {
					code.WriteString(" = ")
					code.WriteString(fmt.Sprintf("%v", variableValue))
				}

				variableValue = nil
				code.WriteRune('\n')
			default:
				var matchedAction actionDefinition
				for identifier, definition := range actions {
					var shortcutsIdentifier = "is.workflow.actions." + definition.identifier
					if shortcutsIdentifier == action.WFWorkflowActionIdentifier || definition.appIdentifier == action.WFWorkflowActionIdentifier {
						matchedAction = *definition
						code.WriteString(fmt.Sprintf("%s(", identifier))
						break
					}
				}
				if matchedAction.identifier != "" {
					var matchedParamsSize = len(matchedAction.parameters)
					for i, param := range matchedAction.parameters {
						if param.key == "" {
							// TODO: Run make functions
							return
						}
						if value, found := action.WFWorkflowActionParameters[param.key]; found {
							code.WriteString(decompValue(value))
						}
						if matchedParamsSize != 1 && matchedParamsSize > i {
							code.WriteRune(',')
						}
					}
					code.WriteString(")\n")
				}
			}
			continue
		}

		code.WriteString(identifier)
		code.WriteRune('(')

		code.WriteString(")\n")
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
			var convErr error
			itemValue, convErr = strconv.Atoi(itemStringValue)
			handle(convErr)
		case itemTypeBool:
			switch itemStringValue {
			case "true":
				itemValue = true
			case "false":
				itemValue = false
			default:
				itemValue = itemStringValue
			}
		case itemTypeText:
			itemValue = strings.Trim(itemStringValue, "\"")
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
	// value["WFSerializationType"].(string)
	var valueType = reflect.TypeOf(value["Value"]).String()
	switch valueType {
	case "map[string]interface {}":
		var Value = value["Value"].(map[string]interface{})
		var attachmentString = Value["string"].(string)
		if _, found := Value["attachmentsByRange"]; found {
			for attachmentRange, a := range Value["attachmentsByRange"].(map[string]interface{}) {
				var position, convErr = strconv.Atoi(strings.TrimPrefix(strings.Split(attachmentRange, ",")[0], "{"))
				handle(convErr)
				var attachment = a.(map[string]interface{})
				// attachment["Type"]
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

func matchAction(identifier string) string {
	for action, data := range actions {
		if data.identifier == identifier {
			return action
		}
	}

	return ""
}
