package toon

// JsonValue represents any valid JSON value in Go.
// This is the root type for all TOON encode/decode operations.
type JsonValue interface{}

// JsonPrimitive represents JSON primitive types.
type JsonPrimitive interface{}

// JsonObject represents a JSON object as a Go map.
type JsonObject map[string]interface{}

// JsonArray represents a JSON array as a Go slice.
type JsonArray []interface{}

// Delimiter represents the character used to separate values in arrays and tabular rows.
type Delimiter rune

const (
	// DelimiterComma is the default comma delimiter.
	DelimiterComma Delimiter = ','
	// DelimiterTab is the tab character delimiter.
	DelimiterTab Delimiter = '\t'
	// DelimiterPipe is the pipe character delimiter.
	DelimiterPipe Delimiter = '|'
)

// EncodeOptions configures TOON encoding behavior.
type EncodeOptions struct {
	// Indent specifies the number of spaces per indentation level.
	// Default: 2
	Indent int

	// Delimiter specifies the character to use for separating values in arrays and tabular rows.
	// Valid values: DelimiterComma (','), DelimiterTab ('\t'), DelimiterPipe ('|')
	// Default: DelimiterComma
	Delimiter Delimiter

	// LengthMarker when true, prefixes array lengths with '#' (e.g., [#3] instead of [3]).
	// Default: false
	LengthMarker bool
}

// ResolvedEncodeOptions represents fully resolved encoding options with all defaults applied.
type ResolvedEncodeOptions struct {
	Indent       int
	Delimiter    Delimiter
	LengthMarker bool
}

// DecodeOptions configures TOON decoding behavior.
type DecodeOptions struct {
	// Indent specifies the expected number of spaces per indentation level.
	// Default: 2
	Indent int

	// Strict when true, enforces strict validation of array lengths and tabular row counts.
	// Default: true
	Strict bool
}

// ResolvedDecodeOptions represents fully resolved decoding options with all defaults applied.
type ResolvedDecodeOptions struct {
	Indent int
	Strict bool
}

// ArrayHeaderInfo contains parsed information from an array header line.
type ArrayHeaderInfo struct {
	// Key is the property name (empty for root arrays).
	Key string

	// Length is the declared array length from the header.
	Length int

	// Delimiter is the delimiter character detected from the header.
	Delimiter Delimiter

	// Fields contains field names for tabular arrays (nil for non-tabular).
	Fields []string

	// HasLengthMarker indicates if the '#' prefix was present in the length.
	HasLengthMarker bool
}

// ParsedLine represents a single parsed line from the input.
type ParsedLine struct {
	// Raw is the original line content.
	Raw string

	// Depth is the nesting depth (0-based).
	Depth int

	// Indent is the number of leading spaces.
	Indent int

	// Content is the line content after removing indentation.
	Content string

	// LineNumber is the 1-based line number in the original input.
	LineNumber int
}

// BlankLineInfo contains information about blank lines encountered during parsing.
type BlankLineInfo struct {
	// LineNumber is the 1-based line number.
	LineNumber int

	// Indent is the indentation level of the blank line.
	Indent int

	// Depth is the nesting depth.
	Depth int
}

// ScanResult contains the result of scanning/tokenizing TOON input.
type ScanResult struct {
	// Lines contains all non-blank parsed lines.
	Lines []ParsedLine

	// BlankLines contains information about blank lines.
	BlankLines []BlankLineInfo
}
