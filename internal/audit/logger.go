package audit

import (
	"fmt"
	"log"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/auth"
)

// LogMutation records a workload mutation for Phase 3 audit trail (stdout; durable store in later phases).
func LogMutation(principal auth.Principal, clusterID, action, resource, result string) string {
	id := fmt.Sprintf("aud-%d", time.Now().UnixNano())
	log.Printf("audit id=%s user=%s cluster=%s action=%s resource=%s result=%s", id, principal.ID, clusterID, action, resource, result)
	return id
}
