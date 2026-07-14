#!/usr/bin/env node
// scripts/check-boilerplates.js
//
// Compiles every boilerplate in problems/**/boilerplate/cpp/*.cpp using
// g++ -fsyntax-only, so we catch missing entry symbols / typos in starter
// templates before users hit them. The dashboard's runtime test runner
// only catches problems that *the user* tries to solve — boilerplates
// that no one has touched can rot silently.
//
// This is a syntax-only check. It does NOT run tests. The boilerplates
// have empty TODO bodies — the goal is "does it parse and link as a
// translation unit", not "does it pass tests".
//
// Usage:
//   node scripts/check-boilerplates.js                # exit 1 on any failure
//   node scripts/check-boilerplates.js --tier 1       # only tier1 problems
//   node scripts/check-boilerplates.js --json
//
// Requires: g++ in PATH with C++17 support.

'use strict';

const fs = require('fs');
const path = require('path');
const { spawnSync } = require('child_process');

const REPO_ROOT = path.resolve(__dirname, '..');
const PROBLEMS_DIR = path.join(REPO_ROOT, 'problems');

const argv = process.argv.slice(2);
const AS_JSON = argv.includes('--json');
const TIER_FLAG = argv.indexOf('--tier');
const FILTER_TIER = TIER_FLAG >= 0 ? argv[TIER_FLAG + 1] : null;

const TIERS = ['tier1-foundation', 'tier2-intermediate', 'tier3-advanced', 'gauntlet']
  .filter(t => !FILTER_TIER || t.startsWith(`tier${FILTER_TIER}`));

// Verify g++ exists.
const which = spawnSync(process.platform === 'win32' ? 'where' : 'which', ['g++']);
if (which.status !== 0) {
  console.error('check-boilerplates: g++ not found in PATH. Skipping.');
  process.exit(1);
}

function listBoilerplates() {
  const out = [];
  for (const tier of TIERS) {
    const tierDir = path.join(PROBLEMS_DIR, tier);
    if (!fs.existsSync(tierDir)) continue;
    for (const problemDir of fs.readdirSync(tierDir)) {
      const cppDir = path.join(tierDir, problemDir, 'boilerplate', 'cpp');
      if (!fs.existsSync(cppDir)) continue;
      // Top-level: interview.cpp / guided.cpp / learning.cpp
      for (const mode of ['interview', 'guided', 'learning']) {
        const f = path.join(cppDir, `${mode}.cpp`);
        if (fs.existsSync(f)) out.push({ problem: problemDir, tier, mode, part: null, file: f });
      }
      // Per-part subdirs (partN/{interview,guided,learning}.cpp)
      for (const entry of fs.readdirSync(cppDir, { withFileTypes: true })) {
        if (!entry.isDirectory() || !/^part\d+$/.test(entry.name)) continue;
        for (const mode of ['interview', 'guided', 'learning']) {
          const f = path.join(cppDir, entry.name, `${mode}.cpp`);
          if (fs.existsSync(f)) out.push({ problem: problemDir, tier, mode, part: entry.name, file: f });
        }
      }
    }
  }
  return out;
}

function compileSyntax(file) {
  // -fsyntax-only: parse + semantic check, no codegen, no link.
  // -DRUNNING_TESTS so #ifndef RUNNING_TESTS guards work as users expect.
  const r = spawnSync('g++', ['-std=c++17', '-fsyntax-only', '-DRUNNING_TESTS', file], { encoding: 'utf8' });
  return { ok: r.status === 0, stderr: r.stderr || r.stdout || '' };
}

const targets = listBoilerplates();
const failures = [];
let passCount = 0;

if (!AS_JSON) {
  process.stdout.write(`check-boilerplates: ${targets.length} files, `);
}

for (const t of targets) {
  const { ok, stderr } = compileSyntax(t.file);
  if (ok) {
    passCount++;
  } else {
    failures.push({
      file: path.relative(REPO_ROOT, t.file),
      problem: t.problem, mode: t.mode, part: t.part,
      // Trim noise — keep first 4 stderr lines, max ~600 chars.
      err: stderr.split(/\r?\n/).slice(0, 4).join('\n').slice(0, 600),
    });
  }
}

if (AS_JSON) {
  console.log(JSON.stringify({ total: targets.length, pass: passCount, fail: failures.length, failures }, null, 2));
} else {
  console.log(`${passCount} pass, ${failures.length} fail`);
  if (failures.length > 0) {
    console.log('');
    for (const f of failures.slice(0, 20)) {
      console.log(`✗ ${f.file}`);
      console.log(`    ${f.err.split('\n').join('\n    ')}\n`);
    }
    if (failures.length > 20) console.log(`... and ${failures.length - 20} more failures`);
  }
}

process.exit(failures.length > 0 ? 1 : 0);
