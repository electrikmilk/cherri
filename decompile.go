/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	plists "howett.net/plist"
	"reflect"
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
				var text = action.WFWorkflowActionParameters["WFTextActionText"]
				if reflect.TypeOf(text).String() == "string" {
					variableValue = fmt.Sprintf("\"%v\"", text)
				}
			case "is.workflow.actions.number":
				variableValue = fmt.Sprintf("%v", action.WFWorkflowActionParameters["WFNumberActionNumber"])
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
					for _, param := range matchedAction.parameters {
						if param.key == "" {
							// TODO: Run make functions
							return
						}
						if value, found := action.WFWorkflowActionParameters[param.key]; found {
							code.WriteString(fmt.Sprintf("%v", value))
						}
						// action.WFWorkflowActionParameters
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

func matchAction(identifier string) string {
	for action, data := range actions {
		if data.identifier == identifier {
			return action
		}
	}

	return ""
}
