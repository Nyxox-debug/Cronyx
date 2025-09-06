package renderers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Nyxox-debug/Cronyx/pkg/cronyx"
	"io/ioutil"
	_ "strings"
	"text/template"

	bf "github.com/russross/blackfriday/v2"
)

type MarkdownRenderer struct{}

func (MarkdownRenderer) Render(ctx context.Context, tplPath string, data cronyx.DataPayload) (cronyx.RenderedDoc, error) {
	// Read template file
	b, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return cronyx.RenderedDoc{}, fmt.Errorf("failed to read template: %w", err)
	}

	// Create template
	tmpl, err := template.New("report").Parse(string(b))
	if err != nil {
		return cronyx.RenderedDoc{}, fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare template data
	templateData := map[string]interface{}{
		"Rows": data.Rows,
		"Data": data.Rows, // alias for convenience
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return cronyx.RenderedDoc{}, fmt.Errorf("failed to execute template: %w", err)
	}

	// Convert markdown to HTML
	html := string(bf.Run(buf.Bytes()))

	return cronyx.RenderedDoc{
		HTML: html,
		Meta: map[string]interface{}{
			"source":     tplPath,
			"rows_count": len(data.Rows),
		},
	}, nil
}
