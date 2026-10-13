/*
 * Copyright (c) Cherri
 */

package main

import "strings"

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

// menuBoxArgSet reports whether an optional argument was actually provided, treating a skipped
// position (nil) the same as an omitted one so its key is left out of the action.
func menuBoxArgSet(args []actionArgument, idx int) bool {
	return len(args) > idx && args[idx].valueType != Nil && args[idx].value != nil
}

// menuBoxSymbolEntity builds the SF Symbol entity Menu Box expects, which is identified by the
// symbol's name.
func menuBoxSymbolEntity(name string) map[string]any {
	return map[string]any{
		"title":      map[string]string{"key": name},
		"subtitle":   map[string]string{"key": name},
		"symbol":     map[string]string{"systemName": name},
		"identifier": name,
	}
}

// menuBoxWeightEntity builds the symbol weight entity. Shortcuts identifies a weight by its
// lowercase name, shows the capitalized name, and attaches the chosen symbol as its preview.
func menuBoxWeightEntity(weight string, symbol string) map[string]any {
	var title = strings.ToUpper(weight[:1]) + weight[1:]

	return map[string]any{
		"title":      map[string]string{"key": title},
		"subtitle":   map[string]string{"key": title},
		"symbol":     map[string]string{"systemName": symbol},
		"identifier": weight,
	}
}

func init() {
	enumerations["menuBoxSymbolWeight"] = []string{"ultralight", "thin", "light", "regular", "medium", "semibold", "bold", "heavy", "black"}
	enumerations["menuBoxSymbolScale"] = []string{"small", "medium", "large"}

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

	actions["createMenuItemFromSymbol"] = &actionDefinition{
		doc: selfDoc{
			title: "Create Menu Item From Symbol",
			description: "Creates a menu item with an SF Symbol icon, for use with Choose From List or Add To Menu Set. " +
				"The symbol and weight are literal values. From the Menu Box app.",
			category:    "thirdparty",
			subcategory: "Menu Box",
		},
		appIdentifier: menuBoxBundleIdentifier,
		identifier:    "CreateSymbolItemIntent",
		appIntent:     menuBoxIntent("CreateSymbolItemIntent"),
		nonMacOnly:    true,
		parameters: []parameterDefinition{
			{
				name:      "title",
				validType: String,
				key:       "title",
			},
			{
				name:      "symbol",
				validType: String,
				literal:   true,
			},
			{
				name:      "subtitle",
				validType: String,
				key:       "subtitle",
				optional:  true,
			},
			{
				name:      "weight",
				validType: String,
				enum:      "menuBoxSymbolWeight",
				literal:   true,
				optional:  true,
			},
			{
				name:      "scale",
				validType: String,
				key:       "scale",
				enum:      "menuBoxSymbolScale",
				literal:   true,
				optional:  true,
			},
			{
				name:      "transparentBackground",
				validType: Bool,
				key:       "transparentBG",
				optional:  true,
			},
		},
		makeParams: func(args []actionArgument) map[string]any {
			var symbol = argumentValue(args, 1).(string)
			var params = map[string]any{
				"title":  argumentValue(args, 0),
				"symbol": menuBoxSymbolEntity(symbol),
			}
			if menuBoxArgSet(args, 2) {
				params["subtitle"] = argumentValue(args, 2)
			}
			if menuBoxArgSet(args, 3) {
				params["weight"] = menuBoxWeightEntity(argumentValue(args, 3).(string), symbol)
			}
			if menuBoxArgSet(args, 4) {
				params["scale"] = argumentValue(args, 4)
			}
			if menuBoxArgSet(args, 5) {
				params["transparentBG"] = argumentValue(args, 5)
			}

			return params
		},
		outputType: Variable,
		outputName: "Menu Item",
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
