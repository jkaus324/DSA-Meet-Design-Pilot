#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
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

// ─── Concrete Formatters ─────────────────────────────────────────────────────

class PlainTextFormatter : public LogFormatter {
public:
    string format(const LogEntry& entry) override {
        // TODO: Return "[TIMESTAMP] [LEVEL] message"
        // Example: "[2024-01-15 10:30:00] [ERROR] Something failed"
        return "";
    }
};

class JsonFormatter : public LogFormatter {
public:
    string format(const LogEntry& entry) override {
        // TODO: Return JSON string
        // Example: {"timestamp":"2024-01-15 10:30:00","level":"ERROR","message":"Something failed"}
        return "";
    }
};

class CsvFormatter : public LogFormatter {
public:
    string format(const LogEntry& entry) override {
        // TODO: Return CSV string
        // Example: 2024-01-15 10:30:00,ERROR,Something failed
        return "";
    }
};

// ─── Logger (Singleton) ─────────────────────────────────────────────────────

class Logger {
    LogLevel minLevel;
    vector<string> history;
    LogFormatter* formatter;
    PlainTextFormatter defaultFormatter;

    Logger() : minLevel(LogLevel::INFO), formatter(&defaultFormatter) {}

public:
    static Logger& getInstance() {
        static Logger inst;
        return inst;
    }

    void log(LogLevel level, const string& message) {
        // TODO: Check if level >= minLevel, return if not

        // TODO: Create a LogEntry with the level, message, and a timestamp

        // TODO: Use formatter->format(entry) to get the formatted string

        // TODO: Push formatted string to history

        // TODO: Print formatted string to cout
    }

    void setLevel(LogLevel level) {
        // TODO: Update minLevel
    }

    void setFormatter(LogFormatter* f) {
        // TODO: Update the active formatter
        // If f is nullptr, revert to defaultFormatter
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
        // TODO: Clear history
    }

    Logger(const Logger&) = delete;
    Logger& operator=(const Logger&) = delete;
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Logger System Part 2 — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
