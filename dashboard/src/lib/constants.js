export const PATTERN_COLORS = {
  Strategy:   { bg: 'rgba(99,102,241,0.1)',  text: 'var(--color-accent)' },
  Observer:   { bg: 'rgba(34,197,94,0.1)',   text: '#22c55e' },
  State:      { bg: 'rgba(245,158,11,0.1)',  text: '#f59e0b' },
  Singleton:  { bg: 'rgba(239,68,68,0.1)',   text: '#ef4444' },
  Factory:    { bg: 'rgba(139,92,246,0.1)',  text: '#8b5cf6' },
  Composite:  { bg: 'rgba(6,182,212,0.1)',   text: '#06b6d4' },
  Builder:    { bg: 'rgba(217,70,239,0.1)',  text: '#d946ef' },
  Comparator: { bg: 'rgba(16,185,129,0.1)',  text: '#10b981' },
  Decorator:  { bg: 'rgba(249,115,22,0.1)',  text: '#f97316' },
  Command:    { bg: 'rgba(99,102,241,0.1)',  text: '#818cf8' },
};

export const TIER_COLORS = {
  1: { bg: 'var(--color-tier-1-bg)', text: 'var(--color-tier-1)', label: 'Foundation' },
  2: { bg: 'var(--color-tier-2-bg)', text: 'var(--color-tier-2)', label: 'Intermediate' },
  3: { bg: 'var(--color-tier-3-bg)', text: 'var(--color-tier-3)', label: 'Advanced' },
};

export const STATUS_ICONS = {
  solved:    { icon: '✓', color: '#01b328' },
  attempted: { icon: '~', color: '#ffb800' },
  unsolved:  { icon: '○', color: '#dfdfdf' },
};

export const STATUS_LABELS = {
  solved:    'Solved',
  attempted: 'Attempted',
  unsolved:  'Unsolved',
};

export const TOPMATE_URL = 'https://topmate.io/jatin_kaushal24';
export const CODEJUNCTION_URL = 'https://topmate.io/jatin_kaushal24/2053177';

// ─── Part status display config (spec §12) ───────────────────────────────────
export const PART_STATUS = {
  locked:    { icon: '🔒', color: '#9ca3af',               bg: 'var(--color-surface-tertiary)', label: 'Locked'      },
  active:    { icon: '●',  color: 'var(--color-accent)',   bg: 'var(--color-accent-light)',     label: 'Current'     },
  attempted: { icon: '●',  color: '#f59e0b',               bg: 'rgba(245,158,11,0.1)',          label: 'In Progress' },
  passed:    { icon: '✅', color: '#22c55e',               bg: 'rgba(34,197,94,0.1)',           label: 'Passed'      },
};

// ─── Test result display config (spec §12) ───────────────────────────────────
export const TEST_RESULT = {
  passed: { icon: '✅', color: '#22c55e' },
  failed: { icon: '❌', color: '#f87171' },
};

export const DIFFICULTY_MODES = {
  interview: {
    label: 'Interview',
    icon: '🔴',
    color: '#dc2626',
    bg: '#fef2f2',
    description: 'Blank slate — design everything from scratch',
  },
  guided: {
    label: 'Guided',
    icon: '🟡',
    color: '#f59e0b',
    bg: '#fffbeb',
    description: 'Structural hints — you name and connect the pieces',
  },
  learning: {
    label: 'Learning',
    icon: '🟢',
    color: '#16a34a',
    bg: '#ecfdf5',
    description: 'Full scaffolding — implement the method bodies',
  },
};

// Reading time estimate per primer (words / 200 wpm)
export const PRIMER_READ_TIME = {
  strategy:  8,
  observer:  7,
  state:     8,
  singleton: 6,
  factory:   7,
  composite: 8,
  builder:   7,
};
