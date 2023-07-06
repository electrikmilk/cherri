/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/electrikmilk/args-parser"
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
			args.Args["unsigned"] = ""
			resetParser()

			main()
			fmt.Print(ansi("PASSED", green) + "\n\n")

			var removePath = relativePath + basename + "_unsigned.shortcut"
			removeShortcutErr := os.Remove(removePath)
			if removeShortcutErr != nil {
				fmt.Println(ansi(fmt.Sprintf("Failed to remove test file %s!\n", removePath), red))
			}
		}
	}
}

func resetParser() {
	lines = []string{}
	chars = []string{}
	char = -1
	idx = 0
	lineIdx = 0
	lineCharIdx = 0
	closureUUIDs = map[int]string{}
	closureTypes = map[int]tokenType{}
	closureIdx = 0
	currentGroupingUUID = ""
	variables = map[string]variableValue{}
	actions = map[string]*actionDefinition{}
	iconColor = "-1263359489"
	iconGlyph = 61440
	inputs = []string{}
	outputs = []string{}
	globals = map[string]variableValue{}
	tokens = []token{}
	included = []string{}
	includes = []include{}
	workflowName = ""
	plist = ""
	menus = map[string][]variableValue{}
	uuids = map[string]string{}
	customActions = map[string]customAction{}
}
