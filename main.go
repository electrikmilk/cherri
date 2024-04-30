/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	"github.com/electrikmilk/args-parser"
	"os"
	"runtime"
	"strings"
	"unicode"
)

var filePath string
var filename string
var basename string
var contents string
var relativePath string
var inputPath string
var outputPath string

const unsignedEnd = "_unsigned.shortcut"
const darwin = runtime.GOOS == "darwin"

func main() {
	if args.Using("help") {
		args.PrintUsage()
		os.Exit(0)
	}

	if args.Using("version") {
		printVersion()
		os.Exit(0)
	}

	if args.Using("action") {
		if args.Value("action") == "" {
			for identifier, definition := range actions {
				setCurrentAction(identifier, definition)
				fmt.Println(generateActionDefinition(parameterDefinition{}, true, true))
				fmt.Print("\n")
			}
		} else {
			printActionDefinitions()
		}
		os.Exit(0)
	}

	filePath = fileArg()
	if len(os.Args) == 1 || filePath == "" {
		printLogo()
		printVersion()
		if !darwin {
			fmt.Println(ansi("\nWarning:", yellow, bold), "Shortcuts compiled on this platform will not run on iOS 15+ or macOS 12+.")
		}
		fmt.Print("\n")
		args.PrintUsage()
		os.Exit(1)
	}

	filename = checkFile(filePath)

	handleFile()

	initParse()

	makePlist()

	createShortcut()
}

func handle(err error) {
	if err == nil {
		return
	}

	fmt.Print(ansi("\nProgram panic :(\n\n", red, bold))
	fmt.Print(ansi("Please report this: https://github.com/electrikmilk/cherri/issues/new\n\n", red))

	if args.Using("debug") {
		panicDebug(err)
	} else {
		panic(err)
	}
}

func printVersion() {
	var color outputType
	if strings.Contains(version, "beta") {
		color = yellow
	} else {
		color = green
	}
	fmt.Printf("Cherri Compiler " + ansi(version, color) + "\n")
}

func printLogo() {
	fmt.Print(ansi("\n                       $$$\n                     $$$$$\n                 $$$$$$$$$\n            $$$$$$$$  $$$$\n        $$$$      $$  $$$$\n      $$$     $$$$$$  $ $$\n     $    $$$$$$$$$  $  $$              \n   $$   $$$$$$$$$   $   $$              \n  $$  $$$$$$$$$$   $    $$              \n  $$$$$$$$$$$     $$     $$    ", green))
	fmt.Print(ansi("$$$$$$   \n", red))
	fmt.Print(ansi("$$$$$$$$$       $$       $$ ", green))
	fmt.Print(ansi("$$$$$$$$$$\n         $$$$$  ", red))
	fmt.Print(ansi("$$$    $$  $  ", green))
	fmt.Print(ansi("$$$$$$$$$\n      $$$$$$$$$ $$$$$$$$$$$$$$$$$$$$$$$$\n     $$$$$$$$$$$$  $$$$$  $$$$$$$$$$$$$$\n    $$$  $$$$$$$$$$$$$$$$  $$$$$$$$$  $$\n    $$   $$$$$$$$$$$$$$$$$ $$$$$$$$$  $$\n    $$  $$$$$$$$$$$$$$$$$$ $$$$$$$$$ $$ \n    $$$ $$$$$$$$$$$$$$$$$$ $$$$$$$$$$$  \n    $$$  $$$$$$$$$$$$$$$$ $$$$$$$$$$$   \n     $$$$$$$$$$$$$$$$$$$  $$$$$$$$$     \n     $$$$$$$$$$$$$$$$$$$                \n       $$$$$$$$$$$$$$$                  \n          $$$$$$$$$$                    \n\n", red))
}

func fileArg() string {
	for _, arg := range os.Args {
		if !strings.Contains(arg, ".cherri") || startsWith("-", arg) {
			continue
		}

		return arg
	}
	return ""
}

// handleFile splits the file argument into parts.
func handleFile() {
	relativePath = strings.Replace(filePath, filename, "", 1)
	var nameParts = strings.Split(filename, ".")
	basename = nameParts[0]
	workflowName = basename

	outputPath = relativePath + workflowName + ".shortcut"
	if args.Using("output") {
		outputPath = args.Value("output")
	}

	var fileBytes, readErr = os.ReadFile(filePath)
	handle(readErr)
	contents = string(fileBytes)
}

// checkFile checks if the file exists and is a .cherri file.
func checkFile(filePath string) (filename string) {
	var file, statErr = os.Stat(filePath)
	if os.IsNotExist(statErr) {
		exit(fmt.Sprintf("File '%s' does not exist!", filePath))
	}
	var nameParts = strings.Split(file.Name(), ".")
	var ext = end(nameParts)
	if ext != "cherri" {
		exit(fmt.Sprintf("File '%s' is not a .cherri file!", filePath))
	}
	return file.Name()
}

func end(slice []string) string {
	return slice[len(slice)-1]
}

func capitalize(s string) string {
	var char = s[0]
	var after, _ = strings.CutPrefix(s, string(char))
	return fmt.Sprintf("%c%s", unicode.ToUpper(rune(char)), after)
}

// startsWith determines if the beginning characters of `substr` match `s`.
func startsWith(s string, substr string) bool {
	var stringChars = []rune(s)
	var subStringChars = []rune(substr)
	for i, char := range subStringChars {
		if char != stringChars[i] {
			return false
		}
	}
	return true
}

func panicDebug(err error) {
	fmt.Println(ansi("###################\n#   DEBUG PANIC   #\n###################\n", bold, red))
	printParsingDebug()
	printPlistGenDebug()
	printCustomActionsDebug()
	printIncludesDebug()
	fmt.Println(ansi("#############################################################\n", bold, red))

	if err != nil {
		panic(err)
	}

	panic("debug")
}

func printActionDefinitions() {
	var identifier = args.Value("action")
	if _, found := actions[identifier]; !found {
		fmt.Println(ansi(fmt.Sprintf("\nAction %s() does not exist or has not yet been defined.", identifier), red))

		switch identifier {
		case "text":
			fmt.Print("\nText actions are abstracted into string statements. For example:\n\n@variable = \"Hello, Cherri!\"\n\n")
			os.Exit(1)
		case "dictionary":
			fmt.Print("\nDictionary actions are abstracted into JSON object statements. For example:\n\n@variable = {\"test\":5\", \"key\":\"value\"}\n\n")
			os.Exit(1)
		}

		var actionSearchResults strings.Builder
		for actionIdentifier, definition := range actions {
			if strings.Contains(strings.ToLower(actionIdentifier), identifier) {
				setCurrentAction(actionIdentifier, definition)
				var definition = generateActionDefinition(parameterDefinition{}, false, false)
				definition, _ = strings.CutPrefix(definition, actionIdentifier)

				var capitalized = capitalize(identifier)
				var lowercase = strings.ToLower(identifier)
				switch {
				case strings.Contains(actionIdentifier, identifier):
					identifier = strings.ReplaceAll(actionIdentifier, identifier, ansi(identifier, red))
				case strings.Contains(actionIdentifier, capitalized):
					identifier = strings.ReplaceAll(actionIdentifier, capitalized, ansi(capitalized, red))
				case strings.Contains(actionIdentifier, lowercase):
					identifier = strings.ReplaceAll(actionIdentifier, lowercase, ansi(lowercase, red))
				}
				actionSearchResults.WriteString(fmt.Sprintf("- %s%s\n", identifier, definition))
			}
		}
		if actionSearchResults.Len() > 0 {
			fmt.Println(ansi("\nThe closest actions are:", yellow, italic, bold))
			fmt.Println(actionSearchResults.String())
		}

		os.Exit(1)
	}
	setCurrentAction(identifier, actions[identifier])
	fmt.Println(generateActionDefinition(parameterDefinition{}, true, true))
}
