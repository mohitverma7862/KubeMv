package stub

import (
	"context"
	"testing"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func TestStubScaleWorkload(t *testing.T) {
	c := &client{}
	result, err := c.ScaleWorkload(context.Background(), kubernetes.ResourceDeployments, "payments", "payment-api", 5)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "SUCCESS" {
		t.Fatalf("unexpected status: %s", result.Status)
	}
}
