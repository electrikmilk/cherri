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
		Description: "Create plist file, print debug and stack traces.",
	})
	args.Register(args.Argument{
		Name:        "output",
		Short:       "o",
		Description: "Set custom output file path.",
	})
	args.Register(args.Argument{
		Name:        "import",
		Short:       "i",
		Description: "Import compiled Shortcut (ignored if unsigned).",
	})
	args.Register(args.Argument{
		Name:        "comments",
		Short:       "c",
		Description: "Create comment actions for text comments (e.g. //, /**/)",
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
		Description:  "Print action definition. Empty prints all definitions.",
		DefaultValue: "",
		ExpectsValue: true,
	})
	args.Register(args.Argument{
		Name:         "glyph",
		Description:  "Search glyphs in the compiler.",
		DefaultValue: "",
		ExpectsValue: true,
	})
}
