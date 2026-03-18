import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import './index.css';

(function () {
  const scripts = document.querySelectorAll('script[data-tenant-id]');
  const script = scripts[scripts.length - 1];
  if (!script) return;

  const tenantId = script.getAttribute('data-tenant-id');

  // Derive the ROI server base URL from where this script was loaded
  const apiBase = script.src.replace(/\/embed\/blueprint-roi\.js.*$/, '');

  // Inject the companion stylesheet
  const link = document.createElement('link');
  link.rel = 'stylesheet';
  link.href = apiBase + '/embed/blueprint-roi.css';
  document.head.appendChild(link);

  const containerId = 'blueprint-roi-' + tenantId;
  let containerEl = document.getElementById(containerId);

  if (!containerEl) {
    containerEl = document.createElement('div');
    containerEl.id = containerId;
    script.after(containerEl);
  }

  createRoot(containerEl).render(<App tenantId={tenantId} apiBase={apiBase} embedded={true} />);
})();
