// Package syscall provides additional functionality on top of the standard syscall library.
package syscall

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"time"
)

// Error wraps syscall errors with additional context
type Error struct {
	Op   string
	Err  error
	Path string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s %s: %v", e.Op, e.Path, e.Err)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// GetPID returns the process ID of the current process
func GetPID() int {
	return syscall.Getpid()
}

// GetPPID returns the parent process ID of the current process
func GetPPID() int {
	return syscall.Getppid()
}

// SetPriority sets the process priority
func SetPriority(pid, priority int) error {
	return syscall.Setpriority(syscall.PRIO_PROCESS, pid, priority)
}

// GetPriority gets the process priority
func GetPriority(pid int) (int, error) {
	return syscall.Getpriority(syscall.PRIO_PROCESS, pid)
}

// CreateLockFile creates a lock file and returns its file descriptor
func CreateLockFile(path string) (int, error) {
	fd, err := syscall.Open(path, syscall.O_CREAT|syscall.O_RDWR, 0666)
	if err != nil {
		return -1, &Error{"open", err, path}
	}

	err = syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		syscall.Close(fd)
		return -1, &Error{"flock", err, path}
	}

	return fd, nil
}

// ReleaseLockFile releases a lock file by file descriptor
func ReleaseLockFile(fd int, path string) error {
	if err := syscall.Flock(fd, syscall.LOCK_UN); err != nil {
		return &Error{"unlock", err, path}
	}

	if err := syscall.Close(fd); err != nil {
		return &Error{"close", err, path}
	}

	return nil
}

// GetSystemInfo returns basic system information
func GetSystemInfo() (string, error) {
	// Using runtime package as a cross-platform alternative to syscall.Uname
	return fmt.Sprintf("%s %s %s", runtime.GOOS, runtime.GOARCH, runtime.Version()), nil
}

// FileExists checks if a file exists using syscalls
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// Timeout executes a function with a timeout
func Timeout(timeout time.Duration, f func() error) error {
	done := make(chan error, 1)
	go func() {
		done <- f()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("operation timed out after %v", timeout)
	}
}

// GetOSType returns the current operating system type
func GetOSType() string {
	return runtime.GOOS
}

// GetCPUCount returns the number of CPUs
func GetCPUCount() int {
	return runtime.NumCPU()
}
