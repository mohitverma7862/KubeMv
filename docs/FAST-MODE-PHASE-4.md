# Fast Mode Phase 4

> Zero deviation from the Fast Mode reference: https://github.com/mohitverma7862/KubeMv/blob/main/docs/FAST-MODE-TUI-REFERENCE-SPEC.md and https://github.com/mohitverma7862/KubeMv/blob/main/docs/assets/kubemv-fast-mode-reference.svg. Match that specification and that image exactly: same logo, same colors, and same style.


Command mode, the command palette, quick help, and the live filter share one registry in the terminal console.

## Added

- `:` command mode dispatches through the command registry. Aliases such as `:po`, `:dp`, and `:q` resolve to the same command as the full name.
- Typing in command mode lists matching commands. Tab fills a unique match and leaves the candidate list when more than one command matches.
- Ctrl-P opens the palette from that same registry.
- `?` replaces the persistent quick-help block with the shortcuts for the current focus. The closed help block stays the reference shortcut list.
- `/` filters live. Tokens combine with AND. `status=`, `ns=` / `namespace=`, `label=value`, `/regex/`, and fuzzy text are supported.

## Not done

- Live Kubernetes watch, real logs, and exec remain later phases.
