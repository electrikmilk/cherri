package main

func defineRawAction() {
	actions["rawAction"] = &actionDefinition{
		parameters: []parameterDefinition{
			{
				name:      "identifier",
				validType: String,
				literal:   true,
			},
			{
				name:      "parameters",
				optional:  true,
				validType: Dict,
				literal:   true,
			},
		},
		check: func(args []actionArgument, _ *actionDefinition) {
			if len(args) == 1 {
				return
			}

			var params, hasParams = args[1].value.(map[string]any)
			if !hasParams {
				return
			}

			// Shortcuts crashes while loading a Shortcut with an empty descriptor, which looks like our output is wrong.
			if params["AppIntentDescriptor"] == "" {
				parserError("An empty 'AppIntentDescriptor' is an invalid Shortcut format and signing will fail. Remove it or give it a valid intent value.")
			}
		},
		makeParams: func(args []actionArgument) map[string]any {
			if len(args) == 1 {
				return map[string]any{}
			}

			var value = args[1].value
			if value == nil {
				return map[string]any{}
			}

			var params = value.(map[string]interface{})
			handleRawParams(params)

			return params
		},
		outputType: Variable,
	}
}

func handleRawParams(params map[string]any) {
	for key, value := range params {
		if reference, isReference := value.(varValue); isReference {
			params[key] = variableValue(reference)
			continue
		}

		params[key] = normalizeRawActionParamValue(value)
	}
}

func normalizeRawActionParamValue(value any) any {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		if !hasBraces(v) {
			return value
		}
		return attachmentValues(v)
	case map[string]any:
		var normalized any = normalizeRawActionDictionary(v)
		return makeDictionaryValue(&normalized, true)
	case []any:
		return makeRawActionArrayValue(v)
	default:
		return value
	}
}

func normalizeRawActionDictionary(value map[string]any) map[string]any {
	var normalized = make(map[string]any, len(value))
	for key, item := range value {
		normalized[key] = normalizeRawActionParamValue(item)
	}
	return normalized
}

func makeRawActionArrayValue(value []any) WFArrayValue {
	var items = make([]WFDictionaryFieldValueItem, 0, len(value))
	for _, item := range value {
		items = append(items, makeDictionaryItem("", normalizeRawActionParamValue(item), true))
	}
	return makeArrayValue(items)
}
