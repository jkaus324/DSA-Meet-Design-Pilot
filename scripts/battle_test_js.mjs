#!/usr/bin/env node
/* battle_test_js.mjs — verify JavaScript reference solutions against the spec
 * test cases, exactly the way dashboard/server.js runs a submission.
 *
 * For each problem with a solution.js, it stages spec.yaml + cases + the JS
 * runner into a temp dir and runs every part, asserting all cases PASS.
 *
 * Usage:
 *   node scripts/battle_test_js.mjs                 # all problems with solution.js
 *   node scripts/battle_test_js.mjs <problem-id>    # one problem (substring match)
 *   node scripts/battle_test_js.mjs --json
 */

import fs from 'fs';
import os from 'os';
import path from 'path';
import { fileURLToPath } from 'url';
import { execFileSync } from 'child_process';
import { createRequire } from 'module';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const REPO_ROOT = path.join(__dirname, '..');
const HARNESS = path.join(REPO_ROOT, 'harness', 'javascript', 'runner.js');
const NODE_MODULES = path.join(REPO_ROOT, 'dashboard', 'node_modules');
const require = createRequire(path.join(NODE_MODULES, 'noop.js'));
const yaml = require('js-yaml');

const argv = process.argv.slice(2);
const asJson = argv.includes('--json');
const filter = argv.find(a => !a.startsWith('--'));

function findProblems() {
  const out = [];
  for (const tier of fs.readdirSync(path.join(REPO_ROOT, 'problems'))) {
    const tierDir = path.join(REPO_ROOT, 'problems', tier);
    if (!fs.statSync(tierDir).isDirectory()) continue;
    for (const id of fs.readdirSync(tierDir)) {
      const dir = path.join(tierDir, id);
      if (!fs.statSync(dir).isDirectory()) continue;
      if (fs.existsSync(path.join(dir, 'spec.yaml'))) out.push({ id, dir });
    }
  }
  return out.sort((a, b) => a.id.localeCompare(b.id));
}

// Parts the SPEC actually declares — this is exactly what the dashboard server
// runs on submit. Orphaned tests/cases/partN.yaml files that the spec doesn't
// declare are not runnable by any language runner (the runner resolves arg
// types from spec.functions), so we don't count them. Falls back to on-disk
// case files only when a spec has no `parts` block.
function specParts(dir) {
  try {
    const spec = yaml.load(fs.readFileSync(path.join(dir, 'spec.yaml'), 'utf8'), { json: true });
    const parts = spec && spec.parts;
    if (parts && typeof parts === 'object') {
      return Object.keys(parts).map(Number).filter(Number.isInteger).sort((a, b) => a - b);
    }
  } catch (_) { /* fall through */ }
  const casesDir = path.join(dir, 'tests', 'cases');
  if (!fs.existsSync(casesDir)) return [];
  return fs.readdirSync(casesDir)
    .map(f => (f.match(/^part(\d+)\.yaml$/) || [])[1])
    .filter(Boolean)
    .map(Number)
    .sort((a, b) => a - b);
}

function testProblem(p) {
  const result = { id: p.id, parts: [], ok: false, error: null, hasSolution: false };
  const solPath = path.join(p.dir, 'solution.js');
  if (!fs.existsSync(solPath)) { result.error = 'no solution.js'; return result; }
  result.hasSolution = true;

  const tmp = fs.mkdtempSync(path.join(os.tmpdir(), `bt-${p.id}-`));
  try {
    fs.copyFileSync(path.join(p.dir, 'spec.yaml'), path.join(tmp, 'spec.yaml'));
    fs.copyFileSync(solPath, path.join(tmp, 'solution.js'));
    fs.copyFileSync(HARNESS, path.join(tmp, 'runner.js'));
    fs.mkdirSync(path.join(tmp, 'tests', 'cases'), { recursive: true });
    for (const f of fs.readdirSync(path.join(p.dir, 'tests', 'cases'))) {
      if (f.endsWith('.yaml')) {
        fs.copyFileSync(path.join(p.dir, 'tests', 'cases', f), path.join(tmp, 'tests', 'cases', f));
      }
    }

    let allOk = true;
    for (const part of specParts(p.dir)) {
      let stdout = '';
      try {
        stdout = execFileSync('node', [path.join(tmp, 'runner.js'), '--problem-dir', tmp, '--part', String(part)],
          { env: { ...process.env, NODE_PATH: NODE_MODULES }, timeout: 15000, encoding: 'utf8' });
      } catch (e) {
        stdout = (e.stdout || '') + (e.stderr || '');
      }
      const summary = (stdout.match(new RegExp(`PART${part}_SUMMARY (\\d+)/(\\d+)`)) || []);
      const passed = Number(summary[1] || 0);
      const total = Number(summary[2] || 0);
      const fails = (stdout.match(/^FAIL .*/gm) || []).map(l => l.slice(5));
      const partOk = total > 0 && passed === total;
      if (!partOk) allOk = false;
      result.parts.push({ part, passed, total, ok: partOk, fails });
    }
    result.ok = allOk && result.parts.length > 0;
  } catch (e) {
    result.error = String(e.message || e);
  } finally {
    fs.rmSync(tmp, { recursive: true, force: true });
  }
  return result;
}

const problems = findProblems().filter(p => !filter || p.id.includes(filter));
const results = problems.map(testProblem);

if (asJson) {
  console.log(JSON.stringify(results, null, 2));
} else {
  let pass = 0, fail = 0, missing = 0;
  for (const r of results) {
    if (!r.hasSolution) { missing++; console.log(`  --  ${r.id}  (no solution.js)`); continue; }
    if (r.ok) { pass++; console.log(`  ✓   ${r.id}  [${r.parts.map(p => `${p.passed}/${p.total}`).join(' ')}]`); }
    else {
      fail++;
      const detail = r.error ? r.error : r.parts.filter(p => !p.ok).map(p => `part${p.part} ${p.passed}/${p.total}${p.fails.length ? ' fails:' + p.fails.join(',') : ''}`).join('; ');
      console.log(`  ✗   ${r.id}  → ${detail}`);
    }
  }
  console.log(`\n${pass} passed, ${fail} failed, ${missing} missing solution.js  (of ${results.length})`);
  process.exit(fail > 0 ? 1 : 0);
}
