package decode

import (
	"fmt"
	"strings"

	"github.com/toon-format/toon-go"
)

// AssertExpectedCount validates that the actual count matches the expected count in strict mode.
func AssertExpectedCount(actual int, expected int, itemType string, options *toon.DecodeOptions) error {
	if options.Strict && actual != expected {
		return fmt.Errorf("expected %d %s, but got %d", expected, itemType, actual)
	}
	return nil
}

// ValidateNoExtraListItems checks that there are no extra list items beyond the expected count.
func ValidateNoExtraListItems(cursor *LineCursor, itemDepth int, expectedCount int) error {
	if cursor.AtEnd() {
		return nil
	}

	nextLine := cursor.Peek()
	if nextLine != nil && nextLine.Depth == itemDepth && strings.HasPrefix(nextLine.Content, toon.ListItemPrefix) {
		return fmt.Errorf("expected %d list array items, but found more", expectedCount)
	}

	return nil
}

// ValidateNoExtraTabularRows checks that there are no extra tabular rows beyond the expected count.
func ValidateNoExtraTabularRows(cursor *LineCursor, rowDepth int, header *toon.ArrayHeaderInfo) error {
	if cursor.AtEnd() {
		return nil
	}

	nextLine := cursor.Peek()
	if nextLine != nil &&
		nextLine.Depth == rowDepth &&
		!strings.HasPrefix(nextLine.Content, toon.ListItemPrefix) &&
		isDataRow(nextLine.Content, header.Delimiter) {
		return fmt.Errorf("expected %d tabular rows, but found more", header.Length)
	}

	return nil
}

// ValidateNoBlankLinesInRange validates that there are no blank lines within a specific line range.
// In strict mode, blank lines inside arrays/tabular rows are not allowed.
func ValidateNoBlankLinesInRange(
	startLine int,
	endLine int,
	blankLines []toon.BlankLineInfo,
	strict bool,
	context string,
) error {
	if !strict {
		return nil
	}

	// Find blank lines within the range
	// Note: We don't filter by depth because ANY blank line between array items is an error,
	// regardless of its indentation level
	for _, blank := range blankLines {
		if blank.LineNumber > startLine && blank.LineNumber < endLine {
			return fmt.Errorf("line %d: blank lines inside %s are not allowed in strict mode", blank.LineNumber, context)
		}
	}

	return nil
}

// isDataRow checks if a line represents a data row (as opposed to a key-value pair) in a tabular array.
func isDataRow(content string, delimiter toon.Delimiter) bool {
	colonPos := strings.IndexRune(content, toon.Colon)
	delimiterPos := strings.IndexRune(content, rune(delimiter))

	// No colon = definitely a data row
	if colonPos == -1 {
		return true
	}

	// Has delimiter and it comes before colon = data row
	if delimiterPos != -1 && delimiterPos < colonPos {
		return true
	}

	// Colon before delimiter or no delimiter = key-value pair
	return false
}
