package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/electrikmilk/args-parser"
	_ "github.com/glebarez/go-sqlite"
	"howett.net/plist"
)

/*
Reads Shortcuts ToolKit SQLite DB to define actions and their parameters for use and more.
*/

var shortcutsDBPath = os.ExpandEnv("$HOME/Library/Shortcuts/ToolKit/Tools-active")
var toolkit *sql.DB
var imported []string

func handleImports() {
	var matches = copyPasteRegex.FindAllStringSubmatch(contents, -1)
	if len(matches) == 0 {
		return
	}

	parseImports()

	if args.Using("debug") {
		printPasteablesDebug()
	}
}

func parseImports() {
	pasteables = make(map[string]string)
	for char != -1 {
		switch {
		case char == '"':
			collectString()
			advanceUntil('\n')
		case commentAhead():
			collectComment()
		case startOfLineTokenAhead(Import):
			var importPath = collectImport()
			importActions(importPath)
		}
		advance()
	}

	resetParse()
}

func collectImport() string {
	var lineRef = newLineReference()
	skipWhitespace()

	if char != '\'' {
		parserError(fmt.Sprintf("Expected raw string ('), got: %c", char))
	}

	advance()

	var path = collectRawString()

	lineRef.replaceLines()

	return path
}

func connectToolkitDB() {
	if toolkit != nil {
		return
	}

	if args.Using("toolkit") {
		shortcutsDBPath = args.Value("toolkit")
	}

	var db, dbErr = sql.Open("sqlite", shortcutsDBPath)
	handle(dbErr)

	toolkit = db
}

var imported []string

func importActions(identifier string) {
	if slices.Contains(imported, identifier) {
		parserError(fmt.Sprintf("Actions and types from '%s' have already been imported.", identifier))
	}

	connectToolkitDB()

	if !appIdentifierRegex.MatchString(identifier) {
		matchApplication(&identifier)
	}

	fmt.Println("### ACTIONS ###")

	var actions, actionsErr = getActions(identifier)
	handle(actionsErr)

	fmt.Println(actions)

	fmt.Println("### ENUMS ###")

	var enums, enumsErr = getEnums(identifier)
	handle(enumsErr)

	fmt.Println(enums)

	imported = append(imported, identifier)
}

type AppInfo struct {
	CFBundleIdentifier string
}

func matchApplication(identifier *string) {
	var apps, readErr = os.ReadDir("/Applications")
	handle(readErr)
	for _, app := range apps {
		var appName = strings.Replace(app.Name(), ".app", "", 1)
		if appName == *identifier {
			var info AppInfo
			var infoBytes, infoErr = os.ReadFile("/Applications/" + app.Name() + "/Contents/Info.plist")
			handle(infoErr)
			var decodeErr = plist.NewDecoder(bytes.NewReader(infoBytes)).Decode(&info)
			handle(decodeErr)

			*identifier = info.CFBundleIdentifier
			return
		}
	}

	parserError(fmt.Sprintf("Could not find '%s' in /Applications/.", *identifier))
}

type actionTool struct {
	rowId                    sql.NullString
	id                       sql.NullString
	toolType                 sql.NullString
	flags                    sql.NullInt64
	visibilityFlags          sql.NullInt64
	requirements             sql.NullString
	authenticationPolicy     sql.NullString
	customIcon               sql.NullString
	deprecationReplacementId sql.NullString
	sourceActionProvider     sql.NullString
	outputTypeInstance       sql.NullString
	sourceContainerId        sql.NullInt64
	attributionContainerId   sql.NullInt64
}

func getActions(idPattern string) ([]actionTool, error) {
	var query = `select * from Tools WHERE id LIKE ?`

	var rows, err = toolkit.Query(query, idPattern+".%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tools []actionTool
	for rows.Next() {
		var tool actionTool

		if err := rows.Scan(
			&tool.rowId,
			&tool.id,
			&tool.toolType,
			&tool.flags,
			&tool.visibilityFlags,
			&tool.requirements,
			&tool.authenticationPolicy,
			&tool.customIcon,
			&tool.deprecationReplacementId,
			&tool.sourceActionProvider,
			&tool.outputTypeInstance,
			&tool.sourceContainerId,
			&tool.attributionContainerId,
		); err != nil {
			return nil, err
		}

		tools = append(tools, tool)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tools, nil
}

type enumerationCase struct {
	typeId             sql.NullString
	locale             sql.NullString
	id                 sql.NullString
	title              sql.NullString
	subtitle           sql.NullString
	altText            sql.NullString
	image              sql.NullString
	snippetPluginModel sql.NullString
	synonyms           sql.NullString
}

func getEnums(typeId string) ([]enumerationCase, error) {
	var query = `select * from EnumerationCases WHERE typeId LIKE ?`

	var rows, err = toolkit.Query(query, typeId+".%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enums []enumerationCase
	for rows.Next() {
		var enum enumerationCase

		if err := rows.Scan(
			&enum.typeId,
			&enum.locale,
			&enum.id,
			&enum.title,
			&enum.subtitle,
			&enum.altText,
			&enum.image,
			&enum.snippetPluginModel,
			&enum.synonyms,
		); err != nil {
			return nil, err
		}

		enums = append(enums, enum)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return enums, nil
}
