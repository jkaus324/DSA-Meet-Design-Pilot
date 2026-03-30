# Problem 020 — Logger System

**Tier:** 2 (Intermediate) | **Patterns:** Strategy, Observer, Factory, Singleton | **DSA:** HashMap, String Parsing, Queue
**Companies:** Amazon | **Time:** 45 minutes

---

## Problem Statement

Design and implement a configurable logging framework. The system must:

1. Support multiple log levels: `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`
2. Filter log messages based on a configured minimum severity level
3. Support pluggable output formats (plain text, JSON, CSV)
4. Support multiple simultaneous output destinations with independent formatters

---

## Before You Code

> Read this section carefully. This is where the design thinking happens.

**The naive approach:** A single `Logger` class with `if-else` chains for formatting and output:

```cpp
void log(Level level, string msg) {
    if (format == "json") { /* ... */ }
    else if (format == "csv") { /* ... */ }
    // ...
    if (output == "console") { cout << formatted; }
    else if (output == "file") { file << formatted; }
}
```

**Why this breaks:** Adding a new format (e.g., XML) or a new destination (e.g., network socket) requires modifying the `log()` method. Every change risks breaking existing behavior.

**The pattern approach:**
- **Singleton**: Ensures a single logger instance across the application (Part 1)
- **Strategy**: Output formatting becomes swappable — each format is a strategy (Part 2)
- **Observer**: Each output destination is an observer that receives log events independently (Part 3)

**The DSA angle:** Log levels map to integer severity values stored in a HashMap for O(1) lookup. Level filtering is a simple integer comparison once you map level names to ordinals.

---

## Data Structures

```cpp
enum class LogLevel { DEBUG = 0, INFO = 1, WARN = 2, ERROR = 3, FATAL = 4 };

struct LogEntry {
    LogLevel level;
    string message;
    string timestamp;  // "YYYY-MM-DD HH:MM:SS" format
};
```

---

## Part 1

**Base requirement — Singleton logger with level filtering**

Implement a `Logger` class that:

1. Is a **Singleton** — only one instance exists, obtained via `Logger::getInstance()`
2. Logs messages with `log(LogLevel level, const string& message)`
3. Filters messages: only logs entries at or above the configured minimum level
4. Output format: `[TIMESTAMP] [LEVEL] message` printed to console (stdout)
5. Default minimum level is `INFO`

| Level | Severity (higher = more severe) |
|-------|--------------------------------|
| DEBUG | 0 |
| INFO  | 1 |
| WARN  | 2 |
| ERROR | 3 |
| FATAL | 4 |

**Rules:**
- `log(DEBUG, "msg")` with level set to `INFO` produces **no output**
- `log(ERROR, "msg")` with level set to `INFO` **does** produce output
- Timestamp can be any consistent format; tests will check level filtering behavior, not exact timestamp

**Entry points (tests will call these):**
```cpp
Logger& Logger::getInstance();
void Logger::log(LogLevel level, const string& message);
void Logger::setLevel(LogLevel level);
LogLevel Logger::getLevel() const;
vector<string> Logger::getLogHistory() const;  // returns all logged messages for testing
```

**What to implement:**
```cpp
class Logger {
    static Logger* instance;
    LogLevel minLevel;
    vector<string> history;  // stores formatted log strings for testability

    Logger();  // private constructor
public:
    static Logger& getInstance();
    void log(LogLevel level, const string& message);
    void setLevel(LogLevel level);
    LogLevel getLevel() const;
    vector<string> getLogHistory() const;
    void clearHistory();

    // Prevent copying
    Logger(const Logger&) = delete;
    Logger& operator=(const Logger&) = delete;
};
```

---

## Part 2

**Extension 1 — Pluggable output formats via Strategy**

The ops team needs logs in different formats depending on the consumer:
- **PlainText**: `[2024-01-15 10:30:00] [ERROR] Something failed`
- **JSON**: `{"timestamp":"2024-01-15 10:30:00","level":"ERROR","message":"Something failed"}`
- **CSV**: `2024-01-15 10:30:00,ERROR,Something failed`

Adding a new format (e.g., XML) must require **zero changes** to the Logger class.

**New entry points:**
```cpp
void Logger::setFormatter(LogFormatter* formatter);
```

**Formatter interface:**
```cpp
class LogFormatter {
public:
    virtual string format(const LogEntry& entry) = 0;
    virtual ~LogFormatter() = default;
};

class PlainTextFormatter : public LogFormatter { ... };
class JsonFormatter     : public LogFormatter { ... };
class CsvFormatter      : public LogFormatter { ... };
```

**Design challenge:** The Logger now delegates formatting to a strategy object. The default formatter is `PlainTextFormatter`. Changing the formatter at runtime immediately affects all subsequent log calls.

---

## Part 3

**Extension 2 — Multiple simultaneous output destinations**

The system now needs to send logs to **multiple destinations simultaneously**:
- Console output
- File output (appends to a log file)
- Any future destination (e.g., network, database)

Each destination can have its **own formatter** — console might use PlainText while the file uses JSON.

**Rules:**
- Register/unregister destinations at runtime
- Failure in one destination must **not** block others — if file write fails, console output still works
- Each destination receives every log entry that passes the level filter

**New entry points:**
```cpp
void Logger::addDestination(LogDestination* dest);
void Logger::removeDestination(LogDestination* dest);
```

**Destination interface:**
```cpp
class LogDestination {
public:
    virtual void write(const LogEntry& entry) = 0;
    virtual string getName() const = 0;
    virtual ~LogDestination() = default;
};

class ConsoleDestination : public LogDestination {
    LogFormatter* formatter;
public:
    ConsoleDestination(LogFormatter* f);
    void write(const LogEntry& entry) override;
    string getName() const override;
};

class FileDestination : public LogDestination {
    LogFormatter* formatter;
    string filename;
public:
    FileDestination(const string& filename, LogFormatter* f);
    void write(const LogEntry& entry) override;
    string getName() const override;
};
```

**Design challenge:** The Logger becomes a subject (publisher) and destinations become observers (subscribers). How do you handle a destination that throws an exception during `write()`?

---

## Running Tests

```bash
./run-tests.sh 020-logger-system cpp
```
