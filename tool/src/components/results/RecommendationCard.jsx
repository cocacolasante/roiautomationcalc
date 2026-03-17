export default function RecommendationCard({ rec, config }) {
  // Normalize between client-side shape and Go API shape (rec.product.name etc.)
  const product = rec.product || {};
  const title = rec.title || product.name || 'Recommended Solution';
  const description = rec.description || rec.specificInsight || product.description || '';
  const icon = rec.icon || product.icon || '🤖';
  const category = rec.category || product.category || '';
  const estimatedROI = rec.estimatedROI ?? rec.estimatedMonthlySaved != null ? (rec.estimatedMonthlySaved || 0) * 12 : null;
  const timeToImplement = rec.timeToImplement || (rec.paybackPeriodMonths ? `~${rec.paybackPeriodMonths} months` : null);
  const ctaUrl = rec.ctaUrl || config?.ctaUrl || '#';
  const ctaText = rec.ctaText || config?.ctaButtonText || 'Learn More';

  return (
    <div style={{
      border: '1.5px solid #e5e7eb',
      borderRadius: 12,
      overflow: 'hidden',
      marginBottom: '0.875rem',
      background: 'white',
    }}>
      <div style={{ padding: '1.125rem 1.25rem' }}>
        <div style={{ display: 'flex', alignItems: 'flex-start', gap: '0.875rem' }}>
          <span style={{ fontSize: '1.75rem', flexShrink: 0 }}>{icon}</span>
          <div style={{ flex: 1 }}>
            <div style={{ fontWeight: 700, fontSize: '1rem', color: '#1a1a1a', marginBottom: '0.25rem' }}>
              {title}
            </div>
            {category && (
              <span style={{ fontSize: '0.75rem', color: '#6b7280', background: '#f3f4f6', padding: '0.125rem 0.5rem', borderRadius: 999 }}>
                {category}
              </span>
            )}
          </div>
          {estimatedROI > 0 && (
            <div style={{ textAlign: 'right', flexShrink: 0 }}>
              <div style={{ fontSize: '1.125rem', fontWeight: 700, color: 'var(--primary, #2563eb)' }}>
                ${Math.round(estimatedROI).toLocaleString()}
              </div>
              <div style={{ fontSize: '0.75rem', color: '#888' }}>est. annual</div>
            </div>
          )}
        </div>

        {description && (
          <p style={{ color: '#555', fontSize: '0.9rem', lineHeight: 1.6, marginTop: '0.75rem' }}>
            {description}
          </p>
        )}

        {timeToImplement && (
          <div style={{ fontSize: '0.8125rem', color: '#888', marginTop: '0.5rem' }}>
            ⏰ {timeToImplement} to implement
          </div>
        )}
      </div>

      <div style={{ borderTop: '1px solid #f0f0f0', padding: '0.875rem 1.25rem', background: '#fafafa' }}>
        <a
          href={ctaUrl}
          target="_blank"
          rel="noopener noreferrer"
          style={{
            display: 'inline-block',
            background: 'var(--primary, #2563eb)',
            color: 'white',
            padding: '0.5rem 1.25rem',
            borderRadius: 7,
            fontSize: '0.875rem',
            fontWeight: 600,
            textDecoration: 'none',
          }}
        >
          {ctaText} →
        </a>
      </div>
    </div>
  );
}
