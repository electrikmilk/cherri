/*
 * Copyright (c) Brandon Jordan
 */

package main

import "strings"

var variables map[string]variableValue

type variableValue struct {
	variableType string
	valueType    tokenType
	value        any
	getAs        string
	coerce       string
	constant     bool
	repeatItem   bool
}

var globals map[string]variableValue

func makeGlobals() {
	if len(globals) != 0 {
		return
	}
	globals = make(map[string]variableValue)
	globals["ShortcutInput"] = variableValue{
		variableType: "ExtensionInput",
		valueType:    String,
		value:        "ShortcutInput",
	}
	globals["CurrentDate"] = variableValue{
		variableType: "CurrentDate",
		valueType:    Date,
		value:        "CurrentDate",
	}
	globals["Clipboard"] = variableValue{
		variableType: "Clipboard",
		valueType:    String,
		value:        "Clipboard",
	}
	globals["Device"] = variableValue{
		variableType: "DeviceDetails",
		valueType:    String,
		value:        "DeviceDetails",
	}
	globals["Ask"] = variableValue{
		variableType: "Ask",
		valueType:    String,
		value:        "Ask",
	}
	globals["RepeatItem"] = variableValue{
		variableType: "Variable",
		valueType:    String,
		value:        "Repeat Item",
	}
	globals["RepeatIndex"] = variableValue{
		variableType: "Variable",
		valueType:    String,
		value:        "Repeat Index",
	}
}

func validReference(identifier string) bool {
	if globalExists(identifier) {
		return true
	}
	if variableExists(identifier) {
		return true
	}

	return false
}

func globalExists(identifier string) bool {
	if _, found := globals[identifier]; found {
		isInputVariable(identifier)
		return true
	}

	return false
}

func variableExists(identifier string) bool {
	identifier = strings.ToLower(identifier)
	if _, found := variables[identifier]; found {
		return true
	}

	return false
}

func variable(variableName string) (exists bool, global bool, variable variableValue) {
	if value, found := globals[variableName]; found {
		return true, true, value
	}

	variableName = strings.ToLower(variableName)
	if value, found := variables[variableName]; found {
		return true, false, value
	}

	return false, false, variableValue{}
}
