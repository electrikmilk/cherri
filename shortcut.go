/*
 * Copyright (c) Cherri
 */

package main

/*
 Shortcut File Format Data Structures
*/

type ShortcutIcon struct {
	WFWorkflowIconGlyphNumber int64
	WFWorkflowIconStartColor  int
}

type Shortcut struct {
	WFWorkflowIcon                       ShortcutIcon
	WFWorkflowActions                    []ShortcutAction
	WFQuickActionSurfaces                []string
	WFWorkflowInputContentItemClasses    []string
	WFWorkflowClientVersion              string
	WFWorkflowMinimumClientVersion       int
	WFWorkflowMinimumClientVersionString string
	WFWorkflowImportQuestions            interface{}
	WFWorkflowTypes                      []string
	WFWorkflowOutputContentItemClasses   []string
	WFWorkflowHasShortcutInputVariables  bool
	WFWorkflowHasOutputFallback          bool
	WFWorkflowNoInputBehavior            WFWorkflowNoInputBehavior
	WFWorkflowName                       string
}

var shortcut Shortcut

type ShortcutAction struct {
	WFWorkflowActionIdentifier string
	WFWorkflowActionParameters map[string]any
}

type WFWorkflowNoInputBehavior struct {
	Name       string
	Parameters map[string]string
}

type Value struct {
	Type                        string
	VariableName                string
	OutputUUID                  string
	OutputName                  string
	Value                       any
	Variable                    any
	WFDictionaryFieldValueItems []WFDictionaryFieldValueItem
	AttachmentsByRange          map[string]Value
	String                      string
	Aggrandizements             []Aggrandizement
}

type Aggrandizement struct {
	Type              string
	CoercionItemClass string
	DictionaryKey     string
	PropertyName      string
	PropertyUserInfo  any
}

type WFDictionaryFieldValueItem struct {
	WFKey      any
	WFItemType int
	WFValue    WFValue
}

type WFValue struct {
	Value               any
	String              string
	WFSerializationType string
}

type WFInput struct {
	Value Value
}

type WFContactFieldValue struct {
	EntryType       int
	SerializedEntry map[string]interface{}
}

type ImageSize struct {
	Value SizeValue
}

type SizeValue struct {
	Unit      string
	Magnitude string
}

type WFMeasurementUnit struct {
	WFNSUnitSymbol any
	Value          SizeValue
}
