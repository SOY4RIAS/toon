package decode

import (
	"testing"

	"github.com/toon-format/toon-go"
)

func TestParseBracketSegment(t *testing.T) {
	tests := []struct {
		name              string
		seg               string
		defaultDelimiter  toon.Delimiter
		expectedLength    int
		expectedDelimiter toon.Delimiter
		expectedMarker    bool
		expectError       bool
	}{
		{
			name:              "simple number",
			seg:               "5",
			defaultDelimiter:  toon.DelimiterComma,
			expectedLength:    5,
			expectedDelimiter: toon.DelimiterComma,
			expectedMarker:    false,
			expectError:       false,
		},
		{
			name:              "with length marker",
			seg:               "#5",
			defaultDelimiter:  toon.DelimiterComma,
			expectedLength:    5,
			expectedDelimiter: toon.DelimiterComma,
			expectedMarker:    true,
			expectError:       false,
		},
		{
			name:              "with tab delimiter",
			seg:               "5\t",
			defaultDelimiter:  toon.DelimiterComma,
			expectedLength:    5,
			expectedDelimiter: toon.DelimiterTab,
			expectedMarker:    false,
			expectError:       false,
		},
		{
			name:              "with pipe delimiter",
			seg:               "5|",
			defaultDelimiter:  toon.DelimiterComma,
			expectedLength:    5,
			expectedDelimiter: toon.DelimiterPipe,
			expectedMarker:    false,
			expectError:       false,
		},
		{
			name:              "with marker and tab delimiter",
			seg:               "#5\t",
			defaultDelimiter:  toon.DelimiterComma,
			expectedLength:    5,
			expectedDelimiter: toon.DelimiterTab,
			expectedMarker:    true,
			expectError:       false,
		},
		{
			name:             "invalid - not a number",
			seg:              "abc",
			defaultDelimiter: toon.DelimiterComma,
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			length, delimiter, hasMarker, err := ParseBracketSegment(tt.seg, tt.defaultDelimiter)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if length != tt.expectedLength {
				t.Errorf("expected length %d, got %d", tt.expectedLength, length)
			}

			if delimiter != tt.expectedDelimiter {
				t.Errorf("expected delimiter %c, got %c", tt.expectedDelimiter, delimiter)
			}

			if hasMarker != tt.expectedMarker {
				t.Errorf("expected marker %v, got %v", tt.expectedMarker, hasMarker)
			}
		})
	}
}

func TestParseDelimitedValues(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter toon.Delimiter
		expected  []string
	}{
		{
			name:      "comma separated",
			input:     "a,b,c",
			delimiter: toon.DelimiterComma,
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "with quoted value containing delimiter",
			input:     "a,\"b,c\",d",
			delimiter: toon.DelimiterComma,
			expected:  []string{"a", "\"b,c\"", "d"},
		},
		{
			name:      "with escaped quote",
			input:     "a,\"b\\\"c\",d",
			delimiter: toon.DelimiterComma,
			expected:  []string{"a", "\"b\\\"c\"", "d"},
		},
		{
			name:      "tab separated",
			input:     "a\tb\tc",
			delimiter: toon.DelimiterTab,
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "pipe separated",
			input:     "a|b|c",
			delimiter: toon.DelimiterPipe,
			expected:  []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDelimitedValues(tt.input, tt.delimiter)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d values, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("value %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestParsePrimitiveToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected interface{}
	}{
		{
			name:     "string",
			token:    "hello",
			expected: "hello",
		},
		{
			name:     "number",
			token:    "42",
			expected: 42.0,
		},
		{
			name:     "negative number",
			token:    "-3.14",
			expected: -3.14,
		},
		{
			name:     "true",
			token:    "true",
			expected: true,
		},
		{
			name:     "false",
			token:    "false",
			expected: false,
		},
		{
			name:     "null",
			token:    "null",
			expected: nil,
		},
		{
			name:     "quoted string",
			token:    "\"hello world\"",
			expected: "hello world",
		},
		{
			name:     "empty string",
			token:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParsePrimitiveToken(tt.token)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %v (%T), got %v (%T)", tt.expected, tt.expected, result, result)
			}
		})
	}
}

func TestParseStringLiteral(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		expected    string
		expectError bool
	}{
		{
			name:     "unquoted",
			token:    "hello",
			expected: "hello",
		},
		{
			name:     "quoted",
			token:    "\"hello\"",
			expected: "hello",
		},
		{
			name:     "quoted with escape",
			token:    "\"hello\\nworld\"",
			expected: "hello\nworld",
		},
		{
			name:        "unterminated quote",
			token:       "\"hello",
			expectError: true,
		},
		{
			name:        "extra characters after quote",
			token:       "\"hello\"world",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseStringLiteral(tt.token)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
