# KubeMv Fast Mode — Master Cursor Execution Prompt

Paste this prompt into Cursor after reading `docs/FAST-MODE-UI-DESIGN.md`.

## ROLE

You are rebuilding KubeMv Fast Mode into a K9s-first Kubernetes operator terminal with a modern, denser and more capable UI. The UI design document is the visual source of truth. Existing architecture must be inspected before modification.

## HARD RULES

1. Do not restyle the existing dashboard into Fast Mode. Build the terminal operator shell; existing dashboard becomes Visual Mode/secondary.
2. Do not add dashboard cards, gradients, oversized charts, marketing UI, excessive rounded containers or decorative animation.
3. Keyboard-first is mandatory. Mouse is optional.
4. Do not duplicate Kubernetes logic. Reuse the existing internal K8s services and ActionManager.
5. Do not bypass ActionManager for mutations.
6. Do not create a permanent AI panel in Fast Mode.
7. Work in one phase at a time. After every phase run tests/build, fix failures, report, and STOP until explicitly told to continue.
8. Prefer small reviewable diffs. List reused and deleted components.
9. Do not claim K9s parity from visual similarity; parity means operator workflows and keyboard interactions.

## PHASE 0 — AUDIT ONLY

Inspect the complete repo: `cmd/kubemv`, `cmd/server`, `internal`, `backend`, `api`, `apps/web`, desktop, tests, AI code and docs. Produce `docs/FAST-MODE-AUDIT.md` with architecture, reusable code paths, conflicting UI code, K8s list/watch/logs/exec/port-forward/discovery abstractions, risks and gaps.

Also produce `docs/FAST-MODE-K9S-PARITY.md` comparing KubeMv against K9s across startup, resource navigation, `:` commands, aliases, `/` filtering, j/k navigation, context/namespace, breadcrumbs, tables, watch, logs, previous logs, exec, attach, describe, YAML, edit, delete, port-forward, events, rollout, node operations, multi-select, help/hotkeys, custom views, CRDs, multi-cluster and reconnect behavior. Classify each PASS/PARTIAL/MISSING/DIFFERENT-BETTER with file paths and target phase. No code changes.

STOP.

## PHASE 1 — UI ARCHITECTURE

Create/adapt `apps/web/src/fast/` with layout, resource, terminal, command, filter, AI and theme boundaries. Reuse existing state/API services. Define component contracts and theme tokens. No backend rewrite.

STOP after tests/build.

## PHASE 2 — TERMINAL SHELL

Implement the Fast Mode route as the default operator console. Build:

- compact global header
- cluster/context/namespace/production/connection/clock state
- resource navigation tabs
- dense main workspace
- right resource inspector
- lower terminal workspace
- relationship tree
- bottom command/status bar

Use realistic mock data initially. Target 120x30 and 100x25. Implement global keyboard dispatcher with focus-aware overlays/inputs.

STOP after tests/build.

## PHASE 3 — RESOURCE TABLE + NAVIGATION

Implement virtualized resource table, sticky headers, compact rows, selection cursor, sorting, column configuration/resizing, horizontal scroll and incremental updates. Add j/k, arrows, g/G, Enter, Esc and Ctrl-R. Navigation history must preserve selection/filter/scroll state.

STOP after tests/build.

## PHASE 4 — COMMAND/FILTER UX

Implement `:` command mode, autocomplete, Tab completion, aliases, `/` live filter, regex/fuzzy/label/status filters, `?` contextual help and Ctrl-P command palette. Generate help/palette from registries instead of duplicating shortcut metadata.

STOP after tests/build.

## PHASE 5 — LIVE KUBERNETES DATA

Connect the existing Kubernetes services. Use watch streaming with reconnect/resourceVersion behavior already supported by backend. Frontend must patch keyed rows instead of rebuilding the table. Show LIVE/RECONNECTING/DISCONNECTED states.

STOP after tests/build.

## PHASE 6 — RESOURCE INTELLIGENCE

Implement selected-resource inspector and relationship tree. Inspector shows health, namespace, node, IP, readiness, restarts, age, QoS, CPU/memory, containers, labels, owner chain and contextual actions. Never display secret values.

STOP after tests/build.

## PHASE 7 — TERMINAL WORKSPACE

Implement Logs, Events, YAML, Describe, Terminal/Exec, Rollout and History tabs using existing backend operations. Logs support follow, previous, timestamps, container selection, search n/N, pause, clear, download, wrap and Ctrl-C. Exec uses xterm.js/WebSocket with resize and multi-container support.

STOP after tests/build.

## PHASE 8 — ACTION UX

Implement pod, deployment, service and node actions through ActionManager only. Add contextual shortcut bars. Port-forward status appears in the bottom bar. Destructive actions require confirmation. Production gets stronger warnings.

STOP after tests/build.

## PHASE 9 — CONTEXT/NAMESPACE/READONLY

Implement `:ctx`, `:ns`, all-namespaces `0`, production context flags and server-enforced readonly mode. UI must clearly show READ ONLY and disable actions with reasons.

STOP after tests/build.

## PHASE 10 — AI OVERLAY

Remove any permanent Fast Mode AI panel. `a`/Ctrl-A opens contextual AI actions. Gather context through the existing backend service: resource status, container state, current/previous logs, events, owner chain, referenced ConfigMap/Secret names and node conditions. Redact secrets/tokens. Return Root Cause, Evidence, Recommended Fix and Commands. Never auto-execute.

STOP after tests/build.

## PHASE 11 — SCALE/PERFORMANCE

Profile 1,000+ pods and 10,000+ resources. Verify virtualization, keyed updates, memoization, filtering and watch behavior. Verify 120x30 and 100x25 layouts. Fix overflow, jitter and unnecessary rerenders.

STOP after tests/build.

## PHASE 12 — POLISH

Compare the implementation with `docs/FAST-MODE-UI-DESIGN.md` and the supplied design references. Correct alignment, density, typography, borders, status hierarchy, spacing, focus states and keyboard affordances. Do not add decorative UI merely to make it look more impressive.

STOP after tests/build.

## PHASE 13 — ACCEPTANCE

Run the complete keyboard-only flow:

Launch → Pods → j/k → Enter → Logs → `/` search → Esc → Describe → YAML → Esc → `:deploy` → filter → deployment → action → confirmation → `:ns` → `:ctx` → Ctrl-P → `?` → Esc back.

Verify core resources, safety, reconnect behavior, readonly enforcement, command registry, resource registry and UI performance. Add/update documentation and changelog.

## TEST COMMANDS

After each implementation phase run the repository's appropriate test suite. At minimum, if present:

```bash
go test ./...
cd apps/web && npm test && npm run build
```

Fix failures before reporting the phase.

## PHASE REPORT

At the end of every phase report:

- What changed, with files
- What was reused
- What was deleted
- Keyboard shortcuts added
- Tests and results
- Known gaps against the design/spec
- Proposed next phase

Then STOP.

## FINAL DESIGN PRINCIPLE

KubeMv must feel like:

`K9s speed + modern terminal IDE + Kubernetes resource intelligence + observability + safe operations + contextual AI`

It must NOT feel like:

`web dashboard + some keyboard shortcuts`.
