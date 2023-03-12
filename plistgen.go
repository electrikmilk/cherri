/*
 * Copyright (c) 2023 Brandon Jordan
 */

package main

import (
	"reflect"
)

const header = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"https://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n<dict>\n"
const footer = "</dict>\n</plist>"

func makePlist() (plist string) {
	uuids = make(map[string]string)
	plist = header
	plist += plistKeyValue("WFQuickActionSurfaces", Array, []string{})

	if noInput.name != "" {
		plist += plistKeyValue("WFWorkflowNoInputBehavior", Dictionary, []plistData{
			{
				key:      "Name",
				dataType: Text,
				value:    noInput.name,
			},
			{
				key:      "Parameters",
				dataType: Dictionary,
				value:    noInput.params,
			},
		})
	}

	for _, tok := range tokens {
		switch tok.typeof {
		case Var, AddTo:
			var setVariableParams = []plistData{
				{
					key:      "WFVariableName",
					dataType: Text,
					value:    tok.ident,
				},
			}
			if tok.value != nil {
				var varUUID = shortcutsUUID()
				makeVariableValue(&tok, &varUUID)
				uuids[tok.ident] = varUUID
				if tok.valueType != Arr {
					if tok.valueType == Variable {
						setVariableParams = append(setVariableParams, variablePlistValue("WFInput", tok.value.(string), tok.ident))
					} else {
						setVariableParams = append(setVariableParams, inputValue("WFInput", tok.ident, varUUID))
					}
					setVariableParams = append(setVariableParams, plistData{
						key:      "WFSerializationType",
						dataType: Text,
						value:    "WFTextTokenAttachment",
					})
				}
			}
			if tok.typeof == Var {
				shortcutActions = append(shortcutActions, makeStdAction("setvariable", setVariableParams))
			} else if tok.typeof == AddTo && tok.valueType != Arr {
				shortcutActions = append(shortcutActions, makeStdAction("appendvariable", setVariableParams))
			}
			if tok.valueType == Arr {
				for _, value := range tok.value.([]interface{}) {
					var UUID = shortcutsUUID()
					var valueType tokenType
					var addToVariableParams []plistData
					var itemIdent string
					switch reflect.TypeOf(value).String() {
					case stringType:
						valueType = String
						itemIdent = "Text"
					case intType:
						valueType = Integer
						itemIdent = "Number"
					case dictType:
						valueType = Dict
						itemIdent = "Dictionary"
					}
					makeVariableValue(&token{
						typeof:    valueType,
						ident:     itemIdent,
						valueType: valueType,
						value:     value,
					}, &UUID)
					addToVariableParams = append(addToVariableParams, inputValue("WFInput", itemIdent, UUID))
					addToVariableParams = append(addToVariableParams, plistData{
						key:      "WFVariableName",
						dataType: Text,
						value:    tok.ident,
					})
					shortcutActions = append(shortcutActions, makeStdAction("appendvariable", addToVariableParams))
				}
			}
		case Comment:
			shortcutActions = append(shortcutActions, makeStdAction("comment", []plistData{
				{
					key:      "WFCommentActionText",
					dataType: Text,
					value:    tok.value,
				},
			}))
		case Action:
			currentAction = tok.value.(action).ident
			callAction(tok.value.(action).args, plistData{}, plistData{})
		case Repeat:
			if tok.valueType == EndClosure {
				shortcutActions = append(shortcutActions, makeStdAction("repeat.count", []plistData{
					{
						key:      "WFControlFlowMode",
						dataType: Number,
						value:    endStatement,
					},
					{
						key:      "GroupingIdentifier",
						dataType: Text,
						value:    currentGroupingUUID,
					},
					{
						key:      "UUID",
						dataType: Text,
						value:    shortcutsUUID(),
					},
				}))
			} else {
				shortcutActions = append(shortcutActions, makeStdAction("repeat.count", []plistData{
					{
						key:      "WFControlFlowMode",
						dataType: Number,
						value:    startStatement,
					},
					paramValue("WFRepeatCount", actionArgument{
						valueType: tok.valueType,
						value:     tok.value,
					}, Integer, Number),
					{
						key:      "GroupingIdentifier",
						dataType: Text,
						value:    currentGroupingUUID,
					},
				}))
			}
		case RepeatWithEach:
			if tok.valueType == EndClosure {
				shortcutActions = append(shortcutActions, makeStdAction("repeat.each", []plistData{
					{
						key:      "WFControlFlowMode",
						dataType: Number,
						value:    endStatement,
					},
					{
						key:      "GroupingIdentifier",
						dataType: Text,
						value:    currentGroupingUUID,
					},
					{
						key:      "UUID",
						dataType: Text,
						value:    shortcutsUUID(),
					},
				}))
			} else {
				shortcutActions = append(shortcutActions, makeStdAction("repeat.each", []plistData{
					{
						key:      "WFControlFlowMode",
						dataType: Number,
						value:    startStatement,
					},
					paramValue("WFInput", actionArgument{
						valueType: tok.valueType,
						value:     tok.value,
					}, Variable, Text),
					{
						key:      "GroupingIdentifier",
						dataType: Text,
						value:    currentGroupingUUID,
					},
				}))
			}
		case Menu:
			var controlFlow = startStatement
			if tok.valueType == EndClosure {
				controlFlow = endStatement
			}
			var menuParams = []plistData{
				{
					key:      "GroupingIdentifier",
					dataType: Text,
					value:    tok.ident,
				},
				{
					key:      "WFControlFlowMode",
					dataType: Number,
					value:    controlFlow,
				},
			}
			if tok.valueType != EndClosure {
				if tok.valueType == Variable {
					menuParams = append(menuParams, paramValue("WFMenuPrompt", actionArgument{
						valueType: tok.valueType,
						value:     tok.value,
					}, String, Text))
				} else {
					menuParams = append(menuParams, plistData{
						key:      "WFMenuPrompt",
						dataType: Text,
						value:    tok.value,
					})
				}
				var menuItemParams = menus[tok.ident]
				var menuItems []string
				for _, item := range menuItemParams {
					if item.valueType == Variable {
						menuItems = append(menuItems, plistValue(Dictionary, []plistData{
							{
								key:      "WFItemType",
								dataType: Number,
								value:    0,
							},
							paramValue("WFValue", actionArgument{
								valueType: tok.valueType,
								value:     tok.value,
							}, String, Text),
						}))
					} else {
						menuItems = append(menuItems, plistValue(Text, item.value.(string)))
					}
				}
				menuParams = append(menuParams, plistData{
					key:      "WFMenuItems",
					dataType: Array,
					value:    menuItems,
				})
			}
			shortcutActions = append(shortcutActions, makeStdAction("choosefrommenu", menuParams))
		case Case:
			shortcutActions = append(shortcutActions, makeStdAction("choosefrommenu", []plistData{
				{
					key:      "GroupingIdentifier",
					dataType: Text,
					value:    tok.ident,
				},
				{
					key:      "WFControlFlowMode",
					dataType: Number,
					value:    statementPart,
				},
				paramValue("WFMenuItemTitle", actionArgument{
					valueType: tok.valueType,
					value:     tok.value,
				}, String, Text),
			}))
		case Conditional:
			var controlFlowMode int
			var conditionalParams = []plistData{
				{
					key:      "GroupingIdentifier",
					dataType: Text,
					value:    tok.ident,
				},
				{
					key:      "UUID",
					dataType: Text,
					value:    shortcutsUUID(),
				},
			}
			switch tok.valueType {
			case If:
				var cond = tok.value.(condition)
				conditionalParams = append(conditionalParams, plistData{
					key:      "WFInput",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "Type",
							dataType: Text,
							value:    "Variable",
						},
						variablePlistValue("Variable", cond.variableOneValue.(string), tok.ident),
					},
				})
				if cond.variableTwoValue != nil {
					condParam("", &conditionalParams, &cond.variableTwoType, cond.variableTwoValue)
				}
				if cond.variableThreeValue != nil {
					condParam("WFAnotherNumber", &conditionalParams, &cond.variableThreeType, cond.variableThreeValue)
				}
				conditionalParams = append(conditionalParams, plistData{
					key:      "WFCondition",
					dataType: Number,
					value:    cond.condition,
				})
				controlFlowMode = startStatement
			case Else:
				controlFlowMode = statementPart
			case EndClosure:
				controlFlowMode = endStatement
			}
			conditionalParams = append(conditionalParams, plistData{
				key:      "WFControlFlowMode",
				dataType: Number,
				value:    controlFlowMode,
			})
			shortcutActions = append(shortcutActions, makeStdAction("conditional", conditionalParams))
		}
	}
	plist += plistKeyValue("WFWorkflowActions", Array, shortcutActions)
	plist += plistKeyValue("WFWorkflowClientVersion", Text, "1146.14")

	if workflowName != "" {
		plist += plistKeyValue("WFWorkflowName", Text, workflowName)
	}

	plist += plistKeyValue("WFWorkflowHasOutputFallback", Boolean, false)
	plist += plistKeyValue("WFWorkflowHasShortcutInputVariables", Boolean, false)

	plist += plistKeyValue("WFWorkflowIcon", Dictionary, []plistData{
		{
			key:      "WFWorkflowIconStartColor",
			dataType: Number,
			value:    iconColor,
		},
		{
			key:      "WFWorkflowIconGlyphNumber",
			dataType: Number,
			value:    iconGlyph,
		},
	})

	if len(questions) > 0 {
		var importQuestions []string
		for _, q := range questions {
			importQuestions = append(importQuestions, plistValue(Dictionary, []plistData{
				{
					key:      "ParameterKey",
					dataType: Text,
					value:    q.parameter,
				},
				{
					key:      "Category",
					dataType: Text,
					value:    "Parameter",
				},
				{
					key:      "ActionIndex",
					dataType: Number,
					value:    q.actionIndex,
				},
				{
					key:      "Text",
					dataType: Text,
					value:    q.text,
				},
				{
					key:      "DefaultValue",
					dataType: Text,
					value:    q.defaultValue,
				},
			}))
		}
		plist += plistKeyValue("WFWorkflowImportQuestions", Array, importQuestions)
	} else {
		plist += plistKeyValue("WFWorkflowImportQuestions", Array, []string{})
	}

	var inputContentItems []string
	if len(inputs) == 0 {
		makeContentItems()
		for _, input := range contentItems {
			inputContentItems = append(inputContentItems, plistValue(Text, input))
		}
	} else {
		for _, input := range inputs {
			inputContentItems = append(inputContentItems, plistValue(Text, input))
		}
	}

	plist += plistKeyValue("WFWorkflowInputContentItemClasses", Array, inputContentItems)

	var outputContentItems []string
	if len(outputs) != 0 {
		makeContentItems()
		for _, output := range outputs {
			outputContentItems = append(outputContentItems, plistValue(Text, output))
		}
	}
	plist += plistKeyValue("WFWorkflowOutputContentItemClasses", Array, outputContentItems)

	plist += plistKeyValue("WFWorkflowMinimumClientVersion", Number, minVersion)
	plist += plistKeyValue("WFWorkflowMinimumClientVersionString", Text, minVersion)
	plist += plistKeyValue("WFWorkflowHasShortcutInputVariables", Boolean, hasShortcutInputVariables)

	var wfWorkflowTypes []string
	if len(types) != 0 {
		for _, wtype := range types {
			wfWorkflowTypes = append(wfWorkflowTypes, plistValue(Text, wtype))
		}
	}
	plist += plistKeyValue("WFWorkflowTypes", Array, wfWorkflowTypes)

	plist += footer

	return
}
