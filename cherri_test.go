/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var currentTest string

func TestCherri(t *testing.T) {
	var files, err = os.ReadDir("examples")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if strings.Contains(file.Name(), ".cherri") {
			currentTest = "examples/" + file.Name()
			fmt.Printf("\nCompiling \033[1m%s\033[0m...\n", currentTest)
			os.Args[1] = currentTest
			args["unsigned"] = ""
			resetParser()
			main()
			fmt.Println("\033[32mPASSED\033[0m")
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
	actions = map[string]actionDefinition{}
	iconColor = "-1263359489"
	iconGlyph = 61440
	inputs = []string{}
	outputs = []string{}
	globals = map[string]variableValue{}
	tokens = []token{}
	included = []string{}
	workflowName = ""
	menus = map[string][]variableValue{}
	uuids = map[string]string{}
}
