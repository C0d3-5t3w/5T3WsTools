// Package path provides extended functionality for path manipulation
// beyond what's available in the standard path and path/filepath packages.
package pathExt

import (
	"os"
	"path/filepath"
	"strings"
)

// IsDir checks if a path exists and is a directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile checks if a path exists and is a regular file
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// EnsureDir ensures that a directory exists, creating it if necessary
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// SplitAll splits a path into all of its components
func SplitAll(path string) []string {
	dir, file := filepath.Split(path)
	if dir == "" && file == "" {
		return []string{}
	}
	if file == "" {
		dir = filepath.Clean(dir)
		parts := SplitAll(filepath.Dir(dir))
		return append(parts, filepath.Base(dir))
	}
	parts := SplitAll(dir)
	return append(parts, file)
}

// ParentDir returns the parent directory of a path
func ParentDir(path string) string {
	return filepath.Dir(path)
}

// Filename returns just the filename without extension
func Filename(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return base[:len(base)-len(ext)]
}

// JoinWithExt joins path elements and adds extension to the last element
func JoinWithExt(ext string, elem ...string) string {
	if len(elem) == 0 {
		return ""
	}

	result := filepath.Join(elem...)
	if ext == "" {
		return result
	}

	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	return result + ext
}

// GetRelativePath returns a relative path that is relative to the base path
func GetRelativePath(basePath, targetPath string) (string, error) {
	return filepath.Rel(basePath, targetPath)
}

// IsSubPath checks if a path is a subpath of another path
func IsSubPath(basePath, targetPath string) (bool, error) {
	if basePath == targetPath {
		return true, nil
	}

	rel, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		return false, err
	}

	return !strings.HasPrefix(rel, ".."), nil
}
