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
type LogFormatter interface {
	Format(entry LogEntry) string
}

// TODO: Implement PlainTextFormatter, JsonFormatter, CsvFormatter (same as Part 2)

// LogDestination is the interface each output target must implement.
// HINT: Each destination knows how to write a log entry somewhere.
// HINT: Each destination owns its own formatter.
// HINT: GetName() is used to identify destinations for removal.
type LogDestination interface {
	Write(entry LogEntry)
	GetName() string
}

// HINT: ConsoleDestination uses its formatter and writes to stdout
// TODO: Implement ConsoleDestination with a LogFormatter field

// HINT: InMemoryDestination uses its formatter and appends to a Lines slice
// TODO: Implement InMemoryDestination with a LogFormatter field and public Lines []string

// Logger is the Singleton logger with multiple destinations (Observer pattern).
// HINT: Replace the single formatter with a []LogDestination
// HINT: log() iterates over all destinations and calls Write() on each
// HINT: Wrap each Write() in a recover() so one failure doesn't block others
// HINT: AddDestination appends, RemoveDestination removes by pointer match
// HINT: SetFormatter is a no-op stub for test compatibility with Part 2
type Logger struct {
	// HINT: minLevel, history, and a slice of LogDestination
}

var instance *Logger

func GetInstance() *Logger {
	// HINT: Create once, return always
	return nil
}

func (l *Logger) Log(level LogLevel, message string) {
	// HINT: Check level >= minLevel
	// HINT: Create LogEntry
	// HINT: For each destination, call Write wrapped in deferred recover
}

func (l *Logger) SetLevel(level LogLevel) {
	// TODO: Update minLevel
}

func (l *Logger) SetFormatter(f LogFormatter) {
	// No-op — kept for Part 2 test compatibility
}

func (l *Logger) AddDestination(dest LogDestination) {
	// TODO: Append dest to destinations
}

func (l *Logger) RemoveDestination(dest LogDestination) {
	// TODO: Remove dest from destinations (compare by pointer identity)
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
