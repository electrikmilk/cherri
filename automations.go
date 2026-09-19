package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type AutomationTrigger struct {
	WFTriggerIdentifier           string
	WFTriggerSerializedParameters map[string]any
	WFTriggerUUID                 string
}

var automationTriggers []AutomationTrigger

var triggerIdentifiers = map[string]string{
	"bluetooth":    "WFBluetoothTrigger",
	"wifi":         "WFWifiTrigger",
	"stageManager": "WFStageManagerTrigger",
	"display":      "WFExternalDisplayTrigger",
	"app":          "WFAppInFocusTrigger",
	"screenshot":   "WFScreenshotTrigger",
	"battery":      "WFBatteryLevelTrigger",
	"charging":     "WFPlugInTrigger",
}

var validScreenshotLocations = []string{"photos", "files", "clipboard"}
var validStageManagerTypes = []string{"on", "off", "both"}
var validWifiConnectionTypes = []string{"joined", "disconnected", "both"}
var validConnectionTypes = []string{"connect", "disconnect", "both"}

func checkTriggerIdentifier(identifier string) {
	if triggerIdentifiers[identifier] == "" {
		var list = makeKeyList("Available triggers:", triggerIdentifiers, identifier)
		parserError(fmt.Sprintf("Invalid trigger identifier '%s'\n\n%s", identifier, list))
	}
}

func checkTriggerValue(validValues []string, value string) {
	if !slices.Contains(validValues, value) {
		var list = makeValueList("Available values:", validValues, value)
		parserError(fmt.Sprintf("Invalid trigger param value '%s'\n\n%s", value, list))
	}
}

func collectTrigger() {
	if !strings.Contains(lookAheadUntil('\n'), " ") {
		var collectIdentifier = collectUntil('\n')
		checkTriggerIdentifier(collectIdentifier)

		parserError("Expected trigger parameter(s), got only identifier")
	}

	var collectIdentifier = collectUntil(' ')
	if collectIdentifier == "" {
		parserError("Expected trigger identifier")
	}
	checkTriggerIdentifier(collectIdentifier)

	advance()

	var triggerParams = make(map[string]any)
	switch collectIdentifier {
	case "screenshot":
		var paramValue = collectUntil('\n')
		var screenshotLocations = strings.Split(paramValue, ",")
		for i, location := range screenshotLocations {
			var trimmedLocation = strings.TrimSpace(location)
			checkTriggerValue(validScreenshotLocations, trimmedLocation)
			screenshotLocations[i] = trimmedLocation
		}

		triggerParams["ScreenshotLocations"] = screenshotLocations
	case "battery":
		var batteryLevel = collectUntil('\n')
		var batteryLevelFloat, err = strconv.ParseFloat(batteryLevel, 64)
		if err != nil {
			parserError("Invalid battery level: " + batteryLevel)
		}
		// -0.01 is to Compensate for battery level being required to be an actual decimal value.
		triggerParams["WFBatteryLevel"] = batteryLevelFloat - 0.01
	case "stageManager":
		var stageManagerType = collectUntil('\n')
		checkTriggerValue(validStageManagerTypes, stageManagerType)
		triggerParams["WFStageManagerType"] = stageManagerType
	case "wifi":
		var connectionType = collectUntil('\n')
		checkTriggerValue(validWifiConnectionTypes, connectionType)
		triggerParams["WFConnectionType"] = connectionType
	case "bluetooth":
		var connectionType = collectUntil('\n')
		checkTriggerValue(validConnectionTypes, connectionType)
		triggerParams["WFBluetoothConnectionType"] = connectionType
	case "display":
		var connectionType = collectUntil('\n')
		checkTriggerValue(validConnectionTypes, connectionType)
		triggerParams["WFConnectionType"] = connectionType
	case "charging":
		var chargingType = collectUntil('\n')
		checkTriggerValue(validConnectionTypes, chargingType)
		triggerParams["WFConnectionType"] = chargingType
	case "app":
		var appIdentifier = collectUntil('\n')

		triggerParams["WFSelectedApps"] = map[string]any{
			"AppIntentDescriptor": map[string]string{
				"TeamIdentifier":   "0000000000",
				"BundleIdentifier": replaceAppID(appIdentifier),
			},
		}
	}

	var triggerUUID = createUUID(&collectIdentifier)

	automationTriggers = append(automationTriggers, AutomationTrigger{
		WFTriggerIdentifier:           triggerIdentifiers[collectIdentifier],
		WFTriggerSerializedParameters: triggerParams,
		WFTriggerUUID:                 triggerUUID,
	})
}

func printAutomationsDebug() {
	fmt.Println(ansi("## AUTOMATIONS ##", bold))
	for _, trigger := range automationTriggers {
		fmt.Printf("Identifier: %s\nParams: %v\nUUID: %s\n---", trigger.WFTriggerIdentifier, trigger.WFTriggerSerializedParameters, trigger.WFTriggerUUID)
	}
	fmt.Print("\n\n")
}
