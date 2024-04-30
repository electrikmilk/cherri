/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/electrikmilk/args-parser"
	"io"
	"net/http"
	"os"
	"os/exec"
)

// sign runs the shortcuts sign command on the unsigned shortcut file.
func sign() {
	if !darwin {
		fmt.Println(ansi("Warning:", bold, yellow), "macOS is required to sign shortcuts. The compiled Shortcut will not run on iOS 15+ or macOS 12+.")

		fmt.Print("\n")
		fmt.Println("However...")
		fmt.Println(ansi("NEW!", red), "Use", ansi("--hubsign", cyan), "to use RoutineHub's remote service to sign the compiled Shortcut.")
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
		if args.Using("debug") {
			fmt.Print(ansi("Failed!\n", red))
		}

		fmt.Printf("%s\n%s\n", ansi("Failed to sign Shortcut using macOS :(", yellow, bold), ansi(stdErr.String(), yellow))
		hubSign()
	}
}

const HubSignURL = "https://hubsign.routinehub.services/sign"

// Sign the Shortcut using RoutineHub's signing service.
func hubSign() {
	if args.Using("debug") {
		fmt.Print("Attempting to sign using HubSign...")
	}

	if !args.Using("no-ansi") {
		fmt.Println(ansi("Attempting to sign using HubSign service...", green))
		fmt.Println(ansi("Shortcut Signing Powered By RoutineHub", dim))
	}

	var payload = map[string]string{
		"shortcutName": basename,
		"shortcut":     compiled,
	}
	var jsonPayload, jsonErr = json.Marshal(payload)
	handle(jsonErr)

	var request, httpErr = http.NewRequest("POST", HubSignURL, bytes.NewReader(jsonPayload))
	handle(httpErr)
	request.Header.Set("Content-Type", "application/json")

	var client = &http.Client{}
	response, resErr := client.Do(request)
	handle(resErr)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		exit(fmt.Sprintf("Failed to sign Shortcut (%s)", response.Status))
	}

	var body, readErr = io.ReadAll(response.Body)
	handle(readErr)

	var writeErr = os.WriteFile(outputPath, body, 0600)
	handle(writeErr)

	if args.Using("debug") {
		fmt.Println(ansi("Done.", green))
	}
}

func removeUnsigned() {
	var _, statErr = os.Stat(inputPath)
	if os.IsNotExist(statErr) {
		return
	}

	if args.Using("debug") {
		fmt.Printf("Removing %s_unsigned.shortcut...", workflowName)
	}

	removeErr := os.Remove(inputPath)
	handle(removeErr)

	if args.Using("debug") {
		fmt.Println(ansi("Done.", green))
	}
}
