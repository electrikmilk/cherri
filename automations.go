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
	"wifi":         "WFWiFiTrigger",
	"stageManager": "WFStageManagerTrigger",
	"display":      "WFExternalDisplayTrigger",
	"app":          "WFAppInFocusTrigger",
	"screenshot":   "WFScreenshotTrigger",
	"battery":      "WFBatteryLevelTrigger",
	"charging":     "WFPlugInTrigger",
	"focus":        "WFUserFocusActivityTrigger",
}

var validScreenshotLocations = []string{"photos", "files", "clipboard"}
var validStageManagerTypes = []string{"on", "off", "both"}

func checkTriggerIdentifier(identifier string) {
	if triggerIdentifiers[identifier] == "" {
		var list = makeKeyList("Available triggers:", triggerIdentifiers, identifier)
		parserError(fmt.Sprintf("Invalid trigger identifier '%s'\n\n%s", identifier, list))
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
			if slices.Contains(validScreenshotLocations, trimmedLocation) {
				screenshotLocations[i] = trimmedLocation
			} else {
				parserError("Invalid screenshot location (not photos, files, or clipboard): " + trimmedLocation)
			}
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
		if slices.Contains(validStageManagerTypes, stageManagerType) {
			parserError("Invalid stage manager type (not on, off, or both): " + stageManagerType)
		}
		triggerParams["WFStageManagerType"] = stageManagerType
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
