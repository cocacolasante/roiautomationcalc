# Architecture Decisions

## ADR-001: Go Backend with Chi Router

**Decision:** Go 1.22 + chi for the HTTP layer.

**Rationale:** Low memory footprint (important for VPS deployments), excellent concurrency for handling many simultaneous audit submissions, single binary deploys, and strong stdlib. Chi was chosen over Gin/Echo for its idiomatic `net/http` compatibility and middleware composability.

---

## ADR-002: PostgreSQL with JSONB for Flexible Audit Data

**Decision:** Primary relational store with JSONB columns for `answers`, `findings`, `recommendations`, `roi_summary`.

**Rationale:** Audit questions and findings evolve frequently. JSONB lets us store arbitrary structured data without migrations while maintaining queryability. Core identity fields (tenant_id, contact info, ROI values) stay as typed columns for indexing and reporting.

---

## ADR-003: Asynq for Background Job Queue

**Decision:** Redis-backed Asynq instead of a managed queue.

**Rationale:** PDF generation via headless Chrome takes 2–5 seconds; CRM sync can fail and needs retry. Asynq provides at-least-once delivery, retry with backoff, task inspection, and a built-in UI — all on the same Redis instance used for rate limiting. No separate queue infrastructure needed.

---

## ADR-004: Claude AI for Audit Analysis

**Decision:** Anthropic `claude-sonnet-4-20250514` for generating findings + executive summary.

**Rationale:** GPT-4o was tested but Claude produces more structured, less verbose JSON consistently. The prompt requests strict JSON output; we use a regex extractor to pull the JSON block even if the model wraps it in markdown fences. Temperature 0 for determinism.

---

## ADR-005: Dual ROI Calculation (Server + Client)

**Decision:** The same ROI math runs in both Go (`internal/audit/calculator.go`) and JavaScript (`tool/src/lib/calculations.js`).

**Rationale:** The live counter on the frontend must update instantly on every slider change — no API round-trip. The server recalculates authoritatively at submission time to prevent manipulation. The JavaScript implementation is intentionally kept as a close translation of the Go code (same formula, same constants).

---

## ADR-006: Multi-Tenancy via `tenant_id` Column

**Decision:** All tables have a `tenant_id` UUID FK; every query is scoped by it.

**Rationale:** Single-database multi-tenancy is simpler to operate and sufficient for the expected scale (hundreds of tenants, not thousands). Row-level security in PostgreSQL could be added later. Each tenant gets their own tool configuration, questions, and lead pipeline.

---

## ADR-007: Embeddable IIFE Bundle

**Decision:** A separate Vite build config (`vite.embed.config.js`) produces a self-contained IIFE bundle.

**Rationale:** Agency clients need to drop one `<script>` tag onto any page. The IIFE wraps React in a closure so it doesn't conflict with any React version already on the host page. The bundle reads `window.BlueprintROIConfig` for tenant ID and API base URL at runtime.

---

## ADR-008: No TypeScript

**Decision:** `.jsx` only, no TypeScript.

**Rationale:** This codebase is designed to be rapidly customized by non-TypeScript developers (agency owners, freelancers). TypeScript adds friction to quick config changes. The domain types are simple enough that JSDoc comments provide sufficient IDE hints.
