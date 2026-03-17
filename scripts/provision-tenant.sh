#!/usr/bin/env bash
set -euo pipefail

SLUG="${1:-}"
NAME="${2:-}"
EMAIL="${3:-}"

if [[ -z "$SLUG" || -z "$NAME" || -z "$EMAIL" ]]; then
  echo "Usage: $0 <slug> <name> <admin-email>"
  echo "  Example: $0 acme-corp 'Acme Corp' owner@acme.com"
  exit 1
fi

ADMIN_KEY="${ADMIN_KEY:-${BP_ADMIN_KEY:-}}"
API_BASE="${API_BASE:-http://localhost:8080}"

if [[ -z "$ADMIN_KEY" ]]; then
  echo "Error: ADMIN_KEY or BP_ADMIN_KEY env var required."
  exit 1
fi

echo "→ Provisioning tenant '$SLUG'..."

RESPONSE=$(curl -sf -X POST "${API_BASE}/api/admin/tenants" \
  -H "Content-Type: application/json" \
  -H "X-Blueprint-Admin-Key: ${ADMIN_KEY}" \
  -d "{\"slug\":\"${SLUG}\",\"name\":\"${NAME}\",\"adminEmail\":\"${EMAIL}\"}")

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

TENANT_ID=$(echo "$RESPONSE" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])" 2>/dev/null || echo "")

if [[ -n "$TENANT_ID" ]]; then
  echo ""
  echo "✓ Tenant created!"
  echo "  ID:    $TENANT_ID"
  echo "  Slug:  $SLUG"
  echo ""
  echo "Next steps:"
  echo "  1. Configure the tool via the admin dashboard"
  echo "  2. Embed the tool: <script src=\"${API_BASE}/embed/blueprint-roi.js\"></script>"
  echo "     with window.BlueprintROIConfig = { tenantId: \"${TENANT_ID}\", apiBase: \"${API_BASE}\" }"
fi
