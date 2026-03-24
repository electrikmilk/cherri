package main

import (
	"regexp"
	"testing"
)

func TestCapitalizeEmptyString(t *testing.T) {
	if got := capitalize(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestSanitizeIdentifierWhitespaceOnly(t *testing.T) {
	specialCharsRegex = regexp.MustCompile("[^a-zA-Z0-9_]+")
	identifier := " "

	sanitizeIdentifier(&identifier)

	if identifier != "" {
		t.Fatalf("expected sanitized identifier to be empty, got %q", identifier)
	}
}
