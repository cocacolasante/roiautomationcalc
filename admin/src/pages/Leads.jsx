import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { listLeads, updateLeadStatus } from '../api/client';

const TENANT_ID = localStorage.getItem('blueprint_tenant_id') || 'default';

export default function Leads() {
  const [leads, setLeads] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState('all');

  const load = () => {
    setLoading(true);
    const params = filter !== 'all' ? { tier: filter } : {};
    listLeads(TENANT_ID, params)
      .then((r) => setLeads(r.leads || []))
      .catch(() => setLeads([]))
      .finally(() => setLoading(false));
  };

  useEffect(load, [filter]);

  const handleStatus = async (leadId, status) => {
    await updateLeadStatus(TENANT_ID, leadId, status);
    load();
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-xl font-bold">Leads</h1>
        <div className="flex gap-2">
          {['all', 'hot', 'warm', 'cold'].map((t) => (
            <button
              key={t}
              onClick={() => setFilter(t)}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                filter === t ? 'bg-primary text-white' : 'bg-white border border-gray-200 text-gray-600 hover:bg-gray-50'
              }`}
            >
              {t.charAt(0).toUpperCase() + t.slice(1)}
            </button>
          ))}
        </div>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        {loading ? (
          <div className="p-5 text-gray-400 text-sm">Loading...</div>
        ) : leads.length === 0 ? (
          <div className="p-5 text-gray-400 text-sm">No leads found.</div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="text-gray-400 text-xs border-b border-gray-100 bg-gray-50">
                <th className="text-left px-5 py-3 font-medium">Contact</th>
                <th className="text-left px-5 py-3 font-medium">Business</th>
                <th className="text-left px-5 py-3 font-medium">Score</th>
                <th className="text-left px-5 py-3 font-medium">Annual Value</th>
                <th className="text-left px-5 py-3 font-medium">Status</th>
                <th className="text-left px-5 py-3 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {leads.map((lead) => (
                <tr key={lead.id} className="border-b border-gray-50 hover:bg-gray-50">
                  <td className="px-5 py-3">
                    <Link to={`/leads/${lead.id}`} className="font-medium text-primary hover:underline">
                      {lead.contactName}
                    </Link>
                    <div className="text-xs text-gray-400">{lead.contactEmail}</div>
                  </td>
                  <td className="px-5 py-3 text-gray-600">{lead.businessName || '—'}</td>
                  <td className="px-5 py-3">
                    <span className={`font-bold ${lead.leadScore >= 70 ? 'text-red-500' : lead.leadScore >= 40 ? 'text-yellow-600' : 'text-gray-500'}`}>
                      {lead.leadScore}
                    </span>
                    <span className={`ml-2 px-1.5 py-0.5 rounded-full text-xs ${
                      lead.leadTier === 'hot' ? 'bg-red-50 text-red-500' :
                      lead.leadTier === 'warm' ? 'bg-yellow-50 text-yellow-600' :
                      'bg-gray-100 text-gray-500'
                    }`}>{lead.leadTier}</span>
                  </td>
                  <td className="px-5 py-3 font-semibold text-green-600">${Math.round(lead.roiAnnualValue || 0).toLocaleString()}</td>
                  <td className="px-5 py-3">
                    <select
                      value={lead.status || 'new'}
                      onChange={(e) => handleStatus(lead.id, e.target.value)}
                      className="text-xs border border-gray-200 rounded px-2 py-1 bg-white"
                    >
                      {['new', 'contacted', 'qualified', 'proposal', 'closed_won', 'closed_lost'].map((s) => (
                        <option key={s} value={s}>{s.replace(/_/g, ' ')}</option>
                      ))}
                    </select>
                  </td>
                  <td className="px-5 py-3">
                    <Link to={`/leads/${lead.id}`} className="text-primary text-xs hover:underline">View →</Link>
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
