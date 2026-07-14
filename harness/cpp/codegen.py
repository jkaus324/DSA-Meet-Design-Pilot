#!/usr/bin/env python3
"""
harness/cpp/codegen.py — emit a per-problem C++ test runner from spec.yaml + tests/cases/*.yaml

C++ has no eval and no convenient runtime YAML. Instead of vendoring yaml-cpp, we generate
a single self-contained runner.cpp at submit time. Each test case becomes inline C++ that
constructs the args, calls the function, and runs the comparison.

Used by dashboard/server.js — at submit time, server.js calls:
    python3 harness/cpp/codegen.py <problem-dir> <part-num> > <tmpdir>/runner.cpp
then compiles `runner.cpp` (which #includes solution.cpp).

Output contract: same PASS / FAIL / PART<N>_SUMMARY lines as the Python and Java runners.
"""
import sys, json, re
from pathlib import Path

import yaml


def cpp_lit(value, type_str, types, factories):
    """Render a YAML-loaded value as a C++ expression for the given declared type."""
    type_str = (type_str or "").strip()

    # list<X>
    m = re.match(r"list<(.+)>", type_str)
    if m:
        inner = m.group(1).strip()
        items = [cpp_lit(v, inner, types, factories) for v in (value or [])]
        if inner == "factory":
            iface = (cpp_lit.factory_iface_cpp or "RankingStrategy")
            return "std::vector<" + iface + "*>{" + ", ".join(items) + "}"
        return "std::vector<" + cpp_type(inner) + ">{" + ", ".join(items) + "}"

    # factory ref
    if type_str == "factory":
        if not isinstance(value, str) or value not in factories:
            raise ValueError(f"factory ref expected, got {value!r}")
        return factories[value]["cpp"]

    # struct
    if type_str in types:
        struct = types[type_str]
        fields = struct["fields"]
        parts = []
        for f in fields:
            fname, ftype = f["name"], f["type"]
            if fname in value:
                parts.append(cpp_lit(value[fname], ftype, types, factories))
            elif "default" in f:
                parts.append(cpp_lit(f["default"], ftype, types, factories))
            else:
                raise ValueError(f"missing field {fname} for {type_str}")
        return f"{type_str}{{{', '.join(parts)}}}"

    # primitives
    if type_str == "string":
        return json.dumps(value)
    if type_str == "bool":
        return "true" if value else "false"
    if type_str == "int":
        return str(int(value))
    if type_str == "float":
        # always render as double literal
        return repr(float(value))

    # fallback — print verbatim
    return str(value)


def cpp_type(type_str):
    if type_str == "string":
        return "std::string"
    if type_str == "bool":
        return "bool"
    if type_str == "int":
        return "int"
    if type_str == "float":
        return "double"
    m = re.match(r"list<(.+)>", type_str)
    if m:
        return "std::vector<" + cpp_type(m.group(1).strip()) + ">"
    return type_str  # custom struct name


def normalize_functions(spec):
    """spec.functions may be list-of-dicts (problem 001) or dict-of-dicts."""
    fns = spec.get("functions", [])
    if isinstance(fns, list):
        return {f["name"]: f for f in fns}
    return fns


def field_extract_cpp(var, field):
    """Read field/getter from a struct value."""
    return f"{var}.{field}"


def emit_test(test, idx, spec, types, factories, fns_by_name):
    name = test["name"]
    fn = fns_by_name[test["call"]]
    param_types = [p["type"] for p in fn.get("params", [])]
    args = test.get("args", [])
    expect_size = test.get("expect_size")
    expect_field = test.get("expect_field")
    expect = test.get("expect")
    expect_equals = test.get("expect_equals")  # scalar return value
    expect_close = test.get("expect_close")    # float scalar (epsilon compare)
    epsilon = test.get("epsilon", 0.001)
    also = test.get("also", [])
    expect_throws = test.get("expect_throws", False)

    arg_exprs = []
    for i, arg in enumerate(args):
        ptype = param_types[i] if i < len(param_types) else "auto"
        arg_exprs.append(cpp_lit(arg, ptype, types, factories))

    lines = []
    lines.append(f"    // {name}")
    lines.append("    try {")
    if expect_throws:
        lines.append("        bool _threw = false;")
        lines.append("        try {")
        lines.append(f"            {test['call']}({', '.join(arg_exprs)});")
        lines.append("        } catch (...) { _threw = true; }")
        lines.append("        if (!_threw) throw std::runtime_error(\"expected throw\");")
    else:
        ret_type = (fn.get("returns") or "").strip()
        if ret_type == "void":
            lines.append(f"        {test['call']}({', '.join(arg_exprs)});")
        else:
            lines.append(f"        auto _result = {test['call']}({', '.join(arg_exprs)});")
        if expect_equals is not None:
            ret_type = fn.get("returns", "")
            lit = cpp_lit(expect_equals, ret_type, types, factories)
            lines.append(f"        if (_result != {lit}) throw std::runtime_error(\"value mismatch\");")
        if expect_close is not None:
            lines.append(f"        if (std::fabs((double)_result - (double){repr(float(expect_close))}) > {repr(float(epsilon))}) throw std::runtime_error(\"close mismatch\");")
        if expect_size is not None:
            lines.append(f"        if ((int)_result.size() != {int(expect_size)}) throw std::runtime_error(\"size mismatch\");")
        if expect is not None:
            if expect_field:
                # extract field, compare elementwise
                lines.append(f"        const std::vector<decltype(_result[0].{expect_field})> _expected = {{")
                lit_items = []
                # infer element type from first value or from struct field type
                # find field type in struct
                ret_inner = re.match(r"list<(.+)>", fn.get("returns", ""))
                struct_t = ret_inner.group(1).strip() if ret_inner else None
                struct_def = types.get(struct_t, {})
                field_type = next((f["type"] for f in struct_def.get("fields", []) if f["name"] == expect_field), "string")
                for v in expect:
                    lit_items.append(cpp_lit(v, field_type, types, factories))
                lines.append("            " + ", ".join(lit_items))
                lines.append("        };")
                lines.append("        for (size_t _i = 0; _i < _expected.size(); ++_i) {")
                lines.append(f"            if (_result[_i].{expect_field} != _expected[_i]) throw std::runtime_error(\"field mismatch at \" + std::to_string(_i));")
                lines.append("        }")
            else:
                # expect is a list of primitive literals matching _result element-by-element
                ret_type = fn.get("returns", "")
                inner_m = re.match(r"list<(.+)>", ret_type)
                inner = inner_m.group(1).strip() if inner_m else ret_type
                lit_items = [cpp_lit(v, inner, types, factories) for v in expect]
                lines.append(f"        const std::vector<{cpp_type(inner)}> _expected = {{ {', '.join(lit_items)} }};")
                lines.append("        for (size_t _i = 0; _i < _expected.size(); ++_i) {")
                lines.append("            if (_result[_i] != _expected[_i]) throw std::runtime_error(\"value mismatch\");")
                lines.append("        }")
        for chk in also:
            i = chk["index"]
            f = chk["field"]
            v = chk["equals"]
            # infer field type
            ret_inner = re.match(r"list<(.+)>", fn.get("returns", ""))
            struct_t = ret_inner.group(1).strip() if ret_inner else None
            struct_def = types.get(struct_t, {})
            field_type = next((fd["type"] for fd in struct_def.get("fields", []) if fd["name"] == f), "string")
            lines.append(f"        if (_result[{int(i)}].{f} != {cpp_lit(v, field_type, types, factories)}) throw std::runtime_error(\"also-check failed\");")
    lines.append(f"        std::cout << \"PASS {name}\" << std::endl;")
    lines.append("        _passed++;")
    lines.append("    } catch (...) {")
    lines.append(f"        std::cout << \"FAIL {name}\" << std::endl;")
    lines.append("        _failed++;")
    lines.append("    }")
    return "\n".join(lines)


def emit(problem_dir, part_num):
    pdir = Path(problem_dir)
    spec = yaml.safe_load((pdir / "spec.yaml").read_text(encoding="utf-8"))
    cases = yaml.safe_load((pdir / "tests" / "cases" / f"part{part_num}.yaml").read_text(encoding="utf-8"))

    types = spec.get("types", {})
    factories = spec.get("factories", {})
    fns_by_name = normalize_functions(spec)
    # Pass the C++ factory interface name through cpp_lit (default RankingStrategy for back-compat).
    cpp_lit.factory_iface_cpp = (spec.get("factory_interface") or {}).get("cpp", "RankingStrategy")

    # solution.cpp lives at <pdir>/solution.cpp; in submit-time tmpdir, both are siblings
    body = []
    body.append("// generated — do not edit. See harness/cpp/codegen.py")
    body.append("#define RUNNING_TESTS")
    body.append('#include "solution.cpp"')
    body.append("#include <iostream>")
    body.append("#include <vector>")
    body.append("#include <string>")
    body.append("#include <stdexcept>")
    body.append("#include <cmath>")
    body.append("using namespace std;")
    body.append("")
    body.append(f"int main() {{")
    body.append("    int _passed = 0, _failed = 0;")
    body.append("")
    for i, t in enumerate(cases):
        body.append(emit_test(t, i, spec, types, factories, fns_by_name))
        body.append("")
    body.append(f"    std::cout << \"PART{part_num}_SUMMARY \" << _passed << \"/\" << (_passed + _failed) << std::endl;")
    body.append("    return _failed > 0 ? 1 : 0;")
    body.append("}")
    return "\n".join(body) + "\n"


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("usage: codegen.py <problem-dir> <part-num>", file=sys.stderr)
        sys.exit(2)
    # Write raw UTF-8 bytes — callers (dashboard/server.js) redirect this
    # process's stdout straight to a .cpp file via the shell, and on Windows
    # a text-mode print() falls back to the cp1252 locale encoding, which
    # mangles non-ASCII characters (e.g. the em-dash below) into invalid
    # UTF-8.
    sys.stdout.buffer.write(emit(sys.argv[1], int(sys.argv[2])).encode("utf-8"))
