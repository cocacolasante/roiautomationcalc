import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';
import App from './App';

;(function () {
  const hash = window.location.hash.slice(1);
  if (!hash) return;
  const params = new URLSearchParams(hash);
  const auth = params.get('auth');
  const tenant = params.get('tenant');
  const name = params.get('name');
  if (auth) {
    localStorage.setItem('blueprint_admin_key', auth);
    if (tenant) localStorage.setItem('blueprint_tenant_id', tenant);
    if (name) localStorage.setItem('blueprint_tenant_name', decodeURIComponent(name));
    window.history.replaceState(null, '', window.location.pathname);
  }
})();

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <App />
  </StrictMode>
);
