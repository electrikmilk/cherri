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

// TODO: Finish standard actions, then combine them all into the makeActions() function.
//  Move the making of the actions map somewhere else.
//  Rename the function to standardActions() and rename the file to actions_standard.go.

// FIXME: Most of these actions that have enumerable values (a set list values),
//  do not check if the value matches and list out the valid values if it doesn't.
//  Use "hash" as an example.

func makeActions() {
	actions = make(map[string]actionDefinition)
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
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "addnewcalendar",
		args: []argumentDefinition{
			{
				field:     "name",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("CalendarName", args, 0),
			}
		},
	}
	actions["addSeconds"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
			{
				field:     "magnitude",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Add", "sec", args)
		},
	}
	actions["addMinutes"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
			{
				field:     "magnitude",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Add", "min", args)
		},
	}
	actions["addHours"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
			{
				field:     "magnitude",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Add", "hr", args)
		},
	}
	actions["addDays"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
			{
				field:     "magnitude",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Add", "days", args)
		},
	}
	actions["addWeeks"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
			{
				field:     "magnitude",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Add", "weeks", args)
		},
	}
	actions["addMonths"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
			{
				field:     "magnitude",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Add", "months", args)
		},
	}
	actions["addYears"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
			{
				field:     "magnitude",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Add", "yr", args)
		},
	}
	actions["subtractSeconds"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
			{
				field:     "magnitude",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "sec", args)
		},
	}
	actions["subtractMinutes"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
			{
				field:     "magnitude",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "min", args)
		},
	}
	actions["subtractHours"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
			{
				field:     "magnitude",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "hr", args)
		},
	}
	actions["subtractDays"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
			{
				field:     "magnitude",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "days", args)
		},
	}
	actions["subtractWeeks"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
			{
				field:     "magnitude",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "weeks", args)
		},
	}
	actions["subtractMonths"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
			{
				field:     "magnitude",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "months", args)
		},
	}
	actions["subtractYears"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "magnitude",
				validType: Integer,
			},
			{
				field:     "date",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Subtract", "yr", args)
		},
	}
	actions["getStartMinute"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Get Start of Minute", "", args)
		},
	}
	actions["getStartHour"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Get Start of Hour", "", args)
		},
	}
	actions["getStartWeek"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Get Start of Week", "", args)
		},
	}
	actions["getStartMonth"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Get Start of Month", "", args)
		},
	}
	actions["getStartYear"] = actionDefinition{
		ident: "adjustdate",
		args: []argumentDefinition{
			{
				field:     "date",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return adjustDate("Get Start of Year", "", args)
		},
	}
}

func contactActions() {
	actions["emailAddress"] = actionDefinition{
		ident: "email",
		args: []argumentDefinition{
			{
				field:     "email",
				validType: String,
				noMax:     true,
			},
		},
		call: func(args []actionArgument) []plistData {
			return contactValue("WFEmailAddress", "emailaddress", args)
		},
	}
	actions["phoneNumber"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "number",
				validType: String,
				noMax:     true,
			},
		},
		call: func(args []actionArgument) []plistData {
			return contactValue("WFPhoneNumber", "phonenumber", args)
		},
	}
	actions["selectContact"] = actionDefinition{
		ident: "selectcontacts",
		args: []argumentDefinition{
			{
				field:     "multiple",
				validType: Bool,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFSelectMultiple", args, 0),
			}
		},
	}
	actions["selectEmailAddress"] = actionDefinition{
		ident: "selectemail",
	}
	actions["selectPhoneNumber"] = actionDefinition{
		ident: "selectphone",
	}
	actions["getFromContact"] = actionDefinition{
		ident: "properties.contacts",
		args: []argumentDefinition{
			{
				field:     "contact",
				validType: String,
			},
			{
				field:     "property",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
				argumentValue("WFContentItemPropertyName", args, 1),
			}
		},
	}
}

func documentActions() {
	actions["getSelectedFiles"] = actionDefinition{ident: "finder.getselectedfiles"}
	actions["extractArchive"] = actionDefinition{
		ident: "unzip",
		args: []argumentDefinition{
			{
				field:     "file",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFArchive", args[0].value.(string)),
			}
		},
	}
	actions["makeArchive"] = actionDefinition{
		ident: "makezip",
		args: []argumentDefinition{
			{
				field:     "name",
				validType: String,
			},
			{
				field:     "format",
				validType: String,
			},
			{
				field:     "files",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFZIPName", args, 0),
				argumentValue("WFArchiveFormat", args, 1),
				variableInput("WFInput", args[2].value.(string)),
			}
		},
	}
	actions["quicklook"] = actionDefinition{
		ident: "previewdocument",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
			}
		},
	}
	actions["fileSize"] = actionDefinition{
		ident: "format.filesize",
		args: []argumentDefinition{
			{
				field:     "file",
				validType: Variable,
			},
			{
				field:     "format",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		args: []argumentDefinition{
			{
				field:     "detail",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFDeviceDetail", args, 0),
			}
		},
	}
	actions["getFileFrom"] = actionDefinition{
		ident: "gettypeaction",
		args: []argumentDefinition{
			{
				field:     "file",
				validType: String,
			},
			{
				field:     "from",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFFileType", args, 0),
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["getObjectOfClass"] = actionDefinition{
		ident: "getclassaction",
		args: []argumentDefinition{
			{
				field:     "class",
				validType: String,
			},
			{
				field:     "from",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("Class", args, 0),
				variableInput("Input", args[1].value.(string)),
			}
		},
	}
	actions["getOnScreenContent"] = actionDefinition{}
	actions["translateFrom"] = actionDefinition{
		ident: "translate",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
			{
				field:     "from",
				validType: String,
			},
			{
				field:     "to",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInputText", args[0].value.(string)),
				argumentValue("WFSelectedFromLanguage", args, 0),
				argumentValue("WFSelectedLanguage", args, 0),
			}
		},
	}
	actions["translate"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
			{
				field:     "to",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		args: []argumentDefinition{
			{
				field:     "input",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["replaceText"] = actionDefinition{
		ident: "text.replace",
		args: []argumentDefinition{
			{
				field:     "find",
				validType: String,
			},
			{
				field:     "replacement",
				validType: String,
			},
			{
				field:     "subject",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return replaceText(true, false, args)
		},
	}
	actions["iReplaceText"] = actionDefinition{
		ident: "text.replace",
		args: []argumentDefinition{
			{
				field:     "find",
				validType: String,
			},
			{
				field:     "replacement",
				validType: String,
			},
			{
				field:     "subject",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return replaceText(false, false, args)
		},
	}
	actions["regReplaceText"] = actionDefinition{
		ident: "text.replace",
		args: []argumentDefinition{
			{
				field:     "expression",
				validType: String,
			},
			{
				field:     "replacement",
				validType: String,
			},
			{
				field:     "subject",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return replaceText(true, true, args)
		},
	}
	actions["iRegReplaceText"] = actionDefinition{
		ident: "text.replace",
		args: []argumentDefinition{
			{
				field:     "expression",
				validType: String,
			},
			{
				field:     "replacement",
				validType: String,
			},
			{
				field:     "subject",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return replaceText(false, true, args)
		},
	}
	actions["uppercase"] = actionDefinition{
		ident: "text.changecase",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return changeCase("UPPERCASE", args)
		},
	}
	actions["lowercase"] = actionDefinition{
		ident: "text.changecase",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return changeCase("lowercase", args)
		},
	}
	actions["titleCase"] = actionDefinition{
		ident: "text.changecase",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return changeCase("Capitalize with Title Case", args)
		},
	}
	actions["capitalize"] = actionDefinition{
		ident: "text.changecase",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return changeCase("Capitalize with sentence case", args)
		},
	}
	actions["capitalizeAll"] = actionDefinition{
		ident: "text.changecase",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return changeCase("Capitalize Every Word", args)
		},
	}
	actions["alternateCase"] = actionDefinition{
		ident: "text.changecase",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return changeCase("cApItAlIzE wItH aLtErNaTiNg cAsE", args)
		},
	}
	actions["correctSpelling"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "text.split",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
			{
				field:     "glue",
				validType: String,
			},
		},
		call: textParts,
	}
	actions["combineText"] = actionDefinition{
		ident: "text.combine",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
			{
				field:     "glue",
				validType: String,
			},
		},
		call: textParts,
	}
}

func locationActions() {
	actions["getCurrentLocation"] = actionDefinition{
		ident: "location",
		call: func(args []actionArgument) []plistData {
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
}

func mediaActions() {
	actions["clearUpNext"] = actionDefinition{}
	actions["getCurrentSong"] = actionDefinition{}
	actions["latestPhotoImport"] = actionDefinition{ident: "getlatestphotoimport"}
	actions["takePhoto"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "showPreview",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     true,
				},
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFCameraCaptureShowPreview", args, 0),
			}
		},
	}
	actions["takePhotos"] = actionDefinition{
		ident: "takephoto",
		args: []argumentDefinition{
			{
				field:     "count",
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
		call: func(args []actionArgument) []plistData {
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
		args: []argumentDefinition{
			{
				field:     "video",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInputMedia", args, 0),
			}
		},
	}
	actions["takeVideo"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "camera",
				validType: String,
				defaultValue: actionArgument{
					valueType: String,
					value:     "Front",
				},
			},
			{
				field:     "quality",
				validType: String,
				defaultValue: actionArgument{
					valueType: String,
					value:     "Medium",
				},
			},
			{
				field:     "startImmediately",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
			},
		},
		call: func(args []actionArgument) []plistData {
			return argumentValues(&args, []paramMap{
				{key: "WFCameraCaptureDevice", idx: 0},
				{key: "WFCameraCaptureQuality", idx: 1},
				{key: "WFRecordingStart", idx: 2},
			})
		},
	}
	actions["setVolume"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "volume",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			args[0].value = fmt.Sprintf("0.%s", args[0].value)
			return []plistData{
				argumentValue("WFVolume", args, 0),
			}
		},
	}
}

func scriptingActions() {
	actions["setBrightness"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "brightness",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			args[0].value = fmt.Sprintf("0.%s", args[0].value)
			return []plistData{
				argumentValue("WFBrightness", args, 0),
			}
		},
	}
	actions["getName"] = actionDefinition{
		ident: "getitemname",
		args: []argumentDefinition{
			{
				field:     "item",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["setName"] = actionDefinition{
		ident: "setitemname",
		args: []argumentDefinition{
			{
				field:     "item",
				validType: Variable,
			},
			{
				field:     "name",
				validType: String,
			},
			{
				field:     "includeFileExtension",
				validType: Bool,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
				argumentValue("WFName", args, 1),
				argumentValue("WFDontIncludeFileExtension", args, 2),
			}
		},
	}
	actions["countItems"] = actionDefinition{
		ident: "count",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return count("Items", args)
		},
	}
	actions["countChars"] = actionDefinition{
		ident: "count",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return count("Characters", args)
		},
	}
	actions["countWords"] = actionDefinition{
		ident: "count",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return count("Words", args)
		},
	}
	actions["countSentences"] = actionDefinition{
		ident: "count",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return count("Sentences", args)
		},
	}
	actions["countLines"] = actionDefinition{
		ident: "count",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return count("Lines", args)
		},
	}
	actions["toggleAppearance"] = actionDefinition{
		ident: "appearance",
		call: func(args []actionArgument) []plistData {
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
		ident: "appearance",
		call: func(args []actionArgument) []plistData {
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
		ident: "appearance",
		call: func(args []actionArgument) []plistData {
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
		ident: "getmyworkflows",
	}
	actions["url"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "url",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
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
	actions["hash"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "type",
				validType: String,
			},
			{
				field:     "input",
				validType: Variable,
			},
		},
		check: func(args []actionArgument) {
			var hashType = strings.ToUpper(getArgValue(args[0]).(string))
			if !contains(hashTypes, hashType) {
				parserError(fmt.Sprintf("Invalid hash type of '%s'. Available hash types: %v", hashType, hashTypes))
			}
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFHashType", args, 0),
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["formatNumber"] = actionDefinition{
		ident: "format.number",
		args: []argumentDefinition{
			{
				field:     "number",
				validType: Integer,
			},
			{
				field:     "decimalPlaces",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFNumber", args, 0),
				argumentValue("WFNumberFormatDecimalPlaces", args, 1),
			}
		},
	}
	actions["randomNumber"] = actionDefinition{
		ident: "number.random",
		args: []argumentDefinition{
			{
				field:     "min",
				validType: Integer,
			},
			{
				field:     "max",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFRandomNumberMinimum", args, 0),
				argumentValue("WFRandomNumberMaximum", args, 1),
			}
		},
	}
	actions["base64Encode"] = actionDefinition{
		ident: "base64encode",
		args: []argumentDefinition{
			{
				field:     "encodeInput",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "base64encode",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "urlencode",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "urlencode",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "showresult",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("Text", args, 0),
			}
		},
	}
	actions["waitToReturn"] = actionDefinition{
		ident: "waittoreturn",
		call: func(args []actionArgument) []plistData {
			return []plistData{}
		},
	}
	actions["notification"] = actionDefinition{
		ident: "notification",
		args: []argumentDefinition{
			{
				field:     "body",
				validType: String,
			},
			{
				field:     "title",
				validType: String,
			},
			{
				field:     "playSound",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     true,
				},
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFNotificationActionBody", args, 0),
				argumentValue("WFNotificationActionTitle", args, 1),
				argumentValue("WFNotificationActionSound", args, 2),
			}
		},
	}
	actions["stop"] = actionDefinition{
		ident: "exit",
	}
	actions["nothing"] = actionDefinition{}
	actions["wait"] = actionDefinition{
		ident: "delay",
		args: []argumentDefinition{
			{
				field:     "seconds",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFDelayTime", args, 0),
			}
		},
	}
	actions["alert"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "alert",
				validType: String,
			},
			{
				field:     "title",
				validType: String,
			},
			{
				field:     "cancelButton",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     true,
				},
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFAlertActionMessage", args, 0),
				argumentValue("WFAlertActionTitle", args, 1),
				argumentValue("WFAlertActionCancelButtonShown", args, 2),
			}
		},
	}
	actions["askForInput"] = actionDefinition{
		ident: "ask",
		args: []argumentDefinition{
			{
				field:     "inputType",
				validType: String,
			},
			{
				field:     "prompt",
				validType: String,
			},
			{
				field:     "defaultValue",
				validType: String,
				defaultValue: actionArgument{
					valueType: String,
					value:     "",
				},
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInputType", args, 0),
				argumentValue("WFAskActionPrompt", args, 1),
				argumentValue("WFAskActionDefaultAnswer", args, 2),
			}
		},
	}
	actions["chooseFromList"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "dictionary",
				validType: Dict,
			},
			{
				field:     "prompt",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
				argumentValue("WFChooseFromListActionPrompt", args, 1),
			}
		},
	}
	actions["chooseMultipleFromList"] = actionDefinition{
		ident: "choosefromlist",
		args: []argumentDefinition{
			{
				field:     "dictionary",
				validType: Dict,
			},
			{
				field:     "prompt",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "getitemtype",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
			}
		},
	}
	actions["getKeys"] = actionDefinition{
		ident: "getvalueforkey",
		args: []argumentDefinition{
			{
				field:     "dictionary",
				validType: Dict,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "getvalueforkey",
		args: []argumentDefinition{
			{
				field:     "dictionary",
				validType: Dict,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "getvalueforkey",
		args: []argumentDefinition{
			{
				field:     "key",
				validType: String,
			},
			{
				field:     "dictionary",
				validType: Dict,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "setvalueforkey",
		args: []argumentDefinition{
			{
				field:     "key",
				validType: String,
			},
			{
				field:     "value",
				validType: Variable,
			},
			{
				field:     "dictionary",
				validType: Dict,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFDictionary", args, 0),
				argumentValue("WFDictionaryKey", args, 1),
				argumentValue("WFDictionaryValue", args, 2),
			}
		},
	}
	actions["open"] = actionDefinition{
		ident: "openworkflow",
		args: []argumentDefinition{
			{
				field:     "shortcutName",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "runworkflow",
		args: []argumentDefinition{
			{
				field:     "shortcutName",
				validType: String,
			},
			{
				field:     "output",
				validType: Variable,
			},
			{
				field:     "isSelf",
				validType: Bool,
				defaultValue: actionArgument{
					valueType: Bool,
					value:     false,
				},
			},
		},
		call: func(args []actionArgument) []plistData {
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
		args: []argumentDefinition{
			{
				field:     "listItem",
				validType: String,
				noMax:     true,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "statistics",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return calculateStatistics("Average", args)
		},
	}
	actions["calcMin"] = actionDefinition{
		ident: "statistics",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return calculateStatistics("Minimum", args)
		},
	}
	actions["calcMax"] = actionDefinition{
		ident: "statistics",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return calculateStatistics("Maximum", args)
		},
	}
	actions["calcSum"] = actionDefinition{
		ident: "statistics",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return calculateStatistics("Sum", args)
		},
	}
	actions["calcMedian"] = actionDefinition{
		ident: "statistics",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return calculateStatistics("Median", args)
		},
	}
	actions["calcMode"] = actionDefinition{
		ident: "statistics",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return calculateStatistics("Mode", args)
		},
	}
	actions["calcRange"] = actionDefinition{
		ident: "statistics",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return calculateStatistics("Range", args)
		},
	}
	actions["calcStdDevi"] = actionDefinition{
		ident: "statistics",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return calculateStatistics("Standard Deviation", args)
		},
	}
	actions["dismissSiri"] = actionDefinition{}
	actions["getLocalIP"] = actionDefinition{
		ident: "getipaddress",
		args: []argumentDefinition{
			{
				field:     "type",
				validType: String,
			},
		},
		check: checkIPType,
		call: func(args []actionArgument) []plistData {
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
		ident: "getipaddress",
		args: []argumentDefinition{
			{
				field:     "type",
				validType: String,
			},
		},
		check: checkIPType,
		call: func(args []actionArgument) []plistData {
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
		ident: "getitemfromlist",
		args: []argumentDefinition{
			{
				field:     "list",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "getitemfromlist",
		args: []argumentDefinition{
			{
				field:     "list",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "getitemfromlist",
		args: []argumentDefinition{
			{
				field:     "list",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "getitemfromlist",
		args: []argumentDefinition{
			{
				field:     "list",
				validType: Variable,
			},
			{
				field:     "index",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
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
		ident: "getitemfromlist",
		args: []argumentDefinition{
			{
				field:     "list",
				validType: Variable,
			},
			{
				field:     "start",
				validType: Integer,
			},
			{
				field:     "end",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
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
		ident: "detect.number",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["getDictionary"] = actionDefinition{
		ident: "detect.dictionary",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["getText"] = actionDefinition{
		ident: "detect.text",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["getNumber"] = actionDefinition{
		ident: "detect.number",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["getAddresses"] = actionDefinition{
		ident: "detect.address",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["getContacts"] = actionDefinition{
		ident: "detect.contacts",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["getDates"] = actionDefinition{
		ident: "detect.date",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["getEmails"] = actionDefinition{
		ident: "detect.emailaddress",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["getImages"] = actionDefinition{
		ident: "detect.images",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["getPhoneNumbers"] = actionDefinition{
		ident: "detect.phonenumber",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["getURLs"] = actionDefinition{
		ident: "detect.link",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["setWallpaper"] = actionDefinition{
		ident: "wallpaper.set",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["startScreensaver"] = actionDefinition{}
	actions["contentGraph"] = actionDefinition{
		ident: "viewresult",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["openXCallbackURL"] = actionDefinition{
		ident: "openxcallbackurl",
		args: []argumentDefinition{
			{
				field:     "url",
				validType: String,
				noMax:     true,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFXCallbackURL", args, 0),
			}
		},
	}
	actions["openCustomXCallbackURL"] = actionDefinition{
		ident: "openxcallbackurl",
		args: []argumentDefinition{
			{
				field:     "url",
				validType: String,
			},
			{
				field:     "successKey",
				validType: String,
			},
			{
				field:     "cancelKey",
				validType: String,
			},
			{
				field:     "errorKey",
				validType: String,
			},
			{
				field:     "successURL",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
			{
				field:     "response",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFOutput", args, 0),
				{
					key:      "WFNoOutputSurfaceBehavior",
					dataType: Text,
					value:    "Respond",
				},
			}
		},
	}
	actions["outputOrClipboard"] = actionDefinition{
		ident: "output",
		call: func(args []actionArgument) []plistData {
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
		ident: "wifi.set",
		args: []argumentDefinition{
			{
				field:     "status",
				validType: Bool,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("OnValue", args, 0),
			}
		},
	}
	actions["setCellularData"] = actionDefinition{
		ident: "cellulardata.set",
		args: []argumentDefinition{
			{
				field:     "status",
				validType: Bool,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("OnValue", args, 0),
			}
		},
	}
	actions["setCellularVoice"] = actionDefinition{
		ident: "cellular.rat.set",
		args: []argumentDefinition{
			{
				field:     "status",
				validType: Bool,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("OnValue", args, 0),
			}
		},
	}
	actions["toggleBluetooth"] = actionDefinition{
		ident: "bluetooth.set",
		call: func(args []actionArgument) []plistData {
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
		ident: "bluetooth.set",
		args: []argumentDefinition{
			{
				field:     "status",
				validType: Bool,
			},
		},
		call: func(args []actionArgument) []plistData {
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
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["round"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "number",
				validType: Integer,
			},
			{
				field:     "roundTo",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return roundingValue("Normal", args)
		},
	}
	actions["roundUp"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "number",
				validType: Integer,
			},
			{
				field:     "roundTo",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return roundingValue("Always Round Up", args)
		},
	}
	actions["roundDown"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "number",
				validType: Integer,
			},
			{
				field:     "roundTo",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return roundingValue("Always Round Down", args)
		},
	}
	actions["runShellScript"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "script",
				validType: String,
			},
			{
				field:     "input",
				validType: Variable,
			},
			{
				field:     "shell",
				validType: String,
				defaultValue: actionArgument{
					valueType: String,
					value:     "/bin/zsh",
				},
			},
			{
				field:     "inputMode",
				validType: String,
				defaultValue: actionArgument{
					valueType: String,
					value:     "to stdin",
				},
			},
		},
		call: func(args []actionArgument) []plistData {
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
		ident: "airdropdocument",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[1].value.(string)),
			}
		},
	}
	actions["share"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "input",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
			}
		},
	}
	actions["copyToClipboard"] = actionDefinition{
		ident: "setclipboard",
		args: []argumentDefinition{
			{
				field:     "value",
				validType: Variable,
			},
			{
				field:     "local",
				validType: Bool,
			},
			{
				field:     "expire",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
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
	actions["openURL"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "url",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
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
}

func customActions() {
	actions["makeVCard"] = actionDefinition{
		ident: "gettext",
		args: []argumentDefinition{
			{
				field:     "title",
				validType: String,
			},
			{
				field:     "subtitle",
				validType: String,
			},
			{
				field:     "url",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			var vcard = "BEGIN:VCARD\nVERSION:3.0\n"
			vcard += "N;CHARSET=utf-8:" + args[0].value.(string) + "\n"
			vcard += "ORG:" + args[1].value.(string) + "\nPHOTO;ENCODING=b:"
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
