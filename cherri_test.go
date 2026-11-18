/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/electrikmilk/args-parser"
	"howett.net/plist"
)

var currentTest string

const signingTestAttempts = 3
const initialSigningTestBackoff = 5 * time.Second

// retryingTransport lets a rate limited signing service recover mid-suite, since the CLI itself fails on the first error.
// Only rate limiting and server errors are retried because any other status will not change on a retry.
type retryingTransport struct {
	next http.RoundTripper
}

func (transport retryingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	var backoff = initialSigningTestBackoff
	for attempt := 1; ; attempt++ {
		var attemptRequest = request.Clone(request.Context())
		if request.GetBody != nil {
			var body, bodyErr = request.GetBody()
			if bodyErr != nil {
				return nil, bodyErr
			}
			attemptRequest.Body = body
		}

		var response, responseErr = transport.next.RoundTrip(attemptRequest)
		if responseErr != nil || attempt == signingTestAttempts {
			return response, responseErr
		}
		if response.StatusCode != http.StatusTooManyRequests && response.StatusCode < http.StatusInternalServerError {
			return response, nil
		}

		_ = response.Body.Close()
		fmt.Println(ansi(fmt.Sprintf("Failed to sign Shortcut (%s). Backing off from %s for %s (attempt %d of %d)...", response.Status, request.URL.Host, backoff, attempt, signingTestAttempts), red))
		time.Sleep(backoff)
		backoff *= 2
	}
}

func TestCherri(_ *testing.T) {
	var defaultTransport = http.DefaultTransport
	http.DefaultTransport = retryingTransport{next: defaultTransport}
	defer func() { http.DefaultTransport = defaultTransport }()

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

// roundTripKnownFailures lists shortcuts whose decompile or recompile exits non-zero, keyed by
// test file name without _unsigned.shortcut. An entry that starts passing must be removed.
var roundTripKnownFailures = map[string]string{
	"ask-prompts": "the prompt of an Ask in a function call argument is inside the dictionary of the call, and Cherri has no syntax for a prompt inside text, an array or a dictionary",
	"calendar":    "addWorldClock() and the other world clock actions build their city parameters in Go and have no decomp hook to read the city back",
	"functions":   "the compiler writes the @items += item of the functions header for an array argument as an Append to Variable of an output named Variable 1 that no action has",
	"menus":       "the compiler writes @items += makeVCard() as an Append to Variable of an output named Variable 1 that no action has",
	"standard":    "the compiler writes the @items += item of the functions header for an array argument as an Append to Variable of an output named Variable 1 that no action has",
}

// roundTripKnownMismatches lists shortcuts that recompile but with a different action sequence,
// keyed like roundTripKnownFailures. An entry that starts matching must be removed.
var roundTripKnownMismatches = map[string]string{
	"numbers": "an expression with an action call compiles to Calculate Expression, but the decompiled expression reads the call's output as a constant, and an expression of two operands compiles to Math",
}

// roundTripRun runs the built compiler and returns whether it exited zero. The compiler runs as a
// separate process because it keeps global state and calls os.Exit on error.
func roundTripRun(binary string, arguments ...string) (string, bool) {
	var output, runErr = exec.Command(binary, arguments...).CombinedOutput()
	return string(output), runErr == nil
}

// roundTripCheckDeterministic decompiles the shortcut again, in another process as the order of a map changes
// with each one, and fails when the code is not the same.
func roundTripCheckDeterministic(t *testing.T, binary string, shortcutPath string, dir string, firstPath string) {
	t.Helper()

	var secondDir = filepath.Join(dir, "dec-again")
	if mkdirErr := os.MkdirAll(secondDir, 0755); mkdirErr != nil {
		t.Fatal(mkdirErr)
	}
	var output, ok = roundTripRun(binary, "--import="+shortcutPath, "--output="+secondDir+"/", "--no-ansi", "--comments")
	if !ok {
		t.Fatalf("second decompile failed:\n%s", roundTripTail(output))
	}

	var first, firstErr = os.ReadFile(firstPath)
	if firstErr != nil {
		t.Fatal(firstErr)
	}
	var second, secondErr = os.ReadFile(filepath.Join(secondDir, filepath.Base(firstPath)))
	if secondErr != nil {
		t.Fatal(secondErr)
	}
	if string(first) != string(second) {
		t.Fatalf("decompiling twice gave different code:\n%s\n---\n%s", first, second)
	}
}

func roundTripTail(output string) string {
	var outputLines = strings.Split(strings.TrimSpace(output), "\n")
	if len(outputLines) > 20 {
		outputLines = outputLines[len(outputLines)-20:]
	}

	return strings.Join(outputLines, "\n")
}

func roundTripIdentifiers(t *testing.T, path string) []string {
	t.Helper()

	var data, readErr = os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("read %s: %s", path, readErr)
	}

	var parsed struct {
		WFWorkflowActions []struct {
			WFWorkflowActionIdentifier string
		}
	}
	var _, plistErr = plist.Unmarshal(data, &parsed)
	if plistErr != nil {
		t.Fatalf("parse %s: %s", path, plistErr)
	}

	var identifiers []string
	for _, action := range parsed.WFWorkflowActions {
		identifiers = append(identifiers, action.WFWorkflowActionIdentifier)
	}

	return identifiers
}

// roundTripDefinitionKeys are the keys of a shortcut that the definitions at the top of its code write.
var roundTripDefinitionKeys = []string{
	"WFWorkflowIcon",
	"WFWorkflowClientVersion",
	"WFQuickActionSurfaces",
	"WFWorkflowInputContentItemClasses",
	"WFWorkflowOutputContentItemClasses",
	"WFWorkflowTypes",
	"WFWorkflowNoInputBehavior",
	"WFWorkflowImportQuestions",
}

// roundTripDefinitions returns the values of roundTripDefinitionKeys as text, with the import questions in a
// fixed order, as the compiler writes them in the order of a map.
func roundTripDefinitions(t *testing.T, path string) map[string]string {
	t.Helper()

	var data, readErr = os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("read %s: %s", path, readErr)
	}
	var shortcutPlist map[string]interface{}
	var _, plistErr = plist.Unmarshal(data, &shortcutPlist)
	if plistErr != nil {
		t.Fatalf("parse %s: %s", path, plistErr)
	}

	var definitions = make(map[string]string)
	for _, key := range roundTripDefinitionKeys {
		var value = shortcutPlist[key]
		if questions, isQuestions := value.([]interface{}); isQuestions && key == "WFWorkflowImportQuestions" {
			var questionTexts []string
			for _, question := range questions {
				questionTexts = append(questionTexts, fmt.Sprint(question))
			}
			sort.Strings(questionTexts)
			value = questionTexts
		}
		// An icon color is a 32 bit value, which the compiler has a signed and an unsigned name for.
		if icon, isIcon := value.(map[string]interface{}); isIcon && key == "WFWorkflowIcon" {
			for iconKey, iconValue := range icon {
				switch number := iconValue.(type) {
				case int64:
					icon[iconKey] = uint32(number)
				case uint64:
					icon[iconKey] = uint32(number)
				}
			}
		}
		definitions[key] = fmt.Sprint(value)
	}

	return definitions
}

// firstDifference describes where two identifier lists diverge, or returns an empty string when equal.
func firstDifference(original []string, recompiled []string) string {
	for i := 0; i < max(len(original), len(recompiled)); i++ {
		var originalID = "(end)"
		if i < len(original) {
			originalID = original[i]
		}
		var recompiledID = "(end)"
		if i < len(recompiled) {
			recompiledID = recompiled[i]
		}
		if originalID != recompiledID {
			return fmt.Sprintf("action %d: %s became %s", i, originalID, recompiledID)
		}
	}

	return ""
}

// roundTripClientVersion gives a shortcut that has client version 900 the default client version. resetParser
// leaves it at 900 while the test shortcuts are compiled as version 27, so decompiling it writes a definition
// of version 16, which a recompile rejects the newer actions of.
func roundTripClientVersion(t *testing.T, data []byte) []byte {
	t.Helper()

	var shortcutPlist map[string]interface{}
	var _, unmarshalErr = plist.Unmarshal(data, &shortcutPlist)
	if unmarshalErr != nil {
		t.Fatal(unmarshalErr)
	}
	if shortcutPlist["WFWorkflowClientVersion"] != "900" {
		return data
	}

	shortcutPlist["WFWorkflowClientVersion"] = defaultClientVersion
	var changed, marshalErr = plist.Marshal(shortcutPlist, plist.XMLFormat)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}

	return changed
}

// TestRoundTrip decompiles each compiler-generated plist, recompiles the result and compares the
// action identifiers with the original. The _unsigned.shortcut files are produced by TestCherriNoSign.
//
// Run independently:
//
//	go test -run '^TestCherriNoSign$' && go test -run '^TestRoundTrip$'
func TestRoundTrip(t *testing.T) {
	var originals, globErr = filepath.Glob("tests/*_unsigned.shortcut")
	if globErr != nil {
		t.Fatal(globErr)
	}
	if len(originals) == 0 {
		t.Skip("no tests/*_unsigned.shortcut files - run TestCherriNoSign first")
	}

	var tempDir = t.TempDir()
	var binary = filepath.Join(tempDir, "cherri")
	var buildOutput, buildErr = exec.Command("go", "build", "-o", binary, ".").CombinedOutput()
	if buildErr != nil {
		t.Fatalf("build compiler: %s\n%s", buildErr, buildOutput)
	}

	for _, original := range originals {
		var name = strings.TrimSuffix(filepath.Base(original), unsignedEnd)
		t.Run(name, func(t *testing.T) {
			// Each subtest works in its own directory under tempDir and runs the compiler as separate
			// processes, so they share nothing that is written. tempDir is removed after all of them finish.
			t.Parallel()

			var dir = filepath.Join(tempDir, name)
			var decompiledDir = filepath.Join(dir, "dec")
			var mkdirErr = os.MkdirAll(decompiledDir, 0755)
			if mkdirErr != nil {
				t.Fatal(mkdirErr)
			}

			var data, readErr = os.ReadFile(original)
			if readErr != nil {
				t.Fatal(readErr)
			}
			data = roundTripClientVersion(t, data)
			var copied = filepath.Join(dir, name+".shortcut")
			var writeErr = os.WriteFile(copied, data, 0600)
			if writeErr != nil {
				t.Fatal(writeErr)
			}

			var knownFailure, failureKnown = roundTripKnownFailures[name]
			var knownMismatch, mismatchKnown = roundTripKnownMismatches[name]

			var failure string
			var decompiledPath = filepath.Join(decompiledDir, strings.ReplaceAll(name, " ", "_")+".cherri")
			// Comments are included so the comment actions of the original are written as comment actions.
			var output, ok = roundTripRun(binary, "--import="+copied, "--output="+decompiledDir+"/", "--no-ansi", "--comments")
			if !ok {
				failure = "decompile failed:\n" + roundTripTail(output)
			}

			if failure == "" {
				roundTripCheckDeterministic(t, binary, copied, dir, decompiledPath)
			}

			if failure == "" {
				output, ok = roundTripRun(binary, decompiledPath, "--skip-sign", "--no-ansi")
				if !ok {
					failure = "recompile failed:\n" + roundTripTail(output)
				}
			}

			if failure != "" {
				if failureKnown {
					t.Logf("known failure (%s): %s", knownFailure, failure)
					return
				}
				t.Fatal(failure)
			}
			if failureKnown {
				t.Fatalf("%s now decompiles and recompiles: remove from roundTripKnownFailures", name)
			}

			// The recompiled shortcut is named by the decompiled #define name, so it is found by suffix.
			var recompiled, recompiledErr = filepath.Glob(filepath.Join(decompiledDir, "*"+unsignedEnd))
			if recompiledErr != nil || len(recompiled) != 1 {
				t.Fatalf("expected one recompiled shortcut in %s, found %v", decompiledDir, recompiled)
			}

			var originalDefinitions = roundTripDefinitions(t, copied)
			var recompiledDefinitions = roundTripDefinitions(t, recompiled[0])
			for _, key := range roundTripDefinitionKeys {
				if originalDefinitions[key] != recompiledDefinitions[key] {
					t.Errorf("recompiled %s differs from the original:\n%s\n---\n%s", key, originalDefinitions[key], recompiledDefinitions[key])
				}
			}

			var difference = firstDifference(roundTripIdentifiers(t, original), roundTripIdentifiers(t, recompiled[0]))
			if difference != "" {
				if mismatchKnown {
					t.Logf("known mismatch (%s): %s", knownMismatch, difference)
					return
				}
				t.Fatalf("recompiled actions differ from the original, %s", difference)
			}
			if mismatchKnown {
				t.Fatalf("%s now round trips: remove from roundTripKnownMismatches", name)
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

func TestCapitalizeEmptyString(t *testing.T) {
	if got := capitalize(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
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

func TestSanitizeIdentifierWhitespaceOnly(t *testing.T) {
	identifier := sanitizeIdentifier(" ")

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
	iosVersion = 27.0
	questions = map[string]*question{}
	hasShortcutInputVariables = false
	tabLevel = 0
	definedWorkflowTypes = []string{}
	inputs = []string{}
	outputs = []string{}
	noInput = map[string]any{}
	tokens = []token{}
	included = []string{}
	lineOrigins = nil
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
	repeatItemIndex = 1
	repeatIndexDepth = 1
}

// TestCompileSubprocess is the body for compileInSubprocess, since parserError exits the process.
func TestCompileSubprocess(t *testing.T) {
	var file = os.Getenv("CHERRI_SUBPROCESS_FILE")
	if file == "" {
		t.Skip("only runs as a subprocess of compileInSubprocess")
	}

	args.Args["no-ansi"] = ""
	args.Args["skip-sign"] = ""
	os.Args[1] = file
	loadStandardActions()
	main()
}

func compileInSubprocess(t *testing.T, file string) string {
	t.Helper()

	var cmd = exec.Command(os.Args[0], "-test.run=^TestCompileSubprocess$")
	cmd.Env = append(os.Environ(), "CHERRI_SUBPROCESS_FILE="+file)
	var output, _ = cmd.CombinedOutput()
	return string(output)
}

func TestNoAnsiIncludeErrorNamesTheIncludedFile(t *testing.T) {
	var output = compileInSubprocess(t, "tests/include-errors/main.cherri")

	var errorPrefix = "Error: Use raw string (') for include file paths (2:"
	var lines = strings.Split(output, "\n")
	var errorIdx = -1
	for i, line := range lines {
		if strings.HasPrefix(line, errorPrefix) {
			errorIdx = i
			break
		}
	}
	if errorIdx == -1 {
		t.Fatalf("missing error line %q in output:\n%s", errorPrefix, output)
	}
	if errorIdx+1 >= len(lines) || lines[errorIdx+1] != "File: tests/include-errors/shared/helper.cherri" {
		t.Fatalf("expected the File line directly after the error line, got output:\n%s", output)
	}
}

func TestNoAnsiMainFileErrorOmitsTheFile(t *testing.T) {
	var output = compileInSubprocess(t, "tests/include-errors/main-error.cherri")

	if !strings.Contains(output, "Error: Use raw string (') for include file paths (1:") {
		t.Fatalf("missing error line in output:\n%s", output)
	}
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "File:") {
			t.Fatalf("unexpected File line %q in output:\n%s", line, output)
		}
	}
}
