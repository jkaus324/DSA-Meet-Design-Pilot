import java.util.*;

// Data class (given).
class ParkOp {
    public String kind;
    public String s1;
    public String s2;
    public String s3;
    public int i1;
    public int i2;
    public int i3;

    public ParkOp(String kind, String s1, String s2, String s3, int i1, int i2, int i3) {
        this.kind = kind;
        this.s1 = s1;
        this.s2 = s2;
        this.s3 = s3;
        this.i1 = i1;
        this.i2 = i2;
        this.i3 = i3;
    }

    public ParkOp(String kind) {
        this(kind, "", "", "", 0, 0, 0);
    }
}

public class Solution {
    public static List<String> parking_simulate(List<ParkOp> ops) {
        // TODO: write your solution
        return null;
    }

}
