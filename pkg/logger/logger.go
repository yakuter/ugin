package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger represents the application logger
type Logger struct {
	logger *logrus.Logger
	file   *os.File
	mu     sync.Mutex
}

// New creates a new logger instance
func New(level string) *Logger {
	l := &Logger{
		logger: logrus.New(),
	}

	// Set log level
	switch strings.ToLower(level) {
	case "debug":
		l.logger.Level = logrus.DebugLevel
	case "info":
		l.logger.Level = logrus.InfoLevel
	case "warn", "warning":
		l.logger.Level = logrus.WarnLevel
	case "error":
		l.logger.Level = logrus.ErrorLevel
	default:
		l.logger.Level = logrus.InfoLevel
	}

	// Set formatter
	l.logger.Formatter = &formatter{}
	l.logger.SetReportCaller(true)

	// Setup output writer
	if err := l.setupWriter(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup logger file: %v\n", err)
		l.logger.Out = os.Stdout
	}

	return l
}

func (l *Logger) setupWriter() error {
	file, err := os.OpenFile("ugin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	l.file = file
	l.logger.Out = io.MultiWriter(os.Stdout, file)
	return nil
}

// Close closes the log file
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	if l.logger.Level >= logrus.DebugLevel {
		entry := l.logger.WithFields(convertToFields(keysAndValues...))
		entry.Debug(msg)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	if l.logger.Level >= logrus.InfoLevel {
		entry := l.logger.WithFields(convertToFields(keysAndValues...))
		entry.Info(msg)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	if l.logger.Level >= logrus.WarnLevel {
		entry := l.logger.WithFields(convertToFields(keysAndValues...))
		entry.Warn(msg)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	if l.logger.Level >= logrus.ErrorLevel {
		entry := l.logger.WithFields(convertToFields(keysAndValues...))
		entry.Error(msg)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	entry := l.logger.WithFields(convertToFields(keysAndValues...))
	entry.Fatal(msg)
}

// convertToFields converts key-value pairs to logrus.Fields
func convertToFields(keysAndValues ...interface{}) logrus.Fields {
	fields := logrus.Fields{}
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key := fmt.Sprint(keysAndValues[i])
			fields[key] = keysAndValues[i+1]
		}
	}
	return fields
}

// formatter implements logrus.Formatter interface
type formatter struct {
	prefix string
}

// Format builds the log message
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var sb bytes.Buffer

	newLine := "\n"
	if runtime.GOOS == "windows" {
		newLine = "\r\n"
	}

	sb.WriteString(strings.ToUpper(entry.Level.String()))
	sb.WriteString(" ")
	sb.WriteString(entry.Time.Format(time.RFC3339))
	sb.WriteString(" ")

	// Add fields if present
	if len(entry.Data) > 0 {
		for key, value := range entry.Data {
			sb.WriteString(fmt.Sprintf("%s=%v ", key, value))
		}
	}

	sb.WriteString(entry.Message)
	sb.WriteString(newLine)

	return sb.Bytes(), nil
}

// Global helper functions for backward compatibility (but discouraged)
var defaultLogger = New("info")

// Debugf logs a debug message (deprecated: use Logger instance instead)
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debug(fmt.Sprintf(format, args...))
}

// Infof logs an info message (deprecated: use Logger instance instead)
func Infof(format string, args ...interface{}) {
	defaultLogger.Info(fmt.Sprintf(format, args...))
}

// Warnf logs a warning message (deprecated: use Logger instance instead)
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warn(fmt.Sprintf(format, args...))
}

// Errorf logs an error message (deprecated: use Logger instance instead)
func Errorf(format string, args ...interface{}) {
	defaultLogger.Error(fmt.Sprintf(format, args...))
}

// Fatalf logs a fatal message and exits (deprecated: use Logger instance instead)
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatal(fmt.Sprintf(format, args...))
}

