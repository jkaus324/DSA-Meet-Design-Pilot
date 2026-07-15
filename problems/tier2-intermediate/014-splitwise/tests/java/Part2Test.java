// Splitwise — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testEqualSplitStrategy() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            mgr.addUser("bob", "Bob");
            mgr.addUser("charlie", "Charlie");
            EqualSplit eq = new EqualSplit();
            mgr.addExpenseWithStrategy("E1", "alice", 300.0,
            Arrays.asList("alice", "bob", "charlie"),  eq, {});
            var bal = mgr.getBalances();
            boolean pass = Math.abs(bal.get("bob").get("alice") - 100.0) < 1e-9
                && Math.abs(bal.get("charlie").get("alice") - 100.0) < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testEqualSplitStrategy");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEqualSplitStrategy (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testExactSplitStrategy() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            mgr.addUser("bob", "Bob");
            mgr.addUser("charlie", "Charlie");
            ExactSplit exact = new ExactSplit();
            mgr.addExpenseWithStrategy("E1", "alice", 300.0,
            Arrays.asList("alice", "bob", "charlie"),  exact, Arrays.asList(100.0, 150.0, 50.0));
            var bal = mgr.getBalances();
            // Bob owes Alice 150, Charlie owes Alice 50
            boolean pass = Math.abs(bal.get("bob").get("alice") - 150.0) < 1e-9
                && Math.abs(bal.get("charlie").get("alice") - 50.0) < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testExactSplitStrategy");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testExactSplitStrategy (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPercentSplitStrategy() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            mgr.addUser("bob", "Bob");
            mgr.addUser("charlie", "Charlie");
            PercentSplit pct = new PercentSplit();
            mgr.addExpenseWithStrategy("E1", "alice", 1000.0,
            Arrays.asList("alice", "bob", "charlie"),  pct, Arrays.asList(50.0, 30.0, 20.0));
            var bal = mgr.getBalances();
            // Alice's share: 500 (cancels), Bob: 300, Charlie: 200
            boolean pass = Math.abs(bal.get("bob").get("alice") - 300.0) < 1e-9
                && Math.abs(bal.get("charlie").get("alice") - 200.0) < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPercentSplitStrategy");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPercentSplitStrategy (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testExactSplitInvalidSum() {
        try {
            ExactSplit exact = new ExactSplit();
            boolean valid = exact.validate(300.0, Arrays.asList("a", "b", "c"), Arrays.asList(100.0, 100.0, 50.0));
            boolean pass = valid == false);  // 250 != 300;
            System.out.println((pass ? "PASS" : "FAIL") + ": testExactSplitInvalidSum");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testExactSplitInvalidSum (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPercentSplitInvalidSum() {
        try {
            PercentSplit pct = new PercentSplit();
            boolean valid = pct.validate(1000.0, Arrays.asList("a", "b"), Arrays.asList(60.0, 60.0));
            boolean pass = valid == false);  // 120 != 100;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPercentSplitInvalidSum");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPercentSplitInvalidSum (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testExactSplitWrongParamCount() {
        try {
            ExactSplit exact = new ExactSplit();
            boolean valid = exact.validate(300.0, Arrays.asList("a", "b", "c"), Arrays.asList(150.0, 150.0));
            boolean pass = valid == false);  // 2 params for 3 participants;
            System.out.println((pass ? "PASS" : "FAIL") + ": testExactSplitWrongParamCount");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testExactSplitWrongParamCount (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMixedStrategies() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            mgr.addUser("bob", "Bob");
            EqualSplit eq = new EqualSplit();
            ExactSplit exact = new ExactSplit();
            mgr.addExpenseWithStrategy("E1", "alice", 200.0,
            Arrays.asList("alice", "bob"),  eq, {});
            // Bob owes Alice 100
            mgr.addExpenseWithStrategy("E2", "bob", 150.0,
            Arrays.asList("alice", "bob"),  exact, Arrays.asList(90.0, 60.0));
            // Alice owes Bob 90, but Bob owes Alice 100 => net Bob owes Alice 10
            var bal = mgr.getBalances();
            boolean pass = Math.abs(bal.get("bob").get("alice") - 10.0) < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testMixedStrategies");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMixedStrategies (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testInvalidStrategyNoExpense() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            mgr.addUser("bob", "Bob");
            ExactSplit exact = new ExactSplit();
            mgr.addExpenseWithStrategy("E1", "alice", 100.0,
            Arrays.asList("alice", "bob"),  exact, Arrays.asList(60.0, 60.0));  // sum=120, invalid
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
            System.out.println((pass ? "PASS" : "FAIL") + ": testInvalidStrategyNoExpense");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testInvalidStrategyNoExpense (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testEqualSplitStrategy()) passed++;
        total++; if (testExactSplitStrategy()) passed++;
        total++; if (testPercentSplitStrategy()) passed++;
        total++; if (testExactSplitInvalidSum()) passed++;
        total++; if (testPercentSplitInvalidSum()) passed++;
        total++; if (testExactSplitWrongParamCount()) passed++;
        total++; if (testMixedStrategies()) passed++;
        total++; if (testInvalidStrategyNoExpense()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
