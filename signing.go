/*
 * Copyright (c) Cherri
 */

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/electrikmilk/args-parser"
)

var signFailed = false

// sign runs the shortcuts sign command on the unsigned shortcut file.
func sign() {
	if !darwin {
		useSigningService()
		return
	}

	var signingMode = "people-who-know-me"
	if args.Using("share") && args.Value("share") == "anyone" {
		signingMode = "anyone"
	}

	if args.Using("debug") {
		fmt.Printf("Signing %s to %s...", inputPath, outputPath)
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
		signFailed = true
		if args.Using("debug") {
			fmt.Print(ansi("Failed!\n", red))
		}

		fmt.Printf("%s\n%s\n", ansi("Failed to sign Shortcut using macOS :(", orange, bold), ansi(stdErr.String(), orange))

		useSigningService()
	}

	if args.Using("debug") {
		fmt.Println(ansi("Done.", green))
	}
}

func removeUnsigned() {
	var _, signedStatErr = os.Stat(fmt.Sprintf("%s%s.shortcut", relativePath, workflowName))
	if os.IsNotExist(signedStatErr) {
		return
	}
	var _, unsignedStatErr = os.Stat(fmt.Sprintf("%s%s%s", relativePath, workflowName, unsignedEnd))
	if os.IsNotExist(unsignedStatErr) {
		return
	}

	if args.Using("debug") {
		fmt.Printf("Removing %s%s...", workflowName, unsignedEnd)
	}

	var removeErr = os.Remove(inputPath)
	handle(removeErr)

	if args.Using("debug") {
		fmt.Println(ansi("Done.", green))
	}
}

// looksLikeSignedShortcut performs quick checks to make sure response is a signed Shortcut.
func looksLikeSignedShortcut(buffer []byte) bool {
	if len(buffer) >= 4 && string(buffer[:4]) == "AEA1" {
		return true
	}
	return false
}
