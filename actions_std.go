/*
 * Copyright (c) 2023 Brandon Jordan
 */

package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// FIXME: Some of these actions have enumerable arguments (a set list values),
//  but the argument value is not being checked against it's possible values.
//  Use the "hash" action as an example.

func standardActions() {
	if len(actions) != 0 {
		return
	}
	calendarActions()
	contactActions()
	documentActions()
	locationActions()
	mediaActions()
	scriptingActions()
	sharingActions()
	webActions()
	customActions()
}

func calendarActions() {
	actions["date"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFDateActionMode",
					dataType: Text,
					value:    "Specified Date",
				},
				argumentValue("WFDateActionDate", args, 0),
			}
		},
	}
	actions["addCalendar"] = actionDefinition{
		identifier: "addnewcalendar",
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
				key:       "CalendarName",
			},
		},
	}
	actions["addSeconds"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Add", "sec", args)
		},
	}
	actions["addMinutes"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Add", "min", args)
		},
	}
	actions["addHours"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Add", "hr", args)
		},
	}
	actions["addDays"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Add", "days", args)
		},
	}
	actions["addWeeks"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Add", "weeks", args)
		},
	}
	actions["addMonths"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Add", "months", args)
		},
	}
	actions["addYears"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Add", "yr", args)
		},
	}
	actions["subtractSeconds"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "sec", args)
		},
	}
	actions["subtractMinutes"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "min", args)
		},
	}
	actions["subtractHours"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "hr", args)
		},
	}
	actions["subtractDays"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "days", args)
		},
	}
	actions["subtractWeeks"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "weeks", args)
		},
	}
	actions["subtractMonths"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "months", args)
		},
	}
	actions["subtractYears"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "magnitude",
				validType: Integer,
			},
			{
				name:      "date",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "yr", args)
		},
	}
	actions["getStartMinute"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Get Start of Minute", "", args)
		},
	}
	actions["getStartHour"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Get Start of Hour", "", args)
		},
	}
	actions["getStartWeek"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Get Start of Week", "", args)
		},
	}
	actions["getStartMonth"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Get Start of Month", "", args)
		},
	}
	actions["getStartYear"] = actionDefinition{
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return adjustDate("Get Start of Year", "", args)
		},
	}
}

func contactActions() {
	var contactProperties = []string{
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
		"Company",
		"Job Title",
		"Department",
		"File Extension",
		"Creation Date",
		"File Path",
		"Last Modified Date",
		"Name",
		"Random",
	}
	var abcSortOrders = []string{"A to Z", "Z to A"}
	actions["filterContacts"] = actionDefinition{
		identifier: "filter.contacts",
		parameters: []parameterDefinition{
			{
				name:      "contacts",
				validType: Variable,
				key:       "WFContentItemInputParameter",
				optional:  false,
			},
			{
				name:      "property",
				validType: String,
				key:       "WFContentItemSortProperty",
				optional:  false,
			},
			{
				name:      "sortOrder",
				validType: String,
				key:       "WFContentItemSortOrder",
				defaultValue: actionArgument{
					valueType: String,
					value:     "A to Z",
				},
				optional: true,
			},
			{
				name:      "limit",
				validType: Integer,
				key:       "WFContentItemLimitNumber",
				optional:  true,
			},
		},
		check: func(args []actionArgument) {
			checkEnum("contact property", contactProperties, args, 1)
			checkEnum("sort order", abcSortOrders, args, 2)
		},
	}
	actions["emailAddress"] = actionDefinition{
		identifier: "email",
		parameters: []parameterDefinition{
			{
				name:      "email",
				validType: String,
				infinite:  true,
			},
		},
		make: func(args []actionArgument) []plistData {
			return contactValue("WFEmailAddress", "emailaddress", args)
		},
	}
	actions["phoneNumber"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: String,
				infinite:  true,
			},
		},
		make: func(args []actionArgument) []plistData {
			return contactValue("WFPhoneNumber", "phonenumber", args)
		},
	}
	actions["selectContact"] = actionDefinition{
		identifier: "selectcontacts",
		parameters: []parameterDefinition{
			{
				name:      "multiple",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
				key: "WFSelectMultiple",
			},
		},
	}
	actions["selectEmailAddress"] = actionDefinition{
		identifier: "selectemail",
	}
	actions["selectPhoneNumber"] = actionDefinition{
		identifier: "selectphone",
	}
	actions["getContactDetail"] = actionDefinition{
		identifier: "properties.contacts",
		parameters: []parameterDefinition{
			{
				name:      "contact",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "property",
				validType: String,
				key:       "WFContentItemPropertyName",
			},
		},
		check: func(args []actionArgument) {
			checkEnum("contact property", contactProperties, args, 1)
		},
	}
}

func documentActions() {
	// FIXME: Writing to locations other than the Shortcuts folder
	actions["createFolder"] = actionDefinition{
		identifier: "file.createfolder",
		parameters: []parameterDefinition{
			{
				name:      "path",
				validType: String,
				key:       "WFFilePath",
			},
		},
	}
	actions["getFolderContents"] = actionDefinition{
		identifier: "file.getfoldercontents",
		parameters: []parameterDefinition{
			{
				name:      "folder",
				validType: Variable,
				key:       "WFFolder",
			},
			{
				key:       "Recursive",
				name:      "recursive",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     true,
				},
				optional: true,
			},
		},
	}
	actions["matchedTextGroupIndex"] = actionDefinition{
		identifier: "text.match.getgroup",
		parameters: []parameterDefinition{
			{
				name:      "matches",
				validType: Variable,
				key:       "matches",
			},
			{
				name:      "index",
				validType: Integer,
				key:       "WFGroupIndex",
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("matches", args[0].value.(string)),
				argumentValue("WFGroupIndex", args, 1),
				{
					key:      "WFGetGroupType",
					dataType: Text,
					value:    "Group At Index",
				},
			}
		},
	}
	actions["getFileFromFolder"] = actionDefinition{
		identifier: "documentpicker.open",
		parameters: []parameterDefinition{
			{
				name:      "folder",
				validType: Variable,
				key:       "WFFile",
			},
			{
				name:      "path",
				validType: String,
				key:       "WFGetFilePath",
			},
			{
				name:      "errorIfNotFound",
				validType: Bool,
				key:       "WFFileErrorIfNotFound",
				defaultValue: actionArgument{
					valueType: Bool,
					value:     true,
				},
				optional: true,
			},
		},
	}
	actions["getFile"] = actionDefinition{
		identifier: "documentpicker.open",
		parameters: []parameterDefinition{
			{
				name:      "path",
				validType: String,
				key:       "WFGetFilePath",
			},
			{
				name:      "errorIfNotFound",
				validType: Bool,
				key:       "WFFileErrorIfNotFound",
				defaultValue: actionArgument{
					valueType: Bool,
					value:     true,
				},
				optional: true,
			},
		},
	}
	actions["markup"] = actionDefinition{
		identifier: "avairyeditphoto",
		parameters: []parameterDefinition{
			{
				name:      "document",
				validType: Variable,
				key:       "WFDocument",
			},
		},
	}
	actions["rename"] = actionDefinition{
		identifier: "file.rename",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFFile",
			},
			{
				name:      "newName",
				validType: String,
				key:       "WFNewFilename",
			},
		},
	}
	actions["reveal"] = actionDefinition{
		identifier: "file.reveal",
		parameters: []parameterDefinition{
			{
				name:      "files",
				validType: Variable,
				key:       "WFFile",
			},
		},
	}
	actions["define"] = actionDefinition{
		identifier: "showdefinition",
		parameters: []parameterDefinition{
			{
				name:      "word",
				validType: String,
				key:       "Word",
			},
		},
	}
	var errorCorrectionLevels = []string{"low", "medium", "quartile", "high"}
	actions["makeQRcode"] = actionDefinition{
		identifier: "generatebarcode",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
				key:       "WFText",
			},
			{
				name:      "errorCorrection",
				validType: String,
				key:       "WFQRErrorCorrectionLevel",
			},
		},
		check: func(args []actionArgument) {
			checkEnum("error correction level", errorCorrectionLevels, args, 0)
		},
	}
	actions["showNote"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "note",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["splitPDF"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "pdf",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["makeHTML"] = actionDefinition{
		identifier: "gethtmlfromrichtext",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "makeFullDocument",
				validType: Bool,
				key:       "WFMakeFullDocument",
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
				optional: true,
			},
		},
	}
	actions["makeMarkdown"] = actionDefinition{
		identifier: "getmarkdownfromrichtext",
		parameters: []parameterDefinition{
			{
				name:      "richText",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["getRichTextFromHTML"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "html",
				validType: Variable,
				key:       "WFHTML",
			},
		},
	}
	actions["getRichTextFromMarkdown"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "markdown",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["print"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["selectFile"] = actionDefinition{
		identifier: "file.select",
		parameters: []parameterDefinition{
			{
				name:      "multiple",
				validType: Bool,
				key:       "SelectMultiple",
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
				optional: true,
			},
		},
	}
	actions["getFileLink"] = actionDefinition{
		identifier: "file.getlink",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFFile",
			},
		},
	}
	actions["getParentDirectory"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["getEmojiName"] = actionDefinition{
		identifier: "getnameofemoji",
		parameters: []parameterDefinition{
			{
				name:      "emoji",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["getFileDetail"] = actionDefinition{
		identifier: "properties.files",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFFolder",
			},
			{
				name:      "detail",
				validType: String,
				key:       "WFContentItemPropertyName",
			},
		},
	}
	actions["deleteFiles"] = actionDefinition{
		identifier: "file.delete",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
			{
				name:      "immediately",
				key:       "WFDeleteImmediatelyDelete",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
				optional: true,
			},
		},
	}
	actions["getTextFromImage"] = actionDefinition{
		identifier: "extracttextfromimage",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFImage",
			},
		},
	}
	actions["connectToServer"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["appendNote"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "note",
				validType: String,
				key:       "WFNote",
			},
			{
				name:      "input",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["addToBooks"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "BooksInput",
			},
		},
	}
	actions["saveFile"] = actionDefinition{
		identifier: "documentpicker.save",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
			},
			{
				name:      "path",
				validType: String,
			},
			{
				name:      "overwrite",
				validType: Bool,
				key:       "WFSaveFileOverwrite",
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
				optional: true,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
				argumentValue("WFFileDestinationPath", args, 1),
				argumentValue("WFSaveFileOverwrite", args, 2),
				{
					key:      "WFAskWhereToSave",
					dataType: Boolean,
					value:    false,
				},
			}
		},
	}
	actions["saveFilePrompt"] = actionDefinition{
		identifier: "documentpicker.save",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "overwrite",
				validType: Bool,
				key:       "WFSaveFileOverwrite",
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
				optional: true,
			},
		},
	}
	actions["getSelectedFiles"] = actionDefinition{identifier: "finder.getselectedfiles"}
	actions["extractArchive"] = actionDefinition{
		identifier: "unzip",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFArchive",
			},
		},
	}
	actions["makeArchive"] = actionDefinition{
		identifier: "makezip",
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
				key:       "WFZIPName",
			},
			{
				name:      "format",
				validType: String,
				key:       "WFArchiveFormat",
			},
			{
				name:      "files",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["quicklook"] = actionDefinition{
		identifier: "previewdocument",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["translateFrom"] = actionDefinition{
		identifier: "text.translate",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
				key:       "WFInputText",
			},
			{
				name:      "from",
				validType: String,
				key:       "WFSelectedFromLanguage",
			},
			{
				name:      "to",
				validType: String,
				key:       "WFSelectedLanguage",
			},
		},
		check: func(args []actionArgument) {
			args[1].value = languageCode(args[1].value.(string))
			args[2].value = languageCode(args[2].value.(string))
		},
	}
	actions["translate"] = actionDefinition{
		identifier: "text.translate",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
			},
			{
				name:      "to",
				validType: String,
			},
		},
		check: func(args []actionArgument) {
			args[1].value = languageCode(args[1].value.(string))
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInputText", args, 0),
				argumentValue("WFSelectedLanguage", args, 1),
				{
					key:      "WFSelectedFromLanguage",
					dataType: Text,
					value:    "Detect Language",
				},
			}
		},
	}
	actions["detectLanguage"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["replaceText"] = actionDefinition{
		identifier: "text.replace",
		parameters: []parameterDefinition{
			{
				name:      "find",
				validType: String,
			},
			{
				name:      "replacement",
				validType: String,
			},
			{
				name:      "subject",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return replaceText(true, false, args)
		},
	}
	actions["iReplaceText"] = actionDefinition{
		identifier: "text.replace",
		parameters: []parameterDefinition{
			{
				name:      "find",
				validType: String,
			},
			{
				name:      "replacement",
				validType: String,
			},
			{
				name:      "subject",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return replaceText(false, false, args)
		},
	}
	actions["regReplaceText"] = actionDefinition{
		identifier: "text.replace",
		parameters: []parameterDefinition{
			{
				name:      "expression",
				validType: String,
			},
			{
				name:      "replacement",
				validType: String,
			},
			{
				name:      "subject",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return replaceText(true, true, args)
		},
	}
	actions["iRegReplaceText"] = actionDefinition{
		identifier: "text.replace",
		parameters: []parameterDefinition{
			{
				name:      "expression",
				validType: String,
			},
			{
				name:      "replacement",
				validType: String,
			},
			{
				name:      "subject",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return replaceText(false, true, args)
		},
	}
	actions["uppercase"] = actionDefinition{
		identifier: "text.changecase",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return changeCase("UPPERCASE", args)
		},
	}
	actions["lowercase"] = actionDefinition{
		identifier: "text.changecase",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return changeCase("lowercase", args)
		},
	}
	actions["titleCase"] = actionDefinition{
		identifier: "text.changecase",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return changeCase("Capitalize with Title Case", args)
		},
	}
	actions["capitalize"] = actionDefinition{
		identifier: "text.changecase",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return changeCase("Capitalize with sentence case", args)
		},
	}
	actions["capitalizeAll"] = actionDefinition{
		identifier: "text.changecase",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return changeCase("Capitalize Every Word", args)
		},
	}
	actions["alternateCase"] = actionDefinition{
		identifier: "text.changecase",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return changeCase("cApItAlIzE wItH aLtErNaTiNg cAsE", args)
		},
	}
	actions["correctSpelling"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Show-text",
					dataType: Boolean,
					value:    true,
				},
				argumentValue("text", args, 0),
			}
		},
	}
	actions["splitText"] = actionDefinition{
		identifier: "text.split",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
			},
			{
				name:      "glue",
				validType: String,
			},
		},
		make: textParts,
	}
	actions["combineText"] = actionDefinition{
		identifier: "text.combine",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
			},
			{
				name:      "glue",
				validType: String,
			},
		},
		make: textParts,
	}
	actions["makeDiskImage"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
			},
			{
				name:      "contents",
				validType: Variable,
			},
			{
				name:      "encrypt",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
				optional: true,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("VolumeName", args, 0),
				variableInput("WFInput", args[1].value.(string)),
				argumentValue("EncryptImage", args, 1),
				{
					key:      "SizeToFit",
					dataType: Boolean,
					value:    true,
				},
			}
		},
		mac:        true,
		minVersion: 15,
	}
	var storageUnits = []string{
		"bytes",
		"KB",
		"MB",
		"GB",
		"TB",
		"PB",
		"EB",
		"ZB",
		"YB",
	}
	actions["makeSizedDiskImage"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
			},
			{
				name:      "contents",
				validType: Variable,
			},
			{
				name:      "size",
				validType: String,
			},
			{
				name:      "encrypt",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
				optional: true,
			},
		},
		check: func(args []actionArgument) {
			var copyArgs = args
			var size = strings.Split(getArgValue(args[2]).(string), " ")
			copyArgs[2] = actionArgument{
				valueType: String,
				value:     size[1],
			}
			checkEnum("disk size", storageUnits, copyArgs, 2)
		},
		make: func(args []actionArgument) []plistData {
			var size = strings.Split(getArgValue(args[2]).(string), " ")
			return []plistData{
				argumentValue("VolumeName", args, 0),
				variableInput("WFInput", args[1].value.(string)),
				{
					key:      "ImageSize",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "Value",
							dataType: Dictionary,
							value: []plistData{
								{
									key:      "Unit",
									dataType: Text,
									value:    size[0],
								},
								{
									key:      "Magnitude",
									dataType: Text,
									value:    size[1],
								},
							},
						},
						{
							key:      "WFSerializationType",
							dataType: Text,
							value:    "WFQuantityFieldValue",
						},
					},
				},
				argumentValue("EncryptImage", args, 3),
				{
					key:      "SizeToFit",
					dataType: Boolean,
					value:    false,
				},
			}
		},
		mac:        true,
		minVersion: 15,
	}
	actions["openFile"] = actionDefinition{
		identifier: "openin",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "prompt",
				validType: Bool,
				key:       "WFOpenInAskWhenRun",
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
				optional: true,
			},
		},
	}
}

func locationActions() {
	actions["getCurrentLocation"] = actionDefinition{
		identifier: "location",
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFLocation",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "isCurrentLocation",
							dataType: Boolean,
							value:    true,
						},
					},
				},
			}
		},
	}
	actions["getAddresses"] = actionDefinition{
		identifier: "detect.address",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["getCurrentWeather"] = actionDefinition{
		identifier: "currentconditions",
	}
	actions["getCurrentWeatherAt"] = actionDefinition{
		identifier: "currentconditions",
		parameters: []parameterDefinition{
			{
				name:      "location",
				validType: Variable,
				key:       "WFWeatherCustomLocation",
			},
		},
	}
	actions["openInMaps"] = actionDefinition{
		identifier: "searchmaps",
		parameters: []parameterDefinition{
			{
				name:      "location",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["streetAddress"] = actionDefinition{
		identifier: "address",
		parameters: []parameterDefinition{
			{
				name:      "addressLine2",
				validType: String,
				key:       "WFAddressLine1",
			},
			{
				name:      "addressLine2",
				validType: String,
				key:       "WFAddressLine2",
			},
			{
				name:      "city",
				validType: String,
				key:       "WFCity",
			},
			{
				name:      "state",
				validType: String,
				key:       "WFState",
			},
			{
				name:      "country",
				validType: String,
				key:       "WFCountry",
			},
			{
				name:      "zipCode",
				validType: Integer,
				key:       "WFPostalCode",
			},
		},
	}
	var weatherDetails = []string{
		"Name",
		"Air Pollutants",
		"Air Quality Category",
		"Air Quality Index",
		"Sunset Time",
		"Sunrise Time",
		"UV Index",
		"Wind Direction",
		"Wind Speed",
		"Precipitation Chance",
		"Precipitation Amount",
		"Pressure",
		"Humidity",
		"Dewpoint",
		"Visibility",
		"Condition",
		"Feels Like",
		"Low",
		"High",
		"Temperature",
		"Location",
		"Date",
	}
	actions["getWeatherDetail"] = actionDefinition{
		identifier: "properties.weather.conditions",
		parameters: []parameterDefinition{
			{
				name:      "weather",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "detail",
				validType: String,
				key:       "WFContentItemPropertyName",
			},
		},
		check: func(args []actionArgument) {
			checkEnum("weather detail", weatherDetails, args, 1)
		},
	}
	var weatherForecastTypes = []string{
		"Daily",
		"Hourly",
	}
	actions["getWeatherForecast"] = actionDefinition{
		identifier: "weather.forecast",
		parameters: []parameterDefinition{
			{
				name:      "type",
				validType: String,
				key:       "WFWeatherForecastType",
			},
		},
		check: func(args []actionArgument) {
			checkEnum("weather forecast type", weatherForecastTypes, args, 0)
		},
	}
	actions["getWeatherForecastAt"] = actionDefinition{
		identifier: "weather.forecast",
		parameters: []parameterDefinition{
			{
				name:      "type",
				validType: String,
				key:       "WFWeatherForecastType",
			},
			{
				name:      "location",
				validType: Variable,
				key:       "WFInput",
			},
		},
		check: func(args []actionArgument) {
			checkEnum("weather forecast type", weatherForecastTypes, args, 0)
		},
	}
	var locationDetails = []string{
		"Name",
		"URL",
		"Label",
		"Phone Number",
		"Region",
		"ZIP Code",
		"State",
		"City",
		"Street",
		"Altitude",
		"Longitude",
		"Latitude",
	}
	actions["getLocationDetail"] = actionDefinition{
		identifier: "properties.locations",
		parameters: []parameterDefinition{
			{
				name:      "location",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "detail",
				validType: String,
				key:       "WFContentItemPropertyName",
			},
		},
		check: func(args []actionArgument) {
			checkEnum("location detail", locationDetails, args, 1)
		},
	}
	actions["getMapsLink"] = actionDefinition{
		identifier: "",
		parameters: []parameterDefinition{
			{
				name:      "location",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["getHalfwayPoint"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "firstLocation",
				validType: Variable,
				key:       "WFGetHalfwayPointFirstLocation",
			},
			{
				name:      "secondLocation",
				validType: Variable,
				key:       "WFGetHalfwayPointSecondLocation",
			},
		},
	}
}

func mediaActions() {
	actions["clearUpNext"] = actionDefinition{}
	actions["getCurrentSong"] = actionDefinition{}
	actions["latestPhotoImport"] = actionDefinition{identifier: "getlatestphotoimport"}
	actions["takePhoto"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "showPreview",
				validType: Bool,
				key:       "WFCameraCaptureShowPreview",
				defaultValue: actionArgument{
					valueType: Bool,
					value:     true,
				},
			},
		},
	}
	actions["takePhotos"] = actionDefinition{
		identifier: "takephoto",
		parameters: []parameterDefinition{
			{
				name:      "count",
				validType: Integer,
			},
		},
		check: func(args []actionArgument) {
			var photos = args[0].value.(int)
			if photos == 0 {
				parserError("Number of photos to take must be greater than zero.")
			}
			if photos < 2 {
				parserError("Use action takePhoto() to take only one photo.")
			}
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFPhotoCount", args, 0),
				{
					key:      "WFCameraCaptureShowPreview",
					dataType: Boolean,
					value:    true,
				},
			}
		},
	}
	actions["trimVideo"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "video",
				validType: Variable,
				key:       "WFInputMedia",
			},
		},
	}
	actions["takeVideo"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "camera",
				validType: String,
				key:       "WFCameraCaptureDevice",
				defaultValue: actionArgument{
					valueType: String,
					value:     "Front",
				},
			},
			{
				name:      "quality",
				validType: String,
				key:       "WFCameraCaptureQuality",
				defaultValue: actionArgument{
					valueType: String,
					value:     "Medium",
				},
			},
			{
				name:      "startImmediately",
				validType: Bool,
				key:       "WFRecordingStart",
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
			},
		},
	}
	actions["setVolume"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "volume",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			args[0].value = fmt.Sprintf("0.%s", args[0].value)
			return []plistData{
				argumentValue("WFVolume", args, 0),
			}
		},
	}
}

func scriptingActions() {
	actions["getObjectOfClass"] = actionDefinition{
		identifier: "getclassaction",
		parameters: []parameterDefinition{
			{
				name:      "class",
				key:       "Class",
				validType: String,
			},
			{
				name:      "from",
				key:       "Input",
				validType: Variable,
			},
		},
	}
	actions["getOnScreenContent"] = actionDefinition{}
	actions["fileSize"] = actionDefinition{
		identifier: "format.filesize",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
			},
			{
				name:      "format",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFFileSizeIncludeUnits",
					dataType: Boolean,
					value:    false,
				},
				argumentValue("WFFileSize", args, 0),
				argumentValue("WFFileSizeFormat", args, 1),
			}
		},
	}
	actions["getDeviceDetail"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "detail",
				key:       "WFDeviceDetail",
				validType: String,
			},
		},
	}
	actions["setBrightness"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "brightness",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			args[0].value = fmt.Sprintf("0.%s", args[0].value)
			return []plistData{
				argumentValue("WFBrightness", args, 0),
			}
		},
	}
	actions["getName"] = actionDefinition{
		identifier: "getitemname",
		parameters: []parameterDefinition{
			{
				name:      "item",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["setName"] = actionDefinition{
		identifier: "setitemname",
		parameters: []parameterDefinition{
			{
				name:      "item",
				key:       "WFInput",
				validType: Variable,
			},
			{
				name:      "name",
				key:       "WFName",
				validType: String,
			},
			{
				name:      "includeFileExtension",
				key:       "WFDontIncludeFileExtension",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
			},
		},
	}
	actions["countItems"] = actionDefinition{
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return count("Items", args)
		},
	}
	actions["countChars"] = actionDefinition{
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return count("Characters", args)
		},
	}
	actions["countWords"] = actionDefinition{
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return count("Words", args)
		},
	}
	actions["countSentences"] = actionDefinition{
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return count("Sentences", args)
		},
	}
	actions["countLines"] = actionDefinition{
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return count("Lines", args)
		},
	}
	actions["toggleAppearance"] = actionDefinition{
		identifier: "appearance",
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "operation",
					dataType: Text,
					value:    "toggle",
				},
			}
		},
	}
	actions["lightMode"] = actionDefinition{
		identifier: "appearance",
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "operation",
					dataType: Text,
					value:    "set",
				},
				{
					key:      "style",
					dataType: Text,
					value:    "light",
				},
			}
		},
	}
	actions["darkMode"] = actionDefinition{
		identifier: "appearance",
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "operation",
					dataType: Text,
					value:    "set",
				},
				{
					key:      "style",
					dataType: Text,
					value:    "dark",
				},
			}
		},
	}
	actions["getBatteryLevel"] = actionDefinition{}
	actions["isCharging"] = actionDefinition{
		identifier: "getbatterylevel",
		minVersion: 16.2,
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Subject",
					dataType: Text,
					value:    "Is Charging",
				},
			}
		},
	}
	actions["connectedToCharger"] = actionDefinition{
		identifier: "getbatterylevel",
		minVersion: 16.2,
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Subject",
					dataType: Text,
					value:    "Is Connected to Charger",
				},
			}
		},
	}
	actions["getShortcuts"] = actionDefinition{
		identifier: "getmyworkflows",
	}
	actions["url"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			var urlItems []string
			for _, item := range args {
				urlItems = append(urlItems, plistValue(Text, item.value))
			}
			return []plistData{
				{
					key:      "Show-WFURLActionURL",
					dataType: Boolean,
					value:    true,
				},
				{
					key:      "WFURLActionURL",
					dataType: Array,
					value:    urlItems,
				},
			}
		},
	}
	actions["addToReadingList"] = actionDefinition{
		identifier: "readinglist",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Show-WFURL",
					dataType: Boolean,
					value:    true,
				},
				argumentValue("WFURL", args, 0),
			}
		},
	}
	var hashTypes = []string{"md5", "sha1", "sha256", "sha512"}
	actions["hash"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "type",
				key:       "WFHashType",
				validType: String,
				defaultValue: actionArgument{
					valueType: "MD5",
					value:     String,
				},
				optional: true,
			},
		},
		check: func(args []actionArgument) {
			checkEnum("hash type", hashTypes, args, 1)
			if args[1].value != nil && args[1].valueType != Variable {
				args[1].value = strings.ToUpper(args[1].value.(string))
			}
		},
	}
	actions["formatNumber"] = actionDefinition{
		identifier: "format.number",
		parameters: []parameterDefinition{
			{
				name:      "number",
				key:       "WFNumber",
				validType: Integer,
			},
			{
				name:      "decimalPlaces",
				key:       "WFNumberFormatDecimalPlaces",
				validType: Integer,
			},
		},
	}
	actions["randomNumber"] = actionDefinition{
		identifier: "number.random",
		parameters: []parameterDefinition{
			{
				name:      "min",
				key:       "WFRandomNumberMinimum",
				validType: Integer,
			},
			{
				name:      "max",
				key:       "WFRandomNumberMaximum",
				validType: Integer,
			},
		},
	}
	actions["base64Encode"] = actionDefinition{
		identifier: "base64encode",
		parameters: []parameterDefinition{
			{
				name:      "encodeInput",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "input",
					dataType: Text,
					value:    "Encode",
				},
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["base64Decode"] = actionDefinition{
		identifier: "base64encode",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFEncodeMode",
					dataType: Text,
					value:    "Decode",
				},
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["urlEncode"] = actionDefinition{
		identifier: "urlencode",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFEncodeMode",
					dataType: Text,
					value:    "Encode",
				},
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["urlDecode"] = actionDefinition{
		identifier: "urlencode",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFEncodeMode",
					dataType: Text,
					value:    "Decode",
				},
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["show"] = actionDefinition{
		identifier: "showresult",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Text",
				validType: Variable,
			},
		},
	}
	actions["waitToReturn"] = actionDefinition{}
	actions["notification"] = actionDefinition{
		identifier: "notification",
		parameters: []parameterDefinition{
			{
				name:      "body",
				validType: String,
				key:       "WFNotificationActionBody",
			},
			{
				name:      "title",
				key:       "WFNotificationActionTitle",
				validType: String,
				optional:  true,
			},
			{
				name:      "playSound",
				key:       "WFNotificationActionSound",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     true,
				},
			},
		},
	}
	actions["stop"] = actionDefinition{
		identifier: "exit",
	}
	actions["nothing"] = actionDefinition{}
	actions["wait"] = actionDefinition{
		identifier: "delay",
		parameters: []parameterDefinition{
			{
				name:      "seconds",
				key:       "WFDelayTime",
				validType: Integer,
			},
		},
	}
	actions["alert"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "alert",
				key:       "WFAlertActionMessage",
				validType: String,
			},
			{
				name:      "title",
				key:       "WFAlertActionTitle",
				validType: String,
				optional:  true,
			},
			{
				name:      "cancelButton",
				key:       "WFAlertActionCancelButtonShown",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     true,
				},
			},
		},
	}
	actions["askForInput"] = actionDefinition{
		identifier: "ask",
		parameters: []parameterDefinition{
			{
				name:      "inputType",
				validType: String,
				key:       "WFInputType",
			},
			{
				name:      "prompt",
				validType: String,
				key:       "WFAskActionPrompt",
			},
			{
				name:      "defaultValue",
				validType: String,
				optional:  true,
				key:       "WFAskActionDefaultAnswer",
			},
		},
	}
	actions["chooseFromList"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "dictionary",
				key:       "WFInput",
				validType: Dict,
			},
			{
				name:      "prompt",
				key:       "WFChooseFromListActionPrompt",
				validType: String,
			},
		},
	}
	actions["chooseMultipleFromList"] = actionDefinition{
		identifier: "choosefromlist",
		parameters: []parameterDefinition{
			{
				name:      "dictionary",
				validType: Dict,
			},
			{
				name:      "prompt",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFChooseFromListActionSelectMultiple",
					dataType: Boolean,
					value:    true,
				},
				argumentValue("WFInput", args, 0),
				argumentValue("WFChooseFromListActionPrompt", args, 1),
				argumentValue("WFChooseFromListActionSelectAll", args, 2),
			}
		},
	}
	actions["getType"] = actionDefinition{
		identifier: "getitemtype",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["getKeys"] = actionDefinition{
		identifier: "getvalueforkey",
		parameters: []parameterDefinition{
			{
				name:      "dictionary",
				validType: Dict,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetDictionaryValueType",
					dataType: Text,
					value:    "All Keys",
				},
				argumentValue("WFInput", args, 0),
			}
		},
	}
	actions["getValues"] = actionDefinition{
		identifier: "getvalueforkey",
		parameters: []parameterDefinition{
			{
				name:      "dictionary",
				validType: Dict,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetDictionaryValueType",
					dataType: Text,
					value:    "All Values",
				},
				argumentValue("WFInput", args, 0),
			}
		},
	}
	actions["getValue"] = actionDefinition{
		identifier: "getvalueforkey",
		parameters: []parameterDefinition{
			{
				name:      "dictionary",
				validType: Dict,
			},
			{
				name:      "key",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetDictionaryValueType",
					dataType: Text,
					value:    "Value",
				},
				argumentValue("WFInput", args, 0),
				argumentValue("WFDictionaryKey", args, 1),
			}
		},
	}
	actions["setValue"] = actionDefinition{
		identifier: "setvalueforkey",
		parameters: []parameterDefinition{
			{
				name:      "dictionary",
				validType: Variable,
				key:       "WFDictionary",
			},
			{
				name:      "key",
				validType: String,
				key:       "WFDictionaryKey",
			},
			{
				name:      "value",
				validType: String,
				key:       "WFDictionaryValue",
			},
		},
	}
	actions["openApp"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "appID",
				validType: String,
			},
		},
		check: replaceAppId,
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFAppIdentifier", args, 0),
				{
					key:      "WFSelectedApp",
					dataType: Dictionary,
					value: []plistData{
						argumentValue("BundleIdentifier", args, 0),
					},
				},
			}
		},
	}
	actions["hideApp"] = actionDefinition{
		identifier: "hide.app",
		parameters: []parameterDefinition{
			{
				name:      "appId",
				validType: String,
			},
		},
		check: replaceAppId,
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFApp",
					dataType: Dictionary,
					value: []plistData{
						argumentValue("BundleIdentifier", args, 0),
					},
				},
			}
		},
	}
	actions["quitApp"] = actionDefinition{
		identifier: "quit.app",
		parameters: []parameterDefinition{
			{
				name:      "appId",
				validType: String,
			},
		},
		check: replaceAppId,
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFApp",
					dataType: Dictionary,
					value: []plistData{
						argumentValue("BundleIdentifier", args, 0),
					},
				},
			}
		},
	}
	actions["killApp"] = actionDefinition{
		identifier: "quit.app",
		parameters: []parameterDefinition{
			{
				name:      "appId",
				validType: String,
			},
		},
		check: replaceAppId,
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFApp",
					dataType: Dictionary,
					value: []plistData{
						argumentValue("BundleIdentifier", args, 0),
					},
				},
				{
					key:      "WFAskToSaveChanges",
					dataType: Boolean,
					value:    false,
				},
			}
		},
	}
	actions["open"] = actionDefinition{
		identifier: "openworkflow",
		parameters: []parameterDefinition{
			{
				name:      "shortcutName",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFWorkflow",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "workflowIdentifier",
							dataType: Text,
							value:    shortcutsUUID(),
						},
						{
							key:      "isSelf",
							dataType: Boolean,
							value:    false,
						},
						argumentValue("workflowName", args, 0),
					},
				},
			}
		},
	}
	actions["run"] = actionDefinition{
		identifier: "runworkflow",
		parameters: []parameterDefinition{
			{
				name:      "shortcutName",
				validType: String,
			},
			{
				name:      "output",
				validType: Variable,
			},
			{
				name:      "isSelf",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFWorkflow",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "workflowIdentifier",
							dataType: Text,
							value:    shortcutsUUID(),
						},
						{
							key:      "isSelf",
							dataType: Boolean,
							value:    false,
						},
						argumentValue("workflowName", args, 0),
					},
				},
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["list"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "listItem",
				validType: String,
				infinite:  true,
			},
		},
		make: func(args []actionArgument) []plistData {
			var listItems []string
			for _, item := range args {
				listItems = append(listItems, plistValue(Text, item.value))
			}
			return []plistData{
				{
					key:      "WFItems",
					dataType: Array,
					value:    listItems,
				},
			}
		},
	}
	actions["calcAverage"] = actionDefinition{
		identifier: "statistics",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return calculateStatistics("Average", args)
		},
	}
	actions["calcMin"] = actionDefinition{
		identifier: "statistics",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return calculateStatistics("Minimum", args)
		},
	}
	actions["calcMax"] = actionDefinition{
		identifier: "statistics",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return calculateStatistics("Maximum", args)
		},
	}
	actions["calcSum"] = actionDefinition{
		identifier: "statistics",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return calculateStatistics("Sum", args)
		},
	}
	actions["calcMedian"] = actionDefinition{
		identifier: "statistics",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return calculateStatistics("Median", args)
		},
	}
	actions["calcMode"] = actionDefinition{
		identifier: "statistics",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return calculateStatistics("Mode", args)
		},
	}
	actions["calcRange"] = actionDefinition{
		identifier: "statistics",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return calculateStatistics("Range", args)
		},
	}
	actions["calcStdDevi"] = actionDefinition{
		identifier: "statistics",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return calculateStatistics("Standard Deviation", args)
		},
	}
	actions["dismissSiri"] = actionDefinition{}
	actions["isOnline"] = actionDefinition{
		identifier: "getipaddress",
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFIPAddressSourceOption",
					dataType: Text,
					value:    "External",
				},
				{
					key:      "WFIPAddressTypeOption",
					dataType: Text,
					value:    "IPv4",
				},
			}
		},
	}
	var ipTypes = []string{"ipv4", "ipv6"}
	actions["getLocalIP"] = actionDefinition{
		identifier: "getipaddress",
		parameters: []parameterDefinition{
			{
				name:      "type",
				validType: String,
				defaultValue: actionArgument{
					valueType: String,
					value:     "IPv4",
				},
				optional: true,
			},
		},
		check: func(args []actionArgument) {
			checkEnum("IP address type", ipTypes, args, 0)
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFIPAddressTypeOption", args, 0),
				{
					key:      "WFIPAddressSourceOption",
					dataType: Text,
					value:    "Local",
				},
			}
		},
	}
	actions["getExternalIP"] = actionDefinition{
		identifier: "getipaddress",
		parameters: []parameterDefinition{
			{
				name:      "type",
				validType: String,
				defaultValue: actionArgument{
					valueType: String,
					value:     "IPv4",
				},
				optional: true,
			},
		},
		check: func(args []actionArgument) {
			checkEnum("IP address type", ipTypes, args, 0)
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFIPAddressTypeOption", args, 0),
				{
					key:      "WFIPAddressSourceOption",
					dataType: Text,
					value:    "External",
				},
			}
		},
	}
	actions["firstListItem"] = actionDefinition{
		identifier: "getitemfromlist",
		parameters: []parameterDefinition{
			{
				name:      "list",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "First Item",
				},
			}
		},
	}
	actions["lastListItem"] = actionDefinition{
		identifier: "getitemfromlist",
		parameters: []parameterDefinition{
			{
				name:      "list",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "Last Item",
				},
			}
		},
	}
	actions["randomListItem"] = actionDefinition{
		identifier: "getitemfromlist",
		parameters: []parameterDefinition{
			{
				name:      "list",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "Random Item",
				},
			}
		},
	}
	actions["getListItem"] = actionDefinition{
		identifier: "getitemfromlist",
		parameters: []parameterDefinition{
			{
				name:      "list",
				validType: Variable,
			},
			{
				name:      "index",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
				argumentValue("WFItemIndex", args, 1),
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "Item at Index",
				},
			}
		},
	}
	actions["getListItems"] = actionDefinition{
		identifier: "getitemfromlist",
		parameters: []parameterDefinition{
			{
				name:      "list",
				validType: Variable,
			},
			{
				name:      "start",
				validType: Integer,
			},
			{
				name:      "end",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
				argumentValue("WFItemRangeStart", args, 1),
				argumentValue("WFItemRangeEnd", args, 2),
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "Items in Range",
				},
			}
		},
	}
	actions["getNumbers"] = actionDefinition{
		identifier: "detect.number",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["getDictionary"] = actionDefinition{
		identifier: "detect.dictionary",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["getText"] = actionDefinition{
		identifier: "detect.text",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["getContacts"] = actionDefinition{
		identifier: "detect.contacts",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["getDates"] = actionDefinition{
		identifier: "detect.date",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["getEmails"] = actionDefinition{
		identifier: "detect.emailaddress",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: String,
			},
		},
	}
	actions["getImages"] = actionDefinition{
		identifier: "detect.images",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["getPhoneNumbers"] = actionDefinition{
		identifier: "detect.phonenumber",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["getURLs"] = actionDefinition{
		identifier: "detect.link",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["getAllWallpapers"] = actionDefinition{
		identifier: "posters.get",
		minVersion: 16.2,
	}
	actions["getWallpaper"] = actionDefinition{
		identifier: "posters.get",
		minVersion: 16.2,
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFPosterType",
					dataType: Text,
					value:    "Current",
				},
			}
		},
	}
	actions["setWallpaper"] = actionDefinition{
		identifier: "wallpaper.set",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["startScreensaver"] = actionDefinition{}
	actions["contentGraph"] = actionDefinition{
		identifier: "viewresult",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["openXCallbackURL"] = actionDefinition{
		identifier: "openxcallbackurl",
		parameters: []parameterDefinition{
			{
				name:      "url",
				key:       "WFXCallbackURL",
				validType: String,
				infinite:  true,
			},
		},
	}
	actions["openCustomXCallbackURL"] = actionDefinition{
		identifier: "openxcallbackurl",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
			},
			{
				name:      "successKey",
				validType: String,
			},
			{
				name:      "cancelKey",
				validType: String,
			},
			{
				name:      "errorKey",
				validType: String,
			},
			{
				name:      "successURL",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			var xCallbackParams = []plistData{
				argumentValue("WFXCallbackURL", args, 0),
				argumentValue("WFXCallbackCustomSuccessKey", args, 1),
				argumentValue("WFXCallbackCustomCancelKey", args, 2),
				argumentValue("WFXCallbackCustomErrorKey", args, 3),
				argumentValue("WFXCallbackCustomSuccessURL", args, 4),
			}
			if args[1].value.(string) != "" || args[2].value.(string) != "" || args[3].value.(string) != "" {
				xCallbackParams = append(xCallbackParams, plistData{
					key:      "WFXCallbackCustomCallbackEnabled",
					dataType: Boolean,
					value:    true,
				})
			}
			if args[4].value.(string) != "" {
				xCallbackParams = append(xCallbackParams, plistData{
					key:      "WFXCallbackCustomSuccessURLEnabled",
					dataType: Boolean,
					value:    true,
				})
			}
			return xCallbackParams
		},
	}
	actions["output"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "output",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFOutput", args, 0),
				{
					key:      "WFNoOutputSurfaceBehavior",
					dataType: Text,
					value:    "Do Nothing",
				},
			}
		},
	}
	actions["mustOutput"] = actionDefinition{
		identifier: "output",
		parameters: []parameterDefinition{
			{
				name:      "output",
				validType: String,
			},
			{
				name:      "response",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFOutput", args, 0),
				argumentValue("WFResponse", args, 1),
				{
					key:      "WFNoOutputSurfaceBehavior",
					dataType: Text,
					value:    "Respond",
				},
			}
		},
	}
	actions["outputOrClipboard"] = actionDefinition{
		identifier: "output",
		parameters: []parameterDefinition{
			{
				name:      "output",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFOutput", args, 0),
				{
					key:      "WFNoOutputSurfaceBehavior",
					dataType: Text,
					value:    "Copy to Clipboard",
				},
			}
		},
	}
	actions["setWifi"] = actionDefinition{
		identifier: "wifi.set",
		parameters: []parameterDefinition{
			{
				name:      "status",
				key:       "OnValue",
				validType: Bool,
			},
		},
	}
	actions["setCellularData"] = actionDefinition{
		identifier: "cellulardata.set",
		parameters: []parameterDefinition{
			{
				name:      "status",
				key:       "OnValue",
				validType: Bool,
			},
		},
	}
	actions["setCellularVoice"] = actionDefinition{
		identifier: "cellular.rat.set",
		parameters: []parameterDefinition{
			{
				name:      "status",
				key:       "OnValue",
				validType: Bool,
			},
		},
	}
	actions["toggleBluetooth"] = actionDefinition{
		identifier: "bluetooth.set",
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "OnValue",
					dataType: Boolean,
					value:    false,
				},
				{
					key:      "operation",
					dataType: Text,
					value:    "toggle",
				},
			}
		},
	}
	actions["setBluetooth"] = actionDefinition{
		identifier: "bluetooth.set",
		parameters: []parameterDefinition{
			{
				name:      "status",
				validType: Bool,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("OnValue", args, 0),
				{
					key:      "operation",
					dataType: Text,
					value:    "set",
				},
			}
		},
	}
	actions["playSound"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["round"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: Integer,
			},
			{
				name:      "roundTo",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return roundingValue("Normal", args)
		},
	}
	actions["roundUp"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: Integer,
			},
			{
				name:      "roundTo",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return roundingValue("Always Round Up", args)
		},
	}
	actions["roundDown"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: Integer,
			},
			{
				name:      "roundTo",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return roundingValue("Always Round Down", args)
		},
	}
	actions["runShellScript"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "script",
				key:       "Script",
				validType: String,
			},
			{
				name:      "input",
				key:       "Input",
				validType: Variable,
			},
			{
				name:      "shell",
				key:       "Shell",
				validType: String,
				defaultValue: actionArgument{
					valueType: String,
					value:     "/bin/zsh",
				},
			},
			{
				name:      "inputMode",
				key:       "InputMode",
				validType: String,
				defaultValue: actionArgument{
					valueType: String,
					value:     "to stdin",
				},
			},
		},
	}
}

func sharingActions() {
	actions["airdrop"] = actionDefinition{
		identifier: "airdropdocument",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["share"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: String,
			},
		},
	}
	actions["copyToClipboard"] = actionDefinition{
		identifier: "setclipboard",
		parameters: []parameterDefinition{
			{
				name:      "value",
				key:       "WFInput",
				validType: Variable,
			},
			{
				name:      "local",
				key:       "WFLocalOnly",
				validType: Bool,
			},
			{
				name:      "expire",
				key:       "WFExpirationDate",
				validType: String,
			},
		},
	}
	actions["getClipboard"] = actionDefinition{}
}

func webActions() {
	actions["getURLHeaders"] = actionDefinition{
		identifier: "url.getheaders",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["openURL"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Show-WFInput",
					dataType: Boolean,
					value:    true,
				},
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["runJavaScriptOnWebpage"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "javascript",
				validType: String,
				key:       "WFJavaScript",
			},
		},
	}
	var engines = []string{"amazon", "bing", "duckduckgo", "ebay", "google", "reddit", "twitter", "yahoo!", "youTube"}
	actions["searchWeb"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "engine",
				validType: String,
				key:       "WFSearchWebDestination",
			},
			{
				name:      "query",
				validType: String,
				key:       "WFInputText",
			},
		},
		check: func(args []actionArgument) {
			checkEnum("search engine", engines, args, 0)
		},
	}
	actions["showWebpage"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFURL",
			},
			{
				name:      "useReader",
				validType: Bool,
				key:       "WFEnterSafariReader",
				optional:  true,
			},
		},
	}
	actions["getRSSFeeds"] = actionDefinition{
		identifier: "rss.extract",
		parameters: []parameterDefinition{
			{
				name:      "urls",
				validType: String,
				key:       "WFURLs",
			},
		},
	}
	actions["getRSS"] = actionDefinition{
		identifier: "rss",
		parameters: []parameterDefinition{
			{
				name:      "items",
				validType: Integer,
				key:       "WFRSSItemQuantity",
			},
			{
				name:      "url",
				validType: String,
				key:       "WFRSSFeedURL",
			},
		},
	}
	var webpageDetails = []string{"page contents", "page selection", "page url", "name"}
	actions["getWebPageDetail"] = actionDefinition{
		identifier: "properties.safariwebpage",
		parameters: []parameterDefinition{
			{
				name:      "webpage",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "detail",
				validType: String,
				key:       "WFContentItemPropertyName",
			},
		},
		check: func(args []actionArgument) {
			checkEnum("webpage detail", webpageDetails, args, 1)
		},
	}
	actions["getArticleDetail"] = actionDefinition{
		identifier: "properties.articles",
		parameters: []parameterDefinition{
			{
				name:      "article",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "detail",
				validType: String,
				key:       "WFContentItemPropertyName",
			},
		},
	}
	actions["getCurrentURL"] = actionDefinition{
		identifier: "safari.geturl",
	}
	actions["getWebpageContents"] = actionDefinition{
		identifier: "getwebpagecontents",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["searchGiphy"] = actionDefinition{
		identifier: "giphy",
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "WFGiphyQuery",
			},
		},
	}
	actions["getGifs"] = actionDefinition{
		identifier: "giphy",
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
			},
			{
				name:      "gifs",
				validType: Integer,
				defaultValue: actionArgument{
					value:     1,
					valueType: Integer,
				},
				optional: true,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGiphyShowPicker",
					dataType: Boolean,
					value:    false,
				},
				argumentValue("WFGiphyQuery", args, 0),
				argumentValue("WFGiphyLimit", args, 1),
			}
		},
	}
	actions["getArticle"] = actionDefinition{
		identifier: "getarticle",
		parameters: []parameterDefinition{
			{
				name:      "webpage",
				validType: String,
				key:       "WFWebPage",
			},
		},
	}
	actions["expandURL"] = actionDefinition{
		identifier: "url.expand",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "URL",
			},
		},
	}
	var urlComponents = []string{"scheme", "user", "password", "host", "port", "path", "query", "fragment"}
	actions["getURLDetail"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFURL",
			},
			{
				name:      "detail",
				validType: String,
				key:       "WFURLComponent",
			},
		},
		check: func(args []actionArgument) {
			checkEnum("URL component", urlComponents, args, 0)
		},
	}
	actions["downloadURL"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFURL",
			},
			{
				name:      "headers",
				validType: Dict,
				key:       "WFHTTPHeaders",
				optional:  true,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFHTTPMethod",
					dataType: Text,
					value:    "GET",
				},
				argumentValue("WFURL", args, 0),
				argumentValue("WFHTTPHeaders", args, 1),
			}
		},
	}
	var httpMethods = []string{"get", "post", "put", "patch", "delete"}
	var httpParams = []parameterDefinition{
		{
			name:      "url",
			validType: String,
		},
		{
			name:      "method",
			validType: String,
			optional:  true,
			defaultValue: actionArgument{
				valueType: String,
				value:     "GET",
			},
		},
		{
			name:      "body",
			validType: Variable,
			optional:  true,
		},
		{
			name:      "headers",
			validType: Dict,
			optional:  true,
		},
	}
	actions["formRequest"] = actionDefinition{
		identifier: "downloadurl",
		parameters: httpParams,
		check: func(args []actionArgument) {
			checkEnum("HTTP method", httpMethods, args, 1)
		},
		make: func(args []actionArgument) []plistData {
			return httpRequest("Form", "WFFormValues", args)
		},
	}
	actions["jsonRequest"] = actionDefinition{
		identifier: "downloadurl",
		parameters: httpParams,
		check: func(args []actionArgument) {
			checkEnum("HTTP method", httpMethods, args, 1)
		},
		make: func(args []actionArgument) []plistData {
			return httpRequest("JSON", "WFJSONValues", args)
		},
	}
	actions["fileRequest"] = actionDefinition{
		identifier: "downloadurl",
		parameters: httpParams,
		check: func(args []actionArgument) {
			checkEnum("HTTP method", httpMethods, args, 1)
		},
		make: func(args []actionArgument) []plistData {
			return httpRequest("File", "WFRequestVariable", args)
		},
	}
}

func customActions() {
	actions["makeVCard"] = actionDefinition{
		identifier: "gettext",
		parameters: []parameterDefinition{
			{
				name:      "title",
				validType: String,
			},
			{
				name:      "subtitle",
				validType: String,
			},
			{
				name:      "url",
				validType: String,
			},
		},
		check: func(args []actionArgument) {
			var image = getArgValue(args[2])
			if reflect.TypeOf(image).String() != stringType {
				parserError("Image path for VCard must be a string literal")
			}
		},
		make: func(args []actionArgument) []plistData {
			var title = args[0].value.(string)
			var subtitle = args[1].value.(string)
			if _, found := variables[title]; found {
				title = "{" + title + "}"
			}
			if _, found := variables[subtitle]; found {
				subtitle = "{" + subtitle + "}"
			}
			var vcard = "BEGIN:VCARD\nVERSION:3.0\n"
			vcard += "N;CHARSET=utf-8:" + title + "\n"
			vcard += "ORG:" + subtitle + "\nPHOTO;ENCODING=b:"
			bytes, readErr := os.ReadFile(getArgValue(args[2]).(string))
			handle(readErr)
			vcard += base64.StdEncoding.EncodeToString(bytes) + "\nEND:VCARD"
			args[0] = actionArgument{
				valueType: String,
				value:     vcard,
			}
			return []plistData{
				argumentValue("WFTextActionText", args, 0),
			}
		},
	}
}

var contactValues []string

func contactValue(key string, contentKit string, args []actionArgument) []plistData {
	contactValues = []string{}
	var entryType int
	switch contentKit {
	case "emailaddress":
		entryType = 2
	case "phonenumber":
		entryType = 1
	}
	for _, item := range args {
		contactValues = append(contactValues, plistDict("", []plistData{
			{
				key:      "EntryType",
				dataType: Number,
				value:    entryType,
			},
			{
				key:      "SerializedEntry",
				dataType: Dictionary,
				value: []plistData{
					{
						key:      "link.contentkit." + contentKit,
						dataType: Text,
						value:    item.value,
					},
				},
			},
		}))
	}
	return []plistData{
		{
			key:      key,
			dataType: Dictionary,
			value: []plistData{
				{
					key:      "Value",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "WFContactFieldValues",
							dataType: Array,
							value:    contactValues,
						},
					},
				},
				{
					key:      "WFSerializationType",
					dataType: Text,
					value:    "WFContactFieldValue",
				},
			},
		},
	}
}

func roundingValue(mode string, args []actionArgument) []plistData {
	switch args[1].value {
	case "1":
		args[1].value = "Ones Place"
	case "10":
		args[1].value = "Tens Place"
	case "100":
		args[1].value = "Hundreds Place"
	case "1000":
		args[1].value = "Thousands"
	case "10000":
		args[1].value = "Ten Thousands"
	case "100000":
		args[1].value = "Hundred Thousands"
	case "1000000":
		args[1].value = "Millions"
	}
	return []plistData{
		{
			key:      "WFRoundMode",
			dataType: Text,
			value:    mode,
		},
		argumentValue("WFInput", args, 0),
		{
			key:      "WFRoundTo",
			dataType: Text,
			value:    args[1].value,
		},
	}
}

func calculateStatistics(operation string, args []actionArgument) []plistData {
	return []plistData{
		{
			key:      "WFStatisticsOperation",
			dataType: Text,
			value:    operation,
		},
		variableInput("WFInput", args[1].value.(string)),
	}
}

func adjustDate(operation string, unit string, args []actionArgument) []plistData {
	var adjustDateParams = []plistData{
		{
			key:      "WFAdjustOperation",
			dataType: Text,
			value:    operation,
		},
		argumentValue("WFDate", args, 0),
	}
	if unit != "" {
		adjustDateParams = append(adjustDateParams, plistData{
			key:      "WFDuration",
			dataType: Dictionary,
			value: []plistData{
				{
					key:      "Value",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "Unit",
							dataType: Text,
							value:    unit,
						},
						argumentValue("Magnitude", args, 1),
					},
				},
				{
					key:      "WFSerializationType",
					dataType: Text,
					value:    "WFQuantityFieldValue",
				},
			},
		})
	}
	return adjustDateParams
}

func changeCase(textCase string, args []actionArgument) []plistData {
	return []plistData{
		{
			key:      "Show-text",
			dataType: Boolean,
			value:    true,
		},
		{
			key:      "WFCaseType",
			dataType: Text,
			value:    textCase,
		},
		argumentValue("text", args, 0),
	}
}

func textParts(args []actionArgument) []plistData {
	return []plistData{
		{
			key:      "Show-text",
			dataType: Boolean,
			value:    true,
		},
		{
			key:      "WFTextSeparator",
			dataType: Text,
			value:    "Custom",
		},
		argumentValue("text", args, 0),
		argumentValue("WFTextCustomSeparator", args, 1),
	}
}

func replaceText(caseSensitive bool, regExp bool, args []actionArgument) []plistData {
	return []plistData{
		{
			key:      "WFReplaceTextCaseSensitive",
			dataType: Boolean,
			value:    caseSensitive,
		},
		{
			key:      "WFReplaceTextRegularExpression",
			dataType: Boolean,
			value:    regExp,
		},
		argumentValue("WFReplaceTextFind", args, 0),
		argumentValue("WFReplaceTextReplace", args, 1),
		argumentValue("WFInput", args, 2),
	}
}

func languageCode(language string) string {
	makeLanguages()
	if _, found := languages[language]; found {
		return languages[language]
	}
	return language
}

func count(countType string, args []actionArgument) []plistData {
	return []plistData{
		{
			key:      "WFCountType",
			dataType: Text,
			value:    countType,
		},
		variableInput("Input", args[0].value.(string)),
	}
}

var appIds map[string]string

func makeAppIds() {
	if len(appIds) != 0 {
		return
	}
	appIds = make(map[string]string)
	appIds["appstore"] = "com.apple.AppStore"
	appIds["files"] = "com.apple.DocumentsApp"
	appIds["shortcuts"] = "is.workflow.my.app"
	appIds["safari"] = "com.apple.mobilesafari"
	appIds["facetime"] = "com.apple.facetime"
	appIds["notes"] = "com.apple.mobilenotes"
	appIds["phone"] = "com.apple.mobilephone"
	appIds["reminders"] = "com.apple.reminders"
	appIds["mail"] = "com.apple.mobilemail"
	appIds["music"] = "com.apple.Music"
	appIds["calendar"] = "com.apple.mobilecal"
	appIds["maps"] = "com.apple.Maps"
	appIds["contacts"] = "com.apple.MobileAddressBook"
	appIds["health"] = "com.apple.Health"
	appIds["photos"] = "com.apple.mobileslideshow"
}

func replaceAppId(args []actionArgument) {
	makeAppIds()
	var id = getArgValue(args[0]).(string)
	if _, found := appIds[id]; found {
		args[0].value = appIds[id]
	}
}

func httpRequest(bodyType string, valuesKey string, args []actionArgument) []plistData {
	return []plistData{
		{
			key:      "WFHTTPBodyType",
			dataType: Text,
			value:    bodyType,
		},
		argumentValue("WFURL", args, 0),
		argumentValue("WFHTTPMethod", args, 1),
		argumentValue(valuesKey, args, 2),
		argumentValue("WFHTTPHeaders", args, 3),
	}
}
