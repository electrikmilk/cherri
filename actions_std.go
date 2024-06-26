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
		validType: Dict,
		optional:  true,
	},
	{
		name:      "headers",
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

var toggleAlarmIntent = appIntent{
	name:                "Clock",
	bundleIdentifier:    "com.apple.clock",
	appIntentIdentifier: "ToggleAlarmIntent",
}

// actions is the data structure that determines every action the compiler knows about.
// The key determines the identifier of the identifier that must be used in the syntax, it's value defines its behavior, etc. using an actionDefinition.
var actions = map[string]*actionDefinition{
	"date": {
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
	},
	"addMinutes": {
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
	},
	"addHours": {
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
	},
	"addDays": {
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
	},
	"addWeeks": {
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
	},
	"addMonths": {
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
	},
	"addYears": {
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
	},
	"subtractSeconds": {
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
	},
	"subtractMinutes": {
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
	},
	"subtractHours": {
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
	},
	"subtractDays": {
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
	},
	"subtractWeeks": {
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
	},
	"subtractMonths": {
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
	},
	"subtractYears": {
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
	},
	"getStartMinute": {
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
	},
	"getStartHour": {
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
	},
	"getStartWeek": {
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
	},
	"getStartMonth": {
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
	},
	"getStartYear": {
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
		appIdentifier: "com.apple.mobiletimer-framework.MobileTimerIntents.MTCreateAlarmIntent",
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
		appIdentifier: "com.apple.clock.DeleteAlarmIntent",
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
		appIdentifier: "com.apple.mobiletimer-framework.MobileTimerIntents.MTToggleAlarmIntent",
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
		appIdentifier: "com.apple.mobiletimer-framework.MobileTimerIntents.MTToggleAlarmIntent",
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
					value:    "toggle",
				},
			}
		},
	},
	"getAlarms": {
		appIdentifier: "com.apple.mobiletimer-framework.MobileTimerIntents.MTGetAlarmsIntent",
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
		appIdentifier: "com.apple.mobilephone.call",
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
		addParams: func(args []actionArgument) []plistData {
			var textArg = args[1]
			if textArg.valueType == Var {
				args[1] = actionArgument{
					valueType: String,
					value:     fmt.Sprintf("^{%s}", textArg.value),
				}
			} else {
				args[1].value = fmt.Sprintf("^%s", textArg.value)
			}

			return []plistData{
				argumentValue("WFMatchTextPattern", args, 1),
			}
		},
	},
	"matchText": {
		identifier: "text.match",
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
		appIdentifier: "com.apple.iBooksX.openin",
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
			},
			{
				name:      "to",
				validType: String,
				key:       "WFSelectedLanguage",
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			args[1].value = languageCode(args[1].value.(string))
			args[2].value = languageCode(args[2].value.(string))
		},
	},
	"translate": {
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
				optional:  true,
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
	},
	"iReplaceText": {
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
	},
	"regReplaceText": {
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
	},
	"iRegReplaceText": {
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
	},
	"uppercase": {
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
	},
	"lowercase": {
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
	},
	"titleCase": {
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
	},
	"capitalize": {
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
	},
	"capitalizeAll": {
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
	},
	"alternateCase": {
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
				validType: String,
			},
			{
				name:      "separator",
				validType: String,
			},
		},
		make: textParts,
	},
	"joinText": {
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
		appIdentifier: "com.apple.ShortcutsActions.TranscribeAudioAction",
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
		make: func(_ []actionArgument) []plistData {
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
				validType: String,
				key:       "WFVolume",
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			if args[0].valueType != Variable {
				args[0].value = fmt.Sprintf("0.%s", args[0].value)
			}
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
	},
	"encodeVideo": {
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
		identifier: "pausemusic",
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
		identifier: "reboot",
		minVersion: 17,
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
				validType: String,
				key:       "WFBrightness",
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			if args[0].valueType != Variable {
				args[0].value = fmt.Sprintf("0.%s", args[0].value)
			}
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
				validType: Variable,
			},
		},
		make: func(args []actionArgument) []plistData {
			return countParams("Items", args)
		},
	},
	"countChars": {
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
	},
	"countWords": {
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
	},
	"countSentences": {
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
	},
	"countLines": {
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
	},
	"toggleAppearance": {
		identifier: "appearance",
		make: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "operation",
					dataType: Text,
					value:    "toggle",
				},
			}
		},
	},
	"lightMode": {
		identifier: "appearance",
		make: func(_ []actionArgument) []plistData {
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
		identifier: "appearance",
		make: func(_ []actionArgument) []plistData {
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
	"getBatteryLevel": {},
	"isCharging": {
		identifier: "getbatterylevel",
		minVersion: 16.2,
		make: func(_ []actionArgument) []plistData {
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
		make: func(_ []actionArgument) []plistData {
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
	},
	"base64Decode": {
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
				validType: String,
				key:       "WFDictionaryValue",
			},
		},
	},
	"openApp": {
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
	},
	"hideApp": {
		identifier: "hide.app",
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
	},
	"quitApp": {
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
	},
	"openShortcut": {
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
			}
		},
	},
	"runSelf": {
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
				argumentValue("WFInput", args, 0),
			}
		},
	},
	"run": {
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
					value:    "External",
				},
			}
		},
	},
	"getFirstItem": {
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
	},
	"getLastItem": {
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
	},
	"getRandomItem": {
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
		identifier: "posters.get",
		minVersion: 16.2,
	},
	"getWallpaper": {
		identifier: "posters.get",
		minVersion: 16.2,
		make: func(_ []actionArgument) []plistData {
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
		identifier: "openxcallbackurl",
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
		make: func(_ []actionArgument) []plistData {
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
		identifier: "dnd.set",
		make: func(_ []actionArgument) []plistData {
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
	"setWifi": {
		identifier: "wifi.set",
		parameters: []parameterDefinition{
			{
				name:      "status",
				key:       "OnValue",
				validType: Bool,
			},
		},
	},
	"setCellularData": {
		identifier: "cellulardata.set",
		parameters: []parameterDefinition{
			{
				name:         "status",
				key:          "OnValue",
				validType:    Bool,
				defaultValue: true,
			},
		},
	},
	"toggleBluetooth": {
		identifier: "bluetooth.set",
		make: func(_ []actionArgument) []plistData {
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
	},
	"setBluetooth": {
		identifier: "bluetooth.set",
		parameters: []parameterDefinition{
			{
				name:      "status",
				validType: Bool,
				key:       "OnValue",
			},
		},
		addParams: func(_ []actionArgument) []plistData {
			return []plistData{
				{
					key:      "operation",
					dataType: Text,
					value:    "set",
				},
			}
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
	},
	"ceil": {
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
	},
	"floor": {
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
	},
	"searchShortcuts": {
		appIdentifier: "com.apple.shortcuts.SearchShortcutsAction",
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
		make: func(args []actionArgument) []plistData {
			return httpRequest("Form", "WFFormValues", args)
		},
	},
	"jsonRequest": {
		identifier: "downloadurl",
		parameters: httpParams,
		make: func(args []actionArgument) []plistData {
			return httpRequest("JSON", "WFJSONValues", args)
		},
	},
	"fileRequest": {
		identifier: "downloadurl",
		parameters: httpParams,
		make: func(args []actionArgument) []plistData {
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
				validType: Arr,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
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
		argumentValue("Input", args, 0),
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
