#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;


// Data class (given — do not modify).
struct SplitOp {
    string kind;
    string s1;
    string s2;
    string s3;
    string s4;
    int i1;
    SplitOp(const string& kind_, const string& s1_ = "", const string& s2_ = "", const string& s3_ = "", const string& s4_ = "", int i1_ = 0)
      : kind(kind_), s1(s1_), s2(s2_), s3(s3_), s4(s4_), i1(i1_) {}
};

vector<string> splitwise_simulate(vector<SplitOp> ops) {
    // TODO: implement this
    return {};
}
