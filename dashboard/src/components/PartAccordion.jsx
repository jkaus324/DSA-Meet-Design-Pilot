import React, { useState } from 'react';
import MarkdownRenderer from './MarkdownRenderer.jsx';

function PartHeader({ part, name, status, expanded, onToggle }) {
  const isLocked   = status === 'locked';
  const isPassed   = status === 'passed';
  const isCurrent  = status === 'active' || status === 'attempted';

  let icon, iconColor;
  if (isPassed)       { icon = '✅'; iconColor = '#01b328'; }
  else if (isLocked)  { icon = '🔒'; iconColor = '#9ca3af'; }
  else                { icon = '●';  iconColor = '#1a90ff'; }

  return (
    <button
      onClick={!isLocked ? onToggle : undefined}
      disabled={isLocked}
      className={`w-full flex items-center gap-3 px-4 py-3 text-left transition-colors ${
        isLocked
          ? 'cursor-not-allowed'
          : 'hover:bg-surface-secondary'
      } ${expanded ? 'bg-surface-secondary' : ''}`}
    >
      <span className="text-base leading-none flex-shrink-0" style={{ color: iconColor }}>
        {icon}
      </span>
      <div className="flex-1 min-w-0">
        <div className={`text-sm font-medium ${isLocked ? 'text-text-tertiary' : 'text-text-primary'}`}>
          Part {part}: {name}
        </div>
        {isPassed && (
          <div className="text-xs text-status-solved mt-0.5">Passed ✓</div>
        )}
        {isCurrent && (
          <div className="text-xs font-medium mt-0.5" style={{ color: '#1a90ff' }}>
            {status === 'attempted' ? 'In progress — keep going' : 'Current challenge'}
          </div>
        )}
        {isLocked && (
          <div className="text-xs text-text-tertiary mt-0.5">
            Complete Part {part - 1} to unlock
          </div>
        )}
      </div>
      {!isLocked && (
        <span className="text-text-tertiary text-xs flex-shrink-0">
          {expanded ? '▾' : '▸'}
        </span>
      )}
    </button>
  );
}

export default function PartAccordion({ parts }) {
  // Default: expand the current active/attempted part
  const defaultExpanded = () => {
    const current = parts.find(p => p.status === 'attempted' || p.status === 'active');
    if (current) return current.part;
    const lastPassed = [...parts].reverse().find(p => p.status === 'passed');
    if (lastPassed) return lastPassed.part;
    return parts[0]?.part || 1;
  };

  const [expandedPart, setExpandedPart] = useState(defaultExpanded());

  const toggle = (partNum) => {
    setExpandedPart(prev => prev === partNum ? null : partNum);
  };

  if (!parts || parts.length === 0) return null;

  return (
    <div className="border border-border rounded-lg overflow-hidden">
      {parts.map((p, i) => (
        <div key={p.part} className={i > 0 ? 'border-t border-border' : ''}>
          <PartHeader
            part={p.part}
            name={p.name}
            status={p.status}
            expanded={expandedPart === p.part}
            onToggle={() => toggle(p.part)}
          />
          {expandedPart === p.part && p.status !== 'locked' && (
            <div className="px-4 pb-4 pt-1 bg-surface-secondary border-t border-border">
              {p.description_html ? (
                <MarkdownRenderer html={p.description_html} />
              ) : (
                <p className="text-sm text-text-tertiary italic">
                  No description available for this part.
                </p>
              )}
            </div>
          )}
        </div>
      ))}
    </div>
  );
}
