// Part 3 Tests — Logger System: Multiple Destinations (Observer)
// Tests registering/unregistering destinations with independent formatters

#include "solution.cpp"
#include <cassert>
#include <iostream>
#include <stdexcept>
using namespace std;

// A test destination that stores formatted log lines in memory
class TestDestination : public LogDestination {
    LogFormatter* formatter;
    string destName;
public:
    vector<string> received;

    TestDestination(const string& name, LogFormatter* f)
        : destName(name), formatter(f) {}

    void write(const LogEntry& entry) override {
        received.push_back(formatter->format(entry));
    }

    string getName() const override { return destName; }
};

// A destination that always throws (to test fault isolation)
class FailingDestination : public LogDestination {
public:
    void write(const LogEntry& entry) override {
        throw runtime_error("destination failed");
    }
    string getName() const override { return "failing"; }
};

int part3_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: Register a single destination and log
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::DEBUG);

        PlainTextFormatter plainFmt;
        TestDestination dest1("test1", &plainFmt);
        logger.addDestination(&dest1);

        logger.log(LogLevel::INFO, "hello dest");
        assert(dest1.received.size() == 1);
        assert(dest1.received[0].find("INFO") != string::npos);
        assert(dest1.received[0].find("hello dest") != string::npos);

        logger.removeDestination(&dest1);
        cout << "PASS test_single_destination" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_single_destination" << endl;
        failed++;
    }

    // Test 2: Multiple destinations receive same log entry
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::DEBUG);

        PlainTextFormatter plainFmt;
        JsonFormatter jsonFmt;
        TestDestination dest1("plain-dest", &plainFmt);
        TestDestination dest2("json-dest", &jsonFmt);
        logger.addDestination(&dest1);
        logger.addDestination(&dest2);

        logger.log(LogLevel::ERROR, "multi dest");

        assert(dest1.received.size() == 1);
        assert(dest2.received.size() == 1);
        // dest1 uses PlainText
        assert(dest1.received[0].find("[ERROR]") != string::npos);
        // dest2 uses JSON
        assert(dest2.received[0].find("\"level\"") != string::npos);

        logger.removeDestination(&dest1);
        logger.removeDestination(&dest2);
        cout << "PASS test_multiple_destinations" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_multiple_destinations" << endl;
        failed++;
    }

    // Test 3: Each destination has independent formatter
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::DEBUG);

        PlainTextFormatter plainFmt;
        CsvFormatter csvFmt;
        TestDestination destPlain("plain", &plainFmt);
        TestDestination destCsv("csv", &csvFmt);
        logger.addDestination(&destPlain);
        logger.addDestination(&destCsv);

        logger.log(LogLevel::WARN, "format test");

        // Plain has brackets, CSV does not
        assert(destPlain.received[0].find("[WARN]") != string::npos);
        assert(destCsv.received[0].find("[") == string::npos);
        assert(destCsv.received[0].find(",") != string::npos);

        logger.removeDestination(&destPlain);
        logger.removeDestination(&destCsv);
        cout << "PASS test_independent_formatters" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_independent_formatters" << endl;
        failed++;
    }

    // Test 4: Failing destination does not block others
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::DEBUG);

        PlainTextFormatter plainFmt;
        FailingDestination failDest;
        TestDestination goodDest("good", &plainFmt);
        logger.addDestination(&failDest);
        logger.addDestination(&goodDest);

        logger.log(LogLevel::ERROR, "fault isolation");

        // goodDest should still receive the message despite failDest throwing
        assert(goodDest.received.size() == 1);
        assert(goodDest.received[0].find("fault isolation") != string::npos);

        logger.removeDestination(&failDest);
        logger.removeDestination(&goodDest);
        cout << "PASS test_fault_isolation" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_fault_isolation" << endl;
        failed++;
    }

    // Test 5: Remove destination stops it from receiving logs
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::DEBUG);

        PlainTextFormatter plainFmt;
        TestDestination dest("removable", &plainFmt);
        logger.addDestination(&dest);

        logger.log(LogLevel::INFO, "before remove");
        assert(dest.received.size() == 1);

        logger.removeDestination(&dest);
        logger.log(LogLevel::INFO, "after remove");
        assert(dest.received.size() == 1);  // still 1, not 2

        cout << "PASS test_remove_destination" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_remove_destination" << endl;
        failed++;
    }

    // Test 6: Level filtering applies before destinations receive
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::ERROR);

        PlainTextFormatter plainFmt;
        TestDestination dest("level-test", &plainFmt);
        logger.addDestination(&dest);

        logger.log(LogLevel::INFO, "filtered");
        logger.log(LogLevel::ERROR, "passes");
        assert(dest.received.size() == 1);
        assert(dest.received[0].find("passes") != string::npos);

        logger.removeDestination(&dest);
        cout << "PASS test_level_filtering_before_destinations" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_level_filtering_before_destinations" << endl;
        failed++;
    }

    // Test 7: No destinations registered — log still works (no crash)
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::DEBUG);
        // No destinations added
        logger.log(LogLevel::INFO, "no destinations");
        // Should not crash, history may or may not record
        cout << "PASS test_no_destinations_no_crash" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_no_destinations_no_crash" << endl;
        failed++;
    }

    cout << "PART3_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
