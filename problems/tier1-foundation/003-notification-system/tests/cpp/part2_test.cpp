// Part 2 Tests — Priority Filtering
// Tests that notifications respect per-user priority preferences

#include <cassert>
#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
using namespace std;

int part2_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: critical event reaches all users regardless of min priority
    try {
        User u1 = {"u1", "u1@test.com", "+1-555-0001", {"email"}};
        unordered_map<string, string> prefs = {{"*", "critical"}};
        // Should not throw
        notify("CRITICAL: System down", "critical", {u1}, prefs);
        cout << "PASS test_critical_reaches_all" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_critical_reaches_all" << endl;
        failed++;
    }

    // Test 2: promotional event blocked when user wants info+ only
    try {
        User u1 = {"u1", "u1@test.com", "+1-555-0001", {"email"}};
        unordered_map<string, string> prefs = {{"*", "info"}};
        // Promotional should be filtered out — no exception expected
        notify("50% off sale!", "promotional", {u1}, prefs);
        cout << "PASS test_promotional_filtered" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_promotional_filtered" << endl;
        failed++;
    }

    // Test 3: empty priority prefs defaults to sending all events
    try {
        User u1 = {"u1", "u1@test.com", "+1-555-0001", {"email"}};
        unordered_map<string, string> emptyPrefs;
        notify("Informational update", "info", {u1}, emptyPrefs);
        cout << "PASS test_empty_prefs_allow_all" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_empty_prefs_allow_all" << endl;
        failed++;
    }

    cout << "PART2_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
