//go:build darwin

/*
 * Copyright (c) Cherri
 */

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/electrikmilk/args-parser"
)

// TestMacOSSigning verifies that native macOS signing works correctly.
// This test only runs on macOS (darwin) and will fail if the `shortcuts sign` command
// fails, indicating that the code is falling back to an external signing service.
func TestMacOSSigning(t *testing.T) {
	args.Args["no-ansi"] = ""
	// Use "anyone" signing mode which is confirmed to work on macOS 14.4+/15.x
	// The "people-who-know-me" mode may have iCloud/network dependencies that fail
	args.Args["share"] = "anyone"

	// Reset signFailed before testing
	signFailed = false

	var files, err = os.ReadDir("tests")
	if err != nil {
		t.Fatalf("Failed to read tests directory: %v", err)
	}

	loadStandardActions()

	// Find the first valid test file to compile
	var testFile string
	for _, file := range files {
		if file.Name() == "decomp_expected.cherri" || file.Name() == "decomp_me.cherri" {
			continue
		}
		if len(file.Name()) > 7 && file.Name()[len(file.Name())-7:] == ".cherri" {
			testFile = fmt.Sprintf("tests/%s", file.Name())
			break
		}
	}

	if testFile == "" {
		t.Skip("No test files found")
	}

	currentTest = testFile
	os.Args[1] = testFile
	fmt.Println("Testing macOS native signing with:", testFile)

	// Compile the shortcut (which includes signing)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Compilation panicked: %v", r)
		}
	}()

	compile()

	// Check if native macOS signing failed
	if signFailed {
		t.Error("macOS native signing failed - code fell back to external signing service. " +
			"The `shortcuts sign` command is not working correctly on this system.")
	} else {
		fmt.Println("✅ macOS native signing succeeded")
	}

	resetParser()
}
