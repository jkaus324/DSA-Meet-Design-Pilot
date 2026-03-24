import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import PatternBadge from './PatternBadge.jsx';
import TierBadge from './TierBadge.jsx';
import { STATUS_ICONS } from '../lib/constants.js';

const COLUMNS = [
  { key: 'status',       label: '',         width: 'w-10', sortable: false },
  { key: 'id',           label: '#',        width: 'w-12', sortable: true },
  { key: 'name',         label: 'Problem',  width: '',     sortable: true },
  { key: 'progress',     label: 'Parts',    width: 'w-20', sortable: false },
  { key: 'patterns',     label: 'Pattern',  width: 'w-40', sortable: false },
  { key: 'tier',         label: 'Tier',     width: 'w-12', sortable: true },
  { key: 'difficulty_mode', label: 'Mode',  width: 'w-10', sortable: false },
  { key: 'companies',    label: 'Company',  width: 'w-24', sortable: false },
  { key: 'time_minutes', label: 'Time',     width: 'w-12', sortable: true },
];

const MODE_DOT = {
  interview: { icon: '🔴', title: 'Interview mode' },
  guided:    { icon: '🟡', title: 'Guided mode' },
  learning:  { icon: '🟢', title: 'Learning mode' },
};

function PartDots({ totalParts, partsPassed, currentPart }) {
  if (!totalParts || totalParts <= 1) return null;

  const dots = [];
  for (let i = 1; i <= totalParts; i++) {
    let color;
    if (i <= partsPassed)    color = '#01b328'; // passed
    else if (i === currentPart) color = '#1a90ff'; // current
    else                     color = '#d1d5db'; // locked
    dots.push(<span key={i} style={{ width: 7, height: 7, borderRadius: '50%', background: color, display: 'inline-block', margin: '0 1px' }} />);
  }

  return (
    <span title={`${partsPassed} of ${totalParts} parts completed`}>
      {dots}
    </span>
  );
}

function sortProblems(problems, { col, dir }) {
  if (!col) return problems;
  return [...problems].sort((a, b) => {
    let av = a[col], bv = b[col];
    if (typeof av === 'string') av = av.toLowerCase();
    if (typeof bv === 'string') bv = bv.toLowerCase();
    if (av < bv) return dir === 'asc' ? -1 : 1;
    if (av > bv) return dir === 'asc' ? 1 : -1;
    return 0;
  });
}

export default function ProblemTable({ problems }) {
  const navigate = useNavigate();
  const [sort, setSort] = useState({ col: 'id', dir: 'asc' });

  const toggleSort = (col) => {
    setSort(s => s.col === col
      ? { col, dir: s.dir === 'asc' ? 'desc' : 'asc' }
      : { col, dir: 'asc' }
    );
  };

  const sorted = sortProblems(problems, sort);

  if (problems.length === 0) {
    return (
      <div className="text-center py-12 text-text-tertiary text-sm">
        No problems match the current filters.
      </div>
    );
  }

  return (
    <div className="border border-border rounded-lg overflow-hidden">
      <table className="w-full table-fixed text-sm">
        <thead className="bg-surface-tertiary border-b border-border">
          <tr>
            {COLUMNS.map(col => (
              <th
                key={col.key}
                className={`${col.width} px-3 py-2 text-left text-xs font-semibold text-text-tertiary uppercase tracking-wide ${
                  col.sortable ? 'cursor-pointer hover:text-text-primary select-none' : ''
                }`}
                onClick={() => col.sortable && toggleSort(col.key)}
              >
                {col.label}
                {col.sortable && sort.col === col.key && (
                  <span className="ml-1">{sort.dir === 'asc' ? '↑' : '↓'}</span>
                )}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-border bg-surface">
          {sorted.map(problem => {
            const { icon, color } = STATUS_ICONS[problem.status] || STATUS_ICONS.unsolved;
            const num = problem.id.split('-')[0];
            const totalParts  = problem.total_parts  || 1;
            const partsPassed = problem.parts_passed || 0;
            const currentPart = problem.current_part || 1;
            return (
              <tr
                key={problem.id}
                onClick={() => navigate(`/problem/${problem.id}`)}
                className="cursor-pointer hover:bg-surface-secondary transition-colors"
              >
                {/* Status */}
                <td className="w-10 px-3 py-2.5">
                  <span className="font-bold text-base" style={{ color }}>{icon}</span>
                </td>
                {/* # */}
                <td className="w-12 px-3 py-2.5 text-text-tertiary font-mono text-xs">
                  {num}
                </td>
                {/* Name */}
                <td className="px-3 py-2.5 font-medium text-text-primary truncate">
                  {problem.name}
                </td>
                {/* Parts progress */}
                <td className="w-20 px-3 py-2.5">
                  {totalParts > 1 ? (
                    <PartDots
                      totalParts={totalParts}
                      partsPassed={partsPassed}
                      currentPart={currentPart}
                    />
                  ) : (
                    <span className="text-xs text-text-tertiary">—</span>
                  )}
                </td>
                {/* Patterns */}
                <td className="w-40 px-3 py-2.5">
                  <div className="flex flex-wrap gap-1">
                    {(problem.patterns || []).map(p => (
                      <PatternBadge key={p} pattern={p} />
                    ))}
                  </div>
                </td>
                {/* Tier */}
                <td className="w-12 px-3 py-2.5">
                  <TierBadge tier={problem.tier} />
                </td>
                {/* Mode */}
                <td className="w-10 px-3 py-2.5 text-center" title={(MODE_DOT[problem.difficulty_mode] || MODE_DOT.interview).title}>
                  <span className="text-sm">
                    {problem.status === 'unsolved'
                      ? <span className="w-2 h-2 inline-block rounded-full bg-border" />
                      : (MODE_DOT[problem.difficulty_mode] || MODE_DOT.interview).icon}
                  </span>
                </td>
                {/* Company */}
                <td className="w-24 px-3 py-2.5 text-text-secondary truncate">
                  {(problem.companies || [])[0] || '—'}
                </td>
                {/* Time */}
                <td className="w-12 px-3 py-2.5 text-text-tertiary">
                  {problem.time_minutes}m
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
