# Fast Mode Phase 2

> Zero deviation from the Fast Mode reference: https://github.com/mohitverma7862/KubeMv/blob/main/docs/FAST-MODE-TUI-REFERENCE-SPEC.md and https://github.com/mohitverma7862/KubeMv/blob/main/docs/assets/kubemv-fast-mode-reference.svg. Match that specification and that image exactly: same logo, same colors, and same style.


Reference-matched terminal shell with sample data. The shell is the only operator console.

## Added

- `FastConsole` mounts `FastShell` as the full operator console.
- Sample cluster data behind `FastDataSource`. The table header says `sample` until Phase 5 connects the API.
- Pod columns follow the reference: NAME, READY, STATUS, RESTARTS, AGE, IP, NODE, CPU, MEM.
- Primary resource row plus a More menu for namespaces, workloads, HPA, PDB, and CRDs. Ingress stays disabled.
- Focus-aware dispatcher: `?`, `/`, `:`, Ctrl-P, `a` / Ctrl-A, and Esc. Inputs and open overlays consume other shortcuts. Esc closes the top overlay, then command mode, then blurs the field.
- `:` submits through the existing `parseColonCommand` aliases.
- Compact layout rules for a short, narrow viewport.

## Preserved

- `features/fast/FastModeWorkspace.tsx` remains on `/workloads` and still calls the API.
- Sign-in lands on the Fast Mode console.

## Not done

Row movement (`j`/`k`, arrows, `g`/`G`, Enter, Ctrl-R), virtualization, live watch, and the pixel pass against the SVG. The AI overlay opens with no investigation loaded.
