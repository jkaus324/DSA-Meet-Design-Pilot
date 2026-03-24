export const PATTERN_COLORS = {
  Strategy:   { bg: '#e8f3ff', text: '#0066cc' },
  Observer:   { bg: '#e6f9ed', text: '#007a1f' },
  State:      { bg: '#fff8e1', text: '#996600' },
  Singleton:  { bg: '#ffeaea', text: '#cc0000' },
  Factory:    { bg: '#f3eeff', text: '#6600cc' },
  Composite:  { bg: '#e6fafa', text: '#006677' },
  Builder:    { bg: '#fdeeff', text: '#880099' },
  Comparator: { bg: '#edfff2', text: '#006622' },
  Decorator:  { bg: '#fff3e6', text: '#993300' },
  Command:    { bg: '#f5eeff', text: '#5500bb' },
};

export const TIER_COLORS = {
  1: { bg: '#e6f9ed', text: '#007a1f', label: 'T1' },
  2: { bg: '#fff8e1', text: '#996600', label: 'T2' },
  3: { bg: '#ffeaea', text: '#cc0000', label: 'T3' },
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
