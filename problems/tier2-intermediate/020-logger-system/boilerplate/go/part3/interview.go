package main

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

// Logger — extend with multiple simultaneous output destinations.
// Each destination:
//   - Has its own formatter
//   - Receives every log entry that passes the level filter
//   - Can be registered/unregistered at runtime
//   - Must not block other destinations if it fails
//
// Think about:
//   - What pattern lets multiple objects receive the same event?
//   - How do you isolate failures between destinations?
//   - How do you identify a destination for removal?
//
// Entry points (must exist for tests):
//   GetInstance() *Logger
//   (*Logger).Log(level LogLevel, message string)
//   (*Logger).SetLevel(level LogLevel)
//   (*Logger).SetFormatter(f LogFormatter)   // no-op stub for compatibility
//   (*Logger).AddDestination(dest LogDestination)
//   (*Logger).RemoveDestination(dest LogDestination)
//   (*Logger).GetLogHistory() []string
//   (*Logger).ClearHistory()
//
// You also need:
//   LogFormatter interface with PlainTextFormatter, JsonFormatter, CsvFormatter
//   LogDestination interface

// LogFormatter formats a LogEntry into a string.
type LogFormatter interface {
	Format(entry LogEntry) string
}

type PlainTextFormatter struct{}
type JsonFormatter struct{}
type CsvFormatter struct{}

func (f *PlainTextFormatter) Format(entry LogEntry) string { return "" }
func (f *JsonFormatter) Format(entry LogEntry) string      { return "" }
func (f *CsvFormatter) Format(entry LogEntry) string       { return "" }

// LogDestination receives and stores/outputs log entries.
type LogDestination interface {
	Write(entry LogEntry)
	GetName() string
}

type Logger struct {
}

func GetInstance() *Logger {
	return nil
}

func (l *Logger) Log(level LogLevel, message string) {
}

func (l *Logger) SetLevel(level LogLevel) {
}

func (l *Logger) SetFormatter(f LogFormatter) {
}

func (l *Logger) AddDestination(dest LogDestination) {
}

func (l *Logger) RemoveDestination(dest LogDestination) {
}

func (l *Logger) GetLevel() LogLevel {
	return INFO
}

func (l *Logger) GetLogHistory() []string {
	return nil
}

func (l *Logger) ClearHistory() {
}
