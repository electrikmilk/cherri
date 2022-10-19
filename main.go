/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

var filePath string
var filename string
var basename string
var contents string

var fileExtension = "cherri"

func main() {
	customUsage = "[FILE] [options]"
	registerArg("share", "s", "Signing mode. [anyone, contacts] [default=contacts]")
	registerArg("bypass", "b", "Bypass macOS check and signing. Resulting shortcut will NOT run on iOS or macOS.")
	registerArg("debug", "d", "Save generated plist. Print debug messages and stack traces.")
	if !arg("bypass") {
		if runtime.GOOS != "darwin" {
			fmt.Println("\033[31m\033[1mNot on macOS!\u001B[0m")
			fmt.Println("\u001B[31mShortcuts can only be signed on macOS!\033[0m")
			os.Exit(1)
		}
	}
	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}
	filePath = os.Args[1]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("\033[31mFile '%s' does not exist!", filePath)
		os.Exit(1)
	}
	var file, statErr = os.Stat(filePath)
	handle(statErr)
	var nameParts = strings.Split(file.Name(), ".")
	var ext = nameParts[len(nameParts)-1]
	if ext != fileExtension {
		fmt.Printf("\033[31mFile '%s' is not a .%s file!", filePath, fileExtension)
		os.Exit(1)
	}
	filename = file.Name()
	basename = nameParts[0]
	var bytes, readErr = os.ReadFile(filePath)
	handle(readErr)
	contents = string(bytes)

	if arg("debug") {
		fmt.Printf("Parsing %s... ", filename)
	}
	parse()
	if arg("debug") {
		fmt.Print("\033[32mdone!\033[0m\n")
	}

	if arg("debug") {
		fmt.Println(tokenSets)
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
		fmt.Print("\033[32mdone!\033[0m\n")
	}

	if arg("debug") {
		fmt.Printf("Creating %s.plist... ", basename)
		plistWriteErr := os.WriteFile(basename+".plist", []byte(plist), 0600)
		handle(plistWriteErr)
		fmt.Print("\033[32mdone!\033[0m\n")
	}

	if arg("debug") {
		fmt.Printf("Creating unsigned %s.shortcut... ", basename)
	}
	shortcutWriteErr := os.WriteFile(basename+"_unsigned.shortcut", []byte(plist), 0600)
	handle(shortcutWriteErr)
	if arg("debug") {
		fmt.Print("\033[32mdone!\033[0m\n")
	}

	if !arg("bypass") {
		sign()
	}

	if !arg("bypass") {
		removeErr := os.Remove(basename + "_unsigned.shortcut")
		handle(removeErr)
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
		"-o", basename+".shortcut",
		"-m", signingMode,
	).Output()
	if signErr != nil {
		if arg("debug") {
			fmt.Print("\033[31mfailed!\033[0m\n")
		}
		fmt.Println("\n\033[31mError: Failed to sign Shortcut, plist may be invalid!\033[0m")
		if len(signBytes) > 0 {
			fmt.Println("shortcuts:", string(signBytes))
		}
		os.Exit(1)
	}
	if arg("debug") {
		fmt.Printf("\033[32mdone!\033[0m\n")
	}
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
