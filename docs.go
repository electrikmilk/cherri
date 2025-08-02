/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/electrikmilk/args-parser"
)

var actionCategories = []string{"basic"}
var currentCategory string

type selfDoc struct {
	title       string
	description string
	category    string
	subcategory string
}

type actionCategory struct {
	title         string
	actions       []string
	subcategories map[string][]string
}

func generateDocs() {
	defineToggleSetActions()
	defineRawAction()
	loadActionsByCategory()
	var cat = args.Value("docs")
	for i, category := range actionCategories {
		if cat != "" && cat != category {
			continue
		}
		var actionCategory = generateCategory(category)
		fmt.Println("#", actionCategory.title)

		slices.Sort(actionCategory.actions)
		fmt.Println(strings.Join(actionCategory.actions, "\n\n---\n"))

		if actionCategory.subcategories != nil {
			printCategories(actionCategory.subcategories)
		}

		if cat == "" && i != 0 {
			fmt.Print("---\n")
		}
	}
}

func printCategories(categories map[string][]string) {
	var keys []string
	for k := range categories {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		var category = categories[k]
		fmt.Printf("\n## %s\n", k)
		slices.Sort(category)
		fmt.Println(strings.Join(category, "\n\n---\n"))
	}
}

func generateCategory(category string) actionCategory {
	var categoryTitle = category
	if categoryTitle == "pdf" {
		categoryTitle = "PDF"
	}
	var cat = actionCategory{
		title: fmt.Sprintf("%s Actions", capitalize(categoryTitle)),
	}
	var subcat = args.Value("subcat")
	for name, def := range actions {
		if def.doc.category != category || (subcat != "" && def.doc.subcategory != subcat) {
			continue
		}

		currentAction = *def
		currentActionIdentifier = name
		var definition = generateActionDefinition(parameterDefinition{}, true)

		if def.doc.subcategory != "" {
			if cat.subcategories == nil {
				cat.subcategories = make(map[string][]string)
			}
			cat.subcategories[def.doc.subcategory] = append(cat.subcategories[def.doc.subcategory], definition)
			continue
		}

		cat.actions = append(cat.actions, definition)
	}
	return cat
}
