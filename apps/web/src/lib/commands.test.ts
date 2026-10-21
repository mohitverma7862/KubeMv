import { describe, expect, it } from 'vitest'
import { parseColonCommand, searchCommands } from './commands'

describe('parseColonCommand', () => {
  it('parses pods with namespace and name filter', () => {
    const result = parseColonCommand(':pods -n payments /api')
    expect(result.type).toBe('resource')
    if (result.type === 'resource') {
      expect(result.kind).toBe('pods')
      expect(result.namespace).toBe('payments')
      expect(result.filter).toBe('api')
    }
  })

  it('parses deployment alias', () => {
    const result = parseColonCommand(':deploy')
    expect(result.type).toBe('resource')
    if (result.type === 'resource') {
      expect(result.kind).toBe('deployments')
    }
  })

  it('parses system commands', () => {
    expect(parseColonCommand(':q')).toEqual({ type: 'quit' })
    expect(parseColonCommand(':help')).toEqual({ type: 'help' })
    expect(parseColonCommand(':clear')).toEqual({ type: 'clear-filter' })
  })
})

describe('searchCommands', () => {
  it('finds resources by alias fragment', () => {
    const hits = searchCommands('po')
    expect(hits.some((c) => c.id === 'pods')).toBe(true)
  })
})
