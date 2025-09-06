# Cronyx: Complete Project Documentation

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture & Design](#architecture--design)
3. [Project Structure](#project-structure)
4. [Core Interfaces](#core-interfaces)
5. [Implementation Details](#implementation-details)
6. [Configuration & Setup](#configuration--setup)
7. [Running the Application](#running-the-application)
8. [Code Walkthrough](#code-walkthrough)
9. [Advanced Usage](#advanced-usage)
10. [Troubleshooting](#troubleshooting)

---

## Project Overview

**Cronyx** is a Go-based automated report generation engine that follows a plugin-based architecture. It's designed to:

- **Load data** from various sources (CSV, databases, APIs)
- **Render templates** with that data (Markdown, HTML)
- **Generate outputs** in different formats (HTML, PDF, Excel)
- **Deliver reports** through multiple channels (Console, Email, Slack)
- **Schedule execution** using cron expressions

### Key Features

- 🔄 **Cron-based scheduling** - Run reports automatically
- 🔌 **Plugin architecture** - Easily extensible components
- 🎯 **Template-driven** - Markdown templates with Go templating
- 📊 **Multiple outputs** - HTML, PDF, Excel support
- 🚀 **Concurrent execution** - Worker pool for parallel processing
- 📦 **Modular design** - Clean separation of concerns

---

## Architecture & Design

Cronyx uses a **pipeline-based architecture** with four main stages:

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌──────────────┐
│ Data        │───▶│ Template     │───▶│ Output      │───▶│ Delivery     │
│ Loading     │    │ Rendering    │    │ Generation  │    │ Distribution │
└─────────────┘    └──────────────┘    └─────────────┘    └──────────────┘
```

### Core Components

1. **Engine**: Central orchestrator managing the pipeline
2. **DataLoaders**: Extract data from various sources
3. **TemplateRenderers**: Process templates with data
4. **OutputGenerators**: Create final output files
5. **DeliveryAdapters**: Distribute generated reports

### Design Patterns Used

- **Strategy Pattern**: Pluggable components for each pipeline stage
- **Registry Pattern**: Component registration and lookup
- **Worker Pool Pattern**: Concurrent job processing
- **Template Pattern**: Consistent execution pipeline

---

## Project Structure

```
Cronyx/
├── cmd/                          # Command-line applications
├── docs/                         # Documentation
│   └── Forme.md                 # Project documentation
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── internal/                    # Private application code
├── pkg/                         # Public library code
│   ├── cronyx/                  # Main package
│   │   ├── delivery/            # Delivery adapters
│   │   │   └── console.go       # Console output delivery
│   │   ├── embed-example/       # Example implementation
│   │   │   ├── data.csv         # Sample data file
│   │   │   ├── main.go          # Main application
│   │   │   └── sample.md        # Template file
│   │   ├── engine.go            # Core engine implementation
│   │   ├── interfaces.go        # Interface definitions
│   │   ├── job.go               # Job data structures
│   │   ├── loaders/             # Data loaders
│   │   │   └── csv.go           # CSV data loader
│   │   ├── outputs/             # Output generators
│   │   │   └── filegen.go       # File output generator
│   │   └── renderers/           # Template renderers
│   │       └── markdown.go      # Markdown renderer
│   └── prompt/                  # Prompt templates
│       └── prompt.txt           # System prompts
└── README.md                    # Project overview
```

### Directory Purpose Explanation

- **`cmd/`**: Entry points for different applications (CLI tools, servers)
- **`internal/`**: Code that's internal to the project (not importable by others)
- **`pkg/`**: Public library code that can be imported by other projects
- **`docs/`**: Project documentation and specifications
- **`embed-example/`**: Complete working example with sample data

---

## Core Interfaces

The power of Cronyx lies in its interface-based design. Let's examine each interface:

### 1. DataLoader Interface

```go
type DataLoader interface {
    Load(ctx context.Context, cfg DataSourceConfig) (DataPayload, error)
}
```

**Purpose**: Abstracts data loading from any source
**Implementations**: CSVLoader (more can be added: DatabaseLoader, APILoader, etc.)

### 2. TemplateRenderer Interface

```go
type TemplateRenderer interface {
    Render(ctx context.Context, tplPath string, data DataPayload) (RenderedDoc, error)
}
```

**Purpose**: Converts templates + data into rendered documents
**Implementations**: MarkdownRenderer (HTMLRenderer, PDFRenderer can be added)

### 3. OutputGenerator Interface

```go
type OutputGenerator interface {
    Generate(ctx context.Context, rendered RenderedDoc, format string) (OutputFile, error)
}
```

**Purpose**: Creates final output files in various formats
**Implementations**: FileOutputGenerator (S3Generator, DatabaseGenerator can be added)

### 4. DeliveryAdapter Interface

```go
type DeliveryAdapter interface {
    Deliver(ctx context.Context, target DeliveryConfig, files []OutputFile) error
}
```

**Purpose**: Distributes generated files to target destinations
**Implementations**: ConsoleDelivery (EmailDelivery, SlackDelivery can be added)

### Supporting Data Structures

```go
// Generic data container
type DataPayload struct {
    Rows []map[string]interface{}  // Structured data
    Raw  []byte                    // Raw data for binary sources
}

// Rendered document with metadata
type RenderedDoc struct {
    HTML    string                 // HTML content
    Content string                 // Raw content
    Meta    map[string]interface{} // Metadata
}

// Output file representation
type OutputFile struct {
    Name string  // File name
    Path string  // File path or URI
    Data []byte  // File content (optional)
}
```

---

## Implementation Details

### Engine Implementation (`engine.go`)

The Engine is the heart of Cronyx. Let's break down its implementation:

```go
type Engine struct {
    cronSched  *cron.Cron                    // Cron scheduler
    Loaders    map[string]DataLoader         // Registered data loaders
    Renderers  map[string]TemplateRenderer   // Registered renderers
    Outputs    map[string]OutputGenerator    // Registered output generators
    Deliveries map[string]DeliveryAdapter    // Registered delivery adapters

    jobQueue   chan ReportJob                // Job queue for workers
    workers    int                           // Number of worker goroutines
    stopCh     chan struct{}                 // Stop signal channel
}
```

#### Key Methods Explained

**1. NewEngine(workers int)**

```go
func NewEngine(workers int) *Engine {
    e := &Engine{
        cronSched:  cron.New(cron.WithSeconds()),  // Enable second-precision cron
        Loaders:    map[string]DataLoader{},       // Initialize empty registries
        Renderers:  map[string]TemplateRenderer{},
        Outputs:    map[string]OutputGenerator{},
        Deliveries: map[string]DeliveryAdapter{},
        jobQueue:   make(chan ReportJob, 100),     // Buffered channel for jobs
        workers:    workers,
        stopCh:     make(chan struct{}),           // Unbuffered stop channel
    }
    return e
}
```

**2. Registration Methods**

```go
func (e *Engine) RegisterLoader(name string, d DataLoader) {
    e.Loaders[name] = d  // Store loader by name for lookup
}
```

These methods populate the component registries for runtime lookup.

**3. Worker Pool Implementation**

```go
func (e *Engine) workerLoop(id int) {
    for {
        select {
        case job := <-e.jobQueue:                    // Receive job from queue
            ctx, cancel := context.WithTimeout(      // Create timeout context
                context.Background(),
                job.Timeout
            )
            _ = e.execute(ctx, job)                   // Execute job
            cancel()                                  // Clean up context
        case <-e.stopCh:                            // Receive stop signal
            return                                    // Exit worker
        }
    }
}
```

**4. Job Execution Pipeline**

```go
func (e *Engine) execute(ctx context.Context, job ReportJob) error {
    // 1. Load data
    dsType := job.DataSource["type"]
    loader, ok := e.Loaders[dsType]
    if !ok {
        return fmt.Errorf("no loader for type %s", dsType)
    }
    data, err := loader.Load(ctx, job.DataSource)
    if err != nil {
        return err
    }

    // 2. Render template
    renderer := e.Renderers["markdown"]  // Could be dynamic based on template
    rendered, err := renderer.Render(ctx, job.TemplatePath, data)
    if err != nil {
        return err
    }

    // 3. Generate outputs
    var files []OutputFile
    for _, fmtName := range job.Outputs {
        outGen, ok := e.Outputs[fmtName]
        if !ok {
            return fmt.Errorf("no output generator for %s", fmtName)
        }
        f, err := outGen.Generate(ctx, rendered, fmtName)
        if err != nil {
            return err
        }
        files = append(files, f)
    }

    // 4. Deliver files
    for _, dCfg := range job.Delivery {
        dtype := dCfg["type"]
        adapter, ok := e.Deliveries[dtype]
        if !ok {
            return fmt.Errorf("no delivery adapter for %s", dtype)
        }
        if err := adapter.Deliver(ctx, dCfg, files); err != nil {
            return err
        }
    }

    return nil
}
```

### CSV Loader Implementation (`loaders/csv.go`)

```go
type CSVLoader struct{}

func (c CSVLoader) Load(ctx context.Context, cfg cronyx.DataSourceConfig) (cronyx.DataPayload, error) {
    path := cfg["path"]                    // Extract file path from config
    f, err := os.Open(path)               // Open CSV file
    if err != nil {
        return cronyx.DataPayload{}, err
    }
    defer f.Close()                       // Ensure file is closed

    r := csv.NewReader(f)                 // Create CSV reader
    headers, err := r.Read()              // Read header row
    if err != nil {
        return cronyx.DataPayload{}, err
    }

    var rows []map[string]interface{}     // Store rows as key-value maps
    for {
        record, err := r.Read()           // Read each data row
        if err == io.EOF {
            break                         // End of file reached
        }
        if err != nil {
            return cronyx.DataPayload{}, err
        }

        row := map[string]interface{}{}   // Create row map
        for i, h := range headers {       // Map each value to its header
            row[h] = record[i]
        }
        rows = append(rows, row)          // Add to results
    }

    return cronyx.DataPayload{Rows: rows}, nil
}
```

### Markdown Renderer Implementation (`renderers/markdown.go`)

```go
func (MarkdownRenderer) Render(ctx context.Context, tplPath string, data cronyx.DataPayload) (cronyx.RenderedDoc, error) {
    // 1. Read template file
    b, err := ioutil.ReadFile(tplPath)
    if err != nil {
        return cronyx.RenderedDoc{}, fmt.Errorf("failed to read template: %w", err)
    }

    // 2. Create Go template
    tmpl, err := template.New("report").Parse(string(b))
    if err != nil {
        return cronyx.RenderedDoc{}, fmt.Errorf("failed to parse template: %w", err)
    }

    // 3. Prepare template data
    templateData := map[string]interface{}{
        "Rows": data.Rows,              // Make data available as .Rows
        "Data": data.Rows,              // Alias for convenience
    }

    // 4. Execute template with data
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, templateData); err != nil {
        return cronyx.RenderedDoc{}, fmt.Errorf("failed to execute template: %w", err)
    }

    // 5. Convert markdown to HTML using blackfriday
    html := string(bf.Run(buf.Bytes()))

    return cronyx.RenderedDoc{
        HTML:    html,
        Content: buf.String(),  // Keep original markdown
        Meta: map[string]interface{}{
            "source":     tplPath,
            "rows_count": len(data.Rows),
        },
    }, nil
}
```

### File Output Generator (`outputs/filegen.go`)

```go
func (g FileOutputGenerator) Generate(ctx context.Context, r cronyx.RenderedDoc, format string) (cronyx.OutputFile, error) {
    // 1. Ensure output directory exists
    if err := os.MkdirAll(g.OutDir, 0755); err != nil {
        return cronyx.OutputFile{}, fmt.Errorf("failed to create output directory: %w", err)
    }

    // 2. Generate unique filename with timestamp
    timestamp := time.Now().Format("20060102_150405")
    filename := fmt.Sprintf("report_%s.%s", timestamp, format)
    outPath := filepath.Join(g.OutDir, filename)

    // 3. Determine content based on format
    var data []byte
    switch format {
    case "html":
        data = []byte(r.HTML)       // Use HTML content
    case "pdf":
        data = []byte(r.HTML)       // Placeholder - would use PDF generator
    case "md":
        data = []byte(r.Content)    // Use original markdown
    default:
        data = []byte(r.HTML)       // Default to HTML
    }

    // 4. Write file to disk
    if err := ioutil.WriteFile(outPath, data, 0644); err != nil {
        return cronyx.OutputFile{}, fmt.Errorf("failed to write output file: %w", err)
    }

    return cronyx.OutputFile{
        Name: filename,
        Path: outPath,
        Data: data,
    }, nil
}
```

---

## Configuration & Setup

### Prerequisites

1. **Go 1.19+** installed on your system
2. **Git** for version control
3. **Text editor** or IDE

### Initial Setup

#### 1. Create Project Directory

```bash
mkdir cronyx-project
cd cronyx-project
```

#### 2. Initialize Go Module

```bash
go mod init github.com/your-username/cronyx
```

#### 3. Install Dependencies

```bash
go get github.com/robfig/cron/v3@v3.0.1
go get github.com/russross/blackfriday/v2@v2.1.0
```

Your `go.mod` should look like:

```go
module github.com/your-username/cronyx

go 1.19

require (
    github.com/robfig/cron/v3 v3.0.1
    github.com/russross/blackfriday/v2 v2.1.0
)
```

#### 4. Create Directory Structure

```bash
mkdir -p pkg/cronyx/{delivery,loaders,outputs,renderers,embed-example}
mkdir -p cmd docs internal
```

### Configuration Files

#### Sample Data (`pkg/cronyx/embed-example/data.csv`)

```csv
name,value,category,date
Product A,150,Electronics,2024-01-15
Product B,200,Clothing,2024-01-16
Product C,75,Books,2024-01-17
Product D,300,Electronics,2024-01-18
Product E,120,Home & Garden,2024-01-19
```

#### Sample Template (`pkg/cronyx/embed-example/sample.md`)

```markdown
# Daily Sales Report

**Generated**: {{.Meta.timestamp}}  
**Total Products**: {{len .Rows}}

## Product Inventory

{{range .Rows}}

### {{.name}}

- **Value**: ${{.value}}
- **Category**: {{.category}}
- **Date**: {{.date}}

---

{{end}}

## Summary by Category

{{$categories := .Categories}}
{{range $categories}}

- **{{.Name}}**: {{.Count}} items, ${{.Total}} total value
  {{end}}

## Report Details

- Generated by Cronyx Report Engine
- Data source: CSV file
- Template: Markdown with Go templating
- Output format: HTML

_This is an automated report. Please do not reply to this document._
```

---

## Running the Application

### Step-by-Step Execution Guide

#### 1. Prepare the Environment

```bash
# Navigate to project root
cd /path/to/your/cronyx-project

# Verify Go installation
go version

# Check dependencies
go mod tidy
go mod verify
```

#### 2. Create the Main Application

Create `pkg/cronyx/embed-example/main.go`:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    cronyx "github.com/your-username/cronyx/pkg/cronyx"
    deliver "github.com/your-username/cronyx/pkg/cronyx/delivery"
    loader "github.com/your-username/cronyx/pkg/cronyx/loaders"
    generate "github.com/your-username/cronyx/pkg/cronyx/outputs"
    render "github.com/your-username/cronyx/pkg/cronyx/renderers"
)

func main() {
    fmt.Println("🚀 Starting Cronyx Report Engine...")

    // 1. Create engine with 4 worker threads
    eng := cronyx.NewEngine(4)
    fmt.Println("✅ Engine created with 4 workers")

    // 2. Register all components
    eng.RegisterLoader("csv", loader.CSVLoader{})
    eng.RegisterRenderer("markdown", render.MarkdownRenderer{})
    eng.RegisterOutput("html", generate.FileOutputGenerator{OutDir: "./out"})
    eng.RegisterDelivery("console", deliver.ConsoleDelivery{})
    fmt.Println("✅ All components registered")

    // 3. Define the report job
    job := cronyx.ReportJob{
        ID:           "daily-sales-001",
        Name:         "daily-sales-report",
        TemplatePath: "sample.md",
        DataSource: cronyx.DataSourceConfig{
            "type": "csv",
            "path": "data.csv",
        },
        Outputs:  []string{"html"},
        Schedule: "@every 15s",  // Run every 15 seconds for demo
        Delivery: []cronyx.DeliveryConfig{
            {"type": "console"},
        },
        Timeout: 30 * time.Second,
    }

    // 4. Test single execution first
    fmt.Println("\n🧪 Testing single job execution...")
    ctx := context.Background()
    if err := eng.TestExecute(ctx, job); err != nil {
        log.Fatalf("❌ Job execution failed: %v", err)
    }
    fmt.Println("✅ Single execution successful!")

    // 5. Add job to cron scheduler
    if err := eng.AddCronJob(job); err != nil {
        log.Fatalf("❌ Failed to add cron job: %v", err)
    }
    fmt.Println("✅ Job added to scheduler")

    // 6. Start the engine
    fmt.Println("\n🔄 Starting scheduled execution...")
    eng.Start()

    // 7. Run for 2 minutes then stop
    fmt.Println("⏰ Running for 2 minutes. Press Ctrl+C to stop early.")
    time.Sleep(2 * time.Minute)

    // 8. Graceful shutdown
    fmt.Println("\n🛑 Shutting down...")
    eng.Stop()
    fmt.Println("👋 Cronyx stopped. Goodbye!")
}
```

#### 3. Execute the Application

**Option A: Run from embed-example directory**

```bash
cd pkg/cronyx/embed-example
go run main.go
```

**Option B: Run from project root**

```bash
# Update paths in main.go first:
# TemplatePath: "pkg/cronyx/embed-example/sample.md"
# DataSource path: "pkg/cronyx/embed-example/data.csv"
# OutDir: "./out"

go run pkg/cronyx/embed-example/main.go
```

#### 4. Expected Output

```
🚀 Starting Cronyx Report Engine...
✅ Engine created with 4 workers
✅ All components registered

🧪 Testing single job execution...
=== Executing job: daily-sales-report ===
✅ Single execution successful!
✅ Job added to scheduler

🔄 Starting scheduled execution...
⏰ Running for 2 minutes. Press Ctrl+C to stop early.
[Cronyx] Delivered file: report_20240315_143022.html (./out/report_20240315_143022.html)
[Cronyx] Delivered file: report_20240315_143037.html (./out/report_20240315_143037.html)
[Cronyx] Delivered file: report_20240315_143052.html (./out/report_20240315_143052.html)
...

🛑 Shutting down...
👋 Cronyx stopped. Goodbye!
```

#### 5. Verify Output Files

```bash
# Check generated files
ls -la out/

# View a generated report
open out/report_20240315_143022.html
# or
cat out/report_20240315_143022.html
```

---

## Code Walkthrough

Let's trace through a complete execution cycle:

### 1. Application Startup

```go
// main.go starts here
func main() {
    // Creates new engine instance
    eng := cronyx.NewEngine(4)
```

**What happens**:

- Creates cron scheduler with second precision
- Initializes empty component registries
- Creates buffered job queue (capacity: 100)
- Sets up worker pool with 4 goroutines

### 2. Component Registration

```go
eng.RegisterLoader("csv", loader.CSVLoader{})
eng.RegisterRenderer("markdown", render.MarkdownRenderer{})
eng.RegisterOutput("html", generate.FileOutputGenerator{OutDir: "./out"})
eng.RegisterDelivery("console", deliver.ConsoleDelivery{})
```

**What happens**:

- Each component is stored in its respective registry map
- Components are looked up by string keys during execution
- This enables runtime plugin selection and extensibility

### 3. Job Definition

```go
job := cronyx.ReportJob{
    ID:           "daily-sales-001",
    Name:         "daily-sales-report",
    TemplatePath: "sample.md",
    DataSource:   cronyx.DataSourceConfig{"type": "csv", "path": "data.csv"},
    Outputs:      []string{"html"},
    Schedule:     "@every 15s",
    Delivery:     []cronyx.DeliveryConfig{{"type": "console"}},
    Timeout:      30 * time.Second,
}
```

**What happens**:

- Job configuration is created with all necessary parameters
- DataSource config tells the engine which loader to use
- Outputs array specifies which formats to generate
- Delivery config specifies where to send the results

### 4. Engine Startup & Worker Pool

```go
eng.Start()
```

**What happens**:

- Starts 4 worker goroutines running `workerLoop()`
- Each worker listens on the job queue channel
- Starts the cron scheduler
- Workers are now ready to process jobs

### 5. Job Execution Pipeline

When a job is triggered (either by cron or manual enqueue):

**Step 1: Data Loading**

```go
// Engine looks up loader by type
dsType := job.DataSource["type"]  // "csv"
loader, ok := e.Loaders[dsType]   // Gets CSVLoader instance

// CSVLoader.Load() is called
data, err := loader.Load(ctx, job.DataSource)
```

**CSVLoader execution**:

1. Opens file at path specified in DataSource
2. Reads CSV headers
3. Converts each row to `map[string]interface{}`
4. Returns DataPayload with all rows

**Step 2: Template Rendering**

```go
renderer := e.Renderers["markdown"]  // Gets MarkdownRenderer
rendered, err := renderer.Render(ctx, job.TemplatePath, data)
```

**MarkdownRenderer execution**:

1. Reads template file (`sample.md`)
2. Parses it as Go template
3. Executes template with data (replaces `{{.Rows}}`, etc.)
4. Converts resulting markdown to HTML using blackfriday
5. Returns RenderedDoc with HTML and metadata

**Step 3: Output Generation**

```go
for _, fmtName := range job.Outputs {  // ["html"]
    outGen, ok := e.Outputs[fmtName]   // Gets FileOutputGenerator
    f, err := outGen.Generate(ctx, rendered, fmtName)
}
```

**FileOutputGenerator execution**:

1. Creates output directory if it doesn't exist
2. Generates timestamped filename
3. Writes HTML content to file
4. Returns OutputFile with path and metadata

**Step 4: Delivery**

```go
for _, dCfg := range job.Delivery {    // [{"type": "console"}]
    adapter, ok := e.Deliveries[dtype] // Gets ConsoleDelivery
    adapter.Deliver(ctx, dCfg, files)
}
```

**ConsoleDelivery execution**:

1. Prints file information to console
2. Could be extended to send emails, post to Slack, etc.

### 6. Concurrent Execution

Multiple workers can process jobs simultaneously:

```
Worker 1: Job A (Data Loading)
Worker 2: Job B (Template Rendering)
Worker 3: Job C (Output Generation)
Worker 4: Idle (waiting for jobs)
```

This enables high throughput for multiple concurrent reports.

---

## Advanced Usage

### Custom Data Loader Example

Create `loaders/database.go`:

```go
package loaders

import (
    "context"
    "database/sql"
    "github.com/your-username/cronyx/pkg/cronyx"
    _ "github.com/lib/pq"  // PostgreSQL driver
)

type DatabaseLoader struct{}

func (d DatabaseLoader) Load(ctx context.Context, cfg cronyx.DataSourceConfig) (cronyx.DataPayload, error) {
    // Get connection details from config
    dsn := cfg["dsn"]
    query := cfg["query"]

    // Connect to database
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return cronyx.DataPayload{}, err
    }
    defer db.Close()

    // Execute query
    rows, err := db.QueryContext(ctx, query)
    if err != nil {
        return cronyx.DataPayload{}, err
    }
    defer rows.Close()

    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        return cronyx.DataPayload{}, err
    }

    // Process rows
    var result []map[string]interface{}
    for rows.Next() {
        // Create slice of interface{} for Scan
        values := make([]interface{}, len(columns))
        pointers := make([]interface{}, len(columns))
        for i := range values {
            pointers[i] = &values[i]
        }

        // Scan row
        if err := rows.Scan(pointers...); err != nil {
            return cronyx.DataPayload{}, err
        }

        // Convert to map
        row := make(map[string]interface{})
        for i, col := range columns {
            row[col] = values[i]
        }
        result = append(result, row)
    }

    return cronyx.DataPayload{Rows: result}, nil
}
```

**Usage**:

```go
// Register the custom loader
eng.RegisterLoader("database", loader.DatabaseLoader{})

// Use in job definition
job := cronyx.ReportJob{
    // ... other fields
    DataSource: cronyx.DataSourceConfig{
        "type": "database",
        "dsn":  "postgres://user:pass@localhost/db",
        "query": "SELECT * FROM sales WHERE date >= CURRENT_DATE - INTERVAL '7 days'",
    },
}
```

### Email Delivery Adapter

Create `delivery/email.go`:

```go
package delivery

import (
    "context"
    "fmt"
    "net/smtp"
    "github.com/your-username/cronyx/pkg/cronyx"
)

type EmailDelivery struct{}

func (e EmailDelivery) Deliver(ctx context.Context, target cronyx.DeliveryConfig, files []cronyx.OutputFile) error {
    // Extract email configuration
    smtpHost := target["smtp_host"]
    smtpPort := target["smtp_port"]
    username := target["username"]
    password := target["password"]
    to := target["to"]
    subject := target["subject"]

    // Create email message
    body := "Please find attached reports:\n\n"
    for _, file := range files {
        body += fmt.Sprintf("- %s\n", file.Name)
    }

    msg := []byte("To: " + to + "\r\n" +
        "Subject: " + subject + "\r\n" +
        "\r\n" +
        body + "\r\n")

    // Send email
    auth := smtp.PlainAuth("", username, password, smtpHost)
    addr := smtpHost + ":" + smtpPort
    return smtp.SendMail(addr, auth, username, []string{to}, msg)
}
```

### PDF Output Generator

Create `outputs/pdf.go`:

```go
package outputs

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "time"
    "github.com/your-username/cronyx/pkg/cronyx"
)

type PDFGenerator struct {
    OutDir string
}

func (p PDFGenerator) Generate(ctx context.Context, r cronyx.RenderedDoc, format string) (cronyx.OutputFile, error) {
    // Create temp HTML file
    timestamp := time.Now().Format("20060102_150405")
    htmlFile := filepath.Join(p.OutDir, fmt.Sprintf("temp_%s.html", timestamp))
    pdfFile := filepath.Join(p.OutDir, fmt.Sprintf("report_%s.pdf", timestamp))

    // Write HTML to temp file
    if err := ioutil.WriteFile(htmlFile, []byte(r.HTML), 0644); err != nil {
        return cronyx.OutputFile{}, err
    }
    defer os.Remove(htmlFile) // Clean up temp file

    // Convert HTML to PDF using wkhtmltopdf
    cmd := exec.CommandContext(ctx, "wkhtmltopdf", htmlFile, pdfFile)
    if err := cmd.Run(); err != nil {
        return cronyx.OutputFile{}, fmt.Errorf("PDF generation failed: %w", err)
    }

    // Read generated PDF
    data, err := ioutil.ReadFile(pdfFile)
    if err != nil {
        return cronyx.OutputFile{}, err
    }

    return cronyx.OutputFile{
        Name: filepath.Base(pdfFile),
        Path: pdfFile,
        Data: data,
    }, nil
}
```

### Complex Template with Functions

Create advanced template `templates/advanced-report.md`:

```markdown
# {{.Title | title}}

**Generated**: {{.Meta.timestamp | formatDate}}  
**Total Records**: {{len .Rows}}  
**Report ID**: {{.ReportID}}

## Executive Summary

{{if gt (len .Rows) 0}}

- **Highest Value**: ${{.Rows | maxValue "value"}}
- **Lowest Value**: ${{.Rows | minValue "value"}}
- **Average Value**: ${{.Rows | avgValue "value" | printf "%.2f"}}
- **Total Value**: ${{.Rows | sumValue "value"}}
  {{else}}
  No data available for this period.
  {{end}}

## Category Breakdown

{{range $category, $items := (.Rows | groupBy "category")}}

### {{$category}}

{{range $items}}

- **{{.name}}**: ${{.value}} ({{.date | formatDate}})
{{end}}
**Category Total**: ${{$items | sumValue "value"}}

---

{{end}}

## Detailed Listings

{{range $index, $row := .Rows}}
{{add $index 1}}. **{{$row.name}}**

- Value: ${{$row.value}}
- Category: {{$row.category}}
- Date: {{$row.date | formatDate}}
  {{if $row.description}}
- Description: {{$row.description}}
  {{end}}
  {{end}}

## Charts and Visualizations

{{if .Meta.include_charts}}

### Sales by Category
```

{{range (.Rows | groupBy "category")}}
{{.Key}}: {{repeat "█" (div .Total 50)}} (${{.Total}})
{{end}}

```

### Trend Analysis
```

{{range (.Rows | groupBy "date" | sort)}}
{{.Key}}: {{repeat "▓" (div (.Items | len) 2)}} ({{.Items | len}} items)
{{end}}

```
{{end}}

---
*Report generated by Cronyx Engine v1.0*
*Data as of {{.Meta.data_timestamp | formatDate}}*
```

### Enhanced Renderer with Template Functions

Create `renderers/advanced_markdown.go`:

```go
package renderers

import (
    "context"
    "fmt"
    "io/ioutil"
    "sort"
    "strconv"
    "strings"
    "text/template"
    "time"
    "bytes"

    "github.com/your-username/cronyx/pkg/cronyx"
    bf "github.com/russross/blackfriday/v2"
)

type AdvancedMarkdownRenderer struct{}

func (AdvancedMarkdownRenderer) Render(ctx context.Context, tplPath string, data cronyx.DataPayload) (cronyx.RenderedDoc, error) {
    // Read template
    b, err := ioutil.ReadFile(tplPath)
    if err != nil {
        return cronyx.RenderedDoc{}, fmt.Errorf("failed to read template: %w", err)
    }

    // Create template with custom functions
    tmpl := template.New("advanced-report").Funcs(template.FuncMap{
        "formatDate": formatDate,
        "title":      strings.Title,
        "add":        func(a, b int) int { return a + b },
        "maxValue":   func(rows []map[string]interface{}, field string) float64 {
            return aggregateFloat(rows, field, func(values []float64) float64 {
                if len(values) == 0 { return 0 }
                max := values[0]
                for _, v := range values[1:] {
                    if v > max { max = v }
                }
                return max
            })
        },
        "minValue": func(rows []map[string]interface{}, field string) float64 {
            return aggregateFloat(rows, field, func(values []float64) float64 {
                if len(values) == 0 { return 0 }
                min := values[0]
                for _, v := range values[1:] {
                    if v < min { min = v }
                }
                return min
            })
        },
        "avgValue": func(rows []map[string]interface{}, field string) float64 {
            return aggregateFloat(rows, field, func(values []float64) float64 {
                if len(values) == 0 { return 0 }
                sum := 0.0
                for _, v := range values {
                    sum += v
                }
                return sum / float64(len(values))
            })
        },
        "sumValue": func(rows []map[string]interface{}, field string) float64 {
            return aggregateFloat(rows, field, func(values []float64) float64 {
                sum := 0.0
                for _, v := range values {
                    sum += v
                }
                return sum
            })
        },
        "groupBy": func(rows []map[string]interface{}, field string) map[string][]map[string]interface{} {
            groups := make(map[string][]map[string]interface{})
            for _, row := range rows {
                if val, ok := row[field]; ok {
                    key := fmt.Sprintf("%v", val)
                    groups[key] = append(groups[key], row)
                }
            }
            return groups
        },
        "repeat": func(str string, count int) string {
            if count < 0 { count = 0 }
            if count > 100 { count = 100 } // Safety limit
            return strings.Repeat(str, count)
        },
        "div": func(a, b int) int {
            if b == 0 { return 0 }
            return a / b
        },
        "printf": fmt.Sprintf,
    })

    // Parse template
    tmpl, err = tmpl.Parse(string(b))
    if err != nil {
        return cronyx.RenderedDoc{}, fmt.Errorf("failed to parse template: %w", err)
    }

    // Prepare enhanced template data
    templateData := map[string]interface{}{
        "Rows": data.Rows,
        "Data": data.Rows,
        "Title": "Sales Report",
        "ReportID": fmt.Sprintf("RPT-%d", time.Now().Unix()),
        "Meta": map[string]interface{}{
            "timestamp": time.Now(),
            "data_timestamp": time.Now().Add(-1 * time.Hour),
            "include_charts": true,
            "rows_count": len(data.Rows),
        },
    }

    // Execute template
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, templateData); err != nil {
        return cronyx.RenderedDoc{}, fmt.Errorf("failed to execute template: %w", err)
    }

    // Convert to HTML
    html := string(bf.Run(buf.Bytes()))

    return cronyx.RenderedDoc{
        HTML:    html,
        Content: buf.String(),
        Meta: templateData["Meta"].(map[string]interface{}),
    }, nil
}

// Helper functions
func formatDate(t time.Time) string {
    return t.Format("January 2, 2006 at 3:04 PM")
}

func aggregateFloat(rows []map[string]interface{}, field string, aggFunc func([]float64) float64) float64 {
    var values []float64
    for _, row := range rows {
        if val, ok := row[field]; ok {
            if fval, err := strconv.ParseFloat(fmt.Sprintf("%v", val), 64); err == nil {
                values = append(values, fval)
            }
        }
    }
    return aggFunc(values)
}
```

### Multi-Format Job Configuration

```go
// Complex job with multiple outputs and deliveries
complexJob := cronyx.ReportJob{
    ID:           "quarterly-sales-001",
    Name:         "quarterly-sales-analysis",
    TemplatePath: "templates/quarterly-report.md",
    DataSource: cronyx.DataSourceConfig{
        "type": "database",
        "dsn":  "postgres://user:pass@localhost/sales_db",
        "query": `
            SELECT
                p.name,
                SUM(s.amount) as value,
                p.category,
                DATE(s.sale_date) as date,
                p.description
            FROM sales s
            JOIN products p ON s.product_id = p.id
            WHERE s.sale_date >= CURRENT_DATE - INTERVAL '3 months'
            GROUP BY p.id, p.name, p.category, DATE(s.sale_date)
            ORDER BY s.sale_date DESC
        `,
    },
    Outputs: []string{"html", "pdf", "excel"},
    Schedule: "0 0 8 1 * *", // First day of every month at 8 AM
    Delivery: []cronyx.DeliveryConfig{
        {
            "type": "email",
            "to": "executives@company.com",
            "subject": "Quarterly Sales Report",
            "smtp_host": "smtp.company.com",
            "smtp_port": "587",
            "username": "reports@company.com",
            "password": "secret123",
        },
        {
            "type": "slack",
            "webhook": "https://hooks.slack.com/services/...",
            "channel": "#sales-reports",
            "message": "📊 New quarterly sales report available",
        },
        {
            "type": "s3",
            "bucket": "company-reports",
            "prefix": "quarterly-sales/",
            "region": "us-west-2",
        },
    },
    Timeout: 5 * time.Minute,
    Labels: map[string]string{
        "department": "sales",
        "priority":   "high",
        "type":       "quarterly",
    },
}
```

---

## Monitoring and Observability

### Job Status Tracking

Add to `job.go`:

```go
type JobStatus string

const (
    StatusPending   JobStatus = "pending"
    StatusRunning   JobStatus = "running"
    StatusSuccess   JobStatus = "success"
    StatusFailed    JobStatus = "failed"
    StatusTimeout   JobStatus = "timeout"
)

type JobExecution struct {
    JobID     string
    RunID     string
    Status    JobStatus
    StartTime time.Time
    EndTime   time.Time
    Error     error
    OutputFiles []OutputFile
    Metrics   JobMetrics
}

type JobMetrics struct {
    DataLoadTime     time.Duration
    RenderTime       time.Duration
    OutputTime       time.Duration
    DeliveryTime     time.Duration
    TotalTime        time.Duration
    DataRows         int
    OutputSizeBytes  int64
}
```

### Enhanced Engine with Monitoring

Add to `engine.go`:

```go
type Engine struct {
    // ... existing fields

    // Monitoring
    executions  map[string]*JobExecution
    execMutex   sync.RWMutex
    metrics     *EngineMetrics
    logger      *log.Logger
}

type EngineMetrics struct {
    JobsCompleted   int64
    JobsFailed      int64
    TotalRunTime    time.Duration
    AverageRunTime  time.Duration
}

func (e *Engine) executeWithMonitoring(ctx context.Context, job ReportJob) error {
    runID := fmt.Sprintf("%s-%d", job.ID, time.Now().Unix())

    // Create execution record
    exec := &JobExecution{
        JobID:     job.ID,
        RunID:     runID,
        Status:    StatusRunning,
        StartTime: time.Now(),
        Metrics:   JobMetrics{},
    }

    // Store execution
    e.execMutex.Lock()
    e.executions[runID] = exec
    e.execMutex.Unlock()

    // Log start
    e.logger.Printf("Starting job %s (run: %s)", job.Name, runID)

    // Execute with timing
    start := time.Now()
    err := e.executeStaged(ctx, job, exec)
    exec.TotalTime = time.Since(start)
    exec.EndTime = time.Now()

    // Update status
    if err != nil {
        exec.Status = StatusFailed
        exec.Error = err
        e.metrics.JobsFailed++
        e.logger.Printf("Job %s failed: %v", job.Name, err)
    } else {
        exec.Status = StatusSuccess
        e.metrics.JobsCompleted++
        e.logger.Printf("Job %s completed successfully in %v", job.Name, exec.TotalTime)
    }

    // Update metrics
    e.updateEngineMetrics(exec)

    return err
}

func (e *Engine) executeStaged(ctx context.Context, job ReportJob, exec *JobExecution) error {
    var err error

    // Stage 1: Data Loading
    start := time.Now()
    dsType := job.DataSource["type"]
    loader, ok := e.Loaders[dsType]
    if !ok {
        return fmt.Errorf("no loader for type %s", dsType)
    }

    data, err := loader.Load(ctx, job.DataSource)
    if err != nil {
        return fmt.Errorf("data loading failed: %w", err)
    }
    exec.Metrics.DataLoadTime = time.Since(start)
    exec.Metrics.DataRows = len(data.Rows)

    // Stage 2: Template Rendering
    start = time.Now()
    renderer := e.Renderers["markdown"]
    rendered, err := renderer.Render(ctx, job.TemplatePath, data)
    if err != nil {
        return fmt.Errorf("template rendering failed: %w", err)
    }
    exec.Metrics.RenderTime = time.Since(start)

    // Stage 3: Output Generation
    start = time.Now()
    var files []OutputFile
    for _, fmtName := range job.Outputs {
        outGen, ok := e.Outputs[fmtName]
        if !ok {
            return fmt.Errorf("no output generator for %s", fmtName)
        }
        f, err := outGen.Generate(ctx, rendered, fmtName)
        if err != nil {
            return fmt.Errorf("output generation failed: %w", err)
        }
        files = append(files, f)
        exec.Metrics.OutputSizeBytes += int64(len(f.Data))
    }
    exec.Metrics.OutputTime = time.Since(start)
    exec.OutputFiles = files

    // Stage 4: Delivery
    start = time.Now()
    for _, dCfg := range job.Delivery {
        dtype := dCfg["type"]
        adapter, ok := e.Deliveries[dtype]
        if !ok {
            return fmt.Errorf("no delivery adapter for %s", dtype)
        }
        if err := adapter.Deliver(ctx, dCfg, files); err != nil {
            return fmt.Errorf("delivery failed: %w", err)
        }
    }
    exec.Metrics.DeliveryTime = time.Since(start)

    return nil
}

// API endpoints for monitoring
func (e *Engine) GetJobStatus(runID string) (*JobExecution, bool) {
    e.execMutex.RLock()
    defer e.execMutex.RUnlock()
    exec, exists := e.executions[runID]
    return exec, exists
}

func (e *Engine) GetEngineMetrics() EngineMetrics {
    return *e.metrics
}

func (e *Engine) GetRecentExecutions(limit int) []*JobExecution {
    e.execMutex.RLock()
    defer e.execMutex.RUnlock()

    var executions []*JobExecution
    for _, exec := range e.executions {
        executions = append(executions, exec)
    }

    // Sort by start time (most recent first)
    sort.Slice(executions, func(i, j int) bool {
        return executions[i].StartTime.After(executions[j].StartTime)
    })

    if limit > 0 && len(executions) > limit {
        executions = executions[:limit]
    }

    return executions
}
```

### Health Check Endpoint

Create `cmd/cronyx-server/main.go`:

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "time"

    cronyx "github.com/your-username/cronyx/pkg/cronyx"
    // ... other imports
)

type Server struct {
    engine *cronyx.Engine
}

func main() {
    // Create and configure engine
    eng := cronyx.NewEngine(4)
    // ... register components

    server := &Server{engine: eng}

    // API routes
    http.HandleFunc("/health", server.healthCheck)
    http.HandleFunc("/metrics", server.getMetrics)
    http.HandleFunc("/executions", server.getExecutions)
    http.HandleFunc("/jobs", server.listJobs)
    http.HandleFunc("/jobs/trigger", server.triggerJob)

    // Start engine
    eng.Start()

    // Start HTTP server
    log.Println("Starting Cronyx server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
    metrics := s.engine.GetEngineMetrics()

    response := map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now(),
        "metrics": metrics,
        "uptime": time.Since(startTime),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (s *Server) getExecutions(w http.ResponseWriter, r *http.Request) {
    limitStr := r.URL.Query().Get("limit")
    limit := 50 // default
    if limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil {
            limit = l
        }
    }

    executions := s.engine.GetRecentExecutions(limit)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "executions": executions,
        "total": len(executions),
    })
}

func (s *Server) triggerJob(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        JobID string `json:"job_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Find and trigger job
    // Implementation depends on how jobs are stored

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status": "triggered",
        "job_id": req.JobID,
    })
}
```

---

## Configuration Management

### Environment-based Configuration

Create `config/config.go`:

```go
package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    Engine   EngineConfig
    Database DatabaseConfig
    SMTP     SMTPConfig
    S3       S3Config
    Logging  LoggingConfig
}

type EngineConfig struct {
    Workers     int
    QueueSize   int
    DefaultTimeout time.Duration
}

type DatabaseConfig struct {
    DSN             string
    MaxConnections  int
    ConnTimeout     time.Duration
}

type SMTPConfig struct {
    Host     string
    Port     string
    Username string
    Password string
    TLS      bool
}

type S3Config struct {
    Region          string
    AccessKeyID     string
    SecretAccessKey string
    DefaultBucket   string
}

type LoggingConfig struct {
    Level  string
    Format string
    Output string
}

func Load() *Config {
    return &Config{
        Engine: EngineConfig{
            Workers:        getEnvInt("CRONYX_WORKERS", 4),
            QueueSize:      getEnvInt("CRONYX_QUEUE_SIZE", 100),
            DefaultTimeout: getEnvDuration("CRONYX_DEFAULT_TIMEOUT", 30*time.Second),
        },
        Database: DatabaseConfig{
            DSN:            getEnv("DATABASE_URL", ""),
            MaxConnections: getEnvInt("DB_MAX_CONNECTIONS", 10),
            ConnTimeout:    getEnvDuration("DB_CONN_TIMEOUT", 5*time.Second),
        },
        SMTP: SMTPConfig{
            Host:     getEnv("SMTP_HOST", "localhost"),
            Port:     getEnv("SMTP_PORT", "587"),
            Username: getEnv("SMTP_USERNAME", ""),
            Password: getEnv("SMTP_PASSWORD", ""),
            TLS:      getEnvBool("SMTP_TLS", true),
        },
        S3: S3Config{
            Region:          getEnv("AWS_REGION", "us-east-1"),
            AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
            SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
            DefaultBucket:   getEnv("S3_DEFAULT_BUCKET", ""),
        },
        Logging: LoggingConfig{
            Level:  getEnv("LOG_LEVEL", "info"),
            Format: getEnv("LOG_FORMAT", "json"),
            Output: getEnv("LOG_OUTPUT", "stdout"),
        },
    }
}

// Helper functions
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        if boolValue, err := strconv.ParseBool(value); err == nil {
            return boolValue
        }
    }
    return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}
```

### Environment File Example

Create `.env`:

```bash
# Engine Configuration
CRONYX_WORKERS=8
CRONYX_QUEUE_SIZE=200
CRONYX_DEFAULT_TIMEOUT=60s

# Database Configuration
DATABASE_URL=postgres://user:password@localhost:5432/reports_db
DB_MAX_CONNECTIONS=20
DB_CONN_TIMEOUT=10s

# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_TLS=true

# AWS S3 Configuration
AWS_REGION=us-west-2
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
S3_DEFAULT_BUCKET=your-reports-bucket

# Logging Configuration
LOG_LEVEL=debug
LOG_FORMAT=json
LOG_OUTPUT=./logs/cronyx.log

# Application Configuration
HTTP_PORT=8080
METRICS_ENABLED=true
HEALTH_CHECK_INTERVAL=30s
```

---

## Troubleshooting

### Common Issues and Solutions

#### 1. **File Not Found Errors**

**Problem**: Template or data files not found

```
Error: failed to read template: open sample.md: no such file or directory
```

**Solutions**:

```bash
# Check current directory
pwd

# List files in current directory
ls -la

# Verify file paths in job configuration
# Make sure paths are relative to where you run the command

# Option 1: Run from correct directory
cd pkg/cronyx/embed-example
go run main.go

# Option 2: Use absolute paths
TemplatePath: "/full/path/to/sample.md"

# Option 3: Use relative paths from project root
TemplatePath: "pkg/cronyx/embed-example/sample.md"
```

#### 2. **Permission Denied Errors**

**Problem**: Cannot create output directory or files

```
Error: failed to create output directory: mkdir out: permission denied
```

**Solutions**:

```bash
# Check permissions
ls -la

# Create directory with proper permissions
mkdir -p out
chmod 755 out

# Run with proper permissions
sudo go run main.go  # Not recommended

# Or fix directory ownership
sudo chown -R $USER:$USER .
```

#### 3. **Module Import Errors**

**Problem**: Cannot import cronyx modules

```
Error: package github.com/your-username/cronyx/pkg/cronyx is not in GOROOT
```

**Solutions**:

```bash
# Initialize go module in project root
go mod init github.com/your-username/cronyx

# Update module paths in all Go files to match go.mod
# Make sure import paths match your actual module name

# Download dependencies
go mod tidy
go mod download

# Clean module cache if needed
go clean -modcache
```

#### 4. **Dependency Issues**

**Problem**: Missing or incompatible dependencies

```
Error: cannot find package "github.com/robfig/cron/v3"
```

**Solutions**:

```bash
# Install specific versions
go get github.com/robfig/cron/v3@v3.0.1
go get github.com/russross/blackfriday/v2@v2.1.0

# Check go.mod file
cat go.mod

# Clean and reinstall
go mod tidy
go mod download
```

#### 5. **Template Execution Errors**

**Problem**: Template fails to render

```
Error: failed to execute template: template: report:5:14: executing "report" at <.InvalidField>: can't evaluate field InvalidField
```

**Solutions**:

```go
// Debug template data
fmt.Printf("Template data: %+v\n", templateData)

// Check available fields in template
{{range $key, $value := .}}
Key: {{$key}}, Value: {{$value}}
{{end}}

// Use conditional checks in templates
{{if .OptionalField}}
{{.OptionalField}}
{{else}}
Field not available
{{end}}

// Verify data structure matches template expectations
```

#### 6. **Cron Expression Errors**

**Problem**: Invalid cron schedule

```
Error: failed to add cron job: expected 5 to 6 fields, found 4: "invalid expression"
```

**Solutions**:

```go
// Valid cron expressions (with seconds support)
"@every 10s"           // Every 10 seconds
"0 */5 * * * *"        // Every 5 minutes
"0 0 8 * * *"          // Every day at 8 AM
"0 0 8 * * MON-FRI"    // Every weekday at 8 AM
"0 0 8 1 * *"          // First day of every month at 8 AM

// Test cron expressions online: https://crontab.guru/
```

#### 7. **Memory Issues with Large Datasets**

**Problem**: Out of memory with large CSV files

```
Error: runtime: out of memory
```

**Solutions**:

```go
// Implement streaming CSV loader
func (c CSVLoader) LoadStreaming(ctx context.Context, cfg cronyx.DataSourceConfig) (<-chan map[string]interface{}, error) {
    // Return channel instead of loading all data at once
}

// Process data in batches
func (c CSVLoader) LoadBatch(ctx context.Context, cfg cronyx.DataSourceConfig, batchSize int) ([]cronyx.DataPayload, error) {
    // Split large files into smaller chunks
}

// Use database pagination
func (d DatabaseLoader) LoadPaginated(ctx context.Context, cfg cronyx.DataSourceConfig) (cronyx.DataPayload, error) {
    query := cfg["query"] + " LIMIT 1000 OFFSET " + cfg["offset"]
    // Process in pages
}
```

#### 8. **Concurrent Execution Issues**

**Problem**: Race conditions or deadlocks

```
Error: fatal error: concurrent map writes
```

**Solutions**:

```go
// Use proper synchronization
type ThreadSafeEngine struct {
    mu     sync.RWMutex
    engine *Engine
}

func (t *ThreadSafeEngine) GetStatus() EngineStatus {
    t.mu.RLock()
    defer t.mu.RUnlock()
    return t.engine.status
}

// Use channels for communication
type SafeJobQueue struct {
    jobs   chan ReportJob
    status chan JobStatus
}
```

### Debugging Techniques

#### 1. **Enable Verbose Logging**

```go
// Add to main.go
import "log"

func main() {
    // Enable detailed logging
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    // Or create custom logger
    logger := log.New(os.Stdout, "[CRONYX] ", log.LstdFlags)

    // Use in engine
    eng := cronyx.NewEngineWithLogger(4, logger)
}
```

#### 2. **Add Debug Prints**

```go
// In execute method
func (e *Engine) execute(ctx context.Context, job ReportJob) error {
    fmt.Printf("DEBUG: Starting job %s\n", job.ID)
    fmt.Printf("DEBUG: DataSource: %+v\n", job.DataSource)

    // ... execution steps with debug prints

    fmt.Printf("DEBUG: Job %s completed\n", job.ID)
    return nil
}
```

#### 3. **Test Individual Components**

```go
// Test CSV loader independently
func testCSVLoader() {
    loader := CSVLoader{}
    data, err := loader.Load(context.Background(), map[string]string{
        "path": "test-data.csv",
    })
    if err != nil {
        fmt.Printf("Loader error: %v\n", err)
        return
    }
    fmt.Printf("Loaded %d rows\n", len(data.Rows))
    for i, row := range data.Rows {
        if i >= 3 { break } // Show first 3 rows only
        fmt.Printf("Row %d: %+v\n", i, row)
    }
}

// Test template rendering
func testTemplateRenderer() {
    renderer := MarkdownRenderer{}
    data := cronyx.DataPayload{
        Rows: []map[string]interface{}{
            {"name": "Test Product", "value": "100", "category": "Test"},
        },
    }

    rendered, err := renderer.Render(context.Background(), "sample.md", data)
    if err != nil {
        fmt.Printf("Renderer error: %v\n", err)
        return
    }
    fmt.Printf("Rendered HTML: %s\n", rendered.HTML[:200]) // First 200 chars
}
```

#### 4. **Use Go's Built-in Profiler**

```go
// Add to main.go for performance debugging
import _ "net/http/pprof"

func main() {
    // Start profiler server
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

    // ... rest of main function
}

// Then access profiler at http://localhost:6060/debug/pprof/
```

#### 5. **Environment-specific Debug Configuration**

```go
// debug.go
//go:build debug
// +build debug

package main

import "fmt"

func debugPrint(format string, args ...interface{}) {
    fmt.Printf("[DEBUG] "+format+"\n", args...)
}

// release.go
//go:build !debug
// +build !debug

package main

func debugPrint(format string, args ...interface{}) {
    // No-op in release builds
}

// Usage in code:
debugPrint("Processing job %s with %d data rows", job.ID, len(data.Rows))

// Build with debug:
// go build -tags debug
```

### Performance Optimization

#### 1. **Optimize Data Loading**

```go
// Efficient CSV loader with streaming
type StreamingCSVLoader struct {
    BufferSize int
}

func (s StreamingCSVLoader) LoadStream(ctx context.Context, cfg cronyx.DataSourceConfig) (<-chan map[string]interface{}, <-chan error) {
    dataCh := make(chan map[string]interface{}, s.BufferSize)
    errCh := make(chan error, 1)

    go func() {
        defer close(dataCh)
        defer close(errCh)

        file, err := os.Open(cfg["path"])
        if err != nil {
            errCh <- err
            return
        }
        defer file.Close()

        reader := csv.NewReader(file)
        headers, err := reader.Read()
        if err != nil {
            errCh <- err
            return
        }

        for {
            record, err := reader.Read()
            if err == io.EOF {
                break
            }
            if err != nil {
                errCh <- err
                return
            }

            row := make(map[string]interface{})
            for i, header := range headers {
                if i < len(record) {
                    row[header] = record[i]
                }
            }

            select {
            case dataCh <- row:
            case <-ctx.Done():
                return
            }
        }
    }()

    return dataCh, errCh
}
```

#### 2. **Template Caching**

```go
// Cached renderer to avoid re-parsing templates
type CachedMarkdownRenderer struct {
    cache map[string]*template.Template
    mutex sync.RWMutex
}

func NewCachedMarkdownRenderer() *CachedMarkdownRenderer {
    return &CachedMarkdownRenderer{
        cache: make(map[string]*template.Template),
    }
}

func (c *CachedMarkdownRenderer) Render(ctx context.Context, tplPath string, data cronyx.DataPayload) (cronyx.RenderedDoc, error) {
    // Check cache first
    c.mutex.RLock()
    tmpl, exists := c.cache[tplPath]
    c.mutex.RUnlock()

    if !exists {
        // Load and parse template
        content, err := ioutil.ReadFile(tplPath)
        if err != nil {
            return cronyx.RenderedDoc{}, err
        }

        tmpl, err = template.New(filepath.Base(tplPath)).Parse(string(content))
        if err != nil {
            return cronyx.RenderedDoc{}, err
        }

        // Cache template
        c.mutex.Lock()
        c.cache[tplPath] = tmpl
        c.mutex.Unlock()
    }

    // Execute cached template
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, map[string]interface{}{
        "Rows": data.Rows,
        "Data": data.Rows,
    }); err != nil {
        return cronyx.RenderedDoc{}, err
    }

    html := string(bf.Run(buf.Bytes()))
    return cronyx.RenderedDoc{HTML: html, Content: buf.String()}, nil
}
```

#### 3. **Connection Pooling for Database Loaders**

```go
type PooledDatabaseLoader struct {
    pool *sql.DB
}

func NewPooledDatabaseLoader(dsn string) (*PooledDatabaseLoader, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    return &PooledDatabaseLoader{pool: db}, nil
}

func (p *PooledDatabaseLoader) Load(ctx context.Context, cfg cronyx.DataSourceConfig) (cronyx.DataPayload, error) {
    query := cfg["query"]

    rows, err := p.pool.QueryContext(ctx, query)
    if err != nil {
        return cronyx.DataPayload{}, err
    }
    defer rows.Close()

    // Process rows efficiently...
    return cronyx.DataPayload{Rows: result}, nil
}
```

---

## Production Deployment

### Docker Configuration

#### Dockerfile

```dockerfile
# Build stage
FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cronyx ./cmd/cronyx-server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/cronyx .
COPY --from=builder /app/templates ./templates/
COPY --from=builder /app/config ./config/

# Install wkhtmltopdf for PDF generation (optional)
RUN apk add --no-cache wkhtmltopdf

EXPOSE 8080

CMD ["./cronyx"]
```

#### docker-compose.yml

```yaml
version: "3.8"

services:
  cronyx:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://cronyx:password@postgres:5432/cronyx_db
      - SMTP_HOST=mailhog
      - SMTP_PORT=1025
      - LOG_LEVEL=info
      - CRONYX_WORKERS=4
    volumes:
      - ./data:/app/data
      - ./output:/app/output
      - ./templates:/app/templates
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

  postgres:
    image: postgres:13
    environment:
      - POSTGRES_DB=cronyx_db
      - POSTGRES_USER=cronyx
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"

  redis:
    image: redis:6-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  mailhog:
    image: mailhog/mailhog
    ports:
      - "1025:1025" # SMTP
      - "8025:8025" # Web UI

volumes:
  postgres_data:
  redis_data:
```

### Kubernetes Deployment

#### deployment.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cronyx
  labels:
    app: cronyx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: cronyx
  template:
    metadata:
      labels:
        app: cronyx
    spec:
      containers:
        - name: cronyx
          image: your-registry/cronyx:latest
          ports:
            - containerPort: 8080
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: cronyx-secrets
                  key: database-url
            - name: SMTP_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: cronyx-secrets
                  key: smtp-password
            - name: AWS_ACCESS_KEY_ID
              valueFrom:
                secretKeyRef:
                  name: cronyx-secrets
                  key: aws-access-key
            - name: AWS_SECRET_ACCESS_KEY
              valueFrom:
                secretKeyRef:
                  name: cronyx-secrets
                  key: aws-secret-key
            - name: LOG_LEVEL
              value: "info"
            - name: CRONYX_WORKERS
              value: "4"
          resources:
            requests:
              memory: "256Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
          volumeMounts:
            - name: templates
              mountPath: /app/templates
            - name: output
              mountPath: /app/output
      volumes:
        - name: templates
          configMap:
            name: cronyx-templates
        - name: output
          persistentVolumeClaim:
            claimName: cronyx-output-pvc

---
apiVersion: v1
kind: Service
metadata:
  name: cronyx-service
spec:
  selector:
    app: cronyx
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cronyx-templates
data:
  daily-report.md: |
    # Daily Report
    Generated: {{.Meta.timestamp}}
    Total Records: {{len .Rows}}

    {{range .Rows}}
    - {{.name}}: {{.value}}
    {{end}}

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: cronyx-output-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi

---
apiVersion: v1
kind: Secret
metadata:
  name: cronyx-secrets
type: Opaque
data:
  database-url: <base64-encoded-database-url>
  smtp-password: <base64-encoded-smtp-password>
  aws-access-key: <base64-encoded-aws-access-key>
  aws-secret-key: <base64-encoded-aws-secret-key>
```

### CI/CD Pipeline (GitHub Actions)

#### .github/workflows/deploy.yml

```yaml
name: Build and Deploy Cronyx

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.19

      - name: Run tests
        run: |
          go mod download
          go test -v ./...

      - name: Run linter
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest

  build-and-push:
    needs: test
    runs-on: ubuntu-latest
    if: github.event_name != 'pull_request'
    steps:
      - name: Checkout
        uses: actions/checkout@v3

      - name: Login to Container Registry
        uses: docker/login-action@v2
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v4
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=sha,prefix={{branch}}-
            type=raw,value=latest,enable={{is_default_branch}}

      - name: Build and push
        uses: docker/build-push-action@v3
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}

  deploy:
    needs: build-and-push
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - name: Deploy to Kubernetes
        run: |
          echo "Deploying to production..."
          # Add your deployment commands here
          # kubectl set image deployment/cronyx cronyx=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
```

### Monitoring and Observability in Production

#### Prometheus Metrics

```go
// metrics.go
package cronyx

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    jobsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cronyx_jobs_total",
            Help: "Total number of jobs processed",
        },
        []string{"job_id", "status"},
    )

    jobDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "cronyx_job_duration_seconds",
            Help: "Job execution duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"job_id"},
    )

    activeJobs = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "cronyx_active_jobs",
        Help: "Number of currently running jobs",
    })

    queueSize = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "cronyx_queue_size",
        Help: "Number of jobs in queue",
    })
)

func (e *Engine) recordJobMetrics(jobID string, duration time.Duration, status string) {
    jobsTotal.WithLabelValues(jobID, status).Inc()
    jobDuration.WithLabelValues(jobID).Observe(duration.Seconds())
}

func (e *Engine) updateQueueMetrics() {
    queueSize.Set(float64(len(e.jobQueue)))
}
```

#### Structured Logging

```go
// logger.go
package cronyx

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewLogger(level string, format string) (*zap.Logger, error) {
    var config zap.Config

    if format == "json" {
        config = zap.NewProductionConfig()
    } else {
        config = zap.NewDevelopmentConfig()
    }

    // Set log level
    switch level {
    case "debug":
        config.Level.SetLevel(zapcore.DebugLevel)
    case "info":
        config.Level.SetLevel(zapcore.InfoLevel)
    case "warn":
        config.Level.SetLevel(zapcore.WarnLevel)
    case "error":
        config.Level.SetLevel(zapcore.ErrorLevel)
    }

    return config.Build()
}

// Usage in engine
func (e *Engine) executeWithLogging(ctx context.Context, job ReportJob) error {
    logger := e.logger.With(
        zap.String("job_id", job.ID),
        zap.String("job_name", job.Name),
        zap.String("run_id", generateRunID()),
    )

    logger.Info("Starting job execution")
    start := time.Now()

    err := e.execute(ctx, job)
    duration := time.Since(start)

    if err != nil {
        logger.Error("Job execution failed",
            zap.Error(err),
            zap.Duration("duration", duration),
        )
    } else {
        logger.Info("Job execution completed",
            zap.Duration("duration", duration),
        )
    }

    return err
}
```

---

## Testing Strategy

### Unit Tests

#### Engine Tests

```go
// engine_test.go
package cronyx_test

import (
    "context"
    "testing"
    "time"

    "github.com/your-username/cronyx/pkg/cronyx"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock implementations
type MockDataLoader struct {
    mock.Mock
}

func (m *MockDataLoader) Load(ctx context.Context, cfg cronyx.DataSourceConfig) (cronyx.DataPayload, error) {
    args := m.Called(ctx, cfg)
    return args.Get(0).(cronyx.DataPayload), args.Error(1)
}

type MockRenderer struct {
    mock.Mock
}

func (m *MockRenderer) Render(ctx context.Context, tplPath string, data cronyx.DataPayload) (cronyx.RenderedDoc, error) {
    args := m.Called(ctx, tplPath, data)
    return args.Get(0).(cronyx.RenderedDoc), args.Error(1)
}

func TestEngineExecution(t *testing.T) {
    // Setup
    eng := cronyx.NewEngine(1)

    mockLoader := &MockDataLoader{}
    mockRenderer := &MockRenderer{}

    eng.RegisterLoader("test", mockLoader)
    eng.RegisterRenderer("markdown", mockRenderer)

    // Mock expectations
    expectedData := cronyx.DataPayload{
        Rows: []map[string]interface{}{
            {"name": "test", "value": "123"},
        },
    }

    mockLoader.On("Load", mock.Anything, mock.Anything).Return(expectedData, nil)
    mockRenderer.On("Render", mock.Anything, mock.Anything, expectedData).Return(
        cronyx.RenderedDoc{HTML: "<h1>Test</h1>", Content: "# Test"},
        nil,
    )

    // Test job
    job := cronyx.ReportJob{
        ID:           "test-job",
        Name:         "test",
        TemplatePath: "test.md",
        DataSource:   cronyx.DataSourceConfig{"type": "test"},
        Outputs:      []string{},
        Delivery:     []cronyx.DeliveryConfig{},
        Timeout:      10 * time.Second,
    }

    // Execute
    ctx := context.Background()
    err := eng.TestExecute(ctx, job)

    // Verify
    assert.NoError(t, err)
    mockLoader.AssertExpectations(t)
    mockRenderer.AssertExpectations(t)
}

func TestJobQueue(t *testing.T) {
    eng := cronyx.NewEngine(2)
    eng.Start()
    defer eng.Stop()

    // Enqueue multiple jobs
    for i := 0; i < 5; i++ {
        job := cronyx.ReportJob{
            ID:      fmt.Sprintf("job-%d", i),
            Name:    fmt.Sprintf("test-job-%d", i),
            Timeout: 5 * time.Second,
        }
        eng.Enqueue(job)
    }

    // Wait for processing
    time.Sleep(1 * time.Second)

    // Verify queue is processed
    assert.Equal(t, 0, len(eng.GetJobQueue()))
}
```

### Integration Tests

#### Full Pipeline Tests

```go
// integration_test.go
package cronyx_test

import (
    "context"
    "io/ioutil"
    "os"
    "path/filepath"
    "testing"
    "time"

    "github.com/your-username/cronyx/pkg/cronyx"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFullPipeline(t *testing.T) {
    // Setup test environment
    tempDir, err := ioutil.TempDir("", "cronyx-test")
    require.NoError(t, err)
    defer os.RemoveAll(tempDir)

    // Create test data
    csvPath := filepath.Join(tempDir, "test-data.csv")
    csvContent := `name,value,category
Product A,100,Electronics
Product B,200,Clothing`

    err = ioutil.WriteFile(csvPath, []byte(csvContent), 0644)
    require.NoError(t, err)

    // Create test template
    tplPath := filepath.Join(tempDir, "test-template.md")
    tplContent := `# Test Report
Total items: {{len .Rows}}
{{range .Rows}}
- {{.name}}: ${{.value}}
{{end}}`

    err = ioutil.WriteFile(tplPath, []byte(tplContent), 0644)
    require.NoError(t, err)

    // Setup engine
    eng := cronyx.NewEngine(1)
    eng.RegisterLoader("csv", cronyx.CSVLoader{})
    eng.RegisterRenderer("markdown", cronyx.MarkdownRenderer{})
    eng.RegisterOutput("html", cronyx.FileOutputGenerator{
        OutDir: filepath.Join(tempDir, "output"),
    })
    eng.RegisterDelivery("console", cronyx.ConsoleDelivery{})

    // Create job
    job := cronyx.ReportJob{
        ID:           "integration-test",
        Name:         "integration-test",
        TemplatePath: tplPath,
        DataSource: cronyx.DataSourceConfig{
            "type": "csv",
            "path": csvPath,
        },
        Outputs:  []string{"html"},
        Delivery: []cronyx.DeliveryConfig{{"type": "console"}},
        Timeout:  30 * time.Second,
    }

    // Execute
    ctx := context.Background()
    err = eng.TestExecute(ctx, job)
    require.NoError(t, err)

    // Verify output file was created
    outputDir := filepath.Join(tempDir, "output")
    files, err := ioutil.ReadDir(outputDir)
    require.NoError(t, err)
    assert.Len(t, files, 1)

    // Verify output content
    outputFile := filepath.Join(outputDir, files[0].Name())
    content, err := ioutil.ReadFile(outputFile)
    require.NoError(t, err)

    htmlContent := string(content)
    assert.Contains(t, htmlContent, "Test Report")
    assert.Contains(t, htmlContent, "Product A")
    assert.Contains(t, htmlContent, "Product B")
    assert.Contains(t, htmlContent, "Total items: 2")
}
```

### Benchmark Tests

```go
// benchmark_test.go
package cronyx_test

import (
    "context"
    "testing"
    "time"

    "github.com/your-username/cronyx/pkg/cronyx"
)

func BenchmarkJobExecution(b *testing.B) {
    eng := setupBenchmarkEngine()
    job := createBenchmarkJob()
    ctx := context.Background()

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        err := eng.TestExecute(ctx, job)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkConcurrentJobs(b *testing.B) {
    eng := cronyx.NewEngine(4)
    eng.Start()
    defer eng.Stop()

    job := createBenchmarkJob()

    b.ResetTimer()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            eng.Enqueue(job)
        }
    })
}

func BenchmarkLargeDataset(b *testing.B) {
    // Test with datasets of different sizes
    sizes := []int{100, 1000, 10000, 100000}

    for _, size := range sizes {
        b.Run(fmt.Sprintf("rows-%d", size), func(b *testing.B) {
            data := generateTestData(size)
            renderer := cronyx.MarkdownRenderer{}

            b.ResetTimer()

            for i := 0; i < b.N; i++ {
                _, err := renderer.Render(context.Background(), "benchmark.md", data)
                if err != nil {
                    b.Fatal(err)
                }
            }
        })
    }
}
```

---

## Security Considerations

### Input Validation and Sanitization

```go
// security.go
package cronyx

import (
    "fmt"
    "path/filepath"
    "regexp"
    "strings"
)

// Validator interface for input validation
type Validator interface {
    Validate() error
}

// Implement validation for job configuration
func (j ReportJob) Validate() error {
    if j.ID == "" {
        return fmt.Errorf("job ID cannot be empty")
    }

    // Validate job ID format (alphanumeric, dash, underscore only)
    validID := regexp.MustCompile(`^[a-zA-Z0-9_-]+# Cronyx: Complete Project Documentation

## Table of Contents
1. [Project Overview](#project-overview)
2. [Architecture & Design](#architecture--design)
3. [Project Structure](#project-structure)
4. [Core Interfaces](#core-interfaces)
5. [Implementation Details](#implementation-details)
6. [Configuration & Setup](#configuration--setup)
7. [Running the Application](#running-the-application)
8. [Code Walkthrough](#code-walkthrough)
9. [Advanced Usage](#advanced-usage)
10. [Troubleshooting](#troubleshooting)

---

## Project Overview

**Cronyx** is a Go-based automated report generation engine that follows a plugin-based architecture. It's designed to:

- **Load data** from various sources (CSV, databases, APIs)
- **Render templates** with that data (Markdown, HTML)
- **Generate outputs** in different formats (HTML, PDF, Excel)
- **Deliver reports** through multiple channels (Console, Email, Slack)
- **Schedule execution** using cron expressions

### Key Features
- 🔄 **Cron-based scheduling** - Run reports automatically
- 🔌 **Plugin architecture** - Easily extensible components
- 🎯 **Template-driven** - Markdown templates with Go templating
- 📊 **Multiple outputs** - HTML, PDF, Excel support
- 🚀 **Concurrent execution** - Worker pool for parallel processing
- 📦 **Modular design** - Clean separation of concerns

---

## Architecture & Design

Cronyx uses a **pipeline-based architecture** with four main stages:

```

┌─────────────┐ ┌──────────────┐ ┌─────────────┐ ┌──────────────┐
│ Data │───▶│ Template │───▶│ Output │───▶│ Delivery │
│ Loading │ │ Rendering │ │ Generation │ │ Distribution │
└─────────────┘ └──────────────┘ └─────────────┘ └──────────────┘

```

### Core Components

1. **Engine**: Central orchestrator managing the pipeline
2. **DataLoaders**: Extract data from various sources
3. **TemplateRenderers**: Process templates with data
4. **OutputGenerators**: Create final output files
5. **DeliveryAdapters**: Distribute generated reports

### Design Patterns Used

- **Strategy Pattern**: Pluggable components for each pipeline stage
- **Registry Pattern**: Component registration and lookup
- **Worker Pool Pattern**: Concurrent job processing
- **Template Pattern**: Consistent execution pipeline

---

## Project Structure

```

Cronyx/
├── cmd/ # Command-line applications
├── docs/ # Documentation
│ └── Forme.md # Project documentation
├── go.mod # Go module definition
├── go.sum # Dependency checksums
├── internal/ # Private application code
├── pkg/ # Public library code
│ ├── cronyx/ # Main package
│ │ ├── delivery/ # Delivery adapters
│ │ │ └── console.go # Console output delivery
│ │ ├── embed-example/ # Example implementation
│ │ │ ├── data.csv # Sample data file
│ │ │ ├── main.go # Main application
│ │ │ └── sample.md # Template file
│ │ ├── engine.go # Core engine implementation
│ │ ├── interfaces.go # Interface definitions
│ │ ├── job.go # Job data structures
│ │ ├── loaders/ # Data loaders
│ │ │ └── csv.go # CSV data loader
│ │ ├── outputs/ # Output generators
│ │ │ └── filegen.go # File output generator
│ │ └── renderers/ # Template renderers
│ │ └── markdown.go # Markdown renderer
│ └── prompt/ # Prompt templates
│ └── prompt.txt # System prompts
└── README.md # Project overview

````

### Directory Purpose Explanation

- **`cmd/`**: Entry points for different applications (CLI tools, servers)
- **`internal/`**: Code that's internal to the project (not importable by others)
- **`pkg/`**: Public library code that can be imported by other projects
- **`docs/`**: Project documentation and specifications
- **`embed-example/`**: Complete working example with sample data

---

## Core Interfaces

The power of Cronyx lies in its interface-based design. Let's examine each interface:

### 1. DataLoader Interface

```go
type DataLoader interface {
    Load(ctx context.Context, cfg DataSourceConfig) (DataPayload, error)
}
````

**Purpose**: Abstracts data loading from any source
**Implementations**: CSVLoader (more can be added: DatabaseLoader, APILoader, etc.)

### 2. TemplateRenderer Interface

```go
type TemplateRenderer interface {
    Render(ctx context.Context, tplPath string, data DataPayload) (RenderedDoc, error)
}
```

**Purpose**: Converts templates + data into rendered documents
**Implementations**: MarkdownRenderer (HTMLRenderer, PDFRenderer can be added)

### 3. OutputGenerator Interface

```go
type OutputGenerator interface {
    Generate(ctx context.Context, rendered RenderedDoc, format string) (OutputFile, error)
}
```

**Purpose**: Creates final output files in various formats
**Implementations**: FileOutputGenerator (S3Generator, DatabaseGenerator can be added)

### 4. DeliveryAdapter Interface

```go
type DeliveryAdapter interface {
    Deliver(ctx context.Context, target DeliveryConfig, files []OutputFile) error
}
```

**Purpose**: Distributes generated files to target destinations
**Implementations**: ConsoleDelivery (EmailDelivery, SlackDelivery can be added)

### Supporting Data Structures

```go
// Generic data container
type DataPayload struct {
    Rows []map[string]interface{}  // Structured data
    Raw  []byte                    // Raw data for binary sources
}

// Rendered document with metadata
type RenderedDoc struct {
    HTML    string                 // HTML content
    Content string                 // Raw content
    Meta    map[string]interface{} // Metadata
}

// Output file representation
type OutputFile struct {
    Name string  // File name
    Path string  // File path or URI
    Data []byte  // File content (optional)
}
```

---

## Implementation Details

### Engine Implementation (`engine.go`)

The Engine is the heart of Cronyx. Let's break down its implementation:

```go
type Engine struct {
    cronSched  *cron.Cron                    // Cron scheduler
    Loaders    map[string]DataLoader         // Registered data loaders
    Renderers  map[string]TemplateRenderer   // Registered renderers
    Outputs    map[string]OutputGenerator    // Registered output generators
    Deliveries map[string]DeliveryAdapter    // Registered delivery adapters

    jobQueue   chan ReportJob                // Job queue for workers
    workers    int                           // Number of worker goroutines
    stopCh     chan struct{}                 // Stop signal channel
}
```

#### Key Methods Explained

**1. NewEngine(workers int)**

```go
func NewEngine(workers int) *Engine {
    e := &Engine{
        cronSched:  cron.New(cron.WithSeconds()),  // Enable second-precision cron
        Loaders:    map[string]DataLoader{},       // Initialize empty registries
        Renderers:  map[string]TemplateRenderer{},
        Outputs:    map[string]OutputGenerator{},
        Deliveries: map[string]DeliveryAdapter{},
        jobQueue:   make(chan ReportJob, 100),     // Buffered channel for jobs
        workers:    workers,
        stopCh:     make(chan struct{}),           // Unbuffered stop channel
    }
    return e
}
```

**2. Registration Methods**

```go
func (e *Engine) RegisterLoader(name string, d DataLoader) {
    e.Loaders[name] = d  // Store loader by name for lookup
}
```

These methods populate the component registries for runtime lookup.

**3. Worker Pool Implementation**

```go
func (e *Engine) workerLoop(id int) {
    for {
        select {
        case job := <-e.jobQueue:                    // Receive job from queue
            ctx, cancel := context.WithTimeout(      // Create timeout context
                context.Background(),
                job.Timeout
            )
            _ = e.execute(ctx, job)                   // Execute job
            cancel()                                  // Clean up context
        case <-e.stopCh:                            // Receive stop signal
            return                                    // Exit worker
        }
    }
}
```

**4. Job Execution Pipeline**

```go
func (e *Engine) execute(ctx context.Context, job ReportJob) error {
    // 1. Load data
    dsType := job.DataSource["type"]
    loader, ok := e.Loaders[dsType]
    if !ok {
        return fmt.Errorf("no loader for type %s", dsType)
    }
    data, err := loader.Load(ctx, job.DataSource)
    if err != nil {
        return err
    }

    // 2. Render template
    renderer := e.Renderers["markdown"]  // Could be dynamic based on template
    rendered, err := renderer.Render(ctx, job.TemplatePath, data)
    if err != nil {
        return err
    }

    // 3. Generate outputs
    var files []OutputFile
    for _, fmtName := range job.Outputs {
        outGen, ok := e.Outputs[fmtName]
        if !ok {
            return fmt.Errorf("no output generator for %s", fmtName)
        }
        f, err := outGen.Generate(ctx, rendered, fmtName)
        if err != nil {
            return err
        }
        files = append(files, f)
    }

    // 4. Deliver files
    for _, dCfg := range job.Delivery {
        dtype := dCfg["type"]
        adapter, ok := e.Deliveries[dtype]
        if !ok {
            return fmt.Errorf("no delivery adapter for %s", dtype)
        }
        if err := adapter.Deliver(ctx, dCfg, files); err != nil {
            return err
        }
    }

    return nil
}
```

### CSV Loader Implementation (`loaders/csv.go`)

```go
type CSVLoader struct{}

func (c CSVLoader) Load(ctx context.Context, cfg cronyx.DataSourceConfig) (cronyx.DataPayload, error) {
    path := cfg["path"]                    // Extract file path from config
    f, err := os.Open(path)               // Open CSV file
    if err != nil {
        return cronyx.DataPayload{}, err
    }
    defer f.Close()                       // Ensure file is closed

    r := csv.NewReader(f)                 // Create CSV reader
    headers, err := r.Read()              // Read header row
    if err != nil {
        return cronyx.DataPayload{}, err
    }

    var rows []map[string]interface{}     // Store rows as key-value maps
    for {
        record, err := r.Read()           // Read each data row
        if err == io.EOF {
            break                         // End of file reached
        }
        if err != nil {
            return cronyx.DataPayload{}, err
        }

        row := map[string]interface{}{}   // Create row map
        for i, h := range headers {       // Map each value to its header
            row[h] = record[i]
        }
        rows = append(rows, row)          // Add to results
    }

    return cronyx.DataPayload{Rows: rows}, nil
}
```

### Markdown Renderer Implementation (`renderers/markdown.go`)

```go
func (MarkdownRenderer) Render(ctx context.Context, tplPath string, data cronyx.DataPayload) (cronyx.RenderedDoc, error) {
    // 1. Read template file
    b, err := ioutil.ReadFile(tplPath)
    if err != nil {
        return cronyx.RenderedDoc{}, fmt.Errorf("failed to read template: %w", err)
    }

    // 2. Create Go template
    tmpl, err := template.New("report").Parse(string(b))
    if err != nil {
        return cronyx.RenderedDoc{}, fmt.Errorf("failed to parse template: %w", err)
    }

    // 3. Prepare template data
    templateData := map[string]interface{}{
        "Rows": data.Rows,              // Make data available as .Rows
        "Data": data.Rows,              // Alias for convenience
    }

    // 4. Execute template with data
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, templateData); err != nil {
        return cronyx.RenderedDoc{}, fmt.Errorf("failed to execute template: %w", err)
    }

    // 5. Convert markdown to HTML using blackfriday
    html := string(bf.Run(buf.Bytes()))

    return cronyx.RenderedDoc{
        HTML:    html,
        Content: buf.String(),  // Keep original markdown
        Meta: map[string]interface{}{
            "source":     tplPath,
            "rows_count": len(data.Rows),
        },
    }, nil
}
```

### File Output Generator (`outputs/filegen.go`)

```go
func (g FileOutputGenerator) Generate(ctx context.Context, r cronyx.RenderedDoc, format string) (cronyx.OutputFile, error) {
    // 1. Ensure output directory exists
    if err := os.MkdirAll(g.OutDir, 0755); err != nil {
        return cronyx.OutputFile{}, fmt.Errorf("failed to create output directory: %w", err)
    }

    // 2. Generate unique filename with timestamp
    timestamp := time.Now().Format("20060102_150405")
    filename := fmt.Sprintf("report_%s.%s", timestamp, format)
    outPath := filepath.Join(g.OutDir, filename)

    // 3. Determine content based on format
    var data []byte
    switch format {
    case "html":
        data = []byte(r.HTML)       // Use HTML content
    case "pdf":
        data = []byte(r.HTML)       // Placeholder - would use PDF generator
    case "md":
        data = []byte(r.Content)    // Use original markdown
    default:
        data = []byte(r.HTML)       // Default to HTML
    }

    // 4. Write file to disk
    if err := ioutil.WriteFile(outPath, data, 0644); err != nil {
        return cronyx.OutputFile{}, fmt.Errorf("failed to write output file: %w", err)
    }

    return cronyx.OutputFile{
        Name: filename,
        Path: outPath,
        Data: data,
    }, nil
}
```

---

## Configuration & Setup

### Prerequisites

1. **Go 1.19+** installed on your system
2. **Git** for version control
3. **Text editor** or IDE

### Initial Setup

#### 1. Create Project Directory

```bash
mkdir cronyx-project
cd cronyx-project
```

#### 2. Initialize Go Module

```bash
go mod init github.com/your-username/cronyx
```

#### 3. Install Dependencies

```bash
go get github.com/robfig/cron/v3@v3.0.1
go get github.com/russross/blackfriday/v2@v2.1.0
```

Your `go.mod` should look like:

```go
module github.com/your-username/cronyx

go 1.19

require (
    github.com/robfig/cron/v3 v3.0.1
    github.com/russross/blackfriday/v2 v2.1.0
)
```

#### 4. Create Directory Structure

```bash
mkdir -p pkg/cronyx/{delivery,loaders,outputs,renderers,embed-example}
mkdir -p cmd docs internal
```

### Configuration Files

#### Sample Data (`pkg/cronyx/embed-example/data.csv`)

```csv
name,value,category,date
Product A,150,Electronics,2024-01-15
Product B,200,Clothing,2024-01-16
Product C,75,Books,2024-01-17
Product D,300,Electronics,2024-01-18
Product E,120,Home & Garden,2024-01-19
```

#### Sample Template (`pkg/cronyx/embed-example/sample.md`)

```markdown
# Daily Sales Report

**Generated**: {{.Meta.timestamp}}  
**Total Products**: {{len .Rows}}

## Product Inventory

{{range .Rows}}

### {{.name}}

- **Value**: ${{.value}}
- **Category**: {{.category}}
- **Date**: {{.date}}

---

{{end}}

## Summary by Category

{{$categories := .Categories}}
{{range $categories}}

- **{{.Name}}**: {{.Count}} items, ${{.Total}} total value
  {{end}}

## Report Details

- Generated by Cronyx Report Engine
- Data source: CSV file
- Template: Markdown with Go templating
- Output format: HTML

_This is an automated report. Please do not reply to this document._
```

---

## Running the Application

### Step-by-Step Execution Guide

#### 1. Prepare the Environment

```bash
# Navigate to project root
cd /path/to/your/cronyx-project

# Verify Go installation
go version

# Check dependencies
go mod tidy
go mod verify
```

#### 2. Create the Main Application

Create `pkg/cronyx/embed-example/main.go`:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    cronyx "github.com/your-username/cronyx/pkg/cronyx"
    deliver "github.com/your-username/cronyx/pkg/cronyx/delivery"
    loader "github.com/your-username/cronyx/pkg/cronyx/loaders"
    generate "github.com/your-username/cronyx/pkg/cronyx/outputs"
    render "github.com/your-username/cronyx/pkg/cronyx/renderers"
)

func main() {
    fmt.Println("🚀 Starting Cronyx Report Engine...")

    // 1. Create engine with 4 worker threads
    eng := cronyx.NewEngine(4)
    fmt.Println("✅ Engine created with 4 workers")

    // 2. Register all components
    eng.RegisterLoader("csv", loader.CSVLoader{})
    eng.RegisterRenderer("markdown", render.MarkdownRenderer{})
    eng.RegisterOutput("html", generate.FileOutputGenerator{OutDir: "./out"})
    eng.RegisterDelivery("console", deliver.ConsoleDelivery{})
    fmt.Println("✅ All components registered")

    // 3. Define the report job
    job := cronyx.ReportJob{
        ID:           "daily-sales-001",
        Name:         "daily-sales-report",
        TemplatePath: "sample.md",
        DataSource: cronyx.DataSourceConfig{
            "type": "csv",
            "path": "data.csv",
        },
        Outputs:  []string{"html"},
        Schedule: "@every 15s",  // Run every 15 seconds for demo
        Delivery: []cronyx.DeliveryConfig{
            {"type": "console"},
        },
        Timeout: 30 * time.Second,
    }

    // 4. Test single execution first
    fmt.Println("\n🧪 Testing single job execution...")
    ctx := context.Background()
    if err := eng.TestExecute(ctx, job); err != nil {
        log.Fatalf("❌ Job execution failed: %v", err)
    }
    fmt.Println("✅ Single execution successful!")

    // 5. Add job to cron scheduler
    if err := eng.AddCronJob(job); err != nil {
        log.Fatalf("❌ Failed to add cron job: %v", err)
    }
    fmt.Println("✅ Job added to scheduler")

    // 6. Start the engine
    fmt.Println("\n🔄 Starting scheduled execution...")
    eng.Start()

    // 7. Run for 2 minutes then stop
    fmt.Println("⏰ Running for 2 minutes. Press Ctrl+C to stop early.")
    time.Sleep(2 * time.Minute)

    // 8. Graceful shutdown
    fmt.Println("\n🛑 Shutting down...")
    eng.Stop()
    fmt.Println("👋 Cronyx stopped. Goodbye!")
}
```

#### 3. Execute the Application

**Option A: Run from embed-example directory**

```bash
cd pkg/cronyx/embed-example
go run main.go
```

**Option B: Run from project root**

```bash
# Update paths in main.go first:
# TemplatePath: "pkg/cronyx/embed-example/sample.md"
# DataSource path: "pkg/cronyx/embed-example/data.csv"
# OutDir: "./out"

go run pkg/cronyx/embed-example/main.go
```

#### 4. Expected Output

```
🚀 Starting Cronyx Report Engine...
✅ Engine created with 4 workers
✅ All components registered

🧪 Testing single job execution...
=== Executing job: daily-sales-report ===
✅ Single execution successful!
✅ Job added to scheduler

🔄 Starting scheduled execution...
⏰ Running for 2 minutes. Press Ctrl+C to stop early.
[Cronyx] Delivered file: report_20240315_143022.html (./out/report_20240315_143022.html)
[Cronyx] Delivered file: report_20240315_143037.html (./out/report_20240315_143037.html)
[Cronyx] Delivered file: report_20240315_143052.html (./out/report_20240315_143052.html)
...

🛑 Shutting down...
👋 Cronyx stopped. Goodbye!
```

#### 5. Verify Output Files

```bash
# Check generated files
ls -la out/

# View a generated report
open out/report_20240315_143022.html
# or
cat out/report_20240315_143022.html
```

---

## Code Walkthrough

Let's trace through a complete execution cycle:

### 1. Application Startup

```go
// main.go starts here
func main() {
    // Creates new engine instance
    eng := cronyx.NewEngine(4)
```

**What happens**:

- Creates cron scheduler with second precision
- Initializes empty component registries
- Creates buffered job queue (capacity: 100)
- Sets up worker pool with 4 goroutines

### 2. Component Registration

```go
eng.RegisterLoader("csv", loader.CSVLoader{})
eng.RegisterRenderer("markdown", render.MarkdownRenderer{})
eng.RegisterOutput("html", generate.FileOutputGenerator{OutDir: "./out"})
eng.RegisterDelivery("console", deliver.ConsoleDelivery{})
```

**What happens**:

- Each component is stored in its respective registry map
- Components are looked up by string keys during execution
- This enables runtime plugin selection and extensibility

### 3. Job Definition

```go
job := cronyx.ReportJob{
    ID:           "daily-sales-001",
    Name:         "daily-sales-report",
    TemplatePath: "sample.md",
    DataSource:   cronyx.DataSourceConfig{"type": "csv", "path": "data.csv"},
    Outputs:      []string{"html"},
    Schedule:     "@every 15s",
    Delivery:     []cronyx.DeliveryConfig{{"type": "console"}},
    Timeout:      30 * time.Second,
}
```

**What happens**:

- Job configuration is created with all necessary parameters
- DataSource config tells the engine which loader to use
- Outputs array specifies which formats to generate
- Delivery config specifies where to send the results

### 4. Engine Startup & Worker Pool

```go
eng.Start()
```

**What happens**:

- Starts 4 worker goroutines running `workerLoop()`
- Each worker listens on the job queue channel
- Starts the cron scheduler
- Workers are now ready to process jobs

### 5. Job Execution Pipeline

When a job is triggered (either by cron or manual enqueue):

**Step 1: Data Loading**

```go
// Engine looks up loader by type
dsType := job.DataSource["type"]  // "csv"
loader, ok := e.Loaders[dsType]   // Gets CSVLoader instance

// CSVLoader.Load() is called
data, err := loader.Load(ctx, job.DataSource)
```

**CSVLoader execution**:

1. Opens file at path specified in DataSource
2. Reads CSV headers
3. Converts each row to `map[string]interface{}`
4. Returns DataPayload with all rows

**Step 2: Template Rendering**

```go
renderer := e.Renderers["markdown"]  // Gets MarkdownRenderer
rendered, err := renderer.Render(ctx, job.TemplatePath, data)
```

**MarkdownRenderer execution**:

1. Reads template file (`sample.md`)
2. Parses it as Go template
3. Executes template with data (replaces `{{.Rows}}`, etc.)
4. Converts resulting markdown to HTML using blackfriday
5. Returns RenderedDoc with HTML and metadata

**Step 3: Output Generation**

```go
for _, fmtName := range job.Outputs {  // ["html"]
    outGen, ok := e.Outputs[fmtName]   // Gets FileOutputGenerator
    f, err := outGen.Generate(ctx, rendered, fmtName)
}
```

**FileOutputGenerator execution**:

1. Creates output directory if it doesn't exist
2. Generates timestamped filename
3. Writes HTML content to file
4. Returns OutputFile with path and metadata

**Step 4: Delivery**

```go
for _, dCfg := range job.Delivery {    // [{"type": "console"}]
    adapter, ok := e.Deliveries[dtype] // Gets ConsoleDelivery
    adapter.Deliver(ctx, dCfg, files)
}
```

**ConsoleDelivery execution**:

1. Prints file information to console
2. Could be extended to send emails, post to Slack, etc.

### 6. Concurrent Execution

Multiple workers can process jobs simultaneously:

```
Worker 1: Job A (Data Loading)
Worker 2: Job B (Template Rendering)
Worker 3: Job C (Output Generation)
Worker 4: Idle (waiting for jobs)
```

This enables high throughput for multiple concurrent reports.

---

## Advanced Usage

### Custom Data Loader Example

Create `loaders/database.go`:

```go
package loaders

import (
    "context"
    "database/sql"
    "github.com/your-username/cronyx/pkg/cronyx"
    _ "github.com/lib/pq"  // PostgreSQL driver
)

type DatabaseLoader struct{}

func (d DatabaseLoader) Load(ctx context.Context, cfg cronyx.DataSourceConfig) (cronyx.DataPayload, error) {
    // Get connection details from config
    dsn := cfg["dsn"]
    query := cfg["query"]

    // Connect to database
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return cronyx.DataPayload{}, err
    }
    defer db.Close()

    // Execute query
    rows, err := db.QueryContext(ctx, query)
    if err != nil {
        return cronyx.DataPayload{}, err
    }
    defer rows.Close()

    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        return cronyx.DataPayload{}, err
    }

    // Process rows
    var result []map[string]interface{}
    for rows.Next() {
        // Create slice of interface{} for Scan
        values := make([]interface{}, len(columns))
        pointers := make([]interface{}, len(columns))
        for i := range values {
            pointers[i] = &values[i]
        }

        // Scan row
        if err := rows.Scan(pointers...); err != nil {
            return cronyx.DataPayload{}, err
        }

        // Convert to map
        row := make(map[string]interface{})
        for i, col := range columns {
            row[col] = values[i]
        }
        result = append(result, row)
    }

    return cronyx.DataPayload{Rows: result}, nil
}
```

**Usage**:

```go
// Register the custom loader
eng.RegisterLoader("database", loader.DatabaseLoader{})

// Use in job definition
job := cronyx.ReportJob{
    // ... other fields
    DataSource: cronyx.DataSourceConfig{
        "type": "database",
        "dsn":  "postgres://user:pass@localhost/db",
        "query": "SELECT * FROM sales WHERE date >= CURRENT_DATE - INTERVAL '7 days'",
    },
}
```

### Email Delivery Adapter

Create `delivery/email.go`:

```go
package delivery

import (
    "context"
    "fmt"
    "net/smtp"
    "github.com/your-username/cronyx/pkg/cronyx"
)

type EmailDelivery struct{}

func (e EmailDelivery) Deliver(ctx context.Context, target cronyx.DeliveryConfig, files []cronyx.OutputFile) error {
    // Extract email configuration
    smtpHost := target["smtp_host"]
    smtpPort := target["smtp_port"]
    username := target["username"]
    password := target["password"]
    to := target["to"]
    subject := target["subject"]

    // Create email message
    body := "Please find attached reports:\n\n"
    for _, file := range files {
        body += fmt.Sprintf("- %s\n", file.Name)
    }

    msg := []byte("To: " + to + "\r\n" +
        "Subject: " + subject + "\r\n" +
        "\r\n" +
        body + "\r\n")

    // Send email
    auth := smtp.PlainAuth("", username, password, smtpHost)
    addr := smtpHost + ":" + smtpPort
    return smtp.SendMail(addr, auth, username, []string{to}, msg)
}
```

### PDF Output Generator

Create `outputs/pdf.go`:

```go
package outputs

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "time"
    "github.com/your-username/cronyx/pkg/cronyx"
)

type PDFGenerator struct {
    OutDir string
}

func (p PDFGenerator) Generate(ctx context.Context, r cronyx.RenderedDoc, format string) (cronyx.OutputFile, error) {
    // Create temp HTML file
    timestamp := time.Now().Format("20060102_150405")
    htmlFile := filepath.Join(p.OutDir, fmt.Sprintf("temp_%s.html", timestamp))
    pdfFile := filepath.Join(p.OutDir, fmt.Sprintf("report_%s.pdf", timestamp))

    // Write HTML to temp file
    if err := ioutil.WriteFile(htmlFile, []byte(r.HTML), 0644); err != nil {
        return cronyx.OutputFile{}, err
    }
    defer os.Remove(htmlFile) // Clean up temp file

    // Convert HTML to PDF using wkhtmltopdf
    cmd := exec.CommandContext(ctx, "wkhtmltopdf", htmlFile, pdfFile)
    if err := cmd.Run(); err != nil {
        return cronyx.OutputFile{}, fmt.Errorf("PDF generation failed: %w", err)
    }

    // Read generated PDF
    data, err := ioutil.ReadFile(pdfFile)
    if err != nil {
        return cronyx.OutputFile{}, err
    }

    return cronyx.OutputFile{
        Name: filepath.Base(pdfFile),
        Path: pdfFile,
        Data: data,
    }, nil
}
```

### Complex Template with Functions

Create advanced template `templates/advanced-report.md`:

```markdown
# {{.Title | title}}

**Generated**: {{.Meta.timestamp | formatDate}}  
**Total Records**: {{len .Rows}}  
**Report ID**: {{.ReportID}}

## Executive Summary

{{if gt (len .Rows) 0}}

- **Highest Value**: ${{.Rows | maxValue "value"}}
- **Lowest Value**: ${{.Rows | minValue "value"}}
- **Average Value**: ${{.Rows | avgValue "value" | printf "%.2f"}}
- **Total Value**: ${{.Rows | sumValue "value"}}
  {{else}}
  No data available for this period.
  {{end}}

## Category Breakdown

{{range $category, $items := (.Rows | groupBy "category")}}

### {{$category}}

{{range $items}}

- **{{.name}}**: ${{.value}} ({{.date | formatDate}})
{{end}}
**Category Total**: ${{$items | sumValue "value"}}

---

{{end}}

## Detailed Listings

{{range $index, $row := .Rows}}
{{add $index 1}}. **{{$row.name}}**

- Value: ${{$row.value}}
- Category: {{$row.category}}
- Date: {{$row.date | formatDate}}
  {{if $row.description}}
- Description: {{$row.description}}
  {{end}}
  {{end}}

## Charts and Visualizations

{{if .Meta.include_charts}}

### Sales by Category
```

{{range (.Rows | groupBy "category")}}
{{.Key}}: {{repeat "█" (div .Total 50)}} (${{.Total}})
{{end}}

```

### Trend Analysis
```

{{range (.Rows | groupBy "date" | sort)}}
{{.Key}}: {{repeat "▓" (div (.Items | len) 2)}} ({{.Items | len}} items)
{{end}}

```
{{end}}

---
*Report generated by Cronyx Engine v1.0*
*Data as of {{.Meta.data_timestamp | formatDate}}*
```

### Enhanced Renderer with Template Functions

Create `renderers/advanced_markdown.go`:

```go
package renderers

import (
    "context"
    "fmt"
    "io/ioutil"
    "sort"
    "strconv"
    "strings"
    "text/template"
    "time"
    "bytes"

    "github.com/your-username/cronyx/pkg/cronyx"
    bf "github.com/russross/blackfriday/v2"
)

type AdvancedMarkdownRenderer struct{}

func (AdvancedMarkdownRenderer) Render(ctx context.Context, tplPath string, data cronyx.DataPayload) (cronyx.RenderedDoc, error) {
    // Read template
    b, err := ioutil.ReadFile(tplPath)
    if err != nil {
        return cronyx.RenderedDoc{}, fmt.Errorf("failed to read template: %w", err)
    }

    // Create template with custom functions
    tmpl := template.New("advanced-report").Funcs(template.FuncMap{
        "formatDate": formatDate,
        "title":      strings.Title,
        "add":        func(a, b int) int { return a + b },
        "maxValue":   func(rows []map[string]interface{}, field string) float64 {
            return aggregateFloat(rows, field, func(values []float64) float64 {
                if len(values) == 0 { return 0 }
                max := values[0]
                for _, v := range values[1:] {
                    if v > max { max = v }
                }
                return max
            })
        },
        "minValue": func(rows []map[string]interface{}, field string) float64 {
            return aggregateFloat(rows, field, func(values []float64) float64 {
                if len(values) == 0 { return 0 }
                min := values[0]
                for _, v := range values[1:] {
                    if v < min { min = v }
                }
                return min
            })
        },
        "avgValue": func(rows []map[string]interface{}, field string) float64 {
            return aggregateFloat(rows, field, func(values []float64) float64 {
                if len(values) == 0 { return 0 }
                sum := 0.0
                for _, v := range values {
                    sum += v
                }
                return sum / float64(len(values))
            })
        },
        "sumValue": func(rows []map[string]interface{}, field string) float64 {
            return aggregateFloat(rows, field, func(values []float64) float64 {
                sum := 0.0
                for _, v := range values {
                    sum += v
                }
                return sum
            })
        },
        "groupBy": func(rows []map[string]interface{}, field string) map[string][]map[string]interface{} {
            groups := make(map[string][]map[string]interface{})
            for _, row := range rows {
                if val, ok := row[field]; ok {
                    key := fmt.Sprintf("%v", val)
                    groups[key] = append(groups[key], row)
                }
            }
            return groups
        },
        "repeat": func(str string, count int) string {
            if count < 0 { count = 0 }
            if count > 100 { count = 100 } // Safety limit
            return strings.Repeat(str, count)
        },
        "div": func(a, b int) int {
            if b == 0 { return 0 }
            return a / b
        },
        "printf": fmt.Sprintf,
    })

    // Parse template
    tmpl, err = tmpl.Parse(string(b))
    if err != nil {
        return cronyx.RenderedDoc{}, fmt.Errorf("failed to parse template: %w", err)
    }

    // Prepare enhanced template data
    templateData := map[string]interface{}{
        "Rows": data.Rows,
        "Data": data.Rows,
        "Title": "Sales Report",
        "ReportID": fmt.Sprintf("RPT-%d", time.Now().Unix()),
        "Meta": map[string]interface{}{
            "timestamp": time.Now(),
            "data_timestamp": time.Now().Add(-1 * time.Hour),
            "include_charts": true,
            "rows_count": len(data.Rows),
        },
    }

    // Execute template
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, templateData); err != nil {
        return cronyx.RenderedDoc{}, fmt.Errorf("failed to execute template: %w", err)
    }

    // Convert to HTML
    html := string(bf.Run(buf.Bytes()))

    return cronyx.RenderedDoc{
        HTML:    html,
        Content: buf.String(),
        Meta: templateData["Meta"].(map[string]interface{}),
    }, nil
}

// Helper functions
func formatDate(t time.Time) string {
    return t.Format("January 2, 2006 at 3:04 PM")
}

func aggregateFloat(rows []map[string]interface{}, field string, aggFunc func([]float64) float64) float64 {
    var values []float64
    for _, row := range rows {
        if val, ok := row[field]; ok {
            if fval, err := strconv.ParseFloat(fmt.Sprintf("%v", val), 64); err == nil {
                values = append(values, fval)
            }
        }
    }
    return aggFunc(values)
}
```

### Multi-Format Job Configuration

```go
// Complex job with multiple outputs and deliveries
complexJob := cronyx.ReportJob{
    ID:           "quarterly-sales-001",
    Name:         "quarterly-sales-analysis",
    TemplatePath: "templates/quarterly-report.md",
    DataSource: cronyx.DataSourceConfig{
        "type": "database",
        "dsn":  "postgres://user:pass@localhost/sales_db",
        "query": `
            SELECT
                p.name,
                SUM(s.amount) as value,
                p.category,
                DATE(s.sale_date) as date,
                p.description
            FROM sales s
            JOIN products p ON s.product_id = p.id
            WHERE s.sale_date >= CURRENT_DATE - INTERVAL '3 months'
            GROUP BY p.id, p.name, p.category, DATE(s.sale_date)
            ORDER BY s.sale_date DESC
        `,
    },
    Outputs: []string{"html", "pdf", "excel"},
    Schedule: "0 0 8 1 * *", // First day of every month at 8 AM
    Delivery: []cronyx.DeliveryConfig{
        {
            "type": "email",
            "to": "executives@company.com",
            "subject": "Quarterly Sales Report",
            "smtp_host": "smtp.company.com",
            "smtp_port": "587",
            "username": "reports@company.com",
            "password": "secret123",
        },
        {
            "type": "slack",
            "webhook": "https://hooks.slack.com/services/...",
            "channel": "#sales-reports",
            "message": "📊 New quarterly sales report available",
        },
        {
            "type": "s3",
            "bucket": "company-reports",
            "prefix": "quarterly-sales/",
            "region": "us-west-2",
        },
    },
    Timeout: 5 * time.Minute,
    Labels: map[string]string{
        "department": "sales",
        "priority":   "high",
        "type":       "quarterly",
    },
}
```

---

## Monitoring and Observability

### Job Status Tracking

Add to `job.go`:

```go
type JobStatus string

const (
    StatusPending   JobStatus = "pending"
    StatusRunning   JobStatus = "running"
    StatusSuccess   JobStatus = "success"
    StatusFailed    JobStatus = "failed"
    StatusTimeout   JobStatus = "timeout"
)

type JobExecution struct {
    JobID     string
    RunID     string
    Status    JobStatus
    StartTime time.Time
    EndTime   time.Time
    Error     error
    OutputFiles []OutputFile
    Metrics   JobMetrics
}

type JobMetrics struct {
    DataLoadTime     time.Duration
    RenderTime       time.Duration
    OutputTime       time.Duration
    DeliveryTime     time.Duration
    TotalTime        time.Duration
    DataRows         int
    OutputSizeBytes  int64
}
```

### Enhanced Engine with Monitoring

Add to `engine.go`:

```go
type Engine struct {
    // ... existing fields

    // Monitoring
    executions  map[string]*JobExecution
    execMutex   sync.RWMutex
    metrics     *EngineMetrics
    logger      *log.Logger
}

type EngineMetrics struct {
    JobsCompleted   int64
    JobsFailed      int64
    TotalRunTime    time.Duration
    AverageRunTime  time.Duration
}

func (e *Engine) executeWithMonitoring(ctx context.Context, job ReportJob) error {
    runID := fmt.Sprintf("%s-%d", job.ID, time.Now().Unix())

    // Create execution record
    exec := &JobExecution{
        JobID:     job.ID,
        RunID:     runID,
        Status:    StatusRunning,
        StartTime: time.Now(),
        Metrics:   JobMetrics{},
    }

    // Store execution
    e.execMutex.Lock()
    e.executions[runID] = exec
    e.execMutex.Unlock()

    // Log start
    e.logger.Printf("Starting job %s (run: %s)", job.Name, runID)

    // Execute with timing
    start := time.Now()
    err := e.executeStaged(ctx, job, exec)
    exec.TotalTime = time.Since(start)
    exec.EndTime = time.Now()

    // Update status
    if err != nil {
        exec.Status = StatusFailed
        exec.Error = err
        e.metrics.JobsFailed++
        e.logger.Printf("Job %s failed: %v", job.Name, err)
    } else {
        exec.Status = StatusSuccess
        e.metrics.JobsCompleted++
        e.logger.Printf("Job %s completed successfully in %v", job.Name, exec.TotalTime)
    }

    // Update metrics
    e.updateEngineMetrics(exec)

    return err
}

func (e *Engine) executeStaged(ctx context.Context, job ReportJob, exec *JobExecution) error {
    var err error

    // Stage 1: Data Loading
    start := time.Now()
    dsType := job.DataSource["type"]
    loader, ok := e.Loaders[dsType]
    if !ok {
        return fmt.Errorf("no loader for type %s", dsType)
    }

    data, err := loader.Load(ctx, job.DataSource)
    if err != nil {
        return fmt.Errorf("data loading failed: %w", err)
    }
    exec.Metrics.DataLoadTime = time.Since(start)
    exec.Metrics.DataRows = len(data.Rows)

    // Stage 2: Template Rendering
    start = time.Now()
    renderer := e.Renderers["markdown"]
    rendered, err := renderer.Render(ctx, job.TemplatePath, data)
    if err != nil {
        return fmt.Errorf("template rendering failed: %w", err)
    }
    exec.Metrics.RenderTime = time.Since(start)

    // Stage 3: Output Generation
    start = time.Now()
    var files []OutputFile
    for _, fmtName := range job.Outputs {
        outGen, ok := e.Outputs[fmtName]
        if !ok {
            return fmt.Errorf("no output generator for %s", fmtName)
        }
        f, err := outGen.Generate(ctx, rendered, fmtName)
        if err != nil {
            return fmt.Errorf("output generation failed: %w", err)
        }
        files = append(files, f)
        exec.Metrics.OutputSizeBytes += int64(len(f.Data))
    }
    exec.Metrics.OutputTime = time.Since(start)
    exec.OutputFiles = files

    // Stage 4: Delivery
    start = time.Now()
    for _, dCfg := range job.Delivery {
        dtype := dCfg["type"]
        adapter, ok := e.Deliveries[dtype]
        if !ok {
            return fmt.Errorf("no delivery adapter for %s", dtype)
        }
        if err := adapter.Deliver(ctx, dCfg, files); err != nil {
            return fmt.Errorf("delivery failed: %w", err)
        }
    }
    exec.Metrics.DeliveryTime = time.Since(start)

    return nil
}

// API endpoints for monitoring
func (e *Engine) GetJobStatus(runID string) (*JobExecution, bool) {
    e.execMutex.RLock()
    defer e.execMutex.RUnlock()
    exec, exists := e.executions[runID]
    return exec, exists
}

func (e *Engine) GetEngineMetrics() EngineMetrics {
    return *e.metrics
}

func (e *Engine) GetRecentExecutions(limit int) []*JobExecution {
    e.execMutex.RLock()
    defer e.execMutex.RUnlock()

    var executions []*JobExecution
    for _, exec := range e.executions {
        executions = append(executions, exec)
    }

    // Sort by start time (most recent first)
    sort.Slice(executions, func(i, j int) bool {
        return executions[i].StartTime.After(executions[j].StartTime)
    })

    if limit > 0 && len(executions) > limit {
        executions = executions[:limit]
    }

    return executions
}
```

### Health Check Endpoint

Create `cmd/cronyx-server/main.go`:

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "time"

    cronyx "github.com/your-username/cronyx/pkg/cronyx"
    // ... other imports
)

type Server struct {
    engine *cronyx.Engine
}

func main() {
    // Create and configure engine
    eng := cronyx.NewEngine(4)
    // ... register components

    server := &Server{engine: eng}

    // API routes
    http.HandleFunc("/health", server.healthCheck)
    http.HandleFunc("/metrics", server.getMetrics)
    http.HandleFunc("/executions", server.getExecutions)
    http.HandleFunc("/jobs", server.listJobs)
    http.HandleFunc("/jobs/trigger", server.triggerJob)

    // Start engine
    eng.Start()

    // Start HTTP server
    log.Println("Starting Cronyx server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
    metrics := s.engine.GetEngineMetrics()

    response := map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now(),
        "metrics": metrics,
        "uptime": time.Since(startTime),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (s *Server) getExecutions(w http.ResponseWriter, r *http.Request) {
    limitStr := r.URL.Query().Get("limit")
    limit := 50 // default
    if limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil {
            limit = l
        }
    }

    executions := s.engine.GetRecentExecutions(limit)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "executions": executions,
        "total": len(executions),
    })
}

func (s *Server) triggerJob(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        JobID string `json:"job_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Find and trigger job
    // Implementation depends on how jobs are stored

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status": "triggered",
        "job_id": req.JobID,
    })
}
```

---

## Configuration Management

### Environment-based Configuration

Create `config/config.go`:

```go
package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    Engine   EngineConfig
    Database DatabaseConfig
    SMTP     SMTPConfig
    S3       S3Config
    Logging  LoggingConfig
}

type EngineConfig struct {
    Workers     int
    QueueSize   int
    DefaultTimeout time.Duration
}

type DatabaseConfig struct {
    DSN             string
    MaxConnections  int
    ConnTimeout     time.Duration
}

type SMTPConfig struct {
    Host     string
    Port     string
    Username string
    Password string
    TLS      bool
}

type S3Config struct {
    Region          string
    AccessKeyID     string
    SecretAccessKey string
    DefaultBucket   string
}

type LoggingConfig struct {
    Level  string
    Format string
    Output string
}

func Load() *Config {
    return &Config{
        Engine: EngineConfig{
            Workers:        getEnvInt("CRONYX_WORKERS", 4),
            QueueSize:      getEnvInt("CRONYX_QUEUE_SIZE", 100),
            DefaultTimeout: getEnvDuration("CRONYX_DEFAULT_TIMEOUT", 30*time.Second),
        },
        Database: DatabaseConfig{
            DSN:            getEnv("DATABASE_URL", ""),
            MaxConnections: getEnvInt("DB_MAX_CONNECTIONS", 10),
            ConnTimeout:    getEnvDuration("DB_CONN_TIMEOUT", 5*time.Second),
        },
        SMTP: SMTPConfig{
            Host:     getEnv("SMTP_HOST", "localhost"),
            Port:     getEnv("SMTP_PORT", "587"),
            Username: getEnv("SMTP_USERNAME", ""),
            Password: getEnv("SMTP_PASSWORD", ""),
            TLS:      getEnvBool("SMTP_TLS", true),
        },
        S3: S3Config{
            Region:          getEnv("AWS_REGION", "us-east-1"),
            AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
            SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
            DefaultBucket:   getEnv("S3_DEFAULT_BUCKET", ""),
        },
        Logging: LoggingConfig{
            Level:  getEnv("LOG_LEVEL", "info"),
            Format: getEnv("LOG_FORMAT", "json"),
            Output: getEnv("LOG_OUTPUT", "stdout"),
        },
    }
}

// Helper functions
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        if boolValue, err := strconv.ParseBool(value); err == nil {
            return boolValue
        }
    }
    return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}
```

### Environment File Example

Create `.env`:

```bash
# Engine Configuration
CRONYX_WORKERS=8
CRONYX_QUEUE_SIZE=200
CRONYX_DEFAULT_TIMEOUT=60s

# Database Configuration
DATABASE_URL=postgres://user:password@localhost:5432/reports_db
DB_MAX_CONNECTIONS=20
DB_CONN_TIMEOUT=10s

# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_TLS=true

# AWS S3 Configuration
AWS_REGION=us-west-2
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
S3_DEFAULT_BUCKET=your-reports-bucket

# Logging Configuration
LOG_LEVEL=debug
LOG_FORMAT=json
LOG_OUTPUT=./logs/cronyx.log

# Application Configuration
HTTP_PORT=8080
METRICS_ENABLED=true
HEALTH_CHECK_INTERVAL=30s
```

---

## Troubleshooting

### Common Issues and Solutions

#### 1. **File Not Found Errors**

**Problem**: Template or data files not found

```
Error: failed to read template: open sample.md: no such file or directory
```

**Solutions**:

```bash
# Check current directory
pwd

# List files in current directory
ls -la

# Verify file paths in job configuration
# Make sure paths are relative to where you run the command

# Option 1: Run from correct directory
cd pkg/cronyx/embed-example
go run main.go

# Option 2: Use absolute paths
TemplatePath: "/full/path/to/sample.md"

# Option 3: Use relative paths from project root
TemplatePath: "pkg/cronyx/embed-example/sample.md"
```

#### 2. **Permission Denied Errors**

**Problem**: Cannot create output directory or files

```
Error: failed to create output directory: mkdir out: permission denied
```

**Solutions**:

```bash
# Check permissions
ls -la

# Create directory with proper permissions
mkdir -p out
chmod 755 out

# Run with proper permissions
sudo go run main.go  # Not recommended

# Or fix directory ownership
sudo chown -R $USER:$USER .
```

#### 3. **Module Import Errors**

**Problem**: Cannot import cronyx modules

```
Error: package github.com/your-username/cronyx/pkg/cronyx is not in GOROOT
```

**Solutions**:

```bash
# Initialize go module in project root
go mod init github.com/your-username/cronyx

# Update module paths in all Go files to match go.mod
# Make sure import paths match your actual module name

# Download dependencies
go mod tidy
go mod download

# Clean module cache if needed
go clean -modcache
```

#### 4. **Dependency Issues**

**Problem**: Missing or incompatible dependencies

```
Error: cannot find package "github.com/robfig/cron/v3"
```

**Solutions**:

```bash
# Install specific versions
go get github.com/robfig/cron/v3@v3.0.1
go get github.com/russross/blackfriday/v2@v2.1.0

# Check go.mod file
cat go.mod

# Clean and reinstall
go mod tidy
go mod download
```

#### 5. **Template Execution Errors**

**Problem**: Template fails to render

```
Error: failed to execute template: template: report:5:14: executing "report" at <.InvalidField>: can't evaluate field InvalidField
```

**Solutions**:

```go
// Debug template data
fmt.Printf("Template data: %+v\n", templateData)

// Check available fields in template
{{range $key, $value := .}}
Key: {{$key}}, Value: {{$value}}
{{end}}

// Use conditional checks in templates
{{if .OptionalField}}
{{.OptionalField}}
{{else}}
Field not available
{{end}}

// Verify data structure matches template expectations
```

#### 6. **Cron Expression Errors**

**Problem**: Invalid cron schedule

```
Error: failed to add cron job: expected 5 to 6 fields, found 4: "invalid expression"
```

**Solutions**:

```go
// Valid cron expressions (with seconds support)
"@every 10s"           // Every 10 seconds
"0 */5 * * * *"        // Every 5 minutes
"0 0 8 * * *"          // Every day at 8 AM
"0 0 8 * * MON-FRI"    // Every weekday at 8 AM
"0 0 8 1 * *"          // First day of every month at 8 AM

// Test cron expressions online: https://crontab.guru/
```

#### 7. **Memory Issues with Large Datasets**

**Problem**: Out of memory with large CSV files

```
Error: runtime: out of memory
```

**Solutions**:

```go
// Implement streaming CSV loader
func (c CSVLoader) LoadStreaming(ctx context.Context, cfg cronyx.DataSourceConfig) (<-chan map[string]interface{}, error) {
    // Return channel instead of loading all data at once
}

// Process data in batches
func (c CSVLoader) LoadBatch(ctx context.Context, cfg cronyx.DataSourceConfig, batchSize int) ([]cronyx.DataPayload, error) {
    // Split large files into smaller chunks
}

// Use database pagination
func (d DatabaseLoader) LoadPaginated(ctx context.Context, cfg cronyx.DataSourceConfig) (cronyx.DataPayload, error) {
    query := cfg["query"] + " LIMIT 1000 OFFSET " + cfg["offset"]
    // Process in pages
}
```

#### 8. **Concurrent Execution Issues**

**Problem**: Race conditions or deadlocks

```
Error: fatal error: concurrent map writes
```

**Solutions**:

```go
// Use proper synchronization
type ThreadSafeEngine struct {
    mu     sync.RWMutex
    engine *Engine
}

func (t *ThreadSafeEngine) GetStatus() EngineStatus {
    t.mu.RLock()
    defer t.mu.RUnlock()
    return t.engine.status
}

// Use channels for communication
type SafeJobQueue struct {
    jobs   chan ReportJob
    status chan JobStatus
}
```

### Debugging Techniques

#### 1. **Enable Verbose Logging**

```go
// Add to main.go
import "log"

func main() {
    // Enable detailed logging
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    // Or create custom logger
    logger := log.New(os.Stdout, "[CRONYX] ", log.LstdFlags)

    // Use in engine
    eng := cronyx.NewEngineWithLogger(4, logger)
}
```

#### 2. **Add Debug Prints**

```go
// In execute method
func (e *Engine) execute(ctx context.Context, job ReportJob) error {
    fmt.Printf("DEBUG: Starting job %s\n", job.ID)
    fmt.Printf("DEBUG: DataSource: %+v\n", job.DataSource)

    // ... execution steps with debug prints

    fmt.Printf("DEBUG: Job %s completed\n", job.ID)
    return nil
}
```

#### 3. **Test Individual Components**

```go
// Test CSV loader independently
func testCSVLoader() {
    loader := CSVLoader{}
    data, err := loader.Load(context.Background(), map[string]string{
        "path": "test-data.csv",
    })
)
    if !validID.MatchString(j.ID) {
        return fmt.Errorf("invalid job ID format: %s", j.ID)
    }

    // Validate template path (prevent directory traversal)
    if err := validatePath(j.TemplatePath); err != nil {
        return fmt.Errorf("invalid template path: %w", err)
    }

    // Validate timeout
    if j.Timeout <= 0 || j.Timeout > 24*time.Hour {
        return fmt.Errorf("invalid timeout: must be between 1s and 24h")
    }

    // Validate data source
    if j.DataSource["type"] == "" {
        return fmt.Errorf("data source type cannot be empty")
    }

    return nil
}

func validatePath(path string) error {
    // Clean the path
    cleanPath := filepath.Clean(path)

    // Check for directory traversal attempts
    if strings.Contains(cleanPath, "..") {
        return fmt.Errorf("path traversal not allowed")
    }

    // Ensure path is relative
    if filepath.IsAbs(cleanPath) {
        return fmt.Errorf("absolute paths not allowed")
    }

    return nil
}

// Sanitize template content to prevent XSS
func sanitizeHTML(html string) string {
    // Use a proper HTML sanitizer like bluemonday in production
    // This is a simplified example

    // Remove script tags
    scriptRegex := regexp.MustCompile(`<script[^>]*>.*?</script>`)
    html = scriptRegex.ReplaceAllString(html, "")

    // Remove on* event handlers
    eventRegex := regexp.MustCompile(`on\w+\s*=\s*"[^"]*"`)
    html = eventRegex.ReplaceAllString(html, "")

    return html
}
```

### Access Control and Authentication

```go
// auth.go
package cronyx

import (
    "context"
    "crypto/subtle"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v4"
)

type AuthConfig struct {
    JWTSecret   string
    APIKeys     map[string]string // API key -> user mapping
    AdminUsers  []string
}

type Claims struct {
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}

// Middleware for API authentication
func (a AuthConfig) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check for API key first
        apiKey := r.Header.Get("X-API-Key")
        if apiKey != "" {
            if user, valid := a.validateAPIKey(apiKey); valid {
                ctx := context.WithValue(r.Context(), "user", user)
                next(w, r.WithContext(ctx))
                return
            }
        }

        // Check for JWT token
        authHeader := r.Header.Get("Authorization")
        if strings.HasPrefix(authHeader, "Bearer ") {
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            if claims, valid := a.validateJWT(tokenString); valid {
                ctx := context.WithValue(r.Context(), "user", claims.Username)
                ctx = context.WithValue(ctx, "roles", claims.Roles)
                next(w, r.WithContext(ctx))
                return
            }
        }

        http.Error(w, "Unauthorized", http.StatusUnauthorized)
    }
}

func (a AuthConfig) validateAPIKey(key string) (string, bool) {
    for apiKey, user := range a.APIKeys {
        if subtle.ConstantTimeCompare([]byte(key), []byte(apiKey)) == 1 {
            return user, true
        }
    }
    return "", false
}

func (a AuthConfig) validateJWT(tokenString string) (*Claims, bool) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.
```
