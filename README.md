# Cronyx 📊

A flexible, developer-friendly Go package for automated report generation and scheduling.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![GoDoc](https://godoc.org/github.com/Nyxox-debug/Cronyx?status.svg)](https://godoc.org/github.com/Nyxox-debug/Cronyx)

## ✨ Features

- **🔌 Plugin Architecture**: Easily extensible with custom data loaders, renderers, outputs, and delivery methods
- **⏰ Cron Scheduling**: Built-in cron scheduler with flexible scheduling options
- **🏗️ Builder Pattern**: Fluent API for creating jobs with sensible defaults
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
    "log"
    "time"

    "github.com/Nyxox-debug/Cronyx/pkg/cronyx"
    "github.com/Nyxox-debug/Cronyx/pkg/cronyx/loaders"
    "github.com/Nyxox-debug/Cronyx/pkg/cronyx/renderers"
    "github.com/Nyxox-debug/Cronyx/pkg/cronyx/outputs"
    "github.com/Nyxox-debug/Cronyx/pkg/cronyx/delivery"
)

func main() {
    // Create and configure engine
    engine := cronyx.NewEngine().
        WithLoader("csv", &loaders.CSVLoader{}).
        WithRenderer("markdown", &renderers.MarkdownRenderer{}).
        WithOutput("html", &outputs.FileOutputGenerator{OutDir: "./reports"}).
        WithDelivery("console", &delivery.ConsoleDelivery{})

    // Create a daily report job
    job, err := cronyx.DailyReport("Sales Report").
        WithTemplate("templates/sales.md").
        WithCSVData("data/sales.csv").
        OutputHTML().
        DeliverToConsole().
        Build()

    if err != nil {
        log.Fatal(err)
    }

    // Start the engine
    engine.Start()
    defer engine.Stop()

    // Schedule the job
    engine.ScheduleJob(job)

    // Keep running
    select {}
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

## 🏗️ Builder Pattern

The builder pattern makes creating jobs intuitive and type-safe:

```go
// Daily report
daily, err := cronyx.DailyReport("Daily Analytics").
    WithTemplate("templates/analytics.md").
    WithDatabaseData("postgres://...", "SELECT * FROM analytics").
    OutputHTML().
    OutputPDF().
    DeliverToEmail("team@company.com", "Daily Analytics Report").
    DeliverToSlack("webhook-url", "#reports").
    WithTimeout(time.Minute * 5).
    WithLabel("priority", "high").
    Build()

// Weekly report
weekly, err := cronyx.WeeklyReport("Weekly Summary").
    WithTemplate("templates/summary.md").
    WithJSONData("data/weekly.json").
    OutputExcel().
    ScheduleWeekly(time.Friday, 17, 0). // Friday 5 PM
    DeliverToS3("reports-bucket", "weekly/").
    Build()

// Custom schedule
custom, err := cronyx.NewJob("Hourly Metrics").
    WithTemplate("templates/metrics.md").
    WithAPIData("https://api.example.com/metrics").
    OutputHTML().
    ScheduleEvery(time.Hour).
    DeliverToConsole().
    Build()
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
- **JSON**: Load data from JSON files
- **Database**: Load data from SQL databases
- **HTTP**: Load data from REST APIs

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
- **Email**: Send via SMTP
- **Slack**: Post to Slack channels
- **S3**: Upload to Amazon S3
- **File System**: Save to local/network drives

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
engine.WithLoader("api", &APILoader{
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

engine.WithOutput("slack-message", &SlackMessageOutput{})
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
    engine := cronyx.NewEngine().
        WithLoader("csv", &loaders.CSVLoader{}).
        WithRenderer("markdown", &renderers.MarkdownRenderer{}).
        WithOutput("html", &outputs.FileOutputGenerator{OutDir: "/tmp"}).
        WithDelivery("console", &delivery.ConsoleDelivery{})

    job := cronyx.NewJob("Test Job").
        WithTemplate("testdata/template.md").
        WithCSVData("testdata/data.csv").
        OutputHTML().
        DeliverToConsole().
        MustBuild() // Panics on error - good for tests

    err := engine.ExecuteJob(context.Background(), job)
    assert.NoError(t, err)
}
```

## 🔒 Error Handling

Cronyx provides typed errors for better error handling:

```go
err := engine.ScheduleJob(job)
if errors.Is(err, cronyx.ErrNoLoader) {
    log.Printf("Data loader not found: %v", err)
} else if errors.Is(err, cronyx.ErrEmptySchedule) {
    log.Printf("Invalid schedule: %v", err)
}
```

## 🎯 Best Practices

1. **Use the builder pattern** for creating jobs - it prevents common mistakes
2. **Set appropriate timeouts** for jobs based on their complexity
3. **Use labels** to organize and filter jobs
4. **Test jobs individually** before scheduling them
5. **Monitor metrics** in production environments
6. **Use structured logging** for better observability
7. **Implement graceful shutdown** in your applications

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
