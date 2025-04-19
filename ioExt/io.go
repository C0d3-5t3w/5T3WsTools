package ioExt

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// ReadFileString reads the entire contents of a file and returns it as a string.
// If the file cannot be read, it returns an error.
func ReadFileString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFileWithDirs writes data to a file, creating all necessary directories.
// If the file already exists, it will be truncated.
func WriteFileWithDirs(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, perm)
}

// CopyWithProgress copies data from src to dst, reporting progress periodically.
// It returns the number of bytes copied and the first error encountered, if any.
func CopyWithProgress(dst io.Writer, src io.Reader, progressFn func(written int64)) (int64, error) {
	buf := make([]byte, 32*1024)
	var written int64

	for {
		nr, err := src.Read(buf)
		if nr > 0 {
			nw, werr := dst.Write(buf[:nr])
			if nw > 0 {
				written += int64(nw)
				if progressFn != nil {
					progressFn(written)
				}
			}
			if werr != nil {
				return written, werr
			}
			if nw != nr {
				return written, io.ErrShortWrite
			}
		}
		if err != nil {
			if err == io.EOF {
				return written, nil
			}
			return written, err
		}
	}
}

// SafeClose attempts to close the closer and returns the error if any.
// It's useful in defer statements where you might want to track the close error.
func SafeClose(c io.Closer) error {
	if c == nil {
		return nil
	}
	return c.Close()
}

// FileExists checks if a file exists and is not a directory.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// AppendToFile appends data to a file, creating the file if it doesn't exist.
func AppendToFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer SafeClose(f)

	_, err = f.Write(data)
	return err
}

// WalkFiles recursively walks through a directory and calls the walkFn for each file.
// Directories are skipped.
func WalkFiles(root string, walkFn func(path string, info fs.FileInfo) error) error {
	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return walkFn(path, info)
		}
		return nil
	})
}

// MultiReadCloser returns an io.ReadCloser that's the logical concatenation of the provided
// input readers. They're read sequentially, and Close() closes all readers.
func MultiReadCloser(readers ...io.ReadCloser) io.ReadCloser {
	return &multiReadCloser{
		readers: readers,
		current: 0,
	}
}

type multiReadCloser struct {
	readers []io.ReadCloser
	current int
}

func (m *multiReadCloser) Read(p []byte) (n int, err error) {
	if m.current >= len(m.readers) {
		return 0, io.EOF
	}

	n, err = m.readers[m.current].Read(p)
	if err == io.EOF {
		m.current++
		if m.current < len(m.readers) {
			return m.Read(p)
		}
	}
	return
}

func (m *multiReadCloser) Close() error {
	var firstErr error
	for _, r := range m.readers {
		if err := r.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
