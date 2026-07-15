#!/usr/bin/env python3
"""
generate_tests.py — CodeJunction YAML spec → C++ + Java test generator
=======================================================================
Reads test specs (YAML) and generates test files for each language.
This eliminates duplication: write test logic ONCE, get C++ and Java for free.

Usage:
  python scripts/generate_tests.py <problem_dir>            # one problem
  python scripts/generate_tests.py --all                    # all problems with specs
  python scripts/generate_tests.py <dir> --force            # overwrite existing

Input:   tests/specs/partN_spec.yaml
Output:  tests/cpp/partN_test.cpp   (same PASS/FAIL/PARTX_SUMMARY format)
         tests/java/PartNTest.java  (same format, Java idioms)

Spec format: see scripts/SPEC_FORMAT.md or docs/SPEC_FORMAT.md
"""

import re, os, sys
from pathlib import Path

try:
    import yaml
except ImportError:
    print("pyyaml not found. Run: pip install pyyaml")
    sys.exit(1)

FORCE = '--force' in sys.argv

# ═══════════════════════════════════════════════════════════════════════════════
# Expression helpers
# ═══════════════════════════════════════════════════════════════════════════════

def to_java_expr(s):
    """Apply Java-specific rewrites to an expression string."""
    s = str(s)
    s = re.sub(r'\bnullptr\b', 'null', s)
    s = re.sub(r'\bstd::', '', s)
    s = re.sub(r'\bbool\b', 'boolean', s)
    s = re.sub(r'\bstring\b', 'String', s)
    # .empty() → .isEmpty()
    s = re.sub(r'\.empty\(\)', '.isEmpty()', s)
    # .back() → .get(var.size()-1) — requires variable name
    s = re.sub(r'(\w+)\.back\(\)', r'\1.get(\1.size()-1)', s)
    # .front() → .get(0)
    s = re.sub(r'(\w+)\.front\(\)', r'\1.get(0)', s)
    # nullptr → null
    s = re.sub(r'\bnullptr\b', 'null', s)
    return s

def is_string_val(s):
    """True if s is a double-quoted String literal (NOT char literal)."""
    s = str(s).strip()
    return s.startswith('"')

def is_char_val(s):
    """True if s is a single-quoted char literal like 'X' or '\\0'."""
    s = str(s).strip()
    return s.startswith("'")

def normalize_assertion(a):
    """
    YAML parses 'true:' and 'false:' keys as Python booleans True/False.
    Normalize them back to strings so downstream code can use string comparisons.
    """
    result = {}
    for k, v in a.items():
        if k is True:
            result['true'] = v
        elif k is False:
            result['false'] = v
        else:
            result[k] = v
    return result

def java_eq(lhs, rhs):
    """Java equality: use .equals() only for double-quoted String literals."""
    lhs_s = to_java_expr(str(lhs))
    rhs_s = to_java_expr(str(rhs))
    # char literals and primitives: use ==
    if is_char_val(lhs) or is_char_val(rhs):
        return f'{lhs_s} == {rhs_s}'
    # String literals: use .equals()
    if is_string_val(lhs) or is_string_val(rhs):
        return f'{lhs_s}.equals({rhs_s})'
    return f'{lhs_s} == {rhs_s}'

def java_ne(lhs, rhs):
    lhs_s = to_java_expr(str(lhs))
    rhs_s = to_java_expr(str(rhs))
    if is_char_val(lhs) or is_char_val(rhs):
        return f'{lhs_s} != {rhs_s}'
    if is_string_val(lhs) or is_string_val(rhs):
        return f'!{lhs_s}.equals({rhs_s})'
    return f'{lhs_s} != {rhs_s}'

# ═══════════════════════════════════════════════════════════════════════════════
# Step translation
# ═══════════════════════════════════════════════════════════════════════════════

def parse_new_decl(s):
    """
    Parse "ClassName var" or "ClassName var(args)" or "ClassName var(arg1, arg2)".
    Returns (ctype, varname, args_string_or_None).
    """
    s = s.strip()
    # With args: "ClassName var(arg1, arg2)"
    m = re.match(r'^(\w+)\s+(\w+)\s*\(([^)]*)\)\s*$', s)
    if m:
        return m.group(1), m.group(2), m.group(3).strip()
    # No args: "ClassName var"
    m = re.match(r'^(\w+)\s+(\w+)\s*$', s)
    if m:
        return m.group(1), m.group(2), None
    return None, None, None


def step_to_cpp(step, indent='        '):
    """Convert a spec step to C++ source line(s)."""
    if step is None:
        return ''

    # new: "ClassName var" or "ClassName var(args)"
    if 'new' in step:
        ctype, var, args = parse_new_decl(step['new'])
        if ctype:
            return f'{indent}{ctype} {var}({args if args else ""});'
        return f'{indent}{step["new"]};'

    # call: "obj.method(args)"  — void call
    if 'call' in step:
        return f'{indent}{step["call"].rstrip(";")};'

    # var: "result = expr" or "Type result = expr"
    if 'var' in step:
        expr = step['var'].strip()
        if '=' in expr:
            lhs, rhs = expr.split('=', 1)
            lhs = lhs.strip(); rhs = rhs.strip()
            if re.match(r'^[a-z]', lhs) and ' ' not in lhs:
                return f'{indent}auto {lhs} = {rhs};'
            return f'{indent}{lhs} = {rhs};'
        return f'{indent}auto {expr};'

    # for: loop step (raw, same in both languages usually)
    if 'for_cpp' in step:
        lines = step['for_cpp'].strip().splitlines()
        return '\n'.join(f'{indent}{l}' for l in lines)

    # raw_cpp: verbatim C++ (skipped in Java)
    if 'raw_cpp' in step:
        lines = step['raw_cpp'].strip().splitlines()
        return '\n'.join(f'{indent}{l}' for l in lines)

    # raw: same code in both languages
    if 'raw' in step:
        lines = step['raw'].strip().splitlines()
        return '\n'.join(f'{indent}{l}' for l in lines)

    # raw_java: Java-only step, skip in C++
    if 'raw_java' in step:
        return ''

    return f'{indent}// UNKNOWN STEP: {step}'


def step_to_java(step, indent='            '):
    """Convert a spec step to Java source line(s)."""
    if step is None:
        return ''

    if 'new' in step:
        ctype, var, args = parse_new_decl(step['new'])
        if ctype:
            args_j = to_java_expr(args) if args else ''
            return f'{indent}{ctype} {var} = new {ctype}({args_j});'
        return f'{indent}{to_java_expr(step["new"])};'

    if 'call' in step:
        return f'{indent}{to_java_expr(step["call"].rstrip(";"))};'

    if 'var' in step:
        expr = step['var'].strip()
        if '=' in expr:
            lhs, rhs = expr.split('=', 1)
            lhs = lhs.strip(); rhs = rhs.strip()
            if re.match(r'^[a-z]', lhs) and ' ' not in lhs:
                return f'{indent}var {lhs} = {to_java_expr(rhs)};'
            return f'{indent}{lhs} = {to_java_expr(rhs)};'
        return f'{indent}var {to_java_expr(expr)};'

    if 'for_java' in step:
        lines = step['for_java'].strip().splitlines()
        return '\n'.join(f'{indent}{l}' for l in lines)

    if 'raw_java' in step:
        lines = step['raw_java'].strip().splitlines()
        return '\n'.join(f'{indent}{l}' for l in lines)

    if 'raw' in step:
        lines = to_java_expr(step['raw']).strip().splitlines()
        return '\n'.join(f'{indent}{l}' for l in lines)

    if 'raw_cpp' in step:
        return ''

    return f'{indent}// UNKNOWN STEP: {step}'


# ═══════════════════════════════════════════════════════════════════════════════
# Assertion translation
# ═══════════════════════════════════════════════════════════════════════════════

def assertion_to_cpp(a, indent='        '):
    """Convert a structured assertion to a C++ assert() statement."""
    a = normalize_assertion(a)
    if 'eq' in a:
        l, r = a['eq']
        return f'{indent}assert({l} == {r});'
    if 'ne' in a:
        l, r = a['ne']
        return f'{indent}assert({l} != {r});'
    if 'lt' in a:
        l, r = a['lt']
        return f'{indent}assert({l} < {r});'
    if 'le' in a:
        l, r = a['le']
        return f'{indent}assert({l} <= {r});'
    if 'gt' in a:
        l, r = a['gt']
        return f'{indent}assert({l} > {r});'
    if 'ge' in a:
        l, r = a['ge']
        return f'{indent}assert({l} >= {r});'
    if 'true' in a:
        return f'{indent}assert({a["true"]});'
    if 'false' in a:
        return f'{indent}assert(!({a["false"]}));'
    if 'size' in a:
        container, n = a['size']
        return f'{indent}assert({container}.size() == {n});'
    if 'empty' in a:
        return f'{indent}assert({a["empty"]}.empty());'
    if 'not_empty' in a:
        return f'{indent}assert(!{a["not_empty"]}.empty());'
    if 'approx' in a:
        expr, val, eps = a['approx']
        return f'{indent}assert(abs(({expr}) - ({val})) < {eps});'
    if 'contains' in a:
        container, val = a['contains']
        return f'{indent}assert(find({container}.begin(), {container}.end(), {val}) != {container}.end());'
    if 'not_contains' in a:
        container, val = a['not_contains']
        return f'{indent}assert(find({container}.begin(), {container}.end(), {val}) == {container}.end());'
    if 'raw' in a:
        return f'{indent}assert({a["raw"]});'
    return f'{indent}// UNKNOWN ASSERTION: {a}'


def assertion_to_java_expr(a):
    """Convert a structured assertion to a Java boolean sub-expression."""
    a = normalize_assertion(a)
    if 'eq' in a:
        return java_eq(a['eq'][0], a['eq'][1])
    if 'ne' in a:
        return java_ne(a['ne'][0], a['ne'][1])
    if 'lt' in a:
        l, r = a['lt']
        return f'{to_java_expr(str(l))} < {to_java_expr(str(r))}'
    if 'le' in a:
        l, r = a['le']
        return f'{to_java_expr(str(l))} <= {to_java_expr(str(r))}'
    if 'gt' in a:
        l, r = a['gt']
        return f'{to_java_expr(str(l))} > {to_java_expr(str(r))}'
    if 'ge' in a:
        l, r = a['ge']
        return f'{to_java_expr(str(l))} >= {to_java_expr(str(r))}'
    if 'true' in a:
        return to_java_expr(str(a['true']))
    if 'false' in a:
        return f'!({to_java_expr(str(a["false"]))})'
    if 'size' in a:
        container, n = a['size']
        return f'{to_java_expr(str(container))}.size() == {n}'
    if 'empty' in a:
        return f'{to_java_expr(str(a["empty"]))}.isEmpty()'
    if 'not_empty' in a:
        return f'!{to_java_expr(str(a["not_empty"]))}.isEmpty()'
    if 'approx' in a:
        expr, val, eps = a['approx']
        return f'Math.abs(({to_java_expr(str(expr))}) - ({val})) < {eps}'
    if 'contains' in a:
        container, val = a['contains']
        return f'{to_java_expr(str(container))}.contains({to_java_expr(str(val))})'
    if 'not_contains' in a:
        container, val = a['not_contains']
        return f'!{to_java_expr(str(container))}.contains({to_java_expr(str(val))})'
    if 'raw' in a:
        return to_java_expr(str(a['raw']))
    return f'true /* UNKNOWN ASSERTION: {a} */'


# ═══════════════════════════════════════════════════════════════════════════════
# Test block generation
# ═══════════════════════════════════════════════════════════════════════════════

def test_block_cpp(test, idx):
    """Generate a single C++ try/catch test block."""
    name = test.get('name', f'test_{idx + 1}')
    comment = test.get('comment', '')
    steps = test.get('steps', [])
    assertions = test.get('assertions', [])
    pass_expr_cpp = test.get('pass_cpp', test.get('pass', None))
    throws_expr = test.get('throws', None)

    lines = []
    if comment:
        lines.append(f'    // {comment}')
    lines.append(f'    try {{')

    for step in steps:
        line = step_to_cpp(step)
        if line:
            lines.append(line)

    if throws_expr:
        # Expect a specific call to throw
        lines.append(f'        bool threw = false;')
        lines.append(f'        try {{ {throws_expr.rstrip(";")}; }}')
        lines.append(f'        catch (...) {{ threw = true; }}')
        lines.append(f'        assert(threw);')
    elif pass_expr_cpp:
        lines.append(f'        assert({pass_expr_cpp});')
    else:
        for a in assertions:
            lines.append(assertion_to_cpp(a))

    lines.append(f'        cout << "PASS {name}" << endl;')
    lines.append(f'        passed++;')
    lines.append(f'    }} catch (...) {{')
    lines.append(f'        cout << "FAIL {name}" << endl;')
    lines.append(f'        failed++;')
    lines.append(f'    }}')
    lines.append('')
    return '\n'.join(lines)


def snake_to_camel(s):
    """Convert test_foo_bar → testFooBar."""
    parts = s.split('_')
    return parts[0] + ''.join(p.capitalize() for p in parts[1:])


def test_block_java(test, idx):
    """Generate a single Java static boolean test method."""
    name = test.get('name', f'test_{idx + 1}')
    comment = test.get('comment', '')
    steps = test.get('steps', [])
    assertions = test.get('assertions', [])
    pass_expr_java = test.get('pass_java', test.get('pass', None))
    throws_expr = test.get('throws', None)

    method_name = snake_to_camel(name)

    lines = []
    lines.append(f'    static boolean {method_name}() {{')
    if comment:
        lines.append(f'        // {comment}')
    lines.append(f'        try {{')

    for step in steps:
        line = step_to_java(step)
        if line:
            lines.append(line)

    if throws_expr:
        throws_j = to_java_expr(throws_expr.rstrip(';'))
        lines.append(f'            boolean threw = false;')
        lines.append(f'            try {{ {throws_j}; }}')
        lines.append(f'            catch (Exception ex) {{ threw = true; }}')
        lines.append(f'            boolean pass = threw;')
    elif pass_expr_java:
        pass_j = to_java_expr(pass_expr_java)
        lines.append(f'            boolean pass = {pass_j};')
    elif assertions:
        exprs = [assertion_to_java_expr(a) for a in assertions]
        if len(exprs) == 1:
            lines.append(f'            boolean pass = {exprs[0]};')
        else:
            first = exprs[0]
            rest = exprs[1:]
            lines.append(f'            boolean pass = {first}')
            for expr in rest[:-1]:
                lines.append(f'                && {expr}')
            lines.append(f'                && {rest[-1]};')
    else:
        lines.append(f'            boolean pass = true;')

    lines.append(f'            System.out.println((pass ? "PASS" : "FAIL") + ": {name}");')
    lines.append(f'            return pass;')
    lines.append(f'        }} catch (Exception e) {{')
    lines.append(f'            System.out.println("FAIL: {name} (exception: " + e.getMessage() + ")");')
    lines.append(f'            return false;')
    lines.append(f'        }}')
    lines.append(f'    }}')
    lines.append('')
    return '\n'.join(lines)


# ═══════════════════════════════════════════════════════════════════════════════
# Full file generation
# ═══════════════════════════════════════════════════════════════════════════════

def generate_cpp_file(spec, display_name, part_num):
    """Generate a complete C++ test file from spec."""
    tests = spec.get('tests', [])
    extra_includes = spec.get('includes', {}).get('cpp', [])
    helpers_cpp = spec.get('helpers_cpp', '')

    lines = []
    lines.append(f'// {display_name} — Part {part_num} Tests')
    lines.append(f'#define RUNNING_TESTS')
    lines.append(f'#include "../../solution.cpp"')
    lines.append(f'#include <cassert>')
    lines.append(f'#include <iostream>')
    lines.append(f'#include <algorithm>')
    for inc in extra_includes:
        lines.append(f'#include {inc}')
    lines.append(f'using namespace std;')
    lines.append('')

    if helpers_cpp:
        lines.append(helpers_cpp.strip())
        lines.append('')

    lines.append(f'int part{part_num}_tests() {{')
    lines.append(f'    int passed = 0, failed = 0;')
    lines.append('')

    for i, test in enumerate(tests):
        lines.append(test_block_cpp(test, i))

    lines.append(f'    cout << "PART{part_num}_SUMMARY " << passed << "/" << (passed + failed) << endl;')
    lines.append(f'    return failed;')
    lines.append(f'}}')
    lines.append('')
    lines.append(f'int main() {{ return part{part_num}_tests(); }}')
    return '\n'.join(lines) + '\n'


def generate_java_file(spec, display_name, part_num):
    """Generate a complete Java test file from spec."""
    tests = spec.get('tests', [])
    extra_imports = spec.get('includes', {}).get('java', [])
    helpers_java = spec.get('helpers_java', '')

    lines = []
    lines.append(f'// {display_name} — Part {part_num} Tests')
    lines.append(f'import java.util.*;')
    lines.append(f'import java.util.stream.*;')
    for imp in extra_imports:
        imp = imp.strip()
        if not imp.startswith('import'):
            lines.append(f'import {imp};')
        else:
            lines.append(f'{imp};' if not imp.endswith(';') else imp)
    lines.append('')
    lines.append(f'class Part{part_num}Test {{')

    if helpers_java:
        for hline in helpers_java.strip().splitlines():
            lines.append(f'    {hline}')
        lines.append('')

    for i, test in enumerate(tests):
        lines.append(test_block_java(test, i))

    # runTests
    lines.append(f'    public static int runTests() {{')
    lines.append(f'        int passed = 0, total = 0;')
    for test in tests:
        name = test.get('name', f'test_{tests.index(test) + 1}')
        method_name = snake_to_camel(name)
        lines.append(f'        total++; if ({method_name}()) passed++;')
    lines.append(f'        System.out.println("PART{part_num}_SUMMARY " + passed + "/" + total);')
    lines.append(f'        return passed;')
    lines.append(f'    }}')
    lines.append(f'}}')
    return '\n'.join(lines) + '\n'


# ═══════════════════════════════════════════════════════════════════════════════
# Problem processing
# ═══════════════════════════════════════════════════════════════════════════════

def process_problem(problem_dir):
    path = Path(problem_dir)
    spec_dir = path / 'tests' / 'specs'

    if not spec_dir.exists():
        print(f'  No tests/specs/ directory found.')
        return 0, 0

    spec_files = sorted(spec_dir.glob('part*_spec.yaml'))
    if not spec_files:
        print(f'  No spec files found in tests/specs/')
        return 0, 0

    # Pretty display name: "052-task-queue" → "Task Queue"
    raw_name = path.name
    parts_of_name = raw_name.split('-')
    display_name = ' '.join(w.capitalize() for w in parts_of_name[1:])

    cpp_dir = path / 'tests' / 'cpp'
    java_dir = path / 'tests' / 'java'
    cpp_dir.mkdir(parents=True, exist_ok=True)
    java_dir.mkdir(parents=True, exist_ok=True)

    generated = 0
    skipped = 0

    for spec_file in spec_files:
        m = re.search(r'part(\d+)_spec', spec_file.name)
        if not m:
            continue
        part_num = int(m.group(1))

        with open(spec_file, encoding='utf-8') as f:
            spec = yaml.safe_load(f)

        if not spec or not spec.get('tests'):
            print(f'  Part {part_num}: empty or no tests, skipping')
            continue

        cpp_out = cpp_dir / f'part{part_num}_test.cpp'
        java_out = java_dir / f'Part{part_num}Test.java'

        # C++
        if FORCE or not cpp_out.exists():
            cpp_code = generate_cpp_file(spec, display_name, part_num)
            cpp_out.write_text(cpp_code, encoding='utf-8')
            print(f'  Part {part_num}: -> {cpp_out.name}')
            generated += 1
        else:
            print(f'  Part {part_num}: skip {cpp_out.name} (exists, use --force to overwrite)')
            skipped += 1

        # Java
        if FORCE or not java_out.exists():
            java_code = generate_java_file(spec, display_name, part_num)
            java_out.write_text(java_code, encoding='utf-8')
            print(f'  Part {part_num}: -> {java_out.name}')
            generated += 1
        else:
            print(f'  Part {part_num}: skip {java_out.name} (exists, use --force to overwrite)')
            skipped += 1

    return generated, skipped


# ═══════════════════════════════════════════════════════════════════════════════
# Entry point
# ═══════════════════════════════════════════════════════════════════════════════

def main():
    args = [a for a in sys.argv[1:] if not a.startswith('--')]

    if '--all' in sys.argv:
        repo_root = Path(__file__).parent.parent
        dirs = sorted([
            d for d in (repo_root / 'problems').rglob('*')
            if d.is_dir() and (d / 'tests' / 'specs').exists()
               and any(d.glob('tests/specs/part*_spec.yaml'))
        ])
        if not dirs:
            print('No problems with tests/specs/partN_spec.yaml found.')
            return
        total_gen = total_skip = 0
        for d in dirs:
            print(f'=== {d.name} ===')
            g, s = process_problem(d)
            total_gen += g; total_skip += s
            print()
        print(f'Done. Generated: {total_gen}  Skipped: {total_skip}')
        return

    if not args:
        print(__doc__)
        sys.exit(0)

    for arg in args:
        print(f'=== Processing: {Path(arg).name} ===')
        process_problem(arg)
        print()


if __name__ == '__main__':
    main()
