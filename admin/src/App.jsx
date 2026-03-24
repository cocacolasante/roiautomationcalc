import { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate, useNavigate } from 'react-router-dom';
import Login from './pages/Login';
import Sidebar from './components/Sidebar';
import Dashboard from './pages/Dashboard';
import Leads from './pages/Leads';
import LeadDetail from './pages/LeadDetail';
import ToolConfig from './pages/ToolConfig';
import Analytics from './pages/Analytics';
import EmbedCode from './pages/EmbedCode';
import Tenants from './pages/Tenants';

function useAuth() {
  const [authed, setAuthed] = useState(!!localStorage.getItem('blueprint_admin_key'));
  const login = (key) => { localStorage.setItem('blueprint_admin_key', key); setAuthed(true); };
  const logout = () => {
    localStorage.removeItem('blueprint_admin_key');
    localStorage.removeItem('blueprint_tenant_id');
    localStorage.removeItem('blueprint_tenant_name');
    setAuthed(false);
  };
  return { authed, login, logout };
}

// After auto-auth, sets tenant context then redirects to dashboard.
function AutoAuthRedirect({ tenantId, tenantName }) {
  const navigate = useNavigate();
  useEffect(() => {
    if (tenantId) {
      localStorage.setItem('blueprint_tenant_id', tenantId);
      if (tenantName) localStorage.setItem('blueprint_tenant_name', decodeURIComponent(tenantName));
    }
    navigate('/dashboard', { replace: true });
  }, []);
  return null;
}

function Layout({ children, logout }) {
  return (
    <div className="flex h-screen overflow-hidden">
      <Sidebar onLogout={logout} />
      <main className="flex-1 overflow-y-auto p-6">{children}</main>
    </div>
  );
}

export default function App() {
  const { authed, login, logout } = useAuth();
  const [autoAuth, setAutoAuth] = useState(null);

  useEffect(() => {
    // Read hash params: #auth=KEY&tenant=TENANT_ID&name=TENANT_NAME
    // The hash is never sent to the server, keeping the key client-side only.
    const hash = window.location.hash.slice(1);
    if (!hash) return;
    const params = new URLSearchParams(hash);
    const authKey = params.get('auth');
    const tenantId = params.get('tenant');
    const tenantName = params.get('name') || '';
    if (!authKey) return;

    // Clear hash immediately so the key doesn't linger in the URL bar.
    history.replaceState(null, '', window.location.pathname);

    fetch('/api/admin/verify', {
      headers: { 'X-Blueprint-Admin-Key': authKey },
    }).then(r => {
      if (r.ok) { login(authKey); setAutoAuth({ tenantId, tenantName }); }
    }).catch(() => {});
  }, []);

  if (!authed) return <Login onLogin={login} />;

  return (
    <BrowserRouter basename="/admin">
      <Layout logout={logout}>
        <Routes>
          {autoAuth && (
            <Route path="*" element={
              <AutoAuthRedirect tenantId={autoAuth.tenantId} tenantName={autoAuth.tenantName} />
            } />
          )}
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/leads" element={<Leads />} />
          <Route path="/leads/:id" element={<LeadDetail />} />
          <Route path="/config" element={<ToolConfig />} />
          <Route path="/analytics" element={<Analytics />} />
          <Route path="/embed" element={<EmbedCode />} />
          <Route path="/tenants" element={<Tenants />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
}
