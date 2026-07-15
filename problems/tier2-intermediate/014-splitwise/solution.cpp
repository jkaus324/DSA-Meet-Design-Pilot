#include <iostream>
#include <memory>
#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>
#include <tuple>
#include <cmath>
using namespace std;

// ─── Data Structures ────────────────────────────────────────────────────────

struct User {
    string id;
    string name;
};

struct Split {
    string userId;
    double amount;
};

struct Expense {
    string id;
    string paidBy;
    double totalAmount;
    vector<Split> splits;
};

// ─── Strategy Interface ────────────────────────────────────────────────────

class SplitStrategy {
public:
    virtual vector<Split> split(double totalAmount,
                                const vector<string>& participants,
                                const vector<double>& params) = 0;
    virtual bool validate(double totalAmount,
                          const vector<string>& participants,
                          const vector<double>& params) = 0;
    virtual ~SplitStrategy() = default;
};

// ─── Concrete Strategies ───────────────────────────────────────────────────

class EqualSplit : public SplitStrategy {
public:
    vector<Split> split(double totalAmount,
                        const vector<string>& participants,
                        const vector<double>& params) override {
        vector<Split> result;
        double share = totalAmount / participants.size();
        for (auto& p : participants)
            result.push_back({p, share});
        return result;
    }
    bool validate(double totalAmount, const vector<string>& participants,
                  const vector<double>& params) override {
        return !participants.empty();
    }
};

class ExactSplit : public SplitStrategy {
public:
    vector<Split> split(double totalAmount,
                        const vector<string>& participants,
                        const vector<double>& params) override {
        vector<Split> result;
        for (size_t i = 0; i < participants.size(); i++)
            result.push_back({participants[i], params[i]});
        return result;
    }
    bool validate(double totalAmount, const vector<string>& participants,
                  const vector<double>& params) override {
        if (params.size() != participants.size()) return false;
        double sum = 0;
        for (double v : params) sum += v;
        return fabs(sum - totalAmount) < 1e-9;
    }
};

class PercentSplit : public SplitStrategy {
public:
    vector<Split> split(double totalAmount,
                        const vector<string>& participants,
                        const vector<double>& params) override {
        vector<Split> result;
        for (size_t i = 0; i < participants.size(); i++)
            result.push_back({participants[i], totalAmount * params[i] / 100.0});
        return result;
    }
    bool validate(double totalAmount, const vector<string>& participants,
                  const vector<double>& params) override {
        if (params.size() != participants.size()) return false;
        double sum = 0;
        for (double v : params) sum += v;
        return fabs(sum - 100.0) < 1e-9;
    }
};

// ─── Expense Manager ───────────────────────────────────────────────────────

class ExpenseManager {
    unordered_map<string, User> users;
    vector<Expense> expenses;
    unordered_map<string, unordered_map<string, double>> balances;

    void updateBalance(const string& debtor, const string& creditor, double amount) {
        if (debtor == creditor) return;
        // Net out: if creditor already owes debtor, reduce that first
        if (balances[creditor].count(debtor) && balances[creditor][debtor] > 0) {
            double offset = min(balances[creditor][debtor], amount);
            balances[creditor][debtor] -= offset;
            amount -= offset;
            if (balances[creditor][debtor] < 1e-9)
                balances[creditor].erase(debtor);
        }
        if (amount > 1e-9)
            balances[debtor][creditor] += amount;
    }

public:
    void addUser(const string& userId, const string& name) {
        users[userId] = {userId, name};
    }

    void addExpense(const string& expenseId, const string& paidBy,
                    double amount, const vector<string>& participants) {
        EqualSplit strategy;
        vector<Split> splits = strategy.split(amount, participants, {});
        Expense expense{expenseId, paidBy, amount, splits};
        expenses.push_back(expense);
        for (auto& s : splits) {
            updateBalance(s.userId, paidBy, s.amount);
        }
    }

    void addExpenseWithStrategy(const string& expenseId, const string& paidBy,
                                double amount, const vector<string>& participants,
                                SplitStrategy* strategy,
                                const vector<double>& params) {
        if (!strategy->validate(amount, participants, params)) return;
        vector<Split> splits = strategy->split(amount, participants, params);
        Expense expense{expenseId, paidBy, amount, splits};
        expenses.push_back(expense);
        for (auto& s : splits) {
            updateBalance(s.userId, paidBy, s.amount);
        }
    }

    unordered_map<string, unordered_map<string, double>> getBalances() const {
        return balances;
    }

    vector<tuple<string, string, double>> simplifyDebts() {
        // Step 1: Compute net balances
        unordered_map<string, double> net;
        for (auto& [debtor, creditors] : balances) {
            for (auto& [creditor, amount] : creditors) {
                net[debtor] -= amount;
                net[creditor] += amount;
            }
        }

        // Step 2: Separate into creditors and debtors
        vector<pair<string, double>> creditors, debtors;
        for (auto& [user, amount] : net) {
            if (amount > 1e-9) creditors.push_back({user, amount});
            else if (amount < -1e-9) debtors.push_back({user, -amount});
        }

        // Step 3: Sort descending by amount
        sort(creditors.begin(), creditors.end(),
             [](auto& a, auto& b) { return a.second > b.second; });
        sort(debtors.begin(), debtors.end(),
             [](auto& a, auto& b) { return a.second > b.second; });

        // Step 4: Greedy matching
        vector<tuple<string, string, double>> transactions;
        int i = 0, j = 0;
        while (i < (int)creditors.size() && j < (int)debtors.size()) {
            double settle = min(creditors[i].second, debtors[j].second);
            transactions.push_back({debtors[j].first, creditors[i].first, settle});
            creditors[i].second -= settle;
            debtors[j].second -= settle;
            if (creditors[i].second < 1e-9) i++;
            if (debtors[j].second < 1e-9) j++;
        }
        return transactions;
    }
};

// ─── Test Entry Points ─────────────────────────────────────────────────────

ExpenseManager manager;

void add_user(const string& userId, const string& name) {
    manager.addUser(userId, name);
}

void add_expense(const string& expenseId, const string& paidBy,
                 double amount, const vector<string>& participants) {
    manager.addExpense(expenseId, paidBy, amount, participants);
}

void add_expense_with_strategy(const string& expenseId, const string& paidBy,
                               double amount, const vector<string>& participants,
                               SplitStrategy* strategy,
                               const vector<double>& params) {
    manager.addExpenseWithStrategy(expenseId, paidBy, amount, participants,
                                  strategy, params);
}

unordered_map<string, unordered_map<string, double>> get_balances() {
    return manager.getBalances();
}

vector<tuple<string, string, double>> simplify_debts() {
    return manager.simplifyDebts();
}

// ─── Ops simulator (used by spec-based tests) ──────────────────────────────
//
// Drives one ExpenseManager through a sequence of operations.
//
// Op fields:
//   "new"                    -> "ok"
//   "add_user"  s1=id s2=name -> "ok"
//   "add_expense"  s1=expId s2=paidBy s3=participantsCSV  i1=amount(int) -> "ok"
//   "add_eq_strat" s1=expId s2=paidBy s3=participantsCSV  i1=amount       -> "ok"
//   "add_exact"    s1=expId s2=paidBy s3=participantsCSV  i1=amount  s4=amountsCSV -> "ok"
//   "add_pct"      s1=expId s2=paidBy s3=participantsCSV  i1=amount  s4=pctCSV     -> "ok"
//   "balance"      s1=debtor s2=creditor                                            -> "<num>" with 2 decimals (or "0.00")
//   "validate_exact" i1=total s3=partsCSV s4=amountsCSV  -> "yes"/"no"
//   "validate_pct"   i1=total s3=partsCSV s4=pctCSV      -> "yes"/"no"
//   "simplify_count"                                  -> int as string
//   "simplify_total_to"   s1=creditor                  -> "<num>"
//   "simplify_total_from" s1=debtor                    -> "<num>"
//   "simplify_unique_pair" s1=debtor s2=creditor       -> "<amount>" or "0.00"

struct SplitOp {
    string kind;
    string s1;
    string s2;
    string s3;
    string s4;
    int    i1;
};

static vector<string> split_csv(const string& s) {
    vector<string> out;
    string cur;
    for (char c : s) {
        if (c == ',') { out.push_back(cur); cur.clear(); }
        else cur.push_back(c);
    }
    if (!cur.empty()) out.push_back(cur);
    return out;
}
static vector<double> split_csv_double(const string& s) {
    auto parts = split_csv(s);
    vector<double> out;
    for (auto& p : parts) out.push_back(stod(p));
    return out;
}
static string num2(double v) {
    char buf[32];
    snprintf(buf, sizeof(buf), "%.2f", v);
    return string(buf);
}

vector<string> splitwise_simulate(vector<SplitOp> ops) {
    vector<string> out;
    unique_ptr<ExpenseManager> mgr(new ExpenseManager());
    for (const auto& op : ops) {
        const string& k = op.kind;
        if (k == "new") {
            mgr.reset(new ExpenseManager());
            out.push_back("ok");
        } else if (k == "add_user") {
            mgr->addUser(op.s1, op.s2);
            out.push_back("ok");
        } else if (k == "add_expense") {
            mgr->addExpense(op.s1, op.s2, (double)op.i1, split_csv(op.s3));
            out.push_back("ok");
        } else if (k == "add_eq_strat") {
            EqualSplit eq;
            mgr->addExpenseWithStrategy(op.s1, op.s2, (double)op.i1, split_csv(op.s3), &eq, {});
            out.push_back("ok");
        } else if (k == "add_exact") {
            ExactSplit ex;
            mgr->addExpenseWithStrategy(op.s1, op.s2, (double)op.i1, split_csv(op.s3), &ex, split_csv_double(op.s4));
            out.push_back("ok");
        } else if (k == "add_pct") {
            PercentSplit pct;
            mgr->addExpenseWithStrategy(op.s1, op.s2, (double)op.i1, split_csv(op.s3), &pct, split_csv_double(op.s4));
            out.push_back("ok");
        } else if (k == "balance") {
            auto b = mgr->getBalances();
            double v = 0.0;
            auto it = b.find(op.s1);
            if (it != b.end()) {
                auto it2 = it->second.find(op.s2);
                if (it2 != it->second.end()) v = it2->second;
            }
            out.push_back(num2(v));
        } else if (k == "validate_exact") {
            ExactSplit ex;
            out.push_back(ex.validate((double)op.i1, split_csv(op.s3), split_csv_double(op.s4)) ? "yes" : "no");
        } else if (k == "validate_pct") {
            PercentSplit pct;
            out.push_back(pct.validate((double)op.i1, split_csv(op.s3), split_csv_double(op.s4)) ? "yes" : "no");
        } else if (k == "simplify_count") {
            out.push_back(to_string((int)mgr->simplifyDebts().size()));
        } else if (k == "simplify_total_to") {
            auto txns = mgr->simplifyDebts();
            double tot = 0;
            for (auto& t : txns) if (get<1>(t) == op.s1) tot += get<2>(t);
            out.push_back(num2(tot));
        } else if (k == "simplify_total_from") {
            auto txns = mgr->simplifyDebts();
            double tot = 0;
            for (auto& t : txns) if (get<0>(t) == op.s1) tot += get<2>(t);
            out.push_back(num2(tot));
        } else if (k == "simplify_unique_pair") {
            auto txns = mgr->simplifyDebts();
            double tot = 0;
            for (auto& t : txns)
                if (get<0>(t) == op.s1 && get<1>(t) == op.s2) tot += get<2>(t);
            out.push_back(num2(tot));
        } else {
            out.push_back("unknown:" + k);
        }
    }
    return out;
}

// ─── Main ──────────────────────────────────────────────────────────────────

#ifndef RUNNING_TESTS
int main() {
    add_user("alice", "Alice");
    add_user("bob", "Bob");
    add_user("charlie", "Charlie");

    add_expense("E1", "alice", 300.0, {"alice", "bob", "charlie"});
    auto bal = get_balances();
    cout << "Bob owes Alice: " << bal["bob"]["alice"] << endl;
    cout << "Charlie owes Alice: " << bal["charlie"]["alice"] << endl;

    return 0;
}
#endif
