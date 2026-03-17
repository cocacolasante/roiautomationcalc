import { NavLink } from 'react-router-dom';
import { useState, useEffect } from 'react';

const NAV = [
  { to: '/dashboard', icon: '📊', label: 'Dashboard' },
  { to: '/leads', icon: '👥', label: 'Leads' },
  { to: '/analytics', icon: '📈', label: 'Analytics' },
  { to: '/config', icon: '⚙️', label: 'Tool Config' },
  { to: '/embed', icon: '🔗', label: 'Embed Code' },
  { to: '/tenants', icon: '🏢', label: 'Tenants' },
];

export default function Sidebar({ onLogout }) {
  const [tenantName, setTenantName] = useState(localStorage.getItem('blueprint_tenant_name'));

  useEffect(() => {
    const onStorage = () => setTenantName(localStorage.getItem('blueprint_tenant_name'));
    window.addEventListener('storage', onStorage);
    return () => window.removeEventListener('storage', onStorage);
  }, []);

  return (
    <aside className="w-56 bg-white border-r border-gray-200 flex flex-col shrink-0">
      <div className="p-5 border-b border-gray-100">
        <div className="font-bold text-lg text-primary">Blueprint ROI</div>
        <div className="text-xs text-gray-400 mt-0.5">Admin Dashboard</div>
        {tenantName && (
          <div className="mt-2 px-2 py-1 bg-blue-50 rounded-md">
            <div className="text-xs text-gray-400">Active tenant</div>
            <div className="text-xs font-semibold text-primary truncate">{tenantName}</div>
          </div>
        )}
      </div>

      <nav className="flex-1 p-3 space-y-0.5">
        {NAV.map(({ to, icon, label }) => (
          <NavLink
            key={to}
            to={to}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors ${
                isActive
                  ? 'bg-blue-50 text-primary font-semibold'
                  : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
              }`
            }
          >
            <span>{icon}</span>
            {label}
          </NavLink>
        ))}
      </nav>

      <div className="p-3 border-t border-gray-100">
        <button
          onClick={onLogout}
          className="w-full text-left px-3 py-2 text-sm text-gray-500 hover:text-red-500 rounded-lg hover:bg-gray-50 transition-colors"
        >
          🚪 Sign Out
        </button>
      </div>
    </aside>
  );
}
