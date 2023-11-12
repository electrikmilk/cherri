/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/electrikmilk/args-parser"
	"github.com/google/uuid"
)

var filePath string
var filename string
var basename string
var contents string
var relativePath string
var inputPath string
var outputPath string

var darwin bool

const unsignedEnd = "_unsigned.shortcut"

func main() {
	darwin = runtime.GOOS == "darwin"

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
			for action := range actions {
				currentAction = action
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

	handleIncludes()

	initParse()

	makePlist()

	createShortcut()
}

func printActionDefinitions() {
	standardActions()
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
		for actionIdentifier := range actions {
			if strings.Contains(strings.ToLower(actionIdentifier), identifier) {
				currentAction = identifier
				var definition = generateActionDefinition(parameterDefinition{}, false, false)
				definition, _ = strings.CutPrefix(definition, actionIdentifier)

				var capitalized = capitalize(identifier)
				var lowercase = strings.ToLower(identifier)
				switch {
				case strings.Contains(actionIdentifier, identifier):
					identifier = strings.ReplaceAll(actionIdentifier, identifier, ansi(identifier, red))
				case strings.Contains(identifier, capitalized):
					identifier = strings.ReplaceAll(actionIdentifier, capitalized, ansi(capitalized, red))
				case strings.Contains(identifier, lowercase):
					identifier = strings.ReplaceAll(actionIdentifier, lowercase, ansi(lowercase, red))
				}
				actionSearchResults.WriteString(fmt.Sprintf("- %s%s\n", identifier, definition))
			}
		}
		if actionSearchResults.Len() > 0 {
			fmt.Println(ansi("\nThe closest identifier(s) are:", yellow, italic, bold))
			fmt.Println(actionSearchResults.String())
		}

		os.Exit(1)
	}
	currentAction = identifier
	fmt.Println(generateActionDefinition(parameterDefinition{}, true, true))
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

// createShortcut writes the Shortcut files to disk and signs them if the unsigned argument is not unused.
func createShortcut() {
	var path = fmt.Sprintf("%s%s", relativePath, workflowName)
	if args.Using("debug") {
		writeFile(path+".plist", workflowName+".plist")
	}
	writeFile(path+unsignedEnd, workflowName+unsignedEnd)

	inputPath = fmt.Sprintf("%s%s%s", relativePath, workflowName, unsignedEnd)

	sign()
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

// writeFile writes plist in bytes to filename.
func writeFile(filename string, debug string) {
	var writeDebugOutput = args.Using("debug")
	if writeDebugOutput {
		fmt.Println("Writing to " + debug + "...")
	}

	writeErr := os.WriteFile(filename, []byte(compiled), 0600)
	handle(writeErr)

	if writeDebugOutput {
		fmt.Println(ansi("Done.", green) + "\n")
	}
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

// sign runs the shortcuts sign command on the unsigned shortcut file.
func sign() {
	if !darwin {
		fmt.Println(ansi("Warning:", bold, yellow), "macOS is required to sign shortcuts. The compiled Shortcut will not run on iOS 15+ or macOS 12+.")
		return
	}

	var signingMode = "people-who-know-me"
	if args.Using("share") && args.Value("share") == "anyone" {
		signingMode = "anyone"
	}

	if args.Using("debug") {
		fmt.Printf("Signing %s to %s...\n", inputPath, outputPath)
	}
	var sign = exec.Command(
		"shortcuts",
		"sign",
		"-i", inputPath,
		"-o", outputPath,
		"-m", signingMode,
	)
	var stdErr bytes.Buffer
	sign.Stderr = &stdErr
	var signErr = sign.Run()
	if signErr != nil {
		if args.Using("debug") {
			fmt.Print(ansi("failed!", red) + "\n")
		}
		exit("Failed to sign Shortcut\n\nshortcuts: " + stdErr.String())
	}
	if args.Using("debug") {
		fmt.Println(ansi("Done.", green) + "\n")
	}

	removeUnsigned()

	if args.Using("import") {
		openShortcut()
	}
}

func removeUnsigned() {
	if args.Using("debug") {
		fmt.Println("Removing " + workflowName + "_unsigned.shortcut...")
	}

	removeErr := os.Remove(inputPath)
	handle(removeErr)

	if args.Using("debug") {
		fmt.Println(ansi("Done.", green))
	}
}

func openShortcut() {
	var _, importErr = exec.Command("open", outputPath).Output()
	handle(importErr)
}

func end(slice []string) string {
	return slice[len(slice)-1]
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func capitalize(s string) string {
	var chars = strings.Split(s, "")
	chars[0] = strings.ToUpper(chars[0])
	return strings.Join(chars, "")
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

func shortcutsUUID() string {
	return strings.ToUpper(uuid.New().String())
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

func splitContents() {
	contents = strings.Join(lines, "\n")
	lines = strings.Split(contents, "\n")
	chars = []rune(contents)
}

func printChar(ch rune) {
	var currentChar string
	switch ch {
	case ' ':
		currentChar = "SPACE"
	case '\t':
		currentChar = "TAB"
	case '\n':
		currentChar = "LF"
	case -1:
		currentChar = "EMPTY"
	default:
		currentChar = string(ch)
	}
	fmt.Printf("%s %d:%d\n", currentChar, lineIdx+1, lineCharIdx)
}

type outputType int

const (
	bold      outputType = 1
	italic    outputType = 3
	underline outputType = 4
	red       outputType = 31
	green     outputType = 32
	yellow    outputType = 33
)

const CSI = "\033["

func ansi(message string, typeOf ...outputType) string {
	if args.Using("no-ansi") {
		return message
	}
	var formattedMessage string
	for _, t := range typeOf {
		formattedMessage += fmt.Sprintf("%s%dm", CSI, t)
	}
	formattedMessage += message + "\033[0m"
	return formattedMessage
}

func exit(message string) {
	fmt.Println(ansi("\nError: "+message+"\n", red))
	if args.Using("debug") {
		panicDebug(nil)
	} else {
		os.Exit(1)
	}
}
