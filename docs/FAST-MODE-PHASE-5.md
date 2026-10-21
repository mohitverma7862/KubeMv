# Fast Mode Phase 5

> Zero deviation from the Fast Mode reference: https://github.com/mohitverma7862/KubeMv/blob/main/docs/FAST-MODE-TUI-REFERENCE-SPEC.md and https://github.com/mohitverma7862/KubeMv/blob/main/docs/assets/kubemv-fast-mode-reference.svg. Match that specification and that image exactly: same logo, same colors, and same style.


The terminal console watches the active resource kind from the current kubeconfig context.

## Added

- A watch session lists the kind, then watches from that resourceVersion. Bookmarks advance the resume version. An expired watch lists again before opening a new watch.
- Failed lists retry with backoff. The header shows `◌ RECONNECTING <attempt> in <delay>`. After the attempt budget it shows `✕ DISCONNECTED` and the error. A healthy watch shows `● LIVE`.
- Ctrl-R starts the watch again. Switching resource or namespace does the same. `:` namespace changes are included.
- Rows are stored by namespace and name. One event updates that row. The table repaints the changed row and the cursor, and does not clear the table.
- The resource title shows the watch age and the current resourceVersion.
- Without a kubeconfig context, the console stays on the sample cluster.

## Not done

- Logs, exec, describe, and the relationship tree remain later phases.
