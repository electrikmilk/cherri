/*
 * Copyright (c) Cherri
 */

package main

/*
 Shortcut File Format Data Structures
*/

type ShortcutIcon struct {
	WFWorkflowIconGlyphNumber int64 `plist:",omitempty"`
	WFWorkflowIconStartColor  int   `plist:",omitempty"`
}

type Shortcut struct {
	WFWorkflowIcon                       ShortcutIcon     `plist:",omitempty"`
	WFWorkflowActions                    []ShortcutAction `plist:",omitempty"`
	WFQuickActionSurfaces                []string         `plist:",omitempty"`
	WFWorkflowInputContentItemClasses    []string         `plist:",omitempty"`
	WFWorkflowClientVersion              string           `plist:",omitempty"`
	WFWorkflowMinimumClientVersion       int              `plist:",omitempty"`
	WFWorkflowMinimumClientVersionString string           `plist:",omitempty"`
	WFWorkflowImportQuestions            interface{}      `plist:",omitempty"`
	WFWorkflowTypes                      []string         `plist:",omitempty"`
	WFWorkflowOutputContentItemClasses   []string         `plist:",omitempty"`
	WFWorkflowHasShortcutInputVariables  bool             `plist:",omitempty"`
	WFWorkflowHasOutputFallback          bool             `plist:",omitempty"`
	WFWorkflowNoInputBehavior            map[string]any   `plist:",omitempty"`
}

var shortcut Shortcut

type ShortcutAction struct {
	WFWorkflowActionIdentifier string
	WFWorkflowActionParameters map[string]any `plist:",omitempty"`
}

type Value struct {
	Type                        string                       `plist:",omitempty"`
	VariableName                string                       `plist:",omitempty"`
	OutputUUID                  string                       `plist:",omitempty"`
	OutputName                  string                       `plist:",omitempty"`
	Value                       any                          `plist:",omitempty"`
	Variable                    any                          `plist:",omitempty"`
	WFDictionaryFieldValueItems []WFDictionaryFieldValueItem `plist:",omitempty"`
	AttachmentsByRange          map[string]Value             `plist:",omitempty"`
	String                      string                       `plist:",omitempty"`
	Aggrandizements             []Aggrandizement             `plist:",omitempty"`
	Prompt                      string                       `plist:",omitempty"`
}

type Aggrandizement struct {
	Type              string `plist:",omitempty"`
	CoercionItemClass string `plist:",omitempty"`
	DictionaryKey     string `plist:",omitempty"`
	PropertyName      string `plist:",omitempty"`
	PropertyUserInfo  any    `plist:",omitempty"`
}

type WFDictionaryFieldValueItem struct {
	WFKey      any `plist:",omitempty"`
	WFItemType int `plist:",omitempty"`
	WFValue    any `plist:",omitempty"`
}

type WFValue struct {
	Value               any    `plist:",omitempty"`
	String              string `plist:",omitempty"`
	WFSerializationType string `plist:",omitempty"`
}

type WFInput struct {
	Value Value `plist:",omitempty"`
}

type WFContactFieldValue struct {
	EntryType       int                    `plist:",omitempty"`
	SerializedEntry map[string]interface{} `plist:",omitempty"`
}

type SizeValue struct {
	Unit      string `plist:",omitempty"`
	Magnitude string `plist:",omitempty"`
}

type WFMeasurementUnit struct {
	WFNSUnitSymbol any       `plist:",omitempty"`
	Value          SizeValue `plist:",omitempty"`
}

type WFTextTokenAttachment struct {
	Value               Value  `plist:",omitempty"`
	WFSerializationType string `plist:",omitempty"`
}

type WFTextTokenString struct {
	Value               WFTextTokenStringValue `plist:",omitempty"`
	WFSerializationType string                 `plist:",omitempty"`
}

type WFTextTokenStringValue struct {
	AttachmentsByRange map[string]Value `plist:"attachmentsByRange,omitempty"`
	String             string           `plist:"string,omitempty"`
}

type WFDictionaryFieldValue struct {
	Value               WFDictionaryFieldValueWrapper `plist:",omitempty"`
	WFSerializationType string                        `plist:",omitempty"`
}

type WFDictionaryFieldValueWrapper struct {
	WFDictionaryFieldValueItems []WFDictionaryFieldValueItem `plist:",omitempty"`
}

type WFContentPredicateTableTemplate struct {
	Value               WFConditionValue `plist:",omitempty"`
	WFSerializationType string           `plist:",omitempty"`
}

type WFConditionValue struct {
	WFActionParameterFilterPrefix    int                `plist:",omitempty"`
	WFActionParameterFilterTemplates []WFConditionParam `plist:",omitempty"`
}

type WFConditionParam struct {
	WFCondition               int             `plist:",omitempty"`
	WFInput                   WFInputVariable `plist:",omitempty"`
	WFConditionalActionString any             `plist:",omitempty"`
	WFNumberValue             any             `plist:",omitempty"`
	WFAnotherNumber           any             `plist:",omitempty"`
}

type WFInputVariable struct {
	Type     string                `plist:",omitempty"`
	Variable WFTextTokenAttachment `plist:",omitempty"`
}

type WFColorValue struct {
	WFColorRepresentationType string  `plist:",omitempty"`
	RedComponent              float64 `plist:"redComponent,omitempty"`
	GreenComponent            float64 `plist:"greenComponent,omitempty"`
	BlueComponent             float64 `plist:"blueComponent,omitempty"`
	AlphaComponent            any     `plist:"alphaComponent,omitempty"`
}

type WFActionReference struct {
	CustomOutputName string `plist:",omitempty"`
	UUID             string `plist:",omitempty"`
}

type WFSetVariableParams struct {
	WFVariableName      string `plist:",omitempty"`
	WFInput             any    `plist:",omitempty"`
	WFSerializationType string `plist:",omitempty"`
	CustomOutputName    string `plist:",omitempty"`
	UUID                string `plist:",omitempty"`
}

type WFConditionalActionParams struct {
	GroupingIdentifier        string `plist:",omitempty"`
	WFControlFlowMode         uint64 `plist:",omitempty"`
	WFInput                   any    `plist:",omitempty"`
	WFCondition               int    `plist:",omitempty"`
	WFConditionalActionString any    `plist:",omitempty"`
	WFNumberValue             any    `plist:",omitempty"`
	WFAnotherNumber           any    `plist:",omitempty"`
	WFConditions              any    `plist:",omitempty"`
	UUID                      string `plist:",omitempty"`
}

type WFMenuParams struct {
	GroupingIdentifier string       `plist:",omitempty"`
	WFControlFlowMode  uint64       `plist:",omitempty"`
	WFMenuPrompt       any          `plist:",omitempty"`
	WFMenuItems        []WFMenuItem `plist:",omitempty"`
	UUID               string       `plist:",omitempty"`
}

type WFMenuItem struct {
	WFItemType int `plist:",omitempty"`
	WFValue    any `plist:",omitempty"`
}

type WFMenuItemParams struct {
	GroupingIdentifier        string `plist:",omitempty"`
	WFControlFlowMode         uint64 `plist:",omitempty"`
	WFMenuItemAttributedTitle any    `plist:",omitempty"`
	WFMenuItemTitle           any    `plist:",omitempty"`
}

type WFRepeatParams struct {
	GroupingIdentifier string `plist:",omitempty"`
	WFControlFlowMode  uint64 `plist:",omitempty"`
	WFRepeatCount      any    `plist:",omitempty"`
	WFInput            any    `plist:",omitempty"`
	UUID               string `plist:",omitempty"`
}

var uuids map[string]string

type dictDataType int

const itemTypeText dictDataType = 0
const itemTypeNumber dictDataType = 3
const itemTypeArray dictDataType = 2
const itemTypeDict dictDataType = 1
const itemTypeBool dictDataType = 4

var noInput map[string]any

var hasShortcutInputVariables = false

// ObjectReplaceChar is a Shortcuts convention to mark the placement of inline variables in a string.
const ObjectReplaceChar = '\uFFFC'
const ObjectReplaceCharStr = "\uFFFC"

var workflowName string

var definitions map[string]any

/* Colors */

var colors = map[string]int{
	"red":        4282601983,
	"darkorange": 4251333119,
	"orange":     4271458815,
	"yellow":     4274264319,
	"green":      4292093695,
	"teal":       431817727,
	"lightblue":  1440408063,
	"blue":       463140863,
	"darkblue":   946986751,
	"violet":     2071128575,
	"purple":     3679049983,
	"pink":       3980825855,
	"darkgray":   255,
	"gray":       3031607807,
	"taupe":      2846468607,
}
var iconColor = 3031607807

var altColors = map[string]int{
	"red":        -12365313,
	"darkorange": -43634177,
	"orange":     -23508481,
	"yellow":     -20702977,
	"green":      -2873601,
	"teal":       -3863149569,
	"lightblue":  -2854559233,
	"blue":       -3831826433,
	"darkblue":   -3347980545,
	"violet":     -2223838721,
	"purple":     -615917313,
	"pink":       -314141441,
	"darkgray":   -4294967041,
	"gray":       -1263359489,
	"taupe":      -1448498689,
}

/* Inputs */

var contentItems = map[string]string{
	"app":         "WFAppStoreAppContentItem",
	"article":     "WFArticleContentItem",
	"contact":     "WFContactContentItem",
	"date":        "WFDateContentItem",
	"email":       "WFEmailAddressContentItem",
	"folder":      "WFFolderContentItem",
	"file":        "WFGenericFileContentItem",
	"image":       "WFImageContentItem",
	"itunes":      "WFiTunesProductContentItem",
	"location":    "WFLocationContentItem",
	"maplink":     "WFDCMapsLinkContentItem",
	"media":       "WFAVAssetContentItem",
	"pdf":         "WFPDFContentItem",
	"phonenumber": "WFPhoneNumberContentItem",
	"richtext":    "WFRichTextContentItem",
	"webpage":     "WFSafariWebPageContentItem",
	"text":        "WFStringContentItem",
	"dictionary":  "WFDictionaryContentItem",
	"number":      "WFNumberContentItem",
	"url":         "WFURLContentItem",
}

var revContentItems map[string]string

func reversedContentItems() map[string]string {
	if len(revContentItems) == 0 {
		var revContentItems = make(map[string]string)
		for key, item := range contentItems {
			revContentItems[item] = key
		}
	}

	return revContentItems
}

var inputs []string
var outputs []string

/* Workflow Types */

var workflowTypes = map[string]string{
	"menubar":       "MenuBar",
	"quickactions":  "QuickActions",
	"sharesheet":    "ActionExtension",
	"notifications": "NCWidget",
	"sleepmode":     "Sleep",
	"watch":         "Watch",
	"onscreen":      "ReceivesOnScreenContent",
	"search":        "WFWorkflowTypeShowInSearch",
	"spotlight":     "WFWorkflowTypeReceivesInputFromSearch",
}
var definedWorkflowTypes []string

/* Quick Actions */

var quickActions = map[string]string{
	"finder":   "Finder",
	"services": "Services",
}
var definedQuickActions []string

/* Versions */

var versions = map[string]string{
	"26":     "4033.0.4.3",
	"18.4":   "3218.0.4.100",
	"18":     "3036.0.4.2",
	"17":     "2106.0.3",
	"16.5":   "900",
	"16.4":   "900",
	"16.3":   "900",
	"16.2":   "900",
	"16":     "900",
	"15.7.2": "800",
	"15":     "800",
	"14":     "700",
	"13":     "600",
	"12":     "500",
}
var clientVersion = versions["26"]
var iosVersion = 26.0

/* Conditionals */

type WFConditions struct {
	conditions                    []condition
	WFActionParameterFilterPrefix int
}

type condition struct {
	condition int
	arguments []actionArgument
}

var conditions = map[tokenType]int{
	Is:             4,
	Not:            5,
	Any:            100,
	Empty:          101,
	Contains:       99,
	DoesNotContain: 999,
	BeginsWith:     8,
	EndsWith:       9,
	GreaterThan:    2,
	GreaterOrEqual: 3,
	LessThan:       0,
	LessOrEqual:    1,
	Between:        1003,
}

var conditionFilterPrefixes = map[tokenType]int{
	Or:  0,
	And: 1,
}

var allowedConditionalTypes = map[tokenType][]tokenType{
	Is:             {String, Integer, Bool, Action},
	Not:            {String, Integer, Bool, Action},
	Any:            {},
	Empty:          {},
	Contains:       {String, Arr},
	DoesNotContain: {String, Arr},
	BeginsWith:     {String},
	EndsWith:       {String},
	GreaterThan:    {Integer, Float},
	GreaterOrEqual: {Integer, Float},
	LessThan:       {Integer, Float},
	LessOrEqual:    {Integer, Float},
	Between:        {Integer, Float},
}

/* Menus */

var menus map[string][]varValue
