// Package tst extends the functionality of the standard Go testing library
// with additional helper functions for assertions, test organization, and reporting.
package testingExt

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// Assert fails the test if the condition is false.
func Assert(t *testing.T, condition bool, msg string, args ...interface{}) {
	t.Helper()
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\nAssertion failed at %s:%d\n"+msg, append([]interface{}{filepath.Base(file), line}, args...)...)
	}
}

// Equals checks if expected and actual are equal and fails the test if they are not.
func Equals(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(1)
		msg := fmt.Sprintf("\nAssertion failed at %s:%d\nExpected: %v\nActual:   %v",
			filepath.Base(file), line, expected, actual)

		if len(msgAndArgs) > 0 {
			if str, ok := msgAndArgs[0].(string); ok {
				msg += "\n" + fmt.Sprintf(str, msgAndArgs[1:]...)
			} else {
				msg += "\n" + fmt.Sprint(msgAndArgs...)
			}
		}
		t.Error(msg)
	}
}

// NotEquals checks if expected and actual are not equal and fails the test if they are.
func NotEquals(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(1)
		msg := fmt.Sprintf("\nAssertion failed at %s:%d\nExpected values to differ, but both are: %v",
			filepath.Base(file), line, expected)

		if len(msgAndArgs) > 0 {
			if str, ok := msgAndArgs[0].(string); ok {
				msg += "\n" + fmt.Sprintf(str, msgAndArgs[1:]...)
			} else {
				msg += "\n" + fmt.Sprint(msgAndArgs...)
			}
		}
		t.Error(msg)
	}
}

// Nil checks if the value is nil and fails the test if it is not.
func Nil(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !isNil(value) {
		_, file, line, _ := runtime.Caller(1)
		msg := fmt.Sprintf("\nAssertion failed at %s:%d\nExpected nil, but got: %v",
			filepath.Base(file), line, value)

		if len(msgAndArgs) > 0 {
			if str, ok := msgAndArgs[0].(string); ok {
				msg += "\n" + fmt.Sprintf(str, msgAndArgs[1:]...)
			} else {
				msg += "\n" + fmt.Sprint(msgAndArgs...)
			}
		}
		t.Error(msg)
	}
}

// NotNil checks if the value is not nil and fails the test if it is nil.
func NotNil(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if isNil(value) {
		_, file, line, _ := runtime.Caller(1)
		msg := fmt.Sprintf("\nAssertion failed at %s:%d\nExpected non-nil value",
			filepath.Base(file), line)

		if len(msgAndArgs) > 0 {
			if str, ok := msgAndArgs[0].(string); ok {
				msg += "\n" + fmt.Sprintf(str, msgAndArgs[1:]...)
			} else {
				msg += "\n" + fmt.Sprint(msgAndArgs...)
			}
		}
		t.Error(msg)
	}
}

// Error checks if the error is not nil and fails the test if it is nil.
func Error(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		msg := fmt.Sprintf("\nAssertion failed at %s:%d\nExpected error, but got nil",
			filepath.Base(file), line)

		if len(msgAndArgs) > 0 {
			if str, ok := msgAndArgs[0].(string); ok {
				msg += "\n" + fmt.Sprintf(str, msgAndArgs[1:]...)
			} else {
				msg += "\n" + fmt.Sprint(msgAndArgs...)
			}
		}
		t.Error(msg)
	}
}

// NoError checks if the error is nil and fails the test if it is not nil.
func NoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		msg := fmt.Sprintf("\nAssertion failed at %s:%d\nExpected no error, but got: %v",
			filepath.Base(file), line, err)

		if len(msgAndArgs) > 0 {
			if str, ok := msgAndArgs[0].(string); ok {
				msg += "\n" + fmt.Sprintf(str, msgAndArgs[1:]...)
			} else {
				msg += "\n" + fmt.Sprint(msgAndArgs...)
			}
		}
		t.Error(msg)
	}
}

// Contains checks if the string contains the substring and fails the test if it does not.
func Contains(t *testing.T, str, substr string, msgAndArgs ...interface{}) {
	t.Helper()
	if !strings.Contains(str, substr) {
		_, file, line, _ := runtime.Caller(1)
		msg := fmt.Sprintf("\nAssertion failed at %s:%d\nExpected substring: %q\nNot found in: %q",
			filepath.Base(file), line, substr, str)

		if len(msgAndArgs) > 0 {
			if str, ok := msgAndArgs[0].(string); ok {
				msg += "\n" + fmt.Sprintf(str, msgAndArgs[1:]...)
			} else {
				msg += "\n" + fmt.Sprint(msgAndArgs...)
			}
		}
		t.Error(msg)
	}
}

// helper function to check if a value is nil
func isNil(value interface{}) bool {
	if value == nil {
		return true
	}

	// For interfaces, need to use reflection
	valueOf := reflect.ValueOf(value)
	switch valueOf.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return valueOf.IsNil()
	}

	return false
}
