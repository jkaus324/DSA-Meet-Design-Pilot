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

// ─── Logger Class ────────────────────────────────────────────────────────────
// HINT: To ensure only one instance, make the constructor private and
// provide a static method that returns a reference to the single instance.

class Logger {
    // HINT: You need a static pointer to hold the single instance
    // HINT: Store the minimum log level and a history of formatted strings

    // HINT: Private constructor prevents external instantiation
    Logger();

public:
    // HINT: This static method creates the instance on first call, returns it on subsequent calls
    static Logger& getInstance();

    // HINT: Compare the incoming level against the minimum level before logging
    // Use static_cast<int> on LogLevel values for comparison
    void log(LogLevel level, const string& message);

    void setLevel(LogLevel level);
    LogLevel getLevel() const;
    vector<string> getLogHistory() const;
    void clearHistory();

    // HINT: Delete these to prevent copying a singleton
    Logger(const Logger&) = delete;
    Logger& operator=(const Logger&) = delete;
};

// ─── Test Entry Points ──────────────────────────────────────────────────────
// Your solution must provide Logger::getInstance() which returns
// the single Logger instance. All operations go through that instance.
// ─────────────────────────────────────────────────────────────────────────────

