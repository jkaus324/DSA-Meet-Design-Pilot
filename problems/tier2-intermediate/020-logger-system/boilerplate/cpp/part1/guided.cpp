#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;


// Data class (given).
struct LogOp {
    string kind;
    string s1;
    string s2;
    int i1;
    LogOp(const string& kind_, const string& s1_ = "", const string& s2_ = "", int i1_ = 0)
      : kind(kind_), s1(s1_), s2(s2_), i1(i1_) {}
};

// HINT: introduce an abstraction so new ranking rules don't change existing code.
// HINT: keep the comparator small — one rule per class.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
vector<string> logger_simulate(vector<LogOp> ops) {
    // TODO: write your solution
    return {};
}
