/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"

	"github.com/electrikmilk/args-parser"
)

var actionCategories = []string{"basic"}
var currentCategory string

type selfDoc struct {
	title       string
	description string
	category    string
}

func generateDocs() {
	defineToggleSetActions()
	defineRawAction()
	loadActionsByCategory()
	args.Args["no-ansi"] = ""
	var cat = args.Value("docs")
	for _, category := range actionCategories {
		if cat != "" && cat != category {
			continue
		}
		if cat == "" {
			if category == "pdf" {
				category = "PDF"
			}
			fmt.Printf("\n# %s Actions\n\n", capitalize(category))
		}

		generateCategory(category)
	}
}

func generateCategory(category string) {
	var i = 0
	for name, def := range actions {
		if def.doc.category != category {
			continue
		}
		if i != 0 {
			fmt.Print("\n---\n\n")
		}

		currentAction = *def
		currentActionIdentifier = name

		fmt.Println(generateActionDefinition(parameterDefinition{}, true))
		i++
	}
}
