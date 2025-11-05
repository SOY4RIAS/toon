package toon

// JsonValue represents any valid JSON value
type JsonValue interface{}

// JsonPrimitive represents a JSON primitive type
type JsonPrimitive interface{}

// JsonObject represents a JSON object
type JsonObject map[string]JsonValue

// JsonArray represents a JSON array
type JsonArray []JsonValue

// EncodeOptions contains options for encoding
type EncodeOptions struct {
	// Indent is the number of spaces per indentation level (default: 2)
	Indent int
	// Delimiter is the delimiter to use for arrays (default: comma)
	Delimiter Delimiter
	// LengthMarker when true adds '#' prefix to array lengths (default: false)
	LengthMarker bool
}

// DecodeOptions contains options for decoding
type DecodeOptions struct {
	// Indent is the number of spaces per indentation level (default: 2)
	Indent int
	// Strict when true enforces strict validation (default: true)
	Strict bool
}

// ArrayHeaderInfo contains parsed array header information
type ArrayHeaderInfo struct {
	Key             *string
	Length          int
	Delimiter       Delimiter
	Fields          []string
	HasLengthMarker bool
}

// ParsedLine represents a parsed line
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

// ResolvedEncodeOptions returns encode options with defaults applied
func ResolvedEncodeOptions(opts *EncodeOptions) EncodeOptions {
	if opts == nil {
		return EncodeOptions{
			Indent:       2,
			Delimiter:    DefaultDelimiter,
			LengthMarker: false,
		}
	}
	result := *opts
	if result.Indent == 0 {
		result.Indent = 2
	}
	if result.Delimiter == 0 {
		result.Delimiter = DefaultDelimiter
	}
	return result
}

// ResolvedDecodeOptions returns decode options with defaults applied
func ResolvedDecodeOptions(opts *DecodeOptions) DecodeOptions {
	if opts == nil {
		return DecodeOptions{
			Indent: 2,
			Strict: true,
		}
	}
	result := *opts
	if result.Indent == 0 {
		result.Indent = 2
	}
	return result
}
