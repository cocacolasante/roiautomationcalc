import { useState } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
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
  const logout = () => { localStorage.removeItem('blueprint_admin_key'); setAuthed(false); };
  return { authed, login, logout };
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

  if (!authed) return <Login onLogin={login} />;

  return (
    <BrowserRouter>
      <Layout logout={logout}>
        <Routes>
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
