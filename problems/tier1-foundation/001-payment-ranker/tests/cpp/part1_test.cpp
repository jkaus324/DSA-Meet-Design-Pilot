// Part 1 Tests — Payment Method Ranker
// Tests the three basic ranking strategies: rewards, low-fee, trust

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

int part1_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: rank_by_rewards — highest cashback first
    try {
        vector<PaymentMethod> methods = {
            {"UPI",           0.01, 0.0, 1000},
            {"Credit Card A", 0.05, 5.0, 500},
            {"Credit Card B", 0.10, 8.0, 300},
        };
        auto ranked = rank_by_rewards(methods);
        assert(ranked.size() == 3);
        assert(ranked[0].name == "Credit Card B"); // 10% cashback
        assert(ranked[1].name == "Credit Card A"); // 5% cashback
        assert(ranked[2].name == "UPI");           // 1% cashback
        cout << "PASS test_rewards_ranking" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_rewards_ranking" << endl;
        failed++;
    }

    // Test 2: rank_by_low_fee — lowest fee first
    try {
        vector<PaymentMethod> methods = {
            {"Debit Card",    0.0,  2.0, 800},
            {"Credit Card A", 0.05, 5.0, 500},
            {"Credit Card B", 0.10, 8.0, 300},
            {"UPI",           0.01, 0.0, 1000},
        };
        auto ranked = rank_by_low_fee(methods);
        assert(ranked.size() == 4);
        assert(ranked[0].name == "UPI");         // 0 fee
        assert(ranked[1].name == "Debit Card");  // 2.0 fee
        assert(ranked[2].name == "Credit Card A"); // 5.0 fee
        assert(ranked[3].name == "Credit Card B"); // 8.0 fee
        cout << "PASS test_low_fee_ranking" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_low_fee_ranking" << endl;
        failed++;
    }

    // Test 3: rank_by_trust — highest usage count first
    try {
        vector<PaymentMethod> methods = {
            {"Credit Card A", 0.05, 5.0, 500},
            {"UPI",           0.01, 0.0, 1000},
            {"Debit Card",    0.0,  2.0, 800},
        };
        auto ranked = rank_by_trust(methods);
        assert(ranked.size() == 3);
        assert(ranked[0].name == "UPI");          // 1000 uses
        assert(ranked[1].name == "Debit Card");   // 800 uses
        assert(ranked[2].name == "Credit Card A"); // 500 uses
        cout << "PASS test_trust_ranking" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_trust_ranking" << endl;
        failed++;
    }

    // Test 4: empty input returns empty
    try {
        vector<PaymentMethod> empty;
        assert(rank_by_rewards(empty).empty());
        assert(rank_by_low_fee(empty).empty());
        assert(rank_by_trust(empty).empty());
        cout << "PASS test_empty_input" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_empty_input" << endl;
        failed++;
    }

    // Test 5: single item returns itself
    try {
        vector<PaymentMethod> single = {{"UPI", 0.01, 0.0, 100}};
        assert(rank_by_rewards(single).size() == 1);
        assert(rank_by_rewards(single)[0].name == "UPI");
        cout << "PASS test_single_item" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_single_item" << endl;
        failed++;
    }

    cout << "PART1_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
