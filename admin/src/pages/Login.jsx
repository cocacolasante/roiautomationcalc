import { useState } from 'react';

export default function Login({ onLogin }) {
  const [key, setKey] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!key.trim()) return;
    setLoading(true);
    setError('');
    try {
      const res = await fetch('/api/admin/verify', {
        headers: { 'X-Blueprint-Admin-Key': key.trim() },
      });
      if (!res.ok) throw new Error('Invalid key');
      onLogin(key.trim());
    } catch {
      setError('Invalid admin key. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-sm border border-gray-200 p-8 w-full max-w-sm">
        <div className="text-center mb-6">
          <div className="text-3xl mb-3">🔐</div>
          <h1 className="text-xl font-bold text-gray-900">Blueprint ROI Admin</h1>
          <p className="text-sm text-gray-500 mt-1">Enter your admin key to continue</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <input
            type="password"
            value={key}
            onChange={(e) => setKey(e.target.value)}
            placeholder="Admin key"
            className="w-full border border-gray-200 rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:border-primary"
            autoFocus
          />
          {error && <p className="text-red-500 text-sm">{error}</p>}
          <button
            type="submit"
            disabled={loading || !key.trim()}
            className="w-full bg-primary text-white font-semibold py-2.5 rounded-lg text-sm hover:bg-blue-700 disabled:opacity-50 transition-colors"
          >
            {loading ? 'Verifying...' : 'Sign In'}
          </button>
        </form>
      </div>
    </div>
  );
}
