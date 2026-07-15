// File Search — Part 3 Tests
import java.util.*;
import java.util.stream.*;

class Part3Test {
    static boolean testSortByName() {
        try {
            SearchByMinSize criteria = new SearchByMinSize(1);
            var results = search_and_sort( root, criteria, "name");
            boolean pass = results.size() == 6
                && results[0].name == "build.sh"
                && results[1].name == "helper.h"
                && results[2].name == "main.cpp"
                && results[3].name == "readme.md"
                && results[4].name == "report.pdf"
                && results[5].name == "utils.cpp";
            System.out.println((pass ? "PASS" : "FAIL") + ": testSortByName");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSortByName (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSortBySize() {
        try {
            SearchByMinSize criteria = new SearchByMinSize(1);
            var results = search_and_sort( root, criteria, "size");
            boolean pass = results.size() == 6
                && results[0].name == "report.pdf");   // 200 KB
                && results[1].name == "utils.cpp");     // 120 KB
                && results[2].name == "main.cpp");      // 50 KB
                && results[3].name == "helper.h");      // 10 KB
                && results[4].name == "readme.md");     // 5 KB
                && results[5].name == "build.sh");      // 2 KB;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSortBySize");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSortBySize (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSearchCppSortBySize() {
        try {
            SearchByExtension criteria = new SearchByExtension("cpp");
            var results = search_and_sort( root, criteria, "size");
            boolean pass = results.size() == 2
                && results[0].name == "utils.cpp");  // 120 KB
                && results[1].name == "main.cpp");   // 50 KB;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSearchCppSortBySize");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSearchCppSortBySize (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSortByExtension() {
        try {
            SearchByMinSize criteria = new SearchByMinSize(1);
            var results = search_and_sort( root, criteria, "extension");
            // Extensions in order: cpp, cpp, h, md, pdf, sh
            boolean pass = results.size() == 6
                && results[0].extension == "cpp"
                && results[1].extension == "cpp"
                && results[2].extension == "h"
                && results[3].extension == "md"
                && results[4].extension == "pdf"
                && results[5].extension == "sh";
            System.out.println((pass ? "PASS" : "FAIL") + ": testSortByExtension");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSortByExtension (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSortEmptyResults() {
        try {
            SearchByExtension criteria = new SearchByExtension("java");
            var results = search_and_sort( root, criteria, "name");
            boolean pass = results.size() == 0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSortEmptyResults");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSortEmptyResults (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testSortByName()) passed++;
        total++; if (testSortBySize()) passed++;
        total++; if (testSearchCppSortBySize()) passed++;
        total++; if (testSortByExtension()) passed++;
        total++; if (testSortEmptyResults()) passed++;
        System.out.println("PART3_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
