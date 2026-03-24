import React, { useState } from 'react';

// ─── Compilation Error Block ─────────────────────────────────────────────────

function CompilationError({ errors }) {
  return (
    <div className="mx-4 mt-3 mb-2 rounded-lg border border-red-400/40 bg-red-500/10 overflow-hidden">
      <div className="px-3 py-2 border-b border-red-400/30 text-xs font-semibold text-red-400 uppercase tracking-wide">
        Compilation Error
      </div>
      <pre className="px-3 py-2 text-xs text-red-400 font-mono whitespace-pre-wrap leading-relaxed overflow-x-auto">
        {errors}
      </pre>
    </div>
  );
}

// ─── Single Part Result ──────────────────────────────────────────────────────

function PartResult({ part, name, passed, total, all_passed, tests }) {
  const [expanded, setExpanded] = useState(!all_passed);

  return (
    <div className="border-b border-border last:border-0">
      <button
        onClick={() => setExpanded(e => !e)}
        className="w-full flex items-center gap-3 px-4 py-2.5 text-left hover:bg-surface-secondary transition-colors"
      >
        <span className="text-sm flex-shrink-0">
          {all_passed ? '✅' : '❌'}
        </span>
        <span className="text-sm font-medium text-text-primary flex-1">
          Part {part}: {name}
        </span>
        <span className={`text-xs font-semibold flex-shrink-0 ${all_passed ? 'text-status-solved' : 'text-red-500'}`}>
          {passed}/{total}
        </span>
        <span className="text-text-tertiary text-xs flex-shrink-0 ml-1">
          {expanded ? '▾' : '▸'}
        </span>
      </button>

      {expanded && tests && tests.length > 0 && (
        <div className="pb-2 pl-10 pr-4">
          {tests.map((test, i) => (
            <div key={i} className="py-1">
              <div className="flex items-center gap-2">
                <span className="text-xs flex-shrink-0">{test.passed ? '✅' : '❌'}</span>
                <span className={`text-xs font-mono ${test.passed ? 'text-text-secondary' : 'text-red-400'}`}>
                  {test.name}
                </span>
              </div>
              {!test.passed && test.error && (
                <div className="ml-5 mt-0.5 text-xs text-red-400/80 font-mono whitespace-pre-wrap leading-relaxed">
                  {test.error}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

// ─── Main TestOutput ─────────────────────────────────────────────────────────

export default function TestOutput({ result, running, submitStatus }) {
  const [collapsed, setCollapsed] = useState(false);

  // Running spinner
  if (running) {
    return (
      <div className="flex items-center gap-2 px-4 py-3 text-sm text-text-secondary">
        <span className="inline-block w-3 h-3 border-2 border-accent border-t-transparent rounded-full animate-spin" />
        {submitStatus || 'Compiling and running...'}
      </div>
    );
  }

  if (!result) {
    return (
      <div className="px-4 py-3 text-sm text-text-tertiary">
        Click <span className="font-semibold text-accent">Submit Part</span> to run tests.
      </div>
    );
  }

  // ── Structured submit result (v3) ──
  if (result.submitted_part !== undefined) {
    const { success, submitted_part, compilation, parts, timed_out, time_ms, runner_available } = result;

    if (runner_available === false) {
      return (
        <div className="px-4 py-3 text-sm">
          <p className="text-red-400 font-medium mb-1">g++ not found on your system.</p>
          <p className="text-text-tertiary text-xs">
            Install g++ to use the test runner, or use "Skip to next part" below.
          </p>
        </div>
      );
    }

    return (
      <div className="flex flex-col">
        {/* Header bar */}
        <div className="flex items-center justify-between px-4 py-2 border-b border-border bg-surface">
          <div className={`flex items-center gap-2 text-xs font-semibold ${success ? 'text-status-solved' : 'text-red-500'}`}>
            <span>{success ? '✅' : '❌'}</span>
            <span>
              {timed_out ? 'Timed out (10s limit)' : success ? 'All tests passed!' : 'Some tests failed'}
            </span>
          </div>
          <div className="flex items-center gap-3">
            {time_ms !== undefined && (
              <span className="text-xs text-text-tertiary">{time_ms}ms</span>
            )}
            <button onClick={() => setCollapsed(c => !c)} className="text-xs text-text-tertiary hover:text-text-primary">
              {collapsed ? '▸' : '▾'}
            </button>
          </div>
        </div>

        {!collapsed && (
          <>
            {/* Compilation error */}
            {compilation && !compilation.success && (
              <CompilationError errors={compilation.errors} />
            )}

            {/* Per-part results */}
            {parts && parts.length > 0 && (
              <div>
                {parts.map(p => (
                  <PartResult key={p.part} {...p} />
                ))}
              </div>
            )}
          </>
        )}
      </div>
    );
  }

  // ── Legacy simple run result ──
  const { success, stage, output, error, time_ms } = result;

  const noCompiler = error && (
    error.includes('not recognized') ||
    error.includes('not found') ||
    error.includes('No such file')
  );

  if (noCompiler) {
    return (
      <div className="px-4 py-3 text-sm">
        <p className="text-red-400 font-medium mb-1">g++ not found on your system.</p>
        <p className="text-text-tertiary text-xs">
          Install g++ to use the Run button:<br />
          • <strong>Windows:</strong> Install MinGW-w64 and add to PATH<br />
          • <strong>Mac:</strong> <code className="bg-surface-tertiary px-1 rounded">xcode-select --install</code><br />
          • <strong>Linux:</strong> <code className="bg-surface-tertiary px-1 rounded">sudo apt install g++</code>
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      <div className={`flex items-center gap-3 px-4 py-2 border-b border-border text-xs font-medium ${
        success ? 'text-status-solved' : 'text-red-500'
      }`}>
        <span>{success ? '✓ Accepted' : stage === 'compile' ? '✗ Compilation Error' : '✗ Runtime Error'}</span>
        {time_ms !== undefined && <span className="text-text-tertiary font-normal">{time_ms}ms</span>}
      </div>
      <div className="flex-1 overflow-y-auto px-4 py-2 font-mono text-xs leading-relaxed">
        {output && (
          <div>
            <div className="text-text-tertiary mb-1 uppercase tracking-wide text-xs">Output</div>
            <pre className="text-text-primary whitespace-pre-wrap">{output}</pre>
          </div>
        )}
        {error && (
          <div className={output ? 'mt-3' : ''}>
            <div className="text-text-tertiary mb-1 uppercase tracking-wide text-xs">
              {stage === 'compile' ? 'Compiler Error' : 'Stderr'}
            </div>
            <pre className="text-red-400 whitespace-pre-wrap">{error}</pre>
          </div>
        )}
        {success && !output && !error && (
          <div className="text-text-tertiary italic">No output produced.</div>
        )}
      </div>
    </div>
  );
}
