export default function SliderInput({ question, value, onChange }) {
  const min = question.minValue ?? 0;
  const max = question.maxValue ?? 40;
  const current = value ?? question.defaultValue ?? 0;

  return (
    <div className="slider-container">
      <input
        type="range"
        className="slider-input"
        min={min}
        max={max}
        step={0.5}
        value={current}
        onChange={(e) => onChange(parseFloat(e.target.value))}
      />
      <div className="slider-value">
        {current} <span style={{ fontSize: '0.75rem', fontWeight: 400, color: '#888' }}>{question.unit || 'hrs/wk'}</span>
      </div>
    </div>
  );
}
