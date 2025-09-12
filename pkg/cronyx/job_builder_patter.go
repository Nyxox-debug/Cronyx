package cronyx

import (
	"crypto/rand"
	"fmt"
	"time"
)

// JobBuilder provides a fluent API for creating jobs
type JobBuilder struct {
	job ReportJob
	err error
}

// NewJob creates a new job builder
func NewJob(name string) *JobBuilder {
	return &JobBuilder{
		job: ReportJob{
			ID:      generateJobID(),
			Name:    name,
			Timeout: 30 * time.Second, // default timeout
			Labels:  make(map[string]string),
		},
	}
}

// WithID sets a custom job ID
func (jb *JobBuilder) WithID(id string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.ID = id
	return jb
}

// WithTemplate sets the template path
func (jb *JobBuilder) WithTemplate(path string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.TemplatePath = path
	return jb
}

// WithCSVData configures CSV data source
func (jb *JobBuilder) WithCSVData(path string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	if jb.job.DataSource == nil {
		jb.job.DataSource = make(DataSourceConfig)
	}
	jb.job.DataSource["type"] = "csv"
	jb.job.DataSource["path"] = path
	return jb
}

// WithJSONData configures JSON data source
func (jb *JobBuilder) WithJSONData(path string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	if jb.job.DataSource == nil {
		jb.job.DataSource = make(DataSourceConfig)
	}
	jb.job.DataSource["type"] = "json"
	jb.job.DataSource["path"] = path
	return jb
}

// WithDatabaseData configures database data source
func (jb *JobBuilder) WithDatabaseData(dsn, query string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	if jb.job.DataSource == nil {
		jb.job.DataSource = make(DataSourceConfig)
	}
	jb.job.DataSource["type"] = "database"
	jb.job.DataSource["dsn"] = dsn
	jb.job.DataSource["query"] = query
	return jb
}

// WithCustomDataSource allows any data source configuration
func (jb *JobBuilder) WithCustomDataSource(config DataSourceConfig) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.DataSource = config
	return jb
}

// OutputHTML adds HTML output
func (jb *JobBuilder) OutputHTML() *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Outputs = append(jb.job.Outputs, "html")
	return jb
}

// OutputPDF adds PDF output
func (jb *JobBuilder) OutputPDF() *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Outputs = append(jb.job.Outputs, "pdf")
	return jb
}

// OutputExcel adds Excel output
func (jb *JobBuilder) OutputExcel() *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Outputs = append(jb.job.Outputs, "xlsx")
	return jb
}

// OutputCSV adds CSV output
func (jb *JobBuilder) OutputCSV() *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Outputs = append(jb.job.Outputs, "csv")
	return jb
}

// WithOutputs sets multiple output formats at once
func (jb *JobBuilder) WithOutputs(formats ...string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Outputs = append(jb.job.Outputs, formats...)
	return jb
}

// DeliverToConsole adds console delivery
func (jb *JobBuilder) DeliverToConsole() *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Delivery = append(jb.job.Delivery, DeliveryConfig{"type": "console"})
	return jb
}

// DeliverToEmail adds email delivery
func (jb *JobBuilder) DeliverToEmail(to, subject string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Delivery = append(jb.job.Delivery, DeliveryConfig{
		"type":    "email",
		"to":      to,
		"subject": subject,
	})
	return jb
}

// DeliverToSlack adds Slack delivery
func (jb *JobBuilder) DeliverToSlack(webhook, channel string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Delivery = append(jb.job.Delivery, DeliveryConfig{
		"type":    "slack",
		"webhook": webhook,
		"channel": channel,
	})
	return jb
}

// DeliverToS3 adds S3 delivery
func (jb *JobBuilder) DeliverToS3(bucket, prefix string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Delivery = append(jb.job.Delivery, DeliveryConfig{
		"type":   "s3",
		"bucket": bucket,
		"prefix": prefix,
	})
	return jb
}

// WithCustomDelivery allows any delivery configuration
func (jb *JobBuilder) WithCustomDelivery(config DeliveryConfig) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Delivery = append(jb.job.Delivery, config)
	return jb
}

// ScheduleDaily schedules job to run daily at specified time
func (jb *JobBuilder) ScheduleDaily(hour, minute int) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		jb.err = fmt.Errorf("invalid time: %d:%d", hour, minute)
		return jb
	}
	jb.job.Schedule = fmt.Sprintf("0 %d %d * * *", minute, hour)
	return jb
}

// ScheduleWeekly schedules job to run weekly
func (jb *JobBuilder) ScheduleWeekly(weekday time.Weekday, hour, minute int) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		jb.err = fmt.Errorf("invalid time: %d:%d", hour, minute)
		return jb
	}
	jb.job.Schedule = fmt.Sprintf("0 %d %d * * %d", minute, hour, weekday)
	return jb
}

// ScheduleMonthly schedules job to run monthly on specified day
func (jb *JobBuilder) ScheduleMonthly(day, hour, minute int) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	if day < 1 || day > 31 || hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		jb.err = fmt.Errorf("invalid date/time: day %d, %d:%d", day, hour, minute)
		return jb
	}
	jb.job.Schedule = fmt.Sprintf("0 %d %d %d * *", minute, hour, day)
	return jb
}

// ScheduleEvery schedules job to run at regular intervals
func (jb *JobBuilder) ScheduleEvery(interval time.Duration) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	if interval < time.Second {
		jb.err = fmt.Errorf("interval too short: %v", interval)
		return jb
	}
	jb.job.Schedule = fmt.Sprintf("@every %s", interval.String())
	return jb
}

// WithCronSchedule sets a custom cron expression
func (jb *JobBuilder) WithCronSchedule(cronExpr string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Schedule = cronExpr
	return jb
}

// WithTimeout sets job execution timeout
func (jb *JobBuilder) WithTimeout(timeout time.Duration) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Timeout = timeout
	return jb
}

// WithLabel adds a label to the job
func (jb *JobBuilder) WithLabel(key, value string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Labels[key] = value
	return jb
}

// WithLabels adds multiple labels to the job
func (jb *JobBuilder) WithLabels(labels map[string]string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	for k, v := range labels {
		jb.job.Labels[k] = v
	}
	return jb
}

// Build creates the final ReportJob
func (jb *JobBuilder) Build() (ReportJob, error) {
	if jb.err != nil {
		return ReportJob{}, jb.err
	}

	// Validation
	if jb.job.Name == "" {
		return ReportJob{}, fmt.Errorf("job name is required")
	}
	if jb.job.TemplatePath == "" {
		return ReportJob{}, fmt.Errorf("template path is required")
	}
	if jb.job.DataSource == nil || jb.job.DataSource["type"] == "" {
		return ReportJob{}, fmt.Errorf("data source is required")
	}
	if len(jb.job.Outputs) == 0 {
		return ReportJob{}, fmt.Errorf("at least one output format is required")
	}
	if len(jb.job.Delivery) == 0 {
		return ReportJob{}, fmt.Errorf("at least one delivery method is required")
	}

	return jb.job, nil
}

// MustBuild creates the job and panics on error (useful for testing)
func (jb *JobBuilder) MustBuild() ReportJob {
	job, err := jb.Build()
	if err != nil {
		panic(err)
	}
	return job
}

// generateJobID creates a random job ID
func generateJobID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return fmt.Sprintf("job-%x", bytes)
}

// Convenience functions for common job patterns

// DailyReport creates a builder for daily reports
func DailyReport(name string) *JobBuilder {
	return NewJob(name).
		ScheduleDaily(9, 0). // 9:00 AM by default
		WithLabel("type", "daily")
}

// WeeklyReport creates a builder for weekly reports
func WeeklyReport(name string) *JobBuilder {
	return NewJob(name).
		ScheduleWeekly(time.Monday, 9, 0). // Monday 9:00 AM by default
		WithLabel("type", "weekly")
}

// MonthlyReport creates a builder for monthly reports
func MonthlyReport(name string) *JobBuilder {
	return NewJob(name).
		ScheduleMonthly(1, 9, 0). // 1st day of month, 9:00 AM
		WithLabel("type", "monthly")
}
