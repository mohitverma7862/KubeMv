package stub

import (
	"context"
	"math"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (c *client) QueryMetrics(_ context.Context, query kubernetes.MetricsQuery) (kubernetes.MetricsQueryResult, error) {
	end := query.End
	if end.IsZero() {
		end = time.Now().UTC()
	}
	start := query.Start
	if start.IsZero() {
		start = end.Add(-1 * time.Hour)
	}
	step := query.Step
	if step <= 0 {
		step = 5 * time.Minute
	}
	series := syntheticSeries(query.Query, start, end, step)
	return kubernetes.MetricsQueryResult{
		ResultType: "matrix",
		Series:     []kubernetes.MetricSeries{series},
	}, nil
}

func (c *client) ListScrapeTargets(_ context.Context, namespace string) ([]kubernetes.ScrapeTarget, error) {
	targets := []kubernetes.ScrapeTarget{
		{
			Namespace:  "payments",
			Kind:       "Deployment",
			Name:       "payment-api",
			Job:        "payment-api",
			Instance:   "10.0.12.4:8080",
			Up:         true,
			LastScrape: time.Now().UTC().Add(-30 * time.Second).Format(time.RFC3339),
			Labels:     map[string]string{"app": "payment-api", "team": "payments"},
		},
		{
			Namespace:  "payments",
			Kind:       "Deployment",
			Name:       "payment-worker",
			Job:        "payment-worker",
			Instance:   "10.0.12.9:8080",
			Up:         true,
			LastScrape: time.Now().UTC().Add(-45 * time.Second).Format(time.RFC3339),
			Labels:     map[string]string{"app": "payment-worker"},
		},
		{
			Namespace:  "platform",
			Kind:       "Deployment",
			Name:       "ingress-nginx-controller",
			Job:        "ingress-nginx",
			Instance:   "10.0.0.12:10254",
			Up:         false,
			LastScrape: time.Now().UTC().Add(-2 * time.Minute).Format(time.RFC3339),
			Labels:     map[string]string{"app": "ingress-nginx"},
		},
	}
	if namespace == "" {
		return targets, nil
	}
	out := make([]kubernetes.ScrapeTarget, 0)
	for _, t := range targets {
		if t.Namespace == namespace {
			out = append(out, t)
		}
	}
	return out, nil
}

func (c *client) GetObservabilityDashboard(_ context.Context, namespace, kind, name string) (kubernetes.ObservabilityDashboard, error) {
	if namespace == "" {
		namespace = "payments"
	}
	if kind == "" {
		kind = "Deployment"
	}
	if name == "" {
		name = "payment-api"
	}
	now := time.Now().UTC()
	start := now.Add(-1 * time.Hour)
	step := 5 * time.Minute
	presets := []kubernetes.ObservabilityPreset{
		{
			ID:    "cpu",
			Title: "CPU cores",
			Unit:  "cores",
			Query: `sum(rate(container_cpu_usage_seconds_total{pod=~"payment-api-.*"}[5m]))`,
			Series: syntheticSeries("cpu", start, now, step),
		},
		{
			ID:    "memory",
			Title: "Memory working set",
			Unit:  "MiB",
			Query: `sum(container_memory_working_set_bytes{pod=~"payment-api-.*"}) / 1024 / 1024`,
			Series: syntheticSeries("memory", start, now, step),
		},
		{
			ID:    "rps",
			Title: "HTTP request rate",
			Unit:  "req/s",
			Query: `sum(rate(http_requests_total{service="payment-api"}[5m]))`,
			Series: syntheticSeries("rps", start, now, step),
		},
		{
			ID:    "errors",
			Title: "5xx error rate",
			Unit:  "%",
			Query: `100 * sum(rate(http_requests_total{service="payment-api",status=~"5.."}[5m])) / sum(rate(http_requests_total{service="payment-api"}[5m]))`,
			Series: syntheticSeries("errors", start, now, step),
		},
	}
	targets, _ := c.ListScrapeTargets(context.Background(), namespace)
	return kubernetes.ObservabilityDashboard{
		Namespace:   namespace,
		Kind:        kind,
		Name:        name,
		Presets:     presets,
		Targets:     targets,
		GrafanaURL:  "https://grafana.example/d/payment-api",
		LokiURL:     "https://loki.example",
		LogDeepLink: "/fast?kind=pods&namespace=" + namespace,
	}, nil
}

func syntheticSeries(seed string, start, end time.Time, step time.Duration) kubernetes.MetricSeries {
	labels := map[string]string{"query": seed}
	points := make([]kubernetes.MetricSample, 0)
	hash := 0
	for _, ch := range seed {
		hash += int(ch)
	}
	for t := start; !t.After(end); t = t.Add(step) {
		minutes := float64(t.Unix() / 60)
		base := 0.35 + float64(hash%10)*0.03
		wave := math.Sin(minutes/12) * 0.12
		spike := 0.0
		if seed == "errors" && int(minutes)%37 == 0 {
			spike = 2.5
		}
		if seed == "rps" {
			base = 120 + float64(hash%20)
			wave = math.Sin(minutes/8) * 15
		}
		if seed == "memory" {
			base = 256 + float64(hash%32)
			wave = math.Cos(minutes/20) * 12
		}
		val := base + wave + spike
		if val < 0 {
			val = 0
		}
		points = append(points, kubernetes.MetricSample{Timestamp: t.UTC(), Value: val})
	}
	return kubernetes.MetricSeries{Labels: labels, Points: points}
}
