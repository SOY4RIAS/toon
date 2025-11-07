package toon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestFixture represents a test fixture file structure
type TestFixture struct {
	Version     string     `json:"version"`
	Category    string     `json:"category"`
	Description string     `json:"description"`
	Tests       []TestCase `json:"tests"`
}

// TestCase represents an individual test case
type TestCase struct {
	Name           string                 `json:"name"`
	Input          interface{}            `json:"input"`
	Expected       interface{}            `json:"expected"`
	ShouldError    bool                   `json:"shouldError,omitempty"`
	Options        map[string]interface{} `json:"options,omitempty"`
	SpecSection    string                 `json:"specSection,omitempty"`
	Note           string                 `json:"note,omitempty"`
	MinSpecVersion string                 `json:"minSpecVersion,omitempty"`
}

// loadFixture loads a test fixture file
func loadFixture(path string) (*TestFixture, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read fixture file: %w", err)
	}

	var fixture TestFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, fmt.Errorf("failed to parse fixture file: %w", err)
	}

	return &fixture, nil
}

// buildOptions converts a map to EncodeOptions or DecodeOptions
func buildEncodeOptions(optMap map[string]interface{}) *EncodeOptions {
	if optMap == nil {
		return nil // Return nil to use default options
	}

	opts := &EncodeOptions{}
	if delimiter, ok := optMap["delimiter"].(string); ok {
		// Convert string delimiter to Delimiter type
		if len(delimiter) > 0 {
			if delimiter == "\\t" || delimiter == "\t" {
				opts.Delimiter = DelimiterTab
			} else {
				opts.Delimiter = Delimiter([]rune(delimiter)[0])
			}
		}
	}
	if indent, ok := optMap["indent"].(float64); ok {
		opts.Indent = int(indent)
	}
	if lengthMarker, ok := optMap["lengthMarker"].(string); ok {
		// Convert string to rune
		if len(lengthMarker) > 0 {
			opts.LengthMarker = []rune(lengthMarker)[0]
		}
	}

	return opts
}

func buildDecodeOptions(optMap map[string]interface{}) *DecodeOptions {
	if optMap == nil {
		return nil // Return nil to use default options (strict: true)
	}

	opts := &DecodeOptions{}

	// Check if strict is explicitly set, otherwise it defaults to false in the struct
	// but we want true as the default, so we need to check if it's in the map
	if strict, ok := optMap["strict"].(bool); ok {
		opts.Strict = strict
	} else {
		// If not specified in options, use default (true)
		opts.Strict = true
	}

	if indent, ok := optMap["indent"].(float64); ok {
		opts.Indent = int(indent)
	}

	return opts
}

// TestConformance_Decode runs all decode conformance tests
func TestConformance_Decode(t *testing.T) {
	fixturesDir := "spec-tests/fixtures/decode"

	// Check if fixtures directory exists
	if _, err := os.Stat(fixturesDir); os.IsNotExist(err) {
		t.Skip("Conformance test fixtures not found. Download fixtures from spec repository.")
		return
	}

	// Get all fixture files
	files, err := filepath.Glob(filepath.Join(fixturesDir, "*.json"))
	if err != nil {
		t.Fatalf("Failed to list fixture files: %v", err)
	}

	if len(files) == 0 {
		t.Skip("No decode fixture files found")
		return
	}

	totalTests := 0
	passedTests := 0
	failedTests := 0
	skippedTests := 0

	for _, file := range files {
		fixture, err := loadFixture(file)
		if err != nil {
			t.Errorf("Failed to load fixture %s: %v", filepath.Base(file), err)
			continue
		}

		t.Run(filepath.Base(file), func(t *testing.T) {
			for _, tc := range fixture.Tests {
				totalTests++

				t.Run(tc.Name, func(t *testing.T) {
					// Skip if minimum spec version is higher than what we support
					if tc.MinSpecVersion != "" && tc.MinSpecVersion > fixture.Version {
						skippedTests++
						t.Skipf("Requires spec version %s (current: %s)", tc.MinSpecVersion, fixture.Version)
						return
					}

					// Get input as string
					var inputStr string
					switch v := tc.Input.(type) {
					case string:
						inputStr = v
					default:
						// For non-string inputs, marshal to JSON first
						data, _ := json.Marshal(v)
						inputStr = string(data)
					}

					// Build options
					opts := buildDecodeOptions(tc.Options)

					// Run decode
					result, err := Decode(inputStr, opts)

					// Check for expected errors
					if tc.ShouldError {
						if err == nil {
							failedTests++
							t.Errorf("Expected error but got none. Result: %v", result)
						} else {
							passedTests++
							t.Logf("✓ Correctly errored: %v", err)
						}
						return
					}

					// Check for unexpected errors
					if err != nil {
						failedTests++
						t.Errorf("Unexpected error: %v", err)
						return
					}

					// Compare result with expected
					if !deepEqual(result, tc.Expected) {
						failedTests++
						t.Errorf("Result mismatch:\nExpected: %#v\nGot:      %#v", tc.Expected, result)
						return
					}

					passedTests++
				})
			}
		})
	}

	// Print summary
	t.Logf("\n=== DECODE CONFORMANCE TEST SUMMARY ===")
	t.Logf("Total:   %d", totalTests)
	t.Logf("Passed:  %d (%.1f%%)", passedTests, float64(passedTests)/float64(totalTests)*100)
	t.Logf("Failed:  %d", failedTests)
	t.Logf("Skipped: %d", skippedTests)
}

// TestConformance_Encode runs all encode conformance tests
func TestConformance_Encode(t *testing.T) {
	fixturesDir := "spec-tests/fixtures/encode"

	// Check if fixtures directory exists
	if _, err := os.Stat(fixturesDir); os.IsNotExist(err) {
		t.Skip("Conformance test fixtures not found. Download fixtures from spec repository.")
		return
	}

	// Get all fixture files
	files, err := filepath.Glob(filepath.Join(fixturesDir, "*.json"))
	if err != nil {
		t.Fatalf("Failed to list fixture files: %v", err)
	}

	if len(files) == 0 {
		t.Skip("No encode fixture files found")
		return
	}

	totalTests := 0
	passedTests := 0
	failedTests := 0
	skippedTests := 0

	for _, file := range files {
		fixture, err := loadFixture(file)
		if err != nil {
			t.Errorf("Failed to load fixture %s: %v", filepath.Base(file), err)
			continue
		}

		t.Run(filepath.Base(file), func(t *testing.T) {
			for _, tc := range fixture.Tests {
				totalTests++

				t.Run(tc.Name, func(t *testing.T) {
					// Skip if minimum spec version is higher than what we support
					if tc.MinSpecVersion != "" && tc.MinSpecVersion > fixture.Version {
						skippedTests++
						t.Skipf("Requires spec version %s (current: %s)", tc.MinSpecVersion, fixture.Version)
						return
					}

					// Build options
					opts := buildEncodeOptions(tc.Options)

					// Run encode
					result, err := Encode(tc.Input, opts)

					// Check for expected errors
					if tc.ShouldError {
						if err == nil {
							failedTests++
							t.Errorf("Expected error but got none. Result: %s", result)
						} else {
							passedTests++
							t.Logf("✓ Correctly errored: %v", err)
						}
						return
					}

					// Check for unexpected errors
					if err != nil {
						failedTests++
						t.Errorf("Unexpected error: %v", err)
						return
					}

					// Get expected output as string
					var expectedStr string
					switch v := tc.Expected.(type) {
					case string:
						expectedStr = v
					default:
						// For non-string expected values, marshal to JSON
						data, _ := json.Marshal(v)
						expectedStr = string(data)
					}

					// Compare result with expected (normalize whitespace)
					resultNorm := normalizeWhitespace(result)
					expectedNorm := normalizeWhitespace(expectedStr)

					if resultNorm != expectedNorm {
						failedTests++
						t.Errorf("Result mismatch:\nExpected:\n%s\n\nGot:\n%s", expectedStr, result)
						return
					}

					passedTests++
				})
			}
		})
	}

	// Print summary
	t.Logf("\n=== ENCODE CONFORMANCE TEST SUMMARY ===")
	t.Logf("Total:   %d", totalTests)
	t.Logf("Passed:  %d (%.1f%%)", passedTests, float64(passedTests)/float64(totalTests)*100)
	t.Logf("Failed:  %d", failedTests)
	t.Logf("Skipped: %d", skippedTests)
}

// deepEqual compares two values for deep equality, handling JSON types
func deepEqual(a, b interface{}) bool {
	// Use reflect.DeepEqual for most cases
	if reflect.DeepEqual(a, b) {
		return true
	}

	// Handle special cases for JSON numbers (float64 vs int)
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	if aVal.Kind() == reflect.Float64 && bVal.Kind() == reflect.Float64 {
		return aVal.Float() == bVal.Float()
	}

	// Handle map comparison
	if aVal.Kind() == reflect.Map && bVal.Kind() == reflect.Map {
		if aVal.Len() != bVal.Len() {
			return false
		}
		for _, key := range aVal.MapKeys() {
			aItem := aVal.MapIndex(key).Interface()
			bItem := bVal.MapIndex(key).Interface()
			if !deepEqual(aItem, bItem) {
				return false
			}
		}
		return true
	}

	// Handle slice comparison
	if aVal.Kind() == reflect.Slice && bVal.Kind() == reflect.Slice {
		if aVal.Len() != bVal.Len() {
			return false
		}
		for i := 0; i < aVal.Len(); i++ {
			if !deepEqual(aVal.Index(i).Interface(), bVal.Index(i).Interface()) {
				return false
			}
		}
		return true
	}

	return false
}

// normalizeWhitespace normalizes whitespace for comparison
func normalizeWhitespace(s string) string {
	// Trim leading/trailing whitespace
	s = strings.TrimSpace(s)
	// Normalize line endings
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return s
}
