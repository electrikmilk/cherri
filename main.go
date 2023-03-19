/*
 * Copyright (c) 2023 Brandon Jordan
 */

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
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
	args.Register("share", "s", "Signing mode. [anyone, contacts] [default=contacts]", true)
	args.Register("unsigned", "u", "Don't sign compiled Shortcut. Will NOT run on iOS or macOS.", false)
	args.Register("debug", "d", "Save generated plist. Print debug messages and stack traces.", false)
	args.Register("output", "o", "Optional output file path. (e.g. /path/to/file.shortcut).", true)
	args.Register("import", "i", "Opens compiled Shortcut after compilation. Ignored if unsigned.", false)
	args.Register("no-ansi", "", "Don't output ANSI escape sequences that format and color the output.", false)
	args.Register("auto-inc", "a", "Automatically include Cherri files in this directory.", false)
	args.CustomUsage = "[FILE]"

	if len(os.Args) <= 1 {
		args.PrintUsage()
	}

	filePath = os.Args[1]
	checkFile(filePath)

	var stat, statErr = os.Stat(filePath)
	handle(statErr)
	filename = stat.Name()

	relativePath = strings.Replace(filePath, filename, "", 1)
	var nameParts = strings.Split(filename, ".")
	basename = nameParts[0]

	var fileBytes, readErr = os.ReadFile(filePath)
	handle(readErr)
	contents = string(fileBytes)

	outputPath = basename + ".shortcut"
	if args.Using("output") {
		outputPath = args.Value("output")
	}

	if args.Using("auto-inc") {
		autoInclude()
	}

	if strings.Contains(contents, "#include") {
		parseIncludes()
	}

	actions = make(map[string]*actionDefinition)

	if strings.Contains(contents, "action") {
		standardActions()
		parseCustomActions()
	}

	if args.Using("debug") {
		fmt.Printf("Parsing %s... ", filename)
	}

	parse()

	if args.Using("debug") {
		fmt.Print(ansi("done!", green) + "\n")
		fmt.Printf("Generating plist... ")
	}

	var plist = makePlist()
	if args.Using("debug") {
		fmt.Print(ansi("done!", green) + "\n")
	}

	if args.Using("debug") {
		fmt.Printf("Creating %s.plist... ", basename)
		plistWriteErr := os.WriteFile(basename+".plist", []byte(plist), 0600)
		handle(plistWriteErr)
		fmt.Print(ansi("done!", green) + "\n")

		fmt.Printf("Creating unsigned %s.shortcut... ", basename)
	}

	var unsignedPath = basename + "_unsigned.shortcut"
	if args.Using("unsigned") {
		unsignedPath = outputPath
	}
	shortcutWriteErr := os.WriteFile(unsignedPath, []byte(plist), 0600)
	handle(shortcutWriteErr)
	if args.Using("debug") {
		fmt.Print(ansi("done!", green) + "\n")
	}

	if !args.Using("unsigned") {
		sign()
	}

	if args.Using("debug") {
		printDebug()
	}

	if args.Using("import") && !args.Using("unsigned") {
		var _, importErr = exec.Command("open", outputPath).Output()
		handle(importErr)
	}
}

func checkFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		exit(fmt.Sprintf("File '%s' does not exist!", filePath))
	}
	var file, statErr = os.Stat(filePath)
	handle(statErr)
	var nameParts = strings.Split(file.Name(), ".")
	var ext = end(nameParts)
	if ext != "cherri" {
		exit(fmt.Sprintf("File '%s' is not a .cherri file!", filePath))
	}
}

func sign() {
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
		"-i", basename+"_unsigned.shortcut",
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
	removeErr := os.Remove(basename + "_unsigned.shortcut")
	handle(removeErr)
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

// startsWith determines if the beginning characters of `substr` match `s`.
func startsWith(s string, substr string) bool {
	var stringChars = strings.Split(s, "")
	var subStringChars = strings.Split(substr, "")
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

		if len(contents) < idx {
			fmt.Println("\nNext Character:")
			printChar(next(1))
		}

		fmt.Println("\nCurrent Line: \n" + lines[lineIdx])
		fmt.Print("\n")

		fmt.Println(ansi("### TOKENS ###", bold))
		fmt.Println(tokens)
		fmt.Print("\n")

		fmt.Println(ansi("### VARIABLES ###", bold))
		fmt.Println(variables)
		fmt.Print("\n")

		fmt.Println(ansi("### MENUS ###", bold))
		fmt.Println(menus)
		fmt.Print("\n")

		fmt.Println(ansi("### IMPORT QUESTIONS ###", bold))
		fmt.Println(questions)
		fmt.Print("\n")

		fmt.Println(ansi("### CUSTOM ACTIONS ###", bold))
		for identifier, customAction := range customActions {
			fmt.Println("identifier: " + identifier)
			fmt.Println("arguments:", customAction.arguments)
			fmt.Println("body:")
			fmt.Println(customAction.body)
		}
		fmt.Print("\n")

		fmt.Println(ansi("### UUIDS ###", bold))
		fmt.Println(uuids)
		fmt.Print("\n")

		fmt.Println(ansi("### INCLUDES ###", bold))
		fmt.Println(includes)
		fmt.Print("\n")
	}
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

func ansi(message string, typeOf outputType) string {
	if args.Using("no-ansi") {
		return message
	}
	return fmt.Sprintf("%s%dm%s", CSI, typeOf, message) + "\033[0m"
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
