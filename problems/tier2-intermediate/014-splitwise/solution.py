"""Splitwise — equal/exact/percent splits + debt simplification."""


class SplitOp:
    def __init__(self, kind, s1="", s2="", s3="", s4="", i1=0):
        self.kind = kind
        self.s1 = s1
        self.s2 = s2
        self.s3 = s3
        self.s4 = s4
        self.i1 = i1


class EqualSplit:
    def split(self, totalAmount, participants, params):
        share = totalAmount / len(participants)
        return [(p, share) for p in participants]

    def validate(self, totalAmount, participants, params):
        return len(participants) > 0


class ExactSplit:
    def split(self, totalAmount, participants, params):
        return [(participants[i], params[i]) for i in range(len(participants))]

    def validate(self, totalAmount, participants, params):
        if len(params) != len(participants):
            return False
        return abs(sum(params) - totalAmount) < 1e-9


class PercentSplit:
    def split(self, totalAmount, participants, params):
        return [(participants[i], totalAmount * params[i] / 100.0) for i in range(len(participants))]

    def validate(self, totalAmount, participants, params):
        if len(params) != len(participants):
            return False
        return abs(sum(params) - 100.0) < 1e-9


class ExpenseManager:
    def __init__(self):
        self.users = {}
        self.expenses = []
        self.balances = {}  # debtor -> {creditor -> amount}

    def addUser(self, userId, name):
        self.users[userId] = name

    def _update_balance(self, debtor, creditor, amount):
        if debtor == creditor:
            return
        # If creditor already owes debtor, offset first
        if creditor in self.balances and debtor in self.balances[creditor] and self.balances[creditor][debtor] > 0:
            offset = min(self.balances[creditor][debtor], amount)
            self.balances[creditor][debtor] -= offset
            amount -= offset
            if self.balances[creditor][debtor] < 1e-9:
                del self.balances[creditor][debtor]
        if amount > 1e-9:
            if debtor not in self.balances:
                self.balances[debtor] = {}
            self.balances[debtor][creditor] = self.balances[debtor].get(creditor, 0.0) + amount

    def addExpense(self, expenseId, paidBy, amount, participants):
        strategy = EqualSplit()
        splits = strategy.split(amount, participants, [])
        self.expenses.append((expenseId, paidBy, amount, splits))
        for uid, amt in splits:
            self._update_balance(uid, paidBy, amt)

    def addExpenseWithStrategy(self, expenseId, paidBy, amount, participants, strategy, params):
        if not strategy.validate(amount, participants, params):
            return
        splits = strategy.split(amount, participants, params)
        self.expenses.append((expenseId, paidBy, amount, splits))
        for uid, amt in splits:
            self._update_balance(uid, paidBy, amt)

    def getBalances(self):
        return self.balances

    def simplifyDebts(self):
        net = {}
        for debtor, creditors in self.balances.items():
            for creditor, amount in creditors.items():
                net[debtor] = net.get(debtor, 0.0) - amount
                net[creditor] = net.get(creditor, 0.0) + amount

        creditors = []
        debtors = []
        for user, amount in net.items():
            if amount > 1e-9:
                creditors.append([user, amount])
            elif amount < -1e-9:
                debtors.append([user, -amount])

        creditors.sort(key=lambda x: -x[1])
        debtors.sort(key=lambda x: -x[1])

        transactions = []
        i, j = 0, 0
        while i < len(creditors) and j < len(debtors):
            settle = min(creditors[i][1], debtors[j][1])
            transactions.append((debtors[j][0], creditors[i][0], settle))
            creditors[i][1] -= settle
            debtors[j][1] -= settle
            if creditors[i][1] < 1e-9:
                i += 1
            if debtors[j][1] < 1e-9:
                j += 1
        return transactions


def _split_csv(s):
    return s.split(",") if s else []


def _split_csv_double(s):
    return [float(x) for x in s.split(",")] if s else []


def _num2(v):
    return f"{v:.2f}"


def splitwise_simulate(ops):
    out = []
    mgr = ExpenseManager()
    for op in ops:
        k = op.kind
        if k == "new":
            mgr = ExpenseManager()
            out.append("ok")
        elif k == "add_user":
            mgr.addUser(op.s1, op.s2)
            out.append("ok")
        elif k == "add_expense":
            mgr.addExpense(op.s1, op.s2, float(op.i1), _split_csv(op.s3))
            out.append("ok")
        elif k == "add_eq_strat":
            mgr.addExpenseWithStrategy(op.s1, op.s2, float(op.i1), _split_csv(op.s3), EqualSplit(), [])
            out.append("ok")
        elif k == "add_exact":
            mgr.addExpenseWithStrategy(op.s1, op.s2, float(op.i1), _split_csv(op.s3), ExactSplit(), _split_csv_double(op.s4))
            out.append("ok")
        elif k == "add_pct":
            mgr.addExpenseWithStrategy(op.s1, op.s2, float(op.i1), _split_csv(op.s3), PercentSplit(), _split_csv_double(op.s4))
            out.append("ok")
        elif k == "balance":
            b = mgr.getBalances()
            v = 0.0
            if op.s1 in b and op.s2 in b[op.s1]:
                v = b[op.s1][op.s2]
            out.append(_num2(v))
        elif k == "validate_exact":
            ok = ExactSplit().validate(float(op.i1), _split_csv(op.s3), _split_csv_double(op.s4))
            out.append("yes" if ok else "no")
        elif k == "validate_pct":
            ok = PercentSplit().validate(float(op.i1), _split_csv(op.s3), _split_csv_double(op.s4))
            out.append("yes" if ok else "no")
        elif k == "simplify_count":
            out.append(str(len(mgr.simplifyDebts())))
        elif k == "simplify_total_to":
            txns = mgr.simplifyDebts()
            tot = sum(t[2] for t in txns if t[1] == op.s1)
            out.append(_num2(tot))
        elif k == "simplify_total_from":
            txns = mgr.simplifyDebts()
            tot = sum(t[2] for t in txns if t[0] == op.s1)
            out.append(_num2(tot))
        elif k == "simplify_unique_pair":
            txns = mgr.simplifyDebts()
            tot = sum(t[2] for t in txns if t[0] == op.s1 and t[1] == op.s2)
            out.append(_num2(tot))
        else:
            out.append("unknown:" + k)
    return out
