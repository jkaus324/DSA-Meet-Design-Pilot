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

func levelToString(level LogLevel) string {
	names := map[LogLevel]string{
		DEBUG: "DEBUG", INFO: "INFO", WARN: "WARN", ERROR: "ERROR", FATAL: "FATAL",
	}
	return names[level]
}

// LogFormatter is the interface each formatter must implement.
// HINT: Each formatter converts a LogEntry into a string.
// The Logger delegates formatting to whichever formatter is active.
type LogFormatter interface {
	Format(entry LogEntry) string
}

// PlainTextFormatter formats as: [TIMESTAMP] [LEVEL] message
// HINT: Concatenate with brackets around timestamp and level
type PlainTextFormatter struct{}

func (f *PlainTextFormatter) Format(entry LogEntry) string {
	// TODO: Return "[TIMESTAMP] [LEVEL] message"
	return ""
}

// JsonFormatter formats as: {"timestamp":"...","level":"...","message":"..."}
// HINT: Build a JSON object string with three keys
type JsonFormatter struct{}

func (f *JsonFormatter) Format(entry LogEntry) string {
	// TODO: Return JSON string with keys: "timestamp", "level", "message"
	return ""
}

// CsvFormatter formats as: TIMESTAMP,LEVEL,message
// HINT: Comma-separated values, no brackets
type CsvFormatter struct{}

func (f *CsvFormatter) Format(entry LogEntry) string {
	// TODO: Return CSV string (no brackets)
	return ""
}

// Logger is the Singleton logger with a pluggable formatter.
// HINT: Add a LogFormatter field to Logger.
// HINT: The default formatter should be PlainTextFormatter.
// HINT: SetFormatter() swaps the active formatter at runtime.
type Logger struct {
	// HINT: same as Part 1 plus a formatter field
}

var instance *Logger

func GetInstance() *Logger {
	// HINT: Create once, return always
	return nil
}

func (l *Logger) Log(level LogLevel, message string) {
	// HINT: Check level >= minLevel, create LogEntry, delegate to formatter.Format(entry)
	// HINT: Append result to history, print to stdout
}

func (l *Logger) SetLevel(level LogLevel) {
	// TODO: Update minLevel
}

func (l *Logger) SetFormatter(formatter LogFormatter) {
	// HINT: If formatter is nil, revert to PlainTextFormatter
}

func (l *Logger) GetLevel() LogLevel {
	return INFO
}

func (l *Logger) GetLogHistory() []string {
	return nil
}

func (l *Logger) ClearHistory() {
	instance = nil
}
