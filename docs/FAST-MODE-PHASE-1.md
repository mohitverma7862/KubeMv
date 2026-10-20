# Fast Mode Phase 1

UI architecture only. The existing Fast page (`features/fast/FastModeWorkspace.tsx`) is unchanged and still calls the API.

## Added

`apps/web/src/fast/` holds the terminal shell contracts:

- `layout/` header, resource nav, splits, shell
- `resource/` table, inspector, owner tree
- `terminal/` logs, events, YAML, describe, exec placeholder, workspace tabs
- `command/` one shortcut registry, command bar, palette, help, autocomplete over the existing `lib/commands.ts` aliases
- `filter/` parser for text, `status=`, labels, and regex
- `ai/` overlay that renders nothing until `open` is true
- `theme/kubemv-terminal.css` with the reference palette

`FastShell` is not mounted on a route. Phase 2 replaces the current Fast screen with this shell and the reference composition.

## Reused

- `ResourceRow`, `ResourceDetail`, and `ResourceKindPath` from `apps/web/src/services/types.ts`
- `parseColonCommand` and `resourceCommands` from `apps/web/src/lib/commands.ts`

## Not done

Keyboard dispatch, virtualization, mock cluster rows, live watches, and the reference pixel pass. Ingress stays disabled because the API has no ingress list.
