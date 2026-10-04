package api

import (
	"net/http"

	"github.com/mohitverma7862/KubeMv/internal/api/handlers"
	"github.com/mohitverma7862/KubeMv/internal/api/middleware"
	"github.com/mohitverma7862/KubeMv/internal/auth"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes/portforward"
)

const APIVersion = "0.5.0-phase4"

type Dependencies struct {
	Authenticator auth.Authenticator
	Connector     kubernetes.Connector
}

func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /healthz", handlers.HealthHandler{Version: APIVersion})
	mux.Handle("GET /api/v1/meta", handlers.MetaHandler{Phase: "4", Version: APIVersion})

	authHandler := handlers.AuthHandler{Authenticator: deps.Authenticator}
	mux.Handle("POST /api/v1/auth/login", http.HandlerFunc(authHandler.Login))

	protected := http.NewServeMux()
	protected.Handle("GET /api/v1/auth/me", http.HandlerFunc(authHandler.Me))

	clusterHandler := handlers.ClustersHandler{Connector: deps.Connector}
	protected.Handle("GET /api/v1/clusters", http.HandlerFunc(clusterHandler.List))
	protected.Handle("GET /api/v1/clusters/{clusterID}/overview", http.HandlerFunc(clusterHandler.Overview))

	resourceHandler := handlers.ResourcesHandler{Connector: deps.Connector}
	protected.Handle("GET /api/v1/clusters/{clusterID}/resources/{kind}", http.HandlerFunc(resourceHandler.List))
	protected.Handle("GET /api/v1/clusters/{clusterID}/resources/{kind}/{namespace}/{name}", http.HandlerFunc(resourceHandler.Get))
	protected.Handle("GET /api/v1/clusters/{clusterID}/crds", http.HandlerFunc(resourceHandler.CRDs))

	podsHandler := handlers.PodsHandler{Connector: deps.Connector}
	protected.Handle("GET /api/v1/clusters/{clusterID}/pods/{namespace}/{name}/containers", http.HandlerFunc(podsHandler.Containers))
	protected.Handle("GET /api/v1/clusters/{clusterID}/pods/{namespace}/{name}/logs", http.HandlerFunc(podsHandler.Logs))

	pfManager := portforward.NewManager(deps.Connector)
	pfHandler := handlers.PortForwardHandler{Manager: pfManager}
	protected.Handle("GET /api/v1/clusters/{clusterID}/portforwards", http.HandlerFunc(pfHandler.List))
	protected.Handle("POST /api/v1/clusters/{clusterID}/portforwards", http.HandlerFunc(pfHandler.Create))
	protected.Handle("DELETE /api/v1/clusters/{clusterID}/portforwards/{id}", http.HandlerFunc(pfHandler.Delete))

	execHandler := handlers.ExecWSHandler{Authenticator: deps.Authenticator, Connector: deps.Connector}
	mux.Handle("GET /api/v1/ws/clusters/{clusterID}/exec", execHandler)

	workloadsHandler := handlers.WorkloadsHandler{Connector: deps.Connector}
	protected.Handle("GET /api/v1/clusters/{clusterID}/workloads/{kind}/{namespace}/{name}/rollout", http.HandlerFunc(workloadsHandler.Rollout))
	protected.Handle("POST /api/v1/clusters/{clusterID}/workloads/{kind}/{namespace}/{name}/scale", http.HandlerFunc(workloadsHandler.Scale))
	protected.Handle("POST /api/v1/clusters/{clusterID}/workloads/{kind}/{namespace}/{name}/restart", http.HandlerFunc(workloadsHandler.Restart))
	protected.Handle("POST /api/v1/clusters/{clusterID}/workloads/deployments/{namespace}/{name}/rollback", http.HandlerFunc(workloadsHandler.Rollback))

	topologyHandler := handlers.TopologyHandler{Connector: deps.Connector}
	protected.Handle("GET /api/v1/clusters/{clusterID}/topology", http.HandlerFunc(topologyHandler.Get))

	mux.Handle("/api/v1/", middleware.RequireAuth(deps.Authenticator)(protected))

	return withCORS(mux)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
