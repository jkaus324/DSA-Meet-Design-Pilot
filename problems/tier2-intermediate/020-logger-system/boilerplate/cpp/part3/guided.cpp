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

// TODO: Implement PlainTextFormatter, JsonFormatter, CsvFormatter (same as Part 2)

// ─── Destination Interface ───────────────────────────────────────────────────
// HINT: Each destination knows how to write a log entry somewhere.
// HINT: Each destination owns its own formatter.
// HINT: getName() is used to identify destinations for removal.

class LogDestination {
public:
    virtual void write(const LogEntry& entry) = 0;
    virtual string getName() const = 0;
    virtual ~LogDestination() = default;
};

// ─── Concrete Destinations ───────────────────────────────────────────────────
// HINT: ConsoleDestination uses its formatter, then writes to cout
// HINT: FileDestination uses its formatter, then appends to a vector<string> (simulating file)

// TODO: Implement ConsoleDestination with a LogFormatter* member
// TODO: Implement FileDestination with a LogFormatter* member and a vector<string> for stored lines

// ─── Logger (Singleton) ─────────────────────────────────────────────────────
// HINT: Replace the single formatter with a vector<LogDestination*>
// HINT: log() iterates over all destinations and calls write() on each
// HINT: Wrap each write() in try/catch so one failure doesn't block others
// HINT: addDestination pushes to the vector, removeDestination erases by pointer match
// HINT: Add setFormatter(LogFormatter*) as empty stub for test compatibility

// class Logger { ... };

