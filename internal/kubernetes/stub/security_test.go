package stub

import (
	"context"
	"testing"
)

func TestSecuritySummary(t *testing.T) {
	c := mustClient(t, "local")
	sum, err := c.GetSecuritySummary(context.Background(), "payments")
	if err != nil {
		t.Fatal(err)
	}
	if sum.Score <= 0 || len(sum.Findings) < 3 {
		t.Fatalf("unexpected summary: score=%d findings=%d", sum.Score, len(sum.Findings))
	}
}
