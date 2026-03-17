import { useEffect, useState } from 'react';
import { getToolConfig, updateToolConfig } from '../api/client';

function Field({ label, name, value, onChange, type = 'text', hint }) {
  return (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-1">{label}</label>
      {hint && <p className="text-xs text-gray-400 mb-1">{hint}</p>}
      <input
        type={type}
        name={name}
        value={value || ''}
        onChange={onChange}
        className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary"
      />
    </div>
  );
}

export default function ToolConfig() {
  const tenantId = localStorage.getItem('blueprint_tenant_id');
  const [config, setConfig] = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    if (!tenantId) { setLoading(false); return; }
    getToolConfig(tenantId)
      .then(setConfig)
      .catch(() => setConfig({}))
      .finally(() => setLoading(false));
  }, [tenantId]);

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setConfig((c) => ({ ...c, [name]: type === 'checkbox' ? checked : value }));
  };

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    try {
      await updateToolConfig(tenantId, config);
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch (err) {
      alert('Save failed: ' + err.message);
    } finally {
      setSaving(false);
    }
  };

  if (!tenantId) return <div className="text-gray-400 p-4">No tenant selected. Go to <a href="/tenants" className="text-primary underline">Tenants</a> and click Configure.</div>;
  if (loading) return <div className="text-gray-400 p-4">Loading...</div>;

  return (
    <div className="max-w-2xl">
      <h1 className="text-xl font-bold mb-6">Tool Configuration</h1>

      <form onSubmit={handleSave} className="space-y-5 bg-white rounded-xl border border-gray-200 p-6">
        <div className="text-sm font-semibold text-gray-500 uppercase tracking-wide">Branding</div>
        <Field label="Tool Title" name="toolTitle" value={config?.toolTitle} onChange={handleChange} />
        <Field label="Tool Subtitle" name="toolSubtitle" value={config?.toolSubtitle} onChange={handleChange} />
        <Field label="Welcome Message" name="welcomeMessage" value={config?.welcomeMessage} onChange={handleChange} />
        <Field label="Company Name" name="companyName" value={config?.companyName} onChange={handleChange} />
        <Field label="Logo URL" name="logoUrl" value={config?.logoUrl} onChange={handleChange} hint="HTTPS URL to your logo image" />
        <Field label="Primary Color" name="primaryColor" value={config?.primaryColor || '#6C63FF'} onChange={handleChange} type="color" />

        <hr className="border-gray-100" />
        <div className="text-sm font-semibold text-gray-500 uppercase tracking-wide">CTA / Conversion</div>
        <Field label="CTA Button Text" name="ctaButtonText" value={config?.ctaButtonText} onChange={handleChange} />
        <Field label="CTA URL" name="ctaUrl" value={config?.ctaUrl} onChange={handleChange} hint="Where the CTA button links to" />
        <Field label="CTA Description" name="ctaDescription" value={config?.ctaDescription} onChange={handleChange} />

        <hr className="border-gray-100" />
        <div className="text-sm font-semibold text-gray-500 uppercase tracking-wide">ROI Settings</div>
        <Field label="Default Hourly Rate ($)" name="defaultHourlyRate" value={config?.defaultHourlyRate} onChange={handleChange} type="number" />

        <div className="flex items-center gap-3 pt-2">
          <button
            type="submit"
            disabled={saving}
            className="bg-primary text-white font-semibold px-5 py-2 rounded-lg text-sm hover:bg-blue-700 disabled:opacity-50 transition-colors"
          >
            {saving ? 'Saving...' : 'Save Changes'}
          </button>
          {saved && <span className="text-green-600 text-sm">✓ Saved!</span>}
        </div>
      </form>
    </div>
  );
}
