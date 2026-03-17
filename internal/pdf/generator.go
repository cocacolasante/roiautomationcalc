package pdf

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/blueprintautomation/blueprint-roi/internal/audit"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Generator renders audit reports to PDF using headless Chrome.
type Generator struct {
	templateDir string
	chromePath  string
}

func NewGenerator(templateDir, chromePath string) *Generator {
	return &Generator{templateDir: templateDir, chromePath: chromePath}
}

type templateData struct {
	BusinessName          string
	CompanyName           string
	LogoUrl               string
	Date                  string
	FormattedAnnualValue  string
	FormattedHoursPerWeek string
	FormattedHoursPerYear string
	LeadScore             int
	LeadTier              string
	ExecutiveSummary      string
	Findings              []audit.Finding
	Recommendations       []*audit.Recommendation
	CTAURL                string
	CTAText               string
	CTADescription        string
}

// GenerateData is additional metadata not stored on the Audit struct needed for rendering.
type GenerateData struct {
	CompanyName    string
	LogoURL        string
	CTAURL         string
	CTAText        string
	CTADescription string
}

// Generate renders the audit report HTML and converts it to PDF bytes.
func (g *Generator) Generate(ctx context.Context, a *audit.Audit, meta GenerateData) ([]byte, error) {
	html, err := g.renderHTML(a, meta)
	if err != nil {
		return nil, fmt.Errorf("render html: %w", err)
	}
	return g.htmlToPDF(ctx, html)
}

func (g *Generator) renderHTML(a *audit.Audit, meta GenerateData) (string, error) {
	tmplPath := filepath.Join(g.templateDir, "pdf", "audit_report.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	roi := a.ROISummary
	var hoursPerWeek, annualValue float64
	if roi != nil {
		hoursPerWeek = roi.TotalHoursPerWeek
		annualValue = roi.RecoverableCostPerYear
	}
	data := templateData{
		BusinessName:          a.BusinessName,
		CompanyName:           meta.CompanyName,
		LogoUrl:               meta.LogoURL,
		Date:                  time.Now().Format("January 2, 2006"),
		FormattedAnnualValue:  formatInt(int(math.Round(annualValue))),
		FormattedHoursPerWeek: formatInt(int(math.Round(hoursPerWeek))),
		FormattedHoursPerYear: formatInt(int(math.Round(hoursPerWeek * 52))),
		LeadScore:             a.LeadScore,
		LeadTier:              a.LeadQuality,
		ExecutiveSummary:      a.ExecutiveSummary,
		Findings:              a.Findings,
		Recommendations:       a.Recommendations,
		CTAURL:                meta.CTAURL,
		CTAText:               meta.CTAText,
		CTADescription:        meta.CTADescription,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}

func (g *Generator) htmlToPDF(ctx context.Context, html string) ([]byte, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.Headless,
	)
	if g.chromePath != "" {
		opts = append(opts, chromedp.ExecPath(g.chromePath))
	}
	if envPath := os.Getenv("CHROME_PATH"); envPath != "" && g.chromePath == "" {
		opts = append(opts, chromedp.ExecPath(envPath))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	chromeCtx, cancel = context.WithTimeout(chromeCtx, 30*time.Second)
	defer cancel()

	var pdfBuf []byte
	if err := chromedp.Run(chromeCtx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.5).
				WithPaperHeight(11).
				WithMarginTop(0.4).
				WithMarginBottom(0.4).
				WithMarginLeft(0.4).
				WithMarginRight(0.4).
				Do(ctx)
			return err
		}),
	); err != nil {
		return nil, fmt.Errorf("chromedp: %w", err)
	}

	return pdfBuf, nil
}

func formatInt(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%d,%03d", n/1000, n%1000)
}
