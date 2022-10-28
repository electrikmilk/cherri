/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

type tokenType string

type token struct {
	col       int
	typeof    tokenType
	ident     string
	valueType tokenType
	value     any
}

var tokens []token

/* Keywords */

const (
	Var            tokenType = "@"
	If             tokenType = "if"
	Else           tokenType = "else"
	EndIf          tokenType = "endif"
	Repeat         tokenType = "repeat"
	RepeatWithEach tokenType = "foreach"
	Menu           tokenType = "menu"
	Case           tokenType = "case"
	Definition     tokenType = "#define"
)

/* Definitions */

const (
	Name    tokenType = "name"
	Color   tokenType = "color"
	Glyph   tokenType = "glyph"
	Inputs  tokenType = "inputs"
	Outputs tokenType = "outputs"
	From    tokenType = "from"
	Version tokenType = "version"
	NoInput tokenType = "noinput"
)

/* No Inputs */

const (
	StopWith     tokenType = "stopwith"
	AskFor       tokenType = "askfor"
	GetClipboard tokenType = "getclipboard"
)

/* Types */

const (
	String      tokenType = "\""
	Integer     tokenType = "0123456789.-"
	Dict        tokenType = "{"
	Arr         tokenType = "["
	Bool        tokenType = "boolean"
	Date        tokenType = "date"
	True        tokenType = "true"
	False       tokenType = "false"
	Comment     tokenType = "comment"
	Expression  tokenType = "expression"
	Variable    tokenType = "variable"
	Action      tokenType = "action"
	Conditional tokenType = "conditional"
)

func typeName(typeOf tokenType) string {
	switch typeOf {
	case String:
		return "string"
	case Integer:
		return "integer"
	case Arr:
		return "array"
	case Dict:
		return "dictionary"
	default:
		return string(typeOf)
	}
}

/* Operators */

const (
	Set            tokenType = "="
	AddTo          tokenType = "+="
	Is             tokenType = "=="         // is <- use these?
	Not            tokenType = "!="         // is not
	Any            tokenType = "value"      // has any value
	Empty          tokenType = "!value"     // does not have any value
	Contains       tokenType = "contains"   // contains
	DoesNotContain tokenType = "!contains"  // does not contain
	BeginsWith     tokenType = "beginsWith" // begins with
	EndsWith       tokenType = "endsWith"   // ends with
	GreaterThan    tokenType = ">"
	GreaterOrEqual tokenType = ">="
	LessThan       tokenType = "<"
	LessOrEqual    tokenType = "<="
	Exclamation    tokenType = "!"
	Between        tokenType = "<>"
	ForwardSlash   tokenType = "/"
	Plus           tokenType = "+"
	Minus          tokenType = "-"
	Multiply       tokenType = "*"
	Divide         tokenType = "/"
	Modulus        tokenType = "%"
	LeftBrace      tokenType = "{"
	RightBrace     tokenType = "}"
)

/* Variables */

var variables map[string]variableValue

type variableValue struct {
	variableType string
	valueType    tokenType
	value        any
	getAs        string
	coerce       string
}

/* Menus */

var menus map[string][]variableValue

/* Conditionals */

type condition struct {
	variableOneType    tokenType
	variableOneValue   any
	condition          string
	variableTwoType    tokenType
	variableTwoValue   any
	variableThreeType  tokenType
	variableThreeValue any
}

var conditions map[tokenType]string

var groupingUUID string

func makeConditions() {
	conditions = make(map[tokenType]string)
	conditions[Is] = "4"
	conditions[Not] = "5"
	conditions[Any] = "100"
	conditions[Empty] = "101"
	conditions[Contains] = "99"
	conditions[DoesNotContain] = "999"
	conditions[BeginsWith] = "8"
	conditions[EndsWith] = "9"
	conditions[GreaterThan] = "2"
	conditions[GreaterOrEqual] = "3"
	conditions[LessThan] = "0"
	conditions[LessOrEqual] = "1"
	conditions[Between] = "1003"
}
