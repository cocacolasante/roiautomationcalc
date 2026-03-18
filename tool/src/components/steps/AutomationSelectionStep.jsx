import { useState } from 'react';

const OPTIONS = [
  { key: 'email',      icon: '📧', label: 'Email & Communication',     description: 'Sorting, drafting replies, follow-up sequences',        defaultHours: 5  },
  { key: 'leads',      icon: '🎯', label: 'Lead Management',            description: 'Lead capture, CRM entry, nurture campaigns',            defaultHours: 4  },
  { key: 'billing',    icon: '💰', label: 'Billing & Invoicing',        description: 'Invoice generation, payment reminders, collections',    defaultHours: 3  },
  { key: 'support',    icon: '🎧', label: 'Customer Support',           description: 'Ticket routing, FAQ responses, status updates',         defaultHours: 5  },
  { key: 'scheduling', icon: '📅', label: 'Scheduling & Appointments',  description: 'Booking, reminders, rescheduling, confirmations',       defaultHours: 3  },
  { key: 'reporting',  icon: '📊', label: 'Reporting & Analytics',      description: 'Weekly reports, KPI dashboards, data exports',          defaultHours: 3  },
  { key: 'social',     icon: '📱', label: 'Social Media',               description: 'Content scheduling, replies, engagement tracking',      defaultHours: 4  },
  { key: 'data_entry', icon: '⚙️', label: 'Data Entry & Operations',   description: 'Manual data transfers, spreadsheet updates, file mgmt', defaultHours: 5  },
];

export default function AutomationSelectionStep({ answers, onAnswer, onNext, onBack }) {
  const [error, setError] = useState('');

  const selectedAreas = answers.automation_areas || [];

  const toggleArea = (key) => {
    const next = selectedAreas.includes(key)
      ? selectedAreas.filter((k) => k !== key)
      : [...selectedAreas, key];
    onAnswer('automation_areas', next);

    if (!selectedAreas.includes(key) && !answers[`automation_hours_${key}`]) {
      const opt = OPTIONS.find((o) => o.key === key);
      onAnswer(`automation_hours_${key}`, opt.defaultHours);
    }
  };

  const handleNext = () => {
    if (selectedAreas.length === 0) {
      setError("Please select at least one area you'd like to automate.");
      return;
    }
    setError('');
    onNext();
  };

  return (
    <div>
      <h2 style={{ fontSize: '1.5rem', fontWeight: 700, marginBottom: '0.5rem' }}>What would you like to automate?</h2>
      <p style={{ color: '#888', marginBottom: '1.75rem' }}>
        Select every area that applies. For each one, adjust how many hours per week your team currently spends on it.
      </p>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(260px, 1fr))', gap: '0.75rem', marginBottom: '1.5rem' }}>
        {OPTIONS.map((opt) => {
          const selected = selectedAreas.includes(opt.key);
          const hours = answers[`automation_hours_${opt.key}`] ?? opt.defaultHours;

          return (
            <div
              key={opt.key}
              className={`radio-option${selected ? ' selected' : ''}`}
              style={{ flexDirection: 'column', alignItems: 'flex-start', gap: '0.5rem', cursor: 'pointer', padding: '1rem' }}
              onClick={() => toggleArea(opt.key)}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.625rem', width: '100%' }}>
                <span style={{ fontSize: '1.375rem', lineHeight: 1 }}>{opt.icon}</span>
                <span style={{ fontWeight: 600, fontSize: '0.9375rem', flex: 1 }}>{opt.label}</span>
                <input type="checkbox" checked={selected} readOnly style={{ width: 18, height: 18, flexShrink: 0 }} onClick={(e) => e.stopPropagation()} onChange={() => toggleArea(opt.key)} />
              </div>

              <div style={{ fontSize: '0.8125rem', color: '#888', paddingLeft: '2rem' }}>{opt.description}</div>

              {selected && (
                <div
                  style={{ width: '100%', paddingTop: '0.625rem', marginTop: '0.25rem', borderTop: '1px solid #e5e7eb' }}
                  onClick={(e) => e.stopPropagation()}
                >
                  <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.8125rem', color: '#555', marginBottom: '0.375rem' }}>
                    <span>Hours/week on this:</span>
                    <strong style={{ color: 'var(--primary)' }}>{hours}h</strong>
                  </div>
                  <input
                    type="range"
                    min={1}
                    max={20}
                    value={hours}
                    onChange={(e) => onAnswer(`automation_hours_${opt.key}`, parseInt(e.target.value))}
                    className="slider-input"
                    style={{ width: '100%' }}
                  />
                  <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.75rem', color: '#bbb', marginTop: '0.25rem' }}>
                    <span>1h</span><span>20h</span>
                  </div>
                </div>
              )}
            </div>
          );
        })}
      </div>

      <div className="question-card">
        <label className="question-label">
          Anything else you'd like to automate?{' '}
          <span style={{ fontWeight: 400, color: '#aaa' }}>(optional)</span>
        </label>
        <textarea
          value={answers.automation_notes || ''}
          onChange={(e) => onAnswer('automation_notes', e.target.value)}
          placeholder="e.g. We manually move orders from Shopify into a spreadsheet every morning, then email each supplier…"
          rows={3}
          style={{ width: '100%', padding: '0.75rem', border: '2px solid #e5e7eb', borderRadius: 8, fontSize: '1rem', resize: 'vertical', outline: 'none' }}
          onFocus={(e) => (e.target.style.borderColor = 'var(--primary)')}
          onBlur={(e) => (e.target.style.borderColor = '#e5e7eb')}
        />
      </div>

      {error && <div style={{ color: '#ef4444', fontSize: '0.875rem', marginBottom: '1rem' }}>{error}</div>}

      <div className="step-actions">
        {onBack && <button className="btn-secondary" onClick={onBack}>Back</button>}
        <button className="btn-primary" onClick={handleNext}>Continue</button>
      </div>
    </div>
  );
}
