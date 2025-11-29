/*
 * Copyright (c) Cherri
 */

package main

import (
	"bufio"
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
var originalContents string
var relativePath string
var inputPath string
var outputPath string

var internalDirectoryPath = os.ExpandEnv("$HOME/.cherri")

const unsignedEnd = "_unsigned.shortcut"
const darwin = runtime.GOOS == "darwin"

func main() {
	filePath = fileArg()
	if filePath != "" {
		filename = checkFile(filePath)

		handleFile()

		initParse()

		generateShortcut()

		createShortcut()
		return
	}

	if args.Using("help") {
		args.PrintUsage()
		os.Exit(0)
	}

	if args.Using("version") {
		printVersion()
		os.Exit(0)
	}

	if args.Using("docs") {
		generateDocs()
		os.Exit(0)
	}

	if args.Using("init") {
		initPackage()
		os.Exit(0)
	}
	if args.Using("add-uri") {
		addUri()
		os.Exit(0)
	}
	if args.Using("install") {
		addPackage()
		os.Exit(0)
	}
	if args.Using("remove") {
		removePackage()
		os.Exit(0)
	}
	if args.Using("package") {
		listPackage()
		os.Exit(0)
	}
	if args.Using("packages") {
		listPackages()
		os.Exit(0)
	}
	if args.Using("tidy") {
		tidyPackage()
		os.Exit(0)
	}

	if args.Using("import") && args.Value("import") != "" {
		var shortcutBytes = importShortcut()
		decompile(shortcutBytes)

		os.Exit(0)
	}

	if args.Using("action") {
		markBuiltins()
		defineRawAction()
		defineToggleSetActions()
		loadStandardActions()
		handleActionSearch()
		os.Exit(0)
	}

	if args.Using("glyph") {
		handleGlyphSearch()
		os.Exit(0)
	}

	printLogo()
	printVersion()
	fmt.Print("\n")
	args.PrintUsage()
	os.Exit(1)
}

func yesNo() bool {
	var input string
	for {
		scan("(y/n) ", &input)
		input = strings.ToLower(input)
		if input == "y" || input == "n" {
			return input == "y"
		}
	}
}

func camelCase(s string) (c string) {
	var lastChar rune
	for i, r := range s {
		switch {
		case unicode.IsLetter(r):
			if i == 0 {
				c += strings.ToLower(string(r))
			} else if unicode.IsLetter(lastChar) {
				c += strings.ToLower(string(r))
			} else {
				c += strings.ToUpper(string(r))
			}
		case r == '_' || unicode.IsDigit(r):
			c += string(r)
		}
		lastChar = r
	}
	return
}

func scan(prompt string, store *string) {
	fmt.Print(prompt)

	var scanner = bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		*store = scanner.Text()
	}

	var scanErr = scanner.Err()
	handle(scanErr)
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

func createInternalDir() {
	if _, statErr := os.Stat(internalDirectoryPath); os.IsNotExist(statErr) {
		var intDirErr = os.Mkdir(internalDirectoryPath, 0777)
		handle(intDirErr)
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
	if len(os.Args) < 2 {
		return ""
	}
	var fileName = os.Args[1]
	if !startsWith("-", fileName) && strings.Contains(fileName, ".cherri") {
		return fileName
	}
	return ""
}

// handleFile splits the file argument into parts.
func handleFile() {
	relativePath = strings.Replace(filePath, filename, "", 1)
	var nameParts = strings.Split(filename, ".")
	basename = nameParts[0]
	workflowName = basename

	var fileBytes, readErr = os.ReadFile(filePath)
	handle(readErr)
	contents = string(fileBytes)
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
	printFunctionsDebug()
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
