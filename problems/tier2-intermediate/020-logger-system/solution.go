package main

import "fmt"

// ─── Levels ──────────────────────────────────────────────────────────────────

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

func levelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return "UNKNOWN"
}

// ─── Log Entry ───────────────────────────────────────────────────────────────

type LogEntry struct {
	Level     LogLevel
	Message   string
	Timestamp string
}

// ─── Formatters (Strategy) ───────────────────────────────────────────────────

type LogFormatter interface {
	Format(entry LogEntry) string
}

type PlainTextFormatter struct{}

func (f *PlainTextFormatter) Format(entry LogEntry) string {
	return fmt.Sprintf("[%s] [%s] %s", entry.Timestamp, levelToString(entry.Level), entry.Message)
}

type JsonFormatter struct{}

func (f *JsonFormatter) Format(entry LogEntry) string {
	return fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s"}`,
		entry.Timestamp, levelToString(entry.Level), entry.Message)
}

type CsvFormatter struct{}

func (f *CsvFormatter) Format(entry LogEntry) string {
	return fmt.Sprintf("%s,%s,%s", entry.Timestamp, levelToString(entry.Level), entry.Message)
}

// ─── Destinations (Observer) ─────────────────────────────────────────────────

type LogDestination interface {
	Write(entry LogEntry)
	GetName() string
}

type ConsoleDestination struct {
	formatter LogFormatter
	name      string
}

func NewConsoleDestination(name string, f LogFormatter) *ConsoleDestination {
	return &ConsoleDestination{name: name, formatter: f}
}

func (d *ConsoleDestination) Write(entry LogEntry) {
	fmt.Println(d.formatter.Format(entry))
}

func (d *ConsoleDestination) GetName() string { return d.name }

type InMemoryDestination struct {
	formatter LogFormatter
	name      string
	Lines     []string
}

func NewInMemoryDestination(name string, f LogFormatter) *InMemoryDestination {
	return &InMemoryDestination{name: name, formatter: f}
}

func (d *InMemoryDestination) Write(entry LogEntry) {
	d.Lines = append(d.Lines, d.formatter.Format(entry))
}

func (d *InMemoryDestination) GetName() string { return d.name }

// ─── Logger (Singleton) ──────────────────────────────────────────────────────

type Logger struct {
	minLevel     LogLevel
	formatter    LogFormatter
	history      []string
	destinations []LogDestination
}

var instance *Logger

func GetInstance() *Logger {
	if instance == nil {
		instance = &Logger{
			minLevel:  INFO,
			formatter: &PlainTextFormatter{},
		}
	}
	return instance
}

func (l *Logger) Log(level LogLevel, message string) {
	if level < l.minLevel {
		return
	}
	entry := LogEntry{Level: level, Message: message, Timestamp: "2024-01-15 10:30:00"}

	for _, dest := range l.destinations {
		func(d LogDestination) {
			defer func() { _ = recover() }()
			d.Write(entry)
		}(dest)
	}

	if l.formatter == nil {
		l.formatter = &PlainTextFormatter{}
	}
	l.history = append(l.history, l.formatter.Format(entry))
}

func (l *Logger) SetLevel(level LogLevel) { l.minLevel = level }

func (l *Logger) SetFormatter(f LogFormatter) {
	if f == nil {
		l.formatter = &PlainTextFormatter{}
		return
	}
	l.formatter = f
}

func (l *Logger) AddDestination(dest LogDestination) {
	l.destinations = append(l.destinations, dest)
}

func (l *Logger) RemoveDestination(dest LogDestination) {
	out := l.destinations[:0]
	for _, d := range l.destinations {
		if d != dest {
			out = append(out, d)
		}
	}
	l.destinations = out
}

func (l *Logger) GetLevel() LogLevel       { return l.minLevel }
func (l *Logger) GetLogHistory() []string  { return l.history }

func (l *Logger) ClearHistory() {
	l.history = nil
	l.destinations = nil
	instance = nil
}
