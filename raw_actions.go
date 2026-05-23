package main

import (
	"regexp"
	"strings"
)

var rawActionVariableValueRegex = regexp.MustCompile(`^\$\{@?.*}$`)
var rawActionQuantityFieldValueSerializationType = "WFQuantityFieldValue"

func defineRawAction() {
	actions["rawAction"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "identifier",
				validType: String,
			},
			{
				name:      "parameters",
				optional:  true,
				validType: Dict,
			},
		},
		makeParams: func(args []actionArgument) map[string]any {
			if len(args) == 1 {
				return map[string]any{}
			}

			var params = getArgValue(args[1]).(map[string]interface{})
			handleRawParams(params)

			return params
		},
	}
}

func handleRawParams(params map[string]any) {
	for key, value := range params {
		params[key] = normalizeRawActionParamValue(value)
	}
}

func normalizeRawActionParamValue(value any) any {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		if !strings.ContainsAny(v, "{}") {
			return value
		}
		if rawActionVariableValueRegex.MatchString(v) {
			return variableValue(varValue{
				value: rawActionVariableIdentifier(v),
			})
		}
		return attachmentValues(v)
	case map[string]any:
		if isRawActionSerializedValue(v, rawActionQuantityFieldValueSerializationType) {
			return makeRawActionQuantityFieldValue(v)
		}
		var normalized any = normalizeRawActionDictionary(v)
		return makeDictionaryValue(&normalized)
	case []any:
		return makeRawActionArrayValue(v)
	default:
		return value
	}
}

func rawActionVariableIdentifier(value string) string {
	var identifier = strings.TrimPrefix(strings.TrimSuffix(value, "}"), "${")
	return strings.TrimPrefix(identifier, "@")
}

func isRawActionSerializedValue(value map[string]any, serializationType string) bool {
	var valueSerializationType, ok = value["WFSerializationType"].(string)
	return ok && valueSerializationType == serializationType
}

func normalizeRawActionDictionary(value map[string]any) map[string]any {
	var normalized = make(map[string]any, len(value))
	for key, item := range value {
		normalized[key] = normalizeRawActionParamValue(item)
	}
	return normalized
}

func makeRawActionQuantityFieldValue(value map[string]any) WFQuantityFieldValue {
	var quantityValue = WFQuantityValue{}
	var rawValue, ok = value["Value"].(map[string]any)
	if ok {
		if magnitude, found := rawValue["Magnitude"]; found {
			quantityValue.Magnitude = normalizeRawActionQuantityValue(magnitude)
		}
		if unit, found := rawValue["Unit"]; found {
			quantityValue.Unit = normalizeRawActionQuantityValue(unit)
		}
	}

	return WFQuantityFieldValue{
		Value:               quantityValue,
		WFSerializationType: rawActionQuantityFieldValueSerializationType,
	}
}

func normalizeRawActionQuantityValue(value any) any {
	var stringValue, ok = value.(string)
	if !ok || !rawActionVariableValueRegex.MatchString(stringValue) {
		return value
	}

	return variableValueWithSerialization(varValue{
		value: rawActionVariableIdentifier(stringValue),
	}, "")
}

func makeRawActionArrayValue(value []any) WFArrayValue {
	var items = make([]WFDictionaryFieldValueItem, 0, len(value))
	for _, item := range value {
		items = append(items, makeDictionaryItem("", normalizeRawActionParamValue(item)))
	}
	return makeArrayValue(items)
}
