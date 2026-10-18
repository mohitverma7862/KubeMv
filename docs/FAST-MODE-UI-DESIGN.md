# KubeMv Fast Mode — Advanced TUI Design Specification

## Goal

Fast Mode is the default KubeMv operator experience. It must feel like an advanced Kubernetes terminal rather than a SaaS dashboard: K9s keyboard-first interaction + modern terminal IDE density + resource intelligence + contextual AI.

> **Design principle:** KubeMv should not be a prettier dashboard. It should be a faster, denser, more informative Kubernetes operator console.

## Visual hierarchy

```text
┌ GLOBAL HEADER: KubeMv | cluster | context | namespace | status | clock ┐
├ RESOURCE NAV: Pods Deployments Services Nodes Events Jobs ...          ┤
├───────────────────────────────────────┬────────────────────────────────┤
│ RESOURCE TABLE                        │ RESOURCE INTELLIGENCE          │
│ dense, live, virtualized              │ health / containers / labels  │
│                                       │ owner / actions                │
├───────────────────────────────────────┼────────────────────────────────┤
│ LOGS / EVENTS / YAML / DESCRIBE       │ RESOURCE RELATIONSHIP TREE     │
│ terminal-style output                 │ owner/dependency graph         │
├───────────────────────────────────────┴────────────────────────────────┤
│ : command | / filter | ? help | Ctrl-P palette | ctx | ns | status    │
└────────────────────────────────────────────────────────────────────────┘
```

## 1. Design language

- Dark operator-console default; Light and High Contrast themes also supported.
- Dense information hierarchy; minimal whitespace.
- Monospace-first typography: JetBrains Mono, Cascadia Code, Consolas, monospace.
- Mostly square panels; maximum radius 4px.
- 1px borders; no heavy shadows.
- No giant cards, marketing sections, decorative gradients, excessive animation, or oversized icons.
- Accent colors communicate state rather than decoration.
- Target 120x30 character-equivalent viewport; remain usable at 100x25.

Suggested tokens:

```text
--km-bg             #050A12
--km-panel          #07111F
--km-panel-alt      #081522
--km-border         #12304A
--km-text           #D7E3F0
--km-muted          #52697D
--km-secondary      #7F96AB
--km-primary        #00B7FF
--km-success        #00D084
--km-warning        #FFB020
--km-error          #FF4D5A
--km-ai             #B56CFF
--km-production     #00D084
```

## 2. Header

48–54px high. Left: KubeMv + Fast Mode. Center: cluster, user/context, namespace. Right: production badge, connection state and clock.

Connection states: `● LIVE`, `◌ RECONNECTING n`, `✕ DISCONNECTED`.

Production contexts display `[PRODUCTION]` and require stronger confirmation for destructive operations.

## 3. Resource navigation

Compact terminal tabs below the header. Examples: `[K] Pods`, `[D] Deployments`, `[S] Services`, `[N] Nodes`, `[I] Ingresses`, `[E] Events`, `[J] Jobs`, `[C] CronJobs`, `[M] ConfigMaps`, `[R] Secrets`, `[>] More`.

Selected resource uses a subtle cyan background/border. Never use large rounded buttons.

## 4. Main resource table

The table dominates the screen. Use virtualization, sticky headers, compact 12–13px text, 32–34px maximum rows, aligned columns, row selection, sorting, resizing, horizontal scroll and incremental updates.

Example columns:

```text
NAME                         READY   STATUS       RESTARTS   AGE   IP
> api-7d9f8c9b6f-2kq7m       1/1     ● Running    0          2h    10.24.1.12
  worker-6f8d7c9b4-xz9qp     1/1     ● Running    1          3h    10.24.2.21
  backup-job-2814500         0/1     ● CrashLoop   5          2h    10.24.2.33
```

Selected row: left cyan cursor + subtle background, not a large selection card.

Status indicators: Running, Pending, Warning, CrashLoopBackOff, Terminating.

## 5. Resource intelligence inspector

Right panel 280–340px. For the selected resource show:

- health/status
- namespace, node, IP, readiness, restarts, age, QoS
- CPU/memory
- containers and images
- ports/mounts/resource limits
- labels as compact pills
- owner chain
- contextual keyboard actions

Never display secret values.

## 6. Resource relationship tree

Bottom-right panel:

```text
Deployment
 └─ ReplicaSet
     └─ Pod
         └─ Container
```

Support keyboard navigation and opening the matching Fast Mode resource view.

## 7. Lower terminal workspace

Tabbed: `[Logs] [Events] [YAML] [Describe] [Terminal] [Rollout] [History]`.

Logs must look like a real terminal: aligned timestamps, levels, follow/pause, previous logs, search `/`, next/previous match `n/N`, wrap toggle, container selector and Ctrl-C to stop streaming.

Events use a compact timeline/table with Normal/Warning/Error severity.

YAML: syntax highlighting, search, line numbers, copy/download.

Describe: structured metadata, containers, conditions and prominent events.

Exec: xterm.js/WebSocket terminal with resize and multi-container support.

## 8. Command bar

Always-visible bottom bar. Inactive:

`: command   / filter   ? help   Ctrl-P palette   ctx   ns   READ ONLY   N forwards`

`:` opens command mode with arrows, Tab, Enter and Esc. `/` opens live filtering. `?` opens contextual help. Ctrl-P opens command palette.

## 9. Command palette

Centered compact terminal overlay. Search commands/actions and display their shortcut. Disabled actions remain visible in read-only mode with a reason.

## 10. AI overlay

AI must never be a permanent panel in Fast Mode. `a` or Ctrl-A opens contextual actions:

- Explain this resource
- Analyze logs
- Analyze events
- Find root cause
- Recommend remediation
- Generate commands

Output: Root Cause → Evidence → Recommended Fix → Commands. Never auto-execute. Secrets/tokens must be redacted.

## 11. Production/safety UX

Destructive actions use explicit confirmation. Strong confirmation for Namespace, ClusterRole, CRD, StatefulSet, PVC, Secret and all production mutations. Read-only mode must be enforced server-side.

Production warning example:

```text
⚠ PRODUCTION ACTION
Delete pod/api-123?
[Y] Continue   [N] Cancel
```

## 12. Responsive behavior

At reduced width, shrink inspector, remove optional columns, collapse lower panels, but always preserve the main resource table and command bar. Do not hide primary Kubernetes state.

## 13. Interaction principles

Mouse is optional. Every important workflow must work by keyboard. j/k and arrows move; g/G first/last; Enter opens; Esc returns; Ctrl-R refreshes; `:` commands; `/` filters; `?` help; Ctrl-P palette.

## 14. Component direction

Prefer:

```text
apps/web/src/fast/
├── layout/        FastHeader ResourceNav Workspace StatusBar SplitLayout
├── resource/      ResourceTable ResourceHeader ResourceRow ResourceInspector ResourceTree
├── terminal/      LogViewer EventViewer YamlViewer DescribeViewer ExecTerminal
├── command/       CommandBar CommandAutocomplete CommandPalette HelpOverlay
├── filter/        FilterBar
├── ai/            AIActionOverlay
└── theme/         kubemv-terminal.css
```

Adapt to existing architecture rather than creating duplicate state or Kubernetes services.

## 15. Performance

Design for 1,000+ pods and 10,000+ resources. Use virtualization, memoized rows, keyed incremental updates, debounced filters and selective rendering. A single watch event must not rerender the entire table.

## 16. Visual acceptance checklist

- dark terminal background
- compact header
- visible cluster/context/namespace
- production indicator
- compact resource navigation
- dense resource table
- selected-row cursor
- live filter
- resource inspector
- labels/containers/owner chain
- keyboard actions
- logs/events workspace
- relationship tree
- command bar
- live connection state
- active port-forward status
- monospace terminal typography
- no dashboard-style card grid

## 17. Product distinction

Fast Mode is KubeMv's primary product surface. Visual Mode is secondary and may contain topology, charts and dashboards. Fast Mode should remain an operator terminal even as features are added.
