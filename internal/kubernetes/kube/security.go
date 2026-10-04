package kube

import (
	"context"
	"fmt"
	"strings"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *liveClient) GetSecuritySummary(ctx context.Context, namespace string) (kubemvk8s.SecuritySummary, error) {
	if namespace == "" {
		namespace = metav1.NamespaceDefault
	}
	rbac := make([]kubemvk8s.RBACBinding, 0)
	rbList, err := c.clientset.RbacV1().RoleBindings(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return kubemvk8s.SecuritySummary{}, err
	}
	for _, item := range rbList.Items {
		rbac = append(rbac, bindingFromRoleBinding(item))
	}
	crbList, err := c.clientset.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		return kubemvk8s.SecuritySummary{}, err
	}
	for _, item := range crbList.Items {
		if bindingAffectsNamespace(item, namespace) {
			rbac = append(rbac, bindingFromClusterRoleBinding(item))
		}
	}

	findings := make([]kubemvk8s.PolicyFinding, 0)
	for _, b := range rbac {
		if strings.EqualFold(b.Risk, "high") {
			findings = append(findings, kubemvk8s.PolicyFinding{
				ID:          "rbac-" + b.Name,
				Severity:    "high",
				Category:    "RBAC",
				Title:       "Privileged cluster role binding",
				Message:     fmt.Sprintf("%s %s references %s", b.Kind, b.Name, b.RoleRef),
				Resource:    b.Kind + "/" + b.Name,
				Remediation: "Restrict to namespace RoleBinding with least privilege.",
			})
		}
	}

	npList, err := c.clientset.NetworkingV1().NetworkPolicies(namespace).List(ctx, metav1.ListOptions{})
	if err == nil && len(npList.Items) == 0 {
		findings = append(findings, kubemvk8s.PolicyFinding{
			ID:          "netpol-missing",
			Severity:    "medium",
			Category:    "Network",
			Title:       "No NetworkPolicy in namespace",
			Message:     "Namespace " + namespace + " has no NetworkPolicies.",
			Resource:    "Namespace/" + namespace,
			Remediation: "Add default-deny and explicit allow NetworkPolicies.",
		})
	}

	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err == nil {
		for _, pod := range pods.Items {
			for _, container := range pod.Spec.Containers {
				sc := container.SecurityContext
				if sc == nil || sc.RunAsNonRoot == nil || !*sc.RunAsNonRoot {
					findings = append(findings, kubemvk8s.PolicyFinding{
						ID:          "pod-nonroot-" + pod.Name + "-" + container.Name,
						Severity:    "medium",
						Category:    "Workload",
						Title:       "Container not enforced runAsNonRoot",
						Message:     fmt.Sprintf("Pod %s container %s lacks runAsNonRoot: true.", pod.Name, container.Name),
						Resource:    "Pod/" + pod.Name,
						Remediation: "Set securityContext.runAsNonRoot and runAsUser on pod or container.",
					})
					break
				}
			}
		}
	}

	stats := kubemvk8s.SecurityStats{}
	for _, b := range rbac {
		if b.Kind == "RoleBinding" {
			stats.RoleBindings++
		} else {
			stats.ClusterRoleBindings++
		}
	}
	for _, f := range findings {
		switch f.Severity {
		case "high":
			stats.HighFindings++
		case "medium":
			stats.MediumFindings++
		default:
			stats.LowFindings++
		}
	}
	score := 100 - stats.HighFindings*15 - stats.MediumFindings*8 - stats.LowFindings*3
	if score < 0 {
		score = 0
	}
	grade := gradeForScore(score)
	return kubemvk8s.SecuritySummary{
		Namespace: namespace,
		Score:     score,
		Grade:     grade,
		Stats:     stats,
		RBAC:      rbac,
		Findings:  findings,
	}, nil
}

func bindingFromRoleBinding(item rbacv1.RoleBinding) kubemvk8s.RBACBinding {
	return kubemvk8s.RBACBinding{
		Kind:      "RoleBinding",
		Namespace: item.Namespace,
		Name:      item.Name,
		RoleRef:   item.RoleRef.Kind + "/" + item.RoleRef.Name,
		Subjects:  formatSubjects(item.Subjects),
		Risk:      riskForRoleRef(item.RoleRef.Name, item.RoleRef.Kind),
	}
}

func bindingFromClusterRoleBinding(item rbacv1.ClusterRoleBinding) kubemvk8s.RBACBinding {
	return kubemvk8s.RBACBinding{
		Kind:     "ClusterRoleBinding",
		Name:     item.Name,
		RoleRef:  item.RoleRef.Kind + "/" + item.RoleRef.Name,
		Subjects: formatSubjects(item.Subjects),
		Risk:     riskForRoleRef(item.RoleRef.Name, item.RoleRef.Kind),
	}
}

func formatSubjects(subjects []rbacv1.Subject) []string {
	out := make([]string, 0, len(subjects))
	for _, s := range subjects {
		out = append(out, fmt.Sprintf("%s:%s", s.Kind, s.Name))
	}
	return out
}

func riskForRoleRef(name, kind string) string {
	lower := strings.ToLower(name)
	if kind == "ClusterRole" && (lower == "cluster-admin" || lower == "admin" || strings.Contains(lower, "admin")) {
		return "high"
	}
	if strings.Contains(lower, "edit") || strings.Contains(lower, "write") {
		return "medium"
	}
	return "low"
}

func bindingAffectsNamespace(crb rbacv1.ClusterRoleBinding, namespace string) bool {
	for _, s := range crb.Subjects {
		if s.Kind == "ServiceAccount" && s.Namespace == namespace {
			return true
		}
	}
	// Include cluster-admin style bindings for visibility in any namespace view.
	if riskForRoleRef(crb.RoleRef.Name, crb.RoleRef.Kind) == "high" {
		return true
	}
	return false
}

func gradeForScore(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}
