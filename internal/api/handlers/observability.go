package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/httputil"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
	"github.com/mohitverma7862/KubeMv/internal/observability/prometheus"
)

type ObservabilityHandler struct {
	Connector       kubernetes.Connector
	PrometheusURL   string
	prometheusClient *prometheus.Client
}

func (h *ObservabilityHandler) prom() *prometheus.Client {
	if h.PrometheusURL == "" {
		return nil
	}
	if h.prometheusClient == nil {
		h.prometheusClient = prometheus.NewClient(h.PrometheusURL)
	}
	return h.prometheusClient
}

func (h *ObservabilityHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		httputil.WriteBadRequest(w, "query is required")
		return
	}
	end := time.Now().UTC()
	start := end.Add(-1 * time.Hour)
	if v := r.URL.Query().Get("start"); v != "" {
		if ts, err := parseUnixTime(v); err == nil {
			start = ts
		}
	}
	if v := r.URL.Query().Get("end"); v != "" {
		if ts, err := parseUnixTime(v); err == nil {
			end = ts
		}
	}
	step := time.Minute
	if v := r.URL.Query().Get("step"); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			step = time.Duration(sec) * time.Second
		}
	}
	mq := kubernetes.MetricsQuery{Query: query, Start: start, End: end, Step: step}

	if prom := h.prom(); prom != nil {
		result, err := prom.QueryRange(r.Context(), mq)
		if err != nil {
			httputil.WriteInternal(w, err.Error())
			return
		}
		httputil.WriteJSON(w, http.StatusOK, result)
		return
	}

	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "cluster connect failed")
		return
	}
	result, err := client.QueryMetrics(r.Context(), mq)
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, result)
}

func (h *ObservabilityHandler) Targets(w http.ResponseWriter, r *http.Request) {
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "cluster connect failed")
		return
	}
	targets, err := client.ListScrapeTargets(r.Context(), r.URL.Query().Get("namespace"))
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, targets)
}

func (h *ObservabilityHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "cluster connect failed")
		return
	}
	dash, err := client.GetObservabilityDashboard(
		r.Context(),
		r.URL.Query().Get("namespace"),
		r.URL.Query().Get("kind"),
		r.URL.Query().Get("name"),
	)
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	if h.PrometheusURL != "" {
		dash.PrometheusURL = h.PrometheusURL
	}
	httputil.WriteJSON(w, http.StatusOK, dash)
}

func parseUnixTime(raw string) (time.Time, error) {
	sec, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0).UTC(), nil
}
