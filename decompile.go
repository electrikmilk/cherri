/*
 * Copyright (c) Cherri
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/electrikmilk/args-parser"
	"howett.net/plist"
)

const (
	UUID          = "UUID"
	ShortcutInput = "ShortcutInput"
	Ask           = "Ask"
)

var code strings.Builder
var specialCharsRegex *regexp.Regexp

func decompile(b []byte) {
	var _, marshalIndexedErr = plist.Unmarshal(b, &shortcut)
	handle(marshalIndexedErr)

	variables = make(map[string]varValue)
	uuids = make(map[string]string)
	standardActionIncludes = make(map[string]string)
	decompiledIncludes = []string{}
	controlFlowOutputs = make(map[string]string)
	controlFlowOutputUUIDs = []string{}
	openLoops = nil
	loopIdentifiers = make(map[string]string)
	currentVariableValue = ""
	currentVariableUUID = ""

	basename = strings.ReplaceAll(basename, " ", "_")
	outputPath = getOutputPath(basename + ".cherri")
	if args.Using("no-ansi") {
		filePath = basename + ".cherri"
	} else {
		filePath = outputPath
	}

	loadBasicStandardActions()
	resetParse()
	firstChar()

	mapVariables()
	mapSplitActions()
	mapIdentifiers()
	mapControlFlowOutputs()
	mapLoopVariables()

	defineName()
	decompileIcon()

	decompileActions()

	if args.Using("debug") {
		printDecompDebug()
	}

	var writeErr = os.WriteFile(outputPath, []byte(code.String()), 0600)
	handle(writeErr)
}

// mapIdentifiers creates a map of variable identifiers and UUIDs that are assigned throughout the Shortcut.
func mapIdentifiers() {
	for _, action := range shortcut.WFWorkflowActions {
		currentAction.identifier = action.WFWorkflowActionIdentifier
		var params = action.WFWorkflowActionParameters
		if action.WFWorkflowActionIdentifier == SetVariableIdentifier || action.WFWorkflowActionIdentifier == AppendVariableIdentifier {
			continue
		}

		if params["UUID"] != nil && params["CustomOutputName"] != nil {
			mapUUID(paramString(&action, "UUID"), paramString(&action, "CustomOutputName"))
		}

		checkParamIdentifiers(params, &action)
	}
}

// controlFlowOutputs maps the GroupingIdentifier of a control flow block to the name its output is referenced by.
var controlFlowOutputs map[string]string
var controlFlowOutputUUIDs []string

var controlFlowIdentifiers = []string{
	"is.workflow.actions.conditional",
	"is.workflow.actions.choosefrommenu",
	"is.workflow.actions.repeat.count",
	"is.workflow.actions.repeat.each",
}

// mapControlFlowOutputs names the output of each control flow block from its end action, which is the only one
// that carries the UUID, so the start action that shares its GroupingIdentifier can be written as a constant.
func mapControlFlowOutputs() {
	for _, action := range shortcut.WFWorkflowActions {
		if !slices.Contains(controlFlowIdentifiers, action.WFWorkflowActionIdentifier) {
			continue
		}
		if controlFlowMode(&action) != endStatement || action.WFWorkflowActionParameters[UUID] == nil {
			continue
		}

		var outputUUID = paramString(&action, UUID)
		var outputName, referenced = uuids[outputUUID]
		if !referenced {
			continue
		}

		controlFlowOutputs[paramString(&action, "GroupingIdentifier")] = outputName
		controlFlowOutputUUIDs = append(controlFlowOutputUUIDs, outputUUID)
	}
}

const repeatIdentifier = "is.workflow.actions.repeat.count"
const repeatEachIdentifier = "is.workflow.actions.repeat.each"

type openLoop struct {
	groupingUUID string
	variableName string
	identifier   string
}

var openLoops []openLoop

// loopIdentifiers maps the GroupingIdentifier of each loop whose variable is referenced to the identifier the
// loop declares, so a loop whose variable is never referenced can still be written with `_`.
var loopIdentifiers map[string]string

var loopVariableRegex = regexp.MustCompile("^Repeat (Index|Item)( [0-9]+)?$")

// mapLoopVariables finds the loops whose variable is referenced before any code is written, because the loop
// identifier is written before the references in its body.
func mapLoopVariables() {
	for _, action := range shortcut.WFWorkflowActions {
		currentAction.identifier = action.WFWorkflowActionIdentifier
		var isLoop = action.WFWorkflowActionIdentifier == repeatIdentifier || action.WFWorkflowActionIdentifier == repeatEachIdentifier
		if isLoop && controlFlowMode(&action) == endStatement {
			popLoop(&action)
			continue
		}

		markLoopReferences(action.WFWorkflowActionParameters)

		if isLoop && controlFlowMode(&action) == startStatement {
			pushLoop(&action)
		}
	}
	openLoops = nil
}

func markLoopReferences(value any) {
	switch typedValue := value.(type) {
	case map[string]interface{}:
		if variableName, isString := typedValue["VariableName"].(string); isString && typedValue["Type"] == "Variable" {
			if loop, found := findOpenLoop(variableName); found {
				loopIdentifiers[loop.groupingUUID] = loop.identifier
			}
		}
		for _, item := range typedValue {
			markLoopReferences(item)
		}
	case []interface{}:
		for _, item := range typedValue {
			markLoopReferences(item)
		}
	}
}

// pushLoop opens a loop. The compiler numbers the variable of a repeat by how many repeats are open, and the
// variable of a repeat with each by how many of those are open, so each kind is counted on its own.
func pushLoop(action *ShortcutAction) openLoop {
	var variableName = "Repeat Index"
	var identifier = "repeatIndex"
	if action.WFWorkflowActionIdentifier == repeatEachIdentifier {
		variableName = "Repeat Item"
		identifier = "repeatItem"
	}

	var depth = 1
	for _, loop := range openLoops {
		if strings.HasPrefix(loop.variableName, variableName) {
			depth++
		}
	}
	if depth > 1 {
		variableName = fmt.Sprintf("%s %d", variableName, depth)
		identifier = fmt.Sprintf("%s%d", identifier, depth)
	}

	var loop = openLoop{
		groupingUUID: paramString(action, "GroupingIdentifier"),
		variableName: variableName,
		identifier:   identifier,
	}
	openLoops = append(openLoops, loop)

	return loop
}

func popLoop(action *ShortcutAction) {
	var groupingUUID = paramString(action, "GroupingIdentifier")
	if len(openLoops) == 0 || openLoops[len(openLoops)-1].groupingUUID != groupingUUID {
		decompError(fmt.Sprintf("Expected loop '%s' to end the innermost open loop", groupingUUID), action)
	}
	openLoops = openLoops[:len(openLoops)-1]
}

func findOpenLoop(variableName string) (openLoop, bool) {
	for i := len(openLoops) - 1; i >= 0; i-- {
		if openLoops[i].variableName == variableName {
			return openLoops[i], true
		}
	}

	return openLoop{}, false
}

// loopVariableReference returns how a variable name that Shortcuts gives a loop is referenced in Cherri.
// Outside a loop that defines it, only the first index and item can be referenced, as globals.
func loopVariableReference(variableName string) (reference string, isLoopVariable bool) {
	if !loopVariableRegex.MatchString(variableName) {
		return "", false
	}
	if loop, found := findOpenLoop(variableName); found {
		return fmt.Sprintf("@%s", loop.identifier), true
	}
	for name, global := range globals {
		if global.value == variableName {
			return name, true
		}
	}

	decompError(fmt.Sprintf("Expected the variable '%s' to be referenced inside a loop that defines it", variableName), currentDecompAction())
	return "", false
}

// loopDeclaration opens a loop and returns the identifier it declares.
func loopDeclaration(action *ShortcutAction) string {
	var loop = pushLoop(action)
	if _, referenced := loopIdentifiers[loop.groupingUUID]; !referenced {
		return "_"
	}

	return loop.identifier
}

// Map out variables in the Shortcut and their UUIDs for later checks.
func mapVariables() {
	for _, action := range shortcut.WFWorkflowActions {
		currentAction.identifier = action.WFWorkflowActionIdentifier
		if action.WFWorkflowActionIdentifier == SetVariableIdentifier || action.WFWorkflowActionIdentifier == AppendVariableIdentifier {
			var varName = paramString(&action, "WFVariableName")
			sanitizeIdentifier(&varName)
			if _, found := variables[varName]; !found {
				variables[varName] = varValue{}
			}

			if action.WFWorkflowActionParameters["WFInput"] != nil {
				var wfInput WFInput
				mapToStruct(action.WFWorkflowActionParameters["WFInput"], &wfInput)
				if _, found := uuids[wfInput.Value.OutputUUID]; !found {
					varUUIDs = append(varUUIDs, wfInput.Value.OutputUUID)
				}
			}
		}
	}
}

func peekActions(peek int) ShortcutAction {
	return shortcut.WFWorkflowActions[actionIndex+peek]
}

// insertCodeComment creates a comment in the decompiled code.
func insertCodeComment(comment string) {
	newCodeLine(fmt.Sprintf("// %s\n", comment))
}

func checkParamIdentifiers(params map[string]interface{}, action *ShortcutAction) {
	for _, value := range params {
		var paramValues, isMap = value.(map[string]interface{})
		if !isMap {
			continue
		}
		checkParamValueAttachments(paramValues)

		// A dictionary value can hold any key, so a "Variable" that is not a dictionary is not an attachment.
		if paramVariable, isMap := paramValues["Variable"].(map[string]interface{}); isMap {
			checkParamValueAttachments(paramVariable)
		}

		if _, found := paramValues["WFConditions"]; found {
			var wfConditions, hasConditions = paramValues["WFConditions"].(map[string]interface{})
			if !hasConditions {
				decompError(fmt.Sprintf("Expected WFConditions to be a dictionary, got %T", paramValues["WFConditions"]), action)
			}
			var conditionsValue, hasValue = wfConditions["Value"].(map[string]interface{})
			if !hasValue {
				decompError("Expected WFConditions to have a Value dictionary", action)
			}
			if conditionsValue["WFActionParameterFilterTemplates"] != nil {
				for _, filterTemplate := range conditionTemplates(conditionsValue, action) {
					if filterTemplate["WFInput"] == nil {
						continue
					}
					var templateInput, hasInput = filterTemplate["WFInput"].(map[string]interface{})
					if !hasInput {
						decompError(fmt.Sprintf("Expected WFInput of a condition to be a dictionary, got %T", filterTemplate["WFInput"]), action)
					}
					checkParamIdentifiers(templateInput, action)
				}
			}
		}
	}
}

// conditionTemplates returns the conditions in the Value of a WFConditions parameter.
func conditionTemplates(conditionsValue map[string]interface{}, action *ShortcutAction) (templates []map[string]interface{}) {
	var templateValues, isSlice = conditionsValue["WFActionParameterFilterTemplates"].([]interface{})
	if !isSlice {
		decompError(fmt.Sprintf("Expected WFActionParameterFilterTemplates to be an array, got %T", conditionsValue["WFActionParameterFilterTemplates"]), action)
	}

	for _, templateValue := range templateValues {
		var template, isMap = templateValue.(map[string]interface{})
		if !isMap {
			decompError(fmt.Sprintf("Expected a condition to be a dictionary, got %T", templateValue), action)
		}
		templates = append(templates, template)
	}

	return
}

// checkParamValueAttachments checks for attachments on values to map them.
func checkParamValueAttachments(param map[string]interface{}) {
	var paramValue, isMap = param["Value"].(map[string]interface{})
	if !isMap {
		return
	}

	var inputValue Value
	mapToStruct(paramValue, &inputValue)
	mapValueReference(inputValue)
	checkConstantUUID(inputValue)

	if inputValue.AttachmentsByRange != nil {
		mapAttachmentIdentifiers(inputValue.AttachmentsByRange)
	}
	if inputValue.WFDictionaryFieldValueItems != nil {
		mapDictionaryValueIdentifiers(inputValue.WFDictionaryFieldValueItems)
	}
}

func mapDictionaryValueIdentifiers(items []WFDictionaryFieldValueItem) {
	for _, item := range items {
		var wfValue, isMap = item.WFValue.(map[string]interface{})
		if !isMap {
			continue
		}
		if wfValue["Value"] == nil {
			continue
		}

		// A file item wraps its attachment one level deeper than other items, which would otherwise be read as a nested dictionary.
		if dictDataType(item.WFItemType) == itemTypeFile {
			var fileValue, isMap = wfValue["Value"].(map[string]interface{})
			if !isMap {
				decompError(fmt.Sprintf("Expected a file dictionary item to hold a dictionary, got %T", wfValue["Value"]), currentDecompAction())
			}
			checkParamValueAttachments(fileValue)
			continue
		}

		checkParamValueAttachments(wfValue)

		switch valueTwo := wfValue["Value"].(type) {
		case map[string]interface{}:
			if itemsValue, isMap := valueTwo["Value"].(map[string]interface{}); isMap {
				var dictionaryItems []WFDictionaryFieldValueItem
				mapToStruct(itemsValue["WFDictionaryFieldValueItems"], &dictionaryItems)
				mapDictionaryValueIdentifiers(dictionaryItems)
			}
		case []interface{}:
			var dictionaryItems = mapDictionaryItems(valueTwo)
			mapDictionaryValueIdentifiers(dictionaryItems)
		}
	}
}

func mapDictionaryItems(items []interface{}) (dictionaryItems []WFDictionaryFieldValueItem) {
	for _, value := range items {
		if _, isMap := value.(map[string]interface{}); !isMap {
			continue
		}

		var dictionaryItem WFDictionaryFieldValueItem
		mapToStruct(value, &dictionaryItem)
		dictionaryItems = append(dictionaryItems, dictionaryItem)
	}
	return
}

// mapAttachmentIdentifiers maps the UUID and output name of an attachment value.
func mapAttachmentIdentifiers(attachments map[string]Value) {
	for _, attachmentValue := range attachments {
		mapValueReference(attachmentValue)
		checkConstantUUID(attachmentValue)
	}
}

func checkConstantUUID(value Value) {
	if value.OutputUUID != "" && value.Type == "ActionOutput" {
		constUUIDs = append(constUUIDs, value.OutputUUID)
	}
}

func mapValueReference(value Value) {
	if value.OutputName == "" || value.OutputUUID == "" {
		return
	}
	mapUUID(value.OutputUUID, value.OutputName)
}

func mapUUID(uuid string, varName string) {
	var outputName string
	if _, found := uuids[uuid]; !found {
		outputName = varName
		sanitizeIdentifier(&outputName)
		uuids[uuid] = checkDuplicateOutputName(outputName)
	}
}

// sanitizeIdentifier strips special characters and replaces dashes with underscores.
func sanitizeIdentifier(identifier *string) {
	if specialCharsRegex == nil {
		specialCharsRegex = regexp.MustCompile("[^a-zA-Z0-9_]+")
	}
	if strings.Contains(*identifier, " ") {
		var words = strings.Split(strings.TrimSpace(*identifier), " ")
		for i, word := range words {
			words[i] = capitalize(word)
		}
		*identifier = strings.Join(words, "")
	}

	*identifier = specialCharsRegex.ReplaceAllString(*identifier, "")
	*identifier = strings.ReplaceAll(*identifier, "-", "_")
}

type actionValue struct {
	identifier string
	definition *actionDefinition
}

var identifierMap map[string][]actionValue

// mapSplitActions creates a map of actions that have been split into a few actions to reduce the number of arguments.
func mapSplitActions() {
	if identifierMap == nil {
		identifierMap = make(map[string][]actionValue)
	}
	for identifier, action := range actions {
		var ident = action.identifier
		if action.identifier == "" {
			ident = identifier
		}

		ident = strings.ToLower(ident)

		var splitActionValue = actionValue{
			identifier: identifier,
			definition: action,
		}
		if identifierMap[ident] != nil && slices.Contains(identifierMap[ident], splitActionValue) {
			continue
		}

		identifierMap[ident] = append(identifierMap[ident], splitActionValue)
	}
	for identifier, actions := range identifierMap {
		if len(actions) < 2 {
			delete(identifierMap, identifier)
			continue
		}
		sort.Slice(actions, func(i, j int) bool {
			return actions[i].identifier < actions[j].identifier
		})
		identifierMap[identifier] = actions
	}
}

// newCodeLine writes s on a new line in the code.
func newCodeLine(s string) {
	if tabLevel > 0 {
		for i := 0; i < tabLevel; i++ {
			code.WriteRune('\t')
		}
	}
	code.WriteString(s)
}

// defineName writes the name of the imported Shortcut to code as a name definition if the basename contains a space.
func defineName() {
	if strings.Contains(basename, " ") {
		newCodeLine(fmt.Sprintf("#define name %s\n", basename))
	}
}

func decompileIcon() {
	var icon = shortcut.WFWorkflowIcon
	var hasIcon = icon.WFWorkflowIconStartColor != iconColor || icon.WFWorkflowIconGlyphNumber != iconGlyph

	if icon.WFWorkflowIconStartColor != iconColor {
		defineColors(icon, colors)
		defineColors(icon, altColors)
	}

	if icon.WFWorkflowIconGlyphNumber != iconGlyph {
		for name, i := range glyphs {
			if icon.WFWorkflowIconGlyphNumber != int64(i) {
				continue
			}

			newCodeLine(fmt.Sprintf("#define glyph %s\n", name))
		}
	}

	if hasIcon {
		newCodeLine("\n")
	}
}

func defineColors(icon ShortcutIcon, colors map[string]int) {
	for name, i := range colors {
		if icon.WFWorkflowIconStartColor != i {
			continue
		}

		newCodeLine(fmt.Sprintf("#define color %s\n", name))
	}
}

var currentVariableValue string

// currentVariableUUID is the UUID of the action that produced currentVariableValue, which only the Set Variable
// referencing that UUID may consume.
var currentVariableUUID string

func decompileActions() {
	for _, action := range shortcut.WFWorkflowActions {
		currentAction.identifier = action.WFWorkflowActionIdentifier
		switch action.WFWorkflowActionIdentifier {
		case "is.workflow.actions.comment":
			decompComment(&action)
		case "is.workflow.actions.gettext":
			decompTextValue(&action)
		case "is.workflow.actions.number":
			var nonLiteral = decompNumberValue(&action)
			if nonLiteral {
				decompAction(&action)
			}
		case "is.workflow.actions.dictionary":
			decompDictionary(&action)
		case "is.workflow.actions.math":
			decompBasicExpression(&action)
		case "is.workflow.actions.calculateexpression":
			decompExpression(&action)
		case SetVariableIdentifier, AppendVariableIdentifier:
			decompVariable(&action)
		case "is.workflow.actions.conditional":
			decompConditional(&action)
		case "is.workflow.actions.repeat.count":
			decompRepeat(&action)
		case "is.workflow.actions.repeat.each":
			decompFor(&action)
		case "is.workflow.actions.choosefrommenu":
			decompMenu(&action)
		case "is.workflow.actions.getvalueforkey":
			var dictionaryKey = action.WFWorkflowActionParameters["WFDictionaryKey"]
			if dictionaryKey != nil && action.WFWorkflowActionParameters[UUID] != nil &&
				!slices.Contains(constUUIDs, paramString(&action, UUID)) {
				decompDictionaryGetValue(&action)
				continue
			}
			fallthrough
		default:
			decompAction(&action)
		}
		actionIndex++
	}
}

var varUUIDs []string
var constUUIDs []string

// checkConstantLiteral writes the current variable value of action as a constant, holds it for the Set Variable
// that references it, or writes action as a call on its own line when nothing references its output.
func checkConstantLiteral(action *ShortcutAction) {
	if action.WFWorkflowActionParameters[UUID] == nil {
		writeUnreferencedCall(action)
		return
	}
	var actionUUID = paramString(action, UUID)
	if slices.Contains(varUUIDs, actionUUID) && !slices.Contains(constUUIDs, actionUUID) {
		currentVariableUUID = actionUUID
		return
	}
	var outputName, referenced = uuids[actionUUID]
	if !referenced {
		writeUnreferencedCall(action)
		return
	}

	newCodeLine(fmt.Sprintf("const %s = ", outputName))
	code.WriteString(currentVariableValue)
	code.WriteRune('\n')
	currentVariableValue = ""
}

// writeUnreferencedCall writes action as a call because a literal on its own line is not valid Cherri.
// A text action is matched by more than one action definition, so it is written with the text alias.
func writeUnreferencedCall(action *ShortcutAction) {
	var call string
	if action.WFWorkflowActionIdentifier == "is.workflow.actions.gettext" {
		call = fmt.Sprintf("text(%s)", currentVariableValue)
	} else {
		call = makeActionCallCode(action)
	}

	code.WriteString(tabbedLine(call))
	code.WriteRune('\n')
	currentVariableValue = ""
}

func makeConstantLiteral(action *ShortcutAction) {
	if _, found := action.WFWorkflowActionParameters[UUID]; !found {
		return
	}
	var actionUUID = paramString(action, UUID)
	if outputName, found := uuids[actionUUID]; found {
		newCodeLine(fmt.Sprintf("const %s = ", outputName))
		code.WriteString(currentVariableValue)
		code.WriteRune('\n')
		currentVariableValue = ""
	}
}

func decompComment(action *ShortcutAction) {
	var commentText = paramString(action, "WFCommentActionText")
	if args.Using("comments") {
		if strings.Contains(commentText, "\n") {
			newCodeLine(fmt.Sprintf("comment('\n%s\n')\n\n", commentText))
		} else {
			newCodeLine(fmt.Sprintf("comment('%s')\n", commentText))
		}
	}

	if strings.Contains(commentText, "\n") {
		newCodeLine("/*\n")
		for _, line := range strings.Split(commentText, "\n") {
			newCodeLine(fmt.Sprintln(line))
		}
		newCodeLine("*/\n")
	} else {
		newCodeLine(fmt.Sprintf("// %s\n\n", commentText))
	}
}

var decompilingText = false

func decompTextValue(action *ShortcutAction) {
	decompilingText = true
	currentVariableValue = decompValue(action.WFWorkflowActionParameters["WFTextActionText"])
	if currentVariableValue == "" {
		currentVariableValue = "\"\""
	} else {
		currentVariableValue = fmt.Sprintf("\"%s\"", escapeString(strings.Trim(currentVariableValue, "\"")))
	}
	checkConstantLiteral(action)
	decompilingText = false
}

func decompNumberValue(action *ShortcutAction) (nonLiteral bool) {
	var value = requiredParam(action, "WFNumberActionNumber")
	if mapValue, isMap := value.(map[string]interface{}); isMap {
		var Value, hasValue = mapValue["Value"].(map[string]interface{})
		if !hasValue {
			decompError(fmt.Sprintf("Expected the parameter 'WFNumberActionNumber' to have a Value dictionary, got %T", mapValue["Value"]), action)
		}
		value = decompValue(value)

		match, _ := regexp.MatchString("[0-9.]", value.(string))
		if !match {
			nonLiteral = true
			return
		}

		if Value["Type"] == "ActionOutput" {
			currentVariableValue = value.(string)
			checkConstantLiteral(action)
		}
		return
	}

	var number any
	var convErr error
	if stringValue, isString := value.(string); isString {
		if strings.Contains(stringValue, ".") {
			number, convErr = strconv.ParseFloat(stringValue, 64)
		} else if strings.Contains(stringValue, ",") {
			number, convErr = strconv.ParseFloat(strings.Replace(stringValue, ",", ".", 1), 64)
		} else if stringValue != "" {
			number, convErr = strconv.Atoi(stringValue)
		}
		handle(convErr)
	} else {
		switch v := value.(type) {
		case uint64:
			number = int(v)
		case float64:
			number = v
		}
	}

	currentVariableValue = decompValue(number)
	checkConstantLiteral(action)

	return
}

func decompBasicExpression(action *ShortcutAction) {
	var input = action.WFWorkflowActionParameters["WFInput"]
	var operand = action.WFWorkflowActionParameters["WFMathOperand"]
	var operation = "+"
	if action.WFWorkflowActionParameters["WFMathOperation"] != nil {
		switch action.WFWorkflowActionParameters["WFMathOperation"] {
		case "×":
			operation = "*"
		case "÷":
			operation = "/"
		default:
			operation = paramString(action, "WFMathOperation")
		}
	}
	var expression = strings.Trim(decompValue(input), "\"") + " " + strings.Trim(operation, "\"") + " " + strings.Trim(decompValue(operand), "\"")
	var varRegex = regexp.MustCompile(`{(.*?)}`)
	currentVariableValue = varRegex.ReplaceAllString(expression, "$1")

	checkConstantLiteral(action)
}

func decompExpression(action *ShortcutAction) {
	var expression string
	switch input := requiredParam(action, "Input").(type) {
	case string:
		expression = input
	case map[string]interface{}:
		expression = strings.Trim(decompValue(input["Value"]), "\"")
	default:
		decompError(fmt.Sprintf("Expected the parameter 'Input' to be a string or a dictionary, got %T", input), action)
	}
	var varRegex = regexp.MustCompile(`{(.*?)}`)
	currentVariableValue = varRegex.ReplaceAllString(expression, "$1")

	checkConstantLiteral(action)
}

func decompVariable(action *ShortcutAction) {
	var variableName = paramString(action, "WFVariableName")
	sanitizeIdentifier(&variableName)
	newCodeLine(fmt.Sprintf("@%s", variableName))

	var inputUUID string
	if action.WFWorkflowActionParameters["WFInput"] != nil {
		var wfInput WFInput
		mapToStruct(action.WFWorkflowActionParameters["WFInput"], &wfInput)
		inputUUID = wfInput.Value.OutputUUID
	}

	if currentVariableValue == "" || inputUUID == "" || inputUUID != currentVariableUUID {
		var decompInput = decompValue(action.WFWorkflowActionParameters["WFInput"])
		if decompInput != "" {
			code.WriteString(fmt.Sprintf(" = %s", decompInput))
		}
		code.WriteRune('\n')
		return
	}

	code.WriteRune(' ')
	if action.WFWorkflowActionIdentifier == AppendVariableIdentifier {
		code.WriteString("+= ")
	} else {
		code.WriteString("= ")
	}
	code.WriteString(currentVariableValue)
	code.WriteRune('\n')

	currentVariableValue = ""
	currentVariableUUID = ""
}

func decompDictionaryGetValue(action *ShortcutAction) {
	var dictionaryValueRef strings.Builder
	dictionaryValueRef.WriteString(decompValue(action.WFWorkflowActionParameters["WFInput"]))

	if action.WFWorkflowActionParameters["WFDictionaryKey"] != nil {
		dictionaryValueRef.WriteRune('[')
		if keyValue, isString := action.WFWorkflowActionParameters["WFDictionaryKey"].(string); isString {
			dictionaryValueRef.WriteRune('\'')
			dictionaryValueRef.WriteString(keyValue)
			dictionaryValueRef.WriteRune('\'')
		} else {
			dictionaryValueRef.WriteString(decompValue(action.WFWorkflowActionParameters["WFDictionaryKey"]))
		}
		dictionaryValueRef.WriteRune(']')
		currentVariableValue = dictionaryValueRef.String()
		checkConstantLiteral(action)
	}
}

func checkControlFlowOutput(action *ShortcutAction) bool {
	switch action.WFWorkflowActionIdentifier {
	case "is.workflow.actions.conditional":
		fallthrough
	case "is.workflow.actions.choosefrommenu":
		fallthrough
	case "is.workflow.actions.repeat.count":
		fallthrough
	case "is.workflow.actions.repeat.each":
		if controlFlowMode(action) == startStatement {
			var outputName, captured = controlFlowOutputs[paramString(action, "GroupingIdentifier")]
			if captured {
				newCodeLine(fmt.Sprintf("const %s = ", outputName))
				return true
			}
		}
	}
	return false
}

func beginStatement(t tokenType, action *ShortcutAction) {
	if tabLevel == 0 {
		newCodeLine("\n")
	}

	var statementKeyword = fmt.Sprintf("%s ", strings.TrimSpace(string(t)))
	if checkControlFlowOutput(action) {
		code.WriteString(statementKeyword)
	} else {
		newCodeLine(statementKeyword)
	}
}

func decompMenu(action *ShortcutAction) {
	if len(menus) == 0 {
		menus = make(map[string][]varValue)
	}
	var groupingUUID = paramString(action, "GroupingIdentifier")
	switch controlFlowMode(action) {
	case startStatement:
		beginStatement(Menu, action)

		menus[groupingUUID] = []varValue{}
		for _, item := range paramSlice(action, "WFMenuItems") {
			menus[groupingUUID] = append(menus[groupingUUID], varValue{value: item})
		}
		code.WriteString(decompValue(action.WFWorkflowActionParameters["WFMenuPrompt"]))
		code.WriteString(" {\n")
		tabLevel++
	case statementPart:
		tabLevel--
		newCodeLine("item ")

		// Cases with an attributed title do not read their option but still consume it.
		var menuItem varValue
		if len(menus[groupingUUID]) > 0 {
			menuItem = menus[groupingUUID][0]
			menus[groupingUUID] = menus[groupingUUID][1:]
		}

		if attributedTitle, found := action.WFWorkflowActionParameters["WFMenuItemAttributedTitle"]; found {
			code.WriteString(decompValue(attributedTitle))
		} else if menuItem.value != nil {
			code.WriteString(decompValue(menuItem.value))
		}
		code.WriteString(":\n")
		tabLevel++
	case endStatement:
		tabLevel--
		newCodeLine("}\n")
	}
}

func decompRepeat(action *ShortcutAction) {
	switch controlFlowMode(action) {
	case startStatement:
		beginStatement(Repeat, action)

		// The count is evaluated outside the loop, so it is decompiled before the loop opens.
		var count = decompValue(action.WFWorkflowActionParameters["WFRepeatCount"])
		code.WriteString(fmt.Sprintf("%s for %s {\n", loopDeclaration(action), count))
		tabLevel++
	case endStatement:
		popLoop(action)
		tabLevel--
		newCodeLine("}\n")
	}
}

func decompFor(action *ShortcutAction) {
	switch controlFlowMode(action) {
	case startStatement:
		beginStatement(RepeatWithEach, action)

		// The input is evaluated outside the loop, so it is decompiled before the loop opens.
		var input = decompValue(action.WFWorkflowActionParameters["WFInput"])
		code.WriteString(fmt.Sprintf("%s in %s {\n", loopDeclaration(action), input))
		tabLevel++
	case endStatement:
		popLoop(action)
		tabLevel--
		newCodeLine("}\n")
	}
}

func decompConditional(action *ShortcutAction) {
	switch controlFlowMode(action) {
	case startStatement:
		beginStatement(If, action)

		if action.WFWorkflowActionParameters["WFConditions"] != nil {
			var conditions = paramMap(action, "WFConditions")
			var conditionValue, hasValue = conditions["Value"].(map[string]interface{})
			if !hasValue {
				decompError("Expected WFConditions to have a Value dictionary", action)
			}
			var intFilterPrefix, convErr = strconv.Atoi(fmt.Sprintf("%d", conditionValue["WFActionParameterFilterPrefix"]))
			handle(convErr)
			var paramFilterPrefix = convertFilterPrefix(intFilterPrefix)
			var templates = conditionTemplates(conditionValue, action)
			var conditionsLen = len(templates)
			for i, condition := range templates {
				decompCondition(condition, action)

				if i != conditionsLen-1 {
					code.WriteString(fmt.Sprintf(" %s ", paramFilterPrefix))
				}
			}
		} else {
			decompCondition(action.WFWorkflowActionParameters, action)
		}

		code.WriteString(" {\n")
		tabLevel++
	case statementPart:
		tabLevel--
		newCodeLine("} else {\n")
		tabLevel++
	case endStatement:
		tabLevel--
		newCodeLine("}\n")
	}
}

// matchConditionOperator reverses makeConditionParam, which swaps the greater and less than codes
// for 1000 and 1001 when the condition compares a date to a duration.
func matchConditionOperator(conditionCode int, condition map[string]interface{}, action *ShortcutAction) tokenType {
	switch conditionCode {
	case 1000, 1001:
		if condition["WFDuration"] == nil {
			decompError(fmt.Sprintf("Expected a WFDuration operand for conditional %d", conditionCode), action)
		}
		if conditionCode == 1000 {
			return GreaterThan
		}
		return LessThan
	}
	for operator, operatorCode := range conditions {
		if operatorCode == conditionCode {
			return operator
		}
	}

	decompError(fmt.Sprintf("Invalid conditional %d", conditionCode), action)
	return ""
}

func decompCondition(condition map[string]interface{}, action *ShortcutAction) {
	// The plist omits WFCondition when it is 0 (LessThan).
	var conditionCode uint64
	if condition["WFCondition"] != nil {
		var wfCondition, isInt = condition["WFCondition"].(uint64)
		if !isInt {
			decompError(fmt.Sprintf("Expected WFCondition to be an integer, got %T", condition["WFCondition"]), action)
		}
		conditionCode = wfCondition
	}

	var conditionalOperator = matchConditionOperator(int(conditionCode), condition, action)
	if conditionalOperator == Empty {
		code.WriteRune('!')
	}

	code.WriteString(decompValue(condition["WFInput"]))

	switch conditionalOperator {
	case Any, Empty:
		return
	case IsToday:
		code.WriteString(fmt.Sprintf(" %s", IsToday))
		return
	}

	code.WriteString(fmt.Sprintf(" %s ", conditionalOperator))
	code.WriteString(decompConditionOperand(condition, []string{"WFDuration", "WFDate", "WFNumberValue", "WFConditionalActionString"}, action))

	if conditionalOperator == Between {
		code.WriteRune(' ')
		code.WriteString(decompConditionOperand(condition, []string{"WFAnotherDate", "WFAnotherNumber"}, action))
	}
}

// decompConditionOperand writes the first of keys present in condition, as routeConditionalValue
// and conditionalParameter write each operand to exactly one of them.
func decompConditionOperand(condition map[string]interface{}, keys []string, action *ShortcutAction) string {
	for _, key := range keys {
		var operand = condition[key]
		if operand == nil {
			continue
		}

		switch key {
		case "WFDuration":
			return decompConditionDuration(operand, action)
		case "WFNumberValue", "WFAnotherNumber":
			return decompConditionNumber(operand)
		default:
			return decompValue(operand)
		}
	}

	decompError(fmt.Sprintf("Expected a condition operand in one of %s", strings.Join(keys, ", ")), action)
	return ""
}

// decompConditionNumber writes a number stored as text, as a duration's Magnitude always is, without
// converting it so floats are not truncated.
func decompConditionNumber(operand any) string {
	var numberText, isText = operand.(string)
	if !isText {
		return decompValue(operand)
	}
	if _, parseErr := strconv.ParseFloat(numberText, 64); parseErr != nil {
		return decompValue(numberText)
	}

	return numberText
}

func decompConditionDuration(operand any, action *ShortcutAction) string {
	var duration, isMap = operand.(map[string]interface{})
	if !isMap {
		decompError(fmt.Sprintf("Expected WFDuration to be a WFQuantityFieldValue, got %T", operand), action)
	}
	var quantity, hasValue = duration["Value"].(map[string]interface{})
	if !hasValue || quantity["Magnitude"] == nil || quantity["Unit"] == nil {
		decompError("Expected WFDuration to have a Value with a Magnitude and a Unit", action)
	}

	return fmt.Sprintf("qty(%s, %s)", decompConditionNumber(quantity["Magnitude"]), decompValue(quantity["Unit"]))
}

type DictionaryActionParameters struct {
	WFItems WFItems
}

type WFItems struct {
	Value Value
}

var decompilingDictionary = false

func decompDictionary(action *ShortcutAction) {
	var params DictionaryActionParameters
	mapToStruct(action.WFWorkflowActionParameters, &params)

	currentVariableValue = decompDictionaryValue(params.WFItems.Value.WFDictionaryFieldValueItems)

	checkConstantLiteral(action)
}

func decompDictionaryValue(items []WFDictionaryFieldValueItem) string {
	decompilingDictionary = true

	var dictionary = decompDictionaryItems(items)
	var jsonBytes, jsonErr = json.MarshalIndent(dictionary, strings.Repeat("\t", tabLevel), "\t")
	handle(jsonErr)

	decompilingDictionary = false

	return bareReferenceMarkerRegex.ReplaceAllString(string(jsonBytes), "$1")
}

// bareReference is a variable reference that is written without quotes. JSON cannot express that, so it is
// marshaled as a string with a NUL prefix that is unquoted afterwards, as no real dictionary text starts with one.
type bareReference string

const bareReferenceMarker = "\x00"

var bareReferenceMarkerRegex = regexp.MustCompile(`"\\u0000([^"]*)"`)

func (reference bareReference) MarshalJSON() ([]byte, error) {
	return json.Marshal(bareReferenceMarker + string(reference))
}

func isReferenceValue(value any) bool {
	if mapValue, isMap := value.(map[string]interface{}); isMap {
		if mapValue["Value"] != nil && mapValue["WFSerializationType"] == "WFTextTokenAttachment" {
			return true
		}
	}

	return false
}

func decompDictionaryItems(items []WFDictionaryFieldValueItem) (dictionary map[string]interface{}) {
	dictionary = make(map[string]interface{})
	for _, item := range items {
		var itemKey = strings.Trim(decompValue(item.WFKey), "\"")
		dictionary[itemKey] = decompDictionaryItem(item)
	}
	return
}

var numericRegex = regexp.MustCompile("^[0-9]+$")

func decompDictionaryItem(item WFDictionaryFieldValueItem) any {
	var itemValueMap, isMap = item.WFValue.(map[string]interface{})
	if !isMap {
		decompError(fmt.Sprintf("Expected a dictionary item value to be a dictionary, got %T", item.WFValue), currentDecompAction())
	}
	var itemStringValue = decompValue(item.WFValue)
	var itemValue any

	if dictDataType(item.WFItemType) == itemTypeFile {
		return bareReference(decompValue(itemValueMap["Value"]))
	}

	if isReferenceValue(itemValueMap["Value"]) {
		return fmt.Sprintf("{%s}", itemStringValue)
	}

	switch dictDataType(item.WFItemType) {
	case itemTypeNumber:
		itemValue = decompValue(item.WFValue)
		if !numericRegex.MatchString(itemValue.(string)) {
			break
		}

		var convErr error
		itemValue, convErr = strconv.ParseInt(itemValue.(string), 10, 64)
		handle(convErr)
	case itemTypeBool:
		var boolValue, isBool = itemValueMap["Value"].(bool)
		if !isBool {
			decompError(fmt.Sprintf("Expected a boolean dictionary item value, got %T", itemValueMap["Value"]), currentDecompAction())
		}
		itemValue = boolValue
	case itemTypeText:
		itemValue = strings.Trim(itemStringValue, "\"")
	case itemTypeArray:
		var arrayItems, isSlice = itemValueMap["Value"].([]interface{})
		if !isSlice {
			decompError(fmt.Sprintf("Expected an array dictionary item value to be an array, got %T", itemValueMap["Value"]), currentDecompAction())
		}
		itemValue = decompArray(arrayItems)
	case itemTypeDict:
		var dictionaryItems []WFDictionaryFieldValueItem
		var value, hasValue = itemValueMap["Value"].(map[string]interface{})
		var nestedValue, hasNestedValue = value["Value"].(map[string]interface{})
		if !hasValue || !hasNestedValue {
			decompError("Expected a dictionary item value to be a dictionary holding a Value dictionary", currentDecompAction())
		}
		mapToStruct(nestedValue["WFDictionaryFieldValueItems"], &dictionaryItems)
		itemValue = decompDictionaryItems(dictionaryItems)
	default:
		itemValue = itemStringValue
	}

	return itemValue
}

func decompArray(items []interface{}) (array []interface{}) {
	for _, item := range items {
		if itemMap, isMap := item.(map[string]interface{}); isMap {
			var fieldValueItem WFDictionaryFieldValueItem
			mapToStruct(itemMap, &fieldValueItem)
			array = append(array, decompDictionaryItem(fieldValueItem))
			continue
		}

		array = append(array, strings.Trim(decompValue(item), "\""))
	}
	return
}

func decompValue(value any) string {
	if value == nil {
		return "nil"
	}

	switch typed := value.(type) {
	case map[string]interface{}:
		return decompValueObject(typed)
	case string:
		return fmt.Sprintf("\"%s\"", escapeString(typed))
	default:
		return fmt.Sprintf("%v", value)
	}
}

func escapeString(value string) string {
	var escapes = map[string]string{
		"\n": "\\n",
		"\t": "\\t",
		"\r": "\\r",
		"\"": `\"`,
	}
	for chr, e := range escapes {
		value = strings.ReplaceAll(value, chr, e)
	}

	return value
}

func decompValueObject(value map[string]interface{}) string {
	if v, found := value["Value"]; found {
		if nested, isMap := v.(map[string]interface{}); isMap {
			value = nested
		}
	}

	if value["WFDictionaryFieldValueItems"] != nil {
		var items []WFDictionaryFieldValueItem
		mapToStruct(value["WFDictionaryFieldValueItems"], &items)

		return decompDictionaryValue(items)
	}

	switch value["Type"] {
	case "Variable":
		if _, found := value["VariableName"]; found {
			var variableName, isString = value["VariableName"].(string)
			if !isString {
				decompError(fmt.Sprintf("Expected VariableName to be a string, got %T", value["VariableName"]), currentDecompAction())
			}
			if loopReference, isLoopVariable := loopVariableReference(variableName); isLoopVariable {
				return loopReference
			}
			sanitizeIdentifier(&variableName)

			return fmt.Sprintf("@%s", variableName)
		}

		var variableValue, hasVariable = value["Variable"].(map[string]interface{})
		if !hasVariable {
			decompError(fmt.Sprintf("Expected a Variable value to have a VariableName or a Variable dictionary, got %T", value["Variable"]), currentDecompAction())
		}
		return decompValue(variableValue["Value"])
	case "ActionOutput":
		if _, found := value["OutputUUID"]; found {
			var outputUUID, isString = value["OutputUUID"].(string)
			if !isString {
				decompError(fmt.Sprintf("Expected OutputUUID to be a string, got %T", value["OutputUUID"]), currentDecompAction())
			}
			var outputName = uuids[outputUUID]
			sanitizeIdentifier(&outputName)

			if value["Aggrandizements"] == nil {
				return outputName
			}
		}
		break
	case "Ask":
		if value["Prompt"] == nil {
			return Ask
		}

		return fmt.Sprintf("Ask: \"%s\"", value["Prompt"])
	case globals[ShortcutInput].variableType:
		return ShortcutInput
	}

	return decompObjectValue(value)
}

func decompObjectValue(valueObj any) string {
	var value, isMap = valueObj.(map[string]interface{})
	if !isMap {
		return fmt.Sprintf("%v", valueObj)
	}

	var attachmentString string
	if value["Value"] != nil {
		var nested, isNestedMap = value["Value"].(map[string]interface{})
		if !isNestedMap {
			return fmt.Sprintf("%v", valueObj)
		}
		value = nested
	}

	if _, found := value["string"]; found {
		var valueString, isString = value["string"].(string)
		if !isString {
			decompError(fmt.Sprintf("Expected string of a text value to be a string, got %T", value["string"]), currentDecompAction())
		}
		attachmentString = valueString
	}

	if value["Aggrandizements"] != nil {
		attachmentString = ObjectReplaceCharStr
		decompAttachmentString(&attachmentString, map[string]any{
			"{0, 1}": value,
		})
		attachmentString = strings.Trim(attachmentString, "\"{}")
	}

	if attachments, found := value["attachmentsByRange"]; found {
		var attachmentsMap, isMap = attachments.(map[string]interface{})
		if !isMap {
			decompError(fmt.Sprintf("Expected attachmentsByRange to be a dictionary, got %T", attachments), currentDecompAction())
		}
		decompAttachmentString(&attachmentString, attachmentsMap)
	}

	return attachmentString
}

// decompAttachmentString substitutes each attachment into the string. Attachment ranges are UTF-16
// code unit offsets, the inverse of mapInlineVariables, so the string is indexed as UTF-16 rather
// than as code points, which would misplace attachments after an emoji.
func decompAttachmentString(attachmentString *string, attachments map[string]interface{}) {
	var originalString = *attachmentString
	var units = utf16.Encode([]rune(originalString))

	var positions []int
	var attachmentsByPosition = make(map[int]any, len(attachments))
	for attachmentRange, a := range attachments {
		var attachmentStart, _, _ = strings.Cut(attachmentRange, ",")
		var attachmentPosition = strings.TrimPrefix(attachmentStart, "{")
		var position, convErr = strconv.Atoi(attachmentPosition)
		handle(convErr)

		positions = append(positions, position)
		attachmentsByPosition[position] = a
	}
	sort.Ints(positions)

	var result strings.Builder
	var consumed int
	var referencesControlFlowOutput bool
	for _, position := range positions {
		if position < consumed || position >= len(units) {
			decompError(fmt.Sprintf("Attachment position %d is outside of the string", position), currentDecompAction())
		}

		var attachment Value
		mapToStruct(attachmentsByPosition[position], &attachment)

		if slices.Contains(controlFlowOutputUUIDs, attachment.OutputUUID) {
			referencesControlFlowOutput = true
		}

		var variableName = attachment.VariableName
		if _, uuid := uuids[attachment.OutputUUID]; uuid {
			variableName = uuids[attachment.OutputUUID]
		}
		sanitizeIdentifier(&variableName)
		if variableName == "" && attachment.Type == globals[ShortcutInput].variableType {
			variableName = ShortcutInput
		}
		if attachment.Type != "Variable" {
			for name, global := range globals {
				if global.variableType == attachment.Type {
					variableName = name
				}
			}
		} else if loopReference, isLoopVariable := loopVariableReference(attachment.VariableName); isLoopVariable {
			variableName = loopReference
		} else {
			variableName = fmt.Sprintf("@%s", variableName)
		}

		if len(attachment.Aggrandizements) != 0 {
			decompAggrandizements(&variableName, attachment.Aggrandizements)
		}

		result.WriteString(string(utf16.Decode(units[consumed:position])))
		result.WriteString(fmt.Sprintf("{%s}", variableName))
		consumed = position + 1
	}
	result.WriteString(string(utf16.Decode(units[consumed:])))

	*attachmentString = escapeString(result.String())

	if !decompilingDictionary && !decompilingText {
		// The output of a control flow block has no type, so the compiler only accepts it inside text.
		if originalString == ObjectReplaceCharStr && !referencesControlFlowOutput {
			*attachmentString = strings.Trim(*attachmentString, "{}")
		} else if !strings.Contains("\"", *attachmentString) {
			*attachmentString = fmt.Sprintf("\"%s\"", *attachmentString)
		}
	}
}

func decompAggrandizements(reference *string, aggrs []Aggrandizement) {
	var index string
	var coerce string
	var revContentItems = reversedContentItems()
	for _, aggr := range aggrs {
		switch aggr.Type {
		case "WFCoercionVariableAggrandizement":
			if _, found := revContentItems[aggr.CoercionItemClass]; found {
				coerce = revContentItems[coerce]
			}
		case "WFDictionaryValueVariableAggrandizement":
			index = aggr.DictionaryKey
		case "WFPropertyVariableAggrandizement":
			index = aggr.PropertyName
		}
	}

	if index != "" {
		*reference = fmt.Sprintf("%s['%s']", *reference, index)
	}
	if index != "" && coerce != "" {
		*reference = fmt.Sprintf("%s.%s", *reference, coerce)
	}
}

// popLine adds line to the top of the generated Cherri code.
func popLine(line string) {
	var saveCode = code.String()
	code.Reset()
	code.WriteString(fmt.Sprintf("%s\n%s", line, saveCode))
}

var macDefinition bool
var setMacDefinition bool

func decompAction(action *ShortcutAction) {
	if skipDecompAction(action) {
		return
	}

	currentVariableValue = makeActionCallCode(action)

	var isConstant, isVariableValue = checkOutputType(action)
	if isConstant && isVariableValue {
		makeConstantLiteral(action)
	} else if isConstant {
		checkConstantLiteral(action)
	} else if isVariableValue {
		currentVariableUUID = paramString(action, UUID)
	} else {
		code.WriteString(tabbedLine(currentVariableValue))
		code.WriteRune('\n')
		currentVariableValue = ""
	}
}

func makeActionCallCode(action *ShortcutAction) string {
	var actionCallCode strings.Builder
	var matchedIdentifier, matchedAction = matchAction(action)
	if matchedIdentifier == "" {
		if !strings.Contains(action.WFWorkflowActionIdentifier, ".") {
			decompError(fmt.Sprintf("Expected the action identifier '%s' to be a dot separated identifier such as 'is.workflow.actions.name'", action.WFWorkflowActionIdentifier), action)
		}

		checkMissingStandardInclude(&action.WFWorkflowActionIdentifier, false)

		matchedIdentifier, matchedAction = matchAction(action)
		if matchedIdentifier == "" {
			actionCallCode.WriteString(makeRawAction(action))
			return actionCallCode.String()
		}
	}
	includeStandardAction(matchedIdentifier)

	if (matchedAction.macOnly || matchedAction.nonMacOnly) && !setMacDefinition {
		macDefinition = matchedAction.macOnly && !matchedAction.nonMacOnly
		popLine(fmt.Sprintf("#define mac %v", macDefinition))
		setMacDefinition = true
	}

	actionCallCode.WriteString(fmt.Sprintf("%s(", matchedIdentifier))

	if matchedAction.makeParams != nil || matchedAction.decomp != nil {
		decompActionCustom(&actionCallCode, &matchedAction, action)
	} else {
		var matchedParamsSize = len(matchedAction.parameters)
		if matchedParamsSize > 0 {
			decompActionArguments(&actionCallCode, &matchedAction, action)
		}
	}

	actionCallCode.WriteString(")")

	return actionCallCode.String()
}

// checkOutputType determines if action output is a constant or a variable.
// If it is a constant we will write a constant statement on a new line to prepend the action.
func checkOutputType(action *ShortcutAction) (isConstant bool, isVariableValue bool) {
	if action.WFWorkflowActionParameters[UUID] == nil {
		return
	}
	var uuid = paramString(action, UUID)
	if _, found := uuids[uuid]; !found {
		return
	}

	isConstant = slices.Contains(constUUIDs, uuid) || !slices.Contains(varUUIDs, uuid)
	isVariableValue = slices.Contains(varUUIDs, uuid)

	return
}

// skipDecompAction skips actions we don't support or when necessary.
func skipDecompAction(action *ShortcutAction) bool {
	var identifier = actionIdentifierEnd(action.WFWorkflowActionIdentifier)
	if identifier == "getvariable" {
		var varName = decompValue(action.WFWorkflowActionParameters["WFVariable"])

		insertCodeComment(fmt.Sprintf("TODO: Get Variable not supported: Assign variable here to '%s'.", varName))
		decompWarning(fmt.Sprintf("Get variable '%s' is not supported. Set a variable to that value instead if something was depending on it's output.", varName))

		return true
	}

	if identifier == "nothing" && actionIndex+1 < len(shortcut.WFWorkflowActions) {
		var nextAction = peekActions(1)
		var controlflowActionIdentifiers = []string{"conditional", "repeat.each", "repeat.count", "choosefrommenu"}
		var nextActionIdentifier = actionIdentifierEnd(nextAction.WFWorkflowActionIdentifier)
		if slices.Contains(controlflowActionIdentifiers, nextActionIdentifier) {
			var controlFlowMode = nextAction.WFWorkflowActionParameters["WFControlFlowMode"]
			if controlFlowMode == endStatement || controlFlowMode == statementPart {
				return true
			}
		}
	}

	return false
}

func actionIdentifierEnd(identifier string) string {
	return strings.Replace(identifier, "is.workflow.actions.", "", 1)
}

func decompActionArguments(actionCallCode *strings.Builder, matchedAction *actionDefinition, action *ShortcutAction) {
	for i, param := range matchedAction.parameters {
		if param.key == "" {
			continue
		}

		var argValue string
		if value, found := action.WFWorkflowActionParameters[param.key]; found {
			argValue = decompValue(value)
		} else if !param.optional {
			argValue = makeDefaultValue(param)
		}

		switch param.validType {
		case Integer:
			if startsWith(Ask, argValue) {
				break
			}
			argValue = strings.Trim(argValue, "\"")
		}

		if argValue != "" {
			if i == 0 {
				actionCallCode.WriteString(argValue)
			} else {
				actionCallCode.WriteString(fmt.Sprintf(", %s", argValue))
			}
		}
	}
}

func makeDefaultValue(param parameterDefinition) string {
	if param.defaultValue != nil {
		if stringValue, isString := param.defaultValue.(string); isString {
			return fmt.Sprintf("\"%s\"", stringValue)
		}

		return fmt.Sprintf("%v", param.defaultValue)
	}

	switch param.validType {
	case Integer:
		return "0"
	case String:
		return "\"\""
	case Arr:
		return "[]"
	case Dict:
		return "{}"
	case Bool:
		return "false"
	}

	return "nil"
}

func makeRawAction(action *ShortcutAction) string {
	var paramsSize = len(action.WFWorkflowActionParameters)
	var onlyUUIDParam = false
	var onlyOutputNameParam = false
	if paramsSize > 1 {
		if _, found := action.WFWorkflowActionParameters[UUID]; found {
			onlyUUIDParam = found
		}
		if _, found := action.WFWorkflowActionParameters["CustomOutputName"]; found {
			onlyOutputNameParam = found
		}
	}
	if paramsSize == 0 || (onlyUUIDParam && onlyOutputNameParam) {
		return fmt.Sprintf("rawAction(\"%s\")", action.WFWorkflowActionIdentifier)
	}

	var rawParams = processRawParameters(action.WFWorkflowActionParameters)
	var jb, jsonErr = json.MarshalIndent(rawParams, strings.Repeat("\t", tabLevel), "\t")
	handle(jsonErr)

	var arguments = strings.Join([]string{fmt.Sprintf("\"%s\"", action.WFWorkflowActionIdentifier), string(jb)}, ", ")

	return fmt.Sprintf("rawAction(%s)", arguments)
}

// processRawParameters returns the parameters to write for a raw action. It builds a new map as
// the action's own parameters are still read after the raw action is written.
func processRawParameters(params map[string]any) map[string]any {
	var rawParams = make(map[string]any, len(params))
	for key, value := range params {
		if key == UUID || key == "CustomOutputName" {
			continue
		}

		if mapValue, isMap := value.(map[string]interface{}); isMap {
			decompilingDictionary = true
			rawParams[key] = decompValueObject(mapValue)
			decompilingDictionary = false
			continue
		}

		rawParams[key] = value
	}

	return rawParams
}

func matchAction(action *ShortcutAction) (name string, definition actionDefinition) {
	for call, def := range actions {
		var identifier = strings.ToLower(call)
		if def.identifier != "" {
			identifier = def.identifier
		}
		var appIdentifier = "is.workflow.actions"
		if def.appIdentifier != "" {
			appIdentifier = def.appIdentifier
		}
		var actionIdentifier = fmt.Sprintf("%s.%s", appIdentifier, identifier)
		if actionIdentifier != action.WFWorkflowActionIdentifier && def.overrideIdentifier != action.WFWorkflowActionIdentifier {
			continue
		}

		definition = *def
		name = call

		if splitActions, found := identifierMap[identifier]; found {
			matchSplitAction(&splitActions, action.WFWorkflowActionParameters, &name, &definition)
			if name != "run" && name != "runSelf" {
				return
			}
			// A shortcut run by name has no WFWorkflow, and a run of itself without input only has isSelf.
			if _, hasWorkflow := action.WFWorkflowActionParameters["WFWorkflow"]; !hasWorkflow {
				return
			}
			var workflow = paramMap(action, "WFWorkflow")
			if _, hasSelf := workflow["isSelf"]; !hasSelf {
				return
			}
			var isSelf, isBool = workflow["isSelf"].(bool)
			if !isBool {
				decompError(fmt.Sprintf("Expected isSelf of WFWorkflow to be a boolean, got %T", workflow["isSelf"]), action)
			}
			if isSelf || action.WFWorkflowActionParameters["WFWorkflowName"] == basename {
				name = "runSelf"
				definition = *actions["runSelf"]
			} else {
				name = "run"
			}
		}

		return
	}
	return
}

type actionMatch struct {
	params uint
	values uint
	action actionValue
}

func matchSplitAction(splitActions *[]actionValue, parameters map[string]any, identifier *string, definition *actionDefinition) {
	if args.Using("debug") {
		fmt.Println("## MATCHING SPLIT ACTIONS ##")
		fmt.Println("parameters", parameters)
	}

	var defaultAction, hasDefaultAction = getDefaultAction(splitActions)
	if hasDefaultAction {
		*identifier = defaultAction.identifier
		*definition = *defaultAction.definition
		if args.Using("debug") {
			fmt.Println("has default action", defaultAction.identifier)
		}
	}

	var matches = getSplitActionMatches(splitActions, parameters)

	if len(matches) == 0 {
		return
	}

	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].params > matches[j].params || matches[i].values > matches[j].values
	})

	if args.Using("debug") {
		for _, match := range matches[0:] {
			fmt.Printf("%s()\n", match.action.identifier)
			fmt.Println("params:", match.params, " values:", match.values)
			fmt.Println("---")
		}
		fmt.Print("\n")
	}

	if !competingMatches(matches) {
		if args.Using("debug") {
			fmt.Println("matches are not competing! using default")
		}
		return
	}

	var matchedAction = matches[0]
	*identifier = matchedAction.action.identifier
	*definition = *matchedAction.action.definition
}

func getSplitActionMatches(splitActions *[]actionValue, parameters map[string]any) (matches []actionMatch) {
	for _, splitAction := range *splitActions {
		var matchedParams, matchedValues = scoreActionMatch(splitAction, splitAction.definition.parameters, parameters)
		if matchedParams == 0 {
			continue
		}

		matches = append(matches, actionMatch{
			params: matchedParams,
			values: matchedValues,
			action: splitAction,
		})
	}
	return
}

// competingMatches determines if the matches for this identifier have more values than 1 matching this action.
func competingMatches(matches []actionMatch) bool {
	var hasCompetingValues bool
	for _, match := range matches {
		if match.values > 1 {
			hasCompetingValues = true
			break
		}
	}
	if hasCompetingValues {
		return true
	}
	for _, match := range matches {
		if match.params > 1 {
			return true
		}
	}
	return false
}

func scoreActionMatch(splitAction actionValue, splitActionParams []parameterDefinition, parameters map[string]any) (matchedParams uint, matchedValues uint) {
	var paramMatches, valueMatches = scoreActionParams(&splitActionParams, parameters)
	matchedParams += paramMatches
	matchedValues += valueMatches

	var splitActionAddParams []parameterDefinition
	for key, value := range splitAction.definition.appendParams {
		splitActionAddParams = append(splitActionAddParams, parameterDefinition{
			key:          key,
			defaultValue: value,
		})
	}
	if splitAction.definition.appendParamsFunc != nil {
		for key, value := range splitAction.definition.appendParamsFunc([]actionArgument{}) {
			splitActionAddParams = append(splitActionAddParams, parameterDefinition{
				key:          key,
				defaultValue: value,
			})
		}
	}

	var addParamMatches, addValueMatches = scoreActionAddParams(&splitActionAddParams, parameters)
	matchedParams += addParamMatches
	matchedValues += addValueMatches

	return
}

func scoreActionParams(splitActionParams *[]parameterDefinition, parameters map[string]any) (matchedParams uint, matchedValues uint) {
	for _, param := range *splitActionParams {
		if param.key == "" {
			continue
		}
		if value, found := parameters[param.key]; found {
			matchedParams++
			if len(param.enum) > 0 && slices.Contains(getEnum(param.enum), fmt.Sprintf("%s", value)) {
				matchedValues++
			}
			if param.defaultValue != nil {
				var defaultValue = fmt.Sprintf("%v", param.defaultValue)
				var rawValue = strings.Trim(decompValue(value), "\"")
				if defaultValue == rawValue && defaultValue != "" && rawValue != "" {
					matchedValues++
				}
			}
		}
	}
	return
}

func scoreActionAddParams(splitActionAddParams *[]parameterDefinition, parameters map[string]any) (matchedParams uint, matchedValues uint) {
	for _, param := range *splitActionAddParams {
		if param.key == "" {
			continue
		}
		if value, found := parameters[param.key]; found {
			matchedParams++

			// Values can hold maps and slices, which cannot be compared with ==.
			if reflect.DeepEqual(param.defaultValue, value) {
				matchedValues++
			}
		}
	}
	return
}

func decompActionCustom(actionCode *strings.Builder, matchedAction *actionDefinition, action *ShortcutAction) {
	var arguments = matchedAction.decomp(action)
	if len(arguments) == 0 {
		return
	}
	actionCode.WriteString(strings.Join(arguments, ", "))
}

func glueToChar(glue string) string {
	switch glue {
	case "New Lines":
		return "\n"
	case "Spaces":
		return " "
	case "Every Character":
		return ""
	default:
		return glue
	}
}

// getDefaultAction gets the default action from a slice of split actions.
func getDefaultAction(splitActions *[]actionValue) (action actionValue, found bool) {
	for _, splitAction := range *splitActions {
		if splitAction.definition.defaultAction {
			action = splitAction
			found = true
			return
		}
	}

	return
}

func printDecompDebug() {
	fmt.Println(ansi("##### DEBUG #####\n", red))

	fmt.Println("### VARIABLES ###")
	printVariables()
	fmt.Print("\n")

	fmt.Println("### VARIABLE UUIDs ###")
	fmt.Println(varUUIDs)
	fmt.Print("\n")

	fmt.Println("### CONSTANT UUIDs ###")
	fmt.Println(constUUIDs)
	fmt.Print("\n")

	fmt.Println("### UUIDS ###")
	for uuid, name := range uuids {
		fmt.Printf("%s | %s\n", uuid, name)
	}
	fmt.Print("\n")
}

func decompWarning(message string) {
	var linesLen = strings.Count(code.String(), "\n")
	fmt.Println(ansi("Warning:", orange, bold), fmt.Sprintf("%s (%s:%d:0)\n", message, filePath, linesLen+1))
}

func decompError(message string, action *ShortcutAction) {
	fmt.Println(ansi(fmt.Sprintf("Error: %s\n\n", message), red, bold))

	fmt.Println("Action identifier:", action.WFWorkflowActionIdentifier)
	var writtenLines = strings.Split(code.String(), "\n")
	var linesLen = len(writtenLines)
	fmt.Printf("\nStopped while writing line %d:\n", linesLen)
	if linesLen >= 2 {
		fmt.Printf("%s %s\n", ansi(fmt.Sprintf("%d |", linesLen-1), dim), ansi(writtenLines[linesLen-2], dim))
	}
	fmt.Printf("%s %s\n", ansi(fmt.Sprintf("%d |", linesLen), dim), ansi(writtenLines[linesLen-1], red))

	os.Exit(1)
}

// currentDecompAction describes the action being decompiled for errors raised where the action is not at hand.
func currentDecompAction() *ShortcutAction {
	return &ShortcutAction{WFWorkflowActionIdentifier: currentAction.identifier}
}

// requiredParam returns the value of key, which the action must have.
func requiredParam(action *ShortcutAction, key string) any {
	var value, found = action.WFWorkflowActionParameters[key]
	if !found {
		decompError(fmt.Sprintf("Expected the parameter '%s' to be present", key), action)
	}

	return value
}

func paramMap(action *ShortcutAction, key string) map[string]interface{} {
	var value = requiredParam(action, key)
	var mapValue, isMap = value.(map[string]interface{})
	if !isMap {
		decompError(fmt.Sprintf("Expected the parameter '%s' to be a dictionary, got %T", key, value), action)
	}

	return mapValue
}

func paramString(action *ShortcutAction, key string) string {
	var value = requiredParam(action, key)
	var stringValue, isString = value.(string)
	if !isString {
		decompError(fmt.Sprintf("Expected the parameter '%s' to be a string, got %T", key, value), action)
	}

	return stringValue
}

func paramSlice(action *ShortcutAction, key string) []interface{} {
	var value = requiredParam(action, key)
	var sliceValue, isSlice = value.([]interface{})
	if !isSlice {
		decompError(fmt.Sprintf("Expected the parameter '%s' to be an array, got %T", key, value), action)
	}

	return sliceValue
}

func controlFlowMode(action *ShortcutAction) uint64 {
	var value = requiredParam(action, "WFControlFlowMode")
	var mode, isInt = value.(uint64)
	if !isInt {
		decompError(fmt.Sprintf("Expected the parameter 'WFControlFlowMode' to be an integer, got %T", value), action)
	}

	return mode
}
