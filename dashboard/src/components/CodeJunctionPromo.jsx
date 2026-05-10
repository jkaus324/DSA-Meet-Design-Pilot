import React, { useState, useEffect } from 'react';
import { Sparkles, X } from 'lucide-react';
import { CODEJUNCTION_URL } from '../lib/constants.js';

const STORAGE_KEY = 'cj-promo-dismissed';

/**
 * Floating promo card on problem pages. Appears immediately on mount
 * (and re-appears on a failed submit if not yet dismissed). One global
 * dismiss persists across all sessions.
 */
export default function CodeJunctionPromo() {
  const [visible, setVisible] = useState(false);
  const [animateIn, setAnimateIn] = useState(false);

  useEffect(() => {
    // If user previously dismissed, never show.
    try {
      if (localStorage.getItem(STORAGE_KEY) === 'true') return;
    } catch {}

    const reveal = () => {
      setVisible(true);
      // Allow the slide-in animation to play after first paint
      requestAnimationFrame(() => setAnimateIn(true));
    };

    reveal();

    // Also surface it on a failed submit, in case it had been hidden.
    const onFail = () => reveal();
    window.addEventListener('cj-submit-failed', onFail);

    return () => {
      window.removeEventListener('cj-submit-failed', onFail);
    };
  }, []);

  const dismiss = () => {
    try { localStorage.setItem(STORAGE_KEY, 'true'); } catch {}
    setAnimateIn(false);
    // Match the transition duration before unmounting
    setTimeout(() => setVisible(false), 200);
  };

  if (!visible) return null;

  return (
    <div
      role="dialog"
      aria-label="CodeJunction Pro promotion"
      style={{
        position: 'fixed',
        bottom: 24,
        right: 24,
        width: 320,
        maxWidth: 'calc(100vw - 48px)',
        zIndex: 60,
        background: 'linear-gradient(135deg, rgba(99,102,241,0.97), rgba(139,92,246,0.97))',
        color: '#fff',
        borderRadius: 14,
        padding: '14px 16px 14px 14px',
        boxShadow: '0 10px 30px rgba(0,0,0,0.25), 0 4px 12px rgba(99,102,241,0.35)',
        backdropFilter: 'blur(8px)',
        opacity: animateIn ? 1 : 0,
        transform: animateIn ? 'translateY(0)' : 'translateY(16px)',
        transition: 'opacity 200ms ease, transform 200ms ease',
      }}
    >
      <button
        onClick={dismiss}
        aria-label="Dismiss"
        title="Don't show this again"
        style={{
          position: 'absolute',
          top: 6,
          right: 6,
          width: 24,
          height: 24,
          display: 'inline-flex',
          alignItems: 'center',
          justifyContent: 'center',
          background: 'rgba(255,255,255,0.12)',
          border: 'none',
          borderRadius: 6,
          color: '#fff',
          cursor: 'pointer',
          transition: 'background 120ms ease',
        }}
        onMouseEnter={(e) => { e.currentTarget.style.background = 'rgba(255,255,255,0.25)'; }}
        onMouseLeave={(e) => { e.currentTarget.style.background = 'rgba(255,255,255,0.12)'; }}
      >
        <X size={14} />
      </button>

      <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 6 }}>
        <Sparkles size={16} />
        <span style={{ fontSize: 11, fontWeight: 700, letterSpacing: 0.6, textTransform: 'uppercase' }}>
          CodeJunction Pro
        </span>
      </div>

      <div style={{ fontSize: 14, fontWeight: 600, lineHeight: 1.35, marginBottom: 4 }}>
        Stuck? Want 80 more problems like this one?
      </div>

      <div style={{ fontSize: 12, opacity: 0.92, lineHeight: 1.45, marginBottom: 12 }}>
        100+ machine coding + LLD problems with dual-view editorials and 9 company tracks.
      </div>

      <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
        <a
          href={CODEJUNCTION_URL}
          target="_blank"
          rel="noopener noreferrer"
          style={{
            flex: 1,
            display: 'inline-flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '7px 12px',
            background: '#fff',
            color: 'var(--color-accent)',
            fontSize: 12,
            fontWeight: 700,
            borderRadius: 8,
            textDecoration: 'none',
            transition: 'transform 120ms ease',
          }}
          onMouseEnter={(e) => { e.currentTarget.style.transform = 'translateY(-1px)'; }}
          onMouseLeave={(e) => { e.currentTarget.style.transform = 'translateY(0)'; }}
        >
          Get it →
        </a>
        <button
          onClick={dismiss}
          style={{
            background: 'transparent',
            border: 'none',
            color: 'rgba(255,255,255,0.85)',
            fontSize: 11,
            cursor: 'pointer',
            padding: '4px 2px',
          }}
        >
          Not now
        </button>
      </div>
    </div>
  );
}
