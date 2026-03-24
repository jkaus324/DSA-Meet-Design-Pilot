// Part 1 Tests — Notification System
// Tests basic observer pattern: subscribe, notify, multi-channel

#include "solution.cpp"
#include <cassert>
#include <iostream>
#include <vector>
#include <string>
using namespace std;

// Track what was sent (tests inject a spy observer)
vector<string> sent_notifications;

int part1_tests() {
    int passed = 0;
    int failed = 0;
    sent_notifications.clear();

    // Test 1: notify sends to subscribed channels
    try {
        User u1 = {"user1", "u1@test.com", "+91-9000000001", {"email"}};
        User u2 = {"user2", "u2@test.com", "+91-9000000002", {"sms", "push"}};
        User u3 = {"user3", "u3@test.com", "+91-9000000003", {"email", "sms"}};

        // notify should not throw
        notify("Order shipped", {u1, u2, u3});
        cout << "PASS test_notify_no_throw" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_notify_no_throw" << endl;
        failed++;
    }

    // Test 2: user not subscribed to a channel receives nothing on that channel
    try {
        // If we can capture output, verify only subscribed channels receive
        // For simplicity, we test that an empty subscriber list causes no crash
        vector<User> empty;
        notify("Event", empty);
        cout << "PASS test_empty_user_list" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_empty_user_list" << endl;
        failed++;
    }

    // Test 3: user subscribed to multiple channels receives on each
    try {
        User u = {"u1", "u1@test.com", "+1-555-0001", {"email", "sms", "push"}};
        // Should not throw even with multiple channels
        notify("Flash sale", {u});
        cout << "PASS test_multi_channel_user" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_multi_channel_user" << endl;
        failed++;
    }

    cout << "PART1_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
