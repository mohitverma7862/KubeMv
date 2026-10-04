import { describe, expect, it } from 'vitest'
import { parseColonCommand } from './commands'

describe('parseColonCommand', () => {
  it('parses pods with namespace and name filter', () => {
    const result = parseColonCommand(':pods -n payments /api')
    expect(result.kind).toBe('pods')
    expect(result.namespace).toBe('payments')
    expect(result.filter).toBe('api')
  })

  it('parses deployment alias', () => {
    const result = parseColonCommand(':deploy')
    expect(result.kind).toBe('deployments')
  })
})
