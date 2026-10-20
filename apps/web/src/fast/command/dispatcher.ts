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

export type DispatchContext = {
  region: FocusRegion
  overlayOpen: boolean
}

/** Focus-aware shortcut routing. Table movement stays out of this phase. */
export function dispatchKey(event: KeyInput, ctx: DispatchContext): DispatchAction | null {
  if (event.altKey) return null
  const key = event.key
  if (key === 'Escape') return { type: 'close' }

  if (ctx.overlayOpen || ctx.region !== 'shell') return null

  if (event.ctrlKey || event.metaKey) {
    if (key === 'p' || key === 'P') return { type: 'open-palette' }
    if (key === 'a' || key === 'A') return { type: 'open-ai' }
    return null
  }

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
