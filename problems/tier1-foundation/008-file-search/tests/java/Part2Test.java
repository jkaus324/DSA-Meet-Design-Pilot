// File Search — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testAndFilterCppAndLarge() {
        try {
            SearchByExtension extCriteria = new SearchByExtension("cpp");
            SearchByMinSize sizeCriteria = new SearchByMinSize(100);
            var results = search_composite( root, { extCriteria,  sizeCriteria}, "AND");
            boolean pass = results.size() == 1
                && results[0].name == "utils.cpp"); // only .cpp file >= 100 KB;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAndFilterCppAndLarge");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAndFilterCppAndLarge (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testOrFilterPdfOrMain() {
        try {
            SearchByExtension extCriteria = new SearchByExtension("pdf");
            SearchByName nameCriteria = new SearchByName("main");
            var results = search_composite( root, { extCriteria,  nameCriteria}, "OR");
            boolean foundMain = false, foundReport = false;
            for (var* f : results) {
            if (f.name == "main.cpp") foundMain = true;
            if (f.name == "report.pdf") foundReport = true;
            }
            boolean pass = results.size() == 2
                && foundMain & foundReport;
            System.out.println((pass ? "PASS" : "FAIL") + ": testOrFilterPdfOrMain");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testOrFilterPdfOrMain (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAndFilterNoMatch() {
        try {
            SearchByExtension extCriteria = new SearchByExtension("h");
            SearchByMinSize sizeCriteria = new SearchByMinSize(50);
            var results = search_composite( root, { extCriteria,  sizeCriteria}, "AND");
            boolean pass = results.size() == 0); // helper.h is only 10 KB;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAndFilterNoMatch");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAndFilterNoMatch (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSingleCriterionOr() {
        try {
            SearchByExtension extCriteria = new SearchByExtension("cpp");
            var composite = search_composite( root, { extCriteria}, "OR");
            var direct = search_by_extension( root, "cpp");
            boolean pass = composite.size() == direct.size();
            System.out.println((pass ? "PASS" : "FAIL") + ": testSingleCriterionOr");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSingleCriterionOr (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testAndFilterSizeAndDot() {
        try {
            SearchByMinSize sizeCriteria = new SearchByMinSize(50);
            SearchByName nameCriteria = new SearchByName(".");
            var results = search_composite( root, { sizeCriteria,  nameCriteria}, "AND");
            // Files >= 50KB with '.' in name: main.cpp(50), utils.cpp(120), report.pdf(200)
            boolean pass = results.size() == 3;
            System.out.println((pass ? "PASS" : "FAIL") + ": testAndFilterSizeAndDot");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testAndFilterSizeAndDot (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testAndFilterCppAndLarge()) passed++;
        total++; if (testOrFilterPdfOrMain()) passed++;
        total++; if (testAndFilterNoMatch()) passed++;
        total++; if (testSingleCriterionOr()) passed++;
        total++; if (testAndFilterSizeAndDot()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
