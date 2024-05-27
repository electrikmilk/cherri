/*
 * Copyright (c) Brandon Jordan
 */

package main

import (
	"fmt"
	"github.com/electrikmilk/args-parser"
	"os"
	"strings"
)

func actionsSearch() {
	var identifier = args.Value("action")
	if _, found := actions[identifier]; !found {
		fmt.Println(ansi(fmt.Sprintf("\nAction '%s(...)' does not exist or has not yet been defined.", identifier), red))

		switch identifier {
		case "text":
			fmt.Print("\nText actions are abstracted into string statements. For example:\n\n@variable = \"Hello, Cherri!\"\n\n")
			os.Exit(1)
		case "dictionary":
			fmt.Print("\nDictionary actions are abstracted into JSON object statements. For example:\n\n@variable = {\"test\":5\", \"key\":\"value\"}\n\n")
			os.Exit(1)
		}

		var actionSearchResults strings.Builder
		for actionIdentifier, definition := range actions {
			if strings.Contains(strings.ToLower(actionIdentifier), identifier) {
				setCurrentAction(actionIdentifier, definition)
				var definition = generateActionDefinition(parameterDefinition{}, false, false)
				definition, _ = strings.CutPrefix(definition, actionIdentifier)

				var capitalized = capitalize(identifier)
				var lowercase = strings.ToLower(identifier)
				switch {
				case strings.Contains(actionIdentifier, identifier):
					identifier = strings.ReplaceAll(actionIdentifier, identifier, ansi(identifier, red))
				case strings.Contains(actionIdentifier, capitalized):
					identifier = strings.ReplaceAll(actionIdentifier, capitalized, ansi(capitalized, red))
				case strings.Contains(actionIdentifier, lowercase):
					identifier = strings.ReplaceAll(actionIdentifier, lowercase, ansi(lowercase, red))
				}
				actionSearchResults.WriteString(fmt.Sprintf("- %s%s\n", identifier, definition))
			}
		}
		if actionSearchResults.Len() > 0 {
			fmt.Println(ansi("\nThe closest actions are:", yellow, italic, bold))
			fmt.Println(actionSearchResults.String())
		}

		os.Exit(1)
	}
	setCurrentAction(identifier, actions[identifier])
	fmt.Println(generateActionDefinition(parameterDefinition{}, true, true))
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
