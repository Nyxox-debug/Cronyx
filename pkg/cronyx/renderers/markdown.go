package renderers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"text/template"
	"time"

	"github.com/Nyxox-debug/Cronyx/pkg/cronyx"
	bf "github.com/russross/blackfriday/v2"
)

type MarkdownRenderer struct{}

func (MarkdownRenderer) Render(ctx context.Context, tplPath string, data cronyx.DataPayload) (cronyx.RenderedDoc, error) {
	// Read template file
	b, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return cronyx.RenderedDoc{}, fmt.Errorf("failed to read template: %w", err)
	}

	// Create template with custom functions
	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"len": func(slice interface{}) int {
			switch v := slice.(type) {
			case []map[string]interface{}:
				return len(v)
			default:
				return 0
			}
		},
	}).Parse(string(b))
	if err != nil {
		return cronyx.RenderedDoc{}, fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare template data with metadata
	templateData := map[string]interface{}{
		"Rows": data.Rows,
		"Data": data.Rows, // alias for convenience
		"Meta": map[string]interface{}{
			"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
			"rows_count": len(data.Rows),
			"source":     tplPath,
		},
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return cronyx.RenderedDoc{}, fmt.Errorf("failed to execute template: %w", err)
	}

	// Convert markdown to HTML
	html := string(bf.Run(buf.Bytes()))

	return cronyx.RenderedDoc{
		HTML:    html,
		Content: buf.String(), // Store original markdown too
		Meta: map[string]interface{}{
			"source":     tplPath,
			"rows_count": len(data.Rows),
			"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
		},
	}, nil
}
