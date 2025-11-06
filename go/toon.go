// Package toon provides encoding and decoding for the TOON (Token-Oriented Object Notation) format.
// TOON is a compact, human-readable serialization format designed for passing structured data
// to Large Language Models with significantly reduced token usage.
package toon

// Decode parses a TOON-formatted string and returns the decoded JSON value.
//
// The input string is parsed according to the TOON specification, with optional
// validation and indentation settings.
//
// Options:
//   - Indent: Number of spaces per indentation level (default: 2)
//   - Strict: Enable strict validation of array lengths and structure (default: true)
//
// Example:
//
//	toon := `items[2]{id,name}:
//	  1,Alice
//	  2,Bob`
//	value, err := Decode(toon, nil)
//	// value = map[string]interface{}{
//	//   "items": []interface{}{
//	//     map[string]interface{}{"id": 1.0, "name": "Alice"},
//	//     map[string]interface{}{"id": 2.0, "name": "Bob"},
//	//   },
//	// }
func Decode(input string, options *DecodeOptions) (JsonValue, error) {
	resolvedOptions := resolveDecodeOptions(options)
	scanResult, err := ToParsedLines(input, resolvedOptions.Indent, resolvedOptions.Strict)
	if err != nil {
		return nil, err
	}

	if len(scanResult.Lines) == 0 {
		return make(JsonObject), nil
	}

	cursor := NewLineCursor(scanResult.Lines, scanResult.BlankLines)
	return DecodeValueFromLines(cursor, resolvedOptions)
}

// Encode converts a Go value to TOON format.
// TODO: Implement encoding
func Encode(input interface{}, options *EncodeOptions) (string, error) {
	// TODO: Implement encoder
	return "", nil
}
