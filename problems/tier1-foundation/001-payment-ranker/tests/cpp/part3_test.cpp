// Part 3 Tests — Easy-Refund Eligibility Filter
// Tests that easy-refund preference correctly reorders results

#include <cassert>
#include <iostream>
using namespace std;

int part3_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: preferEasyRefund=true puts eligible methods first
    try {
        vector<PaymentMethod> methods = {
            {"Card A", 0.10, 5.0, 300, false},  // high cashback, no easy refund
            {"Card B", 0.02, 2.0, 500, true},   // low cashback, easy refund
            {"Card C", 0.05, 3.0, 400, true},   // medium cashback, easy refund
        };
        auto ranked = rank_with_refund_filter(methods, true);
        assert(ranked.size() == 3);
        // Both B and C have easy refund, so they come first (in cashback order)
        assert(ranked[0].easyRefundEligible == true);
        assert(ranked[1].easyRefundEligible == true);
        assert(ranked[2].name == "Card A"); // no easy refund, goes last
        cout << "PASS test_easy_refund_preferred" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_easy_refund_preferred" << endl;
        failed++;
    }

    // Test 2: preferEasyRefund=false should not reorder by refund eligibility
    try {
        vector<PaymentMethod> methods = {
            {"Card A", 0.10, 5.0, 300, false},
            {"Card B", 0.02, 2.0, 500, true},
        };
        auto ranked = rank_with_refund_filter(methods, false);
        // Without refund preference, Card A should still win (higher cashback)
        assert(ranked[0].name == "Card A");
        cout << "PASS test_refund_filter_disabled" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_refund_filter_disabled" << endl;
        failed++;
    }

    // Test 3: all methods have easy refund — order falls back to other criteria
    try {
        vector<PaymentMethod> methods = {
            {"Card A", 0.10, 5.0, 300, true},
            {"Card B", 0.02, 2.0, 500, true},
        };
        auto ranked = rank_with_refund_filter(methods, true);
        // All have easy refund, so tiebreak by cashback
        assert(ranked[0].name == "Card A"); // higher cashback
        cout << "PASS test_all_refund_eligible_tiebreak" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_all_refund_eligible_tiebreak" << endl;
        failed++;
    }

    cout << "PART3_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
