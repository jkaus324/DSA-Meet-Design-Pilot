#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;


// Data class (given).

// HINT: introduce an abstraction so new ranking rules don't change existing code.
// HINT: keep the comparator small — one rule per class.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
void reset_service() {
    // TODO: write your solution
    // nothing to return
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
void init_limiter(int maxRequests, int windowSize) {
    // TODO: write your solution
    // nothing to return
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
bool allow_request_simple(string clientId, int timestamp, string endpoint) {
    // TODO: write your solution
    return {};
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
int get_request_count(string clientId) {
    // TODO: write your solution
    return {};
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
bool allow_request_with_strategy_simple(string algorithm, string clientId, int timestamp, string endpoint) {
    // TODO: write your solution
    return {};
}
