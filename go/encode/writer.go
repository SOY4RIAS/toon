package encode

import (
	"strings"
)

// LineWriter accumulates output lines with proper indentation
type LineWriter struct {
	lines             []string
	indentationString string
}

// NewLineWriter creates a new LineWriter with the specified indentation size
func NewLineWriter(indentSize int) *LineWriter {
	return &LineWriter{
		lines:             make([]string, 0),
		indentationString: strings.Repeat(" ", indentSize),
	}
}

// Push adds a line with the specified depth and content
func (w *LineWriter) Push(depth int, content string) {
	indent := strings.Repeat(w.indentationString, depth)
	w.lines = append(w.lines, indent+content)
}

// PushListItem adds a list item line (with "- " prefix) at the specified depth
func (w *LineWriter) PushListItem(depth int, content string) {
	w.Push(depth, listItemPrefix+content)
}

// String returns the complete output as a single string
func (w *LineWriter) String() string {
	return strings.Join(w.lines, "\n")
}
