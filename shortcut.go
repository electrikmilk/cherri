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
	Prompt                      string
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
	WFValue    any
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

var uuids map[string]string

type dictDataType int

const itemTypeText dictDataType = 0
const itemTypeNumber dictDataType = 3
const itemTypeArray dictDataType = 2
const itemTypeDict dictDataType = 1
const itemTypeBool dictDataType = 4

var noInput WFWorkflowNoInputBehavior

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
var iconColor = -1263359489

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

/* Language Codes */
var languages = []string{
	"ar_AE",
	"zh_CN",
	"zh_TW",
	"nl_NL",
	"en_GB",
	"en_US",
	"fr_FR",
	"de_DE",
	"id_ID",
	"it_IT",
	"jp_JP",
	"ko_KR",
	"pl_PL",
	"pt_BR",
	"ru_RU",
	"es_ES",
	"th_TH",
	"tr_TR",
	"vn_VN",
}

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
