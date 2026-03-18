import { useState, useEffect } from 'react';
import { fetchConfig, submitAudit } from './lib/api';
import { calculateROI } from './lib/calculations';
import { applyBranding, getTenantIdFromEmbed } from './lib/branding';
import AuditWizard from './components/AuditWizard';
import ResultsPage from './components/results/ResultsPage';

function LoadingScreen() {
  return (
    <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', flexDirection: 'column', gap: '1rem' }}>
      <div style={{ width: 48, height: 48, border: '4px solid #e5e7eb', borderTopColor: 'var(--primary)', borderRadius: '50%', animation: 'spin 0.8s linear infinite' }} />
      <p style={{ color: '#888' }}>Loading your audit tool...</p>
    </div>
  );
}

export default function App({ tenantId: propTenantId, apiBase, embedded = false }) {
  const [step, setStep] = useState(0);
  const [answers, setAnswers] = useState({});
  const [config, setConfig] = useState(null);
  const [audit, setAudit] = useState(null);
  const [liveROI, setLiveROI] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const tid = propTenantId || getTenantIdFromEmbed();
    fetchConfig(tid, apiBase)
      .then((cfg) => { setConfig(cfg); applyBranding(cfg); })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [propTenantId, apiBase]);

  useEffect(() => {
    if (config && Object.keys(answers).length > 0) {
      setLiveROI(calculateROI(answers, config));
    }
  }, [answers, config]);

  if (loading) return <LoadingScreen />;
  if (error) return <div style={{ padding: '2rem', color: 'red' }}>Error: {error}</div>;
  if (!config) return <LoadingScreen />;

  if (audit) return <ResultsPage audit={audit} config={config} liveROI={liveROI} />;

  return (
    <div style={{ '--primary': config.primaryColor || '#6C63FF', '--accent': config.accentColor || '#F18F01' }}>
      <AuditWizard
        config={config}
        step={step}
        answers={answers}
        liveROI={liveROI}
        onAnswer={(key, value) => setAnswers((prev) => ({ ...prev, [key]: value }))}
        onNext={() => setStep((s) => s + 1)}
        onBack={() => setStep((s) => s - 1)}
        onComplete={async (contactInfo) => {
          const result = await submitAudit({ ...answers, ...contactInfo }, config.tenantId, apiBase);
          setAudit(result);
        }}
      />
    </div>
  );
}
