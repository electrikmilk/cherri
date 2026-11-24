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

func init() {
	var cities []worldClockCity
	var nameCounts = map[string]int{}
	for _, line := range strings.Split(strings.TrimSpace(worldClockCitiesData), "\n") {
		var identifier, name, found = strings.Cut(line, "\t")
		if !found {
			panic(fmt.Sprintf("Invalid world clock city line '%s'", line))
		}
		cities = append(cities, worldClockCity{identifier: identifier, name: name})
		nameCounts[name]++
	}

	var values []string
	for _, city := range cities {
		var value = city.name
		// Cities that share a name are told apart by their identifier, but Shortcuts shows only the name.
		if nameCounts[city.name] > 1 {
			value = fmt.Sprintf("%s (%s)", city.name, city.identifier)
		}
		worldClockCities[value] = city
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

// decompWorldClockCity recovers the city argument, which is a variable or the name of the city of the identifier
// Shortcuts stores.
func decompWorldClockCity(action *ShortcutAction) map[string]string {
	var value, found = optionalParam(action, "city")
	if !found {
		return nil
	}

	return map[string]string{"city": decompCityValue(value, action)}
}

// decompWorldClockCities recovers the cities argument, which Shortcuts writes as a single city when there is only one.
func decompWorldClockCities(action *ShortcutAction) map[string]string {
	var value, found = optionalParam(action, "entities")
	if !found {
		return nil
	}

	if cities, isArray := value.([]interface{}); isArray {
		var names []interface{}
		for _, city := range cities {
			names = append(names, decompCityName(city, action))
		}

		return map[string]string{"cities": decompLiteral(names)}
	}
	if isReferenceValue(value) {
		return map[string]string{"cities": decompValue(value)}
	}

	return map[string]string{"cities": decompLiteral([]interface{}{decompCityName(value, action)})}
}

// decompCityValue writes a city as the name of its value in the worldClockCity enum, or a variable as it is.
func decompCityValue(value any, action *ShortcutAction) string {
	if isReferenceValue(value) {
		return decompValue(value)
	}

	return decompValue(decompCityName(value, action))
}

// decompCityName finds the enum value of the city with the identifier Shortcuts stores.
func decompCityName(value any, action *ShortcutAction) string {
	var city, isMap = value.(map[string]interface{})
	var identifier, hasIdentifier = city["identifier"].(string)
	if !isMap || !hasIdentifier {
		decompError(fmt.Sprintf("Expected a world clock city to be a dictionary with an identifier, got %T", value), action)
	}

	for name, known := range worldClockCities {
		if known.identifier == identifier {
			return name
		}
	}
	decompError(fmt.Sprintf("Expected the world clock city identifier '%s' to be one that Cherri knows", identifier), action)

	return ""
}
