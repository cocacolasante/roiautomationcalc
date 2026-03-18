package audit

import (
	"math"
	"sort"
)

type Calculator struct{}

func NewCalculator() *Calculator { return &Calculator{} }

func (c *Calculator) Calculate(answers map[string]interface{}, config *ToolConfig) *ROICalculation {
	hourlyRate := config.DefaultHourlyRate
	if hourlyRate == 0 {
		hourlyRate = 35.0
	}
	weeksPerYear := 52.0
	globalAutoPct := config.AutomationSavingPct
	if globalAutoPct == 0 {
		globalAutoPct = 0.75
	}

	var lineItems []ROILineItem
	var totalHours float64
	categoryMap := map[string]*CategoryROI{}

	for _, question := range config.Questions {
		if !question.IsActive {
			continue
		}
		answer, ok := answers[question.Key]
		if !ok {
			continue
		}
		hours := extractHours(question, answer)
		if hours <= 0 {
			continue
		}
		autoPct := question.AutomationPct
		if autoPct == 0 {
			autoPct = globalAutoPct
		}
		annualCost := hours * weeksPerYear * hourlyRate
		recoverableHours := hours * autoPct
		recoverableCost := recoverableHours * weeksPerYear * hourlyRate
		totalHours += hours

		lineItems = append(lineItems, ROILineItem{
			QuestionKey:      question.Key,
			Description:      question.Question,
			HoursPerWeek:     hours,
			AnnualCost:       annualCost,
			CanAutomate:      question.IsAutomatic,
			AutomationPct:    autoPct,
			RecoverableHours: recoverableHours,
			RecoverableCost:  recoverableCost,
		})

		cat, exists := categoryMap[question.Category]
		if !exists {
			cat = &CategoryROI{Category: question.Category}
			categoryMap[question.Category] = cat
		}
		cat.HoursPerWeek += hours
		cat.AnnualCost += annualCost
		cat.RecoverableCost += recoverableCost
		if cat.AutomationPct == 0 {
			cat.AutomationPct = autoPct
		}
		if cat.TopOpportunity == "" {
			cat.TopOpportunity = question.Question
		}
	}

	// Also process automation area selections (from AutomationSelectionStep when no DB questions)
	type automationMeta struct {
		description string
		category    string
		autoPct     float64
	}
	automationOptions := map[string]automationMeta{
		"email":      {"Email & Communication", "email", 0.70},
		"leads":      {"Lead Management", "leads", 0.80},
		"billing":    {"Billing & Invoicing", "billing", 0.85},
		"support":    {"Customer Support", "support", 0.65},
		"scheduling": {"Scheduling & Appointments", "scheduling", 0.80},
		"reporting":  {"Reporting & Analytics", "reporting", 0.75},
		"social":     {"Social Media Management", "social", 0.70},
		"data_entry": {"Data Entry & Operations", "general", 0.80},
	}
	if areasRaw, ok := answers["automation_areas"]; ok {
		if areaSlice, ok := areasRaw.([]interface{}); ok {
			for _, areaRaw := range areaSlice {
				key, ok := areaRaw.(string)
				if !ok {
					continue
				}
				meta, exists := automationOptions[key]
				if !exists {
					continue
				}
				hours := toFloat(answers["automation_hours_"+key])
				if hours <= 0 {
					continue
				}
				annualCost := hours * weeksPerYear * hourlyRate
				recoverableHours := hours * meta.autoPct
				recoverableCost := recoverableHours * weeksPerYear * hourlyRate
				totalHours += hours

				lineItems = append(lineItems, ROILineItem{
					QuestionKey:      "automation_hours_" + key,
					Description:      meta.description,
					HoursPerWeek:     hours,
					AnnualCost:       annualCost,
					CanAutomate:      true,
					AutomationPct:    meta.autoPct,
					RecoverableHours: recoverableHours,
					RecoverableCost:  recoverableCost,
				})

				cat, exists := categoryMap[meta.category]
				if !exists {
					cat = &CategoryROI{Category: meta.category}
					categoryMap[meta.category] = cat
				}
				cat.HoursPerWeek += hours
				cat.AnnualCost += annualCost
				cat.RecoverableCost += recoverableCost
				if cat.AutomationPct == 0 {
					cat.AutomationPct = meta.autoPct
				}
			}
		}
	}

	totalRecoverableHours := totalHours * globalAutoPct
	totalRecoverableCost := totalRecoverableHours * weeksPerYear * hourlyRate
	totalCost := totalHours * weeksPerYear * hourlyRate

	var categories []CategoryROI
	for _, cat := range categoryMap {
		categories = append(categories, *cat)
	}
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].RecoverableCost > categories[j].RecoverableCost
	})

	breakevenMonths := calculateBreakeven(config.Products, totalRecoverableCost)

	return &ROICalculation{
		TotalHoursPerWeek:       math.Round(totalHours*10) / 10,
		TotalHoursPerYear:       math.Round(totalHours * weeksPerYear),
		TotalCostPerYear:        math.Round(totalCost),
		RecoverableHoursPerYear: math.Round(totalRecoverableHours * weeksPerYear),
		RecoverableCostPerYear:  math.Round(totalRecoverableCost),
		BreakevenMonths:         breakevenMonths,
		Categories:              categories,
		HourlyRate:              hourlyRate,
		Breakdown:               lineItems,
	}
}

func extractHours(q *Question, answer interface{}) float64 {
	switch q.InputType {
	case "hours_slider":
		return toFloat(answer)
	case "number", "people_count":
		return toFloat(answer) * q.HoursPerUnit
	case "multiple_choice":
		if strVal, ok := answer.(string); ok {
			for _, opt := range q.Options {
				if opt.Value == strVal {
					return opt.Hours
				}
			}
		}
	case "yes_no":
		if strVal, ok := answer.(string); ok && strVal == "yes" {
			return q.HoursPerUnit
		}
	}
	return 0
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	}
	return 0
}

func calculateBreakeven(products []*Product, annualRecoverable float64) int {
	if len(products) == 0 || annualRecoverable == 0 {
		return 0
	}
	var totalMonthlyCost float64
	for _, p := range products {
		if p.IsActive {
			totalMonthlyCost += p.MonthlyFee
		}
	}
	if totalMonthlyCost == 0 {
		return 0
	}
	monthlyRecoverable := annualRecoverable / 12
	months := math.Ceil(totalMonthlyCost / monthlyRecoverable * 12)
	if months > 60 {
		months = 60
	}
	return int(months)
}
