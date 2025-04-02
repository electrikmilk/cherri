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
	"howett.net/plist"
)

var compiled string

var longEmptyArraySyntax = regexp.MustCompile(`<array>\n(.*?)</array>`)
var longEmptyDictSyntax = regexp.MustCompile(`<dict>\n(.*?)</dict>`)

func generateShortcut() {
	if args.Using("debug") {
		fmt.Print("Generating Shortcut data...")
	}

	shortcut = Shortcut{
		WFWorkflowIcon: ShortcutIcon{
			iconGlyph,
			iconColor,
		},
		WFWorkflowClientVersion:              clientVersion,
		WFWorkflowHasShortcutInputVariables:  hasShortcutInputVariables,
		WFWorkflowImportQuestions:            shortcutImportQuestions(),
		WFWorkflowInputContentItemClasses:    inputContentItems(),
		WFWorkflowOutputContentItemClasses:   outputContentItems(),
		WFWorkflowMinimumClientVersion:       900,
		WFWorkflowMinimumClientVersionString: "900",
		WFWorkflowTypes:                      types,
		WFWorkflowNoInputBehavior:            noInput,
	}

	if workflowName != "" {
		shortcut.WFWorkflowName = workflowName
	}

	generateActions()

	marshalPlist()

	if args.Using("debug") {
		printShortcutGenDebug()
		fmt.Println(ansi("Done.\n", green))
	}

	resetShortcutGen()
}

func marshalPlist() {
	var marshaledPlist, plistErr = plist.MarshalIndent(shortcut, plist.XMLFormat, "\t")
	handle(plistErr)

	compiled = string(marshaledPlist)
	compiled = longEmptyArraySyntax.ReplaceAllString(compiled, "<array/>")
	compiled = longEmptyDictSyntax.ReplaceAllString(compiled, "<dict/>")
}

func resetShortcutGen() {
	tokens = []token{}
	menus = map[string][]varValue{}
	uuids = map[string]string{}
	variables = map[string]varValue{}
	questions = map[string]*question{}
	noInput = WFWorkflowNoInputBehavior{}
	types = []string{}
	inputs = []string{}
	outputs = []string{}
}

func printShortcutGenDebug() {
	fmt.Println(ansi("### SHORTCUT GEN ###", bold) + "\n")

	fmt.Println(ansi("## UUIDS ##", bold))
	fmt.Println(uuids)

	fmt.Print("\n")
}

func generateActions() {
	uuids = make(map[string]string)
	for _, t := range tokens {
		switch t.typeof {
		case Var, AddTo, SubFrom, MultiplyBy, DivideBy:
			makeVariableAction(&t)
		case Comment:
			makeCommentAction(t.value.(string))
		case Action:
			var tokenAction = t.value.(action)
			setCurrentAction(tokenAction.ident, actions[tokenAction.ident])
			makeAction(tokenAction.args, &map[string]any{})
		case Repeat:
			makeRepeatAction(&t)
		case RepeatWithEach:
			makeRepeatEachAction(&t)
		case Menu:
			makeMenuAction(&t)
		case Item:
			makeMenuItemAction(&t)
		case Conditional:
			makeConditionalAction(&t)
		}
	}
}

func makeCommentAction(comment string) {
	buildStdAction("comment", map[string]any{
		"WFCommentActionText": comment,
	})
}

func makeVariableAction(t *token) {
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

		makeVariableValueAction(t, &outputName, &varUUID)
		if t.valueType != Arr {
			if t.typeof == Var && t.valueType == Variable {
				setVariableParams["WFInput"] = variableValue(t.value.(string), t.ident)
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
		shortcutArrayVariable(t)
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

func shortcutArrayVariable(t *token) {
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
		makeVariableValueAction(&token{
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

func makeConditionalAction(t *token) {
	var conditionalParams = map[string]any{
		"GroupingIdentifier": t.ident,
		"UUID":               uuid.New().String(),
	}
	switch t.valueType {
	case If:
		var cond = t.value.(condition)
		conditionalParams["WFInput"] = map[string]any{
			"Type":     "Variable",
			"Variable": variableValue(cond.variableOneValue.(string), t.ident),
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

func makeMenuAction(t *token) {
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

func makeMenuItemAction(t *token) {
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

func makeRepeatAction(t *token) {
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

func makeRepeatEachAction(t *token) {
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

func shortcutImportQuestions() (importQuestions []map[string]any) {
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

func inputContentItems() (inputContentItems []string) {
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

func outputContentItems() (outputContentItems []string) {
	if len(outputs) == 0 {
		return
	}
	for _, output := range outputs {
		outputContentItems = append(outputContentItems, output)
	}

	return
}
