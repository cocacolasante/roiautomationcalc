package ai

type AuditAnalysis struct {
	Executive   string    `json:"executive"`
	Findings    []Finding `json:"findings"`
	TopPriority string    `json:"topPriority"`
	RiskFactors []string  `json:"riskFactors"`
}

type Finding struct {
	Category     string  `json:"category"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Impact       string  `json:"impact"`
	HoursPerWeek float64 `json:"hoursPerWeek"`
	AnnualCost   float64 `json:"annualCost"`
}

type AnalyzeRequest struct {
	BusinessName      string
	BusinessType      string
	TeamSize          int
	Industry          string
	WorkflowBreakdown string
	TotalHours        float64
	AnnualCost        float64
	RecoverableAnnual float64
	RecoverablePct    float64
	ProductsList      string
}
