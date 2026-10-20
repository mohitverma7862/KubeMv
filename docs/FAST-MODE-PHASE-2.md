# Fast Mode Phase 2

Reference-matched terminal shell with sample data. The shell is the operator console at `/fast`. The dashboard stays on its existing routes.

## Added

- `FastConsole` mounts `FastShell` full viewport, outside the dashboard chrome.
- Sample cluster data behind `FastDataSource`. The table header says `sample` until Phase 5 connects the API.
- Pod columns follow the reference: NAME, READY, STATUS, RESTARTS, AGE, IP, NODE, CPU, MEM.
- Primary resource row plus a More menu for namespaces, workloads, HPA, PDB, and CRDs. Ingress stays disabled.
- Focus-aware dispatcher: `?`, `/`, `:`, Ctrl-P, `a` / Ctrl-A, and Esc. Inputs and open overlays consume other shortcuts. Esc closes the top overlay, then command mode, then blurs the field.
- `:` submits through the existing `parseColonCommand` aliases.
- Compact layout rules for a short, narrow viewport.

## Preserved

- `features/fast/FastModeWorkspace.tsx` remains on `/workloads` and still calls the API.
- Overview, Visual Mode, topology, observability, security, GitOps, and assist stay on the dashboard shell.
- Sign-in lands on `/fast`. Visual in the Fast header returns to `/visual`.

## Not done

Row movement (`j`/`k`, arrows, `g`/`G`, Enter, Ctrl-R), virtualization, live watch, and the pixel pass against the SVG. The AI overlay opens with no investigation loaded.
