export default function CTASection({ config, roi }) {
  const ctaUrl = config?.ctaUrl || '#';
  const ctaText = config?.ctaButtonText || 'Book a Free Strategy Call';
  const annualValue = roi?.annualValue ?? roi?.recoverableCostPerYear ?? 0;

  return (
    <div style={{
      background: 'linear-gradient(135deg, var(--primary, #2563eb) 0%, #1d4ed8 100%)',
      borderRadius: 14,
      padding: '2rem 1.5rem',
      textAlign: 'center',
      color: 'white',
      margin: '1.5rem 0',
    }}>
      <div style={{ fontSize: '1.5rem', marginBottom: '0.75rem' }}>🚀</div>
      <h3 style={{ fontSize: '1.25rem', fontWeight: 800, marginBottom: '0.625rem' }}>
        Ready to Capture Your ${Math.round(annualValue).toLocaleString()} Opportunity?
      </h3>
      <p style={{ opacity: 0.85, fontSize: '0.9375rem', lineHeight: 1.6, maxWidth: 440, margin: '0 auto 1.5rem' }}>
        {config?.ctaDescription || "Let's build a custom automation roadmap for your business. No tech jargon, no pressure — just a clear plan to save time and grow faster."}
      </p>
      <a
        href={ctaUrl}
        target="_blank"
        rel="noopener noreferrer"
        style={{
          display: 'inline-block',
          background: 'white',
          color: 'var(--primary, #2563eb)',
          fontWeight: 700,
          fontSize: '1rem',
          padding: '0.875rem 2.25rem',
          borderRadius: 9,
          textDecoration: 'none',
          boxShadow: '0 2px 12px rgba(0,0,0,0.15)',
        }}
      >
        {ctaText} →
      </a>
      <div style={{ opacity: 0.7, fontSize: '0.8125rem', marginTop: '1rem' }}>
        Free consultation · No commitment required
      </div>
    </div>
  );
}
