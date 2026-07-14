// Splitwise — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testSimpleEqualSplit() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            mgr.addUser("bob", "Bob");
            mgr.addUser("charlie", "Charlie");
            mgr.addExpense("E1", "alice", 300.0, Arrays.asList("alice", "bob", "charlie"));
            var bal = mgr.getBalances();
            // Alice should not owe herself
            boolean pass = Math.abs(bal.get("bob").get("alice") - 100.0) < 1e-9
                && Math.abs(bal.get("charlie").get("alice") - 100.0) < 1e-9
                && !bal.get("alice").containsKey("alice") || bal.get("alice").get("alice") < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSimpleEqualSplit");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSimpleEqualSplit (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAccumulatedBalances() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            mgr.addUser("bob", "Bob");
            mgr.addExpense("E1", "alice", 200.0, Arrays.asList("alice", "bob"));
            mgr.addExpense("E2", "alice", 100.0, Arrays.asList("alice", "bob"));
            var bal = mgr.getBalances();
            // Bob owes Alice: 100 + 50 = 150
            boolean pass = Math.abs(bal.get("bob").get("alice") - 150.0) < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAccumulatedBalances");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAccumulatedBalances (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMutualDebtNetting() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            mgr.addUser("bob", "Bob");
            mgr.addExpense("E1", "alice", 200.0, Arrays.asList("alice", "bob"));
            // Bob owes Alice $100
            mgr.addExpense("E2", "bob", 60.0, Arrays.asList("alice", "bob"));
            // Alice owes Bob $30, but Bob already owes Alice $100
            // Net: Bob owes Alice $70
            var bal = mgr.getBalances();
            // Alice should NOT owe Bob anything
            boolean pass = Math.abs(bal.get("bob").get("alice") - 70.0) < 1e-9
                && !bal.get("alice").containsKey("bob") || bal.get("alice").get("bob") < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testMutualDebtNetting");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMutualDebtNetting (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPayerNotInParticipants() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            mgr.addUser("bob", "Bob");
            mgr.addUser("charlie", "Charlie");
            mgr.addExpense("E1", "alice", 200.0, Arrays.asList("bob", "charlie"));
            var bal = mgr.getBalances();
            boolean pass = Math.abs(bal.get("bob").get("alice") - 100.0) < 1e-9
                && Math.abs(bal.get("charlie").get("alice") - 100.0) < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPayerNotInParticipants");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPayerNotInParticipants (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSelfExpenseNoDebt() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            mgr.addExpense("E1", "alice", 100.0, Arrays.asList("alice"));
            var bal = mgr.getBalances();
            boolean noDebts = true;
            for (var _e_bal_ : bal.entrySet()) {
            var debtor = _e_bal_.getKey(); var creditors = _e_bal_.getValue();
            for (var _e_creditors_ : creditors.entrySet()) {
            var creditor = _e_creditors_.getKey(); var amount = _e_creditors_.getValue();
            if (amount > 1e-9) noDebts = false;
            }
            }
            boolean pass = noDebts;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSelfExpenseNoDebt");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSelfExpenseNoDebt (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMultiUserMultiExpense() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("A", "A");
            mgr.addUser("B", "B");
            mgr.addUser("C", "C");
            mgr.addUser("D", "D");
            mgr.addExpense("E1", "A", 400.0, Arrays.asList("A", "B", "C", "D"));
            // Each owes A: $100. B.A:100, C.A:100, D.A:100
            mgr.addExpense("E2", "B", 200.0, Arrays.asList("A", "B", "C", "D"));
            // Each owes B: $50. A.B:50, C.B:50, D.B:50
            // But B owes A 100, and A owes B 50 => net: B owes A 50
            var bal = mgr.getBalances();
            boolean pass = Math.abs(bal.get("B").get("A") - 50.0) < 1e-9
                && Math.abs(bal.get("C").get("A") - 100.0) < 1e-9
                && Math.abs(bal.get("C").get("B") - 50.0) < 1e-9
                && Math.abs(bal.get("D").get("A") - 100.0) < 1e-9
                && Math.abs(bal.get("D").get("B") - 50.0) < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testMultiUserMultiExpense");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMultiUserMultiExpense (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEmptyBalances() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            mgr.addUser("bob", "Bob");
            var bal = mgr.getBalances();
            boolean empty = true;
            for (var _e_bal_ : bal.entrySet()) {
            var k = _e_bal_.getKey(); var v = _e_bal_.getValue();
            for (var _e_v_ : v.entrySet()) {
            var k2 = _e_v_.getKey(); var amt = _e_v_.getValue();
            if (amt > 1e-9) empty = false;
            }
            }
            boolean pass = empty;
            System.out.println((pass ? "PASS" : "FAIL") + ": testEmptyBalances");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEmptyBalances (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testSimpleEqualSplit()) passed++;
        total++; if (testAccumulatedBalances()) passed++;
        total++; if (testMutualDebtNetting()) passed++;
        total++; if (testPayerNotInParticipants()) passed++;
        total++; if (testSelfExpenseNoDebt()) passed++;
        total++; if (testMultiUserMultiExpense()) passed++;
        total++; if (testEmptyBalances()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
