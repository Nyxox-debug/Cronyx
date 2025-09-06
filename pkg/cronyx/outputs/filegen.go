package outputs

import (
	"context"
	"fmt"
	"github.com/Nyxox-debug/Cronyx/pkg/cronyx"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type FileOutputGenerator struct {
	OutDir string
}

func (g FileOutputGenerator) Generate(ctx context.Context, r cronyx.RenderedDoc, format string) (cronyx.OutputFile, error) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(g.OutDir, 0755); err != nil {
		return cronyx.OutputFile{}, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("report_%s.%s", timestamp, format)
	outPath := filepath.Join(g.OutDir, filename)

	var data []byte
	switch format {
	case "html":
		data = []byte(r.HTML)
	case "pdf":
		// Placeholder - in production, use a PDF generator
		data = []byte(r.HTML)
	case "md":
		// Extract original markdown if available, otherwise convert HTML back
		data = []byte(r.HTML)
	default:
		data = []byte(r.HTML)
	}

	if err := ioutil.WriteFile(outPath, data, 0644); err != nil {
		return cronyx.OutputFile{}, fmt.Errorf("failed to write output file: %w", err)
	}

	return cronyx.OutputFile{
		Name: filename,
		Path: outPath,
		Data: data,
	}, nil
}
