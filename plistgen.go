/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/electrikmilk/args-parser"
	"github.com/google/uuid"
	plists "howett.net/plist"
)

const header = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n<dict>\n"
const footer = "</dict>\n</plist>\n"

var plist strings.Builder
var compiled string

var longEmptyArraySyntax = regexp.MustCompile(`<array>\n(.*?)</array>`)

func marshalPlist() {
	// type sparseBundleHeader struct {
	// 	InfoDictionaryVersion string `plist:"CFBundleInfoDictionaryVersion"`
	// 	BandSize              uint64 `plist:"band-size"`
	// 	BackingStoreVersion   int    `plist:"bundle-backingstore-version"`
	// 	DiskImageBundleType   string `plist:"diskimage-bundle-type"`
	// 	Size                  uint64 `plist:"size"`
	// }
	// data := &sparseBundleHeader{
	// 	InfoDictionaryVersion: "6.0",
	// 	BandSize:              8388608,
	// 	Size:                  4 * 1048576 * 1024 * 1024,
	// 	DiskImageBundleType:   "com.apple.diskimage.sparsebundle",
	// 	BackingStoreVersion:   1,
	// }

	generateShortcut()

	var plist, plistErr = plists.MarshalIndent(shortcut, plists.XMLFormat, "\t")
	handle(plistErr)

	compiled = longEmptyArraySyntax.ReplaceAllString(string(plist), "<array/>")

	resetPlistGen()
}

func generateShortcut() {
	if args.Using("debug") {
		fmt.Print("Generating plist data...")
	}
	shortcut = Shortcut{
		WFWorkflowIcon: ShortcutIcon{
			iconGlyph,
			iconColor,
		},
		WFWorkflowClientVersion:              clientVersion,
		WFWorkflowHasShortcutInputVariables:  hasShortcutInputVariables,
		WFWorkflowImportQuestions:            plistImportQuestions(),
		WFWorkflowInputContentItemClasses:    plistInputContentItems(),
		WFWorkflowOutputContentItemClasses:   plistOutputContentItems(),
		WFWorkflowMinimumClientVersion:       900,
		WFWorkflowMinimumClientVersionString: "900",
		WFWorkflowTypes:                      plistWorkflowTypes(),
		WFWorkflowNoInputBehavior:            noInput,
	}

	if workflowName != "" {
		shortcut.WFWorkflowName = workflowName
	}

	generatePlistActions()

	if args.Using("debug") {
		printPlistGenDebug()
		fmt.Println(ansi("Done.\n", green))
	}
}

func resetPlistGen() {
	tabLevel = 0
	tokens = []token{}
	menus = map[string][]variableValue{}
	uuids = map[string]string{}
	variables = map[string]variableValue{}
	questions = map[string]*question{}
	noInput = WFWorkflowNoInputBehavior{}
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

func generatePlistActions() {
	uuids = make(map[string]string)
	for _, t := range tokens {
		switch t.typeof {
		case Var, AddTo, SubFrom, MultiplyBy, DivideBy:
			plistVariable(&t)
		case Comment:
			plistComment(t.value.(string))
		case Action:
			var tokenAction = t.value.(action)
			setCurrentAction(tokenAction.ident, actions[tokenAction.ident])
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
		var varUUID = uuid.New().String()
		var outputName = makeOutputName(t)
		makeVariableAction(t, &outputName, &varUUID)
		uuids[outputName] = varUUID
		if t.valueType != Arr {
			if t.typeof == Var && t.valueType == Variable {
				setVariableParams = append(setVariableParams, variablePlistValue("WFInput", t.value.(string), t.ident))
			} else {
				setVariableParams = append(setVariableParams, inputValue("WFInput", outputName, varUUID))
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

func makeOutputName(token *token) string {
	if variable, found := variables[token.ident]; found {
		if variable.constant {
			return token.ident
		}
	}
	if token.valueType == Var {
		if _, found := variables[token.value.(string)]; found {
			return token.value.(string)
		}
		if _, found := globals[token.value.(string)]; found {
			return token.value.(string)
		}
	}
	var customOutputName = fmt.Sprintf("%sValue", string(token.valueType))
	if customOutputName == "action" {
		customOutputName = token.value.(action).ident
	}

	return checkDuplicateOutputName(customOutputName)
}

func plistArrayVariable(t *token) {
	if t.value == nil {
		return
	}
	for _, value := range t.value.([]interface{}) {
		if value == nil {
			continue
		}
		var UUID = uuid.New().String()
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
		}, &itemIdent, &UUID)
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
	var controlFlowMode uint64
	var conditionalParams = []plistData{
		{
			key:      "GroupingIdentifier",
			dataType: Text,
			value:    t.ident,
		},
		{
			key:      "UUID",
			dataType: Text,
			value:    uuid.New().String(),
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
			value:    uuid.New().String(),
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
			value:    uuid.New().String(),
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

func plistWorkflowTypes() (wfWorkflowTypes []string) {
	if len(types) == 0 {
		return
	}

	for _, workflowType := range types {
		wfWorkflowTypes = append(wfWorkflowTypes, workflowType)
	}
	return
}

func plistInputContentItems() (inputContentItems []string) {
	if len(inputs) == 0 {
		for _, input := range contentItems {
			inputContentItems = append(inputContentItems, input)
		}
		return
	}

	for _, input := range inputs {
		inputContentItems = append(inputContentItems, input)
	}
	return
}

func plistOutputContentItems() (outputContentItems []string) {
	if len(outputs) == 0 {
		return
	}
	for _, output := range outputs {
		outputContentItems = append(outputContentItems, output)
	}

	return
}
