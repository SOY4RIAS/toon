package decode

import (
	"fmt"
	"strings"

	"github.com/toon-format/toon-go"
)

// ScanResult contains the result of scanning input
type ScanResult struct {
	Lines      []toon.ParsedLine
	BlankLines []toon.BlankLineInfo
}

// LineCursor tracks the current position when parsing lines
type LineCursor struct {
	lines      []toon.ParsedLine
	blankLines []toon.BlankLineInfo
	index      int
}

// NewLineCursor creates a new line cursor
func NewLineCursor(lines []toon.ParsedLine, blankLines []toon.BlankLineInfo) *LineCursor {
	return &LineCursor{
		lines:      lines,
		blankLines: blankLines,
		index:      0,
	}
}

// GetBlankLines returns the blank lines information
func (c *LineCursor) GetBlankLines() []toon.BlankLineInfo {
	return c.blankLines
}

// Peek returns the current line without advancing
func (c *LineCursor) Peek() *toon.ParsedLine {
	if c.index >= len(c.lines) {
		return nil
	}
	return &c.lines[c.index]
}

// Next returns the current line and advances
func (c *LineCursor) Next() *toon.ParsedLine {
	if c.index >= len(c.lines) {
		return nil
	}
	line := &c.lines[c.index]
	c.index++
	return line
}

// Current returns the previously consumed line
func (c *LineCursor) Current() *toon.ParsedLine {
	if c.index == 0 {
		return nil
	}
	return &c.lines[c.index-1]
}

// Advance moves to the next line
func (c *LineCursor) Advance() {
	c.index++
}

// AtEnd returns true if there are no more lines
func (c *LineCursor) AtEnd() bool {
	return c.index >= len(c.lines)
}

// Length returns the total number of lines
func (c *LineCursor) Length() int {
	return len(c.lines)
}

// PeekAtDepth returns the current line if it's at the target depth
func (c *LineCursor) PeekAtDepth(targetDepth int) *toon.ParsedLine {
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

// ToParsedLines converts source string into parsed lines
func ToParsedLines(source string, indentSize int, strict bool) (*ScanResult, error) {
	if strings.TrimSpace(source) == "" {
		return &ScanResult{
			Lines:      []toon.ParsedLine{},
			BlankLines: []toon.BlankLineInfo{},
		}, nil
	}

	lines := strings.Split(source, "\n")
	parsed := []toon.ParsedLine{}
	blankLines := []toon.BlankLineInfo{}

	for i, raw := range lines {
		lineNumber := i + 1
		indent := 0

		// Count leading spaces
		for indent < len(raw) && raw[indent] == ' ' {
			indent++
		}

		content := raw[indent:]

		// Track blank lines
		if strings.TrimSpace(content) == "" {
			depth := computeDepthFromIndent(indent, indentSize)
			blankLines = append(blankLines, toon.BlankLineInfo{
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
				return nil, fmt.Errorf("line %d: tabs are not allowed in indentation in strict mode", lineNumber)
			}

			// Check for exact multiples of indentSize
			if indent > 0 && indent%indentSize != 0 {
				return nil, fmt.Errorf("line %d: indentation must be exact multiple of %d, but found %d spaces", lineNumber, indentSize, indent)
			}
		}

		parsed = append(parsed, toon.ParsedLine{
			Raw:        raw,
			Indent:     indent,
			Content:    content,
			Depth:      depth,
			LineNumber: lineNumber,
		})
	}

	return &ScanResult{
		Lines:      parsed,
		BlankLines: blankLines,
	}, nil
}

func computeDepthFromIndent(indentSpaces int, indentSize int) int {
	return indentSpaces / indentSize
}
