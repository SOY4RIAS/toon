// Package toon provides encoding and decoding for the TOON (Terse Object Oriented Notation) format.
//
// TOON is a human-readable data serialization format designed for clarity and compactness.
// It supports objects, arrays, and primitive types with a clean, indented syntax.
//
// Example:
//
//	input := `
//	name: John Doe
//	age: 30
//	active: true
//	tags[3]: web, backend, api
//	`
//
//	result, err := toon.Decode(input, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
package toon

// Decode decodes a TOON-formatted string into a Go value.
//
// The input string is parsed according to TOON syntax rules, and the resulting
// value is returned as an interface{} which can be type-asserted to the appropriate type.
//
// Options can be provided to customize the decoding behavior:
//   - Indent: Number of spaces per indentation level (default: 2)
//   - Strict: Enforce strict validation (default: true)
//
// Returns an error if the input is invalid or cannot be parsed.
func Decode(input string, options *DecodeOptions) (interface{}, error) {
	resolvedOptions := resolveDecodeOptions(options)

	scanResult, err := ToParsedLines(input, resolvedOptions.Indent, resolvedOptions.Strict)
	if err != nil {
		return nil, err
	}

	if len(scanResult.Lines) == 0 {
		return make(map[string]interface{}), nil
	}

	cursor := NewLineCursor(scanResult.Lines, scanResult.BlankLines)
	return DecodeValueFromLines(cursor, resolvedOptions)
}

// resolveDecodeOptions resolves decode options with defaults
func resolveDecodeOptions(options *DecodeOptions) ResolvedDecodeOptions {
	// Default values
	indent := 2
	strict := true

	if options != nil {
		if options.Indent != nil {
			indent = *options.Indent
		}
		if options.Strict != nil {
			strict = *options.Strict
		}
	}

	return ResolvedDecodeOptions{
		Indent: indent,
		Strict: strict,
	}
}
