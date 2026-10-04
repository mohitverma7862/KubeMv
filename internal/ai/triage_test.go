package ai

import "testing"

func TestBuildBundleCrashLoop(t *testing.T) {
	b := BuildBundle(Signals{
		Namespace:   "payments",
		Kind:        "Pod",
		Name:        "payment-api-7d9c8",
		Status:      "CrashLoopBackOff",
		Restarts:    "17",
		RecentEvent: "Back-off restarting failed container",
	})
	if b.Triage.Confidence < 0.7 || len(b.Triage.Hypotheses) == 0 {
		t.Fatalf("expected crashloop triage, got conf=%f hyps=%d", b.Triage.Confidence, len(b.Triage.Hypotheses))
	}
}
