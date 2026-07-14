#!/usr/bin/env node
/* Generic JavaScript test runner driven by spec.yaml + tests/cases/partN.yaml.
 *
 * Invoked from dashboard/server.js as:
 *     node runner.js --problem-dir <tmpdir> --part <N>
 *
 * Emits one `PASS <name>` or `FAIL <name>` per test case, then a
 * `PART<N>_SUMMARY <passed>/<total>` line. This output contract is identical
 * to the Python/Java/C++ runners, so server.js's parseTestOutput consumes it
 * unchanged.
 *
 * Mirrors harness/python/runner.py semantics:
 *   • functions may be a list of {name,...} or a dict keyed by name
 *   • args are deserialized from YAML into solution class instances
 *   • the same assertion vocabulary is supported (expect / expect_equals /
 *     expect_close / expect_size / expect_field / also / expect_throws)
 *
 * Convention (must match scripts/gen_stubs.py JS emitter):
 *   solution.js exports every top-level class and function via CommonJS,
 *   and data-class constructors take positional args in declared field order.
 */

'use strict';

const fs = require('fs');
const path = require('path');
const yaml = require('js-yaml');

// ─── arg parsing ──────────────────────────────────────────────────────────────

function parseArgs(argv) {
  const out = {};
  for (let i = 0; i < argv.length; i++) {
    if (argv[i] === '--problem-dir') out.problemDir = argv[++i];
    else if (argv[i] === '--part') out.part = parseInt(argv[++i], 10);
  }
  if (!out.problemDir || !Number.isInteger(out.part)) {
    console.error('usage: node runner.js --problem-dir <dir> --part <N>');
    process.exit(2);
  }
  return out;
}

// ─── solution loading ───────────────────────────────────────────────────────

function loadSolution(problemDir) {
  const solPath = path.join(problemDir, 'solution.js');
  // require() caches by absolute path; each runner process is fresh so this is fine.
  const mod = require(solPath);
  if (!mod || typeof mod !== 'object') {
    throw new Error('solution.js must export its classes and functions (module.exports = { ... }).');
  }
  return mod;
}

// ─── spec helpers ───────────────────────────────────────────────────────────

// `functions` may be a list of {name,...} or a dict keyed by name.
function normalizeFunctions(spec) {
  const fns = spec.functions || {};
  if (Array.isArray(fns)) {
    const map = {};
    for (const f of fns) map[f.name] = f;
    return map;
  }
  return fns;
}

// Return the inner type for `list<X>`, else null.
function parseListType(t) {
  if (typeof t === 'string' && t.startsWith('list<') && t.endsWith('>')) {
    return t.slice(5, -1);
  }
  return null;
}

// ─── argument deserialization ───────────────────────────────────────────────

/* Instantiate a user-defined struct from a plain object.
 *
 * Fields are passed positionally in declared order. A field present in `value`
 * uses that value; a field absent but carrying a `default` in the spec falls
 * back to that default; an absent field with no default is a hard error. This
 * keeps construction order-stable and matches the positional ctor emitted by
 * gen_stubs.py for JS (and mirrors the C++/Java POJO convention). */
function deserializeStruct(value, typeName, types, factories, solution) {
  const typeDef = types[typeName];
  const fields = typeDef.fields || [];
  const Cls = solution[typeName];
  if (typeof Cls !== 'function') {
    throw new Error(`solution.js does not export class '${typeName}'`);
  }
  const args = [];
  for (const f of fields) {
    const fname = f.name;
    if (Object.prototype.hasOwnProperty.call(value, fname)) {
      args.push(value[fname]);
    } else if (Object.prototype.hasOwnProperty.call(f, 'default')) {
      args.push(f.default);
    } else {
      throw new Error(`missing required field '${fname}' for ${typeName}`);
    }
  }
  return new Cls(...args);
}

// Evaluate a factory expression (e.g. `new RewardsMaximizer()`) with every
// solution export in scope. JS factory exprs are syntactically identical to the
// Java/C++ ones, so we fall back through javascript → java → cpp.
function evalFactory(expr, solution) {
  const names = Object.keys(solution);
  const values = names.map(n => solution[n]);
  // eslint-disable-next-line no-new-func
  const fn = new Function(...names, `return (${expr});`);
  return fn(...values);
}

// Coerce a YAML-loaded value into the actual JS argument.
function deserializeArg(value, declaredType, types, factories, solution) {
  const inner = parseListType(declaredType);
  if (inner !== null) {
    if (!Array.isArray(value)) {
      throw new Error(`expected list for ${declaredType}, got ${typeof value}`);
    }
    return value.map(item => deserializeArg(item, inner, types, factories, solution));
  }

  if (declaredType === 'factory') {
    const entry = factories[value] || {};
    const expr = entry.javascript || entry.js || entry.java || entry.cpp;
    if (!expr) throw new Error(`factory '${value}' has no javascript/java/cpp expression`);
    return evalFactory(expr, solution);
  }

  if (types[declaredType] && value && typeof value === 'object' && !Array.isArray(value)) {
    return deserializeStruct(value, declaredType, types, factories, solution);
  }

  // Primitives (string/int/float/bool) and unknown types pass through.
  return value;
}

// ─── result inspection + comparison ──────────────────────────────────────────

function extractField(obj, field) {
  if (obj === null || obj === undefined) return undefined;
  return obj[field];
}

// Structural deep-equality matching Python's == for the value shapes that
// appear in test YAML (scalars, arrays, plain objects). Numbers compare with a
// tiny tolerance so YAML floats like 0.1 don't trip exact-equality.
function deepEqual(a, b) {
  if (a === b) return true;
  if (typeof a === 'number' && typeof b === 'number') {
    if (Number.isNaN(a) && Number.isNaN(b)) return true;
    return Math.abs(a - b) < 1e-9;
  }
  if (Array.isArray(a) && Array.isArray(b)) {
    if (a.length !== b.length) return false;
    for (let i = 0; i < a.length; i++) if (!deepEqual(a[i], b[i])) return false;
    return true;
  }
  if (a && b && typeof a === 'object' && typeof b === 'object') {
    const ka = Object.keys(a), kb = Object.keys(b);
    if (ka.length !== kb.length) return false;
    for (const k of ka) {
      if (!Object.prototype.hasOwnProperty.call(b, k)) return false;
      if (!deepEqual(a[k], b[k])) return false;
    }
    return true;
  }
  return false;
}

// ─── per-test execution ───────────────────────────────────────────────────────

function runTest(testCase, fnSpecs, types, factories, solution) {
  const name = testCase.name;
  const fnName = testCase.call;
  const fnSpec = fnSpecs[fnName];
  if (!fnSpec) throw new Error(`spec has no function '${fnName}'`);
  const params = fnSpec.params || [];

  const rawArgs = testCase.args || [];
  const args = rawArgs.map((a, i) =>
    deserializeArg(a, params[i].type, types, factories, solution));

  const fn = solution[fnName];
  if (typeof fn !== 'function') {
    throw new Error(`solution.js does not export function '${fnName}'`);
  }

  const expectThrows = testCase.expect_throws || false;
  let result;
  try {
    result = fn(...args);
  } catch (e) {
    return [!!expectThrows, name];
  }
  if (expectThrows) return [false, name];

  if (Object.prototype.hasOwnProperty.call(testCase, 'expect_equals')) {
    if (!deepEqual(result, testCase.expect_equals)) return [false, name];
  }

  if (Object.prototype.hasOwnProperty.call(testCase, 'expect_close')) {
    const eps = testCase.epsilon !== undefined ? testCase.epsilon : 0.001;
    const r = Number(result), e = Number(testCase.expect_close);
    if (!(Math.abs(r - e) <= eps)) return [false, name];
  }

  if (Object.prototype.hasOwnProperty.call(testCase, 'expect_size')) {
    if (!result || result.length !== testCase.expect_size) return [false, name];
  }

  if (Object.prototype.hasOwnProperty.call(testCase, 'expect')) {
    const field = testCase.expect_field;
    const actual = field
      ? Array.from(result, item => extractField(item, field))
      : Array.from(result);
    if (!deepEqual(actual, testCase.expect)) return [false, name];
  }

  for (const check of testCase.also || []) {
    if (extractField(result[check.index], check.field) !== check.equals) {
      // fall back to deep compare for object/array equals
      if (!deepEqual(extractField(result[check.index], check.field), check.equals)) {
        return [false, name];
      }
    }
  }

  return [true, name];
}

// ─── main ─────────────────────────────────────────────────────────────────────

function main() {
  const { problemDir, part } = parseArgs(process.argv.slice(2));

  // Use { json: true } so duplicate mapping keys are tolerated (last value
  // wins) instead of throwing — matching PyYAML's safe_load leniency used by
  // the Python runner. Some specs carry a duplicated `types:` block; the two
  // runners must agree on how that's handled.
  const loadYaml = (file) => yaml.load(fs.readFileSync(file, 'utf8'), { json: true });
  const spec = loadYaml(path.join(problemDir, 'spec.yaml'));
  const casesPath = path.join(problemDir, 'tests', 'cases', `part${part}.yaml`);
  const cases = loadYaml(casesPath) || [];

  const solution = loadSolution(problemDir);
  const fnSpecs = normalizeFunctions(spec);
  const types = spec.types || {};
  const factories = spec.factories || {};

  let passed = 0;
  const total = cases.length;
  for (const testCase of cases) {
    let ok, name;
    try {
      [ok, name] = runTest(testCase, fnSpecs, types, factories, solution);
    } catch (e) {
      ok = false;
      name = (testCase && testCase.name) || '<unnamed>';
      process.stderr.write(String(e && e.stack ? e.stack : e) + '\n');
    }
    process.stdout.write(`${ok ? 'PASS' : 'FAIL'} ${name}\n`);
    if (ok) passed++;
  }

  process.stdout.write(`PART${part}_SUMMARY ${passed}/${total}\n`);
}

main();
