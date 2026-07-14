#include <iostream>
#include <memory>
#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>
#include <functional>
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

// ─── Ops simulator (used by spec-based tests) ──────────────────────────────
//
// The Logger writes to cout from log(); when running the spec runner, that
// would mix with PASS/FAIL lines. We provide free-function wrappers that
// suppress cout while invoking Logger methods, and expose query helpers that
// the simulator uses to verify behavior.

#include <sstream>

static void with_silent_cout(const std::function<void()>& body) {
    std::ostringstream oss;
    auto* old = std::cout.rdbuf(oss.rdbuf());
    body();
    std::cout.rdbuf(old);
}

static LogLevel level_from(const string& s) {
    if (s == "DEBUG") return LogLevel::DEBUG;
    if (s == "INFO")  return LogLevel::INFO;
    if (s == "WARN")  return LogLevel::WARN;
    if (s == "ERROR") return LogLevel::ERROR;
    return LogLevel::FATAL;
}

struct LogOp {
    string kind;
    string s1;
    string s2;
    int    i1;
};

class TestDest : public LogDestination {
    LogFormatter* fmt;
    string n;
public:
    vector<string> received;
    TestDest(const string& name, LogFormatter* f) : fmt(f), n(name) {}
    void write(const LogEntry& entry) override { received.push_back(fmt->format(entry)); }
    string getName() const override { return n; }
};

class FailingDest : public LogDestination {
public:
    void write(const LogEntry&) override { throw runtime_error("fail"); }
    string getName() const override { return "failing"; }
};

vector<string> logger_simulate(vector<LogOp> ops) {
    vector<string> out;
    Logger& logger = Logger::getInstance();
    PlainTextFormatter plain;
    JsonFormatter      json;
    CsvFormatter       csv;

    // Up to 4 destinations + 1 failing slot
    vector<unique_ptr<TestDest>> dests(8);
    unique_ptr<FailingDest> failing;

    auto fmt_for = [&](const string& s) -> LogFormatter* {
        if (s == "json") return &json;
        if (s == "csv")  return &csv;
        if (s == "plain") return &plain;
        return nullptr;
    };

    for (const auto& op : ops) {
        const string& k = op.kind;
        if (k == "reset") {
            with_silent_cout([&]{
                logger.clearHistory();
                logger.setLevel(LogLevel::INFO);
                logger.setFormatter(nullptr);
            });
            for (auto& d : dests) d.reset();
            failing.reset();
            out.push_back("ok");
        } else if (k == "set_level") {
            logger.setLevel(level_from(op.s1));
            out.push_back("ok");
        } else if (k == "set_formatter") {
            with_silent_cout([&]{ logger.setFormatter(fmt_for(op.s1)); });
            out.push_back("ok");
        } else if (k == "log") {
            with_silent_cout([&]{ logger.log(level_from(op.s1), op.s2); });
            out.push_back("ok");
        } else if (k == "history_size") {
            out.push_back(to_string((int)logger.getLogHistory().size()));
        } else if (k == "history_contains") {
            // i1 = index, s1 = needle
            const auto& h = logger.getLogHistory();
            if (op.i1 < 0 || op.i1 >= (int)h.size()) out.push_back("no");
            else out.push_back(h[op.i1].find(op.s1) != string::npos ? "yes" : "no");
        } else if (k == "add_dest") {
            // s1 = "test:<idx>:<formatter>" e.g. "test:0:plain"
            size_t a = op.s1.find(':');
            size_t b = op.s1.find(':', a+1);
            int idx = stoi(op.s1.substr(a+1, b-a-1));
            string fmtName = op.s1.substr(b+1);
            dests[idx].reset(new TestDest("d" + to_string(idx), fmt_for(fmtName)));
            with_silent_cout([&]{ logger.addDestination(dests[idx].get()); });
            out.push_back("ok");
        } else if (k == "rm_dest") {
            int idx = stoi(op.s1.substr(op.s1.find(':')+1));
            if (dests[idx]) {
                with_silent_cout([&]{ logger.removeDestination(dests[idx].get()); });
            }
            out.push_back("ok");
        } else if (k == "add_failing") {
            failing.reset(new FailingDest());
            with_silent_cout([&]{ logger.addDestination(failing.get()); });
            out.push_back("ok");
        } else if (k == "rm_failing") {
            if (failing) with_silent_cout([&]{ logger.removeDestination(failing.get()); });
            out.push_back("ok");
        } else if (k == "dest_size") {
            int idx = stoi(op.s1.substr(op.s1.find(':')+1));
            out.push_back(dests[idx] ? to_string((int)dests[idx]->received.size()) : "0");
        } else if (k == "dest_contains") {
            // s1 = "test:<idx>" + i1 = entry index, s2 = needle
            int idx = stoi(op.s1.substr(op.s1.find(':')+1));
            if (!dests[idx]) { out.push_back("no"); continue; }
            const auto& v = dests[idx]->received;
            if (op.i1 < 0 || op.i1 >= (int)v.size()) out.push_back("no");
            else out.push_back(v[op.i1].find(op.s2) != string::npos ? "yes" : "no");
        } else if (k == "fmt_plain") {
            // s1 = level, s2 = message; assume timestamp "T"
            LogEntry e{level_from(op.s1), op.s2, "T"};
            out.push_back(plain.format(e));
        } else if (k == "fmt_json") {
            LogEntry e{level_from(op.s1), op.s2, "T"};
            out.push_back(json.format(e));
        } else if (k == "fmt_csv") {
            LogEntry e{level_from(op.s1), op.s2, "T"};
            out.push_back(csv.format(e));
        } else if (k == "singleton_check") {
            Logger& a = Logger::getInstance();
            Logger& b = Logger::getInstance();
            out.push_back(&a == &b ? "yes" : "no");
        } else {
            out.push_back("unknown:" + k);
        }
    }
    return out;
}

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
