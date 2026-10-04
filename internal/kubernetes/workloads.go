package kubernetes

// RolloutRevision describes one deployment revision in a rollout.
type RolloutRevision struct {
	Revision int64  `json:"revision"`
	Replicas int32  `json:"replicas"`
	Ready    int32  `json:"ready"`
	Progress int    `json:"progress"`
	Status   string `json:"status"`
}

// RolloutStatus is deployment rollout progress (Phase 3).
type RolloutStatus struct {
	Kind        string            `json:"kind"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Replicas    int32             `json:"replicas"`
	Ready       int32             `json:"ready"`
	Updated     int32             `json:"updated"`
	Available   int32             `json:"available"`
	Strategy    string            `json:"strategy"`
	Revisions   []RolloutRevision `json:"revisions"`
}

// MutationResult is returned from safe write operations.
type MutationResult struct {
	Action   string `json:"action"`
	Risk     string `json:"risk"`
	Status   string `json:"status"`
	Message  string `json:"message"`
	AuditID  string `json:"auditId"`
}
