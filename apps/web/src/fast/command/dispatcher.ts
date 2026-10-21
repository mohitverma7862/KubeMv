export type FocusRegion = 'shell' | 'input' | 'overlay' | 'workspace'

export type KeyInput = {
  key: string
  ctrlKey?: boolean
  metaKey?: boolean
  altKey?: boolean
  shiftKey?: boolean
}

export type DispatchAction =
  | { type: 'open-help' }
  | { type: 'open-palette' }
  | { type: 'open-ai' }
  | { type: 'enter-filter' }
  | { type: 'enter-command' }
  | { type: 'close' }
  | { type: 'back' }
  | { type: 'quit' }
  | { type: 'move-next' }
  | { type: 'move-prev' }
  | { type: 'page-up' }
  | { type: 'page-down' }
  | { type: 'first' }
  | { type: 'last' }
  | { type: 'open' }
  | { type: 'refresh' }
  | { type: 'history-back' }
  | { type: 'history-forward' }
  | { type: 'repeat-command' }
  | { type: 'workspace-tab-next' }
  | { type: 'workspace-tab-prev' }
  | { type: 'resource-action'; action: 'logs' | 'shell' | 'describe' | 'yaml' | 'edit' | 'port-forward' | 'attach' | 'container' | 'delete' }
  | { type: 'nav-digit'; digit: string }
  | { type: 'nav-letter'; letter: string }
  | { type: 'command-clear' }

export type DispatchContext = {
  region: FocusRegion
  overlayOpen: boolean
  commandActive: boolean
}

function isMod(event: KeyInput): boolean {
  return Boolean(event.ctrlKey || event.metaKey)
}

/** Focus-aware shortcut routing modeled on k9s help (general + navigation + resource keys). */
export function dispatchKey(event: KeyInput, ctx: DispatchContext): DispatchAction | null {
  if (event.altKey) return null
  const key = event.key

  if (key === 'Escape') return { type: 'close' }

  if (ctx.overlayOpen) {
    if (key === 'q' || key === 'Q') return { type: 'close' }
    return null
  }

  if (ctx.region === 'input') {
    if (isMod(event) && (key === 'r' || key === 'R')) return { type: 'refresh' }
    if (isMod(event) && (key === 'u' || key === 'U') && ctx.commandActive) return { type: 'command-clear' }
    if (key === 'Escape') return { type: 'close' }
    return null
  }

  if (isMod(event)) {
    if (key === 'r' || key === 'R') return { type: 'refresh' }
    if (key === 'p' || key === 'P') return { type: 'open-palette' }
    if (key === 'a' || key === 'A') return { type: 'open-palette' }
    if (key === 'b' || key === 'B') return { type: 'page-up' }
    if (key === 'f' || key === 'F') return { type: 'page-down' }
    return null
  }

  if (key === 'Tab' && event.shiftKey) return { type: 'workspace-tab-prev' }
  if (key === 'Tab') return { type: 'workspace-tab-next' }

  if (key === '?') return { type: 'open-help' }
  if (key === '/') return { type: 'enter-filter' }
  if (key === ':') return { type: 'enter-command' }
  if (key === '-') return { type: 'repeat-command' }
  if (key === '[') return { type: 'history-back' }
  if (key === ']') return { type: 'history-forward' }

  if (key === 'q') return { type: 'back' }

  if (/^[1-9]$/.test(key)) return { type: 'nav-digit', digit: key }

  if (key === 'l') return { type: 'resource-action', action: 'logs' }
  if (key === 's') return { type: 'resource-action', action: 'shell' }
  if (key === 'd') return { type: 'resource-action', action: 'describe' }
  if (key === 'y') return { type: 'resource-action', action: 'yaml' }
  if (key === 'e') return { type: 'resource-action', action: 'edit' }
  if (key === 'f') return { type: 'resource-action', action: 'port-forward' }
  if (key === 'c') return { type: 'resource-action', action: 'container' }
  if (key === 'a') return { type: 'open-ai' }

  if (key === 'A' && event.shiftKey) return { type: 'resource-action', action: 'attach' }
  if (key === 'D' && event.shiftKey) return { type: 'resource-action', action: 'delete' }

  if (key === 'j' || key === 'ArrowDown') return { type: 'move-next' }
  if (key === 'k' || key === 'ArrowUp') return { type: 'move-prev' }
  if (key === 'PageUp') return { type: 'page-up' }
  if (key === 'PageDown') return { type: 'page-down' }
  if (key === 'g') return { type: 'first' }
  if (key === 'G') return { type: 'last' }
  if (key === 'Enter') return { type: 'open' }

  if (key.length === 1 && key >= 'A' && key <= 'Z') {
    return { type: 'nav-letter', letter: key }
  }

  return null
}

export function focusRegion(target: EventTarget | null): FocusRegion {
  if (!(target instanceof Element)) return 'shell'
  if (target.closest('[data-fast-overlay]')) return 'overlay'
  if (target.closest('[data-fast-workspace]')) return 'workspace'
  if (target.closest('input, textarea, [contenteditable="true"]')) return 'input'
  return 'shell'
}
