// Package sync provides extensions to Go's standard sync and sync/atomic packages.
// It offers additional synchronization primitives, atomic operations, and thread-safe
// data structures beyond what's provided in the standard library.
package syncExt

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Common errors
var (
	ErrTimeout  = errors.New("operation timed out")
	ErrCanceled = errors.New("operation was canceled")
)

// TimeoutMutex extends sync.Mutex with timeout capabilities.
type TimeoutMutex struct {
	mu sync.Mutex
}

// Lock locks the mutex.
func (m *TimeoutMutex) Lock() {
	m.mu.Lock()
}

// Unlock unlocks the mutex.
func (m *TimeoutMutex) Unlock() {
	m.mu.Unlock()
}

// TryLock attempts to lock the mutex and returns immediately.
// It returns true if the lock was acquired.
func (m *TimeoutMutex) TryLock() bool {
	// Use a channel to communicate success
	locked := make(chan bool, 1)

	go func() {
		m.mu.Lock()
		locked <- true
	}()

	select {
	case <-locked:
		return true
	default:
		return false
	}
}

// LockWithTimeout attempts to lock the mutex and times out after the specified duration.
// It returns nil if the lock was acquired, otherwise ErrTimeout.
func (m *TimeoutMutex) LockWithTimeout(timeout time.Duration) error {
	// Use a channel to communicate success
	locked := make(chan bool, 1)

	go func() {
		m.mu.Lock()
		locked <- true
	}()

	select {
	case <-locked:
		return nil
	case <-time.After(timeout):
		return ErrTimeout
	}
}

// LockWithContext attempts to lock the mutex and respects context cancellation.
// It returns nil if the lock was acquired, otherwise ErrCanceled or the context error.
func (m *TimeoutMutex) LockWithContext(ctx context.Context) error {
	// Use a channel to communicate success
	locked := make(chan bool, 1)

	go func() {
		m.mu.Lock()
		locked <- true
	}()

	select {
	case <-locked:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// AtomicBool is a boolean value that can be updated atomically.
type AtomicBool struct {
	value uint32
}

// Set sets the value atomically.
func (b *AtomicBool) Set(value bool) {
	var i uint32 = 0
	if value {
		i = 1
	}
	atomic.StoreUint32(&b.value, i)
}

// Get gets the value atomically.
func (b *AtomicBool) Get() bool {
	return atomic.LoadUint32(&b.value) != 0
}

// Toggle atomically toggles the boolean and returns the new value.
func (b *AtomicBool) Toggle() bool {
	for {
		old := atomic.LoadUint32(&b.value)
		new := uint32(1)
		if old != 0 {
			new = 0
		}
		if atomic.CompareAndSwapUint32(&b.value, old, new) {
			return new != 0
		}
	}
}

// AtomicInt64 extends atomic operations for int64.
type AtomicInt64 struct {
	value int64
}

// Get gets the value atomically.
func (i *AtomicInt64) Get() int64 {
	return atomic.LoadInt64(&i.value)
}

// Set sets the value atomically.
func (i *AtomicInt64) Set(n int64) {
	atomic.StoreInt64(&i.value, n)
}

// Add adds delta to the value atomically and returns the new value.
func (i *AtomicInt64) Add(delta int64) int64 {
	return atomic.AddInt64(&i.value, delta)
}

// Increment increments the value atomically and returns the new value.
func (i *AtomicInt64) Increment() int64 {
	return atomic.AddInt64(&i.value, 1)
}

// Decrement decrements the value atomically and returns the new value.
func (i *AtomicInt64) Decrement() int64 {
	return atomic.AddInt64(&i.value, -1)
}

// CompareAndSwap atomically compares the value with old and,
// if equal, swaps it with new and returns true.
func (i *AtomicInt64) CompareAndSwap(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&i.value, old, new)
}

// Map is a generic thread-safe map using sync.RWMutex.
type Map struct {
	mu   sync.RWMutex
	data map[interface{}]interface{}
}

// NewMap creates a new thread-safe map.
func NewMap() *Map {
	return &Map{
		data: make(map[interface{}]interface{}),
	}
}

// Get retrieves a value from the map.
func (m *Map) Get(key interface{}) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

// Set stores a value in the map.
func (m *Map) Set(key, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Delete removes a key-value pair from the map.
func (m *Map) Delete(key interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

// Len returns the number of items in the map.
func (m *Map) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// ForEach executes a function for each key-value pair in the map.
func (m *Map) ForEach(fn func(key, value interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.data {
		if !fn(k, v) {
			break
		}
	}
}

// Once extends sync.Once to provide a reset capability.
type Once struct {
	done uint32
	m    sync.Mutex
}

// Do executes the function f only once until Reset is called.
func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

// Reset allows the Once to be used again.
func (o *Once) Reset() {
	o.m.Lock()
	defer o.m.Unlock()
	atomic.StoreUint32(&o.done, 0)
}

// WaitGroup extends sync.WaitGroup with timeout capabilities.
type WaitGroup struct {
	wg sync.WaitGroup
}

// Add adds delta to the WaitGroup counter.
func (wg *WaitGroup) Add(delta int) {
	wg.wg.Add(delta)
}

// Done decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() {
	wg.wg.Done()
}

// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	wg.wg.Wait()
}

// WaitWithTimeout waits until the WaitGroup counter is zero or times out.
// It returns true if the wait completed successfully, false if it timed out.
func (wg *WaitGroup) WaitWithTimeout(timeout time.Duration) bool {
	c := make(chan struct{})

	go func() {
		defer close(c)
		wg.wg.Wait()
	}()

	select {
	case <-c:
		return true
	case <-time.After(timeout):
		return false
	}
}
