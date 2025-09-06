# Cronyx

Cronyx is a Go-first, embeddable report generation engine with optional API server.
It provides a pluggable pipeline: DataLoader → TemplateRenderer → OutputGenerator → DeliveryAdapter.

Quickstart (embed):

1. go get github.com/Nyxox-debug/Cronyx/pkg/cronyx
2. Register adapters and create ReportJob in your app (see examples/embed-example)
3. Start engine with engine.Start()

See /examples for runnable demos.
