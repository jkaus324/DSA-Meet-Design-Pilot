#!/usr/bin/env python3
"""stress_test.py — verify every reference solution against the spec test cases,
the exact way dashboard/server.js runs a submission, for ALL languages.

This is the pilot repo's language-agnostic battle-tester. It stages
spec.yaml + tests/cases + the per-language runner into a temp dir and runs
every spec-declared part, asserting all cases PASS.

Usage:
    python3 scripts/stress_test.py                       # all problems, all langs
    python3 scripts/stress_test.py --lang go             # one language
    python3 scripts/stress_test.py --problem 017-lru     # substring match
    python3 scripts/stress_test.py --lang go --problem splitwise
"""
import argparse
import os
import shutil
import subprocess
import sys
import tempfile

import sys
import yaml

REPO_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
HARNESS = os.path.join(REPO_ROOT, "harness")
DASHBOARD_NODE_MODULES = os.path.join(REPO_ROOT, "dashboard", "node_modules")
SNAKEYAML_JAR = os.path.join(HARNESS, "java", "lib", "snakeyaml-2.2.jar")

IS_WINDOWS = os.name == "nt"
# The interpreter that runs this script is always a valid Python; on Windows the
# `python3` alias often doesn't exist, so reuse sys.executable everywhere.
PYTHON_CMD = sys.executable or ("python" if IS_WINDOWS else "python3")
# Native compiled-binary suffix (Windows needs .exe to run the file).
EXE = ".exe" if IS_WINDOWS else ""
ALL_LANGS = ["cpp", "go", "java", "python", "javascript"]


def find_problems():
    out = []
    pdir = os.path.join(REPO_ROOT, "problems")
    for tier in sorted(os.listdir(pdir)):
        tdir = os.path.join(pdir, tier)
        if not os.path.isdir(tdir):
            continue
        for pid in sorted(os.listdir(tdir)):
            d = os.path.join(tdir, pid)
            if os.path.isdir(d) and os.path.exists(os.path.join(d, "spec.yaml")):
                out.append((pid, d))
    return out


def spec_parts(pdir):
    with open(os.path.join(pdir, "spec.yaml")) as f:
        spec = yaml.safe_load(f)
    parts = spec.get("parts")
    if isinstance(parts, dict):
        return sorted(int(k) for k in parts.keys())
    cases_dir = os.path.join(pdir, "tests", "cases")
    if os.path.isdir(cases_dir):
        return sorted(
            int(f[4:-5]) for f in os.listdir(cases_dir)
            if f.startswith("part") and f.endswith(".yaml")
        )
    return []


def stage_common(pdir, tmp):
    shutil.copy(os.path.join(pdir, "spec.yaml"), tmp)
    cases_dst = os.path.join(tmp, "tests", "cases")
    os.makedirs(cases_dst, exist_ok=True)
    src_cases = os.path.join(pdir, "tests", "cases")
    for f in os.listdir(src_cases):
        if f.endswith(".yaml"):
            shutil.copy(os.path.join(src_cases, f), cases_dst)


def summarize(out):
    lines = [l for l in out.strip().split("\n") if l.strip()]
    fails = [l for l in lines if l.startswith("FAIL")]
    summary = next((l for l in lines if "_SUMMARY" in l), lines[-1] if lines else "")
    return summary, fails


def run_cpp(pdir, tmp, part):
    shutil.copy(os.path.join(pdir, "solution.cpp"), tmp)
    cg = subprocess.run(
        [PYTHON_CMD, os.path.join(HARNESS, "cpp", "codegen.py"), pdir, str(part)],
        capture_output=True, text=True,
    )
    if cg.returncode != 0:
        return None, [f"CODEGEN FAIL: {cg.stderr[-400:]}"]
    open(os.path.join(tmp, "r.cpp"), "w").write(cg.stdout)
    out_bin = os.path.join(tmp, "r" + EXE)
    comp = subprocess.run(
        ["g++", "-std=c++17", "-DRUNNING_TESTS", "-o", out_bin,
         os.path.join(tmp, "r.cpp")], capture_output=True, text=True,
    )
    if comp.returncode != 0:
        return None, [f"COMPILE FAIL: {comp.stderr[-800:]}"]
    rr = subprocess.run([out_bin], capture_output=True, text=True, timeout=30)
    return summarize(rr.stdout)


def run_python(pdir, tmp, part):
    shutil.copy(os.path.join(pdir, "solution.py"), tmp)
    shutil.copy(os.path.join(HARNESS, "python", "runner.py"), tmp)
    rr = subprocess.run(
        [PYTHON_CMD, os.path.join(tmp, "runner.py"), "--problem-dir", tmp, "--part", str(part)],
        capture_output=True, text=True, timeout=30,
    )
    return summarize(rr.stdout + rr.stderr if not rr.stdout.strip() else rr.stdout)


def run_javascript(pdir, tmp, part):
    shutil.copy(os.path.join(pdir, "solution.js"), tmp)
    shutil.copy(os.path.join(HARNESS, "javascript", "runner.js"), tmp)
    env = dict(os.environ, NODE_PATH=DASHBOARD_NODE_MODULES)
    rr = subprocess.run(
        ["node", os.path.join(tmp, "runner.js"), "--problem-dir", tmp, "--part", str(part)],
        capture_output=True, text=True, timeout=30, env=env,
    )
    return summarize(rr.stdout)


def run_java(pdir, tmp, part):
    shutil.copy(os.path.join(pdir, "solution.java"), os.path.join(tmp, "Solution.java"))
    shutil.copy(os.path.join(HARNESS, "java", "Runner.java"), tmp)
    cp = f"{tmp}{os.pathsep}{SNAKEYAML_JAR}"  # os.pathsep is ';' on Windows
    comp = subprocess.run(
        ["javac", "-cp", cp, os.path.join(tmp, "Solution.java"), os.path.join(tmp, "Runner.java")],
        capture_output=True, text=True,
    )
    if comp.returncode != 0:
        return None, [f"COMPILE FAIL: {comp.stdout[-800:]}{comp.stderr[-800:]}"]
    rr = subprocess.run(
        ["java", "-cp", cp, "Runner", "--problem-dir", tmp, "--part", str(part)],
        capture_output=True, text=True, timeout=30,
    )
    return summarize(rr.stdout)


def run_go(pdir, tmp, part):
    shutil.copy(os.path.join(pdir, "solution.go"), tmp)
    cg = subprocess.run(
        [PYTHON_CMD, os.path.join(HARNESS, "go", "codegen.py"), pdir, str(part)],
        capture_output=True, text=True,
    )
    if cg.returncode != 0:
        return None, [f"CODEGEN FAIL: {cg.stderr[-400:]}"]
    open(os.path.join(tmp, "runner.go"), "w").write(cg.stdout)
    # Minimal module so `go run` works offline, stdlib-only.
    open(os.path.join(tmp, "go.mod"), "w").write("module cjrun\n\ngo 1.21\n")
    env = dict(os.environ, GOPROXY="off", GOFLAGS="-mod=mod", CGO_ENABLED="0")
    rr = subprocess.run(
        ["go", "run", "."], cwd=tmp, capture_output=True, text=True, timeout=60, env=env,
    )
    if "runner.go" in rr.stderr and "PART" not in rr.stdout:
        return None, [f"BUILD FAIL: {rr.stderr[-800:]}"]
    return summarize(rr.stdout)


RUNNERS = {
    "cpp": run_cpp, "go": run_go, "java": run_java,
    "python": run_python, "javascript": run_javascript,
}
SOLUTION_FILE = {
    "cpp": "solution.cpp", "go": "solution.go", "java": "solution.java",
    "python": "solution.py", "javascript": "solution.js",
}


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--lang", choices=ALL_LANGS, help="restrict to one language")
    ap.add_argument("--problem", help="substring filter on problem id")
    args = ap.parse_args()

    langs = [args.lang] if args.lang else ALL_LANGS
    problems = find_problems()
    if args.problem:
        problems = [(pid, d) for pid, d in problems if args.problem in pid]

    total_ok = total = 0
    failures = []
    for pid, pdir in problems:
        parts = spec_parts(pdir)
        for lang in langs:
            sol = os.path.join(pdir, SOLUTION_FILE[lang])
            if not os.path.exists(sol):
                failures.append(f"{pid} [{lang}]: MISSING {SOLUTION_FILE[lang]}")
                print(f"✗ {pid:38s} {lang:11s} MISSING SOLUTION")
                total += 1
                continue
            ok_all = True
            detail = []
            for part in parts:
                tmp = tempfile.mkdtemp()
                try:
                    # Runtime-YAML runners (python/js/java) read spec + cases
                    # from the problem dir; C++/Go codegen bakes them in.
                    if lang in ("python", "javascript", "java"):
                        stage_common(pdir, tmp)
                    summary, fails = RUNNERS[lang](pdir, tmp, part)
                    if summary is None or fails:
                        ok_all = False
                        detail.append(f"part{part}: {fails[:3]}")
                    elif "/" in summary:
                        p, t = summary.split()[-1].split("/")
                        if p != t:
                            ok_all = False
                            detail.append(f"part{part}: {summary}")
                except Exception as e:
                    ok_all = False
                    detail.append(f"part{part}: EXC {e}")
                finally:
                    shutil.rmtree(tmp, ignore_errors=True)
            total += 1
            if ok_all:
                total_ok += 1
                print(f"✓ {pid:38s} {lang:11s} {len(parts)} parts OK")
            else:
                failures.append(f"{pid} [{lang}]: {detail}")
                print(f"✗ {pid:38s} {lang:11s} {detail}")

    print(f"\n{'='*70}\nRESULT: {total_ok}/{total} (problem×lang) combos fully pass")
    if failures:
        print(f"\n{len(failures)} FAILURES:")
        for f in failures:
            print("  " + f)
        sys.exit(1)
    print("ALL GREEN ✅")


if __name__ == "__main__":
    main()
