const express = require('express');
const path = require('path');
const fs = require('fs');
const { exec } = require('child_process');
const os = require('os');
const yaml = require('js-yaml');
const PYTHON_CMD = os.platform() === 'win32' ? 'python' : 'python3';
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
// PROGRESS_JSON_PATH lets the e2e suite point at a throwaway file so test runs
// don't clobber the user's real progress. Resolved relative to REPO_ROOT when
// the env var is a relative path; absolute paths are honored as-is.
const PROGRESS_JSON  = process.env.PROGRESS_JSON_PATH
  ? (path.isAbsolute(process.env.PROGRESS_JSON_PATH)
      ? process.env.PROGRESS_JSON_PATH
      : path.join(REPO_ROOT, process.env.PROGRESS_JSON_PATH))
  : path.join(REPO_ROOT, 'progress.json');
const PATTERNS_DIR   = path.join(REPO_ROOT, 'patterns');
const DIST_DIR       = path.join(__dirname, 'dist');

// ─── g++ availability check ─────────────────────────────────────────────────

// Auto-detect MinGW on Windows — check common install locations
if (os.platform() === 'win32') {
  const mingwCandidates = [
    'C:\\ProgramData\\mingw64\\mingw64\\bin',
    'C:\\msys64\\mingw64\\bin',
    'C:\\msys64\\ucrt64\\bin',
    'C:\\mingw64\\bin',
    'C:\\MinGW\\bin',
    'C:\\TDM-GCC-64\\bin',
    path.join(process.env.LOCALAPPDATA || '', 'Programs', 'mingw64', 'bin'),
  ];
  for (const candidate of mingwCandidates) {
    if (fs.existsSync(path.join(candidate, 'g++.exe'))) {
      process.env.PATH = candidate + path.delimiter + (process.env.PATH || '');
      break;
    }
  }
}

// ─── Java (javac) availability check ────────────────────────────────────────

// Auto-detect JDK on all platforms — javac is often installed but not on PATH.
// Check JAVA_HOME first, then platform-specific common locations.
(function detectJava() {
  const { execSync } = require('child_process');
  // Quick check: already on PATH?
  try { execSync('javac -version', { stdio: 'ignore' }); return; } catch (_) {}

  const candidates = [];

  // 1. JAVA_HOME (all platforms)
  if (process.env.JAVA_HOME) {
    candidates.push(path.join(process.env.JAVA_HOME, 'bin'));
  }

  if (os.platform() === 'win32') {
    // Common Windows JDK install paths
    const programFiles = process.env.ProgramFiles || 'C:\\Program Files';
    const programFilesX86 = process.env['ProgramFiles(x86)'] || 'C:\\Program Files (x86)';
    for (const base of [programFiles, programFilesX86]) {
      for (const vendor of ['Java', 'Eclipse Adoptium', 'Microsoft', 'Amazon Corretto', 'Zulu', 'BellSoft']) {
        const vendorDir = path.join(base, vendor);
        if (fs.existsSync(vendorDir)) {
          try {
            const entries = fs.readdirSync(vendorDir).filter(e => e.startsWith('jdk'));
            entries.sort().reverse(); // prefer latest version
            for (const e of entries) candidates.push(path.join(vendorDir, e, 'bin'));
          } catch (_) {}
        }
      }
    }
  } else if (os.platform() === 'darwin') {
    // macOS: /usr/libexec/java_home, Homebrew, installed JDKs
    try {
      const jHome = execSync('/usr/libexec/java_home 2>/dev/null').toString().trim();
      if (jHome) candidates.push(path.join(jHome, 'bin'));
    } catch (_) {}
    candidates.push(
      '/opt/homebrew/opt/openjdk/bin',
      '/usr/local/opt/openjdk/bin'
    );
    // Scan /Library/Java/JavaVirtualMachines for installed JDKs
    const libJvm = '/Library/Java/JavaVirtualMachines';
    if (fs.existsSync(libJvm)) {
      try {
        const entries = fs.readdirSync(libJvm).filter(e => e.includes('jdk'));
        entries.sort().reverse();
        for (const e of entries) candidates.push(path.join(libJvm, e, 'Contents', 'Home', 'bin'));
      } catch (_) {}
    }
  } else {
    // Linux: common locations (Dell machines, Ubuntu, Fedora, etc.)
    candidates.push('/usr/lib/jvm/default/bin', '/usr/lib/jvm/java/bin');
    const jvmBase = '/usr/lib/jvm';
    if (fs.existsSync(jvmBase)) {
      try {
        const entries = fs.readdirSync(jvmBase).filter(e => e.includes('jdk') || e.includes('java'));
        entries.sort().reverse();
        for (const e of entries) candidates.push(path.join(jvmBase, e, 'bin'));
      } catch (_) {}
    }
    // SDKMAN
    const home = process.env.HOME || '';
    if (home) {
      candidates.push(path.join(home, '.sdkman', 'candidates', 'java', 'current', 'bin'));
      const jabbaDir = path.join(home, '.jabba', 'jdk');
      if (fs.existsSync(jabbaDir)) {
        try {
          const entries = fs.readdirSync(jabbaDir);
          entries.sort().reverse();
          for (const e of entries) candidates.push(path.join(jabbaDir, e, 'bin'));
        } catch (_) {}
      }
    }
    // Snap-installed OpenJDK
    candidates.push('/snap/openjdk/current/jdk/bin');
  }

  const javacName = os.platform() === 'win32' ? 'javac.exe' : 'javac';
  for (const dir of candidates) {
    if (fs.existsSync(path.join(dir, javacName))) {
      process.env.PATH = dir + path.delimiter + (process.env.PATH || '');
      console.log(`✅ Java auto-detected at ${dir}`);
      return;
    }
  }
})();

// ─── Python availability check ──────────────────────────────────────────────

// On Windows, python3 may not exist but 'python' or 'py' launcher might.
// On Linux/macOS with pyenv or conda, python3 might not be on the default PATH.
(function detectPython() {
  const { execSync } = require('child_process');
  // Quick check: already works?
  try { execSync(`${os.platform() === 'win32' ? 'python' : 'python3'} --version`, { stdio: 'ignore' }); return; } catch (_) {}

  const candidates = [];

  if (os.platform() === 'win32') {
    // Windows: check py launcher, then common install paths
    try { execSync('py -3 --version', { stdio: 'ignore' }); return; } catch (_) {}
    const localAppData = process.env.LOCALAPPDATA || '';
    if (localAppData) {
      const pyDir = path.join(localAppData, 'Programs', 'Python');
      if (fs.existsSync(pyDir)) {
        try {
          const entries = fs.readdirSync(pyDir).filter(e => e.startsWith('Python3'));
          entries.sort().reverse();
          for (const e of entries) candidates.push(path.join(pyDir, e));
        } catch (_) {}
      }
    }
    candidates.push('C:\\Python39', 'C:\\Python310', 'C:\\Python311', 'C:\\Python312', 'C:\\Python313');
  } else {
    // macOS / Linux: pyenv, conda, Homebrew
    const home = process.env.HOME || '';
    if (home) {
      candidates.push(
        path.join(home, '.pyenv', 'shims'),
        path.join(home, 'miniconda3', 'bin'),
        path.join(home, 'anaconda3', 'bin')
      );
    }
    candidates.push('/opt/homebrew/bin', '/usr/local/bin');
  }

  const pyName = os.platform() === 'win32' ? 'python.exe' : 'python3';
  for (const dir of candidates) {
    if (fs.existsSync(path.join(dir, pyName))) {
      process.env.PATH = dir + path.delimiter + (process.env.PATH || '');
      console.log(`✅ Python auto-detected at ${dir}`);
      return;
    }
  }
})();

// Language registry — adding a new language means: add an entry here, write a
// runner under harness/<lang>/, and ensure gen_stubs.py emits boilerplate for it.
const HARNESS_DIR = path.join(REPO_ROOT, 'harness');
const SNAKEYAML_JAR = path.join(HARNESS_DIR, 'java', 'lib', 'snakeyaml-2.2.jar');
const LANGS = {
  cpp: {
    ext: 'cpp',
    boilerplateMode: m => `${m}.cpp`,
    testFileLegacy:  n => `part${n}_test.cpp`,
    detect: 'g++ --version',
  },
  java: {
    ext: 'java',
    boilerplateMode: m => `${m.charAt(0).toUpperCase() + m.slice(1)}.java`,
    testFileLegacy:  n => `Part${n}Test.java`,
    detect: 'javac -version',
  },
  python: {
    ext: 'py',
    boilerplateMode: m => `${m}.py`,
    testFileLegacy:  n => `part${n}_test.py`,
    detect: `${PYTHON_CMD} --version`,
  },
  javascript: {
    ext: 'js',
    boilerplateMode: m => `${m}.js`,
    testFileLegacy:  n => `part${n}_test.js`,
    detect: 'node --version',
  },
  go: {
    ext: 'go',
    boilerplateMode: m => `${m}.go`,
    testFileLegacy:  n => `part${n}_runner.go`,
    detect: 'go version',
  },
};

const runnerAvailable = { cpp: false, java: false, python: false, javascript: false, go: false };

for (const [lang, cfg] of Object.entries(LANGS)) {
  exec(cfg.detect, (err) => {
    runnerAvailable[lang] = !err;
    if (err) {
      console.warn(`⚠️  ${lang}: '${cfg.detect}' failed. Submit for ${lang} disabled.`);
    } else {
      console.log(`✅ ${lang} runner available.`);
    }
  });
}


// ─── Data helpers ───────────────────────────────────────────────────────────

function loadProblems() {
  if (!fs.existsSync(PROBLEMS_YML)) {
    console.error('ERROR: Could not find docs/_data/problems.yml');
    return [];
  }
  try {
    const raw = fs.readFileSync(PROBLEMS_YML, 'utf8');
    const all = yaml.load(raw) || [];
    const TIER_DIR = { 1: 'tier1-foundation', 2: 'tier2-intermediate', 3: 'tier3-advanced' };
    // Filter out problems whose folder doesn't exist on disk
    return all.filter(p => {
      const tierFolder = TIER_DIR[p.tier] || `tier${p.tier}`;
      const derivedPath = p.path || `problems/${tierFolder}/${p.id}`;
      p.path = derivedPath;
      const problemDir = path.join(REPO_ROOT, derivedPath);
      if (!fs.existsSync(problemDir)) {
        return false;
      }
      return true;
    });
  } catch (e) {
    console.error('ERROR: Failed to parse problems.yml:', e.message);
    return [];
  }
}

// Normalize parts: handles both array-of-objects (free repo) and number (mined)
function normalizePartDefs(problem) {
  if (Array.isArray(problem.parts)) return problem.parts;
  const numParts = typeof problem.parts === 'number' ? problem.parts : 1;
  const details = problem.parts_detail || [];
  const defs = [];
  for (let i = 0; i < numParts; i++) {
    defs.push({
      name: details[i] || `Part ${i + 1}`,
      description_marker: `## Part ${i + 1}`,
      test_file: `part${i + 1}_test.cpp`,
      java_test_file: `Part${i + 1}Test.java`,
    });
  }
  return defs;
}

// Returns the code-storage key for a given lang+mode combo.
// C++ uses the legacy bare key (backward compat). Java uses 'java:<mode>'.
function codeKey(lang, mode) {
  // C++ keeps the legacy bare key; everything else namespaces by language.
  return lang === 'cpp' ? mode : `${lang}:${mode}`;
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
    version: 4,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    problems: {},
    primers_read: [],
    activity: [],
    // ── Gamification (§9) ──────────────────────────────────────────────────────
    // Signal Score (XP) is NOT stored — it's derived client-side from passed
    // parts × tier × dsa_depth. Only earned/spent values persist here.
    gamification: {
      hint_tokens: 0,
      unlocked_design: [],     // problem ids whose DESIGN.md was bought
      unlocked_editorial: [],  // problem ids whose editorial was bought early
      drill_stats: { total_correct: 0, best_streak: 0, runs: 0, tokens_earned: 0 },
    },
  };
}

// Ensure the gamification block exists (defensive — also covers the migration).
function ensureGamification(progress) {
  if (!progress.gamification) {
    progress.gamification = {
      hint_tokens: 0,
      unlocked_design: [],
      unlocked_editorial: [],
      drill_stats: { total_correct: 0, best_streak: 0, runs: 0, tokens_earned: 0 },
    };
  }
  const g = progress.gamification;
  if (typeof g.hint_tokens !== 'number') g.hint_tokens = 0;
  if (!Array.isArray(g.unlocked_design)) g.unlocked_design = [];
  if (!Array.isArray(g.unlocked_editorial)) g.unlocked_editorial = [];
  if (!g.drill_stats) g.drill_stats = { total_correct: 0, best_streak: 0, runs: 0, tokens_earned: 0 };
  return g;
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
  // v3 → v4: add the gamification block (idempotent).
  if ((data.version || 1) >= 3 && (data.version || 1) < 4) {
    console.log('Migrating progress.json from v3 to v4 (gamification)...');
    ensureGamification(data);
    data.version = 4;
    saveProgress(data);
    console.log('Migration complete.');
  }
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

  ensureGamification(data);
  data.version = 4;
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

  // Per part, the set of heading prefixes that may introduce that part. We
  // accept all common conventions used in our READMEs:
  //   • "## Part N"                (canonical)
  //   • "## Base Requirement"      (part 1)        + "## Extension N-1"
  //   • "### Base Requirement"     (part 1, nested under "## Parts")
  //                                                + "### Extension N-1"
  const candidatesFor = (idx, def) => {
    const m = [def.description_marker, `## Part ${idx + 1}`, `### Part ${idx + 1}`];
    if (idx === 0) m.push('## Base Requirement', '### Base Requirement');
    else m.push(`## Extension ${idx}`, `### Extension ${idx}`);
    return m;
  };

  const matches = (line, mks) => {
    const t = line.trim();
    return mks.some(m => t.startsWith(m));
  };

  const allCandidates = partDefs.map((def, i) => candidatesFor(i, def));

  const markerIndices = allCandidates.map(mks => {
    for (let i = 0; i < lines.length; i++) {
      if (matches(lines[i], mks)) return i;
    }
    return -1;
  });

  // For section termination: next part marker, OR the next H2 heading that
  // is NOT a part marker (e.g. "## Running Tests", "## Constraints"). This
  // prevents trailing sections from being absorbed into the last part.
  const isPartMarkerLine = (j) => allCandidates.some(mks => matches(lines[j], mks));
  const nextSectionBreak = (start) => {
    for (let j = start + 1; j < lines.length; j++) {
      if (/^##\s/.test(lines[j]) && !isPartMarkerLine(j)) return j;
    }
    return lines.length;
  };

  // Everything before the first found marker = scenario/context.
  // If that scenario ends with a "## Parts" wrapper heading, drop it — the
  // parts list itself will be rendered below.
  const firstMarker = markerIndices.find(i => i !== -1);
  let scenarioLines = firstMarker !== undefined && firstMarker > 0
    ? lines.slice(0, firstMarker)
    : firstMarker === undefined ? lines : [];
  while (scenarioLines.length && /^(\s*|##\s+Parts\s*)$/.test(scenarioLines[scenarioLines.length - 1])) {
    scenarioLines = scenarioLines.slice(0, -1);
  }
  const scenario = scenarioLines.join('\n').trim();

  const sections = [];
  for (let i = 0; i < markerIndices.length; i++) {
    const start = markerIndices[i];
    if (start === -1) {
      sections.push(null);
      continue;
    }
    const nextPart = markerIndices.slice(i + 1).find(idx => idx !== -1);
    const breakAt  = nextSectionBreak(start);
    const end = nextPart !== undefined ? Math.min(nextPart, breakAt) : breakAt;
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
  if (!entry.parts) entry.parts = {};

  // Fill any missing slots from 1..totalParts. Part 1 defaults to 'active' when
  // truly fresh; later parts default to 'locked'. We can't trust that an
  // existing partial dict (e.g. only "1" present from earlier writes) covers
  // all parts — without backfilling, the submit-time unlock check in
  // updateProgressAfterRun would never see part N+1 as 'locked' and would
  // silently skip the unlock.
  for (let i = 1; i <= totalParts; i++) {
    const key = String(i);
    if (!entry.parts[key]) {
      entry.parts[key] = {
        status: i === 1 ? 'active' : 'locked',
        passed_at: null,
        carry_forward: null,
      };
    }
  }
  return entry;
}

function mergeProblemWithProgress(problem, progress, primersRead) {
  const p = progress.problems[problem.id] || {};
  const totalParts = normalizePartDefs(problem).length;
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
    reveal_mode:      p.reveal_mode     || 'interview',
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
  const { notes, difficulty_mode, reveal_mode } = req.body;

  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const progress = loadProgress();
  const entry    = progress.problems[id] || {};

  if (difficulty_mode !== undefined && ['interview', 'guided', 'learning'].includes(difficulty_mode)) {
    entry.difficulty_mode = difficulty_mode;
  }
  if (reveal_mode !== undefined && ['interview', 'study'].includes(reveal_mode)) {
    entry.reveal_mode = reveal_mode;
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

// ── POST /api/problems/:id/reset ─────────────────────────────────────────────
// Wipes all progress for a problem: clears saved code (every lang/mode),
// timestamps, and rebuilds parts to the canonical fresh state (part 1 active,
// the rest locked). Difficulty/reveal mode and notes are preserved — this
// resets PROGRESS, not the user's per-problem preferences.

app.post('/api/problems/:id/reset', (req, res) => {
  const { id } = req.params;

  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const progress = loadProgress();
  const entry    = progress.problems[id];

  // Nothing recorded yet — already in a fresh state.
  if (!entry) {
    const primersRead = progress.primers_read || [];
    return res.json(mergeProblemWithProgress(problem, progress, primersRead));
  }

  // Drop progress fields; keep preferences (difficulty_mode, reveal_mode, notes).
  delete entry.parts;
  delete entry.code;
  delete entry.started_at;
  delete entry.completed_at;

  const totalParts = normalizePartDefs(problem).length;
  ensurePartsInitialized(entry, totalParts);  // rebuilds part 1 active / rest locked

  progress.problems[id] = entry;
  logActivity(progress, 'reset_problem');
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
  const partDefs   = normalizePartDefs(problem);
  const totalParts = partDefs.length;

  ensurePartsInitialized(entry, totalParts);
  const partsProgress = entry.parts || {};
  const revealMode = entry.reveal_mode || 'interview';
  const isStudyMode = revealMode === 'study';

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
    // In study mode, show descriptions for locked parts; in interview mode, hide them
    const hideContent = isLocked && !isStudyMode;
    const section    = sections[i];

    // Count tests — prefer spec-driven tests/cases/partN.yaml when available,
    // else fall back to legacy per-language test files.
    let testCount = 0;
    if (!isLocked) {
      const lang = LANGS[req.query.lang] ? req.query.lang : 'cpp';
      const specCases = path.join(REPO_ROOT, problem.path, 'tests', 'cases', `part${partNum}.yaml`);
      if (fs.existsSync(specCases)) {
        const raw = fs.readFileSync(specCases, 'utf8');
        const arr = yaml.load(raw) || [];
        testCount = Array.isArray(arr) ? arr.length : 0;
      } else {
        const testFileName = lang === 'java' ? def.java_test_file : def.test_file;
        const testPath = path.join(REPO_ROOT, problem.path, 'tests', lang, testFileName || '');
        if (fs.existsSync(testPath)) {
          const testContent = fs.readFileSync(testPath, 'utf8');
          const matches = lang === 'java'
            ? testContent.match(/System\.out\.println\("(PASS|FAIL)\s/g)
            : testContent.match(/cout\s*<<\s*"(PASS|FAIL)\s/g);
          testCount = matches ? matches.length / 2 : 0;
        }
      }
    }

    return {
      part:             partNum,
      name:             def.name,
      status:           partProg.status,
      passed_at:        partProg.passed_at,
      carry_forward:    partProg.carry_forward,
      description_html: hideContent ? null : (section ? marked(section) : null),
      test_count:       testCount,
    };
  });

  res.json({
    id,
    total_parts: totalParts,
    scenario_html: scenario ? marked(scenario) : null,
    parts,
    runner_available: runnerAvailable.cpp,
  });
});

// ── GET /api/problems/:id/starter ────────────────────────────────────────────
// Updated: accepts ?mode= and ?part= query params

app.get('/api/problems/:id/starter', (req, res) => {
  const { id } = req.params;
  const mode = ['interview', 'guided', 'learning'].includes(req.query.mode)
    ? req.query.mode : 'interview';
  const part = parseInt(req.query.part, 10) || 1;
  const lang = LANGS[req.query.lang] ? req.query.lang : 'cpp';

  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const modeFile = LANGS[lang].boilerplateMode;
  const tryFile = (p, m) => path.join(REPO_ROOT, problem.path, 'boilerplate', lang, `part${p}`, modeFile(m));

  let filePath = tryFile(part, mode);
  let fallback = false;

  if (!fs.existsSync(filePath)) { filePath = tryFile(part, 'learning'); fallback = true; }
  if (!fs.existsSync(filePath)) { filePath = tryFile(1, mode);          fallback = true; }
  if (!fs.existsSync(filePath)) { filePath = tryFile(1, 'learning');    fallback = true; }

  const commentChar = lang === 'python' ? '#' : '//';
  const defaultComment = `${commentChar} Starter file not found for mode: ${mode}, part: ${part}\n${commentChar} Write your solution here.\n`;

  const code = fs.existsSync(filePath)
    ? fs.readFileSync(filePath, 'utf8')
    : defaultComment;

  res.json({ lang, mode, part, code, fallback });
});

// ── GET /api/problems/:id/code ────────────────────────────────────────────────

app.get('/api/problems/:id/code', (req, res) => {
  const { id } = req.params;
  const mode = ['interview', 'guided', 'learning'].includes(req.query.mode)
    ? req.query.mode : 'interview';
  const lang = LANGS[req.query.lang] ? req.query.lang : 'cpp';

  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const progress = loadProgress();
  const entry    = progress.problems[id] || {};
  const saved    = (entry.code || {})[codeKey(lang, mode)] || '';

  if (saved) {
    return res.json({ lang, mode, code: saved, is_starter: false });
  }

  const modeFile = LANGS[lang].boilerplateMode;
  const tryFile = (m) => path.join(REPO_ROOT, problem.path, 'boilerplate', lang, 'part1', modeFile(m));
  let filePath = tryFile(mode);
  if (!fs.existsSync(filePath)) filePath = tryFile('learning');

  const commentChar = lang === 'python' ? '#' : '//';
  const code = fs.existsSync(filePath)
    ? fs.readFileSync(filePath, 'utf8')
    : `${commentChar} Write your solution here.\n`;

  res.json({ lang, mode, code, is_starter: true });
});

// ── POST /api/problems/:id/code ───────────────────────────────────────────────

app.post('/api/problems/:id/code', (req, res) => {
  const { id } = req.params;
  const { mode, code, lang = 'cpp' } = req.body;

  if (!['interview', 'guided', 'learning'].includes(mode)) {
    return res.status(400).json({ error: 'mode must be interview, guided, or learning' });
  }
  if (!LANGS[lang]) {
    return res.status(400).json({ error: `lang must be one of ${Object.keys(LANGS).join(', ')}` });
  }
  if (typeof code !== 'string') return res.status(400).json({ error: 'code must be a string' });

  const problems = loadProblems();
  if (!problems.find(p => p.id === id)) return res.status(404).json({ error: `Problem ${id} not found` });

  const progress = loadProgress();
  if (!progress.problems[id]) progress.problems[id] = {};
  if (!progress.problems[id].code) progress.problems[id].code = {};
  progress.problems[id].code[codeKey(lang, mode)] = code;
  logActivity(progress, 'saved_code');
  saveProgress(progress);

  res.json({ saved: true, lang, mode });
});

// ── POST /api/problems/:id/submit ← NEW ──────────────────────────────────────

app.post('/api/problems/:id/submit', (req, res) => {
  const { id } = req.params;
  const { part, mode, code, lang = 'cpp' } = req.body;

  if (typeof code !== 'string') return res.status(400).json({ error: 'code must be a string' });
  if (!Number.isInteger(part) || part < 1) return res.status(400).json({ error: 'part must be a positive integer' });
  if (!LANGS[lang]) return res.status(400).json({ error: `lang must be one of ${Object.keys(LANGS).join(', ')}` });

  if (!runnerAvailable[lang]) {
    return res.status(503).json({
      error: `${lang} runner not available on this host (detected via '${LANGS[lang].detect}').`,
      runner_available: false,
    });
  }
  const isJava = lang === 'java';

  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const partDefs   = normalizePartDefs(problem);
  const totalParts = partDefs.length;

  if (part > totalParts) {
    return res.status(400).json({ error: `Part ${part} does not exist. Problem has ${totalParts} parts.` });
  }

  // Save code first
  const progress = loadProgress();
  if (!progress.problems[id]) progress.problems[id] = {};
  const entry = progress.problems[id];
  if (!entry.code) entry.code = {};
  entry.code[codeKey(lang, mode || 'interview')] = code;
  ensurePartsInitialized(entry, totalParts);

  if (!entry.started_at) entry.started_at = new Date().toISOString();
  if ((entry.parts[String(part)] || {}).status !== 'passed') {
    entry.parts[String(part)] = { ...entry.parts[String(part)], status: 'attempted' };
  }
  saveProgress(progress);

  // Create temp directory
  const timestamp = Date.now();
  const tmpDir    = path.join(os.tmpdir(), `cj-${id}-${lang}-${timestamp}`);
  fs.mkdirSync(tmpDir, { recursive: true });

  const startTime = Date.now();

  // ── Shared: update progress after run ───────────────────────────────────────
  function updateProgressAfterRun(allPassed, timedOut) {
    const freshProgress = loadProgress();
    const freshEntry    = freshProgress.problems[id] || {};
    ensurePartsInitialized(freshEntry, totalParts);
    if (allPassed && !timedOut) {
      freshEntry.parts[String(part)] = {
        status:        'passed',
        passed_at:     new Date().toISOString(),
        carry_forward: freshEntry.parts[String(part)]?.carry_forward ?? null,
      };
      if (part < totalParts) {
        const nextKey  = String(part + 1);
        const nextProg = freshEntry.parts[nextKey];
        if (!nextProg || nextProg.status === 'locked') {
          freshEntry.parts[nextKey] = { status: 'active', passed_at: null, carry_forward: null };
        }
      }
      if (part === totalParts) freshEntry.completed_at = new Date().toISOString();
      logActivity(freshProgress, 'passed_part');
    }
    freshProgress.problems[id] = freshEntry;
    saveProgress(freshProgress);
  }

  // Common reply shape for compile / runtime failures.
  const reply = (compilationOk, errors, parsedParts, timedOut) => {
    const allPassed = compilationOk && parsedParts.every(p => p.all_passed) && !timedOut;
    updateProgressAfterRun(allPassed, timedOut);
    res.json({
      success: allPassed,
      submitted_part: part,
      compilation: { success: compilationOk, errors },
      parts: parsedParts,
      unlocked_next_part: allPassed && part < totalParts,
      timed_out: timedOut || false,
      time_ms: Date.now() - startTime,
      runner_available: true,
    });
  };

  // Spec-driven path — used when problem has spec.yaml + tests/cases/.
  const problemAbsDir = path.join(REPO_ROOT, problem.path);
  const hasSpec = fs.existsSync(path.join(problemAbsDir, 'spec.yaml'))
              && fs.existsSync(path.join(problemAbsDir, 'tests', 'cases'));

  if (hasSpec) {
    runSpecDrivenSubmit({ lang, problem, problemAbsDir, partDefs, part, code, tmpDir, reply });
    return;
  }

  if (lang === 'python' || lang === 'javascript') {
    return res.status(400).json({
      error: `Problem ${id} has no spec.yaml. ${lang} is only supported for spec-driven problems.`,
    });
  }

  if (!isJava) {
    // ── C++ pipeline ─────────────────────────────────────────────────────────
    const solutionFile = path.join(tmpDir, 'solution.cpp');
    fs.writeFileSync(solutionFile, code, 'utf8');

    let combined = '#include "solution.cpp"\n\n';
    for (let i = 1; i <= part; i++) {
      const testFile = partDefs[i - 1].test_file;
      const src      = path.join(REPO_ROOT, problem.path, 'tests', 'cpp', testFile);
      if (fs.existsSync(src)) {
        let content = fs.readFileSync(src, 'utf8');
        content = content.replace(/^\s*#include\s+"(?:\.\.\/)*solution\.cpp"\s*$/m, '// (included above)');
        content = content.replace(/^\s*int\s+main\s*\(\s*\)\s*\{[\s\S]*?\bpart\d+_tests\s*\([^)]*\)[\s\S]*?\}\s*$/m, '// (standalone main stripped by harness)');
        combined += `// --- ${testFile} ---\n${content}\n\n`;
      }
    }
    const partNames = partDefs.slice(0, part).map((_, i) => `part${i + 1}_tests`);
    combined += `\n// --- generated main ---\nint main() {\n  int total_failures = 0;\n  ${partNames.map(fn => `total_failures += ${fn}();`).join('\n  ')}\n  return total_failures > 0 ? 1 : 0;\n}\n`;

    const combinedFile = path.join(tmpDir, 'combined.cpp');
    fs.writeFileSync(combinedFile, combined, 'utf8');

    const outBin     = path.join(tmpDir, os.platform() === 'win32' ? 'runner.exe' : 'runner');
    const compileCmd = `g++ -std=c++17 -DRUNNING_TESTS -o "${outBin}" "${combinedFile}" 2>&1`;

    exec(compileCmd, { timeout: 15000, cwd: tmpDir }, (compileErr, compileOut) => {
      if (compileErr) {
        cleanup(tmpDir);
        return res.json({
          success: false, submitted_part: part,
          compilation: { success: false, errors: compileOut || compileErr.message },
          parts: [], time_ms: Date.now() - startTime, runner_available: true,
        });
      }
      exec(`"${outBin}"`, { timeout: 10000, cwd: tmpDir }, (runErr, stdout) => {
        cleanup(tmpDir);
        const parsedParts = parseTestOutput(stdout || '', partDefs.slice(0, part));
        const allPassed   = parsedParts.every(p => p.all_passed);
        const timedOut    = runErr && runErr.killed;
        updateProgressAfterRun(allPassed, timedOut);
        res.json({
          success: allPassed && !timedOut, submitted_part: part,
          compilation: { success: true, errors: null },
          parts: parsedParts,
          unlocked_next_part: allPassed && !timedOut && part < totalParts,
          timed_out: timedOut || false,
          time_ms: Date.now() - startTime, runner_available: true,
        });
      });
    });

  } else {
    // ── Java pipeline ─────────────────────────────────────────────────────────
    // Write user's solution
    fs.writeFileSync(path.join(tmpDir, 'Solution.java'), code, 'utf8');

    // Copy test files
    for (let i = 1; i <= part; i++) {
      const testFileName = partDefs[i - 1].java_test_file;
      const src = path.join(REPO_ROOT, problem.path, 'tests', 'java', testFileName);
      if (fs.existsSync(src)) {
        fs.copyFileSync(src, path.join(tmpDir, testFileName));
      }
    }

    // Generate Main.java driver
    const testClasses = partDefs.slice(0, part).map((_, i) => `Part${i + 1}Test`);
    const mainJava = [
      'public class Main {',
      '    public static void main(String[] args) {',
      '        int failures = 0;',
      ...testClasses.map(cls => `        failures += ${cls}.runTests();`),
      '        System.exit(failures > 0 ? 1 : 0);',
      '    }',
      '}',
    ].join('\n');
    fs.writeFileSync(path.join(tmpDir, 'Main.java'), mainJava, 'utf8');

    // List all .java files explicitly (no shell glob needed)
    const javaFiles = fs.readdirSync(tmpDir)
      .filter(f => f.endsWith('.java'))
      .map(f => `"${path.join(tmpDir, f)}"`)
      .join(' ');
    const compileCmd = `javac -cp "${tmpDir}" ${javaFiles} 2>&1`;

    exec(compileCmd, { timeout: 20000, cwd: tmpDir }, (compileErr, compileOut) => {
      if (compileErr) {
        cleanup(tmpDir);
        return res.json({
          success: false, submitted_part: part,
          compilation: { success: false, errors: compileOut || compileErr.message },
          parts: [], time_ms: Date.now() - startTime, runner_available: true,
        });
      }
      exec(`java -cp "${tmpDir}" Main`, { timeout: 10000, cwd: tmpDir }, (runErr, stdout) => {
        cleanup(tmpDir);
        const parsedParts = parseTestOutput(stdout || '', partDefs.slice(0, part));
        const allPassed   = parsedParts.every(p => p.all_passed);
        const timedOut    = runErr && runErr.killed;
        updateProgressAfterRun(allPassed, timedOut);
        res.json({
          success: allPassed && !timedOut, submitted_part: part,
          compilation: { success: true, errors: null },
          parts: parsedParts,
          unlocked_next_part: allPassed && !timedOut && part < totalParts,
          timed_out: timedOut || false,
          time_ms: Date.now() - startTime, runner_available: true,
        });
      });
    });
  }
});

function cleanup(dir) {
  try { fs.rmSync(dir, { recursive: true, force: true }); } catch (_) {}
}

// Spec-driven submission — used when problem has spec.yaml + tests/cases/.
// Each language's runner reads the YAML at runtime (Python, Java) or via codegen (C++)
// and emits PASS/FAIL/PART<N>_SUMMARY lines that parseTestOutput consumes unchanged.
function runSpecDrivenSubmit({ lang, problem, problemAbsDir, partDefs, part, code, tmpDir, reply }) {
  // Stage spec + cases for the runner to find.
  fs.copyFileSync(path.join(problemAbsDir, 'spec.yaml'), path.join(tmpDir, 'spec.yaml'));
  const casesDst = path.join(tmpDir, 'tests', 'cases');
  fs.mkdirSync(casesDst, { recursive: true });
  for (const f of fs.readdirSync(path.join(problemAbsDir, 'tests', 'cases'))) {
    if (f.endsWith('.yaml')) fs.copyFileSync(path.join(problemAbsDir, 'tests', 'cases', f), path.join(casesDst, f));
  }

  const t0 = Date.now();
  const driveParts = (compileFn, runFn) => {
    compileFn((compileErr, compileOut) => {
      if (compileErr) {
        cleanup(tmpDir);
        return reply(false, compileOut || compileErr.message, [], false);
      }
      // Run each part 1..N sequentially; collect output and parse.
      let allOut = '';
      let timedOut = false;
      const next = (i) => {
        if (i > part) {
          cleanup(tmpDir);
          return reply(true, null, parseTestOutput(allOut, partDefs.slice(0, part)), timedOut);
        }
        runFn(i, (runErr, stdout) => {
          allOut += '\n' + (stdout || '');
          if (runErr && runErr.killed) timedOut = true;
          next(i + 1);
        });
      };
      next(1);
    });
  };

  if (lang === 'python') {
    fs.writeFileSync(path.join(tmpDir, 'solution.py'), code, 'utf8');
    fs.copyFileSync(path.join(HARNESS_DIR, 'python', 'runner.py'), path.join(tmpDir, 'runner.py'));
    driveParts(
      (cb) => cb(null, ''),  // no compile step
      (i, cb) => exec(`${PYTHON_CMD} "${path.join(tmpDir, 'runner.py')}" --problem-dir "${tmpDir}" --part ${i}`, { timeout: 10000, cwd: tmpDir }, cb),
    );
    return;
  }

  if (lang === 'javascript') {
    fs.writeFileSync(path.join(tmpDir, 'solution.js'), code, 'utf8');
    fs.copyFileSync(path.join(HARNESS_DIR, 'javascript', 'runner.js'), path.join(tmpDir, 'runner.js'));
    // The runner depends on js-yaml; point NODE_PATH at the dashboard's
    // node_modules so the temp-dir copy resolves it (no per-problem install).
    const runEnv = { ...process.env, NODE_PATH: path.join(__dirname, 'node_modules') };
    driveParts(
      (cb) => cb(null, ''),  // interpreted — no compile step
      (i, cb) => exec(`node "${path.join(tmpDir, 'runner.js')}" --problem-dir "${tmpDir}" --part ${i}`, { timeout: 10000, cwd: tmpDir, env: runEnv }, cb),
    );
    return;
  }

  if (lang === 'java') {
    fs.writeFileSync(path.join(tmpDir, 'Solution.java'), code, 'utf8');
    fs.copyFileSync(path.join(HARNESS_DIR, 'java', 'Runner.java'), path.join(tmpDir, 'Runner.java'));
    // path.delimiter is ';' on Windows, ':' on macOS/Linux — never hardcode.
    const cp = `${tmpDir}${path.delimiter}${SNAKEYAML_JAR}`;
    driveParts(
      (cb) => exec(`javac -cp "${cp}" "${path.join(tmpDir, 'Solution.java')}" "${path.join(tmpDir, 'Runner.java')}" 2>&1`, { timeout: 20000, cwd: tmpDir }, cb),
      (i, cb) => exec(`java -cp "${cp}" Runner --problem-dir "${tmpDir}" --part ${i}`, { timeout: 10000, cwd: tmpDir }, cb),
    );
    return;
  }

  if (lang === 'cpp') {
    fs.writeFileSync(path.join(tmpDir, 'solution.cpp'), code, 'utf8');
    // Codegen one runner per part, compile each, run.
    const codegen = path.join(HARNESS_DIR, 'cpp', 'codegen.py');
    let compileErrs = '';
    const compileAll = (cb) => {
      const compileOne = (i) => {
        if (i > part) return cb(null, '');
        const runnerSrc = path.join(tmpDir, `runner${i}.cpp`);
        exec(`${PYTHON_CMD} "${codegen}" "${problemAbsDir}" ${i} > "${runnerSrc}"`, { timeout: 5000 }, (cgErr) => {
          if (cgErr) return cb(cgErr, `codegen failed for part ${i}: ${cgErr.message}`);
          const outBin = path.join(tmpDir, os.platform() === 'win32' ? `r${i}.exe` : `r${i}`);
          exec(`g++ -std=c++17 -DRUNNING_TESTS -o "${outBin}" "${runnerSrc}" 2>&1`, { timeout: 15000, cwd: tmpDir }, (err, out) => {
            if (err) { compileErrs += out || err.message; return cb(err, compileErrs); }
            compileOne(i + 1);
          });
        });
      };
      compileOne(1);
    };
    driveParts(
      compileAll,
      (i, cb) => {
        const outBin = path.join(tmpDir, os.platform() === 'win32' ? `r${i}.exe` : `r${i}`);
        exec(`"${outBin}"`, { timeout: 10000, cwd: tmpDir }, cb);
      },
    );
    return;
  }

  if (lang === 'go') {
    // Go is compiled and has one main() per directory, so each part gets its
    // own subdir (goN/) holding solution.go + a generated runner.go + a minimal
    // go.mod, then `go run .` offline (stdlib-only, no module downloads).
    const codegen = path.join(HARNESS_DIR, 'go', 'codegen.py');
    const partDir = (i) => path.join(tmpDir, `go${i}`);
    let compileErrs = '';
    const stageAll = (cb) => {
      const stageOne = (i) => {
        if (i > part) return cb(null, '');
        const d = partDir(i);
        try {
          fs.mkdirSync(d, { recursive: true });
          fs.writeFileSync(path.join(d, 'solution.go'), code, 'utf8');
          fs.writeFileSync(path.join(d, 'go.mod'), 'module cjrun\n\ngo 1.21\n', 'utf8');
        } catch (e) { return cb(e, `staging failed for part ${i}: ${e.message}`); }
        const runnerSrc = path.join(d, 'runner.go');
        exec(`${PYTHON_CMD} "${codegen}" "${problemAbsDir}" ${i} > "${runnerSrc}"`, { timeout: 5000 }, (cgErr) => {
          if (cgErr) return cb(cgErr, `codegen failed for part ${i}: ${cgErr.message}`);
          stageOne(i + 1);
        });
      };
      stageOne(1);
    };
    const goEnv = { ...process.env, GOPROXY: 'off', GOFLAGS: '-mod=mod', CGO_ENABLED: '0' };
    driveParts(
      // `go run` compiles + runs in one step; compile errors surface at run time
      // per-part, so there's no separate compile pass. We still validate that
      // the build for part `part` is sane by letting run-time report it.
      (cb) => stageAll(cb),
      (i, cb) => exec('go run .', { timeout: 30000, cwd: partDir(i), env: goEnv }, (err, stdout, stderr) => {
        // A build failure yields no PART summary; surface stderr so the user
        // sees the compile error (parseTestOutput treats missing parts as 0/0).
        if (err && !(stdout || '').includes('PART')) {
          compileErrs += stderr || err.message;
          return cb(null, `BUILD ERROR (part ${i}):\n${stderr || err.message}`);
        }
        cb(null, stdout);
      }),
    );
    return;
  }
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
  if (runnerAvailable.cpp) {
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

// ── GET /api/problems/:id/editorial ──────────────────────────────────────────

app.get('/api/problems/:id/editorial', (req, res) => {
  const { id } = req.params;
  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const filePath = path.join(REPO_ROOT, problem.path, 'EDITORIAL.md');
  if (!fs.existsSync(filePath)) return res.json({ available: false, html: null });

  const html = renderMarkdown(filePath);
  res.json({ available: true, html });
});

// ── GET /api/problems/:id/solution ────────────────────────────────────────────

app.get('/api/problems/:id/solution', (req, res) => {
  const { id } = req.params;
  const { lang = 'cpp' } = req.query;
  const problems = loadProblems();
  const problem  = problems.find(p => p.id === id);
  if (!problem) return res.status(404).json({ error: `Problem ${id} not found` });

  const ext = LANGS[lang] ? LANGS[lang].ext : 'cpp';
  const filename = `solution.${ext}`;
  const filePath = path.join(REPO_ROOT, problem.path, filename);
  if (!fs.existsSync(filePath)) return res.json({ available: false, code: null, html: null });

  const code = fs.readFileSync(filePath, 'utf8');
  const langTag = lang === 'python' ? 'python'
    : lang === 'java' ? 'java'
    : lang === 'javascript' ? 'javascript'
    : 'cpp';
  const html = marked('```' + langTag + '\n' + code + '\n```');
  res.json({ available: true, code, html });
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

// ── Gamification (§9) ─────────────────────────────────────────────────────────

// GET /api/gamification — token balance, content unlocks, drill stats.
app.get('/api/gamification', (req, res) => {
  const progress = loadProgress();
  const g = ensureGamification(progress);
  res.json({
    hint_tokens:        g.hint_tokens,
    unlocked_design:    g.unlocked_design,
    unlocked_editorial: g.unlocked_editorial,
    drill_stats:        g.drill_stats,
  });
});

// POST /api/gamification/spend — spend 1 token to unlock design/editorial content.
// Server-authoritative: rejects if balance < 1. Idempotent if already unlocked.
app.post('/api/gamification/spend', (req, res) => {
  const { kind, problemId } = req.body || {};
  if (!['design', 'editorial'].includes(kind)) {
    return res.status(400).json({ error: "kind must be 'design' or 'editorial'" });
  }
  if (typeof problemId !== 'string' || !problemId) {
    return res.status(400).json({ error: 'problemId is required' });
  }
  const progress = loadProgress();
  const g = ensureGamification(progress);
  const list = kind === 'design' ? g.unlocked_design : g.unlocked_editorial;

  // Already unlocked → no charge, just confirm.
  if (list.includes(problemId)) {
    return res.json({ ok: true, already: true, hint_tokens: g.hint_tokens, kind, problemId });
  }
  if (g.hint_tokens < 1) {
    return res.status(402).json({ error: 'Not enough hint tokens', hint_tokens: g.hint_tokens });
  }
  g.hint_tokens -= 1;
  list.push(problemId);
  saveProgress(progress);
  res.json({ ok: true, hint_tokens: g.hint_tokens, kind, problemId });
});

// POST /api/gamification/drill-result — record a drill run; grant a token on a
// 3-correct run. Body: { correct: bool, runStreak: int, runComplete: bool, earnedToken: bool }
app.post('/api/gamification/drill-result', (req, res) => {
  const { correct, runStreak, runComplete, earnedToken } = req.body || {};
  const progress = loadProgress();
  const g = ensureGamification(progress);

  if (correct) g.drill_stats.total_correct = (g.drill_stats.total_correct || 0) + 1;
  if (typeof runStreak === 'number') {
    g.drill_stats.best_streak = Math.max(g.drill_stats.best_streak || 0, runStreak);
  }
  if (runComplete) g.drill_stats.runs = (g.drill_stats.runs || 0) + 1;
  if (earnedToken) {
    g.hint_tokens += 1;
    g.drill_stats.tokens_earned = (g.drill_stats.tokens_earned || 0) + 1;
    logActivity(progress, 'drill_token');
  }
  saveProgress(progress);
  res.json({ ok: true, hint_tokens: g.hint_tokens, drill_stats: g.drill_stats });
});

// GET /api/drills — serve the authored pattern-recognition drill bank.
let _drillsCache = null;
app.get('/api/drills', (req, res) => {
  if (_drillsCache) return res.json(_drillsCache);
  const drillsPath = path.join(REPO_ROOT, 'docs', '_data', 'drills.yml');
  try {
    if (!fs.existsSync(drillsPath)) return res.json({ patterns: {} });
    const raw = fs.readFileSync(drillsPath, 'utf8');
    const parsed = yaml.load(raw) || {};
    _drillsCache = { patterns: parsed };
    res.json(_drillsCache);
  } catch (e) {
    console.error('Failed to load drills.yml:', e.message);
    res.status(500).json({ error: 'Failed to load drills' });
  }
});

// ── GET /api/runner-status ────────────────────────────────────────────────────

app.get('/api/runner-status', (req, res) => {
  res.json({
    available: runnerAvailable.cpp,        // legacy field — kept for older clients
    cpp:        runnerAvailable.cpp,
    java:       runnerAvailable.java,
    python:     runnerAvailable.python,
    javascript: runnerAvailable.javascript,
    go:         runnerAvailable.go,
  });
});

// ── Static frontend ───────────────────────────────────────────────────────────

// Serve whatever static assets exist. The catch-all below re-checks index.html
// on every request rather than only at startup: in `npm run dev`, vite runs
// concurrently and creates dist/ (and dist/assets/) BEFORE it finishes writing
// index.html, so a startup-only guard would lock in a route that then throws
// ENOENT on sendFile until the build lands. Checking per-request also recovers
// from a partial/interrupted prior build (dist exists, index.html missing).
const INDEX_HTML = path.join(DIST_DIR, 'index.html');

app.use(express.static(DIST_DIR));
app.get('*', (req, res) => {
  if (fs.existsSync(INDEX_HTML)) {
    return res.sendFile(INDEX_HTML);
  }
  res
    .status(503)
    .send('<h2>Frontend not built yet — the build may still be running. Refresh in a few seconds.</h2>');
});

// ── Start server ──────────────────────────────────────────────────────────────

const PORT = parseInt(process.env.PORT || '3000', 10);

function startServer(port) {
  const server = app.listen(port, () => {
    console.log(`Dashboard running at http://localhost:${port}`);
    if (process.env.PROGRESS_JSON_PATH) {
      console.log(`📝 Using progress file: ${PROGRESS_JSON}`);
    }
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
