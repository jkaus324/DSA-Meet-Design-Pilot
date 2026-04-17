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

// Logger — extend with pluggable output formats:
//   - PlainText: [TIMESTAMP] [LEVEL] message
//   - JSON: {"timestamp":"...","level":"...","message":"..."}
//   - CSV: TIMESTAMP,LEVEL,message
//
// Think about:
//   - What abstraction lets you swap formatting logic at runtime?
//   - How do you ensure adding a new format requires zero changes to Logger?
//   - What is the default formatter if none is explicitly set?
//
// Entry points (must exist for tests):
//   GetInstance() *Logger
//   (*Logger).Log(level LogLevel, message string)
//   (*Logger).SetLevel(level LogLevel)
//   (*Logger).SetFormatter(formatter LogFormatter)
//   (*Logger).GetLogHistory() []string
//   (*Logger).ClearHistory()
//
// You also need:
//   LogFormatter interface
//   PlainTextFormatter, JsonFormatter, CsvFormatter structs

// LogFormatter is the interface for all log formatters.
type LogFormatter interface {
	Format(entry LogEntry) string
}

type PlainTextFormatter struct{}
type JsonFormatter struct{}
type CsvFormatter struct{}

func (f *PlainTextFormatter) Format(entry LogEntry) string { return "" }
func (f *JsonFormatter) Format(entry LogEntry) string      { return "" }
func (f *CsvFormatter) Format(entry LogEntry) string       { return "" }

type Logger struct {
}

func GetInstance() *Logger {
	return nil
}

func (l *Logger) Log(level LogLevel, message string) {
}

func (l *Logger) SetLevel(level LogLevel) {
}

func (l *Logger) SetFormatter(formatter LogFormatter) {
}

func (l *Logger) GetLevel() LogLevel {
	return INFO
}

func (l *Logger) GetLogHistory() []string {
	return nil
}

func (l *Logger) ClearHistory() {
}
