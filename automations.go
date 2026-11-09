package main

import (
	"fmt"
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
	"screenshot":   "WFScreenshotTrigger",
	"battery":      "WFBatteryLevelTrigger",
	"charger":      "WFPlugInTrigger",
}

var validScreenshotLocations = []string{"photos", "files", "clipboard"}
var validStageManagerStates = []string{"on", "off", "both"}
var validWifiConnectionStatus = []string{"joined", "disconnected", "both"}
var validConnectionStatus = []string{"connect", "disconnect", "both"}

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
		var list = valueList{name: "screenshot location", list: &validScreenshotLocations}
		triggerParams["ScreenshotLocations"] = list.parseList('\n')
	case "battery":
		if char == '>' {
			triggerParams["WFBatteryLevelComparison"] = "rises_above"
			advance()
		} else if char == '<' {
			triggerParams["WFBatteryLevelComparison"] = "falls_below"
			advance()
		}

		var batteryLevel = collectUntil('\n')
		var batteryLevelFloat, err = strconv.ParseFloat(batteryLevel, 64)
		if err != nil {
			parserError("Invalid battery level: " + batteryLevel)
		}
		// -0.01 is to Compensate for battery level being required to be an actual decimal value.
		triggerParams["WFBatteryLevel"] = batteryLevelFloat - 0.01
	case "stageManager":
		var list = valueList{name: "Stage Manager state", list: &validStageManagerStates}
		triggerParams["WFStageManagerType"] = list.parse('\n')
	case "wifi":
		var list = valueList{name: "WiFi connection status", plural: "WiFi connection statuses", list: &validWifiConnectionStatus}
		triggerParams["WFConnectionType"] = list.parse('\n')
	case "bluetooth":
		var list = valueList{name: "Bluetooth connection status", plural: "Bluetooth connection statuses", list: &validConnectionStatus}
		triggerParams["WFBluetoothConnectionType"] = list.parse('\n')
	case "display":
		var list = valueList{name: "display connection status", plural: "display connection statuses", list: &validConnectionStatus}
		triggerParams["WFConnectionType"] = list.parse('\n')
	case "charger":
		var list = valueList{name: "charger connection status", plural: "charger connection statuses", list: &validConnectionStatus}
		triggerParams["WFConnectionType"] = list.parse('\n')
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
