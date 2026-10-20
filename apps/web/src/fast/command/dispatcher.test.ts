import { describe, expect, it } from 'vitest'
import { dispatchKey } from './dispatcher'

describe('dispatchKey', () => {
  it('opens overlays and moves the table from the shell', () => {
    const shell = { region: 'shell' as const, overlayOpen: false }
    expect(dispatchKey({ key: '?' }, shell)).toEqual({ type: 'open-help' })
    expect(dispatchKey({ key: '/' }, shell)).toEqual({ type: 'enter-filter' })
    expect(dispatchKey({ key: ':' }, shell)).toEqual({ type: 'enter-command' })
    expect(dispatchKey({ key: 'p', ctrlKey: true }, shell)).toEqual({ type: 'open-palette' })
    expect(dispatchKey({ key: 'a' }, shell)).toEqual({ type: 'open-ai' })
    expect(dispatchKey({ key: 'j' }, shell)).toEqual({ type: 'move-next' })
    expect(dispatchKey({ key: 'k' }, shell)).toEqual({ type: 'move-prev' })
    expect(dispatchKey({ key: 'ArrowDown' }, shell)).toEqual({ type: 'move-next' })
    expect(dispatchKey({ key: 'g' }, shell)).toEqual({ type: 'first' })
    expect(dispatchKey({ key: 'G' }, shell)).toEqual({ type: 'last' })
    expect(dispatchKey({ key: 'Enter' }, shell)).toEqual({ type: 'open' })
    expect(dispatchKey({ key: 'r', ctrlKey: true }, shell)).toEqual({ type: 'refresh' })
  })

  it('lets an input consume shortcuts and still closes on Esc', () => {
    const input = { region: 'input' as const, overlayOpen: false }
    expect(dispatchKey({ key: '?' }, input)).toBeNull()
    expect(dispatchKey({ key: '/' }, input)).toBeNull()
    expect(dispatchKey({ key: 'p', ctrlKey: true }, input)).toBeNull()
    expect(dispatchKey({ key: 'j' }, input)).toBeNull()
    expect(dispatchKey({ key: 'Escape' }, input)).toEqual({ type: 'close' })
    expect(dispatchKey({ key: 'r', ctrlKey: true }, input)).toEqual({ type: 'refresh' })
  })

  it('closes an open overlay and ignores further shortcuts', () => {
    const overlay = { region: 'shell' as const, overlayOpen: true }
    expect(dispatchKey({ key: 'Escape' }, overlay)).toEqual({ type: 'close' })
    expect(dispatchKey({ key: '?' }, overlay)).toBeNull()
    expect(dispatchKey({ key: 'j' }, overlay)).toBeNull()
    expect(dispatchKey({ key: 'r', ctrlKey: true }, overlay)).toBeNull()
  })
})
