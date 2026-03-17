import SliderInput from './SliderInput';

export default function WorkflowCard({ question, value, onChange }) {
  return (
    <div className="question-card">
      <div className="question-label">{question.question}</div>
      {question.hint && <div className="question-hint">{question.hint}</div>}

      {question.inputType === 'hours_slider' && (
        <SliderInput question={question} value={value} onChange={onChange} />
      )}

      {(question.inputType === 'number' || question.inputType === 'people_count') && (
        <input
          type="number"
          min={question.minValue ?? 0}
          max={question.maxValue}
          value={value ?? ''}
          onChange={(e) => onChange(parseFloat(e.target.value) || 0)}
          style={{ width: '100%', padding: '0.75rem', border: '2px solid #e5e7eb', borderRadius: 8, fontSize: '1rem' }}
          placeholder={`Enter ${question.unit || 'number'}`}
        />
      )}

      {question.inputType === 'yes_no' && (
        <div className="radio-group">
          {['yes', 'no'].map((opt) => (
            <label key={opt} className={`radio-option${value === opt ? ' selected' : ''}`}>
              <input type="radio" name={question.key} value={opt} checked={value === opt} onChange={() => onChange(opt)} />
              {opt.charAt(0).toUpperCase() + opt.slice(1)}
            </label>
          ))}
        </div>
      )}

      {question.inputType === 'multiple_choice' && (
        <div className="radio-group">
          {(question.options || []).map((opt) => (
            <label key={opt.value} className={`radio-option${value === opt.value ? ' selected' : ''}`}>
              <input type="radio" name={question.key} value={opt.value} checked={value === opt.value} onChange={() => onChange(opt.value)} />
              {opt.label}
            </label>
          ))}
        </div>
      )}
    </div>
  );
}
