package outputs

import (
	"context"
	"fmt"
	"github.com/Nyxox-debug/Cronyx/pkg/cronyx"
	"io/ioutil"
	"path/filepath"
	"time"
)

type FileOutputGenerator struct {
	OutDir string
}

func (g FileOutputGenerator) Generate(ctx context.Context, r cronyx.RenderedDoc, format string) (cronyx.OutputFile, error) {
	// For "html" we simply write to file. PDF / XLSX would need additional pipelines.
	filename := fmt.Sprintf("%d.%s", time.Now().Unix(), format)
	outPath := filepath.Join(g.OutDir, filename)
	var data []byte
	if format == "html" {
		data = []byte(r.HTML)
	} else if format == "pdf" {
		// placeholder — in prod call a PDF generator (chromedp or wkhtmltopdf)
		data = []byte(r.HTML) // placeholder
	} else {
		data = []byte(r.HTML)
	}
	if err := ioutil.WriteFile(outPath, data, 0644); err != nil {
		return cronyx.OutputFile{}, err
	}
	return cronyx.OutputFile{Name: filename, Path: outPath, Data: data}, nil
}
