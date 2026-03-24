// Part 2 Tests — Composite Ranking
// Tests composite ranking with multiple chained criteria

#include <cassert>
#include <iostream>
using namespace std;

// Note: solution.cpp is included via the harness (not directly here)
// These tests assume rank_composite() is available

int part2_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: composite ranking — cashback first, then fee as tiebreaker
    try {
        vector<PaymentMethod> methods = {
            {"Card A", 0.10, 8.0, 300},  // 10% cashback, high fee
            {"Card B", 0.10, 3.0, 400},  // 10% cashback, low fee
            {"Card C", 0.05, 1.0, 200},  // 5% cashback
        };
        RewardsMaximizer rewards;
        LowFeeSeeker fee;
        auto ranked = rank_composite(methods, {&rewards, &fee});
        assert(ranked.size() == 3);
        assert(ranked[0].name == "Card B"); // tied cashback, lower fee wins
        assert(ranked[1].name == "Card A"); // tied cashback, higher fee loses
        assert(ranked[2].name == "Card C"); // lower cashback
        cout << "PASS test_composite_cashback_then_fee" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_composite_cashback_then_fee" << endl;
        failed++;
    }

    // Test 2: composite ranking — trust first, then cashback
    try {
        vector<PaymentMethod> methods = {
            {"UPI",    0.01, 0.0, 1000},
            {"Card A", 0.10, 5.0, 200},
            {"Card B", 0.05, 3.0, 1000},  // tied trust with UPI
        };
        TrustBasedRanker trust;
        RewardsMaximizer rewards;
        auto ranked = rank_composite(methods, {&trust, &rewards});
        assert(ranked.size() == 3);
        // UPI and Card B both have 1000 uses — tiebreak by cashback
        // Card B has 5% cashback > UPI's 1%
        assert(ranked[0].name == "Card B");
        assert(ranked[1].name == "UPI");
        assert(ranked[2].name == "Card A");
        cout << "PASS test_composite_trust_then_cashback" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_composite_trust_then_cashback" << endl;
        failed++;
    }

    // Test 3: single criterion composite behaves like that criterion alone
    try {
        vector<PaymentMethod> methods = {
            {"Card X", 0.02, 5.0, 100},
            {"Card Y", 0.08, 3.0, 200},
        };
        RewardsMaximizer rewards;
        auto composite = rank_composite(methods, {&rewards});
        auto direct    = rank_by_rewards(methods);
        assert(composite[0].name == direct[0].name);
        assert(composite[1].name == direct[1].name);
        cout << "PASS test_single_criterion_composite" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_single_criterion_composite" << endl;
        failed++;
    }

    cout << "PART2_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
