/*
 * Copyright (c) Cherri
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/electrikmilk/args-parser"
	"howett.net/plist"
)

type SigningService struct {
	name string
	url  string
	info func() string
	// requestBody returns the request body and its content type, defaults to the JSON body of shortcut-signing-server.
	requestBody func() (*bytes.Buffer, string)
}

// signService can be set before signing, like in tests, to skip choosing a service from the arguments.
var signService SigningService

func chooseSigningService() SigningService {
	if args.Using("signing-server") {
		return customSigningServer()
	}

	return hubSign()
}

func customSigningServer() SigningService {
	return SigningService{
		name: "Custom Signing Server URL",
		url:  args.Value("signing-server"),
		info: func() string {
			return ansi(fmt.Sprintf("Signing using %s", args.Value("signing-server")), cyan)
		},
	}
}

// Sign the Shortcut using a signing service.
func useSigningService() {
	if signService.url == "" {
		signService = chooseSigningService()
	}

	if !args.Using("no-ansi") {
		fmt.Println(ansi(fmt.Sprintf("Signing using %s service...", signService.name), green))
		if signService.info != nil {
			fmt.Println(signService.info())
		}
	}

	var signedShortcut = requestSignedShortcut()

	if !looksLikeSignedShortcut(signedShortcut) {
		exit(fmt.Sprintf("%s did not properly sign Shortcut.", signService.name))
	}

	var writeErr = os.WriteFile(outputPath, signedShortcut, 0600)
	handle(writeErr)

	if args.Using("debug") {
		fmt.Println(ansi("Done.", green))
	}
}

func marshalShortcutPlist() []byte {
	var marshaledPlist, plistErr = plist.Marshal(shortcut, plist.XMLFormat)
	handle(plistErr)

	return marshaledPlist
}

func jsonSigningRequestBody() (*bytes.Buffer, string) {
	var payload = map[string]string{
		"shortcutName": basename,
		"shortcut":     string(marshalShortcutPlist()),
	}
	var jsonPayload, jsonErr = json.Marshal(payload)
	handle(jsonErr)

	return bytes.NewBuffer(jsonPayload), "application/json"
}

func hubSignRequestBody(apiKey string) (*bytes.Buffer, string) {
	var body bytes.Buffer
	var form = multipart.NewWriter(&body)

	handle(form.WriteField("shortcut_name", basename))
	handle(form.WriteField("api_key", apiKey))

	var shortcutFile, fileErr = form.CreateFormFile("shortcut_file", basename+".shortcut")
	handle(fileErr)
	var _, writeErr = shortcutFile.Write(marshalShortcutPlist())
	handle(writeErr)

	handle(form.Close())

	return &body, form.FormDataContentType()
}

func requestSignedShortcut() []byte {
	var requestBody = signService.requestBody
	if requestBody == nil {
		requestBody = jsonSigningRequestBody
	}
	var body, contentType = requestBody()

	var request, httpErr = http.NewRequest("POST", signService.url, body)
	handle(httpErr)
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("User-Agent", fmt.Sprintf("cherri/%s", version))

	var client = &http.Client{
		Timeout: time.Second * 30,
	}
	var response, resErr = client.Do(request)
	handle(resErr)
	defer response.Body.Close()

	// Checked before the content type because error responses can have any content type.
	if response.StatusCode != http.StatusOK {
		exit(fmt.Sprintf("Failed to sign Shortcut (%s)", response.Status))
	}

	var responseContentType = response.Header.Get("Content-Type")
	var allowedContentTypes = []string{"application/octet-stream", "application/x-plist", "application/x-apple-shortcut"}
	if !slices.Contains(allowedContentTypes, responseContentType) {
		exit(fmt.Sprintf("Unsupported response content type: %s", responseContentType))
	}

	var signedBody, readErr = io.ReadAll(response.Body)
	handle(readErr)

	return signedBody
}

// Sign the Shortcut using RoutineHub's signing service.
func hubSign() SigningService {
	var apiKey = args.Value("routinehub-api-key")
	if apiKey == "" {
		exit("Signing using HubSign requires an API key, provide one using the --routinehub-api-key= option.")
	}

	return SigningService{
		name: "HubSign",
		url:  "https://routinehub.co/api/v1/sign-shortcut",
		info: func() string {
			return fmt.Sprintf("HubSign (Shortcut signing service) - Powered By %s", ansi("RoutineHub", red))
		},
		requestBody: func() (*bytes.Buffer, string) {
			return hubSignRequestBody(apiKey)
		},
	}
}
