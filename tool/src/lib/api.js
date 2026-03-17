const API_BASE = import.meta.env.VITE_API_BASE || '/api';

export async function fetchConfig(tenantId) {
  const res = await fetch(`${API_BASE}/v1/config/${tenantId}`);
  if (!res.ok) throw new Error('Failed to fetch config');
  return res.json();
}

export async function submitAudit(data, tenantId) {
  const res = await fetch(`${API_BASE}/v1/audit/${tenantId}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to submit audit');
  return res.json();
}

export async function trackEvent(event) {
  try {
    await fetch(`${API_BASE}/v1/analytics`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(event),
    });
  } catch (_) {
    // Analytics failures are non-fatal
  }
}
