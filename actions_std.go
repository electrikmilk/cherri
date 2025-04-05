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

var toggleAlarmIntent = appIntent{
	name:                "Clock",
	bundleIdentifier:    "com.apple.clock",
	appIntentIdentifier: "ToggleAlarmIntent",
}

// actions is the data structure that determines every action the compiler knows about.
// The key determines the identifier of the identifier that must be used in the syntax, it's value defines its behavior, etc. using an actionDefinition.
var actions = map[string]*actionDefinition{
	"date": {
		category:      Dates,
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "date",
				validType: String,
				key:       "WFDateActionDate",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFDateActionMode": "Specified Date",
			}
		},
		doc: ActionDoc{
			title:       "Date",
			description: "Create date value from `date`.",
			example:     "date(\"October 5, 2022\")",
		},
	},
	"currentDate": {
		category:   Dates,
		identifier: "date",
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFDateActionMode": "Current Date",
			}
		},
	},
	"addCalendar": {
		category:   Calendars,
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
		category:      Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
		category:   Dates,
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
	"editEvent": {
		category:   Calendars,
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
				enum:      editEventDetails,
			},
			{
				name:      "newValue",
				validType: String,
				key:       "WFCalendarEventContentItemStartDate",
			},
		},
	},
	"getEventDetail": {
		identifier: "properties.calendarevents",
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
		},
	},
	"formatTime": {
		category:   Dates,
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFDateFormatStyle": "None",
			}
		},
	},
	"formatDate": {
		category:      Dates,
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFTimeFormatStyle": "None",
			}
		},
	},
	"formatTimestamp": {
		category:   Dates,
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
		category: Calendars,
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
		category: Reminders,
		parameters: []parameterDefinition{
			{
				name:      "reminders",
				validType: Variable,
				key:       "WFInputReminders",
			},
		},
	},
	"showInCalendar": {
		category: Calendars,
		parameters: []parameterDefinition{
			{
				name:      "event",
				validType: Variable,
				key:       "WFEvent",
			},
		},
	},
	"openRemindersList": {
		category:   Reminders,
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
		category:   Clock,
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
		category:      Clock,
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
		category:      Clock,
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
		category:      Clock,
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
		category:      Clock,
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
		category: Clock,
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
	"getAlarms": {
		category:      Clock,
		appIdentifier: "com.apple.mobiletimer-framework",
		identifier:    "MobileTimerIntents.MTGetAlarmsIntent",
	},
	"filterContacts": {
		category:   Contacts,
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
		addParams: func(args []actionArgument) (params map[string]any) {
			if len(args) == 4 {
				return map[string]any{
					"WFContentItemLimitEnabled": true,
				}
			}
			return
		},
	},
	"emailAddress": {
		category:   Email,
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
		category: Phone,
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
	"selectContact": {
		category:   Contacts,
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
		category:   Email,
		identifier: "selectemail",
	},
	"selectPhoneNumber": {
		category:   Phone,
		identifier: "selectphone",
	},
	"getContactDetail": {
		category:   Contacts,
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
		category:      Phone,
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
		category: Email,
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
		category: Messaging,
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
		category:      Phone,
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFFaceTimeType": "Video",
			}
		},
	},
	"newContact": {
		category:   Contacts,
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
		category:   Contacts,
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
		category:   Speech,
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
		category:   Files,
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
		category:      Files,
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
		category:   Files,
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
		category:   Files,
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
		category:   Documents,
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
		category:   Documents,
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
		category: Documents,
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
		category:   Speech,
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
		category:   Files,
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
		category:   TextEditing,
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
			if textArg.valueType == Var {
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
		category:      TextEditing,
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
		category:   TextEditing,
		identifier: "text.match.getgroup",
		parameters: []parameterDefinition{
			{
				name:      "matches",
				validType: Var,
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
		category:   TextEditing,
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
		category:   Files,
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
		category:      Files,
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
		category:   Documents,
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
		category:   Files,
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
		category:   Files,
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
		category:   Documents,
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
		category:   QRCodes,
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
		category:   Notes,
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
		category: Documents,
		parameters: []parameterDefinition{
			{
				name:      "pdf",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"makeHTML": {
		category:   Documents,
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
		category:   Documents,
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
		category: TextEditing,
		parameters: []parameterDefinition{
			{
				name:      "html",
				validType: Variable,
				key:       "WFHTML",
			},
		},
	},
	"getRichTextFromMarkdown": {
		category: TextEditing,
		parameters: []parameterDefinition{
			{
				name:      "markdown",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"print": {
		category: Printing,
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"selectFile": {
		category:      Files,
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
		category:   Files,
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
		category:   Files,
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
		category: Files,
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"getEmojiName": {
		category:   Scripting,
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
		category:   Files,
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
		category:   Files,
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
		category:   Images,
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
		category:   Network,
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
		category: Notes,
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
		category:      Books,
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
		category:   Files,
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
		category:      Files,
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
	"getSelectedFiles": {
		category:   Files,
		identifier: "finder.getselectedfiles",
	},
	"extractArchive": {
		category:   Archives,
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
		category:   Archives,
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
		category:   Previewing,
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
		category:   Translation,
		identifier: "text.translate",
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
				key:       "WFInputText",
			},
			{
				name:      "fromLanguage",
				validType: String,
				key:       "WFSelectedFromLanguage",
				enum:      languages,
			},
			{
				name:      "toLanguage",
				validType: String,
				key:       "WFSelectedLanguage",
				enum:      languages,
			},
		},
	},
	"translate": {
		category:      Translation,
		identifier:    "text.translate",
		defaultAction: true,
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: String,
				key:       "WFInputText",
			},
			{
				name:      "language",
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
	},
	"detectLanguage": {
		category: Translation,
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
		category:   TextEditing,
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
		category:   TextEditing,
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
		category:   TextEditing,
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
	},
	"capitalize": {
		category:   TextEditing,
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
	},
	"capitalizeAll": {
		category:   TextEditing,
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
	},
	"alternateCase": {
		category:   TextEditing,
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
	},
	"correctSpelling": {
		category: TextEditing,
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
	},
	"splitText": {
		category:   TextEditing,
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
		category:   TextEditing,
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
		category: Disk,
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
		category:      Disk,
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
		category:   Files,
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
		category:      Speech,
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
		category:   Location,
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
		category:   Maps,
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
		category:   Weather,
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
		category:   Maps,
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
		category:   Location,
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
		category:   Weather,
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
		category:   Weather,
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
		category:   Location,
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
		category: Maps,
		parameters: []parameterDefinition{
			{
				name:      "location",
				validType: Variable,
				key:       "WFInput",
			},
		},
	},
	"getHalfwayPoint": {
		category: Maps,
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
		category:   Images,
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
	"clearUpNext":    {category: Music},
	"getCurrentSong": {category: Music},
	"getLastImport": {
		category:   Photos,
		identifier: "getlatestphotoimport",
	},
	"getLatestBursts": {
		category:   Photos,
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
		category:   Photos,
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
		category:   Photos,
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
		category:   Photos,
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
		category:   Photos,
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
		category:   Images,
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
		category: Camera,
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
		category: Video,
		parameters: []parameterDefinition{
			{
				name:      "video",
				validType: Variable,
				key:       "WFInputMedia",
			},
		},
	},
	"takeVideo": {
		category: Video,
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
		category: Audio,
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
		category:   Music,
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
		category:      Music,
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
		category:      Music,
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
		category:   Music,
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
		category: Music,
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
		category:   GIFs,
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
		category:   Images,
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
		category:   Images,
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
		category:   Images,
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
		category:   Photos,
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
		category:   Photos,
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
		category:   Images,
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
		category: Photos,
		parameters: []parameterDefinition{
			{
				name:      "photos",
				validType: Variable,
				key:       "photos",
			},
		},
	},
	"removeFromAlbum": {
		category:   Photos,
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
		category:   Music,
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
		category: Music,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFSkipBackBehavior": "Previous Song",
			}
		},
	},
	"skipFwd": {
		category:   Music,
		identifier: "skipforward",
	},
	"searchAppStore": {
		category: Stores,
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "WFSearchTerm",
			},
		},
	},
	"searchPodcasts": {
		category: Podcasts,
		parameters: []parameterDefinition{
			{
				name:      "query",
				validType: String,
				key:       "WFSearchTerm",
			},
		},
	},
	"getPodcasts": {
		category:   Podcasts,
		identifier: "getpodcastsfromlibrary",
	},
	"playPodcast": {
		category: Podcasts,
		parameters: []parameterDefinition{
			{
				name:      "podcast",
				validType: Var,
				key:       "WFPodcastShow",
			},
		},
	},
	"getPodcastDetail": {
		category:   Podcasts,
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
		category: GIFs,
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
		category:   Images,
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
		category: Audio,
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
		category:      Images,
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
		category:      Video,
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
		category:   Audio,
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
		category:   Media,
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
		category:   Media,
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
		category:   Images,
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
		category:   Photos,
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
		category:   Music,
		identifier: "pausemusic",
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFPlayPauseBehavior": "Play",
			}
		},
	},
	"pause": {
		category:   Music,
		identifier: "pausemusic",
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFPlayPauseBehavior": "Pause",
			}
		},
	},
	"togglePlayPause": {
		category:      Music,
		identifier:    "pausemusic",
		defaultAction: true,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFPlayPauseBehavior": "Play/Pause",
			}
		},
	},
	"startShazam": {
		category:   Music,
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
		category:   Stores,
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
		category:      Device,
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
		category:   Device,
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
		category:      Device,
		identifier:    "reboot",
		minVersion:    17,
		defaultAction: true,
	},
	"reboot": {
		category:   Device,
		minVersion: 17,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFShutdownMode": "Restart",
			}
		},
	},
	"sleep": {
		category:   Device,
		minVersion: 17,
		mac:        true,
	},
	"displaySleep": {
		category:   Device,
		minVersion: 17,
		mac:        true,
	},
	"logout": {
		category:   Device,
		minVersion: 17,
		mac:        true,
	},
	"lockScreen": {
		category:   Device,
		minVersion: 17,
	},
	"number": {
		category: Numbers,
		parameters: []parameterDefinition{
			{
				name:      "number",
				validType: Variable,
				key:       "WFNumberActionNumber",
			},
		},
	},
	"getObjectOfClass": {
		category:   Scripting,
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
	"getOnScreenContent": {category: Device},
	"fileSize": {
		category:   Files,
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
		category:   Device,
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
		category: Device,
		parameters: []parameterDefinition{
			{
				name:      "brightness",
				validType: Float,
				key:       "WFBrightness",
			},
		},
	},
	"getName": {
		category:   Scripting,
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
		category:   Scripting,
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
		category: Scripting,
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Input",
				validType: Variable,
			},
		},
		defaultAction: true,
		addParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"WFCountType": "Items",
			}
		},
	},
	"countChars": {
		category:   Scripting,
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Input",
				validType: String,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFCountType": "Characters",
			}
		},
	},
	"countWords": {
		category:   Scripting,
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Input",
				validType: String,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"WFCountType": "Words",
			}
		},
	},
	"countSentences": {
		category:   Scripting,
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Input",
				validType: String,
			},
		},
		addParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"WFCountType": "Sentences",
			}
		},
	},
	"countLines": {
		category:   Scripting,
		identifier: "count",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Input",
				validType: String,
			},
		},

		addParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"WFCountType": "Lines",
			}
		},
	},
	"lightMode": {
		category:   Device,
		identifier: "appearance",
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"operation": "set",
				"style":     "light",
			}
		},
	},
	"darkMode": {
		category:      Device,
		identifier:    "appearance",
		defaultAction: true,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"operation": "set",
				"style":     "dark",
			}
		},
	},
	"getBatteryLevel": {
		category:      Device,
		defaultAction: true,
	},
	"isCharging": {
		category:   Device,
		identifier: "getbatterylevel",
		minVersion: 16.2,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"Subject": "Is Charging",
			}
		},
	},
	"connectedToCharger": {
		category:   Device,
		identifier: "getbatterylevel",
		minVersion: 16.2,
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"Subject": "Is Connected to Charger",
			}
		},
	},
	"getShortcuts": {
		category:   Shortcuts,
		identifier: "getmyworkflows",
	},
	"url": {
		category: URLs,
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
		category:   Safari,
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
		category: Crypto,
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
		category:   Numbers,
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
		category:   Numbers,
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
		category:   Crypto,
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
	},
	"base64Decode": {
		category:   Crypto,
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
	},
	"urlEncode": {
		category:   URLs,
		identifier: "urlencode",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
				key:       "WFInput",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFEncodeMode": "Encode",
			}
		},
		defaultAction: true,
	},
	"urlDecode": {
		category:   URLs,
		identifier: "urlencode",
		parameters: []parameterDefinition{
			{
				name:      "input",
				validType: String,
				key:       "WFInput",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFEncodeMode": "Decode",
			}
		},
	},
	"show": {
		category:   Previewing,
		identifier: "showresult",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "Text",
				validType: String,
			},
		},
	},
	"waitToReturn": {category: ControlFlow},
	"notification": {
		category:   Notifications,
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
		category:   ControlFlow,
		identifier: "exit",
	},
	"comment": {
		category: Noonce,
		parameters: []parameterDefinition{
			{
				name:      "text",
				validType: RawString,
				key:       "WFCommentActionText",
			},
		},
	},
	"nothing": {category: Noonce},
	"wait": {
		category:   ControlFlow,
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
		category: Dialogs,
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
		category:      Dialogs,
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
		category:   Dialogs,
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
		category: Lists,
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
		category:   Scripting,
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
		category:   Dictionaries,
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
	},
	"getValues": {
		category:   Dictionaries,
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
	},
	"getValue": {
		category:      Dictionaries,
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
		category:   Dictionaries,
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
		category:      Apps,
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
		category:   Apps,
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
		category:      Apps,
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
		category:   Apps,
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
		category:   Apps,
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
		category:   Apps,
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
		category:   Apps,
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
		category:      Shortcuts,
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
		category:   Shortcuts,
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
		category:      Shortcuts,
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
		category: Lists,
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
		category:   Math,
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
	},
	"statistic": {
		category:   Math,
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
	"dismissSiri": {category: System},
	"isOnline": {
		category:   Network,
		identifier: "getipaddress",
		make: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFIPAddressSourceOption": "External",
				"WFIPAddressTypeOption":   "IPv4",
			}
		},
	},
	"getLocalIP": {
		category:   Network,
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFIPAddressSourceOption": "Local",
			}
		},
	},
	"getExternalIP": {
		category:      Network,
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFIPAddressSourceOption": "External",
			}
		},
	},
	"getFirstItem": {
		category:      Lists,
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
		category:   Lists,
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
		category:   Lists,
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
		category:   Lists,
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
		category:   Lists,
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
	},
	"getNumbers": {
		category:   Numbers,
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
		category:   Dictionaries,
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
		category:   Contacts,
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
		category:   Dates,
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
		category:   Email,
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
		category:   Phone,
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
		category:   URLs,
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
		category:      System,
		identifier:    "posters.get",
		minVersion:    16.2,
		mac:           false,
		defaultAction: true,
	},
	"getWallpaper": {
		category:   System,
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
		category:   System,
		identifier: "wallpaper.set",
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	},
	"startScreensaver": {
		category: System,
		mac:      true,
	},
	"contentGraph": {
		category:   Previewing,
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
		category:      XCallback,
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
		category:   XCallback,
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
		category:      ControlFlow,
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
		category:   ControlFlow,
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
		category:   ControlFlow,
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
		category:   Settings,
		identifier: "dnd.set",
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"FocusModes": focusModes,
				"Enabled":    1,
			}
		},
	},
	"DNDOff": {
		category:      Settings,
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
		category:      Settings,
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
		category:      Settings,
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
		category: Audio,
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: Variable,
			},
		},
	},
	"round": {
		category:      Math,
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
		category:   Math,
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
		category:   Math,
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
		category: Scripts,
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
		category:      Shortcuts,
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
		category:      Shortcuts,
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
		category:   Sharing,
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
		category: Sharing,
		parameters: []parameterDefinition{
			{
				name:      "input",
				key:       "WFInput",
				validType: String,
			},
		},
	},
	"copyToClipboard": {
		category:   Clipboard,
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
	"getClipboard": {
		category: Clipboard,
	},
	"getURLHeaders": {
		category:   URLs,
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
		category: URLs,
		parameters: []parameterDefinition{
			{
				name:      "url",
				validType: String,
				key:       "WFInput",
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"Show-WFInput": true,
			}
		},
	},
	"runJavaScriptOnWebpage": {
		category: Scripts,
		parameters: []parameterDefinition{
			{
				name:      "javascript",
				validType: String,
				key:       "WFJavaScript",
			},
		},
	},
	"searchWeb": {
		category: Web,
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
		category: Safari,
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
		category:   RSS,
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
		category:   RSS,
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
		category:   Safari,
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
		category:   Articles,
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
		category:   URLs,
		identifier: "safari.geturl",
	},
	"getWebpageContents": {
		category:   Safari,
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
		category:   Giphy,
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
		category:      Giphy,
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
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFGiphyShowPicker": false,
			}
		},
	},
	"getArticle": {
		category:   Articles,
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
		category:   URLs,
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
		category:   URLs,
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
		category:      HTTP,
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
				literal:   true,
			},
		},
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"WFHTTPMethod": "GET",
			}
		},
	},
	"formRequest": {
		category:   HTTP,
		identifier: "downloadurl",
		parameters: httpParams,
		addParams: func(args []actionArgument) map[string]any {
			return httpRequest("Form", "WFFormValues", args)
		},
	},
	"jsonRequest": {
		category:   HTTP,
		identifier: "downloadurl",
		parameters: httpParams,
		addParams: func(args []actionArgument) map[string]any {
			return httpRequest("JSON", "WFJSONValues", args)
		},
	},
	"fileRequest": {
		category:   HTTP,
		identifier: "downloadurl",
		parameters: httpParams,
		addParams: func(args []actionArgument) map[string]any {
			return httpRequest("File", "WFRequestVariable", args)
		},
	},
	"runAppleScript": {
		category: Scripts,
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
		category:   Scripts,
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
		category:   Windows,
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
		category: Windows,
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
		category: Windows,
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
		category:   Measurements,
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
		category:   Measurements,
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
		category:   BuiltIn,
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
				if reflect.TypeOf(image).String() != stringType && args[2].valueType == Variable {
					photo = fmt.Sprintf("{%s}", args[2].value)
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
		category:   BuiltIn,
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
		category:      Contacts,
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
		category:   Music,
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
		category:   Apps,
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
		category:   Shortcuts,
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
		category:   Apps,
		identifier: "filter.apps",
		mac:        true,
		minVersion: 18,
	},
	"setSoundRecognition": {
		category:      Settings,
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
		category:      Settings,
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

var plistTypes = map[string]string{"string": "", "integer": "", "boolean": "", "array": "", "dict": "", "real": ""}

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

			if len(args) > 1 {
				for _, parameterDefinitions := range getArgValue(args[1]).([]interface{}) {
					var definitions = parameterDefinitions.(map[string]interface{})
					if definitions["type"] != nil {
						var paramKey = definitions["key"]
						var paramType = definitions["type"].(string)
						var paramValue = definitions["value"].(string)
						if _, found := plistTypes[paramType]; !found {
							var list = makeKeyList("Available plist types:", plistTypes, paramValue)
							parserError(fmt.Sprintf("Raw action parameter '%s' type '%s' is not a plist type.\n\n%s", paramKey, paramType, list))
						}
					}
				}
			}
		},
		make: func(args []actionArgument) map[string]any {
			var params = make(map[string]any)
			if len(args) == 1 {
				return params
			}
			for _, parameterDefinitions := range getArgValue(args[1]).([]interface{}) {
				var paramKey string
				var paramType dataType
				var rawValue any
				for key, value := range parameterDefinitions.(map[string]interface{}) {
					switch key {
					case "key":
						paramKey = value.(string)
					case "type":
						paramType = dataType(value.(string))
					case "value":
						rawValue = value
					}
				}

				var tokenType = convertDataTypeToTokenType(paramType)
				params[paramKey] = paramValue(actionArgument{
					valueType: tokenType,
					value:     rawValue,
				}, tokenType)
			}

			return params
		},
		doc: ActionDoc{
			title:       "Raw Action",
			description: "Write a raw definition of an action not defined inside of Cherri, in Cherri.",
			example:     "rawAction(\"is.workflow.actions.gettext\", [\n\t{\n\t\t\"key\": \"WFTextActionText\",\n\t\t\"type\": \"string\",\n\t\t\"value\": \"Hello, world!\"\n\t}\n])",
		},
	}
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
		category:      Settings,
		doc: ActionDoc{
			title: "Background Sounds",
		},
	},
	"MediaBackgroundSounds": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleBackgroundSoundsIntent",
		addParams: func(_ []actionArgument) map[string]any {
			return map[string]any{
				"setting": "whenMediaIsPlaying",
			}
		},
		category: Settings,
		doc: ActionDoc{
			title: "Background Sounds (while media is playing)",
		},
	},
	"AutoAnswerCalls": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleAutoAnswerCallsIntent",
		category:      Settings,
		doc: ActionDoc{
			title: "Auto Answer Calls",
		},
	},
	"Appearance": {
		identifier: "appearance",
		category:   Settings,
	},
	"Bluetooth": {
		identifier: "bluetooth.set",
		setKey:     "OnValue",
		category:   Settings,
	},
	"Wifi": {
		identifier: "wifi.set",
		setKey:     "OnValue",
		category:   Settings,
	},
	"CellularData": {
		identifier: "cellulardata.set",
		setKey:     "OnValue",
		category:   Settings,
		doc: ActionDoc{
			title: "Cellular Data",
		},
	},
	"NightShift": {
		identifier: "nightshift.set",
		setKey:     "OnValue",
		category:   Settings,
		doc: ActionDoc{
			title: "Night Shift",
		},
	},
	"TrueTone": {
		identifier: "truetone.set",
		setKey:     "OnValue",
		category:   Settings,
		doc: ActionDoc{
			title: "True Tone",
		},
	},
	"ClassicInvert": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleClassicInvertIntent",
		category:      Settings,
		doc: ActionDoc{
			title: "Classic Invert",
		},
	},
	"ClosedCaptionsSDH": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleCaptionsIntent",
		category:      Settings,
		doc: ActionDoc{
			title: "Closed Captions SDH",
		},
	},
	"ColorFilters": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleColorFiltersIntent",
		category:      Settings,
		doc: ActionDoc{
			title: "Color Filters",
		},
	},
	"Contrast": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleContrastIntent",
		category:      Settings,
	},
	"LEDFlash": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleLEDFlashIntent",
		category:      Settings,
		doc: ActionDoc{
			title: "LED Flash",
		},
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
		category: Settings,
		doc: ActionDoc{
			title: "Left Right Balance",
		},
	},
	"LiveCaptions": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleLiveCaptionsIntent",
		category:      Settings,
		doc: ActionDoc{
			title: "Live Captions",
		},
	},
	"MonoAudio": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleMonoAudioIntent",
		category:      Settings,
		doc: ActionDoc{
			title: "Mono Audio",
		},
	},
	"ReduceMotion": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleReduceMotionIntent",
		category:      Settings,
		doc: ActionDoc{
			title: "Reduce Motion",
		},
	},
	"ReduceTransparency": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleTransparencyIntent",
		category:      Settings,
		doc: ActionDoc{
			title: "Reduce Transparency",
		},
	},
	"SmartInvert": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleSmartInvertIntent",
		category:      Settings,
		doc: ActionDoc{
			title: "Smart Invert",
		},
	},
	"SwitchControl": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleSwitchControlIntent",
		category:      Settings,
		doc: ActionDoc{
			title: "Switch Control",
		},
	},
	"VoiceControl": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleVoiceControlIntent",
		category:      Settings,
		doc: ActionDoc{
			title: "Voice Control",
		},
	},
	"WhitePoint": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleWhitePointIntent",
		category:      Settings,
		doc: ActionDoc{
			title: "White Point",
		},
	},
	"Zoom": {
		appIdentifier: "com.apple.AccessibilityUtilities.AXSettingsShortcuts",
		identifier:    "AXToggleZoomIntent",
		category:      Settings,
	},
}

// ToggleSetActions automates the creation of actions which simply toggle and set a state in the same format.
func makeToggleSetActions() {
	var emptyDoc ActionDoc
	for name, def := range toggleSetActions {
		var docTitle string
		if def.doc == emptyDoc {
			def.doc = ActionDoc{title: name}
			docTitle = name
		} else {
			docTitle = def.doc.title
		}

		var toggleName = fmt.Sprintf("toggle%s", name)
		var toggleDef = def
		toggleDef.addParams = func(_ []actionArgument) map[string]any {
			return map[string]any{
				"operation": "Toggle",
			}
		}
		toggleDef.parameters = nil

		toggleDef.doc.title = fmt.Sprintf("Toggle %s", docTitle)
		toggleDef.doc.description = fmt.Sprintf("Toggles %s on or off depending on its current state.", docTitle)

		actions[toggleName] = &toggleDef

		if name == "Appearance" {
			continue
		}

		var setDef = def
		var setName = fmt.Sprintf("set%s", name)
		setDef.addParams = nil
		var setKey = "state"
		if def.setKey != "" {
			setKey = def.setKey
		}
		setDef.parameters = append([]parameterDefinition{
			{
				name:      "status",
				validType: Bool,
				key:       setKey,
			},
		}, def.parameters...)

		setDef.doc.title = fmt.Sprintf("Set %s", docTitle)
		setDef.doc.description = fmt.Sprintf("Set %s on or off.", docTitle)
		setDef.doc.example = fmt.Sprintf("set%s(true)\nset%s(false)", name, name)

		actions[setName] = &setDef
	}
	toggleSetActions = nil
}
