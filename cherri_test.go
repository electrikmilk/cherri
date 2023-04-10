/*
 * Copyright (c) 2023 Brandon Jordan
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

func TestCherri(t *testing.T) {
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
			args.Args["output"] = file.Name() + "_unsigned.shortcut"
			resetParser()
			main()
			fmt.Print(ansi("PASSED", green) + "\n\n")
			err := os.Remove(args.Args["output"])
			if err != nil {
				fmt.Println(ansi(fmt.Sprintf("Failed to remove test file %s!\n", args.Args["output"]), red))
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
	menus = map[string][]variableValue{}
	uuids = map[string]string{}
	customActions = map[string]customAction{}
}
