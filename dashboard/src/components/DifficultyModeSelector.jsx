import React from 'react';

const MODES = [
  { id: 'interview', label: 'Interview', icon: '🔴', color: '#dc2626', bg: '#fef2f2', desc: 'Blank slate — design everything from scratch' },
  { id: 'guided',    label: 'Guided',    icon: '🟡', color: '#f59e0b', bg: '#fffbeb', desc: 'Structural hints — you name and connect the pieces' },
  { id: 'learning',  label: 'Learning',  icon: '🟢', color: '#16a34a', bg: '#ecfdf5', desc: 'Full scaffolding — implement the method bodies' },
];

export default function DifficultyModeSelector({ mode, onChange }) {
  return (
    <div className="flex items-center gap-2 px-3 py-2 border-b border-border bg-surface flex-shrink-0">
      <span className="text-xs text-text-tertiary font-medium mr-1">Difficulty:</span>
      {MODES.map(m => {
        const active = mode === m.id;
        return (
          <button
            key={m.id}
            onClick={() => !active && onChange(m.id)}
            title={m.desc}
            className="flex items-center gap-1.5 px-3 py-1 rounded-md text-xs font-semibold transition-all border"
            style={
              active
                ? { background: m.color, color: '#fff', borderColor: m.color }
                : { background: 'transparent', color: '#a3a3a3', borderColor: '#e5e5e5' }
            }
          >
            <span>{m.icon}</span>
            {m.label}
          </button>
        );
      })}
    </div>
  );
}
