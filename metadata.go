/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

var workflowName string

/* Colors */

var colors map[string]string
var iconColor = "-1263359489"

func makeColors() {
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
var minimumVersion = "900"

func makeVersions() {
	versions = make(map[string]string)
	versions["16"] = "900"
	versions["15"] = "800"
	versions["14"] = "700"
	versions["13"] = "600"
	versions["12"] = "500"
}

/* Globals */

var globals map[string]variableValue

func makeGlobals() {
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
