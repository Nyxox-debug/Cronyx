package renderers

import (
	"context"
	"github.com/Nyxox-debug/Cronyx/pkg/cronyx"
	"io/ioutil"

	bf "github.com/russross/blackfriday/v2"
)

type MarkdownRenderer struct{}

func (MarkdownRenderer) Render(ctx context.Context, tplPath string, data cronyx.DataPayload) (cronyx.RenderedDoc, error) {
	// Minimal: read template markdown and do naive replacement for keys in data[0]
	b, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return cronyx.RenderedDoc{}, err
	}
	md := string(b)
	// naive templating — you should replace with a proper templating engine (text/template or pongo2)
	// For demo, we just render markdown as-is.
	html := string(bf.Run([]byte(md)))
	return cronyx.RenderedDoc{HTML: html}, nil
}
