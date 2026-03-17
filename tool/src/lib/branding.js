export function applyBranding(config) {
  if (!config) return;
  const root = document.documentElement;
  if (config.primaryColor) root.style.setProperty('--primary', config.primaryColor);
  if (config.accentColor) root.style.setProperty('--accent', config.accentColor);
  if (config.toolTitle) document.title = config.toolTitle;
}

export function getTenantIdFromEmbed() {
  const scripts = document.querySelectorAll('script[data-tenant-id]');
  if (scripts.length > 0) return scripts[scripts.length - 1].getAttribute('data-tenant-id');
  const match = window.location.pathname.match(/\/audit\/([^/]+)/);
  if (match) return match[1];
  const params = new URLSearchParams(window.location.search);
  return params.get('tenant') || 'blueprint';
}
