// Package expvar provides extensions to the standard expvar library
// for exposing application variables for monitoring.
package expvar

import (
	"encoding/json"
	"expvar"
	"sync"
	"sync/atomic"
	"time"
)

// Bool is a boolean variable that satisfies the expvar.Var interface.
type Bool struct {
	value atomic.Value
}

// NewBool creates a new Bool variable.
func NewBool(val bool) *Bool {
	v := &Bool{}
	v.value.Store(val)
	return v
}

// Value returns the current value of the boolean variable.
func (v *Bool) Value() bool {
	return v.value.Load().(bool)
}

// Set sets the value of the boolean variable.
func (v *Bool) Set(val bool) {
	v.value.Store(val)
}

// Toggle toggles the value of the boolean variable.
func (v *Bool) Toggle() {
	v.Set(!v.Value())
}

// String returns the value as a string.
func (v *Bool) String() string {
	if v.Value() {
		return "true"
	}
	return "false"
}

// PublishBool publishes a Bool with the given name.
func PublishBool(name string, val bool) *Bool {
	v := NewBool(val)
	expvar.Publish(name, v)
	return v
}

// Duration is a time.Duration variable that satisfies the expvar.Var interface.
type Duration struct {
	value atomic.Int64
}

// NewDuration creates a new Duration variable.
func NewDuration(val time.Duration) *Duration {
	v := &Duration{}
	v.value.Store(int64(val))
	return v
}

// Value returns the current value of the duration variable.
func (v *Duration) Value() time.Duration {
	return time.Duration(v.value.Load())
}

// Set sets the value of the duration variable.
func (v *Duration) Set(val time.Duration) {
	v.value.Store(int64(val))
}

// Add adds the given duration to the duration variable.
func (v *Duration) Add(val time.Duration) {
	v.value.Add(int64(val))
}

// String returns the value as a string.
func (v *Duration) String() string {
	return v.Value().String()
}

// PublishDuration publishes a Duration with the given name.
func PublishDuration(name string, val time.Duration) *Duration {
	v := NewDuration(val)
	expvar.Publish(name, v)
	return v
}

// Timestamp is a time.Time variable that satisfies the expvar.Var interface.
type Timestamp struct {
	mu    sync.RWMutex
	value time.Time
}

// NewTimestamp creates a new Timestamp variable.
func NewTimestamp(val time.Time) *Timestamp {
	return &Timestamp{value: val}
}

// Now creates a new Timestamp variable with the current time.
func Now() *Timestamp {
	return NewTimestamp(time.Now())
}

// Value returns the current value of the timestamp variable.
func (v *Timestamp) Value() time.Time {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.value
}

// Set sets the value of the timestamp variable.
func (v *Timestamp) Set(val time.Time) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.value = val
}

// String returns the value as a JSON string.
func (v *Timestamp) String() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	b, _ := json.Marshal(v.value)
	return string(b)
}

// PublishTimestamp publishes a Timestamp with the given name.
func PublishTimestamp(name string, val time.Time) *Timestamp {
	v := NewTimestamp(val)
	expvar.Publish(name, v)
	return v
}

// Reset is a helper for Int that supports resetting to zero.
type Reset struct {
	v *expvar.Int
}

// NewReset creates a resettable Int counter.
func NewReset(name string) *Reset {
	return &Reset{v: expvar.NewInt(name)}
}

// Add adds delta to the counter.
func (r *Reset) Add(delta int64) {
	r.v.Add(delta)
}

// Reset resets the counter to zero.
func (r *Reset) Reset() {
	// Read current value, then subtract it to get to 0
	current := r.v.Value()
	r.v.Add(-current)
}

// Value returns the current value.
func (r *Reset) Value() int64 {
	return r.v.Value()
}

// String returns the value as a string.
func (r *Reset) String() string {
	return r.v.String()
}

// PublishReset publishes a resettable counter with the given name.
func PublishReset(name string) *Reset {
	r := &Reset{v: expvar.NewInt(name)}
	return r
}

// GetCounter is a convenience function that gets or creates an expvar.Int.
func GetCounter(name string) *expvar.Int {
	v := expvar.Get(name)
	if v != nil {
		if counter, ok := v.(*expvar.Int); ok {
			return counter
		}
	}
	return expvar.NewInt(name)
}

// Increment increments the named counter by one.
func Increment(name string) {
	GetCounter(name).Add(1)
}

// Decrement decrements the named counter by one.
func Decrement(name string) {
	GetCounter(name).Add(-1)
}

// PublishFunc registers a function that returns an expvar-compatible string.
func PublishFunc(name string, fn func() interface{}) {
	expvar.Publish(name, expvar.Func(func() interface{} {
		return fn()
	}))
}
