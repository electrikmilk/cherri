/*
 * Copyright (c) Cherri
 */

package main

import (
	"crypto/sha1"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/electrikmilk/args-parser"
	"github.com/google/uuid"
)

var variables map[string]varValue

type varValue struct {
	variableType string
	valueType    tokenType
	value        any
	getAs        string
	coerce       string
	constant     bool
	repeatItem   bool
	prompt       string
}

var globals = map[string]varValue{
	ShortcutInput: {
		variableType: "ExtensionInput",
		valueType:    String,
		value:        ShortcutInput,
	},
	"CurrentDate": {
		variableType: "CurrentDate",
		valueType:    Date,
		value:        "CurrentDate",
	},
	"Clipboard": {
		variableType: "Clipboard",
		valueType:    String,
		value:        "Clipboard",
	},
	"Device": {
		variableType: "DeviceDetails",
		valueType:    String,
		value:        "DeviceDetails",
	},
	"Ask": {
		variableType: "Ask",
		valueType:    String,
		value:        "Ask",
	},
	"RepeatItem": {
		variableType: "Variable",
		valueType:    String,
		value:        "Repeat Item",
	},
	"RepeatIndex": {
		variableType: "Variable",
		valueType:    Integer,
		value:        "Repeat Index",
	},
}

func availableIdentifier(identifier *string) {
	if v, found := variables[*identifier]; found {
		if v.constant {
			parserError(fmt.Sprintf("Cannot redefine constant '%s'.", *identifier))
		}
		if v.repeatItem {
			parserError(fmt.Sprintf("Cannot redefine repeat item '%s'.", *identifier))
		}
	}
	if _, found := globals[*identifier]; found {
		parserError(fmt.Sprintf("Cannot redefine global variable '%s'.", *identifier))
	}
	if _, found := questions[*identifier]; found {
		parserError(fmt.Sprintf("Reference conflicts with defined import question '%s'.", *identifier))
	}
}

func createUUIDReference(identifier string) string {
	var actionUUID = createUUID(&identifier)
	uuids[identifier] = actionUUID

	return actionUUID
}

func wrapVariableReference(s *string) {
	if validReference(*s) {
		*s = fmt.Sprintf("{%s}", *s)
	}
}

func validReference(identifier string) bool {
	if _, found := globals[identifier]; found {
		isInputVariable(identifier)
		return true
	}
	if _, found := variables[identifier]; found {
		return true
	}

	return false
}

func getVariableValue(identifier string) (variableValue *varValue, found bool) {
	if value, found := globals[identifier]; found {
		isInputVariable(identifier)
		return &value, true
	}
	if value, found := variables[identifier]; found {
		return &value, true
	}

	return nil, false
}

func makeVariableReferenceString(value varValue) string {
	var identifier strings.Builder
	identifier.WriteString(value.value.(string))

	if value.getAs != "" {
		identifier.WriteString(fmt.Sprintf("['%s']", value.getAs))
	}
	if value.coerce != "" {
		identifier.WriteString(fmt.Sprintf(".%s", value.coerce))
	}

	return identifier.String()
}

var currentOutputName string
var duplicateDelta int

func checkDuplicateOutputName(name string) string {
	if name != currentOutputName {
		currentOutputName = name
		duplicateDelta = 0
	}
	if _, found := uuids[name]; found {
		return checkDuplicateOutputName(duplicateOutputName())
	} else if args.Using("import") {
		for _, outputName := range uuids {
			if outputName == currentOutputName {
				return checkDuplicateOutputName(duplicateOutputName())
			}
		}
	}

	return currentOutputName
}

func duplicateOutputName() string {
	duplicateDelta++

	var revChars = []rune(currentOutputName)
	slices.Reverse(revChars)
	var numChars []rune
	for _, rc := range revChars {
		if rc >= '0' && rc <= '9' {
			numChars = append(numChars, rc)
		} else {
			break
		}
	}

	if len(numChars) != 0 {
		slices.Reverse(numChars)
		var num = string(numChars)
		var endingDelta, numErr = strconv.Atoi(num)
		handle(numErr)
		endingDelta++

		currentOutputName, _ = strings.CutSuffix(currentOutputName, num)
		duplicateDelta = endingDelta
	}

	if !args.Using("import") && !strings.Contains(currentOutputName, " ") {
		currentOutputName = fmt.Sprintf("%s ", currentOutputName)
	}

	currentOutputName = fmt.Sprintf("%s%d", currentOutputName, duplicateDelta)

	return currentOutputName
}

func createUUID(salt *string) string {
	var deterministic = args.Using("derive-uuids")
	if deterministic {
		var seeds = [][]byte{
			[]byte(workflowName),
			[]byte(*salt),
		}
		return deterministicUUID(seeds...)
	}

	return uuid.New().String()
}

// deterministicUUID returns a UUID derived from the given pieces.
func deterministicUUID(pieces ...[]byte) string {
	var hash = sha1.New()
	for i, piece := range pieces {
		if i > 0 {
			hash.Write([]byte{0x00})
		}
		hash.Write(piece)
	}
	var seed = hash.Sum(nil)

	return uuid.NewSHA1(uuid.NameSpaceURL, seed).String()
}
