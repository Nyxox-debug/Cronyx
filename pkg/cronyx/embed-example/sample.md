# Daily Report

Generated on: {{.Meta.timestamp}}

## Data Summary

Total records: {{len .Rows}}

{{range .Rows}}

- **{{.name}}**: {{.value}} (Category: {{.category}})
  {{end}}

## Category Breakdown

{{$categories := index . "categories"}}
Electronics items: {{range .Rows}}{{if eq .category "Electronics"}}{{.name}} ({{.value}}) {{end}}{{end}}

## Additional Information

This report was generated automatically by Cronyx at {{.Meta.timestamp}}.
Source: {{.Meta.source}}
Total records processed: {{.Meta.rows_count}}
