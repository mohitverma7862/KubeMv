package stub

import (
	"context"
	"testing"
)

func TestGitOpsOverview(t *testing.T) {
	c := mustClient(t, "local")
	ov, err := c.GetGitOpsOverview(context.Background(), "payments")
	if err != nil {
		t.Fatal(err)
	}
	if len(ov.Applications) < 2 || len(ov.Drift) == 0 {
		t.Fatalf("unexpected overview apps=%d drift=%d", len(ov.Applications), len(ov.Drift))
	}
}
