package cronyx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Engine struct {
	cronSched  *cron.Cron
	loaders    map[string]DataLoader
	renderers  map[string]TemplateRenderer
	outputs    map[string]OutputGenerator
	deliveries map[string]DeliveryAdapter

	jobQueue chan ReportJob
	workers  int
	stopCh   chan struct{}
}

func NewEngine(workers int) *Engine {
	e := &Engine{
		cronSched:  cron.New(cron.WithSeconds()),
		loaders:    map[string]DataLoader{},
		renderers:  map[string]TemplateRenderer{},
		outputs:    map[string]OutputGenerator{},
		deliveries: map[string]DeliveryAdapter{},
		jobQueue:   make(chan ReportJob, 100),
		workers:    workers,
		stopCh:     make(chan struct{}),
	}
	return e
}

// Register helpers:
func (e *Engine) RegisterLoader(name string, d DataLoader) {
	e.loaders[name] = d
}
func (e *Engine) RegisterRenderer(name string, r TemplateRenderer) {
	e.renderers[name] = r
}
func (e *Engine) RegisterOutput(name string, o OutputGenerator) {
	e.outputs[name] = o
}
func (e *Engine) RegisterDelivery(name string, d DeliveryAdapter) {
	e.deliveries[name] = d
}

// Start scheduler/workers
func (e *Engine) Start() {
	for i := 0; i < e.workers; i++ {
		go e.workerLoop(i)
	}
	e.cronSched.Start()
}

// Stop
func (e *Engine) Stop() {
	ctx := e.cronSched.Stop()
	select {
	case <-ctx.Done():
	case <-time.After(2 * time.Second):
	}
	close(e.stopCh)
}

func (e *Engine) AddCronJob(job ReportJob) error {
	if job.Schedule == "" {
		return errors.New("empty schedule")
	}
	// enqueue on schedule
	_, err := e.cronSched.AddFunc(job.Schedule, func() {
		e.jobQueue <- job
	})
	return err
}

func (e *Engine) Enqueue(job ReportJob) {
	e.jobQueue <- job
}

func (e *Engine) workerLoop(id int) {
	for {
		select {
		case job := <-e.jobQueue:
			ctx, cancel := context.WithTimeout(context.Background(), job.Timeout)
			_ = e.execute(ctx, job)
			cancel()
		case <-e.stopCh:
			return
		}
	}
}

func (e *Engine) execute(ctx context.Context, job ReportJob) error {
	// 1. find loader (based on type in DataSource)
	dsType := job.DataSource["type"]
	loader, ok := e.loaders[dsType]
	if !ok {
		return fmt.Errorf("no loader for type %s", dsType)
	}

	// 2. load
	data, err := loader.Load(ctx, job.DataSource)
	if err != nil {
		return err
	}

	// 3. render (pick renderer from template type; we'll assume "markdown")
	renderer := e.renderers["markdown"]
	rendered, err := renderer.Render(ctx, job.TemplatePath, data)
	if err != nil {
		return err
	}

	// 4. outputs
	var files []OutputFile
	for _, fmtName := range job.Outputs {
		outGen, ok := e.outputs[fmtName]
		if !ok {
			return fmt.Errorf("no output generator for %s", fmtName)
		}
		f, err := outGen.Generate(ctx, rendered, fmtName)
		if err != nil {
			return err
		}
		files = append(files, f)
	}

	// 5. delivery
	for _, dCfg := range job.Delivery {
		dtype := dCfg["type"]
		adapter, ok := e.deliveries[dtype]
		if !ok {
			return fmt.Errorf("no delivery adapter for %s", dtype)
		}
		if err := adapter.Deliver(ctx, dCfg, files); err != nil {
			return err
		}
	}

	return nil
}
