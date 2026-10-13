/*
 * Copyright (c) Cherri
 */

package main

// Actions from Menu Box (dev.alexhay.MenuBox) — a third-party App Intents extension for
// building custom menus. Like the Actions app, every intent shares one signing team and
// bundle; only the AppIntentIdentifier and parameters vary per action.
//
// Menu Box only exists on iOS and iPadOS, so every definition here is nonMacOnly.
const menuBoxTeamIdentifier = "QRVYSKD5TN"
const menuBoxBundleIdentifier = "dev.alexhay.MenuBox"
const menuBoxAppName = "Menu Box"

// menuBoxIntent builds the appIntent descriptor shared by every Menu Box intent, varying only
// by its AppIntentIdentifier.
func menuBoxIntent(appIntentIdentifier string) appIntent {
	return appIntent{
		name:                menuBoxAppName,
		bundleIdentifier:    menuBoxBundleIdentifier,
		appIntentIdentifier: appIntentIdentifier,
		teamIdentifier:      menuBoxTeamIdentifier,
	}
}

func init() {
	actions["quickMenu"] = &actionDefinition{
		doc: selfDoc{
			title: "Quick Menu",
			description: "Builds menu items from marked-up text, for use with Choose From List or Add To Menu Set. " +
				"Each item is a block of `Key: value` lines separated by a blank line. `Title` is required; " +
				"`Sub`, `Icon`, `ID`, `B64`, `Mask`, `Scale` and `Field1`–`Field4` are optional. From the Menu Box app.",
			category:    "thirdparty",
			subcategory: "Menu Box",
		},
		appIdentifier: menuBoxBundleIdentifier,
		identifier:    "QuickMenuIntent",
		appIntent:     menuBoxIntent("QuickMenuIntent"),
		nonMacOnly:    true,
		parameters: []parameterDefinition{
			{
				name:      "menu",
				key:       "text",
				validType: String,
			},
		},
		outputType: Variable,
		outputName: "Quick Menu",
	}
}
