/*
 * Copyright (c) Cherri
 */

package main

var definitions map[string]any

var workflowName string

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
var types []string

/* Versions */

var versions = map[string]string{
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
var clientVersion = "3218.0.4.100"
var iosVersion = 18.4

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

type condition struct {
	variableOneType    tokenType
	variableOneValue   any
	condition          int
	variableTwoType    tokenType
	variableTwoValue   any
	variableThreeType  tokenType
	variableThreeValue any
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

/* Menus */

var menus map[string][]varValue
