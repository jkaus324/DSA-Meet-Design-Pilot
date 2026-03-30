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

// ─── Your Design Starts Here ─────────────────────────────────────────────────
//
// Extend your Logger to support multiple simultaneous output destinations.
// Each destination:
//   - Has its own formatter
//   - Receives every log entry that passes the level filter
//   - Can be registered/unregistered at runtime
//   - Must not block other destinations if it fails
//
// Think about:
//   - What pattern lets multiple objects receive the same event?
//   - How do you isolate failures between destinations?
//   - How do you identify a destination for removal?
//
// Entry points (must exist for tests):
//   Logger& Logger::getInstance();
//   void Logger::log(LogLevel level, const string& message);
//   void Logger::setLevel(LogLevel level);
//   void Logger::setFormatter(LogFormatter* f);  // For test compatibility (no-op in Part 3)
//   void Logger::addDestination(LogDestination* dest);
//   void Logger::removeDestination(LogDestination* dest);
//   vector<string> Logger::getLogHistory() const;
//   void Logger::clearHistory();
//
// You also need:
//   LogFormatter interface with PlainTextFormatter, JsonFormatter, CsvFormatter
//   LogDestination interface with ConsoleDestination, FileDestination
//
// ─────────────────────────────────────────────────────────────────────────────


