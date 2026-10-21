export type ShortcutScope = 'global' | 'table' | 'workspace' | 'overlay'

export type ShortcutDefinition = {
  id: string
  keys: string
  label: string
  scope: ShortcutScope
  section?: 'navigation' | 'general' | 'resource' | 'workspace'
}

/** Single shortcut catalog (k9s-aligned). Help, status bar, and dispatch read this list. */
export const shortcuts: ShortcutDefinition[] = [
  { id: 'move-next', keys: 'j', label: 'Down', scope: 'table', section: 'navigation' },
  { id: 'move-prev', keys: 'k', label: 'Up', scope: 'table', section: 'navigation' },
  { id: 'move-down', keys: '↓', label: 'Down', scope: 'table', section: 'navigation' },
  { id: 'move-up', keys: '↑', label: 'Up', scope: 'table', section: 'navigation' },
  { id: 'first', keys: 'g', label: 'Goto top', scope: 'table', section: 'navigation' },
  { id: 'last', keys: 'G', label: 'Goto bottom', scope: 'table', section: 'navigation' },
  { id: 'page-up', keys: 'Ctrl-B', label: 'Page up', scope: 'table', section: 'navigation' },
  { id: 'page-down', keys: 'Ctrl-F', label: 'Page down', scope: 'table', section: 'navigation' },
  { id: 'page-up-alt', keys: 'PgUp', label: 'Page up', scope: 'table', section: 'navigation' },
  { id: 'page-down-alt', keys: 'PgDn', label: 'Page down', scope: 'table', section: 'navigation' },
  { id: 'history-back', keys: '[', label: 'History back', scope: 'global', section: 'navigation' },
  { id: 'history-fwd', keys: ']', label: 'History forward', scope: 'global', section: 'navigation' },
  { id: 'open', keys: 'Enter', label: 'Drill down / describe', scope: 'table', section: 'navigation' },
  { id: 'back', keys: 'Esc', label: 'Back / clear', scope: 'global', section: 'general' },
  { id: 'back-q', keys: 'q', label: 'Back', scope: 'global', section: 'general' },
  { id: 'refresh', keys: 'Ctrl-R', label: 'Reload', scope: 'global', section: 'general' },
  { id: 'filter', keys: '/', label: 'Filter', scope: 'global', section: 'general' },
  { id: 'command', keys: ':', label: 'Command mode', scope: 'global', section: 'general' },
  { id: 'help', keys: '?', label: 'Help', scope: 'global', section: 'general' },
  { id: 'palette', keys: 'Ctrl-P', label: 'Command palette', scope: 'global', section: 'general' },
  { id: 'aliases', keys: 'Ctrl-A', label: 'Aliases / palette', scope: 'global', section: 'general' },
  { id: 'ai', keys: 'a', label: 'AI overlay', scope: 'global', section: 'general' },
  { id: 'repeat-cmd', keys: '-', label: 'Last command', scope: 'global', section: 'general' },
  { id: 'tab-next', keys: 'Tab', label: 'Next workspace tab', scope: 'workspace', section: 'workspace' },
  { id: 'tab-prev', keys: 'Shift-Tab', label: 'Previous workspace tab', scope: 'workspace', section: 'workspace' },
  { id: 'logs', keys: 'l', label: 'Logs', scope: 'table', section: 'resource' },
  { id: 'shell', keys: 's', label: 'Shell', scope: 'table', section: 'resource' },
  { id: 'describe', keys: 'd', label: 'Describe', scope: 'table', section: 'resource' },
  { id: 'yaml', keys: 'y', label: 'YAML', scope: 'table', section: 'resource' },
  { id: 'edit', keys: 'e', label: 'Edit (YAML)', scope: 'table', section: 'resource' },
  { id: 'port-forward', keys: 'f', label: 'Port forward', scope: 'table', section: 'resource' },
  { id: 'attach', keys: 'A', label: 'Attach', scope: 'table', section: 'resource' },
  { id: 'container', keys: 'c', label: 'Containers', scope: 'table', section: 'resource' },
  { id: 'delete', keys: 'Shift-D', label: 'Delete', scope: 'table', section: 'resource' },
  { id: 'digit-view', keys: '1-9', label: 'Switch resource view', scope: 'global', section: 'navigation' },
  { id: 'nav-letter', keys: 'K,D,S,…', label: 'Resource mnemonic', scope: 'global', section: 'navigation' },
]

export function shortcutById(id: string): ShortcutDefinition | undefined {
  return shortcuts.find((item) => item.id === id)
}

export function statusShortcuts(): ShortcutDefinition[] {
  return shortcuts.filter((item) => item.scope === 'global' && item.section === 'general')
}

export function helpSections(): Record<string, ShortcutDefinition[]> {
  const sections: Record<string, ShortcutDefinition[]> = {
    Navigation: [],
    General: [],
    Resource: [],
    Workspace: [],
  }
  for (const item of shortcuts) {
    const key =
      item.section === 'navigation'
        ? 'Navigation'
        : item.section === 'resource'
          ? 'Resource'
          : item.section === 'workspace'
            ? 'Workspace'
            : 'General'
    sections[key].push(item)
  }
  return sections
}
