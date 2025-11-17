package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/electrikmilk/args-parser"
	_ "github.com/glebarez/go-sqlite"
)

var shortcutsDBPath = os.ExpandEnv("$HOME/Library/Shortcuts/ToolKit/Tools-active")

func init() {
	if args.Using("database") {
		shortcutsDBPath = args.Value("database")
	}

	var db, dbErr = sql.Open("sqlite", shortcutsDBPath)
	handle(dbErr)

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

	var row = db.QueryRow("select * from EnumerationCases WHERE typeId = ?", "com.sindresorhus.Color-Picker.TransportType")

	var enum enumerationCase
	var scanErr = row.Scan(&enum.typeId, &enum.locale, &enum.id, &enum.title, &enum.subtitle, &enum.altText, &enum.image, &enum.snippetPluginModel, &enum.synonyms)
	handle(scanErr)

	fmt.Println(enum)
}
