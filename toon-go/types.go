package toon

// Constants for list markers
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

// Delimiter type represents a delimiter character used in TOON format
type Delimiter rune

// Delimiter constants
const (
	DelimiterComma Delimiter = ','
	DelimiterTab   Delimiter = '\t'
	DelimiterPipe  Delimiter = '|'
)

// DefaultDelimiter is the default delimiter (comma)
const DefaultDelimiter = DelimiterComma

// DecodeOptions represents options for decoding TOON format
type DecodeOptions struct {
	// Indent is the number of spaces per indentation level
	// Default: 2
	Indent *int

	// Strict enforces strict validation of array lengths and tabular row counts
	// Default: true
	Strict *bool
}

// ResolvedDecodeOptions represents resolved decode options with defaults applied
type ResolvedDecodeOptions struct {
	Indent int
	Strict bool
}

// ArrayHeaderInfo contains information parsed from an array header line
type ArrayHeaderInfo struct {
	Key             *string
	Length          int
	Delimiter       Delimiter
	Fields          []string
	HasLengthMarker bool
}

// ParsedLine represents a single parsed line from the input
type ParsedLine struct {
	Raw        string
	Depth      int
	Indent     int
	Content    string
	LineNumber int
}

// BlankLineInfo contains information about a blank line
type BlankLineInfo struct {
	LineNumber int
	Indent     int
	Depth      int
}

// ScanResult contains the result of scanning TOON input
type ScanResult struct {
	Lines      []ParsedLine
	BlankLines []BlankLineInfo
}
