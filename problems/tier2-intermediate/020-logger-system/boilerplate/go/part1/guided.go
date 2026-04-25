package main

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

// Logger is a Singleton logger with level filtering.
// HINT: To ensure only one instance, keep a package-level variable
// and return it from GetInstance on every call.
type Logger struct {
	// HINT: Store the minimum log level and a history of formatted strings
}

// HINT: Declare a package-level variable to hold the single instance
// var instance *Logger

func GetInstance() *Logger {
	// HINT: Create the instance only once (check if nil), then return it
	return nil
}

func (l *Logger) Log(level LogLevel, message string) {
	// HINT: Compare level >= l.minLevel using integer comparison
	// HINT: If below min level, return immediately
	// HINT: Create a timestamp string (e.g. "2024-01-15 10:30:00")
	// HINT: Format as "[TIMESTAMP] [LEVEL] message"
	// HINT: Append formatted string to history
}

func (l *Logger) SetLevel(level LogLevel) {
	// HINT: Update the minimum level
}

func (l *Logger) GetLevel() LogLevel {
	// HINT: Return the current minimum level
	return INFO
}

func (l *Logger) GetLogHistory() []string {
	// HINT: Return the history slice
	return nil
}

func (l *Logger) ClearHistory() {
	// HINT: Reset the history slice to empty
}
