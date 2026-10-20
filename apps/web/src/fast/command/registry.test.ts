import { describe, expect, it } from 'vitest'
import { shortcuts } from './registry'

describe('shortcut registry', () => {
  it('keeps one definition per id and per key', () => {
    const ids = shortcuts.map((item) => item.id)
    const keys = shortcuts.map((item) => item.keys)
    expect(new Set(ids).size).toBe(ids.length)
    expect(new Set(keys).size).toBe(keys.length)
    expect(shortcuts.map((item) => item.keys)).toEqual(
      expect.arrayContaining(['j', 'k', '/', ':', '?', 'Ctrl-P', 'Esc']),
    )
  })
})
