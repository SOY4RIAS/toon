package toon

import (
	"fmt"
)

// Error types for TOON encoding and decoding operations.

// EncodeError represents an error that occurred during encoding.
type EncodeError struct {
	Message string
	Cause   error
}

func (e *EncodeError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("encode error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("encode error: %s", e.Message)
}

func (e *EncodeError) Unwrap() error {
	return e.Cause
}

// DecodeError represents an error that occurred during decoding.
type DecodeError struct {
	Message    string
	LineNumber int
	Cause      error
}

func (e *DecodeError) Error() string {
	if e.LineNumber > 0 {
		if e.Cause != nil {
			return fmt.Sprintf("decode error at line %d: %s: %v", e.LineNumber, e.Message, e.Cause)
		}
		return fmt.Sprintf("decode error at line %d: %s", e.LineNumber, e.Message)
	}
	if e.Cause != nil {
		return fmt.Sprintf("decode error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("decode error: %s", e.Message)
}

func (e *DecodeError) Unwrap() error {
	return e.Cause
}

// ValidationError represents a validation error (used in strict mode).
type ValidationError struct {
	Message    string
	LineNumber int
}

func (e *ValidationError) Error() string {
	if e.LineNumber > 0 {
		return fmt.Sprintf("validation error at line %d: %s", e.LineNumber, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// Common error constructors

func newEncodeError(message string) error {
	return &EncodeError{Message: message}
}

func newEncodeErrorf(format string, args ...interface{}) error {
	return &EncodeError{Message: fmt.Sprintf(format, args...)}
}

func newDecodeError(message string, lineNumber int) error {
	return &DecodeError{Message: message, LineNumber: lineNumber}
}

func newDecodeErrorf(lineNumber int, format string, args ...interface{}) error {
	return &DecodeError{Message: fmt.Sprintf(format, args...), LineNumber: lineNumber}
}

func newValidationError(message string, lineNumber int) error {
	return &ValidationError{Message: message, LineNumber: lineNumber}
}

func newValidationErrorf(lineNumber int, format string, args ...interface{}) error {
	return &ValidationError{Message: fmt.Sprintf(format, args...), LineNumber: lineNumber}
}
