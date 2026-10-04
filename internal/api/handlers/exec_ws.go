package handlers

import (
	"context"
	"io"
	"net/http"
	"strings"
	"github.com/gorilla/websocket"
	"github.com/mohitverma7862/KubeMv/internal/auth"
	"github.com/mohitverma7862/KubeMv/internal/httputil"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
	corev1 "k8s.io/api/core/v1"
	k8sclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ExecWSHandler struct {
	Authenticator auth.Authenticator
	Connector     kubernetes.Connector
}

func (h ExecWSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		httputil.WriteUnauthorized(w, "missing token")
		return
	}
	if _, err := h.Authenticator.ValidateSession(r.Context(), token); err != nil {
		httputil.WriteUnauthorized(w, "invalid session")
		return
	}
	clusterID := r.PathValue("clusterID")
	namespace := r.URL.Query().Get("namespace")
	pod := r.URL.Query().Get("pod")
	container := r.URL.Query().Get("container")
	if namespace == "" || pod == "" {
		httputil.WriteBadRequest(w, "namespace and pod required")
		return
	}

	client, err := h.Connector.Connect(r.Context(), clusterID)
	if err != nil {
		httputil.WriteInternal(w, "cluster connect failed")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	provider, ok := client.(interface{ RESTConfig() *rest.Config })
	if !ok || provider.RESTConfig() == nil {
		runStubExec(conn, namespace, pod, container)
		return
	}
	runLiveExec(r.Context(), provider.RESTConfig(), conn, namespace, pod, container)
}

func runStubExec(conn *websocket.Conn, namespace, pod, container string) {
	_ = conn.WriteMessage(websocket.TextMessage, []byte("KubeMv exec (stub) — connected to "+namespace+"/"+pod+" container="+container+"\r\n"))
	_ = conn.WriteMessage(websocket.TextMessage, []byte("$ "))
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		cmd := strings.TrimSpace(string(msg))
		if cmd == "exit" {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("session closed\r\n"))
			return
		}
		_ = conn.WriteMessage(websocket.TextMessage, []byte("(stub) executed: "+cmd+"\r\n$ "))
	}
}

func runLiveExec(ctx context.Context, cfg *rest.Config, conn *websocket.Conn, namespace, pod, container string) {
	clientset, err := k8sclientset.NewForConfig(cfg)
	if err != nil {
		_ = conn.WriteMessage(websocket.TextMessage, []byte("exec init failed\r\n"))
		return
	}
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   []string{"/bin/sh"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		runStubExec(conn, namespace, pod, container)
		return
	}

	stdinReader, stdinWriter := io.Pipe()
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				stdinWriter.Close()
				return
			}
			_, _ = stdinWriter.Write(msg)
		}
	}()

	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  stdinReader,
		Stdout: &wsWriter{conn: conn},
		Stderr: &wsWriter{conn: conn},
		Tty:    true,
	})
	if err != nil {
		_ = conn.WriteMessage(websocket.TextMessage, []byte("exec ended: "+err.Error()+"\r\n"))
	}
}

type wsWriter struct {
	conn *websocket.Conn
}

func (w *wsWriter) Write(p []byte) (int, error) {
	err := w.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
