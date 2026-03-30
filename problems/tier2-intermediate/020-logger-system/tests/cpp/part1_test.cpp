// Part 1 Tests — Logger System: Singleton + Level Filtering
// Tests the basic logger with level filtering and singleton access

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

int part1_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: Singleton — getInstance returns same instance
    try {
        Logger& logger1 = Logger::getInstance();
        Logger& logger2 = Logger::getInstance();
        assert(&logger1 == &logger2);
        logger1.clearHistory();
        cout << "PASS test_singleton_same_instance" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_singleton_same_instance" << endl;
        failed++;
    }

    // Test 2: Default level is INFO — DEBUG messages are filtered
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::INFO);
        logger.log(LogLevel::DEBUG, "debug msg");
        assert(logger.getLogHistory().empty());
        cout << "PASS test_debug_filtered_at_info_level" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_debug_filtered_at_info_level" << endl;
        failed++;
    }

    // Test 3: INFO message passes at INFO level
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::INFO);
        logger.log(LogLevel::INFO, "info msg");
        assert(logger.getLogHistory().size() == 1);
        // Check that the log contains the level and message
        string entry = logger.getLogHistory()[0];
        assert(entry.find("INFO") != string::npos);
        assert(entry.find("info msg") != string::npos);
        cout << "PASS test_info_passes_at_info_level" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_info_passes_at_info_level" << endl;
        failed++;
    }

    // Test 4: ERROR and FATAL pass at INFO level
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::INFO);
        logger.log(LogLevel::ERROR, "error msg");
        logger.log(LogLevel::FATAL, "fatal msg");
        assert(logger.getLogHistory().size() == 2);
        assert(logger.getLogHistory()[0].find("ERROR") != string::npos);
        assert(logger.getLogHistory()[1].find("FATAL") != string::npos);
        cout << "PASS test_error_fatal_pass_at_info_level" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_error_fatal_pass_at_info_level" << endl;
        failed++;
    }

    // Test 5: Changing level to WARN filters out INFO
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::WARN);
        logger.log(LogLevel::INFO, "should be filtered");
        logger.log(LogLevel::WARN, "should pass");
        logger.log(LogLevel::ERROR, "should also pass");
        assert(logger.getLogHistory().size() == 2);
        assert(logger.getLogHistory()[0].find("WARN") != string::npos);
        assert(logger.getLogHistory()[1].find("ERROR") != string::npos);
        cout << "PASS test_warn_level_filters_info" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_warn_level_filters_info" << endl;
        failed++;
    }

    // Test 6: setLevel to DEBUG allows all messages
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::DEBUG);
        logger.log(LogLevel::DEBUG, "d");
        logger.log(LogLevel::INFO, "i");
        logger.log(LogLevel::WARN, "w");
        logger.log(LogLevel::ERROR, "e");
        logger.log(LogLevel::FATAL, "f");
        assert(logger.getLogHistory().size() == 5);
        cout << "PASS test_debug_level_allows_all" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_debug_level_allows_all" << endl;
        failed++;
    }

    // Test 7: setLevel to FATAL only allows FATAL
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::FATAL);
        logger.log(LogLevel::DEBUG, "d");
        logger.log(LogLevel::INFO, "i");
        logger.log(LogLevel::WARN, "w");
        logger.log(LogLevel::ERROR, "e");
        logger.log(LogLevel::FATAL, "f");
        assert(logger.getLogHistory().size() == 1);
        assert(logger.getLogHistory()[0].find("FATAL") != string::npos);
        cout << "PASS test_fatal_level_only_fatal" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_fatal_level_only_fatal" << endl;
        failed++;
    }

    // Test 8: Log format contains timestamp brackets
    try {
        Logger& logger = Logger::getInstance();
        logger.clearHistory();
        logger.setLevel(LogLevel::INFO);
        logger.log(LogLevel::INFO, "check format");
        string entry = logger.getLogHistory()[0];
        // Should contain brackets around timestamp and level
        assert(entry.find("[") != string::npos);
        assert(entry.find("]") != string::npos);
        assert(entry.find("check format") != string::npos);
        cout << "PASS test_log_format_has_brackets" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_log_format_has_brackets" << endl;
        failed++;
    }

    cout << "PART1_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
