package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"
	"unicode"

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

func importActions(identifier string) {
	if slices.Contains(imported, identifier) {
		parserError(fmt.Sprintf("Actions and types from '%s' have already been imported.", identifier))
	}

	connectToolkitDB()

	if !appIdentifierRegex.MatchString(identifier) {
		matchApplication(&identifier)
	}

	var importedActions, actionsErr = getActions(identifier)
	handle(actionsErr)

	if args.Using("debug") {
		fmt.Println("### IMPORTING ACTIONS ###")
		fmt.Println("Importing actions for", identifier)
	}

	defineImportedActions(identifier, importedActions)

	imported = append(imported, identifier)
	if args.Using("debug") {
		fmt.Println("Imported actions from", identifier)
	}
}

func defineImportedActions(identifier string, importedActions []actionTool) {
	var namespace = end(strings.Split(identifier, "."))
	for _, action := range importedActions {
		var intent = end(strings.Split(action.id.String, "."))
		var trimIntent = strings.TrimSuffix(intent, "Intent")
		var name = fmt.Sprintf("%s_%s", namespace, trimIntent)
		if args.Using("debug") {
			fmt.Println("Action name: ", name)
		}
		var actionIdentifier = camelCase(name)

		var outputType, outputTypeErr = getActionOutputType(action.rowId.String)
		handle(outputTypeErr)

		var actionLocalization, localizeErr = getActionLocalization(action.rowId.String)
		handle(localizeErr)
		var doc = selfDoc{
			title:       actionLocalization.name.String,
			description: actionLocalization.descriptionSummary.String,
		}

		var paramDefs = importParamDefinitions(action.rowId.String, action.id.String)

		actions[actionIdentifier] = &actionDefinition{
			overrideIdentifier: action.id.String,
			parameters:         paramDefs,
			outputType:         outputType,
			doc:                doc,
		}

		if args.Using("debug") {
			fmt.Println("Imported action:", name+"()")
			fmt.Println("Params:", paramDefs)
			fmt.Print("\n")
		}
	}
}

func importParamDefinitions(toolId string, identifier string) (definitions []parameterDefinition) {
	var params, paramsErr = getActionParams(toolId)
	handle(paramsErr)

	var sortedParams = make([]toolParam, len(params))
	copy(sortedParams, params)
	slices.SortFunc(sortedParams, func(a, b toolParam) int {
		return int(a.sortOrder.Int64 - b.sortOrder.Int64)
	})

	for _, param := range sortedParams {
		var def = parameterDefinition{
			key: param.key.String,
		}

		if args.Using("debug") {
			fmt.Println("Param:", def.key)
		}

		var paramName, paramNameErr = getActionParamName(toolId, def.key)
		handle(paramNameErr)
		def.name = paramName

		var paramTokenType, tokenTypeErr = getActionParamType(toolId, def.key)
		handle(tokenTypeErr)
		def.validType = paramTokenType

		var enums, enumErr = getParamEnums(identifier, def.key)
		handle(enumErr)

		defineParamEnums(param, enums, &def)

		definitions = append(definitions, def)
	}

	return
}

func defineParamEnums(param toolParam, enums []enumerationCase, definition *parameterDefinition) {
	var paramEnumerations []string
	for _, enum := range enums {
		paramEnumerations = append(paramEnumerations, enum.title.String)
	}

	if args.Using("debug") {
		fmt.Println("Param Enum:", paramEnumerations)
	}

	if len(paramEnumerations) == 0 {
		return
	}

	var enumName = fmt.Sprintf("%ss", param.key.String)
	definition.enum = enumName
	definition.validType = String

	if _, found := enumerations[enumName]; !found {
		enumerations[enumName] = paramEnumerations

		if args.Using("debug") {
			fmt.Println("Defined enum", enumName)
		}
	}
}

func camelCase(s string) (c string) {
	for i, r := range s {
		if unicode.IsSpace(r) || r == '-' {
			continue
		}
		if i != 0 && unicode.IsUpper(r) {
			c += strings.ToUpper(string(r))
		} else {
			c += strings.ToLower(string(r))
		}
	}
	return
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

type toolParam struct {
	key       sql.NullString
	sortOrder sql.NullInt64
}

func getActionParams(toolId string) ([]toolParam, error) {
	var query = `select key, sortOrder from Parameters WHERE toolId = ?`

	var rows, err = toolkit.Query(query, toolId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var params []toolParam
	for rows.Next() {
		var param toolParam

		if err := rows.Scan(
			&param.key,
			&param.sortOrder,
		); err != nil {
			return nil, err
		}

		params = append(params, param)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return params, nil
}

func getActionParamName(toolId string, key string) (string, error) {
	var query = `select name from ParameterLocalizations WHERE toolId = ? AND key = ? AND locale = 'en' LIMIT 1`
	var row = toolkit.QueryRow(query, toolId, key)
	handle(row.Err())

	var paramName string
	var scanErr = row.Scan(&paramName)
	handle(scanErr)

	paramName = camelCase(paramName)

	return paramName, nil
}

func getActionParamType(toolId string, key string) (tokenType, error) {
	var query = `select typeId from ToolParameterTypes WHERE toolId = ? AND key = ? LIMIT 1`
	var row = toolkit.QueryRow(query, toolId, key)
	handle(row.Err())

	var paramType string
	var scanErr = row.Scan(&paramType)
	handle(scanErr)

	var paramTokenType tokenType
	switch paramType {
	case "string":
		paramTokenType = String
	case "int":
		fallthrough
	case "number":
		paramTokenType = Integer
	case "decimal":
		paramTokenType = Float
	case "bool":
		paramTokenType = Bool
	case "dictionary":
		paramTokenType = Dict
	default:
		paramTokenType = Variable
	}

	return paramTokenType, nil
}

func getActionOutputType(toolId string) (tokenType, error) {
	var query = `select typeIdentifier from ToolOutputTypes WHERE toolId = ? LIMIT 1`
	var row = toolkit.QueryRow(query, toolId)
	handle(row.Err())

	var outputType string
	var scanErr = row.Scan(&outputType)
	handle(scanErr)

	var outputTokenType tokenType
	switch outputType {
	case "string":
		outputTokenType = String
	case "int":
		outputTokenType = Integer
	case "decimal":
		outputTokenType = Float
	case "bool":
		outputTokenType = Bool
	}

	return outputTokenType, nil
}

type actionLocalization struct {
	name               sql.NullString
	descriptionSummary sql.NullString
}

func getActionLocalization(toolId string) (actionLocalization, error) {
	var query = `select name, descriptionSummary from ToolLocalizations WHERE toolId = ? and locale = 'en' LIMIT 1`
	var row = toolkit.QueryRow(query, toolId)
	handle(row.Err())

	var localization actionLocalization
	var scanErr = row.Scan(&localization.name, &localization.descriptionSummary)
	if scanErr != nil {
		return actionLocalization{}, scanErr
	}

	return localization, nil
}

type enumerationCase struct {
	title sql.NullString
}

func getParamEnums(identifier string, key string) ([]enumerationCase, error) {
	var query = `select title from EnumerationCases WHERE typeId = ? AND locale = 'en'`

	var enumIdentifier = fmt.Sprintf("com.apple.shortcuts.%s.%s", identifier, key)

	var rows, err = toolkit.Query(query, enumIdentifier)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enums []enumerationCase
	for rows.Next() {
		var enum enumerationCase

		if err := rows.Scan(
			&enum.title,
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
