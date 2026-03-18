import { useState } from 'react';

export default function EmbedCode() {
  const tenantId = localStorage.getItem('blueprint_tenant_id') || '';
  const [host, setHost] = useState(window.location.origin.replace(':5174', ':8080'));
  const [copied, setCopied] = useState(false);

  const snippet = `<!-- Blueprint ROI Tool -->
<script src="${host}/embed/blueprint-roi.js" data-tenant-id="${tenantId}" defer></script>`;

  const copy = () => {
    navigator.clipboard.writeText(snippet).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  };

  return (
    <div className="max-w-2xl">
      <h1 className="text-xl font-bold mb-6">Embed Code</h1>

      <div className="bg-white rounded-xl border border-gray-200 p-6 mb-4">
        <div className="text-sm font-semibold text-gray-700 mb-2">API / Server Host</div>
        <input
          type="text"
          value={host}
          onChange={(e) => setHost(e.target.value)}
          className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-primary"
          placeholder="https://your-server.com"
        />
        <p className="text-xs text-gray-400 mt-1">The URL where your Blueprint ROI backend is running.</p>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <div className="px-5 py-3 border-b border-gray-100 flex items-center justify-between">
          <span className="text-sm font-semibold text-gray-700">Embed Snippet</span>
          <button
            onClick={copy}
            className="text-sm text-primary font-medium hover:underline"
          >
            {copied ? '✓ Copied!' : 'Copy'}
          </button>
        </div>
        <pre className="p-5 text-xs text-gray-600 overflow-x-auto bg-gray-50 leading-relaxed">{snippet}</pre>
      </div>

      <div className="mt-5 bg-blue-50 border border-blue-100 rounded-xl p-4 text-sm text-blue-700">
        <strong>How to use:</strong> Paste the snippet into any HTML page where you want the audit tool to appear.
        The <code className="bg-blue-100 px-1 rounded">#blueprint-roi</code> div is where the tool will render.
      </div>
    </div>
  );
}
