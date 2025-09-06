package cronyx

import "context"

// DataLoader loads raw data for a job.
type DataLoader interface {
	// Load returns data in a canonical form (e.g., []map[string]interface{} or bytes)
	Load(ctx context.Context, cfg DataSourceConfig) (DataPayload, error)
}

// TemplateRenderer: given template + data returns rendered HTML and/or structured doc.
type TemplateRenderer interface {
	Render(ctx context.Context, tplPath string, data DataPayload) (RenderedDoc, error)
}

// OutputGenerator: takes rendered docs to produce final files (pdf, xlsx, csv)
type OutputGenerator interface {
	Generate(ctx context.Context, rendered RenderedDoc, format string) (OutputFile, error)
}

// DeliveryAdapter: delivers file(s) to a target (email, slack, s3)
type DeliveryAdapter interface {
	Deliver(ctx context.Context, target DeliveryConfig, files []OutputFile) error
}

// Minimal data structures
type DataPayload struct {
	// generic bag — implementers decide representation
	Rows []map[string]interface{}
	Raw  []byte
}

type RenderedDoc struct {
	HTML    string                 // for HTML→PDF pipelines
	Content string                 // raw content
	Meta    map[string]interface{} // metadata
}

type OutputFile struct {
	Name string
	Path string // local path or s3:// uri depending on storage adapter
	Data []byte // optional
}
