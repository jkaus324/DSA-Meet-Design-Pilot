# Test Spec Format — CodeJunction

Each problem part has a `tests/specs/partN_spec.yaml` file.  
Run `python scripts/generate_tests.py <problem_dir>` to produce:
- `tests/cpp/partN_test.cpp`
- `tests/java/PartNTest.java`

Both files produce identical test scenarios in the `PASS/FAIL + PARTX_SUMMARY` format.

---

## Top-Level Fields

```yaml
# Required
tests: [...]          # list of test cases (see below)

# Optional
includes:
  cpp: ["<sstream>"]  # extra #include lines (angle brackets included)
  java: ["java.util.concurrent.*"]  # extra import lines

helpers_cpp: |        # raw C++ helper functions placed before the test function
  bool isValidOrder(...) { ... }

helpers_java: |       # raw Java helper methods placed inside the test class
  static boolean isValidOrder(...) { ... }
```

---

## Test Case Fields

```yaml
tests:
  - name: "snake_case_name"    # Required. Becomes PASS/FAIL label and Java method name.
    comment: "what this tests" # Optional. Emitted as a comment.
    
    steps: [...]               # Setup and operations (see Step Types)
    
    # Pick ONE assertion style:
    assertions: [...]          # Structured assertions (see Assertion Types)
    pass: "bool_expr"          # Free-form boolean expression (same in both langs)
    pass_cpp: "bool_expr"      # C++-specific expression (use with pass_java)
    pass_java: "bool_expr"     # Java-specific expression (use with pass_cpp)
    throws: "expr_that_throws" # Expects this expression to throw (wraps in try/catch)
```

---

## Step Types

```yaml
steps:
  # Construct an object (no-arg or with args)
  - new: "ClassName varName"
  - new: "ClassName varName(arg1, arg2)"

  # Void method call
  - call: "obj.method(args)"

  # Assign return value (type inferred as auto/var)
  - var: "result = obj.getValue()"
  
  # Assign with explicit type
  - var: "int count = obj.size()"

  # Raw code — emitted in BOTH languages (apply minimal transforms)
  - raw: "obj.setFlag(true);"

  # Language-specific raw code
  - raw_cpp: "auto& ref = container.back();"
  - raw_java: "var last = list.get(list.size()-1);"
```

---

## Assertion Types

```yaml
assertions:
  - eq: ["lhs", "rhs"]           # lhs == rhs  (strings get .equals() in Java)
  - ne: ["lhs", "rhs"]           # lhs != rhs
  - lt: ["lhs", "rhs"]           # lhs < rhs
  - le: ["lhs", "rhs"]           # lhs <= rhs
  - gt: ["lhs", "rhs"]           # lhs > rhs
  - ge: ["lhs", "rhs"]           # lhs >= rhs
  - true: "bool_expr"            # assert(bool_expr)
  - false: "bool_expr"           # assert(!bool_expr)
  - size: ["container", 3]       # container.size() == 3
  - empty: "container"           # container is empty
  - not_empty: "container"       # container is not empty
  - approx: ["expr", 3.14, 1e-9] # |expr - 3.14| < 1e-9
  - contains: ["container", "val"] # container contains val
  - not_contains: ["container", "val"]
  - raw: "bool_expr"             # custom expression (same in both langs)
```

---

## Auto-translations applied to Java output

| C++ | Java |
|-----|------|
| `nullptr` | `null` |
| `.empty()` | `.isEmpty()` |
| `var.back()` | `var.get(var.size()-1)` |
| `var.front()` | `var.get(0)` |
| `std::` | `` (removed) |
| `string` | `String` |
| `bool` | `boolean` |
| `"str" == x` | `"str".equals(x)` |
| Constructor `Foo bar(args)` | `Foo bar = new Foo(args)` |
| `auto x = expr` | `var x = expr` |

---

## Full Example

```yaml
# problems/tier1-foundation/073-tic-tac-toe/tests/specs/part1_spec.yaml

tests:
  - name: "make_move_valid"
    comment: "Player X makes a valid move on empty board"
    steps:
      - new: "TicTacToe game(3)"
      - var: "result = game.makeMove(0, 0, 'X')"
    assertions:
      - eq: ["result", "true"]

  - name: "make_move_occupied"
    comment: "Cannot place on already occupied cell"
    steps:
      - new: "TicTacToe game(3)"
      - call: "game.makeMove(0, 0, 'X')"
      - var: "result = game.makeMove(0, 0, 'O')"
    assertions:
      - eq: ["result", "false"]

  - name: "win_row"
    comment: "X wins by filling a row"
    steps:
      - new: "TicTacToe game(3)"
      - call: "game.makeMove(0, 0, 'X')"
      - call: "game.makeMove(0, 1, 'X')"
      - call: "game.makeMove(0, 2, 'X')"
    assertions:
      - eq: ["game.getWinner()", "'X'"]

  - name: "no_winner_empty"
    comment: "Empty board has no winner"
    steps:
      - new: "TicTacToe game(3)"
    assertions:
      - eq: ["game.getWinner()", "'\\0'"]  # null char = no winner
```

---

## Workflow for New Problems

```
1. Write solution.cpp  (C++ reference implementation)
2. Write tests/specs/partN_spec.yaml  (once per part)
3. python scripts/generate_tests.py <problem_dir>
   → generates tests/cpp/partN_test.cpp
   → generates tests/java/PartNTest.java
4. python scripts/cpp_to_java.py <problem_dir>
   → generates solution.java + boilerplate/java/
5. Update docs/_data/problems.yml → add "java" to languages
```

For changes to test logic, edit the spec and re-run step 3 with `--force`.  
The C++ and Java test files are **generated artifacts** — do not edit them directly.
