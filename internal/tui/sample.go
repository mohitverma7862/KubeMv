package tui

import "github.com/mohitverma7862/KubeMv/internal/kubernetes"

func defaultKinds() []Kind {
	return []Kind{
		{Label: "Pods", Mnemonic: "K", Path: kubernetes.ResourcePods, Available: true},
		{Label: "Deployments", Mnemonic: "D", Path: kubernetes.ResourceDeployments, Available: true},
		{Label: "Services", Mnemonic: "S", Path: kubernetes.ResourceServices, Available: true},
		{Label: "Nodes", Mnemonic: "N", Path: kubernetes.ResourceNodes, Available: true},
		{Label: "Ingresses", Mnemonic: "I", Available: false},
		{Label: "Events", Mnemonic: "E", Path: kubernetes.ResourceEvents, Available: true},
		{Label: "Jobs", Mnemonic: "J", Path: kubernetes.ResourceJobs, Available: true},
		{Label: "CronJobs", Mnemonic: "C", Path: kubernetes.ResourceCronJobs, Available: true},
		{Label: "ConfigMaps", Mnemonic: "M", Path: kubernetes.ResourceConfigMaps, Available: true},
		{Label: "Secrets", Mnemonic: "R", Path: kubernetes.ResourceSecrets, Available: true},
		{Label: "StatefulSets", Path: kubernetes.ResourceStatefulSets, Available: true},
		{Label: "DaemonSets", Path: kubernetes.ResourceDaemonSets, Available: true},
		{Label: "Namespaces", Path: kubernetes.ResourceNamespaces, Available: true},
		{Label: "CRDs", Path: kubernetes.ResourceCRDs, Available: true},
	}
}

func row(kind, namespace, name, status, age string, extra map[string]string) kubernetes.ResourceRow {
	return kubernetes.ResourceRow{Kind: kind, Namespace: namespace, Name: name, Status: status, Age: age, Extra: extra}
}

func sampleModel(note string) Model {
	pods := kubernetes.ResourcePods
	api := row("Pod", "production", "api-7c9f8b6d4-k2m9x", "Running", "4h", map[string]string{
		"ready": "1/1", "restarts": "0", "ip": "10.8.1.14", "node": "ip-10-0-1-12",
		"cpu": "42m", "mem": "128Mi", "qos": "Burstable", "limits": "cpu=500m mem=256Mi",
		"ports": "8080", "mounts": "config=configmap:api",
	})
	api.Labels = map[string]string{"app": "api"}
	api.Owners = []kubernetes.OwnerRef{{Kind: "ReplicaSet", Namespace: "production", Name: "api-7c9f8b6d4"}}
	api.Containers = []kubernetes.ContainerRef{{Name: "api", Image: "api:1.2.3"}}
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
		Relations: map[string]kubernetes.ResourceRow{
			"ReplicaSet/production/api-7c9f8b6d4": {
				Kind: "ReplicaSet", Namespace: "production", Name: "api-7c9f8b6d4",
				Owners: []kubernetes.OwnerRef{{Kind: "Deployment", Namespace: "production", Name: "api"}},
			},
		},
		Rows: map[kubernetes.ResourceKindPath][]kubernetes.ResourceRow{
			pods: {
				api,
				row("Pod", "production", "payments-worker-6b8d4c7f9-q1w2e", "CrashLoopBackOff", "27m", map[string]string{"ready": "0/1", "restarts": "14", "ip": "10.8.1.22", "node": "ip-10-0-1-18", "cpu": "8m", "mem": "64Mi"}),
				row("Pod", "production", "web-6d5c4b3a2-h7t1a", "Running", "1d", map[string]string{"ready": "1/1", "restarts": "1", "ip": "10.8.2.4", "node": "ip-10-0-1-12", "cpu": "18m", "mem": "96Mi"}),
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
			kubernetes.ResourceEvents: {
				row("Event", "production", "payments-worker.18a3", "Warning", "1m", map[string]string{"reason": "BackOff", "object": "payments-worker-6b8d4c7f9-q1w2e"}),
				row("Event", "production", "api.9c21", "Normal", "4h", map[string]string{"reason": "Started", "object": "api-7c9f8b6d4-k2m9x"}),
			},
			kubernetes.ResourceJobs: {
				row("Job", "production", "migrate-db", "Complete", "2d", map[string]string{"completions": "1/1"}),
			},
			kubernetes.ResourceCronJobs: {
				row("CronJob", "production", "backup", "Active", "30d", map[string]string{"schedule": "0 2 * * *"}),
			},
			kubernetes.ResourceConfigMaps: {
				row("ConfigMap", "production", "api", "Active", "12d", map[string]string{"keys": "3"}),
			},
			kubernetes.ResourceSecrets: {
				row("Secret", "production", "db", "Active", "12d", map[string]string{"type": "Opaque", "keys": "2"}),
			},
			kubernetes.ResourceStatefulSets: {
				row("StatefulSet", "production", "redis", "Available", "20d", map[string]string{"ready": "1/1"}),
			},
		},
	}
}
