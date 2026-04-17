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

// Logger — design and implement a Singleton Logger that:
//   1. Has only one instance (GetInstance returns the same pointer every call)
//   2. Supports log levels: DEBUG, INFO, WARN, ERROR, FATAL
//   3. Filters messages below the configured minimum level
//   4. Stores formatted log entries in a history slice for testability
//   5. Format: [TIMESTAMP] [LEVEL] message
//
// Think about:
//   - How do you prevent multiple instances in Go?
//   - How do you compare log levels efficiently?
//   - What happens when SetLevel(WARN) is called and then Log(INFO, ...) is called?
//
// Entry points (must exist for tests):
//   GetInstance() *Logger
//   (*Logger).Log(level LogLevel, message string)
//   (*Logger).SetLevel(level LogLevel)
//   (*Logger).GetLevel() LogLevel
//   (*Logger).GetLogHistory() []string
//   (*Logger).ClearHistory()

type Logger struct {
}

func GetInstance() *Logger {
	return nil
}

func (l *Logger) Log(level LogLevel, message string) {
}

func (l *Logger) SetLevel(level LogLevel) {
}

func (l *Logger) GetLevel() LogLevel {
	return INFO
}

func (l *Logger) GetLogHistory() []string {
	return nil
}

func (l *Logger) ClearHistory() {
}
