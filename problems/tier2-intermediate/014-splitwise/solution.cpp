#include <iostream>
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
