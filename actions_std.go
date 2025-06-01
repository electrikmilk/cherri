/*
 * Copyright (c) Cherri
 */

package main

import (
	"encoding/base64"
	"fmt"
	"maps"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/google/uuid"
)

const SetVariableIdentifier = "is.workflow.actions.setvariable"
const AppendVariableIdentifier = "is.workflow.actions.appendvariable"

// FIXME: Some of these actions have a value with a set list values for an arguments,
//  but the argument value is not being checked against its possible values.
//  Use the "hash" action as an example.

var measurementUnitTypes = []string{"Acceleration", "Angle", "Area", "Concentration Mass", "Dispersion", "Duration", "Electric Charge", "Electric Current", "Electric Potential Difference", "V Electric Resistance", "Energy", "Frequency", "Fuel Efficiency", "Illuminance", "Information Storage", "Length", "Mass", "Power", "Pressure", "Speed", "Temperature", "Volume"}
var units map[string][]string
var abcSortOrders = []string{"A to Z", "Z to A"}
var contactDetails = []string{"First Name", "Middle Name", "Last Name", "Birthday", "Prefix", "Suffix", "Nickname", "Phonetic First Name", "Phonetic Last Name", "Phonetic Middle Name", "Company", "Job Title", "Department", "File Extension", "Creation Date", "File Path", "Last Modified Date", "Name", "Random"}
var facetimeCallTypes = []string{"Video", "Audio"}
var stopListening = []string{"After Pause", "After Short Pause", "On Tap"}
var errorCorrectionLevels = []string{"Low", "Medium", "Quartile", "High"}
var archiveTypes = []string{".zip", ".tar.gz", ".tar.bz2", ".tar.xz", ".tar", ".gz", ".cpio", ".iso"}
var storageUnits = []string{"bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
var weatherDetails = []string{"Name", "Air Pollutants", "Air Quality Category", "Air Quality Index", "Sunset Time", "Sunrise Time", "UV Index", "Wind Direction", "Wind Speed", "Precipitation Chance", "Precipitation Amount", "Pressure", "Humidity", "Dewpoint", "Visibility", "Condition", "Feels Like", "Low", "High", "Temperature", "Location", "Date"}
var weatherForecastTypes = []string{"Daily", "Hourly"}
var locationDetails = []string{"Name", "URL", "Label", "Phone Number", "Region", "ZIP Code", "State", "City", "Street", "Altitude", "Longitude", "Latitude"}
var flipDirections = []string{"Horizontal", "Vertical"}
var recordingQualities = []string{"Normal", "Very High"}
var recordingStarts = []string{"On Tap", "Immediately"}
var imageFormats = []string{"TIFF", "GIF", "PNG", "BMP", "PDF", "HEIF"}
var audioFormats = []string{"M4A", "AIFF"}
var speeds = []string{"0.5X", "Normal", "2X"}
var videoSizes = []string{"640×480", "960×540", "1280×720", "1920×1080", "3840×2160", "HEVC 1920×1080", "HEVC 3840x2160", "ProRes 422"}
var selectionTypes = []string{"Window", "Custom"}
var fileSizeUnits = []string{"Closest Unit", "Bytes", "Kilobytes", "Megabytes", "Gigabytes", "Terabytes", "Petabytes", "Exabytes", "Zettabytes", "Yottabytes"}
var deviceDetails = []string{"Device Name", "Device Hostname", "Device Model", "Device Is Watch", "System Version", "Screen Width", "Screen Height", "Current Volume", "Current Brightness", "Current Appearance"}
var hashTypes = []string{"MD5", "SHA1", "SHA256", "SHA512"}
var inputTypes = []string{"Text", "Number", "URL", "Date", "Time", "Date and Time"}
var appSplitRatios = []string{"half", "thirdByTwo"}
var calculationOperations = []string{"x^2", "х^3", "x^у", "e^x", "10^x", "In(x)", "log(x)", "√x", "∛x", "x!", "sin(x)", "cos(X)", "tan(x)", "abs(x)"}
var statisticsOperations = []string{"Average", "Minimum", "Maximum", "Sum", "Median", "Mode", "Range", "Standard Deviation"}
var ipTypes = []string{"IPv4", "IPv6"}
var engines = []string{"Amazon", "Bing", "DuckDuckGo", "eBay", "Google", "Reddit", "Twitter", "Yahoo!", "YouTube"}
var webpageDetails = []string{"Page Contents", "Page Selection", "Page URL", "Name"}
var urlComponents = []string{"Scheme", "User", "Password", "Host", "Port", "Path", "Query", "Fragment"}
var httpMethods = []string{"POST", "PUT", "PATCH", "DELETE"}
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
		enum:      httpMethods,
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
var imageEditingPositions = []string{"Center", "Top Left", "Top Right", "Bottom Left", "Bottom Right", "Custom"}
var sortOrders = []string{"asc", "desc"}
var windowSortings = []string{"Title", "App Name", "Width", "Height", "X Position", "Y Position", "Window Index", "Name", "Random"}
var windowPositions = []string{"Top Left", "Top Center", "Top Right", "Middle Left", "Center", "Middle Right", "Bottom Left", "Bottom Center", "Bottom Right", "Coordinates"}
var windowConfigurations = []string{"Fit Screen", "Top Half", "Bottom Half", "Left Half", "Right Half", "Top Left Quarter", "Top Right Quarter", "Bottom Left Quarter", "Bottom Right Quarter", "Dimensions"}
var focusModes = map[string]any{
	"Identifier":    "com.apple.donotdisturb.mode.default",
	"DisplayString": "Do Not Disturb",
}
var editEventDetails = []string{"Start Date", "End Date", "Is All Day", "Location", "Duration", "My Status", "Attendees", "URL", "Title", "Notes", "Attachments"}
var eventDetails = []string{"Start Date", "End Date", "Is All Day", "Calendar", "Location", "Has Alarms", "Duration", "Is Canceled", "My Status", "Organizer", "Organizer Is Me", "Attendees", "Number of Attendees", "URL", "Title", "Notes", "Attachments", "File Size", "File Extension", "Creation Date", "File Path", "Last Modified Date", "Name"}
var dateFormats = []string{"None", "Short", "Medium", "Long", "Relative", "RFC 2822", "ISO 8601", "Custom"}
var timeFormats = []string{"None", "Short", "Medium", "Long", "Relative"}
var timerDurations = []string{"hr", "min", "sec"}
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
var fileLabels = []string{"red", "orange", "yellow", "green", "blue", "purple", "gray"}
var filesSortBy = []string{"File Size", "File Extension", "Creation Date", "File Path", "Last Modified Date", "Name", "Random"}
var pdfMergeBehaviors = []string{"Append", "Shuffle"}
var cameras = []string{"Front", "Back"}
var cameraQualities = []string{"Low", "Medium", "High"}
var backgroundSounds = []string{"BalancedNoise", "BrightNoise", "DarkNoise", "Ocean", "Rain", "Stream"}
var soundRecognitionOperations = []string{"pause", "activate", "toggle"}
var textSizes = []string{
	"Accessibility Extra Extra Extra Large",
	"Accessibility Extra Extra Large",
	"Accessibility Extra Large",
	"Accessibility Large",
	"Accessibility Medium",
	"Extra Extra Extra Large",
	"Extra Extra Large",
	"Extra Large",
	"Default",
	"Medium",
	"Small",
	"Extra Small",
}
var roundings = []string{"Ones Place", "Tens Place", "Hundreds Place", "Thousands", "Ten Thousands", "Hundred Thousands", "Millions"}
var imageMaskTypes = []string{"Rounded Rectangle", "Ellipse", "Icon"}
var imageDetails = []string{"Album", "Width", "Height", "Date Taken", "Media Type", "Photo Type", "Is a Screenshot", "Is a Screen Recording", "Location", "Duration", "Frame Rate", "Orientation", "Camera Make", "Camera Model", "Metadata Dictionary", "Is Favorite", "File Size", "File Extension", "Creation Date", "File Path", "Last Modified Date", "Name"}
var colorSpaces = []string{"RGB", "Gray"}
var shuffleOptions = []string{"Off", "Songs"}
var repeatOptions = []string{"None", "One", "All"}
var musicDetails = []string{"Title", "Album", "Artist", "Album Artist", "Genre", "Composer", "Date Added", "Media Kind", "Duration", "Play Count", "Track Number", "Disc Number", "Album Artwork", "Is Explicit", "Lyrics", "Release Date", "Comments", "Is Cloud Item", "Skip Count", "Last Played Date", "Rating", "File Path", "Name"}
var wifiNetworkDetails = []string{"Network Name", "BSSID", "Wi-Fi Standard", "RX Rate", "TX Rate", "RSSI", "Noise", "Channel Number", "Hardware MAC Address"}
var cellularNetworkDetails = []string{"Carrier Name", "Radio Technology", "Country Code", "Is Roaming Abroad", "Number of Signal Bars"}

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
	"vibrate":            {mac: false},
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
				enum:         timerDurations,
			},
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
				enum:      contactDetails,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"Mode": "Remove",
			}
		},
	},
	"speak": {
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
	},
	"listen": {
		identifier: "dictatetext",
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
				enum:      languages,
			},
		},
	},
	"prependToFile": {
		identifier: "file.append",
		parameters: []parameterDefinition{
			{
				name:      "filePath",
				validType: String,
				key:       "WFFilePath",
			},
			{
				name:      "text",
				validType: String,
				key:       "WFInput",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFAppendFileWriteMode": "Prepend",
			}
		},
	},
	"appendToFile": {
		identifier:    "file.append",
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "filePath",
				validType: String,
				key:       "WFFilePath",
			},
			{
				name:      "text",
				validType: String,
				key:       "WFInput",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFAppendFileWriteMode": "Append",
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
				enum:      fileLabels,
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
				enum:      filesSortBy,
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
	"optimizePDF": {
		identifier: "compresspdf",
		parameters: []parameterDefinition{
			{
				name:      "pdfFile",
				validType: Variable,
				key:       "WFInput",
			},
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
	"makePDF": {
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "includeMargin",
				validType:    Bool,
				key:          "WFPDFIncludeMargin",
				defaultValue: false,
				optional:     true,
			},
			{
				name:         "mergeBehavior",
				validType:    String,
				key:          "WFPDFDocumentMergeBehavior",
				defaultValue: "Append",
				enum:         pdfMergeBehaviors,
				optional:     true,
			},
		},
	},
	"makeSpokenAudio": {
		identifier: "makespokenaudiofromtext",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
				key:       "WFInput",
			},
			{
				name:      "rate",
				validType: Integer,
				key:       "WFSpeakTextRate",
				optional:  true,
			},
			{
				name:      "pitch",
				validType: Integer,
				key:       "WFSpeakTextPitch",
				optional:  true,
			},
		},
	},
	"createFolder": { // TODO: Writing to locations other than the Shortcuts folder.
		identifier: "file.createfolder",
		parameters: []parameterDefinition{
			{
				name:      "path",
				validType: String,
				key:       "WFFilePath",
			},
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
	"matchText": {
		identifier:    "text.match",
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "regexPattern",
				validType: String,
				key:       "WFMatchTextPattern",
			},
			{
				name:      "text",
				validType: String,
				key:       "text",
			},
			{
				name:         "caseSensitive",
				validType:    Bool,
				key:          "WFMatchTextCaseSensitive",
				defaultValue: true,
				optional:     true,
			},
		},
	},
	"getMatchGroups": {
		identifier: "text.match.getgroup",
		parameters: []parameterDefinition{
			{
				name:      "matches",
				validType: Variable,
				key:       "matches",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFGetGroupType": "All Groups",
			}
		},
	},
	"getMatchGroup": {
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFGetGroupType": "Group At Index",
			}
		},
		defaultAction: true,
	},
	"getFileFromFolder": {
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
	},
	"getFile": {
		identifier:    "documentpicker.open",
		defaultAction: true,
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
	},
	"markup": {
		identifier: "avairyeditphoto",
		parameters: []parameterDefinition{
			{
				name:      "document",
				validType: Variable,
				key:       "WFDocument",
			},
		},
	},
	"rename": {
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
	},
	"reveal": {
		identifier: "file.reveal",
		parameters: []parameterDefinition{
			{
				name:      "files",
				validType: Variable,
				key:       "WFFile",
			},
		},
	},
	"define": {
		identifier: "showdefinition",
		parameters: []parameterDefinition{
			{
				name:      "word",
				validType: String,
				key:       "Word",
			},
		},
	},
	"makeQRCode": {
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
	},
	"openNote": {
		identifier: "shownote",
		parameters: []parameterDefinition{
			{
				name:      "note",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"splitPDF": {
		parameters: []parameterDefinition{
			{
				name:      "pdf",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"makeHTML": {
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
	},
	"makeMarkdown": {
		identifier: "getmarkdownfromrichtext",
		parameters: []parameterDefinition{
			{
				name:      "richText",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"getRichTextFromHTML": {
		parameters: []parameterDefinition{
			{
				name:      "html",
				validType: Variable,
				key:       "WFHTML",
			},
		},
	},
	"getRichTextFromMarkdown": {
		parameters: []parameterDefinition{
			{
				name:      "markdown",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"print": {
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"selectFile": {
		identifier:    "file.select",
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:         "selectMultiple",
				validType:    Bool,
				key:          "SelectMultiple",
				defaultValue: false,
				optional:     true,
			},
		},
	},
	"selectFolder": {
		identifier: "file.select",
		parameters: []parameterDefinition{
			{
				name:         "selectMultiple",
				validType:    Bool,
				key:          "SelectMultiple",
				defaultValue: false,
				optional:     true,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"WFPickingMode": "Folders",
			}
		},
	},
	"getFileLink": {
		identifier: "file.getlink",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFFile",
			},
		},
	},
	"getParentDirectory": {
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"getEmojiName": {
		identifier: "getnameofemoji",
		parameters: []parameterDefinition{
			{
				name:      "emoji",
				validType: String,
				key:       "WFInput",
			},
		},
	},
	"getFileDetail": {
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
	},
	"deleteFiles": {
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
	},
	"getTextFromImage": {
		identifier: "extracttextfromimage",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFImage",
			},
		},
	},
	"connectToServer": {
		identifier: "connecttoservers",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFInput",
			},
		},
	},
	"appendNote": {
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
	},
	"addToBooks": {
		appIdentifier: "com.apple.iBooksX",
		identifier:    "openin",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "BooksInput",
			},
		},
	},
	"saveFile": {
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFAskWhereToSave": false,
			}
		},
	},
	"saveFilePrompt": {
		identifier:    "documentpicker.save",
		defaultAction: true,
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
	},
	"getSelectedFiles": {identifier: "finder.getselectedfiles"},
	"extractArchive": {
		identifier: "unzip",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFArchive",
			},
		},
	},
	"makeArchive": {
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
	},
	"quicklook": {
		identifier: "previewdocument",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"translateFrom": {
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
				enum:      languages,
			},
			{
				name:      "to",
				validType: String,
				key:       "WFSelectedLanguage",
				enum:      languages,
			},
		},
		outputType: String,
	},
	"translate": {
		identifier:    "text.translate",
		defaultAction: true,
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
				optional:  true,
				enum:      languages,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFSelectedFromLanguage": "Detect Language",
			}
		},
		outputType: String,
	},
	"detectLanguage": {
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
				key:       "WFInput",
			},
		},
		outputType: String,
	},
	"replaceText": {
		identifier: "text.replace",
		parameters: []parameterDefinition{
			{
				name:      "find",
				key:       "WFReplaceTextFind",
				validType: String,
			},
			{
				name:      "replacement",
				key:       "WFReplaceTextReplace",
				validType: String,
			},
			{
				name:      "subject",
				key:       "WFInput",
				validType: String,
			},
			{
				name:         "caseSensitive",
				key:          "WFReplaceTextCaseSensitive",
				validType:    Bool,
				defaultValue: true,
				optional:     true,
			},
			{
				name:         "regExp",
				key:          "WFReplaceTextRegularExpression",
				validType:    Bool,
				defaultValue: false,
				optional:     true,
			},
		},
		outputType: String,
	},
	"uppercase": {
		identifier: "text.changecase",
		parameters: []parameterDefinition{
			{
				name:      "text",
				key:       "text",
				validType: String,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return changeCase("UPPERCASE")
		},
		defaultAction: true,
	},
	"lowercase": {
		identifier: "text.changecase",
		parameters: []parameterDefinition{
			{
				name:      "text",
				key:       "text",
				validType: String,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return changeCase("lowercase")
		},
	},
	"titleCase": {
		identifier: "text.changecase",
		parameters: []parameterDefinition{
			{
				name:      "text",
				key:       "text",
				validType: String,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return changeCase("Capitalize with Title Case")
		},
		outputType: String,
	},
	"capitalize": {
		identifier: "text.changecase",
		parameters: []parameterDefinition{
			{
				name:      "text",
				key:       "text",
				validType: String,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return changeCase("Capitalize with sentence case")
		},
		outputType: String,
	},
	"capitalizeAll": {
		identifier: "text.changecase",
		parameters: []parameterDefinition{
			{
				name:      "text",
				key:       "text",
				validType: String,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return changeCase("Capitalize Every Word")
		},
		outputType: String,
	},
	"alternateCase": {
		identifier: "text.changecase",
		parameters: []parameterDefinition{
			{
				name:      "text",
				key:       "text",
				validType: String,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return changeCase("cApItAlIzE wItH aLtErNaTiNg cAsE")
		},
		outputType: String,
	},
	"correctSpelling": {
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
				key:       "text",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"Show-text": true,
			}
		},
		outputType: String,
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
	"makeDiskImage": {
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"SizeToFit": true,
			}
		},
		mac:        true,
		minVersion: 15,
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
				enum: storageUnits,
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
	"openFile": {
		identifier: "openin",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "askWhenRun",
				validType:    Bool,
				key:          "WFOpenInAskWhenRun",
				defaultValue: false,
				optional:     true,
			},
		},
	},
	"transcribeText": {
		appIdentifier: "com.apple.ShortcutsActions",
		identifier:    "TranscribeAudioAction",
		parameters: []parameterDefinition{
			{
				name:      "audioFile",
				validType: Variable,
				key:       "audioFile",
			},
		},
		minVersion: 17,
	},
	"getCurrentLocation": {
		identifier: "location",
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFLocation": map[string]bool{
					"isCurrentLocation": true,
				},
			}
		},
	},
	"getAddresses": {
		identifier: "detect.address",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"getCurrentWeather": {
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
	},
	"openInMaps": {
		identifier: "searchmaps",
		parameters: []parameterDefinition{
			{
				name:      "location",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"streetAddress": {
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
	},
	"getWeatherDetail": {
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
	},
	"getWeatherForecast": {
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
	},
	"getLocationDetail": {
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
	},
	"getMapsLink": {
		identifier: "",
		parameters: []parameterDefinition{
			{
				name:      "location",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"getHalfwayPoint": {
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
	},
	"removeBackground": {
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
	},
	"clearUpNext":    {},
	"getCurrentSong": {},
	"getLastImport":  {identifier: "getlatestphotoimport"},
	"getLatestBursts": {
		parameters: []parameterDefinition{
			{
				name:      "count",
				validType: Integer,
				key:       "WFGetLatestPhotoCount",
			},
		},
	},
	"getLatestLivePhotos": {
		identifier: "getlatestlivephotos",
		parameters: []parameterDefinition{
			{
				name:      "count",
				validType: Integer,
				key:       "WFGetLatestPhotoCount",
			},
		},
	},
	"getLatestScreenshots": {
		identifier: "getlastscreenshot",
		parameters: []parameterDefinition{
			{
				name:      "count",
				validType: Integer,
				key:       "WFGetLatestPhotoCount",
			},
		},
	},
	"getLatestVideos": {
		identifier: "getlastvideo",
		parameters: []parameterDefinition{
			{
				name:      "count",
				validType: Integer,
				key:       "WFGetLatestPhotoCount",
			},
		},
	},
	"getLatestPhotos": {
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
	},
	"maskImage": {
		identifier: "image.mask",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "type",
				validType: String,
				enum:      imageMaskTypes,
				key:       "WFMaskType",
			},
			{
				name:      "radius",
				validType: String,
				key:       "WFMaskCornerRadius",
				optional:  true,
			},
		},
	},
	"customImageMask": {
		identifier: "image.mask",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "customMaskImage",
				validType: Variable,
				key:       "WFCustomMaskImage",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFMaskType": "Custom Image",
			}
		},
	},
	"getImageDetail": {
		identifier: "properties.images",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "detail",
				validType: String,
				enum:      imageDetails,
				key:       "WFContentItemPropertyName",
			},
		},
	},
	"extractImageText": {
		identifier: "extracttextfromimage",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFImage",
			},
		},
	},
	"makeImageFromPDFPage": {
		parameters: []parameterDefinition{
			{
				name:      "pdf",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "colorSpace",
				validType:    String,
				enum:         colorSpaces,
				key:          "WFMakeImageFromPDFPageColorspace",
				optional:     true,
				defaultValue: "RGB",
			},
			{
				name:         "pageResolution",
				validType:    String,
				key:          "WFMakeImageFromPDFPageResolution",
				optional:     true,
				defaultValue: "300",
			},
		},
	},
	"makeImageFromRichText": {
		parameters: []parameterDefinition{
			{
				name:      "pdf",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "width",
				validType: String,
				key:       "WFWidth",
			},
			{
				name:      "height",
				validType: String,
				key:       "WFHeight",
			},
		},
	},
	"getImages": {
		identifier: "detect.images",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	},
	"overlayImage": {
		identifier: "overlayimageonimage",
		parameters: []parameterDefinition{
			{
				name:      "image",
				key:       "WFInput",
				validType: Variable,
			},
			{
				name:      "overlayImage",
				key:       "WFImage",
				validType: Variable,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFShouldShowImageEditor": true,
			}
		},
	},
	"customImageOverlay": {
		identifier: "overlayimageonimage",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "overlayImage",
				validType: Variable,
				key:       "WFImage",
			},
			{
				name:      "width",
				validType: String,
				key:       "WFImageWidth",
				optional:  true,
			},
			{
				name:      "height",
				validType: String,
				key:       "WFImageHeight",
				optional:  true,
			},
			{
				name:         "rotation",
				validType:    String,
				key:          "WFRotation",
				optional:     true,
				defaultValue: "0",
			},
			{
				name:         "opacity",
				validType:    String,
				key:          "WFOverlayImageOpacity",
				defaultValue: "100",
				optional:     true,
			},
			{
				name:         "position",
				validType:    String,
				key:          "WFImagePosition",
				enum:         imageEditingPositions,
				defaultValue: "Center",
				optional:     true,
			},
			{
				name:      "customPositionX",
				validType: String,
				key:       "WFImageX",
				optional:  true,
			},
			{
				name:      "customPositionY",
				validType: String,
				key:       "WFImageY",
				optional:  true,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFShouldShowImageEditor": false,
			}
		},
	},
	"takePhoto": {
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
		check: func(args []actionArgument, _ *actionDefinition) {
			if len(args) == 0 {
				return
			}

			var photos = getArgValue(args[0])
			if photos == "0" {
				parserError("Number of photos to take must be greater than zero.")
			}
		},
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
				enum:      timerDurations,
			},
			{
				name:         "behavior",
				key:          "WFSeekBehavior",
				validType:    String,
				defaultValue: "To Time",
				enum:         []string{"To Time", "Forward By", "Backward By"},
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
	"trimVideo": {
		parameters: []parameterDefinition{
			{
				name:      "video",
				validType: Variable,
				key:       "WFInputMedia",
			},
		},
	},
	"takeVideo": {
		parameters: []parameterDefinition{
			{
				name:         "camera",
				validType:    String,
				key:          "WFCameraCaptureDevice",
				defaultValue: "Front",
				enum:         cameras,
			},
			{
				name:         "quality",
				validType:    String,
				key:          "WFCameraCaptureQuality",
				defaultValue: "High",
				enum:         cameraQualities,
			},
			{
				name:         "recordingStart",
				validType:    String,
				key:          "WFRecordingStart",
				defaultValue: "Immediately",
				enum:         recordingStarts,
			},
		},
	},
	"setVolume": {
		parameters: []parameterDefinition{
			{
				name:      "volume",
				validType: Float,
				key:       "WFVolume",
			},
		},
	},
	"getMusicDetail": {
		identifier: "properties.music",
		parameters: []parameterDefinition{
			{
				name:      "music",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "detail",
				validType: String,
				key:       "WFContentItemPropertyName",
				enum:      musicDetails,
			},
		},
	},
	"addToMusic": {
		identifier: "addtoplaylist",
		parameters: []parameterDefinition{
			{
				name:      "songs",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"addToPlaylist": {
		defaultAction: true,
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
	},
	"playNext": {
		identifier:    "addmusictoupnext",
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "music",
				validType: Variable,
				key:       "WFMusic",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFWhenToPlay": "Next",
			}
		},
	},
	"playLater": {
		identifier: "addmusictoupnext",
		parameters: []parameterDefinition{
			{
				name:      "music",
				validType: Variable,
				key:       "WFMusic",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFWhenToPlay": "Later",
			}
		},
	},
	"createPlaylist": {
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
	},
	"addToGIF": {
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
	},
	"getImageFrames": {
		identifier: "getframesfromimage",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFImage",
			},
		},
	},
	"getPlaylistSongs": {
		identifier: "get.playlist",
		parameters: []parameterDefinition{
			{
				name:      "playlistName",
				validType: Variable,
				key:       "WFPlaylistName",
			},
		},
	},
	"playMusic": {
		parameters: []parameterDefinition{
			{
				name: "music",
				key:  "WFMediaItems",
			},
			{
				name:     "shuffle",
				key:      "WFPlayMusicActionShuffle",
				enum:     shuffleOptions,
				optional: true,
			},
			{
				name:     "repeat",
				key:      "WFPlayMusicActionRepeat",
				enum:     repeatOptions,
				optional: true,
			},
		},
	},
	"makeGIF": {
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "delay",
				validType:    String,
				key:          "WFMakeGIFActionDelayTime",
				defaultValue: "0.3",
				optional:     true,
			},
			{
				name:      "loops",
				validType: Integer,
				key:       "WFMakeGIFActionLoopCount",
				optional:  true,
			},
			{
				name:      "width",
				validType: String,
				key:       "WFMakeGIFActionManualSizeWidth",
				optional:  true,
			},
			{
				name:      "height",
				validType: String,
				key:       "WFMakeGIFActionManualSizeHeight",
				optional:  true,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			var params = make(map[string]any)
			var argsLen = len(args)
			params["WFMakeGIFActionLoopEnabled"] = !(argsLen > 2)
			params["WFMakeGIFActionAutoSize"] = !(argsLen > 3)

			return params
		},
	},
	"convertToJPEG": {
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFImageFormat": "JPEG",
			}
		},
	},
	"combineImages": {
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
		check: func(args []actionArgument, _ *actionDefinition) {
			if len(args) < 2 {
				return
			}
			var combineMode = fmt.Sprintf("%s", args[1].value)
			if strings.ToLower(combineMode) == "grid" {
				args[1].value = "In a Grid"
			}
		},
	},
	"rotateMedia": {
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
	},
	"selectPhotos": {
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
	},
	"createAlbum": {
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
	},
	"cropImage": {
		identifier: "image.crop",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
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
				optional:     true,
				defaultValue: "100",
			},
			{
				name:         "position",
				validType:    String,
				key:          "WFImageCropPosition",
				enum:         imageEditingPositions,
				optional:     true,
				defaultValue: "Center",
			},
			{
				name:      "customPositionX",
				validType: String,
				key:       "WFImageCropX",
				optional:  true,
			},
			{
				name:      "customPositionY",
				validType: String,
				key:       "WFImageCropY",
				optional:  true,
			},
		},
	},
	"deletePhotos": {
		parameters: []parameterDefinition{
			{
				name:      "photos",
				validType: Variable,
				key:       "photos",
			},
		},
	},
	"removeFromAlbum": {
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
	},
	"selectMusic": {
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
	},
	"skipBack": {
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFSkipBackBehavior": "Previous Song",
			}
		},
	},
	"skipFwd": {
		identifier: "skipforward",
	},
	"searchAppStore": {
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "WFSearchTerm",
			},
		},
	},
	"searchPodcasts": {
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "WFSearchTerm",
			},
		},
	},
	"getPodcasts": {
		identifier: "getpodcastsfromlibrary",
	},
	"playPodcast": {
		parameters: []parameterDefinition{
			{
				name:      "podcast",
				validType: Variable,
				key:       "WFPodcastShow",
			},
		},
	},
	"getPodcastDetail": {
		identifier: "properties.podcastshow",
		parameters: []parameterDefinition{
			{
				name:      "podcast",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "detail",
				validType: String,
				key:       "WFContentItemPropertyName",
				enum:      []string{"Feed URL", "Genre", "Episode Count", "Artist", "Store ID", "Store URL", "Artwork", "Artwork URL", "Name"},
			},
		},
	},
	"makeVideoFromGIF": {
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
	},
	"flipImage": {
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
	},
	"recordAudio": {
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
	},
	"resizeImage": {
		identifier: "image.resize",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFImage",
			},
			{
				name:      "width",
				key:       "WFImageResizeWidth",
				validType: String,
			},
			{
				name:      "height",
				key:       "WFImageResizeHeight",
				validType: String,
				optional:  true,
			},
		},
	},
	"resizeImageByPercent": {
		identifier: "image.resize",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFImage",
			},
			{
				name:      "percentage",
				validType: String,
				key:       "WFImageResizePercentage",
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"WFImageResizeKey": "Percentage",
			}
		},
	},
	"resizeImageByLongestEdge": {
		identifier: "image.resize",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFImage",
			},
			{
				name:      "length",
				validType: String,
				key:       "WFImageResizeLength",
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"WFImageResizeKey": "Longest Edge",
			}
		},
	},
	"convertImage": {
		identifier:    "image.convert",
		defaultAction: true,
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
				name:      "quality",
				key:       "WFImageCompressionQuality",
				validType: Float,
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
	},
	"encodeVideo": {
		identifier:    "encodemedia",
		defaultAction: true,
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
	},
	"encodeAudio": {
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFMediaAudioOnly": true,
			}
		},
	},
	"setMetadata": {
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"Metadata": true,
			}
		},
	},
	"stripMediaMetadata": {
		identifier: "encodemedia",
		parameters: []parameterDefinition{
			{
				name:      "media",
				validType: Variable,
				key:       "WFMedia",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"Metadata": true,
			}
		},
	},
	"stripImageMetadata": {
		identifier: "image.convert",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFImagePreserveMetadata": false,
				"WFImageFormat":           "Match Input",
			}
		},
	},
	"savePhoto": {
		identifier: "savetocameraroll",
		parameters: []parameterDefinition{
			{
				name:      "image",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:         "album",
				validType:    String,
				key:          "WFCameraRollSelectedGroup",
				defaultValue: "Recents",
				optional:     true,
			},
		},
	},
	"play": {
		identifier: "pausemusic",
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFPlayPauseBehavior": "Play",
			}
		},
	},
	"pause": {
		identifier: "pausemusic",
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFPlayPauseBehavior": "Pause",
			}
		},
	},
	"togglePlayPause": {
		identifier:    "pausemusic",
		defaultAction: true,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFPlayPauseBehavior": "Play/Pause",
			}
		},
	},
	"startShazam": {
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
	},
	"showIniTunes": {
		identifier: "showinstore",
		parameters: []parameterDefinition{
			{
				name:      "product",
				validType: Variable,
				key:       "WFProduct",
			},
		},
	},
	"takeScreenshot": {
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:         "mainMonitorOnly",
				validType:    Bool,
				key:          "WFTakeScreenshotMainMonitorOnly",
				defaultValue: false,
				optional:     true,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFTakeScreenshotScreenshotType": "Full Screen",
			}
		},
		mac: true,
	},
	"takeInteractiveScreenshot": {
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFTakeScreenshotScreenshotType": "Interactive",
			}
		},
		mac: true,
	},
	"shutdown": {
		defaultAction: true,
		identifier:    "reboot",
		minVersion:    17,
	},
	"reboot": {
		minVersion: 17,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFShutdownMode": "Restart",
			}
		},
	},
	"sleep": {
		minVersion: 17,
		mac:        true,
	},
	"displaySleep": {
		minVersion: 17,
		mac:        true,
	},
	"logout": {
		minVersion: 17,
		mac:        true,
	},
	"lockScreen": {
		minVersion: 17,
	},
	"getOrientation": {
		identifier:    "GetOrientationAction",
		appIdentifier: "com.apple.ShortcutsActions",
		outputType:    String,
	},
	"number": {
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: Variable,
				key:       "WFNumberActionNumber",
			},
		},
		outputType: Integer,
	},
	"text": {
		identifier: "gettext",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
				key:       "WFTextActionText",
			},
		},
		outputType: String,
	},
	"getObjectOfClass": {
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
	},
	"getOnScreenContent": {},
	"fileSize": {
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFFileSizeIncludeUnits": false,
			}
		},
	},
	"getDeviceDetail": {
		identifier: "getdevicedetails",
		parameters: []parameterDefinition{
			{
				name:      "detail",
				key:       "WFDeviceDetail",
				validType: String,
				enum:      deviceDetails,
			},
		},
	},
	"setBrightness": {
		parameters: []parameterDefinition{
			{
				name:      "brightness",
				validType: Float,
				key:       "WFBrightness",
			},
		},
	},
	"getName": {
		identifier: "getitemname",
		parameters: []parameterDefinition{
			{
				name:      "item",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"setName": {
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
	},
	"count": {
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Input",
				validType: Variable,
			},
			{
				name:         "type",
				key:          "WFCountType",
				validType:    String,
				enum:         []string{"Items", "Characters", "Words", "Sentences", "Lines"},
				defaultValue: "Items",
			},
		},
		defaultAction: true,
		addParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"WFCountType": "Items",
			}
		},
		outputType: Integer,
	},
	"lightMode": {
		identifier: "appearance",
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"operation": "set",
				"style":     "light",
			}
		},
	},
	"darkMode": {
		identifier:    "appearance",
		defaultAction: true,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"operation": "set",
				"style":     "dark",
			}
		},
	},
	"getBatteryLevel": {defaultAction: true},
	"isCharging": {
		identifier: "getbatterylevel",
		minVersion: 16.2,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"Subject": "Is Charging",
			}
		},
		outputType: Bool,
	},
	"connectedToCharger": {
		identifier: "getbatterylevel",
		minVersion: 16.2,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"Subject": "Is Connected to Charger",
			}
		},
		outputType: Bool,
	},
	"getShortcuts": {
		identifier: "getmyworkflows",
		outputType: Arr,
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
	"hash": {
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
		outputType: String,
	},
	"formatNumber": {
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
		outputType: Integer,
	},
	"randomNumber": {
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
		outputType: Integer,
	},
	"base64Encode": {
		identifier: "base64encode",
		parameters: []parameterDefinition{
			{
				name:      "encodeInput",
				key:       "WFInput",
				validType: Variable,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"input": "Encode",
			}
		},
		defaultAction: true,
		outputType:    String,
	},
	"base64Decode": {
		identifier: "base64encode",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFEncodeMode": "Decode",
			}
		},
		outputType: String,
	},
	"show": {
		identifier: "showresult",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Text",
				validType: String,
			},
		},
	},
	"waitToReturn": {},
	"showNotification": {
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
			{
				name:      "attachment",
				key:       "WFInput",
				validType: Variable,
				optional:  true,
			},
		},
	},
	"stop": {
		identifier: "exit",
	},
	"comment": {
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: RawString,
				key:       "WFCommentActionText",
			},
		},
	},
	"nothing": {},
	"wait": {
		identifier: "delay",
		parameters: []parameterDefinition{
			{
				name:      "seconds",
				key:       "WFDelayTime",
				validType: Integer,
			},
		},
	},
	"alert": {
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFAlertActionCancelButtonShown": false,
			}
		},
	},
	"confirm": {
		identifier:    "alert",
		defaultAction: true,
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFAlertActionCancelButtonShown": true,
			}
		},
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
				enum:         inputTypes,
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
	"chooseFromList": {
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
	},
	"typeOf": {
		identifier: "getitemtype",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
		outputType: String,
	},
	"getKeys": {
		identifier: "getvalueforkey",
		parameters: []parameterDefinition{
			{
				name:      "dictionary",
				validType: Dict,
				key:       "WFInput",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFGetDictionaryValueType": "All Keys",
			}
		},
		outputType: Arr,
	},
	"getValues": {
		identifier: "getvalueforkey",
		parameters: []parameterDefinition{
			{
				name:      "dictionary",
				validType: Dict,
				key:       "WFInput",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFGetDictionaryValueType": "All Values",
			}
		},
		outputType: Arr,
	},
	"getValue": {
		identifier:    "getvalueforkey",
		defaultAction: true,
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFGetDictionaryValueType": "Value",
			}
		},
	},
	"setValue": {
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
				enum:         appSplitRatios,
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
	"calculate": {
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFMathOperation": "...",
			}
		},
		outputType: Integer,
	},
	"statistic": {
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
	},
	"dismissSiri": {},
	"getFirstItem": {
		identifier:    "getitemfromlist",
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "list",
				key:       "WFInput",
				validType: Variable,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFItemSpecifier": "First Item",
			}
		},
	},
	"getLastItem": {
		identifier: "getitemfromlist",
		parameters: []parameterDefinition{
			{
				name:      "list",
				key:       "WFInput",
				validType: Variable,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"WFItemSpecifier": "Last Item",
			}
		},
	},
	"getRandomItem": {
		identifier: "getitemfromlist",
		parameters: []parameterDefinition{
			{
				name:      "list",
				key:       "WFInput",
				validType: Variable,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFItemSpecifier": "Random Item",
			}
		},
	},
	"getListItem": {
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFItemSpecifier": "Item At Index",
			}
		},
	},
	"getListItems": {
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFItemSpecifier": "Items in Range",
			}
		},
		outputType: Arr,
	},
	"getNumbers": {
		identifier: "detect.number",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
		outputType: Integer,
	},
	"getDictionary": {
		identifier: "detect.dictionary",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
		outputType: Dict,
	},
	"getText": {
		identifier: "detect.text",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
		outputType: String,
	},
	"getAllWallpapers": {
		identifier:    "posters.get",
		minVersion:    16.2,
		mac:           false,
		defaultAction: true,
		outputType:    Arr,
	},
	"getWallpaper": {
		identifier: "posters.get",
		minVersion: 16.2,
		mac:        false,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFPosterType": "Current",
			}
		},
	},
	"setWallpaper": {
		identifier: "wallpaper.set",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	},
	"startScreensaver": {mac: true},
	"contentGraph": {
		identifier: "viewresult",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	},
	"openXCallbackURL": {
		identifier:    "openxcallbackurl",
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "url",
				key:       "WFXCallbackURL",
				validType: String,
				infinite:  true,
			},
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
	"output": {
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "output",
				validType: String,
				key:       "WFOutput",
			},
		},
	},
	"mustOutput": {
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFNoOutputSurfaceBehavior": "Respond",
			}
		},
	},
	"outputOrClipboard": {
		identifier: "output",
		parameters: []parameterDefinition{
			{
				name:      "output",
				validType: String,
				key:       "WFOutput",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFNoOutputSurfaceBehavior": "Copy to Clipboard",
			}
		},
	},
	"DNDOn": {
		identifier: "dnd.set",
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"FocusModes": focusModes,
				"Enabled":    1,
			}
		},
	},
	"DNDOff": {
		identifier:    "dnd.set",
		defaultAction: true,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"FocusModes": focusModes,
				"Enabled":    0,
			}
		},
	},
	"setBackgroundSound": {
		appIdentifier: "com.apple.AccessibilityUtilities",
		identifier:    "AXSettingsShortcuts.AXSetBackgroundSoundIntent",
		parameters: []parameterDefinition{
			{
				name:         "sound",
				validType:    String,
				key:          "backgroundSound",
				defaultValue: "Balanced Noise",
				enum:         backgroundSounds,
				optional:     true,
			},
		},
	},
	"setBackgroundSoundsVolume": {
		appIdentifier: "com.apple.AccessibilityUtilities",
		identifier:    "AXSettingsShortcuts.AXSetBackgroundSoundVolumeIntent",
		parameters: []parameterDefinition{
			{
				name:      "volume",
				validType: Float,
				key:       "volumeValue",
			},
		},
	},
	"playSound": {
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	},
	"round": {
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: Integer,
				key:       "WFInput",
			},
			{
				name:         "roundTo",
				validType:    String,
				key:          "WFRoundTo",
				enum:         roundings,
				optional:     true,
				defaultValue: "Ones Place",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFRoundMode": "Normal",
			}
		},
	},
	"ceil": {
		identifier: "round",
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: Integer,
				key:       "WFInput",
			},
			{
				name:         "roundTo",
				validType:    String,
				key:          "WFRoundTo",
				enum:         roundings,
				optional:     true,
				defaultValue: "Ones Place",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFRoundMode": "Always Round Up",
			}
		},
	},
	"floor": {
		identifier: "round",
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: Integer,
				key:       "WFInput",
			},
			{
				name:         "roundTo",
				validType:    String,
				key:          "WFRoundTo",
				enum:         roundings,
				optional:     true,
				defaultValue: "Ones Place",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFRoundMode": "Always Round Down",
			}
		},
	},
	"runShellScript": {
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
	},
	"makeShortcut": {
		appIdentifier: "com.apple.shortcuts",
		identifier:    "CreateWorkflowAction",
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
	"searchShortcuts": {
		appIdentifier: "com.apple.shortcuts",
		identifier:    "SearchShortcutsAction",
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "searchPhrase",
			},
		},
		minVersion: 16.4,
	},
	"searchPasswords": {
		identifier: "openpasswords",
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "WFShowPasswordsSearchTerm",
			},
		},
	},
	"convertToUSDZ": {
		appIdentifier: "com.apple.HydraUSDAppIntents",
		identifier:    "ConvertToUSDZ",
		parameters: []parameterDefinition{
			{
				name:      "file",
				key:       "file",
				validType: Variable,
			},
		},
	},
	"runSSHScript": {
		parameters: []parameterDefinition{
			{
				name:      "script",
				key:       "WFSSHScript",
				validType: String,
			},
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
			{
				name:      "host",
				key:       "WFSSHHost",
				validType: String,
			},
			{
				name:      "port",
				key:       "WFSSHPort",
				validType: String,
			},
			{
				name:      "user",
				key:       "WFSSHUser",
				validType: String,
			},
			{
				name: "authType",
				key:  "WFSSHAuthenticationType",
				enum: []string{"Password", "SSH Key"},
			},
			{
				name:      "password",
				key:       "WFSSHPassword",
				validType: String,
			},
		},
	},
	"getWifiDetail": {
		identifier: "getwifi",
		parameters: []parameterDefinition{
			{
				name:      "detail",
				validType: String,
				key:       "WFWiFiDetail",
				enum:      wifiNetworkDetails,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"WFNetworkDetailsNetwork": "Wi-Fi",
			}
		},
	},
	"getCellularDetail": {
		identifier: "getwifi",
		mac:        false,
		parameters: []parameterDefinition{
			{
				name:      "detail",
				validType: String,
				key:       "WFCellularDetail",
				enum:      cellularNetworkDetails,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"WFNetworkDetailsNetwork": "Cellular",
			}
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
	"runAppleScript": {
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
	},
	"runJSAutomation": {
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
	},
	"getWindows": {
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
	"moveWindow": {
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
	},
	"resizeWindow": {
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
				enum:      measurementUnitTypes,
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

			makeMeasurementUnits()

			var unitType = value.(string)
			checkEnum(&parameterDefinition{
				name: "measurement unit",
				enum: units[unitType],
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
				enum:      measurementUnitTypes,
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

			makeMeasurementUnits()

			var unitType = value.(string)
			checkEnum(&parameterDefinition{
				name: "unit",
				enum: units[unitType],
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
				enum:      contactDetails,
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
	"getShazamDetail": {
		identifier: "properties.shazam",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "detail",
				validType: String,
				key:       "WFContentItemPropertyName",
				enum:      []string{"Apple Music ID", "Artist", "Title", "Is Explicit", "Lyrics Snippet", "Lyric Snippet Synced", "Artwork", "Video URL", "Shazam URL", "Apple Music URL", "Name"},
			},
		},
	},
	"springBoard": {
		identifier: "openapp",
		make: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFAppIdentifier": "com.apple.springboard",
				"WFSelectedApp": map[string]any{
					"BundleIdentifer": "com.apple.springboard",
					"Name":            "SpringBoard",
					"TeamIdentifer":   "0000000000",
				},
			}
		},
	},
	"getShortcutDetail": {
		identifier: "properties.workflow",
		parameters: []parameterDefinition{
			{
				name:      "shortcut",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "detail",
				validType: String,
				key:       "WFContentItemPropertyName",
				enum:      []string{"Folder", "Icon", "Action Count", "File Size", "File Extension Creation Date", "File Path", "Last Modified Date", "Name"},
			},
		},
	},
	"getApps": {
		identifier: "filter.apps",
		mac:        true,
		minVersion: 18,
	},
	"setSoundRecognition": {
		appIdentifier: "com.apple.AccessibilityUtilities",
		identifier:    "AXSettingsShortcuts.AXToggleSoundDetectionIntent",
		parameters: []parameterDefinition{
			{
				name:         "operation",
				validType:    String,
				key:          "operation",
				defaultValue: "activate",
				enum:         soundRecognitionOperations,
				optional:     true,
			},
		},
	},
	"setTextSize": {
		appIdentifier: "com.apple.AccessibilityUtilities",
		identifier:    "AXSettingsShortcuts.AXSetLargeTextIntent",
		parameters: []parameterDefinition{
			{
				name:      "size",
				validType: String,
				key:       "textSize",
				enum:      textSizes,
			},
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
	// "contacts",
	// "documents",
	// "location",
	// "math",
	// "media",
	// "scripting",
	"sharing",
	// "settings",
	// "shortcuts",
	// "translation",
	"web",
}

func loadStandardActions() {
	includeStandardActions()
	handleIncludes()
	parseActionDefinitions()
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

	maps.Copy(adjustDateParams, map[string]any{
		"WFAdjustOperation": operation,
	})
	if unit == "" {
		return adjustDateParams
	}

	maps.Copy(adjustDateParams, magnitudeValue(unit, args, 1))

	return
}

func magnitudeValue(unit string, args []actionArgument, index int) map[string]any {
	var magnitudeValue = argumentValue(args, index)
	if reflect.TypeOf(magnitudeValue).String() == "[]map[string]any" {
		var value = magnitudeValue.([]map[string]any)
		magnitudeValue = value[0]
	}

	return map[string]any{
		"WFDuration": map[string]any{
			"Value": map[string]any{
				"Unit":      unit,
				"Magnitude": magnitudeValue,
			},
			"WFSerializationType": "WFQuantityFieldValue",
		},
	}
}

func changeCase(textCase string) map[string]any {
	return map[string]any{
		"Show-text":  true,
		"WFCaseType": textCase,
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
