package toon

import (
	"fmt"
	"strings"
)

// ScanResult contains the parsed lines and blank line information
type ScanResult struct {
	Lines      []ParsedLine
	BlankLines []BlankLineInfo
}

// LineCursor provides a cursor for iterating through parsed lines
type LineCursor struct {
	lines      []ParsedLine
	index      int
	blankLines []BlankLineInfo
}

// NewLineCursor creates a new LineCursor
func NewLineCursor(lines []ParsedLine, blankLines []BlankLineInfo) *LineCursor {
	return &LineCursor{
		lines:      lines,
		index:      0,
		blankLines: blankLines,
	}
}

// GetBlankLines returns the blank line information
func (c *LineCursor) GetBlankLines() []BlankLineInfo {
	return c.blankLines
}

// Peek returns the current line without advancing
func (c *LineCursor) Peek() *ParsedLine {
	if c.index >= len(c.lines) {
		return nil
	}
	return &c.lines[c.index]
}

// Next returns the current line and advances the cursor
func (c *LineCursor) Next() *ParsedLine {
	if c.index >= len(c.lines) {
		return nil
	}
	line := &c.lines[c.index]
	c.index++
	return line
}

// Current returns the last line that was read
func (c *LineCursor) Current() *ParsedLine {
	if c.index == 0 {
		return nil
	}
	return &c.lines[c.index-1]
}

// Advance moves the cursor forward by one
func (c *LineCursor) Advance() {
	c.index++
}

// AtEnd returns true if the cursor is at the end
func (c *LineCursor) AtEnd() bool {
	return c.index >= len(c.lines)
}

// Length returns the total number of lines
func (c *LineCursor) Length() int {
	return len(c.lines)
}

// PeekAtDepth returns the next line if it's at the target depth
func (c *LineCursor) PeekAtDepth(targetDepth int) *ParsedLine {
	line := c.Peek()
	if line == nil || line.Depth < targetDepth {
		return nil
	}
	if line.Depth == targetDepth {
		return line
	}
	return nil
}

// HasMoreAtDepth returns true if there are more lines at the target depth
func (c *LineCursor) HasMoreAtDepth(targetDepth int) bool {
	return c.PeekAtDepth(targetDepth) != nil
}

// ToParsedLines parses the source string into lines with metadata
func ToParsedLines(source string, indentSize int, strict bool) (*ScanResult, error) {
	// Trim trailing whitespace and leading/trailing newlines
	// But preserve leading spaces/tabs on first line for validation
	source = strings.TrimRight(source, " \t\n\r")
	source = strings.TrimLeft(source, "\n\r")
	if source == "" {
		return &ScanResult{Lines: []ParsedLine{}, BlankLines: []BlankLineInfo{}}, nil
	}

	lines := strings.Split(source, "\n")
	parsed := []ParsedLine{}
	blankLines := []BlankLineInfo{}

	for i, raw := range lines {
		lineNumber := i + 1
		indent := 0

		// Count leading spaces
		for indent < len(raw) && raw[indent] == Space {
			indent++
		}

		content := raw[indent:]

		// Track blank lines
		if strings.TrimSpace(content) == "" {
			depth := computeDepthFromIndent(indent, indentSize)
			blankLines = append(blankLines, BlankLineInfo{
				LineNumber: lineNumber,
				Indent:     indent,
				Depth:      depth,
			})
			continue
		}

		depth := computeDepthFromIndent(indent, indentSize)

		// Strict mode validation
		if strict {
			// Find the full leading whitespace region (spaces and tabs)
			wsEnd := 0
			for wsEnd < len(raw) && (raw[wsEnd] == Space || raw[wsEnd] == Tab) {
				wsEnd++
			}

			// Check for tabs in leading whitespace (before actual content)
			if strings.Contains(raw[:wsEnd], string(Tab)) {
				return nil, fmt.Errorf("line %d: tabs are not allowed in indentation in strict mode", lineNumber)
			}

			// Check for exact multiples of indentSize
			if indent > 0 && indent%indentSize != 0 {
				return nil, fmt.Errorf("line %d: indentation must be exact multiple of %d, but found %d spaces", lineNumber, indentSize, indent)
			}
		}

		parsed = append(parsed, ParsedLine{
			Raw:        raw,
			Indent:     indent,
			Content:    content,
			Depth:      depth,
			LineNumber: lineNumber,
		})
	}

	return &ScanResult{Lines: parsed, BlankLines: blankLines}, nil
}

func computeDepthFromIndent(indentSpaces int, indentSize int) int {
	return indentSpaces / indentSize
}
