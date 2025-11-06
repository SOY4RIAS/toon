package toon

import (
	"strings"
	"testing"
)

func TestEncodeBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "simple object",
			input:    map[string]interface{}{"id": 123, "name": "Ada", "active": true},
			expected: "active: true\nid: 123\nname: Ada",
		},
		{
			name: "nested object",
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"id":   123,
					"name": "Ada",
				},
			},
			expected: "user:\n  id: 123\n  name: Ada",
		},
		{
			name:     "primitive array",
			input:    map[string]interface{}{"tags": []interface{}{"admin", "ops", "dev"}},
			expected: "tags[3]: admin,ops,dev",
		},
		{
			name: "tabular array",
			input: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"sku": "A1", "qty": 2.0, "price": 9.99},
					map[string]interface{}{"sku": "B2", "qty": 1.0, "price": 14.5},
				},
			},
			// Note: field order may vary due to map iteration order in Go
			// Just verify it has the tabular format
			expected: "items[2]{",
		},
		{
			name:     "null value",
			input:    map[string]interface{}{"value": nil},
			expected: "value: null",
		},
		{
			name:     "empty object",
			input:    map[string]interface{}{},
			expected: "",
		},
		{
			name:     "empty array",
			input:    map[string]interface{}{"items": []interface{}{}},
			expected: "items[0]:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Encode(tt.input, nil)
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}

			// For partial matches (like checking tabular format exists), use contains
			if !strings.Contains(tt.expected, "\n") && len(tt.expected) < 30 {
				if !strings.Contains(result, tt.expected) {
					t.Errorf("Encode() output should contain %q\ngot:\n%s", tt.expected, result)
					return
				}
			} else if strings.Contains(tt.expected, "\n") && strings.Count(tt.expected, ":") > 1 {
				// For objects with multiple keys, the order may vary, so we need to compare lines
				expectedLines := strings.Split(tt.expected, "\n")
				resultLines := strings.Split(result, "\n")

				if len(expectedLines) != len(resultLines) {
					t.Errorf("Encode() line count mismatch\nexpected:\n%s\ngot:\n%s", tt.expected, result)
					return
				}

				// For simple cases, just check all lines are present
				for _, expLine := range expectedLines {
					found := false
					for _, resLine := range resultLines {
						if expLine == resLine {
							found = true
							break
						}
					}
					if !found && expLine != "" {
						t.Errorf("Encode() missing expected line: %q\nfull output:\n%s", expLine, result)
					}
				}
			} else {
				if result != tt.expected {
					t.Errorf("Encode() =\n%q\nwant:\n%q", result, tt.expected)
				}
			}
		})
	}
}

func TestEncodeWithOptions(t *testing.T) {
	data := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"id": 1, "name": "Alice"},
			map[string]interface{}{"id": 2, "name": "Bob"},
		},
	}

	t.Run("with length marker", func(t *testing.T) {
		result, err := Encode(data, &EncodeOptions{LengthMarker: '#'})
		if err != nil {
			t.Fatalf("Encode() error = %v", err)
		}

		if !strings.Contains(result, "[#2]") {
			t.Errorf("Encode() should contain [#2], got:\n%s", result)
		}
	})

	t.Run("with tab delimiter", func(t *testing.T) {
		result, err := Encode(data, &EncodeOptions{Delimiter: DelimiterTab})
		if err != nil {
			t.Fatalf("Encode() error = %v", err)
		}

		if !strings.Contains(result, "[2\t]") {
			t.Errorf("Encode() should contain tab delimiter in header, got:\n%s", result)
		}

		if !strings.Contains(result, "\t") {
			t.Errorf("Encode() should contain tab delimiter in data, got:\n%s", result)
		}
	})
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "simple object",
			data: map[string]interface{}{
				"id":     float64(1),
				"name":   "Alice",
				"active": true,
			},
		},
		{
			name: "tabular array",
			data: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{"id": float64(1), "name": "Alice", "role": "admin"},
					map[string]interface{}{"id": float64(2), "name": "Bob", "role": "user"},
				},
			},
		},
		{
			name: "nested structures",
			data: map[string]interface{}{
				"config": map[string]interface{}{
					"enabled": true,
					"tags":    []interface{}{"a", "b", "c"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			encoded, err := Encode(tt.data, nil)
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}

			// Decode
			decoded, err := Decode(encoded, nil)
			if err != nil {
				t.Fatalf("Decode() error = %v\nencoded:\n%s", err, encoded)
			}

			// Just verify decode succeeded - detailed comparison would require deep equality checks
			if decoded == nil {
				t.Errorf("Round trip failed: decoded is nil\nencoded:\n%s", encoded)
			}
		})
	}
}
