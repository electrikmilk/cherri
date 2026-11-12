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

// actions is the data structure that determines every action the compiler knows about.
// The key determines the identifier of the identifier that must be used in the syntax, it's value defines its behavior, etc. using an actionDefinition.
var actions = map[string]*actionDefinition{
	"createAlarm": {
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
		appendParamsFunc: func(args []actionArgument) (params map[string]any) {
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
	"emailAddress": {
		check: func(args []actionArgument, _ *actionDefinition) {
			if len(args) > 1 && args[0].valueType == Variable {
				parserError("Shortcuts only allows one variable for an email address.")
			}
		},
		makeParams: func(args []actionArgument) map[string]any {
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
		check: func(args []actionArgument, _ *actionDefinition) {
			if len(args) > 1 && args[0].valueType == Variable {
				parserError("Shortcuts only allows one variable for a phone number.")
			}
		},
		makeParams: func(args []actionArgument) map[string]any {
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
		appendParamsFunc: func(args []actionArgument) (params map[string]any) {
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
	"formRequest": {
		check: func(args []actionArgument, _ *actionDefinition) {
			if len(args) < 3 || args[2].valueType != Dict {
				return
			}

			var fields = args[2].value.(map[string]interface{})
			var formFields = make(map[string]interface{}, len(fields))
			for key, value := range fields {
				formFields[key] = value

				if text, isText := value.(string); isText {
					if name, isVariableOnly := variableOnlyValueName(text); isVariableOnly {
						if _, found := getVariableValue(name); !found {
							parserError(fmt.Sprintf("Undefined reference '%s'", name))
						}

						formFields[key] = varValue{value: name}
						continue
					}
				}

				checkNoFileItemSyntax(key, value)
			}

			args[2].value = formFields
		},
	},
	"labelFile": {
		appendParamsFunc: func(args []actionArgument) map[string]any {
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
		appendParamsFunc: func(args []actionArgument) (params map[string]any) {
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
		appendParamsFunc: func(args []actionArgument) map[string]any {
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
	"containsText": {
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
		appendParamsFunc: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{}
			}
			var folderPath = getArgValue(args[0])
			var pathParts = strings.Split(folderPath.(string), "/")
			var fileLocationType = pathParts[0]

			var filename = end(pathParts)
			// Shortcuts stores the subpath relative to the location root, including the folder's own name.
			folderPath = strings.Trim(strings.Join(pathParts[1:], "/"), "/")
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
		appendParamsFunc: textParts,
		decomp:           decompTextParts,
	},
	"joinText": {
		appendParamsFunc: textParts,
		decomp:           decompTextParts,
	},
	"url": {
		makeParams: func(args []actionArgument) map[string]any {
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
		makeParams: func(args []actionArgument) map[string]any {
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
		appendParamsFunc: func(args []actionArgument) map[string]any {
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
		check: func(args []actionArgument, definition *actionDefinition) {
			replaceAppIDs(args, definition)
		},
		appendParamsFunc: func(args []actionArgument) map[string]any {
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
		check: func(args []actionArgument, definition *actionDefinition) {
			replaceAppIDs(args, definition)
		},
		makeParams: func(args []actionArgument) map[string]any {
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
		check:      replaceAppIDs,
		makeParams: makeAllAppsAction,
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompAppAction("WFAppsExcept", action)
		},
	},
	"quitApp": {
		check: func(args []actionArgument, definition *actionDefinition) {
			replaceAppIDs(args, definition)
		},
		makeParams: func(args []actionArgument) map[string]any {
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
		check:      replaceAppIDs,
		makeParams: makeAllAppsAction,
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompAppAction("WFAppsExcept", action)
		},
	},
	"killApp": {
		check: func(args []actionArgument, definition *actionDefinition) {
			replaceAppIDs(args, definition)
		},
		makeParams: func(args []actionArgument) (params map[string]any) {
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
		check:      replaceAppIDs,
		makeParams: makeAllAppsAction,
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompAppAction("WFAppsExcept", action)
		},
	},
	"splitApps": {
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
		appendParamsFunc: func(args []actionArgument) map[string]any {
			if len(args) < 2 {
				return map[string]any{}
			}

			return map[string]any{
				"WFPrimaryAppIdentifier":   appIdentifierValue(args, 0),
				"WFSecondaryAppIdentifier": appIdentifierValue(args, 1),
			}
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
		appendParamsFunc: func(args []actionArgument) map[string]any {
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
		appendParamsFunc: func(args []actionArgument) map[string]any {
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
		makeParams: func(args []actionArgument) map[string]any {
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
		appendParamsFunc: func(args []actionArgument) (xCallbackParams map[string]any) {
			if len(args) == 0 {
				return
			}
			xCallbackParams = make(map[string]any)
			if len(args) >= 4 {
				if args[1].value.(string) != "" || args[2].value.(string) != "" || args[3].value.(string) != "" {
					xCallbackParams["WFXCallbackCustomCallbackEnabled"] = true
				}
			}
			if len(args) == 5 {
				if args[4].value.(string) != "" {
					xCallbackParams["WFXCallbackCustomSuccessURLEnabled"] = true
				}
			}

			return
		},
	},
	"getWindows": {
		appendParamsFunc: func(args []actionArgument) (params map[string]any) {
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
	},
	"convertMeasurement": {
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
		appendParamsFunc: func(args []actionArgument) map[string]any {
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
		appendParamsFunc: func(args []actionArgument) map[string]any {
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
		makeParams: func(args []actionArgument) map[string]any {
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
		makeParams: func(args []actionArgument) map[string]any {
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
		check: func(args []actionArgument, _ *actionDefinition) {
			if args[1].valueType == Variable {
				parserError("Contact detail must be a literal value as it determines which parameter is set")
			}
		},
		appendParamsFunc: func(args []actionArgument) map[string]any {
			if len(args) < 3 {
				return map[string]any{}
			}

			var detailKey = strings.ReplaceAll(getArgValue(args[1]).(string), " ", "")

			return map[string]any{
				"WFContactContentItem" + detailKey: argumentValue(args, 2),
			}
		},
	},
	"setWindowedMultitasking": {
		appendParamsFunc: func(args []actionArgument) map[string]any {
			return appendSetMultitaskingModeParam("Windowed Apps", "windowedApps", "macwindow.on.rectangle")
		},
	},
	"setStageManagerMultitasking": {
		appendParamsFunc: func(args []actionArgument) map[string]any {
			return appendSetMultitaskingModeParam("Stage Manager", "stageManager", "squares.leading.rectangle")
		},
	},
	"setFocusMode": {
		appendParamsFunc: func(args []actionArgument) map[string]any {
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
		appendParamsFunc: func(args []actionArgument) map[string]any {
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
		appendParamsFunc: func(args []actionArgument) map[string]any {
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
	"settings",
	"sharing",
	"shortcuts",
	"storage",
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

// checkMissingStandardInclude returns the include statement that defines the action when parsing,
// and adds the include to the decompiled source otherwise.
func checkMissingStandardInclude(identifier *string, parsing bool) (requiredInclude string) {
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

			return ""
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

		var _, nameErr = getActionNameByIdentifier(identifier)
		if nameErr != nil {
			continue
		}

		var includeStatement = fmt.Sprintf("#include 'actions/%s'", actionInclude)
		if parsing {
			return includeStatement
		}
		popLine(includeStatement)
		break
	}
	return ""
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
		currentAction.definition = *action
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

	var separator any = "\n"
	if len(args) > 1 {
		separator = getArgValue(args[1])
	}

	switch {
	case separator == " ":
		data["WFTextSeparator"] = "Spaces"
	case separator == "\n":
		data["WFTextSeparator"] = "New Lines"
	case separator == "" && currentAction.definition.identifier == "text.split":
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
	var hasCustomSeparator bool
	if action.WFWorkflowActionParameters["WFTextSeparator"] != nil {
		glue = decompReferenceValue(action.WFWorkflowActionParameters["WFTextSeparator"])
	}
	if action.WFWorkflowActionParameters["WFTextCustomSeparator"] != nil {
		glue = decompReferenceValue(action.WFWorkflowActionParameters["WFTextCustomSeparator"])
		hasCustomSeparator = true
	}

	if glue != "" || hasCustomSeparator {
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

func appArgument(args []actionArgument) (appArg []map[string]any) {
	for _, arg := range args {
		if arg.valueType != Variable {
			appArg = append(appArg, map[string]any{
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

func appIdentifierValue(args []actionArgument, index int) any {
	if args[index].valueType == Variable {
		return argumentValue(args, index)
	}

	return map[string]any{
		"BundleIdentifier": argumentValue(args, index),
	}
}

// A value that looks like a file item but cannot be one would otherwise silently compile to text with a literal $.
func checkNoFileItemSyntax(key string, value any) {
	switch v := value.(type) {
	case string:
		if strings.HasPrefix(v, "${") {
			parserError(fmt.Sprintf("Form request field '%s' must be only a single variable reference to be a file, e.g. \"${@variable}\"", key))
		}
	case map[string]interface{}:
		for nestedKey, nestedValue := range v {
			checkNoFileItemSyntax(nestedKey, nestedValue)
		}
	case []interface{}:
		for _, item := range v {
			checkNoFileItemSyntax(key, item)
		}
	}
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

func makeAllAppsAction(args []actionArgument) (params map[string]any) {
	if len(args) == 0 {
		return map[string]any{}
	}

	params = make(map[string]any)

	if args[0].valueType != Variable {
		params["WFAppsExcept"] = appArgument(args)
	} else {
		params["WFAppsExcept"] = argumentValue(args, 0)
	}

	return
}

func appendSetMultitaskingModeParam(title, identifier, symbol string) map[string]any {
	return map[string]any{
		"mode": map[string]any{
			"identifier": identifier,
			"subtitle": map[string]string{
				"key": title,
			},
			"symbol": map[string]string{
				"systemName": symbol,
			},
			"title": map[string]string{
				"key": title,
			},
			"value": identifier,
		},
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
