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

// ─── Logger (Singleton) ─────────────────────────────────────────────────────

class Logger {
    static Logger* instance;
    LogLevel minLevel;
    vector<string> history;

    Logger() : minLevel(LogLevel::INFO) {}

public:
    static Logger& getInstance() {
        // TODO: Create instance if it doesn't exist, then return it
        // Use: if (!instance) instance = new Logger();
        static Logger inst;
        return inst;
    }

    void log(LogLevel level, const string& message) {
        // TODO: Check if level >= minLevel using static_cast<int>
        // If level is below minLevel, return immediately (do nothing)

        // TODO: Create a timestamp string (use "2024-01-15 10:30:00" for simplicity)

        // TODO: Format as "[TIMESTAMP] [LEVEL] message"
        // Use levelToString() to convert LogLevel to string

        // TODO: Push formatted string to history vector

        // TODO: Print formatted string to cout
    }

    void setLevel(LogLevel level) {
        // TODO: Update minLevel
    }

    LogLevel getLevel() const {
        // TODO: Return current minLevel
        return LogLevel::INFO;
    }

    vector<string> getLogHistory() const {
        // TODO: Return the history vector
        return {};
    }

    void clearHistory() {
        // TODO: Clear the history vector
    }

    Logger(const Logger&) = delete;
    Logger& operator=(const Logger&) = delete;
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Logger System — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
