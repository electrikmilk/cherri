/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/electrikmilk/args-parser"
	"howett.net/plist"
)

// createShortcut writes the Shortcut files to disk and signs them if the unsigned argument is not unused.
func createShortcut() {
	var path = fmt.Sprintf("%s%s", relativePath, workflowName)
	if args.Using("debug") {
		writeShortcut(path+".plist", workflowName+".plist")
	}
	writeShortcut(path+unsignedEnd, workflowName+unsignedEnd)

	inputPath = fmt.Sprintf("%s%s%s", relativePath, workflowName, unsignedEnd)

	if !args.Using("skip-sign") {
		switch {
		case args.Using("signing-server"):
			useSigningService(&SigningService{
				name: "Custom Signing Server URL",
				url:  args.Value("signing-server"),
			})
		case args.Using("hubsign"):
			useHubSign()
		default:
			sign()
		}

		removeUnsigned()

		if args.Using("open") {
			openShortcut()
		}
	}
}

// writeShortcut encodes shortcut by writing plist data at path.
func writeShortcut(path string, debug string) {
	var writeDebugOutput = args.Using("debug")
	if writeDebugOutput {
		fmt.Printf("Writing to %s...", debug)
	}

	var unsignedFile, unsignedFileErr = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	handle(unsignedFileErr)
	defer unsignedFile.Close()

	var plistEncoder = plist.NewEncoder(unsignedFile)
	if args.Using("debug") {
		plistEncoder.Indent("\t")
	}

	var encodeErr = plistEncoder.Encode(shortcut)
	handle(encodeErr)

	if writeDebugOutput {
		fmt.Println(ansi("Done.", green))
	}
}

func openShortcut() {
	var _, importErr = exec.Command("open", outputPath).Output()
	handle(importErr)
}

func printChar(ch rune, chLineIdx int, chLineCharIdx int) {
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
	fmt.Printf("%s %d:%d\n", currentChar, chLineIdx+1, chLineCharIdx+1)
}

type outputType int

const (
	bold      outputType = 1
	dim       outputType = 2
	italic    outputType = 3
	underline outputType = 4
	red       outputType = 31
	green     outputType = 32
	yellow    outputType = 93
	orange    outputType = 33
	magenta   outputType = 95
	blue      outputType = 34
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
