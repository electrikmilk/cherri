/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"regexp"
	"strings"
)

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

// menuBoxColourPreset is one entry of the native colour list Menu Box offers for symbol and
// background colours.
type menuBoxColourPreset struct {
	name string
	hex  string
}

var menuBoxColourPresets = []menuBoxColourPreset{
	{"Yellow", "#ffcc01"},
	{"Pink", "#ff2d55"},
	{"White", "#ffffff"},
	{"Indigo", "#5856d7"},
	{"Purple", "#af52de"},
	{"Grey", "#808080"},
	{"Red", "#ff3a2f"},
	{"Green", "#33c759"},
	{"Orange", "#ff9502"},
	{"Blue", "#017aff"},
	{"Black", "#000000"},
	{"Teal", "#5ac8fa"},
}

// menuBoxHexColour matches a six digit hex colour, which Menu Box reads with or without the
// leading #. The value is passed through exactly as written.
var menuBoxHexColour = regexp.MustCompile(`^#?[0-9a-fA-F]{6}$`)

// menuBoxFindColourPreset looks a preset up by name, ignoring case. Menu Box's own docs accept
// "Gray" as well as "Grey" for its Quick Menu colour names, so both resolve to the Grey preset.
func menuBoxFindColourPreset(name string) (menuBoxColourPreset, bool) {
	if strings.EqualFold(name, "gray") {
		name = "grey"
	}
	for _, preset := range menuBoxColourPresets {
		if strings.EqualFold(preset.name, name) {
			return preset, true
		}
	}

	return menuBoxColourPreset{}, false
}

// menuBoxApplyColour writes one colour slot (symbol or background) as a hex. A preset name is
// written as that preset's hex because Menu Box ignores its own preset entity: with both hex fields
// empty, every preset in the picker rendered the default blue, even in a shortcut built by hand in
// the Shortcuts app. The same colour as a hex renders correctly, so the hex is the path that works.
// Anything else is already a hex and is passed through exactly as written.
func menuBoxApplyColour(params map[string]any, value string, hexKey string) {
	if preset, found := menuBoxFindColourPreset(value); found {
		params[hexKey] = strings.TrimPrefix(preset.hex, "#")

		return
	}

	params[hexKey] = value
}

// menuBoxCheckColour rejects a colour that is neither a preset name nor a six digit hex, since
// Menu Box would silently ignore it.
func menuBoxCheckColour(args []actionArgument, idx int, label string) {
	if !menuBoxArgSet(args, idx) {
		return
	}

	var value, isString = getArgValue(args[idx]).(string)
	if !isString {
		return
	}
	if _, found := menuBoxFindColourPreset(value); found || menuBoxHexColour.MatchString(value) {
		return
	}

	var names []string
	for _, preset := range menuBoxColourPresets {
		names = append(names, preset.name)
	}
	parserError(fmt.Sprintf("Invalid %s '%s'. Use a six digit hex such as 00ff00 or one of: %s", label, value, strings.Join(names, ", ")))
}

// menuBoxApplyDataFields writes the optional data fields (Field 1 to 4). Add Data is the switch that
// reveals them in the editor, so it is written whenever any field is, as Shortcuts does.
func menuBoxApplyDataFields(params map[string]any, args []actionArgument, firstIdx int) {
	for offset, key := range []string{"f1", "f2", "f3", "f4"} {
		if menuBoxArgSet(args, firstIdx+offset) {
			params["addData"] = true
			params[key] = argumentValue(args, firstIdx+offset)
		}
	}
}

func menuBoxDataFieldParameters() []parameterDefinition {
	var fields []parameterDefinition
	for _, name := range []string{"field1", "field2", "field3", "field4"} {
		fields = append(fields, parameterDefinition{
			name:      name,
			validType: String,
			optional:  true,
		})
	}

	return fields
}

func init() {
	enumerations["menuBoxSymbolWeight"] = []string{"ultralight", "thin", "light", "regular", "medium", "semibold", "bold", "heavy", "black"}
	enumerations["menuBoxIconScale"] = []string{"small", "medium", "large"}
	// "none" is the default mask and is only written when set explicitly, which is how its raw
	// value was seen in a device export (Create Menu Item From Emoji).
	enumerations["menuBoxIconMask"] = []string{"none", "circle", "roundRect"}

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
				"The symbol, weight, scale, mask and colours are literal values. A colour is either the name of one of " +
				"Menu Box's preset colours (Yellow, Pink, White, Indigo, Purple, Grey, Red, Green, Orange, Blue, Black, Teal) " +
				"or a six digit hex with or without the leading `#`. A preset is written as its hex. With no colours " +
				"the item is a blue symbol on a white rounded square, and a symbol colour alone also sits on white. " +
				"A background colour is hidden while " +
				"`transparentBackground` is on. Turn `renderIcon` off to skip rendering the icon until the item is added " +
				"to a Menu Set, which speeds up the shortcut when you are building a set. From the Menu Box app.",
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
				enum:      "menuBoxIconScale",
				literal:   true,
				optional:  true,
			},
			{
				name:      "transparentBackground",
				validType: Bool,
				key:       "transparentBG",
				optional:  true,
			},
			{
				name:      "mask",
				validType: String,
				key:       "mask",
				enum:      "menuBoxIconMask",
				literal:   true,
				optional:  true,
			},
			{
				name:      "symbolColour",
				validType: String,
				literal:   true,
				optional:  true,
			},
			{
				name:      "backgroundColour",
				validType: String,
				literal:   true,
				optional:  true,
			},
			{
				name:      "renderIcon",
				validType: Bool,
				key:       "renderIcon",
				optional:  true,
			},
			{name: "field1", validType: String, optional: true},
			{name: "field2", validType: String, optional: true},
			{name: "field3", validType: String, optional: true},
			{name: "field4", validType: String, optional: true},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			menuBoxCheckColour(args, 7, "symbol colour")
			menuBoxCheckColour(args, 8, "background colour")
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
			if menuBoxArgSet(args, 6) {
				params["mask"] = argumentValue(args, 6)
			}

			// Setting either colour switches Menu Box to its custom colours. A slot left unset is
			// written as an empty hex, the shape Shortcuts produces once the pre-filled hex is cleared.
			var symbolColourSet = menuBoxArgSet(args, 7)
			var backgroundColourSet = menuBoxArgSet(args, 8)
			if symbolColourSet || backgroundColourSet {
				params["colourChoices"] = "custom"
				params["symbolColourHex"] = ""
				params["bgColourHex"] = ""
			}
			if symbolColourSet {
				menuBoxApplyColour(params, argumentValue(args, 7).(string), "symbolColourHex")
			}
			if backgroundColourSet {
				menuBoxApplyColour(params, argumentValue(args, 8).(string), "bgColourHex")
			}
			if menuBoxArgSet(args, 9) {
				params["renderIcon"] = argumentValue(args, 9)
			}
			menuBoxApplyDataFields(params, args, 10)

			return params
		},
		outputType: Variable,
		outputName: "Menu Item",
	}

	actions["createMenuItemFromEmoji"] = &actionDefinition{
		doc: selfDoc{
			title: "Create Menu Item From Emoji",
			description: "Creates a menu item with an emoji icon, for use with Choose From List or Add To Menu Set. " +
				"The emoji should be a single emoji, letter or digit. Transparent Background is on by default for " +
				"emoji, so a `backgroundColour` (a six digit hex with or without the leading `#`) is ignored until " +
				"`transparentBackground` is turned off. The mask defaults to none. Fields 1 to 4 hold data that " +
				"isn't shown in the menu but can be read from the chosen item. Turn `renderIcon` off to skip " +
				"rendering the icon until the item is added to a Menu Set. From the Menu Box app.",
			category:    "thirdparty",
			subcategory: "Menu Box",
		},
		appIdentifier: menuBoxBundleIdentifier,
		identifier:    "CreateEmojiItemIntent",
		appIntent:     menuBoxIntent("CreateEmojiItemIntent"),
		nonMacOnly:    true,
		parameters: append([]parameterDefinition{
			{
				name:      "title",
				validType: String,
				key:       "title",
			},
			{
				name:      "emoji",
				validType: String,
				key:       "emoji",
			},
			{
				name:      "subtitle",
				validType: String,
				key:       "subtitle",
				optional:  true,
			},
			{
				name:      "scale",
				validType: String,
				key:       "scale",
				enum:      "menuBoxIconScale",
				literal:   true,
				optional:  true,
			},
			{
				name:      "transparentBackground",
				validType: Bool,
				key:       "transparentBG",
				optional:  true,
			},
			{
				name:      "mask",
				validType: String,
				key:       "mask",
				enum:      "menuBoxIconMask",
				literal:   true,
				optional:  true,
			},
			{
				name:      "backgroundColour",
				validType: String,
				literal:   true,
				optional:  true,
			},
			{
				name:      "renderIcon",
				validType: Bool,
				key:       "renderIcon",
				optional:  true,
			},
		}, menuBoxDataFieldParameters()...),
		check: func(args []actionArgument, _ *actionDefinition) {
			menuBoxCheckColour(args, 6, "background colour")
		},
		makeParams: func(args []actionArgument) map[string]any {
			var params = map[string]any{
				"title": argumentValue(args, 0),
				"emoji": argumentValue(args, 1),
			}
			if menuBoxArgSet(args, 2) {
				params["subtitle"] = argumentValue(args, 2)
			}
			if menuBoxArgSet(args, 3) {
				params["scale"] = argumentValue(args, 3)
			}
			if menuBoxArgSet(args, 4) {
				params["transparentBG"] = argumentValue(args, 4)
			}
			if menuBoxArgSet(args, 5) {
				params["mask"] = argumentValue(args, 5)
			}

			// A background colour switches Menu Box to its custom colours, the only part of them an
			// emoji has. Emoji items have no colour picker, so a preset name is written as its hex.
			if menuBoxArgSet(args, 6) {
				params["colourChoice"] = "custom"
				menuBoxApplyColour(params, argumentValue(args, 6).(string), "bgColourHex")
			}
			if menuBoxArgSet(args, 7) {
				params["renderIcon"] = argumentValue(args, 7)
			}
			menuBoxApplyDataFields(params, args, 8)

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
