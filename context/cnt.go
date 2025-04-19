// Package cnt provides extensions to the standard library context package
package cnt

import (
	"context"
	"errors"
	"time"
)

// Common errors
var (
	ErrValueNotFound = errors.New("value not found in context")
	ErrNoDeadline    = errors.New("context has no deadline")
)

// WithTimeoutIfNone creates a context with a timeout only if the parent context doesn't already have one
func WithTimeoutIfNone(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if _, ok := parent.Deadline(); ok {
		// Parent already has a deadline, no need to create another one
		return parent, func() {}
	}
	return context.WithTimeout(parent, timeout)
}

// MergeContexts creates a new context that inherits cancellation from multiple contexts
// The cancellation of any parent context will cancel the resulting context
func MergeContexts(parents ...context.Context) (context.Context, context.CancelFunc) {
	if len(parents) == 0 {
		return context.Background(), func() {}
	}

	if len(parents) == 1 {
		return parents[0], func() {}
	}

	// Use the first context as the base
	ctx, cancel := context.WithCancel(parents[0])

	// Monitor cancellations of all parent contexts
	for _, parent := range parents {
		go func(p context.Context) {
			<-p.Done()
			cancel()
		}(parent)
	}

	return ctx, cancel
}

// GetStringValue retrieves a string value from the context or returns an error if not found
func GetStringValue(ctx context.Context, key interface{}) (string, error) {
	value := ctx.Value(key)
	if value == nil {
		return "", ErrValueNotFound
	}

	str, ok := value.(string)
	if !ok {
		return "", errors.New("value is not a string")
	}

	return str, nil
}

// GetStringValueWithDefault retrieves a string value from the context or returns the default value
func GetStringValueWithDefault(ctx context.Context, key interface{}, defaultValue string) string {
	value, err := GetStringValue(ctx, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// RemainingTime returns the amount of time remaining before the context's deadline
func RemainingTime(ctx context.Context) (time.Duration, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return 0, ErrNoDeadline
	}
	return time.Until(deadline), nil
}

// IsDeadlineApproaching checks if the context's deadline is approaching within the given threshold
func IsDeadlineApproaching(ctx context.Context, threshold time.Duration) (bool, error) {
	remaining, err := RemainingTime(ctx)
	if err != nil {
		return false, err
	}
	return remaining <= threshold, nil
}

// WithValues adds multiple key-value pairs to a context in a single call
func WithValues(parent context.Context, keyVals ...interface{}) context.Context {
	if len(keyVals)%2 != 0 {
		panic("WithValues requires an even number of arguments")
	}

	ctx := parent
	for i := 0; i < len(keyVals); i += 2 {
		ctx = context.WithValue(ctx, keyVals[i], keyVals[i+1])
	}

	return ctx
}

// WithCancel creates a cancellable context with an optional onCancel callback
func WithCancel(parent context.Context, onCancel func()) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	if onCancel == nil {
		return ctx, cancel
	}

	cancelFunc := func() {
		cancel()
		onCancel()
	}

	return ctx, cancelFunc
}
