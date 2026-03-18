import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getLead, updateLeadStatus } from '../api/client';

const TENANT_ID = localStorage.getItem('blueprint_tenant_id') || 'default';

export default function LeadDetail() {
  const { id } = useParams();
  const [lead, setLead] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getLead(TENANT_ID, id)
      .then(setLead)
      .catch(() => setLead(null))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) return <div className="text-gray-400 p-4">Loading...</div>;
  if (!lead) return <div className="p-4 text-red-500">Lead not found.</div>;

  const { audit } = lead;
  const roi = audit?.roiSummary;

  return (
    <div className="max-w-3xl">
      <div className="flex items-center gap-3 mb-6">
        <Link to="/leads" className="text-gray-400 hover:text-gray-600 text-sm">← Leads</Link>
        <span className="text-gray-300">/</span>
        <span className="font-semibold">{lead.contactName}</span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <div className="text-sm font-semibold text-gray-500 mb-3">Contact Info</div>
          <div className="space-y-2 text-sm">
            <div><span className="text-gray-400">Name:</span> <span className="font-medium">{lead.contactName}</span></div>
            <div><span className="text-gray-400">Email:</span> <a href={`mailto:${lead.contactEmail}`} className="text-primary">{lead.contactEmail}</a></div>
            {lead.contactPhone && <div><span className="text-gray-400">Phone:</span> {lead.contactPhone}</div>}
            {lead.businessName && <div><span className="text-gray-400">Business:</span> {lead.businessName}</div>}
          </div>
        </div>
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <div className="text-sm font-semibold text-gray-500 mb-3">Lead Score</div>
          <div className="flex items-end gap-3">
            <div className={`text-4xl font-bold ${lead.leadScore >= 70 ? 'text-red-500' : lead.leadScore >= 40 ? 'text-yellow-500' : 'text-gray-400'}`}>
              {lead.leadScore}
            </div>
            <div className="mb-1">
              <span className={`px-2 py-0.5 rounded-full text-xs font-semibold ${
                lead.leadTier === 'hot' ? 'bg-red-50 text-red-500' :
                lead.leadTier === 'warm' ? 'bg-yellow-50 text-yellow-600' :
                'bg-gray-100 text-gray-500'
              }`}>{lead.leadTier}</span>
            </div>
          </div>
          {roi && (
            <div className="mt-3 text-sm">
              <div className="text-green-600 font-bold text-lg">${Math.round(roi.recoverableCostPerYear || roi.annualValue || 0).toLocaleString()}/yr</div>
              <div className="text-gray-400">{Math.round(roi.totalHoursPerWeek || roi.hoursPerWeek || 0)} hrs/week recoverable</div>
            </div>
          )}
          <div className="mt-3">
            <select
              value={lead.status || 'new'}
              onChange={(e) => updateLeadStatus(TENANT_ID, id, e.target.value).then(() => setLead({ ...lead, status: e.target.value }))}
              className="text-sm border border-gray-200 rounded-lg px-3 py-2 bg-white w-full"
            >
              {['new', 'contacted', 'qualified', 'proposal', 'closed_won', 'closed_lost'].map((s) => (
                <option key={s} value={s}>{s.replace(/_/g, ' ')}</option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {audit?.executiveSummary && (
        <div className="bg-white rounded-xl border border-gray-200 p-5 mb-4">
          <div className="text-sm font-semibold text-gray-500 mb-2">Executive Summary</div>
          <p className="text-sm text-gray-600 leading-relaxed">{audit.executiveSummary}</p>
        </div>
      )}

      {audit?.findings?.length > 0 && (
        <div className="bg-white rounded-xl border border-gray-200 p-5 mb-4">
          <div className="text-sm font-semibold text-gray-500 mb-3">Findings</div>
          <ul className="space-y-2">
            {audit.findings.map((f, i) => (
              <li key={i} className="text-sm">
                <span className="font-medium">{f.icon} {f.title}</span>
                <span className={`ml-2 text-xs px-1.5 rounded-full ${f.priority === 'high' ? 'bg-red-50 text-red-500' : 'bg-yellow-50 text-yellow-600'}`}>{f.priority}</span>
                <p className="text-gray-500 mt-0.5">{f.description}</p>
              </li>
            ))}
          </ul>
        </div>
      )}

      {lead.answers && (
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <div className="text-sm font-semibold text-gray-500 mb-3">Raw Answers</div>
          <pre className="text-xs text-gray-500 overflow-x-auto bg-gray-50 rounded p-3">
            {JSON.stringify(lead.answers, null, 2)}
          </pre>
        </div>
      )}
    </div>
  );
}
