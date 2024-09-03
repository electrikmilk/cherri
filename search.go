/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"github.com/electrikmilk/args-parser"
	"log"
	"os"
	"strings"
)

// Initialize the logger
var logger = log.New(os.Stdout, "DEBUG",log.Ldate|log.Ltime|log.Lshortfile)

// actionSearch searches for action based on a provided identifier
func actionsSearch() {
	// Get the value of the "action" argument
	var identifier = args.Value("action")
    logger.Printf("Searching for action: %s", identifier)

	// check if the action is exists in the 'action' map
	if _, found := actions[identifier]; !found {
		fmt.Println(ansi(fmt.Sprintf("\nAction '%s(...)' does not exist or has not yet been defined.", identifier), red))

		// special case handling for "text" and "dictionary" actions
		switch identifier {
		case "text":
			// Log the specific help message for text actions
			logger.Println("Providing help for 'text' action")
			fmt.Print("\nText actions are abstracted into string statements. For example:\n\n@variable = \"Hello, Cherri!\"\n\n")
			os.Exit(1)
		case "dictionary":
			// Log the specific help message for dictionary actions
			logger.Println("Providing help for 'dictionary' action")
			fmt.Print("\nDictionary actions are abstracted into JSON object statements. For example:\n\n@variable = {\"test\":5\", \"key\":\"value\"}\n\n")
			os.Exit(1)
		}

		// Search for similar actions if the exact action is not found
		var actionSearchResults strings.Builder
		for actionIdentifier, definition := range actions {
            logger.Printf("Checking if '%s' matches '%s'",actionIdentifier,identifier)
			var matched, result = matchString(&actionIdentifier, &identifier)
			if matched {
				logger.Printf("Matched action: %s", actionIdentifier)
				setCurrentAction(actionIdentifier, definition)
				var definition = generateActionDefinition(parameterDefinition{}, false, false)
				// Remove the actionIdentifier prefix from the defination
				definition, _ = strings.CutPrefix(definition, actionIdentifier)

				// Add the matched result to the search results
				actionSearchResults.WriteString(fmt.Sprintf("- %s%s\n", result, definition))
			}
		}
		// Print the closest actions if any matches are found
		if actionSearchResults.Len() > 0 {
			logger.Println("Closest actions found")
			fmt.Println(ansi("\nThe closest actions are:", yellow, italic, bold))
			fmt.Println(actionSearchResults.String())
		}else{
			logger.Println("No close actions found")
		}

		os.Exit(1)
	}
	// If the action is found, set it as the current action
	logger.Printf("Action '%s' found, setting as current action", identifier)
	setCurrentAction(identifier, actions[identifier])

	// Generate and print the action defination
	fmt.Println(generateActionDefinition(parameterDefinition{}, true, true))
}

// glyphsSearch for a glyph based on a provided identifier
func glyphsSearch() {
	// Get the value of the "glyph" argument
	var identifier = args.Value("glyph")
    logger.Printf("Searching for glyph: %s", identifier)

	var searchResults strings.Builder
	// Iterate through available glyph to find matches
	for glyphIdentifier := range glyphs {
		logger.Printf("Checking if glyph '%s' matches '%s'",glyphIdentifier,identifier)
		var matched, result = matchString(&glyphIdentifier, &identifier)
		if matched {
			logger.Printf("Matched glyph: %s", glyphIdentifier)
			// Add matched glyphs to the search results
			searchResults.WriteString(fmt.Sprintf("- %s\n", result))
		}
	}
	// Print the closest glyphs if any matches are found
	if searchResults.Len() > 0 {
		logger.Println("Closest glyphs found")
		fmt.Println(ansi("\nThe closest glyphs are:", yellow, italic, bold))
		fmt.Println(searchResults.String())
	}else{
		logger.Println("No close glyphs found")
	}
}

// matchString checks if the search string matches or is contained in the subject string

func matchString(subject *string, search *string) (matched bool, result string) {
	// check for exact match
	if *subject == *search {
		result = ansi(*subject, red)
		matched = true
		logger.Printf("Exact match found: %s", *subject)
		return
	}

	// Check for partial match (case-insensitive)
	if strings.Contains(strings.ToLower(*subject), strings.ToLower(*search)) {
		matched = true
		logger.Printf("Partial match found: %s contains %s", *subject, *search)
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
	}else{
		logger.Printf("No match found for: %s", *search)
	}

	return
}
