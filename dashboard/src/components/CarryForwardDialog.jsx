import React from 'react';

export default function CarryForwardDialog({ partNum, partName, mode, onContinue, onStartFresh }) {
  const modeLabel = mode === 'interview' ? 'Interview' : mode === 'guided' ? 'Guided' : 'Learning';

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/50" />

      {/* Dialog */}
      <div className="relative bg-surface border border-border rounded-2xl shadow-2xl max-w-lg w-full p-6 z-10">
        <div className="text-2xl mb-2">🎉</div>
        <h2 className="text-lg font-bold text-text-primary mb-1">
          Part {partNum - 1} complete!
        </h2>
        <p className="text-sm text-text-secondary mb-5">
          Part {partNum} — <span className="font-medium">{partName}</span> is now unlocked.
          How would you like to proceed?
        </p>

        <div className="grid grid-cols-2 gap-3 mb-5">
          {/* Continue with my code */}
          <button
            onClick={onContinue}
            className="flex flex-col text-left p-4 rounded-xl border-2 border-status-solved/40 bg-status-solved/5 hover:bg-status-solved/10 transition-colors group"
          >
            <div className="text-sm font-semibold text-text-primary mb-1.5 group-hover:text-status-solved transition-colors">
              Continue with my code
            </div>
            <div className="text-xs text-text-secondary leading-relaxed">
              Extend your Part {partNum - 1} solution — this is how real interviews work.
            </div>
          </button>

          {/* Start fresh */}
          <button
            onClick={onStartFresh}
            className="flex flex-col text-left p-4 rounded-xl border-2 border-border hover:border-border-hover bg-surface-secondary hover:bg-surface-tertiary transition-colors"
          >
            <div className="text-sm font-semibold text-text-primary mb-1.5">
              Start fresh
            </div>
            <div className="text-xs text-text-secondary leading-relaxed">
              Load Part {partNum} {modeLabel} starter. Use if your Part {partNum - 1} design needs a rewrite.
            </div>
          </button>
        </div>

        <p className="text-xs text-text-tertiary text-center">
          ✧ Recommended: Continue with your code
        </p>
      </div>
    </div>
  );
}
