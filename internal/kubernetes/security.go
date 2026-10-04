package kubernetes

// SecuritySummary aggregates RBAC visibility and policy hints for a namespace scope.
type SecuritySummary struct {
	Namespace string          `json:"namespace"`
	Score     int             `json:"score"`
	Grade     string          `json:"grade"`
	Stats     SecurityStats   `json:"stats"`
	RBAC      []RBACBinding   `json:"rbac"`
	Findings  []PolicyFinding `json:"findings"`
}

type SecurityStats struct {
	RoleBindings        int `json:"roleBindings"`
	ClusterRoleBindings int `json:"clusterRoleBindings"`
	HighFindings        int `json:"highFindings"`
	MediumFindings      int `json:"mediumFindings"`
	LowFindings         int `json:"lowFindings"`
}

type RBACBinding struct {
	Kind      string   `json:"kind"`
	Namespace string   `json:"namespace,omitempty"`
	Name      string   `json:"name"`
	RoleRef   string   `json:"roleRef"`
	Subjects  []string `json:"subjects"`
	Risk      string   `json:"risk"`
}

type PolicyFinding struct {
	ID          string `json:"id"`
	Severity    string `json:"severity"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Message     string `json:"message"`
	Resource    string `json:"resource"`
	Remediation string `json:"remediation"`
}
