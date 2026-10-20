import { describe, expect, it } from 'vitest'
import { dispatchKey } from './dispatcher'

describe('dispatchKey', () => {
  it('opens overlays from the shell and leaves table keys alone', () => {
    const shell = { region: 'shell' as const, overlayOpen: false }
    expect(dispatchKey({ key: '?' }, shell)).toEqual({ type: 'open-help' })
    expect(dispatchKey({ key: '/' }, shell)).toEqual({ type: 'enter-filter' })
    expect(dispatchKey({ key: ':' }, shell)).toEqual({ type: 'enter-command' })
    expect(dispatchKey({ key: 'p', ctrlKey: true }, shell)).toEqual({ type: 'open-palette' })
    expect(dispatchKey({ key: 'a' }, shell)).toEqual({ type: 'open-ai' })
    expect(dispatchKey({ key: 'j' }, shell)).toBeNull()
    expect(dispatchKey({ key: 'k' }, shell)).toBeNull()
  })

  it('lets an input consume shortcuts and still closes on Esc', () => {
    const input = { region: 'input' as const, overlayOpen: false }
    expect(dispatchKey({ key: '?' }, input)).toBeNull()
    expect(dispatchKey({ key: '/' }, input)).toBeNull()
    expect(dispatchKey({ key: 'p', ctrlKey: true }, input)).toBeNull()
    expect(dispatchKey({ key: 'Escape' }, input)).toEqual({ type: 'close' })
  })

  it('closes an open overlay and ignores further shortcuts', () => {
    const overlay = { region: 'shell' as const, overlayOpen: true }
    expect(dispatchKey({ key: 'Escape' }, overlay)).toEqual({ type: 'close' })
    expect(dispatchKey({ key: '?' }, overlay)).toBeNull()
    expect(dispatchKey({ key: 'a' }, overlay)).toBeNull()
  })
})
