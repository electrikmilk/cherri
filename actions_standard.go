/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

// FIXME: Some of these actions have enumerable arguments (a set list values),
//  but the argument value is not being checked against it's possible values.
//  Use the "hash" action as an example.

func standardActions() {
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
		stdIdentifier: "addnewcalendar",
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("CalendarName", args, 0),
			}
		},
	}
	actions["addSeconds"] = actionDefinition{
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
		stdIdentifier: "adjustdate",
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
	actions["emailAddress"] = actionDefinition{
		stdIdentifier: "email",
		parameters: []parameterDefinition{
			{
				name:      "email",
				validType: String,
				noMax:     true,
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
				noMax:     true,
			},
		},
		make: func(args []actionArgument) []plistData {
			return contactValue("WFPhoneNumber", "phonenumber", args)
		},
	}
	actions["selectContact"] = actionDefinition{
		stdIdentifier: "selectcontacts",
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
		stdIdentifier: "selectemail",
	}
	actions["selectPhoneNumber"] = actionDefinition{
		stdIdentifier: "selectphone",
	}
	actions["getFromContact"] = actionDefinition{
		stdIdentifier: "properties.contacts",
		parameters: []parameterDefinition{
			{
				name:      "contact",
				validType: String,
			},
			{
				name:      "property",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
				argumentValue("WFContentItemPropertyName", args, 1),
			}
		},
	}
}

func documentActions() {
	// FIXME: Writing to locations other than the Shortcuts folder
	actions["createShortcutsFolder"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "path",
				validType: String,
				key:       "WFFilePath",
			},
		},
	}
	actions["getFolderContents"] = actionDefinition{
		stdIdentifier: "file.getfoldercontents",
		parameters: []parameterDefinition{
			{
				name:      "folder",
				validType: Variable,
			},
			{
				name:      "recursive",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     true,
				},
				optional: true,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFFolder", args[0].value.(string)),
				argumentValue("Recursive", args, 1),
			}
		},
	}
	actions["matchedTextGroupIndex"] = actionDefinition{
		stdIdentifier: "text.match.getgroup",
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
		stdIdentifier: "documentpicker.open",
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
	actions["markup"] = actionDefinition{
		stdIdentifier: "avairyeditphoto",
		parameters: []parameterDefinition{
			{
				name:      "document",
				validType: Variable,
				key:       "WFDocument",
			},
		},
	}
	actions["rename"] = actionDefinition{
		stdIdentifier: "file.rename",
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
		stdIdentifier: "file.reveal",
		parameters: []parameterDefinition{
			{
				name:      "files",
				validType: Variable,
				key:       "WFFile",
			},
		},
	}
	actions["define"] = actionDefinition{
		stdIdentifier: "showdefinition",
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
		stdIdentifier: "generatebarcode",
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
			if args[0].value != nil {
				checkEnum("error correction level", args[0], errorCorrectionLevels)
			}
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
		stdIdentifier: "gethtmlfromrichtext",
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
		stdIdentifier: "getmarkdownfromrichtext",
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
		stdIdentifier: "file.select",
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
		stdIdentifier: "file.getlink",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFFile", args[0].value.(string)),
			}
		},
	}
	actions["getParentDirectory"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["getEmojiName"] = actionDefinition{
		stdIdentifier: "getnameofemoji",
		parameters: []parameterDefinition{
			{
				name:      "emoji",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["getFileDetail"] = actionDefinition{
		stdIdentifier: "properties.files",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
			},
			{
				name:      "detail",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFFolder", args[0].value.(string)),
				argumentValue("WFContentItemPropertyName", args, 1),
			}
		},
	}
	actions["deleteFiles"] = actionDefinition{
		stdIdentifier: "file.delete",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
			{
				name:      "immediately",
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
				variableInput("WFInput", args[0].value.(string)),
				argumentValue("WFDeleteImmediatelyDelete", args, 1),
			}
		},
	}
	actions["getTextFromImage"] = actionDefinition{
		stdIdentifier: "extracttextfromimage",
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
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("BooksInput", args[0].value.(string)),
			}
		},
	}
	actions["saveFile"] = actionDefinition{
		stdIdentifier: "documentpicker.save",
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
	actions["getSelectedFiles"] = actionDefinition{stdIdentifier: "finder.getselectedfiles"}
	actions["extractArchive"] = actionDefinition{
		stdIdentifier: "unzip",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFArchive", args[0].value.(string)),
			}
		},
	}
	actions["makeArchive"] = actionDefinition{
		stdIdentifier: "makezip",
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
			},
			{
				name:      "format",
				validType: String,
			},
			{
				name:      "files",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFZIPName", args, 0),
				argumentValue("WFArchiveFormat", args, 1),
				variableInput("WFInput", args[2].value.(string)),
			}
		},
	}
	actions["quicklook"] = actionDefinition{
		stdIdentifier: "previewdocument",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["translateFrom"] = actionDefinition{
		stdIdentifier: "translate",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
			},
			{
				name:      "from",
				validType: String,
			},
			{
				name:      "to",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInputText", args[0].value.(string)),
				argumentValue("WFSelectedFromLanguage", args, 0),
				argumentValue("WFSelectedLanguage", args, 0),
			}
		},
	}
	actions["translate"] = actionDefinition{
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
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInputText", args[0].value.(string)),
				{
					key:      "WFSelectedFromLanguage",
					dataType: Text,
					value:    "Detect Language",
				},
				argumentValue("WFSelectedLanguage", args, 0),
			}
		},
	}
	actions["detectLanguage"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["replaceText"] = actionDefinition{
		stdIdentifier: "text.replace",
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
		stdIdentifier: "text.replace",
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
		stdIdentifier: "text.replace",
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
		stdIdentifier: "text.replace",
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
		stdIdentifier: "text.changecase",
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
		stdIdentifier: "text.changecase",
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
		stdIdentifier: "text.changecase",
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
		stdIdentifier: "text.changecase",
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
		stdIdentifier: "text.changecase",
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
		stdIdentifier: "text.changecase",
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
		stdIdentifier: "text.split",
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
		stdIdentifier: "text.combine",
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
}

func locationActions() {
	actions["getCurrentLocation"] = actionDefinition{
		stdIdentifier: "location",
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
		stdIdentifier: "detect.address",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["getCurrentWeather"] = actionDefinition{
		stdIdentifier: "currentconditions",
	}
	actions["getCurrentWeatherAt"] = actionDefinition{
		stdIdentifier: "currentconditions",
		parameters: []parameterDefinition{
			{
				name:      "location",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFWeatherCustomLocation", args, 0),
			}
		},
	}
	actions["openInMaps"] = actionDefinition{
		stdIdentifier: "searchmaps",
		parameters: []parameterDefinition{
			{
				name:      "location",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
			}
		},
	}
	actions["streetAddress"] = actionDefinition{
		stdIdentifier: "address",
		parameters: []parameterDefinition{
			{
				name:      "addressLine2",
				validType: String,
			},
			{
				name:      "addressLine2",
				validType: String,
			},
			{
				name:      "city",
				validType: String,
			},
			{
				name:      "state",
				validType: String,
			},
			{
				name:      "country",
				validType: String,
			},
			{
				name:      "zipCode",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFAddressLine1", args, 0),
				argumentValue("WFAddressLine2", args, 1),
				argumentValue("WFCity", args, 2),
				argumentValue("WFState", args, 3),
				argumentValue("WFCountry", args, 4),
				argumentValue("WFPostalCode", args, 5),
			}
		},
	}
	actions["getWeatherDetail"] = actionDefinition{
		stdIdentifier: "properties.weather.conditions",
		parameters: []parameterDefinition{
			{
				name:      "weather",
				validType: Variable,
			},
			{
				name:      "property",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
				argumentValue("WFContentItemPropertyName", args, 1),
			}
		},
	}
	actions["getWeatherForcast"] = actionDefinition{
		stdIdentifier: "weather.forecast",
		parameters: []parameterDefinition{
			{
				name:      "type",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFWeatherForecastType", args, 0),
			}
		},
	}
	actions["getWeatherForcastAt"] = actionDefinition{
		stdIdentifier: "weather.forecast",
		parameters: []parameterDefinition{
			{
				name:      "type",
				validType: String,
			},
			{
				name:      "location",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFWeatherForecastType", args, 0),
				argumentValue("WFInput", args, 1),
			}
		},
	}
	actions["getLocationDetail"] = actionDefinition{
		stdIdentifier: "properties.locations",
		parameters: []parameterDefinition{
			{
				name:      "location",
				validType: Variable,
			},
			{
				name:      "property",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
				argumentValue("WFContentItemPropertyName", args, 1),
			}
		},
	}
	actions["getMapsLink"] = actionDefinition{
		stdIdentifier: "",
		parameters: []parameterDefinition{
			{
				name:      "location",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
			}
		},
	}
	actions["getHalfwayPoint"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "firstLocation",
				validType: Variable,
			},
			{
				name:      "secondLocation",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFGetHalfwayPointFirstLocation", args, 0),
				argumentValue("WFGetHalfwayPointSecondLocation", args, 1),
			}
		},
	}
}

func mediaActions() {
	actions["clearUpNext"] = actionDefinition{}
	actions["getCurrentSong"] = actionDefinition{}
	actions["latestPhotoImport"] = actionDefinition{stdIdentifier: "getlatestphotoimport"}
	actions["takePhoto"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "showPreview",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     true,
				},
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFCameraCaptureShowPreview", args, 0),
			}
		},
	}
	actions["takePhotos"] = actionDefinition{
		stdIdentifier: "takephoto",
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
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInputMedia", args, 0),
			}
		},
	}
	actions["takeVideo"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "camera",
				validType: String,
				defaultValue: actionArgument{
					valueType: String,
					value:     "Front",
				},
			},
			{
				name:      "quality",
				validType: String,
				defaultValue: actionArgument{
					valueType: String,
					value:     "Medium",
				},
			},
			{
				name:      "startImmediately",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
			},
		},
		make: func(args []actionArgument) []plistData {
			return argumentValues(&args, "WFCameraCaptureDevice", "WFCameraCaptureQuality", "WFRecordingStart")
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
		stdIdentifier: "getclassaction",
		parameters: []parameterDefinition{
			{
				name:      "class",
				validType: String,
			},
			{
				name:      "from",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("Class", args, 0),
				variableInput("Input", args[1].value.(string)),
			}
		},
	}
	actions["getOnScreenContent"] = actionDefinition{}
	actions["fileSize"] = actionDefinition{
		stdIdentifier: "format.filesize",
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
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFDeviceDetail", args, 0),
			}
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
		stdIdentifier: "getitemname",
		parameters: []parameterDefinition{
			{
				name:      "item",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["setName"] = actionDefinition{
		stdIdentifier: "setitemname",
		parameters: []parameterDefinition{
			{
				name:      "item",
				validType: Variable,
			},
			{
				name:      "name",
				validType: String,
			},
			{
				name:      "includeFileExtension",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
				argumentValue("WFName", args, 1),
				argumentValue("WFDontIncludeFileExtension", args, 2),
			}
		},
	}
	actions["countItems"] = actionDefinition{
		stdIdentifier: "count",
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
		stdIdentifier: "count",
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
		stdIdentifier: "count",
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
		stdIdentifier: "count",
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
		stdIdentifier: "count",
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
		stdIdentifier: "appearance",
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
		stdIdentifier: "appearance",
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
		stdIdentifier: "appearance",
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
	actions["getShortcuts"] = actionDefinition{
		stdIdentifier: "getmyworkflows",
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
		stdIdentifier: "readinglist",
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
			},
			{
				name:      "type",
				validType: String,
				defaultValue: actionArgument{
					valueType: "MD5",
					value:     String,
				},
				optional: true,
			},
		},
		check: func(args []actionArgument) {
			if args[1].value != nil {
				checkEnum("hash type", args[1], hashTypes)
				args[1].value = strings.ToUpper(args[1].value.(string))
			}
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
				argumentValue("WFHashType", args, 1),
			}
		},
	}
	actions["formatNumber"] = actionDefinition{
		stdIdentifier: "format.number",
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: Integer,
			},
			{
				name:      "decimalPlaces",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFNumber", args, 0),
				argumentValue("WFNumberFormatDecimalPlaces", args, 1),
			}
		},
	}
	actions["randomNumber"] = actionDefinition{
		stdIdentifier: "number.random",
		parameters: []parameterDefinition{
			{
				name:      "min",
				validType: Integer,
			},
			{
				name:      "max",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFRandomNumberMinimum", args, 0),
				argumentValue("WFRandomNumberMaximum", args, 1),
			}
		},
	}
	actions["base64Encode"] = actionDefinition{
		stdIdentifier: "base64encode",
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
		stdIdentifier: "base64encode",
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
		stdIdentifier: "urlencode",
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
		stdIdentifier: "urlencode",
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
		stdIdentifier: "showresult",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("Text", args, 0),
			}
		},
	}
	actions["waitToReturn"] = actionDefinition{}
	actions["notification"] = actionDefinition{
		stdIdentifier: "notification",
		parameters: []parameterDefinition{
			{
				name:      "body",
				validType: String,
			},
			{
				name:      "title",
				validType: String,
				optional:  true,
			},
			{
				name:      "playSound",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     true,
				},
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFNotificationActionBody", args, 0),
				argumentValue("WFNotificationActionTitle", args, 1),
				argumentValue("WFNotificationActionSound", args, 2),
			}
		},
	}
	actions["stop"] = actionDefinition{
		stdIdentifier: "exit",
	}
	actions["nothing"] = actionDefinition{}
	actions["wait"] = actionDefinition{
		stdIdentifier: "delay",
		parameters: []parameterDefinition{
			{
				name:      "seconds",
				validType: Integer,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFDelayTime", args, 0),
			}
		},
	}
	actions["alert"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "alert",
				validType: String,
			},
			{
				name:      "title",
				validType: String,
				optional:  true,
			},
			{
				name:      "cancelButton",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     true,
				},
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFAlertActionMessage", args, 0),
				argumentValue("WFAlertActionTitle", args, 1),
				argumentValue("WFAlertActionCancelButtonShown", args, 2),
			}
		},
	}
	actions["askForInput"] = actionDefinition{
		stdIdentifier: "ask",
		parameters: []parameterDefinition{
			{
				name:      "inputType",
				validType: String,
			},
			{
				name:      "prompt",
				validType: String,
			},
			{
				name:      "defaultValue",
				validType: String,
				optional:  true,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInputType", args, 0),
				argumentValue("WFAskActionPrompt", args, 1),
				argumentValue("WFAskActionDefaultAnswer", args, 2),
			}
		},
	}
	actions["chooseFromList"] = actionDefinition{
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
				argumentValue("WFInput", args, 0),
				argumentValue("WFChooseFromListActionPrompt", args, 1),
			}
		},
	}
	actions["chooseMultipleFromList"] = actionDefinition{
		stdIdentifier: "choosefromlist",
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
		stdIdentifier: "getitemtype",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
			}
		},
	}
	actions["getKeys"] = actionDefinition{
		stdIdentifier: "getvalueforkey",
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
		stdIdentifier: "getvalueforkey",
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
		stdIdentifier: "getvalueforkey",
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
		stdIdentifier: "setvalueforkey",
		parameters: []parameterDefinition{
			{
				name:      "key",
				validType: String,
			},
			{
				name:      "value",
				validType: Variable,
			},
			{
				name:      "dictionary",
				validType: Dict,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFDictionary", args, 0),
				argumentValue("WFDictionaryKey", args, 1),
				argumentValue("WFDictionaryValue", args, 2),
			}
		},
	}
	actions["open"] = actionDefinition{
		stdIdentifier: "openworkflow",
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
		stdIdentifier: "runworkflow",
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
				noMax:     true,
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
		stdIdentifier: "statistics",
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
		stdIdentifier: "statistics",
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
		stdIdentifier: "statistics",
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
		stdIdentifier: "statistics",
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
		stdIdentifier: "statistics",
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
		stdIdentifier: "statistics",
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
		stdIdentifier: "statistics",
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
		stdIdentifier: "statistics",
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
		stdIdentifier: "getipaddress",
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
		stdIdentifier: "getipaddress",
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
			checkEnum("IP address type", args[0], ipTypes)
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
		stdIdentifier: "getipaddress",
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
			checkEnum("IP address type", args[0], ipTypes)
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
		stdIdentifier: "getitemfromlist",
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
		stdIdentifier: "getitemfromlist",
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
		stdIdentifier: "getitemfromlist",
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
		stdIdentifier: "getitemfromlist",
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
		stdIdentifier: "getitemfromlist",
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
		stdIdentifier: "detect.number",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["getDictionary"] = actionDefinition{
		stdIdentifier: "detect.dictionary",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["getText"] = actionDefinition{
		stdIdentifier: "detect.text",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["getContacts"] = actionDefinition{
		stdIdentifier: "detect.contacts",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["getDates"] = actionDefinition{
		stdIdentifier: "detect.date",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["getEmails"] = actionDefinition{
		stdIdentifier: "detect.emailaddress",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["getImages"] = actionDefinition{
		stdIdentifier: "detect.images",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["getPhoneNumbers"] = actionDefinition{
		stdIdentifier: "detect.phonenumber",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["getURLs"] = actionDefinition{
		stdIdentifier: "detect.link",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["setWallpaper"] = actionDefinition{
		stdIdentifier: "wallpaper.set",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["startScreensaver"] = actionDefinition{}
	actions["contentGraph"] = actionDefinition{
		stdIdentifier: "viewresult",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["openXCallbackURL"] = actionDefinition{
		stdIdentifier: "openxcallbackurl",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				noMax:     true,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFXCallbackURL", args, 0),
			}
		},
	}
	actions["openCustomXCallbackURL"] = actionDefinition{
		stdIdentifier: "openxcallbackurl",
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
		stdIdentifier: "output",
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
		stdIdentifier: "output",
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
		stdIdentifier: "wifi.set",
		parameters: []parameterDefinition{
			{
				name:      "status",
				validType: Bool,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("OnValue", args, 0),
			}
		},
	}
	actions["setCellularData"] = actionDefinition{
		stdIdentifier: "cellulardata.set",
		parameters: []parameterDefinition{
			{
				name:      "status",
				validType: Bool,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("OnValue", args, 0),
			}
		},
	}
	actions["setCellularVoice"] = actionDefinition{
		stdIdentifier: "cellular.rat.set",
		parameters: []parameterDefinition{
			{
				name:      "status",
				validType: Bool,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("OnValue", args, 0),
			}
		},
	}
	actions["toggleBluetooth"] = actionDefinition{
		stdIdentifier: "bluetooth.set",
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
		stdIdentifier: "bluetooth.set",
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
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
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
				validType: String,
			},
			{
				name:      "input",
				validType: Variable,
			},
			{
				name:      "shell",
				validType: String,
				defaultValue: actionArgument{
					valueType: String,
					value:     "/bin/zsh",
				},
			},
			{
				name:      "inputMode",
				validType: String,
				defaultValue: actionArgument{
					valueType: String,
					value:     "to stdin",
				},
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("Script", args, 0),
				variableInput("Input", args[1].value.(string)),
				argumentValue("Shell", args, 2),
				argumentValue("InputMode", args, 3),
			}
		},
	}
}

func sharingActions() {
	actions["airdrop"] = actionDefinition{
		stdIdentifier: "airdropdocument",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["share"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
			}
		},
	}
	actions["copyToClipboard"] = actionDefinition{
		stdIdentifier: "setclipboard",
		parameters: []parameterDefinition{
			{
				name:      "value",
				validType: Variable,
			},
			{
				name:      "local",
				validType: Bool,
			},
			{
				name:      "expire",
				validType: String,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
				argumentValue("WFLocalOnly", args, 1),
				argumentValue("WFExpirationDate", args, 2),
			}
		},
	}
	actions["getClipboard"] = actionDefinition{}
}

func webActions() {
	actions["getURLHeaders"] = actionDefinition{
		stdIdentifier: "url.getheaders",
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
			checkEnum("search engine", args[0], engines)
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
		stdIdentifier: "rss.extract",
		parameters: []parameterDefinition{
			{
				name:      "urls",
				validType: String,
				key:       "WFURLs",
			},
		},
	}
	actions["getRSS"] = actionDefinition{
		stdIdentifier: "rss",
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
		stdIdentifier: "properties.safariwebpage",
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
			checkEnum("webpage detail", args[1], webpageDetails)
		},
	}
	actions["getArticleDetail"] = actionDefinition{
		stdIdentifier: "properties.articles",
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
		stdIdentifier: "safari.geturl",
	}
	actions["getWebpageContents"] = actionDefinition{
		stdIdentifier: "getwebpagecontents",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["searchGiphy"] = actionDefinition{
		stdIdentifier: "giphy",
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "WFGiphyQuery",
			},
		},
	}
	actions["getGifs"] = actionDefinition{
		stdIdentifier: "giphy",
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
		stdIdentifier: "getarticle",
		parameters: []parameterDefinition{
			{
				name:      "webpage",
				validType: String,
				key:       "WFWebPage",
			},
		},
	}
	actions["expandURL"] = actionDefinition{
		stdIdentifier: "url.expand",
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
			checkEnum("URL component", args[0], urlComponents)
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
	var bodyTypes = []string{"json", "form", "file"}
	var httpMethods = []string{"get", "post", "put", "patch", "delete"}
	actions["httpRequest"] = actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
			},
			{
				name:      "method",
				validType: String,
			},
			{
				name:      "body",
				validType: Dict,
				optional:  true,
			},
			{
				name:      "bodyType",
				validType: String,
				optional:  true,
			},
			{
				name:      "headers",
				validType: Dict,
				optional:  true,
			},
		},
		check: func(args []actionArgument) {
			if args[1].value != nil {
				checkEnum("HTTP method", args[1], httpMethods)
			}
			if args[3].value != nil {
				checkEnum("HTTP body type", args[3], bodyTypes)
			}
		},
		make: func(args []actionArgument) []plistData {
			return argumentValues(&args, "WFURL", "WFHTTPMethod", "WFFormValues", "WFHTTPBodyType", "WFHTTPHeaders")
		},
	}
}

func customActions() {
	actions["makeVCard"] = actionDefinition{
		stdIdentifier: "gettext",
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
