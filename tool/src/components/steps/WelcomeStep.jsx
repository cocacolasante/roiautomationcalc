export default function WelcomeStep({ config, onNext }) {
  return (
    <div style={{ textAlign: 'center', paddingTop: '2rem' }}>
      {config.logoUrl && (
        <img src={config.logoUrl} alt={config.companyName} style={{ height: 56, marginBottom: '2rem' }} />
      )}
      <h1 style={{ fontSize: 'clamp(1.75rem, 4vw, 2.5rem)', fontWeight: 800, color: '#1a1a1a', marginBottom: '1rem' }}>
        {config.toolTitle || 'Free Automation Audit'}
      </h1>
      {config.toolSubtitle && (
        <p style={{ fontSize: '1.1875rem', color: '#555', maxWidth: 540, margin: '0 auto 1.5rem', lineHeight: 1.6 }}>
          {config.toolSubtitle}
        </p>
      )}
      {config.welcomeMessage && (
        <p style={{ color: '#888', maxWidth: 500, margin: '0 auto 2.5rem', lineHeight: 1.6 }}>
          {config.welcomeMessage}
        </p>
      )}
      <div style={{ display: 'flex', gap: '1.5rem', justifyContent: 'center', flexWrap: 'wrap', marginBottom: '2.5rem', color: '#888', fontSize: '0.875rem' }}>
        <span>⏱ Takes ~5 minutes</span>
        <span>🔒 100% confidential</span>
        <span>📊 Personalized report</span>
      </div>
      <button className="btn-primary" onClick={onNext} style={{ fontSize: '1.125rem', padding: '1rem 3rem' }}>
        Start My Free Audit →
      </button>
    </div>
  );
}
