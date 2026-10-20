import { describe, expect, it } from 'vitest'
import type { ResourceRow } from '../../services/types'
import { reconcileRows, sortRows, visibleWindow } from './tableModel'
import type { TableColumn } from '../types'

const column: TableColumn = { id: 'name', header: 'NAME', value: (row) => row.name }
const age: TableColumn = { id: 'age', header: 'AGE', value: (row) => row.age }

function pod(name: string, ageValue: string): ResourceRow {
  return { kind: 'Pod', namespace: 'production', name, status: 'Running', age: ageValue }
}

describe('table model', () => {
  it('windows a tall list', () => {
    expect(visibleWindow(500, 0, 96, 32, 1)).toEqual({ start: 0, end: 4 })
    expect(visibleWindow(500, 320, 96, 32, 1)).toEqual({ start: 9, end: 14 })
  })

  it('sorts age chronologically and reconciles unchanged rows', () => {
    const rows = [pod('b', '4h'), pod('a', '2m')]
    expect(sortRows(rows, [column, age], { columnId: 'age', direction: 'asc' }).map((row) => row.name)).toEqual(['a', 'b'])

    const updated = [pod('b', '4h'), { ...pod('a', '2m'), status: 'Pending' }]
    const reconciled = reconcileRows(rows, updated)
    expect(reconciled[0]).toBe(rows[0])
    expect(reconciled[1]).not.toBe(rows[1])
    expect(reconciled[1].status).toBe('Pending')
  })
})
