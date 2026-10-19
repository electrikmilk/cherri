/*
 * Copyright (c) Cherri
 */

package main

// Actions by @sindresorhus (com.sindresorhus.Actions) — a third-party App Intents
// extension. Every intent shares the same signing team, bundle, and display name;
// only the AppIntentIdentifier and parameters vary per action.
//
// Prior to this file, these were only reachable through rawAction(), which has no
// declared parameters/output and forces every call site to hand-build the
// AppIntentDescriptor dict and cast the result before use. These definitions give
// each intent real Cherri syntax, type-checked parameters, and a proper output
// label, matching how first-party intents (see toggleAlarmIntent, etc. above) are
// defined.
const sindresorhusTeamIdentifier = "YG56YK5RN5"
const sindresorhusBundleIdentifier = "com.sindresorhus.Actions"
const sindresorhusAppName = "Actions"

// sindresorhusIntent builds the appIntent descriptor shared by every Actions app
// intent, varying only by its AppIntentIdentifier.
func sindresorhusIntent(appIntentIdentifier string) appIntent {
	return appIntent{
		name:                sindresorhusAppName,
		bundleIdentifier:    sindresorhusBundleIdentifier,
		appIntentIdentifier: appIntentIdentifier,
		teamIdentifier:      sindresorhusTeamIdentifier,
	}
}

func init() {
	actions["generateUUID"] = &actionDefinition{
		doc: selfDoc{
			title:       "Generate UUID",
			description: "Generates a random UUID string. From the Actions app by @sindresorhus.",
			category:    "thirdparty",
			subcategory: "Actions App",
		},
		appIdentifier: sindresorhusBundleIdentifier,
		identifier:    "GenerateUUIDIntent",
		appIntent:     sindresorhusIntent("GenerateUUIDIntent"),
		outputType:    Variable,
		outputName:    "UUID",
	}

	actions["randomBoolean"] = &actionDefinition{
		doc: selfDoc{
			title:       "Random Boolean",
			description: "Generates a random true/false value. From the Actions app by @sindresorhus.",
			category:    "thirdparty",
			subcategory: "Actions App",
		},
		appIdentifier: sindresorhusBundleIdentifier,
		identifier:    "RandomBooleanIntent",
		appIntent:     sindresorhusIntent("RandomBooleanIntent"),
		outputType:    Variable,
		outputName:    "Random Boolean",
	}

	// Named isDeviceOnline, not isOnline — that name is already taken by the
	// built-in network check (actions/network.cherri, is.workflow.actions.getipaddress).
	actions["isDeviceOnline"] = &actionDefinition{
		doc: selfDoc{
			title:       "Is Online",
			description: "Checks whether the device currently has an internet connection. From the Actions app by @sindresorhus.",
			category:    "thirdparty",
			subcategory: "Actions App",
		},
		appIdentifier: sindresorhusBundleIdentifier,
		identifier:    "IsOnlineIntent",
		appIntent:     sindresorhusIntent("IsOnlineIntent"),
		outputType:    Variable,
		outputName:    "Is Online",
	}

	actions["getBatteryState"] = &actionDefinition{
		doc: selfDoc{
			title:       "Get Battery State",
			description: "Returns the device's current battery state. From the Actions app by @sindresorhus.",
			category:    "thirdparty",
			subcategory: "Actions App",
		},
		appIdentifier: sindresorhusBundleIdentifier,
		identifier:    "GetBatteryStateIntent",
		appIntent:     sindresorhusIntent("GetBatteryStateIntent"),
		outputType:    Variable,
		outputName:    "Is Battery State",
	}

	actions["formatNumberCompact"] = &actionDefinition{
		doc: selfDoc{
			title:       "Format Number (Compact)",
			description: "Formats a number using compact notation (e.g. 1.2K). From the Actions app by @sindresorhus.",
			category:    "thirdparty",
			subcategory: "Actions App",
		},
		appIdentifier: sindresorhusBundleIdentifier,
		identifier:    "FormatNumberCompactIntent",
		appIntent:     sindresorhusIntent("FormatNumberCompactIntent"),
		parameters: []parameterDefinition{
			{
				name:      "number",
				key:       "number",
				validType: Integer,
			},
		},
		outputType: Variable,
		outputName: "Formatted Number",
	}

	actions["transformText"] = &actionDefinition{
		doc: selfDoc{
			title:       "Transform Text",
			description: "Transforms text (e.g. slugify, uppercase). From the Actions app by @sindresorhus.",
			category:    "thirdparty",
			subcategory: "Actions App",
		},
		appIdentifier: sindresorhusBundleIdentifier,
		identifier:    "TransformTextIntent",
		appIntent:     sindresorhusIntent("TransformTextIntent"),
		parameters: []parameterDefinition{
			{
				name:      "transformation",
				key:       "transformation",
				validType: String,
			},
			{
				name:      "text",
				key:       "text",
				validType: String,
			},
		},
		outputType: Variable,
		outputName: "Transformed Text",
	}

	actions["truncateNumber"] = &actionDefinition{
		doc: selfDoc{
			title:       "Truncate Number",
			description: "Truncates a number to a given number of decimal places. From the Actions app by @sindresorhus.",
			category:    "thirdparty",
			subcategory: "Actions App",
		},
		appIdentifier: sindresorhusBundleIdentifier,
		// Unlike the other intents in this file, the real AppIntentIdentifier has no
		// "Intent" suffix — verified against the reference plist.
		identifier: "TruncateNumber",
		appIntent:  sindresorhusIntent("TruncateNumber"),
		parameters: []parameterDefinition{
			{
				name:      "number",
				key:       "number",
				validType: Integer,
			},
			{
				name:      "decimalPlaces",
				key:       "decimalPlaces",
				validType: Integer,
			},
		},
		outputType: Variable,
		outputName: "Truncated Number",
	}

	actions["randomText"] = &actionDefinition{
		doc: selfDoc{
			title:       "Random Text",
			description: "Generates random text. From the Actions app by @sindresorhus.",
			category:    "thirdparty",
			subcategory: "Actions App",
		},
		appIdentifier: sindresorhusBundleIdentifier,
		identifier:    "RandomTextIntent",
		appIntent:     sindresorhusIntent("RandomTextIntent"),
		outputType:    Variable,
		outputName:    "Random Text",
	}

	// Named trimTextWhitespace, not trimWhitespace — that name is already taken by
	// the built-in text action (actions/text.cherri, is.workflow.actions.text.trimwhitespace).
	actions["trimTextWhitespace"] = &actionDefinition{
		doc: selfDoc{
			title:       "Trim Whitespace",
			description: "Trims leading and trailing whitespace from text. From the Actions app by @sindresorhus.",
			category:    "thirdparty",
			subcategory: "Actions App",
		},
		appIdentifier: sindresorhusBundleIdentifier,
		identifier:    "TrimWhitespaceIntent",
		appIntent:     sindresorhusIntent("TrimWhitespaceIntent"),
		parameters: []parameterDefinition{
			{
				name:      "text",
				key:       "text",
				validType: String,
			},
		},
		outputType: Variable,
		outputName: "Text with Trimmed Whitespace",
	}

	actions["formatTextList"] = &actionDefinition{
		doc: selfDoc{
			title:       "Format Text List",
			description: "Formats a list of text items into a single readable string (e.g. \"a, b, and c\"). From the Actions app by @sindresorhus.",
			category:    "thirdparty",
			subcategory: "Actions App",
		},
		appIdentifier: sindresorhusBundleIdentifier,
		identifier:    "FormatTextListIntent",
		appIntent:     sindresorhusIntent("FormatTextListIntent"),
		parameters: []parameterDefinition{
			{
				name:      "items",
				key:       "items",
				validType: String,
			},
		},
		outputType: Variable,
		outputName: "Formatted Text",
	}
}
