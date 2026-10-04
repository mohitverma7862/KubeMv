package audit

import (
	"context"
	"strings"
	"time"
)

type Result string

const (
	ResultSuccess Result = "success"
	ResultFailure Result = "failure"
	ResultDenied  Result = "denied"
)

// Event is an internal audit record. Phase 0 keeps these inside the process.
// The administration phase exposes a reviewed audit read API.
type Event struct {
	ID         string
	Time       time.Time
	Actor      string
	Action     string
	Result     Result
	TargetType string
	TargetID   string
	Detail     string
}

// Recorder accepts security-relevant actions. Implementations must drop secret
// material instead of storing it.
type Recorder interface {
	Record(ctx context.Context, event Event)
	Events() []Event
}

var secretMarkers = []string{
	"-----begin",
	"client-key-data",
	"client-certificate-data",
	"password=",
	"token:",
	"bearer ",
	"kubeconfig",
}

func redact(detail string) string {
	if len(detail) > 300 {
		detail = detail[:300]
	}
	lowered := strings.ToLower(detail)
	for _, marker := range secretMarkers {
		if strings.Contains(lowered, marker) {
			return "[redacted]"
		}
	}
	return detail
}
