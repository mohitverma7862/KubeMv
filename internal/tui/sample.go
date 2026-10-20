package tui

import "github.com/mohitverma7862/KubeMv/internal/kubernetes"

func defaultKinds() []Kind {
	return []Kind{
		{Label: "Pods", Path: kubernetes.ResourcePods},
		{Label: "Deployments", Path: kubernetes.ResourceDeployments},
		{Label: "Services", Path: kubernetes.ResourceServices},
		{Label: "Nodes", Path: kubernetes.ResourceNodes},
		{Label: "Events", Path: kubernetes.ResourceEvents},
		{Label: "Jobs", Path: kubernetes.ResourceJobs},
		{Label: "CronJobs", Path: kubernetes.ResourceCronJobs},
		{Label: "ConfigMaps", Path: kubernetes.ResourceConfigMaps},
		{Label: "Secrets", Path: kubernetes.ResourceSecrets},
	}
}

func row(kind, namespace, name, status, age string, extra map[string]string) kubernetes.ResourceRow {
	return kubernetes.ResourceRow{Kind: kind, Namespace: namespace, Name: name, Status: status, Age: age, Extra: extra}
}

func sampleModel(note string) Model {
	pods := kubernetes.ResourcePods
	return Model{
		Cluster:    "prod-cluster",
		Context:    "prod",
		Namespace:  "production",
		User:       "admin",
		Production: true,
		ReadOnly:   true,
		Live:       false,
		Note:       note,
		Kinds:      defaultKinds(),
		Rows: map[kubernetes.ResourceKindPath][]kubernetes.ResourceRow{
			pods: {
				row("Pod", "production", "api-7c9f8b6d4-k2m9x", "Running", "4h", map[string]string{"ready": "1/1", "restarts": "0", "ip": "10.8.1.14", "node": "ip-10-0-1-12"}),
				row("Pod", "production", "payments-worker-6b8d4c7f9-q1w2e", "CrashLoopBackOff", "27m", map[string]string{"ready": "0/1", "restarts": "14", "ip": "10.8.1.22", "node": "ip-10-0-1-18"}),
				row("Pod", "production", "web-6d5c4b3a2-h7t1a", "Running", "1d", map[string]string{"ready": "1/1", "restarts": "1", "ip": "10.8.2.4", "node": "ip-10-0-1-12"}),
			},
			kubernetes.ResourceDeployments: {
				row("Deployment", "production", "api", "Available", "12d", map[string]string{"ready": "3/3"}),
				row("Deployment", "production", "payments-worker", "Progressing", "12d", map[string]string{"ready": "0/2"}),
			},
			kubernetes.ResourceServices: {
				row("Service", "production", "api", "ClusterIP", "12d", map[string]string{"ip": "10.96.0.21"}),
			},
			kubernetes.ResourceNodes: {
				row("Node", "", "ip-10-0-1-12", "Ready", "40d", nil),
			},
		},
	}
}
