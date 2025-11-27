/*
 * Copyright (c) Cherri
 */

package main

import "github.com/electrikmilk/args-parser"

func init() {
	args.CustomUsage = "[FILE]"
	args.Register(args.Argument{
		Name:        "version",
		Short:       "v",
		Description: "Print version information.",
	})
	args.Register(args.Argument{
		Name:        "help",
		Short:       "h",
		Description: "Print this usage information.",
	})
	args.Register(args.Argument{
		Name:         "share",
		Short:        "s",
		Description:  "Set the signing mode.",
		Values:       []string{"anyone", "contacts"},
		DefaultValue: "contacts",
		ExpectsValue: true,
	})
	args.Register(args.Argument{
		Name:        "debug",
		Short:       "d",
		Description: "Create plist file, preprocessed file, print debug and stack traces.",
	})
	args.Register(args.Argument{
		Name:        "derive-uuids",
		Description: "Output deterministic UUIDs.",
	})
	args.Register(args.Argument{
		Name:         "output",
		Short:        "o",
		Description:  "Set custom output file path.",
		DefaultValue: "",
		ExpectsValue: true,
	})
	args.Register(args.Argument{
		Name:        "comments",
		Short:       "c",
		Description: "Include comment actions in compiled Shortcut or import.",
	})
	args.Register(args.Argument{
		Name:        "open",
		Description: "Open compiled Shortcut (ignored if unsigned).",
	})
	args.Register(args.Argument{
		Name:         "import",
		Description:  "[BETA] Import Shortcut from an iCloud link or file path and convert to Cherri.",
		ExpectsValue: true,
	})
	args.Register(args.Argument{
		Name:         "init",
		Description:  "Create a Cherri package. Pattern: @{github_username}/{package-name}",
		ExpectsValue: true,
	})
	args.Register(args.Argument{
		Name:         "install",
		Description:  "Install a package.",
		ExpectsValue: true,
	})
	args.Register(args.Argument{
		Name:         "remove",
		Description:  "Remove a package.",
		ExpectsValue: true,
	})
	args.Register(args.Argument{
		Name:        "package",
		Description: "List package info.",
	})
	args.Register(args.Argument{
		Name:        "packages",
		Description: "List installed packages.",
	})
	args.Register(args.Argument{
		Name:        "tidy",
		Description: "Re-install all packages.",
	})
	args.Register(args.Argument{
		Name:         "toolkit",
		Description:  "Path to Shortcuts ToolKit SQLite database.",
		ExpectsValue: true,
	})
	args.Register(args.Argument{
		Name:         "toolkit-locale",
		Description:  "Set custom locale to get action data for.",
		ExpectsValue: true,
		DefaultValue: "en",
	})
	args.Register(args.Argument{
		Name:        "no-toolkit",
		Description: "Do not use the Shortcuts toolkit DB to decompile non-standard actions.",
	})
	args.Register(args.Argument{
		Name:         "signing-server",
		ExpectsValue: true,
		Description:  "Sign the compiled Shortcut using a remote signing service that runs https://github.com/scaxyz/shortcut-signing-server.",
	})
	args.Register(args.Argument{
		Name:        "hubsign",
		Description: "Sign the compiled Shortcut using RoutineHub's remote signing service.",
	})
	args.Register(args.Argument{
		Name:        "no-ansi",
		Description: "Don't output ANSI escape sequences that format the output.",
	})
	args.Register(args.Argument{
		Name:        "skip-sign",
		Description: "Do not sign the compiled Shortcut.",
	})
	args.Register(args.Argument{
		Name:         "action",
		Description:  "Search for available actions. Empty prints all definitions.",
		DefaultValue: "",
		ExpectsValue: true,
	})
	args.Register(args.Argument{
		Name:         "glyph",
		Description:  "Search for available glyphs.",
		DefaultValue: "",
		ExpectsValue: true,
	})

	for _, actionCat := range actionIncludes {
		actionCategories = append(actionCategories, actionCat)
	}
	args.Register(args.Argument{
		Name:         "docs",
		Description:  "Generate action documentation, optionally by category.",
		DefaultValue: "",
		ExpectsValue: true,
		Values:       actionCategories,
	})
	args.Register(args.Argument{
		Name:         "subcat",
		Description:  "Filter action documentation category by subcategory.",
		DefaultValue: "",
		ExpectsValue: true,
	})
}
