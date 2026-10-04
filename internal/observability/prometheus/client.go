package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

// Client proxies read-only Prometheus HTTP API calls.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type apiResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]json.RawMessage `json:"values"`
		} `json:"result"`
	} `json:"data"`
	Error string `json:"error"`
}

func (c *Client) QueryRange(ctx context.Context, q kubernetes.MetricsQuery) (kubernetes.MetricsQueryResult, error) {
	if c.BaseURL == "" {
		return kubernetes.MetricsQueryResult{}, fmt.Errorf("prometheus url not configured")
	}
	u, err := url.Parse(c.BaseURL + "/api/v1/query_range")
	if err != nil {
		return kubernetes.MetricsQueryResult{}, err
	}
	params := url.Values{}
	params.Set("query", q.Query)
	params.Set("start", fmt.Sprintf("%d", q.Start.Unix()))
	params.Set("end", fmt.Sprintf("%d", q.End.Unix()))
	step := q.Step
	if step <= 0 {
		step = time.Minute
	}
	params.Set("step", fmt.Sprintf("%d", int(step.Seconds())))
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return kubernetes.MetricsQueryResult{}, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return kubernetes.MetricsQueryResult{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return kubernetes.MetricsQueryResult{}, fmt.Errorf("prometheus returned %s", resp.Status)
	}
	var body apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return kubernetes.MetricsQueryResult{}, err
	}
	if body.Status != "success" {
		return kubernetes.MetricsQueryResult{}, fmt.Errorf("prometheus error: %s", body.Error)
	}
	out := kubernetes.MetricsQueryResult{ResultType: body.Data.ResultType}
	for _, row := range body.Data.Result {
		series := kubernetes.MetricSeries{Labels: row.Metric}
		for _, pair := range row.Values {
			if len(pair) < 2 {
				continue
			}
			ts, val, err := parseSample(pair[0], pair[1])
			if err != nil {
				continue
			}
			series.Points = append(series.Points, kubernetes.MetricSample{Timestamp: ts, Value: val})
		}
		out.Series = append(out.Series, series)
	}
	return out, nil
}

func parseSample(tsRaw, valRaw json.RawMessage) (time.Time, float64, error) {
	var tsFloat float64
	if err := json.Unmarshal(tsRaw, &tsFloat); err != nil {
		return time.Time{}, 0, err
	}
	var valStr string
	if err := json.Unmarshal(valRaw, &valStr); err != nil {
		return time.Time{}, 0, err
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return time.Time{}, 0, err
	}
	return time.Unix(int64(tsFloat), 0).UTC(), val, nil
}
