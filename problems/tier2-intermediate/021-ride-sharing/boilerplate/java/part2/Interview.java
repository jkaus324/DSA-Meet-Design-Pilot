import java.util.*;

// Data class (given).
class RideOp {
    public String kind;
    public String s1;
    public String s2;
    public String s3;
    public String s4;
    public int i1;
    public int i2;

    public RideOp(String kind, String s1, String s2, String s3, String s4, int i1, int i2) {
        this.kind = kind;
        this.s1 = s1;
        this.s2 = s2;
        this.s3 = s3;
        this.s4 = s4;
        this.i1 = i1;
        this.i2 = i2;
    }

    public RideOp(String kind) {
        this(kind, "", "", "", "", 0, 0);
    }
}

public class Solution {
    public static List<String> ride_simulate(List<RideOp> ops) {
        // TODO: write your solution
        return null;
    }

}
