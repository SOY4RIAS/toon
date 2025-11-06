package toon

// List markers
const (
	// ListItemMarker is the hyphen character used to denote list items.
	ListItemMarker = '-'

	// ListItemPrefix is the full prefix for list items (hyphen + space).
	ListItemPrefix = "- "
)

// Structural characters
const (
	// Comma is the comma character.
	Comma = ','

	// Colon is the colon character used to separate keys from values.
	Colon = ':'

	// Space is the space character.
	Space = ' '

	// Pipe is the pipe character.
	Pipe = '|'

	// Hash is the hash/pound character.
	Hash = '#'
)

// Brackets and braces
const (
	// OpenBracket is the opening square bracket.
	OpenBracket = '['

	// CloseBracket is the closing square bracket.
	CloseBracket = ']'

	// OpenBrace is the opening curly brace.
	OpenBrace = '{'

	// CloseBrace is the closing curly brace.
	CloseBrace = '}'
)

// Literals
const (
	// NullLiteral is the string representation of null.
	NullLiteral = "null"

	// TrueLiteral is the string representation of true.
	TrueLiteral = "true"

	// FalseLiteral is the string representation of false.
	FalseLiteral = "false"
)

// Escape characters
const (
	// Backslash is the backslash escape character.
	Backslash = '\\'

	// DoubleQuote is the double quote character.
	DoubleQuote = '"'

	// Newline is the newline character.
	Newline = '\n'

	// CarriageReturn is the carriage return character.
	CarriageReturn = '\r'

	// Tab is the tab character.
	Tab = '\t'
)

// DefaultDelimiter is the default delimiter used for encoding.
const DefaultDelimiter = DelimiterComma

// DefaultIndent is the default number of spaces per indentation level.
const DefaultIndent = 2
