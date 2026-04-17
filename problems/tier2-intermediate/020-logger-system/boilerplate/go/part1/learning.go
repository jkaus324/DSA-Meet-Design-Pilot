package main

import "fmt"

// LogLevel defines severity levels for log messages.
type LogLevel int

const (
	DEBUG LogLevel = iota // 0
	INFO  LogLevel = iota // 1
	WARN  LogLevel = iota // 2
	ERROR LogLevel = iota // 3
	FATAL LogLevel = iota // 4
)

// LogEntry holds all data for a single log record.
type LogEntry struct {
	Level     LogLevel
	Message   string
	Timestamp string
}

// levelToString converts a LogLevel to its string name.
func levelToString(level LogLevel) string {
	names := map[LogLevel]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		FATAL: "FATAL",
	}
	return names[level]
}

// Logger is the Singleton logger.
type Logger struct {
	minLevel LogLevel
	history  []string
}

// instance holds the single Logger. It is created lazily in GetInstance.
var instance *Logger

func GetInstance() *Logger {
	// TODO: If instance == nil, create a new Logger with minLevel = INFO
	// TODO: Return instance
	return instance
}

func (l *Logger) Log(level LogLevel, message string) {
	// TODO: Check if level >= l.minLevel; if not, return immediately

	// TODO: Create a timestamp string (use "2024-01-15 10:30:00" for simplicity)
	timestamp := "2024-01-15 10:30:00"

	// TODO: Format as "[TIMESTAMP] [LEVEL] message"
	formatted := fmt.Sprintf("[%s] [%s] %s", timestamp, levelToString(level), message)

	// TODO: Append formatted to l.history
	// TODO: Print formatted to stdout
	_ = formatted
}

func (l *Logger) SetLevel(level LogLevel) {
	// TODO: Update l.minLevel
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
	// TODO: Reset l.history to an empty slice
	// Also reset instance so tests start clean
	instance = nil
}
