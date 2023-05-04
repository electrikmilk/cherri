/*
 * Copyright (c) Brandon Jordan
 */

package main

var workflowName string

/* Colors */

var colors map[string]string
var iconColor = "-1263359489"

func makeColors() {
	if len(colors) != 0 {
		return
	}
	colors = make(map[string]string)
	colors["red"] = "4282601983"
	colors["darkorange"] = "4251333119"
	colors["orange"] = "4271458815"
	colors["yellow"] = "4274264319"
	colors["green"] = "4292093695"
	colors["teal"] = "431817727"
	colors["lightblue"] = "1440408063"
	colors["blue"] = "463140863"
	colors["darkblue"] = "946986751"
	colors["violet"] = "2071128575"
	colors["purple"] = "3679049983"
	colors["pink"] = "3980825855"
	colors["taupe"] = "3031607807"
	colors["gray"] = "2846468607"
	colors["darkgray"] = "255"
}

/* Inputs */

var contentItems map[string]string
var inputs []string
var outputs []string

func makeContentItems() {
	if len(contentItems) != 0 {
		return
	}
	contentItems = make(map[string]string)
	contentItems["app"] = "WFAppStoreAppContentItem"
	contentItems["article"] = "WFArticleContentItem"
	contentItems["contact"] = "WFContactContentItem"
	contentItems["date"] = "WFDateContentItem"
	contentItems["email"] = "WFEmailAddressContentItem"
	contentItems["folder"] = "WFFolderContentItem"
	contentItems["file"] = "WFGenericFileContentItem"
	contentItems["image"] = "WFImageContentItem"
	contentItems["itunes"] = "WFiTunesProductContentItem"
	contentItems["location"] = "WFLocationContentItem"
	contentItems["maplink"] = "WFDCMapsLinkContentItem"
	contentItems["media"] = "WFAVAssetContentItem"
	contentItems["pdf"] = "WFPDFContentItem"
	contentItems["phonenumber"] = "WFPhoneNumberContentItem"
	contentItems["richtext"] = "WFRichTextContentItem"
	contentItems["webpage"] = "WFSafariWebPageContentItem"
	contentItems["text"] = "WFStringContentItem"
	contentItems["number"] = "WFNumberContentItem"
	contentItems["url"] = "WFURLContentItem"
}

/* Workflow Types */

var workflowTypes map[string]string
var types []string

func makeWorkflowTypes() {
	if len(workflowTypes) != 0 {
		return
	}
	workflowTypes = make(map[string]string)
	workflowTypes["menubar"] = "MenuBar"
	workflowTypes["quickactions"] = "QuickActions"
	workflowTypes["sharesheet"] = "ActionExtension"
	workflowTypes["notifications"] = "NCWidget"
	workflowTypes["sleepmode"] = "Sleep"
	workflowTypes["watchkit"] = "WatchKit"
	workflowTypes["watch"] = "Watch"
	workflowTypes["onscreen"] = "ReceivesOnScreenContent"
}

/* Versions */

var versions map[string]string
var minVersion = "900"
var iosVersion = 16.2

func makeVersions() {
	if len(versions) != 0 {
		return
	}
	versions = make(map[string]string)
	versions["16.4"] = "900"
	versions["16.3"] = "900"
	versions["16.2"] = "900"
	versions["16"] = "900"
	versions["15.7.2"] = "800"
	versions["15"] = "800"
	versions["14"] = "700"
	versions["13"] = "600"
	versions["12"] = "500"
}

/* Languages */

var languages map[string]string

func makeLanguages() {
	if len(languages) != 0 {
		return
	}
	languages = make(map[string]string)
	languages["Arabic"] = "ar_AE"
	languages["Mandarin Chinese (Mainland)"] = "zh_CN"
	languages["Mandarin Chinese (Taiwan)"] = "zh_TW"
	languages["Dutch"] = "nl_NL"
	languages["English (UK)"] = "en_GB"
	languages["English (US)"] = "en_US"
	languages["French"] = "fr_FR"
	languages["German"] = "de_DE"
	languages["Indonesian"] = "id_ID"
	languages["Italian"] = "it_IT"
	languages["Japanese"] = "jp_JP"
	languages["Korean"] = "ko_KR"
	languages["Polish"] = "pl_PL"
	languages["Portuguese (Brazil)"] = "pt_BR"
	languages["Russian"] = "ru_RU"
	languages["Spanish (Spain)"] = "es_ES"
	languages["Thai"] = "th_TH"
	languages["Turkish"] = "tr_TR"
	languages["Vietnamese"] = "vn_VN"
}

/* Import Questions */

var questions map[string]question

type question struct {
	parameter    string
	actionIndex  int
	text         string
	defaultValue string
}

/* Variables */

var variables map[string]variableValue

type variableValue struct {
	variableType string
	valueType    tokenType
	value        any
	getAs        string
	coerce       string
}

/* Globals */

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

/* Conditionals */

type condition struct {
	variableOneType    tokenType
	variableOneValue   any
	condition          string
	variableTwoType    tokenType
	variableTwoValue   any
	variableThreeType  tokenType
	variableThreeValue any
}

var conditions map[tokenType]string

func makeConditions() {
	if len(conditions) != 0 {
		return
	}
	conditions = make(map[tokenType]string)
	conditions[Is] = "4"
	conditions[Not] = "5"
	conditions[Any] = "100"
	conditions[Empty] = "101"
	conditions[Contains] = "99"
	conditions[DoesNotContain] = "999"
	conditions[BeginsWith] = "8"
	conditions[EndsWith] = "9"
	conditions[GreaterThan] = "2"
	conditions[GreaterOrEqual] = "3"
	conditions[LessThan] = "0"
	conditions[LessOrEqual] = "1"
	conditions[Between] = "1003"
}

/* Menus */

var menus map[string][]variableValue
