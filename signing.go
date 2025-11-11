/*
 * Copyright (c) Cherri
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"time"

	"github.com/electrikmilk/args-parser"
)

var signFailed = false
var signingServiceFailed = false
var backoff = 10

// sign runs the shortcuts sign command on the unsigned shortcut file.
func sign() {
	if !darwin {
		useHubSign()
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

		useHubSign()
	}
}

func useHubSign() {
	var hubSignService = hubSign()
	useSigningService(&hubSignService)
}

type SigningService struct {
	name string
	url  string
	info func() string
}

// Sign the Shortcut using a signing service.
func useSigningService(service *SigningService) {
	if signingServiceFailed {
		fmt.Println(ansi(fmt.Sprintf("Backing off from %s", service.name), red))
		for i := 5; i > 0; i-- {
			fmt.Printf("%d seconds...\r", i)
			time.Sleep(1 * time.Second)
		}
		fmt.Print("\n\n")
	}

	if !args.Using("no-ansi") {
		fmt.Println(ansi(fmt.Sprintf("Signing using %s service...", service.name), green))
		if service.info != nil {
			fmt.Println(service.info())
		}
	}

	var payload = map[string]string{
		"shortcutName": basename,
		"shortcut":     compiled,
	}
	var jsonPayload, jsonErr = json.Marshal(payload)
	handle(jsonErr)

	var request, httpErr = http.NewRequest("POST", service.url, bytes.NewReader(jsonPayload))
	handle(httpErr)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", fmt.Sprintf("cherri/%s", version))

	var client = &http.Client{
		Timeout: time.Second * 30,
	}
	var response, resErr = client.Do(request)
	handle(resErr)
	defer response.Body.Close()

	var responseContentType = response.Header.Get("Content-Type")
	var allowedContentTypes = []string{"application/octet-stream", "application/x-plist", "application/x-apple-shortcut"}
	if !slices.Contains(allowedContentTypes, responseContentType) {
		exit(fmt.Sprintf("Unsupported response type: %s", responseContentType))
	}

	if response.StatusCode != http.StatusOK {
		signingServiceFailed = true
		backoff += 10
		fmt.Println(ansi(fmt.Sprintf("Failed to sign Shortcut (%s)", response.Status), red))
		return
	}

	signingServiceFailed = false

	if backoff > 10 {
		backoff -= 10
	} else {
		backoff = 10
	}

	var body, readErr = io.ReadAll(response.Body)
	handle(readErr)

	if !looksLikeSignedShortcut(body) {
		exit("Signing server response does not look like a Shortcut file.")
	}

	var writeErr = os.WriteFile(outputPath, body, 0600)
	handle(writeErr)

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

	removeErr := os.Remove(inputPath)
	handle(removeErr)

	if args.Using("debug") {
		fmt.Println(ansi("Done.", green))
	}
}

// Sign the Shortcut using RoutineHub's signing service.
func hubSign() SigningService {
	return SigningService{
		name: "HubSign",
		url:  "https://hubsign.routinehub.services/sign",
		info: func() string {
			return fmt.Sprintf("Shortcut Signing Powered By %s", ansi("RoutineHub", red))
		},
	}
}
