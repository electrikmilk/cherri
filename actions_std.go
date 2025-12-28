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

	"github.com/electrikmilk/args-parser"
	"github.com/google/uuid"
)

const SetVariableIdentifier = "is.workflow.actions.setvariable"
const AppendVariableIdentifier = "is.workflow.actions.appendvariable"

var weekdays = []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

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
	"returnToHomescreen": {nonMacOnly: true},
	"createAlarm": {
		doc: selfDoc{
			title:       "Create Alarm",
			description: "Creates an alarm at a specific time with a name, snooze allowance, and applicable weekdays.",
			category:    "calendar",
			subcategory: "Alarms",
		},
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
		doc: selfDoc{
			title:       "Delete Alarm",
			description: "Deletes an alarm.",
			category:    "calendar",
			subcategory: "Alarms",
		},
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
		doc: selfDoc{
			title:       "Turn On Alarm",
			description: "Turn on an alarm.",
			category:    "calendar",
			subcategory: "Alarms",
		},
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
		defaultAction: true,
	},
	"turnOffAlarm": {
		doc: selfDoc{
			title:       "Turn Off Alarm",
			description: "Turn off an alarm.",
			category:    "calendar",
			subcategory: "Alarms",
		},
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
		doc: selfDoc{
			title:       "Toggle Alarm",
			description: "Toggle an alarm.",
			category:    "calendar",
			subcategory: "Alarms",
		},
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
		doc: selfDoc{
			title:       "Email Address",
			description: "Create an email address value.",
			category:    "contacts",
			subcategory: "Email",
		},
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
		doc: selfDoc{
			title:       "Phone Number",
			description: "Create a phone number value.",
			category:    "contacts",
			subcategory: "Phone",
		},
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
		doc: selfDoc{
			title:       "Add New Contact",
			description: "Create a new contact.",
			category:    "contacts",
		},
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
		doc: selfDoc{
			title:       "Remove Contact Detail",
			description: "Remove detail from contact.",
			category:    "contacts",
		},
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
		doc: selfDoc{
			title:       "Label File",
			description: "Label a file.",
			category:    "documents",
			subcategory: "Files & Folders",
		},
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
		doc: selfDoc{
			title:       "Filter Files",
			description: "Filter the provided files with various filters.",
			category:    "documents",
			subcategory: "Files & Folders",
		},
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
		doc: selfDoc{
			title:       "Get PDF Text",
			description: "Get text from PDF.",
			category:    "pdf",
		},
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
		doc: selfDoc{
			title:       "Get Folder Contacts",
			description: "Get contents of folder.",
			category:    "documents",
			subcategory: "Files & Folders",
		},
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
		doc: selfDoc{
			title:       "Contains Text",
			description: "Uses Match Text to check if text is within subject.",
			category:    "text",
			subcategory: "Text Editing",
		},
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
		doc: selfDoc{
			title:       "Get File From Folder",
			description: "Get a file from a folder.",
			category:    "documents",
			subcategory: "Files & Folders",
		},
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
		doc: selfDoc{
			title:       "Split Text",
			description: "Split text by a separator.",
			category:    "text",
			subcategory: "Text Editing",
		},
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
		doc: selfDoc{
			title:       "Join Text",
			description: "Join text by a combiner.",
			category:    "text",
			subcategory: "Text Editing",
		},
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
	"url": {
		doc: selfDoc{
			title:       "URL",
			description: "Create a URL value.",
			category:    "web",
			subcategory: "URLs",
		},
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
		doc: selfDoc{
			title:       "Add to Reading List",
			description: "Add a link to the reading list.",
			category:    "web",
			subcategory: "Safari",
		},
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
		doc: selfDoc{
			title:       "Ask for Input",
			description: "Ask for input with prompt, with optional inputType and defaultValue.",
			category:    "basic",
			subcategory: "Notifications",
		},
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
		doc: selfDoc{
			title:       "Open App",
			description: "Open an app.",
			category:    "scripting",
			subcategory: "Apps",
		},
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
		doc: selfDoc{
			title:       "Hide App",
			description: "Hide an app.",
			category:    "scripting",
			subcategory: "Apps",
		},
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
		doc: selfDoc{
			title:       "Hide Apps",
			description: "Hide multiple apps. Allows exception.",
			category:    "scripting",
			subcategory: "Apps",
		},
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
		doc: selfDoc{
			title:       "Qut App",
			description: "Quit an app.",
			category:    "scripting",
			subcategory: "Apps",
		},
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
		doc: selfDoc{
			title:       "Quit All Apps",
			description: "Quits all apps. Allows exceptions.",
			category:    "scripting",
			subcategory: "Apps",
		},
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
		doc: selfDoc{
			title:       "Kill App",
			description: "Kill an app.",
			warning:     "This will not ask to save changes!",
			category:    "scripting",
			subcategory: "Apps",
		},
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
		doc: selfDoc{
			title:       "Kill All Apps",
			description: "Kills all apps. Allows exceptions.",
			warning:     "This will quit all the apps running on the device without asking to save changes!",
			category:    "scripting",
			subcategory: "Apps",
		},
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
		doc: selfDoc{
			title:       "Split Apps",
			description: "Split apps across the screen.",
			category:    "scripting",
			subcategory: "Apps",
		},
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
		doc: selfDoc{
			title:       "Open Shortcut",
			description: "Open a shortcut in the Shortcuts app.",
			category:    "shortcuts",
		},
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
		doc: selfDoc{
			title:       "Run Self",
			description: "Run the current Shortcut with optional output.",
			category:    "shortcuts",
		},
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
	"list": {
		doc: selfDoc{
			title:       "List",
			description: "Create a list.",
			category:    "scripting",
			subcategory: "Lists",
		},
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
		doc: selfDoc{
			title:    "Open Custom X-Callback URL",
			category: "web",
		},
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
		doc: selfDoc{
			title:    "Create Shortcut Link",
			category: "shortcuts",
		},
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
	"getWindows": {
		doc: selfDoc{
			title:       "Get Windows",
			category:    "mac",
			subcategory: "Windows",
		},
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
		macOnly: true,
	},
	"convertMeasurement": {
		doc: selfDoc{
			title:       "Convert Measurement",
			category:    "scripting",
			subcategory: "Measurement",
		},
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
		doc: selfDoc{
			title:       "Create Measurement",
			category:    "scripting",
			subcategory: "Measurement",
		},
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
				name: fmt.Sprintf("%s unit", unitType),
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
		doc: selfDoc{
			title:    "Make VCard",
			category: "builtin",
		},
		identifier: "gettext",
		parameters: []parameterDefinition{
			{
				name:      "title",
				validType: String,
				literal:   true,
			},
			{
				name:      "subtitle",
				validType: String,
				literal:   true,
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
	"embedFile": {
		doc: selfDoc{
			title:       "Base 64 Embed File",
			description: "Embed file at path as base 64 text.",
			category:    "builtin",
		},
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
		doc: selfDoc{
			title:    "Update Contact",
			category: "contacts",
		},
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
	"setFocusMode": {
		doc: selfDoc{
			title:       "Set Focus Mode",
			description: "Set a default focus mode on or off. If setting to on, optionally set until with optional arguments for time or event.",
			category:    "settings",
			subcategory: "Notifications",
		},
		identifier: "dnd.set",
		parameters: []parameterDefinition{
			{
				name:         "focusMode",
				validType:    String,
				defaultValue: "Do Not Disturb",
				enum:         "focusModes",
				optional:     true,
			},
			{
				name:         "until",
				validType:    String,
				key:          "AssertionType",
				defaultValue: "Turned Off",
				enum:         "focusUntil",
				optional:     true,
			},
			{
				name:      "time",
				validType: String,
				key:       "Time",
				optional:  true,
			},
			{
				name:      "event",
				validType: Variable,
				key:       "Event",
				optional:  true,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			if len(args) > 0 {
				var mode = getArgValue(args[0]).(string)
				if fm, found := focusModes[mode]; found {
					return map[string]any{
						"FocusModes": fm,
					}
				}
			}

			return map[string]any{}
		},
	},
	"toggleFocusMode": {
		doc: selfDoc{
			title:       "Toggle Focus Mode",
			description: "Toggle a focus mode.",
			category:    "settings",
			subcategory: "Notifications",
		},
		identifier: "dnd.set",
		parameters: []parameterDefinition{
			{
				name:         "focusMode",
				validType:    String,
				defaultValue: "Do Not Disturb",
				enum:         "focusModes",
				optional:     true,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			var params = map[string]any{
				"Operation": "Toggle",
			}

			if len(args) > 0 {
				var mode = getArgValue(args[0]).(string)
				if fm, found := focusModes[mode]; found {
					params["FocusModes"] = fm
				}
			}

			return params
		},
	},
	"generateImage": {
		doc: selfDoc{
			title:       "Create Image using Image Playground",
			description: "Generate an Image with a prompt using the Image Playground app.",
			category:    "intelligence",
			subcategory: "Image Playground",
		},
		appIdentifier: "com.apple.GenerativePlaygroundApp",
		identifier:    "GenerateImageIntent",
		parameters: []parameterDefinition{
			{
				name:      "prompt",
				validType: String,
				key:       "prompt",
			},
			{
				name:      "image",
				validType: Variable,
				key:       "image",
				optional:  true,
			},
			{
				name:         "style",
				validType:    String,
				defaultValue: "animation",
				optional:     true,
				enum:         "imagePlaygroundStyle",
			},
			{
				name:         "saveToPlayground",
				validType:    String,
				key:          "saveToPlayground",
				enum:         "saveToPlaygroundBehavior",
				defaultValue: "always",
				optional:     true,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			if len(args) < 3 {
				return map[string]any{}
			}

			var title string
			var subtitle string
			var style = argumentValue(args, 2)
			if args[2].valueType == String {
				style = getArgValue(args[2])
				switch style {
				case "chatgpt":
					style = "z_external_provider"
					title = "ChatGPT"
					subtitle = "ChatGPT"
				case "chatgpt_oil_painting":
					style = "z_external_provider_1"
					title = "Oil Painting (ChatGPT)"
					subtitle = "Oil Painting (ChatGPT)"
				case "chatgpt_watercolor":
					style = "z_external_provider_2"
					title = "Watercolor (ChatGPT)"
					subtitle = "Watercolor (ChatGPT)"
				case "chatgpt_vector":
					style = "z_external_provider_3"
					title = "Vector (ChatGPT)"
					subtitle = "Vector (ChatGPT)"
				case "chatgpt_anime":
					style = "z_external_provider_4"
					title = "Anime (ChatGPT)"
					subtitle = "Anime (ChatGPT)"
				case "chatgpt_print":
					style = "z_external_provider_5"
					title = "Print (ChatGPT)"
					subtitle = "Print (ChatGPT)"
				}
			}

			return map[string]any{
				"style": map[string]any{
					"identifier": style,
					"subtitle":   subtitle,
					"title":      title,
				},
			}
		},
	},
}

type focusMode struct {
	DisplayString string
	Identifier    string
}

var focusModes = map[string]focusMode{
	"Do Not Disturb": {
		DisplayString: "Do Not Disturb",
		Identifier:    "com.apple.donotdisturb.mode.default",
	},
	"Personal": {
		DisplayString: "Personal",
		Identifier:    "com.apple.focus.personal-time",
	},
	"Work": {
		DisplayString: "Work",
		Identifier:    "com.apple.focus.work",
	},
	"Sleep": {
		DisplayString: "Sleep",
		Identifier:    "com.apple.focus.sleep-mode",
	},
	"Driving": {
		DisplayString: "Driving",
		Identifier:    "com.apple.donotdisturb.mode.driving",
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

func handleRawParams(params map[string]any) {
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

var actionIncludes = []string{
	"a11y",
	"calendar",
	"contacts",
	"crypto",
	"device",
	"documents",
	"dropbox",
	"images",
	"location",
	"intelligence",
	"mac",
	"math",
	"media",
	"music",
	"network",
	"pdf",
	"photos",
	"scripting",
	"settings",
	"sharing",
	"shortcuts",
	"text",
	"translation",
	"web",
}

var includedStandardActions bool
var includedBasicStandardActions bool

func loadStandardActions() {
	includeStandardActions()
	includeBasicStandardActions()
	handleIncludes()
	handleActionDefinitions()
	resetParse()
}

func loadBasicStandardActions() {
	includeBasicStandardActions()
	handleIncludes()
	handleActionDefinitions()
}

func loadActionsByCategory() {
	for _, actionInclude := range actionCategories {
		lines = append(lines, fmt.Sprintf("#include 'actions/%s'\n", actionInclude))
		resetParse()
		handleIncludes()
		currentCategory = actionInclude
		handleActionDefinitions()

		included = []string{}
		includes = []include{}
		lines = []string{}
		tokens = []token{}

		resetParse()
	}
}

func includeBasicStandardActions() {
	if includedBasicStandardActions {
		return
	}
	lines = append([]string{"#include 'actions/basic'\n"}, lines...)
	resetParse()
	includedBasicStandardActions = true
}

func includeStandardActions() {
	if includedStandardActions {
		return
	}
	var standardIncludes []string
	for _, actionInclude := range actionIncludes {
		standardIncludes = append(standardIncludes, fmt.Sprintf("#include 'actions/%s'\n", actionInclude))
	}
	lines = append(standardIncludes, lines...)
	resetParse()
}

func checkMissingStandardInclude(identifier *string, parsing bool) {
	if !parsing && !args.Using("no-toolkit") {
		connectToolkitDB()
		var identifiers = strings.Split(*identifier, ".")
		identifiers = append(identifiers[:3], identifiers[4:]...)
		var baseIdentifier = strings.Join(identifiers, ".")

		var containerId, containerErr = getContainerIdByIdentifier(&baseIdentifier)
		if containerErr == nil {
			var containerName, containerErr = getContainerName(&containerId)
			importActions(baseIdentifier)
			if containerErr == nil {
				popLine(fmt.Sprintf("#import '%s'", containerName))
			} else {
				popLine(fmt.Sprintf("#import '%s'", baseIdentifier))
			}

			return
		}
	}

	for _, actionInclude := range actionIncludes {
		if slices.Contains(included, fmt.Sprintf("actions/%s", actionInclude)) {
			continue
		}
		lines = append([]string{fmt.Sprintf("#include 'actions/%s'\n", actionInclude)}, lines...)
		resetParse()
		handleIncludes()
		handleActionDefinitions()

		if !parsing {
			mapSplitActions()
		}

		var name, nameErr = getActionNameByIdentifier(identifier)
		if nameErr != nil {
			continue
		}

		var includeStatement = fmt.Sprintf("#include 'actions/%s'", actionInclude)
		if parsing {
			exit(fmt.Sprintf("Action '%s()' requires include:\n\n%s", name, includeStatement))
		} else {
			popLine(includeStatement)
			break
		}
	}
	return
}

func getActionNameByIdentifier(identifier *string) (name string, err error) {
	if actions[*identifier] == nil {
		var actionName, found = findActionByIdentifier(identifier)
		if !found {
			return "", fmt.Errorf("action name for '%s' not found", *identifier)
		}
		name = actionName
	} else {
		name = *identifier
	}
	return name, nil
}

func findActionByIdentifier(identifier *string) (name string, found bool) {
	for actionName, action := range actions {
		currentAction = *action
		var actionIdentifier = getFullActionIdentifier()
		if *identifier == actionIdentifier {
			return actionName, true
		}
	}
	return "", false
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
	case separator == "" && currentAction.identifier == "text.split":
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

var appIdentifierRegex = regexp.MustCompile(`^(.*?)\.(.*?)\.(.*?)$`)

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

// toggleSetActions are actions which all are state based and so can either be toggled or set in the same format.
var toggleSetActions = map[string]actionDefinition{
	"BackgroundSounds": {
		doc: selfDoc{
			title:    "Background Sounds",
			category: "a11y",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleBackgroundSoundsIntent",
	},
	"MediaBackgroundSounds": {
		doc: selfDoc{
			title:    "Media Background Sounds",
			category: "a11y",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleBackgroundSoundsIntent",
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"setting": "whenMediaIsPlaying",
			}
		},
	},
	"AutoAnswerCalls": {
		doc: selfDoc{
			title:    "Auto Answer Calls",
			category: "a11y",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleAutoAnswerCallsIntent",
	},
	"Appearance": {
		doc: selfDoc{
			category:    "settings",
			subcategory: "Appearance",
		},
		identifier: "appearance",
	},
	"Bluetooth": {
		doc: selfDoc{
			category:    "settings",
			subcategory: "Wireless",
		},
		identifier: "bluetooth.set",
		setKey:     "OnValue",
	},
	"Wifi": {
		doc: selfDoc{
			category:    "settings",
			subcategory: "Wireless",
		},
		identifier: "wifi.set",
		setKey:     "OnValue",
	},
	"CellularData": {
		doc: selfDoc{
			title:       "Cellular Data",
			category:    "settings",
			subcategory: "Wireless",
		},
		identifier: "cellulardata.set",
		setKey:     "OnValue",
	},
	"NightShift": {
		doc: selfDoc{
			title:       "Night Shift",
			category:    "settings",
			subcategory: "Display",
		},
		identifier: "nightshift.set",
		setKey:     "OnValue",
	},
	"TrueTone": {
		doc: selfDoc{
			title:       "True Tone",
			category:    "settings",
			subcategory: "Display",
		},
		identifier: "truetone.set",
		setKey:     "OnValue",
	},
	"AirplaneMode": {
		doc: selfDoc{
			title:    "Airplane Mode",
			category: "device",
		},
		identifier: "airplanemode.set",
		setKey:     "OnValue",
	},
	"ClassicInvert": {
		doc: selfDoc{
			title:       "Classic Invert",
			category:    "a11y",
			subcategory: "Vision",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleClassicInvertIntent",
	},
	"ClosedCaptionsSDH": {
		doc: selfDoc{
			title:       "Closed Captions SDH",
			category:    "a11y",
			subcategory: "Hearing",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleCaptionsIntent",
	},
	"ColorFilters": {
		doc: selfDoc{
			title:       "Color Filters",
			category:    "a11y",
			subcategory: "Vision",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleColorFiltersIntent",
	},
	"Contrast": {
		doc: selfDoc{
			category:    "ally",
			subcategory: "Vision",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleContrastIntent",
	},
	"LEDFlash": {
		doc: selfDoc{
			title:       "LED Flash",
			category:    "a11y",
			subcategory: "Hearing",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleLEDFlashIntent",
	},
	"LeftRightBalance": {
		doc: selfDoc{
			title:       "Left-Right Balance",
			category:    "a11y",
			subcategory: "Hearing",
		},
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
		doc: selfDoc{
			title:       "Live Captions",
			category:    "a11y",
			subcategory: "Hearing",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleLiveCaptionsIntent",
	},
	"MonoAudio": {
		doc: selfDoc{
			title:       "Mono Audio",
			category:    "a11y",
			subcategory: "Hearing",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleMonoAudioIntent",
	},
	"ReduceMotion": {
		doc: selfDoc{
			title:       "Reduce Motion",
			category:    "a11y",
			subcategory: "Vision",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleReduceMotionIntent",
	},
	"ReduceTransparency": {
		doc: selfDoc{
			title:       "Reduce Transparency",
			category:    "a11y",
			subcategory: "Vision",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleTransparencyIntent",
	},
	"SmartInvert": {
		doc: selfDoc{
			title:       "Smart Invert",
			category:    "a11y",
			subcategory: "Vision",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleSmartInvertIntent",
	},
	"SwitchControl": {
		doc: selfDoc{
			title:    "Switch Control",
			category: "a11y",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleSwitchControlIntent",
	},
	"VoiceControl": {
		doc: selfDoc{
			title:       "Voice Control",
			category:    "a11y",
			subcategory: "Vision",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleVoiceControlIntent",
	},
	"WhitePoint": {
		doc: selfDoc{
			title:       "White Point",
			category:    "a11y",
			subcategory: "Vision",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleWhitePointIntent",
	},
	"Zoom": {
		doc: selfDoc{
			category:    "a11y",
			subcategory: "Vision",
		},
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleZoomIntent",
	},
	"StageManager": {
		doc: selfDoc{
			title:    "Stage Manager",
			category: "settings",
		},
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
		var docTitle = def.doc.title
		var toggleName = fmt.Sprintf("toggle%s", name)
		def.addParams = func(_ []actionArgument) map[string]any {
			return map[string]any{
				"operation": "toggle",
			}
		}
		if docTitle != "" {
			def.doc.title = fmt.Sprintf("Toggle %s", docTitle)
		} else {
			def.doc.title = fmt.Sprintf("Toggle %s", name)
		}
		def.defaultAction = false

		var toggleDef = def
		actions[toggleName] = &toggleDef

		if name == "Appearance" {
			continue
		}

		var setName = fmt.Sprintf("set%s", name)
		def.defaultAction = true
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
		if docTitle != "" {
			def.doc.title = fmt.Sprintf("Set %s", docTitle)
		} else {
			def.doc.title = fmt.Sprintf("Set %s", name)
		}

		var setDef = def
		actions[setName] = &setDef
	}
	toggleSetActions = nil
}
