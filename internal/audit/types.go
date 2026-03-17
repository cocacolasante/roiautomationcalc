package audit

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID               uuid.UUID        `json:"id"`
	TenantID         uuid.UUID        `json:"tenantId"`
	Key              string           `json:"key"`
	Category         string           `json:"category"`
	Question         string           `json:"question"`
	Hint             string           `json:"hint"`
	InputType        string           `json:"inputType"`
	Options          []QuestionOption `json:"options"`
	MinValue         float64          `json:"minValue"`
	MaxValue         float64          `json:"maxValue"`
	DefaultValue     float64          `json:"defaultValue"`
	Unit             string           `json:"unit"`
	HoursPerUnit     float64          `json:"hoursPerUnit"`
	IsAutomatic      bool             `json:"isAutomatic"`
	AutomationPct    float64          `json:"automationPct"`
	RelevantProducts []string         `json:"relevantProducts"`
	SortOrder        int              `json:"sortOrder"`
	IsActive         bool             `json:"isActive"`
	IsRequired       bool             `json:"isRequired"`
}

type QuestionOption struct {
	Label string  `json:"label"`
	Value string  `json:"value"`
	Hours float64 `json:"hours"`
}

type Product struct {
	ID                   uuid.UUID         `json:"id"`
	TenantID             uuid.UUID         `json:"tenantId"`
	Key                  string            `json:"key"`
	Name                 string            `json:"name"`
	Description          string            `json:"description"`
	Category             string            `json:"category"`
	SetupFee             float64           `json:"setupFee"`
	MonthlyFee           float64           `json:"monthlyFee"`
	ROIDescription       string            `json:"roiDescription"`
	MinWeeklyHours       float64           `json:"minWeeklyHours"`
	RelevantKeys         []string          `json:"relevantKeys"`
	RequiredAnswers      map[string]string `json:"requiredAnswers"`
	EstHoursSavedPerWeek float64           `json:"estHoursSavedPerWeek"`
	TimeToROIMonths      int               `json:"timeToRoiMonths"`
	SortOrder            int               `json:"sortOrder"`
	IsActive             bool              `json:"isActive"`
}

type ToolConfig struct {
	TenantID            string             `json:"tenantId"`
	CompanyName         string             `json:"companyName"`
	LogoURL             string             `json:"logoUrl"`
	PrimaryColor        string             `json:"primaryColor"`
	AccentColor         string             `json:"accentColor"`
	DomainName          string             `json:"domainName"`
	ToolTitle           string             `json:"toolTitle"`
	ToolSubtitle        string             `json:"toolSubtitle"`
	WelcomeMessage      string             `json:"welcomeMessage"`
	CTAText             string             `json:"ctaText"`
	CTAUrl              string             `json:"ctaUrl"`
	ThankyouMessage     string             `json:"thankyouMessage"`
	DefaultHourlyRate   float64            `json:"defaultHourlyRate"`
	IndustryHourlyRate  map[string]float64 `json:"industryHourlyRate"`
	WorkingHoursPerYear int                `json:"workingHoursPerYear"`
	AutomationSavingPct float64            `json:"automationSavingPct"`
	Questions           []*Question        `json:"questions"`
	Products            []*Product         `json:"products"`
	RequireContact      bool               `json:"requireContact"`
	CTASendReport       bool               `json:"ctaSendReport"`
	SendNotifications   bool               `json:"sendNotifications"`
	NotifyEmail         string             `json:"notifyEmail"`
	NotifyDiscord       string             `json:"notifyDiscord"`
	CRMAdapter          string             `json:"crmAdapter"`
	GoogleTagID         string             `json:"googleTagId"`
	MetaPixelID         string             `json:"metaPixelId"`
}

type AuditRequest struct {
	TenantID     uuid.UUID              `json:"tenantId"`
	BusinessName string                 `json:"businessName"`
	ContactName  string                 `json:"contactName"`
	ContactEmail string                 `json:"contactEmail"`
	ContactPhone string                 `json:"contactPhone"`
	TeamSize     int                    `json:"teamSize"`
	Industry     string                 `json:"industry"`
	Answers      map[string]interface{} `json:"answers"`
	SourceURL    string                 `json:"sourceUrl"`
	UTMSource    string                 `json:"utmSource"`
	UTMCampaign  string                 `json:"utmCampaign"`
	SessionID    string                 `json:"sessionId"`
}

type Audit struct {
	ID                   uuid.UUID              `json:"id"`
	TenantID             uuid.UUID              `json:"tenantId"`
	BusinessName         string                 `json:"businessName"`
	ContactName          string                 `json:"contactName"`
	ContactEmail         string                 `json:"contactEmail"`
	ContactPhone         string                 `json:"contactPhone"`
	TeamSize             int                    `json:"teamSize"`
	Industry             string                 `json:"industry"`
	Answers              map[string]interface{} `json:"answers"`
	ROISummary           *ROICalculation        `json:"roiSummary"`
	Findings             []Finding              `json:"findings"`
	ExecutiveSummary     string                 `json:"executiveSummary"`
	TopPriority          string                 `json:"topPriority"`
	Recommendations      []*Recommendation      `json:"recommendations"`
	LeadScore            int                    `json:"leadScore"`
	LeadQuality          string                 `json:"leadQuality"`
	EstimatedAnnualValue float64                `json:"estimatedAnnualValue"`
	Status               string                 `json:"status"`
	PDFPath              string                 `json:"pdfPath,omitempty"`
	PDFGeneratedAt       *time.Time             `json:"pdfGeneratedAt,omitempty"`
	SourceURL            string                 `json:"sourceUrl"`
	UTMSource            string                 `json:"utmSource"`
	UTMCampaign          string                 `json:"utmCampaign"`
	CRMContactID         string                 `json:"crmContactId,omitempty"`
	CRMSynced            bool                   `json:"crmSynced"`
	CreatedAt            time.Time              `json:"createdAt"`
	UpdatedAt            time.Time              `json:"updatedAt"`
}

type ROICalculation struct {
	TotalHoursPerWeek       float64       `json:"totalHoursPerWeek"`
	TotalHoursPerYear       float64       `json:"totalHoursPerYear"`
	TotalCostPerYear        float64       `json:"totalCostPerYear"`
	RecoverableHoursPerYear float64       `json:"recoverableHoursPerYear"`
	RecoverableCostPerYear  float64       `json:"recoverableCostPerYear"`
	BreakevenMonths         int           `json:"breakevenMonths"`
	Categories              []CategoryROI `json:"categories"`
	HourlyRate              float64       `json:"hourlyRate"`
	Breakdown               []ROILineItem `json:"breakdown"`
}

type CategoryROI struct {
	Category        string  `json:"category"`
	HoursPerWeek    float64 `json:"hoursPerWeek"`
	AnnualCost      float64 `json:"annualCost"`
	RecoverableCost float64 `json:"recoverableCost"`
	AutomationPct   float64 `json:"automationPct"`
	TopOpportunity  string  `json:"topOpportunity"`
}

type ROILineItem struct {
	QuestionKey      string  `json:"questionKey"`
	Description      string  `json:"description"`
	HoursPerWeek     float64 `json:"hoursPerWeek"`
	AnnualCost       float64 `json:"annualCost"`
	CanAutomate      bool    `json:"canAutomate"`
	AutomationPct    float64 `json:"automationPct"`
	RecoverableHours float64 `json:"recoverableHours"`
	RecoverableCost  float64 `json:"recoverableCost"`
}

type Finding struct {
	Category     string  `json:"category"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Impact       string  `json:"impact"`
	HoursPerWeek float64 `json:"hoursPerWeek"`
	AnnualCost   float64 `json:"annualCost"`
}

type Recommendation struct {
	Product               *Product `json:"product"`
	RelevanceScore        float64  `json:"relevanceScore"`
	EstimatedHoursSaved   float64  `json:"estimatedHoursSaved"`
	EstimatedMonthlySaved float64  `json:"estimatedMonthlySaved"`
	PaybackPeriodMonths   int      `json:"paybackPeriodMonths"`
	SpecificInsight       string   `json:"specificInsight"`
	Priority              string   `json:"priority"`
}

type Lead struct {
	ID           uuid.UUID `json:"id"`
	AuditID      uuid.UUID `json:"auditId"`
	TenantID     uuid.UUID `json:"tenantId"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Company      string    `json:"company"`
	LeadScore    int       `json:"leadScore"`
	EstimatedROI float64   `json:"estimatedRoi"`
	TopProducts  []string  `json:"topProducts"`
	CRMContactID string    `json:"crmContactId,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (roi *ROICalculation) FindByKey(key string) *ROILineItem {
	for i := range roi.Breakdown {
		if roi.Breakdown[i].QuestionKey == key {
			return &roi.Breakdown[i]
		}
	}
	return nil
}
