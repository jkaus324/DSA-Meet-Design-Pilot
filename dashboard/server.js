const express = require('express');
const path = require('path');
const fs = require('fs');
const { exec } = require('child_process');
const os = require('os');
const yaml = require('js-yaml');
const { marked } = require('marked');
const { markedHighlight } = require('marked-highlight');
const hljs = require('highlight.js');

// ─── Markdown setup ─────────────────────────────────────────────────────────

marked.use(
  markedHighlight({
    langPrefix: 'hljs language-',
    highlight(code, lang) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext';
      return hljs.highlight(code, { language }).value;
    },
  })
);

// ─── Paths ──────────────────────────────────────────────────────────────────

const REPO_ROOT      = path.join(__dirname, '..');
const PROBLEMS_YML   = path.join(REPO_ROOT, 'docs', '_data', 'problems.yml');
const PROGRESS_JSON  = path.join(REPO_ROOT, 'progress.json');
const PATTERNS_DIR   = path.join(REPO_ROOT, 'patterns');
const DIST_DIR       = path.join(__dirname, 'dist');

// ─── g++ availability check ─────────────────────────────────────────────────

// Add MinGW to PATH if it exists (Windows/Chocolatey installs)
const mingwPath = 'C:\\ProgramData\\mingw64\\mingw64\\bin';
if (fs.existsSync(path.join(mingwPath, 'g++.exe'))) {
  process.env.PATH = mingwPath + path.delimiter + (process.env.PATH || '');
}

let testRunnerAvailable = false;
exec('g++ --version', (err) => {
  if (err) {
    console.warn('⚠️  g++ not found. Test runner / Submit will not work.');
    console.warn('   Install: sudo apt install g++ (Linux) | xcode-select --install (Mac) | MinGW (Windows)');
    testRunnerAvailable = false;
  } else {
    testRunnerAvailable = true;
    console.log('✅ g++ found — test runner available.');
  }
});

// ─── Data helpers ───────────────────────────────────────────────────────────

function loadProblems() {
  if (!fs.existsSync(PROBLEMS_YML)) {
    console.error('ERROR: Could not find docs/_data/problems.yml');
    return [];
  }
  try {
    const raw = fs.readFileSync(PROBLEMS_YML, 'utf8');
    const all = yaml.load(raw) || [];
    // Filter out problems whose folder doesn't exist on disk
    return all.filter(p => {
      const problemDir = path.join(REPO_ROOT, p.path);
      if (!fs.existsSync(problemDir)) {
        console.warn(`WARNING: Problem "${p.id}" registered in problems.yml but folder missing: ${p.path}`);
        return false;
      }
      return true;
    });
  } catch (e) {
    console.error('ERROR: Failed to parse problems.yml:', e.message);
    return [];
  }
}

function loadProgress() {
  if (!fs.existsSync(PROGRESS_JSON)) {
    return createEmptyProgress();
  }
  try {
    const raw = fs.readFileSync(PROGRESS_JSON, 'utf8');
    let data = JSON.parse(raw);
    data = migrateProgress(data);
    return data;
  } catch (e) {
    const backupPath = PROGRESS_JSON + '.bak';
    console.error(`progress.json was corrupted. Backed up to progress.json.bak.`);
    try { fs.copyFileSync(PROGRESS_JSON, backupPath); } catch (_) {}
    const fresh = createEmptyProgress();
    saveProgress(fresh);
    return fresh;
  }
}

function createEmptyProgress() {
  return {
    version: 3,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    problems: {},
    primers_read: [],
    activity: [],
  };
}

function logActivity(progress, actionType) {
  if (!progress.activity) progress.activity = [];
  const today = new Date().toISOString().slice(0, 10);
  let entry = progress.activity.find(a => a.date === today);
  if (!entry) {
    entry = { date: today, actions: [], count: 0 };
    progress.activity.push(entry);
  }
  entry.actions.push(actionType);
  entry.count = entry.actions.length;
}

function calculateActivityStreak(progress) {
  const activityDates = new Set((progress.activity || []).filter(a => a.count > 0).map(a => a.date));

  // Also include legacy dates from problems
  for (const p of Object.values(progress.problems)) {
    if (p.started_at) activityDates.add(p.started_at.slice(0, 10));
    if (p.completed_at) activityDates.add(p.completed_at.slice(0, 10));
    if (p.parts) {
      for (const part of Object.values(p.parts)) {
        if (part.passed_at) activityDates.add(part.passed_at.slice(0, 10));
      }
    }
  }

  if (activityDates.size === 0) return 0;

  const sorted = [...activityDates].sort().reverse();
  const today = new Date().toISOString().slice(0, 10);
  const yesterday = new Date(Date.now() - 86400000).toISOString().slice(0, 10);

  if (sorted[0] !== today && sorted[0] !== yesterday) return 0;

  let streak = 1;
  for (let i = 1; i < sorted.length; i++) {
    const prev = new Date(sorted[i - 1]);
    const curr = new Date(sorted[i]);
    const diff = (prev - curr) / 86400000;
    if (diff === 1) streak++;
    else break;
  }
  return streak;
}

// ─── v2 → v3 migration ──────────────────────────────────────────────────────

function migrateProgress(data) {
  if ((data.version || 1) >= 3) return data;

  console.log('Migrating progress.json from v2 to v3...');
  const problems = loadProblems();

  for (const [id, entry] of Object.entries(data.problems || {})) {
    const problem = problems.find(p => p.id === id);
    const partCount = (problem && problem.parts) ? problem.parts.length : 1;

    // Build parts object from old status
    const oldStatus = entry.status || 'unsolved';
    const parts = {};

    if (oldStatus === 'solved') {
      // All parts passed
      for (let i = 1; i <= partCount; i++) {
        parts[String(i)] = {
          status: 'passed',
          passed_at: entry.completed_at || new Date().toISOString(),
          carry_forward: true,
        };
      }
    } else if (oldStatus === 'attempted') {
      parts['1'] = { status: 'attempted', passed_at: null, carry_forward: null };
      for (let i = 2; i <= partCount; i++) {
        parts[String(i)] = { status: 'locked', passed_at: null, carry_forward: null };
      }
    } else {
      // unsolved → part 1 active, rest locked
      parts['1'] = { status: 'active', passed_at: null, carry_forward: null };
      for (let i = 2; i <= partCount; i++) {
        parts[String(i)] = { status: 'locked', passed_at: null, carry_forward: null };
      }
    }

    delete entry.status;
    delete entry.ext1;
    delete entry.ext2;
    entry.parts = parts;
  }

  data.version = 3;
  saveProgress(data);
  console.log('Migration complete.');
  return data;
}

function saveProgress(data) {
  data.updated_at = new Date().toISOString();
  fs.writeFileSync(PROGRESS_JSON, JSON.stringify(data, null, 2), 'utf8');
}

function loadPrimers() {
  if (!fs.existsSync(PATTERNS_DIR)) return [];
  return fs.readdirSync(PATTERNS_DIR)
    .filter(f => f.endsWith('.md') && f !== 'README.md')
    .map(f => ({
      name: path.basename(f, '.md'),
      file: path.join(PATTERNS_DIR, f),
    }));
}

function renderMarkdown(filePath) {
  if (!fs.existsSync(filePath)) return null;
  const raw = fs.readFileSync(filePath, 'utf8');
  return marked(raw);
}

function extractBeforeYouCode(rawMarkdown) {
  const lines = rawMarkdown.split('\n');
  let start = -1, end = -1;
  for (let i = 0; i < lines.length; i++) {
    if (start === -1 && /^#+\s*Before You Code/i.test(lines[i])) {
      start = i;
    } else if (start !== -1 && i > start && /^#+\s/.test(lines[i])) {
      end = i;
      break;
    }
  }
  if (start === -1) return null;
  return marked(lines.slice(start, end === -1 ? undefined : end).join('\n'));
}

function stripBeforeYouCode(rawMarkdown) {
  const lines = rawMarkdown.split('\n');
  let start = -1, end = -1;
  for (let i = 0; i < lines.length; i++) {
    if (start === -1 && /^#+\s*Before You Code/i.test(lines[i])) {
      start = i;
    } else if (start !== -1 && i > start && /^#+\s/.test(lines[i])) {
      end = i;
      break;
    }
  }
  if (start === -1) return rawMarkdown;
  const kept = [...lines.slice(0, start), ...(end === -1 ? [] : lines.slice(end))];
  return kept.join('\n');
}

// ─── Split README by part markers ───────────────────────────────────────────

function splitReadmeByParts(rawMarkdown, partDefs) {
  if (!partDefs || partDefs.length === 0) return { scenario: rawMarkdown, sections: [rawMarkdown] };

  const lines = rawMarkdown.split('\n');
  const markerIndices = partDefs.map(p => {
    for (let i = 0; i < lines.length; i++) {
      if (lines[i].trim().startsWith(p.description_marker)) return i;
    }
    return -1;
  });

  // Everything before the first found marker = scenario/context
  const firstMarker = markerIndices.find(i => i !== -1);
  const scenarioLines = firstMarker !== undefined && firstMarker > 0
    ? lines.slice(0, firstMarker)
    : lines;
  const scenario = scenarioLines.join('\n').trim();

  const sections = [];
  for (let i = 0; i < markerIndices.length; i++) {
    const start = markerIndices[i];
    if (start === -1) {
      sections.push(null);
      continue;
    }
    const end = markerIndices[i + 1] !== undefined && markerIndices[i + 1] !== -1
      ? markerIndices[i + 1]
      : lines.length;
    sections.push(lines.slice(start, end).join('\n'));
  }

  return { scenario, sections };
}

// ─── Derive problem status from parts ───────────────────────────────────────

function deriveStatus(partsObj, totalParts) {
  if (!partsObj || totalParts === 0) return 'unsolved';
  const partsList = Object.values(partsObj);
  const allPassed = partsList.length === totalParts && partsList.every(p => p.status === 'passed');
  if (allPassed) return 'solved';
  const anyActive = partsList.some(p => p.status === 'attempted' || p.status === 'passed');
  if (anyActive) return 'attempted';
  return 'unsolved';
}

function ensurePartsInitialized(entry, totalParts) {
  if (!entry.parts) {
    entry.parts = {};
    entry.parts['1'] = { status: 'active', passed_at: null, carry_forward: null };
    for (let i = 2; i <= totalParts; i++) {
      entry.parts[String(i)] = { status: 'locked', passed_at: null, carry_forward: null };
    }
  }
  return entry;
}

function mergeProblemWithProgress(problem, progress, primersRead) {
  const p = progress.problems[problem.id] || {};
  const totalParts = (problem.parts || []).length || 1;
  const partsProgress = p.parts || null;

  ensurePartsInitialized(p, totalParts);
  const status = deriveStatus(p.parts, totalParts);

  // Compute parts_passed and current_part
  let partsPassed = 0;
  let currentPart = 1;
  if (p.parts) {
    for (let i = 1; i <= totalParts; i++) {
      const ps = (p.parts[String(i)] || {}).status;
      if (ps === 'passed') {
        partsPassed++;
      } else if (ps !== 'locked' && currentPart === 1) {
        currentPart = i;
      }
    }
    // If all passed, currentPart = totalParts (last)
    if (partsPassed === totalParts) currentPart = totalParts;
    else {
      // Find first non-passed, non-locked
      for (let i = 1; i <= totalParts; i++) {
        const ps = (p.parts[String(i)] || {}).status;
        if (ps !== 'passed' && ps !== 'locked') { currentPart = i; break; }
        if (ps === 'locked') { currentPart = i - 1 || 1; break; }
      }
    }
  }

  const primerName = problem.prerequisite_primer;
  return {
    ...problem,
    status,
    difficulty_mode:  p.difficulty_mode || 'interview',
    started_at:       p.started_at      || null,
    completed_at:     p.completed_at    || null,
    notes:            p.notes           || '',
    total_parts:      totalParts,
    parts_passed:     partsPassed,
    current_part:     partsPassed === totalParts ? null : currentPart,
    primer_read:      primerName ? (primersRead || []).includes(primerName) : null,
  };
}

function buildSummary(mergedProblems) {
  const total     = mergedProblems.length;
  const solved    = mergedProblems.filter(p => p.status === 'solved').length;
  const attempted = mergedProblems.filter(p => p.status === 'attempted').length;
  const unsolved  = total - solved - attempted;
  return { total, solved, attempted, unsolved };
}

function calculateStreak(progress) {
  const dates = new Set();
  for (const p of Object.values(progress.problems)) {
    if (p.started_at)   dates.add(p.started_at.slice(0, 10));
    if (p.completed_at) dates.add(p.completed_at.slice(0, 10));
    if (p.parts) {
      for (const part of Object.values(p.parts)) {
        if (part.passed_at) dates.add(part.passed_at.slice(0, 10));
      }
    }
  }
  if (dates.size === 0) return { current_days: 0, last_activity: null };

  const sorted = [...dates].sort().reverse();
  const lastActivity = sorted[0];

  const today = new Date().toISOString().slice(0, 10);
  const yesterday = new Date(Date.now() - 86400000).toISOString().slice(0, 10);

  if (lastActivity !== today && lastActivity !== yesterday) {
    return { current_days: 0, last_activity: lastActivity };
  }

  let streak = 1;
  for (let i = 1; i < sorted.length; i++) {
    const prev = new Date(sorted[i - 1]);
    const curr = new Date(sorted[i]);
    const diff = (prev - curr) / 86400000;
    if (diff === 1) streak++;
    else break;
  }
  return { current_days: streak, last_activity: lastActivity };
}

// ─── Express app ────────────────────────────────────────────────────────────

const app = express();
app.use(express.json());

// ── GET /api/problems ────────────────────────────────────────────────────────

app.get('/api/problems', (req, res) => {
  const problems    = loadProblems();
  const progress    = loadProgress();
  const primersRead = progress.primers_read || [];
  const merged      = problems.map(p => mergeProblemWithProgress(p, progress, primersRead));
  const summary     = buildSummary(merged);
  res.json({ problems: merged, summary });
});

// ── POST /api/problems/:id/status ────────────────────────────────────────────
// v3: only accepts notes and difficulty_mode (status is derived from parts)

app.post('/api/problems/:id/status', (req, res) => {
  const { id } = req.params;
  const { notes, difficulty_mode } = req.body;

  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const progress = loadProgress();
  const entry    = progress.problems[id] || {};

  if (difficulty_mode !== undefined && ['interview', 'guided', 'learning'].includes(difficulty_mode)) {
    entry.difficulty_mode = difficulty_mode;
  }
  if (notes !== undefined) {
    entry.notes = notes;
  }

  const totalParts = (problem.parts || []).length || 1;
  ensurePartsInitialized(entry, totalParts);
  progress.problems[id] = entry;
  saveProgress(progress);

  const primersRead = progress.primers_read || [];
  res.json(mergeProblemWithProgress(problem, progress, primersRead));
});

// ── GET /api/problems/:id/parts ← NEW ────────────────────────────────────────

app.get('/api/problems/:id/parts', (req, res) => {
  const { id } = req.params;
  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const progress   = loadProgress();
  const entry      = progress.problems[id] || {};
  const partDefs   = problem.parts || [{ name: 'Complete', description_marker: '## Part 1', test_file: 'part1_test.cpp' }];
  const totalParts = partDefs.length;

  ensurePartsInitialized(entry, totalParts);
  const partsProgress = entry.parts || {};

  // Read README and split by markers
  const readmePath = path.join(REPO_ROOT, problem.path, 'README.md');
  let rawReadme = '';
  if (fs.existsSync(readmePath)) {
    rawReadme = fs.readFileSync(readmePath, 'utf8');
    rawReadme = stripBeforeYouCode(rawReadme);
  }
  const { scenario, sections } = splitReadmeByParts(rawReadme, partDefs);

  const parts = partDefs.map((def, i) => {
    const partNum    = i + 1;
    const partKey    = String(partNum);
    const partProg   = partsProgress[partKey] || { status: 'locked', passed_at: null, carry_forward: null };
    const isLocked   = partProg.status === 'locked';
    const section    = sections[i];

    // Count tests in test file
    let testCount = 0;
    if (!isLocked) {
      const testPath = path.join(REPO_ROOT, problem.path, 'tests', 'cpp', def.test_file);
      if (fs.existsSync(testPath)) {
        const testContent = fs.readFileSync(testPath, 'utf8');
        const matches = testContent.match(/cout\s*<<\s*"(PASS|FAIL)\s/g);
        testCount = matches ? matches.length / 2 : 0; // PASS + FAIL per test
      }
    }

    return {
      part:             partNum,
      name:             def.name,
      status:           partProg.status,
      passed_at:        partProg.passed_at,
      carry_forward:    partProg.carry_forward,
      description_html: isLocked ? null : (section ? marked(section) : null),
      test_count:       testCount,
    };
  });

  res.json({
    id,
    total_parts: totalParts,
    scenario_html: scenario ? marked(scenario) : null,
    parts,
    runner_available: testRunnerAvailable,
  });
});

// ── GET /api/problems/:id/starter ────────────────────────────────────────────
// Updated: accepts ?mode= and ?part= query params

app.get('/api/problems/:id/starter', (req, res) => {
  const { id } = req.params;
  const mode = ['interview', 'guided', 'learning'].includes(req.query.mode)
    ? req.query.mode : 'interview';
  const part = parseInt(req.query.part, 10) || 1;

  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const tryFile = (p, m) => path.join(REPO_ROOT, problem.path, 'boilerplate', 'cpp', `part${p}`, `${m}.cpp`);

  let filePath = tryFile(part, mode);
  let fallback = false;

  if (!fs.existsSync(filePath)) {
    filePath = tryFile(part, 'learning');
    fallback = true;
  }
  if (!fs.existsSync(filePath)) {
    filePath = tryFile(1, mode);
    fallback = true;
  }
  if (!fs.existsSync(filePath)) {
    filePath = tryFile(1, 'learning');
    fallback = true;
  }

  const code = fs.existsSync(filePath)
    ? fs.readFileSync(filePath, 'utf8')
    : `// Starter file not found for mode: ${mode}, part: ${part}\n// Write your solution here.\n`;

  res.json({ mode, part, code, fallback });
});

// ── GET /api/problems/:id/code ────────────────────────────────────────────────

app.get('/api/problems/:id/code', (req, res) => {
  const { id } = req.params;
  const mode = ['interview', 'guided', 'learning'].includes(req.query.mode)
    ? req.query.mode : 'interview';

  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const progress  = loadProgress();
  const entry     = progress.problems[id] || {};
  const saved     = (entry.code || {})[mode] || '';

  if (saved) {
    return res.json({ mode, code: saved, is_starter: false });
  }

  // Fall back to part1 starter
  const tryFile = (m) => path.join(REPO_ROOT, problem.path, 'boilerplate', 'cpp', 'part1', `${m}.cpp`);
  let filePath = tryFile(mode);
  if (!fs.existsSync(filePath)) filePath = tryFile('learning');

  const code = fs.existsSync(filePath)
    ? fs.readFileSync(filePath, 'utf8')
    : `// Write your solution here.\n`;

  res.json({ mode, code, is_starter: true });
});

// ── POST /api/problems/:id/code ───────────────────────────────────────────────

app.post('/api/problems/:id/code', (req, res) => {
  const { id } = req.params;
  const { mode, code } = req.body;

  if (!['interview', 'guided', 'learning'].includes(mode)) {
    return res.status(400).json({ error: 'mode must be interview, guided, or learning' });
  }
  if (typeof code !== 'string') return res.status(400).json({ error: 'code must be a string' });

  const problems = loadProblems();
  if (!problems.find(p => p.id === id)) return res.status(404).json({ error: `Problem ${id} not found` });

  const progress = loadProgress();
  if (!progress.problems[id]) progress.problems[id] = {};
  if (!progress.problems[id].code) progress.problems[id].code = {};
  progress.problems[id].code[mode] = code;
  logActivity(progress, 'saved_code');
  saveProgress(progress);

  res.json({ saved: true, mode });
});

// ── POST /api/problems/:id/submit ← NEW ──────────────────────────────────────

app.post('/api/problems/:id/submit', (req, res) => {
  const { id } = req.params;
  const { part, mode, code } = req.body;

  if (typeof code !== 'string') return res.status(400).json({ error: 'code must be a string' });
  if (!Number.isInteger(part) || part < 1) return res.status(400).json({ error: 'part must be a positive integer' });

  if (!testRunnerAvailable) {
    return res.status(503).json({
      error: 'g++ not available. Install g++ to use the test runner.',
      runner_available: false,
    });
  }

  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const partDefs   = problem.parts || [{ name: 'Complete', description_marker: '## Part 1', test_file: 'part1_test.cpp' }];
  const totalParts = partDefs.length;

  if (part > totalParts) {
    return res.status(400).json({ error: `Part ${part} does not exist. Problem has ${totalParts} parts.` });
  }

  // Save code first
  const progress = loadProgress();
  if (!progress.problems[id]) progress.problems[id] = {};
  const entry = progress.problems[id];
  if (!entry.code) entry.code = {};
  entry.code[mode || 'interview'] = code;
  ensurePartsInitialized(entry, totalParts);

  // Set started_at if not already set
  if (!entry.started_at) {
    entry.started_at = new Date().toISOString();
  }
  // Mark current part as attempted
  if ((entry.parts[String(part)] || {}).status !== 'passed') {
    entry.parts[String(part)] = { ...entry.parts[String(part)], status: 'attempted' };
  }

  saveProgress(progress);

  // Create temp directory
  const timestamp = Date.now();
  const tmpDir    = path.join(os.tmpdir(), `dsa-md-${id}-${timestamp}`);
  fs.mkdirSync(tmpDir, { recursive: true });

  const solutionFile = path.join(tmpDir, 'solution.cpp');
  fs.writeFileSync(solutionFile, code, 'utf8');

  // Build a single combined .cpp to avoid multi-TU duplicate symbol issues.
  // Tests #include "solution.cpp", so we include it once, then inline all test
  // code (with their #include stripped) and a main() driver.
  let combined = '#include "solution.cpp"\n\n';
  for (let i = 1; i <= part; i++) {
    const testFile = partDefs[i - 1].test_file;
    const src      = path.join(REPO_ROOT, problem.path, 'tests', 'cpp', testFile);
    if (fs.existsSync(src)) {
      let content = fs.readFileSync(src, 'utf8');
      content = content.replace(/^\s*#include\s+"solution\.cpp"\s*$/m, '// (included above)');
      combined += `// --- ${testFile} ---\n${content}\n\n`;
    }
  }

  // Append main() driver
  const partNames = partDefs.slice(0, part).map((_, i) => `part${i + 1}_tests`);
  combined += `
// --- generated main ---
int main() {
  int total_failures = 0;
  ${partNames.map(fn => `total_failures += ${fn}();`).join('\n  ')}
  return total_failures > 0 ? 1 : 0;
}
`;
  const combinedFile = path.join(tmpDir, 'combined.cpp');
  fs.writeFileSync(combinedFile, combined, 'utf8');

  const outBin   = path.join(tmpDir, os.platform() === 'win32' ? 'runner.exe' : 'runner');
  const compileCmd = `g++ -std=c++17 -DRUNNING_TESTS -o "${outBin}" "${combinedFile}" 2>&1`;

  const startTime = Date.now();

  exec(compileCmd, { timeout: 15000, cwd: tmpDir }, (compileErr, compileOut) => {
    if (compileErr) {
      cleanup(tmpDir);
      return res.json({
        success: false,
        submitted_part: part,
        compilation: { success: false, errors: compileOut || compileErr.message },
        parts: [],
        time_ms: Date.now() - startTime,
        runner_available: true,
      });
    }

    exec(`"${outBin}"`, { timeout: 10000, cwd: tmpDir }, (runErr, stdout, stderr) => {
      cleanup(tmpDir);

      const output = stdout || '';
      const parsedParts = parseTestOutput(output, partDefs.slice(0, part));

      const allPassed = parsedParts.every(p => p.all_passed);
      const timedOut  = runErr && runErr.killed;

      // Update progress based on results
      const freshProgress = loadProgress();
      const freshEntry    = freshProgress.problems[id] || {};
      ensurePartsInitialized(freshEntry, totalParts);

      if (allPassed && !timedOut) {
        freshEntry.parts[String(part)] = {
          status:        'passed',
          passed_at:     new Date().toISOString(),
          carry_forward: freshEntry.parts[String(part)]?.carry_forward ?? null,
        };
        // Unlock next part
        if (part < totalParts) {
          const nextKey = String(part + 1);
          if ((freshEntry.parts[nextKey] || {}).status === 'locked') {
            freshEntry.parts[nextKey] = { status: 'active', passed_at: null, carry_forward: null };
          }
        }
        // If last part, set completed_at
        if (part === totalParts) {
          freshEntry.completed_at = new Date().toISOString();
        }
        logActivity(freshProgress, 'passed_part');
      }

      freshProgress.problems[id] = freshEntry;
      saveProgress(freshProgress);

      res.json({
        success:         allPassed && !timedOut,
        submitted_part:  part,
        compilation:     { success: true, errors: null },
        parts:           parsedParts,
        unlocked_next_part: allPassed && !timedOut && part < totalParts,
        timed_out:       timedOut || false,
        time_ms:         Date.now() - startTime,
        runner_available: true,
      });
    });
  });
});

function cleanup(dir) {
  try { fs.rmSync(dir, { recursive: true, force: true }); } catch (_) {}
}

function parseTestOutput(stdout, partDefs) {
  const lines = stdout.split('\n').map(l => l.trim()).filter(Boolean);
  const results = [];

  for (let i = 0; i < partDefs.length; i++) {
    const partNum  = i + 1;
    const tests    = [];
    let passed     = 0;
    let total      = 0;
    let summaryParsed = false;

    for (const line of lines) {
      if (line.startsWith(`PART${partNum}_SUMMARY`)) {
        const match = line.match(/(\d+)\/(\d+)/);
        if (match) {
          passed = parseInt(match[1], 10);
          total  = parseInt(match[2], 10);
          summaryParsed = true;
        }
      } else if (line.startsWith('PASS ') || line.startsWith('FAIL ')) {
        const isPassed = line.startsWith('PASS ');
        const name     = line.slice(5).trim();
        tests.push({ name, passed: isPassed });
      }
    }

    if (!summaryParsed) {
      passed = tests.filter(t => t.passed).length;
      total  = tests.length;
    }

    results.push({
      part:       partNum,
      name:       partDefs[i].name,
      passed,
      total,
      all_passed: total > 0 && passed === total,
      tests,
    });
  }

  return results;
}

// ── POST /api/problems/:id/parts/:part/carry-forward ← NEW ───────────────────

app.post('/api/problems/:id/parts/:part/carry-forward', (req, res) => {
  const { id, part } = req.params;
  const { carry_forward } = req.body;
  const partNum = parseInt(part, 10);

  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const progress = loadProgress();
  const entry    = progress.problems[id] || {};
  const totalParts = (problem.parts || []).length || 1;
  ensurePartsInitialized(entry, totalParts);

  if (!entry.parts[String(partNum)]) {
    return res.status(400).json({ error: `Part ${partNum} not found` });
  }

  entry.parts[String(partNum)].carry_forward = !!carry_forward;
  progress.problems[id] = entry;
  saveProgress(progress);

  res.json({ part: partNum, carry_forward: !!carry_forward });
});

// ── POST /api/problems/:id/parts/:part/skip ← NEW ────────────────────────────

app.post('/api/problems/:id/parts/:part/skip', (req, res) => {
  if (testRunnerAvailable) {
    return res.status(403).json({ error: 'Skip is only available when g++ is not installed.' });
  }

  const { id, part } = req.params;
  const partNum = parseInt(part, 10);

  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const progress = loadProgress();
  const entry    = progress.problems[id] || {};
  const totalParts = (problem.parts || []).length || 1;
  ensurePartsInitialized(entry, totalParts);

  entry.parts[String(partNum)] = {
    status:       'passed',
    passed_at:    new Date().toISOString(),
    carry_forward: entry.parts[String(partNum)]?.carry_forward ?? null,
  };

  if (partNum < totalParts) {
    entry.parts[String(partNum + 1)] = { status: 'active', passed_at: null, carry_forward: null };
  } else {
    entry.completed_at = new Date().toISOString();
  }

  if (!entry.started_at) entry.started_at = new Date().toISOString();

  logActivity(progress, 'passed_part');
  progress.problems[id] = entry;
  saveProgress(progress);

  res.json({
    skipped:            true,
    part:               partNum,
    next_part_unlocked: partNum < totalParts ? partNum + 1 : null,
  });
});

// ── GET /api/problems/:id/readme ─────────────────────────────────────────────

app.get('/api/problems/:id/readme', (req, res) => {
  const { id } = req.params;
  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const filePath = path.join(REPO_ROOT, problem.path, 'README.md');
  if (!fs.existsSync(filePath)) return res.status(404).json({ error: `README.md not found for problem ${id}` });

  const raw = fs.readFileSync(filePath, 'utf8');
  const before_you_code = extractBeforeYouCode(raw);
  const html = marked(stripBeforeYouCode(raw));

  res.json({ html, before_you_code });
});

// ── GET /api/problems/:id/design ─────────────────────────────────────────────

app.get('/api/problems/:id/design', (req, res) => {
  const { id } = req.params;
  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const filePath = path.join(REPO_ROOT, problem.path, 'DESIGN.md');
  const html = renderMarkdown(filePath);
  if (!html) return res.status(404).json({ error: `DESIGN.md not found for problem ${id}` });

  res.json({ html });
});

// ── GET /api/problems/:id/ai-prompt ──────────────────────────────────────────

app.get('/api/problems/:id/ai-prompt', (req, res) => {
  const { id } = req.params;
  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const filePath = path.join(REPO_ROOT, problem.path, 'AI_REVIEW_PROMPT.md');
  if (!fs.existsSync(filePath)) {
    return res.status(404).json({ error: `AI_REVIEW_PROMPT.md not found for problem ${id}` });
  }
  const markdown = fs.readFileSync(filePath, 'utf8');
  res.json({ markdown });
});

// ── GET /api/primers ──────────────────────────────────────────────────────────

app.get('/api/primers', (req, res) => {
  const primers  = loadPrimers();
  const progress = loadProgress();
  const read     = new Set(progress.primers_read || []);

  const result = primers.map(p => ({ name: p.name, read: read.has(p.name) }));
  const summary = { total: result.length, read: result.filter(p => p.read).length };
  res.json({ primers: result, summary });
});

// ── GET /api/primers/:name ────────────────────────────────────────────────────

app.get('/api/primers/:name', (req, res) => {
  const { name } = req.params;
  const primers  = loadPrimers();
  const primer   = primers.find(p => p.name === name);
  if (!primer) return res.status(404).json({ error: `Primer '${name}' not found` });

  const html     = renderMarkdown(primer.file);
  const progress = loadProgress();
  const read     = (progress.primers_read || []).includes(name);

  res.json({ name, html, read });
});

// ── POST /api/primers/:name/read ─────────────────────────────────────────────

app.post('/api/primers/:name/read', (req, res) => {
  const { name } = req.params;
  const primers  = loadPrimers();
  if (!primers.find(p => p.name === name)) {
    return res.status(404).json({ error: `Primer '${name}' not found` });
  }

  const progress = loadProgress();
  if (!progress.primers_read) progress.primers_read = [];
  if (!progress.primers_read.includes(name)) {
    progress.primers_read.push(name);
    logActivity(progress, 'read_primer');
  }
  saveProgress(progress);
  res.json({ name, read: true });
});

// ── GET /api/activity ─────────────────────────────────────────────────────────

app.get('/api/activity', (req, res) => {
  const progress = loadProgress();
  const activity = progress.activity || [];
  const streak = calculateActivityStreak(progress);
  res.json({ activity, streak });
});

// ── GET /api/stats ────────────────────────────────────────────────────────────

app.get('/api/stats', (req, res) => {
  const problems    = loadProblems();
  const progress    = loadProgress();
  const primersRead = progress.primers_read || [];
  const merged      = problems.map(p => mergeProblemWithProgress(p, progress, primersRead));

  const overall = buildSummary(merged);
  overall.percent_complete = overall.total > 0
    ? Math.round((overall.solved / overall.total) * 100)
    : 0;
  overall.attempted = merged.filter(p => p.status === 'attempted').length;
  overall.unsolved  = overall.total - overall.solved - overall.attempted;

  // by_tier
  const by_tier = { 1: null, 2: null, 3: null };
  for (const tier of [1, 2, 3]) {
    const tier_problems = merged.filter(p => p.tier === tier);
    by_tier[tier] = buildSummary(tier_problems);
  }

  // by_pattern
  const by_pattern = {};
  for (const p of merged) {
    for (const pattern of (p.patterns || [])) {
      if (!by_pattern[pattern]) by_pattern[pattern] = { solved: 0, attempted: 0, total: 0 };
      by_pattern[pattern].total++;
      if (p.status === 'solved')    by_pattern[pattern].solved++;
      if (p.status === 'attempted') by_pattern[pattern].attempted++;
    }
  }

  // by_difficulty_mode
  const by_difficulty_mode = {
    interview: { solved: 0, attempted: 0 },
    guided:    { solved: 0, attempted: 0 },
    learning:  { solved: 0, attempted: 0 },
  };
  for (const p of merged) {
    const m = p.difficulty_mode || 'interview';
    if (by_difficulty_mode[m]) {
      if (p.status === 'solved')    by_difficulty_mode[m].solved++;
      if (p.status === 'attempted') by_difficulty_mode[m].attempted++;
    }
  }

  // parts_stats
  let totalPartsAcrossAll = 0;
  let partsPassed         = 0;
  let fullyCompleted      = 0;
  let partiallyCompleted  = 0;

  for (const p of merged) {
    totalPartsAcrossAll += p.total_parts;
    partsPassed         += p.parts_passed;
    if (p.status === 'solved')    fullyCompleted++;
    if (p.status === 'attempted') partiallyCompleted++;
  }

  const parts_stats = {
    total_parts_across_all_problems: totalPartsAcrossAll,
    parts_passed:                    partsPassed,
    avg_parts_per_solve:             fullyCompleted > 0 ? Math.round((partsPassed / fullyCompleted) * 10) / 10 : 0,
    problems_fully_completed:        fullyCompleted,
    problems_partially_completed:    partiallyCompleted,
  };

  // primers
  const primers  = loadPrimers();
  const readSet  = new Set(progress.primers_read || []);
  const primersStats = { total: primers.length, read: [...readSet].length };

  // streak
  const streak = calculateStreak(progress);

  res.json({ overall, by_tier, by_pattern, by_difficulty_mode, parts_stats, primers: primersStats, streak });
});

// ── GET /api/roadmap/progress ─────────────────────────────────────────────────
// Returns problem statuses, primers read, and manual roadmap checks for the roadmap page

app.get('/api/roadmap/progress', (req, res) => {
  const problems    = loadProblems();
  const progress    = loadProgress();
  const primersRead = progress.primers_read || [];
  const merged      = problems.map(p => mergeProblemWithProgress(p, progress, primersRead));

  const problemStatuses = {};
  for (const p of merged) {
    problemStatuses[p.id] = {
      status:       p.status,
      parts_passed: p.parts_passed,
      total_parts:  p.total_parts,
      current_part: p.current_part,
    };
  }

  res.json({
    problems:       problemStatuses,
    primers_read:   primersRead,
    manual_checks:  progress.roadmap_checks || {},
  });
});

// ── POST /api/roadmap/check ──────────────────────────────────────────────────
// Persist manual checkbox toggles for non-problem roadmap tasks

app.post('/api/roadmap/check', (req, res) => {
  const { key, checked } = req.body;
  if (typeof key !== 'string') return res.status(400).json({ error: 'key must be a string' });

  const progress = loadProgress();
  if (!progress.roadmap_checks) progress.roadmap_checks = {};
  progress.roadmap_checks[key] = !!checked;
  saveProgress(progress);

  res.json({ key, checked: !!checked });
});

// ── GET /api/runner-status ────────────────────────────────────────────────────

app.get('/api/runner-status', (req, res) => {
  res.json({ available: testRunnerAvailable });
});

// ── Static frontend ───────────────────────────────────────────────────────────

if (fs.existsSync(DIST_DIR)) {
  app.use(express.static(DIST_DIR));
  app.get('*', (req, res) => {
    res.sendFile(path.join(DIST_DIR, 'index.html'));
  });
} else {
  app.get('/', (req, res) => {
    res.send('<h2>Frontend not built. Run <code>npm install</code> first.</h2>');
  });
}

// ── Start server ──────────────────────────────────────────────────────────────

const PORT = parseInt(process.env.PORT || '3000', 10);

function startServer(port) {
  const server = app.listen(port, () => {
    console.log(`Dashboard running at http://localhost:${port}`);
    import('open').then(({ default: open }) => {
      open(`http://localhost:${port}`);
    }).catch(() => {});
  });

  server.on('error', (err) => {
    if (err.code === 'EADDRINUSE') {
      if (port === PORT) {
        console.log(`Port ${port} is in use. Trying port ${port + 1}...`);
        startServer(port + 1);
      } else {
        console.error(`Port ${port} is also in use. Try: PORT=3002 npm start`);
        process.exit(1);
      }
    } else {
      throw err;
    }
  });
}

startServer(PORT);
