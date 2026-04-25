package main

import "fmt"

// LogLevel defines severity levels for log messages.
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO  LogLevel = iota
	WARN  LogLevel = iota
	ERROR LogLevel = iota
	FATAL LogLevel = iota
)

// LogEntry holds all data for a single log record.
type LogEntry struct {
	Level     LogLevel
	Message   string
	Timestamp string
}

func levelToString(level LogLevel) string {
	names := map[LogLevel]string{
		DEBUG: "DEBUG", INFO: "INFO", WARN: "WARN", ERROR: "ERROR", FATAL: "FATAL",
	}
	return names[level]
}

// LogFormatter is the interface each formatter must implement.
type LogFormatter interface {
	Format(entry LogEntry) string
}

// PlainTextFormatter formats as: [TIMESTAMP] [LEVEL] message
type PlainTextFormatter struct{}

func (f *PlainTextFormatter) Format(entry LogEntry) string {
	// TODO: Return "[TIMESTAMP] [LEVEL] message"
	// Example: "[2024-01-15 10:30:00] [ERROR] Something failed"
	return fmt.Sprintf("[%s] [%s] %s", entry.Timestamp, levelToString(entry.Level), entry.Message)
}

// JsonFormatter formats as: {"timestamp":"...","level":"...","message":"..."}
type JsonFormatter struct{}

func (f *JsonFormatter) Format(entry LogEntry) string {
	// TODO: Return JSON string
	// Example: {"timestamp":"2024-01-15 10:30:00","level":"ERROR","message":"Something failed"}
	return fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s"}`,
		entry.Timestamp, levelToString(entry.Level), entry.Message)
}

// CsvFormatter formats as: TIMESTAMP,LEVEL,message
type CsvFormatter struct{}

func (f *CsvFormatter) Format(entry LogEntry) string {
	// TODO: Return CSV string
	// Example: 2024-01-15 10:30:00,ERROR,Something failed
	return fmt.Sprintf("%s,%s,%s", entry.Timestamp, levelToString(entry.Level), entry.Message)
}

// Logger is the Singleton logger with a pluggable formatter.
type Logger struct {
	minLevel        LogLevel
	history         []string
	formatter       LogFormatter
	defaultFormatter *PlainTextFormatter
}

var instance *Logger

func GetInstance() *Logger {
	// TODO: If instance == nil, create Logger with minLevel=INFO and PlainTextFormatter
	// TODO: Return instance
	return instance
}

func (l *Logger) Log(level LogLevel, message string) {
	// TODO: Check if level >= l.minLevel, return if not

	// TODO: Create a LogEntry with level, message, and timestamp "2024-01-15 10:30:00"

	// TODO: Use l.formatter.Format(entry) to get formatted string

	// TODO: Append formatted string to l.history

	// TODO: Print formatted string to stdout
}

func (l *Logger) SetLevel(level LogLevel) {
	// TODO: Update l.minLevel
}

func (l *Logger) SetFormatter(formatter LogFormatter) {
	// TODO: If formatter is nil, set l.formatter = l.defaultFormatter
	// TODO: Otherwise set l.formatter = formatter
}

func (l *Logger) GetLevel() LogLevel {
	// TODO: Return l.minLevel
	return INFO
}

func (l *Logger) GetLogHistory() []string {
	// TODO: Return l.history
	return nil
}

func (l *Logger) ClearHistory() {
	// TODO: Reset l.history to empty slice
	// Reset singleton so tests start fresh
	instance = nil
}
