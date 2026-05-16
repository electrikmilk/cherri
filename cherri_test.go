/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/electrikmilk/args-parser"
)

var currentTest string

func TestCherri(_ *testing.T) {
	args.Args["no-ansi"] = ""
	args.Args["comments"] = ""
	var files, err = os.ReadDir("tests")
	if err != nil {
		fmt.Println(ansi("FAILED: unable to read tests directory", red))
		panic(err)
	}
	loadStandardActions()
	for _, file := range files {
		if !strings.Contains(file.Name(), ".cherri") || file.Name() == "decomp-expected.cherri" || file.Name() == "decomp-me.cherri" {
			continue
		}
		currentTest = fmt.Sprintf("tests/%s", file.Name())
		os.Args[1] = currentTest
		fmt.Println(ansi(currentTest, underline, bold))

		compile()

		fmt.Println(ansi("✅  PASSED", green, bold))
		fmt.Print("\n")

		resetParser()

		if signFailed {
			fmt.Println(ansi("Using remote service HubSign", cyan, bold))
			for i := 5; i > 0; i-- {
				fmt.Print(ansi(fmt.Sprintf("Respectfully waiting %d second(s) between tests...\r", i), cyan))
				time.Sleep(1 * time.Second)
			}
			fmt.Print("\n")
		}
	}
}

func TestCherriNoSign(t *testing.T) {
	args.Args["skip-sign"] = ""
	TestCherri(t)
}

func TestPackages(t *testing.T) {
	args.Args["no-ansi"] = ""

	if _, statErr := os.Stat("info.plist"); !os.IsNotExist(statErr) {
		var removeErr = os.Remove("info.plist")
		handle(removeErr)
	}

	if _, statErr := os.Stat("./packages"); !os.IsNotExist(statErr) {
		var removeDirErr = os.RemoveAll("./packages")
		handle(removeDirErr)
	}

	args.Args["init"] = "@electrikmilk/package-test"
	initPackage()
	delete(args.Args, "init")

	args.Args["install"] = "https://github.com/electrikmilk/package-example"
	input := []byte("y")
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}

	// Restore stdin right after the test.
	defer func(v *os.File) { os.Stdin = v }(os.Stdin)
	os.Stdin = r

	addPackage()
	fmt.Println("You entered:", string(input))
	delete(args.Args, "install")

	listPackage()

	listPackages()

	args.Args["remove"] = "@electrikmilk/package-example"
	removePackage()
	delete(args.Args, "remove")
}

func TestDecomp(t *testing.T) {
	defer resetParser()

	fmt.Println("Decompiling...")
	args.Args["import"] = "tests/decomp-me.plist"
	decompile(importShortcut(args.Value("import")))

	fmt.Println("Comparing to expected...")
	var bytes, readErr = os.ReadFile("tests/decomp-expected.cherri")
	handle(readErr)

	if code.String() != string(bytes) {
		fmt.Println(ansi("Does not match expected!", red, bold))
		t.Fail()
		return
	}
	fmt.Print(ansi("✅  PASSED", green, bold) + "\n\n")
}

// TestRoundTrip decompiles compiler-generated plists to verify the full compile→decompile
// pipeline. The _unsigned.shortcut files are produced by TestCherriNoSign; individual
// sub-tests skip gracefully when the file is absent rather than failing.
//
// Run independently:
//
//	go test -run TestCherriNoSign && go test -run TestRoundTrip
func TestRoundTrip(t *testing.T) {
	args.Args["no-ansi"] = ""

	// Chosen because their plist structures are within the decompiler's action handlers.
	var candidates = []string{
		"tests/calc_unsigned.shortcut",
		"tests/numbers_unsigned.shortcut",
		"tests/repeats_unsigned.shortcut",
		"tests/conditionals_unsigned.shortcut",
		"tests/dictionary_unsigned.shortcut",
		"tests/variables_unsigned.shortcut",
	}

	for _, plistPath := range candidates {
		t.Run(plistPath, func(t *testing.T) {
			defer resetParser()

			if _, statErr := os.Stat(plistPath); os.IsNotExist(statErr) {
				t.Skipf("compile output absent — run TestCherriNoSign first: %s", plistPath)
			}

			// Direct decompiler output to /dev/null so no .cherri files land in
			// tests/, which would be picked up and compiled by TestCherriNoSign.
			args.Args["output"] = os.DevNull
			args.Args["import"] = plistPath
			decompile(importShortcut(args.Value("import")))
			delete(args.Args, "output")

			if code.Len() == 0 {
				t.Errorf("decompile of %s produced empty output", plistPath)
			}
		})
	}
}

func TestActionIdentifiers(t *testing.T) {
	args.Args["no-ansi"] = ""
	args.Args["skip-sign"] = ""
	loadStandardActions()

	currentTest = "tests/zz-action-identifiers.cherri"
	os.Args[1] = currentTest

	compile()

	var expected = []string{
		"is.workflow.actions.shortid",
		"is.workflow.actions.two.parts",
		"is.workflow.actions.text.match.getgroup",
		"notion.id.CreatePageIntent",
		"com.apple.facetime.facetime",
	}

	var actual []string
	for _, a := range shortcut.WFWorkflowActions {
		actual = append(actual, a.WFWorkflowActionIdentifier)
	}

	if len(actual) != len(expected) {
		t.Fatalf("Expected %d actions, got %d: %v", len(expected), len(actual), actual)
	}

	for i, ident := range expected {
		if actual[i] != ident {
			t.Errorf("Action %d: got %q, want %q", i, actual[i], ident)
		}
	}

	resetParser()
}

func TestRawActionQuantityFieldValue(t *testing.T) {
	defer resetParser()
	args.Args["no-ansi"] = ""
	args.Args["skip-sign"] = ""

	currentTest = "tests/raw-action-health-quantity.cherri"
	os.Args[1] = currentTest
	compile()

	var action = shortcut.WFWorkflowActions[len(shortcut.WFWorkflowActions)-1]
	if action.WFWorkflowActionIdentifier != "is.workflow.actions.health.quantity.log" {
		t.Fatalf("expected health quantity log action, got %q", action.WFWorkflowActionIdentifier)
	}

	var quantity, quantityOK = action.WFWorkflowActionParameters["WFQuantitySampleQuantity"].(WFQuantityFieldValue)
	if !quantityOK {
		t.Fatalf("expected quantity sample to be WFQuantityFieldValue, got %T", action.WFWorkflowActionParameters["WFQuantitySampleQuantity"])
	}
	var magnitude, magnitudeOK = quantity.Value.Magnitude.(Value)
	if !magnitudeOK {
		t.Fatalf("expected quantity magnitude to be variable Value, got %T", quantity.Value.Magnitude)
	}
	if magnitude.OutputName != "high" || magnitude.Type != "ActionOutput" {
		t.Fatalf("expected quantity magnitude to reference high, got %#v", magnitude)
	}

	var additionalQuantity, additionalQuantityOK = action.WFWorkflowActionParameters["WFQuantitySampleAdditionalQuantity"].(WFQuantityFieldValue)
	if !additionalQuantityOK {
		t.Fatalf("expected additional quantity to be WFQuantityFieldValue, got %T", action.WFWorkflowActionParameters["WFQuantitySampleAdditionalQuantity"])
	}
	var additionalMagnitude, additionalMagnitudeOK = additionalQuantity.Value.Magnitude.(Value)
	if !additionalMagnitudeOK {
		t.Fatalf("expected additional quantity magnitude to be variable Value, got %T", additionalQuantity.Value.Magnitude)
	}
	if additionalMagnitude.OutputName != "low" || additionalMagnitude.Type != "ActionOutput" {
		t.Fatalf("expected additional quantity magnitude to reference low, got %#v", additionalMagnitude)
	}
}

func TestHealthQuantityActions(t *testing.T) {
	defer resetParser()
	args.Args["no-ansi"] = ""
	args.Args["skip-sign"] = ""
	loadStandardActions()

	currentTest = "tests/health.cherri"
	os.Args[1] = currentTest
	compile()

	if len(shortcut.WFWorkflowActions) < 7 {
		t.Fatalf("expected health actions to compile, got %d actions", len(shortcut.WFWorkflowActions))
	}

	var bloodPressure = findCompiledAction("is.workflow.actions.health.quantity.log", "WFQuantitySampleAdditionalQuantity")
	if bloodPressure == nil {
		t.Fatal("expected blood pressure action")
	}
	if bloodPressure.WFWorkflowActionIdentifier != "is.workflow.actions.health.quantity.log" {
		t.Fatalf("expected blood pressure action, got %q", bloodPressure.WFWorkflowActionIdentifier)
	}
	if sampleType := bloodPressure.WFWorkflowActionParameters["WFQuantitySampleType"]; sampleType != "Systolic Blood Pressure" {
		t.Fatalf("expected blood pressure sample type, got %#v", sampleType)
	}
	if _, ok := bloodPressure.WFWorkflowActionParameters["WFQuantitySampleQuantity"].(WFQuantityFieldValue); !ok {
		t.Fatalf("expected systolic quantity field, got %T", bloodPressure.WFWorkflowActionParameters["WFQuantitySampleQuantity"])
	}
	if _, ok := bloodPressure.WFWorkflowActionParameters["WFQuantitySampleAdditionalQuantity"].(WFQuantityFieldValue); !ok {
		t.Fatalf("expected diastolic quantity field, got %T", bloodPressure.WFWorkflowActionParameters["WFQuantitySampleAdditionalQuantity"])
	}

	var steps = findCompiledAction("is.workflow.actions.health.quantity.log", "WFQuantitySampleDate")
	if steps == nil {
		t.Fatal("expected dated steps action")
	}
	if sampleType := steps.WFWorkflowActionParameters["WFQuantitySampleType"]; sampleType != "Steps" {
		t.Fatalf("expected steps sample type, got %#v", sampleType)
	}

	var heartRate = findCompiledActionWithParamValue("is.workflow.actions.health.quantity.log", "WFQuantitySampleType", "Heart Rate")
	if heartRate == nil {
		t.Fatal("expected heart rate action")
	}
	if _, ok := heartRate.WFWorkflowActionParameters["WFQuantitySampleQuantity"].(WFQuantityFieldValue); !ok {
		t.Fatalf("expected heart rate quantity field, got %T", heartRate.WFWorkflowActionParameters["WFQuantitySampleQuantity"])
	}

	if findCompiledAction("is.workflow.actions.filter.health.quantity", "WFContentItemInputParameter") == nil {
		t.Fatal("expected find health samples action")
	}
	if findCompiledAction("is.workflow.actions.properties.health.quantity", "WFContentItemPropertyName") == nil {
		t.Fatal("expected get health sample detail action")
	}
	if findCompiledAction("com.apple.Health.OpenViewIntent", "target") == nil {
		t.Fatal("expected open health view action")
	}
	if findCompiledAction("com.apple.Health.OpenDataTypeIntent", "target") == nil {
		t.Fatal("expected open health data action")
	}
	if findCompiledAction("com.apple.Health.OpenCategoryIntent", "target") == nil {
		t.Fatal("expected open health category action")
	}
	if findCompiledAction("com.apple.Health.OpenRecordsIntent", "target") == nil {
		t.Fatal("expected open health records action")
	}
	if findCompiledAction("com.apple.Health.OpenSearchIntent", "searchPhrase") == nil {
		t.Fatal("expected open health search action")
	}
	if findCompiledAction("com.apple.Health.OpenSleepScheduleIntentV2", "AppIntentDescriptor") == nil {
		t.Fatal("expected open sleep schedule action")
	}
	if findCompiledAction("com.apple.Health.OpenTabIntent", "target") == nil {
		t.Fatal("expected open health tab action")
	}
	if findCompiledAction("is.workflow.actions.health.workout.log", "WFWorkoutCaloriesQuantity") == nil {
		t.Fatal("expected log workout action")
	}
	if findCompiledAction("is.workflow.actions.health.workout.log", "WFWorkoutDuration") == nil {
		t.Fatal("expected log workout duration")
	}
	if findCompiledAction("is.workflow.actions.health.workout.log", "WFWorkoutDistanceQuantity") == nil {
		t.Fatal("expected log workout distance")
	}
	if findCompiledAction("is.workflow.actions.health.quantity.log", "WFCategorySampleEnumeration") == nil {
		t.Fatal("expected log health category action")
	}
	if findCompiledAction("com.apple.ShortcutsActions.GetPhysicalActivity", "AppIntentDescriptor") == nil {
		t.Fatal("expected get physical activity action")
	}
}

func findCompiledAction(identifier string, requiredParam string) *ShortcutAction {
	for i := range shortcut.WFWorkflowActions {
		var action = &shortcut.WFWorkflowActions[i]
		if action.WFWorkflowActionIdentifier != identifier {
			continue
		}
		if _, found := action.WFWorkflowActionParameters[requiredParam]; found {
			return action
		}
	}
	return nil
}

func findCompiledActionWithParamValue(identifier string, param string, value any) *ShortcutAction {
	for i := range shortcut.WFWorkflowActions {
		var action = &shortcut.WFWorkflowActions[i]
		if action.WFWorkflowActionIdentifier != identifier {
			continue
		}
		if action.WFWorkflowActionParameters[param] == value {
			return action
		}
	}
	return nil
}

func TestDecompileRawActionQuantityFieldValue(t *testing.T) {
	defer resetParser()
	uuids = map[string]string{
		"high-uuid": "high",
		"low-uuid":  "low",
	}
	var action = ShortcutAction{
		WFWorkflowActionIdentifier: "is.workflow.actions.health.quantity.log",
		WFWorkflowActionParameters: map[string]any{
			"UUID": "action-uuid",
			"WFQuantitySampleQuantity": map[string]any{
				"Value": map[string]any{
					"Magnitude": map[string]any{
						"OutputName": "high",
						"OutputUUID": "high-uuid",
						"Type":       "ActionOutput",
					},
					"Unit": "mmHg",
				},
				"WFSerializationType": "WFQuantityFieldValue",
			},
			"WFQuantitySampleAdditionalQuantity": map[string]any{
				"Value": map[string]any{
					"Magnitude": map[string]any{
						"OutputName": "low",
						"OutputUUID": "low-uuid",
						"Type":       "ActionOutput",
					},
					"Unit": "mmHg",
				},
				"WFSerializationType": "WFQuantityFieldValue",
			},
			"WFQuantitySampleType": "Systolic Blood Pressure",
		},
	}

	var rawAction = makeRawAction(&action)
	if !strings.Contains(rawAction, `"Magnitude": "${high}"`) {
		t.Fatalf("expected decompiled raw action to reference high, got:\n%s", rawAction)
	}
	if !strings.Contains(rawAction, `"Magnitude": "${low}"`) {
		t.Fatalf("expected decompiled raw action to reference low, got:\n%s", rawAction)
	}
	if !strings.Contains(rawAction, `"WFSerializationType": "WFQuantityFieldValue"`) {
		t.Fatalf("expected decompiled raw action to keep quantity serialization, got:\n%s", rawAction)
	}
}

func TestScoreActionParamsSkipsInvalidEnumValue(t *testing.T) {
	enumerations["testHTTPMethod"] = []string{"POST", "PUT", "PATCH", "DELETE"}
	defer delete(enumerations, "testHTTPMethod")

	var parameters = []parameterDefinition{
		{
			key:  "WFHTTPMethod",
			enum: "testHTTPMethod",
		},
	}
	var matchedParams, matchedValues = scoreActionParams(&parameters, map[string]any{
		"WFHTTPMethod": "GET",
	})

	if matchedParams != 0 || matchedValues != 0 {
		t.Fatalf("expected invalid enum value not to match, got params=%d values=%d", matchedParams, matchedValues)
	}
}

func TestDecompileTextPartsKeepsEmptyCustomSeparator(t *testing.T) {
	defer resetParser()
	var arguments = decompTextParts(&ShortcutAction{
		WFWorkflowActionParameters: map[string]any{
			"text":                  "items",
			"WFTextSeparator":       "Custom",
			"WFTextCustomSeparator": "",
		},
	})

	if len(arguments) != 2 || arguments[1] != `""` {
		t.Fatalf("expected empty custom separator argument, got %#v", arguments)
	}
}

func TestCapitalizeEmptyString(t *testing.T) {
	if got := capitalize(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestSanitizeIdentifierWhitespaceOnly(t *testing.T) {
	identifier := " "
	sanitizeIdentifier(&identifier)

	if identifier != "" {
		t.Fatalf("expected sanitized identifier to be empty, got %q", identifier)
	}
}

func compile() {
	defer func() {
		if recover() != nil {
			panicDebug(nil)
		}
	}()

	main()
}

func resetParser() {
	lines = []string{}
	chars = []rune{}
	char = -1
	idx = 0
	lineIdx = 0
	lineCharIdx = -1
	controlFlowGroups = map[int]controlFlowGroup{}
	groupingIdx = 0
	variables = map[string]varValue{}
	iconColor = -1263359489
	iconGlyph = 61440
	clientVersion = "900"
	iosVersion = 26.0
	questions = map[string]*question{}
	hasShortcutInputVariables = false
	tabLevel = 0
	definedWorkflowTypes = []string{}
	inputs = []string{}
	outputs = []string{}
	noInput = map[string]any{}
	tokens = []token{}
	included = []string{}
	includes = []include{}
	workflowName = ""
	menus = map[string][]varValue{}
	uuids = map[string]string{}
	functions = map[string]*function{}
	shortcut = Shortcut{}
	actionIndex = 0
	code.Reset()
	varUUIDs = nil
	constUUIDs = nil
	identifierMap = nil
	currentVariableValue = ""
	decompilingText = false
	decompilingDictionary = false
	macDefinition = false
	setMacDefinition = false
	appIds = nil
	pasteables = nil
	usedEnums = nil
	usingFunctions = false
	currentCategory = ""
}
