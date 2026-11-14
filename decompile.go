/*
 * Copyright (c) Cherri
 */

package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
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
	currentVariableOperator = ""
	decompiledImports = nil
	decompiledDefinitions = nil
	questionIdentifiers = make(map[questionKey]string)
	openMenuItems = make(map[string]bool)
	appendedVariables = make(map[string]bool)
	assignedVariables = make(map[string]bool)

	// The Shortcut keeps its name for the name definition and runSelf, only the file name has no spaces.
	var outputBasename = strings.ReplaceAll(basename, " ", "_")
	outputPath = getOutputPath(outputBasename + ".cherri")
	if args.Using("no-ansi") {
		filePath = outputBasename + ".cherri"
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
	decompileDefinitions()
	decompileQuestions()

	decompileActions()
	checkQuestionsWritten()
	writeHeader()

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
			if params["WFInput"] == nil {
				continue
			}
			var wfInput WFInput
			mapToStruct(params["WFInput"], &wfInput)
			if len(wfInput.Value.Aggrandizements) != 0 {
				mapValueReference(wfInput.Value)
			}
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
	// assigned is whether the loop starts with the Set Variable that the compiler writes for its variable.
	assigned bool
}

var openLoops []openLoop

// loopIdentifiers maps the GroupingIdentifier of each loop whose variable is referenced to the identifier the
// loop declares, so a loop whose variable is never referenced can still be written with `_`.
var loopIdentifiers map[string]string

var loopVariableRegex = regexp.MustCompile("^Repeat (Index|Item)( [0-9]+)?$")

// mapLoopVariables finds the loops whose variable is referenced before any code is written, because the loop
// identifier is written before the references in its body.
func mapLoopVariables() {
	for index, action := range shortcut.WFWorkflowActions {
		currentAction.identifier = action.WFWorkflowActionIdentifier
		var isLoop = action.WFWorkflowActionIdentifier == repeatIdentifier || action.WFWorkflowActionIdentifier == repeatEachIdentifier
		if isLoop && controlFlowMode(&action) == endStatement {
			popLoop(&action)
			continue
		}

		markLoopReferences(action.WFWorkflowActionParameters)

		if isLoop && controlFlowMode(&action) == startStatement {
			pushLoop(&action, index)
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

// pushLoop opens the loop at index. The compiler numbers the variable of a repeat by how many repeats are open,
// and the variable of a repeat with each by how many of those are open, so each kind is counted on its own.
func pushLoop(action *ShortcutAction, index int) openLoop {
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

	// The compiler sets the variable the loop declares to the variable Shortcuts names, so the loop is
	// declared with that variable, and writing the Set Variable too would have the compiler write it twice.
	var assignedIdentifier, assigned = loopVariableAssignment(index, variableName)
	if assigned {
		identifier = assignedIdentifier
	}

	var loop = openLoop{
		groupingUUID: paramString(action, "GroupingIdentifier"),
		variableName: variableName,
		identifier:   identifier,
		assigned:     assigned,
	}
	openLoops = append(openLoops, loop)

	return loop
}

// loopVariableAssignment finds the variable that the Set Variable right after the loop at index sets to
// variableName, which is the variable the loop declares.
func loopVariableAssignment(index int, variableName string) (identifier string, found bool) {
	if index+1 >= len(shortcut.WFWorkflowActions) {
		return
	}

	var next = shortcut.WFWorkflowActions[index+1]
	if next.WFWorkflowActionIdentifier != SetVariableIdentifier || next.WFWorkflowActionParameters["WFInput"] == nil {
		return
	}

	var wfInput WFInput
	mapToStruct(next.WFWorkflowActionParameters["WFInput"], &wfInput)
	var reference = wfInput.Value
	if reference.Type != "Variable" || reference.VariableName != variableName || len(reference.Aggrandizements) != 0 {
		return
	}

	identifier = paramString(&next, "WFVariableName")
	sanitizeIdentifier(&identifier)

	return identifier, true
}

// isLoopVariableAssignment is whether the Set Variable at the current action is the one that the loop
// declaration just written stands for.
func isLoopVariableAssignment() bool {
	if actionIndex == 0 || len(openLoops) == 0 {
		return false
	}

	var previous = peekActions(-1)
	var isLoop = previous.WFWorkflowActionIdentifier == repeatIdentifier || previous.WFWorkflowActionIdentifier == repeatEachIdentifier
	if !isLoop || controlFlowMode(&previous) != startStatement {
		return false
	}

	var innermost = openLoops[len(openLoops)-1]
	return innermost.assigned && innermost.groupingUUID == paramString(&previous, "GroupingIdentifier")
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
	var loop = pushLoop(action, actionIndex)
	if _, referenced := loopIdentifiers[loop.groupingUUID]; !referenced && !loop.assigned {
		return "_"
	}

	return loop.identifier
}

// Map out variables in the Shortcut and their UUIDs for later checks.
func mapVariables() {
	for index, action := range shortcut.WFWorkflowActions {
		currentAction.identifier = action.WFWorkflowActionIdentifier
		if action.WFWorkflowActionIdentifier == SetVariableIdentifier || action.WFWorkflowActionIdentifier == AppendVariableIdentifier {
			var varName = paramString(&action, "WFVariableName")
			sanitizeIdentifier(&varName)
			if _, found := variables[varName]; !found {
				variables[varName] = varValue{}
			}
			if action.WFWorkflowActionIdentifier == AppendVariableIdentifier {
				appendedVariables[varName] = true
			}

			if action.WFWorkflowActionParameters["WFInput"] != nil {
				var wfInput WFInput
				mapToStruct(action.WFWorkflowActionParameters["WFInput"], &wfInput)
				// An output read through aggrandizements stays a constant, since folding its action
				// into the Set Variable would drop them.
				var _, found = uuids[wfInput.Value.OutputUUID]
				if found || len(wfInput.Value.Aggrandizements) != 0 {
					continue
				}

				// An action is only folded into the Set Variable written right after it, as folding one that
				// has other actions between would move it after them.
				if index > 0 && shortcut.WFWorkflowActions[index-1].WFWorkflowActionParameters[UUID] == wfInput.Value.OutputUUID {
					varUUIDs = append(varUUIDs, wfInput.Value.OutputUUID)
					continue
				}
				mapValueReference(wfInput.Value)
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

// identifierMap maps the full identifier that an action writes to the actions that write it, sorted by
// name, so which one of several that share an identifier is picked never depends on the order of a map.
var identifierMap map[string][]actionValue

// mapSplitActions indexes the actions by identifier. It is rebuilt from every action defined so far, as actions
// are defined while the Shortcut is decompiled.
func mapSplitActions() {
	identifierMap = make(map[string][]actionValue)
	for name, action := range actions {
		var identifier = fullActionIdentifier(name, action)
		identifierMap[identifier] = append(identifierMap[identifier], actionValue{
			identifier: name,
			definition: action,
		})
	}
	for _, candidates := range identifierMap {
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].identifier < candidates[j].identifier
		})
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

// decompiledImports are the includes and imports the code needs, and decompiledDefinitions the definitions it
// is written with. Both are written above the code, the imports first, once the code is known.
var decompiledImports []string
var decompiledDefinitions []string

// writeHeader writes the imports and then the definitions above the code. The blank line written before a
// block is trimmed from the start of the code even without a header, as the compiler rejects a file that
// starts with a blank line.
func writeHeader() {
	var body = strings.TrimLeft(code.String(), "\n")
	code.Reset()

	var header = append(slices.Clone(decompiledImports), decompiledDefinitions...)
	if len(header) == 0 {
		code.WriteString(body)
		return
	}

	code.WriteString(strings.Join(header, "\n"))
	code.WriteString("\n\n")
	code.WriteString(body)
}

// defineName defines the name of the imported Shortcut if the basename contains a space.
func defineName() {
	if strings.Contains(basename, " ") {
		decompiledDefinitions = append(decompiledDefinitions, fmt.Sprintf("#define name %s", basename))
	}
}

func decompileIcon() {
	var icon = shortcut.WFWorkflowIcon

	if icon.WFWorkflowIconStartColor != iconColor {
		if name, found := firstNameWithValue(icon.WFWorkflowIconStartColor, colors, altColors); found {
			decompiledDefinitions = append(decompiledDefinitions, fmt.Sprintf("#define color %s", name))
		}
	}

	if icon.WFWorkflowIconGlyphNumber != iconGlyph {
		if name, found := firstNameWithValue(int(icon.WFWorkflowIconGlyphNumber), glyphs); found {
			decompiledDefinitions = append(decompiledDefinitions, fmt.Sprintf("#define glyph %s", name))
		}
	}
}

// defaultClientVersion is the client version that is written when none is defined.
var defaultClientVersion = clientVersion

// defaultWorkflowTypes is what Shortcuts is set to be used in when no workflow types are defined.
var defaultWorkflowTypes = []string{"Watch", "WFWorkflowTypeShowInSearch"}

// definitionError stops on a part of the Shortcut that holds a definition and is not what the compiler writes it as.
func definitionError(message string, key string) {
	decompError(message, &ShortcutAction{WFWorkflowActionIdentifier: key})
}

// decompileDefinitions writes the definitions that the compiler writes the keys of the Shortcut from, left out
// when they are what is written without the definition.
func decompileDefinitions() {
	defineInputs()
	defineList("outputs", "WFWorkflowOutputContentItemClasses", shortcut.WFWorkflowOutputContentItemClasses, contentItems)

	var workflowTypeNames = slices.Sorted(slices.Values(shortcut.WFWorkflowTypes))
	if !slices.Equal(workflowTypeNames, slices.Sorted(slices.Values(defaultWorkflowTypes))) {
		defineList("from", "WFWorkflowTypes", shortcut.WFWorkflowTypes, workflowTypes)
	}
	defineList("quickactions", "WFQuickActionSurfaces", shortcut.WFQuickActionSurfaces, quickActions)

	defineNoInput()
	defineVersion()
}

// defineInputs writes the content items that the Shortcut accepts, which are all of them when not defined.
func defineInputs() {
	var accepted []string
	for _, class := range shortcut.WFWorkflowInputContentItemClasses {
		if _, found := reversedContentItems()[class]; found {
			accepted = append(accepted, class)
		}
	}
	if len(accepted) == len(contentItems) {
		warnUnnamedValues("WFWorkflowInputContentItemClasses", shortcut.WFWorkflowInputContentItemClasses, contentItems)
		return
	}

	defineList("inputs", "WFWorkflowInputContentItemClasses", shortcut.WFWorkflowInputContentItemClasses, contentItems)
}

// defineList writes the definition of a list of values by the names the compiler reads them from.
func defineList(definition string, key string, values []string, names map[string]string) {
	var listed = warnUnnamedValues(key, values, names)
	if len(listed) == 0 {
		return
	}

	decompiledDefinitions = append(decompiledDefinitions, fmt.Sprintf("#define %s %s", definition, strings.Join(listed, ", ")))
}

// warnUnnamedValues returns the names of values, and warns about the values Cherri has no name for, which
// cannot be written.
func warnUnnamedValues(key string, values []string, names map[string]string) (listed []string) {
	var valueNames = make(map[string]string, len(names))
	for name, value := range names {
		valueNames[value] = name
	}

	for _, value := range values {
		var name, found = valueNames[value]
		if !found {
			decompWarning(fmt.Sprintf("'%s' of %s has no name in Cherri and is left out.", value, key))
			continue
		}
		listed = append(listed, name)
	}

	return
}

// defineNoInput writes what the Shortcut does without input, which is the inverse of collectNoInputDefinition.
func defineNoInput() {
	if len(shortcut.WFWorkflowNoInputBehavior) == 0 {
		return
	}

	var key = "WFWorkflowNoInputBehavior"
	var name, isText = shortcut.WFWorkflowNoInputBehavior["Name"].(string)
	if !isText {
		definitionError(fmt.Sprintf("Expected %s to have a Name that is text, got %T", key, shortcut.WFWorkflowNoInputBehavior["Name"]), key)
	}

	var parameters map[string]interface{}
	if rawParameters, hasParameters := shortcut.WFWorkflowNoInputBehavior["Parameters"]; hasParameters {
		var isMap bool
		parameters, isMap = rawParameters.(map[string]interface{})
		if !isMap {
			definitionError(fmt.Sprintf("Expected the Parameters of %s to be a dictionary, got %T", key, rawParameters), key)
		}
	}

	switch name {
	case "WFWorkflowNoInputBehaviorShowError":
		var message, hasMessage = parameters["Error"].(string)
		if !hasMessage {
			definitionError(fmt.Sprintf("Expected the Parameters of %s to have an Error that is text, got %T", key, parameters["Error"]), key)
		}
		decompiledDefinitions = append(decompiledDefinitions, fmt.Sprintf("#define noinput %s \"%s\"", StopWith, escapeString(message)))
	case "WFWorkflowNoInputBehaviorAskForInput":
		var class, hasClass = parameters["ItemClass"].(string)
		if !hasClass {
			definitionError(fmt.Sprintf("Expected the Parameters of %s to have an ItemClass that is text, got %T", key, parameters["ItemClass"]), key)
		}
		var listed = warnUnnamedValues(key, []string{class}, contentItems)
		if len(listed) == 0 {
			return
		}
		decompiledDefinitions = append(decompiledDefinitions, fmt.Sprintf("#define noinput %s %s", AskFor, listed[0]))
	case "WFWorkflowNoInputBehaviorGetClipboard":
		decompiledDefinitions = append(decompiledDefinitions, fmt.Sprintf("#define noinput %s", GetClipboard))
	default:
		definitionError(fmt.Sprintf("Expected the Name of %s to be one of WFWorkflowNoInputBehaviorShowError, WFWorkflowNoInputBehaviorAskForInput or WFWorkflowNoInputBehaviorGetClipboard, got '%s'", key, name), key)
	}
}

// defineVersion writes the version that the client version of the Shortcut is written for. Several versions
// have the same client version, of which the shortest name is written.
func defineVersion() {
	if shortcut.WFWorkflowClientVersion == "" || shortcut.WFWorkflowClientVersion == defaultClientVersion {
		return
	}

	var versionName string
	for name, versionClient := range versions {
		if versionClient != shortcut.WFWorkflowClientVersion {
			continue
		}
		if versionName == "" || len(name) < len(versionName) || (len(name) == len(versionName) && name < versionName) {
			versionName = name
		}
	}
	if versionName == "" {
		decompWarning(fmt.Sprintf("The client version '%s' has no version in Cherri and is left out.", shortcut.WFWorkflowClientVersion))
		return
	}

	decompiledDefinitions = append(decompiledDefinitions, fmt.Sprintf("#define version %s", versionName))
}

// questionKey is the argument an import question holds the value of.
type questionKey struct {
	actionIndex  int
	parameterKey string
}

// questionIdentifiers names the import question of an argument, and writtenQuestions are those whose argument
// has been written.
var questionIdentifiers map[questionKey]string
var writtenQuestions = make(map[questionKey]bool)

// decompileQuestions writes a question declaration for each import question and names it for the argument that
// takes it. The compiler writes an empty value for the argument, so only the question holds one.
func decompileQuestions() {
	writtenQuestions = make(map[questionKey]bool)

	var questions = slices.Clone(shortcut.WFWorkflowImportQuestions)
	sort.SliceStable(questions, func(i, j int) bool {
		if questions[i].ActionIndex != questions[j].ActionIndex {
			return questions[i].ActionIndex < questions[j].ActionIndex
		}
		return questions[i].ParameterKey < questions[j].ParameterKey
	})

	var number = 1
	for _, question := range questions {
		var defaultValue, isText = question.DefaultValue.(string)
		if question.Category != "Parameter" || !isText || question.ParameterKey == "" {
			decompWarning(fmt.Sprintf("The import question '%s' is not a question about an argument with a text default value, which Cherri cannot write, and is left out.", question.Text))
			continue
		}

		var identifier = freshIdentifier("question", &number)
		questionIdentifiers[questionKey{question.ActionIndex, question.ParameterKey}] = identifier
		decompiledDefinitions = append(decompiledDefinitions, fmt.Sprintf("#question %s \"%s\" \"%s\"", identifier, escapeString(question.Text), escapeString(defaultValue)))
	}
}

// freshIdentifier returns the prefix followed by the next number that is not the name of a variable or a
// constant, as the compiler rejects a question named like one.
func freshIdentifier(prefix string, number *int) string {
	for {
		var identifier = fmt.Sprintf("%s%d", prefix, *number)
		*number++

		var taken bool
		for _, name := range uuids {
			taken = taken || name == identifier
		}
		if _, isVariable := variables[identifier]; !isVariable && !taken {
			return identifier
		}
	}
}

// questionArgument returns the question that holds the value of the parameter of the current action.
func questionArgument(parameterKey string) (identifier string, isQuestion bool) {
	var key = questionKey{actionIndex, parameterKey}
	identifier, isQuestion = questionIdentifiers[key]
	if isQuestion {
		writtenQuestions[key] = true
	}

	return
}

// checkQuestionsWritten warns about the import questions that no argument was written for, which are left out.
func checkQuestionsWritten() {
	for key, identifier := range questionIdentifiers {
		if !writtenQuestions[key] {
			decompWarning(fmt.Sprintf("The import question '%s' is for the parameter '%s' of action %d, which is not written as an argument.", identifier, key.parameterKey, key.actionIndex))
		}
	}
}

// firstNameWithValue returns the alphabetically first name that has the value, so that when names share a value
// the one that is written never depends on the order of a map.
func firstNameWithValue(value int, namedValues ...map[string]int) (firstName string, found bool) {
	for _, names := range namedValues {
		for name, namedValue := range names {
			if namedValue == value && (!found || name < firstName) {
				firstName = name
				found = true
			}
		}
	}

	return
}

var currentVariableValue string

// currentVariableUUID is the UUID of the action that produced currentVariableValue, which only the Set Variable
// referencing that UUID may consume.
var currentVariableUUID string

func decompileActions() {
	for index, action := range shortcut.WFWorkflowActions {
		// Set from the loop so an action that is written without reaching the end of it does not shift
		// the actions that are peeked at.
		actionIndex = index
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
			if isLoopVariableAssignment() {
				continue
			}
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

// decompComment writes a comment action as a comment action when comments are included, which is the only
// way for it to compile back into the Shortcut, and as a code comment otherwise.
func decompComment(action *ShortcutAction) {
	var commentText = paramString(action, "WFCommentActionText")
	if args.Using("comments") {
		newCodeLine(fmt.Sprintf("comment(%s)\n", quoteRawString(commentText)))
		return
	}

	// A comment that holds the end of a multi-line comment cannot be written as one.
	if strings.Contains(commentText, "\n") && !strings.Contains(commentText, "*/") {
		newCodeLine("/*\n")
		for _, line := range strings.Split(commentText, "\n") {
			newCodeLine(fmt.Sprintln(line))
		}
		newCodeLine("*/\n")
		return
	}

	for _, line := range strings.Split(commentText, "\n") {
		newCodeLine(fmt.Sprintf("// %s\n", line))
	}
	code.WriteRune('\n')
}

var decompilingText = false

// decompilingTextParameter is set while writing the argument of a text parameter, which keeps text that is only a
// reference as text but, unlike a Text action, can be written as an Ask with a prompt.
var decompilingTextParameter bool

func decompTextValue(action *ShortcutAction) {
	decompilingText = true
	// Shortcuts leaves out the text of a Text action that is empty.
	if action.WFWorkflowActionParameters["WFTextActionText"] == nil {
		currentVariableValue = quoteString("", false)
	} else {
		currentVariableValue = decompValue(action.WFWorkflowActionParameters["WFTextActionText"])
	}
	checkConstantLiteral(action)
	decompilingText = false
}

var numberLiteralRegex = regexp.MustCompile(`^-?[0-9]+(\.[0-9]+)?$`)

func decompNumberValue(action *ShortcutAction) (nonLiteral bool) {
	var value = requiredParam(action, "WFNumberActionNumber")
	if mapValue, isMap := value.(map[string]interface{}); isMap {
		if _, hasValue := mapValue["Value"].(map[string]interface{}); !hasValue {
			decompError(fmt.Sprintf("Expected the parameter 'WFNumberActionNumber' to have a Value dictionary, got %T", mapValue["Value"]), action)
		}

		var literal = strings.Trim(decompValue(value), "\"")
		if !numberLiteralRegex.MatchString(literal) {
			nonLiteral = true
			return
		}

		currentVariableValue = literal
		checkConstantLiteral(action)
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
		case "+", "-":
			operation = paramString(action, "WFMathOperation")
		case "…":
			// A scientific operation is stored in a second key that an infix expression has no way to
			// write, so it is written as the calculate() call that makes it.
			if action.WFWorkflowActionParameters["WFScientificMathOperation"] == nil {
				decompError("Expected the parameter 'WFScientificMathOperation' to name the scientific operation of the math action", action)
			}
			decompAction(action)
			return
		default:
			decompError(fmt.Sprintf("Expected the parameter 'WFMathOperation' to be '+', '-', '×', '÷' or '…', got '%v'", action.WFWorkflowActionParameters["WFMathOperation"]), action)
		}
	}
	var inputText = strings.Trim(decompValue(input), "\"")
	var operandText = strings.Trim(decompValue(operand), "\"")
	var expression = inputText + " " + strings.Trim(operation, "\"") + " " + operandText
	var varRegex = regexp.MustCompile(`{(.*?)}`)
	currentVariableValue = varRegex.ReplaceAllString(expression, "$1")

	checkConstantLiteral(action)

	var compound, isCompound = compoundAssignment(action, inputText, operation, operandText)
	if isCompound {
		currentVariableValue = operandText
		currentVariableOperator = compound
	}
}

// currentVariableOperator is the operator that the Set Variable consuming currentVariableValue is written
// with when the value is the operand of a compound assignment.
var currentVariableOperator string

var integerLiteralRegex = regexp.MustCompile(`^-?[0-9]+$`)

// compoundAssignment reverses variableValueModifier: a math action that has a variable as its first operand,
// an integer as its second and whose output is only taken by the Set Variable right after it that sets that
// variable is what x += 3 compiles to. The compiler only writes that for an integer operand, and only appends
// to a variable that is an array instead, so neither is written as one.
func compoundAssignment(action *ShortcutAction, inputText string, operation string, operandText string) (operator string, isCompound bool) {
	if action.WFWorkflowActionParameters[UUID] == nil || actionIndex+1 >= len(shortcut.WFWorkflowActions) {
		return
	}
	var actionUUID = paramString(action, UUID)
	if currentVariableUUID != actionUUID || !integerLiteralRegex.MatchString(operandText) {
		return
	}

	var next = peekActions(1)
	if next.WFWorkflowActionIdentifier != SetVariableIdentifier {
		return
	}
	var variableName = paramString(&next, "WFVariableName")
	sanitizeIdentifier(&variableName)
	if inputText != "@"+variableName || appendedVariables[variableName] {
		return
	}

	return operation + "=", true
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

// appendedVariables are the variables that an Append to Variable adds to, and assignedVariables those that
// a Set Variable or Append to Variable has written so far.
var appendedVariables map[string]bool
var assignedVariables map[string]bool

func decompVariable(action *ShortcutAction) {
	var variableName = paramString(action, "WFVariableName")
	sanitizeIdentifier(&variableName)
	newCodeLine(fmt.Sprintf("@%s", variableName))

	var firstAssignment = !assignedVariables[variableName]
	assignedVariables[variableName] = true

	// The compiler only appends to a variable that was an array where it was first written, and otherwise
	// sets it, so the first write of a variable that is appended to declares the type.
	if firstAssignment && appendedVariables[variableName] && action.WFWorkflowActionIdentifier == SetVariableIdentifier &&
		action.WFWorkflowActionParameters["WFInput"] == nil {
		code.WriteString(": array\n")
		return
	}

	var inputUUID string
	if action.WFWorkflowActionParameters["WFInput"] != nil {
		var wfInput WFInput
		mapToStruct(action.WFWorkflowActionParameters["WFInput"], &wfInput)
		inputUUID = wfInput.Value.OutputUUID
	}

	if currentVariableValue == "" || inputUUID == "" || inputUUID != currentVariableUUID {
		var decompInput = decompValue(action.WFWorkflowActionParameters["WFInput"])
		if decompInput != "" {
			code.WriteString(fmt.Sprintf(" %s %s", assignmentOperator(action), decompInput))
		}
		code.WriteRune('\n')
		return
	}

	// Assigning an empty dictionary compiles to a Nothing before the Set Variable, so only a type
	// declaration compiles to a Set Variable that takes the dictionary.
	if currentVariableValue == emptyDictionary && action.WFWorkflowActionIdentifier == SetVariableIdentifier {
		code.WriteString(": dictionary\n")
		currentVariableValue = ""
		currentVariableUUID = ""
		return
	}

	var operator = assignmentOperator(action)
	if currentVariableOperator != "" {
		operator = currentVariableOperator
		currentVariableOperator = ""
	}
	code.WriteString(fmt.Sprintf(" %s ", operator))
	code.WriteString(currentVariableValue)
	code.WriteRune('\n')

	currentVariableValue = ""
	currentVariableUUID = ""
}

func assignmentOperator(action *ShortcutAction) string {
	if action.WFWorkflowActionIdentifier == AppendVariableIdentifier {
		return "+="
	}

	return "="
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
		// The body of the item before this one is one level in from the item.
		if openMenuItems[groupingUUID] {
			tabLevel--
		}
		openMenuItems[groupingUUID] = true
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
			code.WriteString(decompValue(menuItemValue(menuItem.value, action)))
		}
		code.WriteString(":\n")
		tabLevel++
	case endStatement:
		if openMenuItems[groupingUUID] {
			tabLevel--
		}
		delete(openMenuItems, groupingUUID)
		tabLevel--
		newCodeLine("}\n")
	}
}

// openMenuItems are the menus that have an item whose body is open.
var openMenuItems = make(map[string]bool)

// menuItemValue is the title of a menu item, which the compiler writes as a WFMenuItem dictionary
// holding it in WFValue, and Shortcuts as a string.
func menuItemValue(item any, action *ShortcutAction) any {
	var menuItem, isMap = item.(map[string]interface{})
	if !isMap {
		return item
	}
	var value, found = menuItem["WFValue"]
	if !found {
		decompError("Expected a menu item dictionary to hold its title in WFValue", action)
	}

	return value
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
			return decompQuantity(operand, "WFDuration", action)
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

func decompQuantity(operand any, key string, action *ShortcutAction) string {
	var duration, isMap = operand.(map[string]interface{})
	if !isMap {
		decompError(fmt.Sprintf("Expected %s to be a WFQuantityFieldValue, got %T", key, operand), action)
	}
	var quantity, hasValue = duration["Value"].(map[string]interface{})
	if !hasValue || quantity["Magnitude"] == nil || quantity["Unit"] == nil {
		decompError(fmt.Sprintf("Expected %s to have a Value with a Magnitude and a Unit", key), action)
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

const emptyDictionary = "{}"

func decompDictionary(action *ShortcutAction) {
	var params DictionaryActionParameters
	mapToStruct(action.WFWorkflowActionParameters, &params)

	currentVariableValue = decompDictionaryValue(params.WFItems.Value.WFDictionaryFieldValueItems)

	checkConstantLiteral(action)
}

func decompDictionaryValue(items []WFDictionaryFieldValueItem) string {
	return decompLiteral(decompDictionaryItems(items))
}

// enterLiteral keeps the text inside of an array or a dictionary as it is until the returned function is called,
// as it is escaped when the literal is written as JSON, and restores what was set so literals can be nested.
func enterLiteral() (leave func()) {
	var wasLiteral = decompilingDictionary
	decompilingDictionary = true

	return func() {
		decompilingDictionary = wasLiteral
	}
}

// bareReference is a variable reference that is written in a literal without quotes.
type bareReference string

func isReferenceValue(value any) bool {
	if mapValue, isMap := value.(map[string]interface{}); isMap {
		if mapValue["Value"] != nil && mapValue["WFSerializationType"] == "WFTextTokenAttachment" {
			return true
		}
	}

	return false
}

// literalItem is a key and value of a dictionary literal.
type literalItem struct {
	key   string
	value any
}

// literalDictionary is a dictionary literal that keeps the order of its items, as the order the items
// are in is the order they show in the Shortcut, which a Go map would not keep.
type literalDictionary []literalItem

func decompDictionaryItems(items []WFDictionaryFieldValueItem) (dictionary literalDictionary) {
	defer enterLiteral()()
	for _, item := range items {
		var itemKey = decompValue(item.WFKey)
		dictionary = append(dictionary, literalItem{key: itemKey, value: decompDictionaryItem(item)})
	}
	return
}

// jsonNumberRegex is a number as JSON writes it, which has no leading zero.
var jsonNumberRegex = regexp.MustCompile(`^-?(0|[1-9][0-9]*)(\.[0-9]+)?$`)

// decompNumberItem writes the text of a number item as a number if it is one, keeping a decimal or a negative
// as it is. A number with a leading zero is written as the number it is, as JSON has no number like it.
func decompNumberItem(text string) (number json.Number, isNumber bool) {
	if !numberLiteralRegex.MatchString(text) {
		return "", false
	}
	if jsonNumberRegex.MatchString(text) {
		return json.Number(text), true
	}

	var value, parseErr = strconv.ParseFloat(text, 64)
	handle(parseErr)

	return json.Number(strconv.FormatFloat(value, 'f', -1, 64)), true
}

func decompDictionaryItem(item WFDictionaryFieldValueItem) any {
	var itemValueMap, isMap = item.WFValue.(map[string]interface{})
	if !isMap {
		decompError(fmt.Sprintf("Expected a dictionary item value to be a dictionary, got %T", item.WFValue), currentDecompAction())
	}
	var itemValue any

	if dictDataType(item.WFItemType) == itemTypeFile {
		return bareReference(decompValue(itemValueMap["Value"]))
	}

	if isReferenceValue(itemValueMap["Value"]) {
		return fmt.Sprintf("{%s}", decompValue(item.WFValue))
	}

	switch dictDataType(item.WFItemType) {
	case itemTypeNumber:
		var numberText = decompValue(item.WFValue)
		var number, isNumber = decompNumberItem(numberText)
		if !isNumber {
			itemValue = numberText
			break
		}
		itemValue = number
	case itemTypeBool:
		var boolValue, isBool = itemValueMap["Value"].(bool)
		if !isBool {
			decompError(fmt.Sprintf("Expected a boolean dictionary item value, got %T", itemValueMap["Value"]), currentDecompAction())
		}
		itemValue = boolValue
	case itemTypeText:
		itemValue = decompValue(item.WFValue)
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
		itemValue = decompValue(item.WFValue)
	}

	return itemValue
}

func decompArray(items []interface{}) (array []interface{}) {
	defer enterLiteral()()
	for _, item := range items {
		if itemMap, isMap := item.(map[string]interface{}); isMap {
			var fieldValueItem WFDictionaryFieldValueItem
			mapToStruct(itemMap, &fieldValueItem)
			array = append(array, decompDictionaryItem(fieldValueItem))
			continue
		}

		array = append(array, decompValue(item))
	}
	return
}

func decompValue(value any) string {
	if value == nil {
		return "nil"
	}

	switch typedValue := value.(type) {
	case map[string]interface{}:
		return decompValueObject(typedValue)
	case string:
		return decompTextLiteral(typedValue, false)
	case []interface{}:
		return decompLiteral(decompArray(typedValue))
	case bool, int, int64, uint64, float64:
		return fmt.Sprintf("%v", typedValue)
	default:
		decompError(fmt.Sprintf("Expected a value to be text, a number, a boolean, an array or a dictionary, got %T", value), currentDecompAction())
		return ""
	}
}

// decompLiteral writes an array or a dictionary as a Cherri literal, which is JSON, so a Go value is never
// written with its Go formatting.
func decompLiteral(value any) string {
	var literal strings.Builder
	writeLiteral(&literal, value, tabLevel)

	return literal.String()
}

// writeLiteral writes value as indented JSON that keeps the order of a dictionary. Each line after the first
// starts at level tabs, as the literal is written where code is indented that much.
func writeLiteral(literal *strings.Builder, value any, level int) {
	var inner = strings.Repeat("\t", level+1)
	var outer = strings.Repeat("\t", level)

	switch typedValue := value.(type) {
	case literalDictionary:
		if len(typedValue) == 0 {
			literal.WriteString(emptyDictionary)
			return
		}

		literal.WriteString("{\n")
		for i, item := range typedValue {
			literal.WriteString(inner)
			writeLiteralScalar(literal, item.key)
			literal.WriteString(": ")
			writeLiteral(literal, item.value, level+1)
			if i < len(typedValue)-1 {
				literal.WriteRune(',')
			}
			literal.WriteRune('\n')
		}
		literal.WriteString(outer + "}")
	case []interface{}:
		if len(typedValue) == 0 {
			literal.WriteString("[]")
			return
		}

		literal.WriteString("[\n")
		for i, item := range typedValue {
			literal.WriteString(inner)
			writeLiteral(literal, item, level+1)
			if i < len(typedValue)-1 {
				literal.WriteRune(',')
			}
			literal.WriteRune('\n')
		}
		literal.WriteString(outer + "]")
	case bareReference:
		literal.WriteString(string(typedValue))
	default:
		writeLiteralScalar(literal, value)
	}
}

// writeLiteralScalar writes text, a number or a boolean as JSON. Text keeps <, > and & as they are, which
// the encoder would escape as it assumes the JSON is put in HTML.
func writeLiteralScalar(literal *strings.Builder, value any) {
	var encoded strings.Builder
	var encoder = json.NewEncoder(&encoded)
	encoder.SetEscapeHTML(false)
	handle(encoder.Encode(value))

	literal.WriteString(strings.TrimSuffix(encoded.String(), "\n"))
}

// stringEscaper is one replacer, rather than a chain of replacements, so the backslashes it adds are
// never escaped again.
var stringEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "\t", `\t`, "\r", `\r`)

// escapeString escapes raw text for the inside of a double quoted Cherri string, the inverse of collectString.
func escapeString(value string) string {
	return stringEscaper.Replace(value)
}

// quoteString is the only place text becomes a Cherri string literal, so every caller passes the
// text as it is in the Shortcut and never trims or escapes it first. Text that has its inline
// references already written as {reference} is passed with hasReferences.
//
// Cherri has no escape for a brace in a double quoted string, and the compiler reads the braces in
// one as inline references, so text with a literal brace is written as a raw string.
func quoteString(text string, hasReferences bool) string {
	if hasReferences || !strings.ContainsAny(text, "{}") {
		return fmt.Sprintf("\"%s\"", escapeString(text))
	}

	return quoteRawString(text)
}

// quoteRawString writes text as a raw string, which only has an escape for the single quote, so a
// backslash at the end of the text would escape the closing quote.
func quoteRawString(text string) string {
	if strings.HasSuffix(text, `\`) {
		decompError(fmt.Sprintf("Expected text that is written as a raw string not to end with a backslash, got '%s'", text), currentDecompAction())
	}

	return fmt.Sprintf("'%s'", strings.ReplaceAll(text, "'", `\'`))
}

// decompTextLiteral writes text as it is written in the current context. Inside an array or a dictionary it
// stays as it is because it is escaped when the literal is written as JSON.
func decompTextLiteral(text string, hasReferences bool) string {
	if decompilingDictionary {
		if !hasReferences && strings.ContainsAny(text, "{}") {
			decompError(fmt.Sprintf("Expected text in an array or a dictionary to have no literal brace, as Cherri reads it as an inline reference, got '%s'", text), currentDecompAction())
		}

		return text
	}

	return quoteString(text, hasReferences)
}

func decompValueObject(value map[string]interface{}) string {
	// The compiler writes an empty dictionary for an argument given as nil.
	if len(value) == 0 {
		return "nil"
	}

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

	var valueType, hasType = value["Type"]
	if !hasType {
		return decompObjectValue(value)
	}

	// A condition input wraps its reference in a Variable dictionary instead of naming it.
	if _, named := value["VariableName"]; valueType == "Variable" && !named {
		var variableValue, hasVariable = value["Variable"].(map[string]interface{})
		if !hasVariable {
			decompError(fmt.Sprintf("Expected a Variable value to have a VariableName or a Variable dictionary, got %T", value["Variable"]), currentDecompAction())
		}
		return decompValue(variableValue["Value"])
	}

	for _, key := range []string{"Type", "VariableName", "OutputUUID", "OutputName", "Prompt"} {
		if field, found := value[key]; found {
			if _, isString := field.(string); !isString {
				decompError(fmt.Sprintf("Expected %s of a reference to be a string, got %T", key, field), currentDecompAction())
			}
		}
	}
	if aggrandizements, found := value["Aggrandizements"]; found {
		if _, isArray := aggrandizements.([]interface{}); !isArray {
			decompError(fmt.Sprintf("Expected Aggrandizements of a reference to be an array, got %T", aggrandizements), currentDecompAction())
		}
	}

	var reference Value
	mapToStruct(value, &reference)

	return decompReference(reference)
}

// decompReference is the only place a plist reference becomes Cherri, so references written as
// whole values and references inside text resolve globals and aggrandizements the same way.
func decompReference(reference Value) string {
	var name string
	switch reference.Type {
	case "Variable":
		if reference.VariableName == "" {
			decompError("Expected a Variable reference to have a VariableName", currentDecompAction())
		}
		if loopReference, isLoopVariable := loopVariableReference(reference.VariableName); isLoopVariable {
			name = loopReference
			break
		}
		name = reference.VariableName
		sanitizeIdentifier(&name)
		name = fmt.Sprintf("@%s", name)
	case "ActionOutput":
		var outputName, found = uuids[reference.OutputUUID]
		if !found {
			decompError(fmt.Sprintf("Expected the output UUID '%s' to belong to an earlier action", reference.OutputUUID), currentDecompAction())
		}
		sanitizeIdentifier(&outputName)
		name = outputName
	default:
		name = decompGlobal(reference.Type)
		if reference.Type == Ask && reference.Prompt != "" {
			name = fmt.Sprintf("%s: %s", name, quoteRawString(reference.Prompt))
		}
	}

	decompAggrandizements(&name, reference.Aggrandizements)

	return name
}

// decompGlobal skips the globals of type Variable, the repeat variables, because Shortcuts names them.
func decompGlobal(referenceType string) string {
	for name, global := range globals {
		if global.variableType != "Variable" && global.variableType == referenceType {
			return name
		}
	}

	var globalTypes []string
	for _, global := range globals {
		if global.variableType != "Variable" {
			globalTypes = append(globalTypes, global.variableType)
		}
	}
	sort.Strings(globalTypes)
	decompError(fmt.Sprintf("Expected a reference of type Variable, ActionOutput or a global (%s), got '%s'", strings.Join(globalTypes, ", "), referenceType), currentDecompAction())

	return ""
}

func decompObjectValue(valueObj any) string {
	var value, isMap = valueObj.(map[string]interface{})
	if !isMap {
		return decompValue(valueObj)
	}

	// Shortcuts writes empty text as a text token string without a Value.
	if value["Value"] == nil && value["WFSerializationType"] == "WFTextTokenString" {
		return decompTextLiteral("", false)
	}
	if value["Value"] != nil {
		var nestedValue, isNestedMap = value["Value"].(map[string]interface{})
		if !isNestedMap {
			return decompValue(value["Value"])
		}
		value = nestedValue
	}

	var _, hasString = value["string"]
	var _, hasAttachments = value["attachmentsByRange"]
	if !hasString && !hasAttachments {
		var keys = slices.Sorted(maps.Keys(value))
		decompError(fmt.Sprintf("Expected a dictionary value to be text, a variable reference or a dictionary, got one with the keys: %s", strings.Join(keys, ", ")), currentDecompAction())
	}

	var text string
	if _, found := value["string"]; found {
		var valueString, isString = value["string"].(string)
		if !isString {
			decompError(fmt.Sprintf("Expected string of a text value to be a string, got %T", value["string"]), currentDecompAction())
		}
		text = valueString
	}

	var attachments, found = value["attachmentsByRange"]
	if !found {
		return decompTextLiteral(text, false)
	}
	var attachmentsMap, isAttachmentsMap = attachments.(map[string]interface{})
	if !isAttachmentsMap {
		decompError(fmt.Sprintf("Expected attachmentsByRange to be a dictionary, got %T", attachments), currentDecompAction())
	}

	return decompAttachmentString(text, attachmentsMap)
}

// decompAttachmentString substitutes each attachment into the text and writes it as text. Attachment
// ranges are UTF-16 code unit offsets, the inverse of mapInlineVariables, so the text is indexed as
// UTF-16 rather than as code points, which would misplace attachments after an emoji.
func decompAttachmentString(text string, attachments map[string]interface{}) string {
	var units = utf16.Encode([]rune(text))

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
	var isWholeValue = text == ObjectReplaceCharStr
	var lastReference string
	var hasPrompt bool
	for _, position := range positions {
		if position < consumed || position >= len(units) {
			decompError(fmt.Sprintf("Attachment position %d is outside of the string", position), currentDecompAction())
		}

		var attachment Value
		mapToStruct(attachmentsByPosition[position], &attachment)
		// Cherri can only write the prompt of an Ask on a whole value, not inside of text or an array or a dictionary.
		var keepsPrompt = isWholeValue && !decompilingDictionary && !decompilingText
		if attachment.Prompt != "" && !keepsPrompt {
			decompError(fmt.Sprintf("Expected the Ask in text to have no prompt, as Cherri can only write one on an Ask that is a whole value, got '%s'", attachment.Prompt), currentDecompAction())
		}
		if slices.Contains(controlFlowOutputUUIDs, attachment.OutputUUID) {
			referencesControlFlowOutput = true
		}

		lastReference = decompReference(attachment)
		hasPrompt = attachment.Prompt != ""
		result.WriteString(decompLiteralText(units[consumed:position]))
		result.WriteString(fmt.Sprintf("{%s}", lastReference))
		consumed = position + 1
	}
	result.WriteString(decompLiteralText(units[consumed:]))

	// The output of a control flow block has no type, so the compiler only accepts it inside text.
	// A text parameter keeps the text so the reference is converted to text, which a bare variable of
	// another type is not.
	var keepsText = decompilingText || decompilingTextParameter
	if isWholeValue && !referencesControlFlowOutput && !decompilingDictionary && (!keepsText || hasPrompt) {
		return lastReference
	}

	return decompTextLiteral(result.String(), true)
}

// decompLiteralText decodes text between two attachments, which Cherri can only write when it has no brace
// because the compiler would read the braces with the references as one reference.
func decompLiteralText(units []uint16) string {
	var text = string(utf16.Decode(units))
	if strings.ContainsAny(text, "{}") {
		decompError(fmt.Sprintf("Expected text with inline references to have no literal brace, got '%s'", text), currentDecompAction())
	}

	return text
}

// decompAggrandizements writes the aggrandizements in plist order. The compiler reads at most one
// key or property followed by at most one coercion, so any other order cannot be written as Cherri.
func decompAggrandizements(reference *string, aggrandizements []Aggrandizement) {
	var revContentItems = reversedContentItems()
	var keyed, coerced bool
	for _, aggrandizement := range aggrandizements {
		switch aggrandizement.Type {
		case "WFDictionaryValueVariableAggrandizement", "WFPropertyVariableAggrandizement":
			if keyed || coerced {
				decompError("Expected at most one dictionary key or property read on a reference, before any coercion", currentDecompAction())
			}
			var key = aggrandizement.DictionaryKey
			if aggrandizement.Type == "WFPropertyVariableAggrandizement" {
				key = aggrandizement.PropertyName
			}
			if key == "" || strings.Contains(key, "'") {
				decompError(fmt.Sprintf("Expected a dictionary key or property name that is not empty and has no single quote, got '%s'", key), currentDecompAction())
			}
			*reference = fmt.Sprintf("%s['%s']", *reference, key)
			keyed = true
		case "WFCoercionVariableAggrandizement":
			if coerced {
				decompError("Expected at most one coercion on a reference", currentDecompAction())
			}
			var coerce, found = revContentItems[aggrandizement.CoercionItemClass]
			if !found {
				decompError(fmt.Sprintf("Expected a coercion to a known content item class, got '%s'", aggrandizement.CoercionItemClass), currentDecompAction())
			}
			*reference = fmt.Sprintf("%s.%s", *reference, coerce)
			coerced = true
		default:
			decompError(fmt.Sprintf("Expected a coercion, dictionary key or property aggrandizement, got '%s'", aggrandizement.Type), currentDecompAction())
		}
	}
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
		decompiledDefinitions = append(decompiledDefinitions, fmt.Sprintf("#define mac %v", macDefinition))
		setMacDefinition = true
	}

	actionCallCode.WriteString(fmt.Sprintf("%s(", matchedIdentifier))

	if matchedAction.makeParams != nil || matchedAction.decomp != nil {
		decompActionCustom(&actionCallCode, matchedIdentifier, &matchedAction, action)
	} else {
		var matchedParamsSize = len(matchedAction.parameters)
		if matchedParamsSize > 0 {
			decompActionArguments(&actionCallCode, matchedIdentifier, &matchedAction, action)
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

	return identifier == "nothing" && isCompilerNothing()
}

// isCompilerNothing reports whether the current Nothing action is one the compiler inserts, which is
// written again when the code is compiled, so it must not be written. Each is a Nothing that:
//   - follows the end of a block, which the compiler always adds so no output is passed on
//   - comes before the end of a block, or before an item of a menu, that does not capture its output
//   - comes before a Set Variable with no input, which a variable set to nil is compiled to
func isCompilerNothing() bool {
	if actionIndex > 0 {
		var previous = peekActions(-1)
		if isControlFlowAction(&previous) && controlFlowMode(&previous) == endStatement {
			return true
		}
	}
	if actionIndex+1 >= len(shortcut.WFWorkflowActions) {
		return false
	}

	var next = peekActions(1)
	if next.WFWorkflowActionIdentifier == SetVariableIdentifier {
		return next.WFWorkflowActionParameters["WFInput"] == nil
	}
	if !isControlFlowAction(&next) {
		return false
	}

	var _, capturesOutput = controlFlowOutputs[paramString(&next, "GroupingIdentifier")]
	switch controlFlowMode(&next) {
	case endStatement:
		return !capturesOutput
	case statementPart:
		// The compiler adds no Nothing before an else.
		return next.WFWorkflowActionIdentifier == "is.workflow.actions.choosefrommenu" && !capturesOutput
	}

	return false
}

func isControlFlowAction(action *ShortcutAction) bool {
	return slices.Contains(controlFlowIdentifiers, action.WFWorkflowActionIdentifier)
}

func actionIdentifierEnd(identifier string) string {
	return strings.Replace(identifier, "is.workflow.actions.", "", 1)
}

func decompActionArguments(actionCallCode *strings.Builder, name string, matchedAction *actionDefinition, action *ShortcutAction) {
	actionCallCode.WriteString(strings.Join(decompParameters(name, matchedAction.parameters, action, nil), ", "))
}

// decompHookName and decompHookDefinition are the action whose decomp hook is running. Hooks cannot
// look their own definition up in actions, as they are part of its initialization.
var decompHookName string
var decompHookDefinition *actionDefinition

// decompHookArguments is decompParameters for a decomp hook, which passes the values it recovered
// in recovered by parameter name.
func decompHookArguments(action *ShortcutAction, recovered map[string]string) []string {
	return decompParameters(decompHookName, decompHookDefinition.parameters, action, recovered)
}

// decompParameters writes one argument per parameter, as arguments are matched to parameters by
// position. An absent parameter that can be skipped is written as nil so later arguments keep their
// position, and the skipped ones left at the end are trimmed. A decomp hook passes the values it
// recovered in recovered, which are used in place of the plist, and a parameter without a plist key
// can only be recovered by a hook.
func decompParameters(name string, parameters []parameterDefinition, action *ShortcutAction, recovered map[string]string) []string {
	var arguments []string
	var written int
	for _, param := range parameters {
		var skippable = param.optional || param.defaultValue != nil
		if identifier, isQuestion := questionArgument(param.key); isQuestion {
			arguments = append(arguments, identifier)
			written = len(arguments)
			continue
		}
		if argValue, isRecovered := recovered[param.name]; isRecovered {
			arguments = append(arguments, argValue)
			if argValue != "nil" || !skippable {
				written = len(arguments)
			}
			continue
		}

		var value, found = action.WFWorkflowActionParameters[param.key]
		switch {
		case param.key == "" && recovered == nil:
			decompError(fmt.Sprintf("Expected %s() to have a decomp hook, as its parameter '%s' has no plist key to read it from", name, param.name), action)
		case param.key == "" && !skippable:
			decompError(fmt.Sprintf("Expected the decomp hook of %s() to recover its required parameter '%s'", name, param.name), action)
		case (param.key == "" || !found) && skippable:
			arguments = append(arguments, "nil")
			continue
		case !found && param.literal:
			arguments = append(arguments, makeDefaultValue(param))
		case !found:
			// The compiler writes nothing for a required argument given as nil.
			arguments = append(arguments, "nil")
		default:
			arguments = append(arguments, decompArgument(param, value, action))
			if end(arguments) == "nil" && skippable {
				continue
			}
		}
		written = len(arguments)
	}

	return arguments[:written]
}

// decompArgument writes the plist value of a parameter, which for a color or a quantity is written
// as the color() or qty() Cherri reads it from.
func decompArgument(param parameterDefinition, value any, action *ShortcutAction) string {
	if isReferenceValue(value) {
		return decompValue(value)
	}

	switch param.validType {
	case Quantity:
		var quantity = decompQuantity(value, param.key, action)
		checkDecompEnum(param, value.(map[string]interface{})["Value"].(map[string]interface{})["Unit"], action)
		return quantity
	case Color:
		return decompColor(value, param.key, action)
	}

	checkDecompEnum(param, value, action)

	decompilingTextParameter = param.validType == String
	var argValue = decompValue(value)
	decompilingTextParameter = false
	if param.validType == Integer && !startsWith(Ask, argValue) {
		argValue = strings.Trim(argValue, "\"")
	}

	return argValue
}

// checkDecompEnum stops on a text value that is not in the enum of param, as Cherri would reject it.
func checkDecompEnum(param parameterDefinition, value any, action *ShortcutAction) {
	var text, isText = value.(string)
	if param.enum == "" || !isText {
		return
	}

	var enum = getEnum(param.enum)
	if !slices.Contains(enum, text) {
		decompError(fmt.Sprintf("Expected the value of '%s' to be one of: %s, got '%s'", param.name, strings.Join(enum, ", "), text), action)
	}
}

// decompColor writes a WFColorValue as color(). Shortcuts leaves a component out when it is 0.
func decompColor(value any, key string, action *ShortcutAction) string {
	var colorValue, isMap = value.(map[string]interface{})
	if !isMap || colorValue["WFColorRepresentationType"] != "WFColorRepresentationTypeCGColor" {
		decompError(fmt.Sprintf("Expected %s to be a CGColor color value", key), action)
	}

	var components []string
	for _, component := range []string{"redComponent", "greenComponent", "blueComponent", "alphaComponent"} {
		var number float64
		switch componentValue := colorValue[component].(type) {
		case nil:
		case float64:
			number = componentValue
		case uint64:
			number = float64(componentValue)
		default:
			decompError(fmt.Sprintf("Expected %s of %s to be a number, got %T", component, key, componentValue), action)
		}
		// Color components are float literals, which need a decimal point to be read as a float.
		components = append(components, strconv.FormatFloat(number, 'f', -1, 64))
		if !strings.Contains(end(components), ".") {
			components[len(components)-1] += ".0"
		}
	}

	return fmt.Sprintf("color(%s)", strings.Join(components, ", "))
}

func makeDefaultValue(param parameterDefinition) string {
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
			var leave = enterLiteral()
			rawParams[key] = decompValueObject(mapValue)
			leave()
			continue
		}

		rawParams[key] = value
	}

	return rawParams
}

func matchAction(action *ShortcutAction) (name string, definition actionDefinition) {
	var candidates = identifierMap[action.WFWorkflowActionIdentifier]
	if len(candidates) == 0 {
		return
	}

	name = candidates[0].identifier
	definition = *candidates[0].definition
	if len(candidates) == 1 {
		return
	}

	matchSplitAction(&candidates, action, &name, &definition)
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

	return
}

type actionMatch struct {
	params  uint
	values  uint
	missing uint
	action  actionValue
}

// matchSplitAction picks which of the actions sharing an identifier wrote the parameters. Of the actions
// that could have written them, the one that matches the most parameters, then the most values, then
// misses the fewest of its fixed parameters is picked, and on a tie the default action is preferred.
func matchSplitAction(splitActions *[]actionValue, action *ShortcutAction, identifier *string, definition *actionDefinition) {
	var parameters = action.WFWorkflowActionParameters
	if args.Using("debug") {
		fmt.Println("## MATCHING SPLIT ACTIONS ##")
		fmt.Println("parameters", parameters)
	}

	var matches []actionMatch
	var names []string
	for _, splitAction := range *splitActions {
		names = append(names, splitAction.identifier)
		var missing, conflicts = splitActionConflicts(splitAction, parameters)
		if conflicts {
			continue
		}

		var matchedParams, matchedValues = scoreActionMatch(splitAction, splitAction.definition.parameters, parameters)
		matches = append(matches, actionMatch{
			params:  matchedParams,
			values:  matchedValues,
			missing: missing,
			action:  splitAction,
		})
	}
	if len(matches) == 0 {
		decompError(fmt.Sprintf("Expected the parameters to fit one of the actions %s", strings.Join(names, ", ")), action)
	}

	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].params != matches[j].params {
			return matches[i].params > matches[j].params
		}
		if matches[i].values != matches[j].values {
			return matches[i].values > matches[j].values
		}
		if matches[i].missing != matches[j].missing {
			return matches[i].missing < matches[j].missing
		}
		if matches[i].action.definition.defaultAction != matches[j].action.definition.defaultAction {
			return matches[i].action.definition.defaultAction
		}
		return matches[i].action.identifier < matches[j].action.identifier
	})

	if args.Using("debug") {
		for _, match := range matches {
			fmt.Printf("%s()\n", match.action.identifier)
			fmt.Println("params:", match.params, " values:", match.values, " missing:", match.missing)
			fmt.Println("---")
		}
		fmt.Print("\n")
	}

	*identifier = matches[0].action.identifier
	*definition = *matches[0].action.definition
}

// splitActionConflicts reports whether splitAction cannot have written the parameters: a parameter it
// always writes holds another value, or one of its variable parameters holds a plain dictionary, which
// a variable is never written as. It also counts the parameters it always writes that are missing, as
// Shortcuts leaves a parameter out when it holds its default.
func splitActionConflicts(splitAction actionValue, parameters map[string]any) (missing uint, conflicts bool) {
	for key, value := range splitAction.definition.appendParams {
		var parameter, found = parameters[key]
		if !found {
			missing++
			continue
		}
		if !samePlistValue(value, parameter) {
			return missing, true
		}
	}

	for _, param := range splitAction.definition.parameters {
		if param.key == "" || param.validType != Variable {
			continue
		}
		var dictionary, isDictionary = parameters[param.key].(map[string]interface{})
		if !isDictionary {
			continue
		}
		if _, serialized := dictionary["WFSerializationType"]; !serialized {
			return missing, true
		}
	}

	return missing, false
}

// samePlistValue compares a value written by Cherri with one read from a plist, which holds the same
// data in other Go types, such as uint64 for an int, so both are compared as JSON.
func samePlistValue(written any, read any) bool {
	var writtenJSON, writtenErr = json.Marshal(written)
	handle(writtenErr)
	var readJSON, readErr = json.Marshal(read)
	handle(readErr)

	return string(writtenJSON) == string(readJSON)
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
			if param.defaultValue != nil && samePlistValue(param.defaultValue, value) {
				matchedValues++
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

			if samePlistValue(param.defaultValue, value) {
				matchedValues++
			}
		}
	}
	return
}

func decompActionCustom(actionCode *strings.Builder, name string, matchedAction *actionDefinition, action *ShortcutAction) {
	if matchedAction.decomp == nil {
		decompError(fmt.Sprintf("Expected %s() to have a decomp hook, as it builds its parameters with makeParams", name), action)
	}

	decompHookName = name
	decompHookDefinition = matchedAction
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
