#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;


// Data class (given).
struct TwitterOp {
    string kind;
    int i1;
    int i2;
    TwitterOp(const string& kind_, int i1_ = 0, int i2_ = 0)
      : kind(kind_), i1(i1_), i2(i2_) {}
};

// HINT: introduce an abstraction so new ranking rules don't change existing code.
// HINT: keep the comparator small — one rule per class.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
vector<string> twitter_simulate(vector<TwitterOp> ops) {
    // TODO: write your solution
    return {};
}
