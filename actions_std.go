/*
 * Copyright (c) Cherri
 */

package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/google/uuid"
)

const SetVariableIdentifier = "is.workflow.actions.setvariable"
const AppendVariableIdentifier = "is.workflow.actions.appendvariable"

var weekdays = []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

var httpParams = []parameterDefinition{
	{
		name:      "url",
		key:       "WFURL",
		validType: String,
	},
	{
		name:      "method",
		key:       "WFHTTPMethod",
		validType: String,
		optional:  true,
		enum:      "httpMethod",
	},
	{
		name:      "body",
		validType: Dict,
		optional:  true,
		literal:   true,
	},
	{
		name:      "headers",
		key:       "WFHTTPHeaders",
		validType: Dict,
		optional:  true,
		literal:   true,
	},
}
var fileLabelsMap = map[string]int{
	"red":    6,
	"orange": 7,
	"yellow": 5,
	"green":  2,
	"blue":   4,
	"purple": 3,
	"gray":   1,
}

var toggleAlarmIntent = appIntent{
	name:                "Clock",
	bundleIdentifier:    "com.apple.clock",
	appIntentIdentifier: "ToggleAlarmIntent",
}

var createShortcutiCloudLink = appIntent{
	name:                "Shortcuts",
	bundleIdentifier:    "com.apple.shortcuts",
	appIntentIdentifier: "CreateShortcutiCloudLinkAction",
}

// actions is the data structure that determines every action the compiler knows about.
// The key determines the identifier of the identifier that must be used in the syntax, it's value defines its behavior, etc. using an actionDefinition.
var actions = map[string]*actionDefinition{
	"returnToHomescreen": {mac: false},
	"addSeconds": {
		identifier:    "adjustdate",
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Add", "sec", args)
		},
	},
	"addMinutes": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Add", "min", args)
		},
	},
	"addHours": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Add", "hr", args)
		},
	},
	"addDays": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Add", "days", args)
		},
	},
	"addWeeks": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Add", "weeks", args)
		},
	},
	"addMonths": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Add", "months", args)
		},
	},
	"addYears": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Add", "yr", args)
		},
	},
	"subtractSeconds": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Subtract", "sec", args)
		},
	},
	"subtractMinutes": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Subtract", "min", args)
		},
	},
	"subtractHours": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Subtract", "hr", args)
		},
	},
	"subtractDays": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Subtract", "days", args)
		},
	},
	"subtractWeeks": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Subtract", "weeks", args)
		},
	},
	"subtractMonths": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Subtract", "months", args)
		},
	},
	"subtractYears": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
			{
				name:      "magnitude",
				validType: Integer,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Subtract", "yr", args)
		},
	},
	"getStartMinute": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Get Start of Minute", "", args)
		},
	},
	"getStartHour": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Get Start of Hour", "", args)
		},
	},
	"getStartWeek": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Get Start of Week", "", args)
		},
	},
	"getStartMonth": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Get Start of Month", "", args)
		},
	},
	"getStartYear": {
		identifier: "adjustdate",
		parameters: []parameterDefinition{
			{
				name:      "date",
				key:       "WFDate",
				validType: String,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return adjustDate("Get Start of Year", "", args)
		},
	},
	"startTimer": {
		identifier: "timer.start",
		parameters: []parameterDefinition{
			{
				name:      "magnitude",
				validType: Integer,
			},
			{
				name:         "unit",
				validType:    String,
				defaultValue: "min",
				enum:         "timerDuration",
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"WFDuration": magnitudeValue(argumentValue(args, 1), args, 0),
			}
		},
	},
	"createAlarm": {
		appIdentifier: "com.apple.mobiletimer-framework",
		identifier:    "MobileTimerIntents.MTCreateAlarmIntent",
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
				key:       "name",
			},
			{
				name:      "time",
				validType: String,
				key:       "dateComponents",
			},
			{
				name:         "allowsSnooze",
				validType:    Bool,
				key:          "allowsSnooze",
				defaultValue: true,
				optional:     true,
			},
			{
				name:      "repeatWeekdays",
				validType: Arr,
				optional:  true,
			},
		},
		appIntent: appIntent{
			name:                "Clock",
			bundleIdentifier:    "com.apple.clock",
			appIntentIdentifier: "CreateAlarmIntent",
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			if len(args) < 4 {
				return
			}

			var repeatDays = getArgValue(args[3])
			for _, day := range repeatDays.([]interface{}) {
				if !slices.Contains(weekdays, strings.ToLower(day.(string))) {
					parserError(fmt.Sprintf("Invalid repeat weekday for alarm '%s'", day))
				}
			}
		},
		addParams: func(args []actionArgument) (params map[string]any) {
			if len(args) < 4 {
				return
			}

			var repeatDays = getArgValue(args[3])
			var repeats []map[string]any
			for _, day := range repeatDays.([]interface{}) {
				var dayStr = day.(string)
				var dayLower = strings.ToLower(dayStr)
				var dayCap = capitalize(dayStr)

				repeats = append(repeats, map[string]any{
					"identifier": dayLower,
					"value":      dayLower,
					"title": map[string]string{
						"key": dayCap,
					},
					"subtitle": map[string]string{
						"key": dayCap,
					},
				})
			}

			return map[string]any{
				"repeats": repeats,
			}
		},
	},
	"deleteAlarm": {
		appIdentifier: "com.apple.clock",
		identifier:    "DeleteAlarmIntent",
		appIntent: appIntent{
			name:                "Clock",
			bundleIdentifier:    "com.apple.clock",
			appIntentIdentifier: "DeleteAlarmIntent",
		},
		parameters: []parameterDefinition{
			{
				name:      "alarm",
				validType: Variable,
				key:       "entities",
			},
		},
	},
	"turnOnAlarm": {
		defaultAction: true,
		appIdentifier: "com.apple.mobiletimer-framework",
		identifier:    "MobileTimerIntents.MTToggleAlarmIntent",
		appIntent:     toggleAlarmIntent,
		parameters: []parameterDefinition{
			{
				name:      "alarm",
				validType: Variable,
				key:       "alarm",
			},
			{
				name:         "showWhenRun",
				validType:    Bool,
				key:          "ShowWhenRun",
				defaultValue: true,
				optional:     true,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"state": 1,
			}
		},
	},
	"turnOffAlarm": {
		appIdentifier: "com.apple.mobiletimer-framework",
		identifier:    "MobileTimerIntents.MTToggleAlarmIntent",
		appIntent:     toggleAlarmIntent,
		parameters: []parameterDefinition{
			{
				name:      "alarm",
				validType: Variable,
				key:       "alarm",
			},
			{
				name:         "showWhenRun",
				validType:    Bool,
				key:          "ShowWhenRun",
				defaultValue: true,
				optional:     true,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"state": 0,
			}
		},
	},
	"toggleAlarm": {
		parameters: []parameterDefinition{
			{
				name:      "alarm",
				validType: Variable,
				key:       "alarm",
			},
			{
				name:         "showWhenRun",
				validType:    Bool,
				key:          "ShowWhenRun",
				defaultValue: true,
				optional:     true,
			},
		},
		appIntent: toggleAlarmIntent,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"operation": "Toggle",
			}
		},
	},
	"emailAddress": {
		identifier: "email",
		parameters: []parameterDefinition{
			{
				name:      "email",
				validType: String,
				infinite:  true,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			if len(args) > 1 && args[0].valueType == Variable {
				parserError("Shortcuts only allows one variable for an email address.")
			}
		},
		make: func(args []actionArgument) map[string]any {
			if args[0].valueType == Variable {
				return map[string]any{
					"WFEmailAddress": argumentValue(args, 0),
				}
			}

			return map[string]any{
				"WFEmailAddress": contactValue(emailAddress, args),
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompContactValue(action, "WFEmailAddress", emailAddress)
		},
	},
	"phoneNumber": {
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: String,
				infinite:  true,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			if len(args) > 1 && args[0].valueType == Variable {
				parserError("Shortcuts only allows one variable for a phone number.")
			}
		},
		make: func(args []actionArgument) map[string]any {
			if args[0].valueType == Variable {
				return map[string]any{
					"WFPhoneNumber": argumentValue(args, 0),
				}
			}

			return map[string]any{
				"WFPhoneNumber": contactValue(phoneNumber, args),
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompContactValue(action, "WFPhoneNumber", phoneNumber)
		},
	},
	"newContact": {
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
		addParams: func(args []actionArgument) (params map[string]any) {
			params = make(map[string]any)
			if len(args) >= 3 {
				if args[2].valueType == Variable {
					params["WFContactPhoneNumbers"] = argumentValue(args, 2)
				} else {
					params["WFContactPhoneNumbers"] = contactValue(phoneNumber, []actionArgument{args[2]})
				}
			}

			if len(args) >= 4 {
				if args[3].valueType == Variable {
					params["WFContactEmails"] = argumentValue(args, 3)
				} else {
					params["WFContactEmails"] = contactValue(emailAddress, []actionArgument{args[3]})
				}
			}

			return
		},
	},
	"removeContactDetail": {
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
				enum:      "contactDetails",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"Mode": "Remove",
			}
		},
	},
	"labelFile": {
		identifier: "file.label",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "color",
				validType: String,
				optional:  false,
				enum:      "fileLabel",
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{}
			}

			var color = strings.ToLower(getArgValue(args[1]).(string))

			return map[string]any{
				"WFLabelColorNumber": fileLabelsMap[color],
			}
		},
	},
	"filterFiles": {
		identifier: "filter.files",
		parameters: []parameterDefinition{
			{
				name:      "files",
				validType: Variable,
				key:       "WFContentItemInputParameter",
			},
			{
				name:      "limit",
				validType: Integer,
				key:       "WFContentItemLimitNumber",
				optional:  true,
			},
			{
				name:      "sortBy",
				validType: String,
				key:       "WFContentItemSortProperty",
				enum:      "filesSortBy",
				optional:  true,
			},
			{
				name:      "orderBy",
				validType: String,
				key:       "WFContentItemSortOrder",
				enum:      "fileOrderings",
				optional:  true,
			},
		},
		addParams: func(args []actionArgument) (params map[string]any) {
			if len(args) == 0 {
				return map[string]any{}
			}

			if len(args) != 1 {
				return map[string]any{
					"WFContentItemLimitEnabled": true,
				}
			}

			return
		},
	},
	"getPDFText": {
		identifier: "gettextfrompdf",
		parameters: []parameterDefinition{
			{
				name:      "pdfFile",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "richText",
				validType:    Bool,
				defaultValue: false,
				optional:     true,
			},
			{
				name:         "combinePages",
				validType:    Bool,
				key:          "WFCombinePages",
				defaultValue: true,
				optional:     true,
			},
			{
				name:      "headerText",
				validType: String,
				key:       "WFGetTextFromPDFPageHeader",
				optional:  true,
			},
			{
				name:      "footerText",
				validType: String,
				key:       "WFGetTextFromPDFPageFooter",
				optional:  true,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{}
			}

			if len(args) > 1 {
				var richText = getArgValue(args[1]).(bool)
				if richText {
					return map[string]any{
						"WFGetTextFromPDFTextType": "Rich Text",
					}
				}
			}

			return map[string]any{
				"WFGetTextFromPDFTextType": "Text",
			}
		},
	},
	"getFolderContents": {
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
				defaultValue: false,
				optional:     true,
			},
		},
	},
	"containsText": {
		identifier: "text.match",
		parameters: []parameterDefinition{
			{
				name:      "subject",
				validType: String,
				key:       "text",
			},
			{
				name:      "text",
				key:       "WFMatchTextPattern",
				validType: String,
			},
			{
				name:         "caseSensitive",
				validType:    Bool,
				key:          "WFMatchTextCaseSensitive",
				defaultValue: true,
				optional:     true,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			var textArg = args[1]
			if textArg.valueType == Variable {
				args[1] = actionArgument{
					valueType: String,
					value:     fmt.Sprintf("^{%s}", textArg.value),
				}
			} else {
				args[1].value = fmt.Sprintf("^%s", textArg.value)
			}
		},
	},
	"getFileFromFolder": {
		identifier: "documentpicker.open",
		parameters: []parameterDefinition{
			{
				name:      "folder",
				validType: String,
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
		addParams: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{}
			}
			var folderPath = getArgValue(args[0])
			var pathParts = strings.Split(folderPath.(string), "/")
			var fileLocationType = pathParts[0]

			var filename = end(pathParts)
			slices.Delete(pathParts, 0, len(pathParts)-1)
			folderPath = strings.Trim(strings.Join(pathParts, "/"), "/")
			var fileLocation = map[string]any{
				"relativeSubpath": folderPath,
			}

			if fileLocationType == "~" {
				fileLocation["WFFileLocationType"] = "Home"
			}

			return map[string]any{
				"WFFile": map[string]any{
					"fileLocation": fileLocation,
					"filename":     filename,
					"displayName":  filename,
				},
			}
		},
	},
	"splitText": {
		identifier: "text.split",
		parameters: []parameterDefinition{
			{
				name:      "text",
				key:       "text",
				validType: String,
			},
			{
				name:         "separator",
				validType:    String,
				defaultValue: "\n",
			},
		},
		addParams:  textParts,
		decomp:     decompTextParts,
		outputType: Arr,
	},
	"joinText": {
		identifier: "text.combine",
		parameters: []parameterDefinition{
			{
				name:      "text",
				key:       "text",
				validType: Variable,
			},
			{
				name:         "glue",
				validType:    String,
				defaultValue: "\n",
			},
		},
		addParams:  textParts,
		decomp:     decompTextParts,
		outputType: String,
	},
	"makeSizedDiskImage": {
		identifier:    "makediskimage",
		defaultAction: true,
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
				key:          "EncryptImage",
				validType:    Bool,
				defaultValue: false,
				optional:     true,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			var size = strings.Split(getArgValue(args[2]).(string), " ")
			var storageUnitArg = actionArgument{
				valueType: String,
				value:     size[1],
			}
			checkEnum(&parameterDefinition{
				name: "disk size",
				enum: "storageUnit",
			}, &storageUnitArg)
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			var imageSize = ImageSize{
				Value: SizeValue{
					Unit:      "GB",
					Magnitude: "1",
				},
			}
			mapToStruct(action.WFWorkflowActionParameters["ImageSize"], &imageSize)
			var size = fmt.Sprintf("\"%s %s\"", imageSize.Value.Magnitude, imageSize.Value.Unit)

			return []string{
				decompValue(action.WFWorkflowActionParameters["VolumeName"]),
				decompValue(action.WFWorkflowActionParameters["WFInput"]),
				size,
				decompValue(action.WFWorkflowActionParameters["EncryptImage"]),
			}
		},
		addParams: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{}
			}

			var size = strings.Split(getArgValue(args[2]).(string), " ")

			return map[string]any{
				"SizeToFit": false,
				"ImageSize": map[string]any{
					"Value": map[string]any{
						"Unit":      size[0],
						"Magnitude": size[1],
					},
					"WFSerializationType": "WFQuantityFieldValue",
				},
			}
		},
		mac:        true,
		minVersion: 15,
	},
	"seek": {
		parameters: []parameterDefinition{
			{
				name:      "magnitude",
				validType: Integer,
			},
			{
				name:      "duration",
				validType: String,
				enum:      "timerDuration",
			},
			{
				name:         "behavior",
				key:          "WFSeekBehavior",
				validType:    String,
				defaultValue: "To Time",
				enum:         "seekBehavior",
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{}
			}

			return map[string]any{
				"WFTimeInterval": map[string]any{
					"Value": map[string]any{
						"Magnitude": argumentValue(args, 0),
						"Unit":      argumentValue(args, 1),
					},
					"WFSerializationType": "WFQuantityFieldValue",
				},
			}
		},
	},
	"url": {
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				infinite:  true,
			},
		},
		make: func(args []actionArgument) map[string]any {
			var urlItems []any
			for _, item := range args {
				urlItems = append(urlItems, paramValue(item, String))
			}

			return map[string]any{
				"Show-WFURLActionURL": true,
				"WFURLActionURL":      urlItems,
			}
		},
		decomp: decompInfiniteURLAction,
	},
	"addToReadingList": {
		identifier: "readinglist",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				infinite:  true,
			},
		},
		make: func(args []actionArgument) map[string]any {
			var urlItems []any
			for _, item := range args {
				urlItems = append(urlItems, paramValue(item, String))
			}

			return map[string]any{
				"Show-WFURLActionURL": true,
				"WFURL":               urlItems,
			}
		},
		decomp: decompInfiniteURLAction,
	},
	"prompt": {
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
				enum:         "inputType",
				optional:     true,
				defaultValue: "Text",
			},
			{
				name:      "defaultValue",
				validType: String,
				optional:  true,
			},
			{
				name:         "multiline",
				validType:    String,
				key:          "WFAllowsMultilineText",
				optional:     true,
				defaultValue: true,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			if len(args) < 3 {
				return map[string]any{}
			}
			var defaultAnswer = map[string]any{
				"WFAskActionDefaultAnswer": argumentValue(args, 2),
			}
			if getArgValue(args[1]) == "Number" {
				defaultAnswer["WFAskActionDefaultAnswerNumber"] = paramValue(args[2], Integer)
			}

			return defaultAnswer
		},
	},
	"openApp": {
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "appID",
				validType: String,
				key:       "WFAppIdentifier",
			},
		},
		check: func(args []actionArgument, definition *actionDefinition) {
			replaceAppIDs(args, definition)
		},
		addParams: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{}
			}

			if args[0].valueType == Variable {
				return map[string]any{
					"WFSelectedApp": argumentValue(args, 0),
				}
			}

			return map[string]any{
				"WFSelectedApp": map[string]any{
					"BundleIdentifier": argumentValue(args, 0),
				},
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompAppAction("WFAppIdentifier", action)
		},
	},
	"hideApp": {
		identifier:    "hide.app",
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "appID",
				validType: String,
			},
		},
		check: func(args []actionArgument, definition *actionDefinition) {
			replaceAppIDs(args, definition)
		},
		make: func(args []actionArgument) map[string]any {
			if args[0].valueType == Variable {
				return map[string]any{
					"WFApp": argumentValue(args, 0),
				}
			}

			return map[string]any{
				"WFApp": map[string]any{
					"BundleIdentifier": argumentValue(args, 0),
				},
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompAppAction("WFApp", action)
		},
	},
	"hideAllApps": {
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
		make: func(args []actionArgument) map[string]any {
			if args[0].valueType != Variable {
				return map[string]any{
					"WFApp": argumentValue(args, 0),
				}
			}

			return map[string]any{
				"WFApp": map[string]any{
					"BundleIdentifier": argumentValue(args, 0),
				},
			}
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFHideAppMode": "All Apps",
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompAppAction("WFAppsExcept", action)
		},
	},
	"quitApp": {
		identifier:    "quit.app",
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "appID",
				validType: String,
			},
		},
		check: func(args []actionArgument, definition *actionDefinition) {
			replaceAppIDs(args, definition)
		},
		make: func(args []actionArgument) map[string]any {
			if args[0].valueType == Variable {
				return map[string]any{
					"WFApp": argumentValue(args, 0),
				}
			}

			return map[string]any{
				"WFApp": map[string]any{
					"BundleIdentifier": argumentValue(args, 0),
				},
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompAppAction("WFApp", action)
		},
	},
	"quitAllApps": {
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
		make: func(args []actionArgument) (params map[string]any) {
			params = make(map[string]any)
			if args[0].valueType != Variable {
				params["WFAppsExcept"] = apps(args)
			} else {
				params["WFAppsExcept"] = argumentValue(args, 0)
			}

			return
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFQuitAppMode": "All Apps",
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompAppAction("WFAppsExcept", action)
		},
	},
	"killApp": {
		identifier: "quit.app",
		parameters: []parameterDefinition{
			{
				name:      "appID",
				validType: String,
			},
		},
		check: func(args []actionArgument, definition *actionDefinition) {
			replaceAppIDs(args, definition)
		},
		make: func(args []actionArgument) (params map[string]any) {
			params = make(map[string]any)

			params["WFAskToSaveChanges"] = false

			if args[0].valueType == Variable {
				params["WFApp"] = argumentValue(args, 0)
				return
			}

			params["WFApp"] = map[string]any{
				"BundleIdentifier": argumentValue(args, 0),
			}

			return
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompAppAction("WFApp", action)
		},
	},
	"killAllApps": {
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
		make: func(args []actionArgument) (params map[string]any) {
			params = map[string]any{
				"WFQuitAppMode":      "All Apps",
				"WFAskToSaveChanges": false,
			}

			if args[0].valueType != Variable {
				params["WFAppsExcept"] = apps(args)
			} else {
				params["WFAppsExcept"] = argumentValue(args, 0)
			}

			return
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompAppAction("WFAppsExcept", action)
		},
	},
	"splitApps": {
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
				key:          "WFAppRatio",
				validType:    String,
				optional:     true,
				enum:         "appSplitRatio",
				defaultValue: "half",
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
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
		addParams: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{}
			}

			var params = make(map[string]any)
			if args[0].valueType == Variable {
				params["WFPrimaryAppIdentifier"] = argumentValue(args, 0)
			} else {
				params["WFPrimaryAppIdentifier"] = map[string]any{
					"BundleIdentifier": argumentValue(args, 0),
				}
			}

			if args[0].valueType == Variable {
				params["WFSecondaryAppIdentifier"] = argumentValue(args, 0)
			} else {
				params["WFSecondaryAppIdentifier"] = map[string]any{
					"BundleIdentifier": argumentValue(args, 0),
				}
			}

			return params
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			var splitRatio = decompValue(action.WFWorkflowActionParameters["WFAppRatio"])

			var ratio = "half"
			switch splitRatio {
			case "½ + ½":
				ratio = "half"
			case "⅓ + ⅔":
				ratio = "thirdByTwo"
			}

			arguments = append(arguments, fmt.Sprintf("\"%s\"", ratio))
			arguments = append(arguments, decompAppAction("WFPrimaryAppIdentifier", action)...)
			arguments = append(arguments, decompAppAction("WFSecondaryAppIdentifier", action)...)

			return
		},
	},
	"openShortcut": {
		appIdentifier: "com.apple.shortcuts",
		identifier:    "OpenWorkflowAction",
		parameters: []parameterDefinition{
			{
				name:      "shortcutName",
				validType: String,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{}
			}

			return map[string]any{
				"target": map[string]any{
					"title": argumentValue(args, 0),
				},
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			if action.WFWorkflowActionParameters["target"] != nil {
				var workflow = action.WFWorkflowActionParameters["target"].(map[string]interface{})
				var title = workflow["title"].(map[string]interface{})
				arguments = append(arguments, decompValue(title["key"]))
			}
			return
		},
	},
	"runSelf": {
		identifier: "runworkflow",
		parameters: []parameterDefinition{
			{
				name:      "output",
				key:       "WFInput",
				validType: Variable,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{
					"isSelf": true,
				}
			}

			return map[string]any{
				"WFWorkflow": map[string]any{
					"workflowIdentifier": uuid.New().String(),
					"isSelf":             true,
					"workflowName":       workflowName,
				},
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			if action.WFWorkflowActionParameters["WFInput"] != nil {
				arguments = append(arguments, decompValue(action.WFWorkflowActionParameters["WFInput"]))
			}
			return
		},
	},
	"run": {
		identifier:    "runworkflow",
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "shortcutName",
				validType: String,
			},
			{
				name:      "output",
				validType: Variable,
				key:       "WFInput",
				optional:  true,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{
					"isSelf": false,
				}
			}

			return map[string]any{
				"WFWorkflow": map[string]any{
					"workflowIdentifier": uuid.New().String(),
					"isSelf":             false,
					"workflowName":       argumentValue(args, 0),
				},
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			var workflow = action.WFWorkflowActionParameters["WFWorkflow"].(map[string]any)
			if workflow["isSelf"] != nil && !workflow["isSelf"].(bool) {
				arguments = append(arguments, decompValue(workflow["workflowName"]))
			}
			if action.WFWorkflowActionParameters["WFInput"] != nil {
				arguments = append(arguments, decompValue(action.WFWorkflowActionParameters["WFInput"]))
			}

			return
		},
	},
	"list": {
		parameters: []parameterDefinition{
			{
				name:      "listItem",
				validType: String,
				infinite:  true,
			},
		},
		make: func(args []actionArgument) map[string]any {
			var listItems []map[string]any
			for _, item := range args {
				listItems = append(listItems, map[string]any{
					"WFItemType": 0,
					"WFValue":    paramValue(item, String),
				})
			}

			return map[string]any{
				"WFItems": listItems,
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			var listItems = action.WFWorkflowActionParameters["WFItems"].([]interface{})
			for _, item := range listItems {
				var itemValue = item
				if reflect.TypeOf(item).String() != "string" {
					itemValue = item.(map[string]interface{})["WFValue"]
				}
				arguments = append(arguments, decompValue(itemValue))
			}
			return
		},
	},
	"openCustomXCallbackURL": {
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
		addParams: func(args []actionArgument) (xCallbackParams map[string]any) {
			if len(args) == 0 {
				return
			}
			xCallbackParams = make(map[string]any)
			if args[1].value.(string) != "" || args[2].value.(string) != "" || args[3].value.(string) != "" {
				xCallbackParams["WFXCallbackCustomCallbackEnabled"] = true
			}
			if args[4].value.(string) != "" {
				xCallbackParams["WFXCallbackCustomSuccessURLEnabled"] = true
			}

			return
		},
	},
	"createShortcutLink": {
		appIdentifier: "com.apple.shortcuts",
		identifier:    "CreateShortcutiCloudLinkAction",
		appIntent:     createShortcutiCloudLink,
		parameters: []parameterDefinition{
			{
				name:      "shortcut",
				key:       "shortcut",
				validType: Variable,
			},
		},
	},
	"formRequest": {
		identifier: "downloadurl",
		parameters: httpParams,
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompHTTPAction("WFFormValues", action)
		},
		addParams: func(args []actionArgument) map[string]any {
			return httpRequest("Form", "WFFormValues", args)
		},
	},
	"jsonRequest": {
		identifier: "downloadurl",
		parameters: httpParams,
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompHTTPAction("WFJSONValues", action)
		},
		addParams: func(args []actionArgument) map[string]any {
			return httpRequest("JSON", "WFJSONValues", args)
		},
	},
	"fileRequest": {
		identifier: "downloadurl",
		parameters: httpParams,
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompHTTPAction("WFRequestVariable", action)
		},
		addParams: func(args []actionArgument) map[string]any {
			return httpRequest("File", "WFRequestVariable", args)
		},
	},
	"getWindows": {
		identifier: "filter.windows",
		parameters: []parameterDefinition{
			{
				name:      "sortBy",
				validType: String,
				key:       "WFContentItemSortProperty",
				enum:      "windowSorting",
				optional:  true,
			},
			{
				name:      "orderBy",
				validType: String,
				key:       "WFContentItemSortOrder",
				enum:      "sortOrder",
				optional:  true,
			},
			{
				name:      "limit",
				validType: Integer,
				key:       "WFContentItemLimitNumber",
				optional:  true,
			},
		},
		addParams: func(args []actionArgument) (params map[string]any) {
			if len(args) == 0 {
				return map[string]any{}
			}

			if args[2].value != nil {
				params = map[string]any{
					"WFContentItemLimitEnabled": true,
				}
			}
			return
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			if args[1].value != nil {
				var alphabetic = []string{"Title", "App Name", "Name", "Random"}
				var numeric = []string{"Width", "Height", "X Position", "Y Position", "Window Index"}
				var sortBy = getArgValue(args[0]).(string)
				var orderBy = getArgValue(args[1]).(string)
				if sortBy != "Random" {
					if slices.Contains(alphabetic, sortBy) {
						switch orderBy {
						case "asc":
							args[1].value = "A to Z"
						case "desc":
							args[1].value = "Z to A"
						}
					} else if slices.Contains(numeric, sortBy) {
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
	},
	"convertMeasurement": {
		identifier: "measurement.convert",
		parameters: []parameterDefinition{
			{
				name:      "measurement",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "unitType",
				validType: String,
				key:       "WFMeasurementUnitType",
				enum:      "measurementUnitType",
			},
			{
				name:      "unit",
				validType: String,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			var value = getArgValue(args[1])
			if reflect.TypeOf(value).Kind() != reflect.String {
				return
			}

			var unitType = value.(string)
			checkEnum(&parameterDefinition{
				name: "measurement unit",
				enum: unitType,
			}, &args[2])
		},
		addParams: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{
					"isSelf": false,
				}
			}

			return map[string]any{
				"WFMeasurementUnit": map[string]any{
					"WFNSUnitType":   argumentValue(args, 1),
					"WFNSUnitSymbol": argumentValue(args, 2),
				},
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			arguments = append(arguments,
				decompValue(action.WFWorkflowActionParameters["WFInput"]),
				decompValue(action.WFWorkflowActionParameters["WFMeasurementUnitType"]),
			)

			if action.WFWorkflowActionParameters["WFMeasurementUnit"] != nil {
				var measurementUnit WFMeasurementUnit
				mapToStruct(action.WFWorkflowActionParameters["WFMeasurementUnit"].(map[string]interface{}), &measurementUnit)

				arguments = append(arguments, decompValue(measurementUnit.WFNSUnitSymbol))
			}

			return
		},
	},
	"measurement": {
		identifier: "measurement.create",
		parameters: []parameterDefinition{
			{
				name:      "magnitude",
				validType: String,
			},
			{
				name:      "unitType",
				validType: String,
				key:       "WFMeasurementUnitType",
				enum:      "measurementUnitType",
			},
			{
				name:      "unit",
				validType: String,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			var value = getArgValue(args[1])
			if reflect.TypeOf(value).String() != "string" {
				return
			}

			var unitType = value.(string)
			checkEnum(&parameterDefinition{
				name: "unit",
				enum: unitType,
			}, &args[2])
		},
		addParams: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{}
			}

			return map[string]any{
				"WFMeasurementUnit": map[string]any{
					"Value": map[string]any{
						"Magnitude": argumentValue(args, 0),
						"Unit":      argumentValue(args, 2),
					},
					"WFSerializationType": "WFQuantityFieldValue",
				},
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			if action.WFWorkflowActionParameters["WFMeasurementUnit"] != nil {
				var measurementUnit WFMeasurementUnit
				mapToStruct(action.WFWorkflowActionParameters["WFMeasurementUnit"].(map[string]interface{}), &measurementUnit)

				arguments = append(arguments,
					decompValue(measurementUnit.Value.Magnitude),
					decompValue(action.WFWorkflowActionParameters["WFMeasurementUnitType"]),
					decompValue(measurementUnit.Value.Unit),
				)
			}

			return
		},
	},
	"makeVCard": {
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
				name:      "base64Image",
				validType: String,
				optional:  true,
			},
		},
		make: func(args []actionArgument) map[string]any {
			var title = args[0].value.(string)
			var subtitle = args[1].value.(string)
			wrapVariableReference(&title)
			wrapVariableReference(&subtitle)

			var vcard strings.Builder
			vcard.WriteString(fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nN;CHARSET=utf-8:%s\nORG:%s\n", title, subtitle))

			if len(args) > 2 {
				var photo string
				var image = getArgValue(args[2])
				if reflect.TypeOf(image).Kind() != reflect.String && args[2].valueType == Variable {
					photo = fmt.Sprintf("{%s}", makeVariableReferenceString(args[2].value.(varValue)))
				} else {
					photo = getArgValue(args[2]).(string)
				}

				if photo != "" {
					vcard.WriteString(fmt.Sprintf("PHOTO;ENCODING=b:%s\n", photo))
				}
			}

			vcard.WriteString("END:VCARD")
			args[0] = actionArgument{
				valueType: String,
				value:     vcard.String(),
			}

			return map[string]any{
				"WFTextActionText": argumentValue(args, 0),
			}
		},
	},
	"base64File": {
		identifier: "gettext",
		parameters: []parameterDefinition{
			{
				name:      "filePath",
				validType: String,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			var file = getArgValue(args[0])
			if args[0].valueType == Variable && reflect.TypeOf(file).Kind() != reflect.String {
				parserError("File path must be a string literal")
			}
			if _, err := os.Stat(file.(string)); os.IsNotExist(err) {
				parserError(fmt.Sprintf("File '%s' does not exist!", file))
			}
		},
		make: func(args []actionArgument) map[string]any {
			var file = getArgValue(args[0]).(string)
			var bytes, readErr = os.ReadFile(file)
			handle(readErr)
			var encodedFile = base64.StdEncoding.EncodeToString(bytes)

			return map[string]any{
				"WFTextActionText": encodedFile,
			}
		},
	},
	"updateContact": {
		identifier:    "setters.contacts",
		defaultAction: true,
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
				enum:      "contactDetail",
			},
			{
				name:      "value",
				validType: String,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			var contactDetail = args[1].value.(string)
			var contactDetailKey = strings.ReplaceAll(contactDetail, " ", "")
			currentAction.parameters[2].key = "WFContactContentItem" + contactDetailKey
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"Mode": "Set",
			}
		},
	},
}

var useAttachmentAsVariableValueRegex = regexp.MustCompile(`^\$\{[a-zA-Z0-9]+}$`)

func defineRawAction() {
	actions["rawAction"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "identifier",
				validType: String,
			},
			{
				name:      "parameters",
				optional:  true,
				validType: Dict,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			actions["rawAction"].overrideIdentifier = getArgValue(args[0]).(string)
		},
		make: func(args []actionArgument) map[string]any {
			if len(args) == 1 {
				return map[string]any{}
			}

			var params = getArgValue(args[1]).(map[string]interface{})
			handleRawParams(params)

			return params
		},
	}
}

func handleRawParams(params map[string]interface{}) {
	for key, value := range params {
		if reflect.TypeOf(value).Kind() != reflect.String || !strings.ContainsAny(value.(string), "{}") {
			continue
		}
		if useAttachmentAsVariableValueRegex.MatchString(value.(string)) {
			params[key] = variableValue(varValue{
				value: strings.Trim(value.(string), "${}"),
			})
			continue
		}
		params[key] = attachmentValues(value.(string))
	}
}

var defaultActionIncludes = []string{
	"basic",
	"calendar",
	"contacts",
	"documents",
	"device",
	"location",
	"math",
	"media",
	"scripting",
	"sharing",
	"shortcuts",
	"translation",
	"web",
}

func loadStandardActions() {
	includeStandardActions()
	handleIncludes()
	handleActionDefinitions()
}

var includedStandardActions bool

func includeStandardActions() {
	if includedStandardActions {
		return
	}
	var actionIncludes []string
	for _, actionInclude := range defaultActionIncludes {
		actionIncludes = append(actionIncludes, fmt.Sprintf("#include 'actions/%s'\n", actionInclude))
	}
	lines = append(actionIncludes, lines...)
	resetParse()
	includedStandardActions = true
}

type contentKit string

var emailAddress contentKit = "emailaddress"
var phoneNumber contentKit = "phonenumber"

func contactValue(contentKit contentKit, args []actionArgument) map[string]any {
	var contactValues []map[string]any
	var entryType int
	switch contentKit {
	case emailAddress:
		entryType = 2
	case phoneNumber:
		entryType = 1
	}
	for _, item := range args {
		contactValues = append(contactValues, map[string]any{
			"EntryType": entryType,
			"SerializedEntry": map[string]any{
				"link.contentkit." + string(contentKit): item.value,
			},
		})
	}

	return map[string]any{
		"Value": map[string]any{
			"WFContactFieldValues": contactValues,
		},
		"WFSerializationType": "WFContactFieldValue",
	}
}

func decompContactValue(action *ShortcutAction, key string, contentKit contentKit) (arguments []string) {
	if action.WFWorkflowActionParameters[key] == nil {
		return
	}

	var wfValue WFValue
	mapToStruct(action.WFWorkflowActionParameters[key], &wfValue)

	if wfValue.WFSerializationType == "WFTextTokenString" {
		return []string{decompValue(action.WFWorkflowActionParameters[key])}
	}

	var value = wfValue.Value.(map[string]interface{})
	var contactFieldValues []WFContactFieldValue
	mapToStruct(value["WFContactFieldValues"].([]interface{}), &contactFieldValues)

	for _, contactFieldValue := range contactFieldValues {
		var serializedKey = fmt.Sprintf("link.contentkit.%s", contentKit)
		var email = contactFieldValue.SerializedEntry[serializedKey]
		arguments = append(arguments, fmt.Sprintf("\"%s\"", email.(string)))
	}

	return
}

func adjustDate(operation string, unit string, args []actionArgument) (adjustDateParams map[string]any) {
	if len(args) == 0 {
		return map[string]any{}
	}

	adjustDateParams = map[string]any{
		"WFAdjustOperation": operation,
	}
	if unit == "" {
		return adjustDateParams
	}

	adjustDateParams["WFDuration"] = magnitudeValue(unit, args, 1)

	return
}

func magnitudeValue(unit any, args []actionArgument, index int) map[string]any {
	var magnitudeValue = argumentValue(args, index)
	if reflect.TypeOf(magnitudeValue).String() == "[]map[string]any" {
		var value = magnitudeValue.([]map[string]any)
		magnitudeValue = value[0]
	}

	return map[string]any{
		"Value": map[string]any{
			"Unit":      unit,
			"Magnitude": magnitudeValue,
		},
		"WFSerializationType": "WFQuantityFieldValue",
	}
}

func textParts(args []actionArgument) map[string]any {
	if len(args) == 0 {
		return map[string]any{}
	}

	var data = map[string]any{
		"Show-text": true,
	}

	var separator = getArgValue(args[1])
	switch {
	case separator == " ":
		data["WFTextSeparator"] = "Spaces"
	case separator == "\n":
		data["WFTextSeparator"] = "New Lines"
	case separator == "" && currentAction.identifier == "splitText":
		data["WFTextSeparator"] = "Every Character"
	default:
		data["WFTextSeparator"] = "Custom"
		data["WFTextCustomSeparator"] = argumentValue(args, 1)
	}

	return data
}

func decompReferenceValue(paramValue any) string {
	if isReferenceValue(paramValue) {
		return decompValue(paramValue)
	}

	return paramValue.(string)
}

func decompTextParts(action *ShortcutAction) (arguments []string) {
	arguments = append(arguments, decompValue(action.WFWorkflowActionParameters["text"]))

	var glue string
	if action.WFWorkflowActionParameters["WFTextSeparator"] != nil {
		glue = decompReferenceValue(action.WFWorkflowActionParameters["WFTextSeparator"])
	}
	if action.WFWorkflowActionParameters["WFTextCustomSeparator"] != nil {
		glue = decompReferenceValue(action.WFWorkflowActionParameters["WFTextCustomSeparator"])
	}

	if glue != "" {
		arguments = append(arguments, fmt.Sprintf("\"%s\"", glueToChar(glue)))
	}

	return
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

func apps(args []actionArgument) (apps []map[string]any) {
	for _, arg := range args {
		if arg.valueType != Variable {
			apps = append(apps, map[string]any{
				"BundleIdentifier": arg.value,
				"TeamIdentifier":   "0000000000",
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

func replaceAppIDs(args []actionArgument, _ *actionDefinition) {
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

func httpRequest(bodyType string, valuesKey string, args []actionArgument) map[string]any {
	var params = map[string]any{"WFHTTPBodyType": bodyType}

	if len(args) > 0 {
		params[valuesKey] = argumentValue(args, 2)
	}

	return params
}

func decompAppAction(key string, action *ShortcutAction) (arguments []string) {
	if action.WFWorkflowActionParameters[key] != nil {
		switch reflect.TypeOf(action.WFWorkflowActionParameters[key]).Kind() {
		case reflect.String:
			return append(arguments, decompValue(action.WFWorkflowActionParameters[key]))
		case reflect.Map:
			for key, bundle := range action.WFWorkflowActionParameters[key].(map[string]interface{}) {
				if key == "BundleIdentifier" {
					arguments = append(arguments, fmt.Sprintf("\"%s\"", bundle))
				}
			}
		case reflect.Array:
			for _, app := range action.WFWorkflowActionParameters[key].([]interface{}) {
				var bundleIdentifer = app.(map[string]interface{})["BundleIdentifier"]
				arguments = append(arguments, fmt.Sprintf("\"%s\"", bundleIdentifer))
			}
		default:
			decompError("Unknown app value type", action)
		}
	}

	return
}

func decompInfiniteURLAction(action *ShortcutAction) (arguments []string) {
	var urlValueType = reflect.TypeOf(action.WFWorkflowActionParameters["WFURLActionURL"]).Kind()
	if urlValueType == reflect.Map || urlValueType == reflect.String {
		return append(arguments, decompValue(action.WFWorkflowActionParameters["WFURLActionURL"]))
	}

	var urls = action.WFWorkflowActionParameters["WFURLActionURL"].([]interface{})
	for _, url := range urls {
		arguments = append(arguments, decompValue(url))
	}

	return
}

func decompHTTPAction(key string, action *ShortcutAction) (arguments []string) {
	if action.WFWorkflowActionParameters["WFURL"] != nil {
		arguments = append(arguments, decompValue(action.WFWorkflowActionParameters["WFURL"]))
	}
	if action.WFWorkflowActionParameters["WFHTTPMethod"] != nil {
		arguments = append(arguments, decompValue(action.WFWorkflowActionParameters["WFHTTPMethod"]))
	}
	if action.WFWorkflowActionParameters[key] != nil {
		arguments = append(arguments, decompValue(action.WFWorkflowActionParameters[key]))
	}
	if action.WFWorkflowActionParameters["WFHTTPHeaders"] != nil {
		arguments = append(arguments, decompValue(action.WFWorkflowActionParameters["WFHTTPHeaders"]))
	}

	return
}

// toggleSetActions are actions which all are state based and so can either be toggled or set in the same format.
var toggleSetActions = map[string]actionDefinition{
	"BackgroundSounds": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleBackgroundSoundsIntent",
	},
	"MediaBackgroundSounds": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleBackgroundSoundsIntent",
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"setting": "whenMediaIsPlaying",
			}
		},
	},
	"AutoAnswerCalls": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleAutoAnswerCallsIntent",
	},
	"Appearance": {
		identifier: "appearance",
	},
	"Bluetooth": {
		identifier: "bluetooth.set",
		setKey:     "OnValue",
	},
	"Wifi": {
		identifier: "wifi.set",
		setKey:     "OnValue",
	},
	"CellularData": {
		identifier: "cellulardata.set",
		setKey:     "OnValue",
	},
	"NightShift": {
		identifier: "nightshift.set",
		setKey:     "OnValue",
	},
	"TrueTone": {
		identifier: "truetone.set",
		setKey:     "OnValue",
	},
	"AirplaneMode": {
		identifier: "airplanemode.set",
		setKey:     "OnValue",
	},
	"ClassicInvert": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleClassicInvertIntent",
	},
	"ClosedCaptionsSDH": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleCaptionsIntent",
	},
	"ColorFilters": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleColorFiltersIntent",
	},
	"Contrast": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleContrastIntent",
	},
	"LEDFlash": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleLEDFlashIntent",
	},
	"LeftRightBalance": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXSetLeftRightBalanceIntent",
		parameters: []parameterDefinition{
			{
				name:      "value",
				validType: Integer,
				key:       "value",
				optional:  true,
			},
		},
	},
	"LiveCaptions": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleLiveCaptionsIntent",
	},
	"MonoAudio": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleMonoAudioIntent",
	},
	"ReduceMotion": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleReduceMotionIntent",
	},
	"ReduceTransparency": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleTransparencyIntent",
	},
	"SmartInvert": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleSmartInvertIntent",
	},
	"SwitchControl": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleSwitchControlIntent",
	},
	"VoiceControl": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleVoiceControlIntent",
	},
	"WhitePoint": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleWhitePointIntent",
	},
	"Zoom": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleZoomIntent",
	},
	"StageManager": {
		identifier: "stagemanager.set",
		parameters: []parameterDefinition{
			{
				name:         "showDock",
				key:          "showDock",
				validType:    Bool,
				defaultValue: true,
			},
			{
				name:         "showRecentApps",
				key:          "showRecentApps",
				validType:    Bool,
				defaultValue: true,
			},
		},
	},
}

// defineToggleSetActions automates the creation of actions which simply toggle and set a state in the same format.
func defineToggleSetActions() {
	for name, def := range toggleSetActions {
		var toggleName = fmt.Sprintf("toggle%s", name)
		def.addParams = func(_ []actionArgument) map[string]any {
			return map[string]any{
				"operation": "toggle",
			}
		}
		var toggleDef = def
		def.defaultAction = false
		actions[toggleName] = &toggleDef

		if name == "Appearance" {
			continue
		}

		var setName = fmt.Sprintf("set%s", name)
		def.addParams = nil
		def.defaultAction = true
		var setKey = "state"
		if def.setKey != "" {
			setKey = def.setKey
		}
		def.parameters = append([]parameterDefinition{
			{
				name:      "status",
				validType: Bool,
				key:       setKey,
			},
		}, def.parameters...)
		var setDef = def
		actions[setName] = &setDef
	}
	toggleSetActions = nil
}
