# Client Handoff Guide

How to deploy Blueprint ROI for a new client and hand it off.

## Step 1: Deploy Infrastructure

### Option A: Docker Compose (recommended for VPS)

```bash
# Clone the repo on your server
git clone ... && cd blueprint-roi

# Copy and fill in env vars
cp .env.example .env
nano .env  # set ADMIN_KEY, ANTHROPIC_API_KEY, DISCORD_WEBHOOK_URL

# Start everything
make docker-up

# Run migrations + seed
docker compose exec server ./server -migrate
make seed
```

### Option B: Managed Services

- **Database**: Supabase, Neon, or Railway PostgreSQL
- **Redis**: Upstash (free tier works)
- **Server**: Railway, Fly.io, or Render
- Set env vars in the platform dashboard

---

## Step 2: Provision the Client Tenant

```bash
export ADMIN_KEY=your-admin-key
export API_BASE=https://your-api.example.com

make provision SLUG=client-slug NAME="Client Business Name" EMAIL=client@example.com
```

This creates the tenant and returns a `tenantId`. Save it.

---

## Step 3: Configure the Tool

1. Open the admin dashboard: `https://your-admin.example.com`
2. Sign in with your admin key
3. Go to **Tool Config**
4. Fill in:
   - **Tool Title** — e.g. "Free Automation Audit for Acme Corp"
   - **Logo URL** — client's logo
   - **Primary Color** — client's brand color
   - **CTA Button Text** — e.g. "Book a Free Call with Sarah"
   - **CTA URL** — client's Calendly or contact page
   - **Default Hourly Rate** — appropriate for their industry

---

## Step 4: Generate Embed Code

1. Admin dashboard → **Embed Code**
2. Copy the snippet
3. Send to client (or add it yourself to their site)

---

## Step 5: Connect n8n Workflows (optional)

Import the workflows from `n8n-workflows/`:
1. `01_new_lead_notification.json` — Discord + email alerts
2. `03_followup_email_sequence.json` — Automated follow-up
3. Choose `02` (HubSpot) or `04` (GoHighLevel) for CRM

Configure the webhook URL in tenant settings:
```bash
curl -X PUT https://your-api.example.com/api/admin/tenants/{tenantId} \
  -H "X-Blueprint-Admin-Key: $ADMIN_KEY" \
  -H "Content-Type: application/json" \
  -d '{"crmAdapter":"webhook","webhookUrl":"https://your-n8n.com/webhook/blueprint-lead"}'
```

---

## Step 6: Test End-to-End

1. Open the embed URL (or `http://localhost:5173` for dev)
2. Complete the full audit with test data
3. Verify:
   - [ ] Results page shows ROI breakdown
   - [ ] Discord notification received
   - [ ] Lead appears in admin dashboard
   - [ ] Follow-up email sent (if n8n connected)
   - [ ] CRM contact created (if configured)

---

## Ongoing Admin Access

Give the client a read-only view by creating a second tenant with limited scoping, or share the admin dashboard login directly.

**Admin dashboard URL:** `https://your-admin.example.com`
**Admin key:** set in `.env` as `ADMIN_KEY`

---

## Pricing / Licensing

Each client deployment is a separate tenant. The admin key is shared across all tenants (it's an internal tool). For client-specific admin access, consider adding per-tenant admin key support in a future version.
