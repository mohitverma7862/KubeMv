# KubeMv Fast Mode — Master Cursor Execution Prompt

> Zero deviation from the Fast Mode reference: https://github.com/mohitverma7862/KubeMv/blob/main/docs/FAST-MODE-TUI-REFERENCE-SPEC.md and https://github.com/mohitverma7862/KubeMv/blob/main/docs/assets/kubemv-fast-mode-reference.svg. Match that specification and that image exactly: same logo, same colors, and same style.


> **Read first:** `docs/FAST-MODE-UI-DESIGN.md`
>
> **Visual source of truth:** `docs/assets/kubemv-fast-mode-reference.svg`

Paste this prompt into Cursor after the UI design document is available.

## ROLE

You are rebuilding KubeMv Fast Mode into a K9s-first Kubernetes operator terminal with a modern, dense, capable TUI.

The supplied reference design is the visual source of truth. Do not reinterpret it as a dashboard or invent a different layout. Reproduce its information hierarchy, density, panel structure, keyboard model and terminal styling, then connect real Kubernetes functionality behind it.

## HARD RULES

1. Read `docs/FAST-MODE-UI-DESIGN.md` before changing code.
2. Treat `docs/assets/kubemv-fast-mode-reference.svg` as the visual source of truth.
3. Do not restyle the existing dashboard into Fast Mode. Fast Mode is a new terminal operator shell; the existing dashboard remains Visual Mode/secondary.
4. Do not add dashboard cards, marketing sections, large whitespace, gradients, excessive rounded containers, decorative animation or oversized icons.
5. Do not reinterpret the reference layout. Preserve header → resource nav → dominant table + inspector → lower terminal workspace/tree → command/status bar.
6. Keyboard-first is mandatory. Mouse is optional.
7. Do not duplicate Kubernetes logic. Reuse the existing internal Kubernetes services, resource abstractions and ActionManager.
8. All mutations must use the approved ActionManager/action boundary.
9. Do not create a permanent AI panel in Fast Mode.
10. Work in one phase at a time.
11. After every implementation phase run tests/build, fix failures, report the phase, and STOP until explicitly told to continue.
12. Prefer small reviewable diffs. List reused and deleted components.
13. Do not claim K9s parity from visual similarity; parity means equivalent operator workflows and keyboard interactions.
14. Do not hard-code fake Kubernetes behavior once backend integration begins.
15. Never display Secret values, tokens or credentials in the UI or AI context.
16. Preserve the existing application's working functionality unless a change is explicitly required by this Fast Mode design.

## PHASE 0 — AUDIT ONLY

Inspect the complete repo:

- `cmd/kubemv`
- `cmd/server`
- `internal`
- `backend`
- `api`
- `apps/web`
- desktop
- tests
- AI code
- documentation

Determine:

- existing frontend architecture
- reusable Kubernetes services
- list/watch abstractions
- discovery/CRD support
- logs/previous logs
- exec/attach
- describe/YAML/edit/delete
- port-forward
- context/namespace handling
- ActionManager/mutation boundary
- existing routing
- existing dashboard code that should become Visual Mode
- existing keyboard/input abstractions
- current AI integration

Produce/update:

- `docs/FAST-MODE-AUDIT.md`
- `docs/FAST-MODE-K9S-PARITY.md`

The parity matrix must compare KubeMv against K9s across startup, resource navigation, `:` commands, aliases, `/` filtering, j/k navigation, context/namespace, breadcrumbs, tables, watch, logs, previous logs, exec, attach, describe, YAML, edit, delete, port-forward, events, rollout, node operations, multi-select, help/hotkeys, custom views, CRDs, multi-cluster and reconnect behavior.

Classify each as `PASS`, `PARTIAL`, `MISSING`, or `DIFFERENT-BETTER`, with file paths and target implementation phase.

**No code changes in Phase 0.**

STOP.

## PHASE 1 — UI ARCHITECTURE

Create/adapt `apps/web/src/fast/` without duplicating backend Kubernetes state:

```text
fast/
├── layout/
├── resource/
├── terminal/
├── command/
├── filter/
├── ai/
└── theme/
```

Define component contracts and theme tokens based on `FAST-MODE-UI-DESIGN.md`.

The design must support:

- compact global header
- resource navigation
- main resource table
- right inspector
- lower terminal workspace
- relationship tree
- quick help
- bottom command/status bar

No backend rewrite.

STOP after tests/build.

## PHASE 2 — REFERENCE-MATCHED TERMINAL SHELL

Implement Fast Mode as the default operator console.

Match the supplied reference composition:

1. approximately 48px-class global header
2. approximately 50px resource navigation row
3. dominant left resource table
4. right resource inspector
5. lower Logs/Events/YAML/Describe/Terminal/Rollout/History workspace
6. lower relationship tree/quick-help area
7. always-visible bottom command/status bar

Header must expose:

- KubeMv
- Fast Mode (K9s-Parity)
- cluster
- user/context
- namespace
- production state
- connection state
- clock

Resource nav must include the primary Kubernetes resources shown in the reference.

Use realistic mock data initially, but keep interfaces ready for real backend data.

Target 120×30 and remain usable at 100×25.

Implement the global keyboard dispatcher with focus-aware overlays/inputs.

STOP after tests/build.

## PHASE 3 — RESOURCE TABLE + NAVIGATION

Implement the dominant resource table:

- virtualization
- sticky headers
- compact rows
- selected-row cursor/highlight
- sorting
- column configuration/resizing
- horizontal scroll
- incremental keyed updates
- resource counts
- active filter indicator
- live/watch state

For Pods, support the reference-level information model:

`NAME | READY | STATUS | RESTARTS | AGE | IP | NODE | CPU | MEM`

Implement:

- j/k
- arrows
- g/G
- Enter
- Esc
- Ctrl-R

Navigation history must preserve selection/filter/scroll state where appropriate.

STOP after tests/build.

## PHASE 4 — COMMAND/FILTER UX

Implement:

- `:` command mode
- autocomplete
- Tab completion
- command aliases
- `/` live filter
- field/label/status filters
- regex/fuzzy filtering as supported
- `?` contextual help
- Ctrl-P command palette

Command, shortcut and help metadata must come from shared registries. Do not duplicate shortcut definitions across components.

STOP after tests/build.

## PHASE 5 — LIVE KUBERNETES DATA

Connect the existing Kubernetes services.

Use existing watch/resourceVersion/reconnect behavior where available.

Frontend requirements:

- keyed row patching
- no whole-table rebuild for a single event
- visible LIVE state
- RECONNECTING state with retry information
- DISCONNECTED state
- watch age/refresh indicator

STOP after tests/build.

## PHASE 6 — RESOURCE INTELLIGENCE

Implement the right-side inspector and lower relationship tree.

Inspector shows, when available:

- health/status
- namespace
- node
- IP/address
- readiness
- restarts
- age
- QoS
- CPU/memory
- containers/images
- ports/mounts/limits
- labels
- owner chain
- contextual actions

The relationship tree should support compact ownership navigation:

`Deployment → ReplicaSet → Pod → Container`

Never expose Secret values.

STOP after tests/build.

## PHASE 7 — TERMINAL WORKSPACE

Implement:

- Logs
- Events
- YAML
- Describe
- Terminal/Exec
- Rollout
- History

Logs must support:

- follow
- pause
- previous
- timestamps
- container selection
- search
- n/N match navigation
- wrap
- Ctrl-C to stop streaming
- clear/download where supported

Exec must use the existing secure backend/WebSocket path where available, with resize and multi-container support.

STOP after tests/build.

## PHASE 8 — ACTION UX

Implement pod/deployment/service/node actions through ActionManager only.

Support the relevant operator workflows from the parity matrix, including where supported:

- describe
- edit
- YAML
- delete
- logs
- previous logs
- exec
- attach
- port-forward
- rollout actions
- node operations

Destructive operations require confirmation.

Production operations require stronger confirmation.

Port-forward status must appear in the bottom status bar.

STOP after tests/build.

## PHASE 9 — CONTEXT / NAMESPACE / READONLY

Implement:

- `:ctx`
- `:ns`
- all-namespaces mode
- production context detection
- explicit READ ONLY state
- server-enforced readonly mode

UI must not merely hide forbidden actions; the backend must reject unauthorized mutations.

STOP after tests/build.

## PHASE 10 — AI OVERLAY

Remove any permanent Fast Mode AI panel.

`a`/Ctrl-A opens contextual AI actions.

Gather only safe context through the existing backend service:

- selected resource status
- container state
- current/previous logs
- events
- owner chain
- referenced ConfigMap/Secret names, never secret values
- node conditions

Redact secrets/tokens/credentials.

Return:

```text
Root Cause
Evidence
Recommended Fix
Commands
```

Never auto-execute commands.

STOP after tests/build.

## PHASE 11 — SCALE / PERFORMANCE

Profile:

- 1,000+ pods
- 10,000+ resources
- high-frequency watch updates
- large log streams

Verify:

- virtualization
- keyed updates
- memoization
- debounced filtering
- selective rendering
- fast j/k navigation
- no whole-table rerender from a single watch event
- no whole-page rerender when selection changes
- 120×30 layout
- 100×25 layout

Fix overflow, jitter and unnecessary rerenders.

STOP after tests/build.

## PHASE 12 — VISUAL REFERENCE POLISH

This phase is mandatory.

Compare the implementation directly against:

`docs/assets/kubemv-fast-mode-reference.svg`

and the detailed rules in:

`docs/FAST-MODE-UI-DESIGN.md`

Check and correct:

- overall panel geometry
- header height/density
- resource navigation spacing
- table width/height
- inspector width
- lower workspace proportions
- relationship tree placement
- bottom command bar
- typography
- line heights
- border thickness
- color hierarchy
- selection highlight
- status colors
- filter presentation
- live indicator
- production badge
- shortcut/help placement
- focus states

**Do not add new visual concepts during this phase. Match the reference.**

STOP after tests/build.

## PHASE 13 — ACCEPTANCE

Run the complete keyboard-only flow:

```text
Launch
→ Pods
→ j/k
→ Enter
→ Logs
→ / search
→ Esc
→ Describe
→ YAML
→ Esc
→ command mode
→ deployment
→ filter
→ deployment selection
→ action
→ confirmation
→ :ns
→ :ctx
→ Ctrl-P
→ ?
→ Esc
```

Verify:

- visual reference alignment
- core Kubernetes resources
- K9s operator workflows
- safety confirmations
- readonly enforcement
- reconnect behavior
- command registry
- resource registry
- performance
- no Secret value leakage
- no permanent AI panel
- no dashboard-style Fast Mode regressions

Update documentation and changelog as appropriate.

## TEST COMMANDS

After each implementation phase run the repository's appropriate test suite. At minimum, if present:

```bash
go test ./...
cd apps/web && npm test && npm run build
```

Fix all failures before reporting the phase.

## PHASE REPORT

At the end of every phase report:

- What changed, with files
- What was reused
- What was deleted
- Keyboard shortcuts added
- Tests and results
- Known gaps against `FAST-MODE-UI-DESIGN.md`
- Known gaps against `kubemv-fast-mode-reference.svg`
- Proposed next phase

Then STOP.

## FINAL DESIGN PRINCIPLE

KubeMv Fast Mode must feel like:

`K9s speed + modern terminal IDE density + Kubernetes resource intelligence + observability + safe operations + contextual AI`

It must NOT feel like:

`web dashboard + some keyboard shortcuts`

**The reference image defines what the UI should look and feel like. The existing architecture defines how the functionality should work. Do not mix those responsibilities or invent a third design.**