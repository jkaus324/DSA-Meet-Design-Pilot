// Elevator System — Solution (Java)
import java.util.*;

class ElevOp {
    public String kind;
    public String s1;
    public int i1;
    public int i2;

    public ElevOp(String kind, String s1, int i1, int i2) {
        this.kind = kind;
        this.s1 = s1;
        this.i1 = i1;
        this.i2 = i2;
    }
}

enum ElevatorState { IDLE, MOVING_UP, MOVING_DOWN, DOOR_OPEN }

enum Direction { UP, DOWN, NONE }

class Elevator {
    int id;
    int currentFloor;
    ElevatorState state;
    Direction currentDirection;
    TreeSet<Integer> upRequests = new TreeSet<>();
    TreeSet<Integer> downRequests = new TreeSet<>();

    public Elevator() { this(0); }

    public Elevator(int id) {
        this.id = id;
        this.currentFloor = 0;
        this.state = ElevatorState.IDLE;
        this.currentDirection = Direction.NONE;
    }

    public int getId() { return id; }
    public int getCurrentFloor() { return currentFloor; }
    public ElevatorState getState() { return state; }
    public Direction getCurrentDirection() { return currentDirection; }
    public int getPendingCount() { return upRequests.size() + downRequests.size(); }

    public void addRequest(int floor, Direction direction) {
        if (floor == currentFloor && state == ElevatorState.IDLE) {
            state = ElevatorState.DOOR_OPEN;
            return;
        }

        if (direction == Direction.UP) {
            upRequests.add(floor);
        } else if (direction == Direction.DOWN) {
            downRequests.add(floor);
        } else {
            if (floor > currentFloor) upRequests.add(floor);
            else downRequests.add(floor);
        }

        if (state == ElevatorState.IDLE) {
            if (!upRequests.isEmpty() && (downRequests.isEmpty() ||
                Math.abs(upRequests.first() - currentFloor) <= Math.abs(downRequests.last() - currentFloor))) {
                currentDirection = Direction.UP;
                state = ElevatorState.MOVING_UP;
            } else {
                currentDirection = Direction.DOWN;
                state = ElevatorState.MOVING_DOWN;
            }
        }
    }

    public void step() {
        switch (state) {
            case IDLE:
                break;
            case MOVING_UP:
                currentFloor++;
                if (upRequests.contains(currentFloor)) {
                    upRequests.remove(currentFloor);
                    state = ElevatorState.DOOR_OPEN;
                }
                break;
            case MOVING_DOWN:
                currentFloor--;
                if (downRequests.contains(currentFloor)) {
                    downRequests.remove(currentFloor);
                    state = ElevatorState.DOOR_OPEN;
                }
                break;
            case DOOR_OPEN:
                if (currentDirection == Direction.UP) {
                    if (!upRequests.isEmpty()) state = ElevatorState.MOVING_UP;
                    else if (!downRequests.isEmpty()) { currentDirection = Direction.DOWN; state = ElevatorState.MOVING_DOWN; }
                    else { currentDirection = Direction.NONE; state = ElevatorState.IDLE; }
                } else if (currentDirection == Direction.DOWN) {
                    if (!downRequests.isEmpty()) state = ElevatorState.MOVING_DOWN;
                    else if (!upRequests.isEmpty()) { currentDirection = Direction.UP; state = ElevatorState.MOVING_UP; }
                    else { currentDirection = Direction.NONE; state = ElevatorState.IDLE; }
                } else {
                    if (!upRequests.isEmpty()) { currentDirection = Direction.UP; state = ElevatorState.MOVING_UP; }
                    else if (!downRequests.isEmpty()) { currentDirection = Direction.DOWN; state = ElevatorState.MOVING_DOWN; }
                    else state = ElevatorState.IDLE;
                }
                break;
        }
    }
}

interface DispatchStrategy {
    int selectElevator(List<Elevator> elevators, int requestFloor, Direction requestDirection);
}

class NearestFirst implements DispatchStrategy {
    @Override
    public int selectElevator(List<Elevator> elevators, int requestFloor, Direction requestDirection) {
        int bestIdx = 0;
        int bestScore = Integer.MAX_VALUE;
        final int PENALTY = 10000;

        for (int i = 0; i < elevators.size(); i++) {
            Elevator e = elevators.get(i);
            int dist = Math.abs(e.getCurrentFloor() - requestFloor);
            int score;
            ElevatorState st = e.getState();
            Direction dir = e.getCurrentDirection();

            if (st == ElevatorState.IDLE || dir == Direction.NONE) {
                score = dist;
            } else if (dir == Direction.UP && requestDirection == Direction.UP
                       && e.getCurrentFloor() <= requestFloor) {
                score = dist;
            } else if (dir == Direction.DOWN && requestDirection == Direction.DOWN
                       && e.getCurrentFloor() >= requestFloor) {
                score = dist;
            } else {
                score = dist + PENALTY;
            }

            if (score < bestScore) {
                bestScore = score;
                bestIdx = i;
            }
        }
        return bestIdx;
    }
}

class LeastLoaded implements DispatchStrategy {
    @Override
    public int selectElevator(List<Elevator> elevators, int requestFloor, Direction requestDirection) {
        int bestIdx = 0;
        int bestCount = Integer.MAX_VALUE;
        for (int i = 0; i < elevators.size(); i++) {
            int cnt = elevators.get(i).getPendingCount();
            if (cnt < bestCount) {
                bestCount = cnt;
                bestIdx = i;
            }
        }
        return bestIdx;
    }
}

class ElevatorSystem {
    List<Elevator> elevators = new ArrayList<>();
    DispatchStrategy strategy = null;

    public void addElevator(int id) {
        elevators.add(new Elevator(id));
    }

    public void setDispatchStrategy(DispatchStrategy s) {
        this.strategy = s;
    }

    public Elevator getElevator(int index) {
        if (index < 0 || index >= elevators.size()) return null;
        return elevators.get(index);
    }

    public int getElevatorCount() { return elevators.size(); }

    public void addRequest(int floor, Direction direction) {
        if (elevators.isEmpty()) return;
        int idx = 0;
        if (strategy != null) {
            idx = strategy.selectElevator(elevators, floor, direction);
        }
        elevators.get(idx).addRequest(floor, direction);
    }

    public void step() {
        for (Elevator e : elevators) e.step();
    }
}

public class Solution {
    private static Direction dirFrom(String s) {
        if ("up".equals(s)) return Direction.UP;
        if ("down".equals(s)) return Direction.DOWN;
        return Direction.NONE;
    }

    public static List<String> elevator_simulate(List<ElevOp> ops) {
        List<String> out = new ArrayList<>();
        Elevator single = null;
        ElevatorSystem sys = null;
        NearestFirst nf = new NearestFirst();
        LeastLoaded ll = new LeastLoaded();

        for (ElevOp op : ops) {
            String k = op.kind;
            if ("new_elev".equals(k)) {
                single = new Elevator();
                sys = null;
                out.add("ok");
            } else if ("new_sys".equals(k)) {
                sys = new ElevatorSystem();
                single = null;
                out.add("ok");
            } else if ("add_elev".equals(k)) {
                sys.addElevator(op.i1);
                out.add("ok");
            } else if ("set_strategy".equals(k)) {
                if ("nearest".equals(op.s1)) sys.setDispatchStrategy(nf);
                else if ("least_loaded".equals(op.s1)) sys.setDispatchStrategy(ll);
                out.add("ok");
            } else if ("req".equals(k)) {
                single.addRequest(op.i1, dirFrom(op.s1));
                out.add("ok");
            } else if ("sys_req".equals(k)) {
                sys.addRequest(op.i1, dirFrom(op.s1));
                out.add("ok");
            } else if ("elev_req".equals(k)) {
                sys.getElevator(op.i1).addRequest(op.i2, dirFrom(op.s1));
                out.add("ok");
            } else if ("step".equals(k)) {
                single.step();
                out.add("ok");
            } else if ("sys_step".equals(k)) {
                sys.step();
                out.add("ok");
            } else if ("elev_step".equals(k)) {
                sys.getElevator(op.i1).step();
                out.add("ok");
            } else if ("floor".equals(k)) {
                out.add(Integer.toString(single.getCurrentFloor()));
            } else if ("elev_floor".equals(k)) {
                out.add(Integer.toString(sys.getElevator(op.i1).getCurrentFloor()));
            } else if ("state".equals(k)) {
                out.add(single.getState().name());
            } else if ("elev_state".equals(k)) {
                out.add(sys.getElevator(op.i1).getState().name());
            } else if ("elev_pending".equals(k)) {
                out.add(Integer.toString(sys.getElevator(op.i1).getPendingCount()));
            } else if ("count".equals(k)) {
                out.add(Integer.toString(sys.getElevatorCount()));
            } else if ("elev_null".equals(k)) {
                out.add(sys.getElevator(op.i1) == null ? "yes" : "no");
            } else {
                out.add("unknown:" + k);
            }
        }
        return out;
    }
}
