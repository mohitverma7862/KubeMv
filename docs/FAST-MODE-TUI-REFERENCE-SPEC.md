# KubeMv Fast Mode — Exact TUI Reference Specification

**Status:** Design source of truth  
**Reference image:** `docs/assets/kubemv-fast-mode-reference.svg`  

## Objective

Build KubeMv Fast Mode to match the supplied KubeMv terminal UI reference—not the previous generic dashboard concept.

The target experience is **K9s operator speed + dense terminal IDE layout + Kubernetes resource intelligence + safe operations + contextual AI**.

Fast Mode must feel like an operator console first. It must **not** feel like a web dashboard with keyboard shortcuts.

## Visual source of truth

The reference establishes this exact hierarchy:

1. Global header
2. Resource navigation strip
3. Main resource table
4. Selected-resource inspector
5. Lower terminal workspace
6. Resource relationship tree
7. Quick Help
8. Command/status bar

Do not introduce persistent panels that change this hierarchy.

## Canvas and density

- Reference composition: 1536×1024
- Minimum supported terminal-like viewport: 100×25
- Preferred compact viewport: 120×30
- No page-level scrolling in normal operation
- Dense information presentation
- 1px panel borders
- Maximum 4px radius
- No oversized cards
- No marketing/dashboard spacing
- No decorative gradients, glassmorphism, or large charts
- Hide secondary metadata before shrinking primary controls

## Global header

Approximately 48px high.

Left: KubeMv, separator, `Fast Mode (K9s-Parity)`.

Context: cluster, user, context, namespace.

Safety: highly visible `PRODUCTION` badge. Production state must affect destructive confirmations.

Right: connection/time indicator and live clock.

Header values are live state, not decoration.

## Resource navigation

Second horizontal strip with:

`Pods`, `Deployments`, `Services`, `Nodes`, `Ingresses`, `Events`, `Jobs`, `CronJobs`, `ConfigMaps`, `Secrets`, `More`.

Every resource has a compact keyboard mnemonic. Active resource is highlighted. Mouse is optional. CRDs are dynamically discoverable.

## Main resource table

The table is the primary workspace.

Pod columns:

`NAME | READY | STATUS | RESTARTS | AGE | IP | NODE | CPU | MEM`

Requirements:

- sticky header
- dense rows
- selected-row highlight
- semantic status indicator + text
- sortable/configurable columns
- resize where practical
- horizontal navigation
- virtualization
- stable row keys
- incremental row patches
- never redraw thousands of rows for one watch event

Health states: Running, Pending, CrashLoopBackOff, Terminating, Failed, Unknown.

## Filter

Filter lives inside the resource workspace. Reference: `/filter: status=running`.

`/` enters filter mode; Esc exits. Support live filtering, fuzzy matching, field filters, label/status/namespace filters, optional regex, and persistence across drill-down/back navigation.

## Selected-resource inspector

Right-side contextual panel displays selected resource name/status and, for Pods: namespace, node, IP, readiness, restarts, age, QoS, CPU, memory, containers, labels and owner.

Owner chain example: `Deployment: api`.

**Security:** never expose Secret values, service-account tokens, bearer tokens, private keys, or Kubernetes credential material. Secret resources show metadata and hidden values unless an explicitly authorized secure workflow exists.

## Contextual actions

Reference Pod shortcuts:

`l` Logs, `s` Shell, `a` Attach, `d` Describe, `e` Edit, `y` YAML, `f` Port Forward, `c` Container, `D` Delete.

All mutations must go through one backend `ActionManager`. UI components may not implement Kubernetes mutations directly.

## Lower terminal workspace

Required tabs:

`Logs | Events | YAML | Describe | Terminal | Rollout | History`

Active tab is clearly highlighted.

Logs support follow, Ctrl-C stop, current/previous container, container selection, search, next/previous match, timestamps, wrapping, pause, clear and bounded rendering for huge streams.

Terminal should use a proper terminal implementation such as xterm.js where appropriate and support exec, resize, interactive stdin/stdout, multi-container selection and reconnect/disconnect handling.

## Resource relationship tree

Show Kubernetes ownership:

```text
Deployment (api)
└─ ReplicaSet (api-7d9f8c9b6f)
   └─ Pod (api-7d9f8c9b6f-2kq7m)
      └─ Containers
         └─ api
```

Selectable, keyboard navigable, drill-down capable, based on Kubernetes metadata and lazy-loaded when necessary.

## Quick Help

Persistent compact contextual help. Reference:

```text
j/k       Select row
g/G       First/Last
Enter     Drill down
Esc       Back
Ctrl-R    Refresh
/         Filter
:         Command mode
?         Help
Ctrl-P    Palette
```

Help should adapt to the active workspace.

## Command bar

Bottom bar remains available when no modal owns focus. Reference:

`: command   / filter   ? help   ctx context   ns namespace   R readonly   3 forwards`

Support command mode, autocomplete, Tab completion, history, Esc cancellation, aliases, validation and clear errors. Do not turn command mode into unrestricted shell execution.

## Command registry

Use one registry for command parsing, aliases, autocomplete, help, keyboard metadata and documentation. Examples: `:q`, `:quit`, `:pods`, `:deploy`, `:svc`, `:nodes`, `:events`, `:ns`, `:ctx`, `:readonly`.

## Keyboard model

Primary keyboard behavior:

```text
j/k       move selection
g/G       first / last
Enter     drill down
Esc       back/cancel
/         filter
:         command
?         help
Ctrl-P    command palette
Ctrl-R    refresh
```

Focus ownership: modal → text input → terminal → global dispatcher. Avoid browser shortcuts accidentally triggering destructive operations.

## Production safety

Show `PRODUCTION` prominently. Destructive actions require explicit confirmation with resource and namespace. Bulk delete/edit requires stronger confirmation. Never allow a single accidental key to perform a high-risk production operation.

Example:

```text
DELETE POD
production / api-7d9f8c9b6f-2kq7m
This action cannot be automatically undone.
Type DELETE to confirm.
```

## Readonly

Readonly is a first-class state. Show `R readonly`. Disable destructive operations and explain why. Enforcement must exist server-side; hiding UI controls is not sufficient.

## Live connection

Header states: Live, Reconnecting, Disconnected, Error.

Use keyed resource updates, resourceVersion-aware backend watch, reconnect/backoff and stale-state indication. Do not redraw the entire table for one update.

## Port forwarding

Bottom status bar shows active forwards, e.g. `3 forwards`. Selecting it opens a compact list with target/local ports and status, with stop/restart controls.

## Kubernetes architecture

Reuse shared abstractions. Target separation:

```text
internal/k8s/
  service/
  discovery/
  watch/
  actions/
internal/session/

apps/web/src/fast/
  registry/
  nav/
  filter/
  selection/
  state/
  views/
  overlays/
  terminal/
  inspector/
  tree/
  theme/
```

Do not create a second Kubernetes client stack for Fast Mode.

## AI design

AI is contextual, never a permanent sidebar. Open with `a` or `Ctrl-A`.

Gather selected resource, status, events, current/previous logs, container state, owner chain, referenced ConfigMap names, referenced Secret names and node conditions. Redact secrets/tokens.

Response:

```text
ROOT CAUSE
...
EVIDENCE
...
RECOMMENDED FIX
...
COMMANDS
...
```

AI never auto-executes commands; mutations always require explicit confirmation.

## Visual system

Recommended font stack: JetBrains Mono, Cascadia Mono, Consolas, monospace.

Approximate sizes: header 14–16px; table 12–14px; metadata 11–13px; command/status 11–13px.

Palette semantics: dark navy/black base, cyan/blue focus, green healthy, yellow warning, red failure, purple AI/context.

No huge rounded cards, large gradients, oversized shadows, decorative illustrations, excessive whitespace or marketing copy.

## Responsive behavior

At 120×30 retain table, inspector and command bar while reducing secondary metadata.

At 100×25 reduce columns and collapse lower panels intelligently; inspector may become a compact overlay. Preserve keyboard navigation.

Priority at very small sizes:

1. resource table
2. selected resource
3. command bar
4. safety/connection state

## Performance

Target 1,000+ pods and 10,000+ resources with virtualization, keyed updates, bounded logs, incremental watch updates and no blocking UI calls.

## Unwanted Fast Mode flows to remove

- dashboard-first landing page
- oversized KPI cards
- marketing home screen
- permanent AI sidebar
- duplicate Kubernetes data layer
- separate mutation implementation per view
- unnecessary navigation pages between resource and details
- modal chains for routine read-only operations
- decorative animations
- excessive confirmation for harmless reads
- browser-like page scrolling
- hidden production state
- exposing Secret contents by default

Visual Mode may retain dashboard capabilities, but Fast Mode remains the operator console.

## Visual acceptance

The running UI must immediately resemble this structure:

```text
┌──────────────────────────────────────────────────────────────────────────┐
│ KubeMv | Fast Mode | cluster | user | context | namespace | PRODUCTION  │
├──────────────────────────────────────────────────────────────────────────┤
│ Pods Deployments Services Nodes Ingresses Events Jobs ...                │
├───────────────────────────────────────────────┬──────────────────────────┤
│ RESOURCE TABLE                                │ Selected Resource        │
│                                               │ Inspector + Actions      │
├───────────────────────────────────────────────┴──────────────────────────┤
│ Logs | Events | YAML | Describe | Terminal | Rollout | History           │
├──────────────────────────────────────┬────────┬──────────────────────────┤
│ Terminal/log output                  │ Tree   │ Quick Help               │
├──────────────────────────────────────┴────────┴──────────────────────────┤
│ : command / filter ? help ctx ns readonly forwards                       │
└──────────────────────────────────────────────────────────────────────────┘
```

If Fast Mode resembles a conventional web dashboard, it fails visual acceptance.

## Cursor execution rules

1. Read this specification before implementing Fast Mode UI.
2. Inspect the repository architecture.
3. Reuse existing Kubernetes services.
4. Implement one phase at a time.
5. Run backend/frontend tests after every phase.
6. Fix failures before stopping.
7. Report files changed, reused code, removed code, remaining gaps and test results.
8. Stop after each phase.

Suggested validation:

```bash
go test ./...
cd apps/web
npm test
npm run build
```

## Definition of done

Fast Mode is complete only when the reference hierarchy is matched, K9s-like keyboard workflows work, resource switching/filtering is fast, live watch is stable, logs/exec/describe/YAML work, relationships are visible, production safety and server-side readonly are enforced, AI is contextual and non-executing, Secrets are protected, 1,000+ resources remain usable, and Fast Mode remains separate from Visual Mode.

**This document and `docs/assets/kubemv-fast-mode-reference.svg` are the visual and behavioral source of truth for Cursor implementation.**
