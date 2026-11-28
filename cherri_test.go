/*
 * Copyright (c) Cherri
 */

package main

import (
	"context"
	"errors"
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
	if os.Getenv("CHERRI_SUBPROCESS_DEBUG") != "" {
		args.Args["debug"] = ""
	}
	os.Args[1] = file
	loadStandardActions()
	main()
}

// runSubprocess runs a single test of this binary in its own process, for code that exits the process or reads stdin.
func runSubprocess(t *testing.T, subprocessTest string, directory string, stdin string, environment ...string) (output string, exitCode int) {
	t.Helper()

	var executable, executableErr = os.Executable()
	if executableErr != nil {
		t.Fatal(executableErr)
	}

	var homeSupplied bool
	for _, variable := range environment {
		if strings.HasPrefix(variable, "HOME=") {
			homeSupplied = true
		}
	}
	if !homeSupplied {
		environment = append(environment, "HOME="+t.TempDir())
	}

	// yesNo loops forever when stdin reaches EOF, so a prompt that is not answered must not hang the suite.
	var ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var cmd = exec.CommandContext(ctx, executable, "-test.run=^"+subprocessTest+"$")
	cmd.Dir = directory
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Env = append(os.Environ(), environment...)
	cmd.WaitDelay = 5 * time.Second

	var combined, runErr = cmd.CombinedOutput()
	if ctx.Err() != nil {
		t.Errorf("subprocess %s timed out", subprocessTest)
		return string(combined), -1
	}
	if runErr != nil {
		var exitErr *exec.ExitError
		if !errors.As(runErr, &exitErr) {
			t.Fatal(runErr)
		}
	}
	return string(combined), cmd.ProcessState.ExitCode()
}

func compileInSubprocess(t *testing.T, file string) string {
	t.Helper()

	var output, _ = runSubprocess(t, "TestCompileSubprocess", "", "", "CHERRI_SUBPROCESS_FILE="+file)
	return output
}

// TestInstallSubprocess is the body for TestInstallDeclinedTrustInstallsNothing, since exit ends the process.
func TestInstallSubprocess(t *testing.T) {
	var signature = os.Getenv("CHERRI_SUBPROCESS_INSTALL")
	if signature == "" {
		t.Skip("only runs as a subprocess of runSubprocess")
	}

	args.Args["no-ansi"] = ""
	args.Args["install"] = signature
	addPackage()
}

func TestInstallDeclinedTrustInstallsNothing(t *testing.T) {
	var projectDir = t.TempDir()
	var home = t.TempDir()

	var manifest = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
<dict>
	<key>Name</key>
	<string>project</string>
	<key>User</key>
	<string>me</string>
	<key>Packages</key>
	<array/>
</dict>
</plist>
`)
	var manifestPath = filepath.Join(projectDir, "info.plist")
	if writeErr := os.WriteFile(manifestPath, manifest, 0644); writeErr != nil {
		t.Fatal(writeErr)
	}

	var output, exitCode = runSubprocess(t, "TestInstallSubprocess", projectDir, "n\n", "CHERRI_SUBPROCESS_INSTALL=@u/a", "HOME="+home)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(output, "Package @u/a was not trusted and was not installed.") {
		t.Errorf("missing not trusted message in output:\n%s", output)
	}
	if _, statErr := os.Stat(filepath.Join(projectDir, "packages")); !os.IsNotExist(statErr) {
		t.Errorf("expected no packages directory, stat error: %v", statErr)
	}
	var after, readErr = os.ReadFile(manifestPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(after) != string(manifest) {
		t.Errorf("info.plist changed:\n%s", after)
	}
	if _, statErr := os.Stat(filepath.Join(home, ".cherri", "trusted.plist")); !os.IsNotExist(statErr) {
		t.Errorf("expected no trusted.plist, stat error: %v", statErr)
	}
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

// requireErrorWithFile fails unless output has a line starting with errorPrefix directly followed by the File line for file.
func requireErrorWithFile(t *testing.T, output string, errorPrefix string, file string) {
	t.Helper()

	var outputLines = strings.Split(output, "\n")
	for i, line := range outputLines {
		if !strings.HasPrefix(line, errorPrefix) {
			continue
		}
		if i+1 < len(outputLines) && outputLines[i+1] == "File: "+file {
			return
		}
		t.Fatalf("expected the File line for %s directly after the error line, got output:\n%s", file, output)
	}
	t.Fatalf("missing error line %q in output:\n%s", errorPrefix, output)
}

func requireNoFileLine(t *testing.T, output string) {
	t.Helper()

	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "File:") {
			t.Fatalf("unexpected File line %q in output:\n%s", line, output)
		}
	}
}

func TestIncludeCycleIsReportedAsAlreadyIncluded(t *testing.T) {
	var output = compileInSubprocess(t, "tests/include-errors/cycle/main.cherri")

	requireErrorWithFile(t, output,
		"Error: Path 'tests/include-errors/cycle/a.cherri' has already been included. (1:",
		"tests/include-errors/cycle/b.cherri")
}

func TestIncludeMainFileCycleIsReportedAsAlreadyIncluded(t *testing.T) {
	var output = compileInSubprocess(t, "tests/include-errors/main-cycle/main.cherri")

	requireErrorWithFile(t, output,
		"Error: Path 'tests/include-errors/main-cycle/main.cherri' has already been included. (1:",
		"tests/include-errors/main-cycle/nested/a.cherri")
}

func TestIncludeDuplicateNamesTheIncludingFile(t *testing.T) {
	var output = compileInSubprocess(t, "tests/include-errors/duplicate/main.cherri")

	if !strings.Contains(output, "Error: Path 'tests/include-errors/duplicate/x.cherri' has already been included. (2:") {
		t.Fatalf("missing duplicate error in output:\n%s", output)
	}
	if strings.Contains(output, "Include chain") {
		t.Fatalf("unexpected include chain in output:\n%s", output)
	}
	requireNoFileLine(t, output)
}

func TestIncludeErrorInDeepFileNamesThatFile(t *testing.T) {
	var output = compileInSubprocess(t, "tests/include-errors/deep/main.cherri")

	requireErrorWithFile(t, output,
		"Error: Use raw string (') for include file paths (2:",
		"tests/include-errors/deep/mid/two/three.cherri")
}

func TestIncludeErrorInMainAfterNestedScanIsInMain(t *testing.T) {
	var output = compileInSubprocess(t, "tests/include-errors/after-nested/main.cherri")

	if !strings.Contains(output, "Error: Use raw string (') for include file paths (2:") {
		t.Fatalf("missing error line in output:\n%s", output)
	}
	requireNoFileLine(t, output)
}

func TestIncludeMissingFileNamesTheInclude(t *testing.T) {
	var output = compileInSubprocess(t, "tests/include-errors/missing/main.cherri")

	requireErrorWithFile(t, output,
		"Error: File 'tests/include-errors/missing/nested/nothere.cherri' does not exist! (2:",
		"tests/include-errors/missing/nested/a.cherri")
}

func TestIncludeProjectDirectoryNamedActionsIsNotStdlib(t *testing.T) {
	var output, exitCode, _ = compilePackageFixture(t, "include-errors/actions-dir", false)

	requireCompiledCleanly(t, output, exitCode)
}

func TestPackageIncludeOutsidePackageIsError(t *testing.T) {
	var output = compileInSubprocess(t, "tests/include-isolation/escape/main.cherri")

	requireErrorWithFile(t, output,
		"Error: Include '../../escape.cherri' resolves to 'tests/include-isolation/escape/packages/escape.cherri', outside package '@u/a'. Packages are isolated and cannot include files outside their own directory. (1:",
		"tests/include-isolation/escape/packages/@u/a/main.cherri")
}

func TestPackageIncludeIntoSiblingPackageNamePrefixIsError(t *testing.T) {
	var output = compileInSubprocess(t, "tests/include-isolation/sibling/main.cherri")

	requireErrorWithFile(t, output,
		"Error: Include '../ab/x.cherri' resolves to 'tests/include-isolation/sibling/packages/@u/ab/x.cherri', outside package '@u/a'. Packages are isolated and cannot include files outside their own directory. (1:",
		"tests/include-isolation/sibling/packages/@u/a/main.cherri")
}

func TestPackageIncludeInsidePackageCompiles(t *testing.T) {
	var output, exitCode, _ = compilePackageFixture(t, "include-isolation/inside", false)

	requireCompiledCleanly(t, output, exitCode)
}

func TestPackageIncludeOfStdlibCategoryIsNotStdlib(t *testing.T) {
	var output = compileInSubprocess(t, "tests/include-isolation/stdlib/main.cherri")

	requireErrorWithFile(t, output,
		"Error: File 'tests/include-isolation/stdlib/packages/@u/a/actions/network' does not exist! (1:",
		"tests/include-isolation/stdlib/packages/@u/a/main.cherri")
}

// compilePackageFixture compiles the project in tests/<fixture> with that directory as the working directory, and removes
// the files the compile leaves behind.
func compilePackageFixture(t *testing.T, fixture string, debug bool) (output string, exitCode int, processedSource string) {
	t.Helper()

	var fixtureDirectory, absErr = filepath.Abs(filepath.Join("tests", fixture))
	if absErr != nil {
		t.Fatal(absErr)
	}

	var processedPath = filepath.Join(fixtureDirectory, "main_processed.cherri")
	t.Cleanup(func() {
		_ = os.Remove(processedPath)
		_ = os.Remove(filepath.Join(fixtureDirectory, "main.plist"))
		_ = os.Remove(filepath.Join(fixtureDirectory, "main_unsigned.shortcut"))
	})

	var environment = []string{"CHERRI_SUBPROCESS_FILE=main.cherri"}
	if debug {
		environment = append(environment, "CHERRI_SUBPROCESS_DEBUG=1")
	}
	output, exitCode = runSubprocess(t, "TestCompileSubprocess", fixtureDirectory, "", environment...)
	if !debug {
		return
	}

	var processed, readErr = os.ReadFile(processedPath)
	if readErr != nil {
		t.Fatalf("expected the processed source to be written: %s\n%s", readErr, output)
	}
	return output, exitCode, string(processed)
}

func requireCompiledCleanly(t *testing.T, output string, exitCode int) {
	t.Helper()

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d\n%s", exitCode, output)
	}
	// The debug output lists enum names such as QRCodeErrorCorrection, so match the compiler's own error prefix.
	if strings.Contains(output, "Error: ") {
		t.Errorf("unexpected error in output:\n%s", output)
	}
}

func TestPackageDependenciesAreIncludedTransitively(t *testing.T) {
	var output, exitCode, _ = compilePackageFixture(t, "package-deps", false)

	requireCompiledCleanly(t, output, exitCode)
}

func TestPackageDependenciesComeBeforeDependents(t *testing.T) {
	var output, exitCode, processed = compilePackageFixture(t, "package-order", true)

	requireCompiledCleanly(t, output, exitCode)

	var unrelated = strings.Index(processed, "// marker:@w/c")
	var dependency = strings.Index(processed, "// marker:@u/a")
	var dependent = strings.Index(processed, "// marker:@v/b")
	if unrelated == -1 || dependency == -1 || dependent == -1 {
		t.Fatalf("expected all three package markers in the processed source:\n%s", processed)
	}
	if !(unrelated < dependency && dependency < dependent) {
		t.Errorf("expected @w/c, then @u/a, then @v/b, got offsets %d, %d, %d", unrelated, dependency, dependent)
	}
}

func TestPackageDependencyCycleCompiles(t *testing.T) {
	var output, exitCode, processed = compilePackageFixture(t, "package-cycle", true)

	requireCompiledCleanly(t, output, exitCode)

	for _, marker := range []string{"// marker:@u/a", "// marker:@v/b"} {
		if count := strings.Count(processed, marker); count != 1 {
			t.Errorf("expected %q once in the processed source, found %d", marker, count)
		}
	}
}

func TestPackageDependencyWithPathTraversalIsRejected(t *testing.T) {
	var output, exitCode, _ = compilePackageFixture(t, "package-traversal", false)

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d\n%s", exitCode, output)
	}
	for _, expected := range []string{"x/../../..", "evil", "packages/@u/a/info.plist", "cannot be . or .."} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected %q in output:\n%s", expected, output)
		}
	}

	// The entry would have been cloned to the path it cleans to, which is next to the fixture.
	var escapedPath = filepath.Join("tests", "evil")
	if _, statErr := os.Stat(escapedPath); !os.IsNotExist(statErr) {
		t.Errorf("expected nothing at %s, stat error: %v", escapedPath, statErr)
		_ = os.RemoveAll(escapedPath)
	}
	var installedPackageFile = filepath.Join("tests", "package-traversal", "packages", "@u", "a", "main.cherri")
	if _, statErr := os.Stat(installedPackageFile); statErr != nil {
		t.Errorf("expected the fixture package to be left alone: %v", statErr)
	}
}
