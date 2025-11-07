package toon

import (
	"fmt"
	"strings"
)

// AssertExpectedCount checks if actual count matches expected in strict mode
func AssertExpectedCount(actual, expected int, itemType string, opts ResolvedDecodeOptions) error {
	if opts.Strict && actual != expected {
		return fmt.Errorf("expected %d %s, but got %d", expected, itemType, actual)
	}
	return nil
}

// ValidateNoExtraListItems checks for extra list items beyond expected count
func ValidateNoExtraListItems(cursor *LineCursor, itemDepth int, expectedCount int) error {
	if cursor.AtEnd() {
		return nil
	}

	nextLine := cursor.Peek()
	if nextLine != nil && nextLine.Depth == itemDepth && strings.HasPrefix(nextLine.Content, ListItemPrefix) {
		return fmt.Errorf("expected %d list array items, but found more", expectedCount)
	}
	return nil
}

// ValidateNoExtraTabularRows checks for extra tabular rows beyond expected count
func ValidateNoExtraTabularRows(cursor *LineCursor, rowDepth int, header *ArrayHeaderInfo) error {
	if cursor.AtEnd() {
		return nil
	}

	nextLine := cursor.Peek()
	if nextLine != nil &&
		nextLine.Depth == rowDepth &&
		!strings.HasPrefix(nextLine.Content, ListItemPrefix) &&
		isDataRow(nextLine.Content, header.Delimiter) {
		return fmt.Errorf("expected %d tabular rows, but found more", header.Length)
	}
	return nil
}

// ValidateNoBlankLinesInRange checks for blank lines in a specific range
func ValidateNoBlankLinesInRange(startLine, endLine int, blankLines []BlankLineInfo, strict bool, context string) error {
	if !strict {
		return nil
	}

	// Find blank lines within the range
	for _, blank := range blankLines {
		if blank.LineNumber > startLine && blank.LineNumber < endLine {
			return fmt.Errorf("line %d: blank lines inside %s are not allowed in strict mode", blank.LineNumber, context)
		}
	}

	return nil
}

// isDataRow checks if a line is a data row (vs key-value pair)
func isDataRow(content string, delimiter Delimiter) bool {
	colonPos := strings.Index(content, ":")
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
