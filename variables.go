/*
 * Copyright (c) Cherri
 */

package main

import "fmt"

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

var globals = map[string]variableValue{
	"ShortcutInput": {
		variableType: "ExtensionInput",
		valueType:    String,
		value:        "ShortcutInput",
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
		valueType:    String,
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
		parserError(fmt.Sprintf("Variable conflicts with defined import question '%s'.", *identifier))
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

func getVariableValue(identifier string) (*variableValue, bool) {
	if value, found := globals[identifier]; found {
		return &value, true
	}
	if value, found := variables[identifier]; found {
		return &value, true
	}

	return nil, false
}
