import { useEffect, useRef, useState } from 'react';

function AnimatedNumber({ target, duration = 1200 }) {
  const [display, setDisplay] = useState(0);
  const raf = useRef(null);

  useEffect(() => {
    const start = Date.now();
    const from = 0;
    const tick = () => {
      const elapsed = Date.now() - start;
      const t = Math.min(elapsed / duration, 1);
      const eased = 1 - Math.pow(1 - t, 3);
      setDisplay(Math.round(from + (target - from) * eased));
      if (t < 1) raf.current = requestAnimationFrame(tick);
    };
    raf.current = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(raf.current);
  }, [target, duration]);

  return <>${display.toLocaleString()}</>;
}

export default function ROISummaryHero({ roi, audit }) {
  if (!roi) return null;

  // Support both client-side calc field names and Go API field names
  const annualValue = roi.annualValue ?? roi.recoverableCostPerYear ?? 0;
  const hoursPerWeek = roi.hoursPerWeek ?? roi.totalHoursPerWeek ?? 0;
  const leadScore = audit?.leadScore ?? roi.leadScore ?? 0;
  const leadTier = audit?.leadQuality ?? audit?.leadTier ?? roi.leadTier ?? '';

  const tierColor = { hot: '#ef4444', warm: '#f59e0b', cold: '#6b7280' }[leadTier] || '#6b7280';
  const tierLabel = { hot: 'High Priority', warm: 'Good Fit', cold: 'Exploring' }[leadTier] || leadTier;

  return (
    <div style={{ textAlign: 'center', padding: '2.5rem 1rem 2rem', borderBottom: '1px solid #f0f0f0' }}>
      <div style={{ fontSize: '0.875rem', color: '#888', marginBottom: '0.5rem', textTransform: 'uppercase', letterSpacing: '0.05em' }}>
        Your Estimated Annual Savings
      </div>
      <div style={{ fontSize: 'clamp(2.5rem, 8vw, 4rem)', fontWeight: 800, color: 'var(--primary, #2563eb)', lineHeight: 1.1, marginBottom: '0.5rem' }}>
        <AnimatedNumber target={Math.round(annualValue)} />
      </div>
      <div style={{ color: '#555', fontSize: '1rem', marginBottom: '1.75rem' }}>
        by automating <strong>{Math.round(hoursPerWeek)} hrs/week</strong> of manual work
      </div>

      <div style={{ display: 'flex', gap: '1rem', justifyContent: 'center', flexWrap: 'wrap' }}>
        <div style={{ background: '#f8fafc', borderRadius: 12, padding: '0.75rem 1.25rem', minWidth: 120 }}>
          <div style={{ fontSize: '1.375rem', fontWeight: 700, color: '#1a1a1a' }}>{Math.round(hoursPerWeek || 0)}</div>
          <div style={{ fontSize: '0.8125rem', color: '#888' }}>Hours/Week</div>
        </div>
        <div style={{ background: '#f8fafc', borderRadius: 12, padding: '0.75rem 1.25rem', minWidth: 120 }}>
          <div style={{ fontSize: '1.375rem', fontWeight: 700, color: '#1a1a1a' }}>{Math.round((hoursPerWeek || 0) * 52)}</div>
          <div style={{ fontSize: '0.8125rem', color: '#888' }}>Hours/Year</div>
        </div>
        <div style={{ background: '#f8fafc', borderRadius: 12, padding: '0.75rem 1.25rem', minWidth: 120 }}>
          <div style={{ fontSize: '1.375rem', fontWeight: 700, color: tierColor }}>{leadScore || 0}</div>
          <div style={{ fontSize: '0.8125rem', color: '#888' }}>ROI Score</div>
        </div>
      </div>

      {leadTier && (
        <div style={{ marginTop: '1rem' }}>
          <span style={{ background: tierColor + '20', color: tierColor, fontSize: '0.8125rem', fontWeight: 600, padding: '0.25rem 0.75rem', borderRadius: 999 }}>
            {tierLabel}
          </span>
        </div>
      )}
    </div>
  );
}
