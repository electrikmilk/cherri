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
	"strings"
	"testing"
	"time"

	"github.com/electrikmilk/args-parser"
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
	iosVersion = 27.2
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

func TestCRLFSourceAndIncludeCompile(t *testing.T) {
	var directory = t.TempDir()
	for _, name := range []string{"includes.cherri", "include.cherri", "definitions.cherri"} {
		var source, readErr = os.ReadFile(filepath.Join("tests", name))
		if readErr != nil {
			t.Fatal(readErr)
		}
		var crlfSource = strings.ReplaceAll(normalizeLineEndings(string(source)), "\n", "\r\n")
		if writeErr := os.WriteFile(filepath.Join(directory, name), []byte(crlfSource), 0600); writeErr != nil {
			t.Fatal(writeErr)
		}
	}

	var output = compileInSubprocess(t, filepath.Join(directory, "includes.cherri"))

	var compiled, globErr = filepath.Glob(filepath.Join(directory, "*"+unsignedEnd))
	if globErr != nil {
		t.Fatal(globErr)
	}
	if len(compiled) != 1 {
		t.Fatalf("expected one compiled Shortcut, got output:\n%s", output)
	}
}
