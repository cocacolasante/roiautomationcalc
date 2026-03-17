import { useState } from 'react';

const PRIORITY_STYLE = {
  high: { color: '#dc2626', bg: '#fef2f2', label: 'High Impact' },
  medium: { color: '#d97706', bg: '#fffbeb', label: 'Medium Impact' },
  low: { color: '#6b7280', bg: '#f9fafb', label: 'Low Impact' },
};

function FindingItem({ finding, defaultOpen }) {
  const [open, setOpen] = useState(defaultOpen || false);
  // Go API uses `impact` field; client-side AI analysis may use `priority`
  const priority = finding.priority || finding.impact || 'medium';
  const p = PRIORITY_STYLE[priority] || PRIORITY_STYLE.medium;

  return (
    <div style={{ border: '1px solid #e5e7eb', borderRadius: 10, overflow: 'hidden', marginBottom: '0.625rem' }}>
      <button
        onClick={() => setOpen((o) => !o)}
        style={{
          width: '100%', textAlign: 'left', padding: '1rem 1.25rem',
          background: open ? '#fafafa' : 'white',
          border: 'none', cursor: 'pointer',
          display: 'flex', alignItems: 'center', gap: '0.75rem',
        }}
      >
        <span style={{ fontSize: '1.25rem', flexShrink: 0 }}>{finding.icon || '🔍'}</span>
        <div style={{ flex: 1 }}>
          <div style={{ fontWeight: 600, fontSize: '0.9375rem', color: '#1a1a1a' }}>{finding.title}</div>
          <span style={{ fontSize: '0.75rem', color: p.color, background: p.bg, padding: '0.125rem 0.5rem', borderRadius: 999 }}>{p.label}</span>
        </div>
        <span style={{ color: '#aaa', fontSize: '1.125rem', transform: open ? 'rotate(180deg)' : 'none', transition: 'transform 0.2s' }}>▾</span>
      </button>
      {open && (
        <div style={{ padding: '0 1.25rem 1.25rem', borderTop: '1px solid #f0f0f0' }}>
          <p style={{ color: '#555', fontSize: '0.9375rem', lineHeight: 1.6, marginTop: '0.875rem', marginBottom: '0.5rem' }}>
            {finding.description}
          </p>
          {(finding.estimatedHours || finding.hoursPerWeek) > 0 && (
            <div style={{ fontSize: '0.875rem', color: '#2563eb', fontWeight: 600 }}>
              ⏱ ~{Math.round(finding.estimatedHours || finding.hoursPerWeek)} hrs/week recoverable
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default function FindingsAccordion({ findings }) {
  if (!findings || findings.length === 0) return null;

  return (
    <div style={{ padding: '1.5rem 1.25rem' }}>
      <h3 style={{ fontSize: '1.125rem', fontWeight: 700, marginBottom: '1rem' }}>Key Findings</h3>
      {findings.map((f, i) => (
        <FindingItem key={i} finding={f} defaultOpen={i === 0} />
      ))}
    </div>
  );
}
