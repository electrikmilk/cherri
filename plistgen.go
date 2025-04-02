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

var plist strings.Builder
var compiled string

var longEmptyArraySyntax = regexp.MustCompile(`<array>\n(.*?)</array>`)
var longEmptyDictSyntax = regexp.MustCompile(`<dict>\n(.*?)</dict>`)

func marshalPlist() {
	generateShortcut()

	var plist, plistErr = plists.MarshalIndent(shortcut, plists.XMLFormat, "\t")
	handle(plistErr)

	compiled = longEmptyArraySyntax.ReplaceAllString(string(plist), "<array/>")
	compiled = longEmptyDictSyntax.ReplaceAllString(compiled, "<dict/>")

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
			plistAction(tokenAction.args, &map[string]any{})
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
	buildStdAction("comment", map[string]any{
		"WFCommentActionText": comment,
	})
}

func plistVariable(t *token) {
	var setVariableParams = map[string]any{
		"WFVariableName": t.ident,
	}

	if t.value != nil {
		var outputName = makeOutputName(t)
		var varUUID string
		if uuids[outputName] == "" {
			varUUID = uuid.New().String()
			uuids[outputName] = varUUID
		} else {
			varUUID = uuids[outputName]
		}

		makeVariableAction(t, &outputName, &varUUID)
		if t.valueType != Arr {
			if t.typeof == Var && t.valueType == Variable {
				setVariableParams["WFInput"] = variablePlistValue(t.value.(string), t.ident)
			} else {
				setVariableParams["WFInput"] = inputValue(outputName, varUUID)
			}

			setVariableParams["WFSerializationType"] = "WFTextTokenAttachment"
		}
	}

	if t.typeof != Var {
		if variables[t.ident].valueType != Arr {
			buildStdAction("setvariable", setVariableParams)
			return
		}

		buildStdAction("appendvariable", setVariableParams)
		return
	}

	if v, found := variables[t.ident]; found {
		if v.constant {
			return
		}
	}
	buildStdAction("setvariable", setVariableParams)

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
	var typeOfToken = string(token.valueType)
	if typeOfToken == "action" {
		typeOfToken = token.value.(action).ident
	}

	var customOutputName = fmt.Sprintf("%s%sOutput", strings.ToTitle(string(typeOfToken[0])), typeOfToken[1:])

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

		buildStdAction("appendvariable", map[string]any{
			"WFInput":        inputValue(itemIdent, UUID),
			"WFVariableName": t.ident,
		})
	}
}

func plistConditional(t *token) {
	var conditionalParams = map[string]any{
		"GroupingIdentifier": t.ident,
		"UUID":               uuid.New().String(),
	}
	switch t.valueType {
	case If:
		var cond = t.value.(condition)
		conditionalParams["WFInput"] = map[string]any{
			"Type":     "Variable",
			"Variable": variablePlistValue(cond.variableOneValue.(string), t.ident),
		}
		if cond.variableTwoValue != nil {
			conditionalParameter("", conditionalParams, &cond.variableTwoType, cond.variableTwoValue)
		}
		if cond.variableThreeValue != nil {
			conditionalParameter("WFAnotherNumber", conditionalParams, &cond.variableThreeType, cond.variableThreeValue)
		}
		conditionalParams["WFCondition"] = cond.condition
		conditionalParams["WFControlFlowMode"] = startStatement
	case Else:
		conditionalParams["WFControlFlowMode"] = statementPart
	case EndClosure:
		conditionalParams["WFControlFlowMode"] = endStatement
	}

	buildStdAction("conditional", conditionalParams)
}

func plistMenu(t *token) {
	var controlFlowMode = startStatement
	if t.valueType == EndClosure {
		controlFlowMode = endStatement
	}
	var menuParams = map[string]any{
		"GroupingIdentifier": t.ident,
		"WFControlFlowMode":  controlFlowMode,
	}
	if t.valueType != EndClosure {
		if t.valueType != Nil {
			menuParams["WFMenuPrompt"] = paramValue(actionArgument{
				valueType: t.valueType,
				value:     t.value,
			}, String)
		}
		var menuItemParams = menus[t.ident]
		var menuItems []map[string]any
		for _, item := range menuItemParams {
			menuItems = append(menuItems, map[string]any{
				"WFItemType": 0,
				"WFValue": paramValue(actionArgument{
					valueType: item.valueType,
					value:     item.value,
				}, String),
			})
		}

		menuParams["WFMenuItems"] = menuItems
	}

	buildStdAction("choosefrommenu", menuParams)
}

func plistMenuItem(t *token) {
	buildStdAction("choosefrommenu", map[string]any{
		"GroupingIdentifier": t.ident,
		"WFControlFlowMode":  statementPart,
		"WFMenuItemAttributedTitle": paramValue(actionArgument{
			valueType: t.valueType,
			value:     t.value,
		}, String),
		"WFMenuItemTitle": paramValue(actionArgument{
			valueType: t.valueType,
			value:     t.value,
		}, String),
	})
}

func plistRepeat(t *token) {
	var controlFlowMode = startStatement
	if t.valueType == EndClosure {
		controlFlowMode = endStatement
	}
	var repeatParams = map[string]any{
		"WFControlFlowMode":  controlFlowMode,
		"GroupingIdentifier": t.ident,
		"UUID":               uuid.New().String(),
	}
	if controlFlowMode == startStatement {
		repeatParams["WFRepeatCount"] = paramValue(actionArgument{
			valueType: t.valueType,
			value:     t.value,
		}, Integer)
	}

	buildStdAction("repeat.count", repeatParams)
}

func plistRepeatEach(t *token) {
	var controlFlowMode = startStatement
	if t.valueType == EndClosure {
		controlFlowMode = endStatement
	}
	var repeatEachParams = map[string]any{
		"WFControlFlowMode":  controlFlowMode,
		"GroupingIdentifier": t.ident,
		"UUID":               uuid.New().String(),
	}
	if controlFlowMode == startStatement {
		repeatEachParams["WFInput"] = paramValue(actionArgument{
			valueType: t.valueType,
			value:     t.value,
		}, Integer)
	}

	buildStdAction("repeat.each", repeatEachParams)
}

func plistImportQuestions() (importQuestions []map[string]any) {
	if len(questions) == 0 {
		return
	}

	for _, q := range questions {
		importQuestions = append(importQuestions, map[string]any{
			"ParameterKey": q.parameter,
			"Category":     "Parameter",
			"ActionIndex":  q.actionIndex,
			"Text":         q.text,
			"DefaultValue": q.defaultValue,
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
