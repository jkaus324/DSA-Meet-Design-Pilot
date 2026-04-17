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
	return fmt.Sprintf("[%s] [%s] %s", entry.Timestamp, levelToString(entry.Level), entry.Message)
}

// JsonFormatter formats as: {"timestamp":"...","level":"...","message":"..."}
type JsonFormatter struct{}

func (f *JsonFormatter) Format(entry LogEntry) string {
	// TODO: Return JSON string
	return fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s"}`,
		entry.Timestamp, levelToString(entry.Level), entry.Message)
}

// CsvFormatter formats as: TIMESTAMP,LEVEL,message
type CsvFormatter struct{}

func (f *CsvFormatter) Format(entry LogEntry) string {
	// TODO: Return CSV string
	return fmt.Sprintf("%s,%s,%s", entry.Timestamp, levelToString(entry.Level), entry.Message)
}

// LogDestination is the interface each output target must implement.
type LogDestination interface {
	Write(entry LogEntry)
	GetName() string
}

// ConsoleDestination writes formatted log entries to stdout.
type ConsoleDestination struct {
	formatter LogFormatter
	name      string
}

func NewConsoleDestination(name string, f LogFormatter) *ConsoleDestination {
	return &ConsoleDestination{name: name, formatter: f}
}

func (d *ConsoleDestination) Write(entry LogEntry) {
	// TODO: Use d.formatter.Format(entry) and print to stdout
	fmt.Println(d.formatter.Format(entry))
}

func (d *ConsoleDestination) GetName() string {
	// TODO: Return d.name
	return d.name
}

// InMemoryDestination stores formatted log entries in a Lines slice.
type InMemoryDestination struct {
	formatter LogFormatter
	name      string
	Lines     []string // public for test access
}

func NewInMemoryDestination(name string, f LogFormatter) *InMemoryDestination {
	return &InMemoryDestination{name: name, formatter: f}
}

func (d *InMemoryDestination) Write(entry LogEntry) {
	// TODO: Use d.formatter.Format(entry) and append to d.Lines
}

func (d *InMemoryDestination) GetName() string {
	// TODO: Return d.name
	return d.name
}

// Logger is the Singleton logger that broadcasts to multiple destinations.
type Logger struct {
	minLevel     LogLevel
	history      []string
	destinations []LogDestination
}

var instance *Logger

func GetInstance() *Logger {
	// TODO: If instance == nil, create Logger with minLevel=INFO
	// TODO: Return instance
	return instance
}

func (l *Logger) Log(level LogLevel, message string) {
	// TODO: Check if level >= l.minLevel, return if not

	// TODO: Create LogEntry with level, message, timestamp "2024-01-15 10:30:00"

	// TODO: For each destination, call destination.Write(entry)
	//   Wrap each Write in a goroutine-safe recover to isolate failures:
	//   func() { defer func() { recover() }(); dest.Write(entry) }()

	// TODO: Also store plain-text formatted entry in l.history
	plain := &PlainTextFormatter{}
	_ = plain
}

func (l *Logger) SetLevel(level LogLevel) {
	// TODO: Update l.minLevel
}

func (l *Logger) SetFormatter(f LogFormatter) {
	// No-op — kept for Part 2 test compatibility
}

func (l *Logger) AddDestination(dest LogDestination) {
	// TODO: Append dest to l.destinations
}

func (l *Logger) RemoveDestination(dest LogDestination) {
	// TODO: Remove dest from l.destinations by pointer identity
	// Iterate and rebuild the slice without the matching entry
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
	// TODO: Reset l.history and l.destinations to empty
	// Reset singleton so tests start fresh
	instance = nil
}
