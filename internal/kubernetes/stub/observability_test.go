package stub

import (
	"context"
	"testing"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func TestObservabilityDashboard(t *testing.T) {
	c := mustClient(t, "local")
	dash, err := c.GetObservabilityDashboard(context.Background(), "payments", "Deployment", "payment-api")
	if err != nil {
		t.Fatal(err)
	}
	if len(dash.Presets) < 4 {
		t.Fatalf("expected presets, got %d", len(dash.Presets))
	}
	if dash.Presets[0].Series.Points == nil || len(dash.Presets[0].Series.Points) == 0 {
		t.Fatal("expected series points")
	}
}

func TestQueryMetrics(t *testing.T) {
	c := mustClient(t, "local")
	end := time.Now().UTC()
	res, err := c.QueryMetrics(context.Background(), kubernetes.MetricsQuery{
		Query: "up",
		Start: end.Add(-30 * time.Minute),
		End:   end,
		Step:  time.Minute,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Series) == 0 || len(res.Series[0].Points) == 0 {
		t.Fatal("expected metric points")
	}
}

func mustClient(t *testing.T, id string) kubernetes.ClusterClient {
	conn := NewConnector()
	client, err := conn.Connect(context.Background(), id)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	return client
}
