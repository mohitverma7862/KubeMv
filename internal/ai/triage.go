package ai

import (
	"strings"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

// Signals are collected from the cluster (or stub) before rule-based triage.
type Signals struct {
	Namespace   string
	Kind        string
	Name        string
	Status      string
	Restarts    string
	RecentEvent string
	LogSnippet  string
	ReadyRatio  string
}

// BuildBundle produces assist output without calling external LLM providers (Phase 8).
func BuildBundle(s Signals) kubernetes.AssistBundle {
	resource := kubernetes.AssistResource{
		Namespace: s.Namespace,
		Kind:      s.Kind,
		Name:      s.Name,
	}
	signals := []string{}
	if s.Status != "" {
		signals = append(signals, "status="+s.Status)
	}
	if s.Restarts != "" {
		signals = append(signals, "restarts="+s.Restarts)
	}
	if s.RecentEvent != "" {
		signals = append(signals, "event="+s.RecentEvent)
	}
	if s.ReadyRatio != "" {
		signals = append(signals, "ready="+s.ReadyRatio)
	}

	lower := strings.ToLower(s.Status + " " + s.RecentEvent + " " + s.LogSnippet)
	hypotheses := []kubernetes.Hypothesis{}
	runbookID := "generic-health"
	runbookTitle := "General workload health check"
	steps := []kubernetes.RunbookStep{
		{Order: 1, Title: "Confirm current status", Command: kubectlGet(s), Caution: "Read-only"},
		{Order: 2, Title: "Inspect recent events", Command: kubectlEvents(s), Caution: "Read-only"},
	}
	summary := "No critical signals detected; continue standard health checks."
	confidence := 0.55

	if strings.Contains(lower, "crashloop") || strings.Contains(lower, "backoff") {
		summary = "Workload shows crash-loop behavior — likely failing container start or missing dependency."
		confidence = 0.82
		runbookID = "crashloop-triage"
		runbookTitle = "CrashLoopBackOff triage"
		hypotheses = []kubernetes.Hypothesis{
			{Title: "Application exit on boot", Likelihood: "high", Evidence: "CrashLoopBackOff with restarts"},
			{Title: "Missing ConfigMap/Secret key", Likelihood: "medium", Evidence: "Common after manifest drift"},
			{Title: "Upstream DB unreachable", Likelihood: "medium", Evidence: "Payments stack dependency pattern"},
		}
		steps = append(steps,
			kubernetes.RunbookStep{Order: 3, Title: "Fetch last container logs", Command: kubectlLogs(s), Caution: "Read-only"},
			kubernetes.RunbookStep{Order: 4, Title: "Compare live vs Git manifest", Command: "# open GitOps drift view for " + s.Name, Caution: "No mutation"},
		)
	}
	if strings.Contains(lower, "imagepull") {
		summary = "Image pull failures detected — verify registry auth and image tag."
		confidence = 0.88
		runbookID = "image-pull"
		runbookTitle = "ImagePullBackOff remediation"
		hypotheses = []kubernetes.Hypothesis{
			{Title: "Invalid or removed image tag", Likelihood: "high", Evidence: "ImagePullBackOff"},
		}
	}

	hooks := []kubernetes.AutomationHook{
		{
			ID:               "restart-workload",
			Title:            "Rolling restart workload",
			Description:      "Triggers a controlled rollout restart after operator approval.",
			Risk:             "medium",
			RequiresApproval: true,
			DryRunSupported:  true,
		},
		{
			ID:               "scale-buffer",
			Title:            "Temporary scale-up (+1 replica)",
			Description:      "Adds one replica to absorb crash noise while investigating.",
			Risk:             "medium",
			RequiresApproval: true,
			DryRunSupported:  true,
		},
		{
			ID:               "rollback-deployment",
			Title:            "Rollback to previous revision",
			Description:      "Rolls deployment back one revision when Git sync is healthy.",
			Risk:             "high",
			RequiresApproval: true,
			DryRunSupported:  true,
		},
	}

	return kubernetes.AssistBundle{
		Resource: resource,
		Triage: kubernetes.TriageResult{
			Summary:    summary,
			Confidence: confidence,
			Hypotheses: hypotheses,
			Signals:    signals,
			Disclaimer: "Phase 8 uses on-cluster signals and rule-based triage only; no external LLM calls.",
		},
		Runbook: kubernetes.Runbook{ID: runbookID, Title: runbookTitle, Steps: steps},
		Hooks:   hooks,
	}
}

func kubectlGet(s Signals) string {
	if strings.EqualFold(s.Kind, "Pod") {
		return "kubectl -n " + s.Namespace + " describe pod " + s.Name
	}
	return "kubectl -n " + s.Namespace + " describe " + strings.ToLower(s.Kind) + " " + s.Name
}

func kubectlEvents(s Signals) string {
	return "kubectl -n " + s.Namespace + " get events --field-selector involvedObject.name=" + s.Name
}

func kubectlLogs(s Signals) string {
	pod := s.Name
	if !strings.EqualFold(s.Kind, "Pod") {
		pod = s.Name + "-*"
	}
	return "kubectl -n " + s.Namespace + " logs " + pod + " --tail=200"
}

// DryRunHook returns planned actions without mutating the cluster.
func DryRunHook(hookID string, s Signals) kubernetes.HookDryRunResult {
	actions := []string{}
	msg := "Dry-run only — no changes applied."
	switch hookID {
	case "restart-workload":
		actions = []string{
			"POST /workloads/" + strings.ToLower(s.Kind) + "/" + s.Namespace + "/" + s.Name + "/restart",
		}
		msg = "Would request rolling restart after approval."
	case "scale-buffer":
		actions = []string{
			"POST /workloads/deployments/" + s.Namespace + "/" + s.Name + "/scale { replicas: +1 }",
		}
		msg = "Would scale deployment up by one replica."
	case "rollback-deployment":
		actions = []string{
			"POST /workloads/deployments/" + s.Namespace + "/" + s.Name + "/rollback?revision=previous",
		}
		msg = "Would rollback deployment to previous revision."
	default:
		return kubernetes.HookDryRunResult{
			HookID:  hookID,
			Status:  "unknown",
			Message: "Unknown automation hook",
		}
	}
	return kubernetes.HookDryRunResult{
		HookID:         hookID,
		Status:         "planned",
		Message:        msg,
		PlannedActions: actions,
		AuditID:        "dryrun-" + hookID,
	}
}
