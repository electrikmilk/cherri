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
	WFWorkflowIcon                       ShortcutIcon
	WFWorkflowActions                    []ShortcutAction
	WFQuickActionSurfaces                []string
	WFWorkflowInputContentItemClasses    []string
	WFWorkflowClientVersion              string
	WFWorkflowMinimumClientVersion       int
	WFWorkflowMinimumClientVersionString string
	WFWorkflowImportQuestions            []WFQuestion
	WFWorkflowTypes                      []string
	WFWorkflowOutputContentItemClasses   []string
	WFWorkflowHasShortcutInputVariables  bool
	WFWorkflowHasOutputFallback          bool
	WFWorkflowNoInputBehavior            map[string]any
	WFWorkflowTriggers                   []AutomationTrigger
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
	WFItemType int
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
	Value               any    `plist:",omitempty"`
	WFSerializationType string `plist:",omitempty"`
}

type WFDictionaryFieldValueWrapper struct {
	WFDictionaryFieldValueItems []WFDictionaryFieldValueItem `plist:",omitempty"`
}

type WFArrayValue struct {
	WFSerializationType string                       `plist:",omitempty"`
	Value               []WFDictionaryFieldValueItem `plist:",omitempty"`
}

type WFBoolValue struct {
	WFSerializationType string `plist:",omitempty"`
	Value               bool   `plist:"Value"` // false is the zero value; omitempty would suppress it
}

type WFQuantityValue struct {
	Magnitude any `plist:",omitempty"`
	Unit      any `plist:",omitempty"`
}

type WFQuantityFieldValue struct {
	Value               WFQuantityValue `plist:",omitempty"`
	WFSerializationType string          `plist:",omitempty"`
}

type WFMenuItem struct {
	WFItemType int
	WFValue    any `plist:",omitempty"`
}

type WFContentPredicateTableTemplate struct {
	Value               WFConditionValue `plist:",omitempty"`
	WFSerializationType string           `plist:",omitempty"`
}

type WFConditionValue struct {
	WFActionParameterFilterPrefix    int
	WFActionParameterFilterTemplates []WFConditionParam `plist:",omitempty"`
}

type WFConditionParam struct {
	WFCondition               int
	WFInput                   WFInputVariable `plist:",omitempty"`
	WFConditionalActionString any             `plist:",omitempty"`
	WFNumberValue             any             `plist:",omitempty"`
	WFAnotherNumber           any             `plist:",omitempty"`
	WFDate                    any             `plist:",omitempty"`
	WFAnotherDate             any             `plist:",omitempty"`
	WFDuration                any             `plist:",omitempty"`
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

var uuids map[string]string

type dictDataType int

const itemTypeText dictDataType = 0
const itemTypeNumber dictDataType = 3
const itemTypeArray dictDataType = 2
const itemTypeDict dictDataType = 1
const itemTypeBool dictDataType = 4
const itemTypeFile dictDataType = 5

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
	"installedapp":  "WFAppContentItem",
	"app":           "WFAppStoreAppContentItem",
	"article":       "WFArticleContentItem",
	"bool":          "WFBooleanContentItem",
	"contact":       "WFContactContentItem",
	"date":          "WFDateContentItem",
	"display":       "WFDisplayContentItem",
	"email":         "WFEmailAddressContentItem",
	"eventattendee": "WFEKParticipantContentItem",
	"folder":        "WFFolderContentItem",
	"file":          "WFGenericFileContentItem",
	"image":         "WFImageContentItem",
	"itunes":        "WFiTunesProductContentItem",
	"location":      "WFLocationContentItem",
	"maplink":       "WFDCMapsLinkContentItem",
	"media":         "WFAVAssetContentItem",
	"photomedia":    "WFPhotoMediaContentItem",
	"place":         "WFMKMapItemContentItem",
	"pdf":           "WFPDFContentItem",
	"phonenumber":   "WFPhoneNumberContentItem",
	"richtext":      "WFRichTextContentItem",
	"webpage":       "WFSafariWebPageContentItem",
	"text":          "WFStringContentItem",
	"dictionary":    "WFDictionaryContentItem",
	"number":        "WFNumberContentItem",
	"url":           "WFURLContentItem",
	"vcard":         "WFVCardContentItem",
	"window":        "WFWindowContentItem",
}

var defaultContentItems = []string{
	"installedapp",
	"app",
	"article",
	"contact",
	"date",
	"email",
	"folder",
	"file",
	"image",
	"itunes",
	"location",
	"maplink",
	"media",
	"pdf",
	"phonenumber",
	"richtext",
	"webpage",
	"text",
	"dictionary",
	"number",
	"url",
}

var revContentItems map[string]string

func reversedContentItems() map[string]string {
	if len(revContentItems) == 0 {
		revContentItems = make(map[string]string)
		for key, contentItemType := range contentItems {
			revContentItems[contentItemType] = key
		}
	}

	return revContentItems
}

var globalProperties = map[string]valueList{
	"Device": {
		name:   "device property",
		plural: "device properties",
		list: &[]string{
			"OS",
			"Name",
			"Is Locked",
			"Device Type",
			"Device Is Watch",
			"Device Model",
			"Device Hostname",
			"System Version",
			"System Build Number",
			"Screen Width",
			"Screen Height",
			"Current Volume",
			"Current Brightness",
			"Current Appearance",
		},
	},
	"CurrentApp": {
		name:   "app property",
		plural: "app properties",
		list: &[]string{
			"Is Running",
			"Is Hidden",
			"Is Frontmost",
			"Icon",
			"Bundle Identifier",
			"Process Identifier",
			"Bundle Path",
			"Launch Date",
			"Windows",
			"Displays",
		},
	},
	"CurrentDate": {
		name:   "date property",
		plural: "date properties",
		list: &[]string{
			"Date",
			"Time",
			"Name",
		},
	},
}

var fileProperties = []string{
	"File Size",
	"File Extension",
	"File Path",
	"Date Created",
	"Date Modified",
	"Name",
}

var nameProperties = []string{"Name"}

var labelProperties = []string{"Label", "Name"}

var folderProperties = append([]string{
	"Contents",
	"Number of Items",
}, fileProperties...)

var mediaProperties = append([]string{
	"Duration",
	"Frame Rate",
	"Album",
}, fileProperties...)

var imageProperties = append([]string{
	"Width",
	"Height",
	"Date Taken",
	"Time Taken",
	"Media Type",
	"Photo Type",
	"Is a Screenshot",
	"Is a Screen Recording",
	"Location",
	"Orientation",
	"Camera Make",
	"Camera Model",
	"Metadata Dictionary",
	"Is Favorite",
}, mediaProperties...)

var contactProperties = append([]string{
	"First Name",
	"Middle Name",
	"Last Name",
	"Birthday",
	"Prefix",
	"Suffix",
	"Nickname",
	"Phonetic First Name",
	"Phonetic Last Name",
	"Phonetic Middle Name",
	"Relationship to Me",
	"Contact Photo",
	"Has Photo",
	"URLs",
	"Company",
	"Job Title",
	"Department",
	"Creation Date",
	"Last Modified Date",
	"Group",
}, fileProperties...)

// text, number, bool and date are left out on purpose. Dictionary values and list items are coerced to them, and the
// type of such a reference is not known when parsing, so a key on it is a dictionary key we cannot validate.
var contentTypeProperties = map[string]valueList{
	"installedapp": {list: globalProperties["CurrentApp"].list},
	"app": {
		list: &[]string{
			"Artist",
			"Price",
			"Currency Code",
			"Formatted Price",
			"Release Date",
			"Category",
			"Description",
			"Rating",
			"# of Ratings",
			"Rating (This Version)",
			"# of Ratings (This Version)",
			"Version",
			"Last Updated",
			"Release Notes",
			"Content Rating",
			"Minimum OS Version",
			"Is Universal",
			"Supports Game Center",
			"Supported Devices",
			"Screenshot URLs",
			"iPad Screenshot URLs",
			"Download Size",
			"Bundle Identifier",
			"Store ID",
			"Store URL",
			"Artwork",
			"Artwork URL",
			"Name",
		},
	},
	"article": {
		list: &[]string{
			"Author",
			"Body",
			"Excerpt",
			"Main Image URL",
			"Name",
			"Number of Words",
			"Published Date",
			"Title",
			"URL",
		},
	},
	"contact":       {list: &contactProperties},
	"display":       {list: &[]string{"Is Main Display", "Is Built-in Display", "Scale", "Apps", "Windows"}},
	"email":         {list: &labelProperties},
	"eventattendee": {list: &[]string{"Status", "Role", "Type", "Is Me", "Email Address"}},
	"folder":        {list: &folderProperties},
	"file":          {list: &fileProperties},
	"image":         {list: &imageProperties},
	"itunes": {
		list: &[]string{
			"Artist",
			"Artwork",
			"Artwork URL",
			"Currency Code",
			"Description",
			"Duration",
			"Formatted Price",
			"Genre",
			"Is Explicit",
			"Name",
			"Price",
			"Release Date",
			"Store ID",
			"Store URL",
			"Streamable",
		},
	},
	"location": {
		list: &[]string{
			"Full Street Address",
			"Latitude",
			"Longitude",
			"Altitude",
			"Street",
			"City",
			"State",
			"Zip Code",
			"Region",
			"Phone Number",
		},
	},
	"maplink":     {list: &nameProperties},
	"media":       {list: &mediaProperties},
	"photomedia":  {list: &mediaProperties},
	"place":       {list: &nameProperties},
	"pdf":         {list: &fileProperties},
	"phonenumber": {list: &labelProperties},
	"richtext":    {list: &fileProperties},
	"webpage":     {list: &[]string{"Name", "Page Contents", "Page Selection", "Page URL"}},
	"url":         {list: &nameProperties},
	"vcard":       {list: &fileProperties},
	"window": {list: &[]string{
		"App Name",
		"Is Hidden",
		"Image",
		"Width",
		"Height",
		"X Position",
		"Y Position",
		"Window Index",
		"Display",
		"App",
	}},
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
	"27.2":   "5114.0.8.3",
	"27":     "5037.0.17",
	"26.4":   "4528.0.4.2",
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
var clientVersion = versions["27.2"]
var iosVersion = 27.2

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
	IsToday:        1002,
}

var conditionFilterPrefixes = map[tokenType]int{
	Or:  0,
	And: 1,
}

var allowedConditionalTypes = map[tokenType][]tokenType{
	Is:             {String, Integer, Bool, Action, Date},
	Not:            {String, Integer, Bool, Action, Date},
	Any:            {},
	Empty:          {},
	Contains:       {String, Arr},
	DoesNotContain: {String, Arr},
	BeginsWith:     {String},
	EndsWith:       {String},
	GreaterThan:    {Integer, Float, Quantity, Date},
	GreaterOrEqual: {Integer, Float, Date},
	LessThan:       {Integer, Float, Quantity, Date},
	LessOrEqual:    {Integer, Float, Date},
	Between:        {Integer, Float, Date},
	IsToday:        {},
}

/* Menus */

var menus map[string][]varValue
