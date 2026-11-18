/*
 * Copyright (c) Cherri
 */

package main

import (
	"encoding/base64"
	"fmt"
	"maps"
	"os"
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

			// The decompiler writes nil for a skipped optional argument, which means no repeat days.
			var repeatDays, _ = getArgValue(args[3]).([]interface{})
			for _, day := range repeatDays {
				if !slices.Contains(weekdays, strings.ToLower(day.(string))) {
					parserError(fmt.Sprintf("Invalid repeat weekday for alarm '%s'", day))
				}
			}
		},
		appendParamsFunc: func(args []actionArgument) (params map[string]any) {
			if len(args) < 4 {
				return
			}

			var repeatDays, _ = getArgValue(args[3]).([]interface{})
			var repeats []map[string]any
			for _, day := range repeatDays {
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
		decomp: func(action *ShortcutAction) (arguments []string) {
			var recovered = make(map[string]string)
			if _, found := optionalParam(action, "repeats"); !found {
				return decompHookArguments(action, recovered)
			}

			var weekdays []interface{}
			for _, repeat := range paramSlice(action, "repeats") {
				var repeatDay, isMap = repeat.(map[string]interface{})
				var day, isString = repeatDay["value"].(string)
				if !isMap || !isString {
					decompError("Expected each of repeats to be a dictionary with a value string", action)
				}
				weekdays = append(weekdays, day)
			}
			recovered["repeatWeekdays"] = decompLiteral(weekdays)

			return decompHookArguments(action, recovered)
		},
	},
	"addWorldClock": {
		appendParamsFunc: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return nil
			}

			return map[string]any{
				"city": worldClockCityValue(args[0]),
			}
		},
	},
	"removeWorldClock": {
		check: checkWorldClockCities,
		appendParamsFunc: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return nil
			}

			return map[string]any{
				"entities": worldClockCitiesValue(args[0]),
			}
		},
	},
	"getCityTime": {
		appendParamsFunc: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return nil
			}

			return map[string]any{
				"city": worldClockCityValue(args[0]),
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
		decomp: func(action *ShortcutAction) (arguments []string) {
			var recovered = make(map[string]string)
			for name, contact := range map[string]struct {
				key        string
				contentKit contentKit
			}{
				"phoneNumber":  {"WFContactPhoneNumbers", phoneNumber},
				"emailAddress": {"WFContactEmails", emailAddress},
			} {
				var values = decompContactValue(action, contact.key, contact.contentKit)
				if len(values) > 1 {
					decompError(fmt.Sprintf("Expected the parameter '%s' to hold one value, got %d", contact.key, len(values)), action)
				}
				if len(values) == 1 {
					recovered[name] = values[0]
				}
			}

			return decompHookArguments(action, recovered)
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
		decomp: func(action *ShortcutAction) (arguments []string) {
			var colorNumber = expectType[uint64](requiredParam(action, "WFLabelColorNumber"), "WFLabelColorNumber", "an integer", action)
			var color, found = firstKeyWithValue(int(colorNumber), fileLabelsMap)
			if !found {
				decompError(fmt.Sprintf("Expected WFLabelColorNumber to be a file label color number, got %v", colorNumber), action)
			}

			return decompHookArguments(action, map[string]string{"color": decompValue(color)})
		},
	},
	"filterFiles": {
		derivedKeys: []string{"WFContentItemLimitEnabled"},
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
		check: func(args []actionArgument, _ *actionDefinition) {
			checkLiteralArgument[bool](args, 1, "Rich text must be a literal value, as it determines which text type is written.")
		},
		appendParamsFunc: func(args []actionArgument) map[string]any {
			if len(args) == 0 {
				return map[string]any{}
			}

			if len(args) > 1 {
				// A nil richText is left out like a false one.
				if richText, isBool := getArgValue(args[1]).(bool); isBool && richText {
					return map[string]any{
						"WFGetTextFromPDFTextType": "Rich Text",
					}
				}
			}

			return map[string]any{
				"WFGetTextFromPDFTextType": "Text",
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			var recovered = make(map[string]string)
			if textType, found := optionalParam(action, "WFGetTextFromPDFTextType"); found {
				// Text is written whether richText is false or left out, so only Rich Text needs the argument.
				switch textType {
				case "Rich Text":
					recovered["richText"] = "true"
				case "Text":
				default:
					decompError(fmt.Sprintf("Expected WFGetTextFromPDFTextType to be 'Text' or 'Rich Text', got '%v'", textType), action)
				}
			}

			return decompHookArguments(action, recovered)
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
		decomp: func(action *ShortcutAction) (arguments []string) {
			var file = paramMap(action, "WFFile")
			var fileLocation, hasLocation = file["fileLocation"].(map[string]interface{})
			if !hasLocation {
				decompError(fmt.Sprintf("Expected WFFile to have a fileLocation dictionary, got %T", file["fileLocation"]), action)
			}
			var subpath, hasSubpath = fileLocation["relativeSubpath"].(string)
			if !hasSubpath {
				decompError(fmt.Sprintf("Expected the fileLocation of WFFile to have a relativeSubpath string, got %T", fileLocation["relativeSubpath"]), action)
			}

			// Only a home folder location is stored, the first part of any other folder path is not, so
			// Shortcuts stands in for it.
			var root = "Shortcuts"
			if locationType, found := fileLocation["WFFileLocationType"]; found {
				if locationType != "Home" {
					decompError(fmt.Sprintf("Expected WFFileLocationType of WFFile to be 'Home', got '%v'", locationType), action)
				}
				root = "~"
			}

			var folder = root
			if subpath != "" {
				folder = root + "/" + subpath
			}

			return decompHookArguments(action, map[string]string{
				"folder": decompValue(folder),
			})
		},
	},
	"splitText": {
		derivedKeys:      []string{"Show-text"},
		appendParamsFunc: textParts,
		decomp:           decompTextParts,
	},
	"joinText": {
		derivedKeys:      []string{"Show-text"},
		appendParamsFunc: textParts,
		decomp:           decompTextParts,
	},
	"url": {
		derivedKeys: []string{"Show-WFURLActionURL"},
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
		decomp: decompInfiniteURLAction("WFURLActionURL"),
	},
	"addToReadingList": {
		derivedKeys: []string{"Show-WFURLActionURL"},
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
		decomp: decompInfiniteURLAction("WFURL"),
	},
	"prompt": {
		derivedKeys: []string{"WFAskActionDefaultAnswerNumber"},
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
		decomp: func(action *ShortcutAction) (arguments []string) {
			var recovered = make(map[string]string)
			if defaultAnswer, found := optionalParam(action, "WFAskActionDefaultAnswer"); found {
				recovered["defaultValue"] = decompValue(defaultAnswer)
			}

			return decompHookArguments(action, recovered)
		},
	},
	"openApp": {
		derivedKeys: []string{"WFSelectedApp"},
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
				var ratioName, _ = args[2].value.(string)
				if ratio, found := splitAppRatios[ratioName]; found {
					args[2].value = ratio
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
			// Both apps are required, as the arguments are written by position.
			requiredParam(action, "WFPrimaryAppIdentifier")
			requiredParam(action, "WFSecondaryAppIdentifier")
			arguments = append(arguments, decompAppAction("WFPrimaryAppIdentifier", action)...)
			arguments = append(arguments, decompAppAction("WFSecondaryAppIdentifier", action)...)

			// Shortcuts leaves the ratio out when it is the default.
			var splitRatio, hasRatio = optionalParam(action, "WFAppRatio")
			if !hasRatio {
				return
			}

			switch ratio := splitRatio.(type) {
			case string:
				var ratioName, found = firstKeyWithValue(ratio, splitAppRatios)
				if !found {
					decompError(fmt.Sprintf("Expected WFAppRatio to be '½ + ½' or '⅓ + ⅔', got '%s'", ratio), action)
				}
				arguments = append(arguments, doubleQuoted(ratioName))
			case map[string]interface{}:
				arguments = append(arguments, decompValue(ratio))
			default:
				decompError(fmt.Sprintf("Expected the parameter 'WFAppRatio' to be a string or a variable, got %T", splitRatio), action)
			}

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
			if readParam(action, "target") != nil {
				arguments = append(arguments, decompShortcutTitle(action))
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
			if readParam(action, "WFInput") == nil {
				return []string{"nil"}
			}

			return []string{decompValue(readParam(action, "WFInput"))}
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
			for _, item := range paramSlice(action, "WFItems") {
				var itemValue = item
				switch listItem := item.(type) {
				case string:
				case map[string]interface{}:
					itemValue = listItem["WFValue"]
				default:
					decompError(fmt.Sprintf("Expected an item of WFItems to be a string or a dictionary, got %T", item), action)
				}
				arguments = append(arguments, decompValue(itemValue))
			}
			return
		},
	},
	"openCustomXCallbackURL": {
		derivedKeys: []string{"WFXCallbackCustomCallbackEnabled", "WFXCallbackCustomSuccessURLEnabled"},
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
		derivedKeys: []string{"WFContentItemLimitEnabled"},
		appendParamsFunc: func(args []actionArgument) (params map[string]any) {
			if len(args) == 0 {
				return map[string]any{}
			}

			if len(args) > 2 && args[2].value != nil {
				params = map[string]any{
					"WFContentItemLimitEnabled": true,
				}
			}
			return
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			if len(args) > 1 && args[1].value != nil {
				// A variable order is written as it is, but a literal one is written by the kind of property.
				if args[1].valueType != Variable {
					checkLiteralArgument[string](args, 0, "Sort property must be a literal value when the sort order is, as it determines how the order is written.")
				}
				// A nil sortBy matches no sort property, so the order is kept as is.
				var sortBy, _ = getArgValue(args[0]).(string)
				var orderBy, _ = getArgValue(args[1]).(string)
				for _, sortOrder := range windowSortOrders {
					if !slices.Contains(sortOrder.properties, sortBy) {
						continue
					}
					switch orderBy {
					case "asc":
						args[1].value = sortOrder.ascending
					case "desc":
						args[1].value = sortOrder.descending
					}
					break
				}
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			var recovered = make(map[string]string)
			// Any other sort order is read as is, a variable or the asc or desc kept for a random sort.
			var storedOrder, _ = readParam(action, "WFContentItemSortOrder").(string)
			for _, sortOrder := range windowSortOrders {
				switch storedOrder {
				case sortOrder.ascending:
					recovered["orderBy"] = decompValue("asc")
				case sortOrder.descending:
					recovered["orderBy"] = decompValue("desc")
				}
			}

			return decompHookArguments(action, recovered)
		},
	},
	"convertMeasurement": {
		check: func(args []actionArgument, _ *actionDefinition) {
			var unitType, isString = getArgValue(args[1]).(string)
			if !isString {
				return
			}

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
				decompValue(readParam(action, "WFInput")),
				decompValue(readParam(action, "WFMeasurementUnitType")),
			)

			if readParam(action, "WFMeasurementUnit") != nil {
				var measurementUnit WFMeasurementUnit
				decompMapToStruct(paramMap(action, "WFMeasurementUnit"), &measurementUnit, "WFMeasurementUnit", action)

				arguments = append(arguments, decompValue(measurementUnit.WFNSUnitSymbol))
			}

			return
		},
	},
	"measurement": {
		check: func(args []actionArgument, _ *actionDefinition) {
			var unitType, isString = getArgValue(args[1]).(string)
			if !isString {
				return
			}

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
			if readParam(action, "WFMeasurementUnit") != nil {
				var measurementUnit WFMeasurementUnit
				decompMapToStruct(paramMap(action, "WFMeasurementUnit"), &measurementUnit, "WFMeasurementUnit", action)

				arguments = append(arguments,
					decompValue(measurementUnit.Value.Magnitude),
					decompValue(readParam(action, "WFMeasurementUnitType")),
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
				var imageValue, isString = getArgValue(args[2]).(string)
				if !isString && args[2].valueType == Variable {
					photo = fmt.Sprintf("{%s}", makeVariableReferenceString(args[2].value.(varValue)))
				} else {
					photo = imageValue
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
			var file, isString = getArgValue(args[0]).(string)
			if args[0].valueType == Variable && !isString {
				parserError("File path must be a string literal")
			}
			if _, err := os.Stat(file); os.IsNotExist(err) {
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
		decomp: func(action *ShortcutAction) (arguments []string) {
			var detailKey = "WFContactContentItem" + strings.ReplaceAll(paramString(action, "WFContentItemPropertyName"), " ", "")

			return decompHookArguments(action, map[string]string{
				"value": decompValue(requiredParam(action, detailKey)),
			})
		},
	},
	"setWindowedMultitasking": {
		derivedKeys: []string{"mode"},
		appendParamsFunc: func(args []actionArgument) map[string]any {
			return appendSetMultitaskingModeParam("Windowed Apps", "windowedApps", "macwindow.on.rectangle")
		},
	},
	"setStageManagerMultitasking": {
		derivedKeys: []string{"mode"},
		appendParamsFunc: func(args []actionArgument) map[string]any {
			return appendSetMultitaskingModeParam("Stage Manager", "stageManager", "squares.leading.rectangle")
		},
	},
	"setFocusMode": {
		check:            checkFocusModeLiteral,
		appendParamsFunc: focusModeParams,
		decomp:           decompFocusModeAction,
	},
	"toggleFocusMode": {
		check:            checkFocusModeLiteral,
		appendParamsFunc: focusModeParams,
		decomp:           decompFocusModeAction,
	},
	"generateImage": {
		appendParamsFunc: func(args []actionArgument) map[string]any {
			if len(args) < 3 {
				return map[string]any{}
			}

			var title string
			var style = argumentValue(args, 2)
			if args[2].valueType == String {
				style = getArgValue(args[2])
				if externalStyle, found := externalImageStyles[style.(string)]; found {
					style = externalStyle.identifier
					title = externalStyle.title
				}
			}

			return map[string]any{
				"style": map[string]any{
					"identifier": style,
					"subtitle":   title,
					"title":      title,
				},
			}
		},
		decomp: func(action *ShortcutAction) (arguments []string) {
			var recovered = make(map[string]string)
			if _, found := optionalParam(action, "style"); !found {
				return decompHookArguments(action, recovered)
			}

			var identifier = paramMap(action, "style")["identifier"]
			recovered["style"] = decompValue(identifier)
			for name, externalStyle := range externalImageStyles {
				if identifier == externalStyle.identifier {
					recovered["style"] = decompValue(name)
				}
			}

			return decompHookArguments(action, recovered)
		},
	},
}

type externalImageStyle struct {
	identifier string
	title      string
}

// externalImageStyles maps the Image Playground styles of external providers to the identifier and title
// Shortcuts stores for them, and is read both ways so a decompiled style gets its Cherri name back.
var externalImageStyles = map[string]externalImageStyle{
	"chatgpt":              {"z_external_provider", "ChatGPT"},
	"chatgpt_oil_painting": {"z_external_provider_1", "Oil Painting (ChatGPT)"},
	"chatgpt_watercolor":   {"z_external_provider_2", "Watercolor (ChatGPT)"},
	"chatgpt_vector":       {"z_external_provider_3", "Vector (ChatGPT)"},
	"chatgpt_anime":        {"z_external_provider_4", "Anime (ChatGPT)"},
	"chatgpt_print":        {"z_external_provider_5", "Print (ChatGPT)"},
}

// splitAppRatios maps the ratios of splitApps to the text Shortcuts stores for them, and is read both ways so a
// decompiled ratio gets its Cherri name back.
var splitAppRatios = map[string]string{
	"half":       "½ + ½",
	"thirdByTwo": "⅓ + ⅔",
}

// windowSortOrder is the sort orders Shortcuts stores for an ascending and a descending sort by one of the
// properties, as the order is named for the kind of property.
type windowSortOrder struct {
	properties []string
	ascending  string
	descending string
}

// windowSortOrders is read both ways so a decompiled sort order gets its Cherri name back. A random sort has no
// order.
var windowSortOrders = []windowSortOrder{
	{[]string{"Title", "App Name", "Name"}, "A to Z", "Z to A"},
	{[]string{"Width", "Height", "X Position", "Y Position", "Window Index"}, "Biggest First", "Smallest First"},
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

// decompFocusMode recovers the focusMode argument from FocusModes, which is left out for the default.
func decompFocusMode(action *ShortcutAction) map[string]string {
	var recovered = make(map[string]string)
	if _, found := optionalParam(action, "FocusModes"); !found {
		return recovered
	}

	var modes = paramMap(action, "FocusModes")
	for name, mode := range focusModes {
		if modes["Identifier"] == mode.Identifier {
			recovered["focusMode"] = decompValue(name)
			return recovered
		}
	}

	decompError(fmt.Sprintf("Expected the Identifier of FocusModes to identify a known focus mode, got '%v'", modes["Identifier"]), action)
	return recovered
}

func checkFocusModeLiteral(args []actionArgument, _ *actionDefinition) {
	checkLiteralArgument[string](args, 0, "Focus mode must be a literal value, as it determines which focus is set.")
}

func focusModeParams(args []actionArgument) map[string]any {
	if len(args) > 0 {
		// A nil focusMode is left out for the default.
		var mode, _ = getArgValue(args[0]).(string)
		if fm, found := focusModes[mode]; found {
			return map[string]any{
				"FocusModes": fm,
			}
		}
	}

	return map[string]any{}
}

func decompFocusModeAction(action *ShortcutAction) []string {
	return decompHookArguments(action, decompFocusMode(action))
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
		prependGeneratedLines(fmt.Sprintf("#include 'actions/%s'\n", actionInclude))
		resetParse()
		handleIncludes()
		currentCategory = actionInclude
		handleActionDefinitions()

		included = []string{}
		lineOrigins = nil
		lines = []string{}
		tokens = []token{}

		resetParse()
	}
}

func includeBasicStandardActions() {
	if includedBasicStandardActions {
		return
	}
	prependGeneratedLines("#include 'actions/basic'\n")
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
	prependGeneratedLines(standardIncludes...)
	resetParse()
}

// standardActionIncludes maps the name of an action to the standard include that defines it. Looking for an include
// leaves every category it passed loaded, so a later action from one of them resolves without checkMissingStandardInclude
// and its include is found here instead.
var standardActionIncludes = map[string]string{}

// checkMissingStandardInclude returns the include statement that defines the action when parsing, and loads the
// definitions of the action otherwise. Standard actions win over actions imported from the ToolKit database.
func checkMissingStandardInclude(identifier *string, parsing bool) (requiredInclude string) {
	for _, actionInclude := range actionIncludes {
		if slices.Contains(included, fmt.Sprintf("actions/%s", actionInclude)) {
			continue
		}

		var definedBefore = maps.Clone(actions)
		prependGeneratedLines(fmt.Sprintf("#include 'actions/%s'\n", actionInclude))
		resetParse()
		handleIncludes()
		handleActionDefinitions()

		for name, definition := range actions {
			if definedBefore[name] != definition {
				standardActionIncludes[name] = actionInclude
			}
		}

		if !parsing {
			mapSplitActions()
		}

		var _, nameErr = getActionNameByIdentifier(identifier)
		if nameErr != nil {
			continue
		}

		if parsing {
			return fmt.Sprintf("#include 'actions/%s'", actionInclude)
		}
		return ""
	}

	if !parsing {
		importToolkitActions(identifier)
		mapSplitActions()
	}

	return ""
}

var decompiledIncludes []string

// includeStandardAction adds the include that defines the action to the decompiled source once.
func includeStandardAction(name string) {
	var actionInclude, found = standardActionIncludes[name]
	if !found || slices.Contains(decompiledIncludes, actionInclude) {
		return
	}

	decompiledIncludes = append(decompiledIncludes, actionInclude)
	decompiledImports = append(decompiledImports, fmt.Sprintf("#include 'actions/%s'", actionInclude))
}

// importToolkitActions imports the actions of the app that the action belongs to, which is the identifier without the
// name of the action. The identifier must have at least two parts.
func importToolkitActions(identifier *string) {
	if args.Using("no-toolkit") {
		return
	}

	var identifiers = strings.Split(*identifier, ".")
	var baseIdentifier = strings.Join(identifiers[:len(identifiers)-1], ".")
	if slices.Contains(imported, baseIdentifier) {
		return
	}

	connectToolkitDB()

	var containerId, containerErr = getContainerIdByIdentifier(&baseIdentifier)
	if containerErr != nil {
		return
	}

	importActions(baseIdentifier)

	var containerName, nameErr = getContainerName(&containerId)
	if nameErr != nil {
		decompiledImports = append(decompiledImports, fmt.Sprintf("#import '%s'", baseIdentifier))
		return
	}
	decompiledImports = append(decompiledImports, fmt.Sprintf("#import '%s'", containerName))
}

func getActionNameByIdentifier(identifier *string) (name string, err error) {
	if !callableAction(*identifier) {
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
	if readParam(action, key) == nil {
		return
	}

	var wfValue WFValue
	decompMapToStruct(readParam(action, key), &wfValue, key, action)

	if wfValue.WFSerializationType == "WFTextTokenString" {
		return []string{decompValue(readParam(action, key))}
	}

	var value, hasValue = wfValue.Value.(map[string]interface{})
	if !hasValue {
		decompError(fmt.Sprintf("Expected the parameter '%s' to have a Value dictionary, got %T", key, wfValue.Value), action)
	}
	var fieldValues, hasFieldValues = value["WFContactFieldValues"].([]interface{})
	if !hasFieldValues {
		decompError(fmt.Sprintf("Expected the parameter '%s' to have a WFContactFieldValues array, got %T", key, value["WFContactFieldValues"]), action)
	}

	var contactFieldValues []WFContactFieldValue
	decompMapToStruct(fieldValues, &contactFieldValues, "WFContactFieldValues", action)

	for _, contactFieldValue := range contactFieldValues {
		var serializedKey = fmt.Sprintf("link.contentkit.%s", contentKit)
		var email, isString = contactFieldValue.SerializedEntry[serializedKey].(string)
		if !isString {
			decompError(fmt.Sprintf("Expected a contact field value to have a text entry '%s', got %T", serializedKey, contactFieldValue.SerializedEntry[serializedKey]), action)
		}
		arguments = append(arguments, quoteString(email, false))
	}

	return
}

// checkLiteralArgument stops on a variable argument at index that is not set to a literal of type T, as an
// argument that decides which parameters are written cannot be read from a value only known when the
// Shortcut runs. A nil argument is not a variable, so it is still left out.
func checkLiteralArgument[T any](args []actionArgument, index int, message string) {
	if len(args) <= index || args[index].valueType != Variable {
		return
	}

	if _, isLiteral := getArgValue(args[index]).(T); !isLiteral {
		parserError(message)
	}
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

// decompSeparator writes a separator as a string literal, or as the text or reference it holds. Only a named
// separator, and not a custom one, is converted to the character it names, as a custom separator is the
// text it holds.
func decompSeparator(paramValue any, named bool, action *ShortcutAction) string {
	if _, isMap := paramValue.(map[string]interface{}); isMap {
		return decompValue(paramValue)
	}

	var text, isString = paramValue.(string)
	if !isString {
		decompError(fmt.Sprintf("Expected a text separator to be text or a variable reference, got %T", paramValue), action)
	}
	if named {
		text = glueToChar(text)
	}

	return quoteString(text, false)
}

func decompTextParts(action *ShortcutAction) (arguments []string) {
	arguments = append(arguments, decompValue(readParam(action, "text")))

	var separator = readParam(action, "WFTextSeparator")
	if customSeparator, isCustom := optionalParam(action, "WFTextCustomSeparator"); isCustom && customSeparator != nil {
		arguments = append(arguments, decompSeparator(customSeparator, false, action))
		return
	}
	if separator != nil && separator != "" {
		arguments = append(arguments, decompSeparator(separator, true, action))
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

// decompAppAction writes the apps of key, which is a bundle identifier, a variable, or for the actions that take
// every app except some, an array of bundle identifiers that is left out when there are no exceptions.
func decompAppAction(key string, action *ShortcutAction) (arguments []string) {
	var value, found = optionalParam(action, key)
	if !found {
		return
	}

	switch app := value.(type) {
	case string:
		return append(arguments, decompValue(app))
	case map[string]interface{}:
		return append(arguments, decompApp(app, key, action))
	case []interface{}:
		for _, item := range app {
			arguments = append(arguments, decompApp(expectType[map[string]interface{}](item, fmt.Sprintf("an app of '%s'", key), "a dictionary", action), key, action))
		}
		return
	}

	decompError(fmt.Sprintf("Expected the parameter '%s' to be a string, a dictionary or an array of dictionaries, got %T", key, value), action)
	return
}

// decompApp writes an app that is either a bundle identifier dictionary or a variable.
func decompApp(app map[string]interface{}, key string, action *ShortcutAction) string {
	if bundleIdentifier, found := app["BundleIdentifier"]; found {
		return decompValue(bundleIdentifier)
	}
	if app["Value"] == nil {
		decompError(fmt.Sprintf("Expected an app of '%s' to have a BundleIdentifier or a variable Value", key), action)
	}

	return decompValue(app)
}

// decompInfiniteURLAction writes the URLs of key, which is a single value or an array of values.
func decompInfiniteURLAction(key string) func(action *ShortcutAction) []string {
	return func(action *ShortcutAction) (arguments []string) {
		var value = requiredParam(action, key)
		switch urls := value.(type) {
		case string, map[string]interface{}:
			return append(arguments, decompValue(urls))
		case []interface{}:
			for _, url := range urls {
				arguments = append(arguments, decompValue(url))
			}
			return
		}

		decompError(fmt.Sprintf("Expected the parameter '%s' to be a string, a dictionary or an array, got %T", key, value), action)
		return
	}
}

// decompShortcutTitle writes the title of the target of an Open Shortcut action. Shortcuts stores the title as a
// localized string dictionary, while a shortcut name compiled by Cherri is a plain or interpolated string.
func decompShortcutTitle(action *ShortcutAction) string {
	var target = paramMap(action, "target")
	switch title := target["title"].(type) {
	case string:
		return decompValue(title)
	case map[string]interface{}:
		if key, found := title["key"]; found {
			return decompValue(key)
		}
		if title["Value"] == nil {
			decompError("Expected the title of the target to have a key or a variable Value", action)
		}

		return decompValue(title)
	default:
		decompError(fmt.Sprintf("Expected the title of the target to be a string or a dictionary, got %T", title), action)
		return ""
	}
}
