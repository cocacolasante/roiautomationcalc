function getBase(apiBase) {
  if (apiBase) return `${apiBase}/api`;
  return import.meta.env.VITE_API_BASE || '/api';
}

export async function fetchConfig(tenantId, apiBase) {
  const res = await fetch(`${getBase(apiBase)}/v1/config/${tenantId}`);
  if (!res.ok) throw new Error('Failed to fetch config');
  return res.json();
}

export async function submitAudit(data, tenantId, apiBase) {
  const res = await fetch(`${getBase(apiBase)}/v1/audit/${tenantId}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to submit audit');
  return res.json();
}

export async function trackEvent(event, apiBase) {
  try {
    await fetch(`${getBase(apiBase)}/v1/analytics`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(event),
    });
  } catch (_) {
    // Analytics failures are non-fatal
  }
}
