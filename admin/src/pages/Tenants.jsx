import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { listTenants, createTenant, deleteTenant } from '../api/client';

function DeleteConfirmModal({ tenant, onClose, onDeleted }) {
  const [loading, setLoading] = useState(false);
  async function handleDelete() {
    setLoading(true);
    try {
      await deleteTenant(tenant.id);
      onDeleted(tenant.id);
      onClose();
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-sm p-6">
        <h3 className="text-base font-semibold text-gray-900 mb-2">Delete Tenant</h3>
        <p className="text-sm text-gray-600 mb-5">
          Are you sure you want to delete <strong>{tenant.name}</strong>? This cannot be undone.
        </p>
        <div className="flex gap-2">
          <button onClick={onClose} className="flex-1 px-4 py-2 text-sm text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200">Cancel</button>
          <button onClick={handleDelete} disabled={loading} className="flex-1 px-4 py-2 text-sm text-white bg-red-600 rounded-lg hover:bg-red-700 disabled:opacity-50">
            {loading ? 'Deleting...' : 'Delete'}
          </button>
        </div>
      </div>
    </div>
  );
}

export default function Tenants() {
  const navigate = useNavigate();
  const [tenants, setTenants] = useState([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [form, setForm] = useState({ name: '', slug: '', adminEmail: '' });
  const [showForm, setShowForm] = useState(false);
  const [error, setError] = useState('');
  const [deleteTarget, setDeleteTarget] = useState(null);

  const load = () => {
    listTenants()
      .then((r) => setTenants(r.tenants || r || []))
      .catch(() => setTenants([]))
      .finally(() => setLoading(false));
  };

  useEffect(load, []);

  const handleCreate = async (e) => {
    e.preventDefault();
    setCreating(true);
    setError('');
    try {
      await createTenant(form);
      setForm({ name: '', slug: '', adminEmail: '' });
      setShowForm(false);
      load();
    } catch (err) {
      setError(err.message);
    } finally {
      setCreating(false);
    }
  };

  return (
    <div className="max-w-3xl">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-xl font-bold">Tenants</h1>
        <button
          onClick={() => setShowForm((v) => !v)}
          className="bg-primary text-white text-sm font-semibold px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
        >
          + New Tenant
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} className="bg-white rounded-xl border border-gray-200 p-5 mb-4 space-y-3">
          <div className="text-sm font-semibold text-gray-700 mb-2">Create Tenant</div>
          {['name', 'slug', 'adminEmail'].map((field) => (
            <input
              key={field}
              placeholder={field === 'adminEmail' ? 'Admin Email' : field.charAt(0).toUpperCase() + field.slice(1)}
              value={form[field]}
              onChange={(e) => setForm({ ...form, [field]: e.target.value })}
              className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary"
            />
          ))}
          {error && <p className="text-red-500 text-sm">{error}</p>}
          <div className="flex gap-2">
            <button type="submit" disabled={creating} className="bg-primary text-white text-sm font-semibold px-4 py-2 rounded-lg disabled:opacity-50">
              {creating ? 'Creating...' : 'Create'}
            </button>
            <button type="button" onClick={() => setShowForm(false)} className="text-sm text-gray-500 px-4 py-2 rounded-lg hover:bg-gray-50">
              Cancel
            </button>
          </div>
        </form>
      )}

      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        {loading ? (
          <div className="p-5 text-gray-400 text-sm">Loading...</div>
        ) : tenants.length === 0 ? (
          <div className="p-5 text-gray-400 text-sm">No tenants yet.</div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="text-gray-400 text-xs border-b border-gray-100 bg-gray-50">
                <th className="text-left px-5 py-3 font-medium">Name</th>
                <th className="text-left px-5 py-3 font-medium">Slug / ID</th>
                <th className="text-left px-5 py-3 font-medium">Status</th>
                <th className="px-5 py-3"></th>
              </tr>
            </thead>
            <tbody>
              {tenants.map((t) => (
                <tr key={t.id} className="border-b border-gray-50 hover:bg-gray-50">
                  <td className="px-5 py-3 font-medium">{t.name}</td>
                  <td className="px-5 py-3 text-gray-400 font-mono text-xs">{t.slug || t.id}</td>
                  <td className="px-5 py-3">
                    <span className={`px-2 py-0.5 rounded-full text-xs font-semibold ${t.isActive ? 'bg-green-50 text-green-600' : 'bg-gray-100 text-gray-500'}`}>
                      {t.isActive ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td className="px-5 py-3 text-right">
                    <div className="flex items-center justify-end gap-3">
                      <button
                        onClick={() => { localStorage.setItem('blueprint_tenant_id', t.id); localStorage.setItem('blueprint_tenant_name', t.name); navigate('/config'); }}
                        className="text-xs text-primary font-semibold hover:underline"
                      >
                        Configure
                      </button>
                      <button
                        onClick={() => setDeleteTarget(t)}
                        className="text-xs text-red-500 hover:text-red-700 font-semibold"
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
      {deleteTarget && (
        <DeleteConfirmModal
          tenant={deleteTarget}
          onClose={() => setDeleteTarget(null)}
          onDeleted={(id) => setTenants(prev => prev.filter(t => t.id !== id))}
        />
      )}
    </div>
  );
}
