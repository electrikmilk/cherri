/*
 * Copyright (c) Cherri
 */

package main

import (
	_ "embed"
	"fmt"
	"strings"
)

// The ToolKit database only describes a world clock city as the result of a query, so the list of cities is taken from
// the Localizable_Cities table in AppSupport.framework, whose keys are the identifiers that Shortcuts writes for a city.
//
//go:embed world_clock_cities.tsv
var worldClockCitiesData string

const worldClockCityEnum = "worldClockCity"

type worldClockCity struct {
	identifier string
	name       string
}

// worldClockCities maps a value of the worldClockCity enum to its city.
var worldClockCities = map[string]worldClockCity{}

// buildWorldClockCities registers the worldClockCity enum, which must exist before the actions that use it are defined.
func buildWorldClockCities() {
	if len(worldClockCities) > 0 {
		return
	}

	var values []string
	for _, line := range strings.Split(strings.TrimSpace(worldClockCitiesData), "\n") {
		var identifier, rest, found = strings.Cut(line, "\t")
		if !found {
			panic(fmt.Sprintf("Invalid world clock city line '%s'", line))
		}
		var name, region, _ = strings.Cut(rest, "\t")

		var value = name
		if region != "" {
			value = fmt.Sprintf("%s, %s", name, region)
		}
		if _, duplicate := worldClockCities[value]; duplicate {
			panic(fmt.Sprintf("Duplicate world clock city '%s' (%s)", value, identifier))
		}
		worldClockCities[value] = worldClockCity{identifier: identifier, name: name}
		values = append(values, value)
	}
	enumerations[worldClockCityEnum] = values
}

// checkWorldClockCities validates the names in an array of cities, since an array has no enum to check its values against.
func checkWorldClockCities(args []actionArgument, _ *actionDefinition) {
	if len(args) == 0 || args[0].valueType == Variable {
		return
	}

	var names, isArray = getArgValue(args[0]).([]interface{})
	if !isArray {
		return
	}
	for _, name := range names {
		var value, isString = name.(string)
		if !isString {
			parserError("World clock cities must be text")
		}
		if _, found := worldClockCities[value]; !found {
			parserError(fmt.Sprintf("Invalid world clock city '%s'", value))
		}
	}
}

// worldClockCityValue is the plist value of a city argument, which is a variable or a value of the worldClockCity enum.
func worldClockCityValue(argument actionArgument) any {
	if argument.valueType == Variable {
		return paramValue(argument, Variable)
	}

	var value, isString = getArgValue(argument).(string)
	if !isString {
		parserError("A world clock city must be text or a variable")
	}
	var city, found = worldClockCities[value]
	if !found {
		parserError(fmt.Sprintf("Invalid world clock city '%s'", value))
	}

	return map[string]any{
		"identifier": city.identifier,
		"title": map[string]string{
			"key": city.name,
		},
		"subtitle": map[string]string{
			"key": city.name,
		},
	}
}

// worldClockCitiesValue is the plist value of an array of cities, which Shortcuts writes as a single city when there is only one.
func worldClockCitiesValue(argument actionArgument) any {
	if argument.valueType == Variable {
		return paramValue(argument, Variable)
	}

	var names, isArray = getArgValue(argument).([]interface{})
	if !isArray {
		parserError("World clock cities must be an array of city names or a variable")
	}

	var cities []any
	for _, name := range names {
		cities = append(cities, worldClockCityValue(actionArgument{valueType: String, value: name}))
	}
	if len(cities) == 1 {
		return cities[0]
	}

	return cities
}
