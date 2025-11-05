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
	Backslash      = '\\'
	DoubleQuote    = '"'
	Newline        = '\n'
	CarriageReturn = '\r'
	Tab            = '\t'
)

// Delimiter represents a valid delimiter type
type Delimiter rune

// Delimiter constants
const (
	DelimiterComma Delimiter = Comma
	DelimiterTab   Delimiter = Tab
	DelimiterPipe  Delimiter = Pipe
)

// DefaultDelimiter is the default delimiter (comma)
const DefaultDelimiter = DelimiterComma
