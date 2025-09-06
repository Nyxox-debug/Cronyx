package main

import (
	"time"

	cronyx "github.com/Nyxox-debug/Cronyx/pkg/cronyx"
	deliver "github.com/Nyxox-debug/Cronyx/pkg/cronyx/delivery"
	loader "github.com/Nyxox-debug/Cronyx/pkg/cronyx/loaders"
	generate "github.com/Nyxox-debug/Cronyx/pkg/cronyx/outputs"
	render "github.com/Nyxox-debug/Cronyx/pkg/cronyx/renderers"
)

func main() {
	eng := cronyx.NewEngine(4)
	eng.RegisterLoader("csv", loader.CSVLoader{})
	eng.RegisterRenderer("markdown", render.MarkdownRenderer{})
	eng.RegisterOutput("html", generate.FileOutputGenerator{OutDir: "./out"})
	eng.RegisterDelivery("console", deliver.ConsoleDelivery{})

	job := cronyx.ReportJob{
		ID:           "job1",
		Name:         "daily-sample",
		TemplatePath: "examples/embed-example/sample.md",
		DataSource:   cronyx.DataSourceConfig{"type": "csv", "path": "examples/embed-example/data.csv"},
		Outputs:      []string{"html"},
		Schedule:     "@every 10s",
		Delivery:     []cronyx.DeliveryConfig{{"type": "console"}},
		Timeout:      30 * time.Second,
	}

	eng.AddCronJob(job)
	eng.Start()

	// run for 1 minute
	time.Sleep(1 * time.Minute)
	eng.Stop()
}
