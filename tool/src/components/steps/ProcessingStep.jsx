import { useEffect, useState } from 'react';

const MESSAGES = [
  'Analyzing your workflow data...',
  'Calculating automation potential...',
  'Identifying quick wins...',
  'Matching recommended solutions...',
  'Generating your personalized audit...',
  'Almost ready...',
];

export default function ProcessingStep() {
  const [msgIndex, setMsgIndex] = useState(0);
  const [dots, setDots] = useState('');

  useEffect(() => {
    const msgTimer = setInterval(() => {
      setMsgIndex((i) => Math.min(i + 1, MESSAGES.length - 1));
    }, 1800);
    const dotTimer = setInterval(() => {
      setDots((d) => (d.length >= 3 ? '' : d + '.'));
    }, 400);
    return () => { clearInterval(msgTimer); clearInterval(dotTimer); };
  }, []);

  return (
    <div style={{ textAlign: 'center', padding: '3rem 1rem' }}>
      <div style={{ marginBottom: '2rem' }}>
        <div style={{
          width: 64, height: 64, margin: '0 auto 1.5rem',
          border: '4px solid #e5e7eb',
          borderTopColor: 'var(--primary, #2563eb)',
          borderRadius: '50%',
          animation: 'spin 0.8s linear infinite',
        }} />
        <style>{`@keyframes spin { to { transform: rotate(360deg); } }`}</style>
      </div>

      <h2 style={{ fontSize: '1.5rem', fontWeight: 700, marginBottom: '0.75rem' }}>
        Building Your Report
      </h2>
      <p style={{ color: '#888', fontSize: '1rem', minHeight: '1.5rem' }}>
        {MESSAGES[msgIndex]}{dots}
      </p>

      <div style={{ marginTop: '2.5rem', display: 'flex', gap: '0.5rem', justifyContent: 'center' }}>
        {MESSAGES.map((_, i) => (
          <div
            key={i}
            style={{
              width: i === msgIndex ? 24 : 8,
              height: 8,
              borderRadius: 4,
              background: i <= msgIndex ? 'var(--primary, #2563eb)' : '#e5e7eb',
              transition: 'all 0.3s ease',
            }}
          />
        ))}
      </div>
    </div>
  );
}
