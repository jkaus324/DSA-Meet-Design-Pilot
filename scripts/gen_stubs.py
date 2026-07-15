#!/usr/bin/env python3
"""Generate per-language boilerplate stubs from spec.yaml.
Layout: boilerplate/<lang>/part<N>/{interview,guided,learning}.<ext>
Modes: Interview (sigs+TODO), Guided (sigs+HINT), Learning (full scaffold).
Cumulative: partN carries forward all functions from prior parts.
"""

import argparse, re, sys
from pathlib import Path
import yaml

LANGS = ("cpp", "java", "python", "javascript", "go")
EXT = {"cpp": "cpp", "java": "java", "python": "py", "javascript": "js", "go": "go"}
MODES = ("interview", "guided", "learning")

def log(msg): print(msg, file=sys.stderr)


# ─── spec helpers ────────────────────────────────────────────────────────────

def normalize_functions(spec):
    fns = spec.get("functions", [])
    if isinstance(fns, list):
        return {f["name"]: f for f in fns}
    return fns


def parts_in_order(spec):
    raw = spec.get("parts", {}) or {}
    keys = sorted(raw.keys(), key=lambda k: int(k))
    return [(int(k), raw[k]) for k in keys]


def cumulative_fn_names(spec, part_num):
    """Carry-forward: return ordered fn names for partN, including all prior parts."""
    seen = []
    for n, p in parts_in_order(spec):
        if n > part_num:
            break
        for fn in p.get("functions", []) or []:
            if fn not in seen:
                seen.append(fn)
    return seen


def needs_factory_iface(spec, fn_names):
    """True if any function in scope takes a factory or list<factory> param."""
    fns = normalize_functions(spec)
    for name in fn_names:
        fn = fns.get(name, {})
        for p in fn.get("params", []) or []:
            t = p.get("type", "")
            if t == "factory" or parse_list(t) == "factory":
                return True
    return False


def parse_list(t):
    m = re.match(r"\s*list<(.+)>\s*$", t or "")
    return m.group(1).strip() if m else None


def factory_class(expr):
    """Extract class name from `new Foo()` / `Foo()`."""
    s = (expr or "").strip()
    s = s[4:].strip() if s.startswith("new ") else s
    m = re.match(r"([A-Za-z_][A-Za-z0-9_]*)", s)
    return m.group(1) if m else None


def factory_iface(spec, lang):
    """Optional spec field; falls back to heuristic name."""
    iface = (spec.get("factory_interface") or {}).get(lang)
    if iface:
        return iface
    # Heuristic: look at problem_id, no — just default to "Strategy".
    return "Strategy"


# ─── type rendering ──────────────────────────────────────────────────────────

def cpp_type(t, spec):
    inner = parse_list(t)
    if inner is not None:
        if inner == "factory":
            return f"vector<{factory_iface(spec, 'cpp')}*>"
        return f"vector<{cpp_type(inner, spec)}>"
    if t == "string":
        return "string"
    if t == "bool":
        return "bool"
    if t == "int":
        return "int"
    if t == "float":
        return "double"
    return t  # struct name passes through


def java_type(t, spec):
    inner = parse_list(t)
    if inner is not None:
        if inner == "factory":
            return f"List<{factory_iface(spec, 'java')}>"
        return f"List<{java_type(inner, spec)}>"
    if t == "string":
        return "String"
    if t == "bool":
        return "boolean"
    if t == "int":
        return "int"
    if t == "float":
        return "double"
    return t


def py_type(t, spec):
    inner = parse_list(t)
    if inner is not None:
        if inner == "factory":
            return f"list[{factory_iface(spec, 'python')}]"
        return f"list[{py_type(inner, spec)}]"
    if t == "string":
        return "str"
    if t == "bool":
        return "bool"
    if t == "int":
        return "int"
    if t == "float":
        return "float"
    return t


def js_type(t, spec):
    """JS is dynamically typed; we render types only as JSDoc hints."""
    inner = parse_list(t)
    if inner is not None:
        if inner == "factory":
            return f"{factory_iface(spec, 'javascript')}[]"
        return f"{js_type(inner, spec)}[]"
    if t == "string":
        return "string"
    if t == "bool":
        return "boolean"
    if t == "int":
        return "number"
    if t == "float":
        return "number"
    return t


def go_type(t, spec):
    inner = parse_list(t)
    if inner is not None:
        if inner == "factory":
            return f"[]{factory_iface(spec, 'go')}"
        return f"[]{go_type(inner, spec)}"
    if t == "string":
        return "string"
    if t == "bool":
        return "bool"
    if t == "int":
        return "int"
    if t == "float":
        return "float64"
    if t in ("void", "", None):
        return ""
    return t  # struct name passes through


def go_zero(t, spec):
    """Zero value expression for a declared return type (compiles cleanly)."""
    inner = parse_list(t)
    if inner is not None:
        return "nil"
    if t == "string":
        return '""'
    if t == "bool":
        return "false"
    if t == "int":
        return "0"
    if t == "float":
        return "0.0"
    return f"{t}{{}}"  # struct zero value


# ─── per-language signature emitters ─────────────────────────────────────────

def cpp_sig(fn, spec):
    ret = cpp_type(fn.get("returns", "void"), spec)
    params = ", ".join(
        f"{cpp_type(p['type'], spec)} {p['name']}" for p in fn.get("params", []) or []
    )
    return f"{ret} {fn['name']}({params})"


def java_sig(fn, spec):
    ret = java_type(fn.get("returns", "void"), spec)
    params = ", ".join(
        f"{java_type(p['type'], spec)} {p['name']}" for p in fn.get("params", []) or []
    )
    return f"public static {ret} {fn['name']}({params})"


def py_sig(fn, spec):
    params = ", ".join(p["name"] for p in fn.get("params", []) or [])
    return f"def {fn['name']}({params}):"


def js_sig(fn, spec):
    params = ", ".join(p["name"] for p in fn.get("params", []) or [])
    return f"function {fn['name']}({params})"


def go_sig(fn, spec):
    params = ", ".join(
        f"{p['name']} {go_type(p['type'], spec)}" for p in fn.get("params", []) or []
    )
    ret = go_type(fn.get("returns", "void"), spec)
    ret = f" {ret}" if ret else ""
    return f"func {fn['name']}({params}){ret}"


# ─── hint generation (generic, pattern-name-free) ────────────────────────────

def hint_for_fn(fn):
    name = fn["name"]
    has_factory_list = any(
        parse_list(p.get("type", "")) == "factory" for p in fn.get("params", []) or []
    )
    if has_factory_list:
        return "HINT: think about how to compose multiple criteria into a single decision."
    if any("filter" in name.lower() or "prefer" in p.get("name", "").lower()
           for p in fn.get("params", []) or []):
        return "HINT: a boolean flag changes ranking — handle it as a separate piece you can chain."
    return "HINT: pick the field that defines 'better' for this ranking and compare the two."


# ─── per-mode/per-lang emitters ──────────────────────────────────────────────

def header_block_cpp():
    return (
        "#include <iostream>\n"
        "#include <vector>\n"
        "#include <string>\n"
        "#include <algorithm>\n"
        "using namespace std;\n"
    )


def cpp_struct(types, spec):
    out = []
    for tname, tdef in (types or {}).items():
        out.append(f"struct {tname} {{")
        for f in tdef.get("fields", []):
            out.append(f"    {cpp_type(f['type'], spec)} {f['name']};")
        # constructor with full args (with defaults for fields that have one)
        ctor_params = []
        ctor_inits = []
        for f in tdef.get("fields", []):
            cpp_t = cpp_type(f["type"], spec)
            ref = f"const {cpp_t}&" if cpp_t in ("string",) else cpp_t
            piece = f"{ref} {f['name']}_"
            if "default" in f:
                if f["type"] == "bool":
                    piece += " = " + ("true" if f["default"] else "false")
                elif f["type"] == "string":
                    piece += f' = "{f["default"]}"'
                else:
                    piece += f" = {f['default']}"
            ctor_params.append(piece)
            ctor_inits.append(f"{f['name']}({f['name']}_)")
        out.append(f"    {tname}({', '.join(ctor_params)})")
        out.append(f"      : {', '.join(ctor_inits)} {{}}")
        out.append("};")
        out.append("")
    return "\n".join(out)


def java_pojo(types, spec):
    out = []
    for tname, tdef in (types or {}).items():
        fields = tdef.get("fields", [])
        out.append(f"class {tname} {{")
        for f in fields:
            out.append(f"    public {java_type(f['type'], spec)} {f['name']};")
        out.append("")
        # full ctor
        ctor_params = ", ".join(f"{java_type(f['type'], spec)} {f['name']}" for f in fields)
        out.append(f"    public {tname}({ctor_params}) {{")
        for f in fields:
            out.append(f"        this.{f['name']} = {f['name']};")
        out.append("    }")
        # short ctor (skip trailing default-having fields)
        short = [f for f in fields if "default" not in f]
        if len(short) != len(fields):
            short_params = ", ".join(f"{java_type(f['type'], spec)} {f['name']}" for f in short)
            short_args = []
            for f in fields:
                if "default" in f:
                    if f["type"] == "bool":
                        short_args.append("true" if f["default"] else "false")
                    elif f["type"] == "string":
                        short_args.append(f'"{f["default"]}"')
                    else:
                        short_args.append(str(f["default"]))
                else:
                    short_args.append(f["name"])
            out.append("")
            out.append(f"    public {tname}({short_params}) {{")
            out.append(f"        this({', '.join(short_args)});")
            out.append("    }")
        out.append("}")
        out.append("")
    return "\n".join(out)


def py_dataclass(types, spec):
    out = []
    for tname, tdef in (types or {}).items():
        fields = tdef.get("fields", [])
        params = []
        for f in fields:
            piece = f["name"]
            if "default" in f:
                d = f["default"]
                if f["type"] == "string":
                    piece += f'="{d}"'
                elif f["type"] == "bool":
                    piece += "=" + ("True" if d else "False")
                else:
                    piece += f"={d}"
            params.append(piece)
        out.append(f"class {tname}:")
        out.append(f"    def __init__(self, {', '.join(params)}):")
        for f in fields:
            out.append(f"        self.{f['name']} = {f['name']}")
        out.append("")
    return "\n".join(out)


def js_value_literal(f):
    """Render a default value as a JS literal."""
    d = f["default"]
    if f["type"] == "string":
        return f'"{d}"'
    if f["type"] == "bool":
        return "true" if d else "false"
    return str(d)


def js_class(types, spec):
    """ES class with a positional constructor in declared field order — this is
    the contract the JS runner relies on when deserializing test args."""
    out = []
    for tname, tdef in (types or {}).items():
        fields = tdef.get("fields", [])
        params = []
        for f in fields:
            piece = f["name"]
            if "default" in f:
                piece += f" = {js_value_literal(f)}"
            params.append(piece)
        out.append(f"class {tname} {{")
        out.append(f"  constructor({', '.join(params)}) {{")
        for f in fields:
            out.append(f"    this.{f['name']} = {f['name']};")
        out.append("  }")
        out.append("}")
        out.append("")
    return "\n".join(out)


def go_struct(types, spec):
    """Go structs with fields in declared order — the contract the Go codegen
    runner relies on (positional composite literals in declared field order)."""
    out = []
    for tname, tdef in (types or {}).items():
        out.append(f"type {tname} struct {{")
        for f in tdef.get("fields", []):
            out.append(f"\t{f['name']} {go_type(f['type'], spec)}")
        out.append("}")
        out.append("")
    return "\n".join(out)


# ─── interface-only declarations (used by interview/guided when needed) ─────

def cpp_factory_iface_decl(spec):
    iface = factory_iface(spec, "cpp")
    return (
        f"// Forward declaration so signatures compile; design and implement your own.\n"
        f"class {iface};\n"
    )


def java_factory_iface_decl(spec):
    iface = factory_iface(spec, "java")
    return (
        f"// Marker interface so signatures compile; you supply the methods.\n"
        f"interface {iface} {{}}\n"
    )


# ─── learning-mode scaffolds: empty factory classes + interface ──────────────

def factory_element_type(spec):
    """The struct type a factory/strategy compares — the element type of the
    list param that sits alongside a `list<factory>` param. Falls back to the
    first declared struct type, or None if there are no struct types."""
    fns = normalize_functions(spec)
    types = spec.get("types", {}) or {}
    for fn in fns.values():
        params = fn.get("params", []) or []
        if any(parse_list(p.get("type", "")) == "factory" for p in params):
            for p in params:
                inner = parse_list(p.get("type", ""))
                if inner and inner != "factory" and inner in types:
                    return inner
    return next(iter(types), None)


def cpp_learning_scaffold(spec, fn_names):
    """Only meaningful for factory/strategy problems. Returns '' when the spec
    declares no factories, so we don't emit a bogus interface for problems that
    have no strategy concept (and don't reference a type that may not exist)."""
    facs = spec.get("factories") or {}
    if not facs:
        return ""
    iface = factory_iface(spec, "cpp")
    elem = factory_element_type(spec) or "auto"
    out = []
    out.append(f"// {iface} — strategy interface (given). Implement compare() in each concrete type.")
    out.append(f"class {iface} {{")
    out.append("public:")
    out.append(f"    virtual bool compare(const {elem}& a, const {elem}& b) = 0;")
    out.append(f"    virtual ~{iface}() = default;")
    out.append("};")
    out.append("")
    out.append("// Concrete strategies — fill in compare() bodies.")
    for key, exprs in facs.items():
        cls = factory_class(exprs.get("cpp", ""))
        if not cls:
            continue
        out.append(f"class {cls} : public {iface} {{")
        out.append("public:")
        out.append(f"    bool compare(const {elem}& a, const {elem}& b) override {{")
        out.append("        // TODO: implement this")
        out.append("        return false;")
        out.append("    }")
        out.append("};")
        out.append("")
    return "\n".join(out)


def java_learning_scaffold(spec, fn_names):
    """Only meaningful for factory/strategy problems — returns '' otherwise so we
    don't emit an interface referencing a type that may not exist."""
    facs = spec.get("factories") or {}
    if not facs:
        return ""
    iface = factory_iface(spec, "java")
    elem = factory_element_type(spec) or "Object"
    out = []
    out.append(f"interface {iface} {{")
    out.append(f"    boolean compare({elem} a, {elem} b);")
    out.append("}")
    out.append("")
    for key, exprs in facs.items():
        cls = factory_class(exprs.get("java", ""))
        if not cls:
            continue
        out.append(f"class {cls} implements {iface} {{")
        out.append("    @Override")
        out.append(f"    public boolean compare({elem} a, {elem} b) {{")
        out.append("        // TODO: implement this")
        out.append("        return false;")
        out.append("    }")
        out.append("}")
        out.append("")
    return "\n".join(out)


def py_learning_scaffold(spec, fn_names):
    iface = factory_iface(spec, "python")
    facs = spec.get("factories") or {}
    out = []
    out.append("from abc import ABC, abstractmethod")
    out.append("")
    out.append(f"class {iface}(ABC):")
    out.append("    @abstractmethod")
    out.append("    def compare(self, a, b):")
    out.append('        """Return True iff `a` ranks strictly before `b`."""')
    out.append("")
    if facs:
        for key, exprs in facs.items():
            cls = factory_class(exprs.get("python", ""))
            if not cls:
                continue
            out.append(f"class {cls}({iface}):")
            out.append("    def compare(self, a, b):")
            out.append("        # TODO: implement this")
            out.append("        return False")
            out.append("")
    return "\n".join(out)


def js_learning_scaffold(spec, fn_names):
    iface = factory_iface(spec, "javascript")
    facs = spec.get("factories") or {}
    out = []
    out.append(f"// {iface} — base strategy. Subclasses implement compare().")
    out.append(f"class {iface} {{")
    out.append("  // Return true iff `a` ranks strictly before `b`.")
    out.append("  compare(a, b) { throw new Error('not implemented'); }")
    out.append("}")
    out.append("")
    if facs:
        out.append("// Concrete strategies — fill in compare() bodies.")
        for key, exprs in facs.items():
            cls = factory_class(exprs.get("javascript") or exprs.get("java", ""))
            if not cls:
                continue
            out.append(f"class {cls} extends {iface} {{")
            out.append("  compare(a, b) {")
            out.append("    // TODO: implement this")
            out.append("    return false;")
            out.append("  }")
            out.append("}")
            out.append("")
    return "\n".join(out)


# ─── full file emitters ──────────────────────────────────────────────────────

def emit_cpp(spec, part_num, mode):
    fn_names = cumulative_fn_names(spec, part_num)
    fns = normalize_functions(spec)
    types = spec.get("types", {}) or {}
    out = [header_block_cpp(), ""]

    if mode == "interview":
        out.append("// Data class (given).")
        out.append(cpp_struct(types, spec))
        if needs_factory_iface(spec, fn_names):
            out.append(cpp_factory_iface_decl(spec))
        out.append(f"// TODO: design and implement your solution.")
        out.append(f"// Required free functions:")
        for name in fn_names:
            out.append(f"//   {cpp_sig(fns[name], spec)};")
        out.append("")
        for name in fn_names:
            sig = cpp_sig(fns[name], spec)
            out.append(f"{sig} {{")
            out.append("    // TODO: write your solution")
            out.append(cpp_body_return(fns[name], spec))
            out.append("}")
            out.append("")
    elif mode == "guided":
        out.append("// Data class (given).")
        out.append(cpp_struct(types, spec))
        if needs_factory_iface(spec, fn_names):
            out.append(cpp_factory_iface_decl(spec))
        out.append("// HINT: introduce an abstraction so new ranking rules don't change existing code.")
        out.append("// HINT: keep the comparator small — one rule per class.")
        out.append("")
        for name in fn_names:
            fn = fns[name]
            out.append(f"// {hint_for_fn(fn)}")
            out.append(f"{cpp_sig(fn, spec)} {{")
            out.append("    // TODO: write your solution")
            out.append(cpp_body_return(fn, spec))
            out.append("}")
            out.append("")
    elif mode == "learning":
        out.append("// Data class (given — do not modify).")
        out.append(cpp_struct(types, spec))
        scaffold = cpp_learning_scaffold(spec, fn_names)
        if scaffold:
            out.append(scaffold)
        for name in fn_names:
            fn = fns[name]
            sig = cpp_sig(fn, spec)
            out.append(f"{sig} {{")
            out.append("    // TODO: implement this")
            out.append(cpp_body_return(fn, spec))
            out.append("}")
            out.append("")
    return "\n".join(out).rstrip() + "\n"


def emit_java(spec, part_num, mode):
    fn_names = cumulative_fn_names(spec, part_num)
    fns = normalize_functions(spec)
    types = spec.get("types", {}) or {}
    out = ["import java.util.*;", ""]

    if mode == "interview":
        out.append("// Data class (given).")
        out.append(java_pojo(types, spec))
        if needs_factory_iface(spec, fn_names):
            out.append(java_factory_iface_decl(spec))
        out.append("public class Solution {")
        for name in fn_names:
            sig = java_sig(fns[name], spec)
            out.append(f"    {sig} {{")
            out.append("        // TODO: write your solution")
            out.append(java_body_return(fns[name], spec))
            out.append("    }")
            out.append("")
        out.append("}")
    elif mode == "guided":
        out.append("// Data class (given).")
        out.append(java_pojo(types, spec))
        if needs_factory_iface(spec, fn_names):
            out.append(java_factory_iface_decl(spec))
        out.append("// HINT: introduce an abstraction so new ranking rules don't change existing code.")
        out.append("public class Solution {")
        for name in fn_names:
            fn = fns[name]
            out.append(f"    // {hint_for_fn(fn)}")
            out.append(f"    {java_sig(fn, spec)} {{")
            out.append("        // TODO: write your solution")
            out.append(java_body_return(fn, spec))
            out.append("    }")
            out.append("")
        out.append("}")
    elif mode == "learning":
        out.append("// Data class (given — do not modify).")
        out.append(java_pojo(types, spec))
        scaffold = java_learning_scaffold(spec, fn_names)
        if scaffold:
            out.append(scaffold)
        out.append("public class Solution {")
        for name in fn_names:
            fn = fns[name]
            out.append(f"    {java_sig(fn, spec)} {{")
            out.append("        // TODO: implement this")
            out.append(java_body_return(fn, spec))
            out.append("    }")
            out.append("")
        out.append("}")
    return "\n".join(out).rstrip() + "\n"


def emit_python(spec, part_num, mode):
    fn_names = cumulative_fn_names(spec, part_num)
    fns = normalize_functions(spec)
    types = spec.get("types", {}) or {}
    out = []

    if mode == "interview":
        out.append("# Data class (given).")
        out.append(py_dataclass(types, spec))
        for name in fn_names:
            out.append(py_sig(fns[name], spec))
            out.append("    # TODO: write your solution")
            out.append("    return methods" if any(p["name"] == "methods" for p in fns[name].get("params", [])) else "    return None")
            out.append("")
    elif mode == "guided":
        out.append("# Data class (given).")
        out.append(py_dataclass(types, spec))
        out.append("# HINT: introduce an abstraction so new ranking rules don't change existing code.")
        out.append("")
        for name in fn_names:
            fn = fns[name]
            out.append(f"# {hint_for_fn(fn)}")
            out.append(py_sig(fn, spec))
            out.append("    # TODO: write your solution")
            out.append("    return methods" if any(p["name"] == "methods" for p in fn.get("params", [])) else "    return None")
            out.append("")
    elif mode == "learning":
        out.append("# Data class (given — do not modify).")
        out.append(py_dataclass(types, spec))
        out.append(py_learning_scaffold(spec, fn_names))
        for name in fn_names:
            fn = fns[name]
            out.append(py_sig(fn, spec))
            out.append("    # TODO: implement this")
            out.append("    return methods" if any(p["name"] == "methods" for p in fn.get("params", [])) else "    return None")
            out.append("")
    return "\n".join(out).rstrip() + "\n"


def js_exports(spec, fn_names, mode="learning"):
    """Emit the CommonJS export block the JS runner relies on: every data class
    and all free functions, plus (in learning mode only) the factory interface +
    concrete strategy classes.

    The factory interface and strategy classes are only *defined* by the learning
    scaffold. In interview/guided mode the user designs their own, so exporting
    those names there would reference undefined identifiers and make the file
    crash on require() before any test runs — hence the mode gate."""
    types = spec.get("types", {}) or {}
    facs = spec.get("factories") or {}
    names = list(types.keys())
    if mode == "learning":
        if needs_factory_iface(spec, fn_names) or facs:
            iface = factory_iface(spec, "javascript")
            if iface not in names:
                names.append(iface)
        for key, exprs in facs.items():
            cls = factory_class(exprs.get("javascript") or exprs.get("java", ""))
            if cls and cls not in names:
                names.append(cls)
    names.extend(n for n in fn_names if n not in names)
    inner = ", ".join(names)
    hint = (
        "// If you add classes (e.g. strategy subclasses), add them here too.\n"
        if mode != "learning" else ""
    )
    return (
        "// ── Export everything the test runner needs (do not remove) ──\n"
        f"{hint}"
        f"module.exports = {{ {inner} }};"
    )


def emit_javascript(spec, part_num, mode):
    fn_names = cumulative_fn_names(spec, part_num)
    fns = normalize_functions(spec)
    types = spec.get("types", {}) or {}
    out = []

    if mode == "interview":
        out.append("// Data class (given).")
        out.append(js_class(types, spec))
        out.append("// TODO: design and implement your solution.")
        out.append("// Required functions:")
        for name in fn_names:
            out.append(f"//   {js_sig(fns[name], spec)}")
        out.append("")
        for name in fn_names:
            out.append(f"{js_sig(fns[name], spec)} {{")
            out.append("  // TODO: write your solution")
            out.append("  return methods;" if any(p["name"] == "methods" for p in fns[name].get("params", [])) else "  return null;")
            out.append("}")
            out.append("")
    elif mode == "guided":
        out.append("// Data class (given).")
        out.append(js_class(types, spec))
        out.append("// HINT: introduce an abstraction so new rules don't change existing code.")
        out.append("")
        for name in fn_names:
            fn = fns[name]
            out.append(f"// {hint_for_fn(fn)}")
            out.append(f"{js_sig(fn, spec)} {{")
            out.append("  // TODO: write your solution")
            out.append("  return methods;" if any(p["name"] == "methods" for p in fn.get("params", [])) else "  return null;")
            out.append("}")
            out.append("")
    elif mode == "learning":
        out.append("// Data class (given — do not modify).")
        out.append(js_class(types, spec))
        out.append(js_learning_scaffold(spec, fn_names))
        for name in fn_names:
            fn = fns[name]
            out.append(f"{js_sig(fn, spec)} {{")
            out.append("  // TODO: implement this")
            out.append("  return methods;" if any(p["name"] == "methods" for p in fn.get("params", [])) else "  return null;")
            out.append("}")
            out.append("")

    out.append(js_exports(spec, fn_names, mode))
    return "\n".join(out).rstrip() + "\n"


def cpp_body_return(fn, spec, indent="    "):
    """Return statement for a C++ stub. `void` → no return; a `methods` param
    echoes it; otherwise a value-initialized default (`return {};`)."""
    ret = (fn.get("returns") or "void").strip()
    if cpp_type(ret, spec) in ("", "void"):
        return f"{indent}// nothing to return"
    if any(p["name"] == "methods" for p in fn.get("params", []) or []):
        return f"{indent}return methods;"
    return f"{indent}return {{}};"


def java_body_return(fn, spec, indent="        "):
    """Return statement for a Java stub. `void` → no return; a `methods` param
    echoes it; otherwise `return null;` (or a zero for primitives)."""
    ret = (fn.get("returns") or "void").strip()
    jt = java_type(ret, spec)
    if jt in ("", "void"):
        return f"{indent}// nothing to return"
    if any(p["name"] == "methods" for p in fn.get("params", []) or []):
        return f"{indent}return methods;"
    if jt == "int":
        return f"{indent}return 0;"
    if jt == "double":
        return f"{indent}return 0.0;"
    if jt == "boolean":
        return f"{indent}return false;"
    return f"{indent}return null;"


def go_body_return(fn, spec):
    """Return statement for a Go stub. `void` → bare return; a param named
    `methods` (the common list echo) → return it; else the zero value."""
    ret = fn.get("returns", "void")
    if go_type(ret, spec) == "":
        return "\treturn"
    if any(p["name"] == "methods" for p in fn.get("params", []) or []):
        # only valid if the return type matches methods' type; fall back to zero
        for p in fn.get("params", []) or []:
            if p["name"] == "methods" and p["type"] == ret:
                return "\treturn methods"
    return f"\treturn {go_zero(ret, spec)}"


def emit_go(spec, part_num, mode):
    fn_names = cumulative_fn_names(spec, part_num)
    fns = normalize_functions(spec)
    types = spec.get("types", {}) or {}
    out = ["package main", ""]

    # A factory/strategy interface when functions in scope take list<factory>.
    iface_decl = []
    if needs_factory_iface(spec, fn_names):
        iface = factory_iface(spec, "go")
        iface_decl = [
            f"// {iface} — implement this interface with your own strategy types.",
            f"type {iface} interface {{",
            "\t// TODO: define the method(s) your strategies share.",
            "}",
            "",
        ]

    if mode == "interview":
        out.append("// Data class (given).")
        out.append(go_struct(types, spec))
        out.extend(iface_decl)
        out.append("// TODO: design and implement your solution.")
        out.append("// Required free functions:")
        for name in fn_names:
            out.append(f"//   {go_sig(fns[name], spec)}")
        out.append("")
        for name in fn_names:
            out.append(f"{go_sig(fns[name], spec)} {{")
            out.append("\t// TODO: write your solution")
            out.append(go_body_return(fns[name], spec))
            out.append("}")
            out.append("")
    elif mode == "guided":
        out.append("// Data class (given).")
        out.append(go_struct(types, spec))
        out.extend(iface_decl)
        out.append("// HINT: introduce an abstraction so new rules don't change existing code.")
        out.append("")
        for name in fn_names:
            fn = fns[name]
            out.append(f"// {hint_for_fn(fn)}")
            out.append(f"{go_sig(fn, spec)} {{")
            out.append("\t// TODO: write your solution")
            out.append(go_body_return(fn, spec))
            out.append("}")
            out.append("")
    elif mode == "learning":
        out.append("// Data class (given — do not modify).")
        out.append(go_struct(types, spec))
        out.extend(iface_decl)
        for name in fn_names:
            fn = fns[name]
            out.append(f"{go_sig(fn, spec)} {{")
            out.append("\t// TODO: implement this")
            out.append(go_body_return(fn, spec))
            out.append("}")
            out.append("")
    return "\n".join(out).rstrip() + "\n"


EMITTERS = {"cpp": emit_cpp, "java": emit_java, "python": emit_python, "javascript": emit_javascript, "go": emit_go}


def filename(lang, mode):
    """Java capitalizes; cpp/python lower."""
    base = mode.capitalize() if lang == "java" else mode
    return f"{base}.{EXT[lang]}"


# ─── per-problem driver ──────────────────────────────────────────────────────

def generate_problem(pdir, langs, force):
    pdir = Path(pdir)
    spec_path = pdir / "spec.yaml"
    if not spec_path.exists():
        log(f"skip {pdir}: no spec.yaml")
        return
    spec = yaml.safe_load(spec_path.read_text(encoding="utf-8"))
    parts = parts_in_order(spec)
    if not parts:
        log(f"skip {pdir}: no parts")
        return

    for lang in langs:
        if lang not in EMITTERS:
            log(f"skip lang {lang}: unsupported")
            continue
        for part_num, _ in parts:
            partdir = pdir / "boilerplate" / lang / f"part{part_num}"
            partdir.mkdir(parents=True, exist_ok=True)
            for mode in MODES:
                fpath = partdir / filename(lang, mode)
                if fpath.exists() and not force:
                    log(f"skip {fpath} (exists; use --force to overwrite)")
                    continue
                content = EMITTERS[lang](spec, part_num, mode)
                fpath.write_text(content, encoding="utf-8")
                log(f"wrote {fpath}")


def find_all_problems(root):
    root = Path(root)
    return sorted(p.parent for p in root.rglob("spec.yaml"))


# ─── CLI ─────────────────────────────────────────────────────────────────────

def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("problem_dir", nargs="?", help="Path to one problem dir.")
    ap.add_argument("--all", action="store_true", help="Process every problem with spec.yaml.")
    ap.add_argument("--lang", choices=LANGS, help="Restrict to one language.")
    ap.add_argument("--force", action="store_true", help="Overwrite existing stubs.")
    ap.add_argument("--root", default=None, help="Root for --all (default: cwd).")
    args = ap.parse_args()

    langs = (args.lang,) if args.lang else LANGS

    if args.all:
        root = Path(args.root) if args.root else Path(__file__).resolve().parent.parent / "problems"
        for pdir in find_all_problems(root):
            generate_problem(pdir, langs, args.force)
        return

    if not args.problem_dir:
        ap.error("provide <problem-dir> or --all")
    generate_problem(args.problem_dir, langs, args.force)


if __name__ == "__main__":
    main()
