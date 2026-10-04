# Design system

The operator UI is a dense control surface: a brass accent on a graphite field, IBM Plex Sans for text, and IBM Plex Mono for status, roles, and shortcuts. Dark is the default. Light is a token swap on `document.documentElement.dataset.theme`.

## Chrome

```text
Top bar     product, phase, command trigger, operator
Navigation  Overview, Clusters, Plugins, Settings
Main        the active page
Status      API reachability and connection honesty
```

Tokens live in `desktop/src/styles/tokens.css`. Components use the `km-` classes in `desktop/src/styles/global.css`: panels, tables, fields, badges, buttons, and dialogs.

## Keyboard

| Input | Action |
| --- | --- |
| `Ctrl/Cmd K` | Command palette |
| `?` | Shortcut help |
| `Esc` | Close palette or help |
| `g` then `o`, `c`, `p`, or `s` | Navigate the four foundation pages |

Chords are ignored while typing in a field. Navigation only targets screens this phase implements.

## Responsive

Below 800px the navigation becomes a horizontal bar and the login narrative column hides. Editing workflows stay limited to the registry form. The layout is read-friendly on a narrow viewport; it is not a separate mobile product.
