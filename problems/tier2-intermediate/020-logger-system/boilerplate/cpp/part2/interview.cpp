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
// Extend your Logger to support pluggable output formats:
//   - PlainText: [TIMESTAMP] [LEVEL] message
//   - JSON: {"timestamp":"...","level":"...","message":"..."}
//   - CSV: TIMESTAMP,LEVEL,message
//
// Think about:
//   - What abstraction lets you swap formatting logic at runtime?
//   - How do you ensure adding a new format (e.g., XML) requires
//     zero changes to the Logger?
//   - What is the default formatter if none is explicitly set?
//
// Entry points (must exist for tests):
//   Logger& Logger::getInstance();
//   void Logger::log(LogLevel level, const string& message);
//   void Logger::setLevel(LogLevel level);
//   void Logger::setFormatter(LogFormatter* formatter);
//   vector<string> Logger::getLogHistory() const;
//   void Logger::clearHistory();
//
// ─────────────────────────────────────────────────────────────────────────────


