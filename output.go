/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/electrikmilk/args-parser"
)

// writeFile writes plist in bytes to filename.
func writeFile(filename string, debug string) {
	var writeDebugOutput = args.Using("debug")
	if writeDebugOutput {
		fmt.Printf("Writing to %s...", debug)
	}

	writeErr := os.WriteFile(filename, []byte(compiled), 0600)
	handle(writeErr)

	if writeDebugOutput {
		fmt.Println(ansi("Done.", green))
	}
}

// createShortcut writes the Shortcut files to disk and signs them if the unsigned argument is not unused.
func createShortcut() {
	var path = fmt.Sprintf("%s%s", relativePath, workflowName)
	if args.Using("debug") {
		writeFile(path+".plist", workflowName+".plist")
	}
	writeFile(path+unsignedEnd, workflowName+unsignedEnd)

	inputPath = fmt.Sprintf("%s%s%s", relativePath, workflowName, unsignedEnd)

	if !args.Using("skip-sign") {
		if args.Using("hubsign") {
			hubSign()
		} else {
			sign()
		}

		removeUnsigned()
	}

	if args.Using("import") {
		openShortcut()
	}
}

func openShortcut() {
	var _, importErr = exec.Command("open", outputPath).Output()
	handle(importErr)
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
	fmt.Printf("%s %d:%d\n", currentChar, lineIdx+1, lineCharIdx+1)
}

type outputType int

const (
	bold      outputType = 1
	dim       outputType = 2
	italic    outputType = 3
	underline outputType = 4
	red       outputType = 31
	green     outputType = 32
	yellow    outputType = 33
	cyan      outputType = 36
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
