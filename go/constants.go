package toon

// List markers
const (
	ListItemMarker = '-'
	ListItemPrefix = "- "
)

// Structural characters
const (
	Comma = ','
	Colon = ':'
	Space = ' '
	Pipe  = '|'
	Hash  = '#'
)

// Brackets and braces
const (
	OpenBracket  = '['
	CloseBracket = ']'
	OpenBrace    = '{'
	CloseBrace   = '}'
)

// Literals
const (
	NullLiteral  = "null"
	TrueLiteral  = "true"
	FalseLiteral = "false"
)

// Escape characters
const (
	Backslash       = '\\'
	DoubleQuote     = '"'
	Newline         = '\n'
	CarriageReturn  = '\r'
	Tab             = '\t'
)

// Delimiter is a type for array/table delimiters
type Delimiter rune

// Delimiter constants
const (
	DelimiterComma Delimiter = Comma
	DelimiterTab   Delimiter = Tab
	DelimiterPipe  Delimiter = Pipe
)

// DefaultDelimiter is the default delimiter (comma)
const DefaultDelimiter = DelimiterComma

// DelimiterToString converts a Delimiter to its string representation
func DelimiterToString(d Delimiter) string {
	switch d {
	case DelimiterTab:
		return "\\t"
	default:
		return string(d)
	}
}
