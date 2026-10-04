package kubernetes

// GitOpsOverview summarizes Git-backed applications, drift, and pipeline activity.
type GitOpsOverview struct {
	Namespace    string              `json:"namespace"`
	Stats        GitOpsStats         `json:"stats"`
	Applications []GitOpsApplication `json:"applications"`
	Drift        []GitOpsDrift       `json:"drift"`
	Pipelines    []PipelineRun       `json:"pipelines"`
}

type GitOpsStats struct {
	Applications int `json:"applications"`
	Synced       int `json:"synced"`
	OutOfSync    int `json:"outOfSync"`
	DriftItems   int `json:"driftItems"`
}

type GitOpsApplication struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Provider     string `json:"provider"`
	Repository   string `json:"repository"`
	Revision     string `json:"revision"`
	Path         string `json:"path"`
	SyncStatus   string `json:"syncStatus"`
	Health       string `json:"health"`
	LastSyncedAt string `json:"lastSyncedAt"`
}

type GitOpsDrift struct {
	ID         string `json:"id"`
	Severity   string `json:"severity"`
	Resource   string `json:"resource"`
	Field      string `json:"field"`
	GitValue   string `json:"gitValue"`
	LiveValue  string `json:"liveValue"`
	Suggestion string `json:"suggestion"`
}

type PipelineRun struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Trigger   string `json:"trigger"`
	Status    string `json:"status"`
	Commit    string `json:"commit"`
	StartedAt string `json:"startedAt"`
	URL       string `json:"url"`
}
