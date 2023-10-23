/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var currentTest string

func TestCherri(_ *testing.T) {
	var files, err = os.ReadDir("examples")
	if err != nil {
		fmt.Println(ansi("FAILED", red))
		panic(err)
	}
	for _, file := range files {
		if strings.Contains(file.Name(), ".cherri") {
			currentTest = "examples/" + file.Name()
			fmt.Printf("Compiling %s...\n", ansi(currentTest, bold))
			os.Args[1] = currentTest
			main()
			fmt.Print(ansi("PASSED", green) + "\n\n")

			resetParser()
		}
	}
}

func TestSingleFile(_ *testing.T) {
	var currentTest = "examples/menus.cherri"
	fmt.Printf("Compiling %s...\n", ansi(currentTest, bold))
	os.Args[1] = currentTest
	main()
	fmt.Print(ansi("PASSED", green) + "\n\n")
}

func TestActionList(_ *testing.T) {
	standardActions()
	for identifier := range actions {
		fmt.Println("{label: '" + identifier + "', type: 'function', detail: 'action'},")
	}
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
