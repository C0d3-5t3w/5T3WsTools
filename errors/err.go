// Package errors extends the standard errors package with additional utilities.
package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Standard error handling from the errors package
var (
	Is     = errors.Is
	As     = errors.As
	New    = errors.New
	Unwrap = errors.Unwrap
)

// StackFrame represents a single stack frame
type StackFrame struct {
	File     string
	Function string
	Line     int
}

// Error extends the standard error interface with stack trace capability
type Error struct {
	message    string
	innerError error
	stack      []StackFrame
	context    map[string]interface{}
}

// NewWithStack creates a new error with stack trace information
func NewWithStack(message string) *Error {
	return &Error{
		message: message,
		stack:   captureStack(2), // Skip this function and caller
		context: make(map[string]interface{}),
	}
}

// Wrap wraps an existing error with additional message and stack trace
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &Error{
		message:    message,
		innerError: err,
		stack:      captureStack(2), // Skip this function and caller
		context:    make(map[string]interface{}),
	}
}

// WrapWithContext wraps an error with message and context data
func WrapWithContext(err error, message string, key string, value interface{}) error {
	if err == nil {
		return nil
	}

	e := &Error{
		message:    message,
		innerError: err,
		stack:      captureStack(2),
		context:    make(map[string]interface{}),
	}

	e.context[key] = value
	return e
}

// WithContext adds context to an error
func WithContext(err error, key string, value interface{}) error {
	if err == nil {
		return nil
	}

	var e *Error
	if errors.As(err, &e) {
		e.context[key] = value
		return e
	}

	ne := &Error{
		message:    err.Error(),
		innerError: err,
		stack:      captureStack(2),
		context:    make(map[string]interface{}),
	}

	ne.context[key] = value
	return ne
}

// Error returns the error message
func (e *Error) Error() string {
	if e.innerError != nil {
		return fmt.Sprintf("%s: %v", e.message, e.innerError)
	}
	return e.message
}

// Unwrap returns the wrapped error
func (e *Error) Unwrap() error {
	return e.innerError
}

// Stack returns the error's stack trace
func (e *Error) Stack() []StackFrame {
	return e.stack
}

// Context returns the context value for a given key
func (e *Error) Context(key string) (interface{}, bool) {
	val, ok := e.context[key]
	return val, ok
}

// GetAllContext returns all context values
func (e *Error) GetAllContext() map[string]interface{} {
	// Return a copy to prevent modification of the original
	result := make(map[string]interface{})
	for k, v := range e.context {
		result[k] = v
	}
	return result
}

// FormatStack returns a formatted stack trace string
func (e *Error) FormatStack() string {
	var sb strings.Builder

	for i, frame := range e.stack {
		sb.WriteString(fmt.Sprintf("%d: %s\n   %s:%d\n",
			i, frame.Function, frame.File, frame.Line))
	}

	return sb.String()
}

// IsType checks if an error (or any error in its chain) is of a specific type
func IsType[T error](err error) bool {
	var target T
	return errors.As(err, &target)
}

// captureStack captures the current stack trace
func captureStack(skip int) []StackFrame {
	const depth = 32
	var pcs [depth]uintptr

	// Skip the first 'skip' frames (including this one)
	n := runtime.Callers(skip, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	stack := make([]StackFrame, 0, n)

	for {
		frame, more := frames.Next()
		stack = append(stack, StackFrame{
			File:     frame.File,
			Function: frame.Function,
			Line:     frame.Line,
		})

		if !more {
			break
		}
	}

	return stack
}

// Join concatenates multiple errors into a single error.
// A direct re-export from Go 1.20+ for backward compatibility
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// MustNot panics if the error is not nil
func MustNot(err error) {
	if err != nil {
		panic(err)
	}
}

// GetRootCause returns the innermost error in an error chain
func GetRootCause(err error) error {
	for err != nil {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
	return nil
}
