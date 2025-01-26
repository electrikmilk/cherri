/*
 * Copyright (c) Cherri
 */

package main

import (
	"encoding/json"
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
	var files, err = os.ReadDir("tests")
	if err != nil {
		fmt.Println(ansi("FAILED: unable to read tests directory", red))
		panic(err)
	}
	for _, file := range files {
		if !strings.Contains(file.Name(), ".cherri") {
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

func TestSingleFile(_ *testing.T) {
	currentTest = "tests/conditionals.cherri"
	fmt.Printf("⚙️ Compiling %s...\n", ansi(currentTest, bold))
	os.Args[1] = currentTest
	main()
	fmt.Print(ansi("✅  PASSED", green, bold) + "\n\n")
}

func TestActionList(_ *testing.T) {
	for identifier := range actions {
		fmt.Println("{label: '" + identifier + "', type: 'function', detail: 'action'},")
	}
}

func TestGlyphList(_ *testing.T) {
	var data, jsonErr = json.Marshal(glyphs)
	handle(jsonErr)

	fmt.Println(string(data))
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
	groupingUUIDs = map[int]string{}
	groupingTypes = map[int]tokenType{}
	groupingIdx = 0
	variables = map[string]variableValue{}
	iconColor = -1263359489
	iconGlyph = 61440
	clientVersion = "900"
	iosVersion = 16.5
	questions = map[string]*question{}
	hasShortcutInputVariables = false
	tabLevel = 0
	types = []string{}
	inputs = []string{}
	outputs = []string{}
	noInput = noInputParams{}
	tokens = []token{}
	included = []string{}
	includes = []include{}
	workflowName = ""
	plist.Reset()
	menus = map[string][]variableValue{}
	uuids = map[string]string{}
	customActions = map[string]*customAction{}
}
