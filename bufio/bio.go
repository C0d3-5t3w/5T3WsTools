// Package bufio provides extensions to the standard bufio package
package bufio

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// EnhancedReader extends bufio.Reader with additional functionality
type EnhancedReader struct {
	*bufio.Reader
}

// NewEnhancedReader creates and returns a new EnhancedReader
func NewEnhancedReader(r io.Reader) *EnhancedReader {
	return &EnhancedReader{bufio.NewReader(r)}
}

// ReadAllLines reads all lines from the reader and returns them as a slice of strings
func (er *EnhancedReader) ReadAllLines() ([]string, error) {
	var lines []string
	for {
		line, err := er.ReadString('\n')
		if err != nil && err != io.EOF {
			return lines, err
		}

		line = strings.TrimRight(line, "\r\n")
		lines = append(lines, line)

		if err == io.EOF {
			break
		}
	}
	return lines, nil
}

// EnhancedWriter extends bufio.Writer with additional functionality
type EnhancedWriter struct {
	*bufio.Writer
}

// NewEnhancedWriter creates and returns a new EnhancedWriter
func NewEnhancedWriter(w io.Writer) *EnhancedWriter {
	return &EnhancedWriter{bufio.NewWriter(w)}
}

// WriteLines writes a slice of strings to the writer, each followed by a newline
func (ew *EnhancedWriter) WriteLines(lines []string) error {
	for _, line := range lines {
		_, err := ew.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return ew.Flush()
}

// ReadFirstNBytes reads exactly n bytes from r and returns them as a string.
// If fewer than n bytes are available, it returns an error.
func ReadFirstNBytes(r io.Reader, n int) (string, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// ReadUntilDelimiter reads from r until the delimiter is encountered
// and returns the content as a string (without the delimiter).
func ReadUntilDelimiter(r *bufio.Reader, delim byte) (string, error) {
	var buf bytes.Buffer
	for {
		b, err := r.ReadByte()
		if err != nil {
			return buf.String(), err
		}
		if b == delim {
			return buf.String(), nil
		}
		buf.WriteByte(b)
	}
}

// ScanForPattern is a split function for a Scanner that returns each text section
// that matches the pattern provided.
func ScanForPattern(pattern []byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// Look for pattern in the data
		if i := bytes.Index(data, pattern); i >= 0 {
			return i + len(pattern), data[:i], nil
		}

		// If we're at EOF, we have a final, non-terminated line
		if atEOF {
			return len(data), data, nil
		}

		// Request more data
		return 0, nil, nil
	}
}

// CreateScanner creates a new bufio.Scanner with a larger buffer for handling
// long lines or binary data.
func CreateScanner(r io.Reader, maxCapacity int) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	return scanner
}

// IsPrefix checks if the buffer begins with the specified prefix
func IsPrefix(reader *bufio.Reader, prefix []byte) (bool, error) {
	buf, err := reader.Peek(len(prefix))
	if err != nil {
		return false, err
	}
	return bytes.Equal(buf, prefix), nil
}
