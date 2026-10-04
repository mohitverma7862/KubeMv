package kubernetes

import "time"

// MetricsQuery mirrors a subset of Prometheus HTTP API parameters.
type MetricsQuery struct {
	Query string
	Start time.Time
	End   time.Time
	Step  time.Duration
}

type MetricSample struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type MetricSeries struct {
	Labels map[string]string `json:"labels"`
	Points []MetricSample    `json:"points"`
}

type MetricsQueryResult struct {
	ResultType string         `json:"resultType"`
	Series     []MetricSeries `json:"series"`
}

type ScrapeTarget struct {
	Namespace  string            `json:"namespace"`
	Kind       string            `json:"kind"`
	Name       string            `json:"name"`
	Job        string            `json:"job"`
	Instance   string            `json:"instance"`
	Up         bool              `json:"up"`
	LastScrape string            `json:"lastScrape"`
	Labels     map[string]string `json:"labels"`
}

type ObservabilityPreset struct {
	ID    string        `json:"id"`
	Title string        `json:"title"`
	Unit  string        `json:"unit"`
	Query string        `json:"query"`
	Series MetricSeries `json:"series"`
}

type ObservabilityDashboard struct {
	Namespace      string                `json:"namespace"`
	Kind           string                `json:"kind"`
	Name           string                `json:"name"`
	Presets        []ObservabilityPreset `json:"presets"`
	Targets        []ScrapeTarget        `json:"targets"`
	PrometheusURL  string                `json:"prometheusUrl,omitempty"`
	GrafanaURL     string                `json:"grafanaUrl,omitempty"`
	LokiURL        string                `json:"lokiUrl,omitempty"`
	LogDeepLink    string                `json:"logDeepLink,omitempty"`
}
