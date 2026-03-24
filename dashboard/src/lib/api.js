const BASE = '/api';

async function request(url, options = {}) {
  const res = await fetch(url, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || `Request failed: ${res.status}`);
  }
  return res.json();
}

export const api = {
  getProblems() {
    return request(`${BASE}/problems`);
  },

  // v3: only notes and difficulty_mode accepted (status is derived from parts)
  updateStatus(id, body) {
    return request(`${BASE}/problems/${id}/status`, {
      method: 'POST',
      body: JSON.stringify(body),
    });
  },

  getProblemReadme(id) {
    return request(`${BASE}/problems/${id}/readme`);
  },

  getProblemDesign(id) {
    return request(`${BASE}/problems/${id}/design`);
  },

  getProblemAiPrompt(id) {
    return request(`${BASE}/problems/${id}/ai-prompt`);
  },

  // v3: get per-part descriptions and progress
  getProblemParts(id) {
    return request(`${BASE}/problems/${id}/parts`);
  },

  getPrimers() {
    return request(`${BASE}/primers`);
  },

  getPrimer(name) {
    return request(`${BASE}/primers/${name}`);
  },

  markPrimerRead(name) {
    return request(`${BASE}/primers/${name}/read`, { method: 'POST' });
  },

  getStats() {
    return request(`${BASE}/stats`);
  },

  getStarter(id, mode, part = 1) {
    return request(`${BASE}/problems/${id}/starter?mode=${mode}&part=${part}`);
  },

  getCode(id, mode) {
    return request(`${BASE}/problems/${id}/code?mode=${mode}`);
  },

  saveCode(id, mode, code) {
    return request(`${BASE}/problems/${id}/code`, {
      method: 'POST',
      body: JSON.stringify({ mode, code }),
    });
  },

  // v3: submit code for test-validated part progression
  submitPart(id, part, mode, code) {
    return request(`${BASE}/problems/${id}/submit`, {
      method: 'POST',
      body: JSON.stringify({ part, mode, code }),
    });
  },

  // v3: record carry-forward choice when new part unlocks
  setCarryForward(id, part, carryForward) {
    return request(`${BASE}/problems/${id}/parts/${part}/carry-forward`, {
      method: 'POST',
      body: JSON.stringify({ carry_forward: carryForward }),
    });
  },

  // v3: skip part (only when g++ unavailable)
  skipPart(id, part) {
    return request(`${BASE}/problems/${id}/parts/${part}/skip`, { method: 'POST' });
  },

  runnerStatus() {
    return request(`${BASE}/runner-status`);
  },
};
