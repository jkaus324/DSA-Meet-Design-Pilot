import java.util.*;

// Data class (given).
class LruOp {
    public String kind;
    public int i1;
    public int i2;
    public int i3;
    public int i4;

    public LruOp(String kind, int i1, int i2, int i3, int i4) {
        this.kind = kind;
        this.i1 = i1;
        this.i2 = i2;
        this.i3 = i3;
        this.i4 = i4;
    }

    public LruOp(String kind) {
        this(kind, 0, 0, 0, 0);
    }
}

// HINT: introduce an abstraction so new ranking rules don't change existing code.
public class Solution {
    // HINT: pick the field that defines 'better' for this ranking and compare the two.
    public static List<String> lru_simulate(List<LruOp> ops) {
        // TODO: write your solution
        return null;
    }

}
