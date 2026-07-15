// Lru Cache — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testBasicPutGet() {
        try {
            LRUCache cache = new LRUCache(2);
            cache.put(1, 10);
            cache.put(2, 20);
            boolean pass = cache.get(1) == 10
                && cache.get(2) == 20;
            System.out.println((pass ? "PASS" : "FAIL") + ": testBasicPutGet");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testBasicPutGet (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testGetNonexistent() {
        try {
            LRUCache cache = new LRUCache(2);
            cache.put(1, 10);
            boolean pass = cache.get(99) == -1
                && cache.get(2) == -1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testGetNonexistent");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testGetNonexistent (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEvictionOnCapacity() {
        try {
            LRUCache cache = new LRUCache(2);
            cache.put(1, 10);
            cache.put(2, 20);
            cache.put(3, 30);  // evicts key 1 (LRU)
            boolean pass = cache.put(1, = -1));  // evicted
                && cache.get(2) == 20
                && cache.get(3) == 30;
            System.out.println((pass ? "PASS" : "FAIL") + ": testEvictionOnCapacity");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEvictionOnCapacity (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testGetUpdatesRecency() {
        try {
            LRUCache cache = new LRUCache(2);
            cache.put(1, 10);
            cache.put(2, 20);
            cache.get(1);       // key 1 is now most recent
            cache.put(3, 30);   // evicts key 2 (now LRU)
            boolean pass = cache.put(1, = 10));  // still present
                && cache.put(2, = -1));  // evicted
                && cache.get(3) == 30;
            System.out.println((pass ? "PASS" : "FAIL") + ": testGetUpdatesRecency");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testGetUpdatesRecency (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPutUpdatesValue() {
        try {
            LRUCache cache = new LRUCache(2);
            cache.put(1, 10);
            cache.put(1, 100);  // update value
            boolean pass = cache.get(1) == 100;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPutUpdatesValue");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPutUpdatesValue (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPutUpdatesRecency() {
        try {
            LRUCache cache = new LRUCache(2);
            cache.put(1, 10);
            cache.put(2, 20);
            cache.put(1, 100);  // key 1 updated — now most recent
            cache.put(3, 30);   // evicts key 2 (now LRU)
            boolean pass = cache.get(1) == 100
                && cache.put(2, = -1));  // evicted
                && cache.get(3) == 30;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPutUpdatesRecency");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPutUpdatesRecency (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCapacityOne() {
        try {
            LRUCache cache = new LRUCache(1);
            cache.put(1, 10);
            cache.put(2, 20);   // evicts key 1
            boolean pass = cache.get(1) == 10
                && cache.get(1) == -1
                && cache.get(2) == 20;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCapacityOne");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCapacityOne (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMultipleEvictions() {
        try {
            LRUCache cache = new LRUCache(3);
            cache.put(1, 10);
            cache.put(2, 20);
            cache.put(3, 30);
            cache.put(4, 40);  // evicts 1
            cache.put(5, 50);  // evicts 2
            boolean pass = cache.get(1) == -1
                && cache.get(2) == -1
                && cache.get(3) == 30
                && cache.get(4) == 40
                && cache.get(5) == 50;
            System.out.println((pass ? "PASS" : "FAIL") + ": testMultipleEvictions");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMultipleEvictions (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testBasicPutGet()) passed++;
        total++; if (testGetNonexistent()) passed++;
        total++; if (testEvictionOnCapacity()) passed++;
        total++; if (testGetUpdatesRecency()) passed++;
        total++; if (testPutUpdatesValue()) passed++;
        total++; if (testPutUpdatesRecency()) passed++;
        total++; if (testCapacityOne()) passed++;
        total++; if (testMultipleEvictions()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
