import { describe, expect, it } from 'vitest'
import type { ResourceRow } from '../../services/types'
import { applyFilter, parseFilter } from './FilterEngine'

const rows: ResourceRow[] = [
  { kind: 'Pod', namespace: 'payments', name: 'api', status: 'Running', age: '1h', labels: { app: 'api' } },
  { kind: 'Pod', namespace: 'payments', name: 'worker', status: 'Pending', age: '2m', labels: { app: 'worker' } },
]

describe('parseFilter', () => {
  it('reads status, labels, and text', () => {
    const parsed = parseFilter('status=Running app=api api')
    expect(parsed.status).toBe('Running')
    expect(parsed.labels.app).toBe('api')
    expect(parsed.text).toBe('api')
  })

  it('rejects an invalid regex without throwing', () => {
    expect(parseFilter('/(/').invalidRegex).toBe('(')
  })
})

describe('applyFilter', () => {
  it('keeps rows that match every term', () => {
    expect(applyFilter(rows, 'status=Running app=api').map((row) => row.name)).toEqual(['api'])
  })
})
