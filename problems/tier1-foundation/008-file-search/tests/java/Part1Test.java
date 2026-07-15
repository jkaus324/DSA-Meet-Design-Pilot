// File Search — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testSearchByExtensionCpp() {
        try {
            var results = search_by_extension( root, "cpp");
            // Should find main.cpp and utils.cpp (order depends on DFS)
            boolean foundMain = false, foundUtils = false;
            for (var* f : results) {
            if (f.name == "main.cpp") foundMain = true;
            if (f.name == "utils.cpp") foundUtils = true;
            }
            boolean pass = results.size() == 2
                && foundMain & foundUtils;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSearchByExtensionCpp");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSearchByExtensionCpp (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSearchByExtensionPdf() {
        try {
            var results = search_by_extension( root, "pdf");
            boolean pass = results.size() == 1
                && results[0].name == "report.pdf";
            System.out.println((pass ? "PASS" : "FAIL") + ": testSearchByExtensionPdf");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSearchByExtensionPdf (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSearchByExtensionNoMatch() {
        try {
            var results = search_by_extension( root, "java");
            boolean pass = results.size() == 0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSearchByExtensionNoMatch");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSearchByExtensionNoMatch (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSearchBySize100() {
        try {
            var results = search_by_size( root, 100);
            boolean foundUtils = false, foundReport = false;
            for (var* f : results) {
            if (f.name == "utils.cpp") foundUtils = true;
            if (f.name == "report.pdf") foundReport = true;
            }
            boolean pass = results.size() == 2
                && foundUtils & foundReport;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSearchBySize100");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSearchBySize100 (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSearchBySizeExact() {
        try {
            var results = search_by_size( root, 200);
            boolean pass = results.size() == 1
                && results[0].name == "report.pdf";
            System.out.println((pass ? "PASS" : "FAIL") + ": testSearchBySizeExact");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSearchBySizeExact (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSearchByNameMain() {
        try {
            var results = search_by_name( root, "main");
            boolean pass = results.size() == 1
                && results[0].name == "main.cpp";
            System.out.println((pass ? "PASS" : "FAIL") + ": testSearchByNameMain");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSearchByNameMain (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSearchByNameReport() {
        try {
            var results = search_by_name( root, "report");
            boolean pass = results.size() == 1
                && results[0].name == "report.pdf";
            System.out.println((pass ? "PASS" : "FAIL") + ": testSearchByNameReport");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSearchByNameReport (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEmptyTree() {
        try {
            FileNode emptyRoot{"empty", 0, "", true, {}};
            boolean pass = search_by_extension( emptyRoot, "cpp").isEmpty()
                && search_by_size( emptyRoot, 10).isEmpty()
                && search_by_name( emptyRoot, "test").isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testEmptyTree");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEmptyTree (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSingleFileTree() {
        try {
            FileNode singleDir{"root", 0, "", true, {}};
            FileNode singleFile{"test.txt", 30, "txt", false, {}};
            singleDir.children.add( singleFile);
            var results = search_by_extension( singleDir, "txt");
            boolean pass = results.size() == 1
                && results[0].name == "test.txt";
            System.out.println((pass ? "PASS" : "FAIL") + ": testSingleFileTree");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSingleFileTree (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testSearchByExtensionCpp()) passed++;
        total++; if (testSearchByExtensionPdf()) passed++;
        total++; if (testSearchByExtensionNoMatch()) passed++;
        total++; if (testSearchBySize100()) passed++;
        total++; if (testSearchBySizeExact()) passed++;
        total++; if (testSearchByNameMain()) passed++;
        total++; if (testSearchByNameReport()) passed++;
        total++; if (testEmptyTree()) passed++;
        total++; if (testSingleFileTree()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
