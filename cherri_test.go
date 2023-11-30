/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
)

var currentTest string

func TestCherri(t *testing.T) {
	makeActionsTest()
	runTests(func(file os.DirEntry) {
		compile()
	})
}

// Uses less resources than TestCherri
func TestProfile(t *testing.T) {
	runTests(func(file os.DirEntry) {
		compile()
	})
}

func TestSingleFile(_ *testing.T) {
	currentTest = "tests/conditionals.cherri"
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

func runTests(handle func(file os.DirEntry)) {
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

		handle(file)

		fmt.Println(ansi("✅  PASSED", green, bold))
		fmt.Print("\n")

		resetParser()
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

func makeActionsTest() {
	standardActions()
	var actionsTest strings.Builder
	actionsTest.WriteString("/*\nThis is a generated test of all the actions the compiler supports.\nDO NOT RUN THIS, it is a random list of actions that would likely do bad things.\n*/\n#define mac true\n@emptyVar = nil\n")
	for identifier, definition := range actions {
		actionsTest.WriteString(identifier)
		actionsTest.WriteRune('(')

		if definition.parameters != nil {
			var paramsSize = len(definition.parameters) - 1
			for i, param := range definition.parameters {
				var paramValue string
				switch param.validType {
				case String:
					paramValue = "Test"
					if param.enum != nil {
						paramValue = param.enum[0]
					}
					var appIDParams = []string{"appID", "except", "firstAppID", "secondAppID"}
					switch {
					case identifier == "makeVCard" && param.name == "imagePath":
						paramValue = "assets/cherri_icon.png"
					case identifier == "makeSizedDiskImage" && param.name == "size":
						var randInt = rand.Intn(10)
						if randInt < 1 {
							randInt = 1
						}
						paramValue = fmt.Sprintf("%d GB", randInt)
					case (identifier == "convertMeasurement" || identifier == "measurement") && param.name == "unit":
						paramValue = "g-force"
					case contains(appIDParams, param.name):
						paramValue = "shortcuts"
					case identifier == "rawAction":
						paramValue = "is.workflow.actions.alert"
					case identifier == "setVolume" || identifier == "setBrightness":
						paramValue = "10"
					case strings.Contains(strings.ToLower(param.key), "language"):
						paramValue = "Arabic"
					}
					paramValue = fmt.Sprintf("\"%s\"", paramValue)
				case Integer:
					var randInt = rand.Intn(10)
					if randInt == 0 {
						randInt = 1
					}
					paramValue = fmt.Sprintf("%d", randInt)
				case Bool:
					var randInt = rand.Intn(1)
					if randInt == 1 {
						paramValue = "true"
					} else {
						paramValue = "false"
					}
				case Var:
					paramValue = "emptyVar"
				case Arr:
					paramValue = "[]"

					if identifier == "rawAction" {
						paramValue = "[{\"key\":\"WFAlertActionMessage\",\"type\":\"string\",\"value\":\"Hello, world!\"},{\"key\":\"WFAlertActionTitle\",\"type\":\"string\",\"value\":\"Alert\"}]"
					}
				case Dict:
					paramValue = "{}"
				}
				if i < paramsSize {
					paramValue = fmt.Sprintf("%s, ", paramValue)
				}
				if param.infinite {
					var infiniteArgs = rand.Intn(5)
					if infiniteArgs > 1 {
						paramValue = fmt.Sprintf("%s, ", paramValue)
						paramValue = strings.Repeat(paramValue, infiniteArgs)
						paramValue = strings.Trim(paramValue, ", ")
					}
				}
				actionsTest.WriteString(paramValue)
			}
		}
		actionsTest.WriteString(")\n")
	}

	var writeErr = os.WriteFile("tests/actions.cherri", []byte(actionsTest.String()), 0600)
	handle(writeErr)
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
	usedActions = []string{}
	menus = map[string][]variableValue{}
	uuids = map[string]string{}
	customActions = map[string]customAction{}
}
