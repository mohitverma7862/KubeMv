package stub

import (
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func stubWorkloadRows(kind kubernetes.ResourceKindPath, ns string) []kubernetes.ResourceRow {
	ns = defaultNS(ns)
	switch kind {
	case kubernetes.ResourceStatefulSets:
		return []kubernetes.ResourceRow{
			{Kind: "StatefulSet", Namespace: ns, Name: "payment-db", Status: "Healthy", Age: "40d", Extra: map[string]string{"ready": "3/3"}},
		}
	case kubernetes.ResourceDaemonSets:
		return []kubernetes.ResourceRow{
			{Kind: "DaemonSet", Namespace: "platform", Name: "node-exporter", Status: "Healthy", Age: "90d", Extra: map[string]string{"ready": "3/3"}},
		}
	case kubernetes.ResourceJobs:
		return []kubernetes.ResourceRow{
			{Kind: "Job", Namespace: ns, Name: "migrate-20261004", Status: "Complete", Age: "1h", Extra: map[string]string{"succeeded": "1"}},
		}
	case kubernetes.ResourceCronJobs:
		return []kubernetes.ResourceRow{
			{Kind: "CronJob", Namespace: ns, Name: "backup-hourly", Status: "Active", Age: "120d", Extra: map[string]string{"schedule": "0 * * * *"}},
		}
	case kubernetes.ResourceHPA:
		return []kubernetes.ResourceRow{
			{Kind: "HorizontalPodAutoscaler", Namespace: ns, Name: "payment-api", Status: "Active", Age: "30d", Extra: map[string]string{"targets": "2-10", "current": "3"}},
		}
	case kubernetes.ResourcePDB:
		return []kubernetes.ResourceRow{
			{Kind: "PodDisruptionBudget", Namespace: ns, Name: "payment-api-pdb", Status: "Healthy", Age: "30d", Extra: map[string]string{"minAvailable": "2"}},
		}
	default:
		return nil
	}
}
