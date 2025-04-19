// Package os extends the standard os library with additional functionality
package os

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileExists checks if a file or directory exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	if !FileExists(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// CopyDir recursively copies a directory tree
func CopyDir(src, dst string) error {
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(dst, sourceInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err = CopyDir(sourcePath, destPath); err != nil {
				return err
			}
		} else {
			if err = CopyFile(sourcePath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// HumanFileSize returns a human-readable file size
func HumanFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

// RemoveContents deletes all contents of a directory without removing the directory itself
func RemoveContents(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		} else {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}
	return nil
}

// ReadFileAsString reads file contents as a string
func ReadFileAsString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteStringToFile writes a string to a file
func WriteStringToFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// IsDirectory checks if a path is a directory
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}
