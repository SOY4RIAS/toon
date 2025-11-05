package types

// ParsedLine represents a single parsed line from the input.
type ParsedLine struct {
	// Raw is the original line content.
	Raw string

	// Depth is the nesting depth (0-based).
	Depth int

	// Indent is the number of leading spaces.
	Indent int

	// Content is the line content after removing indentation.
	Content string

	// LineNumber is the 1-based line number in the original input.
	LineNumber int
}

// BlankLineInfo contains information about blank lines encountered during parsing.
type BlankLineInfo struct {
	// LineNumber is the 1-based line number.
	LineNumber int

	// Indent is the indentation level of the blank line.
	Indent int

	// Depth is the nesting depth.
	Depth int
}

// ScanResult contains the result of scanning/tokenizing TOON input.
type ScanResult struct {
	// Lines contains all non-blank parsed lines.
	Lines []ParsedLine

	// BlankLines contains information about blank lines.
	BlankLines []BlankLineInfo
}
