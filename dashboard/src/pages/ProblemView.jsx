import React, { useState, useEffect, useRef, useCallback } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { api } from '../lib/api.js';
import MarkdownRenderer from '../components/MarkdownRenderer.jsx';
import PatternBadge from '../components/PatternBadge.jsx';
import TierBadge from '../components/TierBadge.jsx';
import CopyCommand from '../components/CopyCommand.jsx';
import CodeEditor from '../components/CodeEditor.jsx';
import TestOutput from '../components/TestOutput.jsx';
import DifficultyModeSelector from '../components/DifficultyModeSelector.jsx';
import PartAccordion from '../components/PartAccordion.jsx';
import PartProgressBar from '../components/PartProgressBar.jsx';
import CarryForwardDialog from '../components/CarryForwardDialog.jsx';

const TABS_LEFT  = ['Description', 'Notes'];
const TABS_RIGHT = ['Code', 'Design', 'AI Prompt'];

// Design tab gating: Interview = all parts passed, Guided = Part 1 passed, Learning = always
function canViewDesign(mode, parts) {
  if (mode === 'learning') return true;
  if (!parts || parts.length === 0) return false;
  if (mode === 'guided') return parts[0]?.status === 'passed';
  // interview: all parts passed
  return parts.every(p => p.status === 'passed');
}

function getCurrentPart(parts) {
  if (!parts) return 1;
  const active = parts.find(p => p.status === 'attempted' || p.status === 'active');
  if (active) return active.part;
  const allPassed = parts.every(p => p.status === 'passed');
  if (allPassed) return parts.length;
  return 1;
}

export default function ProblemView({ onProgressChange }) {
  const { id } = useParams();
  const navigate = useNavigate();

  const [problem,       setProblem]       = useState(null);
  const [parts,         setParts]         = useState([]);
  const [scenarioHtml,  setScenarioHtml]  = useState(null);
  const [runnerAvail,   setRunnerAvail]   = useState(true);
  const [design,        setDesign]        = useState(null);
  const [aiPrompt,      setAiPrompt]      = useState(null);
  const [loading,       setLoading]       = useState(true);
  const [error,         setError]         = useState(null);

  const [notes,  setNotes]  = useState('');
  const [saving, setSaving] = useState(false);
  const notesTimer = useRef(null);

  const [code,        setCode]        = useState('');
  const [isDirty,     setIsDirty]     = useState(false);
  const [codeSaved,   setCodeSaved]   = useState(false);
  const [submitResult, setSubmitResult] = useState(null);
  const [submitting,   setSubmitting]   = useState(false);
  const [submitStatus, setSubmitStatus] = useState('');

  const [mode, setMode] = useState('interview');

  const [leftTab,  setLeftTab]  = useState('Description');
  const [rightTab, setRightTab] = useState('Code');

  const [designLoading, setDesignLoading] = useState(false);
  const [aiLoading,     setAiLoading]     = useState(false);

  // Carry-forward dialog state
  const [carryDialog, setCarryDialog] = useState(null); // { partNum, partName }

  const containerRef = useRef(null);
  const [leftWidth, setLeftWidth] = useState(45);
  const dragging = useRef(false);

  const onMouseDown = () => { dragging.current = true; };
  useEffect(() => {
    const onMove = (e) => {
      if (!dragging.current || !containerRef.current) return;
      const rect = containerRef.current.getBoundingClientRect();
      setLeftWidth(Math.max(25, Math.min(70, ((e.clientX - rect.left) / rect.width) * 100)));
    };
    const onUp = () => { dragging.current = false; };
    window.addEventListener('mousemove', onMove);
    window.addEventListener('mouseup', onUp);
    return () => { window.removeEventListener('mousemove', onMove); window.removeEventListener('mouseup', onUp); };
  }, []);

  // ── Load problem + parts ──────────────────────────────────────────────────

  const refreshParts = useCallback(() => {
    return api.getProblemParts(id).then(data => {
      setParts(data.parts || []);
      setScenarioHtml(data.scenario_html || null);
      setRunnerAvail(data.runner_available !== false);
    });
  }, [id]);

  useEffect(() => {
    setLoading(true);
    setSubmitResult(null);
    setDesign(null);
    setAiPrompt(null);
    setLeftTab('Description');
    setRightTab('Code');
    setIsDirty(false);
    setCarryDialog(null);

    Promise.all([api.getProblems(), api.getProblemParts(id)])
      .then(([problemsData, partsData]) => {
        const p = (problemsData.problems || []).find(x => x.id === id);
        if (!p) throw new Error(`Problem '${id}' not found`);
        setProblem(p);
        setNotes(p.notes || '');
        setParts(partsData.parts || []);
        setScenarioHtml(partsData.scenario_html || null);
        setRunnerAvail(partsData.runner_available !== false);

        const initialMode = p.difficulty_mode || 'interview';
        setMode(initialMode);
        return api.getCode(id, initialMode);
      })
      .then(codeData => { setCode(codeData.code); setLoading(false); })
      .catch(err => { setError(err.message); setLoading(false); });
  }, [id]);

  // ── Mode switch ───────────────────────────────────────────────────────────

  const handleModeChange = (newMode) => {
    const doSwitch = () => {
      if (isDirty) api.saveCode(id, mode, code).catch(console.error);
      api.updateStatus(id, { difficulty_mode: newMode })
        .then(() => onProgressChange?.())
        .catch(console.error);
      setMode(newMode);
      setIsDirty(false);
      setSubmitResult(null);
      api.getCode(id, newMode).then(d => setCode(d.code)).catch(console.error);
    };

    if (isDirty) {
      if (window.confirm(`Your code for ${mode} mode will be saved. Load ${newMode} mode?`)) doSwitch();
    } else {
      doSwitch();
    }
  };

  // ── Notes ─────────────────────────────────────────────────────────────────

  const saveNotes = useCallback((value) => {
    setSaving(true);
    api.updateStatus(id, { notes: value }).then(() => setSaving(false)).catch(() => setSaving(false));
  }, [id]);

  const handleNotesChange = (e) => {
    const val = e.target.value;
    setNotes(val);
    clearTimeout(notesTimer.current);
    notesTimer.current = setTimeout(() => saveNotes(val), 1500);
  };

  // ── Code editor ───────────────────────────────────────────────────────────

  const handleCodeChange = (val) => { setCode(val); setIsDirty(true); };

  const handleSaveCode = () => {
    api.saveCode(id, mode, code)
      .then(() => { setCodeSaved(true); setIsDirty(false); setTimeout(() => setCodeSaved(false), 2000); })
      .catch(console.error);
  };

  // ── Submit Part ───────────────────────────────────────────────────────────

  const handleSubmit = async () => {
    const currentPart = getCurrentPart(parts);
    setSubmitting(true);
    setSubmitResult(null);
    setSubmitStatus('Compiling...');

    // Auto-save code first
    api.saveCode(id, mode, code).catch(console.error);

    try {
      const result = await api.submitPart(id, currentPart, mode, code);
      setSubmitResult(result);
      setSubmitting(false);
      setSubmitStatus('');

      if (result.success && result.unlocked_next_part) {
        // Refresh parts to get updated statuses
        await refreshParts();
        // Refresh problem for updated status
        const problemsData = await api.getProblems();
        const updated = (problemsData.problems || []).find(x => x.id === id);
        if (updated) setProblem(updated);
        onProgressChange?.();

        // Show carry-forward dialog for next part
        const nextPartNum  = currentPart + 1;
        const nextPartDef  = parts.find(p => p.part === nextPartNum);
        if (nextPartDef) {
          setCarryDialog({ partNum: nextPartNum, partName: nextPartDef.name });
        }
      } else if (result.success) {
        // All parts done — refresh
        await refreshParts();
        const problemsData = await api.getProblems();
        const updated = (problemsData.problems || []).find(x => x.id === id);
        if (updated) setProblem(updated);
        onProgressChange?.();
      }
    } catch (err) {
      setSubmitResult({
        success: false,
        submitted_part: getCurrentPart(parts),
        compilation: { success: false, errors: err.message },
        parts: [],
        runner_available: runnerAvail,
      });
      setSubmitting(false);
      setSubmitStatus('');
    }
  };

  // ── Skip Part (when g++ unavailable) ─────────────────────────────────────

  const handleSkipPart = async () => {
    const currentPart = getCurrentPart(parts);
    try {
      await api.skipPart(id, currentPart);
      await refreshParts();
      const problemsData = await api.getProblems();
      const updated = (problemsData.problems || []).find(x => x.id === id);
      if (updated) setProblem(updated);
      onProgressChange?.();

      const nextPartNum = currentPart + 1;
      const nextPartDef = parts.find(p => p.part === nextPartNum);
      if (nextPartDef) {
        setCarryDialog({ partNum: nextPartNum, partName: nextPartDef.name });
      }
    } catch (err) {
      console.error('Skip failed:', err);
    }
  };

  // ── Carry-forward dialog actions ──────────────────────────────────────────

  const handleCarryContinue = async () => {
    await api.setCarryForward(id, carryDialog.partNum, true).catch(console.error);
    await refreshParts();
    setCarryDialog(null);
  };

  const handleCarryFresh = async () => {
    await api.setCarryForward(id, carryDialog.partNum, false).catch(console.error);
    // Load fresh starter for this part
    try {
      const starter = await api.getStarter(id, mode, carryDialog.partNum);
      setCode(starter.code);
      setIsDirty(false);
    } catch (err) {
      console.error('Failed to load starter:', err);
    }
    await refreshParts();
    setCarryDialog(null);
  };

  // ── Design / AI Prompt tab ────────────────────────────────────────────────

  const handleRightTab = (tab) => {
    if (tab === 'Design' && !canViewDesign(mode, parts)) return;
    setRightTab(tab);
    if (tab === 'Design' && !design) {
      setDesignLoading(true);
      api.getProblemDesign(id)
        .then(d => { setDesign(d.html); setDesignLoading(false); })
        .catch(() => { setDesign('<p>DESIGN.md not found.</p>'); setDesignLoading(false); });
    }
    if (tab === 'AI Prompt' && !aiPrompt) {
      setAiLoading(true);
      api.getProblemAiPrompt(id)
        .then(d => { setAiPrompt(d.markdown); setAiLoading(false); })
        .catch(() => { setAiPrompt('AI_REVIEW_PROMPT.md not found.'); setAiLoading(false); });
    }
  };

  // ─────────────────────────────────────────────────────────────────────────

  if (loading) return <div className="flex items-center justify-center h-full text-text-tertiary text-sm">Loading...</div>;
  if (error)   return <div className="flex items-center justify-center h-full text-red-500 text-sm">{error}</div>;

  const num              = id.split('-')[0];
  const designUnlocked   = canViewDesign(mode, parts);
  const command          = `./run-tests.sh ${id} ${(problem.languages || ['cpp'])[0]}`;
  const currentPartNum   = getCurrentPart(parts);
  const allPartsPassed   = parts.length > 0 && parts.every(p => p.status === 'passed');
  const totalParts       = problem.total_parts || 1;
  const partsPassed      = problem.parts_passed || 0;

  return (
    <>
      {/* Carry-forward dialog */}
      {carryDialog && (
        <CarryForwardDialog
          partNum={carryDialog.partNum}
          partName={carryDialog.partName}
          mode={mode}
          onContinue={handleCarryContinue}
          onStartFresh={handleCarryFresh}
        />
      )}

      <div ref={containerRef} className="flex h-full" style={{ height: 'calc(100vh - 3rem)' }}>

        {/* ── LEFT PANEL ── */}
        <div className="flex flex-col border-r border-border overflow-hidden" style={{ width: `${leftWidth}%` }}>

          {/* Problem header */}
          <div className="flex-shrink-0 px-4 pt-4 pb-3 border-b border-border bg-surface">
            <button onClick={() => navigate('/')} className="flex items-center gap-1 text-xs text-text-tertiary hover:text-text-primary mb-2 group">
              <span className="group-hover:-translate-x-0.5 transition-transform">←</span>
              Problems
            </button>
            <h1 className="text-base font-bold text-text-primary leading-snug">{num}. {problem.name}</h1>
            <div className="flex flex-wrap items-center gap-1.5 mt-1.5">
              <TierBadge tier={problem.tier} />
              {(problem.patterns || []).map(p => <PatternBadge key={p} pattern={p} />)}
              <span className="text-xs text-text-tertiary">{problem.time_minutes} min</span>
              <span className="text-xs text-text-tertiary">· {(problem.companies || []).join(', ')}</span>
            </div>

            {/* Parts progress bar */}
            {totalParts > 1 && (
              <div className="flex items-center gap-2 mt-2">
                <span className="text-xs text-text-tertiary">Parts:</span>
                <PartProgressBar
                  totalParts={totalParts}
                  partsPassed={partsPassed}
                  currentPart={currentPartNum}
                  size="sm"
                />
              </div>
            )}

            {/* Primer badge */}
            {problem.prerequisite_primer && (
              <div className="mt-1.5">
                {problem.primer_read === false ? (
                  <Link to={`/primer/${problem.prerequisite_primer}`}
                    className="inline-flex items-center gap-1.5 text-xs px-2 py-0.5 rounded-full border border-status-attempted/50 bg-status-attempted/10 text-status-attempted hover:opacity-80">
                    <span>⚠</span> Read {problem.prerequisite_primer} primer first →
                  </Link>
                ) : problem.primer_read === true ? (
                  <span className="inline-flex items-center gap-1.5 text-xs px-2 py-0.5 rounded-full border border-status-solved/40 bg-status-solved/10 text-status-solved">
                    <span>✓</span> {problem.prerequisite_primer} primer read
                  </span>
                ) : (
                  <Link to={`/primer/${problem.prerequisite_primer}`} className="text-xs text-accent hover:underline">
                    Primer: {problem.prerequisite_primer}
                  </Link>
                )}
              </div>
            )}
          </div>

          {/* Tab bar */}
          <div className="flex border-b border-border bg-surface flex-shrink-0">
            {TABS_LEFT.map(tab => (
              <button key={tab} onClick={() => setLeftTab(tab)}
                className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
                  leftTab === tab ? 'border-accent text-accent' : 'border-transparent text-text-secondary hover:text-text-primary'
                }`}>{tab}</button>
            ))}
          </div>

          {/* Left pane content */}
          <div className="flex-1 overflow-y-auto bg-surface">
            {leftTab === 'Description' && (
              <div className="px-4 py-4 space-y-4">
                {/* Scenario / full problem context (always visible) */}
                {scenarioHtml && (
                  <div>
                    <MarkdownRenderer html={scenarioHtml} />
                  </div>
                )}

                {/* Parts accordion */}
                {parts.length > 0 ? (
                  <div>
                    <p className="text-xs font-semibold text-text-tertiary uppercase tracking-wide mb-2">
                      Parts
                    </p>
                    <PartAccordion parts={parts} key={parts.map(p => p.status).join(',')} />
                  </div>
                ) : (
                  <p className="text-sm text-text-tertiary">Loading parts...</p>
                )}

                {/* Terminal command */}
                <div>
                  <p className="text-xs font-semibold text-text-tertiary uppercase tracking-wide mb-2">Terminal</p>
                  <CopyCommand command={command} />
                </div>
              </div>
            )}

            {leftTab === 'Notes' && (
              <div className="px-5 py-4">
                <textarea value={notes} onChange={handleNotesChange}
                  onBlur={() => { clearTimeout(notesTimer.current); saveNotes(notes); }}
                  rows={12} placeholder="Your notes, observations, approach..."
                  className="w-full px-3 py-2 text-sm border border-border rounded-lg resize-y focus:outline-none focus:border-accent text-text-primary placeholder-text-tertiary bg-surface"
                  style={{ background: 'var(--color-surface)' }} />
                <div className="flex justify-between items-center mt-2">
                  <span className="text-xs text-text-tertiary">{saving ? 'Saving...' : 'Auto-saves on blur'}</span>
                  <button onClick={() => saveNotes(notes)} className="text-xs px-3 py-1 bg-accent text-white rounded hover:opacity-90">Save</button>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* ── DRAG HANDLE ── */}
        <div onMouseDown={onMouseDown} className="w-1 flex-shrink-0 cursor-col-resize hover:bg-accent transition-colors bg-border" style={{ userSelect: 'none' }} />

        {/* ── RIGHT PANEL ── */}
        <div className="flex flex-col flex-1 min-w-0 overflow-hidden">

          <DifficultyModeSelector mode={mode} onChange={handleModeChange} />

          {/* Tab bar + Save button */}
          <div className="flex-shrink-0 flex items-center justify-between px-3 border-b border-border bg-surface" style={{ minHeight: '42px' }}>
            <div className="flex">
              {TABS_RIGHT.map(tab => {
                const locked = tab === 'Design' && !designUnlocked;
                let lockTitle = '';
                if (locked) {
                  if (mode === 'interview') lockTitle = 'Pass all parts to unlock';
                  else if (mode === 'guided') lockTitle = 'Pass Part 1 to unlock';
                }
                return (
                  <button key={tab} onClick={() => handleRightTab(tab)}
                    title={lockTitle}
                    className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
                      rightTab === tab ? 'border-accent text-accent'
                        : locked ? 'border-transparent text-text-tertiary cursor-not-allowed'
                        : 'border-transparent text-text-secondary hover:text-text-primary'
                    }`}>
                    {tab}{locked ? ' 🔒' : ''}
                  </button>
                );
              })}
            </div>

            {rightTab === 'Code' && (
              <button onClick={handleSaveCode}
                className="text-xs px-3 py-1 rounded border border-border text-text-secondary hover:border-border-hover hover:text-text-primary transition-colors">
                {codeSaved ? '✓ Saved' : isDirty ? 'Save*' : 'Save'}
              </button>
            )}
          </div>

          {/* Code tab */}
          {rightTab === 'Code' && (
            <div className="flex-1 flex flex-col min-h-0 overflow-hidden">
              {/* Editor */}
              <div className="flex-1 min-h-0" style={{ background: 'var(--color-surface)' }}>
                <CodeEditor value={code} onChange={handleCodeChange} language="cpp" />
              </div>

              {/* Submit button */}
              <div className="flex-shrink-0 border-t border-border px-3 py-2 bg-surface flex items-center gap-2">
                {allPartsPassed ? (
                  <div className="flex-1 py-2 text-sm font-semibold text-center rounded-lg"
                    style={{ background: '#e6f9ed', color: '#01b328' }}>
                    ✅ All Parts Complete
                  </div>
                ) : (
                  <>
                    <button
                      onClick={handleSubmit}
                      disabled={submitting || !runnerAvail}
                      title={!runnerAvail ? 'g++ required. Install g++ or use Skip.' : ''}
                      className="flex-1 py-2 text-sm font-semibold rounded-lg text-white transition-opacity disabled:opacity-50"
                      style={{ background: submitting ? '#888' : '#01b328' }}
                    >
                      {submitting
                        ? (submitStatus || 'Submitting...')
                        : `▶ Submit Part ${currentPartNum}`}
                    </button>
                    {!runnerAvail && (
                      <button
                        onClick={handleSkipPart}
                        className="text-xs text-text-tertiary hover:text-text-secondary underline flex-shrink-0"
                        title="Manually unlock next part (g++ not available)"
                      >
                        Skip →
                      </button>
                    )}
                  </>
                )}
              </div>

              {/* Test output */}
              <div className="flex-shrink-0 border-t border-border overflow-y-auto"
                style={{ maxHeight: '220px', background: 'var(--color-surface-secondary)' }}>
                <TestOutput result={submitResult} running={submitting} submitStatus={submitStatus} />
              </div>
            </div>
          )}

          {/* Design tab */}
          {rightTab === 'Design' && (
            <div className="flex-1 overflow-y-auto px-5 py-4 bg-surface">
              {!designUnlocked ? (
                <div className="flex flex-col items-center justify-center h-full text-center">
                  <div className="text-4xl mb-3">🔒</div>
                  <p className="text-text-secondary text-sm font-medium">
                    {mode === 'interview'
                      ? 'Pass all parts to unlock the design walkthrough'
                      : 'Pass Part 1 to unlock the design walkthrough'}
                  </p>
                  <p className="text-text-tertiary text-xs mt-1">
                    Submit your solution to earn access
                  </p>
                </div>
              ) : designLoading ? (
                <div className="text-text-tertiary text-sm">Loading...</div>
              ) : (
                <MarkdownRenderer html={design} />
              )}
            </div>
          )}

          {/* AI Prompt tab */}
          {rightTab === 'AI Prompt' && (
            <div className="flex-1 overflow-y-auto px-5 py-4 bg-surface">
              {aiLoading ? (
                <div className="text-text-tertiary text-sm">Loading...</div>
              ) : aiPrompt ? (
                <div>
                  <p className="text-xs text-text-tertiary mb-2">Click inside to select all, then copy to your AI tool.</p>
                  <textarea readOnly value={aiPrompt}
                    className="w-full text-sm font-mono border border-border rounded-lg resize-y focus:outline-none px-3 py-2 text-text-primary"
                    style={{ background: 'var(--color-surface-tertiary)', minHeight: '400px' }}
                    onClick={e => e.target.select()} />
                </div>
              ) : null}
            </div>
          )}
        </div>
      </div>
    </>
  );
}
