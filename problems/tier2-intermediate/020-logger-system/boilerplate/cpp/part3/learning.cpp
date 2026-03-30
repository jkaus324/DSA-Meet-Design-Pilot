#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

enum class LogLevel { DEBUG = 0, INFO = 1, WARN = 2, ERROR = 3, FATAL = 4 };

static string levelToString(LogLevel level) {
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

// ─── Formatter Interface ─────────────────────────────────────────────────────

class LogFormatter {
public:
    virtual string format(const LogEntry& entry) = 0;
    virtual ~LogFormatter() = default;
};

class PlainTextFormatter : public LogFormatter {
public:
    string format(const LogEntry& entry) override {
        // TODO: Return "[TIMESTAMP] [LEVEL] message"
        return "";
    }
};

class JsonFormatter : public LogFormatter {
public:
    string format(const LogEntry& entry) override {
        // TODO: Return {"timestamp":"...","level":"...","message":"..."}
        return "";
    }
};

class CsvFormatter : public LogFormatter {
public:
    string format(const LogEntry& entry) override {
        // TODO: Return timestamp,LEVEL,message
        return "";
    }
};

// ─── Destination Interface ───────────────────────────────────────────────────

class LogDestination {
public:
    virtual void write(const LogEntry& entry) = 0;
    virtual string getName() const = 0;
    virtual ~LogDestination() = default;
};

// ─── Concrete Destinations ───────────────────────────────────────────────────

class ConsoleDestination : public LogDestination {
    LogFormatter* formatter;
public:
    ConsoleDestination(LogFormatter* f) : formatter(f) {}

    void write(const LogEntry& entry) override {
        // TODO: Use formatter->format(entry) and print to cout
    }

    string getName() const override {
        // TODO: Return "console"
        return "";
    }
};

class InMemoryDestination : public LogDestination {
    LogFormatter* formatter;
    string name;
public:
    vector<string> lines;  // public for test access

    InMemoryDestination(const string& name, LogFormatter* f)
        : name(name), formatter(f) {}

    void write(const LogEntry& entry) override {
        // TODO: Use formatter->format(entry) and push to lines vector
    }

    string getName() const override {
        // TODO: Return the name
        return "";
    }
};

// ─── Logger (Singleton) ─────────────────────────────────────────────────────

class Logger {
    LogLevel minLevel;
    vector<string> history;
    vector<LogDestination*> destinations;

    Logger() : minLevel(LogLevel::INFO) {}

public:
    static Logger& getInstance() {
        static Logger inst;
        return inst;
    }

    void log(LogLevel level, const string& message) {
        // TODO: Check if level >= minLevel, return if not

        // TODO: Create a LogEntry with level, message, and timestamp

        // TODO: For each destination, call write(entry) wrapped in try/catch
        // Failure in one destination must NOT prevent others from receiving the entry

        // TODO: Also store formatted entry in history (use PlainText format for history)
    }

    void setLevel(LogLevel level) {
        // TODO: Update minLevel
    }

    void setFormatter(LogFormatter* f) {
        // TODO: Not used in Part 3 (for test compatibility with Part 2)
    }

    void addDestination(LogDestination* dest) {
        // TODO: Add dest to the destinations vector
    }

    void removeDestination(LogDestination* dest) {
        // TODO: Remove dest from the destinations vector
        // Use std::remove and erase idiom
    }

    LogLevel getLevel() const {
        // TODO: Return minLevel
        return LogLevel::INFO;
    }

    vector<string> getLogHistory() const {
        // TODO: Return history
        return {};
    }

    void clearHistory() {
        // TODO: Clear history and destinations
    }

    Logger(const Logger&) = delete;
    Logger& operator=(const Logger&) = delete;
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Logger System Part 3 — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
