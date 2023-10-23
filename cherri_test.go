/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"testing"
)

var currentTest string

func TestProfile(t *testing.T) {
	testFiles(func(file os.DirEntry) bool {
		compile()

		return true
	})
}

func TestCherri(t *testing.T) {
	testFiles(func(file os.DirEntry) bool {
		compile()

		return matchesExpected(t)
	})
}

func TestSingleFile(_ *testing.T) {
	currentTest = "examples/conditionals.cherri"
	fmt.Printf("⚙️ Compiling %s...\n", ansi(currentTest, bold))
	os.Args[1] = currentTest
	main()
	fmt.Print(ansi("✅  PASSED", green, bold) + "\n\n")
}

func TestActionList(_ *testing.T) {
	standardActions()
	for identifier := range actions {
		fmt.Println("{label: '" + identifier + "', type: 'function', detail: 'action'},")
	}
}

func testFiles(handle func(file os.DirEntry) bool) {
	var files, err = os.ReadDir("examples")
	if err != nil {
		fmt.Println(ansi("FAILED: unable to read examples directory", red))
		panic(err)
	}
	for _, file := range files {
		if !strings.Contains(file.Name(), ".cherri") {
			continue
		}
		currentTest = fmt.Sprintf("examples/%s", file.Name())

		fmt.Println(ansi(currentTest, underline, bold))

		if handle(file) {
			fmt.Println(ansi("✅  PASSED", green, bold))
		}

		fmt.Print("\n")

		resetParser()
	}
}

func compile() {
	os.Args[1] = currentTest
	defer func() {
		if recover() != nil {
			fmt.Println(ansi("‼️DID NOT COMPILE", bold, green))
		}
	}()

	main()

	fmt.Println(ansi("☑️ COMPILED", bold))
}

func matchesExpected(_ *testing.T) bool {
	var expectedPlist = fmt.Sprintf("examples/%s_expected.plist", basename)
	var _, statErr = os.Stat(expectedPlist)
	if os.IsNotExist(statErr) {
		fmt.Println(ansi("Test has no exported plist to compare against.", yellow))
		return true
	}
	var expectedBytes, readErr = os.ReadFile(expectedPlist)
	handle(readErr)

	var xmlErr error

	var compiled interface{}
	xmlErr = xml.Unmarshal([]byte(plist.String()), &compiled)
	handle(xmlErr)

	var expected interface{}
	xmlErr = xml.Unmarshal(expectedBytes, &expected)
	handle(xmlErr)

	if expected != compiled {
		fmt.Print(ansi("‼️ DOES NOT MATCH EXPECTED", red, bold) + "\n")
		return false
	}

	fmt.Println(ansi("☑️ MATCHES EXPECTED", bold))
	return true
}

func resetParser() {
	lines = []string{}
	chars = []rune{}
	char = -1
	idx = 0
	lineIdx = 0
	lineCharIdx = 0
	groupingUUIDs = map[int]string{}
	groupingTypes = map[int]tokenType{}
	groupingIdx = 0
	variables = map[string]variableValue{}
	actions = map[string]*actionDefinition{}
	iconColor = "-1263359489"
	iconGlyph = 61440
	minVersion = "900"
	iosVersion = 16.5
	questions = map[string]*question{}
	hasShortcutInputVariables = false
	tabLevel = 0
	types = []string{}
	inputs = []string{}
	outputs = []string{}
	globals = map[string]variableValue{}
	noInput = noInputParams{}
	tokens = []token{}
	included = []string{}
	includes = []include{}
	workflowName = ""
	plist.Reset()
	menus = map[string][]variableValue{}
	uuids = map[string]string{}
	customActions = map[string]customAction{}
}
