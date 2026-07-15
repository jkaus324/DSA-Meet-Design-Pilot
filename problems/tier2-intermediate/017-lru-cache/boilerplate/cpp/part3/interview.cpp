#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;


// Data class (given).
struct LruOp {
    string kind;
    int i1;
    int i2;
    int i3;
    int i4;
    LruOp(const string& kind_, int i1_ = 0, int i2_ = 0, int i3_ = 0, int i4_ = 0)
      : kind(kind_), i1(i1_), i2(i2_), i3(i3_), i4(i4_) {}
};

// TODO: design and implement your solution.
// Required free functions:
//   vector<string> lru_simulate(vector<LruOp> ops);

vector<string> lru_simulate(vector<LruOp> ops) {
    // TODO: write your solution
    return {};
}
