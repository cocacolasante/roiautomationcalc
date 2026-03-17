import { useState } from 'react';

const industries = [
  'Technology', 'Marketing/Agency', 'Legal', 'Medical/Healthcare', 'Accounting/Finance',
  'Real Estate', 'Consulting', 'E-commerce/Retail', 'Construction/Trades', 'Other',
];

export default function BusinessProfileStep({ answers, onAnswer, onNext, onBack }) {
  const [errors, setErrors] = useState({});

  const handleNext = () => {
    const errs = {};
    if (!answers.business_name) errs.business_name = 'Required';
    if (!answers.team_size) errs.team_size = 'Required';
    if (Object.keys(errs).length > 0) { setErrors(errs); return; }
    onNext();
  };

  const fieldStyle = (key) => ({
    width: '100%', padding: '0.75rem',
    border: `2px solid ${errors[key] ? '#ef4444' : '#e5e7eb'}`,
    borderRadius: 8, fontSize: '1rem',
  });

  return (
    <div>
      <h2 style={{ fontSize: '1.5rem', fontWeight: 700, marginBottom: '0.5rem' }}>Tell us about your business</h2>
      <p style={{ color: '#888', marginBottom: '2rem' }}>This helps us tailor the audit to your situation.</p>

      <div className="question-card">
        <label className="question-label">Business Name</label>
        <input type="text" value={answers.business_name || ''} onChange={(e) => onAnswer('business_name', e.target.value)} placeholder="Your company name" style={fieldStyle('business_name')} />
        {errors.business_name && <div style={{ color: '#ef4444', fontSize: '0.875rem', marginTop: 4 }}>{errors.business_name}</div>}
      </div>

      <div className="question-card">
        <label className="question-label">How many people are on your team? (including you)</label>
        <input type="number" min={1} value={answers.team_size || ''} onChange={(e) => onAnswer('team_size', parseInt(e.target.value) || '')} placeholder="e.g. 5" style={fieldStyle('team_size')} />
        {errors.team_size && <div style={{ color: '#ef4444', fontSize: '0.875rem', marginTop: 4 }}>{errors.team_size}</div>}
      </div>

      <div className="question-card">
        <label className="question-label">Industry</label>
        <select value={answers.industry || ''} onChange={(e) => onAnswer('industry', e.target.value)} style={{ ...fieldStyle('industry'), background: 'white' }}>
          <option value="">Select your industry</option>
          {industries.map((ind) => <option key={ind} value={ind.toLowerCase().replace(/[/ ]/g, '_')}>{ind}</option>)}
        </select>
      </div>

      <div className="step-actions">
        {onBack && <button className="btn-secondary" onClick={onBack}>Back</button>}
        <button className="btn-primary" onClick={handleNext}>Continue</button>
      </div>
    </div>
  );
}
