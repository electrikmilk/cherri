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
		WFWorkflowMinimumClientVersion:       900,
		WFWorkflowMinimumClientVersionString: "900",
		WFWorkflowTypes:                      types,
		WFWorkflowNoInputBehavior:            noInput,
	}

	waitFor(
		func() {
			shortcut.WFWorkflowInputContentItemClasses = generateInputContentItems()
		},
		func() {
			shortcut.WFWorkflowOutputContentItemClasses = generateOutputContentItems()
		},
		func() {
			shortcut.WFWorkflowImportQuestions = generateImportQuestions()
		},
	)

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
		case Variable, AddTo, SubFrom, MultiplyBy, DivideBy:
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
			if t.typeof == Variable && t.valueType == Variable {
				setVariableParams["WFInput"] = variableValue(t.value.(varValue))
			} else {
				setVariableParams["WFInput"] = inputValue(outputName, varUUID)
			}

			setVariableParams["WFSerializationType"] = "WFTextTokenAttachment"
		}
	}

	if t.typeof != Variable {
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
		makeArrayVariable(t)
	}
}

func makeOutputName(token *token) string {
	if variable, found := variables[token.ident]; found {
		if variable.constant {
			return token.ident
		}
	}
	if token.valueType == Variable {
		var identifer = token.value.(varValue).value.(string)
		if validReference(identifer) {
			return identifer
		}
	}
	var typeOfToken = string(token.valueType)
	if typeOfToken == "action" {
		typeOfToken = token.value.(action).ident
	}

	var customOutputName = fmt.Sprintf("%s%s", strings.ToTitle(string(typeOfToken[0])), typeOfToken[1:])

	return checkDuplicateOutputName(customOutputName)
}

func makeArrayVariable(t *token) {
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
		conditionalParams["WFControlFlowMode"] = startStatement

		var cond = t.value.(WFConditions)
		conditionalParams["WFConditions"] = makeConditions(&cond)
	case Else:
		conditionalParams["WFControlFlowMode"] = statementPart
	case EndClosure:
		conditionalParams["WFControlFlowMode"] = endStatement
	}

	buildStdAction("conditional", conditionalParams)
}

func makeConditions(wfConditions *WFConditions) map[string]any {
	var filterTemplates []map[string]any
	for _, condition := range wfConditions.conditions {
		var conditionParams = map[string]any{
			"WFCondition": condition.condition,
			"WFInput": map[string]any{
				"Type":     "Variable",
				"Variable": variableValue(condition.arguments[0].value.(varValue)),
			},
		}

		if len(condition.arguments) > 1 {
			var argumentTwo = condition.arguments[1]
			conditionalParameter("", conditionParams, &argumentTwo.valueType, argumentTwo.value)
		}
		if len(condition.arguments) > 2 {
			var argumentThree = condition.arguments[2]
			conditionalParameter("WFAnotherNumber", conditionParams, &argumentThree.valueType, argumentThree.value)
		}

		filterTemplates = append(filterTemplates, conditionParams)
	}

	var conditionParams = map[string]any{
		"WFSerializationType": "WFContentPredicateTableTemplate",
		"Value": map[string]any{
			"WFActionParameterFilterPrefix":    wfConditions.WFActionParameterFilterPrefix,
			"WFActionParameterFilterTemplates": filterTemplates,
		},
	}

	return conditionParams
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

func generateImportQuestions() (importQuestions []map[string]any) {
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

func generateInputContentItems() (inputContentItems []string) {
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

func generateOutputContentItems() (outputContentItems []string) {
	if len(outputs) == 0 {
		return
	}
	for _, output := range outputs {
		outputContentItems = append(outputContentItems, output)
	}

	return
}
