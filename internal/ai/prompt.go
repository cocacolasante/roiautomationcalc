package ai

import "fmt"

func BuildAnalyzePrompt(req *AnalyzeRequest) string {
	return fmt.Sprintf(`You are a business automation consultant analyzing a small business's workflow data.

BUSINESS: %s (%s)
TEAM SIZE: %d people
INDUSTRY: %s

WORKFLOW HOURS BREAKDOWN:
%s

ROI CALCULATION:
Total automatable hours: %.1f/week
Annual cost of manual work: $%.0f
Recoverable with automation: $%.0f/year (%.0f%%)

PRODUCTS BEING RECOMMENDED:
%s

Generate a personalized automation audit analysis.

For each significant workflow finding, describe:
1. What the finding is (in business terms, not tech terms)
2. Why it matters for their specific business
3. The specific impact (time and money)
4. What a fix looks like

For the executive summary:
- Open by showing you understand their specific situation
- Quantify the opportunity clearly ($X/year in recoverable time)
- Frame automation as winning back time, not just saving money
- Create urgency without being pushy
- Keep it conversational

Tone: Professional but accessible. Use "your team" not "the company."
Business owner is probably not technical. Avoid jargon.

Respond ONLY with valid JSON:
{
  "executive": "2-3 paragraph executive summary",
  "findings": [
    {
      "category": "email|leads|billing|support|scheduling|reporting",
      "title": "Short finding title",
      "description": "What is happening and why it is a problem",
      "impact": "high|medium|low",
      "hoursPerWeek": 8.5,
      "annualCost": 15470
    }
  ],
  "topPriority": "The single highest-ROI opportunity in 1-2 sentences",
  "riskFactors": ["What happens if they do not automate (optional, 1-2 items)"]
}`,
		req.BusinessName, req.BusinessType, req.TeamSize, req.Industry,
		req.WorkflowBreakdown,
		req.TotalHours, req.AnnualCost, req.RecoverableAnnual, req.RecoverablePct,
		req.ProductsList,
	)
}
