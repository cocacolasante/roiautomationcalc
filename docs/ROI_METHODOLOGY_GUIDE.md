# ROI Methodology Guide

## How We Calculate Automation ROI

### Core Formula

```
Annual Value = Σ (hours_per_week × automation_pct × hourly_rate × 52)
```

For each workflow category question:

1. **Hours per week** — entered by the user (slider 0–40)
2. **Automation percentage** — configured per question (how much of those hours can realistically be automated; defaults vary by category)
3. **Hourly rate** — configured per tenant (default $75/hr; users can optionally specify their own)
4. **52 weeks** — annualized

### Category Automation Percentages

Default values based on industry research. All configurable per question in the database.

| Category | Default Automation % | Rationale |
|----------|---------------------|-----------|
| Email / Communication | 70% | Triage, templates, routing — high automation ceiling |
| Lead Management | 80% | Data entry, scoring, initial outreach — very automatable |
| Billing / Invoicing | 75% | Invoice generation, payment reminders, reconciliation |
| Scheduling | 85% | Near-fully automatable with tools like Calendly |
| Reporting | 65% | Data pulls automatable; analysis less so |
| Customer Support | 60% | FAQs/status automatable; complex issues need humans |
| Social Media | 55% | Scheduling easy; creative work stays manual |
| General Operations | 60% | Wide variance; conservative estimate |

### Lead Scoring (0–100)

Four dimensions contribute to the lead score:

| Dimension | Max Points | How Scored |
|-----------|-----------|------------|
| Annual ROI value | 40 | $0=0, $5k=10, $20k=20, $50k=30, $100k+=40 |
| Hours per week | 20 | 0=0, 5=5, 10=10, 20=15, 30+=20 |
| Team size | 20 | 1=2, 2=5, 5=10, 10=15, 20+=20 |
| Contact completeness | 20 | name+email=10, +phone=15, +business=20 |

**Tiers:**
- **Hot** (70–100): High ROI, large team, complete contact — priority follow-up
- **Warm** (40–69): Moderate opportunity
- **Cold** (0–39): Early-stage exploration

---

## Conservative vs Aggressive Estimates

The default estimates are intentionally conservative. Actual automation savings often exceed these numbers because:

1. **Cognitive switching costs** — not captured (context-switching between tasks)
2. **Error correction time** — manual processes generate errors that require fixing
3. **Compound effects** — faster lead response rates increase close rates

When presenting to clients, position these as **floor estimates**.

---

## Hourly Rate Guidance

The default $75/hr represents a blended rate for a small business team. Adjust based on:

- **Knowledge workers / consultants**: $100–150/hr
- **Agency staff**: $50–80/hr
- **Service business ops**: $35–60/hr

The rate is configured per tenant in the tool config and can be updated via the admin dashboard.
