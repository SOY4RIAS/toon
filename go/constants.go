package toon

// Character constants
const (
	Space        = ' '
	Tab          = '\t'
	Backslash    = '\\'
	DoubleQuote  = '"'
	Colon        = ':'
	Comma        = ','
	Pipe         = '|'
	Hash         = '#'
	Hyphen       = '-'
	OpenBracket  = '['
	CloseBracket = ']'
	OpenBrace    = '{'
	CloseBrace   = '}'
)

// List item markers
const (
	ListItemMarker = '-'
	ListItemPrefix = "- "
)

// Literal constants
const (
	TrueLiteral  = "true"
	FalseLiteral = "false"
	NullLiteral  = "null"
)

// Delimiter type represents array value delimiters
type Delimiter rune

// Available delimiters
const (
	DelimiterComma Delimiter = ','
	DelimiterTab   Delimiter = '\t'
	DelimiterPipe  Delimiter = '|'
)

// DefaultDelimiter is the default delimiter for arrays
const DefaultDelimiter = DelimiterComma

// IsValidDelimiter checks if a delimiter is valid
func IsValidDelimiter(d Delimiter) bool {
	return d == DelimiterComma || d == DelimiterTab || d == DelimiterPipe
}
