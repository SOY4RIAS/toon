package toon

// JsonValue represents any valid JSON value
type JsonValue interface{}

// JsonPrimitive can be string, number, bool, or null
type JsonPrimitive interface{}

// JsonObject is a map representing a JSON object
type JsonObject map[string]JsonValue

// JsonArray is a slice representing a JSON array
type JsonArray []JsonValue

// EncodeOptions contains options for encoding
type EncodeOptions struct {
	// Indent is the number of spaces per indentation level (default: 2)
	Indent int

	// Delimiter to use for tabular array rows and inline primitive arrays (default: comma)
	Delimiter Delimiter

	// LengthMarker is an optional marker to prefix array lengths in headers
	// When set to '#', arrays render as [#N] instead of [N]
	LengthMarker rune // '#' or 0 for false
}

// DecodeOptions contains options for decoding
type DecodeOptions struct {
	// Indent is the number of spaces per indentation level (default: 2)
	Indent int

	// Strict enables strict validation of array lengths and tabular row counts (default: true)
	Strict bool
}

// ResolvedEncodeOptions contains resolved encode options with defaults applied
type ResolvedEncodeOptions struct {
	Indent       int
	Delimiter    Delimiter
	LengthMarker rune
}

// ResolvedDecodeOptions contains resolved decode options with defaults applied
type ResolvedDecodeOptions struct {
	Indent int
	Strict bool
}

// ArrayHeaderInfo contains parsed array header information
type ArrayHeaderInfo struct {
	Key             string    // Optional key name
	Length          int       // Array length
	Delimiter       Delimiter // Delimiter used
	Fields          []string  // Optional field names for tabular arrays
	HasLengthMarker bool      // Whether the length marker (#) is present
}

// ParsedLine represents a parsed line with metadata
type ParsedLine struct {
	Raw        string // Original line text
	Depth      int    // Indentation depth level
	Indent     int    // Number of spaces
	Content    string // Line content without leading spaces
	LineNumber int    // 1-based line number
}

// BlankLineInfo contains information about blank lines
type BlankLineInfo struct {
	LineNumber int // 1-based line number
	Indent     int // Number of spaces
	Depth      int // Indentation depth level
}

// resolveEncodeOptions applies defaults to encode options
func resolveEncodeOptions(opts *EncodeOptions) ResolvedEncodeOptions {
	if opts == nil {
		return ResolvedEncodeOptions{
			Indent:       2,
			Delimiter:    DefaultDelimiter,
			LengthMarker: 0,
		}
	}

	indent := opts.Indent
	if indent == 0 {
		indent = 2
	}

	delimiter := opts.Delimiter
	if delimiter == 0 {
		delimiter = DefaultDelimiter
	}

	return ResolvedEncodeOptions{
		Indent:       indent,
		Delimiter:    delimiter,
		LengthMarker: opts.LengthMarker,
	}
}

// resolveDecodeOptions applies defaults to decode options
func resolveDecodeOptions(opts *DecodeOptions) ResolvedDecodeOptions {
	if opts == nil {
		return ResolvedDecodeOptions{
			Indent: 2,
			Strict: true,
		}
	}

	indent := opts.Indent
	if indent == 0 {
		indent = 2
	}

	return ResolvedDecodeOptions{
		Indent: indent,
		Strict: opts.Strict,
	}
}
