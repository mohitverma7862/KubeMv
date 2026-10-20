# KubeMv Fast Mode — Advanced TUI Design Specification

> Zero deviation from the Fast Mode reference: https://github.com/mohitverma7862/KubeMv/blob/main/docs/FAST-MODE-TUI-REFERENCE-SPEC.md and https://github.com/mohitverma7862/KubeMv/blob/main/docs/assets/kubemv-fast-mode-reference.svg. Match that specification and that image exactly: same logo, same colors, and same style.


> **Visual source of truth:** [`assets/kubemv-fast-mode-reference.svg`](assets/kubemv-fast-mode-reference.svg)
>
> **Implementation contract:** [`FAST-MODE-CURSOR-PROMPT.md`](FAST-MODE-CURSOR-PROMPT.md)
>
> **Important:** This specification describes the supplied KubeMv Fast Mode reference design. The reference image is not inspiration; it is the visual source of truth. Implementations must reproduce its layout, density, hierarchy, interaction model and terminal character. Do not reinterpret it as a dashboard, card grid or generic admin UI.

## 0. Source-of-truth contract

The reference asset is a **1536×1024** visual target showing the intended desktop/operator-console composition.

Use it for visual comparison during Phase 12 and acceptance testing.

Required behavior:

- Match the major panel geometry and hierarchy before adding feature-specific UI.
- Preserve the main resource table as the dominant surface.
- Preserve the right-side inspector and lower workspace/tree structure.
- Preserve the always-visible bottom command/status area.
- Preserve compact terminal typography and dense information presentation.
- Do not substitute rounded SaaS cards, large whitespace, dashboards, charts or decorative widgets.
- Do not move AI into a permanent panel.
- When responsive constraints make the exact geometry impossible, degrade by collapsing optional information—not by replacing the operator-console model.

## 1. Target composition

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ KubeMv | Fast Mode | cluster | user/context | namespace | PROD | clock      │
├──────────────────────────────────────────────────────────────────────────────┤
│ Pods | Deployments | Services | Nodes | Ingresses | Events | Jobs | ...     │
├──────────────────────────────────────────────────────┬───────────────────────┤
│                                                      │ Resource Inspector    │
│                  MAIN RESOURCE TABLE                 │ health               │
│                                                      │ containers           │
│                  dense + live + filterable           │ labels               │
│                                                      │ owner/actions        │
├──────────────────────────────────────────────────────┼───────────────────────┤
│ Logs | Events | YAML | Describe | Terminal | Rollout │ Resource Relationship │
│                                                      │ Deployment            │
│ terminal-style workspace                             │  └ ReplicaSet         │
│                                                      │     └ Pod             │
├──────────────────────────────────────────────────────┴───────────────────────┤
│ : command | / filter | ? help | Ctrl-P palette | ctx | ns | READONLY | ... │
└──────────────────────────────────────────────────────────────────────────────┘
```

## 2. Design language

Fast Mode is a **terminal operator console**, not a web dashboard.

- Dark operator-console default.
- Light and High Contrast themes are supported without changing layout.
- Monospace-first typography: JetBrains Mono, Cascadia Mono/Code, Consolas, monospace.
- Compact information density and minimal whitespace.
- Mostly square panels; maximum radius 4px.
- 1px borders; no heavy shadows.
- No giant cards, marketing sections, decorative gradients, excessive animation or oversized icons.
- Color communicates state, focus and safety rather than decoration.
- Keyboard-first; mouse is optional.
- Target 120×30 character-equivalent viewport; remain usable at 100×25.

### Reference palette

Use these values as the starting theme tokens and keep semantic meanings stable:

```text
--km-bg             #020817
--km-header         #020B19
--km-panel          #031022
--km-panel-alt      #061226
--km-border         #123E73
--km-border-active  #1761A8
--km-text           #B9C8E8
--km-muted          #88A4D0
--km-primary        #4CC9FF
--km-success        #18D59B
--km-warning        #F0DC55
--km-error          #FF5C68
--km-ai             #C58CFF
--km-production     #18D59B
--km-selection      #08437A
```

Colors may be theme-mapped for terminal capability, but the semantic hierarchy must remain equivalent.

## 3. Reference geometry

At the reference 1536×1024 composition:

- Global header: approximately 48px high.
- Resource navigation: approximately 50px high.
- Main workspace begins at approximately y=108px.
- Main table occupies approximately 1202px of the left side.
- Right inspector is approximately 288–334px wide depending on content.
- Main table/inspector workspace is approximately 555px high in the reference.
- Lower workspace contains the terminal workspace on the left and relationship/quick-help content on the right.
- Bottom command/status bar is always visible.
- Footer shortcuts remain discoverable without opening help.

These dimensions are a visual target, not hard-coded pixel requirements. Use responsive layout primitives while preserving the same proportions and hierarchy.

## 4. Global header

The header is compact and information-dense.

Reference content:

```text
⬡ KubeMv | Fast Mode (K9s-Parity) | ◈ prod-cluster | user: admin
| context: prod | namespace: production | PRODUCTION | clock
```

Required elements:

- KubeMv product identity.
- `Fast Mode (K9s-Parity)` mode indicator.
- Cluster name.
- User/context information where available.
- Namespace.
- Production safety badge.
- Connection state.
- Clock.

Connection states:

```text
● LIVE
◌ RECONNECTING n
✕ DISCONNECTED
```

Production contexts must visibly communicate that destructive actions require stronger confirmation.

## 5. Resource navigation

Directly below the header, render compact resource tabs. The reference shows:

```text
[K] Pods   [D] Deployments   [S] Services   [N] Nodes
[I] Ingresses   [E] Events   [J] Jobs   [C] CronJobs
[M] ConfigMaps   [R] Secrets   [>] More
```

Rules:

- Active resource has a subtle blue/cyan background and border.
- Tabs remain compact and keyboard reachable.
- Do not turn tabs into large rounded buttons.
- More resources are available through `>`/command mode without crowding the primary row.

## 6. Main resource table

The table is the visual center of Fast Mode.

Reference Pods columns:

```text
NAME | READY | STATUS | RESTARTS | AGE | IP | NODE | CPU | MEM
```

Required characteristics:

- Virtualized rendering.
- Sticky header.
- Compact 12–13px equivalent text.
- Approximately 31–34px maximum row height in the reference desktop layout.
- Strong column alignment.
- Sort indicators.
- Horizontal scrolling when necessary.
- Column visibility/configuration.
- Incremental row updates.
- Selection cursor at the left edge.
- Selected row uses a subtle dark-blue background rather than a card.
- Status uses semantic indicators.

Example status hierarchy:

```text
● Running             success
● Pending             warning
● CrashLoopBackOff    error
● Terminating         warning
```

The table must remain useful with hundreds or thousands of resources.

## 7. Filter and live state

The reference places the active filter inside the table header area:

```text
⚑ /filter: status=running                         ◷ 23s  ● Live
```

Required:

- `/` enters filter mode.
- Active filter remains visible.
- Filter can be cleared quickly.
- Live/watch age is visible.
- Connection/watch state is visible.
- Filtering must not require leaving the resource table.

Supported filter directions should include label, field, status, fuzzy and regex filtering as appropriate to the existing architecture.

## 8. Resource intelligence inspector

The right panel is a compact inspector, not a dashboard card.

For the selected resource show, when available:

- Health/status.
- Namespace.
- Node.
- IP/address.
- Readiness.
- Restarts.
- Age.
- QoS.
- CPU/memory.
- Containers and images.
- Ports/mounts/resource limits.
- Labels.
- Owner chain.
- Contextual actions/shortcuts.

Never display Secret values, tokens or credentials.

The inspector must update with the selected table row without rebuilding the whole page.

## 9. Lower terminal workspace

The reference uses a lower workspace with terminal-oriented tabs:

```text
[Logs] [Events] [YAML] [Describe] [Terminal] [Rollout] [History]
```

### Logs

- Terminal-style rendering.
- Timestamp and level alignment.
- Follow/pause.
- Previous logs.
- Container selector.
- `/` search.
- `n`/`N` next/previous match.
- Wrap toggle.
- Ctrl-C stops streaming.
- Clear/download where supported.

### Events

Compact timeline/table with Normal/Warning/Error severity and timestamps.

### YAML

Syntax highlighting, line numbers, search, copy/download and safe editing integration where supported.

### Describe

Structured Kubernetes metadata, conditions, containers and important events.

### Terminal / Exec

xterm-style interactive terminal over the existing secure backend/WebSocket abstraction, including resize and multi-container support.

### Rollout / History

Resource-specific rollout status, revisions and history without leaving Fast Mode.

## 10. Resource relationship tree

The bottom-right area visualizes ownership/dependency relationships compactly:

```text
Deployment
 └─ ReplicaSet
     └─ Pod
         └─ Container
```

Requirements:

- Keyboard navigable.
- Selecting/opening a node navigates to its matching resource view.
- Show enough context to understand ownership without becoming a full graph dashboard.
- Preserve compact terminal tree styling.

## 11. Quick Help

The reference includes persistent lightweight help near the lower workspace. It should expose the most important keys without consuming the whole screen.

Reference-level shortcuts:

```text
j/k       move
 g/G      first/last
Enter     open
Esc       back
Ctrl-R    refresh
/         filter
:         command
?         help
Ctrl-P    command palette
```

Full help remains available through `?`.

## 12. Command/status bar

Always-visible at the bottom.

Reference pattern:

```text
: command   / filter   ? help   Ctrl-P palette
ctx context   ns namespace   R readonly   N forwards
KubeMv version
```

Rules:

- `:` opens command mode.
- `/` opens live filter mode.
- `?` opens contextual help.
- Ctrl-P opens the command palette.
- `ctx` changes context.
- `ns` changes namespace.
- Read-only state is explicit.
- Active port-forwards are visible.
- The bar must remain compact enough for 100×25.

## 13. Command palette

Centered compact terminal overlay.

- Search commands/actions.
- Show shortcut beside each command.
- Respect current resource/context/permissions.
- Disabled actions remain discoverable with the reason.
- Use the same command registry as `:` mode and help; do not duplicate metadata.

## 14. Keyboard interaction model

Every important operator workflow must work without a mouse.

```text
j / k           move selection
↑ / ↓           move selection
g / G           first / last
Enter           open/detail
Esc             back/close
Ctrl-R          refresh
/               filter/search
:               command mode
?               help
Ctrl-P          command palette
a / Ctrl-A      contextual AI overlay
```

Focus must be visible. Text inputs and overlays may temporarily consume global shortcuts, but Esc must provide a predictable return path.

## 15. AI overlay

AI is contextual and must never be a permanent Fast Mode panel.

`a` or Ctrl-A opens a compact overlay with actions such as:

- Explain resource.
- Analyze logs/events.
- Find likely root cause.
- Recommend remediation.
- Generate safe commands.

Output structure:

```text
Root Cause
Evidence
Recommended Fix
Commands
```

Never auto-execute generated commands. Redact secrets, tokens and credentials before sending context to an AI service.

## 16. Production and safety UX

All mutations go through the existing `ActionManager`/approved action boundary.

- Destructive operations require explicit confirmation.
- Production operations require stronger confirmation.
- Namespace, ClusterRole, CRD, StatefulSet, PVC, Secret and other high-impact resources require elevated warnings.
- Read-only mode must be enforced server-side, not only by hiding buttons.
- Disabled actions explain why they are unavailable.

## 17. Responsive behavior

When the terminal/browser width decreases:

1. Preserve the main resource table.
2. Preserve the command/status bar.
3. Shrink the inspector.
4. Remove optional columns before removing core Kubernetes state.
5. Collapse lower workspace panels when necessary.
6. Keep keyboard navigation intact.

Do not transform Fast Mode into a mobile/dashboard card layout.

## 18. Architecture/component direction

Adapt to the existing architecture rather than creating duplicate Kubernetes services or state.

```text
apps/web/src/fast/
├── layout/        FastHeader ResourceNav Workspace StatusBar SplitLayout
├── resource/      ResourceTable ResourceHeader ResourceRow ResourceInspector ResourceTree
├── terminal/      LogViewer EventViewer YamlViewer DescribeViewer ExecTerminal
├── command/       CommandBar CommandAutocomplete CommandPalette HelpOverlay
├── filter/        FilterBar FilterEngine
├── ai/            AIActionOverlay
└── theme/         kubemv-terminal.css
```

Backend/service responsibilities must remain centralized in the existing Kubernetes abstraction and mutation/action layer.

## 19. Performance requirements

Design and test for:

- 1,000+ pods.
- 10,000+ resources.
- High-frequency Kubernetes watch updates.
- Fast keyboard navigation.
- Large log streams.

Required techniques:

- Virtualized tables.
- Memoized/keyed rows.
- Incremental row patches.
- Debounced filtering.
- Selective rendering.
- No full-table rerender for a single watch event.
- No full-page rerender when selection changes.

## 20. Exact visual acceptance checklist

Before Fast Mode is considered visually complete, compare the implementation against `docs/assets/kubemv-fast-mode-reference.svg`.

### Layout

- [ ] Compact 48px-class header.
- [ ] Compact resource navigation directly below header.
- [ ] Main table is the dominant area.
- [ ] Right inspector is visible beside the table.
- [ ] Lower terminal workspace is visible.
- [ ] Relationship tree is visible.
- [ ] Bottom command/status bar is always visible.

### Styling

- [ ] Dark navy/black terminal background.
- [ ] Thin blue borders.
- [ ] Cyan primary accent.
- [ ] Green success state.
- [ ] Yellow warning state.
- [ ] Red error state.
- [ ] Purple contextual/AI accent.
- [ ] Monospace typography.
- [ ] Compact spacing.
- [ ] Minimal border radius.
- [ ] No dashboard cards.
- [ ] No oversized icons.
- [ ] No decorative gradients/animations.

### Information hierarchy

- [ ] Cluster/context/namespace visible.
- [ ] Production state visible.
- [ ] Live/watch state visible.
- [ ] Filter visible when active.
- [ ] Dense Pods/resource table.
- [ ] CPU/MEM/Node information available where supported.
- [ ] Selected-row cursor and highlight.
- [ ] Resource inspector.
- [ ] Owner relationship tree.
- [ ] Logs/Events/YAML/Describe/Terminal workspace.
- [ ] Quick Help.
- [ ] Command/status bar.

### Interaction

- [ ] Keyboard-only operation works.
- [ ] j/k and arrows navigate.
- [ ] Enter opens.
- [ ] Esc returns.
- [ ] `/` filters.
- [ ] `:` commands.
- [ ] `?` help.
- [ ] Ctrl-P palette.
- [ ] Context/namespace switching works.
- [ ] Read-only state is obvious and enforced.
- [ ] Production destructive actions require stronger confirmation.

## 21. Product distinction

Fast Mode is KubeMv's primary operator surface.

Visual Mode may provide topology, charts and dashboards, but Fast Mode must remain:

```text
K9s speed
+ modern terminal IDE density
+ Kubernetes resource intelligence
+ observability
+ safe operations
+ contextual AI
```

It must never regress into:

```text
web dashboard + some keyboard shortcuts
```

## 22. Implementation rule

**Do not reinterpret the reference design. Reproduce its visual language and information hierarchy first, then implement the underlying Kubernetes functionality behind it.**