# n8n Workflows

Six ready-to-import n8n workflows for the Blueprint ROI system.

## Workflows

| # | File | Purpose |
|---|------|---------|
| 01 | `01_new_lead_notification.json` | Discord + email alert on every new lead; priority email for hot leads |
| 02 | `02_hubspot_crm_sync.json` | Upsert contact + create deal in HubSpot |
| 03 | `03_followup_email_sequence.json` | Automated D0 audit results + D2 follow-up email |
| 04 | `04_ghl_contact_sync.json` | Create contact + opportunity in GoHighLevel |
| 05 | `05_pdf_delivery.json` | Fetch PDF from API and email it to the lead |
| 06 | `06_weekly_lead_digest.json` | Monday 8am digest email with weekly pipeline summary |

## Import Instructions

1. In n8n, go to **Workflows → Import from File**
2. Select the `.json` file
3. Configure credentials (Discord, email SMTP, HubSpot API key, etc.)
4. Set environment variables in n8n settings:
   - `API_BASE` — Blueprint ROI backend URL (e.g. `https://api.yourdomain.com`)
   - `ADMIN_BASE_URL` — Admin dashboard URL
   - `TENANT_ID` — Your tenant ID
   - `DISCORD_WEBHOOK_URL` — Discord channel webhook
   - `SALES_EMAIL` — Email address to receive lead alerts
   - `BOOKING_URL` — Calendly / booking link

## Webhook Trigger URLs

After importing, each workflow exposes a webhook. Pass the URL to Blueprint ROI via the tenant CRM configuration:

```
POST /api/admin/tenants/{tenantId}
{
  "crmAdapter": "webhook",
  "webhookUrl": "https://your-n8n.com/webhook/blueprint-lead"
}
```

The server's webhook CRM adapter will POST the full lead payload to that URL on every new audit submission.
