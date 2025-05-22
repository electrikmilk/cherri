/*
 * Copyright (c) Cherri
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"unicode"

	"github.com/electrikmilk/args-parser"
)

var filePath string
var filename string
var basename string
var contents string
var relativePath string
var inputPath string
var outputPath string

const unsignedEnd = "_unsigned.shortcut"
const darwin = runtime.GOOS == "darwin"

func main() {
	if args.Using("help") {
		args.PrintUsage()
		os.Exit(0)
	}

	if args.Using("version") {
		printVersion()
		os.Exit(0)
	}

	if args.Using("import") && args.Value("import") != "" {
		var shortcutBytes = importShortcut()
		decompile(shortcutBytes)

		os.Exit(0)
	}

	if args.Using("action") {
		if args.Value("action") == "" {
			for identifier, definition := range actions {
				setCurrentAction(identifier, definition)
				fmt.Println(generateActionDefinition(parameterDefinition{}, true, true))
				fmt.Print("\n---\n\n")
			}
		} else {
			actionsSearch()
		}
		os.Exit(0)
	}

	if args.Using("glyph") {
		if args.Value("glyph") == "" {
			fmt.Println("Search all usable glyphs at https://glyphs.cherrilang.org.")
		} else {
			glyphsSearch()
		}
		os.Exit(0)
	}

	filePath = fileArg()
	if len(os.Args) == 1 || filePath == "" {
		printLogo()
		printVersion()
		if !darwin {
			fmt.Println(ansi("\nWarning:", yellow, bold), "Shortcuts compiled on this platform will not run on iOS 15+ or macOS 12+.")
		}
		fmt.Print("\n")
		args.PrintUsage()
		os.Exit(1)
	}

	filename = checkFile(filePath)

	handleFile()

	initParse()

	generateShortcut()

	createShortcut()
}

func handle(err error) {
	if err == nil {
		return
	}

	fmt.Print(ansi("\nProgram panic :(\n\n", red, bold))
	fmt.Print(ansi("Please report this: https://github.com/electrikmilk/cherri/issues/new\n\n", red))

	if args.Using("debug") {
		panicDebug(err)
	} else {
		panic(err)
	}
}

func printLogo() {
	fmt.Print(ansi("\n           %############                      \n           %#################                 \n           %############*######               \n            ## #############**#*              \n            ##    ############****            \n            ##%     %#%                       \n", green))
	fmt.Print(ansi("             #####", red))
	fmt.Print(ansi("    %##    ####             \n         ###****######  #############         \n        ##**######################***#        \n       ############################*+*#       \n      #############################***#       \n       #############################*##       \n       ################################       \n        ##############  ##############        \n           #########      #########           \n\n", red))
}

func printVersion() {
	var color outputType
	if strings.Contains(version, "beta") {
		color = yellow
	} else {
		color = green
	}
	fmt.Println("Cherri Compiler", ansi(version, color))
}

func fileArg() string {
	for _, arg := range os.Args {
		if !strings.Contains(arg, ".cherri") || startsWith("-", arg) {
			continue
		}

		return arg
	}
	return ""
}

// handleFile splits the file argument into parts.
func handleFile() {
	relativePath = strings.Replace(filePath, filename, "", 1)
	var nameParts = strings.Split(filename, ".")
	basename = nameParts[0]
	workflowName = basename

	outputPath = getOutputPath(relativePath + workflowName + ".shortcut")

	var fileBytes, readErr = os.ReadFile(filePath)
	handle(readErr)
	contents = string(fileBytes)
}

func getOutputPath(defaultPath string) string {
	if args.Using("output") {
		return args.Value("output")
	}

	return defaultPath
}

// checkFile checks if the file exists and is a .cherri file.
func checkFile(filePath string) (filename string) {
	var file, statErr = os.Stat(filePath)
	if os.IsNotExist(statErr) {
		exit(fmt.Sprintf("File '%s' does not exist!", filePath))
	}
	var nameParts = strings.Split(file.Name(), ".")
	var ext = end(nameParts)
	if ext != "cherri" {
		exit(fmt.Sprintf("File '%s' is not a .cherri file!", filePath))
	}
	return file.Name()
}

func end(slice []string) string {
	return slice[len(slice)-1]
}

func capitalize(s string) string {
	var char = s[0]
	var after, _ = strings.CutPrefix(s, string(char))
	return fmt.Sprintf("%c%s", unicode.ToUpper(rune(char)), after)
}

// startsWith determines if the beginning characters of `substr` match `s`.
func startsWith(s string, substr string) bool {
	var subStringChars = []rune(substr)
	var stringChars = []rune(s)
	var stringSize = len(s)
	var start string
	for i, char := range subStringChars {
		if stringSize < i+1 {
			break
		}
		if char != stringChars[i] {
			if len(start) > 0 {
				break
			}
			return false
		}
		start = fmt.Sprintf("%s%c", start, char)
	}

	return start == s
}

func lineReport(label string) {
	fmt.Printf("--- %s ---\n", label)
	if idx != 0 {
		fmt.Println("Previous Character:")
		var prevChar = prev(1)
		if prevChar != '\n' {
			printChar(prevChar, lineIdx, lineCharIdx-1)
		} else {
			printChar(prevChar, lineIdx-1, len(lines[lineIdx-1]))
		}
	}

	fmt.Println("\nCurrent Character:")
	printChar(char, lineIdx, lineCharIdx)
	fmt.Print("\n")

	if len(contents) > idx+1 {
		fmt.Println("Next Character:")
		var nextChar = next(1)
		if char != '\n' {
			printChar(nextChar, lineIdx, lineCharIdx+1)
		} else {
			printChar(nextChar, lineIdx+1, 0)
		}
		fmt.Print("\n")
	}

	if len(lines) > lineIdx {
		fmt.Printf("Current Line:\n%s\n", lines[lineIdx])
	}
}

func panicDebug(err error) {
	fmt.Println(ansi("###################\n#   DEBUG PANIC   #\n###################\n", bold, red))
	printParsingDebug()
	printShortcutGenDebug()
	printCustomActionsDebug()
	printIncludesDebug()
	fmt.Println(ansi("#############################################################\n", bold, red))

	if err != nil {
		panic(err)
	}

	panic("debug")
}

// Converts a map[string]interface{} to a matching struct data type.
func mapToStruct(data any, structure any) {
	var jsonStr, marshalErr = json.Marshal(data)
	handle(marshalErr)

	var jsonErr = json.Unmarshal(jsonStr, &structure)
	if jsonErr != nil {
		fmt.Println("Tried to map to struct, but it was not a struct!", data)
		handle(jsonErr)
	}
}

// waitFor takes functions and uses a WaitGroup to wait for them all to finish.
func waitFor(functions ...func()) {
	var wg sync.WaitGroup
	wg.Add(len(functions))
	for _, function := range functions {
		go func() {
			defer wg.Done()
			function()
		}()
	}
	wg.Wait()
}
