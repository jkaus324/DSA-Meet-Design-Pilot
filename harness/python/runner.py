#!/usr/bin/env python3
"""Generic Python test runner driven by spec.yaml + tests/cases/partN.yaml.

Invoked from dashboard/server.js as:
    python3 runner.py --problem-dir <tmpdir> --part <N>

Emits one `PASS <name>` or `FAIL <name>` per test case, then a
`PART<N>_SUMMARY <passed>/<total>` line.
"""

import argparse
import importlib.util
import os
import sys
import traceback

import yaml


def load_solution(problem_dir):
    sol_path = os.path.join(problem_dir, "solution.py")
    spec = importlib.util.spec_from_file_location("solution", sol_path)
    module = importlib.util.module_from_spec(spec)
    # Make the directory importable so any relative imports work.
    sys.path.insert(0, problem_dir)
    spec.loader.exec_module(module)
    return module


def normalize_functions(spec):
    """`functions` may be a list of {name,...} or a dict keyed by name."""
    fns = spec.get("functions", {})
    if isinstance(fns, list):
        return {f["name"]: f for f in fns}
    return fns


def parse_list_type(t):
    """Return the inner type for `list<X>`, else None."""
    if isinstance(t, str) and t.startswith("list<") and t.endswith(">"):
        return t[5:-1]
    return None


def deserialize_struct(value, type_name, types, factories, solution):
    """Instantiate a user-defined struct from a dict.

    Required fields (no `default` in spec) must be present in `value` and are
    passed positionally in declared order. Optional fields are passed as kwargs
    only when supplied — this lets callers omit any subset and rely on the
    class's default values, regardless of field ordering.
    """
    type_def = types[type_name]
    fields = type_def["fields"]
    cls = getattr(solution, type_name)
    pos_args = []
    kwargs = {}
    in_optional = False
    for f in fields:
        fname = f["name"]
        is_optional = "default" in f
        if not is_optional:
            if fname not in value:
                raise KeyError(f"missing required field '{fname}' for {type_name}")
            if in_optional:
                # Required fields after optional ones force kwarg-only mode.
                kwargs[fname] = value[fname]
            else:
                pos_args.append(value[fname])
        else:
            in_optional = True
            if fname in value:
                kwargs[fname] = value[fname]
    return cls(*pos_args, **kwargs)


def deserialize_arg(value, declared_type, types, factories, solution):
    """Coerce a YAML-loaded value into the actual Python argument."""
    inner = parse_list_type(declared_type)
    if inner is not None:
        if not isinstance(value, list):
            raise TypeError(f"expected list for {declared_type}, got {type(value)}")
        return [
            deserialize_arg(item, inner, types, factories, solution)
            for item in value
        ]

    if declared_type == "factory":
        # Factory key — eval the python expression in the solution namespace.
        expr = factories[value]["python"]
        return eval(expr, solution.__dict__)

    if declared_type in types and isinstance(value, dict):
        return deserialize_struct(value, declared_type, types, factories, solution)

    # Primitives (string/int/float/bool) and unknown types pass through.
    return value


def extract_field(obj, field):
    if isinstance(obj, dict):
        return obj.get(field)
    return getattr(obj, field)


def run_test(case, fn_specs, types, factories, solution):
    name = case["name"]
    fn_name = case["call"]
    fn_spec = fn_specs[fn_name]
    params = fn_spec["params"]

    raw_args = case.get("args", [])
    args = [
        deserialize_arg(raw_args[i], params[i]["type"], types, factories, solution)
        for i in range(len(raw_args))
    ]

    fn = getattr(solution, fn_name)

    expect_throws = case.get("expect_throws", False)
    try:
        result = fn(*args)
    except Exception:
        if expect_throws:
            return True, name
        return False, name

    if expect_throws:
        return False, name

    if "expect_equals" in case:
        if result != case["expect_equals"]:
            return False, name

    if "expect_close" in case:
        eps = case.get("epsilon", 0.001)
        try:
            if abs(float(result) - float(case["expect_close"])) > eps:
                return False, name
        except Exception:
            return False, name

    if "expect_size" in case:
        if len(result) != case["expect_size"]:
            return False, name

    if "expect" in case:
        field = case.get("expect_field")
        actual = (
            [extract_field(item, field) for item in result] if field else list(result)
        )
        if actual != case["expect"]:
            return False, name

    for check in case.get("also", []) or []:
        idx = check["index"]
        field = check["field"]
        expected = check["equals"]
        if extract_field(result[idx], field) != expected:
            return False, name

    return True, name


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--problem-dir", required=True)
    parser.add_argument("--part", required=True, type=int)
    args = parser.parse_args()

    with open(os.path.join(args.problem_dir, "spec.yaml")) as f:
        spec = yaml.safe_load(f)
    cases_path = os.path.join(
        args.problem_dir, "tests", "cases", f"part{args.part}.yaml"
    )
    with open(cases_path) as f:
        cases = yaml.safe_load(f) or []

    solution = load_solution(args.problem_dir)
    fn_specs = normalize_functions(spec)
    types = spec.get("types", {}) or {}
    factories = spec.get("factories", {}) or {}

    passed = 0
    total = len(cases)
    for case in cases:
        try:
            ok, name = run_test(case, fn_specs, types, factories, solution)
        except Exception:
            ok = False
            name = case.get("name", "<unnamed>")
            traceback.print_exc(file=sys.stderr)
        print(f"{'PASS' if ok else 'FAIL'} {name}")
        if ok:
            passed += 1

    print(f"PART{args.part}_SUMMARY {passed}/{total}")


if __name__ == "__main__":
    main()
