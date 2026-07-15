import java.util.*;

// Data class (given).
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

    public SplitOp(String kind) {
        this(kind, "", "", "", "", 0);
    }
}

// HINT: introduce an abstraction so new ranking rules don't change existing code.
public class Solution {
    // HINT: pick the field that defines 'better' for this ranking and compare the two.
    public static List<String> splitwise_simulate(List<SplitOp> ops) {
        // TODO: write your solution
        return null;
    }

}
