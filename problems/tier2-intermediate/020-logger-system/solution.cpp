#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>
using namespace std;

// ─── Data Model ─────────────────────────────────────────────────────────────

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

// ─── Formatter Interface (Strategy) ─────────────────────────────────────────

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

// ─── Destination Interface (Observer) ───────────────────────────────────────

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

    string getName() const override {
        return "console";
    }
};

class InMemoryDestination : public LogDestination {
    LogFormatter* formatter;
    string name;
public:
    vector<string> lines;

    InMemoryDestination(const string& name, LogFormatter* f)
        : name(name), formatter(f) {}

    void write(const LogEntry& entry) override {
        lines.push_back(formatter->format(entry));
    }

    string getName() const override {
        return name;
    }
};

// ─── Logger (Singleton) ────────────────────────────────────────────────────

class Logger {
    LogLevel minLevel;
    vector<string> history;
    vector<LogDestination*> destinations;
    LogFormatter* formatter;
    PlainTextFormatter defaultFormatter;

    Logger() : minLevel(LogLevel::INFO), formatter(&defaultFormatter) {}

public:
    static Logger& getInstance() {
        static Logger inst;
        return inst;
    }

    void log(LogLevel level, const string& message) {
        if (static_cast<int>(level) < static_cast<int>(minLevel)) return;

        LogEntry entry{level, message, "2024-01-15 10:30:00"};

        // Notify all destinations (fault isolation: catch exceptions)
        for (auto* dest : destinations) {
            try {
                dest->write(entry);
            } catch (...) {
                // Failure in one destination must not block others
            }
        }

        // Store in history using the active formatter
        string formatted = formatter->format(entry);
        history.push_back(formatted);
        cout << formatted << endl;
    }

    void setLevel(LogLevel level) {
        minLevel = level;
    }

    void setFormatter(LogFormatter* f) {
        if (f) {
            formatter = f;
        } else {
            formatter = &defaultFormatter;
        }
    }

    void addDestination(LogDestination* dest) {
        destinations.push_back(dest);
    }

    void removeDestination(LogDestination* dest) {
        destinations.erase(
            std::remove(destinations.begin(), destinations.end(), dest),
            destinations.end()
        );
    }

    LogLevel getLevel() const {
        return minLevel;
    }

    vector<string> getLogHistory() const {
        return history;
    }

    void clearHistory() {
        history.clear();
        destinations.clear();
    }

    Logger(const Logger&) = delete;
    Logger& operator=(const Logger&) = delete;
};

#ifndef RUNNING_TESTS
int main() {
    Logger& logger = Logger::getInstance();
    logger.setLevel(LogLevel::DEBUG);

    logger.log(LogLevel::DEBUG, "Application starting");
    logger.log(LogLevel::INFO, "Server listening on port 8080");
    logger.log(LogLevel::WARN, "Memory usage above 80%");
    logger.log(LogLevel::ERROR, "Failed to connect to database");
    logger.log(LogLevel::FATAL, "Unrecoverable error — shutting down");

    return 0;
}
#endif
