// Package nt provides extensions to Go's standard net and net/http packages
package nt

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"
)

// DefaultTimeout is the default timeout for HTTP requests
const DefaultTimeout = 30 * time.Second

// Client wraps http.Client with additional functionality
type Client struct {
	*http.Client
	DefaultHeaders map[string]string
	RetryCount     int
	RetryDelay     time.Duration
}

// NewClient creates a new extended HTTP client
func NewClient(timeout time.Duration, retryCount int, retryDelay time.Duration) *Client {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	return &Client{
		Client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		DefaultHeaders: make(map[string]string),
		RetryCount:     retryCount,
		RetryDelay:     retryDelay,
	}
}

// Get performs an HTTP GET request with retries and default headers
func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Apply default headers
	for key, value := range c.DefaultHeaders {
		req.Header.Set(key, value)
	}

	return c.DoWithRetries(req)
}

// Post performs an HTTP POST request with retries and default headers
func (c *Client) Post(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	// Apply default headers
	for key, value := range c.DefaultHeaders {
		if key != "Content-Type" { // Don't override content-type if explicitly set
			req.Header.Set(key, value)
		}
	}

	return c.DoWithRetries(req)
}

// DoWithRetries performs an HTTP request with configured retry logic
func (c *Client) DoWithRetries(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	attempts := c.RetryCount + 1 // Initial attempt plus retries
	for i := 0; i < attempts; i++ {
		resp, err = c.Do(req)

		// If successful or context canceled, return immediately
		if err == nil || req.Context().Err() != nil {
			return resp, err
		}

		// Don't sleep after the last attempt
		if i < attempts-1 {
			time.Sleep(c.RetryDelay)
		}
	}

	return resp, err
}

// IsTemporaryError checks if a network error is temporary
func IsTemporaryError(err error) bool {
	if tempErr, ok := err.(interface{ Temporary() bool }); ok {
		return tempErr.Temporary()
	}
	return false
}

// WaitForPort waits for a TCP port to be available
func WaitForPort(address string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, time.Second)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return &net.OpError{
		Op:   "dial",
		Net:  "tcp",
		Addr: nil,
		Err:  &timeoutError{},
	}
}

// IsPortOpen checks if a TCP port is open
func IsPortOpen(address string) bool {
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// timeoutError implements net.Error interface
type timeoutError struct{}

func (e *timeoutError) Error() string   { return "operation timed out" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }
