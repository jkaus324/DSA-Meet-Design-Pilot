import java.util.*;

// Data class (given).
class LogOp {
    public String kind;
    public String s1;
    public String s2;
    public int i1;

    public LogOp(String kind, String s1, String s2, int i1) {
        this.kind = kind;
        this.s1 = s1;
        this.s2 = s2;
        this.i1 = i1;
    }

    public LogOp(String kind) {
        this(kind, "", "", 0);
    }
}

// HINT: introduce an abstraction so new ranking rules don't change existing code.
public class Solution {
    // HINT: pick the field that defines 'better' for this ranking and compare the two.
    public static List<String> logger_simulate(List<LogOp> ops) {
        // TODO: write your solution
        return null;
    }

}
