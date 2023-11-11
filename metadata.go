/*
 * Copyright (c) Brandon Jordan
 */

package main

var definitions map[string]any

var workflowName string

/* Colors */

var colors map[string]string
var iconColor = "-1263359489"

func makeColors() {
	if len(colors) != 0 {
		return
	}
	colors = map[string]string{
		"red":        "4282601983",
		"darkorange": "4251333119",
		"orange":     "4271458815",
		"yellow":     "4274264319",
		"green":      "4292093695",
		"teal":       "431817727",
		"lightblue":  "1440408063",
		"blue":       "463140863",
		"darkblue":   "946986751",
		"violet":     "2071128575",
		"purple":     "3679049983",
		"pink":       "3980825855",
		"taupe":      "3031607807",
		"gray":       "2846468607",
		"darkgray":   "255",
	}
}

/* Inputs */

var contentItems map[string]string
var inputs []string
var outputs []string

func makeContentItems() {
	if len(contentItems) != 0 {
		return
	}
	contentItems = map[string]string{
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
		"number":      "WFNumberContentItem",
		"url":         "WFURLContentItem",
	}
}

/* Workflow Types */

var workflowTypes map[string]string
var types []string

func makeWorkflowTypes() {
	if len(workflowTypes) != 0 {
		return
	}
	workflowTypes = map[string]string{
		"menubar":       "MenuBar",
		"quickactions":  "QuickActions",
		"sharesheet":    "ActionExtension",
		"notifications": "NCWidget",
		"sleepmode":     "Sleep",
		"watch":         "Watch",
		"onscreen":      "ReceivesOnScreenContent",
	}
}

/* Versions */

var versions map[string]string
var minVersion = "900"
var iosVersion = 17.0

func makeVersions() {
	if len(versions) != 0 {
		return
	}
	versions = map[string]string{
		"17":     "2038.0.2.4",
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
}

/* Languages */

var languages map[string]string

func makeLanguages() {
	if len(languages) != 0 {
		return
	}
	languages = map[string]string{
		"Arabic":                      "ar_AE",
		"Mandarin Chinese (Mainland)": "zh_CN",
		"Mandarin Chinese (Taiwan)":   "zh_TW",
		"Dutch":                       "nl_NL",
		"English (UK)":                "en_GB",
		"English (US)":                "en_US",
		"French":                      "fr_FR",
		"German":                      "de_DE",
		"Indonesian":                  "id_ID",
		"Italian":                     "it_IT",
		"Japanese":                    "jp_JP",
		"Korean":                      "ko_KR",
		"Polish":                      "pl_PL",
		"Portuguese (Brazil)":         "pt_BR",
		"Russian":                     "ru_RU",
		"Spanish (Spain)":             "es_ES",
		"Thai":                        "th_TH",
		"Turkish":                     "tr_TR",
		"Vietnamese":                  "vn_VN",
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
	conditions = map[tokenType]string{
		Is:             "4",
		Not:            "5",
		Any:            "100",
		Empty:          "101",
		Contains:       "99",
		DoesNotContain: "999",
		BeginsWith:     "8",
		EndsWith:       "9",
		GreaterThan:    "2",
		GreaterOrEqual: "3",
		LessThan:       "0",
		LessOrEqual:    "1",
		Between:        "1003",
	}
}

/* Menus */

var menus map[string][]variableValue
