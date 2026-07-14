// Lru Cache — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testTtlExpiry() {
        try {
            LRUCache cache = new LRUCache(5);
            cache.put(1, 10, 1000, 60);  // expires at 1060
            boolean pass = cache.put(1, 1030, = 10));  // still valid
                && cache.put(1, 1061, = -1));  // expired;
            System.out.println((pass ? "PASS" : "FAIL") + ": testTtlExpiry");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testTtlExpiry (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testNoTtlNeverExpires() {
        try {
            LRUCache cache = new LRUCache(5);
            cache.put(1, 10, 1000, 0);  // no TTL
            boolean pass = cache.put(1, 999999, = 10));  // never expires;
            System.out.println((pass ? "PASS" : "FAIL") + ": testNoTtlNeverExpires");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testNoTtlNeverExpires (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testDeleteExisting() {
        try {
            LRUCache cache = new LRUCache(5);
            cache.put(1, 10, 1000);
            boolean pass = cache.deleteKey(1) == true
                && cache.get(1, 1001) == -1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testDeleteExisting");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDeleteExisting (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testDeleteNonexistent() {
        try {
            LRUCache cache = new LRUCache(5);
            boolean pass = cache.deleteKey(99) == false;
            System.out.println((pass ? "PASS" : "FAIL") + ": testDeleteNonexistent");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testDeleteNonexistent (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSizeTracking() {
        try {
            LRUCache cache = new LRUCache(5);
            cache.put(1, 10, 1000);
            cache.put(2, 20, 1000);
            cache.deleteKey(1);
            boolean pass = cache.size() == 0
                && cache.size() == 2
                && cache.size() == 1;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSizeTracking");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSizeTracking (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testLazyExpiryOnGet() {
        try {
            LRUCache cache = new LRUCache(5);
            cache.put(1, 10, 1000, 30);  // expires at 1030
            cache.put(2, 20, 1000, 60);  // expires at 1060
            // Access after key 1 expired but before key 2
            boolean pass = cache.put(1, 1035, = -1));  // expired, removed
                && cache.put(2, 1035, = 20));  // still valid;
            System.out.println((pass ? "PASS" : "FAIL") + ": testLazyExpiryOnGet");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testLazyExpiryOnGet (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testPutRefreshesTtl() {
        try {
            LRUCache cache = new LRUCache(5);
            cache.put(1, 10, 1000, 30);  // expires at 1030
            cache.put(1, 20, 1020, 30);  // new TTL, expires at 1050
            boolean pass = cache.put(1, 1035, = 20));  // still valid (new TTL
                && cache.put(1, 1051, = -1));  // now expired;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPutRefreshesTtl");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPutRefreshesTtl (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testCapacityWithExpired() {
        try {
            LRUCache cache = new LRUCache(2);
            cache.put(1, 10, 1000, 10);  // expires at 1010
            cache.put(2, 20, 1000, 10);  // expires at 1010
            // Both expired — put new entries
            cache.put(3, 30, 1020);  // should evict expired key 1
            cache.put(4, 40, 1020);  // should evict expired key 2
            boolean pass = cache.get(3, 1020) == 30
                && cache.get(4, 1020) == 40;
            System.out.println((pass ? "PASS" : "FAIL") + ": testCapacityWithExpired");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testCapacityWithExpired (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testTtlExpiry()) passed++;
        total++; if (testNoTtlNeverExpires()) passed++;
        total++; if (testDeleteExisting()) passed++;
        total++; if (testDeleteNonexistent()) passed++;
        total++; if (testSizeTracking()) passed++;
        total++; if (testLazyExpiryOnGet()) passed++;
        total++; if (testPutRefreshesTtl()) passed++;
        total++; if (testCapacityWithExpired()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
