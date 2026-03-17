import { useState } from 'react';

export default function ContactStep({ answers, onAnswer, onNext, onBack, onComplete }) {
  const [errors, setErrors] = useState({});

  const validate = () => {
    const errs = {};
    if (!answers.contact_name) errs.contact_name = 'Name is required';
    if (!answers.contact_email) errs.contact_email = 'Email is required';
    else if (!/\S+@\S+\.\S+/.test(answers.contact_email)) errs.contact_email = 'Invalid email';
    return errs;
  };

  const handleSubmit = async () => {
    const errs = validate();
    if (Object.keys(errs).length > 0) { setErrors(errs); return; }
    onNext(); // Move to processing step
    try {
      await onComplete({
        contactName: answers.contact_name,
        contactEmail: answers.contact_email,
        contactPhone: answers.contact_phone || '',
        businessName: answers.business_name || '',
        teamSize: parseInt(answers.team_size) || 0,
        industry: answers.industry || '',
        answers,
      });
    } catch (e) {
      console.error('Audit submission failed:', e);
    }
  };

  const fieldStyle = (key) => ({
    width: '100%', padding: '0.75rem',
    border: `2px solid ${errors[key] ? '#ef4444' : '#e5e7eb'}`,
    borderRadius: 8, fontSize: '1rem',
  });

  return (
    <div>
      <h2 style={{ fontSize: '1.5rem', fontWeight: 700, marginBottom: '0.5rem' }}>Almost there!</h2>
      <p style={{ color: '#888', marginBottom: '2rem' }}>Where should we send your personalized audit report?</p>

      <div className="question-card">
        <label className="question-label">Your Name *</label>
        <input type="text" value={answers.contact_name || ''} onChange={(e) => onAnswer('contact_name', e.target.value)} placeholder="First Last" style={fieldStyle('contact_name')} />
        {errors.contact_name && <div style={{ color: '#ef4444', fontSize: '0.875rem', marginTop: 4 }}>{errors.contact_name}</div>}
      </div>

      <div className="question-card">
        <label className="question-label">Email Address *</label>
        <input type="email" value={answers.contact_email || ''} onChange={(e) => onAnswer('contact_email', e.target.value)} placeholder="you@company.com" style={fieldStyle('contact_email')} />
        {errors.contact_email && <div style={{ color: '#ef4444', fontSize: '0.875rem', marginTop: 4 }}>{errors.contact_email}</div>}
      </div>

      <div className="question-card">
        <label className="question-label">Phone (optional)</label>
        <input type="tel" value={answers.contact_phone || ''} onChange={(e) => onAnswer('contact_phone', e.target.value)} placeholder="+1 (555) 000-0000" style={{ width: '100%', padding: '0.75rem', border: '2px solid #e5e7eb', borderRadius: 8, fontSize: '1rem' }} />
      </div>

      <div style={{ fontSize: '0.8125rem', color: '#aaa', marginBottom: '1.5rem' }}>🔒 Your information is kept private. No spam, ever.</div>

      <div className="step-actions">
        {onBack && <button className="btn-secondary" onClick={onBack}>Back</button>}
        <button className="btn-primary" onClick={handleSubmit}>Generate My Audit Report →</button>
      </div>
    </div>
  );
}
