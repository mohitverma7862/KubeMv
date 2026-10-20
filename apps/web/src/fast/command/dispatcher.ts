export type FocusRegion = 'shell' | 'input' | 'overlay'

export type KeyInput = {
  key: string
  ctrlKey?: boolean
  metaKey?: boolean
  altKey?: boolean
}

export type DispatchAction =
  | { type: 'open-help' }
  | { type: 'open-palette' }
  | { type: 'open-ai' }
  | { type: 'enter-filter' }
  | { type: 'enter-command' }
  | { type: 'close' }
  | { type: 'move-next' }
  | { type: 'move-prev' }
  | { type: 'first' }
  | { type: 'last' }
  | { type: 'open' }
  | { type: 'refresh' }

export type DispatchContext = {
  region: FocusRegion
  overlayOpen: boolean
}

/** Focus-aware shortcut routing. Inputs and overlays keep Esc as the way back. */
export function dispatchKey(event: KeyInput, ctx: DispatchContext): DispatchAction | null {
  if (event.altKey) return null
  const key = event.key
  if (key === 'Escape') return { type: 'close' }
  if ((event.ctrlKey || event.metaKey) && (key === 'r' || key === 'R') && !ctx.overlayOpen) {
    return { type: 'refresh' }
  }

  if (ctx.overlayOpen || ctx.region !== 'shell') return null

  if (event.ctrlKey || event.metaKey) {
    if (key === 'p' || key === 'P') return { type: 'open-palette' }
    if (key === 'a' || key === 'A') return { type: 'open-ai' }
    return null
  }

  if (key === 'j' || key === 'ArrowDown') return { type: 'move-next' }
  if (key === 'k' || key === 'ArrowUp') return { type: 'move-prev' }
  if (key === 'g') return { type: 'first' }
  if (key === 'G') return { type: 'last' }
  if (key === 'Enter') return { type: 'open' }

  if (key === '?') return { type: 'open-help' }
  if (key === '/') return { type: 'enter-filter' }
  if (key === ':') return { type: 'enter-command' }
  if (key === 'a') return { type: 'open-ai' }
  return null
}

export function focusRegion(target: EventTarget | null): FocusRegion {
  if (!(target instanceof Element)) return 'shell'
  if (target.closest('[data-fast-overlay]')) return 'overlay'
  if (target.closest('input, textarea, [contenteditable="true"]')) return 'input'
  return 'shell'
}
