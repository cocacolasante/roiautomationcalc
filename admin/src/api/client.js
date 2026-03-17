const BASE = '/api';

function getAdminKey() {
  return localStorage.getItem('blueprint_admin_key') || '';
}

async function request(path, options = {}) {
  const res = await fetch(`${BASE}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'X-Blueprint-Admin-Key': getAdminKey(),
      ...(options.headers || {}),
    },
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }
  const ct = res.headers.get('content-type') || '';
  return ct.includes('application/json') ? res.json() : res.text();
}

// Auth
export const verifyAdminKey = (key) => {
  localStorage.setItem('blueprint_admin_key', key);
  return request('/admin/verify');
};

// Tenants
export const listTenants = () => request('/admin/tenants');
export const createTenant = (data) => request('/admin/tenants', { method: 'POST', body: JSON.stringify(data) });
export const getTenant = (id) => request(`/admin/tenants/${id}`);
export const updateTenant = (id, data) => request(`/admin/tenants/${id}`, { method: 'PUT', body: JSON.stringify(data) });

// Leads
export const listLeads = (tenantId, params = {}) => {
  const qs = new URLSearchParams(params).toString();
  return request(`/admin/tenants/${tenantId}/leads${qs ? '?' + qs : ''}`);
};
export const getLead = (tenantId, leadId) => request(`/admin/tenants/${tenantId}/leads/${leadId}`);
export const updateLeadStatus = (tenantId, leadId, status) =>
  request(`/admin/tenants/${tenantId}/leads/${leadId}/status`, { method: 'PATCH', body: JSON.stringify({ status }) });

// Analytics
export const getAnalytics = (tenantId, params = {}) => {
  const qs = new URLSearchParams(params).toString();
  return request(`/admin/tenants/${tenantId}/analytics${qs ? '?' + qs : ''}`);
};

// Tool Config
export const getToolConfig = (tenantId) => request(`/admin/tenants/${tenantId}/config`);
export const updateToolConfig = (tenantId, config) =>
  request(`/admin/tenants/${tenantId}/config`, { method: 'PUT', body: JSON.stringify(config) });
