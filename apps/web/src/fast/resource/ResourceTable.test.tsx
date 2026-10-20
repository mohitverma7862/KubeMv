import { fireEvent, render, screen } from '@testing-library/react'
import { useState } from 'react'
import { describe, expect, it } from 'vitest'
import type { ResourceRow } from '../../services/types'
import { columnsFor } from './columns'
import { ResourceTable } from './ResourceTable'

function pod(index: number): ResourceRow {
  return {
    kind: 'Pod',
    namespace: 'production',
    name: `pod-${index}`,
    status: 'Running',
    age: `${index}m`,
    extra: { ready: '1/1', restarts: '0', ip: '10.0.0.1', node: 'node-a', cpu: '1m', mem: '1Mi' },
  }
}

function Harness() {
  const rows = Array.from({ length: 80 }, (_, index) => pod(index))
  const [hidden, setHidden] = useState<string[]>([])
  const [sort, setSort] = useState<{ columnId: string; direction: 'asc' | 'desc' } | null>(null)
  const [width, setWidth] = useState<Record<string, number>>({})
  return (
    <ResourceTable
      columns={columnsFor('pods')}
      rows={rows}
      totalCount={rows.length}
      selectedKey="production/Pod/pod-0"
      filter=""
      liveLabel="● Live · sample"
      watchAge="◷ 0s"
      viewportHeight={64}
      sort={sort}
      onSort={(columnId) =>
        setSort((current) => (current?.columnId === columnId && current.direction === 'asc' ? { columnId, direction: 'desc' } : { columnId, direction: 'asc' }))
      }
      hiddenColumnIds={hidden}
      onToggleColumn={(columnId) =>
        setHidden((current) => (current.includes(columnId) ? current.filter((id) => id !== columnId) : [...current, columnId]))
      }
      columnWidths={width}
      onColumnWidth={(columnId, next) => setWidth((current) => ({ ...current, [columnId]: next }))}
    />
  )
}

describe('ResourceTable', () => {
  it('virtualizes rows and configures columns', () => {
    render(<Harness />)
    expect(screen.getByTestId('fast-columns-head')).toHaveTextContent('NAME')
    expect(screen.getByTestId('fast-columns-head')).toHaveTextContent('CPU')
    expect(screen.getByTestId('fast-count')).toHaveTextContent('80')
    expect(screen.getAllByTestId('fast-row').length).toBeLessThan(20)
    expect(screen.getByText('pod-0')).toBeInTheDocument()
    expect(screen.queryByText('pod-40')).not.toBeInTheDocument()

    fireEvent.click(screen.getByRole('button', { name: 'cols' }))
    fireEvent.click(screen.getByRole('checkbox', { name: 'IP' }))
    expect(screen.getByTestId('fast-columns-head')).not.toHaveTextContent('IP')

    fireEvent.click(screen.getByRole('columnheader', { name: 'NAME' }))
    expect(screen.getByRole('columnheader', { name: 'NAME ▲' })).toHaveAttribute('data-sort', 'asc')

    const handle = document.querySelector('.km-fast-resize')
    expect(handle).toBeTruthy()
    fireEvent.pointerDown(handle as Element, { clientX: 100 })
    fireEvent.pointerMove(window, { clientX: 140 })
    fireEvent.pointerUp(window)
    expect(screen.getByRole('columnheader', { name: 'NAME ▲' })).toHaveStyle({ width: '320px' })
  })
})
