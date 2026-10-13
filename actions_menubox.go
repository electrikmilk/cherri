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

// menuBoxSetEntity builds the menu set entity Menu Box expects. Shortcuts resolves it by
// identifier (the set's UUID, which is device-specific); the title and subtitle are only the
// name shown on the parameter chip in the editor.
func menuBoxSetEntity(name string, identifier string) map[string]any {
	return map[string]any{
		"title":      map[string]string{"key": name},
		"subtitle":   map[string]string{"key": name},
		"identifier": identifier,
	}
}

func init() {
	actions["loadMenuSet"] = &actionDefinition{
		doc: selfDoc{
			title: "Load Menu Set",
			description: "Loads the menu items of a saved menu set, for use with Choose From List. The set is " +
				"identified by its name and its set ID (long-press the set in Menu Box to copy the ID). " +
				"Set IDs belong to one Menu Box library, so a literal ID will not resolve on another person's device. " +
				"From the Menu Box app.",
			category:    "thirdparty",
			subcategory: "Menu Box",
		},
		appIdentifier: menuBoxBundleIdentifier,
		identifier:    "ChooseMenuSetIntent",
		appIntent:     menuBoxIntent("ChooseMenuSetIntent"),
		nonMacOnly:    true,
		parameters: []parameterDefinition{
			{
				name:      "setName",
				validType: String,
				literal:   true,
			},
			{
				name:      "setID",
				validType: String,
				literal:   true,
			},
			{
				name:      "loadAlternate",
				validType: Bool,
				key:       "loadAlt",
				optional:  true,
			},
		},
		makeParams: func(args []actionArgument) map[string]any {
			var params = map[string]any{
				"menuSet": menuBoxSetEntity(argumentValue(args, 0).(string), argumentValue(args, 1).(string)),
			}
			// Shortcuts only writes loadAlt once the toggle has been touched, so an omitted argument
			// leaves the key out entirely.
			if len(args) > 2 {
				params["loadAlt"] = argumentValue(args, 2)
			}

			return params
		},
		outputType: Variable,
		outputName: "Menu Items",
	}

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
