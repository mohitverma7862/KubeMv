export type ShortcutScope = 'global' | 'table' | 'logs' | 'overlay'

export type ShortcutDefinition = {
  id: string
  keys: string
  label: string
  scope: ShortcutScope
}

/** Single shortcut catalog. Help, the status bar, and later dispatchers read this list. */
export const shortcuts: ShortcutDefinition[] = [
  { id: 'move-next', keys: 'j', label: 'Next row', scope: 'table' },
  { id: 'move-prev', keys: 'k', label: 'Previous row', scope: 'table' },
  { id: 'first', keys: 'g', label: 'First row', scope: 'table' },
  { id: 'last', keys: 'G', label: 'Last row', scope: 'table' },
  { id: 'open', keys: 'Enter', label: 'Open', scope: 'table' },
  { id: 'back', keys: 'Esc', label: 'Back', scope: 'global' },
  { id: 'refresh', keys: 'Ctrl-R', label: 'Refresh', scope: 'global' },
  { id: 'filter', keys: '/', label: 'Filter', scope: 'global' },
  { id: 'command', keys: ':', label: 'Command', scope: 'global' },
  { id: 'help', keys: '?', label: 'Help', scope: 'global' },
  { id: 'palette', keys: 'Ctrl-P', label: 'Palette', scope: 'global' },
  { id: 'ai', keys: 'a', label: 'AI overlay', scope: 'global' },
]

export function shortcutById(id: string): ShortcutDefinition | undefined {
  return shortcuts.find((item) => item.id === id)
}

export function statusShortcuts(): ShortcutDefinition[] {
  return shortcuts.filter((item) => item.scope === 'global')
}
