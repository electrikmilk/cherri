/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	"github.com/electrikmilk/args-parser"
	"reflect"
	"strings"
)

const header = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n<dict>\n"
const footer = "</dict>\n</plist>\n"

var plist strings.Builder
var compiled string

func makePlist() {
	if args.Using("debug") {
		fmt.Println("Generating plist...")
	}

	tabLevel = 0
	uuids = make(map[string]string)
	plist.WriteString(header)

	appendPlist([]plistData{
		{
			key:      "WFQuickActionSurfaces",
			dataType: Array,
		},
	})

	plistActions()

	appendPlist([]plistData{
		{
			key:      "WFWorkflowClientVersion",
			dataType: Text,
			value:    "2038.0.2.4",
		},
		{
			key:      "WFWorkflowHasOutputFallback",
			dataType: Boolean,
			value:    false,
		},
		{
			key:      "WFWorkflowHasShortcutInputVariables",
			dataType: Boolean,
			value:    hasShortcutInputVariables,
		},
		{
			key:      "WFWorkflowIcon",
			dataType: Dictionary,
			value: []plistData{
				{
					key:      "WFWorkflowIconGlyphNumber",
					dataType: Number,
					value:    iconGlyph,
				},
				{
					key:      "WFWorkflowIconStartColor",
					dataType: Number,
					value:    iconColor,
				},
			},
		},
		{
			key:      "WFWorkflowImportQuestions",
			dataType: Array,
			value:    plistImportQuestions(),
		},
		{
			key:      "WFWorkflowInputContentItemClasses",
			dataType: Array,
			value:    plistInputContentItems(),
		},
		{
			key:      "WFWorkflowMinimumClientVersion",
			dataType: Number,
			value:    minVersion,
		},
		{
			key:      "WFWorkflowMinimumClientVersionString",
			dataType: Text,
			value:    minVersion,
		},
		{
			key:      "WFWorkflowOutputContentItemClasses",
			dataType: Array,
			value:    plistOutputContentItems(),
		},
		{
			key:      "WFWorkflowTypes",
			dataType: Array,
			value:    plistWorkflowTypes(),
		},
	})

	if noInput.name != "" {
		appendPlist([]plistData{
			{
				key:      "WFWorkflowNoInputBehavior",
				dataType: Dictionary,
				value: []plistData{
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
				},
			},
		})
	}

	if workflowName != "" {
		appendPlist([]plistData{
			{
				key:      "WFWorkflowName",
				dataType: Text,
				value:    workflowName,
			},
		})
	}

	plist.WriteString(footer)

	if args.Using("debug") {
		printPlistGenDebug()
		fmt.Println(ansi("Done.", green) + "\n")
	}

	compiled = plist.String()
	plist.Reset()
	tabLevel = 0
	tokens = []token{}
	menus = map[string][]variableValue{}
	uuids = map[string]string{}
	variables = map[string]variableValue{}
	actions = map[string]*actionDefinition{}
	questions = map[string]*question{}
	globals = map[string]variableValue{}
	noInput = noInputParams{}
	types = []string{}
	inputs = []string{}
	outputs = []string{}
}

func printPlistGenDebug() {
	fmt.Println(ansi("### PLIST GEN ###", bold) + "\n")

	fmt.Println(ansi("## UUIDS ##", bold))
	fmt.Println(uuids)

	fmt.Print("\n")
}

func plistActions() {
	tabLevel++
	var tabs = strings.Repeat("\t", tabLevel)
	plist.WriteString(tabs + "<key>WFWorkflowActions</key>\n" + tabs + "<array>\n")
	for _, t := range tokens {
		switch t.typeof {
		case Var, AddTo, SubFrom, MultiplyBy, DivideBy:
			plistVariable(&t)
		case Comment:
			plistComment(t.value.(string))
		case Action:
			var tokenAction = t.value.(action)
			currentAction = tokenAction.ident
			plistAction(tokenAction.args, &plistData{}, &plistData{})
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
	plist.WriteString(strings.Repeat("\t", tabLevel) + "</array>\n")
	tabLevel--
}

func plistComment(comment string) {
	appendPlist(makeStdAction("comment", []plistData{
		{
			key:      "WFCommentActionText",
			dataType: Text,
			value:    comment,
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
		makeVariableAction(t, &varUUID)
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

	if t.typeof != Var {
		if variables[t.ident].valueType != Arr {
			appendPlist(makeStdAction("setvariable", setVariableParams))
			return
		}

		appendPlist(makeStdAction("appendvariable", setVariableParams))
		return
	}

	if v, found := variables[t.ident]; found {
		if v.constant {
			return
		}
	}
	appendPlist(makeStdAction("setvariable", setVariableParams))

	if t.valueType == Arr {
		plistArrayVariable(t)
	}
}

func plistArrayVariable(t *token) {
	if t.value == nil {
		return
	}
	for _, value := range t.value.([]interface{}) {
		if value == nil {
			continue
		}
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
		makeVariableAction(&token{
			typeof:    valueType,
			ident:     itemIdent,
			valueType: valueType,
			value:     value,
		}, &UUID)
		addToVariableParams = append(addToVariableParams,
			inputValue("WFInput", itemIdent, UUID),
			plistData{
				key:      "WFVariableName",
				dataType: Text,
				value:    t.ident,
			},
		)
		appendPlist(makeStdAction("appendvariable", addToVariableParams))
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
	appendPlist(makeStdAction("conditional", conditionalParams))
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
		if t.valueType != Nil {
			menuParams = append(menuParams, paramValue("WFMenuPrompt", actionArgument{
				valueType: t.valueType,
				value:     t.value,
			}, String, Text))
		}
		var menuItemParams = menus[t.ident]
		var menuItems []plistData
		for _, item := range menuItemParams {
			menuItems = append(menuItems, plistData{
				dataType: Dictionary,
				value: []plistData{
					{
						key:      "WFItemType",
						dataType: Number,
						value:    0,
					},
					paramValue("WFValue", actionArgument{
						valueType: item.valueType,
						value:     item.value,
					}, String, Text),
				},
			})
		}
		menuParams = append(menuParams, plistData{
			key:      "WFMenuItems",
			dataType: Array,
			value:    menuItems,
		})
	}
	appendPlist(makeStdAction("choosefrommenu", menuParams))
}

func plistMenuItem(t *token) {
	appendPlist(makeStdAction("choosefrommenu", []plistData{
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
	var controlFlowMode = startStatement
	if t.valueType == EndClosure {
		controlFlowMode = endStatement
	}
	var repeatData = []plistData{
		{
			key:      "WFControlFlowMode",
			dataType: Number,
			value:    controlFlowMode,
		},
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
	if controlFlowMode == startStatement {
		repeatData = append(repeatData, paramValue("WFRepeatCount", actionArgument{
			valueType: t.valueType,
			value:     t.value,
		}, Integer, Number))
	}
	appendPlist(makeStdAction("repeat.count", repeatData))
}

func plistRepeatEach(t *token) {
	var controlFlowMode = startStatement
	if t.valueType == EndClosure {
		controlFlowMode = endStatement
	}
	var repeatData = []plistData{
		{
			key:      "WFControlFlowMode",
			dataType: Number,
			value:    controlFlowMode,
		},
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
	if controlFlowMode == startStatement {
		repeatData = append(repeatData, paramValue("WFInput", actionArgument{
			valueType: t.valueType,
			value:     t.value,
		}, Variable, Text))
	}
	appendPlist(makeStdAction("repeat.each", repeatData))
}

func plistImportQuestions() (importQuestions []plistData) {
	if len(questions) == 0 {
		return
	}
	for _, q := range questions {
		importQuestions = append(importQuestions, plistData{
			dataType: Dictionary,
			value: []plistData{
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
			},
		})
	}
	return
}

func plistWorkflowTypes() (wfWorkflowTypes []plistData) {
	if len(types) == 0 {
		return []plistData{}
	}

	for _, workflowType := range types {
		wfWorkflowTypes = append(wfWorkflowTypes, plistData{
			dataType: Text,
			value:    workflowType,
		})
	}
	return
}

func plistInputContentItems() (inputContentItems []plistData) {
	if len(inputs) == 0 {
		makeContentItems()
		for _, input := range contentItems {
			inputContentItems = append(inputContentItems, plistData{
				dataType: Text,
				value:    input,
			})
		}
		return
	}

	for _, input := range inputs {
		inputContentItems = append(inputContentItems, plistData{
			dataType: Text,
			value:    input,
		})
	}
	return
}

func plistOutputContentItems() (outputContentItems []plistData) {
	if len(outputs) == 0 {
		return
	}
	makeContentItems()
	for _, output := range outputs {
		outputContentItems = append(outputContentItems, plistData{
			dataType: Text,
			value:    output,
		})
	}

	return
}
