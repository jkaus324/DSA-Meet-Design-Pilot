#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;


// Data class (given).
struct Op {
    string kind;
    string s1;
    string s2;
    string s3;
    int i1;
    int i2;
    int i3;
    Op(const string& kind_, const string& s1_ = "", const string& s2_ = "", const string& s3_ = "", int i1_ = 0, int i2_ = 0, int i3_ = 0)
      : kind(kind_), s1(s1_), s2(s2_), s3(s3_), i1(i1_), i2(i2_), i3(i3_) {}
};

// TODO: design and implement your solution.
// Required free functions:
//   vector<string> meeting_simulate(vector<Op> ops);

vector<string> meeting_simulate(vector<Op> ops) {
    // TODO: write your solution
    return {};
}
