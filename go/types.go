package toon

// JsonValue represents any valid JSON value
type JsonValue interface{}

// JsonPrimitive represents a JSON primitive value
type JsonPrimitive interface{}

// JsonObject represents a JSON object
type JsonObject map[string]JsonValue

// JsonArray represents a JSON array
type JsonArray []JsonValue

// EncodeOptions configures the encoding behavior
type EncodeOptions struct {
	// Indent is the number of spaces per indentation level (default: 2)
	Indent int
	// Delimiter is the delimiter to use for arrays (default: comma)
	Delimiter Delimiter
	// LengthMarker adds a '#' prefix to array lengths when true
	LengthMarker bool
}

// DecodeOptions configures the decoding behavior
type DecodeOptions struct {
	// Indent is the expected number of spaces per indentation level (default: 2)
	Indent int
	// Strict enables strict validation when true (default: true)
	Strict bool
}

// ArrayHeaderInfo contains parsed array header information
type ArrayHeaderInfo struct {
	Key             string
	Length          int
	Delimiter       Delimiter
	Fields          []string
	HasLengthMarker bool
}

// ParsedLine represents a single parsed line with metadata
type ParsedLine struct {
	Raw        string
	Depth      int
	Indent     int
	Content    string
	LineNumber int
}

// BlankLineInfo tracks blank line information
type BlankLineInfo struct {
	LineNumber int
	Indent     int
	Depth      int
}
