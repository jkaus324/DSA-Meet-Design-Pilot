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
// HINT: Each formatter converts a LogEntry into a string.
// The Logger delegates formatting to whichever formatter is active.

class LogFormatter {
public:
    virtual string format(const LogEntry& entry) = 0;
    virtual ~LogFormatter() = default;
};

// ─── Concrete Formatters ─────────────────────────────────────────────────────
// TODO: Implement format() for each:
//   - PlainText: [timestamp] [LEVEL] message
//   - JSON: {"timestamp":"...","level":"...","message":"..."}
//   - CSV: timestamp,LEVEL,message

class PlainTextFormatter : public LogFormatter {
    // HINT: Concatenate with brackets around timestamp and level
};

class JsonFormatter : public LogFormatter {
    // HINT: Build a JSON object string with three keys
};

class CsvFormatter : public LogFormatter {
    // HINT: Comma-separated values, no brackets
};

// ─── Logger (Singleton) ─────────────────────────────────────────────────────
// HINT: Add a LogFormatter* member to the Logger.
// HINT: The default formatter should be PlainTextFormatter.
// HINT: setFormatter() swaps the active formatter at runtime.
// HINT: log() now uses formatter->format(entry) instead of hardcoded formatting.

// class Logger { ... };

