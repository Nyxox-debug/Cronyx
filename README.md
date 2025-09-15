# Cronyx 📊

A flexible, developer-friendly Go package for automated report generation and scheduling.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![GoDoc](https://godoc.org/github.com/Nyxox-debug/Cronyx?status.svg)](https://godoc.org/github.com/Nyxox-debug/Cronyx)

## ✨ Features

- **🔌 Plugin Architecture**: Easily extensible with custom data loaders, renderers, outputs, and delivery methods
- **⏰ Cron Scheduling**: Built-in cron scheduler with flexible scheduling options
- **🏗️ Flexible Job Definition**: Simple struct-based job configuration with sensible defaults
- **🔄 Worker Pool**: Concurrent job execution with configurable workers
- **📊 Metrics**: Optional job execution metrics and monitoring
- **🎯 Context Support**: Full context support for timeouts and cancellation
- **🔧 Type Safety**: Strong typing throughout the API
- **📝 Rich Templates**: Go template engine with custom functions

## 🚀 Quick Start

### Installation

```bash
go get github.com/Nyxox-debug/Cronyx
```

### Basic Usage

```go
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
	eng.RegisterLoader("csv", loader.CSVLoader{})
	eng.RegisterRenderer("markdown", render.MarkdownRenderer{})
	eng.RegisterOutput("html", generate.FileOutputGenerator{OutDir: "./out"})
	eng.RegisterDelivery("console", deliver.ConsoleDelivery{})

	// Define job using struct literal
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
```

## 📖 Core Concepts

### Engine

The main orchestrator that manages workers, scheduling, and component registration.

### Jobs

Define what data to load, how to render it, what outputs to generate, and where to deliver them.

### Components

- **Loaders**: Fetch data from various sources (CSV, JSON, databases, APIs)
- **Renderers**: Transform data using templates (Markdown, HTML)
- **Outputs**: Generate final files (HTML, PDF, Excel, CSV)
- **Delivery**: Send results to destinations (email, Slack, S3, console)

## 🏗️ Job Configuration

Jobs are defined using the `ReportJob` struct with flexible configuration options:

```go
// Daily report example
dailyJob := cronyx.ReportJob{
	ID:           "daily-analytics",
	Name:         "Daily Analytics Report",
	TemplatePath: "templates/analytics.md",
	DataSource: cronyx.DataSourceConfig{
		"type": "database",
		"connection": "postgres://user:pass@localhost/db",
		"query": "SELECT * FROM analytics WHERE date = CURRENT_DATE",
	},
	Outputs:  []string{"html", "pdf"},
	Schedule: "0 9 * * *", // 9 AM daily
	Delivery: []cronyx.DeliveryConfig{
		{"type": "email", "to": "team@company.com", "subject": "Daily Analytics Report"},
		{"type": "slack", "webhook": "webhook-url", "channel": "#reports"},
	},
	Timeout: 5 * time.Minute,
	Labels:  map[string]string{"priority": "high", "team": "analytics"},
}

// Weekly report example
weeklyJob := cronyx.ReportJob{
	ID:           "weekly-summary",
	Name:         "Weekly Summary",
	TemplatePath: "templates/summary.md",
	DataSource: cronyx.DataSourceConfig{
		"type": "json",
		"path": "data/weekly.json",
	},
	Outputs:  []string{"excel"},
	Schedule: "0 17 * * 5", // Friday 5 PM
	Delivery: []cronyx.DeliveryConfig{
		{"type": "s3", "bucket": "reports-bucket", "prefix": "weekly/"},
	},
	Timeout: 2 * time.Minute,
}

// Custom schedule example
hourlyJob := cronyx.ReportJob{
	ID:           "hourly-metrics",
	Name:         "Hourly Metrics",
	TemplatePath: "templates/metrics.md",
	DataSource: cronyx.DataSourceConfig{
		"type": "http",
		"url": "https://api.example.com/metrics",
		"headers": map[string]string{"Authorization": "Bearer token"},
	},
	Outputs:  []string{"html"},
	Schedule: "0 * * * *", // Every hour
	Delivery: []cronyx.DeliveryConfig{
		{"type": "console"},
	},
	Timeout: 30 * time.Second,
}
```

## 🔧 Configuration

```go
config := &cronyx.Config{
    Workers:        8,                    // Number of worker goroutines
    QueueSize:      200,                  // Job queue buffer size
    DefaultTimeout: time.Minute * 5,      // Default job timeout
    EnableMetrics:  true,                 // Enable job metrics
    Logger:         &CustomLogger{},      // Custom logger
}

engine := cronyx.NewEngineWithConfig(config)
```

## 📊 Built-in Components

### Loaders

- **CSV**: Load data from CSV files
  ```go
  DataSource: cronyx.DataSourceConfig{"type": "csv", "path": "data.csv"}
  ```
- **JSON**: Load data from JSON files
  ```go
  DataSource: cronyx.DataSourceConfig{"type": "json", "path": "data.json"}
  ```
- **Database**: Load data from SQL databases
  ```go
  DataSource: cronyx.DataSourceConfig{
      "type": "database",
      "connection": "postgres://user:pass@localhost/db",
      "query": "SELECT * FROM table",
  }
  ```
- **HTTP**: Load data from REST APIs
  ```go
  DataSource: cronyx.DataSourceConfig{
      "type": "http",
      "url": "https://api.example.com/data",
      "headers": map[string]string{"Authorization": "Bearer token"},
  }
  ```

### Renderers

- **Markdown**: Render using Markdown templates with Go templating
- **HTML**: Direct HTML template rendering

### Outputs

- **HTML**: Generate HTML files
- **PDF**: Generate PDF reports
- **Excel**: Generate Excel spreadsheets
- **CSV**: Generate CSV exports

### Delivery

- **Console**: Print to console/logs
  ```go
  Delivery: []cronyx.DeliveryConfig{{"type": "console"}}
  ```
- **Email**: Send via SMTP
  ```go
  Delivery: []cronyx.DeliveryConfig{{
      "type": "email",
      "to": "recipient@example.com",
      "subject": "Report Title",
      "smtp_host": "smtp.gmail.com",
      "smtp_port": "587",
      "username": "sender@gmail.com",
      "password": "password",
  }}
  ```
- **Slack**: Post to Slack channels
  ```go
  Delivery: []cronyx.DeliveryConfig{{
      "type": "slack",
      "webhook": "https://hooks.slack.com/...",
      "channel": "#reports",
  }}
  ```
- **S3**: Upload to Amazon S3
  ```go
  Delivery: []cronyx.DeliveryConfig{{
      "type": "s3",
      "bucket": "my-reports-bucket",
      "prefix": "reports/",
      "region": "us-west-2",
  }}
  ```

## 🎨 Custom Components

### Custom Data Loader

```go
type APILoader struct {
    BaseURL string
    APIKey  string
}

func (a *APILoader) Load(ctx context.Context, cfg cronyx.DataSourceConfig) (cronyx.DataPayload, error) {
    endpoint := cfg["endpoint"]
    // Make HTTP request, parse response
    // Return cronyx.DataPayload{Rows: data}
}

// Register it
engine.RegisterLoader("api", &APILoader{
    BaseURL: "https://api.example.com",
    APIKey:  "your-api-key",
})
```

### Custom Output Generator

```go
type SlackMessageOutput struct {
    WebhookURL string
}

func (s *SlackMessageOutput) Generate(ctx context.Context, rendered cronyx.RenderedDoc, format string) (cronyx.OutputFile, error) {
    // Convert rendered content to Slack message format
    // Return cronyx.OutputFile with Slack-formatted content
}

engine.RegisterOutput("slack-message", &SlackMessageOutput{})
```

## 📝 Templates

Templates use Go's `text/template` with additional functions:

```markdown
# {{.Meta.timestamp}} Report

## Summary

Total records: {{len .Rows}}

## Details

{{range .Rows}}

- **{{.name}}**: {{.value}} ({{.category}})
  {{end}}

## Metrics

{{$electronics := 0}}
{{$clothing := 0}}
{{range .Rows}}
{{if eq .category "Electronics"}}
{{$electronics = add $electronics .value}}
{{else if eq .category "Clothing"}}  
{{$clothing = add $clothing .value}}
{{end}}
{{end}}

Electronics Total: {{$electronics}}
Clothing Total: {{$clothing}}
```

## 📈 Monitoring & Metrics

```go
config := &cronyx.Config{EnableMetrics: true}
engine := cronyx.NewEngineWithConfig(config)

// After running jobs
metrics := engine.GetMetrics()
fmt.Printf("Total jobs: %d\n", metrics.TotalJobs)
fmt.Printf("Success rate: %.2f%%\n",
    float64(metrics.SuccessfulJobs)/float64(metrics.TotalJobs)*100)
fmt.Printf("Average duration: %v\n", metrics.AvgDuration)
```

## 🧪 Testing

```go
func TestMyJob(t *testing.T) {
    engine := cronyx.NewEngine(1)
    engine.RegisterLoader("csv", &loader.CSVLoader{})
    engine.RegisterRenderer("markdown", &render.MarkdownRenderer{})
    engine.RegisterOutput("html", &generate.FileOutputGenerator{OutDir: "/tmp"})
    engine.RegisterDelivery("console", &deliver.ConsoleDelivery{})

    job := cronyx.ReportJob{
        ID:           "test-job",
        Name:         "Test Job",
        TemplatePath: "testdata/template.md",
        DataSource:   cronyx.DataSourceConfig{"type": "csv", "path": "testdata/data.csv"},
        Outputs:      []string{"html"},
        Delivery:     []cronyx.DeliveryConfig{{"type": "console"}},
        Timeout:      30 * time.Second,
    }

    err := engine.TestExecute(context.Background(), job)
    assert.NoError(t, err)
}
```

## 🔒 Error Handling

Cronyx provides typed errors for better error handling:

```go
err := engine.AddCronJob(job)
if errors.Is(err, cronyx.ErrNoLoader) {
    log.Printf("Data loader not found: %v", err)
} else if errors.Is(err, cronyx.ErrInvalidSchedule) {
    log.Printf("Invalid schedule: %v", err)
}
```

## 🎯 Best Practices

1. **Use descriptive job IDs and names** to make monitoring easier
2. **Set appropriate timeouts** for jobs based on their complexity
3. **Use labels** to organize and filter jobs
4. **Test jobs individually** with `TestExecute` before scheduling them
5. **Monitor metrics** in production environments
6. **Use structured logging** for better observability
7. **Implement graceful shutdown** in your applications
8. **Validate data source configurations** before creating jobs

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Setup

```bash
git clone https://github.com/Nyxox-debug/Cronyx.git
cd Cronyx
go mod tidy
go test ./...
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [robfig/cron](https://github.com/robfig/cron) for cron scheduling
- [russross/blackfriday](https://github.com/russross/blackfriday) for Markdown rendering
- Go community for excellent tooling and libraries
