/*
 * Copyright (c) Brandon Jordan
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
		Name:         "action",
		Short:        "a",
		Description:  "Print an action's definition. Leave empty to print all action definitions.",
		DefaultValue: "",
		ExpectsValue: true,
	})
	args.Register(args.Argument{
		Name:         "share",
		Short:        "s",
		Description:  "Set the Shortcuts signing mode, passed to the `shortcuts` binary.",
		Values:       []string{"anyone", "contacts"},
		DefaultValue: "contacts",
		ExpectsValue: true,
	})
	args.Register(args.Argument{
		Name:        "debug",
		Short:       "d",
		Description: "Save the generated plist and print debug messages and stack traces.",
	})
	args.Register(args.Argument{
		Name:        "output",
		Short:       "o",
		Description: "Optional output file path. (e.g. path/to/file.shortcut).",
	})
	args.Register(args.Argument{
		Name:        "import",
		Short:       "i",
		Description: "Opens compiled Shortcut after compilation (ignored if unsigned).",
	})
	args.Register(args.Argument{
		Name:        "comments",
		Short:       "c",
		Description: "Include comments in the compiled Shortcut.",
	})
	args.Register(args.Argument{
		Name:        "no-ansi",
		Description: "Don't output ANSI escape sequences that format and color the output.",
	})
}
