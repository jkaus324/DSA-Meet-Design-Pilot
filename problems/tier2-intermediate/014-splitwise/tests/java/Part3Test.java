// Splitwise — Part 3 Tests
import java.util.*;
import java.util.stream.*;

class Part3Test {
    static boolean testSimpleTwoPerson() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            mgr.addUser("bob", "Bob");
            mgr.addExpense("E1", "alice", 200.0, Arrays.asList("alice", "bob"));
            // Bob owes Alice 100
            var txns = mgr.simplifyDebts();
            boolean pass = txns.size() == 1
                && get<0>(txns[0]) == "bob"
                && get<1>(txns[0]) == "alice"
                && Math.abs(get<2>(txns[0]) - 100.0) < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSimpleTwoPerson");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSimpleTwoPerson (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testThreePersonSimplification() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("A", "A");
            mgr.addUser("B", "B");
            mgr.addUser("C", "C");
            // A pays $300 for A,B,C => B owes A 100, C owes A 100
            mgr.addExpense("E1", "A", 300.0, Arrays.asList("A", "B", "C"));
            var txns = mgr.simplifyDebts();
            // Net: A is owed 200, B owes 100, C owes 100
            // Simplified: B.A 100, C.A 100 (already minimal)
            // Verify total flow to A is 200
            double totalToA = 0;
            for (var t : txns) {
            if (get<1>(t) == "A") totalToA += get<2>(t);
            }
            boolean pass = txns.size() == 2
                && Math.abs(totalToA - 200.0) < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testThreePersonSimplification");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testThreePersonSimplification (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCircularDebtSimplification() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("A", "A");
            mgr.addUser("B", "B");
            mgr.addUser("C", "C");
            // A pays $300 for A,B,C => B owes A 100, C owes A 100
            mgr.addExpense("E1", "A", 300.0, Arrays.asList("A", "B", "C"));
            // B pays $300 for A,B,C => A owes B 100, C owes B 100
            mgr.addExpense("E2", "B", 300.0, Arrays.asList("A", "B", "C"));
            // Net: A owed 200, owes 100 => net +100
            //      B owed 200, owes 100 => net +100
            //      C owes 200 => net -200
            // Simplified: C pays A 100, C pays B 100
            var txns = mgr.simplifyDebts();
            double totalFromC = 0;
            for (var t : txns) {
            if (get<0>(t) == "C") totalFromC += get<2>(t);
            }
            boolean pass = txns.size() == 2
                && Math.abs(totalFromC - 200.0) < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCircularDebtSimplification");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCircularDebtSimplification (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNoDebtsEmpty() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            var txns = mgr.simplifyDebts();
            boolean pass = txns.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testNoDebtsEmpty");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNoDebtsEmpty (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCancellingDebts() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("alice", "Alice");
            mgr.addUser("bob", "Bob");
            mgr.addExpense("E1", "alice", 200.0, Arrays.asList("alice", "bob"));
            mgr.addExpense("E2", "bob", 200.0, Arrays.asList("alice", "bob"));
            // Bob owes Alice 100, Alice owes Bob 100 => net zero
            var txns = mgr.simplifyDebts();
            boolean pass = txns.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testCancellingDebts");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCancellingDebts (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFourPersonSimplification() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("A", "A");
            mgr.addUser("B", "B");
            mgr.addUser("C", "C");
            mgr.addUser("D", "D");
            // A pays 400 for all => B,C,D each owe A 100
            mgr.addExpense("E1", "A", 400.0, Arrays.asList("A", "B", "C", "D"));
            // Net: A=+300, B=-100, C=-100, D=-100
            var txns = mgr.simplifyDebts();
            double totalToA = 0;
            for (var t : txns) {
            totalToA += get<2>(t);
            }
            boolean pass = txns.size() == 3
                && get<1>(t) == "A");  // All payments go to A
                && Math.abs(totalToA - 300.0) < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFourPersonSimplification");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFourPersonSimplification (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSimplificationPreservesNet() {
        try {
            ExpenseManager mgr = new ExpenseManager();
            mgr.addUser("A", "A");
            mgr.addUser("B", "B");
            mgr.addUser("C", "C");
            mgr.addExpense("E1", "A", 600.0, Arrays.asList("A", "B", "C"));
            mgr.addExpense("E2", "B", 300.0, Arrays.asList("A", "B", "C"));
            // After E1: B.A:200, C.A:200
            // After E2: A.B:100 (netted: B.A:100), C.B:100
            // Net: A=+300, B=-100+200=+100, hmm let me recalc
            // E1: A pays 600, split 200 each. B owes A 200, C owes A 200.
            // E2: B pays 300, split 100 each. A owes B 100, C owes B 100.
            // After netting: B owes A 100, C owes A 200, C owes B 100
            // Net: A=+300, B=+0(owes100 to A, owed 100 by C), C=-300
            // Actually: Net A = +200-100=+100+200 = let me use formula
            // A net = (200+200) - (100) = +300
            // B net = (100) - (200) = actually...
            // B net = (owed to B: 100 from C) - (B owes A: 100) = 0
            // C net = 0 - (200 to A + 100 to B) = -300
            // Simplified: C.A:300 (if B net is 0)
            // But wait, B might have some transactions. Let me check the greedy.
            // creditors: A(+300), debtors: C(300). B is zero.
            // One transaction: C.A 300
            var txns = mgr.simplifyDebts();
            // B's net should be zero after netting
            double totalFromC = 0, totalToA = 0;
            for (var t : txns) {
            if (get<0>(t) == "C") totalFromC += get<2>(t);
            if (get<1>(t) == "A") totalToA += get<2>(t);
            }
            boolean pass = Math.abs(totalFromC - 300.0) < 1e-9
                && Math.abs(totalToA - 300.0) < 1e-9;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSimplificationPreservesNet");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSimplificationPreservesNet (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testSimpleTwoPerson()) passed++;
        total++; if (testThreePersonSimplification()) passed++;
        total++; if (testCircularDebtSimplification()) passed++;
        total++; if (testNoDebtsEmpty()) passed++;
        total++; if (testCancellingDebts()) passed++;
        total++; if (testFourPersonSimplification()) passed++;
        total++; if (testSimplificationPreservesNet()) passed++;
        System.out.println("PART3_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
