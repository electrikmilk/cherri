/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var filePath string
var filename string
var basename string
var contents string
var relativePath string
var outputPath string

var included []string

const fileExtension = "cherri"

func main() {
	registerArg("share", "s", "Signing mode. [anyone, contacts] [default=contacts]")
	registerArg("unsigned", "u", "Don't sign compiled Shortcut. Will NOT run on iOS or macOS.")
	registerArg("debug", "d", "Save generated plist. Print debug messages and stack traces.")
	registerArg("output", "o", "Optional output file path. (e.g. /path/to/file.shortcut).")
	registerArg("import", "i", "Opens compiled Shortcut after compilation. Ignored if unsigned.")
	registerArg("no-ansi", "a", "Don't output ANSI escape sequences that format and color the output.")
	if len(os.Args) <= 1 {
		usage()
		os.Exit(1)
	}
	filePath = os.Args[1]
	checkFile(filePath)
	var pathParts = strings.Split(filePath, "/")
	filename = end(pathParts)
	relativePath = strings.Replace(filePath, filename, "", 1)
	var nameParts = strings.Split(filename, ".")
	basename = nameParts[0]
	var bytes, readErr = os.ReadFile(filePath)
	handle(readErr)
	contents = string(bytes)

	outputPath = basename + ".shortcut"
	if arg("output") {
		outputPath = argValue("output")
	}

	if strings.Contains(contents, "#include") {
		parseIncludes()
	}

	if arg("debug") {
		fmt.Printf("Parsing %s... ", filename)
	}
	parse()
	if arg("debug") {
		fmt.Print(ansi("done!", green) + "\n")
	}

	if arg("debug") {
		fmt.Println(tokens)
		fmt.Print("\n")
		fmt.Println(variables)
		fmt.Print("\n")
		fmt.Println(menus)
		fmt.Print("\n")
	}

	if arg("debug") {
		fmt.Printf("Generating plist... ")
	}
	var plist = makePlist()
	if arg("debug") {
		fmt.Print(ansi("done!", green) + "\n")
	}

	if arg("debug") {
		fmt.Printf("Creating %s.plist... ", basename)
		plistWriteErr := os.WriteFile(basename+".plist", []byte(plist), 0600)
		handle(plistWriteErr)
		fmt.Print(ansi("done!", green) + "\n")
	}

	if arg("debug") {
		fmt.Printf("Creating unsigned %s.shortcut... ", basename)
	}
	shortcutWriteErr := os.WriteFile(basename+"_unsigned.shortcut", []byte(plist), 0600)
	handle(shortcutWriteErr)
	if arg("debug") {
		fmt.Print(ansi("done!", green) + "\n")
	}

	if !arg("unsigned") {
		sign()
	}

	if arg("import") && !arg("unsigned") {
		var _, importErr = exec.Command("open", outputPath).Output()
		handle(importErr)
	}
}

func parseIncludes() {
	lines = strings.Split(contents, "\n")
	for l, line := range lines {
		lineIdx = l
		if !strings.Contains(line, "#include") {
			continue
		}
		r := regexp.MustCompile("\"(.*?)\"")
		var includePath = strings.Trim(r.FindString(line), "\"")
		if includePath == "" {
			parserError("No path inside of include")
		}
		if !strings.Contains(includePath, "..") {
			includePath = relativePath + includePath
		}
		if contains(included, includePath) {
			parserError(fmt.Sprintf("File '%s' has already been included.", includePath))
		}
		checkFile(includePath)
		bytes, readErr := os.ReadFile(includePath)
		handle(readErr)
		lines[l] = string(bytes)
		included = append(included, includePath)
	}
	contents = strings.Join(lines, "\n")
	lineIdx = 0
	if strings.Contains(contents, "#include") {
		parseIncludes()
	}
}

func checkFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		exit(fmt.Sprintf("File at path '%s' does not exist!", filePath))
	}
	var file, statErr = os.Stat(filePath)
	handle(statErr)
	var nameParts = strings.Split(file.Name(), ".")
	var ext = end(nameParts)
	if ext != fileExtension {
		exit(fmt.Sprintf("File '%s' is not a .%s file!", filePath, fileExtension))
	}
}

func sign() {
	var signingMode = "people-who-know-me"
	if arg("share") {
		if argValue("share") == "anyone" {
			signingMode = "anyone"
		}
	}
	if arg("debug") {
		fmt.Printf("Signing %s.shortcut... ", basename)
	}
	var signBytes, signErr = exec.Command(
		"shortcuts",
		"sign",
		"-i", basename+"_unsigned.shortcut",
		"-o", outputPath,
		"-m", signingMode,
	).Output()
	if signErr != nil {
		if arg("debug") {
			fmt.Print(ansi("failed!", red) + "\n")
		}
		fmt.Println("\n" + ansi("Error: Failed to sign Shortcut, plist may be invalid!", red))
		if len(signBytes) > 0 {
			fmt.Println("shortcuts:", string(signBytes))
		}
		os.Exit(1)
	}
	if arg("debug") {
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

func shortcutsUUID() string {
	return strings.ToUpper(uuid.New().String())
}

type outputType int

const (
	bold      outputType = 1
	dim       outputType = 2
	underline outputType = 4
	red       outputType = 31
	green     outputType = 32
	yellow    outputType = 33
)

const CSI = "\033["

func ansi(message string, typeOf outputType) string {
	if arg("no-ansi") {
		return message
	}
	return fmt.Sprintf("%s%dm%s", CSI, typeOf, message) + "\033[0m"
}

func exit(message string) {
	fmt.Println("\nError: " + ansi(message, red) + "\n")
	if arg("debug") {
		panic("debug")
	} else {
		os.Exit(1)
	}
}
