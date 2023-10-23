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
var outputPath string

var included []string

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
		printActionDefinitions()
		os.Exit(0)
	}

	filePath = fileArg()
	if len(os.Args) == 1 || filePath == "" {
		printLogo()
		printVersion()
		if runtime.GOOS != "darwin" {
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

	if args.Using("debug") {
		fmt.Print(ansi("done!", green) + "\n")
		fmt.Printf("Generating plist... ")
	}

	makePlist()

	if args.Using("debug") {
		fmt.Print(ansi("done!", green) + "\n")
	}

	createShortcut()

	if args.Using("debug") {
		printDebug()
	}
}

func printActionDefinitions() {
	standardActions()
	var action = args.Value("action")
	if action != "" {
		if _, found := actions[action]; !found {
			exit(fmt.Sprintf("Action %s() does not exist or has not yet been defined.", action))
		}
		currentAction = action
		fmt.Println(generateActionDefinition(parameterDefinition{}, true, true))
		os.Exit(0)
	}
	for action := range actions {
		currentAction = action
		fmt.Println(generateActionDefinition(parameterDefinition{}, true, true) + "\n")
	}
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
	if args.Using("debug") {
		writeFile(relativePath+workflowName+".plist", fmt.Sprintf("Creating %s.plist", workflowName))
	}
	writeFile(relativePath+workflowName+"_unsigned.shortcut", fmt.Sprintf("Creating unsigned %s.shortcut", workflowName))

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
	if args.Using("debug") {
		fmt.Print(debug + "... ")
	}
	writeErr := os.WriteFile(filename, []byte(plist), 0600)
	handle(writeErr)
	if args.Using("debug") {
		fmt.Print(ansi("done!", green) + "\n")
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
	if runtime.GOOS != "darwin" {
		fmt.Println(ansi("Warning:", bold, yellow), "macOS is required to sign shortcuts. The compiled Shortcut will not run on iOS 15+ or macOS 12+.")
		return
	}
	var signingMode = "people-who-know-me"
	if args.Using("share") && args.Value("share") == "anyone" {
		signingMode = "anyone"
	}
	if args.Using("debug") {
		fmt.Printf("Signing %s.shortcut... ", basename)
	}
	var sign = exec.Command(
		"shortcuts",
		"sign",
		"-i", relativePath+workflowName+"_unsigned.shortcut",
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
		fmt.Print(ansi("done!", green) + "\n")
	}

	removeErr := os.Remove(fmt.Sprintf("%s%s_unsigned.shortcut", relativePath, workflowName))
	handle(removeErr)

	if args.Using("import") {
		var _, importErr = exec.Command("open", outputPath).Output()
		handle(importErr)
	}
}

func end(slice []string) string {
	return slice[len(slice)-1]
}

func handle(err error) {
	if err != nil {
		var message = fmt.Sprintf("%s", err)
		fmt.Println("\n" + ansi("Error: "+message, red) + "\n")
		if args.Using("debug") {
			printDebug()
			panic(err)
		} else {
			os.Exit(1)
		}
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

func printPlistData(data []plistData) {
	var saveTabLevel = tabLevel
	tabLevel = 0
	for _, d := range data {
		fmt.Println(plistKeyValue(d.key, d.dataType, d.value))
	}
	tabLevel = saveTabLevel
}

func printDebug() {
	if args.Using("debug") {
		fmt.Println(ansi("#############\n#   DEBUG   #\n#############\n", red))

		fmt.Println(ansi("### PARSING ###", bold))

		if idx != 0 {
			fmt.Println("Previous Character:")
			printChar(prev(1))
		}

		fmt.Println("\nCurrent Character:")
		printChar(char)

		if len(contents) > idx+1 {
			fmt.Println("\nNext Character:")
			printChar(next(1))
		}

		if len(lines) > lineIdx {
			fmt.Println("\nCurrent Line: \n" + lines[lineIdx])
		}

		fmt.Print("\n")

		fmt.Println(ansi("### PARSING ###", bold) + "\n")

		fmt.Println(ansi("## TOKENS ##", bold))
		printTokens(tokens)
		fmt.Print("\n")

		fmt.Println(ansi("## TOKEN CHARS ##", bold))
		fmt.Println(tokenChars)
		fmt.Print("\n")

		fmt.Println(ansi("## DEFINITIONS ##", bold))
		fmt.Println("Name: " + workflowName)
		fmt.Println("Color: " + iconColor)
		fmt.Printf("Glyph: %d\n", iconGlyph)
		fmt.Printf("Inputs: %v\n", inputs)
		fmt.Printf("Outputs: %v\n", outputs)
		fmt.Printf("Workflows: %v\n", types)
		fmt.Printf("No Input: %v\n", noInput)
		fmt.Printf("macOS Only: %v\n", isMac)
		fmt.Printf("Mininum Version: %s\n", minVersion)
		fmt.Printf("iOS Version: %.1f\n", iosVersion)
		fmt.Print("\n")

		fmt.Println(ansi("## VARIABLES ##", bold))
		fmt.Println(variables)
		fmt.Print("\n")

		fmt.Println(ansi("## MENUS ##", bold))
		fmt.Println(menus)
		fmt.Print("\n")

		fmt.Println(ansi("## IMPORT QUESTIONS ##", bold))
		fmt.Println(questions)
		fmt.Print("\n")

		fmt.Println(ansi("## CUSTOM ACTIONS ##", bold))
		for identifier, customAction := range customActions {
			fmt.Println("identifier: " + identifier)
			fmt.Println("body:")
			fmt.Println(customAction.body)
		}
		fmt.Print("\n")

		fmt.Println(ansi("## INCLUDES ##", bold))
		fmt.Println(includes)
		fmt.Print("\n")

		fmt.Println(ansi("## UUIDS ##", bold))
		fmt.Println(uuids)
		fmt.Print("\n")
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

func splitContents() {
	contents = strings.Join(lines, "\n")
	lines = strings.Split(contents, "\n")
	chars = strings.Split(contents, "")
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
	fmt.Println(ansi("\nError: "+message, red) + "\n")
	if args.Using("debug") {
		printDebug()
		panic("debug")
	} else {
		os.Exit(1)
	}
}
