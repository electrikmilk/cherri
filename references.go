package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/electrikmilk/args-parser"
	"howett.net/plist"
)

// extractedMediaReference holds the state for a single extracted reference.
type extractedMediaReference struct {
	action     string
	identifier string
	index      int
	value      any
}

// extractedReferences holds references extracted from a workflow.
var extractedReferences []extractedMediaReference

// references holds parsed references.
var references = make(map[string]map[string]any)

// extractReferences extracts input references from a workflow with unique identifiers and outputs hash identifier media reference syntax.
func extractReferences(b []byte) {
	var _, marshalIndexedErr = plist.Unmarshal(b, &shortcut)
	handle(marshalIndexedErr)

	for i, action := range shortcut.WFWorkflowActions {
		if len(action.WFWorkflowActionParameters) == 0 {
			continue
		}
		extractParameterReferences(i, action.WFWorkflowActionIdentifier, action.WFWorkflowActionParameters)
	}

	if args.Using("output") {
		writeRefsToFile()
	} else {
		for _, ref := range extractedReferences {
			printMediaReference(&ref)
		}
	}
}

// writeRefsToFile writes extracted references to a file determined by the --output flag.
func writeRefsToFile() {
	var refsFile strings.Builder
	for _, ref := range extractedReferences {
		refsFile.WriteString(makeRefCode(&ref))
	}

	var writeRefsErr = os.WriteFile(args.Value("output"), []byte(refsFile.String()), 0644)
	handle(writeRefsErr)
}

// extractParameterReferences extracts references from a workflow action parameter.
func extractParameterReferences(index int, identifier string, params map[string]interface{}) {
	for _, value := range params {
		var ref = extractedMediaReference{
			action: identifier,
			index:  index + 1,
		}
		if value == nil || reflect.TypeOf(value).Kind() != reflect.Map {
			continue
		}
		var valueMap = value.(map[string]interface{})
		if valueMap["fileLocation"] != nil {
			ref.value = extractFileReference(&ref, valueMap)
		}
		if valueMap["persistentIdentifier"] != nil {
			ref.value = extractMediaReference(&ref, valueMap)
		}
		if valueMap["CustomOutputName"] != nil {
			ref.identifier = valueMap["CustomOutputName"].(string)
		}

		sanitizeIdentifier(&ref.identifier)

		if skipDuplicateReference(ref) {
			continue
		}

		extractedReferences = append(extractedReferences, ref)
	}
	return
}

func skipDuplicateReference(ref extractedMediaReference) bool {
	for _, existingRef := range extractedReferences {
		existingIdentifier := existingRef.identifier
		if existingIdentifier == ref.identifier {
			return true
		}
	}
	return false
}

type WFFile struct {
	Filename     string       `json:"filename,omitempty" plist:"filename,omitempty"`
	DisplayName  string       `json:"displayName,omitempty" plist:"displayName,omitempty"`
	FileLocation FileLocation `json:"fileLocation,omitempty" plist:"fileLocation,omitempty"`
}

type FileLocation struct {
	WFFileLocationType   string
	CrossDeviceItemID    string `json:"crossDeviceItemID,omitempty" plist:"crossDeviceItemID,omitempty"`
	FileProviderDomainID string `json:"fileProviderDomainID" plist:"fileProviderDomainID"`
	RelativeSubpath      string `json:"relativeSubpath" plist:"relativeSubpath"`
}

// extractFileReference extracts a file reference from a workflow action parameter.
func extractFileReference(ref *extractedMediaReference, values map[string]interface{}) (file WFFile) {
	mapToStruct(values, &file)
	ref.identifier = file.DisplayName
	return
}

type WFMediaItem struct {
	PersistentIdentifier uint64 `json:"persistentIdentifier" plist:"persistentIdentifier"`
	ItemName             string `json:"itemName" plist:"itemName"`
	Type                 string `json:"type" plist:"type"`
}

// extractMediaReference extracts a media reference from a workflow action parameter.
func extractMediaReference(ref *extractedMediaReference, values map[string]interface{}) (media WFMediaItem) {
	mapToStruct(values, &media)
	ref.identifier = media.ItemName
	return
}

// makeReferenceHash returns a unique hash for a reference.
func makeReferenceHash(ref *extractedMediaReference) string {
	var jsonBytes, marshalErr = json.Marshal(ref.value)
	handle(marshalErr)

	return base64.StdEncoding.EncodeToString(jsonBytes)
}

// decodeReferenceHash decodes a reference hash into a map to later insert into a Shortcut action parameter.
func decodeReferenceHash(hash string) (ref map[string]any, err error) {
	ref = make(map[string]any)
	var decodedBytes, decodeErr = base64.StdEncoding.DecodeString(hash)
	if decodeErr != nil {
		return nil, fmt.Errorf("could not decode hashed JSON: %s", decodeErr)
	}

	fmt.Println(string(decodedBytes))

	var unmarshalErr = json.Unmarshal(decodedBytes, &ref)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("could not unmarshal decoded JSON: %s", unmarshalErr)
	}

	return ref, nil
}

// makeRefCode returns the reference in Cherri reference syntax.
func makeRefCode(ref *extractedMediaReference) string {
	return fmt.Sprintf("#ref %s %s\n", ref.identifier, makeReferenceHash(ref))
}

// printMediaReference prints a media reference to stdout.
func printMediaReference(ref *extractedMediaReference) {
	if args.Using("debug") {
		fmt.Println(ansi(fmt.Sprintf("%s (Action %d, %s)\n", ref.identifier, ref.index, ref.action), bold))
		fmt.Println(ansi(fmt.Sprintf("Value: %v", ref.value), green))
	}

	fmt.Println(ansi(makeRefCode(ref), magenta))

	if args.Using("debug") {
		fmt.Println(ansi("---\n", dim))
	}
}
