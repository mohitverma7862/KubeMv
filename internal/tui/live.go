package tui

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
	"github.com/rivo/tview"
)

type watchRequest struct {
	kind      kubernetes.ResourceKindPath
	namespace string
	refresh   int
}

type watchCtrl struct {
	mu  sync.Mutex
	req watchRequest
	ch  chan struct{}
}

func newWatchCtrl() *watchCtrl {
	return &watchCtrl{ch: make(chan struct{}, 1)}
}

func (c *watchCtrl) Update(model *Model) {
	if model.watch == nil {
		return
	}
	next := watchRequest{kind: model.kind().Path, namespace: model.Namespace, refresh: model.RefreshCount}
	c.mu.Lock()
	same := c.req == next
	c.req = next
	c.mu.Unlock()
	if same {
		return
	}
	select {
	case c.ch <- struct{}{}:
	default:
	}
}

func (c *watchCtrl) current() watchRequest {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.req
}

func serveWatch(ctx context.Context, app *tview.Application, model *Model, ctrl *watchCtrl, paint func()) {
	if model.watch == nil {
		return
	}
	var cancel context.CancelFunc
	var epoch atomic.Int64
	for {
		select {
		case <-ctx.Done():
			if cancel != nil {
				cancel()
			}
			return
		case <-ctrl.ch:
			if cancel != nil {
				cancel()
			}
			req := ctrl.current()
			source, ok := model.OpenWatch(req.kind, req.namespace)
			if !ok {
				cancel = nil
				continue
			}
			runCtx, runCancel := context.WithCancel(ctx)
			cancel = runCancel
			id := epoch.Add(1)
			session := kubernetes.NewSession(8, nil, nil)
			go func(session *kubernetes.Session, id int64, req watchRequest, source kubernetes.WatchSource) {
				go func() {
					for range session.Notify() {
						status, rows := session.Read()
						app.QueueUpdateDraw(func() {
							if epoch.Load() != id {
								return
							}
							if model.Rows == nil {
								model.Rows = map[kubernetes.ResourceKindPath][]kubernetes.ResourceRow{}
							}
							model.Rows[req.kind] = rows
							model.Connection = string(status.State)
							model.Live = status.State == kubernetes.ConnectionLive
							model.ReconnectAttempt = status.Attempt
							model.RetryAfter = status.RetryIn
							model.SyncedAt = status.SyncedAt
							model.ResourceVersion = status.ResourceVersion
							model.WatchError = status.Err
							model.clamp()
							paint()
						})
					}
				}()
				session.Run(runCtx, source)
			}(session, id, req, source)
		}
	}
}
