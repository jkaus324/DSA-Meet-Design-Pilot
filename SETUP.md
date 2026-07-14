# Setup — runs on macOS, Windows, and Linux

This project is fully cross-platform. Paths and toolchain separators are resolved per-OS at
runtime (Node `path.delimiter` / Python `os.pathsep`), so there are **no hardcoded path issues** —
the same source works on a Mac, a Windows laptop, or a Dell/Linux box.

## 1. Prerequisites

**Required for the dashboard:**
- **Node.js 18+** — https://nodejs.org (`node --version`)

**Install only the language(s) you want to solve in.** The dashboard auto-detects which are
present and disables submit for the rest (you can still browse and write code):

| Language | Install | Verify |
|----------|---------|--------|
| C++      | g++ with C++17 (macOS: Xcode CLT · Windows: MinGW-w64 or MSYS2 · Linux: `build-essential`) | `g++ --version` |
| Go       | Go 1.21+ — https://go.dev/dl | `go version` |
| Java     | JDK 17+ (Temurin/OpenJDK) | `javac -version` |
| Python   | Python 3.9+ — https://python.org | `python3 --version` (Windows: `python --version`) |
| JavaScript | (uses the Node you already installed) | `node --version` |

> **Windows note:** the app detects common install locations for MinGW (`C:\msys64\mingw64\bin`,
> `C:\ProgramData\mingw64\...`), the JDK, and Python even if they aren't on your `PATH`.

## 2. Install & run

```bash
cd dashboard
npm install          # installs dependencies for YOUR platform (do NOT copy node_modules between machines)
npm run dev          # dashboard → http://localhost:5173 , API → http://localhost:3000
```

On Windows the exact same commands work in PowerShell or Command Prompt.

`npm install` runs a `postinstall` build automatically. For a production run instead of dev:

```bash
npm run build
npm start
```

## 3. (Optional) Verify every solution in every language

```bash
python3 scripts/stress_test.py            # all 20 problems × all 5 languages
python3 scripts/stress_test.py --lang go  # one language
python3 scripts/stress_test.py --problem 017-lru   # one problem
```

(On Windows, use `python` instead of `python3` if that's your launcher.)

## Why `node_modules` is not included

The zip deliberately omits `dashboard/node_modules` and `dashboard/dist`. `node_modules` contains
**platform-specific native binaries** (esbuild, rollup, fsevents), so copying it from one machine
to another would break the build. Running `npm install` on each machine fetches the right binaries
for that OS/CPU — this is what makes the project portable across Mac, Windows, and Dell/Linux.
