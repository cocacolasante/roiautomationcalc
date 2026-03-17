export default function ProgressBar({ progress }) {
  return (
    <div style={{ height: 4, background: '#e5e7eb' }}>
      <div style={{ height: '100%', background: 'var(--primary)', width: `${progress}%`, transition: 'width 0.4s ease' }} />
    </div>
  );
}
