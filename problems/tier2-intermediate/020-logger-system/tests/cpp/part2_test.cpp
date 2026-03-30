// Part 2 Tests — Logger System: Pluggable Formatters (Strategy)
// Tests PlainText, JSON, and CSV formatters with runtime swapping

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

int part2_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: PlainTextFormatter output format
    try {
        PlainTextFormatter fmt;
        LogEntry entry{LogLevel::ERROR, "disk full", "2024-01-15 10:30:00"};
        string result = fmt.format(entry);
        assert(result.find("[2024-01-15 10:30:00]") != string::npos);
        assert(result.find("[ERROR]") != string::npos);
        assert(result.find("disk full") != string::npos);
        cout << "PASS test_plaintext_format" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_plaintext_format" << endl;
        failed++;
    }

    // Test 2: JsonFormatter output format
    try {
        JsonFormatter fmt;
        LogEntry entry{LogLevel::WARN, "low memory", "2024-01-15 10:30:00"};
        string result = fmt.format(entry);
        assert(result.find("\"timestamp\"") != string::npos);
        assert(result.find("\"level\"") != string::npos);
        assert(result.find("\"message\"") != string::npos);
        assert(result.find("WARN") != string::npos);
        assert(result.find("low memory") != string::npos);
        assert(result.find("2024-01-15 10:30:00") != string::npos);
        cout << "PASS test_json_format" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_json_format" << endl;
        failed++;
    }

    // Test 3: CsvFormatter output format
    try {
        CsvFormatter fmt;
        LogEntry entry{LogLevel::INFO, "server started", "2024-01-15 10:30:00"};
        string result = fmt.format(entry);
        assert(result.find(",") != string::npos);
        assert(result.find("INFO") != string::npos);
        assert(result.find("server started") != string::npos);
        assert(result.find("2024-01-15 10:30:00") != string::npos);
        // Should NOT have brackets (that's PlainText)
        assert(result.find("[") == string::npos);
        cout << "PASS test_csv_format" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_csv_format" << endl;
        failed++;
    }

    // Test 4: Logger uses PlainText by default
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::DEBUG);
        logger.setFormatter(nullptr);  // reset to default
        logger.log(LogLevel::INFO, "default format");
        string entry = logger.getLogHistory()[0];
        assert(entry.find("[") != string::npos);
        assert(entry.find("[INFO]") != string::npos);
        cout << "PASS test_default_plaintext_formatter" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_default_plaintext_formatter" << endl;
        failed++;
    }

    // Test 5: Switch to JSON formatter at runtime
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::DEBUG);
        JsonFormatter jsonFmt;
        logger.setFormatter(&jsonFmt);
        logger.log(LogLevel::ERROR, "json test");
        string entry = logger.getLogHistory()[0];
        assert(entry.find("\"level\"") != string::npos);
        assert(entry.find("\"message\"") != string::npos);
        assert(entry.find("json test") != string::npos);
        // Reset formatter
        logger.setFormatter(nullptr);
        cout << "PASS test_switch_to_json_formatter" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_switch_to_json_formatter" << endl;
        failed++;
    }

    // Test 6: Switch to CSV formatter at runtime
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::DEBUG);
        CsvFormatter csvFmt;
        logger.setFormatter(&csvFmt);
        logger.log(LogLevel::WARN, "csv test");
        string entry = logger.getLogHistory()[0];
        assert(entry.find(",") != string::npos);
        assert(entry.find("WARN") != string::npos);
        assert(entry.find("[") == string::npos);  // no brackets in CSV
        // Reset formatter
        logger.setFormatter(nullptr);
        cout << "PASS test_switch_to_csv_formatter" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_switch_to_csv_formatter" << endl;
        failed++;
    }

    // Test 7: Level filtering still works with custom formatter
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::ERROR);
        JsonFormatter jsonFmt;
        logger.setFormatter(&jsonFmt);
        logger.log(LogLevel::INFO, "should be filtered");
        logger.log(LogLevel::ERROR, "should pass");
        assert(logger.getLogHistory().size() == 1);
        assert(logger.getLogHistory()[0].find("should pass") != string::npos);
        // Reset formatter
        logger.setFormatter(nullptr);
        cout << "PASS test_level_filtering_with_custom_formatter" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_level_filtering_with_custom_formatter" << endl;
        failed++;
    }

    cout << "PART2_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
