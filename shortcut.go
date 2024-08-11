/*
 * Copyright (c) Cherri
 */

package main

type ShortcutIcon struct {
	WFWorkflowIconGlyphNumber int64
	WFWorkflowIconStartColor  int
}

type Shortcut struct {
	WFWorkflowIcon                      ShortcutIcon
	WFWorkflowActions                   []ShortcutAction
	WFQuickActionSurfaces               []string
	WFWorkflowInputContentItemClasses   []string
	WFWorkflowClientVersion             string
	WFWorkflowMinimumClientVersion      int
	WFWorkflowImportQuestions           interface{}
	WFWorkflowTypes                     []string
	WFWorkflowOutputContentItemClasses  []string
	WFWorkflowHasShortcutInputVariables bool
	WFWorkflowHasOutputFallback         bool
}

type ShortcutAction struct {
	WFWorkflowActionIdentifier string
	WFWorkflowActionParameters map[string]any
}

type GenericShortcut struct {
	WFWorkflowActions []GenericShortcutAction
}

type GenericShortcutAction struct {
	WFWorkflowActionIdentifier string
	WFWorkflowActionParameters GenericActionParameters
}

type GenericActionParameters struct {
	WFVariableName   string
	CustomOutputName string
	UUID             string
	WFInput          WFInput
}

type WFInput struct {
	Value      Value
	Variable   VariableValue
	OutputName string
	OutputUUID string
}

type VariableValue struct {
	Value Value
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
	Value  any
	String string
}

type ArrayValue struct {
	WFValue    any
	WFItemType int
}
