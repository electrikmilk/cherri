/*
 * Copyright (c) Cherri
 */

package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/electrikmilk/args-parser"
)

func handleActionSearch() {
	if args.Value("action") == "" {
		for identifier, definition := range actions {
			setCurrentAction(identifier, definition)
			printActionDebugInfo()
			fmt.Println(generateActionDefinition(parameterDefinition{}, true))
			fmt.Print("\n---\n")
		}
		return
	}

	actionsSearch()

	os.Exit(1)
}

func actionsSearch() {
	var identifier = args.Value("action")
	if _, found := actions[identifier]; found {
		setCurrentAction(identifier, actions[identifier])
		printActionDebugInfo()
		fmt.Println(generateActionDefinition(parameterDefinition{}, true))
		return
	}

	fmt.Println(ansi(fmt.Sprintf("\nAction '%s(...)' does not exist or has not yet been defined.", identifier), red))

	switch identifier {
	case "text":
		fmt.Print("\nText actions are abstracted into string statements. For example:\n\n@variable = \"Hello, Cherri!\"\n\n")
		os.Exit(1)
	case "dictionary":
		fmt.Print("\nDictionary actions are abstracted into JSON object statements. For example:\n\n@variable = {\"test\":5\", \"key\":\"value\"}\n\n")
	}

	var actionSearchResults strings.Builder
	for actionIdentifier, definition := range actions {
		var matched, result = matchString(&actionIdentifier, &identifier)
		var matchedDoc bool
		var docResult string
		if definition.doc.title != "" {
			matchedDoc, docResult = matchString(&definition.doc.title, &identifier)
		}
		if matched || matchedDoc {
			setCurrentAction(actionIdentifier, definition)
			var definition = generateActionDefinition(parameterDefinition{}, false)
			definition, _ = strings.CutPrefix(definition, actionIdentifier)
			if !matched {
				result = docResult
			}

			actionSearchResults.WriteString(fmt.Sprintf("\n- %s\n%s\n\n", result, definition))
		}
	}
	if actionSearchResults.Len() > 0 {
		fmt.Println(ansi("\nThe closest actions are:", yellow, italic, bold))
		fmt.Println(actionSearchResults.String())
	}
}

type debugParameter struct {
	Key  string   `json:"key"`
	Enum []string `json:"enum,omitempty"`
}

type debugParameters struct {
	Parameters []debugParameter `json:"parameters"`
	Fixed      []string         `json:"fixed,omitempty"`
	Dynamic    bool             `json:"dynamic"`
}

// printActionDebugInfo prints what is needed to match the current action against the Shortcuts ToolKit DB.
// The identifier is part of the debug action definition printed after this.
func printActionDebugInfo() {
	if !args.Using("debug") {
		return
	}

	printActionParameters()
}

// printActionParameters prints the plist keys the current action maps from its arguments as JSON.
// Dynamic actions build some of their keys in Go from the arguments, so their keys cannot be listed completely.
func printActionParameters() {
	var definition = currentAction.definition
	var debugInfo = debugParameters{
		Parameters: []debugParameter{},
		Fixed:      slices.Sorted(maps.Keys(definition.appendParams)),
		Dynamic:    definition.makeParams != nil || definition.appendParamsFunc != nil,
	}

	for _, param := range definition.parameters {
		if param.key == "" {
			continue
		}
		debugInfo.Parameters = append(debugInfo.Parameters, debugParameter{Key: param.key, Enum: enumerations[param.enum]})
	}

	var encoded, encodeErr = json.Marshal(debugInfo)
	handle(encodeErr)
	fmt.Println("parameters:", string(encoded))
}

func handleGlyphSearch() {
	if args.Value("glyph") == "" {
		fmt.Println(ansi("You can generate Shortcut icon code using: https://glyphs.cherrilang.org.\n", cyan))
		if args.Value("action") == "" {
			for identifier := range glyphs {
				fmt.Println(identifier)
			}
			return
		}
	}

	glyphsSearch()
}

func glyphsSearch() {
	var identifier = args.Value("glyph")
	var searchResults strings.Builder
	for glyphIdentifier := range glyphs {
		var matched, result = matchString(&glyphIdentifier, &identifier)
		if matched {
			searchResults.WriteString(fmt.Sprintf("- %s\n", result))
		}
	}
	if searchResults.Len() > 0 {
		fmt.Println(ansi("\nThe closest glyphs are:", yellow, italic, bold))
		fmt.Println(searchResults.String())
	}
}

func matchString(subject *string, search *string) (matched bool, result string) {
	if *subject == *search {
		result = ansi(*subject, red)
		matched = true
		return
	}

	if strings.Contains(strings.ToLower(*subject), strings.ToLower(*search)) {
		matched = true
		if args.Using("no-ansi") {
			result = *subject
			return
		}
		var capitalized = capitalize(*search)
		var lowercase = strings.ToLower(*search)
		switch {
		case strings.Contains(*subject, *search):
			result = strings.ReplaceAll(*subject, *search, ansi(*search, red))
		case strings.Contains(*subject, capitalized):
			result = strings.ReplaceAll(*subject, capitalized, ansi(capitalized, red))
		case strings.Contains(*subject, lowercase):
			result = strings.ReplaceAll(*subject, lowercase, ansi(lowercase, red))
		}
	}

	return
}
