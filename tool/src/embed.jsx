import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import './index.css';

(function () {
  const scripts = document.querySelectorAll('script[data-tenant-id]');
  const script = scripts[scripts.length - 1];
  if (!script) return;

  const tenantId = script.getAttribute('data-tenant-id');
  const containerId = 'blueprint-roi-' + tenantId;
  let containerEl = document.getElementById(containerId);

  if (!containerEl) {
    containerEl = document.createElement('div');
    containerEl.id = containerId;
    script.after(containerEl);
  }

  createRoot(containerEl).render(<App tenantId={tenantId} embedded={true} />);
})();
