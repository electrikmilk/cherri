/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// FIXME: Some of these actions have a value with a set list values for an arguments,
//  but the argument value is not being checked against its possible values.
//  Use the "hash" action as an example.

var measurementUnitTypes = []string{"Acceleration", "Angle", "Area", "Concentration Mass", "Dispersion", "Duration", "Electric Charge", "Electric Current", "Electric Potential Difference", "V Electric Resistance", "Energy", "Frequency", "Fuel Efficiency", "Illuminance", "Information Storage", "Length", "Mass", "Power", "Pressure", "Speed", "Temperature", "Volume"}
var units map[string][]string

func standardActions() {
	if len(actions) != 0 {
		return
	}
	actions = make(map[string]*actionDefinition)
	calendarActions()
	contactActions()
	documentActions()
	locationActions()
	mediaActions()
	scriptingActions()
	sharingActions()
	webActions()
	builtinActions()
}

func calendarActions() {
	actions["date"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
				key:       "WFDateActionDate",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFDateActionMode",
					dataType: Text,
					value:    "Specified Date",
				},
			}
		},
	}
	actions["addCalendar"] = &actionDefinition{
		identifier: "addnewcalendar",
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
				key:       "CalendarName",
			},
		},
	}
	actions["addSeconds"] = &actionDefinition{
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
	actions["addMinutes"] = &actionDefinition{
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
	actions["addHours"] = &actionDefinition{
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
	actions["addDays"] = &actionDefinition{
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
	actions["addWeeks"] = &actionDefinition{
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
	actions["addMonths"] = &actionDefinition{
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
	actions["addYears"] = &actionDefinition{
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
	actions["subtractSeconds"] = &actionDefinition{
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
	actions["subtractMinutes"] = &actionDefinition{
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
	actions["subtractHours"] = &actionDefinition{
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
	actions["subtractDays"] = &actionDefinition{
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
	actions["subtractWeeks"] = &actionDefinition{
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
	actions["subtractMonths"] = &actionDefinition{
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
	actions["subtractYears"] = &actionDefinition{
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
	actions["getStartMinute"] = &actionDefinition{
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
	actions["getStartHour"] = &actionDefinition{
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
	actions["getStartWeek"] = &actionDefinition{
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
	actions["getStartMonth"] = &actionDefinition{
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
	actions["getStartYear"] = &actionDefinition{
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
	var contactDetails = []string{"First Name", "Middle Name", "Last Name", "Birthday", "Prefix", "Suffix", "Nickname", "Phonetic First Name", "Phonetic Last Name", "Phonetic Middle Name", "Company", "Job Title", "Department", "File Extension", "Creation Date", "File Path", "Last Modified Date", "Name", "Random"}
	var abcSortOrders = []string{"A to Z", "Z to A"}
	actions["filterContacts"] = &actionDefinition{
		identifier: "filter.contacts",
		parameters: []parameterDefinition{
			{
				name:      "contacts",
				validType: Variable,
				key:       "WFContentItemInputParameter",
			},
			{
				name:      "sortByProperty",
				validType: String,
				key:       "WFContentItemSortProperty",
				enum:      contactDetails,
				optional:  true,
			},
			{
				name:         "sortOrder",
				validType:    String,
				key:          "WFContentItemSortOrder",
				defaultValue: "A to Z",
				enum:         abcSortOrders,
				optional:     true,
			},
			{
				name:      "limit",
				validType: Integer,
				key:       "WFContentItemLimitNumber",
				optional:  true,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			if len(args) == 4 {
				return []plistData{
					{
						key:      "WFContentItemLimitEnabled",
						dataType: Boolean,
						value:    true,
					},
				}
			}
			return []plistData{}
		},
	}
	actions["emailAddress"] = &actionDefinition{
		identifier: "email",
		parameters: []parameterDefinition{
			{
				name:      "email",
				validType: String,
				infinite:  true,
			},
		},
		check: func(args []actionArgument) {
			if len(args) > 1 && args[0].valueType == Variable {
				parserError("Shortcuts only allows one variable for an email address.")
			}
		},
		make: func(args []actionArgument) []plistData {
			if args[0].valueType == Variable {
				return []plistData{
					argumentValue("WFEmailAddress", args, 0),
				}
			}

			return []plistData{
				contactValue("WFEmailAddress", emailAddress, args),
			}
		},
	}
	actions["phoneNumber"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: String,
				infinite:  true,
			},
		},
		check: func(args []actionArgument) {
			if len(args) > 1 && args[0].valueType == Variable {
				parserError("Shortcuts only allows one variable for a phone number.")
			}
		},
		make: func(args []actionArgument) []plistData {
			if args[0].valueType == Variable {
				return []plistData{
					argumentValue("WFPhoneNumber", args, 0),
				}
			}
			return []plistData{
				contactValue("WFPhoneNumber", phoneNumber, args),
			}
		},
	}
	actions["selectContact"] = &actionDefinition{
		identifier: "selectcontacts",
		parameters: []parameterDefinition{
			{
				name:         "multiple",
				validType:    Bool,
				defaultValue: false,
				key:          "WFSelectMultiple",
			},
		},
	}
	actions["selectEmailAddress"] = &actionDefinition{
		identifier: "selectemail",
	}
	actions["selectPhoneNumber"] = &actionDefinition{
		identifier: "selectphone",
	}
	actions["getContactDetail"] = &actionDefinition{
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
				enum:      contactDetails,
			},
		},
	}
	actions["call"] = &actionDefinition{
		appIdentifier: "com.apple.mobilephone.call",
		parameters: []parameterDefinition{
			{
				name:      "contact",
				validType: Variable,
				key:       "WFCallContact",
			},
		},
	}
	actions["sendEmail"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "contact",
				validType: Variable,
				key:       "WFSendEmailActionToRecipients",
			},
			{
				name:      "from",
				validType: String,
				key:       "WFSendEmailActionFrom",
			},
			{
				name:      "subject",
				validType: String,
				key:       "WFSendEmailActionSubject",
			},
			{
				name:      "body",
				validType: String,
				key:       "WFSendEmailActionInputAttachments",
			},
			{
				name:         "prompt",
				validType:    Bool,
				key:          "WFSendEmailActionShowComposeSheet",
				defaultValue: true,
			},
			{
				name:         "draft",
				validType:    Bool,
				key:          "WFSendEmailActionSaveAsDraft",
				defaultValue: false,
			},
		},
	}
	actions["sendMessage"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "contact",
				validType: Variable,
				key:       "WFSendMessageActionRecipients",
			},
			{
				name:      "message",
				validType: String,
				key:       "WFSendMessageContent",
			},
			{
				name:         "prompt",
				validType:    Bool,
				key:          "ShowWhenRun",
				defaultValue: true,
			},
		},
	}
	var facetimeCallTypes = []string{"Video", "Audio"}
	actions["facetimeCall"] = &actionDefinition{
		appIdentifier: "com.apple.facetime.facetime",
		parameters: []parameterDefinition{
			{
				name:      "contact",
				validType: Variable,
				key:       "WFFaceTimeContact",
			},
			{
				name:         "type",
				validType:    String,
				key:          "WFFaceTimeType",
				defaultValue: "Video",
				enum:         facetimeCallTypes,
				optional:     true,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFFaceTimeType",
					dataType: Text,
					value:    "Video",
				},
			}
		},
	}
	actions["newContact"] = &actionDefinition{
		identifier: "addnewcontact",
		parameters: []parameterDefinition{
			{
				name:      "firstName",
				validType: String,
				key:       "WFContactFirstName",
			},
			{
				name:      "lastName",
				validType: String,
				key:       "WFContactLastName",
			},
			{
				name:      "phoneNumber",
				validType: String,
			},
			{
				name:      "emailAddress",
				validType: String,
			},
			{
				name:      "company",
				validType: String,
				key:       "WFContactCompany",
			},
			{
				name:      "notes",
				validType: String,
				key:       "WFContactNotes",
			},
			{
				name:         "prompt",
				validType:    Bool,
				key:          "ShowWhenRun",
				defaultValue: false,
				optional:     true,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			var plistDataArray = []plistData{}

			if len(args) >= 3 {
				if args[2].valueType == Variable {
					plistDataArray = append(plistDataArray, argumentValue("WFContactPhoneNumbers", args, 2))
				} else {
					plistDataArray = append(plistDataArray, contactValue("WFContactPhoneNumbers", phoneNumber, []actionArgument{args[2]}))
				}
			}

			if len(args) >= 4 {
				if args[3].valueType == Variable {
					plistDataArray = append(plistDataArray, argumentValue("WFContactEmails", args, 3))
				} else {
					plistDataArray = append(plistDataArray, contactValue("WFContactEmails", emailAddress, []actionArgument{args[3]}))
				}
			}

			return plistDataArray
		},
	}
	actions["updateContact"] = &actionDefinition{
		identifier: "setters.contacts",
		parameters: []parameterDefinition{
			{
				name:      "contact",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "detail",
				validType: String,
				key:       "WFContentItemPropertyName",
				enum:      contactDetails,
			},
			{
				name:      "value",
				validType: String,
			},
		},
		check: func(args []actionArgument) {
			var contactDetail = args[1].value.(string)
			var contactDetailKey = strings.ReplaceAll(contactDetail, " ", "")
			actions[currentAction].parameters[2].key = "WFContactContentItem" + contactDetailKey
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Mode",
					dataType: Text,
					value:    "Set",
				},
			}
		},
	}
	actions["removeContactDetail"] = &actionDefinition{
		identifier: "setters.contacts",
		parameters: []parameterDefinition{
			{
				name:      "contact",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "detail",
				validType: String,
				key:       "WFContentItemPropertyName",
				enum:      contactDetails,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Mode",
					dataType: Text,
					value:    "Remove",
				},
			}
		},
	}
}

func documentActions() {
	actions["speak"] = &actionDefinition{
		identifier: "speaktext",
		parameters: []parameterDefinition{
			{
				name:      "prompt",
				validType: String,
				key:       "WFText",
			},
			{
				name:         "waitUntilFinished",
				validType:    Bool,
				key:          "WFSpeakTextWait",
				defaultValue: true,
			},
			{
				name:      "language",
				validType: String,
				key:       "WFSpeakTextLanguage",
				optional:  true,
			},
		},
	}
	var stopListening = []string{"After Pause", "After Short Pause", "On Tap"}
	actions["listen"] = &actionDefinition{
		identifier: "dictatetext",
		check: func(args []actionArgument) {
			if len(args) != 2 {
				return
			}
			args[1].value = languageCode(getArgValue(args[1]).(string))
		},
		parameters: []parameterDefinition{
			{
				name:         "stopListening",
				validType:    String,
				key:          "WFDictateTextStopListening",
				defaultValue: "After Pause",
				enum:         stopListening,
				optional:     true,
			},
			{
				name:      "language",
				validType: String,
				key:       "WFSpeechLanguage",
				optional:  true,
			},
		},
	}
	// TODO: Writing to locations other than the Shortcuts folder.
	actions["createFolder"] = &actionDefinition{
		identifier: "file.createfolder",
		parameters: []parameterDefinition{
			{
				name:      "path",
				validType: String,
				key:       "WFFilePath",
			},
		},
	}
	actions["getFolderContents"] = &actionDefinition{
		identifier: "file.getfoldercontents",
		parameters: []parameterDefinition{
			{
				name:      "folder",
				validType: Variable,
				key:       "WFFolder",
			},
			{
				key:          "Recursive",
				name:         "recursive",
				validType:    Bool,
				defaultValue: true,
				optional:     true,
			},
		},
	}
	actions["matchedTextGroupIndex"] = &actionDefinition{
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
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetGroupType",
					dataType: Text,
					value:    "Group At Index",
				},
			}
		},
	}
	actions["getFileFromFolder"] = &actionDefinition{
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
				name:         "errorIfNotFound",
				validType:    Bool,
				key:          "WFFileErrorIfNotFound",
				defaultValue: true,
				optional:     true,
			},
		},
	}
	actions["getFile"] = &actionDefinition{
		identifier: "documentpicker.open",
		parameters: []parameterDefinition{
			{
				name:      "path",
				validType: String,
				key:       "WFGetFilePath",
			},
			{
				name:         "errorIfNotFound",
				validType:    Bool,
				key:          "WFFileErrorIfNotFound",
				defaultValue: true,
				optional:     true,
			},
		},
	}
	actions["markup"] = &actionDefinition{
		identifier: "avairyeditphoto",
		parameters: []parameterDefinition{
			{
				name:      "document",
				validType: Variable,
				key:       "WFDocument",
			},
		},
	}
	actions["rename"] = &actionDefinition{
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
	actions["reveal"] = &actionDefinition{
		identifier: "file.reveal",
		parameters: []parameterDefinition{
			{
				name:      "files",
				validType: Variable,
				key:       "WFFile",
			},
		},
	}
	actions["define"] = &actionDefinition{
		identifier: "showdefinition",
		parameters: []parameterDefinition{
			{
				name:      "word",
				validType: String,
				key:       "Word",
			},
		},
	}
	var errorCorrectionLevels = []string{"Low", "Medium", "Quartile", "High"}
	actions["makeQRCode"] = &actionDefinition{
		identifier: "generatebarcode",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
				key:       "WFText",
			},
			{
				name:         "errorCorrection",
				validType:    String,
				key:          "WFQRErrorCorrectionLevel",
				enum:         errorCorrectionLevels,
				optional:     true,
				defaultValue: "Medium",
			},
		},
	}
	actions["openNote"] = &actionDefinition{
		identifier: "shownote",
		parameters: []parameterDefinition{
			{
				name:      "note",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["splitPDF"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "pdf",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["makeHTML"] = &actionDefinition{
		identifier: "gethtmlfromrichtext",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "makeFullDocument",
				validType:    Bool,
				key:          "WFMakeFullDocument",
				defaultValue: false,
				optional:     true,
			},
		},
	}
	actions["makeMarkdown"] = &actionDefinition{
		identifier: "getmarkdownfromrichtext",
		parameters: []parameterDefinition{
			{
				name:      "richText",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["getRichTextFromHTML"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "html",
				validType: Variable,
				key:       "WFHTML",
			},
		},
	}
	actions["getRichTextFromMarkdown"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "markdown",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["print"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["selectFile"] = &actionDefinition{
		identifier: "file.select",
		parameters: []parameterDefinition{
			{
				name:         "multiple",
				validType:    Bool,
				key:          "SelectMultiple",
				defaultValue: false,
				optional:     true,
			},
		},
	}
	actions["getFileLink"] = &actionDefinition{
		identifier: "file.getlink",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFFile",
			},
		},
	}
	actions["getParentDirectory"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["getEmojiName"] = &actionDefinition{
		identifier: "getnameofemoji",
		parameters: []parameterDefinition{
			{
				name:      "emoji",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["getFileDetail"] = &actionDefinition{
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
	actions["deleteFiles"] = &actionDefinition{
		identifier: "file.delete",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
			{
				name:         "immediately",
				key:          "WFDeleteImmediatelyDelete",
				validType:    Bool,
				defaultValue: false,
				optional:     true,
			},
		},
	}
	actions["getTextFromImage"] = &actionDefinition{
		identifier: "extracttextfromimage",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFImage",
			},
		},
	}
	actions["connectToServer"] = &actionDefinition{
		identifier: "connecttoservers",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["appendNote"] = &actionDefinition{
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
	actions["addToBooks"] = &actionDefinition{
		appIdentifier: "com.apple.iBooksX.openin",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "BooksInput",
			},
		},
	}
	actions["saveFile"] = &actionDefinition{
		identifier: "documentpicker.save",
		parameters: []parameterDefinition{
			{
				name:      "path",
				validType: String,
				key:       "WFFileDestinationPath",
			},
			{
				name:      "content",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "overwrite",
				validType:    Bool,
				key:          "WFSaveFileOverwrite",
				defaultValue: false,
				optional:     true,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFAskWhereToSave",
					dataType: Boolean,
					value:    false,
				},
			}
		},
	}
	actions["saveFilePrompt"] = &actionDefinition{
		identifier: "documentpicker.save",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "overwrite",
				validType:    Bool,
				key:          "WFSaveFileOverwrite",
				defaultValue: false,
				optional:     true,
			},
		},
	}
	actions["getSelectedFiles"] = &actionDefinition{identifier: "finder.getselectedfiles"}
	actions["extractArchive"] = &actionDefinition{
		identifier: "unzip",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFArchive",
			},
		},
	}
	var archiveTypes = []string{".zip", ".tar.gz", ".tar.bz2", ".tar.xz", ".tar", ".gz", ".cpio", ".iso"}
	actions["makeArchive"] = &actionDefinition{
		identifier: "makezip",
		parameters: []parameterDefinition{
			{
				name:      "files",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "format",
				validType:    String,
				key:          "WFArchiveFormat",
				enum:         archiveTypes,
				optional:     true,
				defaultValue: ".zip",
			},
			{
				name:      "name",
				validType: String,
				key:       "WFZIPName",
				optional:  true,
			},
		},
	}
	actions["quicklook"] = &actionDefinition{
		identifier: "previewdocument",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["translateFrom"] = &actionDefinition{
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
	actions["translate"] = &actionDefinition{
		identifier: "text.translate",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
				key:       "WFInputText",
			},
			{
				name:      "to",
				validType: String,
				key:       "WFSelectedLanguage",
			},
		},
		check: func(args []actionArgument) {
			if args[1].valueType != Variable {
				args[1].value = languageCode(getArgValue(args[1]).(string))
			}
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFSelectedFromLanguage",
					dataType: Text,
					value:    "Detect Language",
				},
			}
		},
	}
	actions["detectLanguage"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["replaceText"] = &actionDefinition{
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
	actions["iReplaceText"] = &actionDefinition{
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
	actions["regReplaceText"] = &actionDefinition{
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
	actions["iRegReplaceText"] = &actionDefinition{
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
	actions["uppercase"] = &actionDefinition{
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
	actions["lowercase"] = &actionDefinition{
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
	actions["titleCase"] = &actionDefinition{
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
	actions["capitalize"] = &actionDefinition{
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
	actions["capitalizeAll"] = &actionDefinition{
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
	actions["alternateCase"] = &actionDefinition{
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
	actions["correctSpelling"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
				key:       "text",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Show-text",
					dataType: Boolean,
					value:    true,
				},
			}
		},
	}
	actions["splitText"] = &actionDefinition{
		identifier: "text.split",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
			},
			{
				name:      "separator",
				validType: String,
			},
		},
		make: textParts,
	}
	actions["joinText"] = &actionDefinition{
		identifier: "text.combine",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: Variable,
			},
			{
				name:      "glue",
				validType: String,
			},
		},
		make: textParts,
	}
	actions["makeDiskImage"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
				key:       "VolumeName",
			},
			{
				name:      "contents",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "encrypt",
				validType:    Bool,
				key:          "EncryptImage",
				optional:     true,
				defaultValue: false,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
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
	var storageUnits = []string{"bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	actions["makeSizedDiskImage"] = &actionDefinition{
		identifier: "makediskimage",
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
				key:       "VolumeName",
			},
			{
				name:      "contents",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "size",
				validType:    String,
				defaultValue: "1 GB",
			},
			{
				name:         "encrypt",
				validType:    Bool,
				defaultValue: false,
				optional:     true,
			},
		},
		check: func(args []actionArgument) {
			var size = strings.Split(getArgValue(args[2]).(string), " ")
			var storageUnitArg = actionArgument{
				valueType: String,
				value:     size[1],
			}
			checkEnum(parameterDefinition{
				name: "disk size",
				enum: storageUnits,
			}, storageUnitArg)
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
	actions["openFile"] = &actionDefinition{
		identifier: "openin",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "prompt",
				validType:    Bool,
				key:          "WFOpenInAskWhenRun",
				defaultValue: false,
				optional:     true,
			},
		},
	}
	actions["transcribeText"] = &actionDefinition{
		appIdentifier: "com.apple.ShortcutsActions.TranscribeAudioAction",
		parameters: []parameterDefinition{
			{
				name:      "audioFile",
				validType: Variable,
				key:       "audioFile",
			},
		},
		minVersion: 17,
	}
}

func locationActions() {
	actions["getCurrentLocation"] = &actionDefinition{
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
	actions["getAddresses"] = &actionDefinition{
		identifier: "detect.address",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["getCurrentWeather"] = &actionDefinition{
		identifier: "weather.currentconditions",
		parameters: []parameterDefinition{
			{
				name:         "location",
				validType:    Variable,
				key:          "WFWeatherCustomLocation",
				defaultValue: "Current Location",
				optional:     true,
			},
		},
	}
	actions["openInMaps"] = &actionDefinition{
		identifier: "searchmaps",
		parameters: []parameterDefinition{
			{
				name:      "location",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["streetAddress"] = &actionDefinition{
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
	var weatherDetails = []string{"Name", "Air Pollutants", "Air Quality Category", "Air Quality Index", "Sunset Time", "Sunrise Time", "UV Index", "Wind Direction", "Wind Speed", "Precipitation Chance", "Precipitation Amount", "Pressure", "Humidity", "Dewpoint", "Visibility", "Condition", "Feels Like", "Low", "High", "Temperature", "Location", "Date"}
	actions["getWeatherDetail"] = &actionDefinition{
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
				enum:      weatherDetails,
			},
		},
	}
	var weatherForecastTypes = []string{
		"Daily",
		"Hourly",
	}
	actions["getWeatherForecast"] = &actionDefinition{
		identifier: "weather.forecast",
		parameters: []parameterDefinition{
			{
				name:         "type",
				validType:    String,
				key:          "WFWeatherForecastType",
				enum:         weatherForecastTypes,
				optional:     true,
				defaultValue: "Daily",
			},
			{
				name:         "location",
				validType:    Variable,
				key:          "WFInput",
				defaultValue: "Current Location",
				optional:     true,
			},
		},
	}
	var locationDetails = []string{"Name", "URL", "Label", "Phone Number", "Region", "ZIP Code", "State", "City", "Street", "Altitude", "Longitude", "Latitude"}
	actions["getLocationDetail"] = &actionDefinition{
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
				enum:      locationDetails,
			},
		},
	}
	actions["getMapsLink"] = &actionDefinition{
		identifier: "",
		parameters: []parameterDefinition{
			{
				name:      "location",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["getHalfwayPoint"] = &actionDefinition{
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
	actions["removeBackground"] = &actionDefinition{
		identifier: "image.removebackground",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "crop",
				validType:    Bool,
				key:          "WFCropToBounds",
				defaultValue: false,
				optional:     true,
			},
		},
	}
	actions["clearUpNext"] = &actionDefinition{}
	actions["getCurrentSong"] = &actionDefinition{}
	actions["getLastImport"] = &actionDefinition{identifier: "getlatestphotoimport"}
	actions["getLatestBursts"] = &actionDefinition{
		identifier: "getlatestbursts",
		parameters: []parameterDefinition{
			{
				name:      "count",
				validType: Integer,
				key:       "WFGetLatestPhotoCount",
			},
		},
	}
	actions["getLatestLivePhotos"] = &actionDefinition{
		identifier: "getlatestlivephotos",
		parameters: []parameterDefinition{
			{
				name:      "count",
				validType: Integer,
				key:       "WFGetLatestPhotoCount",
			},
		},
	}
	actions["getLatestScreenshots"] = &actionDefinition{
		identifier: "getlastscreenshot",
		parameters: []parameterDefinition{
			{
				name:      "count",
				validType: Integer,
				key:       "WFGetLatestPhotoCount",
			},
		},
	}
	actions["getLatestVideos"] = &actionDefinition{
		identifier: "getlastvideo",
		parameters: []parameterDefinition{
			{
				name:      "count",
				validType: Integer,
				key:       "WFGetLatestPhotoCount",
			},
		},
	}
	actions["getLatestPhotos"] = &actionDefinition{
		identifier: "getlastphoto",
		parameters: []parameterDefinition{
			{
				name:      "count",
				validType: Integer,
				key:       "WFGetLatestPhotoCount",
			},
			{
				name:         "includeScreenshots",
				validType:    Bool,
				key:          "WFGetLatestPhotosActionIncludeScreenshots",
				defaultValue: true,
				optional:     true,
			},
		},
	}
	actions["getImages"] = &actionDefinition{
		identifier: "detect.images",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["takePhoto"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:         "count",
				validType:    Integer,
				key:          "WFPhotoCount",
				defaultValue: 1,
			},
			{
				name:         "showPreview",
				validType:    Bool,
				key:          "WFCameraCaptureShowPreview",
				defaultValue: true,
			},
		},
		check: func(args []actionArgument) {
			if len(args) == 0 {
				return
			}

			var photos = getArgValue(args[0])
			if photos == "0" {
				parserError("Number of photos to take must be greater than zero.")
			}
		},
	}
	actions["trimVideo"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "video",
				validType: Variable,
				key:       "WFInputMedia",
			},
		},
	}
	actions["takeVideo"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:         "camera",
				validType:    String,
				key:          "WFCameraCaptureDevice",
				defaultValue: "Front",
			},
			{
				name:         "quality",
				validType:    String,
				key:          "WFCameraCaptureQuality",
				defaultValue: "Medium",
			},
			{
				name:         "startImmediately",
				validType:    Bool,
				key:          "WFRecordingStart",
				defaultValue: false,
			},
		},
	}
	actions["setVolume"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "volume",
				validType: String,
				key:       "WFVolume",
			},
		},
		check: func(args []actionArgument) {
			if args[0].valueType != Variable {
				args[0].value = fmt.Sprintf("0.%s", args[0].value)
			}
		},
	}
	actions["addToMusic"] = &actionDefinition{
		identifier: "addtoplaylist",
		parameters: []parameterDefinition{
			{
				name:      "songs",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["addToPlaylist"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "playlistName",
				validType: String,
				key:       "WFPlaylistName",
			},
			{
				name:      "songs",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["playNext"] = &actionDefinition{
		identifier: "addmusictoupnext",
		parameters: []parameterDefinition{
			{
				name:      "music",
				validType: Variable,
				key:       "WFMusic",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFWhenToPlay",
					dataType: Text,
					value:    "Next",
				},
			}
		},
	}
	actions["playLater"] = &actionDefinition{
		identifier: "addmusictoupnext",
		parameters: []parameterDefinition{
			{
				name:      "music",
				validType: Variable,
				key:       "WFMusic",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFWhenToPlay",
					dataType: Text,
					value:    "Later",
				},
			}
		},
	}
	actions["createPlaylist"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "title",
				validType: String,
				key:       "WFPlaylistName",
			},
			{
				name:      "songs",
				validType: Variable,
				key:       "WFPlaylistItems",
				optional:  true,
			},
			{
				name:      "description",
				validType: String,
				key:       "WFPlaylistDescription",
				optional:  true,
			},
			{
				name:      "author",
				validType: String,
				key:       "WFPlaylistAuthor",
				optional:  true,
			},
		},
	}
	actions["addToGIF"] = &actionDefinition{
		identifier: "addframetogif",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: String,
				key:       "WFImage",
			},
			{
				name:      "gif",
				validType: String,
				key:       "WFInputGIF",
			},
			{
				name:         "delay",
				validType:    String,
				key:          "WFGIFDelayTime",
				optional:     true,
				defaultValue: "0.25",
			},
			{
				name:         "autoSize",
				validType:    Bool,
				key:          "WFGIFAutoSize",
				defaultValue: true,
				optional:     true,
			},
			{
				name:      "width",
				validType: String,
				key:       "WFGIFManualSizeWidth",
				optional:  true,
			},
			{
				name:      "height",
				validType: String,
				key:       "WFGIFManualSizeHeight",
				optional:  true,
			},
		},
	}
	actions["convertToJPEG"] = &actionDefinition{
		identifier: "image.convert",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "compressionQuality",
				validType: Integer,
				key:       "WFImageCompressionQuality",
				optional:  true,
			},
			{
				name:         "preserveMetadata",
				validType:    Bool,
				key:          "WFImagePreserveMetadata",
				optional:     true,
				defaultValue: true,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFImageFormat",
					dataType: Text,
					value:    "JPEG",
				},
			}
		},
	}
	actions["combineImages"] = &actionDefinition{
		identifier: "image.combine",
		parameters: []parameterDefinition{
			{
				name:      "images",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "mode",
				validType:    String,
				key:          "WFImageCombineMode",
				defaultValue: "vertically",
				optional:     true,
			},
			{
				name:         "spacing",
				validType:    Integer,
				key:          "WFImageCombineSpacing",
				defaultValue: 1,
				optional:     true,
			},
		},
		check: func(args []actionArgument) {
			if len(args) < 2 {
				return
			}
			var combineMode = fmt.Sprintf("%s", args[1].value)
			if strings.ToLower(combineMode) == "grid" {
				args[1].value = "In a Grid"
			}
		},
	}
	actions["rotateImage"] = &actionDefinition{
		identifier: "image.rotate",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFImage",
			},
			{
				name:      "degrees",
				validType: String,
				key:       "WFImageRotateAmount",
			},
		},
	}
	actions["selectPhotos"] = &actionDefinition{
		identifier: "selectphoto",
		parameters: []parameterDefinition{
			{
				name:         "selectMultiple",
				validType:    Bool,
				key:          "WFSelectMultiplePhotos",
				defaultValue: false,
				optional:     true,
			},
		},
	}
	actions["createAlbum"] = &actionDefinition{
		identifier: "photos.createalbum",
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
				key:       "AlbumName",
			},
			{
				name:      "images",
				validType: Variable,
				key:       "WFInput",
				optional:  true,
			},
		},
	}
	var cropPositions = []string{"Center", "Top Left", "Top Right", "Bottom Left", "Bottom Right", "Custom"}
	actions["cropImage"] = &actionDefinition{
		identifier: "image.crop",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "position",
				validType:    String,
				key:          "WFImageCropPosition",
				optional:     true,
				defaultValue: "Center",
			},
			{
				name:         "width",
				validType:    String,
				key:          "WFImageCropWidth",
				optional:     true,
				defaultValue: "100",
			},
			{
				name:         "height",
				validType:    String,
				key:          "WFImageCropHeight",
				enum:         cropPositions,
				optional:     true,
				defaultValue: "100",
			},
		},
	}
	actions["deletePhotos"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "photos",
				validType: Variable,
				key:       "photos",
			},
		},
	}
	actions["removeFromAlbum"] = &actionDefinition{
		identifier: "removefromalbum",
		parameters: []parameterDefinition{
			{
				name:      "photo",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "album",
				validType: String,
				key:       "WFRemoveAlbumSelectedGroup",
			},
		},
	}
	actions["selectMusic"] = &actionDefinition{
		identifier: "exportsong",
		parameters: []parameterDefinition{
			{
				name:         "selectMultiple",
				validType:    Bool,
				key:          "WFExportSongActionSelectMultiple",
				defaultValue: false,
				optional:     true,
			},
		},
	}
	actions["skipBack"] = &actionDefinition{
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFSkipBackBehavior",
					dataType: Text,
					value:    "Previous Song",
				},
			}
		},
	}
	actions["skipFwd"] = &actionDefinition{
		identifier: "skipforward",
	}
	actions["searchAppStore"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "WFSearchTerm",
			},
		},
	}
	actions["searchPodcasts"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "WFSearchTerm",
			},
		},
	}
	actions["makeVideoFromGIF"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "gif",
				validType: Variable,
				key:       "WFInputGIF",
			},
			{
				name:         "loops",
				validType:    Integer,
				key:          "WFMakeVideoFromGIFActionLoopCount",
				defaultValue: 1,
				optional:     true,
			},
		},
	}
	var flipDirections = []string{"Horizontal", "Vertical"}
	actions["flipImage"] = &actionDefinition{
		identifier: "image.flip",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "direction",
				validType: Variable,
				key:       "WFImageFlipDirection",
				enum:      flipDirections,
			},
		},
	}
	var recordingQualities = []string{"Normal", "Very High"}
	var recordingStarts = []string{"On Tap", "Immediately"}
	actions["recordAudio"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:         "quality",
				validType:    String,
				key:          "WFRecordingCompression",
				defaultValue: "Normal",
				enum:         recordingQualities,
				optional:     true,
			},
			{
				name:         "start",
				validType:    String,
				key:          "WFRecordingStart",
				defaultValue: "On Tap",
				enum:         recordingStarts,
				optional:     true,
			},
		},
	}
	var imageFormats = []string{"TIFF", "GIF", "PNG", "BMP", "PDF", "HEIF"}
	actions["convertImage"] = &actionDefinition{
		identifier: "image.convert",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "format",
				validType: String,
				key:       "WFImageFormat",
				enum:      imageFormats,
			},
			{
				name:         "preserveMetadata",
				validType:    Bool,
				key:          "WFImagePreserveMetadata",
				optional:     true,
				defaultValue: true,
			},
		},
	}
	var audioFormats = []string{"M4A", "AIFF"}
	var speeds = []string{"0.5X", "Normal", "2X"}
	var videoSizes = []string{"640×480", "960×540", "1280×720", "1920×1080", "3840×2160", "HEVC 1920×1080", "HEVC 3840x2160", "ProRes 422"}
	actions["encodeVideo"] = &actionDefinition{
		identifier: "encodemedia",
		parameters: []parameterDefinition{
			{
				name:      "video",
				validType: Variable,
				key:       "WFMedia",
			},
			{
				name:         "size",
				validType:    String,
				key:          "WFMediaSize",
				defaultValue: "Passthrough",
				enum:         videoSizes,
				optional:     true,
			},
			{
				name:         "speed",
				validType:    String,
				key:          "WFMediaCustomSpeed",
				defaultValue: "Normal",
				optional:     true,
				enum:         speeds,
			},
			{
				name:         "preserveTransparency",
				validType:    Bool,
				key:          "WFMediaPreserveTransparency",
				defaultValue: false,
				optional:     true,
			},
		},
	}
	actions["encodeAudio"] = &actionDefinition{
		identifier: "encodemedia",
		parameters: []parameterDefinition{
			{
				name:      "audio",
				validType: Variable,
				key:       "WFMedia",
			},
			{
				name:         "format",
				validType:    String,
				key:          "WFMediaAudioFormat",
				defaultValue: "M4A",
				enum:         audioFormats,
				optional:     true,
			},
			{
				name:         "speed",
				validType:    String,
				key:          "WFMediaCustomSpeed",
				defaultValue: "Normal",
				optional:     true,
				enum:         speeds,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			var params = []plistData{
				{
					key:      "WFMediaAudioOnly",
					dataType: Boolean,
					value:    true,
				},
			}
			return params
		},
	}
	actions["setMetadata"] = &actionDefinition{
		identifier: "encodemedia",
		parameters: []parameterDefinition{
			{
				name:      "media",
				validType: Variable,
				key:       "WFMedia",
			},
			{
				name:      "artwork",
				validType: Variable,
				key:       "WFMetadataArtwork",
				optional:  true,
			},
			{
				name:      "title",
				validType: String,
				key:       "WFMetadataTitle",
				optional:  true,
			},
			{
				name:      "artist",
				validType: String,
				key:       "WFMetadataArtist",
				optional:  true,
			},
			{
				name:      "album",
				validType: String,
				key:       "WFMetadataAlbum",
				optional:  true,
			},
			{
				name:      "genre",
				validType: String,
				key:       "WFMetadataGenre",
				optional:  true,
			},
			{
				name:      "year",
				validType: String,
				key:       "WFMetadataYear",
				optional:  true,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Metadata",
					dataType: Boolean,
					value:    true,
				},
			}
		},
	}
	actions["stripMediaMetadata"] = &actionDefinition{
		identifier: "encodemedia",
		parameters: []parameterDefinition{
			{
				name:      "media",
				validType: Variable,
				key:       "WFMedia",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Metadata",
					dataType: Boolean,
					value:    true,
				},
			}
		},
	}
	actions["stripImageMetadata"] = &actionDefinition{
		identifier: "image.convert",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFImagePreserveMetadata",
					dataType: Boolean,
					value:    false,
				},
				{
					key:      "WFImageFormat",
					dataType: Text,
					value:    "Match Input",
				},
			}
		},
	}
	actions["savePhoto"] = &actionDefinition{
		identifier: "savetocameraroll",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFCameraRollSelectedGroup",
					dataType: Text,
					value:    "Recents",
				},
			}
		},
	}
	actions["play"] = &actionDefinition{
		identifier: "pausemusic",
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFPlayPauseBehavior",
					dataType: Text,
					value:    "Play",
				},
			}
		},
	}
	actions["pause"] = &actionDefinition{
		identifier: "pausemusic",
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFPlayPauseBehavior",
					dataType: Text,
					value:    "Pause",
				},
			}
		},
	}
	actions["togglePlayPause"] = &actionDefinition{
		identifier: "pausemusic",
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFPlayPauseBehavior",
					dataType: Text,
					value:    "Play/Pause",
				},
			}
		},
	}
	actions["startShazam"] = &actionDefinition{
		identifier: "shazamMedia",
		parameters: []parameterDefinition{
			{
				name:         "show",
				validType:    Bool,
				key:          "WFShazamMediaActionShowWhenRun",
				defaultValue: true,
				optional:     true,
			},
			{
				name:         "showError",
				validType:    Bool,
				key:          "WFShazamMediaActionErrorIfNotRecognized",
				defaultValue: true,
				optional:     true,
			},
		},
	}
	actions["showIniTunes"] = &actionDefinition{
		identifier: "showinstore",
		parameters: []parameterDefinition{
			{
				name:      "product",
				validType: Variable,
				key:       "WFProduct",
			},
		},
	}
	actions["takeScreenshot"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:         "mainMonitorOnly",
				validType:    Bool,
				key:          "WFTakeScreenshotMainMonitorOnly",
				defaultValue: false,
				optional:     true,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFTakeScreenshotScreenshotType",
					dataType: Text,
					value:    "Full Screen",
				},
			}
		},
		mac: true,
	}
	var selectionTypes = []string{"Window", "Custom"}
	actions["takeInteractiveScreenshot"] = &actionDefinition{
		identifier: "takescreenshot",
		parameters: []parameterDefinition{
			{
				name:         "selection",
				validType:    String,
				key:          "WFTakeScreenshotActionInteractiveSelectionType",
				defaultValue: "Window",
				enum:         selectionTypes,
				optional:     true,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFTakeScreenshotScreenshotType",
					dataType: Text,
					value:    "Interactive",
				},
			}
		},
		mac: true,
	}
}

func scriptingActions() {
	actions["shutdown"] = &actionDefinition{
		identifier: "reboot",
		minVersion: 17,
	}
	actions["reboot"] = &actionDefinition{
		minVersion: 17,
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFShutdownMode",
					dataType: Text,
					value:    "Restart",
				},
			}
		},
	}
	actions["sleep"] = &actionDefinition{
		minVersion: 17,
		mac:        true,
	}
	actions["displaySleep"] = &actionDefinition{
		minVersion: 17,
		mac:        true,
	}
	actions["logout"] = &actionDefinition{
		minVersion: 17,
		mac:        true,
	}
	actions["lockScreen"] = &actionDefinition{
		minVersion: 17,
	}
	actions["number"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: Variable,
				key:       "WFNumberActionNumber",
			},
		},
	}
	actions["getObjectOfClass"] = &actionDefinition{
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
	actions["getOnScreenContent"] = &actionDefinition{}
	var fileSizeUnits = []string{"Closest Unit", "Bytes", "Kilobytes", "Megabytes", "Gigabytes", "Terabytes", "Petabytes", "Exabytes", "Zettabytes", "Yottabytes"}
	actions["fileSize"] = &actionDefinition{
		identifier: "format.filesize",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFFileSize",
			},
			{
				name:      "format",
				validType: String,
				key:       "WFFileSizeFormat",
				enum:      fileSizeUnits,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFFileSizeIncludeUnits",
					dataType: Boolean,
					value:    false,
				},
			}
		},
	}
	var deviceDetails = []string{"Device Name", "Device Hostname", "Device Model", "Device Is Watch", "System Version", "Screen Width", "Screen Height", "Current Volume", "Current Brightness", "Current Appearance"}
	actions["getDeviceDetail"] = &actionDefinition{
		identifier: "getdevicedetails",
		parameters: []parameterDefinition{
			{
				name:      "detail",
				key:       "WFDeviceDetail",
				validType: String,
				enum:      deviceDetails,
			},
		},
	}
	actions["setBrightness"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "brightness",
				validType: String,
				key:       "WFBrightness",
			},
		},
		check: func(args []actionArgument) {
			if args[0].valueType != Variable {
				args[0].value = fmt.Sprintf("0.%s", args[0].value)
			}
		},
	}
	actions["getName"] = &actionDefinition{
		identifier: "getitemname",
		parameters: []parameterDefinition{
			{
				name:      "item",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["setName"] = &actionDefinition{
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
				name:         "includeFileExtension",
				key:          "WFDontIncludeFileExtension",
				validType:    Bool,
				optional:     true,
				defaultValue: false,
			},
		},
	}
	actions["count"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return countParams("Items", args)
		},
	}
	actions["countChars"] = &actionDefinition{
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return countParams("Characters", args)
		},
	}
	actions["countWords"] = &actionDefinition{
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return countParams("Words", args)
		},
	}
	actions["countSentences"] = &actionDefinition{
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return countParams("Sentences", args)
		},
	}
	actions["countLines"] = &actionDefinition{
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return countParams("Lines", args)
		},
	}
	actions["toggleAppearance"] = &actionDefinition{
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
	actions["lightMode"] = &actionDefinition{
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
	actions["darkMode"] = &actionDefinition{
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
	actions["getBatteryLevel"] = &actionDefinition{}
	actions["isCharging"] = &actionDefinition{
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
	actions["connectedToCharger"] = &actionDefinition{
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
	actions["getShortcuts"] = &actionDefinition{
		identifier: "getmyworkflows",
	}
	actions["url"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				infinite:  true,
			},
		},
		make: func(args []actionArgument) []plistData {
			var urlItems []plistData
			for _, item := range args {
				urlItems = append(urlItems, paramValue("", item, String, Text))
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
	actions["addToReadingList"] = &actionDefinition{
		identifier: "readinglist",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				infinite:  true,
			},
		},
		make: func(args []actionArgument) []plistData {
			var urlItems []plistData
			for _, item := range args {
				urlItems = append(urlItems, paramValue("", item, String, Text))
			}
			return []plistData{
				{
					key:      "Show-WFURLActionURL",
					dataType: Boolean,
					value:    true,
				},
				{
					key:      "WFURL",
					dataType: Array,
					value:    urlItems,
				},
			}
		},
	}
	var hashTypes = []string{"MD5", "SHA1", "SHA256", "SHA512"}
	actions["hash"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "type",
				key:          "WFHashType",
				enum:         hashTypes,
				validType:    String,
				defaultValue: "MD5",
				optional:     true,
			},
		},
	}
	actions["formatNumber"] = &actionDefinition{
		identifier: "format.number",
		parameters: []parameterDefinition{
			{
				name:      "number",
				key:       "WFNumber",
				validType: Integer,
			},
			{
				name:         "decimalPlaces",
				key:          "WFNumberFormatDecimalPlaces",
				validType:    Integer,
				optional:     true,
				defaultValue: 2,
			},
		},
	}
	actions["randomNumber"] = &actionDefinition{
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
	actions["base64Encode"] = &actionDefinition{
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
	actions["base64Decode"] = &actionDefinition{
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
	actions["urlEncode"] = &actionDefinition{
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
	actions["urlDecode"] = &actionDefinition{
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
	actions["show"] = &actionDefinition{
		identifier: "showresult",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Text",
				validType: String,
			},
		},
	}
	actions["waitToReturn"] = &actionDefinition{}
	actions["notification"] = &actionDefinition{
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
				name:         "playSound",
				key:          "WFNotificationActionSound",
				validType:    Bool,
				defaultValue: true,
			},
		},
	}
	actions["stop"] = &actionDefinition{
		identifier: "exit",
	}
	actions["comment"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
				key:       "WFCommentActionText",
			},
		},
	}
	actions["nothing"] = &actionDefinition{}
	actions["wait"] = &actionDefinition{
		identifier: "delay",
		parameters: []parameterDefinition{
			{
				name:      "seconds",
				key:       "WFDelayTime",
				validType: Integer,
			},
		},
	}
	actions["alert"] = &actionDefinition{
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
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFAlertActionCancelButtonShown",
					dataType: Boolean,
					value:    false,
				},
			}
		},
	}
	actions["confirm"] = &actionDefinition{
		identifier: "alert",
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
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFAlertActionCancelButtonShown",
					dataType: Boolean,
					value:    true,
				},
			}
		},
	}
	var inputTypes = []string{"Text", "Number", "URL", "Date", "Time", "Date and Time"}
	actions["prompt"] = &actionDefinition{
		identifier: "ask",
		parameters: []parameterDefinition{
			{
				name:      "prompt",
				validType: String,
				key:       "WFAskActionPrompt",
			},
			{
				name:         "inputType",
				validType:    String,
				key:          "WFInputType",
				enum:         inputTypes,
				optional:     true,
				defaultValue: "Text",
			},
			{
				name:      "defaultValue",
				validType: String,
				optional:  true,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			if len(args) < 3 {
				return []plistData{}
			}
			var defaultAnswer = []plistData{
				argumentValue("WFAskActionDefaultAnswer", args, 2),
			}
			if getArgValue(args[1]) == "Number" {
				defaultAnswer = append(defaultAnswer, paramValue("WFAskActionDefaultAnswerNumber", args[2], Integer, Number))
			}

			return defaultAnswer
		},
	}
	actions["chooseFromList"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "list",
				key:       "WFInput",
				validType: Variable,
			},
			{
				name:      "prompt",
				key:       "WFChooseFromListActionPrompt",
				validType: String,
				optional:  true,
			},
			{
				name:         "selectMultiple",
				key:          "WFChooseFromListActionSelectMultiple",
				validType:    Bool,
				optional:     true,
				defaultValue: false,
			},
			{
				name:         "selectAll",
				key:          "WFChooseFromListActionSelectAll",
				validType:    Bool,
				optional:     true,
				defaultValue: false,
			},
		},
	}
	actions["typeOf"] = &actionDefinition{
		identifier: "getitemtype",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["getKeys"] = &actionDefinition{
		identifier: "getvalueforkey",
		parameters: []parameterDefinition{
			{
				name:      "dictionary",
				validType: Dict,
				key:       "WFInput",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetDictionaryValueType",
					dataType: Text,
					value:    "All Keys",
				},
			}
		},
	}
	actions["getValues"] = &actionDefinition{
		identifier: "getvalueforkey",
		parameters: []parameterDefinition{
			{
				name:      "dictionary",
				validType: Dict,
				key:       "WFInput",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetDictionaryValueType",
					dataType: Text,
					value:    "All Values",
				},
			}
		},
	}
	actions["getValue"] = &actionDefinition{
		identifier: "getvalueforkey",
		parameters: []parameterDefinition{
			{
				name:      "dictionary",
				validType: Dict,
				key:       "WFInput",
			},
			{
				name:      "key",
				validType: String,
				key:       "WFDictionaryKey",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetDictionaryValueType",
					dataType: Text,
					value:    "Value",
				},
			}
		},
	}
	actions["setValue"] = &actionDefinition{
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
	actions["openApp"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "appID",
				validType: String,
			},
		},
		check: func(args []actionArgument) {
			replaceAppIDs(args)
		},
		make: func(args []actionArgument) (params []plistData) {
			params = []plistData{
				argumentValue("WFAppIdentifier", args, 0),
			}

			if args[0].valueType == Variable {
				params = append(params, argumentValue("WFSelectedApp", args, 0))
			} else {
				params = append(params, plistData{
					key:      "WFSelectedApp",
					dataType: Dictionary,
					value: []plistData{
						argumentValue("BundleIdentifier", args, 0),
					},
				})
			}

			return
		},
	}
	actions["hideApp"] = &actionDefinition{
		identifier: "hide.app",
		parameters: []parameterDefinition{
			{
				name:      "appID",
				validType: String,
			},
		},
		check: func(args []actionArgument) {
			replaceAppIDs(args)
		},
		make: func(args []actionArgument) []plistData {
			if args[0].valueType == Variable {
				return []plistData{
					argumentValue("WFApp", args, 0),
				}
			} else {
				return []plistData{
					{
						key:      "WFApp",
						dataType: Dictionary,
						value: []plistData{
							argumentValue("BundleIdentifier", args, 0),
						},
					},
				}
			}
		},
	}
	actions["hideAllApps"] = &actionDefinition{
		identifier: "hide.app",
		parameters: []parameterDefinition{
			{
				name:      "except",
				validType: String,
				optional:  true,
				infinite:  true,
			},
		},
		check: replaceAppIDs,
		make: func(args []actionArgument) (params []plistData) {
			params = []plistData{
				{
					key:      "WFHideAppMode",
					dataType: Text,
					value:    "All Apps",
				},
			}

			if args[0].valueType != Variable {
				params = append(params, plistData{
					key:      "WFAppsExcept",
					dataType: Array,
					value:    apps(args),
				})
			} else {
				params = append(params, argumentValue("WFAppsExcept", args, 0))
			}

			return
		},
	}
	actions["quitApp"] = &actionDefinition{
		identifier: "quit.app",
		parameters: []parameterDefinition{
			{
				name:      "appID",
				validType: String,
			},
		},
		check: func(args []actionArgument) {
			replaceAppIDs(args)
		},
		make: func(args []actionArgument) []plistData {
			if args[0].valueType == Variable {
				return []plistData{
					argumentValue("WFApp", args, 0),
				}
			} else {
				return []plistData{
					{
						key:      "WFApp",
						dataType: Dictionary,
						value: []plistData{
							argumentValue("BundleIdentifier", args, 0),
						},
					},
				}
			}
		},
	}
	actions["quitAllApps"] = &actionDefinition{
		identifier: "quit.app",
		parameters: []parameterDefinition{
			{
				name:      "except",
				validType: String,
				optional:  true,
				infinite:  true,
			},
		},
		check: replaceAppIDs,
		make: func(args []actionArgument) (params []plistData) {
			params = []plistData{
				{
					key:      "WFQuitAppMode",
					dataType: Text,
					value:    "All Apps",
				},
			}

			if args[0].valueType != Variable {
				params = append(params, plistData{
					key:      "WFAppsExcept",
					dataType: Array,
					value:    apps(args),
				})
			} else {
				params = append(params, argumentValue("WFAppsExcept", args, 0))
			}

			return
		},
	}
	actions["killApp"] = &actionDefinition{
		identifier: "quit.app",
		parameters: []parameterDefinition{
			{
				name:      "appID",
				validType: String,
			},
		},
		check: func(args []actionArgument) {
			replaceAppIDs(args)
		},
		make: func(args []actionArgument) (params []plistData) {
			params = []plistData{
				{
					key:      "WFAskToSaveChanges",
					dataType: Boolean,
					value:    false,
				},
			}

			if args[0].valueType == Variable {
				return []plistData{
					argumentValue("WFApp", args, 0),
				}
			} else {
				return []plistData{
					{
						key:      "WFApp",
						dataType: Dictionary,
						value: []plistData{
							argumentValue("BundleIdentifier", args, 0),
						},
					},
				}
			}
		},
	}
	actions["killAllApps"] = &actionDefinition{
		identifier: "quit.app",
		parameters: []parameterDefinition{
			{
				name:      "except",
				validType: String,
				optional:  true,
				infinite:  true,
			},
		},
		check: replaceAppIDs,
		make: func(args []actionArgument) (params []plistData) {
			params = []plistData{
				{
					key:      "WFQuitAppMode",
					dataType: Text,
					value:    "All Apps",
				},
				{
					key:      "WFAskToSaveChanges",
					dataType: Boolean,
					value:    false,
				},
			}

			if args[0].valueType != Variable {
				params = append(params, plistData{
					key:      "WFAppsExcept",
					dataType: Array,
					value:    apps(args),
				})
			} else {
				params = append(params, argumentValue("WFAppsExcept", args, 0))
			}

			return
		},
	}
	var appSplitRatios = []string{"half", "thirdByTwo"}
	actions["splitApps"] = &actionDefinition{
		identifier: "splitscreen",
		parameters: []parameterDefinition{
			{
				name:      "firstAppID",
				validType: String,
			},
			{
				name:      "secondAppID",
				validType: String,
			},
			{
				name:         "ratio",
				validType:    String,
				optional:     true,
				enum:         appSplitRatios,
				defaultValue: "half",
			},
		},
		check: func(args []actionArgument) {
			if args[0].valueType != Variable {
				args[0].value = replaceAppID(getArgValue(args[0]).(string))
			}
			if args[1].valueType != Variable {
				args[1].value = replaceAppID(getArgValue(args[1]).(string))
			}
			if len(args) > 2 {
				if args[2].valueType == Variable {
					return
				}
				switch args[2].value {
				case "half":
					args[2].value = "½ + ½"
				case "thirdByTwo":
					args[2].value = "⅓ + ⅔"
				}
			}
		},
		make: func(args []actionArgument) (params []plistData) {
			params = []plistData{
				argumentValue("WFAppRatio", args, 2),
			}

			if args[0].valueType == Variable {
				params = append(params, argumentValue("WFPrimaryAppIdentifier", args, 0))
			} else {
				params = append(params, plistData{
					key:      "WFPrimaryAppIdentifier",
					dataType: Dictionary,
					value: []plistData{
						argumentValue("BundleIdentifier", args, 0),
					},
				})
			}

			if args[0].valueType == Variable {
				params = append(params, argumentValue("WFSecondaryAppIdentifier", args, 0))
			} else {
				params = append(params, plistData{
					key:      "WFSecondaryAppIdentifier",
					dataType: Dictionary,
					value: []plistData{
						argumentValue("BundleIdentifier", args, 0),
					},
				})
			}

			return
		},
	}
	actions["open"] = &actionDefinition{
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
	actions["runSelf"] = &actionDefinition{
		identifier: "runworkflow",
		parameters: []parameterDefinition{
			{
				name:      "output",
				validType: Variable,
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
							value:    true,
						},
						{
							key:      "workflowName",
							dataType: Text,
							value:    workflowName,
						},
					},
				},
				argumentValue("WFInput", args, 0),
			}
		},
	}
	actions["run"] = &actionDefinition{
		identifier: "runworkflow",
		parameters: []parameterDefinition{
			{
				name:      "shortcutName",
				validType: String,
			},
			{
				name:      "output",
				validType: Variable,
				optional:  true,
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
				argumentValue("WFInput", args, 1),
			}
		},
	}
	actions["list"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "listItem",
				validType: String,
				infinite:  true,
			},
		},
		make: func(args []actionArgument) []plistData {
			var listItems []plistData
			for _, item := range args {
				listItems = append(listItems, plistData{
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "WFItemType",
							dataType: Number,
							value:    0,
						},
						paramValue("WFValue", item, String, Text),
					},
				})
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
	var calculationOperations = []string{"x^2", "х^3", "x^у", "e^x", "10^x", "In(x)", "log(x)", "√x", "∛x", "x!", "sin(x)", "cos(X)", "tan(x)", "abs(x)"}
	actions["calculate"] = &actionDefinition{
		identifier: "math",
		parameters: []parameterDefinition{
			{
				name:      "operation",
				validType: String,
				enum:      calculationOperations,
				key:       "WFScientificMathOperation",
			},
			{
				name:      "operandOne",
				validType: Integer,
				key:       "WFInput",
			},
			{
				name:      "operandTwo",
				validType: Integer,
				key:       "WFMathOperand",
				optional:  true,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFMathOperation",
					dataType: Text,
					value:    "...",
				},
			}
		},
	}
	var statisticsOperations = []string{"Average", "Minimum", "Maximum", "Sum", "Median", "Mode", "Range", "Standard Deviation"}
	actions["statistic"] = &actionDefinition{
		identifier: "statistics",
		parameters: []parameterDefinition{
			{
				name:      "operation",
				validType: String,
				key:       "WFStatisticsOperation",
				enum:      statisticsOperations,
			},
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
	actions["dismissSiri"] = &actionDefinition{}
	actions["isOnline"] = &actionDefinition{
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
	var ipTypes = []string{"IPv4", "IPv6"}
	actions["getLocalIP"] = &actionDefinition{
		identifier: "getipaddress",
		parameters: []parameterDefinition{
			{
				name:         "type",
				validType:    String,
				key:          "WFIPAddressTypeOption",
				enum:         ipTypes,
				defaultValue: "IPv4",
				optional:     true,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFIPAddressSourceOption",
					dataType: Text,
					value:    "Local",
				},
			}
		},
	}
	actions["getExternalIP"] = &actionDefinition{
		identifier: "getipaddress",
		parameters: []parameterDefinition{
			{
				name:         "type",
				validType:    String,
				key:          "WFIPAddressTypeOption",
				enum:         ipTypes,
				defaultValue: "IPv4",
				optional:     true,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFIPAddressSourceOption",
					dataType: Text,
					value:    "External",
				},
			}
		},
	}
	actions["firstListItem"] = &actionDefinition{
		identifier: "getitemfromlist",
		parameters: []parameterDefinition{
			{
				name:      "list",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "First Item",
				},
			}
		},
	}
	actions["lastListItem"] = &actionDefinition{
		identifier: "getitemfromlist",
		parameters: []parameterDefinition{
			{
				name:      "list",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "Last Item",
				},
			}
		},
	}
	actions["randomListItem"] = &actionDefinition{
		identifier: "getitemfromlist",
		parameters: []parameterDefinition{
			{
				name:      "list",
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				variableInput("WFInput", args[0].value.(string)),
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "Random Item",
				},
			}
		},
	}
	actions["getListItem"] = &actionDefinition{
		identifier: "getitemfromlist",
		parameters: []parameterDefinition{
			{
				name:      "list",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "index",
				validType: Integer,
				key:       "WFItemIndex",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "Item At Index",
				},
			}
		},
	}
	actions["getListItems"] = &actionDefinition{
		identifier: "getitemfromlist",
		parameters: []parameterDefinition{
			{
				name:      "list",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "start",
				validType: Integer,
				key:       "WFItemRangeStart",
			},
			{
				name:      "end",
				validType: Integer,
				key:       "WFItemRangeEnd",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "Items in Range",
				},
			}
		},
	}
	actions["getNumbers"] = &actionDefinition{
		identifier: "detect.number",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["getDictionary"] = &actionDefinition{
		identifier: "detect.dictionary",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["getText"] = &actionDefinition{
		identifier: "detect.text",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["getContacts"] = &actionDefinition{
		identifier: "detect.contacts",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["getDates"] = &actionDefinition{
		identifier: "detect.date",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["getEmails"] = &actionDefinition{
		identifier: "detect.emailaddress",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: String,
			},
		},
	}
	actions["getPhoneNumbers"] = &actionDefinition{
		identifier: "detect.phonenumber",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["getURLs"] = &actionDefinition{
		identifier: "detect.link",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["getAllWallpapers"] = &actionDefinition{
		identifier: "posters.get",
		minVersion: 16.2,
	}
	actions["getWallpaper"] = &actionDefinition{
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
	actions["setWallpaper"] = &actionDefinition{
		identifier: "wallpaper.set",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["startScreensaver"] = &actionDefinition{mac: true}
	actions["contentGraph"] = &actionDefinition{
		identifier: "viewresult",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["openXCallbackURL"] = &actionDefinition{
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
	actions["openCustomXCallbackURL"] = &actionDefinition{
		identifier: "openxcallbackurl",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFXCallbackURL",
			},
			{
				name:      "successKey",
				validType: String,
				key:       "WFXCallbackCustomSuccessKey",
				optional:  true,
			},
			{
				name:      "cancelKey",
				validType: String,
				key:       "WFXCallbackCustomCancelKey",
				optional:  true,
			},
			{
				name:      "errorKey",
				validType: String,
				key:       "WFXCallbackCustomErrorKey",
				optional:  true,
			},
			{
				name:      "successURL",
				validType: String,
				key:       "WFXCallbackCustomSuccessURL",
				optional:  true,
			},
		},
		addParams: func(args []actionArgument) (xCallbackParams []plistData) {
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
			return
		},
	}
	actions["output"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "output",
				validType: String,
				key:       "WFOutput",
			},
		},
	}
	actions["mustOutput"] = &actionDefinition{
		identifier: "output",
		parameters: []parameterDefinition{
			{
				name:      "output",
				validType: String,
				key:       "WFOutput",
			},
			{
				name:      "response",
				validType: String,
				key:       "WFResponse",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFNoOutputSurfaceBehavior",
					dataType: Text,
					value:    "Respond",
				},
			}
		},
	}
	actions["outputOrClipboard"] = &actionDefinition{
		identifier: "output",
		parameters: []parameterDefinition{
			{
				name:      "output",
				validType: String,
				key:       "WFOutput",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFNoOutputSurfaceBehavior",
					dataType: Text,
					value:    "Copy to Clipboard",
				},
			}
		},
	}
	var focusModes = plistData{
		key:      "FocusModes",
		dataType: Dictionary,
		value: []plistData{
			{
				key:      "Identifier",
				dataType: Text,
				value:    "com.apple.donotdisturb.mode.default",
			},
			{
				key:      "DisplayString",
				dataType: Text,
				value:    "Do Not Disturb",
			},
		},
	}
	actions["DNDOn"] = &actionDefinition{
		identifier: "dnd.set",
		make: func(args []actionArgument) []plistData {
			return []plistData{
				focusModes,
				{
					key:      "Enabled",
					dataType: Number,
					value:    1,
				},
			}
		},
	}
	actions["DNDOff"] = &actionDefinition{
		identifier: "dnd.set",
		make: func(args []actionArgument) []plistData {
			return []plistData{
				focusModes,
				{
					key:      "Enabled",
					dataType: Number,
					value:    0,
				},
			}
		},
	}
	actions["setWifi"] = &actionDefinition{
		identifier: "wifi.set",
		parameters: []parameterDefinition{
			{
				name:      "status",
				key:       "OnValue",
				validType: Bool,
			},
		},
	}
	actions["setCellularData"] = &actionDefinition{
		identifier: "cellulardata.set",
		parameters: []parameterDefinition{
			{
				name:         "status",
				key:          "OnValue",
				validType:    Bool,
				defaultValue: true,
			},
		},
	}
	actions["toggleBluetooth"] = &actionDefinition{
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
	actions["setBluetooth"] = &actionDefinition{
		identifier: "bluetooth.set",
		parameters: []parameterDefinition{
			{
				name:      "status",
				validType: Bool,
				key:       "OnValue",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "operation",
					dataType: Text,
					value:    "set",
				},
			}
		},
	}
	actions["playSound"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["round"] = &actionDefinition{
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
	actions["ceil"] = &actionDefinition{
		identifier: "round",
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
	actions["floor"] = &actionDefinition{
		identifier: "round",
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
	actions["runShellScript"] = &actionDefinition{
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
				name:         "shell",
				key:          "Shell",
				validType:    String,
				defaultValue: "/bin/zsh",
			},
			{
				name:         "inputMode",
				key:          "InputMode",
				validType:    String,
				defaultValue: "to stdin",
			},
		},
	}
	actions["makeShortcut"] = &actionDefinition{
		appIdentifier: "com.apple.shortcuts.CreateWorkflowAction",
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
				key:       "name",
			},
			{
				name:         "open",
				validType:    Bool,
				key:          "OpenWhenRun",
				defaultValue: true,
				optional:     true,
			},
		},
		minVersion: 16.4,
	}
	actions["searchShortcuts"] = &actionDefinition{
		appIdentifier: "com.apple.shortcuts.SearchShortcutsAction",
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "searchPhrase",
			},
		},
		minVersion: 16.4,
	}
	var shortcutDetails = []string{"Folder", "Icon", "Action Count", "File Size", "File Extension", "Creation Date", "File Path", "Last Modified Date", "Name"}
	actions["shortcutDetail"] = &actionDefinition{
		identifier: "properties.workflow",
		parameters: []parameterDefinition{
			{
				name:      "detail",
				validType: String,
				key:       "WFContentItemPropertyName",
				enum:      shortcutDetails,
			},
			{
				name:      "shortcut",
				validType: Variable,
				key:       "WFInput",
			},
		},
	}
}

func sharingActions() {
	actions["airdrop"] = &actionDefinition{
		identifier: "airdropdocument",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	}
	actions["share"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: String,
			},
		},
	}
	actions["copyToClipboard"] = &actionDefinition{
		identifier: "setclipboard",
		parameters: []parameterDefinition{
			{
				name:      "value",
				key:       "WFInput",
				validType: Variable,
			},
			{
				name:         "local",
				key:          "WFLocalOnly",
				validType:    Bool,
				optional:     true,
				defaultValue: false,
			},
			{
				name:      "expire",
				key:       "WFExpirationDate",
				validType: String,
				optional:  true,
			},
		},
	}
	actions["getClipboard"] = &actionDefinition{}
}

func webActions() {
	actions["getURLHeaders"] = &actionDefinition{
		identifier: "url.getheaders",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["openURL"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFInput",
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Show-WFInput",
					dataType: Boolean,
					value:    true,
				},
			}
		},
	}
	actions["runJavaScriptOnWebpage"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "javascript",
				validType: String,
				key:       "WFJavaScript",
			},
		},
	}
	var engines = []string{"Amazon", "Bing", "DuckDuckGo", "eBay", "Google", "Reddit", "Twitter", "Yahoo!", "YouTube"}
	actions["searchWeb"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "engine",
				validType: String,
				key:       "WFSearchWebDestination",
				enum:      engines,
			},
			{
				name:      "query",
				validType: String,
				key:       "WFInputText",
			},
		},
	}
	actions["showWebpage"] = &actionDefinition{
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
	actions["getRSSFeeds"] = &actionDefinition{
		identifier: "rss.extract",
		parameters: []parameterDefinition{
			{
				name:      "urls",
				validType: String,
				key:       "WFURLs",
			},
		},
	}
	actions["getRSS"] = &actionDefinition{
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
	var webpageDetails = []string{"Page Contents", "Page Selection", "Page URL", "Name"}
	actions["getWebPageDetail"] = &actionDefinition{
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
				enum:      webpageDetails,
			},
		},
	}
	actions["getArticleDetail"] = &actionDefinition{
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
	actions["getCurrentURL"] = &actionDefinition{
		identifier: "safari.geturl",
	}
	actions["getWebpageContents"] = &actionDefinition{
		identifier: "getwebpagecontents",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFInput",
			},
		},
	}
	actions["searchGiphy"] = &actionDefinition{
		identifier: "giphy",
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "WFGiphyQuery",
			},
		},
	}
	actions["getGifs"] = &actionDefinition{
		identifier: "giphy",
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "WFGiphyQuery",
			},
			{
				name:         "gifs",
				validType:    Integer,
				key:          "WFGiphyLimit",
				defaultValue: 1,
				optional:     true,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGiphyShowPicker",
					dataType: Boolean,
					value:    false,
				},
			}
		},
	}
	actions["getArticle"] = &actionDefinition{
		identifier: "getarticle",
		parameters: []parameterDefinition{
			{
				name:      "webpage",
				validType: String,
				key:       "WFWebPage",
			},
		},
	}
	actions["expandURL"] = &actionDefinition{
		identifier: "url.expand",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "URL",
			},
		},
	}
	var urlComponents = []string{"Scheme", "User", "Password", "Host", "Port", "Path", "Query", "Fragment"}
	actions["getURLDetail"] = &actionDefinition{
		identifier: "geturlcomponent",
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
				enum:      urlComponents,
			},
		},
	}
	actions["downloadURL"] = &actionDefinition{
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
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFHTTPMethod",
					dataType: Text,
					value:    "GET",
				},
			}
		},
	}
	var httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	var httpParams = []parameterDefinition{
		{
			name:      "url",
			validType: String,
		},
		{
			name:         "method",
			validType:    String,
			optional:     true,
			enum:         httpMethods,
			defaultValue: "GET",
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
	actions["formRequest"] = &actionDefinition{
		identifier: "downloadurl",
		parameters: httpParams,
		make: func(args []actionArgument) []plistData {
			return httpRequest("Form", "WFFormValues", args)
		},
	}
	actions["jsonRequest"] = &actionDefinition{
		identifier: "downloadurl",
		parameters: httpParams,
		make: func(args []actionArgument) []plistData {
			return httpRequest("JSON", "WFJSONValues", args)
		},
	}
	actions["fileRequest"] = &actionDefinition{
		identifier: "downloadurl",
		parameters: httpParams,
		make: func(args []actionArgument) []plistData {
			return httpRequest("File", "WFRequestVariable", args)
		},
	}
	actions["runAppleScript"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "Input",
			},
			{
				name:      "script",
				validType: String,
				key:       "Script",
			},
		},
		mac: true,
	}
	actions["runJSAutomation"] = &actionDefinition{
		identifier: "runjavascriptforautomation",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "Input",
			},
			{
				name:      "script",
				validType: String,
				key:       "Script",
			},
		},
		mac: true,
	}
	var sortOrders = []string{"asc", "desc"}
	var windowSortings = []string{"Title", "App Name", "Width", "Height", "X Position", "Y Position", "Window Index", "Name", "Random"}
	actions["getWindows"] = &actionDefinition{
		identifier: "filter.windows",
		parameters: []parameterDefinition{
			{
				name:      "sortBy",
				validType: String,
				key:       "WFContentItemSortProperty",
				enum:      windowSortings,
				optional:  true,
			},
			{
				name:      "orderBy",
				validType: String,
				key:       "WFContentItemSortOrder",
				enum:      sortOrders,
				optional:  true,
			},
			{
				name:      "limit",
				validType: Integer,
				key:       "WFContentItemLimitNumber",
				optional:  true,
			},
		},
		addParams: func(args []actionArgument) (params []plistData) {
			if args[2].value != nil {
				params = append(params, plistData{
					key:      "WFContentItemLimitEnabled",
					dataType: Boolean,
					value:    true,
				})
			}
			return
		},
		check: func(args []actionArgument) {
			if args[1].value != nil {
				var alphabetic = []string{"Title", "App Name", "Name", "Random"}
				var numeric = []string{"Width", "Height", "X Position", "Y Position", "Window Index"}
				var sortBy = getArgValue(args[0]).(string)
				var orderBy = getArgValue(args[1]).(string)
				if sortBy != "Random" {
					if contains(alphabetic, sortBy) {
						switch orderBy {
						case "asc":
							args[1].value = "A to Z"
						case "desc":
							args[1].value = "Z to A"
						}
					} else if contains(numeric, sortBy) {
						switch orderBy {
						case "asc":
							args[1].value = "Biggest First"
						case "desc":
							args[1].value = "Smallest First"
						}
					}
				}
			}
		},
		mac: true,
	}
	var windowPositions = []string{"Top Left", "Top Center", "Top Right", "Middle Left", "Center", "Middle Right", "Bottom Left", "Bottom Center", "Bottom Right", "Coordinates"}
	actions["moveWindow"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "window",
				validType: Variable,
				key:       "WFWindow",
			},
			{
				name:      "position",
				validType: String,
				key:       "WFPosition",
				enum:      windowPositions,
			},
			{
				name:         "bringToFront",
				validType:    Bool,
				key:          "WFBringToFront",
				defaultValue: true,
				optional:     true,
			},
		},
		mac: true,
	}
	var windowConfigurations = []string{"Fit Screen", "Top Half", "Bottom Half", "Left Half", "Right Half", "Top Left Quarter", "Top Right Quarter", "Bottom Left Quarter", "Bottom Right Quarter", "Dimensions"}
	actions["resizeWindow"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "window",
				validType: Variable,
				key:       "WFWindow",
			},
			{
				name:      "configuration",
				validType: String,
				key:       "WFConfiguration",
				enum:      windowConfigurations,
				optional:  false,
				infinite:  false,
			},
		},
		mac: true,
	}
	actions["convertMeasurement"] = &actionDefinition{
		identifier: "measurement.convert",
		parameters: []parameterDefinition{
			{
				name:      "measurement",
				validType: Variable,
			},
			{
				name:      "unitType",
				validType: String,
				enum:      measurementUnitTypes,
			},
			{
				name:      "unit",
				validType: String,
			},
		},
		check: func(args []actionArgument) {
			var value = getArgValue(args[1])
			if reflect.TypeOf(value).String() != "string" {
				return
			}

			makeMeasurementUnits()

			var unitType = value.(string)
			checkEnum(parameterDefinition{
				name: "measurement unit",
				enum: units[unitType],
			}, args[2])
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
				{
					key:      "WFMeasurementUnit",
					dataType: Dictionary,
					value: []plistData{
						argumentValue("WFNSUnitSymbol", args, 2),
						argumentValue("WFNSUnitType", args, 1),
					},
				},
				argumentValue("WFMeasurementUnitType", args, 1),
			}
		},
	}
	actions["measurement"] = &actionDefinition{
		identifier: "measurement.create",
		parameters: []parameterDefinition{
			{
				name:      "magnitude",
				validType: String,
			},
			{
				name:      "unitType",
				validType: String,
				enum:      measurementUnitTypes,
			},
			{
				name:      "unit",
				validType: String,
			},
		},
		check: func(args []actionArgument) {
			var value = getArgValue(args[1])
			if reflect.TypeOf(value).String() != "string" {
				return
			}

			makeMeasurementUnits()

			var unitType = value.(string)
			checkEnum(parameterDefinition{
				name: "unit",
				enum: units[unitType],
			}, args[2])
		},
		make: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFMeasurementUnit",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "Value",
							dataType: Dictionary,
							value: []plistData{
								argumentValue("Magnitude", args, 0),
								argumentValue("Unit", args, 2),
							},
						},
						{
							key:      "WFSerializationType",
							dataType: Text,
							value:    "WFQuantityFieldValue",
						},
					},
				},
				argumentValue("WFMeasurementUnitType", args, 1),
			}
		},
	}
}

func builtinActions() {
	actions["rawAction"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "identifier",
				validType: String,
			},
			{
				name:      "parameters",
				validType: Arr,
			},
		},
		check: func(args []actionArgument) {
			actions["rawAction"].appIdentifier = getArgValue(args[0]).(string)
		},
		make: func(args []actionArgument) (params []plistData) {
			for _, parameterDefinitions := range getArgValue(args[1]).([]interface{}) {
				var paramKey string
				var paramType string
				var paramValue string
				for key, value := range parameterDefinitions.(map[string]interface{}) {
					switch key {
					case "key":
						paramKey = value.(string)
					case "type":
						paramType = value.(string)
					case "value":
						paramValue = value.(string)
					}
				}
				params = append(params, plistData{
					key:      paramKey,
					dataType: plistDataType(paramType),
					value:    paramValue,
				})
			}
			return
		},
	}
	actions["makeVCard"] = &actionDefinition{
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
				name:      "imagePath",
				validType: String,
				optional:  true,
			},
		},
		check: func(args []actionArgument) {
			if len(args) != 3 {
				return
			}

			var image = getArgValue(args[2])
			if reflect.TypeOf(image).String() != stringType {
				parserError("Image path for VCard must be a string literal")
			}
			var iconFile = getArgValue(args[2]).(string)
			if _, err := os.Stat(iconFile); os.IsNotExist(err) {
				parserError(fmt.Sprintf("File '%s' does not exist!", iconFile))
			}
		},
		make: func(args []actionArgument) []plistData {
			var title = args[0].value.(string)
			var subtitle = args[1].value.(string)
			wrapVariableReference(&title)
			wrapVariableReference(&subtitle)
			var vcard strings.Builder
			vcard.WriteString(fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nN;CHARSET=utf-8:%s\nORG:%s\n", title, subtitle))
			if len(args) > 2 {
				var iconFile = getArgValue(args[2]).(string)
				var bytes, readErr = os.ReadFile(iconFile)
				handle(readErr)
				vcard.WriteString(fmt.Sprintf("PHOTO;ENCODING=b:%s\n", base64.StdEncoding.EncodeToString(bytes)))
			}
			vcard.WriteString("END:VCARD")
			args[0] = actionArgument{
				valueType: String,
				value:     vcard.String(),
			}
			return []plistData{
				argumentValue("WFTextActionText", args, 0),
			}
		},
	}
	actions["base64File"] = &actionDefinition{
		identifier: "gettext",
		parameters: []parameterDefinition{
			{
				name:      "filePath",
				validType: String,
			},
		},
		check: func(args []actionArgument) {
			var file = getArgValue(args[0])
			if args[0].valueType == Variable && reflect.TypeOf(file).String() != stringType {
				parserError("File path must be a string literal")
			}
			if _, err := os.Stat(file.(string)); os.IsNotExist(err) {
				parserError(fmt.Sprintf("File '%s' does not exist!", file))
			}
		},
		make: func(args []actionArgument) []plistData {
			var file = getArgValue(args[0]).(string)
			var bytes, readErr = os.ReadFile(file)
			handle(readErr)
			var encodedFile = base64.StdEncoding.EncodeToString(bytes)

			return []plistData{
				{
					key:      "WFTextActionText",
					dataType: Text,
					value:    encodedFile,
				},
			}
		},
	}
}

var contactValues []plistData

type contentKit string

var emailAddress contentKit = "emailaddress"
var phoneNumber contentKit = "phonenumber"

func contactValue(key string, contentKit contentKit, args []actionArgument) plistData {
	contactValues = []plistData{}
	var entryType int
	switch contentKit {
	case emailAddress:
		entryType = 2
	case phoneNumber:
		entryType = 1
	}
	for _, item := range args {
		contactValues = append(contactValues, plistData{
			dataType: Dictionary,
			value: []plistData{
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
							key:      "link.contentkit." + string(contentKit),
							dataType: Text,
							value:    item.value,
						},
					},
				},
			},
		})
	}
	return plistData{
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

func adjustDate(operation string, unit string, args []actionArgument) (adjustDateParams []plistData) {
	adjustDateParams = []plistData{
		{
			key:      "WFAdjustOperation",
			dataType: Text,
			value:    operation,
		},
		argumentValue("WFDate", args, 0),
	}
	if unit == "" {
		return adjustDateParams
	}

	var magnitudeValue = argumentValue("Magnitude", args, 1)
	if magnitudeValue.dataType == Dictionary {
		var value = magnitudeValue.value.([]plistData)
		magnitudeValue.dataType = Dictionary
		magnitudeValue.value = value[0].value
	}
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
					magnitudeValue,
				},
			},
			{
				key:      "WFSerializationType",
				dataType: Text,
				value:    "WFQuantityFieldValue",
			},
		},
	})

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
	var data = []plistData{
		{
			key:      "Show-text",
			dataType: Boolean,
			value:    true,
		},

		argumentValue("text", args, 0),
	}

	var separator = getArgValue(args[1])
	switch {
	case separator == " ":
		data = append(data, plistData{
			key:      "WFTextSeparator",
			dataType: Text,
			value:    "Spaces",
		})
	case separator == "\n":
		data = append(data, plistData{
			key:      "WFTextSeparator",
			dataType: Text,
			value:    "New Lines",
		})
	case separator == "" && currentAction == "splitText":
		data = append(data, plistData{
			key:      "WFTextSeparator",
			dataType: Text,
			value:    "Every Character",
		})
	default:
		data = append(data,
			plistData{
				key:      "WFTextSeparator",
				dataType: Text,
				value:    "Custom",
			},
			argumentValue("WFTextCustomSeparator", args, 1),
		)
	}

	return data
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
	if lang, found := languages[language]; found {
		return lang
	}

	parserError(fmt.Sprintf("Unknown language '%s'", language))
	return ""
}

func countParams(countType string, args []actionArgument) []plistData {
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
	appIds = map[string]string{
		"appstore":  "com.apple.AppStore",
		"files":     "com.apple.DocumentsApp",
		"shortcuts": "is.workflow.my.app",
		"safari":    "com.apple.mobilesafari",
		"facetime":  "com.apple.facetime",
		"notes":     "com.apple.mobilenotes",
		"phone":     "com.apple.mobilephone",
		"reminders": "com.apple.reminders",
		"mail":      "com.apple.mobilemail",
		"music":     "com.apple.Music",
		"calendar":  "com.apple.mobilecal",
		"maps":      "com.apple.Maps",
		"contacts":  "com.apple.MobileAddressBook",
		"health":    "com.apple.Health",
		"photos":    "com.apple.mobileslideshow",
	}
}

func apps(args []actionArgument) (apps []plistData) {
	for _, arg := range args {
		if arg.valueType != Variable {
			apps = append(apps, plistData{
				dataType: Dictionary,
				value: []plistData{
					{
						key:      "BundleIdentifier",
						dataType: Text,
						value:    arg.value,
					},
					{
						key:      "TeamIdentifier",
						dataType: Text,
						value:    "0000000000",
					},
				},
			})
		}
	}
	return
}

var appIdentifierRegex = regexp.MustCompile(`^([A-Za-z][A-Za-z\d_]*\.)+[A-Za-z][A-Za-z\d_]*$`)

func replaceAppID(id string) string {
	makeAppIds()
	if appID, found := appIds[id]; found {
		return appID
	}

	var matches = appIdentifierRegex.FindAllString(id, -1)
	if len(matches) == 0 {
		parserError(fmt.Sprintf("Invalid app bundle identifier: %s", id))
	}
	return id
}

func replaceAppIDs(args []actionArgument) {
	if len(args) >= 1 {
		for a := range args {
			if args[a].valueType == Variable {
				continue
			}

			var id = getArgValue(args[a]).(string)
			args[a].value = replaceAppID(id)
		}
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
