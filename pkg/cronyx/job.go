package cronyx

import "time"

type ReportJob struct {
	ID           string
	Name         string
	TemplatePath string
	DataSource   DataSourceConfig
	Outputs      []string
	Schedule     string
	Delivery     []DeliveryConfig
	Timeout      time.Duration
	Labels       map[string]string
}

// DataSourceConfig is generic; specific loaders will parse it.
type DataSourceConfig map[string]string

type DeliveryConfig map[string]string
