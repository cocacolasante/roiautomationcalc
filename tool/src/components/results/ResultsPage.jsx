import ROISummaryHero from './ROISummaryHero';
import ROIChart from './ROIChart';
import FindingsAccordion from './FindingsAccordion';
import RecommendationCard from './RecommendationCard';
import CTASection from './CTASection';

export default function ResultsPage({ audit, config, liveROI }) {
  if (!audit) return null;

  // Go API serializes ROICalculation as `roiSummary`.
  // If the server returned an empty calculation (all zeros — happens when no DB questions are configured),
  // fall back to the client-side liveROI calculated from automation area selections.
  const serverRoi = audit.roiSummary || audit.roiCalculation;
  const hasServerRoi = (serverRoi?.recoverableCostPerYear ?? serverRoi?.annualValue ?? 0) > 0;
  const roi = hasServerRoi ? serverRoi : liveROI;
  const { findings, recommendations, executiveSummary } = audit;

  return (
    <div>
      <ROISummaryHero roi={roi} audit={audit} />

      {executiveSummary && (
        <div style={{ padding: '1.5rem 1.25rem', borderBottom: '1px solid #f0f0f0' }}>
          <h3 style={{ fontSize: '1.125rem', fontWeight: 700, marginBottom: '0.75rem' }}>Executive Summary</h3>
          <p style={{ color: '#555', lineHeight: 1.7, fontSize: '0.9375rem' }}>{executiveSummary}</p>
        </div>
      )}

      {roi?.categories?.length > 0 && (
        <div style={{ borderBottom: '1px solid #f0f0f0' }}>
          <ROIChart categories={roi.categories} />
        </div>
      )}

      {findings?.length > 0 && (
        <div style={{ borderBottom: '1px solid #f0f0f0' }}>
          <FindingsAccordion findings={findings} />
        </div>
      )}

      {recommendations?.length > 0 && (
        <div style={{ padding: '1.5rem 1.25rem', borderBottom: '1px solid #f0f0f0' }}>
          <h3 style={{ fontSize: '1.125rem', fontWeight: 700, marginBottom: '1rem' }}>Recommended Solutions</h3>
          {recommendations.map((rec, i) => (
            <RecommendationCard key={i} rec={rec} config={config} />
          ))}
        </div>
      )}

      <div style={{ padding: '0 1.25rem 1.5rem' }}>
        <CTASection config={config} roi={roi} />
      </div>

      <div style={{ padding: '1rem 1.25rem', textAlign: 'center', fontSize: '0.8rem', color: '#bbb', borderTop: '1px solid #f0f0f0' }}>
        Powered by Blueprint Automation · Results are estimates based on industry averages
      </div>
    </div>
  );
}
