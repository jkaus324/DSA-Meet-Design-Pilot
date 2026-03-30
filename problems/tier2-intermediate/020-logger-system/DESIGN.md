# Design Walkthrough — Logger System

> This file is the answer guide. Only read after you've attempted the problem.

---

## The Core Design Decisions

This problem combines **four patterns** that each solve a specific concern:

```
Logger (Singleton)
    ├── LogFormatter* (Strategy — swappable formatting)
    │       ├── PlainTextFormatter
    │       ├── JsonFormatter
    │       └── CsvFormatter
    └── vector<LogDestination*> (Observer — multiple destinations)
            ├── ConsoleDestination (has its own LogFormatter*)
            └── FileDestination    (has its own LogFormatter*)
```

**Why Singleton?** A logging framework is a cross-cutting concern. Multiple instances would create confusion about configuration (which level? which formatter?). A single access point via `getInstance()` ensures consistency.

**Why Strategy for formatters?** The formatting algorithm varies independently from the Logger. PlainText, JSON, CSV are all valid representations of the same data. Encapsulating each as a separate class means adding XML requires only a new class — no existing code changes.

**Why Observer for destinations?** Each destination independently decides how to handle a log entry. The Logger doesn't know (or care) whether it's writing to console, file, or a network socket. Destinations register themselves and receive notifications.

---

## Reference Implementation

### Part 1 — Singleton Logger with Level Filtering

```cpp
#include <vector>
#include <string>
#include <iostream>
#include <unordered_map>
using namespace std;

enum class LogLevel { DEBUG = 0, INFO = 1, WARN = 2, ERROR = 3, FATAL = 4 };

string levelToString(LogLevel level) {
    static unordered_map<int, string> names = {
        {0, "DEBUG"}, {1, "INFO"}, {2, "WARN"}, {3, "ERROR"}, {4, "FATAL"}
    };
    return names[static_cast<int>(level)];
}

struct LogEntry {
    LogLevel level;
    string message;
    string timestamp;
};

class Logger {
    static Logger* instance;
    LogLevel minLevel;
    vector<string> history;

    Logger() : minLevel(LogLevel::INFO) {}
public:
    static Logger& getInstance() {
        if (!instance) instance = new Logger();
        return *instance;
    }

    void log(LogLevel level, const string& message) {
        if (static_cast<int>(level) < static_cast<int>(minLevel)) return;

        string timestamp = "2024-01-15 10:30:00";  // simplified
        string formatted = "[" + timestamp + "] [" + levelToString(level) + "] " + message;
        history.push_back(formatted);
        cout << formatted << endl;
    }

    void setLevel(LogLevel level) { minLevel = level; }
    LogLevel getLevel() const { return minLevel; }
    vector<string> getLogHistory() const { return history; }
    void clearHistory() { history.clear(); }

    Logger(const Logger&) = delete;
    Logger& operator=(const Logger&) = delete;
};

Logger* Logger::instance = nullptr;
```

### Part 2 — Strategy Formatters

```cpp
class LogFormatter {
public:
    virtual string format(const LogEntry& entry) = 0;
    virtual ~LogFormatter() = default;
};

class PlainTextFormatter : public LogFormatter {
public:
    string format(const LogEntry& entry) override {
        return "[" + entry.timestamp + "] [" + levelToString(entry.level) + "] " + entry.message;
    }
};

class JsonFormatter : public LogFormatter {
public:
    string format(const LogEntry& entry) override {
        return "{\"timestamp\":\"" + entry.timestamp +
               "\",\"level\":\"" + levelToString(entry.level) +
               "\",\"message\":\"" + entry.message + "\"}";
    }
};

class CsvFormatter : public LogFormatter {
public:
    string format(const LogEntry& entry) override {
        return entry.timestamp + "," + levelToString(entry.level) + "," + entry.message;
    }
};
```

The Logger stores a `LogFormatter*` and delegates: `string formatted = formatter->format(entry);`. The default is `PlainTextFormatter`.

### Part 3 — Observer Destinations

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
    ConsoleDestination(LogFormatter* f) : formatter(f) {}
    void write(const LogEntry& entry) override {
        cout << formatter->format(entry) << endl;
    }
    string getName() const override { return "console"; }
};
```

The Logger's `log()` method becomes:

```cpp
void log(LogLevel level, const string& message) {
    if (static_cast<int>(level) < static_cast<int>(minLevel)) return;

    LogEntry entry{level, message, getCurrentTimestamp()};
    for (auto* dest : destinations) {
        try {
            dest->write(entry);
        } catch (...) {
            // Failure in one destination must not block others
        }
    }
}
```

---

## Key Structural Decisions

### Why not just use `cout` directly?

Direct `cout` usage couples the Logger to a single output. When the interviewer asks "now also log to a file," you'd need to modify `log()`. With the Observer pattern, you just register a new destination.

### Why does each destination have its own formatter?

In production systems, console output is human-readable (PlainText) while file output is machine-parseable (JSON). This flexibility comes naturally when each destination owns its formatter rather than the Logger owning a single global formatter.

### Why try/catch around each destination?

A failing file write (disk full, permissions) should never crash the application or prevent console logging. The try/catch per destination provides fault isolation — a critical production concern.

---

## Pattern Interaction Diagram

```
User calls Logger::getInstance().log(ERROR, "disk full")
    │
    ├─ Level check: ERROR (3) >= INFO (1)? Yes → proceed
    │
    ├─ Create LogEntry { ERROR, "disk full", "2024-01-15 10:30:00" }
    │
    └─ For each registered LogDestination:
        ├─ ConsoleDestination
        │   └─ PlainTextFormatter::format(entry) → "[...] [ERROR] disk full"
        │   └─ cout << formatted
        │
        └─ FileDestination
            └─ JsonFormatter::format(entry) → {"level":"ERROR",...}
            └─ file << formatted
```

---

## Interview Tips

1. **Start with Singleton.** It's the simplest entry point and shows you understand the access pattern for a logger.
2. **Introduce Strategy when asked about formats.** Don't pre-engineer it — wait for the extension. Interviewers want to see you refactor.
3. **Name the patterns aloud.** Say "I'll use Strategy for the formatter" and "Observer for destinations."
4. **Discuss trade-offs.** Singleton makes testing harder (global state). Mention that dependency injection is the production alternative.
5. **The HashMap for levels** is a minor DSA point but shows you think about efficient lookups vs. switch statements.

---

## Common Interview Follow-ups

- *"How would you make this thread-safe?"* → Mutex around the `log()` method; or use a lock-free queue where producers enqueue and a single consumer writes.
- *"How would you add log rotation?"* → A `RotatingFileDestination` that checks file size before writing and rotates when a threshold is exceeded.
- *"How would you support structured logging (key-value pairs)?"* → Extend `LogEntry` with a `map<string, string> metadata` field; formatters decide how to serialize it.
- *"What if you need async logging?"* → Producer-consumer pattern: `log()` enqueues to a thread-safe queue; a background thread dequeues and writes to destinations.
