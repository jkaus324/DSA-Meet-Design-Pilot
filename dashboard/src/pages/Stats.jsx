import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { api } from '../lib/api.js';

function StatusDot({ color, label, value }) {
  return (
    <div className="flex items-center gap-2">
      <span className="w-2.5 h-2.5 rounded-full flex-shrink-0" style={{ background: color }} />
      <span className="text-sm text-text-secondary">{label}</span>
      <span className="text-sm font-semibold text-text-primary ml-auto">{value}</span>
    </div>
  );
}

function TierRow({ tier, stats }) {
  const { solved = 0, attempted = 0, total = 0 } = stats || {};
  const unsolved = total - solved - attempted;
  const solvedPct   = total > 0 ? (solved   / total) * 100 : 0;
  const attemptedPct = total > 0 ? (attempted / total) * 100 : 0;
  const unsolvedPct  = total > 0 ? (unsolved  / total) * 100 : 0;

  return (
    <div className="py-3 border-b border-border last:border-0">
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm font-medium text-text-primary">Tier {tier}</span>
        <span className="text-xs text-text-secondary">
          {solved}/{total}
          {total > 0 && <span className="ml-1 text-text-tertiary">({Math.round(solvedPct)}%)</span>}
        </span>
      </div>
      {/* Stacked bar */}
      <div className="h-2 rounded-full overflow-hidden bg-surface-tertiary flex">
        <div className="h-full" style={{ width: `${solvedPct}%`, background: '#01b328', transition: 'width 0.5s ease' }} />
        <div className="h-full" style={{ width: `${attemptedPct}%`, background: '#ffb800', transition: 'width 0.5s ease' }} />
        <div className="h-full" style={{ width: `${unsolvedPct}%`, background: '#dfdfdf', transition: 'width 0.5s ease' }} />
      </div>
      <div className="flex gap-4 mt-1.5 text-xs text-text-tertiary">
        <span>{solved} solved</span>
        <span>{attempted} attempted</span>
        <span>{unsolved} unsolved</span>
      </div>
    </div>
  );
}

function PatternRow({ rank, pattern, data }) {
  return (
    <div className="flex items-center gap-3 py-2 border-b border-border last:border-0">
      <span className="text-xs text-text-tertiary w-4 text-right flex-shrink-0">{rank}.</span>
      <span className="text-sm text-text-primary flex-1">{pattern}</span>
      <span className="text-sm font-semibold text-text-primary">{data.solved}</span>
      <span className="text-xs text-text-tertiary">solved</span>
    </div>
  );
}

export default function Stats() {
  const navigate = useNavigate();
  const [stats,   setStats]   = useState(null);
  const [primers, setPrimers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error,   setError]   = useState(null);

  useEffect(() => {
    Promise.all([api.getStats(), api.getPrimers()])
      .then(([s, p]) => { setStats(s); setPrimers(p.primers || []); setLoading(false); })
      .catch(err => { setError(err.message); setLoading(false); });
  }, []);

  if (loading) return <div className="text-center py-16 text-text-tertiary text-sm">Loading...</div>;
  if (error)   return <div className="text-center py-16 text-red-500 text-sm">{error}</div>;

  const { overall, by_tier, by_pattern, by_difficulty_mode, parts_stats, streak } = stats;

  const patternEntries = Object.entries(by_pattern || {})
    .sort((a, b) => b[1].solved - a[1].solved);

  const formatDate = (d) => d
    ? new Date(d).toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })
    : 'No activity yet';

  const primersRead = primers.filter(p => p.read).length;

  return (
    <div className="max-w-xl mx-auto">

      {/* Page header */}
      <div className="mb-6">
        <h1 className="text-xl font-bold text-text-primary">Stats</h1>
        <p className="text-sm text-text-secondary mt-0.5">
          Progress breakdown across tiers, patterns, primers, and streak
        </p>
      </div>

      {/* ── Primary progress card ─────────────────────────────────────── */}
      <div className="bg-surface border border-border rounded-xl p-5 mb-4">
        <div className="flex items-center gap-4 mb-4">
          <div>
            <div className="text-4xl font-bold text-text-primary tracking-tight">
              {overall.percent_complete}%
            </div>
            <div className="text-sm text-text-secondary mt-0.5">
              {overall.solved}/{overall.total} solved
            </div>
          </div>
        </div>

        {/* Progress bar */}
        <div className="h-2.5 rounded-full bg-surface-tertiary overflow-hidden mb-4">
          <div
            className="h-full rounded-full"
            style={{
              width: `${overall.percent_complete}%`,
              background: '#01b328',
              transition: 'width 0.5s ease',
            }}
          />
        </div>

        {/* Status breakdown */}
        <div className="space-y-2">
          <StatusDot color="#01b328" label="Solved"   value={overall.solved} />
          <StatusDot color="#ffb800" label="Attempted" value={overall.attempted} />
          <StatusDot color="#dfdfdf" label="Unsolved"  value={overall.unsolved} />
        </div>
      </div>

      {/* ── Tier breakdown ────────────────────────────────────────────── */}
      <div className="bg-surface border border-border rounded-xl p-5 mb-4">
        <div className="flex items-center justify-between mb-1">
          <h2 className="text-sm font-semibold text-text-primary">Tier breakdown</h2>
          <span className="text-xs text-text-tertiary">Stacked by status</span>
        </div>
        {/* Legend */}
        <div className="flex gap-4 mb-3 text-xs text-text-tertiary">
          <span className="flex items-center gap-1.5"><span className="w-2 h-2 rounded-full inline-block" style={{background:'#01b328'}}/> Solved</span>
          <span className="flex items-center gap-1.5"><span className="w-2 h-2 rounded-full inline-block" style={{background:'#ffb800'}}/> Attempted</span>
          <span className="flex items-center gap-1.5"><span className="w-2 h-2 rounded-full inline-block" style={{background:'#dfdfdf'}}/> Unsolved</span>
        </div>
        <TierRow tier={1} stats={by_tier[1]} />
        <TierRow tier={2} stats={by_tier[2]} />
        <TierRow tier={3} stats={by_tier[3]} />
      </div>

      {/* ── Patterns practiced ────────────────────────────────────────── */}
      {patternEntries.length > 0 && (
        <div className="bg-surface border border-border rounded-xl p-5 mb-4">
          <h2 className="text-sm font-semibold text-text-primary mb-1">Patterns practiced</h2>
          <p className="text-xs text-text-tertiary mb-3">
            Patterns you've solved problems for, ranked by count
          </p>
          {patternEntries.map(([pattern, data], i) => (
            <PatternRow key={pattern} rank={i + 1} pattern={pattern} data={data} />
          ))}
        </div>
      )}

      {/* ── Difficulty mode distribution ─────────────────────────────── */}
      {by_difficulty_mode && (
        <div className="bg-surface border border-border rounded-xl p-5 mb-4">
          <h2 className="text-sm font-semibold text-text-primary mb-1">Difficulty mode</h2>
          <p className="text-xs text-text-tertiary mb-3">How you solved problems by mode</p>
          <div className="space-y-3">
            {[
              { id: 'interview', label: '🔴 Interview', color: '#dc2626' },
              { id: 'guided',    label: '🟡 Guided',    color: '#f59e0b' },
              { id: 'learning',  label: '🟢 Learning',  color: '#16a34a' },
            ].map(({ id, label, color }) => {
              const d = (by_difficulty_mode || {})[id] || { solved: 0, attempted: 0 };
              return (
                <div key={id} className="flex items-center gap-3">
                  <span className="text-xs text-text-secondary w-28 flex-shrink-0">{label}</span>
                  <div className="flex-1 flex items-center gap-2">
                    <div className="flex-1 h-2 rounded-full bg-surface-tertiary overflow-hidden flex">
                      <div className="h-full rounded-full" style={{ width: d.solved > 0 ? `${(d.solved / Math.max(overall.total, 1)) * 100}%` : '0%', background: color, transition: 'width 0.5s ease' }} />
                    </div>
                    <span className="text-xs text-text-secondary w-16 flex-shrink-0">
                      {d.solved} solved{d.attempted > 0 ? `, ${d.attempted} att.` : ''}
                    </span>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* ── Parts completion ──────────────────────────────────────────── */}
      {parts_stats && parts_stats.total_parts_across_all_problems > 0 && (
        <div className="bg-surface border border-border rounded-xl p-5 mb-4">
          <h2 className="text-sm font-semibold text-text-primary mb-1">Parts progress</h2>
          <p className="text-xs text-text-tertiary mb-3">
            How far you've progressed through each problem's parts
          </p>
          <div className="h-2.5 rounded-full bg-surface-tertiary overflow-hidden mb-3">
            <div
              className="h-full rounded-full"
              style={{
                width: `${Math.round((parts_stats.parts_passed / parts_stats.total_parts_across_all_problems) * 100)}%`,
                background: '#1a90ff',
                transition: 'width 0.5s ease',
              }}
            />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <div className="text-2xl font-bold text-text-primary">
                {parts_stats.parts_passed}
                <span className="text-sm font-normal text-text-tertiary ml-1">
                  / {parts_stats.total_parts_across_all_problems}
                </span>
              </div>
              <div className="text-xs text-text-secondary mt-0.5">parts passed</div>
            </div>
            <div>
              <div className="text-2xl font-bold text-text-primary">
                {parts_stats.problems_fully_completed}
              </div>
              <div className="text-xs text-text-secondary mt-0.5">fully completed</div>
            </div>
            <div>
              <div className="text-lg font-semibold text-text-primary">
                {parts_stats.problems_partially_completed}
              </div>
              <div className="text-xs text-text-secondary mt-0.5">in progress</div>
            </div>
            <div>
              <div className="text-lg font-semibold text-text-primary">
                {parts_stats.avg_parts_per_solve > 0 ? parts_stats.avg_parts_per_solve : '—'}
              </div>
              <div className="text-xs text-text-secondary mt-0.5">avg parts / solve</div>
            </div>
          </div>
        </div>
      )}

      {/* ── Primers ───────────────────────────────────────────────────── */}
      <div className="bg-surface border border-border rounded-xl p-5 mb-4">
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-sm font-semibold text-text-primary">Primers</h2>
          <span className="text-sm font-semibold" style={{ color: '#1a90ff' }}>
            {primersRead}/{primers.length} read
          </span>
        </div>
        <div className="h-1.5 rounded-full bg-surface-tertiary overflow-hidden mb-4">
          <div
            className="h-full rounded-full"
            style={{
              width: primers.length > 0 ? `${(primersRead / primers.length) * 100}%` : '0%',
              background: '#1a90ff',
              transition: 'width 0.5s ease',
            }}
          />
        </div>
        <div className="flex flex-wrap gap-2">
          {primers.map(p => (
            <Link
              key={p.name}
              to={`/primer/${p.name}`}
              className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border transition-colors ${
                p.read
                  ? 'border-transparent text-white'
                  : 'border-border text-text-secondary hover:border-border-hover hover:text-text-primary bg-surface'
              }`}
              style={p.read ? { background: '#01b328' } : {}}
            >
              {p.read && <span>✓</span>}
              {p.name}
            </Link>
          ))}
        </div>
      </div>

      {/* ── Activity / Streak ─────────────────────────────────────────── */}
      <div className="bg-surface border border-border rounded-xl p-5 mb-4">
        <h2 className="text-sm font-semibold text-text-primary mb-3">Activity</h2>
        <div className="flex items-start gap-8">
          <div>
            <div className="text-3xl font-bold text-text-primary">{streak.current_days}</div>
            <div className="text-xs text-text-tertiary mt-0.5">day streak</div>
          </div>
          <div className="flex-1">
            <div className="text-sm font-medium text-text-primary">{formatDate(streak.last_activity)}</div>
            <div className="text-xs text-text-tertiary mt-0.5">last activity</div>
            {streak.current_days === 0 && (
              <p className="text-xs text-text-secondary mt-2">
                Solve a problem today to start your streak.
              </p>
            )}
            {streak.current_days > 0 && (
              <p className="text-xs text-text-secondary mt-2">
                A streak counts if you solved or attempted any problem that day.
              </p>
            )}
          </div>
        </div>

        {/* CTA */}
        <button
          onClick={() => navigate('/')}
          className="mt-4 w-full py-2 text-sm font-medium rounded-lg border border-border text-text-secondary hover:text-text-primary hover:border-border-hover transition-colors"
        >
          Open a problem →
        </button>
      </div>

      {/* Footer */}
      <p className="text-center text-xs text-text-tertiary pb-6">
        Local dashboard · Runs on localhost · Progress saved to progress.json
      </p>
    </div>
  );
}
