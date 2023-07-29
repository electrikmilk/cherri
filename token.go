/*
 * Copyright (c) Brandon Jordan
 */

package main

type tokenType string

type token struct {
	typeof    tokenType
	ident     string
	valueType tokenType
	value     any
}

var tokens []token

var tokenChars map[tokenType][]string

/* Keywords */

const (
	Var            tokenType = "variable"
	Constant       tokenType = "const"
	If             tokenType = "if"
	Else           tokenType = "else"
	EndClosure     tokenType = "endif"
	Repeat         tokenType = "repeat "
	RepeatWithEach tokenType = "for "
	Menu           tokenType = "menu"
	Item           tokenType = "item"
	Definition     tokenType = "#define"
	Import         tokenType = "#import"
	Question       tokenType = "#question"
	CustomAction   tokenType = "action"
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
	Mac     tokenType = "mac"
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
	At             tokenType = "@"
	Set            tokenType = "="
	AddTo          tokenType = "+="
	Is             tokenType = "=="         // is
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
