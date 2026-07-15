#!/usr/bin/env python3
"""
harness/go/codegen.py — emit a per-problem Go test runner from spec.yaml + tests/cases/*.yaml

Go is compiled and (deliberately) has no runtime-YAML dependency in this harness, so — exactly
like the C++ harness — we generate a self-contained `runner.go` at submit time. Each test case
becomes inline Go that constructs the args, calls the function, and runs the comparison.

The generated `runner.go` is `package main` and is compiled alongside the user's `solution.go`
(also `package main`). Because they share a package, the runner can call the spec's snake_case
free functions and read unexported (lowercase) struct fields directly — no exported-name dance.

Used by dashboard/server.js — at submit time, server.js writes solution.go + a minimal go.mod
into a tmpdir, then:
    python3 harness/go/codegen.py <problem-dir> <part-num> > <tmpdir>/runner.go
and runs `go run .` in that dir (offline, stdlib-only).

Output contract: same PASS / FAIL / PART<N>_SUMMARY lines as the Python/Java/C++/JS runners.
"""
import sys
import json
import re
from pathlib import Path

import yaml


# ─── type + literal rendering ────────────────────────────────────────────────

def go_type(type_str, spec):
    type_str = (type_str or "").strip()
    if type_str == "string":
        return "string"
    if type_str == "bool":
        return "bool"
    if type_str == "int":
        return "int"
    if type_str == "float":
        return "float64"
    m = re.match(r"list<(.+)>", type_str)
    if m:
        inner = m.group(1).strip()
        if inner == "factory":
            return "[]" + factory_iface_go(spec)
        return "[]" + go_type(inner, spec)
    if type_str == "void" or type_str == "":
        return ""
    return type_str  # custom struct name


def factory_iface_go(spec):
    fi = (spec.get("factory_interface") or {})
    return fi.get("go") or fi.get("java") or "Strategy"


def go_lit(value, type_str, types, factories, spec):
    """Render a YAML-loaded value as a Go expression for the given declared type."""
    type_str = (type_str or "").strip()

    # list<X>
    m = re.match(r"list<(.+)>", type_str)
    if m:
        inner = m.group(1).strip()
        items = [go_lit(v, inner, types, factories, spec) for v in (value or [])]
        return go_type(type_str, spec) + "{" + ", ".join(items) + "}"

    # factory ref
    if type_str == "factory":
        if not isinstance(value, str) or value not in factories:
            raise ValueError(f"factory ref expected, got {value!r}")
        entry = factories[value]
        expr = entry.get("go")
        if not expr:
            raise ValueError(f"factory '{value}' has no `go:` expression in spec")
        return expr

    # struct — positional, all fields in declared order (value or spec default)
    if type_str in types:
        struct = types[type_str]
        fields = struct["fields"]
        parts = []
        for f in fields:
            fname, ftype = f["name"], f["type"]
            if isinstance(value, dict) and fname in value:
                parts.append(go_lit(value[fname], ftype, types, factories, spec))
            elif "default" in f:
                parts.append(go_lit(f["default"], ftype, types, factories, spec))
            else:
                raise ValueError(f"missing field {fname} for {type_str}")
        return f"{type_str}{{{', '.join(parts)}}}"

    # primitives
    if type_str == "string":
        return json.dumps(value)  # valid Go string literal
    if type_str == "bool":
        return "true" if value else "false"
    if type_str == "int":
        return str(int(value))
    if type_str == "float":
        return repr(float(value))  # always a decimal literal

    # fallback
    return json.dumps(value) if isinstance(value, str) else str(value)


def normalize_functions(spec):
    fns = spec.get("functions", [])
    if isinstance(fns, list):
        return {f["name"]: f for f in fns}
    return fns


def field_type_in_struct(types, struct_name, field_name, default="string"):
    struct = types.get(struct_name, {})
    for f in struct.get("fields", []):
        if f["name"] == field_name:
            return f["type"]
    return default


def scalar_eq_needs_epsilon(ret_type):
    return (ret_type or "").strip() == "float"


# ─── per-test emission ───────────────────────────────────────────────────────

def emit_test(test, spec, types, factories, fns_by_name, state):
    """Return (lines, needs_math). state tracks cross-test flags."""
    name = test["name"]
    fn = fns_by_name[test["call"]]
    param_types = [p["type"] for p in fn.get("params", [])]
    args = test.get("args", [])
    ret_type = (fn.get("returns") or "void").strip()

    expect_size = test.get("expect_size")
    expect_field = test.get("expect_field")
    expect = test.get("expect")
    expect_equals = test.get("expect_equals")
    expect_close = test.get("expect_close")
    epsilon = test.get("epsilon", 0.001)
    also = test.get("also", []) or []
    expect_throws = test.get("expect_throws", False)

    arg_exprs = []
    for i, arg in enumerate(args):
        ptype = param_types[i] if i < len(param_types) else "int"
        arg_exprs.append(go_lit(arg, ptype, types, factories, spec))
    call = f"{test['call']}({', '.join(arg_exprs)})"

    needs_math = False
    body = []
    has_checks = any(x is not None for x in (expect, expect_equals, expect_close, expect_size)) or bool(also)

    if expect_throws:
        # No panic-based "throws" cases appear in these problems, but support it:
        body.append("\t\t\t_threw := false")
        body.append("\t\t\tfunc() {")
        body.append("\t\t\t\tdefer func() { if recover() != nil { _threw = true } }()")
        body.append(f"\t\t\t\t_ = {call}")
        body.append("\t\t\t}()")
        body.append('\t\t\tif !_threw { panic("expected throw") }')
        return body, needs_math

    if ret_type == "void":
        body.append(f"\t\t\t{call}")
        return body, needs_math

    if not has_checks:
        body.append(f"\t\t\t_ = {call}")
        return body, needs_math

    body.append(f"\t\t\t_result := {call}")

    if expect_equals is not None:
        lit = go_lit(expect_equals, ret_type, types, factories, spec)
        if scalar_eq_needs_epsilon(ret_type):
            needs_math = True
            body.append(f"\t\t\tif math.Abs(float64(_result) - float64({lit})) > 1e-9 {{ panic(\"value mismatch\") }}")
        else:
            body.append(f"\t\t\tif _result != {lit} {{ panic(\"value mismatch\") }}")

    if expect_close is not None:
        needs_math = True
        body.append(f"\t\t\tif math.Abs(float64(_result) - {repr(float(expect_close))}) > {repr(float(epsilon))} {{ panic(\"close mismatch\") }}")

    if expect_size is not None:
        body.append(f"\t\t\tif len(_result) != {int(expect_size)} {{ panic(\"size mismatch\") }}")

    if expect is not None:
        inner_m = re.match(r"list<(.+)>", ret_type)
        inner = inner_m.group(1).strip() if inner_m else ret_type
        if expect_field:
            field_type = field_type_in_struct(types, inner, expect_field)
            lits = [go_lit(v, field_type, types, factories, spec) for v in expect]
            body.append(f"\t\t\t_expected := []{go_type(field_type, spec)}{{{', '.join(lits)}}}")
            body.append("\t\t\tif len(_result) != len(_expected) { panic(\"length mismatch\") }")
            body.append("\t\t\tfor _i := range _expected {")
            body.append(f"\t\t\t\tif _result[_i].{expect_field} != _expected[_i] {{ panic(\"field mismatch\") }}")
            body.append("\t\t\t}")
        else:
            lits = [go_lit(v, inner, types, factories, spec) for v in expect]
            body.append(f"\t\t\t_expected := []{go_type(inner, spec)}{{{', '.join(lits)}}}")
            body.append("\t\t\tif len(_result) != len(_expected) { panic(\"length mismatch\") }")
            body.append("\t\t\tfor _i := range _expected {")
            body.append("\t\t\t\tif _result[_i] != _expected[_i] { panic(\"value mismatch\") }")
            body.append("\t\t\t}")

    for chk in also:
        idx = int(chk["index"])
        field = chk["field"]
        val = chk["equals"]
        inner_m = re.match(r"list<(.+)>", ret_type)
        inner = inner_m.group(1).strip() if inner_m else ret_type
        field_type = field_type_in_struct(types, inner, field)
        lit = go_lit(val, field_type, types, factories, spec)
        body.append(f"\t\t\tif _result[{idx}].{field} != {lit} {{ panic(\"also-check failed\") }}")

    return body, needs_math


def emit(problem_dir, part_num):
    pdir = Path(problem_dir)
    spec = yaml.safe_load((pdir / "spec.yaml").read_text(encoding="utf-8"))
    cases = yaml.safe_load((pdir / "tests" / "cases" / f"part{part_num}.yaml").read_text(encoding="utf-8")) or []

    types = spec.get("types", {}) or {}
    factories = spec.get("factories", {}) or {}
    fns_by_name = normalize_functions(spec)

    test_blocks = []
    needs_math = False
    state = {}
    for t in cases:
        lines, nm = emit_test(t, spec, types, factories, fns_by_name, state)
        needs_math = needs_math or nm
        name = t["name"]
        block = []
        block.append(f'\tif _run("{name}", func() {{')
        block.extend(lines)
        block.append("\t}) { _passed++ } else { _failed++ }")
        test_blocks.append("\n".join(block))

    out = []
    out.append("// generated — do not edit. See harness/go/codegen.py")
    out.append("package main")
    out.append("")
    out.append("import (")
    out.append('\t"fmt"')
    if needs_math:
        out.append('\t"math"')
    out.append(")")
    out.append("")
    # _run executes a test closure, recovering panics into a FAIL.
    out.append("func _run(name string, fn func()) bool {")
    out.append("\tok := true")
    out.append("\tfunc() {")
    out.append("\t\tdefer func() { if recover() != nil { ok = false } }()")
    out.append("\t\tfn()")
    out.append("\t}()")
    out.append('\tif ok {')
    out.append('\t\tfmt.Println("PASS " + name)')
    out.append("\t} else {")
    out.append('\t\tfmt.Println("FAIL " + name)')
    out.append("\t}")
    out.append("\treturn ok")
    out.append("}")
    out.append("")
    out.append("func main() {")
    out.append("\t_passed, _failed := 0, 0")
    if not cases:
        out.append("\t_ = _run")  # keep _run referenced when there are no tests
    out.append("")
    for b in test_blocks:
        out.append(b)
        out.append("")
    out.append(f'\tfmt.Printf("PART{part_num}_SUMMARY %d/%d\\n", _passed, _passed+_failed)')
    out.append("\tif _failed > 0 { }")  # no-op; keep exit 0 (parser reads summary)
    out.append("}")
    return "\n".join(out) + "\n"


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("usage: codegen.py <problem-dir> <part-num>", file=sys.stderr)
        sys.exit(2)
    # Write raw UTF-8 bytes — callers (dashboard/server.js) redirect this
    # process's stdout straight to a .go file via the shell, and on Windows
    # a text-mode print() falls back to the cp1252 locale encoding, which
    # mangles non-ASCII characters (e.g. the em-dash below) into invalid
    # UTF-8. Go's compiler rejects invalid UTF-8 source outright.
    sys.stdout.buffer.write(emit(sys.argv[1], int(sys.argv[2])).encode("utf-8"))
