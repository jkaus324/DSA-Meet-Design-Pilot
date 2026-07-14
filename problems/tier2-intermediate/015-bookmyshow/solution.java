// BookMyShow — Solution (Java)
import java.util.*;

class ShowOp {
    public String kind;
    public String s1;
    public String s2;
    public String s3;
    public String s4;
    public int i1;
    public int i2;
    public int i3;

    public ShowOp(String kind, String s1, String s2, String s3, String s4, int i1, int i2, int i3) {
        this.kind = kind;
        this.s1 = s1;
        this.s2 = s2;
        this.s3 = s3;
        this.s4 = s4;
        this.i1 = i1;
        this.i2 = i2;
        this.i3 = i3;
    }
}

enum SeatStatus { AVAILABLE, LOCKED, BOOKED }

class Seat {
    public int row, col;
    public SeatStatus status;
    public String lockedBy;
    public long lockExpiry;
    public String bookedBy;

    public Seat(int row, int col) {
        this.row = row;
        this.col = col;
        this.status = SeatStatus.AVAILABLE;
        this.lockedBy = "";
        this.lockExpiry = 0;
        this.bookedBy = "";
    }
}

class Show {
    public String id;
    public String theaterId;
    public String movie;
    public String time;
    public int rows;
    public int cols;
    public Seat[][] seats;

    public Show(String id, String theaterId, String movie, String time, int rows, int cols) {
        this.id = id;
        this.theaterId = theaterId;
        this.movie = movie;
        this.time = time;
        this.rows = rows;
        this.cols = cols;
        this.seats = new Seat[rows][cols];
        for (int r = 0; r < rows; r++)
            for (int c = 0; c < cols; c++)
                this.seats[r][c] = new Seat(r, c);
    }
}

class Theater {
    public String id;
    public String name;
    public String city;
    public Theater(String id, String name, String city) {
        this.id = id; this.name = name; this.city = city;
    }
}

class Booking {
    public String id;
    public String showId;
    public String userId;
    public List<int[]> seatPositions;
    public Booking(String id, String showId, String userId, List<int[]> sp) {
        this.id = id; this.showId = showId; this.userId = userId; this.seatPositions = sp;
    }
}

class SeatLock {
    public String id;
    public String showId;
    public String userId;
    public List<int[]> seatPositions;
    public long expiry;
    public boolean confirmed;
    public boolean released;

    public SeatLock(String id, String showId, String userId, List<int[]> sp, long expiry) {
        this.id = id; this.showId = showId; this.userId = userId;
        this.seatPositions = sp; this.expiry = expiry;
        this.confirmed = false; this.released = false;
    }
}

class BookingSystem {
    Map<String, Theater> theaters = new LinkedHashMap<>();
    Map<String, Show> shows = new LinkedHashMap<>();
    Map<String, Booking> bookings = new LinkedHashMap<>();
    Map<String, SeatLock> locks = new LinkedHashMap<>();
    Map<String, TreeSet<String>> cityMovies = new LinkedHashMap<>();
    int lockCounter = 0;

    private void expireSeat(Seat seat, long currentTime) {
        if (seat.status == SeatStatus.LOCKED && currentTime >= seat.lockExpiry) {
            seat.status = SeatStatus.AVAILABLE;
            seat.lockedBy = "";
            seat.lockExpiry = 0;
        }
    }

    public void addTheater(String theaterId, String name, String city) {
        theaters.put(theaterId, new Theater(theaterId, name, city));
    }

    public void addShow(String showId, String theaterId, String movie, String time, int rows, int cols) {
        Show show = new Show(showId, theaterId, movie, time, rows, cols);
        shows.put(showId, show);
        Theater t = theaters.get(theaterId);
        if (t != null) {
            cityMovies.computeIfAbsent(t.city, k -> new TreeSet<>()).add(movie);
        }
    }

    public List<String> searchMovies(String city) {
        TreeSet<String> set = cityMovies.get(city);
        if (set == null) return new ArrayList<>();
        return new ArrayList<>(set);
    }

    public List<int[]> getAvailableSeats(String showId, long currentTime) {
        Show show = shows.get(showId);
        if (show == null) return new ArrayList<>();
        List<int[]> result = new ArrayList<>();
        for (int r = 0; r < show.rows; r++) {
            for (int c = 0; c < show.cols; c++) {
                expireSeat(show.seats[r][c], currentTime);
                if (show.seats[r][c].status == SeatStatus.AVAILABLE)
                    result.add(new int[]{r, c});
            }
        }
        return result;
    }

    public boolean bookSeats(String bookingId, String showId, List<int[]> seatPositions,
                              String userId, long currentTime) {
        Show show = shows.get(showId);
        if (show == null) return false;
        for (int[] pos : seatPositions) {
            int r = pos[0], c = pos[1];
            if (r < 0 || r >= show.rows || c < 0 || c >= show.cols) return false;
            expireSeat(show.seats[r][c], currentTime);
            if (show.seats[r][c].status != SeatStatus.AVAILABLE) return false;
        }
        for (int[] pos : seatPositions) {
            int r = pos[0], c = pos[1];
            show.seats[r][c].status = SeatStatus.BOOKED;
            show.seats[r][c].bookedBy = userId;
        }
        bookings.put(bookingId, new Booking(bookingId, showId, userId, seatPositions));
        return true;
    }

    public String lockSeats(String showId, List<int[]> seatPositions, String userId,
                             int ttlMinutes, long currentTime) {
        Show show = shows.get(showId);
        if (show == null) return "";
        for (int[] pos : seatPositions) {
            int r = pos[0], c = pos[1];
            if (r < 0 || r >= show.rows || c < 0 || c >= show.cols) return "";
            expireSeat(show.seats[r][c], currentTime);
            if (show.seats[r][c].status != SeatStatus.AVAILABLE) return "";
        }
        String lockId = "LOCK_" + (++lockCounter);
        long expiry = currentTime + (long) ttlMinutes * 60;
        for (int[] pos : seatPositions) {
            int r = pos[0], c = pos[1];
            show.seats[r][c].status = SeatStatus.LOCKED;
            show.seats[r][c].lockedBy = userId;
            show.seats[r][c].lockExpiry = expiry;
        }
        locks.put(lockId, new SeatLock(lockId, showId, userId, seatPositions, expiry));
        return lockId;
    }

    public boolean confirmBooking(String lockId, long currentTime) {
        SeatLock lock = locks.get(lockId);
        if (lock == null) return false;
        if (lock.confirmed || lock.released) return false;
        if (currentTime >= lock.expiry) {
            Show show = shows.get(lock.showId);
            if (show != null) {
                for (int[] pos : lock.seatPositions) expireSeat(show.seats[pos[0]][pos[1]], currentTime);
            }
            lock.released = true;
            return false;
        }
        Show show = shows.get(lock.showId);
        if (show == null) return false;
        for (int[] pos : lock.seatPositions) {
            int r = pos[0], c = pos[1];
            show.seats[r][c].status = SeatStatus.BOOKED;
            show.seats[r][c].bookedBy = lock.userId;
            show.seats[r][c].lockedBy = "";
            show.seats[r][c].lockExpiry = 0;
        }
        String bookingId = "BK_" + lockId;
        bookings.put(bookingId, new Booking(bookingId, lock.showId, lock.userId, lock.seatPositions));
        lock.confirmed = true;
        return true;
    }

    public boolean releaseLock(String lockId, long currentTime) {
        SeatLock lock = locks.get(lockId);
        if (lock == null) return false;
        if (lock.confirmed || lock.released) return false;
        Show show = shows.get(lock.showId);
        if (show != null) {
            for (int[] pos : lock.seatPositions) {
                int r = pos[0], c = pos[1];
                show.seats[r][c].status = SeatStatus.AVAILABLE;
                show.seats[r][c].lockedBy = "";
                show.seats[r][c].lockExpiry = 0;
            }
        }
        lock.released = true;
        return true;
    }

    public List<int[]> findContiguousSeats(String showId, int n, long currentTime) {
        Show show = shows.get(showId);
        if (show == null) return new ArrayList<>();
        for (int r = 0; r < show.rows; r++) {
            int count = 0, start = -1;
            for (int c = 0; c < show.cols; c++) {
                expireSeat(show.seats[r][c], currentTime);
                if (show.seats[r][c].status == SeatStatus.AVAILABLE) {
                    if (count == 0) start = c;
                    count++;
                    if (count == n) {
                        List<int[]> result = new ArrayList<>();
                        for (int i = start; i < start + n; i++) result.add(new int[]{r, i});
                        return result;
                    }
                } else {
                    count = 0;
                    start = -1;
                }
            }
        }
        return new ArrayList<>();
    }
}

public class Solution {
    private static List<int[]> parseSeats(String s) {
        List<int[]> out = new ArrayList<>();
        if (s == null || s.isEmpty()) return out;
        for (String tok : s.split(";")) {
            if (tok.isEmpty()) continue;
            int comma = tok.indexOf(',');
            if (comma < 0) continue;
            int r = Integer.parseInt(tok.substring(0, comma));
            int c = Integer.parseInt(tok.substring(comma + 1));
            out.add(new int[]{r, c});
        }
        return out;
    }

    public static List<String> show_simulate(List<ShowOp> ops) {
        List<String> out = new ArrayList<>();
        BookingSystem sys = new BookingSystem();
        String[] lockSlots = new String[32];
        for (int i = 0; i < 32; i++) lockSlots[i] = "";
        List<int[]> lastContig = new ArrayList<>();

        for (ShowOp op : ops) {
            String k = op.kind;
            if ("new".equals(k)) {
                sys = new BookingSystem();
                for (int i = 0; i < 32; i++) lockSlots[i] = "";
                lastContig = new ArrayList<>();
                out.add("ok");
            } else if ("add_theater".equals(k)) {
                sys.addTheater(op.s1, op.s2, op.s3); out.add("ok");
            } else if ("add_show".equals(k)) {
                sys.addShow(op.s1, op.s2, op.s3, op.s4, op.i1, op.i2); out.add("ok");
            } else if ("movies_count".equals(k)) {
                out.add(Integer.toString(sys.searchMovies(op.s1).size()));
            } else if ("movies_contains".equals(k)) {
                List<String> m = sys.searchMovies(op.s1);
                out.add(m.contains(op.s2) ? "yes" : "no");
            } else if ("available_count".equals(k)) {
                out.add(Integer.toString(sys.getAvailableSeats(op.s1, (long) op.i1).size()));
            } else if ("available_has".equals(k)) {
                List<int[]> v = sys.getAvailableSeats(op.s1, (long) op.i1);
                boolean found = false;
                for (int[] p : v) if (p[0] == op.i2 && p[1] == op.i3) { found = true; break; }
                out.add(found ? "yes" : "no");
            } else if ("book".equals(k)) {
                out.add(sys.bookSeats(op.s1, op.s2, parseSeats(op.s3), op.s4, (long) op.i1) ? "ok" : "fail");
            } else if ("lock".equals(k)) {
                String lid = sys.lockSeats(op.s1, parseSeats(op.s2), op.s3, op.i1, (long) op.i2);
                if (op.i3 >= 0 && op.i3 < lockSlots.length) lockSlots[op.i3] = lid;
                out.add(lid);
            } else if ("lock_at".equals(k)) {
                out.add(lockSlots[op.i3]);
            } else if ("confirm".equals(k)) {
                out.add(sys.confirmBooking(lockSlots[op.i3], (long) op.i1) ? "ok" : "fail");
            } else if ("release".equals(k)) {
                out.add(sys.releaseLock(lockSlots[op.i3], (long) op.i1) ? "ok" : "fail");
            } else if ("release_id".equals(k)) {
                out.add(sys.releaseLock(op.s1, (long) op.i1) ? "ok" : "fail");
            } else if ("find_contig".equals(k)) {
                lastContig = sys.findContiguousSeats(op.s1, op.i1, (long) op.i2);
                out.add(Integer.toString(lastContig.size()));
            } else if ("contig_at".equals(k)) {
                if (op.i1 < 0 || op.i1 >= lastContig.size()) out.add("");
                else out.add(lastContig.get(op.i1)[0] + "," + lastContig.get(op.i1)[1]);
            } else {
                out.add("unknown:" + k);
            }
        }
        return out;
    }
}
