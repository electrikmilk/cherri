/*
 * Copyright (c) Cherri
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
	"time"
)

var signFailed = false
var hubSignFailed = false
var hubSignBackoff = 10

// sign runs the shortcuts sign command on the unsigned shortcut file.
func sign() {
	if !darwin {
		fmt.Println(ansi("Warning:", bold, yellow), "macOS is required to sign shortcuts. The compiled Shortcut will not run on iOS 15+ or macOS 12+.")

		if !args.Using("no-ansi") {
			fmt.Print("\n")
			fmt.Println("However...")
			fmt.Println(ansi("NEW!", red), "Use", ansi("--hubsign", cyan), "to use RoutineHub's remote service to sign the compiled Shortcut.")
		}
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

		fmt.Printf("%s\n%s\n", ansi("Failed to sign Shortcut using macOS :(", yellow, bold), ansi(stdErr.String(), yellow))
		hubSign()
	}
}

const HubSignURL = "https://hubsign.routinehub.services/sign"

// Sign the Shortcut using RoutineHub's signing service.
func hubSign() {
	if hubSignFailed {
		fmt.Println(ansi("Backing off from HubSign", red))
		for i := 5; i > 0; i-- {
			fmt.Printf("%d seconds...\r", i)
			time.Sleep(1 * time.Second)
		}
		fmt.Print("\n\n")
	}

	if !args.Using("no-ansi") {
		fmt.Println(ansi("Signing using HubSign service...", green))
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
		hubSignFailed = true
		hubSignBackoff += 10
		fmt.Println(ansi(fmt.Sprintf("Failed to sign Shortcut (%s)", response.Status), red))
		return
	}

	if hubSignBackoff > 10 {
		hubSignBackoff -= 10
	} else {
		hubSignBackoff = 10
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
	var _, signedStatErr = os.Stat(fmt.Sprintf("%s%s.shortcut", relativePath, workflowName))
	if os.IsNotExist(signedStatErr) {
		return
	}
	var _, unsignedStatErr = os.Stat(fmt.Sprintf("%s%s%s", relativePath, workflowName, unsignedEnd))
	if os.IsNotExist(unsignedStatErr) {
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
