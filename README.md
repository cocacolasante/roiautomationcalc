# Blueprint ROI — Automation Audit Tool

A production-ready, multi-tenant ROI Calculator & Audit Tool for automation agencies.

## What It Does

1. **Client fills out** a 5-minute wizard assessing their manual workflows
2. **Live ROI counter** animates as they answer
3. **Claude AI generates** a personalized audit with findings + recommendations
4. **Results page** shows annual savings, category breakdown, and recommended solutions
5. **Lead captured** in admin dashboard, synced to CRM, notification sent

---

## Stack

| Layer | Tech |
|-------|------|
| Backend | Go 1.22 + chi + pgx/v5 + asynq |
| AI | Anthropic Claude (claude-sonnet-4-20250514) |
| Database | PostgreSQL 16 (JSONB for flexible audit data) |
| Queue | Redis + Asynq |
| Audit Tool | React 18 + Vite (embeddable IIFE) |
| Admin Dashboard | React 18 + Vite + TailwindCSS |
| PDF | Headless Chrome (chromedp) |
| n8n Workflows | 6 ready-to-import automation workflows |

---

## Quick Start

```bash
# 1. Clone and configure
git clone <repo> && cd blueprint-roi
cp .env.example .env
# Edit .env — set ADMIN_KEY and ANTHROPIC_API_KEY at minimum

# 2. Start infrastructure
make docker-up

# 3. Seed default tenant + questions
make seed

# 4. Open the tool
open http://localhost:5173   # Audit tool
open http://localhost:5174   # Admin dashboard
```

---

## Development

```bash
# Backend only (requires local postgres + redis)
make server      # API server on :8080
make worker      # Background job worker

# Frontend
cd tool && npm install && npm run dev    # Audit tool on :5173
cd admin && npm install && npm run dev  # Admin on :5174

# Build embeddable bundle
make build-embed  # outputs tool/dist-embed/blueprint-roi.js
```

---

## Embed on Any Website

```html
<div id="blueprint-roi"></div>
<script>
  window.BlueprintROIConfig = {
    tenantId: "YOUR_TENANT_ID",
    apiBase: "https://your-api.example.com"
  };
</script>
<script src="https://your-api.example.com/embed/blueprint-roi.js" defer></script>
```

See [docs/EMBED_GUIDE.md](docs/EMBED_GUIDE.md) for full instructions.

---

## Provisioning a New Client

```bash
export ADMIN_KEY=your-admin-key
make provision SLUG=client-slug NAME="Client Name" EMAIL=client@example.com
```

See [docs/CLIENT_HANDOFF.md](docs/CLIENT_HANDOFF.md) for the complete onboarding flow.

---

## n8n Workflows

6 workflows in `n8n-workflows/`:

- **01** — Discord + email alert on new lead
- **02** — HubSpot CRM sync
- **03** — Automated follow-up email sequence
- **04** — GoHighLevel contact sync
- **05** — PDF delivery via email
- **06** — Weekly lead digest

See [n8n-workflows/README.md](n8n-workflows/README.md).

---

## API Reference

```
GET  /api/health
GET  /api/v1/config/{tenantId}       — Fetch tool config + questions
POST /api/v1/audit/{tenantId}        — Submit audit (generates AI analysis)
GET  /api/v1/audit/{auditId}         — Fetch audit result
POST /api/v1/analytics               — Track analytics event

# Admin (X-Blueprint-Admin-Key header required)
GET  /api/admin/tenants
POST /api/admin/tenants
GET  /api/admin/tenants/{id}
PUT  /api/admin/tenants/{id}
GET  /api/admin/tenants/{id}/leads
GET  /api/admin/tenants/{id}/analytics
GET  /api/admin/tenants/{id}/config
PUT  /api/admin/tenants/{id}/config
```

---

## Documentation

- [Architecture Decisions](docs/ARCHITECTURE_DECISIONS.md)
- [ROI Methodology](docs/ROI_METHODOLOGY_GUIDE.md)
- [Embed Guide](docs/EMBED_GUIDE.md)
- [Client Handoff](docs/CLIENT_HANDOFF.md)

---

## Environment Variables

See [.env.example](.env.example) for all variables with descriptions.

Key required variables:
- `DATABASE_URL` — PostgreSQL connection string
- `REDIS_URL` — Redis connection string
- `ADMIN_KEY` — Admin API authentication key
- `ANTHROPIC_API_KEY` — For AI audit generation

---

## Blueprint Command Integration

This product supports the [Blueprint Command](../portals/) partner/reseller portal. Partners can provision and manage tenants through a shared admin API.

### Database Changes

Three columns were added to the `tenants` table via migration:

| Column | Type | Description |
|--------|------|-------------|
| `partner_id` | `UUID` (nullable) | ID of the reseller partner who owns this tenant |
| `managed_by` | `TEXT` (default `'bpa'`) | Who manages this tenant (`'bpa'` or `'partner'`) |
| `client_id` | `TEXT` (nullable) | External CRM or client reference ID |

### API Key Types

| Key Prefix | Access Level |
|------------|-------------|
| `bpa_super_{random}` | Super admin — full access across all tenants |
| `bpa_partner_{uuid}_{random}` | Partner — scoped to their own tenants only |
| (regular key) | Standard admin key validated against `BLUEPRINT_ADMIN_KEY` env var |

### Partner-Scoped Endpoints

All admin endpoints respect partner scoping. When a `bpa_partner_` key is used, list endpoints automatically filter to only return that partner's tenants.

### Stats Endpoint

```
GET /api/admin/tenants/:id/stats
```

Called by Blueprint Command to populate client portal usage dashboards.

**Example response:**
```json
{
  "product_key": "roiautomationcalc",
  "tenant_id": "550e8400-e29b-41d4-a716-446655440000",
  "period": "last_30_days",
  "stats": {
    "calculations_run": 58,
    "reports_generated": 41,
    "avg_roi_percent": 312
  },
  "summary": "Your ROI calculator ran 58 analyses and generated 41 reports this month."
}
```

### Applying Migrations

```bash
./scripts/apply-reseller-migrations.sh
```
