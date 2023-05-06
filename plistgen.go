/*
 * Copyright (c) Brandon Jordan
 */

package main

import "reflect"

const header = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"https://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n<dict>\n"
const footer = "</dict>\n</plist>"

var plist string

func makePlist() {
	uuids = make(map[string]string)
	plist = header

	plist += plistKeyValue("WFWorkflowHasOutputFallback", Boolean, false)
	plist += plistKeyValue("WFWorkflowMinimumClientVersion", Number, minVersion)
	plist += plistKeyValue("WFWorkflowMinimumClientVersionString", Text, minVersion)
	plist += plistKeyValue("WFWorkflowHasShortcutInputVariables", Boolean, hasShortcutInputVariables)

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

	plistActions()

	plist += plistKeyValue("WFWorkflowClientVersion", Text, "1146.14")

	if workflowName != "" {
		plist += plistKeyValue("WFWorkflowName", Text, workflowName)
	}

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
		plistImportQuestions()
	} else {
		plist += plistKeyValue("WFWorkflowImportQuestions", Array, []string{})
	}

	plistContentItems()
	plistWorkflowTypes()

	plist += footer
}

func plistActions() {
	for _, t := range tokens {
		switch t.typeof {
		case Var, AddTo:
			plistVariable(&t)
		case Comment:
			plistComment(&t)
		case Action:
			currentAction = t.value.(action).ident
			plistAction(t.value.(action).args, plistData{}, plistData{})
		case Repeat:
			plistRepeat(&t)
		case RepeatWithEach:
			plistRepeatEach(&t)
		case Menu:
			plistMenu(&t)
		case Item:
			plistMenuItem(&t)
		case Conditional:
			plistConditional(&t)
		}
	}
	plist += plistKeyValue("WFWorkflowActions", Array, shortcutActions)
}

func plistComment(t *token) {
	shortcutActions = append(shortcutActions, makeStdAction("comment", []plistData{
		{
			key:      "WFCommentActionText",
			dataType: Text,
			value:    t.value,
		},
	}))
}

func plistVariable(t *token) {
	var setVariableParams = []plistData{
		{
			key:      "WFVariableName",
			dataType: Text,
			value:    t.ident,
		},
	}
	if t.value != nil {
		var varUUID = shortcutsUUID()
		makeVariableValue(t, &varUUID)
		uuids[t.ident] = varUUID
		if t.valueType != Arr {
			if t.valueType == Variable {
				setVariableParams = append(setVariableParams, variablePlistValue("WFInput", t.value.(string), t.ident))
			} else {
				setVariableParams = append(setVariableParams, inputValue("WFInput", t.ident, varUUID))
			}
			setVariableParams = append(setVariableParams, plistData{
				key:      "WFSerializationType",
				dataType: Text,
				value:    "WFTextTokenAttachment",
			})
		}
	}
	if t.typeof == Var {
		shortcutActions = append(shortcutActions, makeStdAction("setvariable", setVariableParams))
	} else if t.typeof == AddTo && t.valueType != Arr {
		shortcutActions = append(shortcutActions, makeStdAction("appendvariable", setVariableParams))
	}
	if t.valueType == Arr {
		plistArrayVariable(t)
	}
}

func plistArrayVariable(t *token) {
	for _, value := range t.value.([]interface{}) {
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
			value:    t.ident,
		})
		shortcutActions = append(shortcutActions, makeStdAction("appendvariable", addToVariableParams))
	}
}

func plistConditional(t *token) {
	var controlFlowMode int
	var conditionalParams = []plistData{
		{
			key:      "GroupingIdentifier",
			dataType: Text,
			value:    t.ident,
		},
		{
			key:      "UUID",
			dataType: Text,
			value:    shortcutsUUID(),
		},
	}
	switch t.valueType {
	case If:
		var cond = t.value.(condition)
		conditionalParams = append(conditionalParams, plistData{
			key:      "WFInput",
			dataType: Dictionary,
			value: []plistData{
				{
					key:      "Type",
					dataType: Text,
					value:    "Variable",
				},
				variablePlistValue("Variable", cond.variableOneValue.(string), t.ident),
			},
		})
		if cond.variableTwoValue != nil {
			conditionalParameter("", &conditionalParams, &cond.variableTwoType, cond.variableTwoValue)
		}
		if cond.variableThreeValue != nil {
			conditionalParameter("WFAnotherNumber", &conditionalParams, &cond.variableThreeType, cond.variableThreeValue)
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

func plistMenu(t *token) {
	var controlFlow = startStatement
	if t.valueType == EndClosure {
		controlFlow = endStatement
	}
	var menuParams = []plistData{
		{
			key:      "GroupingIdentifier",
			dataType: Text,
			value:    t.ident,
		},
		{
			key:      "WFControlFlowMode",
			dataType: Number,
			value:    controlFlow,
		},
	}
	if t.valueType != EndClosure {
		if t.valueType == Variable {
			menuParams = append(menuParams, paramValue("WFMenuPrompt", actionArgument{
				valueType: t.valueType,
				value:     t.value,
			}, String, Text))
		} else {
			menuParams = append(menuParams, plistData{
				key:      "WFMenuPrompt",
				dataType: Text,
				value:    t.value,
			})
		}
		var menuItemParams = menus[t.ident]
		var menuItems []string
		for _, item := range menuItemParams {
			menuItems = append(menuItems, plistValue(Dictionary, []plistData{
				{
					key:      "WFItemType",
					dataType: Number,
					value:    0,
				},
				paramValue("WFValue", actionArgument{
					valueType: item.valueType,
					value:     item.value,
				}, String, Text),
			}))
		}
		menuParams = append(menuParams, plistData{
			key:      "WFMenuItems",
			dataType: Array,
			value:    menuItems,
		})
	}
	shortcutActions = append(shortcutActions, makeStdAction("choosefrommenu", menuParams))
}

func plistMenuItem(t *token) {
	shortcutActions = append(shortcutActions, makeStdAction("choosefrommenu", []plistData{
		{
			key:      "GroupingIdentifier",
			dataType: Text,
			value:    t.ident,
		},
		{
			key:      "WFControlFlowMode",
			dataType: Number,
			value:    statementPart,
		},
		paramValue("WFMenuItemAttributedTitle", actionArgument{
			valueType: t.valueType,
			value:     t.value,
		}, String, Text),
		paramValue("WFMenuItemTitle", actionArgument{
			valueType: t.valueType,
			value:     t.value,
		}, String, Text),
	}))
}

func plistRepeat(t *token) {
	if t.valueType == EndClosure {
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
				valueType: t.valueType,
				value:     t.value,
			}, Integer, Number),
			{
				key:      "GroupingIdentifier",
				dataType: Text,
				value:    currentGroupingUUID,
			},
		}))
	}
}

func plistRepeatEach(t *token) {
	if t.valueType == EndClosure {
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
				valueType: t.valueType,
				value:     t.value,
			}, Variable, Text),
			{
				key:      "GroupingIdentifier",
				dataType: Text,
				value:    currentGroupingUUID,
			},
		}))
	}
}

var importQuestions []string

func plistImportQuestions() {
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
}

func plistWorkflowTypes() {
	var wfWorkflowTypes []string
	if len(types) != 0 {
		for _, wtype := range types {
			wfWorkflowTypes = append(wfWorkflowTypes, plistValue(Text, wtype))
		}
	}
	plist += plistKeyValue("WFWorkflowTypes", Array, wfWorkflowTypes)
}

func plistContentItems() {
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
}
