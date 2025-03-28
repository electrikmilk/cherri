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
var httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
var httpParams = []parameterDefinition{
	{
		name:      "url",
		key:       "WFURL",
		validType: String,
	},
	{
		name:         "method",
		key:          "WFHTTPMethod",
		validType:    String,
		optional:     true,
		enum:         httpMethods,
		defaultValue: "GET",
	},
	{
		name:      "body",
		validType: Dict,
		optional:  true,
	},
	{
		name:      "headers",
		key:       "WFHTTPHeaders",
		validType: Dict,
		optional:  true,
	},
}
var cropPositions = []string{"Center", "Top Left", "Top Right", "Bottom Left", "Bottom Right", "Custom"}
var sortOrders = []string{"asc", "desc"}
var windowSortings = []string{"Title", "App Name", "Width", "Height", "X Position", "Y Position", "Window Index", "Name", "Random"}
var windowPositions = []string{"Top Left", "Top Center", "Top Right", "Middle Left", "Center", "Middle Right", "Bottom Left", "Bottom Center", "Bottom Right", "Coordinates"}
var windowConfigurations = []string{"Fit Screen", "Top Half", "Bottom Half", "Left Half", "Right Half", "Top Left Quarter", "Top Right Quarter", "Bottom Left Quarter", "Bottom Right Quarter", "Dimensions"}
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
var eventDetails = []string{"Start Date", "End Date", "Is All Day", "Location", "Duration", "My Status", "Attendees", "URL", "Title", "Notes", "Attachments"}
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

var toggleAlarmIntent = appIntent{
	name:                "Clock",
	bundleIdentifier:    "com.apple.clock",
	appIntentIdentifier: "ToggleAlarmIntent",
}

// actions is the data structure that determines every action the compiler knows about.
// The key determines the identifier of the identifier that must be used in the syntax, it's value defines its behavior, etc. using an actionDefinition.
var actions = map[string]*actionDefinition{
	"date": {
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
				key:       "WFDateActionDate",
			},
		},
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFDateActionMode",
					dataType: Text,
					value:    "Specified Date",
				},
			}
		},
	},
	"currentDate": {
		identifier: "date",
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFDateActionMode",
					dataType: Text,
					value:    "Current Date",
				},
			}
		},
	},
	"addCalendar": {
		identifier: "addnewcalendar",
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
				key:       "CalendarName",
			},
		},
	},
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
			return adjustDate("Get Start of Year", "", args)
		},
	},
	"editEvent": {
		identifier: "setters.calendarevents",
		parameters: []parameterDefinition{
			{
				name:      "event",
				validType: Variable,
				key:       "WFInput",
			},
			{
				name:      "detail",
				validType: String,
				key:       "WFContentItemPropertyName",
				enum:      eventDetails,
			},
			{
				name:      "newValue",
				validType: String,
				key:       "WFCalendarEventContentItemStartDate",
			},
		},
	},
	"formatTime": {
		identifier: "format.date",
		parameters: []parameterDefinition{
			{
				name:      "time",
				validType: Variable,
				key:       "WFDate",
			},
			{
				name:         "timeFormat",
				validType:    String,
				key:          "WFTimeFormatStyle",
				defaultValue: "Short",
				enum:         timeFormats,
				optional:     true,
			},
		},
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFDateFormatStyle",
					dataType: Text,
					value:    "None",
				},
			}
		},
	},
	"formatDate": {
		identifier:    "format.date",
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: Variable,
				key:       "WFDate",
			},
			{
				name:         "dateFormat",
				validType:    String,
				key:          "WFDateFormatStyle",
				defaultValue: "Short",
				enum:         dateFormats,
				optional:     true,
			},
		},
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFTimeFormatStyle",
					dataType: Text,
					value:    "None",
				},
			}
		},
	},
	"formatTimestamp": {
		identifier: "format.date",
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: Variable,
				key:       "WFDate",
			},
			{
				name:         "dateFormat",
				validType:    String,
				key:          "WFDateFormatStyle",
				defaultValue: "Short",
				enum:         dateFormats,
				optional:     true,
			},
			{
				name:         "timeFormat",
				validType:    String,
				key:          "WFTimeFormatStyle",
				defaultValue: "Short",
				enum:         timeFormats,
				optional:     true,
			},
		},
	},
	"removeEvents": {
		parameters: []parameterDefinition{
			{
				name:      "events",
				validType: Variable,
				key:       "WFInputEvents",
			},
			{
				name:         "includeFutureEvents",
				validType:    Bool,
				key:          "WFCalendarIncludeFutureEvents",
				defaultValue: false,
				optional:     true,
			},
		},
	},
	"removeReminders": {
		parameters: []parameterDefinition{
			{
				name:      "reminders",
				validType: Variable,
				key:       "WFInputReminders",
			},
		},
	},
	"showInCalendar": {
		parameters: []parameterDefinition{
			{
				name:      "event",
				validType: Variable,
				key:       "WFEvent",
			},
		},
	},
	"openRemindersList": {
		identifier: "showlist",
		parameters: []parameterDefinition{
			{
				name:      "list",
				validType: Variable,
				key:       "WFList",
			},
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
		addParams: func(args []actionArgument) []plistData {
			if len(args) < 4 {
				return []plistData{}
			}

			var repeatDays = getArgValue(args[3])
			var repeats []plistData
			for _, day := range repeatDays.([]interface{}) {
				var dayStr = day.(string)
				var dayLower = strings.ToLower(dayStr)
				var dayCap = capitalize(dayStr)

				repeats = append(repeats, plistData{
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "value",
							dataType: Text,
							value:    dayLower,
						},
						{
							key:      "title",
							dataType: Dictionary,
							value: []plistData{
								{
									key:      "key",
									dataType: Text,
									value:    dayCap,
								},
							},
						},
						{
							key:      "identifier",
							dataType: Text,
							value:    dayLower,
						},
						{
							key:      "subtitle",
							dataType: Dictionary,
							value: []plistData{
								{
									key:      "key",
									dataType: Text,
									value:    dayCap,
								},
							},
						},
					},
				})
			}

			return []plistData{
				{
					key:      "repeats",
					dataType: Array,
					value:    repeats,
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "state",
					dataType: Number,
					value:    1,
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "state",
					dataType: Number,
					value:    0,
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "operation",
					dataType: Text,
					value:    "Toggle",
				},
			}
		},
	},
	"getAlarms": {
		appIdentifier: "com.apple.mobiletimer-framework",
		identifier:    "MobileTimerIntents.MTGetAlarmsIntent",
	},
	"filterContacts": {
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
		decomp: func(action *ShortcutAction) (arguments []string) {
			return decompContactValue(action, "WFPhoneNumber", phoneNumber)
		},
	},
	"selectContact": {
		identifier: "selectcontacts",
		parameters: []parameterDefinition{
			{
				name:         "multiple",
				validType:    Bool,
				defaultValue: false,
				key:          "WFSelectMultiple",
			},
		},
	},
	"selectEmailAddress": {
		identifier: "selectemail",
	},
	"selectPhoneNumber": {
		identifier: "selectphone",
	},
	"getContactDetail": {
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
	},
	"call": {
		appIdentifier: "com.apple.mobilephone",
		identifier:    "call",
		parameters: []parameterDefinition{
			{
				name:      "contact",
				validType: Variable,
				key:       "WFCallContact",
			},
		},
	},
	"sendEmail": {
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
	},
	"sendMessage": {
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
	},
	"facetimeCall": {
		appIdentifier: "com.apple.facetime",
		identifier:    "facetime",
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFFaceTimeType",
					dataType: Text,
					value:    "Video",
				},
			}
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
		addParams: func(args []actionArgument) (params []plistData) {
			if len(args) >= 3 {
				if args[2].valueType == Variable {
					params = append(params, argumentValue("WFContactPhoneNumbers", args, 2))
				} else {
					params = append(params, contactValue("WFContactPhoneNumbers", phoneNumber, []actionArgument{args[2]}))
				}
			}

			if len(args) >= 4 {
				if args[3].valueType == Variable {
					params = append(params, argumentValue("WFContactEmails", args, 3))
				} else {
					params = append(params, contactValue("WFContactEmails", emailAddress, []actionArgument{args[3]}))
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Mode",
					dataType: Text,
					value:    "Remove",
				},
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
		check: func(args []actionArgument, _ *actionDefinition) {
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFAppendFileWriteMode",
					dataType: Boolean,
					value:    "Prepend",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFAppendFileWriteMode",
					dataType: Boolean,
					value:    "Append",
				},
			}
		},
	},
	"labelFile": {
		identifier: "file.label",
		parameters: []parameterDefinition{
			{
				name:      "file",
				validType: Var,
				key:       "WFInput",
			},
			{
				name:      "color",
				validType: String,
				optional:  false,
				enum:      fileLabels,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			var color = strings.ToLower(getArgValue(args[1]).(string))

			return []plistData{
				{
					key:      "WFLabelColorNumber",
					dataType: Number,
					value:    fileLabelsMap[color],
				},
			}
		},
	},
	"filterFiles": {
		identifier: "filter.files",
		parameters: []parameterDefinition{
			{
				name:      "files",
				validType: Var,
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
		addParams: func(args []actionArgument) []plistData {
			if len(args) != 1 {
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
	},
	"optimizePDF": {
		identifier: "compresspdf",
		parameters: []parameterDefinition{
			{
				name:      "pdfFile",
				validType: Var,
				key:       "WFInput",
			},
		},
	},
	"getPDFText": {
		identifier: "gettextfrompdf",
		parameters: []parameterDefinition{
			{
				name:      "pdfFile",
				validType: Var,
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
		addParams: func(args []actionArgument) []plistData {
			if len(args) != 1 {
				var richText = getArgValue(args[1]).(bool)
				if richText {
					return []plistData{
						{
							key:      "WFGetTextFromPDFTextType",
							dataType: Text,
							value:    "Rich Text",
						},
					}
				}
			}

			return []plistData{
				{
					key:      "WFGetTextFromPDFTextType",
					dataType: Text,
					value:    "Text",
				},
			}
		},
	},
	"makePDF": {
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Var,
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
			if textArg.valueType == Var {
				args[1] = actionArgument{
					valueType: String,
					value:     fmt.Sprintf("^{%s}", textArg.value),
				}
			} else {
				args[1].value = fmt.Sprintf("^%s", textArg.value)
			}
		},
		addParams: func(args []actionArgument) []plistData {
			if len(args) == 0 {
				return []plistData{}
			}

			return []plistData{
				argumentValue("WFMatchTextPattern", args, 1),
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
				validType: Var,
				key:       "matches",
			},
		},
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetGroupType",
					dataType: Text,
					value:    "All Groups",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetGroupType",
					dataType: Text,
					value:    "Group At Index",
				},
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
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFPickingMode",
					dataType: Text,
					value:    "Folders",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFAskWhereToSave",
					dataType: Boolean,
					value:    false,
				},
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
				enum:      languagesList,
			},
			{
				name:      "to",
				validType: String,
				key:       "WFSelectedLanguage",
				enum:      languagesList,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			args[1].value = languageCode(args[1].value.(string))
			args[2].value = languageCode(args[2].value.(string))
		},
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
				enum:      languagesList,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			if len(args) < 2 {
				return
			}

			if args[1].valueType != Variable {
				args[1].value = languageCode(getArgValue(args[1]).(string))
			}
		},
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFSelectedFromLanguage",
					dataType: Text,
					value:    "Detect Language",
				},
			}
		},
	},
	"detectLanguage": {
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
				key:       "WFInput",
			},
		},
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
		addParams: func(_ []actionArgument) []plistData {
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
		addParams: func(_ []actionArgument) []plistData {
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
		addParams: func(_ []actionArgument) []plistData {
			return changeCase("Capitalize with Title Case")
		},
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
		addParams: func(_ []actionArgument) []plistData {
			return changeCase("Capitalize with sentence case")
		},
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
		addParams: func(_ []actionArgument) []plistData {
			return changeCase("Capitalize Every Word")
		},
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
		addParams: func(_ []actionArgument) []plistData {
			return changeCase("cApItAlIzE wItH aLtErNaTiNg cAsE")
		},
	},
	"correctSpelling": {
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
				key:       "text",
			},
		},
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Show-text",
					dataType: Boolean,
					value:    true,
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
		addParams: textParts,
		decomp:    decompTextParts,
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
		addParams: textParts,
		decomp:    decompTextParts,
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
		addParams: func(_ []actionArgument) []plistData {
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
		addParams: func(args []actionArgument) []plistData {
			var data = []plistData{
				{
					key:      "SizeToFit",
					dataType: Boolean,
					value:    false,
				},
			}

			if len(args) == 0 {
				return append(data, plistData{
					key:      "ImageSize",
					dataType: Dictionary,
					value:    []plistData{},
				})
			}

			var size = strings.Split(getArgValue(args[2]).(string), " ")
			return append(data, plistData{
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
			})
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
				name:         "prompt",
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
		addParams: func(_ []actionArgument) []plistData {
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
		identifier: "getlatestbursts",
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFWhenToPlay",
					dataType: Text,
					value:    "Next",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFWhenToPlay",
					dataType: Text,
					value:    "Later",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFImageFormat",
					dataType: Text,
					value:    "JPEG",
				},
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
	"rotateImage": {
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
				name:         "position",
				validType:    String,
				key:          "WFImageCropPosition",
				enum:         cropPositions,
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
				optional:     true,
				defaultValue: "100",
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFSkipBackBehavior",
					dataType: Text,
					value:    "Previous Song",
				},
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
				validType: Var,
				key:       "WFPodcastShow",
			},
		},
	},
	"getPodcastDetail": {
		identifier: "properties.podcastshow",
		parameters: []parameterDefinition{
			{
				name:      "podcast",
				validType: Var,
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
		addParams: func(_ []actionArgument) []plistData {
			var params = []plistData{
				{
					key:      "WFMediaAudioOnly",
					dataType: Boolean,
					value:    true,
				},
			}
			return params
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Metadata",
					dataType: Boolean,
					value:    true,
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Metadata",
					dataType: Boolean,
					value:    true,
				},
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
		addParams: func(_ []actionArgument) []plistData {
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFPlayPauseBehavior",
					dataType: Text,
					value:    "Play",
				},
			}
		},
	},
	"pause": {
		identifier: "pausemusic",
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFPlayPauseBehavior",
					dataType: Text,
					value:    "Pause",
				},
			}
		},
	},
	"togglePlayPause": {
		identifier:    "pausemusic",
		defaultAction: true,
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFPlayPauseBehavior",
					dataType: Text,
					value:    "Play/Pause",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFTakeScreenshotScreenshotType",
					dataType: Text,
					value:    "Full Screen",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFTakeScreenshotScreenshotType",
					dataType: Text,
					value:    "Interactive",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFShutdownMode",
					dataType: Text,
					value:    "Restart",
				},
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
	"number": {
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: Variable,
				key:       "WFNumberActionNumber",
			},
		},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFFileSizeIncludeUnits",
					dataType: Boolean,
					value:    false,
				},
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
		},
		defaultAction: true,
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFCountType",
					dataType: Text,
					value:    "Items",
				},
			}
		},
	},
	"countChars": {
		identifier:    "count",
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Input",
				validType: String,
			},
		},
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFCountType",
					dataType: Text,
					value:    "Characters",
				},
			}
		},
	},
	"countWords": {
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Input",
				validType: String,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFCountType",
					dataType: Text,
					value:    "Words",
				},
			}
		},
	},
	"countSentences": {
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Input",
				validType: String,
			},
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFCountType",
					dataType: Text,
					value:    "Sentences",
				},
			}
		},
	},
	"countLines": {
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Input",
				validType: String,
			},
		},

		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFCountType",
					dataType: Text,
					value:    "Lines",
				},
			}
		},
	},
	"lightMode": {
		identifier: "appearance",
		addParams: func(_ []actionArgument) []plistData {
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
	},
	"darkMode": {
		identifier:    "appearance",
		defaultAction: true,
		addParams: func(_ []actionArgument) []plistData {
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
	},
	"getBatteryLevel": {defaultAction: true},
	"isCharging": {
		identifier: "getbatterylevel",
		minVersion: 16.2,
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Subject",
					dataType: Text,
					value:    "Is Charging",
				},
			}
		},
	},
	"connectedToCharger": {
		identifier: "getbatterylevel",
		minVersion: 16.2,
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Subject",
					dataType: Text,
					value:    "Is Connected to Charger",
				},
			}
		},
	},
	"getShortcuts": {
		identifier: "getmyworkflows",
	},
	"url": {
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "input",
					dataType: Text,
					value:    "Encode",
				},
			}
		},
		defaultAction: true,
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFEncodeMode",
					dataType: Text,
					value:    "Decode",
				},
			}
		},
	},
	"urlEncode": {
		identifier: "urlencode",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
				key:       "WFInput",
			},
		},
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFEncodeMode",
					dataType: Text,
					value:    "Encode",
				},
			}
		},
		defaultAction: true,
	},
	"urlDecode": {
		identifier: "urlencode",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
				key:       "WFInput",
			},
		},
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFEncodeMode",
					dataType: Text,
					value:    "Decode",
				},
			}
		},
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
	"notification": {
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
	},
	"stop": {
		identifier: "exit",
	},
	"comment": {
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFAlertActionCancelButtonShown",
					dataType: Boolean,
					value:    false,
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFAlertActionCancelButtonShown",
					dataType: Boolean,
					value:    true,
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetDictionaryValueType",
					dataType: Text,
					value:    "All Keys",
				},
			}
		},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetDictionaryValueType",
					dataType: Text,
					value:    "All Values",
				},
			}
		},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetDictionaryValueType",
					dataType: Text,
					value:    "Value",
				},
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
				validType: Variable,
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
		addParams: func(args []actionArgument) (params []plistData) {
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
		make: func(args []actionArgument) []plistData {
			if args[0].valueType == Variable {
				return []plistData{
					argumentValue("WFApp", args, 0),
				}
			}

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
		make: func(args []actionArgument) (params []plistData) {
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFHideAppMode",
					dataType: Text,
					value:    "All Apps",
				},
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
		make: func(args []actionArgument) []plistData {
			if args[0].valueType == Variable {
				return []plistData{
					argumentValue("WFApp", args, 0),
				}
			}

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
		make: func(args []actionArgument) (params []plistData) {
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFQuitAppMode",
					dataType: Text,
					value:    "All Apps",
				},
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
		make: func(args []actionArgument) (params []plistData) {
			params = []plistData{
				{
					key:      "WFAskToSaveChanges",
					dataType: Boolean,
					value:    false,
				},
			}

			if args[0].valueType == Variable {
				return append(params, argumentValue("WFApp", args, 0))
			}

			return append(params, plistData{
				key:      "WFApp",
				dataType: Dictionary,
				value: []plistData{
					argumentValue("BundleIdentifier", args, 0),
				},
			})
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
		addParams: func(args []actionArgument) (params []plistData) {
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
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "target",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "title",
							dataType: Text,
							value: []plistData{
								argumentValue("key", args, 0),
							},
						},
					},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFWorkflow",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "workflowIdentifier",
							dataType: Text,
							value:    uuid.New().String(),
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
							value:    uuid.New().String(),
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
		decomp: func(action *ShortcutAction) (arguments []string) {
			var workflow = action.WFWorkflowActionParameters["WFWorkflow"].(map[string]any)
			if !workflow["isSelf"].(bool) {
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFMathOperation",
					dataType: Text,
					value:    "...",
				},
			}
		},
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
	"isOnline": {
		identifier: "getipaddress",
		make: func(_ []actionArgument) []plistData {
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
	},
	"getLocalIP": {
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFIPAddressSourceOption",
					dataType: Text,
					value:    "Local",
				},
			}
		},
	},
	"getExternalIP": {
		identifier:    "getipaddress",
		defaultAction: true,
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFIPAddressSourceOption",
					dataType: Text,
					value:    "External",
				},
			}
		},
	},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "First Item",
				},
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
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "Last Item",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "Random Item",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "Item At Index",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFItemSpecifier",
					dataType: Text,
					value:    "Items in Range",
				},
			}
		},
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
	},
	"getContacts": {
		identifier: "detect.contacts",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	},
	"getDates": {
		identifier: "detect.date",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	},
	"getEmails": {
		identifier: "detect.emailaddress",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: String,
			},
		},
	},
	"getPhoneNumbers": {
		identifier: "detect.phonenumber",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	},
	"getURLs": {
		identifier: "detect.link",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
				key:       "WFInput",
			},
		},
	},
	"getAllWallpapers": {
		identifier:    "posters.get",
		minVersion:    16.2,
		mac:           false,
		defaultAction: true,
	},
	"getWallpaper": {
		identifier: "posters.get",
		minVersion: 16.2,
		mac:        false,
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFPosterType",
					dataType: Text,
					value:    "Current",
				},
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
		addParams: func(args []actionArgument) (xCallbackParams []plistData) {
			if len(args) == 0 {
				return
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFNoOutputSurfaceBehavior",
					dataType: Text,
					value:    "Respond",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFNoOutputSurfaceBehavior",
					dataType: Text,
					value:    "Copy to Clipboard",
				},
			}
		},
	},
	"DNDOn": {
		identifier: "dnd.set",
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				focusModes,
				{
					key:      "Enabled",
					dataType: Number,
					value:    1,
				},
			}
		},
	},
	"DNDOff": {
		identifier:    "dnd.set",
		defaultAction: true,
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				focusModes,
				{
					key:      "Enabled",
					dataType: Number,
					value:    0,
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFRoundMode",
					dataType: Text,
					value:    "Normal",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFRoundMode",
					dataType: Text,
					value:    "Always Round Up",
				},
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFRoundMode",
					dataType: Text,
					value:    "Always Round Down",
				},
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
	"airdrop": {
		identifier: "airdropdocument",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	},
	"share": {
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: String,
			},
		},
	},
	"copyToClipboard": {
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
	},
	"getClipboard": {},
	"getURLHeaders": {
		identifier: "url.getheaders",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFInput",
			},
		},
	},
	"openURL": {
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFInput",
			},
		},
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Show-WFInput",
					dataType: Boolean,
					value:    true,
				},
			}
		},
	},
	"runJavaScriptOnWebpage": {
		parameters: []parameterDefinition{
			{
				name:      "javascript",
				validType: String,
				key:       "WFJavaScript",
			},
		},
	},
	"searchWeb": {
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
	},
	"showWebpage": {
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
	},
	"getRSSFeeds": {
		identifier: "rss.extract",
		parameters: []parameterDefinition{
			{
				name:      "urls",
				validType: String,
				key:       "WFURLs",
			},
		},
	},
	"getRSS": {
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
	},
	"getWebPageDetail": {
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
	},
	"getArticleDetail": {
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
	},
	"getCurrentURL": {
		identifier: "safari.geturl",
	},
	"getWebpageContents": {
		identifier: "getwebpagecontents",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFInput",
			},
		},
	},
	"searchGiphy": {
		identifier: "giphy",
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "WFGiphyQuery",
			},
		},
	},
	"getGifs": {
		identifier:    "giphy",
		defaultAction: true,
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGiphyShowPicker",
					dataType: Boolean,
					value:    false,
				},
			}
		},
	},
	"getArticle": {
		identifier: "getarticle",
		parameters: []parameterDefinition{
			{
				name:      "webpage",
				validType: String,
				key:       "WFWebPage",
			},
		},
	},
	"expandURL": {
		identifier: "url.expand",
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "URL",
			},
		},
	},
	"getURLDetail": {
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
	},
	"downloadURL": {
		defaultAction: true,
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFHTTPMethod",
					dataType: Text,
					value:    "GET",
				},
			}
		},
	},
	"formRequest": {
		identifier: "downloadurl",
		parameters: httpParams,
		addParams: func(args []actionArgument) []plistData {
			return httpRequest("Form", "WFFormValues", args)
		},
	},
	"jsonRequest": {
		identifier: "downloadurl",
		parameters: httpParams,
		addParams: func(args []actionArgument) []plistData {
			return httpRequest("JSON", "WFJSONValues", args)
		},
	},
	"fileRequest": {
		identifier: "downloadurl",
		parameters: httpParams,
		addParams: func(args []actionArgument) []plistData {
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
			if reflect.TypeOf(value).String() != stringType {
				return
			}

			makeMeasurementUnits()

			var unitType = value.(string)
			checkEnum(&parameterDefinition{
				name: "measurement unit",
				enum: units[unitType],
			}, &args[2])
		},
		addParams: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFMeasurementUnit",
					dataType: Dictionary,
					value: []plistData{
						argumentValue("WFNSUnitType", args, 1),
						argumentValue("WFNSUnitSymbol", args, 2),
					},
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
		addParams: func(args []actionArgument) []plistData {
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
				name:      "imagePath",
				validType: String,
				optional:  true,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			if len(args) != 3 {
				return
			}

			var image = getArgValue(args[2])
			if reflect.TypeOf(image).String() == stringType {
				var iconFile = getArgValue(args[2]).(string)
				if _, err := os.Stat(iconFile); os.IsNotExist(err) {
					parserError(fmt.Sprintf("File '%s' does not exist!", iconFile))
				}
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
				var photo string
				var image = getArgValue(args[2])
				if reflect.TypeOf(image).String() != stringType && args[2].valueType == Variable {
					photo = fmt.Sprintf("{%s}", args[2].value)
				} else {
					var iconFile = getArgValue(args[2]).(string)
					var bytes, readErr = os.ReadFile(iconFile)
					handle(readErr)
					photo = base64.StdEncoding.EncodeToString(bytes)
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
			return []plistData{
				argumentValue("WFTextActionText", args, 0),
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
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Mode",
					dataType: Text,
					value:    "Set",
				},
			}
		},
	},
	"getShazamDetail": {
		identifier: "properties.shazam",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Var,
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
		make: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFAppIdentifier",
					dataType: Text,
					value:    "com.apple.springboard",
				},
				{
					key:      "WFSelectedApp",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "BundleIdentifier",
							dataType: Text,
							value:    "com.apple.springboard",
						},
						{
							key:      "Name",
							dataType: Text,
							value:    "SpringBoard",
						},
						{
							key:      "TeamIdentifier",
							dataType: Text,
							value:    "0000000000",
						},
					},
				},
			}
		},
	},
	"getShortcutDetail": {
		identifier: "properties.workflow",
		parameters: []parameterDefinition{
			{
				name:      "shortcut",
				validType: Var,
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

func rawAction() {
	actions["rawAction"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "identifier",
				validType: String,
			},
			{
				name:      "parameters",
				optional:  true,
				validType: Arr,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			actions["rawAction"].overrideIdentifier = getArgValue(args[0]).(string)
		},
		make: func(args []actionArgument) (params []plistData) {
			for _, parameterDefinitions := range getArgValue(args[1]).([]interface{}) {
				var paramKey string
				var paramType plistDataType
				var rawValue any
				for key, value := range parameterDefinitions.(map[string]interface{}) {
					switch key {
					case "key":
						paramKey = value.(string)
					case "type":
						paramType = plistDataType(value.(string))
					case "value":
						rawValue = value
					}
				}

				var tokenType = convertPlistTypeToken(paramType)
				params = append(params, paramValue(paramKey, actionArgument{
					valueType: tokenType,
					value:     rawValue,
				}, tokenType, paramType))
			}
			return
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

func adjustDate(operation string, unit string, args []actionArgument) (adjustDateParams []plistData) {
	adjustDateParams = []plistData{
		{
			key:      "WFAdjustOperation",
			dataType: Text,
			value:    operation,
		},
	}
	if unit == "" {
		return adjustDateParams
	}

	adjustDateParams = append(adjustDateParams, magnitudeValue(unit, args, 1))

	return adjustDateParams
}

func magnitudeValue(unit string, args []actionArgument, index int) plistData {
	var magnitudeValue = argumentValue("Magnitude", args, index)
	if magnitudeValue.dataType == Dictionary {
		var value = magnitudeValue.value.([]plistData)
		magnitudeValue.dataType = Dictionary
		magnitudeValue.value = value[0].value
	}

	return plistData{
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
	}
}

func changeCase(textCase string) []plistData {
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
	}
}

func textParts(args []actionArgument) []plistData {
	var data = []plistData{
		{
			key:      "Show-text",
			dataType: Boolean,
			value:    true,
		},
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
	case separator == "" && currentAction.identifier == "splitText":
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

func decompTextParts(action *ShortcutAction) (arguments []string) {
	arguments = append(arguments, decompValue(action.WFWorkflowActionParameters["text"]))
	var glue string
	if action.WFWorkflowActionParameters["WFTextSeparator"] != nil {
		glue = action.WFWorkflowActionParameters["WFTextSeparator"].(string)
		if glue == "New Lines" {
			return
		}
	}
	if action.WFWorkflowActionParameters["WFTextCustomSeparator"] != nil {
		glue = action.WFWorkflowActionParameters["WFTextCustomSeparator"].(string)
	}
	if glue != "" {
		arguments = append(arguments, fmt.Sprintf("\"%s\"", glueToChar(glue)))
	}

	return
}

func languageCode(language string) string {
	if lang, found := languages[language]; found {
		return lang
	}

	parserError(fmt.Sprintf("Unknown language '%s'", language))
	return ""
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

func httpRequest(bodyType string, valuesKey string, args []actionArgument) (params []plistData) {
	params = []plistData{
		{
			key:      "WFHTTPBodyType",
			dataType: Text,
			value:    bodyType,
		},
	}
	if len(args) > 0 {
		params = append(params, argumentValue(valuesKey, args, 2))
	}
	return
}

func decompAppAction(key string, action *ShortcutAction) (arguments []string) {
	if action.WFWorkflowActionParameters[key] != nil {
		var appsType = reflect.TypeOf(action.WFWorkflowActionParameters[key])
		if appsType.Kind() == reflect.String {
			return append(arguments, decompValue(action.WFWorkflowActionParameters[key]))
		}

		if appsType.String() == dictType {
			for key, bundle := range action.WFWorkflowActionParameters[key].(map[string]interface{}) {
				if key == "BundleIdentifier" {
					arguments = append(arguments, fmt.Sprintf("\"%s\"", bundle))
				}
			}
		} else if appsType.String() == "[]interface {}" {
			for _, app := range action.WFWorkflowActionParameters[key].([]interface{}) {
				var bundleIdentifer = app.(map[string]interface{})["BundleIdentifier"]
				arguments = append(arguments, fmt.Sprintf("\"%s\"", bundleIdentifer))
			}
		}
	}

	return
}

func decompInfiniteURLAction(action *ShortcutAction) (arguments []string) {
	var urlValueType = reflect.TypeOf(action.WFWorkflowActionParameters["WFURLActionURL"]).String()
	if urlValueType == dictType || urlValueType == "string" {
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
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleBackgroundSoundsIntent",
	},
	"MediaBackgroundSounds": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleBackgroundSoundsIntent",
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "setting",
					dataType: Text,
					value:    "whenMediaIsPlaying",
				},
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
}

// ToggleSetActions automates the creation of actions which simply toggle and set a state in the same format.
func ToggleSetActions() {
	for name, def := range toggleSetActions {
		var toggleName = fmt.Sprintf("toggle%s", name)
		def.addParams = func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "operation",
					dataType: Text,
					value:    "Toggle",
				},
			}
		}
		var toggleDef = def
		toggleDef.parameters = nil
		actions[toggleName] = &toggleDef

		if name == "Appearance" {
			continue
		}

		var setName = fmt.Sprintf("set%s", name)
		def.addParams = nil
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
