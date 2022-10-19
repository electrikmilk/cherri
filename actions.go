/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func makeActions() {
	actions = make(map[string]actionDefinition)
	actions["wait"] = actionDefinition{
		ident: "delay",
		args: []argumentDefinition{
			{
				field:     "seconds",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFDelayTime", args, 0),
			}
		},
	}
	actions["alert"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "alert",
				validType: String,
			},
			{
				field:     "title",
				validType: String,
			},
			{
				field:     "cancelButton",
				validType: Bool,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFAlertActionMessage", args, 0),
				argumentValue("WFAlertActionTitle", args, 1),
				argumentValue("WFAlertActionCancelButtonShown", args, 2),
			}
		},
	}
	actions["share"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "input",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
			}
		},
	}
	actions["askForInput"] = actionDefinition{
		ident: "ask",
		args: []argumentDefinition{
			{
				field:     "inputType",
				validType: String,
			},
			{
				field:     "prompt",
				validType: String,
			},
			{
				field:     "defaultValue",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInputType", args, 0),
				argumentValue("WFAskActionPrompt", args, 1),
				argumentValue("WFAskActionDefaultAnswer", args, 2),
			}
		},
	}
	actions["getBatteryLevel"] = actionDefinition{}
	actions["getCurrentLocation"] = actionDefinition{
		ident: "location",
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFLocation",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "isCurrentLocation",
							dataType: Boolean,
							value:    true,
						},
					},
				},
			}
		},
	}
	actions["splitText"] = actionDefinition{
		ident: "text.split",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
			{
				field:     "glue",
				validType: String,
			},
		},
		call: textParts,
	}
	actions["combineText"] = actionDefinition{
		ident: "text.combine",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
			{
				field:     "glue",
				validType: String,
			},
		},
		call: textParts,
	}
	actions["getType"] = actionDefinition{
		ident: "getitemtype",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
			}
		},
	}
	actions["copyToClipboard"] = actionDefinition{
		ident: "setclipboard",
		args: []argumentDefinition{
			{
				field:     "value",
				validType: Variable,
			},
			{
				field:     "local",
				validType: Bool,
			},
			{
				field:     "expire",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
				argumentValue("WFLocalOnly", args, 1),
				argumentValue("WFExpirationDate", args, 2),
			}
		},
	}
	actions["getClipboard"] = actionDefinition{}
	actions["chooseFromList"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "dictionary",
				validType: Dict,
			},
			{
				field:     "prompt",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFInput", args, 0),
				argumentValue("WFChooseFromListActionPrompt", args, 1),
			}
		},
	}
	actions["chooseMultipleFromList"] = actionDefinition{
		ident: "choosefromlist",
		args: []argumentDefinition{
			{
				field:     "dictionary",
				validType: Dict,
			},
			{
				field:     "prompt",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFChooseFromListActionSelectMultiple",
					dataType: Boolean,
					value:    true,
				},
				argumentValue("WFInput", args, 0),
				argumentValue("WFChooseFromListActionPrompt", args, 1),
				argumentValue("WFChooseFromListActionSelectAll", args, 2),
			}
		},
	}
	actions["replaceText"] = actionDefinition{
		ident: "text.replace",
		args: []argumentDefinition{
			{
				field:     "find",
				validType: String,
			},
			{
				field:     "replacement",
				validType: String,
			},
			{
				field:     "subject",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return replaceText(true, false, args)
		},
	}
	actions["iReplaceText"] = actionDefinition{
		ident: "text.replace",
		args: []argumentDefinition{
			{
				field:     "find",
				validType: String,
			},
			{
				field:     "replacement",
				validType: String,
			},
			{
				field:     "subject",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return replaceText(false, false, args)
		},
	}
	actions["regReplaceText"] = actionDefinition{
		ident: "text.replace",
		args: []argumentDefinition{
			{
				field:     "expression",
				validType: String,
			},
			{
				field:     "replacement",
				validType: String,
			},
			{
				field:     "subject",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return replaceText(true, true, args)
		},
	}
	actions["iRegReplaceText"] = actionDefinition{
		ident: "text.replace",
		args: []argumentDefinition{
			{
				field:     "expression",
				validType: String,
			},
			{
				field:     "replacement",
				validType: String,
			},
			{
				field:     "subject",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return replaceText(false, true, args)
		},
	}
	actions["uppercase"] = actionDefinition{
		ident: "text.changecase",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return changeCase("UPPERCASE", args)
		},
	}
	actions["lowercase"] = actionDefinition{
		ident: "text.changecase",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return changeCase("lowercase", args)
		},
	}
	actions["titleCase"] = actionDefinition{
		ident: "text.changecase",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return changeCase("Capitalize with Title Case", args)
		},
	}
	actions["capitalize"] = actionDefinition{
		ident: "text.changecase",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return changeCase("Capitalize with sentence case", args)
		},
	}
	actions["capitalizeAll"] = actionDefinition{
		ident: "text.changecase",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return changeCase("Capitalize Every Word", args)
		},
	}
	actions["alternateCase"] = actionDefinition{
		ident: "text.changecase",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return changeCase("cApItAlIzE wItH aLtErNaTiNg cAsE", args)
		},
	}
	actions["quicklook"] = actionDefinition{
		ident: "previewdocument",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				inputValue("WFInput", args[0].value.(string), ""),
			}
		},
	}
	actions["translateFrom"] = actionDefinition{
		ident: "translate",
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
			{
				field:     "from",
				validType: String,
			},
			{
				field:     "to",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				inputValue("WFInputText", args[0].value.(string), ""),
				argumentValue("WFSelectedFromLanguage", args, 0),
				argumentValue("WFSelectedLanguage", args, 0),
			}
		},
	}
	actions["translate"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
			{
				field:     "to",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				inputValue("WFInputText", args[0].value.(string), ""),
				{
					key:      "WFSelectedFromLanguage",
					dataType: Text,
					value:    "Detect Language",
				},
				argumentValue("WFSelectedLanguage", args, 0),
			}
		},
	}
	actions["detectLanguage"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "input",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				inputValue("WFInput", args[0].value.(string), ""),
			}
		},
	}
	actions["correctSpelling"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "text",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Show-text",
					dataType: Boolean,
					value:    true,
				},
				argumentValue("text", args, 0),
			}
		},
	}
	actions["getShortcuts"] = actionDefinition{
		ident: "getmyworkflows",
	}
	actions["getKeys"] = actionDefinition{
		ident: "getvalueforkey",
		args: []argumentDefinition{
			{
				field:     "dictionary",
				validType: Dict,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetDictionaryValueType",
					dataType: Text,
					value:    "All Keys",
				},
				argumentValue("WFInput", args, 0),
			}
		},
	}
	actions["getValues"] = actionDefinition{
		ident: "getvalueforkey",
		args: []argumentDefinition{
			{
				field:     "dictionary",
				validType: Dict,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetDictionaryValueType",
					dataType: Text,
					value:    "All Values",
				},
				argumentValue("WFInput", args, 0),
			}
		},
	}
	actions["getValue"] = actionDefinition{
		ident: "getvalueforkey",
		args: []argumentDefinition{
			{
				field:     "dictionary",
				validType: Dict,
			},
			{
				field:     "key",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFGetDictionaryValueType",
					dataType: Text,
					value:    "Value",
				},
				argumentValue("WFInput", args, 0),
				argumentValue("WFDictionaryKey", args, 1),
			}
		},
	}
	actions["open"] = actionDefinition{
		ident: "openworkflow",
		args: []argumentDefinition{
			{
				field:     "shortcutName",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFWorkflow",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "workflowIdentifier",
							dataType: Text,
							value:    shortcutsUUID(),
						},
						{
							key:      "isSelf",
							dataType: Boolean,
							value:    false,
						},
						argumentValue("workflowName", args, 0),
					},
				},
			}
		},
	}
	actions["run"] = actionDefinition{
		ident: "runworkflow",
		args: []argumentDefinition{
			{
				field:     "shortcutName",
				validType: String,
			},
			{
				field:     "output",
				validType: Variable,
			},
			{
				field:     "isSelf",
				validType: Bool,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFWorkflow",
					dataType: Dictionary,
					value: []plistData{
						{
							key:      "workflowIdentifier",
							dataType: Text,
							value:    shortcutsUUID(),
						},
						{
							key:      "isSelf",
							dataType: Boolean,
							value:    false,
						},
						argumentValue("workflowName", args, 0),
					},
				},
				inputValue("WFInput", args[1].value.(string), ""),
			}
		},
	}
	actions["list"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "listItem",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			var listItems []string
			for _, item := range args {
				listItems = append(listItems, plistValue(Text, item.value))
			}
			return []plistData{
				{
					key:      "WFItems",
					dataType: Array,
					value:    listItems,
				},
			}
		},
	}
	actions["url"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "url",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			var urlItems []string
			for _, item := range args {
				urlItems = append(urlItems, plistValue(Text, item.value))
			}
			return []plistData{
				{
					key:      "Show-WFURLActionURL",
					dataType: Boolean,
					value:    true,
				},
				{
					key:      "WFURLActionURL",
					dataType: Array,
					value:    urlItems,
				},
			}
		},
	}
	actions["hash"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "type",
				validType: String,
			},
			{
				field:     "input",
				validType: Variable,
			},
		},
		check: func(args []actionArgument) {
			var hashType = strings.ToUpper(getArgValue(args[0]).(string))
			if !contains(hashTypes, hashType) {
				parserError(fmt.Sprintf("Invalid hash type of '%s'. Available hash types: %v", hashType, hashTypes))
			}
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFHashType", args, 0),
				inputValue("WFInput", args[1].value.(string), ""),
			}
		},
	}
	actions["formatNumber"] = actionDefinition{
		ident: "format.number",
		args: []argumentDefinition{
			{
				field:     "number",
				validType: Integer,
			},
			{
				field:     "decimalPlaces",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFNumber", args, 0),
				argumentValue("WFNumberFormatDecimalPlaces", args, 1),
			}
		},
	}
	actions["randomNumber"] = actionDefinition{
		ident: "number.random",
		args: []argumentDefinition{
			{
				field:     "min",
				validType: Integer,
			},
			{
				field:     "max",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFRandomNumberMinimum", args, 0),
				argumentValue("WFRandomNumberMaximum", args, 1),
			}
		},
	}
	actions["base64Encode"] = actionDefinition{
		ident: "base64encode",
		args: []argumentDefinition{
			{
				field:     "encodeInput",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFEncodeMode",
					dataType: Text,
					value:    "Encode",
				},
				inputValue("WFInput", args[0].value.(string), ""),
			}
		},
	}
	actions["base64Decode"] = actionDefinition{
		ident: "base64encode",
		args: []argumentDefinition{
			{
				field:     "decodeInput",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "WFEncodeMode",
					dataType: Text,
					value:    "Decode",
				},
				inputValue("WFInput", args[0].value.(string), ""),
			}
		},
	}
	actions["extractArchive"] = actionDefinition{
		ident: "unzip",
		args: []argumentDefinition{
			{
				field:     "file",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				inputValue("WFArchive", args[0].value.(string), ""),
			}
		},
	}
	actions["makeArchive"] = actionDefinition{
		ident: "makezip",
		args: []argumentDefinition{
			{
				field:     "name",
				validType: String,
			},
			{
				field:     "format",
				validType: String,
			},
			{
				field:     "files",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFZIPName", args, 0),
				argumentValue("WFArchiveFormat", args, 1),
				inputValue("WFInput", args[2].value.(string), ""),
			}
		},
	}
	actions["show"] = actionDefinition{
		ident: "showresult",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("Text", args, 0),
			}
		},
	}
	actions["waitToReturn"] = actionDefinition{
		ident: "waittoreturn",
		call: func(args []actionArgument) []plistData {
			return []plistData{}
		},
	}
	actions["notification"] = actionDefinition{
		ident: "notification",
		args: []argumentDefinition{
			{
				field:     "body",
				validType: String,
			},
			{
				field:     "title",
				validType: String,
			},
			{
				field:     "sound",
				validType: Bool,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				argumentValue("WFNotificationActionBody", args, 0),
				argumentValue("WFNotificationActionTitle", args, 1),
				argumentValue("WFNotificationActionSound", args, 2),
			}
		},
	}
	actions["stop"] = actionDefinition{
		ident: "exit",
	}
	actions["nothing"] = actionDefinition{}
	actions["setVolume"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "volume",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			args[0].value = fmt.Sprintf("0.%s", args[0].value)
			return []plistData{
				argumentValue("WFVolume", args, 0),
			}
		},
	}
	actions["setBrightness"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "brightness",
				validType: Integer,
			},
		},
		call: func(args []actionArgument) []plistData {
			args[0].value = fmt.Sprintf("0.%s", args[0].value)
			return []plistData{
				argumentValue("WFBrightness", args, 0),
			}
		},
	}
	actions["getName"] = actionDefinition{
		ident: "getitemname",
		args: []argumentDefinition{
			{
				field:     "item",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				inputValue("WFInput", args[0].value.(string), ""),
			}
		},
	}
	actions["setName"] = actionDefinition{
		ident: "setitemname",
		args: []argumentDefinition{
			{
				field:     "item",
				validType: Variable,
			},
			{
				field:     "name",
				validType: String,
			},
			{
				field:     "includeFileExtension",
				validType: Bool,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				inputValue("WFInput", args[0].value.(string), ""),
				argumentValue("WFName", args, 1),
				argumentValue("WFDontIncludeFileExtension", args, 2),
			}
		},
	}
	actions["countItems"] = actionDefinition{
		ident: "count",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return count("Items", args)
		},
	}
	actions["countChars"] = actionDefinition{
		ident: "count",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return count("Characters", args)
		},
	}
	actions["countWords"] = actionDefinition{
		ident: "count",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return count("Words", args)
		},
	}
	actions["countSentences"] = actionDefinition{
		ident: "count",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return count("Sentences", args)
		},
	}
	actions["countLines"] = actionDefinition{
		ident: "count",
		args: []argumentDefinition{
			{
				field:     "input",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return count("Lines", args)
		},
	}
	actions["toggleAppearance"] = actionDefinition{
		ident: "appearance",
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "operation",
					dataType: Text,
					value:    "toggle",
				},
			}
		},
	}
	actions["lightMode"] = actionDefinition{
		ident: "appearance",
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "operation",
					dataType: Text,
					value:    "set",
				},
				{
					key:      "style",
					dataType: Text,
					value:    "light",
				},
			}
		},
	}
	actions["darkMode"] = actionDefinition{
		ident: "appearance",
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "operation",
					dataType: Text,
					value:    "set",
				},
				{
					key:      "style",
					dataType: Text,
					value:    "dark",
				},
			}
		},
	}
	actions["openURL"] = actionDefinition{
		args: []argumentDefinition{
			{
				field:     "url",
				validType: Variable,
			},
		},
		call: func(args []actionArgument) []plistData {
			return []plistData{
				{
					key:      "Show-WFInput",
					dataType: Boolean,
					value:    true,
				},
				inputValue("WFInput", args[0].value.(string), ""),
			}
		},
	}
	actions["makeVCard"] = actionDefinition{
		ident: "gettext",
		args: []argumentDefinition{
			{
				field:     "title",
				validType: String,
			},
			{
				field:     "subtitle",
				validType: String,
			},
			{
				field:     "url",
				validType: String,
			},
		},
		call: func(args []actionArgument) []plistData {
			var vcard = "BEGIN:VCARD\nVERSION:3.0\n"
			vcard += "N;CHARSET=utf-8:" + args[0].value.(string) + "\n"
			vcard += "ORG:" + args[1].value.(string) + "\nPHOTO;ENCODING=b:"
			bytes, readErr := os.ReadFile(getArgValue(args[2]).(string))
			handle(readErr)
			vcard += base64.StdEncoding.EncodeToString(bytes) + "\nEND:VCARD"
			args[0] = actionArgument{
				valueType: String,
				value:     vcard,
			}
			return []plistData{
				argumentValue("WFTextActionText", args, 0),
			}
		},
	}
}
