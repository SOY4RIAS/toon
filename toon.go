package toon

import (
	"github.com/SOY4RIAS/toon/decode"
	"github.com/SOY4RIAS/toon/encode"
)

// Encode converts a Go value to TOON format.
// Returns the TOON-formatted string or an error if encoding fails.
func Encode(value interface{}, opts *EncodeOptions) (string, error) {
	// Normalize the value
	normalized := encode.NormalizeValue(value)

	// Resolve options
	resolved := resolveEncodeOptions(opts)

	// Encode the value
	result := encode.EncodeValue(normalized, resolved.Indent, rune(resolved.Delimiter), resolved.LengthMarker)

	return result, nil
}

// Decode parses a TOON-formatted string and returns the corresponding Go value.
// Returns the decoded value or an error if parsing fails.
func Decode(input string, opts *DecodeOptions) (interface{}, error) {
	resolved := resolveDecodeOptions(opts)

	// Scan the input into lines
	scanResult, err := decode.ToParsedLines(input, resolved.Indent, resolved.Strict)
	if err != nil {
		return nil, err
	}

	if len(scanResult.Lines) == 0 {
		return map[string]interface{}{}, nil
	}

	// Create cursor and decode
	cursor := decode.NewLineCursor(scanResult.Lines, scanResult.BlankLines)
	return decode.DecodeValueFromLines(cursor, resolved.Strict)
}

func resolveEncodeOptions(opts *EncodeOptions) *ResolvedEncodeOptions {
	if opts == nil {
		return &ResolvedEncodeOptions{
			Indent:       DefaultIndent,
			Delimiter:    DefaultDelimiter,
			LengthMarker: false,
		}
	}

	indent := opts.Indent
	if indent == 0 {
		indent = DefaultIndent
	}

	delimiter := opts.Delimiter
	if delimiter == 0 {
		delimiter = DefaultDelimiter
	}

	return &ResolvedEncodeOptions{
		Indent:       indent,
		Delimiter:    delimiter,
		LengthMarker: opts.LengthMarker,
	}
}

func resolveDecodeOptions(opts *DecodeOptions) *ResolvedDecodeOptions {
	if opts == nil {
		return &ResolvedDecodeOptions{
			Indent: DefaultIndent,
			Strict: true,
		}
	}

	indent := opts.Indent
	if indent == 0 {
		indent = DefaultIndent
	}

	return &ResolvedDecodeOptions{
		Indent: indent,
		Strict: opts.Strict,
	}
}
