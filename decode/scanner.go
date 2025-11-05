package decode

import (
	"fmt"
	"strings"

	"github.com/SOY4RIAS/toon/internal/types"
)

// LineCursor provides cursor-based navigation through parsed lines.
type LineCursor struct {
	lines      []types.ParsedLine
	index      int
	blankLines []types.BlankLineInfo
}

// NewLineCursor creates a new line cursor.
func NewLineCursor(lines []types.ParsedLine, blankLines []types.BlankLineInfo) *LineCursor {
	return &LineCursor{
		lines:      lines,
		index:      0,
		blankLines: blankLines,
	}
}

// GetBlankLines returns information about blank lines.
func (c *LineCursor) GetBlankLines() []types.BlankLineInfo {
	return c.blankLines
}

// Peek returns the current line without advancing the cursor.
func (c *LineCursor) Peek() *types.ParsedLine {
	if c.index >= len(c.lines) {
		return nil
	}
	return &c.lines[c.index]
}

// Next returns the current line and advances the cursor.
func (c *LineCursor) Next() *types.ParsedLine {
	if c.index >= len(c.lines) {
		return nil
	}
	line := &c.lines[c.index]
	c.index++
	return line
}

// Current returns the previously consumed line (the line before the cursor).
func (c *LineCursor) Current() *types.ParsedLine {
	if c.index <= 0 {
		return nil
	}
	return &c.lines[c.index-1]
}

// Advance moves the cursor forward one position.
func (c *LineCursor) Advance() {
	c.index++
}

// AtEnd returns true if the cursor is at or past the end of the lines.
func (c *LineCursor) AtEnd() bool {
	return c.index >= len(c.lines)
}

// Length returns the total number of lines.
func (c *LineCursor) Length() int {
	return len(c.lines)
}

// PeekAtDepth returns the current line if it's at the target depth, nil otherwise.
func (c *LineCursor) PeekAtDepth(targetDepth int) *types.ParsedLine {
	line := c.Peek()
	if line == nil || line.Depth < targetDepth {
		return nil
	}
	if line.Depth == targetDepth {
		return line
	}
	return nil
}

// HasMoreAtDepth returns true if there are more lines at the target depth.
func (c *LineCursor) HasMoreAtDepth(targetDepth int) bool {
	return c.PeekAtDepth(targetDepth) != nil
}

// ToParsedLines scans the source string and returns parsed lines.
func ToParsedLines(source string, indentSize int, strict bool) (types.ScanResult, error) {
	if strings.TrimSpace(source) == "" {
		return types.ScanResult{Lines: []types.ParsedLine{}, BlankLines: []types.BlankLineInfo{}}, nil
	}

	lines := strings.Split(source, "\n")
	parsed := make([]types.ParsedLine, 0, len(lines))
	blankLines := make([]types.BlankLineInfo, 0)

	for i, raw := range lines {
		lineNumber := i + 1

		// Count leading spaces
		indent := 0
		for indent < len(raw) && raw[indent] == ' ' {
			indent++
		}

		content := raw[indent:]

		// Track blank lines
		if strings.TrimSpace(content) == "" {
			depth := computeDepthFromIndent(indent, indentSize)
			blankLines = append(blankLines, types.BlankLineInfo{
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
			for wsEnd < len(raw) && (raw[wsEnd] == ' ' || raw[wsEnd] == '\t') {
				wsEnd++
			}

			// Check for tabs in leading whitespace (before actual content)
			if strings.Contains(raw[:wsEnd], "\t") {
				return types.ScanResult{}, fmt.Errorf("line %d: tabs are not allowed in indentation in strict mode", lineNumber)
			}

			// Check for exact multiples of indentSize
			if indent > 0 && indent%indentSize != 0 {
				return types.ScanResult{}, fmt.Errorf("line %d: indentation must be exact multiple of %d, but found %d spaces", lineNumber, indentSize, indent)
			}
		}

		parsed = append(parsed, types.ParsedLine{
			Raw:        raw,
			Indent:     indent,
			Content:    content,
			Depth:      depth,
			LineNumber: lineNumber,
		})
	}

	return types.ScanResult{Lines: parsed, BlankLines: blankLines}, nil
}

func computeDepthFromIndent(indentSpaces int, indentSize int) int {
	return indentSpaces / indentSize
}
