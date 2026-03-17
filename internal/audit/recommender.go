package audit

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/blueprintautomation/blueprint-roi/internal/tenant"
)

type Recommender struct{}

func NewRecommender() *Recommender { return &Recommender{} }

var priorityOrder = map[string]int{"high": 3, "medium": 2, "low": 1}

func (r *Recommender) Recommend(answers map[string]interface{}, roi *ROICalculation, t *tenant.Tenant) []*Recommendation {
	_ = t
	var recommendations []*Recommendation

	for _, product := range roi.getProducts() {
		score := r.scoreProduct(product, answers, roi)
		if score < 0.3 {
			continue
		}
		monthlyROI := product.EstHoursSavedPerWeek * 52 / 12 * roi.HourlyRate
		rec := &Recommendation{
			Product:               product,
			RelevanceScore:        score,
			EstimatedHoursSaved:   product.EstHoursSavedPerWeek,
			EstimatedMonthlySaved: monthlyROI,
			PaybackPeriodMonths:   product.TimeToROIMonths,
			SpecificInsight:       r.generateInsight(product, answers, roi),
			Priority:              r.classifyPriority(score, roi),
		}
		recommendations = append(recommendations, rec)
	}

	sort.Slice(recommendations, func(i, j int) bool {
		pi := priorityOrder[recommendations[i].Priority]
		pj := priorityOrder[recommendations[j].Priority]
		if pi != pj {
			return pi > pj
		}
		return recommendations[i].RelevanceScore > recommendations[j].RelevanceScore
	})
	return recommendations
}

func (r *Recommender) scoreProduct(product *Product, answers map[string]interface{}, roi *ROICalculation) float64 {
	score := 0.0
	for _, key := range product.RelevantKeys {
		answer, ok := answers[key]
		if !ok {
			continue
		}
		if item := roi.FindByKey(key); item != nil && item.HoursPerWeek >= product.MinWeeklyHours {
			score += 0.3
		}
		for reqKey, reqValue := range product.RequiredAnswers {
			if reqKey == key && matchesRequiredValue(answer, reqValue) {
				score += 0.4
			}
		}
	}
	return math.Min(score, 1.0)
}

func matchesRequiredValue(answer interface{}, required string) bool {
	strAnswer := fmt.Sprintf("%v", answer)
	if strings.HasPrefix(required, ">") {
		var threshold float64
		fmt.Sscanf(required[1:], "%f", &threshold)
		return toFloat(answer) > threshold
	}
	if strings.HasPrefix(required, "<") {
		var threshold float64
		fmt.Sscanf(required[1:], "%f", &threshold)
		return toFloat(answer) < threshold
	}
	return strAnswer == required
}

func (r *Recommender) generateInsight(product *Product, _ map[string]interface{}, roi *ROICalculation) string {
	for _, key := range product.RelevantKeys {
		if item := roi.FindByKey(key); item != nil && item.HoursPerWeek > 0 {
			return fmt.Sprintf("Based on your %.1f hours/week on related tasks, this solution could significantly reduce your team's workload.", item.HoursPerWeek)
		}
	}
	return product.ROIDescription
}

func (r *Recommender) classifyPriority(score float64, roi *ROICalculation) string {
	if score >= 0.7 && roi.RecoverableCostPerYear >= 25000 {
		return "high"
	}
	if score >= 0.5 || roi.RecoverableCostPerYear >= 10000 {
		return "medium"
	}
	return "low"
}

// getProducts is a helper — in practice populated from DB via ToolConfig
func (roi *ROICalculation) getProducts() []*Product {
	return []*Product{}
}
