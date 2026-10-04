package stub

import (
	"context"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (c *client) GetSecuritySummary(_ context.Context, namespace string) (kubernetes.SecuritySummary, error) {
	if namespace == "" {
		namespace = "payments"
	}
	rbac := []kubernetes.RBACBinding{
		{
			Kind:      "RoleBinding",
			Namespace: namespace,
			Name:      "payment-api-deployer",
			RoleRef:   "Role/payment-api-admin",
			Subjects:  []string{"User:ci-bot", "Group:payments-sre"},
			Risk:      "medium",
		},
		{
			Kind:      "RoleBinding",
			Namespace: namespace,
			Name:      "payment-api-readonly",
			RoleRef:   "Role/view",
			Subjects:  []string{"Group:payments-dev"},
			Risk:      "low",
		},
		{
			Kind:     "ClusterRoleBinding",
			Name:     "payments-breakglass",
			RoleRef:  "ClusterRole/cluster-admin",
			Subjects: []string{"User:breakglass@corp.example"},
			Risk:     "high",
		},
	}
	findings := []kubernetes.PolicyFinding{
		{
			ID:          "rbac-cluster-admin-binding",
			Severity:    "high",
			Category:    "RBAC",
			Title:       "Cluster-admin binding present",
			Message:     "ClusterRoleBinding payments-breakglass grants cluster-admin to a user subject.",
			Resource:    "ClusterRoleBinding/payments-breakglass",
			Remediation: "Replace with namespace-scoped RoleBinding and least-privilege ClusterRole.",
		},
		{
			ID:          "pod-root-filesystem",
			Severity:    "high",
			Category:    "Workload",
			Title:       "Container may run as root",
			Message:     "Pod payment-api-7d9c8 has no runAsNonRoot security context.",
			Resource:    "Pod/payment-api-7d9c8",
			Remediation: "Set pod or container securityContext.runAsNonRoot: true and runAsUser.",
		},
		{
			ID:          "netpol-missing",
			Severity:    "medium",
			Category:    "Network",
			Title:       "No NetworkPolicy for namespace",
			Message:     "Namespace payments has no default deny or workload-scoped NetworkPolicies.",
			Resource:    "Namespace/payments",
			Remediation: "Add default-deny NetworkPolicy and explicit allow rules per workload.",
		},
		{
			ID:          "secret-plaintext-mount",
			Severity:    "medium",
			Category:    "Secrets",
			Title:       "Secret exposed as env var",
			Message:     "Deployment payment-api mounts payment-api-secrets keys as environment variables.",
			Resource:    "Deployment/payment-api",
			Remediation: "Prefer mounted volumes or external secret stores; rotate credentials regularly.",
		},
		{
			ID:          "image-tag-latest",
			Severity:    "low",
			Category:    "Supply chain",
			Title:       "Floating image tag in use",
			Message:     "Container payment-api uses image tag :latest.",
			Resource:    "Deployment/payment-api",
			Remediation: "Pin immutable image digests or semver tags in manifests.",
		},
	}
	stats := kubernetes.SecurityStats{
		RoleBindings:        2,
		ClusterRoleBindings: 1,
		HighFindings:        2,
		MediumFindings:      2,
		LowFindings:         1,
	}
	return kubernetes.SecuritySummary{
		Namespace: namespace,
		Score:     62,
		Grade:     "C",
		Stats:     stats,
		RBAC:      rbac,
		Findings:  findings,
	}, nil
}
