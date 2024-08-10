/*
 * Copyright (c) Cherri
 */

package main

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

type ShortcutIcon struct {
	WFWorkflowIconGlyphNumber int64
	WFWorkflowIconStartColor  int
}

type ShortcutAction struct {
	WFWorkflowActionIdentifier string
	WFWorkflowActionParameters map[string]any
}

type GenericActionParameters struct {
	WFVariableName   string
	CustomOutputName string
	UUID             string
	WFInput          WFInput
}

type WFInput struct {
	Value    Value
	Variable VariableValue
}

type VariableValue struct {
	Value Value
}

type Value struct {
	Value        any
	Type         any
	VariableName string
	OutputUUID   string
	OutputName   string
	Variable     any
}
