// Package logs provides extended functionality on top of Go's standard log package
package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// Log levels
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Level names for output formatting
var levelNames = map[int]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// Logger extends the standard log package with levels and formatting
type Logger struct {
	level      int
	stdLogger  *log.Logger
	timeFormat string
	showCaller bool
}

// NewLogger creates a new Logger instance
func NewLogger(out io.Writer, prefix string, flag int, level int) *Logger {
	if out == nil {
		out = os.Stderr
	}
	return &Logger{
		level:      level,
		stdLogger:  log.New(out, prefix, flag),
		timeFormat: "2006-01-02 15:04:05",
		showCaller: true,
	}
}

// DefaultLogger returns a logger with sensible defaults
func DefaultLogger() *Logger {
	return NewLogger(os.Stderr, "", log.LstdFlags, INFO)
}

// SetLevel changes the current logging level
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// SetTimeFormat sets the format for timestamps
func (l *Logger) SetTimeFormat(format string) {
	l.timeFormat = format
}

// SetShowCaller enables/disables showing caller information
func (l *Logger) SetShowCaller(show bool) {
	l.showCaller = show
}

// formatMessage formats a log message with level, timestamp and caller info if enabled
func (l *Logger) formatMessage(level int, v ...interface{}) string {
	ts := time.Now().Format(l.timeFormat)
	levelName := levelNames[level]
	msg := fmt.Sprint(v...)
	parts := []string{ts, levelName, msg}

	if l.showCaller {
		_, file, line, ok := runtime.Caller(2) // Skip two frames to get the actual caller
		if ok {
			parts := strings.Split(file, "/")
			file = parts[len(parts)-1]
			caller := fmt.Sprintf("%s:%d", file, line)
			parts = append(parts, caller)
		}
	}

	return strings.Join(parts, " | ")
}

// log logs a message at the specified level
func (l *Logger) log(level int, v ...interface{}) {
	if level >= l.level {
		l.stdLogger.Println(l.formatMessage(level, v...))
	}
}

// logf logs a formatted message at the specified level
func (l *Logger) logf(level int, format string, v ...interface{}) {
	if level >= l.level {
		l.log(level, fmt.Sprintf(format, v...))
	}
}

// Debug logs a message at DEBUG level
func (l *Logger) Debug(v ...interface{}) {
	l.log(DEBUG, v...)
}

// Debugf logs a formatted message at DEBUG level
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logf(DEBUG, format, v...)
}

// Info logs a message at INFO level
func (l *Logger) Info(v ...interface{}) {
	l.log(INFO, v...)
}

// Infof logs a formatted message at INFO level
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logf(INFO, format, v...)
}

// Warn logs a message at WARN level
func (l *Logger) Warn(v ...interface{}) {
	l.log(WARN, v...)
}

// Warnf logs a formatted message at WARN level
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logf(WARN, format, v...)
}

// Error logs a message at ERROR level
func (l *Logger) Error(v ...interface{}) {
	l.log(ERROR, v...)
}

// Errorf logs a formatted message at ERROR level
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logf(ERROR, format, v...)
}

// Fatal logs a message at FATAL level and then exits
func (l *Logger) Fatal(v ...interface{}) {
	l.log(FATAL, v...)
	os.Exit(1)
}

// Fatalf logs a formatted message at FATAL level and then exits
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logf(FATAL, format, v...)
	os.Exit(1)
}

// Global logger instance for package-level functions
var defaultLogger = DefaultLogger()

// SetDefaultLogger changes the global default logger
func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}

// Global functions that use the default logger

// Debug logs a message at DEBUG level using the default logger
func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

// Debugf logs a formatted message at DEBUG level using the default logger
func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

// Info logs a message at INFO level using the default logger
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Infof logs a formatted message at INFO level using the default logger
func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

// Warn logs a message at WARN level using the default logger
func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}

// Warnf logs a formatted message at WARN level using the default logger
func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

// Error logs a message at ERROR level using the default logger
func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

// Errorf logs a formatted message at ERROR level using the default logger
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

// Fatal logs a message at FATAL level using the default logger and then exits
func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}

// Fatalf logs a formatted message at FATAL level using the default logger and then exits
func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}

// SetLevel sets the level of the default logger
func SetLevel(level int) {
	defaultLogger.SetLevel(level)
}

// SetTimeFormat sets the time format of the default logger
func SetTimeFormat(format string) {
	defaultLogger.SetTimeFormat(format)
}

// SetShowCaller enables/disables showing caller info in the default logger
func SetShowCaller(show bool) {
	defaultLogger.SetShowCaller(show)
}
