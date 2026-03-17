# Embed Guide

How to add the Blueprint ROI audit tool to any website.

## Quick Start

Paste this snippet before `</body>` on any page:

```html
<!-- Blueprint ROI Audit Tool -->
<div id="blueprint-roi"></div>
<script>
  window.BlueprintROIConfig = {
    tenantId: "YOUR_TENANT_ID",
    apiBase: "https://your-api.example.com"
  };
</script>
<script src="https://your-api.example.com/embed/blueprint-roi.js" defer></script>
```

That's it. The tool will render inside `#blueprint-roi`.

---

## Configuration Options

All options are set via `window.BlueprintROIConfig`:

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `tenantId` | string | ✅ | Your tenant ID (from admin dashboard) |
| `apiBase` | string | ✅ | Backend API URL |
| `containerId` | string | | Custom container ID (default: `blueprint-roi`) |

Visual branding (colors, logo, CTA text) is configured in the admin dashboard under **Tool Config** — no code changes needed.

---

## Serving the Bundle

The embed bundle is built to `tool/dist-embed/blueprint-roi.js`. Serve it statically:

```nginx
location /embed/ {
    root /app/tool/dist-embed;
    add_header Access-Control-Allow-Origin *;
    add_header Cache-Control "public, max-age=86400";
}
```

Or configure the Go server to serve it:
```go
// Already handled in cmd/server/main.go
r.Handle("/embed/*", http.StripPrefix("/embed", http.FileServer(http.Dir("tool/dist-embed"))))
```

---

## Building the Bundle

```bash
cd tool
npm install
npm run build:embed   # outputs to tool/dist-embed/blueprint-roi.js
```

Or via Makefile:
```bash
make build-embed
```

---

## WordPress

1. Install a plugin that allows custom HTML in footers (e.g. "Header and Footer Scripts")
2. Paste the snippet in the footer section
3. Or add directly to your theme's `footer.php` before `</body>`

## Webflow

1. Go to **Project Settings → Custom Code → Footer Code**
2. Paste the snippet
3. Publish your site

## Squarespace

1. Go to **Settings → Advanced → Code Injection → Footer**
2. Paste the snippet

---

## Security Notes

- The embed bundle makes requests to `apiBase` — ensure CORS is configured correctly
- The `tenantId` is public (it's used to fetch config, not authenticate)
- Rate limiting is applied server-side (60 req/min per IP per tenant)
