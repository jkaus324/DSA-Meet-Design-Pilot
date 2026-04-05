import React, { useState, useEffect, useRef, useCallback } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { api } from '../lib/api.js';
import MarkdownRenderer from '../components/MarkdownRenderer.jsx';
import PatternBadge from '../components/PatternBadge.jsx';
import TierBadge from '../components/TierBadge.jsx';
import CopyCommand from '../components/CopyCommand.jsx';
import CodeEditor from '../components/CodeEditor.jsx';
import TestOutput from '../components/TestOutput.jsx';
import AiReviewPanel from '../components/AiReviewPanel.jsx';
import DifficultyModeSelector from '../components/DifficultyModeSelector.jsx';
import StepperTimeline from '../components/StepperTimeline.jsx';
import PartProgressBar from '../components/PartProgressBar.jsx';
import CarryForwardDialog from '../components/CarryForwardDialog.jsx';
import { SkeletonProblemView } from '../components/SkeletonLoader.jsx';
import { useToast } from '../components/Toast.jsx';

const KBD_STYLE = {
  display: 'inline-flex',
  alignItems: 'center',
  padding: '1px 6px',
  borderRadius: '9999px',
  fontSize: '10px',
  fontWeight: 500,
  lineHeight: '18px',
  background: 'var(--color-surface-tertiary)',
  color: 'var(--color-text-tertiary)',
  border: '1px solid var(--color-border)',
  marginLeft: '6px',
  whiteSpace: 'nowrap',
};

const TABS_LEFT  = ['Description', 'Notes'];
const TABS_RIGHT = ['Code', 'Design', 'AI Prompt'];

function canViewDesign(mode, parts) {
  if (mode === 'learning') return true;
  if (!parts || parts.length === 0) return false;
  if (mode === 'guided') return parts[0]?.status === 'passed';
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

// ── Chevron-left icon ──────────────────────────────────────────────────────────
function ChevronLeft() {
  return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
      <path d="M10 12L6 8l4-4" />
    </svg>
  );
}

// ── Arrow icons for prev/next navigation ───────────────────────────────────────
function ArrowLeftIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
      <path d="M8.5 10.5L5 7l3.5-3.5" />
    </svg>
  );
}

function ArrowRightIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
      <path d="M5.5 3.5L9 7l-3.5 3.5" />
    </svg>
  );
}

// ── Pause / Play icons for timer ───────────────────────────────────────────────
function PauseIcon() {
  return (
    <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
      <rect x="2" y="1.5" width="3" height="9" rx="0.8" />
      <rect x="7" y="1.5" width="3" height="9" rx="0.8" />
    </svg>
  );
}

function PlayIcon() {
  return (
    <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
      <path d="M3 1.5l7 4.5-7 4.5V1.5z" />
    </svg>
  );
}

// ── Status badge ───────────────────────────────────────────────────────────────
function StatusBadge({ status }) {
  const cfg = {
    solved:    { label: 'Solved',      bg: 'rgba(34,197,94,0.12)',  color: '#22c55e', dot: '#22c55e' },
    attempted: { label: 'In Progress', bg: 'rgba(245,158,11,0.12)', color: '#f59e0b', dot: '#f59e0b' },
    unsolved:  { label: 'Not Started', bg: 'var(--color-surface-tertiary)', color: 'var(--color-text-tertiary)', dot: 'var(--color-border)' },
  };
  const c = cfg[status] || cfg.unsolved;
  return (
    <span
      className="inline-flex items-center gap-1.5 text-xs font-medium px-2.5 py-1 rounded-full"
      style={{ background: c.bg, color: c.color }}
    >
      <span className="w-1.5 h-1.5 rounded-full flex-shrink-0" style={{ background: c.dot }} />
      {c.label}
    </span>
  );
}

// ── Timer hook ─────────────────────────────────────────────────────────────────
function useStopwatch(problemId) {
  const startRef = useRef(Date.now());
  const [elapsed, setElapsed] = useState(0);
  const [paused, setPaused] = useState(false);
  const frozenElapsed = useRef(0);

  useEffect(() => {
    startRef.current = Date.now();
    setElapsed(0);
    setPaused(false);
    frozenElapsed.current = 0;
  }, [problemId]);

  useEffect(() => {
    if (paused) return;
    const interval = setInterval(() => {
      setElapsed(Math.floor((Date.now() - startRef.current) / 1000));
    }, 1000);
    return () => clearInterval(interval);
  }, [paused]);

  const toggle = useCallback(() => {
    if (paused) {
      // Resuming: shift start time so elapsed stays continuous
      startRef.current = Date.now() - (frozenElapsed.current * 1000);
      setPaused(false);
    } else {
      // Pausing: freeze the current elapsed value
      frozenElapsed.current = elapsed;
      setPaused(true);
    }
  }, [paused, elapsed]);

  const mins = String(Math.floor(elapsed / 60)).padStart(2, '0');
  const secs = String(elapsed % 60).padStart(2, '0');

  return { display: `${mins}:${secs}`, paused, toggle };
}

export default function ProblemView({ onProgressChange }) {
  const { id } = useParams();
  const navigate = useNavigate();
  const toast = useToast();

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

  const [code,         setCode]         = useState('');
  const [isDirty,      setIsDirty]      = useState(false);
  const [codeSaved,    setCodeSaved]    = useState(false);
  const [submitResult, setSubmitResult] = useState(null);
  const [submitting,   setSubmitting]   = useState(false);
  const [submitStatus, setSubmitStatus] = useState('');

  const [mode, setMode] = useState('interview');

  const [leftTab,  setLeftTab]  = useState('Description');
  const [rightTab, setRightTab] = useState('Code');

  const [designLoading, setDesignLoading] = useState(false);
  const [aiLoading,     setAiLoading]     = useState(false);

  const [carryDialog, setCarryDialog] = useState(null);

  // ── Problem list for prev/next navigation ──────────────────────────────────
  const [problemList, setProblemList] = useState([]);

  // ── Timer ──────────────────────────────────────────────────────────────────
  const timer = useStopwatch(id);

  // ── Resizable test output panel ────────────────────────────────────────────
  const [executionPanelHeight, setExecutionPanelHeight] = useState(240);
  const [showExecutionPanel, setShowExecutionPanel] = useState(true);
  const testDragging = useRef(false);
  const testDragStartY = useRef(0);
  const testDragStartH = useRef(240);

  // ── Reset code confirmation ────────────────────────────────────────────────
  const [showResetConfirm, setShowResetConfirm] = useState(false);

  const containerRef = useRef(null);
  const editorRef = useRef(null);
  const [leftWidth, setLeftWidth] = useState(42);
  const dragging = useRef(false);

  const onMouseDown = () => { dragging.current = true; };
  useEffect(() => {
    const onMove = (e) => {
      if (!dragging.current || !containerRef.current) return;
      const rect = containerRef.current.getBoundingClientRect();
      setLeftWidth(Math.max(25, Math.min(65, ((e.clientX - rect.left) / rect.width) * 100)));
    };
    const onUp = () => { dragging.current = false; };
    window.addEventListener('mousemove', onMove);
    window.addEventListener('mouseup', onUp);
    return () => { window.removeEventListener('mousemove', onMove); window.removeEventListener('mouseup', onUp); };
  }, []);

  // ── Test panel vertical resize ─────────────────────────────────────────────
  useEffect(() => {
    const onMove = (e) => {
      if (!testDragging.current) return;
      const delta = testDragStartY.current - e.clientY;
      setExecutionPanelHeight(Math.max(120, Math.min(600, testDragStartH.current + delta)));
    };
    const onUp = () => { testDragging.current = false; };
    window.addEventListener('mousemove', onMove);
    window.addEventListener('mouseup', onUp);
    return () => { window.removeEventListener('mousemove', onMove); window.removeEventListener('mouseup', onUp); };
  }, []);

  const onTestDragStart = (e) => {
    testDragging.current = true;
    testDragStartY.current = e.clientY;
    testDragStartH.current = executionPanelHeight;
  };

  // Escape key -- close open dialogs
  useEffect(() => {
    const handler = (e) => {
      if (e.key === 'Escape') {
        if (showResetConfirm) { setShowResetConfirm(false); e.preventDefault(); return; }
        if (carryDialog) { setCarryDialog(null); e.preventDefault(); }
      }
    };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, [carryDialog, showResetConfirm]);

  // ── Prev/Next keyboard shortcuts (Alt+Left / Alt+Right) ───────────────────
  useEffect(() => {
    if (!problemList.length || !id) return;
    const handler = (e) => {
      if (e.altKey && e.key === 'ArrowLeft') {
        e.preventDefault();
        const idx = problemList.findIndex(p => p.id === id);
        if (idx > 0) navigate(`/problem/${problemList[idx - 1].id}`);
      }
      if (e.altKey && e.key === 'ArrowRight') {
        e.preventDefault();
        const idx = problemList.findIndex(p => p.id === id);
        if (idx < problemList.length - 1) navigate(`/problem/${problemList[idx + 1].id}`);
      }
    };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, [problemList, id, navigate]);

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
    setShowResetConfirm(false);

    Promise.all([api.getProblems(), api.getProblemParts(id)])
      .then(([problemsData, partsData]) => {
        const allProblems = problemsData.problems || [];
        setProblemList(allProblems);
        const p = allProblems.find(x => x.id === id);
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

  const handleModeChange = (newMode) => {
    const doSwitch = () => {
      if (isDirty) api.saveCode(id, mode, code).catch(console.error);
      api.updateStatus(id, { difficulty_mode: newMode }).then(() => onProgressChange?.()).catch(console.error);
      setMode(newMode);
      setIsDirty(false);
      setSubmitResult(null);
      api.getCode(id, newMode).then(d => setCode(d.code)).catch(console.error);
    };

    const hasPartsProgress = parts.some(p => p.status === 'passed' || p.status === 'attempted');
    if (hasPartsProgress) {
      const modeLabel = m => ({ interview: 'Interview', guided: 'Guided', learning: 'Learning' }[m] || m);
      const msg =
        `Switching from ${modeLabel(mode)} \u2192 ${modeLabel(newMode)} mode will reset your parts progress for this problem.\n\n` +
        `Your ${modeLabel(mode)} mode code is saved and can be restored by switching back.\n\nContinue?`;
      if (window.confirm(msg)) doSwitch();
    } else if (isDirty) {
      if (window.confirm(`Your code for ${mode} mode will be saved. Load ${newMode} mode?`)) doSwitch();
    } else {
      doSwitch();
    }
  };

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

  const handleFormat = useCallback(() => {
    if (editorRef.current) {
      editorRef.current.getAction('editor.action.formatDocument').run();
      toast.success('Code formatted');
    }
  }, [toast]);

  const handleCodeChange = (val) => { setCode(val); setIsDirty(true); };

  const handleSaveCode = useCallback(() => {
    api.saveCode(id, mode, code)
      .then(() => {
        setCodeSaved(true);
        setIsDirty(false);
        setTimeout(() => setCodeSaved(false), 2000);
        toast.success('Code saved');
      })
      .catch((err) => {
        console.error(err);
        toast.error('Failed to save code');
      });
  }, [id, mode, code, toast]);

  const handleSubmit = useCallback(async () => {
    const currentPart = getCurrentPart(parts);
    setSubmitting(true);
    setSubmitResult(null);
    setSubmitStatus('Compiling...');
    api.saveCode(id, mode, code).catch(console.error);
    try {
      const result = await api.submitPart(id, currentPart, mode, code);
      setSubmitResult(result);
      setSubmitting(false);
      setSubmitStatus('');
      if (result.success) {
        toast.success(`Part ${currentPart} passed!`);
      } else {
        toast.error(`Tests failed for Part ${currentPart}`);
      }
      if (result.success && result.unlocked_next_part) {
        await refreshParts();
        const problemsData = await api.getProblems();
        const updated = (problemsData.problems || []).find(x => x.id === id);
        if (updated) setProblem(updated);
        onProgressChange?.();
        const nextPartNum = currentPart + 1;
        const nextPartDef = parts.find(p => p.part === nextPartNum);
        if (nextPartDef) setCarryDialog({ partNum: nextPartNum, partName: nextPartDef.name });
      } else if (result.success) {
        await refreshParts();
        const problemsData = await api.getProblems();
        const updated = (problemsData.problems || []).find(x => x.id === id);
        if (updated) setProblem(updated);
        onProgressChange?.();
      }
    } catch (err) {
      setSubmitResult({ success: false, submitted_part: getCurrentPart(parts), compilation: { success: false, errors: err.message }, parts: [], runner_available: runnerAvail });
      setSubmitting(false);
      setSubmitStatus('');
      toast.error('Submission failed');
    }
  }, [parts, id, mode, code, toast, refreshParts, onProgressChange, runnerAvail]);

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
      if (nextPartDef) setCarryDialog({ partNum: nextPartNum, partName: nextPartDef.name });
    } catch (err) { console.error('Skip failed:', err); }
  };

  const handleCarryContinue = async () => {
    await api.setCarryForward(id, carryDialog.partNum, true).catch(console.error);
    await refreshParts();
    setCarryDialog(null);
  };

  const handleCarryFresh = async () => {
    await api.setCarryForward(id, carryDialog.partNum, false).catch(console.error);
    try {
      const starter = await api.getStarter(id, mode, carryDialog.partNum);
      setCode(starter.code);
      setIsDirty(false);
    } catch (err) { console.error('Failed to load starter:', err); }
    await refreshParts();
    setCarryDialog(null);
  };

  // ── Reset code handler ─────────────────────────────────────────────────────
  const handleResetCode = useCallback(async () => {
    const currentPart = getCurrentPart(parts);
    try {
      const starter = await api.getStarter(id, mode, currentPart);
      setCode(starter.code);
      setIsDirty(false);
      setShowResetConfirm(false);
      toast.success(`Code reset to starter for Part ${currentPart}`);
    } catch (err) {
      console.error('Failed to load starter:', err);
      toast.error('Failed to reset code');
      setShowResetConfirm(false);
    }
  }, [id, mode, parts, toast]);

  const handleRightTab = (tab) => {
    // Allow clicking Design even when locked (for blurred preview)
    if (tab === 'Design' && !canViewDesign(mode, parts)) {
      setRightTab(tab);
      if (!design) {
        setDesignLoading(true);
        api.getProblemDesign(id)
          .then(d => { setDesign(d.html); setDesignLoading(false); })
          .catch(() => { setDesign('<p>DESIGN.md not found.</p>'); setDesignLoading(false); });
      }
      return;
    }
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

  if (loading) {
    return <SkeletonProblemView />;
  }
  if (error) {
    return (
      <div className="flex items-center justify-center h-full">
        <p className="text-red-500 text-sm">{error}</p>
      </div>
    );
  }

  const num            = id.split('-')[0];
  const designUnlocked = canViewDesign(mode, parts);
  const command        = `./run-tests.sh ${id} ${(problem.languages || ['cpp'])[0]}`;
  const currentPartNum = getCurrentPart(parts);
  const allPartsPassed = parts.length > 0 && parts.every(p => p.status === 'passed');
  const totalParts     = problem.total_parts || 1;
  const partsPassed    = problem.parts_passed || 0;

  // ── Prev/Next problem computation ──────────────────────────────────────────
  const currentIdx = problemList.findIndex(p => p.id === id);
  const prevProblem = currentIdx > 0 ? problemList[currentIdx - 1] : null;
  const nextProblem = currentIdx >= 0 && currentIdx < problemList.length - 1 ? problemList[currentIdx + 1] : null;

  // ── Line / Char count ──────────────────────────────────────────────────────
  const lineCount = code ? code.split('\n').length : 0;
  const charCount = code ? code.length : 0;

  return (
    <>
      {carryDialog && (
        <CarryForwardDialog
          partNum={carryDialog.partNum}
          partName={carryDialog.partName}
          mode={mode}
          onContinue={handleCarryContinue}
          onStartFresh={handleCarryFresh}
        />
      )}

      {/* Reset Code Confirmation Dialog */}
      {showResetConfirm && (
        <div
          style={{
            position: 'fixed', inset: 0, zIndex: 1000,
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            background: 'rgba(0,0,0,0.5)',
          }}
          onClick={() => setShowResetConfirm(false)}
        >
          <div
            style={{
              background: 'var(--color-surface)',
              border: '1px solid var(--color-border)',
              borderRadius: 12,
              padding: '24px',
              maxWidth: 380,
              width: '90%',
              boxShadow: '0 8px 32px rgba(0,0,0,0.3)',
            }}
            onClick={e => e.stopPropagation()}
          >
            <p style={{ color: 'var(--color-text-primary)', fontSize: 14, fontWeight: 600, marginBottom: 8 }}>
              Reset Code?
            </p>
            <p style={{ color: 'var(--color-text-secondary)', fontSize: 13, lineHeight: 1.5, marginBottom: 20 }}>
              Reset to starter code for Part {currentPartNum}? Your current code will be lost.
            </p>
            <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end' }}>
              <button
                onClick={() => setShowResetConfirm(false)}
                style={{
                  padding: '6px 14px', fontSize: 12, fontWeight: 500, borderRadius: 8,
                  background: 'var(--color-surface-tertiary)', color: 'var(--color-text-secondary)',
                  border: '1px solid var(--color-border)', cursor: 'pointer',
                }}
              >
                Cancel
              </button>
              <button
                onClick={handleResetCode}
                style={{
                  padding: '6px 14px', fontSize: 12, fontWeight: 500, borderRadius: 8,
                  background: 'rgba(239,68,68,0.12)', color: '#f87171',
                  border: '1px solid rgba(239,68,68,0.3)', cursor: 'pointer',
                }}
              >
                Reset Code
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="flex flex-col" style={{ height: 'calc(100vh - 52px)' }}>

        {/* ── TOP TOOLBAR ─────────────────────────────────────────────────── */}
        <div
          className="flex-shrink-0 flex items-center gap-3 px-4 border-b border-border"
          style={{ height: '48px', background: 'var(--color-surface)' }}
        >
          {/* Back */}
          <button
            onClick={() => navigate('/')}
            className="flex items-center gap-1 text-text-tertiary hover:text-text-primary transition-colors text-xs font-medium flex-shrink-0"
          >
            <ChevronLeft />
            <span className="hidden sm:inline">Problems</span>
          </button>

          <span className="text-border flex-shrink-0" style={{ width: 1, height: 16, background: 'var(--color-border)' }} />

          {/* Previous problem button */}
          <button
            onClick={() => prevProblem && navigate(`/problem/${prevProblem.id}`)}
            disabled={!prevProblem}
            title={prevProblem ? `\u2190 ${prevProblem.id.split('-')[0]}. ${prevProblem.name}` : 'No previous problem'}
            className="flex items-center justify-center flex-shrink-0 transition-colors"
            style={{
              width: 26, height: 26, borderRadius: 6,
              color: prevProblem ? 'var(--color-text-secondary)' : 'var(--color-text-tertiary)',
              background: 'var(--color-surface-tertiary)',
              border: '1px solid var(--color-border)',
              cursor: prevProblem ? 'pointer' : 'not-allowed',
              opacity: prevProblem ? 1 : 0.4,
            }}
          >
            <ArrowLeftIcon />
          </button>

          {/* Problem title + badges */}
          <div className="flex items-center gap-2 flex-1 min-w-0">
            <span className="text-xs font-mono text-text-tertiary flex-shrink-0">{num}.</span>
            <h1 className="text-sm font-semibold text-text-primary truncate">{problem.name}</h1>
            <div className="hidden md:flex items-center gap-1.5 flex-shrink-0">
              <TierBadge tier={problem.tier} />
              {(problem.patterns || []).slice(0, 2).map(p => <PatternBadge key={p} pattern={p} />)}
            </div>
            {/* Company tags */}
            {(problem.companies || []).length > 0 && (
              <div className="hidden lg:flex items-center gap-1 flex-shrink-0">
                {problem.companies.map(company => (
                  <Link
                    key={company}
                    to={`/problems?company=${encodeURIComponent(company)}`}
                    className="inline-flex items-center px-2 py-0.5 rounded transition-colors"
                    style={{
                      background: 'var(--color-surface-tertiary)',
                      color: 'var(--color-text-secondary)',
                      border: '1px solid var(--color-border)',
                      fontSize: 10,
                      fontWeight: 500,
                      lineHeight: '16px',
                      textDecoration: 'none',
                    }}
                    onMouseEnter={e => {
                      e.currentTarget.style.background = 'var(--color-accent-light)';
                      e.currentTarget.style.color = 'var(--color-accent)';
                      e.currentTarget.style.borderColor = 'var(--color-accent)';
                    }}
                    onMouseLeave={e => {
                      e.currentTarget.style.background = 'var(--color-surface-tertiary)';
                      e.currentTarget.style.color = 'var(--color-text-secondary)';
                      e.currentTarget.style.borderColor = 'var(--color-border)';
                    }}
                  >
                    {company}
                  </Link>
                ))}
              </div>
            )}
            <StatusBadge status={problem.status} />
          </div>

          {/* Next problem button */}
          <button
            onClick={() => nextProblem && navigate(`/problem/${nextProblem.id}`)}
            disabled={!nextProblem}
            title={nextProblem ? `${nextProblem.id.split('-')[0]}. ${nextProblem.name} \u2192` : 'No next problem'}
            className="flex items-center justify-center flex-shrink-0 transition-colors"
            style={{
              width: 26, height: 26, borderRadius: 6,
              color: nextProblem ? 'var(--color-text-secondary)' : 'var(--color-text-tertiary)',
              background: 'var(--color-surface-tertiary)',
              border: '1px solid var(--color-border)',
              cursor: nextProblem ? 'pointer' : 'not-allowed',
              opacity: nextProblem ? 1 : 0.4,
            }}
          >
            <ArrowRightIcon />
          </button>

          {/* Parts progress (center) */}
          {totalParts > 1 && (
            <div className="flex items-center gap-2 flex-shrink-0">
              <PartProgressBar
                totalParts={totalParts}
                partsPassed={partsPassed}
                currentPart={currentPartNum}
                size="sm"
              />
              <span className="text-xs text-text-tertiary hidden sm:inline">
                {partsPassed}/{totalParts}
              </span>
            </div>
          )}

          <span style={{ width: 1, height: 16, background: 'var(--color-border)', flexShrink: 0 }} />

          {/* Timer / Stopwatch */}
          <div
            className="flex items-center gap-1.5 flex-shrink-0"
            title="Time spent on this problem (Alt+Left/Right to navigate)"
          >
            <span
              className="text-xs font-mono tabular-nums"
              style={{ color: 'var(--color-text-tertiary)', minWidth: 38 }}
            >
              {timer.display}
            </span>
            <button
              onClick={timer.toggle}
              title={timer.paused ? 'Resume timer' : 'Pause timer'}
              className="flex items-center justify-center transition-colors"
              style={{
                width: 20, height: 20, borderRadius: 4,
                color: timer.paused ? 'var(--color-accent)' : 'var(--color-text-tertiary)',
                background: 'transparent',
                border: 'none',
                cursor: 'pointer',
              }}
            >
              {timer.paused ? <PlayIcon /> : <PauseIcon />}
            </button>
          </div>

          <span style={{ width: 1, height: 16, background: 'var(--color-border)', flexShrink: 0 }} />

          {/* Difficulty mode selector */}
          <div className="flex-shrink-0">
            <DifficultyModeSelector mode={mode} onChange={handleModeChange} compact />
          </div>
        </div>

        {/* ── SPLIT PANELS ────────────────────────────────────────────────── */}
        <div ref={containerRef} className="flex flex-1 min-h-0">

          {/* LEFT PANEL */}
          <div
            className="flex flex-col min-h-0 overflow-hidden"
            style={{ width: `${leftWidth}%`, borderRight: '1px solid var(--color-border)' }}
          >
            {/* Tab bar */}
            <div
              className="flex flex-shrink-0 border-b border-border px-1"
              style={{ background: 'var(--color-surface)', minHeight: '40px', alignItems: 'stretch' }}
            >
              {TABS_LEFT.map(tab => (
                <button
                  key={tab}
                  onClick={() => setLeftTab(tab)}
                  className="px-4 text-xs font-semibold transition-colors relative"
                  style={{
                    color: leftTab === tab ? 'var(--color-accent)' : 'var(--color-text-tertiary)',
                    borderBottom: leftTab === tab ? '2px solid var(--color-accent)' : '2px solid transparent',
                    marginBottom: -1,
                  }}
                >
                  {tab}
                </button>
              ))}

              {/* Primer badge in tab bar */}
              {problem.prerequisite_primer && problem.primer_read === false && (
                <Link
                  to={`/primer/${problem.prerequisite_primer}`}
                  className="ml-auto flex items-center gap-1 px-3 text-xs font-medium self-center rounded-full mr-1"
                  style={{ color: '#f59e0b', background: 'rgba(245,158,11,0.1)' }}
                >
                  {'\u26A0'} {problem.prerequisite_primer}
                </Link>
              )}
            </div>

            {/* Content */}
            <div className="flex-1 overflow-y-auto" style={{ background: 'var(--color-surface)' }}>
              {leftTab === 'Description' && (
                <div className="px-4 py-4 space-y-5">
                  {/* Scenario */}
                  {scenarioHtml && <MarkdownRenderer html={scenarioHtml} />}

                  {/* Parts accordion */}
                  {parts.length > 0 && (
                    <div>
                      {parts.length > 1 && (
                        <p className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-2.5">
                          Parts
                        </p>
                      )}
                      <StepperTimeline parts={parts} key={parts.map(p => p.status).join(',')} />
                    </div>
                  )}

                </div>
              )}

              {leftTab === 'Notes' && (
                <div className="px-4 py-4">
                  <textarea
                    value={notes}
                    onChange={handleNotesChange}
                    onBlur={() => { clearTimeout(notesTimer.current); saveNotes(notes); }}
                    rows={14}
                    placeholder="Your notes, observations, approach..."
                    className="w-full px-3 py-2.5 text-sm border border-border rounded-lg resize-y focus:outline-none focus:border-accent text-text-primary placeholder-text-tertiary transition-colors"
                    style={{ background: 'var(--color-surface-secondary)', color: 'var(--color-text-primary)', fontFamily: 'inherit' }}
                  />
                  <div className="flex justify-between items-center mt-2">
                    <span className="text-xs text-text-tertiary">{saving ? 'Saving\u2026' : 'Auto-saves on blur'}</span>
                    <button
                      onClick={() => saveNotes(notes)}
                      className="text-xs px-3 py-1.5 rounded-lg font-medium text-white transition-opacity hover:opacity-90"
                      style={{ background: 'var(--color-accent)' }}
                    >
                      Save
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* DRAG HANDLE */}
          <div
            onMouseDown={onMouseDown}
            className="w-1 flex-shrink-0 cursor-col-resize transition-colors"
            style={{ background: 'var(--color-border)' }}
            onMouseEnter={e => { e.currentTarget.style.background = 'var(--color-accent)'; }}
            onMouseLeave={e => { e.currentTarget.style.background = 'var(--color-border)'; }}
          />

          {/* RIGHT PANEL */}
          <div className="flex flex-col flex-1 min-w-0 min-h-0 overflow-hidden">

            {/* Tab bar */}
            <div
              className="flex flex-shrink-0 items-stretch border-b border-border px-1"
              style={{ background: 'var(--color-surface)', minHeight: '40px' }}
            >
              {TABS_RIGHT.map(tab => {
                const locked = tab === 'Design' && !designUnlocked;
                let lockTitle = '';
                if (locked) {
                  lockTitle = mode === 'interview' ? 'Pass all parts to unlock (click to preview)' : 'Pass Part 1 to unlock (click to preview)';
                }
                return (
                  <button
                    key={tab}
                    onClick={() => handleRightTab(tab)}
                    title={lockTitle}
                    className="px-4 text-xs font-semibold transition-colors relative"
                    style={{
                      color: rightTab === tab ? 'var(--color-accent)' : 'var(--color-text-tertiary)',
                      borderBottom: rightTab === tab ? '2px solid var(--color-accent)' : '2px solid transparent',
                      marginBottom: -1,
                      cursor: 'pointer',
                      opacity: locked && rightTab !== tab ? 0.5 : 1,
                    }}
                  >
                    {tab}{locked ? ' \uD83D\uDD12' : ''}
                  </button>
                );
              })}

              {/* Save + Reset buttons (right side) */}
              {rightTab === 'Code' && (
                <div className="ml-auto flex items-center gap-2 pr-2">
                   {/* Reset Code button -- subtle, not prominent */}
                   <button
                    onClick={() => setShowResetConfirm(true)}
                    className="flex items-center gap-1 text-xs px-2 py-1 rounded-lg font-medium transition-all"
                    style={{
                      color: 'var(--color-text-tertiary)',
                      background: 'transparent',
                      border: '1px solid transparent',
                      fontSize: 11,
                    }}
                    onMouseEnter={e => {
                      e.currentTarget.style.color = '#f87171';
                      e.currentTarget.style.background = 'rgba(239,68,68,0.08)';
                      e.currentTarget.style.borderColor = 'rgba(239,68,68,0.2)';
                    }}
                    onMouseLeave={e => {
                      e.currentTarget.style.color = 'var(--color-text-tertiary)';
                      e.currentTarget.style.background = 'transparent';
                      e.currentTarget.style.borderColor = 'transparent';
                    }}
                    title="Reset to starter code"
                  >
                    Reset
                  </button>

                  {/* Format button */}
                  <button
                    onClick={handleFormat}
                    className="flex items-center gap-1 text-xs px-2 py-1 rounded-lg font-medium transition-all"
                    style={{
                      color: 'var(--color-text-tertiary)',
                      background: 'transparent',
                      border: '1px solid transparent',
                      fontSize: 11,
                    }}
                    onMouseEnter={e => {
                      e.currentTarget.style.color = 'var(--color-accent)';
                      e.currentTarget.style.background = 'var(--color-accent-light)';
                      e.currentTarget.style.borderColor = 'var(--color-accent)';
                    }}
                    onMouseLeave={e => {
                      e.currentTarget.style.color = 'var(--color-text-tertiary)';
                      e.currentTarget.style.background = 'transparent';
                      e.currentTarget.style.borderColor = 'transparent';
                    }}
                    title="Format code (Shift+Alt+F)"
                  >
                    Format
                  </button>

                  {/* Save button */}
                  <button
                    onClick={handleSaveCode}
                    className="flex items-center gap-1.5 text-xs px-3 py-1 rounded-lg border border-border font-medium transition-all hover:border-border-hover"
                    style={{
                      color: codeSaved ? '#22c55e' : isDirty ? 'var(--color-text-primary)' : 'var(--color-text-tertiary)',
                      background: 'var(--color-surface-tertiary)',
                    }}
                  >
                    {codeSaved ? '\u2713 Saved' : isDirty ? '\u25CF Save' : 'Save'}
                    <span style={KBD_STYLE}>Ctrl+S</span>
                  </button>
                </div>
              )}
            </div>

            {/* Code tab */}
            {rightTab === 'Code' && (
              <div className="flex-1 flex flex-col min-h-0 overflow-hidden">
                {/* Editor area */}
                <div className="flex-1 min-h-0 relative" style={{ background: 'var(--color-surface)' }}>
                  <CodeEditor 
                    value={code} 
                    onChange={handleCodeChange} 
                    language="cpp" 
                    onSave={handleSaveCode} 
                    onSubmit={handleSubmit} 
                    onMount={(editor) => { editorRef.current = editor; }}
                  />
                  
                  {/* Floating Show Console button (only when hidden) */}
                  {!showExecutionPanel && (
                    <button
                      onClick={() => setShowExecutionPanel(true)}
                      className="absolute shadow-lg flex items-center gap-2 px-3 py-1.5 rounded-lg border border-border font-semibold transition-all hover:border-accent hover:text-accent"
                      style={{
                        bottom: 12,
                        right: 12,
                        background: 'var(--color-surface-secondary)',
                        color: 'var(--color-text-tertiary)',
                        fontSize: 11,
                        zIndex: 20,
                      }}
                    >
                      <span>▲</span> Show Console
                    </button>
                  )}
                </div>

                {/* Line/Char count footer */}
                <div
                  className="flex-shrink-0 flex items-center justify-between px-3"
                  style={{
                    height: 24,
                    background: 'var(--color-surface)',
                    borderTop: '1px solid var(--color-border)',
                  }}
                >
                  <span
                    style={{
                      fontSize: 10,
                      color: 'var(--color-text-tertiary)',
                      fontFamily: "'Fira Code', 'SF Mono', monospace",
                      letterSpacing: '0.02em',
                    }}
                  >
                    Lines: {lineCount} | Chars: {charCount}
                  </span>
                </div>

                {showExecutionPanel && (
                  <>
                    {/* Test output drag handle (horizontal) - now at the top of the region */}
                    <div
                      onMouseDown={onTestDragStart}
                      className="flex-shrink-0"
                      style={{
                        height: 6,
                        cursor: 'row-resize',
                        background: 'var(--color-border)',
                        borderTop: '1px solid var(--color-border)',
                        zIndex: 10,
                      }}
                      onMouseEnter={e => { e.currentTarget.style.background = 'var(--color-accent)'; }}
                      onMouseLeave={e => { e.currentTarget.style.background = 'var(--color-border)'; }}
                    />

                    {/* Resizable Execution Region (Submit button bar + Test output + Run locally command) */}
                    <div
                      className="flex-shrink-0 flex flex-col overflow-hidden"
                      style={{ height: executionPanelHeight, background: 'var(--color-surface)' }}
                    >
                      {/* Submit button bar */}
                      <div
                        className="flex-shrink-0 flex items-center gap-2 px-3 py-2.5"
                        style={{ borderBottom: '1px solid var(--color-border)', background: 'var(--color-surface)' }}
                      >
                        {allPartsPassed ? (
                          <div
                            className="flex-1 py-2 text-sm font-semibold text-center rounded-xl"
                            style={{ background: 'rgba(34,197,94,0.12)', color: '#22c55e' }}
                          >
                            All Parts Complete
                          </div>
                        ) : (
                          <>
                            <button
                              onClick={handleSubmit}
                              disabled={submitting || !runnerAvail}
                              title={!runnerAvail ? 'g++ required \u2014 install g++ or use Skip' : ''}
                              className="flex-1 py-2.5 text-sm font-semibold rounded-xl text-white transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                              style={{
                                background: submitting
                                  ? 'var(--color-border)'
                                  : 'linear-gradient(135deg, #6366f1 0%, #4f46e5 100%)',
                                boxShadow: submitting ? 'none' : '0 2px 8px rgba(99,102,241,0.35)',
                              }}
                            >
                              {submitting ? (
                                <span className="flex items-center justify-center gap-2">
                                  <span className="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                                  {submitStatus || 'Submitting\u2026'}
                                </span>
                              ) : (
                                <span className="flex items-center justify-center gap-2">
                                  {`\u25B6 Submit Part ${currentPartNum}`}
                                  <span style={{ ...KBD_STYLE, background: 'rgba(255,255,255,0.15)', color: 'rgba(255,255,255,0.7)', borderColor: 'rgba(255,255,255,0.2)', marginLeft: 0 }}>Ctrl+Enter</span>
                                </span>
                              )}
                            </button>
                            {!runnerAvail && (
                              <button
                                onClick={handleSkipPart}
                                className="text-xs text-text-tertiary hover:text-text-secondary underline flex-shrink-0 transition-colors"
                                title="Manually unlock next part (g++ not available)"
                              >
                                Skip {'\u2192'}
                              </button>
                            )}
                          </>
                        )}

                        <button
                          onClick={() => setShowExecutionPanel(false)}
                          className="ml-auto flex items-center justify-center p-1 hover:text-text-primary transition-colors"
                          style={{
                            color: 'var(--color-text-tertiary)',
                            background: 'transparent',
                            border: 'none',
                            cursor: 'pointer',
                            marginLeft: 'auto',
                          }}
                          title="Hide console"
                        >
                          <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M1 1l12 12M13 1L1 13" />
                          </svg>
                        </button>
                      </div>

                      {/* Test output (stretches to fill space) */}
                      <div
                        className="flex-1 overflow-y-auto"
                        style={{ background: 'var(--color-surface-secondary)' }}
                      >
                        <TestOutput result={submitResult} running={submitting} submitStatus={submitStatus} />
                      </div>

                      {/* Run locally command */}
                      <div
                        className="flex-shrink-0 border-t border-border px-3 py-2"
                        style={{ background: 'var(--color-surface)' }}
                      >
                        <CopyCommand command={command} />
                      </div>
                    </div>
                  </>
                )}
              </div>
            )}

            {/* Design tab */}
            {rightTab === 'Design' && (
              <div className="flex-1 overflow-y-auto px-5 py-4" style={{ background: 'var(--color-surface)', position: 'relative' }}>
                {!designUnlocked ? (
                  /* Blurred preview when locked */
                  <div style={{ position: 'relative', minHeight: '100%' }}>
                    {/* Blurred design content behind */}
                    {designLoading ? (
                      <div className="text-text-tertiary text-sm">Loading{'\u2026'}</div>
                    ) : design ? (
                      <div
                        style={{
                          filter: 'blur(8px)',
                          pointerEvents: 'none',
                          userSelect: 'none',
                          WebkitUserSelect: 'none',
                          opacity: 0.6,
                        }}
                      >
                        <MarkdownRenderer html={design} />
                      </div>
                    ) : (
                      <div
                        style={{
                          filter: 'blur(8px)',
                          pointerEvents: 'none',
                          opacity: 0.6,
                          padding: '20px',
                        }}
                      >
                        <div style={{ background: 'var(--color-surface-tertiary)', height: 24, width: '60%', borderRadius: 6, marginBottom: 16 }} />
                        <div style={{ background: 'var(--color-surface-tertiary)', height: 16, width: '90%', borderRadius: 4, marginBottom: 10 }} />
                        <div style={{ background: 'var(--color-surface-tertiary)', height: 16, width: '75%', borderRadius: 4, marginBottom: 10 }} />
                        <div style={{ background: 'var(--color-surface-tertiary)', height: 16, width: '85%', borderRadius: 4, marginBottom: 10 }} />
                        <div style={{ background: 'var(--color-surface-tertiary)', height: 16, width: '40%', borderRadius: 4 }} />
                      </div>
                    )}

                    {/* Overlay message */}
                    <div
                      style={{
                        position: 'absolute',
                        inset: 0,
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                        justifyContent: 'center',
                        zIndex: 2,
                      }}
                    >
                      <div
                        style={{
                          background: 'var(--color-surface)',
                          border: '1px solid var(--color-border)',
                          borderRadius: 16,
                          padding: '28px 32px',
                          textAlign: 'center',
                          boxShadow: '0 4px 24px rgba(0,0,0,0.12)',
                          maxWidth: 340,
                        }}
                      >
                        <div
                          style={{
                            width: 48, height: 48, borderRadius: 14,
                            display: 'flex', alignItems: 'center', justifyContent: 'center',
                            fontSize: 24, background: 'var(--color-surface-tertiary)',
                            margin: '0 auto 12px',
                          }}
                        >
                          {'\uD83D\uDD12'}
                        </div>
                        <p style={{ color: 'var(--color-text-primary)', fontSize: 14, fontWeight: 600, marginBottom: 6 }}>
                          Design walkthrough locked
                        </p>
                        <p style={{ color: 'var(--color-text-tertiary)', fontSize: 12, lineHeight: 1.5 }}>
                          Complete the problem to unlock this design walkthrough
                        </p>
                      </div>
                    </div>
                  </div>
                ) : designLoading ? (
                  <div className="text-text-tertiary text-sm">Loading{'\u2026'}</div>
                ) : (
                  <MarkdownRenderer html={design} />
                )}
              </div>
            )}

            {/* AI Prompt tab */}
            {rightTab === 'AI Prompt' && (
              <div className="flex-1 min-h-0 overflow-hidden">
                <AiReviewPanel
                  promptMarkdown={aiPrompt}
                  code={code}
                  problemName={problem?.name}
                  loading={aiLoading}
                />
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  );
}
