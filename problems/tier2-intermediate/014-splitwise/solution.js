'use strict';

/* Splitwise — equal/exact/percent splits + debt simplification. */

class SplitOp {
  constructor(kind, s1 = '', s2 = '', s3 = '', s4 = '', i1 = 0) {
    this.kind = kind;
    this.s1 = s1;
    this.s2 = s2;
    this.s3 = s3;
    this.s4 = s4;
    this.i1 = i1;
  }
}

class EqualSplit {
  split(totalAmount, participants, params) {
    const share = totalAmount / participants.length;
    return participants.map(p => [p, share]);
  }
  validate(totalAmount, participants, params) {
    return participants.length > 0;
  }
}

class ExactSplit {
  split(totalAmount, participants, params) {
    const res = [];
    for (let i = 0; i < participants.length; i++) res.push([participants[i], params[i]]);
    return res;
  }
  validate(totalAmount, participants, params) {
    if (params.length !== participants.length) return false;
    const s = params.reduce((a, b) => a + b, 0);
    return Math.abs(s - totalAmount) < 1e-9;
  }
}

class PercentSplit {
  split(totalAmount, participants, params) {
    const res = [];
    for (let i = 0; i < participants.length; i++) {
      res.push([participants[i], (totalAmount * params[i]) / 100.0]);
    }
    return res;
  }
  validate(totalAmount, participants, params) {
    if (params.length !== participants.length) return false;
    const s = params.reduce((a, b) => a + b, 0);
    return Math.abs(s - 100.0) < 1e-9;
  }
}

class ExpenseManager {
  constructor() {
    this.users = new Map();
    this.expenses = [];
    this.balances = new Map(); // debtor -> Map(creditor -> amount)
  }

  addUser(userId, name) {
    this.users.set(userId, name);
  }

  _update_balance(debtor, creditor, amount) {
    if (debtor === creditor) return;
    // If creditor already owes debtor, offset first
    if (
      this.balances.has(creditor) &&
      this.balances.get(creditor).has(debtor) &&
      this.balances.get(creditor).get(debtor) > 0
    ) {
      const cMap = this.balances.get(creditor);
      const offset = Math.min(cMap.get(debtor), amount);
      cMap.set(debtor, cMap.get(debtor) - offset);
      amount -= offset;
      if (cMap.get(debtor) < 1e-9) {
        cMap.delete(debtor);
      }
    }
    if (amount > 1e-9) {
      if (!this.balances.has(debtor)) {
        this.balances.set(debtor, new Map());
      }
      const dMap = this.balances.get(debtor);
      const cur = dMap.has(creditor) ? dMap.get(creditor) : 0.0;
      dMap.set(creditor, cur + amount);
    }
  }

  addExpense(expenseId, paidBy, amount, participants) {
    const strategy = new EqualSplit();
    const splits = strategy.split(amount, participants, []);
    this.expenses.push([expenseId, paidBy, amount, splits]);
    for (const [uid, amt] of splits) {
      this._update_balance(uid, paidBy, amt);
    }
  }

  addExpenseWithStrategy(expenseId, paidBy, amount, participants, strategy, params) {
    if (!strategy.validate(amount, participants, params)) return;
    const splits = strategy.split(amount, participants, params);
    this.expenses.push([expenseId, paidBy, amount, splits]);
    for (const [uid, amt] of splits) {
      this._update_balance(uid, paidBy, amt);
    }
  }

  getBalances() {
    return this.balances;
  }

  simplifyDebts() {
    const net = new Map();
    for (const [debtor, creditors] of this.balances) {
      for (const [creditor, amount] of creditors) {
        net.set(debtor, (net.has(debtor) ? net.get(debtor) : 0.0) - amount);
        net.set(creditor, (net.has(creditor) ? net.get(creditor) : 0.0) + amount);
      }
    }

    const creditors = [];
    const debtors = [];
    for (const [user, amount] of net) {
      if (amount > 1e-9) creditors.push([user, amount]);
      else if (amount < -1e-9) debtors.push([user, -amount]);
    }

    creditors.sort((a, b) => b[1] - a[1]);
    debtors.sort((a, b) => b[1] - a[1]);

    const transactions = [];
    let i = 0;
    let j = 0;
    while (i < creditors.length && j < debtors.length) {
      const settle = Math.min(creditors[i][1], debtors[j][1]);
      transactions.push([debtors[j][0], creditors[i][0], settle]);
      creditors[i][1] -= settle;
      debtors[j][1] -= settle;
      if (creditors[i][1] < 1e-9) i += 1;
      if (debtors[j][1] < 1e-9) j += 1;
    }
    return transactions;
  }
}

function _split_csv(s) {
  return s ? s.split(',') : [];
}

function _split_csv_double(s) {
  return s ? s.split(',').map(x => parseFloat(x)) : [];
}

function _num2(v) {
  return v.toFixed(2);
}

function splitwise_simulate(ops) {
  const out = [];
  let mgr = new ExpenseManager();
  for (const op of ops) {
    const k = op.kind;
    if (k === 'new') {
      mgr = new ExpenseManager();
      out.push('ok');
    } else if (k === 'add_user') {
      mgr.addUser(op.s1, op.s2);
      out.push('ok');
    } else if (k === 'add_expense') {
      mgr.addExpense(op.s1, op.s2, Number(op.i1), _split_csv(op.s3));
      out.push('ok');
    } else if (k === 'add_eq_strat') {
      mgr.addExpenseWithStrategy(op.s1, op.s2, Number(op.i1), _split_csv(op.s3), new EqualSplit(), []);
      out.push('ok');
    } else if (k === 'add_exact') {
      mgr.addExpenseWithStrategy(op.s1, op.s2, Number(op.i1), _split_csv(op.s3), new ExactSplit(), _split_csv_double(op.s4));
      out.push('ok');
    } else if (k === 'add_pct') {
      mgr.addExpenseWithStrategy(op.s1, op.s2, Number(op.i1), _split_csv(op.s3), new PercentSplit(), _split_csv_double(op.s4));
      out.push('ok');
    } else if (k === 'balance') {
      const b = mgr.getBalances();
      let v = 0.0;
      if (b.has(op.s1) && b.get(op.s1).has(op.s2)) {
        v = b.get(op.s1).get(op.s2);
      }
      out.push(_num2(v));
    } else if (k === 'validate_exact') {
      const ok = new ExactSplit().validate(Number(op.i1), _split_csv(op.s3), _split_csv_double(op.s4));
      out.push(ok ? 'yes' : 'no');
    } else if (k === 'validate_pct') {
      const ok = new PercentSplit().validate(Number(op.i1), _split_csv(op.s3), _split_csv_double(op.s4));
      out.push(ok ? 'yes' : 'no');
    } else if (k === 'simplify_count') {
      out.push(String(mgr.simplifyDebts().length));
    } else if (k === 'simplify_total_to') {
      const txns = mgr.simplifyDebts();
      let tot = 0.0;
      for (const t of txns) if (t[1] === op.s1) tot += t[2];
      out.push(_num2(tot));
    } else if (k === 'simplify_total_from') {
      const txns = mgr.simplifyDebts();
      let tot = 0.0;
      for (const t of txns) if (t[0] === op.s1) tot += t[2];
      out.push(_num2(tot));
    } else if (k === 'simplify_unique_pair') {
      const txns = mgr.simplifyDebts();
      let tot = 0.0;
      for (const t of txns) if (t[0] === op.s1 && t[1] === op.s2) tot += t[2];
      out.push(_num2(tot));
    } else {
      out.push('unknown:' + k);
    }
  }
  return out;
}

module.exports = { SplitOp, splitwise_simulate };
