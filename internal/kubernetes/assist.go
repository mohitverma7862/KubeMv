package kubernetes

// AssistBundle is the unified AI-assist response for a workload or pod context.
type AssistBundle struct {
	Resource AssistResource `json:"resource"`
	Triage   TriageResult   `json:"triage"`
	Runbook  Runbook        `json:"runbook"`
	Hooks    []AutomationHook `json:"hooks"`
}

type AssistResource struct {
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Name      string `json:"name"`
}

type TriageResult struct {
	Summary     string       `json:"summary"`
	Confidence  float64      `json:"confidence"`
	Hypotheses  []Hypothesis `json:"hypotheses"`
	Signals     []string     `json:"signals"`
	Disclaimer  string       `json:"disclaimer"`
}

type Hypothesis struct {
	Title      string `json:"title"`
	Likelihood string `json:"likelihood"`
	Evidence   string `json:"evidence"`
}

type Runbook struct {
	ID    string         `json:"id"`
	Title string         `json:"title"`
	Steps []RunbookStep  `json:"steps"`
}

type RunbookStep struct {
	Order     int    `json:"order"`
	Title     string `json:"title"`
	Command   string `json:"command,omitempty"`
	Caution   string `json:"caution,omitempty"`
}

type AutomationHook struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Risk             string `json:"risk"`
	RequiresApproval bool   `json:"requiresApproval"`
	DryRunSupported  bool   `json:"dryRunSupported"`
}

type HookDryRunResult struct {
	HookID         string   `json:"hookId"`
	Status         string   `json:"status"`
	Message        string   `json:"message"`
	PlannedActions []string `json:"plannedActions"`
	AuditID        string   `json:"auditId"`
}
