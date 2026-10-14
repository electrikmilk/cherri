/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"regexp"
	"slices"
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

// menuBoxQueryIntent builds the descriptor for an entity query such as Find Menu Sets. Shortcuts marks
// these as requiring the app to be installed.
func menuBoxQueryIntent(appIntentIdentifier string) appIntent {
	var intent = menuBoxIntent(appIntentIdentifier)
	intent.requiresAppInstallation = true

	return intent
}

// The numbers Shortcuts writes for each filter operator, all read from device exports. Dates and counts
// share one comparison scheme (0 is less, 1 less or equal, 2 greater, 3 greater or equal, 4 equal, 5 not
// equal); text has its own, with contains at 99. Date Modified reuses the date numbers, which the export
// showed for four of its six operators.
var menuBoxNameOperators = map[string]int{
	"is":         4,
	"isNot":      5,
	"contains":   99,
	"beginsWith": 8,
	"endsWith":   9,
}

var menuBoxDateOperators = map[string]int{
	"exactly":    4,
	"notExactly": 5,
	"after":      2,
	"onOrAfter":  3,
	"before":     0,
	"onOrBefore": 1,
}

var menuBoxCountOperators = map[string]int{
	"is":                 4,
	"isNot":              5,
	"greaterThan":        2,
	"greaterThanOrEqual": 3,
	"lessThan":           0,
	"lessThanOrEqual":    1,
}

// menuBoxFindOrders maps the order names used in Cherri to the display strings Shortcuts stores for them.
// The labels depend on the sort: Name offers A to Z and Z to A, the date sorts offer Latest First and
// Oldest First (see menuBoxFindOrderSort).
var menuBoxFindOrders = map[string]string{
	"latestFirst": "Latest First",
	"oldestFirst": "Oldest First",
	"aToZ":        "A to Z",
	"zToA":        "Z to A",
}

// menuBoxFindOrderSort records which sorts each order belongs to, because Shortcuts only offers the
// matching pair for the chosen sort.
var menuBoxFindOrderSort = map[string][]string{
	"latestFirst": {"dateCreated", "dateModified"},
	"oldestFirst": {"dateCreated", "dateModified"},
	"aToZ":        {"name"},
	"zToA":        {"name"},
}

// menuBoxFindOrdersFor lists the orders that apply to a sort, in a stable order for error messages.
func menuBoxFindOrdersFor(sortBy string) []string {
	var orders []string
	for order, sorts := range menuBoxFindOrderSort {
		if slices.Contains(sorts, sortBy) {
			orders = append(orders, order)
		}
	}
	slices.Sort(orders)

	return orders
}

// Argument positions of Find Menu Sets. Each filter sits next to its operator.
const (
	menuBoxFindArgName = iota
	menuBoxFindArgNameOperator
	menuBoxFindArgItemTitle
	menuBoxFindArgCreated
	menuBoxFindArgCreatedOperator
	menuBoxFindArgModified
	menuBoxFindArgModifiedOperator
	menuBoxFindArgItemCount
	menuBoxFindArgItemCountOperator
	menuBoxFindArgMatchAny
	menuBoxFindArgSortBy
	menuBoxFindArgOrder
	menuBoxFindArgLimit
	menuBoxFindArgInput
)

// menuBoxFindFilter describes one row of the Find filter: the property Shortcuts names it by, the
// arguments holding its value and operator, and the key its value is stored under.
type menuBoxFindFilter struct {
	property        string
	valueKey        string
	valueArg        int
	valueName       string
	operatorArg     int
	operatorName    string
	operators       map[string]int
	defaultOperator string
}

// menuBoxFindFilters lists the filter rows in the order they are written. Menu Item Titles only offers
// contains, so it has no operator argument.
var menuBoxFindFilters = []menuBoxFindFilter{
	{"name", "String", menuBoxFindArgName, "name", menuBoxFindArgNameOperator, "nameOperator", menuBoxNameOperators, "contains"},
	{"menuItemTitles", "String", menuBoxFindArgItemTitle, "itemTitle", -1, "", menuBoxNameOperators, "contains"},
	{"dateCreated", "Date", menuBoxFindArgCreated, "createdDate", menuBoxFindArgCreatedOperator, "createdOperator", menuBoxDateOperators, "exactly"},
	{"dateModified", "Date", menuBoxFindArgModified, "modifiedDate", menuBoxFindArgModifiedOperator, "modifiedOperator", menuBoxDateOperators, "exactly"},
	{"itemCount", "Number", menuBoxFindArgItemCount, "itemCount", menuBoxFindArgItemCountOperator, "itemCountOperator", menuBoxCountOperators, "is"},
}

// menuBoxFilterTemplate builds one row of a Find filter. Removable is always true for rows added in the editor.
func menuBoxFilterTemplate(property string, operator int, values map[string]any) map[string]any {
	return map[string]any{
		"Property":  property,
		"Operator":  operator,
		"Values":    values,
		"Removable": true,
	}
}

// menuBoxExistingSetEntity builds the entity Add To Menu Set uses to point at a set that already
// exists. Unlike menuBoxSetEntity it has no identifier: the set is found by its name alone, carried
// in "value", so it resolves on any device that has a set with that name.
func menuBoxExistingSetEntity(name string) map[string]any {
	return map[string]any{
		"value":    name,
		"title":    map[string]string{"key": name},
		"subtitle": map[string]string{"key": name},
	}
}

// menuBoxTagEntity builds the tag entity Menu Box expects. Like a set, a tag is resolved by its
// identifier (a UUID that belongs to one library); the symbol is always the tag glyph.
func menuBoxTagEntity(name string, identifier string) map[string]any {
	return map[string]any{
		"title":      map[string]string{"key": name},
		"subtitle":   map[string]string{"key": name},
		"symbol":     map[string]string{"systemName": "tag"},
		"identifier": identifier,
	}
}

// menuBoxCheckPair rejects half of a name and ID pair. An entity needs both to resolve, so giving
// only one is always a mistake.
func menuBoxCheckPair(args []actionArgument, firstIdx int, label string) {
	if menuBoxArgSet(args, firstIdx) != menuBoxArgSet(args, firstIdx+1) {
		parserError(fmt.Sprintf("The %s needs both its name and its ID, or neither.", label))
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

// menuBoxApplyDataFields writes the optional data fields (Field 1 to 4) under the given keys. Add Data
// is the switch that reveals them in the editor, so it is written whenever any field is, as Shortcuts
// does. The keys differ per intent: the emoji action uses f1 to f4 and the symbol action field1 to
// field4, so they are passed in rather than assumed to match.
func menuBoxApplyDataFields(params map[string]any, args []actionArgument, firstIdx int, keys []string) {
	for offset, key := range keys {
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

// menuBoxImageParameters lists the parameters the two image actions share, in call order. Only the image
// source differs between them: Create Menu Item From Image takes a file and Create Menu Item From Base64
// takes a base64 string, which are two modes of the one Shortcuts action (Image Type) with mutually
// exclusive inputs. Cherri can't tell a file variable from a text variable at compile time, so each
// mode is its own action and the mode is implied rather than passed.
func menuBoxImageParameters(source parameterDefinition) []parameterDefinition {
	return append([]parameterDefinition{
		{
			name:      "title",
			validType: String,
			key:       "title",
		},
		source,
		{
			name:      "subtitle",
			validType: String,
			key:       "subtitle",
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
	}, menuBoxDataFieldParameters()...)
}

// menuBoxMakeImageParams builds the parameters for either image action. The image is used as it is, so
// unlike the symbol and emoji actions there is no scale, background or render option, only the mask.
// Image Type is always written: Shortcuts leaves it out while it is untouched, but it was seen as
// "file" and "base64" in device exports, so writing it doesn't rely on Menu Box's default. The data
// fields use field1 to field4 (as the symbol action does, not the emoji action's f1 to f4).
func menuBoxMakeImageParams(args []actionArgument, sourceKey string, imageType string) map[string]any {
	var params = map[string]any{
		"title":     argumentValue(args, 0),
		sourceKey:   argumentValue(args, 1),
		"imageType": imageType,
	}
	if menuBoxArgSet(args, 2) {
		params["subtitle"] = argumentValue(args, 2)
	}
	if menuBoxArgSet(args, 3) {
		params["mask"] = argumentValue(args, 3)
	}
	menuBoxApplyDataFields(params, args, 4, []string{"field1", "field2", "field3", "field4"})

	return params
}

func init() {
	enumerations["menuBoxSymbolWeight"] = []string{"ultralight", "thin", "light", "regular", "medium", "semibold", "bold", "heavy", "black"}
	enumerations["menuBoxIconScale"] = []string{"small", "medium", "large"}
	// "none" is the default mask and is only written when set explicitly, which is how its raw
	// value was seen in a device export (Create Menu Item From Emoji).
	enumerations["menuBoxIconMask"] = []string{"none", "circle", "roundRect"}
	// Every Sort By option of Load Menu Set, each read from its own device export.
	enumerations["menuBoxLoadSortBy"] = []string{"title", "subtitle", "field1", "field2", "field3", "field4", "dateCreated", "dateModified"}
	enumerations["menuBoxLoadOrder"] = []string{"ascending", "descending"}
	enumerations["menuBoxNameOperator"] = []string{"is", "isNot", "contains", "beginsWith", "endsWith"}
	enumerations["menuBoxDateOperator"] = []string{"exactly", "notExactly", "after", "onOrAfter", "before", "onOrBefore"}
	enumerations["menuBoxCountOperator"] = []string{"is", "isNot", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"}
	// Random is capitalised because that is how Shortcuts writes it, unlike the name and date properties.
	enumerations["menuBoxFindSortBy"] = []string{"name", "dateCreated", "dateModified", "Random"}
	enumerations["menuBoxFindOrder"] = []string{"latestFirst", "oldestFirst", "aToZ", "zToA"}

	actions["loadMenuSet"] = &actionDefinition{
		doc: selfDoc{
			title: "Load Menu Set",
			description: "Loads the menu items of a saved menu set, for use with Choose From List, optionally sorted by " +
				"`sortBy` (title, subtitle, field1 to field4, dateCreated or dateModified) in `order` (ascending, the default, or " +
				"descending). The set is " +
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
			{
				name:      "alternateSetName",
				validType: String,
				literal:   true,
				optional:  true,
			},
			{
				name:      "alternateSetID",
				validType: String,
				literal:   true,
				optional:  true,
			},
			{
				name:      "sortBy",
				validType: String,
				enum:      "menuBoxLoadSortBy",
				literal:   true,
				optional:  true,
			},
			{
				name:      "order",
				validType: String,
				enum:      "menuBoxLoadOrder",
				literal:   true,
				optional:  true,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			menuBoxCheckPair(args, 3, "alternate set")
		},
		makeParams: func(args []actionArgument) map[string]any {
			var params = map[string]any{
				"menuSet": menuBoxSetEntity(argumentValue(args, 0).(string), argumentValue(args, 1).(string)),
			}
			// Shortcuts only writes loadAlt once the toggle has been touched, so an omitted argument
			// leaves the key out entirely.
			if menuBoxArgSet(args, 2) {
				params["loadAlt"] = argumentValue(args, 2)
			}
			if menuBoxArgSet(args, 3) {
				params["alternateSet"] = menuBoxSetEntity(argumentValue(args, 3).(string), argumentValue(args, 4).(string))
			}
			if menuBoxArgSet(args, 5) {
				params["sortBy"] = argumentValue(args, 5)
			}
			// Ascending is the default and an untouched Order leaves the key out of an export, so only
			// descending is written. That keeps the raw spelling of ascending out of the plist, since it
			// has never been seen.
			if menuBoxArgSet(args, 6) && argumentValue(args, 6).(string) == "descending" {
				params["order"] = "descending"
			}

			return params
		},
		outputType: Variable,
		outputName: "Menu Items",
	}

	actions["exportMenuSet"] = &actionDefinition{
		doc: selfDoc{
			title: "Export Menu Set",
			description: "Exports a saved menu set as Menu Set Data, which Import Menu Sets reads back in. The set is " +
				"identified by its name and its set ID (long-press the set in Menu Box to copy the ID). " +
				"Set IDs belong to one Menu Box library, so a literal ID will not resolve on another person's device. " +
				"From the Menu Box app.",
			category:    "thirdparty",
			subcategory: "Menu Box",
		},
		appIdentifier: menuBoxBundleIdentifier,
		// Unlike the other intents this one has no "Intent" suffix, exactly as Shortcuts writes it.
		identifier: "ExportMenuSets",
		appIntent:  menuBoxIntent("ExportMenuSets"),
		nonMacOnly: true,
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
		},
		makeParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"menuSets": menuBoxSetEntity(argumentValue(args, 0).(string), argumentValue(args, 1).(string)),
			}
		},
		outputType: Variable,
		outputName: "Menu Set Data",
	}

	actions["importMenuSets"] = &actionDefinition{
		doc: selfDoc{
			title: "Import Menu Sets",
			description: "Imports menu sets from Menu Set Data, such as the output of Export Menu Set. " +
				"The data is passed as a variable or constant, not as a text literal. " +
				"`overwrite` replaces sets that already exist in your library, so use it with care. " +
				"From the Menu Box app.",
			category:    "thirdparty",
			subcategory: "Menu Box",
		},
		appIdentifier: menuBoxBundleIdentifier,
		identifier:    "ImportMenuSetsIntent",
		appIntent:     menuBoxIntent("ImportMenuSetsIntent"),
		nonMacOnly:    true,
		parameters: []parameterDefinition{
			{
				// A variable, because the Export output is a generic variable and Shortcuts stores it
				// as a bare attachment, as seen in a device export.
				name:      "menuSetData",
				validType: Variable,
				key:       "exports",
			},
			{
				name:      "overwrite",
				validType: Bool,
				key:       "shouldOverwrite",
				optional:  true,
			},
		},
		makeParams: func(args []actionArgument) map[string]any {
			var params = map[string]any{
				"exports": argumentValue(args, 0),
			}
			// Written only once the toggle has been touched, like every other Menu Box switch.
			if menuBoxArgSet(args, 1) {
				params["shouldOverwrite"] = argumentValue(args, 1)
			}

			return params
		},
	}

	// Menu Box's Add To Menu Set is one intent with a Set Type switch, "new" or "existing", that reveals
	// different parameters. It is split in two here, like the two image actions, so each one only
	// offers the parameters that have been seen with its mode.
	actions["createMenuSet"] = &actionDefinition{
		doc: selfDoc{
			title: "Create Menu Set",
			description: "Creates a new menu set from menu items, such as the output of Quick Menu or Load Menu Set. " +
				"This is Menu Box's Add To Menu Set action with Set Type set to new. " +
				"`overwrite` replaces a set that already has the same name, so use it with care. " +
				"An optional tag is identified by its name and its tag ID, which belongs to one Menu Box library, " +
				"so a literal ID will not resolve on another person's device. From the Menu Box app.",
			category:    "thirdparty",
			subcategory: "Menu Box",
		},
		appIdentifier: menuBoxBundleIdentifier,
		identifier:    "CreateMenuSetIntent",
		appIntent:     menuBoxIntent("CreateMenuSetIntent"),
		nonMacOnly:    true,
		parameters: []parameterDefinition{
			{
				name:      "menuItems",
				validType: Variable,
			},
			{
				name:      "setName",
				validType: String,
			},
			{
				name:      "overwrite",
				validType: Bool,
				key:       "shouldOverwrite",
				optional:  true,
			},
			{
				name:      "tagName",
				validType: String,
				literal:   true,
				optional:  true,
			},
			{
				name:      "tagID",
				validType: String,
				literal:   true,
				optional:  true,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			menuBoxCheckPair(args, 3, "tag")
		},
		makeParams: func(args []actionArgument) map[string]any {
			var params = map[string]any{
				"setType":   "new",
				"menuItems": argumentValue(args, 0),
				"setName":   argumentValue(args, 1),
			}
			// Written only once the toggle has been touched, like every other Menu Box switch.
			if menuBoxArgSet(args, 2) {
				params["shouldOverwrite"] = argumentValue(args, 2)
			}
			if menuBoxArgSet(args, 3) {
				params["tags"] = menuBoxTagEntity(argumentValue(args, 3).(string), argumentValue(args, 4).(string))
			}

			return params
		},
	}

	actions["addToMenuSet"] = &actionDefinition{
		doc: selfDoc{
			title: "Add To Menu Set",
			description: "Adds menu items to a menu set that already exists, such as the output of Quick Menu or " +
				"Load Menu Set. This is Menu Box's Add To Menu Set action with Set Type set to existing. " +
				"The set is found by its name, so unlike Load Menu Set no set ID is needed. From the Menu Box app.",
			category:    "thirdparty",
			subcategory: "Menu Box",
		},
		appIdentifier: menuBoxBundleIdentifier,
		identifier:    "CreateMenuSetIntent",
		appIntent:     menuBoxIntent("CreateMenuSetIntent"),
		nonMacOnly:    true,
		parameters: []parameterDefinition{
			{
				name:      "menuItems",
				validType: Variable,
			},
			{
				name:      "setName",
				validType: String,
				literal:   true,
			},
		},
		makeParams: func(args []actionArgument) map[string]any {
			return map[string]any{
				"setType":         "existing",
				"menuItems":       argumentValue(args, 0),
				"existingSetName": menuBoxExistingSetEntity(argumentValue(args, 1).(string)),
			}
		},
	}

	actions["findMenuSets"] = &actionDefinition{
		doc: selfDoc{
			title: "Find Menu Sets",
			description: "Finds menu sets, optionally filtered by name, by the titles of the items inside them, by date " +
				"created, by date modified and by how many items they hold. Each filter takes an operator beside it: " +
				"`nameOperator` (is, isNot, contains, the default, beginsWith or endsWith), `createdOperator` and " +
				"`modifiedOperator` (exactly, the default, notExactly, after, onOrAfter, before or onOrBefore) and " +
				"`itemCountOperator` (is, the default, isNot, greaterThan, greaterThanOrEqual, lessThan or lessThanOrEqual). " +
				"`itemTitle` matches sets with an item whose title contains it. `createdDate` and `modifiedDate` take a " +
				"date from a Date action. Filters are combined with all by default, or any with `matchAny`. Results are " +
				"sorted by `sortBy` (name, dateCreated, dateModified or Random) in `order` (aToZ or zToA for name, latestFirst " +
				"or oldestFirst for the dates, which Random does not use), and `limit` caps how many are returned. Pass a " +
				"previous Find's output as `input` to filter that list instead of every set. Only one filter per property " +
				"is possible here; use `rawAction()` for more. From the Menu Box app.",
			category:    "thirdparty",
			subcategory: "Menu Box",
		},
		appIdentifier: menuBoxBundleIdentifier,
		identifier:    "MenuSetAppEntity",
		appIntent:     menuBoxQueryIntent("MenuSetAppEntity"),
		nonMacOnly:    true,
		parameters: []parameterDefinition{
			{name: "name", validType: String, literal: true, optional: true},
			{name: "nameOperator", validType: String, enum: "menuBoxNameOperator", literal: true, optional: true},
			{name: "itemTitle", validType: String, literal: true, optional: true},
			{name: "createdDate", validType: Variable, optional: true},
			{name: "createdOperator", validType: String, enum: "menuBoxDateOperator", literal: true, optional: true},
			{name: "modifiedDate", validType: Variable, optional: true},
			{name: "modifiedOperator", validType: String, enum: "menuBoxDateOperator", literal: true, optional: true},
			{name: "itemCount", validType: Integer, literal: true, optional: true},
			{name: "itemCountOperator", validType: String, enum: "menuBoxCountOperator", literal: true, optional: true},
			{name: "matchAny", validType: Bool, optional: true},
			{name: "sortBy", validType: String, enum: "menuBoxFindSortBy", literal: true, optional: true},
			{name: "order", validType: String, enum: "menuBoxFindOrder", literal: true, optional: true},
			{name: "limit", validType: Integer, optional: true},
			{name: "input", validType: Variable, optional: true},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			// An operator without a value would be a row Shortcuts shows empty, which nobody means to write.
			for _, filter := range menuBoxFindFilters {
				if filter.operatorArg >= 0 && menuBoxArgSet(args, filter.operatorArg) && !menuBoxArgSet(args, filter.valueArg) {
					parserError(fmt.Sprintf("%s needs %s to compare against.", filter.operatorName, filter.valueName))
				}
			}
			if !menuBoxArgSet(args, menuBoxFindArgOrder) {
				return
			}
			// The Order row only appears once a sort is chosen.
			if !menuBoxArgSet(args, menuBoxFindArgSortBy) {
				parserError("order needs a sortBy to apply to.")
			}
			var sortBy = argumentValue(args, menuBoxFindArgSortBy).(string)
			// Random has no direction; the Order row disappears for it in Shortcuts.
			if sortBy == "Random" {
				parserError("order does not apply when sorting by Random.")
			}
			var order = argumentValue(args, menuBoxFindArgOrder).(string)
			if !slices.Contains(menuBoxFindOrderSort[order], sortBy) {
				parserError(fmt.Sprintf("order '%s' does not apply when sorting by %s, use %s.", order, sortBy, strings.Join(menuBoxFindOrdersFor(sortBy), " or ")))
			}
		},
		makeParams: func(args []actionArgument) map[string]any {
			var params = map[string]any{}

			var templates []map[string]any
			for _, filter := range menuBoxFindFilters {
				if !menuBoxArgSet(args, filter.valueArg) {
					continue
				}
				var operator = filter.defaultOperator
				if filter.operatorArg >= 0 && menuBoxArgSet(args, filter.operatorArg) {
					operator = argumentValue(args, filter.operatorArg).(string)
				}
				var value = argumentValue(args, filter.valueArg)
				if filter.valueKey == "Number" {
					// A count is stored as a string, as the device export showed.
					value = fmt.Sprint(value)
				}
				templates = append(templates, menuBoxFilterTemplate(filter.property, filter.operators[operator], map[string]any{"Unit": 4, filter.valueKey: value}))
			}
			if len(templates) > 0 {
				// All is 1 and Any is 0, as seen in device exports.
				var prefix = 1
				if menuBoxArgSet(args, menuBoxFindArgMatchAny) && argumentValue(args, menuBoxFindArgMatchAny) == true {
					prefix = 0
				}
				params["WFContentItemFilter"] = map[string]any{
					"Value": map[string]any{
						"WFActionParameterFilterPrefix":    prefix,
						"WFContentPredicateBoundedDate":    false,
						"WFActionParameterFilterTemplates": templates,
					},
					"WFSerializationType": "WFContentPredicateTableTemplate",
				}
			}

			if menuBoxArgSet(args, menuBoxFindArgSortBy) {
				params["WFContentItemSortProperty"] = argumentValue(args, menuBoxFindArgSortBy)
			}
			if menuBoxArgSet(args, menuBoxFindArgOrder) {
				params["WFContentItemSortOrder"] = menuBoxFindOrders[argumentValue(args, menuBoxFindArgOrder).(string)]
			}
			if menuBoxArgSet(args, menuBoxFindArgLimit) {
				params["WFContentItemLimitEnabled"] = true
				params["WFContentItemLimitNumber"] = argumentValue(args, menuBoxFindArgLimit)
			}
			if menuBoxArgSet(args, menuBoxFindArgInput) {
				params["WFContentItemInputParameter"] = argumentValue(args, menuBoxFindArgInput)
			}

			return params
		},
		outputType: Variable,
		outputName: "Menu Sets",
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
			menuBoxApplyDataFields(params, args, 10, []string{"field1", "field2", "field3", "field4"})

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
			menuBoxApplyDataFields(params, args, 8, []string{"f1", "f2", "f3", "f4"})

			return params
		},
		outputType: Variable,
		outputName: "Menu Item",
	}

	actions["createMenuItemFromImage"] = &actionDefinition{
		doc: selfDoc{
			title: "Create Menu Item From Image",
			description: "Creates a menu item whose icon is an image file, for use with Choose From List or Add To Menu Set. " +
				"The image must be a variable holding an image or file, such as the output of Get File, Select Photos or " +
				"Base64 Decode. A file picked straight in the Shortcuts editor is saved as a reference tied to one " +
				"person's iCloud account, so it can't be written from code. The image is used as it is: the only thing " +
				"you can change is the `mask`, which defaults to none. Fields 1 to 4 hold data that isn't shown in " +
				"the menu but can be read from the chosen item. To use a base64 string instead, see " +
				"`createMenuItemFromBase64`. From the Menu Box app.",
			category:    "thirdparty",
			subcategory: "Menu Box",
		},
		appIdentifier: menuBoxBundleIdentifier,
		identifier:    "CreateImageItemIntent",
		appIntent:     menuBoxIntent("CreateImageItemIntent"),
		nonMacOnly:    true,
		parameters: menuBoxImageParameters(parameterDefinition{
			name:      "image",
			validType: Variable,
			key:       "imageFile",
		}),
		makeParams: func(args []actionArgument) map[string]any {
			return menuBoxMakeImageParams(args, "imageFile", "file")
		},
		outputType: Variable,
		outputName: "Menu Item",
	}

	actions["createMenuItemFromBase64"] = &actionDefinition{
		doc: selfDoc{
			title: "Create Menu Item From Base64",
			description: "Creates a menu item whose icon is a base64 encoded image, for use with Choose From List or " +
				"Add To Menu Set. The string is plain base64 with no `data:` prefix, and can be text or a variable such " +
				"as the output of a Text or Base64 Encode action. Menu Box stores icons this way itself, so a Menu Set " +
				"export holds the same string. The image is used as it is: the only thing you can change is the `mask`, " +
				"which defaults to none. Fields 1 to 4 hold data that isn't shown in the menu but can be read from the " +
				"chosen item. To use an image file instead, see `createMenuItemFromImage`. From the Menu Box app.",
			category:    "thirdparty",
			subcategory: "Menu Box",
		},
		appIdentifier: menuBoxBundleIdentifier,
		identifier:    "CreateImageItemIntent",
		appIntent:     menuBoxIntent("CreateImageItemIntent"),
		nonMacOnly:    true,
		parameters: menuBoxImageParameters(parameterDefinition{
			name:      "base64",
			validType: String,
			key:       "base64",
		}),
		makeParams: func(args []actionArgument) map[string]any {
			return menuBoxMakeImageParams(args, "base64", "base64")
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
