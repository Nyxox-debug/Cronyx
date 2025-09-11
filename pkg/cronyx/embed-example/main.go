package main

import (
	"context"
	"fmt"
	"log"
	"time"

	cronyx "github.com/Nyxox-debug/Cronyx/pkg/cronyx"
	deliver "github.com/Nyxox-debug/Cronyx/pkg/cronyx/delivery"
	loader "github.com/Nyxox-debug/Cronyx/pkg/cronyx/loaders"
	generate "github.com/Nyxox-debug/Cronyx/pkg/cronyx/outputs"
	render "github.com/Nyxox-debug/Cronyx/pkg/cronyx/renderers"
)

func main() {
	// Create engine with 4 workers
	eng := cronyx.NewEngine(4)

	// Register components
	// NOTE: Add other methods
	eng.RegisterLoader("csv", loader.CSVLoader{})
	eng.RegisterRenderer("markdown", render.MarkdownRenderer{})
	eng.RegisterOutput("html", generate.FileOutputGenerator{OutDir: "./out"})
	eng.RegisterDelivery("console", deliver.ConsoleDelivery{})

	// Define job
	job := cronyx.ReportJob{
		ID:           "job1",
		Name:         "daily-sample",
		TemplatePath: "sample.md",
		DataSource:   cronyx.DataSourceConfig{"type": "csv", "path": "data.csv"},
		Outputs:      []string{"html"},
		Schedule:     "@every 10s",
		Delivery:     []cronyx.DeliveryConfig{{"type": "console"}},
		Timeout:      30 * time.Second,
	}

	// Test job execution once
	fmt.Println("=== Testing job execution ===")
	ctx := context.Background()
	if err := eng.TestExecute(ctx, job); err != nil {
		log.Printf("Job execution failed: %v", err)
		return
	}

	// Add to cron scheduler
	if err := eng.AddCronJob(job); err != nil {
		log.Printf("Failed to add cron job: %v", err)
		return
	}

	// Start engine
	fmt.Println("Starting Cronyx engine...")
	eng.Start()

	// Run for 1 minute then stop
	time.Sleep(1 * time.Minute)

	fmt.Println("Stopping engine...")
	eng.Stop()
}
