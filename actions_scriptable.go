/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"strings"
)

// Actions from Scriptable (dk.simonbs.Scriptable), a third-party JavaScript runner. Like Pythonista and
// unlike the App Intents providers in this package, its actions carry no AppIntentDescriptor and use plain
// camelCase parameter keys, so only the workflow identifier is needed.
//
// Scriptable only exists on iOS and iPadOS, so every definition here is nonMacOnly. Everything below was read
// from device exports; Shortcuts omits any parameter that was never touched, so an untouched action is just
// its UUID and the Cherri definitions only write what the caller passed.
const scriptableBundleIdentifier = "dk.simonbs.Scriptable"

// Argument positions shared by the two Run actions.
const (
	scriptableRunArgScript = iota
	scriptableRunArgParameter
	scriptableRunArgTexts
	scriptableRunArgURLs
	scriptableRunArgImages
	scriptableRunArgFiles
)

// scriptableRunDoc builds the doc shared by the two Run actions, which only differ in where the code comes from.
func scriptableRunDoc(title string, description string) selfDoc {
	return selfDoc{
		title: title,
		description: description + " `parameter` is passed to the script as `args.shortcutParameter`, and `texts`, `urls`, " +
			"`images` and `files` as its `plainTexts`, `urls`, `images` and `fileURLs` lists. `runInApp` opens Scriptable to " +
			"run it. From the Scriptable app.",
		category:    "thirdparty",
		subcategory: "Scriptable",
	}
}

// scriptableRunParameters lists the parameters shared by both Run actions after the code source. Shortcuts only
// accepts a variable for the Parameter, Images and Files inputs; Texts and URLs are lists of typed items that are
// written by scriptableRunParams because Shortcuts writes a lone item as a bare value and two or more as a list.
func scriptableRunParameters(source parameterDefinition) []parameterDefinition {
	return []parameterDefinition{
		source,
		{
			name:      "parameter",
			validType: Variable,
			key:       "parameter",
			optional:  true,
		},
		{
			name:      "texts",
			validType: Arr,
			literal:   true,
			optional:  true,
		},
		{
			name:      "urls",
			validType: Arr,
			literal:   true,
			optional:  true,
		},
		{
			name:      "images",
			validType: Variable,
			key:       "images",
			optional:  true,
		},
		{
			name:      "files",
			validType: Variable,
			key:       "files",
			optional:  true,
		},
		{
			name:         "runInApp",
			validType:    Bool,
			key:          "runInApp",
			defaultValue: false,
			optional:     true,
		},
		{
			name:         "showWhenRun",
			validType:    Bool,
			key:          "ShowWhenRun",
			defaultValue: true,
			optional:     true,
		},
	}
}

// scriptableItemList builds one of the Texts or URLs lists. Each item is either plain text or a text token.
func scriptableItemList(arg actionArgument) any {
	if arg.valueType != Arr {
		return nil
	}

	var items = arg.value.([]interface{})
	var values = make([]any, 0, len(items))
	for _, item := range items {
		values = append(values, attachmentValues(fmt.Sprint(item)))
	}

	switch len(values) {
	case 0:
		return nil
	case 1:
		return values[0]
	}

	return values
}

// scriptableRunParams builds the parameters of the two Run actions that automatic handling cannot.
func scriptableRunParams(args []actionArgument, build func(arg actionArgument) (key string, value any)) map[string]any {
	var params = make(map[string]any)
	if len(args) > scriptableRunArgScript {
		var key, value = build(args[scriptableRunArgScript])
		params[key] = value
	}
	if len(args) > scriptableRunArgTexts {
		if texts := scriptableItemList(args[scriptableRunArgTexts]); texts != nil {
			params["texts"] = texts
		}
	}
	if len(args) > scriptableRunArgURLs {
		if urls := scriptableItemList(args[scriptableRunArgURLs]); urls != nil {
			params["urls"] = urls
		}
	}

	return params
}

func init() {
	actions["runScriptableScript"] = &actionDefinition{
		doc: scriptableRunDoc(
			"Run Scriptable Script",
			"Runs a script from the Scriptable library by `script` name, without its .js extension.",
		),
		appIdentifier: scriptableBundleIdentifier,
		identifier:    "ParameterizedRunScriptIntent",
		nonMacOnly:    true,
		defaultAction: true,
		parameters: scriptableRunParameters(parameterDefinition{
			name:      "script",
			validType: String,
		}),
		check: func(args []actionArgument, _ *actionDefinition) {
			if len(args) > scriptableRunArgScript && args[scriptableRunArgScript].valueType == String &&
				strings.ContainsAny(args[scriptableRunArgScript].value.(string), "{}") {
				parserError("Scriptable's script picker only accepts a script name or a whole variable, not text with inline variables.")
			}
		},
		appendParamsFunc: func(args []actionArgument) map[string]any {
			return scriptableRunParams(args, func(arg actionArgument) (string, any) {
				if arg.valueType != Variable {
					return "fileName", arg.value
				}

				// Scriptable's own Ask Each Time carries no VariableName, unlike the one Cherri writes for
				// first-party actions, so drop it to match what the Shortcuts app produces.
				var attachment = variableValue(arg.value.(varValue))
				if token, ok := attachment.(WFTextTokenAttachment); ok && token.Value.Type == Ask {
					token.Value.VariableName = ""
					return "fileName", token
				}

				return "fileName", attachment
			})
		},
		outputType: Variable,
		outputName: "Output",
	}

	actions["runScriptableInlineScript"] = &actionDefinition{
		doc: scriptableRunDoc(
			"Run Scriptable Inline Script",
			"Runs JavaScript written into the action instead of a library script. Use a raw string ('...') for `source` "+
				"so braces in the code are not read as variables.",
		),
		appIdentifier: scriptableBundleIdentifier,
		identifier:    "RunScriptInlineIntent",
		nonMacOnly:    true,
		parameters: scriptableRunParameters(parameterDefinition{
			name:      "source",
			validType: String,
		}),
		appendParamsFunc: func(args []actionArgument) map[string]any {
			return scriptableRunParams(args, func(arg actionArgument) (string, any) {
				return "script", paramValue(arg, String)
			})
		},
		outputType: Variable,
		outputName: "Output",
	}

	actions["refreshAllScriptableWidgets"] = &actionDefinition{
		doc: selfDoc{
			title:       "Refresh All Scriptable Widgets",
			description: "Refreshes every Scriptable widget. From the Scriptable app.",
			category:    "thirdparty",
			subcategory: "Scriptable",
		},
		appIdentifier: scriptableBundleIdentifier,
		identifier:    "RefreshAllWidgetsIntent",
		nonMacOnly:    true,
	}

	actions["createScriptableFileBookmark"] = &actionDefinition{
		doc: selfDoc{
			title: "Create Scriptable File Bookmark",
			description: "Creates a Scriptable file bookmark named `name` for `file`, which scripts can then read " +
				"with `FileManager.bookmarkedPath`. From the Scriptable app.",
			category:    "thirdparty",
			subcategory: "Scriptable",
		},
		appIdentifier: scriptableBundleIdentifier,
		identifier:    "CreateFileBookmarkIntent",
		nonMacOnly:    true,
		parameters: []parameterDefinition{
			{
				name:      "name",
				validType: String,
				key:       "name",
			},
			{
				name:      "file",
				validType: Variable,
				key:       "location",
			},
		},
	}
}
