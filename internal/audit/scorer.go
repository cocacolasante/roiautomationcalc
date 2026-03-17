package audit

import "github.com/blueprintautomation/blueprint-roi/internal/tenant"

type Scorer struct{}

func NewScorer() *Scorer { return &Scorer{} }

func (s *Scorer) Score(answers map[string]interface{}, roi *ROICalculation, _ *tenant.Tenant) (int, string) {
	score := 0

	switch {
	case roi.RecoverableCostPerYear >= 100000:
		score += 40
	case roi.RecoverableCostPerYear >= 50000:
		score += 30
	case roi.RecoverableCostPerYear >= 25000:
		score += 20
	case roi.RecoverableCostPerYear >= 10000:
		score += 10
	default:
		score += 5
	}

	switch {
	case roi.TotalHoursPerWeek >= 20:
		score += 20
	case roi.TotalHoursPerWeek >= 10:
		score += 15
	case roi.TotalHoursPerWeek >= 5:
		score += 10
	default:
		score += 5
	}

	if ts, ok := answers["team_size"]; ok {
		switch {
		case toFloat(ts) >= 10:
			score += 20
		case toFloat(ts) >= 5:
			score += 15
		case toFloat(ts) >= 2:
			score += 10
		default:
			score += 5
		}
	}

	if email, ok := answers["contact_email"]; ok && email != "" {
		score += 10
	}
	if phone, ok := answers["contact_phone"]; ok && phone != "" {
		score += 10
	}

	if score > 100 {
		score = 100
	}

	quality := "cold"
	switch {
	case score >= 70:
		quality = "hot"
	case score >= 40:
		quality = "warm"
	}

	return score, quality
}
