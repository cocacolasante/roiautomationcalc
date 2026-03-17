import { useEffect, useState } from 'react';
import { getAnalytics } from '../api/client';

function Bar({ label, value, max, color = '#2563eb' }) {
  const pct = max > 0 ? Math.round((value / max) * 100) : 0;
  return (
    <div className="flex items-center gap-3 text-sm">
      <div className="w-28 text-gray-500 text-right truncate">{label}</div>
      <div className="flex-1 bg-gray-100 rounded-full h-2.5 overflow-hidden">
        <div style={{ width: `${pct}%`, background: color }} className="h-full rounded-full transition-all" />
      </div>
      <div className="w-10 font-semibold text-gray-700">{value}</div>
    </div>
  );
}

export default function Analytics() {
  const tenantId = localStorage.getItem('blueprint_tenant_id');
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!tenantId) { setLoading(false); return; }
    getAnalytics(tenantId)
      .then(setData)
      .catch(() => setData(null))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="text-gray-400 p-4">Loading...</div>;
  if (!data) return <div className="text-gray-400 p-4">No analytics data available.</div>;

  const funnel = data.funnel || {};
  const funnelSteps = [
    { label: 'Tool Loads', value: funnel.loads || 0 },
    { label: 'Started', value: funnel.started || 0 },
    { label: 'Completed', value: funnel.completed || 0 },
    { label: 'Submitted', value: data.totalLeads || 0 },
  ];
  const maxFunnel = funnelSteps[0]?.value || 1;

  return (
    <div>
      <h1 className="text-xl font-bold mb-6">Analytics</h1>

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <div className="text-sm text-gray-400 mb-1">Total Leads</div>
          <div className="text-2xl font-bold text-primary">{data.totalLeads || 0}</div>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <div className="text-sm text-gray-400 mb-1">Completion Rate</div>
          <div className="text-2xl font-bold text-green-600">
            {funnel.loads ? Math.round(((data.totalLeads || 0) / funnel.loads) * 100) : 0}%
          </div>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <div className="text-sm text-gray-400 mb-1">Avg Lead Score</div>
          <div className="text-2xl font-bold text-yellow-500">{Math.round(data.avgLeadScore || 0)}</div>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <div className="text-sm text-gray-400 mb-1">Avg Annual ROI</div>
          <div className="text-2xl font-bold">${Math.round((data.avgRoiValue || 0) / 1000)}k</div>
        </div>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 p-5 mb-4">
        <div className="font-semibold text-sm mb-4">Conversion Funnel</div>
        <div className="space-y-3">
          {funnelSteps.map((s) => (
            <Bar key={s.label} label={s.label} value={s.value} max={maxFunnel} />
          ))}
        </div>
      </div>

      {data.tierBreakdown && (
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <div className="font-semibold text-sm mb-4">Lead Tier Breakdown</div>
          <div className="space-y-3">
            <Bar label="Hot" value={data.tierBreakdown.hot || 0} max={data.totalLeads || 1} color="#ef4444" />
            <Bar label="Warm" value={data.tierBreakdown.warm || 0} max={data.totalLeads || 1} color="#f59e0b" />
            <Bar label="Cold" value={data.tierBreakdown.cold || 0} max={data.totalLeads || 1} color="#6b7280" />
          </div>
        </div>
      )}
    </div>
  );
}
