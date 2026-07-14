// Splitwise — Solution (Java)
import java.util.*;

class SplitOp {
    public String kind;
    public String s1;
    public String s2;
    public String s3;
    public String s4;
    public int i1;

    public SplitOp(String kind, String s1, String s2, String s3, String s4, int i1) {
        this.kind = kind;
        this.s1 = s1;
        this.s2 = s2;
        this.s3 = s3;
        this.s4 = s4;
        this.i1 = i1;
    }
}

class User {
    public String id;
    public String name;
    public User(String id, String name) { this.id = id; this.name = name; }
}

class Split {
    public String userId;
    public double amount;
    public Split(String userId, double amount) { this.userId = userId; this.amount = amount; }
}

class Expense {
    public String id;
    public String paidBy;
    public double totalAmount;
    public List<Split> splits;

    public Expense(String id, String paidBy, double totalAmount, List<Split> splits) {
        this.id = id;
        this.paidBy = paidBy;
        this.totalAmount = totalAmount;
        this.splits = splits;
    }
}

interface SplitStrategy {
    List<Split> split(double totalAmount, List<String> participants, List<Double> params);
    boolean validate(double totalAmount, List<String> participants, List<Double> params);
}

class EqualSplit implements SplitStrategy {
    @Override
    public List<Split> split(double totalAmount, List<String> participants, List<Double> params) {
        List<Split> result = new ArrayList<>();
        double share = totalAmount / participants.size();
        for (String p : participants) result.add(new Split(p, share));
        return result;
    }
    @Override
    public boolean validate(double totalAmount, List<String> participants, List<Double> params) {
        return !participants.isEmpty();
    }
}

class ExactSplit implements SplitStrategy {
    @Override
    public List<Split> split(double totalAmount, List<String> participants, List<Double> params) {
        List<Split> result = new ArrayList<>();
        for (int i = 0; i < participants.size(); i++)
            result.add(new Split(participants.get(i), params.get(i)));
        return result;
    }
    @Override
    public boolean validate(double totalAmount, List<String> participants, List<Double> params) {
        if (params.size() != participants.size()) return false;
        double sum = 0;
        for (double v : params) sum += v;
        return Math.abs(sum - totalAmount) < 1e-9;
    }
}

class PercentSplit implements SplitStrategy {
    @Override
    public List<Split> split(double totalAmount, List<String> participants, List<Double> params) {
        List<Split> result = new ArrayList<>();
        for (int i = 0; i < participants.size(); i++)
            result.add(new Split(participants.get(i), totalAmount * params.get(i) / 100.0));
        return result;
    }
    @Override
    public boolean validate(double totalAmount, List<String> participants, List<Double> params) {
        if (params.size() != participants.size()) return false;
        double sum = 0;
        for (double v : params) sum += v;
        return Math.abs(sum - 100.0) < 1e-9;
    }
}

class DebtTxn {
    public String debtor;
    public String creditor;
    public double amount;
    public DebtTxn(String d, String c, double a) { debtor = d; creditor = c; amount = a; }
}

class ExpenseManager {
    Map<String, User> users = new LinkedHashMap<>();
    List<Expense> expenses = new ArrayList<>();
    Map<String, Map<String, Double>> balances = new LinkedHashMap<>();

    private void updateBalance(String debtor, String creditor, double amount) {
        if (debtor.equals(creditor)) return;

        Map<String, Double> credInner = balances.get(creditor);
        if (credInner != null && credInner.containsKey(debtor) && credInner.get(debtor) > 0) {
            double offset = Math.min(credInner.get(debtor), amount);
            credInner.put(debtor, credInner.get(debtor) - offset);
            amount -= offset;
            if (credInner.get(debtor) < 1e-9) credInner.remove(debtor);
        }
        if (amount > 1e-9) {
            balances.computeIfAbsent(debtor, k -> new LinkedHashMap<>())
                    .merge(creditor, amount, Double::sum);
        }
    }

    public void addUser(String userId, String name) {
        users.put(userId, new User(userId, name));
    }

    public void addExpense(String expenseId, String paidBy, double amount, List<String> participants) {
        EqualSplit strategy = new EqualSplit();
        List<Split> splits = strategy.split(amount, participants, new ArrayList<>());
        expenses.add(new Expense(expenseId, paidBy, amount, splits));
        for (Split s : splits) updateBalance(s.userId, paidBy, s.amount);
    }

    public void addExpenseWithStrategy(String expenseId, String paidBy, double amount,
                                        List<String> participants, SplitStrategy strategy,
                                        List<Double> params) {
        if (!strategy.validate(amount, participants, params)) return;
        List<Split> splits = strategy.split(amount, participants, params);
        expenses.add(new Expense(expenseId, paidBy, amount, splits));
        for (Split s : splits) updateBalance(s.userId, paidBy, s.amount);
    }

    public Map<String, Map<String, Double>> getBalances() { return balances; }

    public List<DebtTxn> simplifyDebts() {
        Map<String, Double> net = new LinkedHashMap<>();
        for (Map.Entry<String, Map<String, Double>> e : balances.entrySet()) {
            String debtor = e.getKey();
            for (Map.Entry<String, Double> c : e.getValue().entrySet()) {
                String creditor = c.getKey();
                double amount = c.getValue();
                net.merge(debtor, -amount, Double::sum);
                net.merge(creditor, amount, Double::sum);
            }
        }

        List<double[]> creditors = new ArrayList<>();
        List<double[]> debtors = new ArrayList<>();
        List<String> credKeys = new ArrayList<>();
        List<String> debtKeys = new ArrayList<>();
        for (Map.Entry<String, Double> e : net.entrySet()) {
            if (e.getValue() > 1e-9) {
                credKeys.add(e.getKey());
                creditors.add(new double[]{e.getValue()});
            } else if (e.getValue() < -1e-9) {
                debtKeys.add(e.getKey());
                debtors.add(new double[]{-e.getValue()});
            }
        }

        // Sort descending by amount keeping keys in sync
        Integer[] credIdx = new Integer[creditors.size()];
        for (int i = 0; i < credIdx.length; i++) credIdx[i] = i;
        Arrays.sort(credIdx, (a, b) -> Double.compare(creditors.get(b)[0], creditors.get(a)[0]));
        Integer[] debtIdx = new Integer[debtors.size()];
        for (int i = 0; i < debtIdx.length; i++) debtIdx[i] = i;
        Arrays.sort(debtIdx, (a, b) -> Double.compare(debtors.get(b)[0], debtors.get(a)[0]));

        List<String> sortedCredKeys = new ArrayList<>();
        List<Double> sortedCredVals = new ArrayList<>();
        for (int i : credIdx) { sortedCredKeys.add(credKeys.get(i)); sortedCredVals.add(creditors.get(i)[0]); }
        List<String> sortedDebtKeys = new ArrayList<>();
        List<Double> sortedDebtVals = new ArrayList<>();
        for (int i : debtIdx) { sortedDebtKeys.add(debtKeys.get(i)); sortedDebtVals.add(debtors.get(i)[0]); }

        List<DebtTxn> txns = new ArrayList<>();
        int i = 0, j = 0;
        while (i < sortedCredVals.size() && j < sortedDebtVals.size()) {
            double settle = Math.min(sortedCredVals.get(i), sortedDebtVals.get(j));
            txns.add(new DebtTxn(sortedDebtKeys.get(j), sortedCredKeys.get(i), settle));
            sortedCredVals.set(i, sortedCredVals.get(i) - settle);
            sortedDebtVals.set(j, sortedDebtVals.get(j) - settle);
            if (sortedCredVals.get(i) < 1e-9) i++;
            if (sortedDebtVals.get(j) < 1e-9) j++;
        }
        return txns;
    }
}

public class Solution {
    private static List<String> splitCsv(String s) {
        List<String> out = new ArrayList<>();
        if (s == null || s.isEmpty()) return out;
        StringBuilder cur = new StringBuilder();
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            if (c == ',') { out.add(cur.toString()); cur.setLength(0); }
            else cur.append(c);
        }
        out.add(cur.toString());
        return out;
    }

    private static List<Double> splitCsvDouble(String s) {
        List<Double> out = new ArrayList<>();
        for (String p : splitCsv(s)) out.add(Double.parseDouble(p));
        return out;
    }

    private static String num2(double v) {
        return String.format(java.util.Locale.US, "%.2f", v);
    }

    public static List<String> splitwise_simulate(List<SplitOp> ops) {
        List<String> out = new ArrayList<>();
        ExpenseManager mgr = new ExpenseManager();
        for (SplitOp op : ops) {
            String k = op.kind;
            if ("new".equals(k)) {
                mgr = new ExpenseManager();
                out.add("ok");
            } else if ("add_user".equals(k)) {
                mgr.addUser(op.s1, op.s2);
                out.add("ok");
            } else if ("add_expense".equals(k)) {
                mgr.addExpense(op.s1, op.s2, (double) op.i1, splitCsv(op.s3));
                out.add("ok");
            } else if ("add_eq_strat".equals(k)) {
                mgr.addExpenseWithStrategy(op.s1, op.s2, (double) op.i1, splitCsv(op.s3),
                        new EqualSplit(), new ArrayList<>());
                out.add("ok");
            } else if ("add_exact".equals(k)) {
                mgr.addExpenseWithStrategy(op.s1, op.s2, (double) op.i1, splitCsv(op.s3),
                        new ExactSplit(), splitCsvDouble(op.s4));
                out.add("ok");
            } else if ("add_pct".equals(k)) {
                mgr.addExpenseWithStrategy(op.s1, op.s2, (double) op.i1, splitCsv(op.s3),
                        new PercentSplit(), splitCsvDouble(op.s4));
                out.add("ok");
            } else if ("balance".equals(k)) {
                Map<String, Map<String, Double>> b = mgr.getBalances();
                double v = 0.0;
                Map<String, Double> inner = b.get(op.s1);
                if (inner != null && inner.containsKey(op.s2)) v = inner.get(op.s2);
                out.add(num2(v));
            } else if ("validate_exact".equals(k)) {
                ExactSplit ex = new ExactSplit();
                out.add(ex.validate((double) op.i1, splitCsv(op.s3), splitCsvDouble(op.s4)) ? "yes" : "no");
            } else if ("validate_pct".equals(k)) {
                PercentSplit pct = new PercentSplit();
                out.add(pct.validate((double) op.i1, splitCsv(op.s3), splitCsvDouble(op.s4)) ? "yes" : "no");
            } else if ("simplify_count".equals(k)) {
                out.add(Integer.toString(mgr.simplifyDebts().size()));
            } else if ("simplify_total_to".equals(k)) {
                double tot = 0;
                for (DebtTxn t : mgr.simplifyDebts()) if (t.creditor.equals(op.s1)) tot += t.amount;
                out.add(num2(tot));
            } else if ("simplify_total_from".equals(k)) {
                double tot = 0;
                for (DebtTxn t : mgr.simplifyDebts()) if (t.debtor.equals(op.s1)) tot += t.amount;
                out.add(num2(tot));
            } else if ("simplify_unique_pair".equals(k)) {
                double tot = 0;
                for (DebtTxn t : mgr.simplifyDebts())
                    if (t.debtor.equals(op.s1) && t.creditor.equals(op.s2)) tot += t.amount;
                out.add(num2(tot));
            } else {
                out.add("unknown:" + k);
            }
        }
        return out;
    }
}
