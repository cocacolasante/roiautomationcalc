import WorkflowCard from '../ui/WorkflowCard';

const CATEGORY_LABELS = {
  email: '📧 Email & Communication',
  leads: '🎯 Lead Management',
  billing: '💰 Billing & Invoicing',
  support: '🎧 Customer Support',
  scheduling: '📅 Scheduling',
  reporting: '📊 Reporting',
  social: '📱 Social Media',
  general: '⚙️ General Operations',
};

export default function WorkflowStep({ questions, answers, categoryLabel, onAnswer, onNext, onBack }) {
  if (!questions || questions.length === 0) {
    return (
      <div>
        <p>No questions in this section.</p>
        <div className="step-actions">
          {onBack && <button className="btn-secondary" onClick={onBack}>Back</button>}
          <button className="btn-primary" onClick={onNext}>Continue</button>
        </div>
      </div>
    );
  }

  const category = questions[0]?.category || 'general';
  const label = CATEGORY_LABELS[category] || categoryLabel || category;
  const requiredAnswered = questions
    .filter((q) => q.isRequired)
    .every((q) => answers[q.key] !== undefined && answers[q.key] !== '');

  return (
    <div>
      <h2 style={{ fontSize: '1.5rem', fontWeight: 700, marginBottom: '0.5rem' }}>{label}</h2>
      <p style={{ color: '#888', marginBottom: '2rem' }}>Answer as accurately as you can — estimates are fine.</p>
      <div>
        {questions.map((q) => (
          <WorkflowCard key={q.key} question={q} value={answers[q.key]} onChange={(val) => onAnswer(q.key, val)} />
        ))}
      </div>
      <div className="step-actions">
        {onBack && <button className="btn-secondary" onClick={onBack}>Back</button>}
        <button className="btn-primary" onClick={onNext} disabled={!requiredAnswered}>Continue</button>
      </div>
    </div>
  );
}
