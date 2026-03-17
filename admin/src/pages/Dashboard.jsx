import { useEffect, useState } from 'react';
import { listLeads, getAnalytics } from '../api/client';

function StatCard({ label, value, sub, color = 'text-primary' }) {
  return (
    <div className="bg-white rounded-xl border border-gray-200 p-5">
      <div className="text-sm text-gray-500 mb-1">{label}</div>
      <div className={`text-2xl font-bold ${color}`}>{value}</div>
      {sub && <div className="text-xs text-gray-400 mt-1">{sub}</div>}
    </div>
  );
}

export default function Dashboard() {
  const tenantId = localStorage.getItem('blueprint_tenant_id');
  const [leads, setLeads] = useState([]);
  const [analytics, setAnalytics] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!tenantId) { setLoading(false); return; }
    Promise.all([
      listLeads(tenantId, { limit: 5 }).catch(() => ({ leads: [] })),
      getAnalytics(tenantId).catch(() => null),
    ]).then(([l, a]) => {
      setLeads(l.leads || []);
      setAnalytics(a);
    }).finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="text-gray-400 p-4">Loading...</div>;

  const hot = leads.filter((l) => l.leadTier === 'hot').length;
  const totalValue = leads.reduce((s, l) => s + (l.roiAnnualValue || 0), 0);

  return (
    <div>
      <h1 className="text-xl font-bold mb-6">Dashboard</h1>

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <StatCard label="Total Leads" value={analytics?.totalLeads ?? leads.length} />
        <StatCard label="Hot Leads" value={hot} color="text-red-500" />
        <StatCard label="Avg ROI Score" value={analytics?.avgLeadScore ? Math.round(analytics.avgLeadScore) : '—'} />
        <StatCard label="Total Pipeline Value" value={`$${Math.round(totalValue / 1000)}k`} color="text-green-600" />
      </div>

      <div className="bg-white rounded-xl border border-gray-200">
        <div className="px-5 py-4 border-b border-gray-100 font-semibold text-sm">Recent Leads</div>
        {leads.length === 0 ? (
          <div className="p-5 text-gray-400 text-sm">No leads yet.</div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="text-gray-400 text-xs border-b border-gray-100">
                <th className="text-left px-5 py-3 font-medium">Name</th>
                <th className="text-left px-5 py-3 font-medium">Email</th>
                <th className="text-left px-5 py-3 font-medium">Score</th>
                <th className="text-left px-5 py-3 font-medium">ROI</th>
                <th className="text-left px-5 py-3 font-medium">Tier</th>
              </tr>
            </thead>
            <tbody>
              {leads.map((lead) => (
                <tr key={lead.id} className="border-b border-gray-50 hover:bg-gray-50">
                  <td className="px-5 py-3 font-medium">{lead.contactName}</td>
                  <td className="px-5 py-3 text-gray-500">{lead.contactEmail}</td>
                  <td className="px-5 py-3">{lead.leadScore}</td>
                  <td className="px-5 py-3 text-green-600 font-semibold">${Math.round(lead.roiAnnualValue || 0).toLocaleString()}</td>
                  <td className="px-5 py-3">
                    <span className={`px-2 py-0.5 rounded-full text-xs font-semibold ${
                      lead.leadTier === 'hot' ? 'bg-red-50 text-red-500' :
                      lead.leadTier === 'warm' ? 'bg-yellow-50 text-yellow-600' :
                      'bg-gray-100 text-gray-500'
                    }`}>{lead.leadTier}</span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
