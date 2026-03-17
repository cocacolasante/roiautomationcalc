export default function ROIChart({ categories }) {
  if (!categories || categories.length === 0) return null;

  // Support both client-side (annualValue) and Go API (recoverableCost) field names
  const getValue = (cat) => cat.annualValue ?? cat.recoverableCost ?? 0;
  const sorted = [...categories].sort((a, b) => getValue(b) - getValue(a)).slice(0, 6);
  const max = getValue(sorted[0]) || 1;

  const CATEGORY_ICONS = {
    email: '📧', leads: '🎯', billing: '💰', support: '🎧',
    scheduling: '📅', reporting: '📊', social: '📱', general: '⚙️',
  };

  return (
    <div style={{ padding: '1.5rem 1.25rem' }}>
      <h3 style={{ fontSize: '1.125rem', fontWeight: 700, marginBottom: '1.25rem' }}>Savings by Category</h3>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
        {sorted.map((cat) => {
          const val = getValue(cat);
          const pct = Math.round((val / max) * 100);
          const icon = CATEGORY_ICONS[cat.category] || '⚙️';
          return (
            <div key={cat.category}>
              <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.875rem', marginBottom: '0.25rem' }}>
                <span style={{ fontWeight: 500 }}>{icon} {cat.label || cat.category}</span>
                <span style={{ color: '#555', fontWeight: 600 }}>${Math.round(val).toLocaleString()}/yr</span>
              </div>
              <div style={{ height: 10, background: '#f0f0f0', borderRadius: 5, overflow: 'hidden' }}>
                <div style={{
                  height: '100%',
                  width: `${pct}%`,
                  background: 'var(--primary, #2563eb)',
                  borderRadius: 5,
                  transition: 'width 0.8s ease',
                }} />
              </div>
              <div style={{ fontSize: '0.75rem', color: '#999', marginTop: 2 }}>
                {Math.round(cat.hoursPerWeek || 0)} hrs/wk
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
