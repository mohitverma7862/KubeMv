# Fast Mode Phase 3

> Zero deviation from the Fast Mode reference: https://github.com/mohitverma7862/KubeMv/blob/main/docs/FAST-MODE-TUI-REFERENCE-SPEC.md and https://github.com/mohitverma7862/KubeMv/blob/main/docs/assets/kubemv-fast-mode-reference.svg. Match that specification and that image exactly: same logo, same colors, and same style.


The resource table is the operator surface: virtualized rows, sorting, column layout, and keyboard movement. Sample data is still the source. Live watch remains Phase 5.

## Added

- Windowed rows at 32px, with a sticky column header and horizontal scroll.
- Selected row uses the dark-blue highlight and a left cursor.
- Column sort (ascending, descending, off), drag resize, and show/hide. NAME stays visible.
- Pod columns remain NAME, READY, STATUS, RESTARTS, AGE, IP, NODE, CPU, MEM.
- Header shows the active filter, `shown/total` when a filter is set, watch age, and the live label.
- `reconcileRows` keeps object identity when a row's contents do not change.
- `j`/`k` and arrows, `g`/`G`, Enter (describe), Esc (back), and Ctrl-R (refresh age).
- Leaving a resource saves selection, filter, sort, scroll, and column layout. Esc restores that view.

## Not done

Live Kubernetes watches, command autocomplete, and the SVG pixel pass. Ctrl-R reloads the current snapshot and resets the watch age; it does not open a cluster connection.
