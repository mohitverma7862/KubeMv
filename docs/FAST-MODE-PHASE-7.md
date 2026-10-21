# Fast Mode Phase 7

> Zero deviation from the Fast Mode reference: https://github.com/mohitverma7862/KubeMv/blob/main/docs/FAST-MODE-TUI-REFERENCE-SPEC.md and https://github.com/mohitverma7862/KubeMv/blob/main/docs/assets/kubemv-fast-mode-reference.svg. Match that specification and that image exactly: same logo, same colors, and same style.


The lower workspace tabs load logs, events, YAML, describe, a shell, rollout status, and visit history.

## Added

- Logs keep a bounded buffer. `f` follows, `p` pauses, `P` switches to the previous container log, `t` shows timestamps, `w` wraps, `c` cycles containers, `'` searches, and `n`/`N` move between matches. Ctrl-C stops follow. `C` clears. `O` writes the visible log to a private file.
- Events, YAML, and describe load from the cluster for the selected resource. Secret YAML drops data and string data. Describe hides secret values.
- `s` opens a shell through the same SPDY exec path as the websocket handler. The stream accepts terminal resize. Read-only mode rejects exec.
- Rollout shows deployment replica counts and revisions.
- History lists the resources visited in this console.

## Not done

- Edit, delete, attach, and port-forward remain Phase 8. The console still starts read-only, so a live shell waits until that mode allows it.
