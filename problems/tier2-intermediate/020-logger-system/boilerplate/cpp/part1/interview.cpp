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

// ─── Your Design Starts Here ─────────────────────────────────────────────────
//
// Design and implement a Logger that:
//   1. Is a Singleton — only one instance exists
//   2. Supports log levels: DEBUG, INFO, WARN, ERROR, FATAL
//   3. Filters messages based on a configured minimum level
//   4. Stores formatted log entries in a history for testability
//   5. Format: [TIMESTAMP] [LEVEL] message
//
// Think about:
//   - How do you prevent multiple instances?
//   - How do you compare log levels efficiently?
//   - What happens when setLevel(WARN) is called and then log(INFO, ...) is called?
//
// Entry points (must exist for tests):
//   Logger& Logger::getInstance();
//   void Logger::log(LogLevel level, const string& message);
//   void Logger::setLevel(LogLevel level);
//   LogLevel Logger::getLevel() const;
//   vector<string> Logger::getLogHistory() const;
//   void Logger::clearHistory();
//
// ─────────────────────────────────────────────────────────────────────────────


