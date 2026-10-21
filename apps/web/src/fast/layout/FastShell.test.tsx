import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { FastShell } from './FastShell'
import { referenceNav } from '../nav'
import type { FastHeaderModel } from '../types'

const header: FastHeaderModel = {
  clusterName: 'prod-cluster',
  userLabel: 'admin',
  contextName: 'prod',
  namespace: 'production',
  production: true,
  readOnly: true,
  connection: 'disconnected',
  clock: '12:00:00',
}

describe('FastShell', () => {
  it('renders the terminal regions without inventing resources', () => {
    render(
      <FastShell
        header={header}
        nav={referenceNav}
        columns={[{ id: 'name', header: 'NAME', value: (row) => row.name }]}
        rows={[]}
        filter=""
        workspaceTab="logs"
      />,
    )
    expect(screen.getByTestId('fast-brand')).toHaveTextContent('KubeMv')
    expect(screen.getByTestId('fast-header')).toHaveTextContent('Fast Mode (K9s-Parity)')
    expect(screen.getByTestId('fast-header')).toHaveTextContent('◈ prod-cluster')
    expect(screen.getByTestId('fast-header')).toHaveTextContent('PRODUCTION')
    expect(screen.getByTestId('fast-header')).toHaveTextContent('READ ONLY')
    expect(screen.getByTestId('fast-connection')).toHaveTextContent('DISCONNECTED')
    expect(screen.getByRole('button', { name: '[I] Ingresses' })).toBeDisabled()
    expect(screen.getByTestId('fast-table')).toHaveTextContent('No resources loaded.')
    expect(screen.getByTestId('fast-inspector')).toHaveTextContent('No resource selected.')
    expect(screen.getByTestId('fast-workspace')).toBeInTheDocument()
    expect(screen.getByTestId('fast-tree')).toHaveTextContent('No owner chain.')
    expect(screen.getByTestId('fast-status')).toHaveTextContent('Ctrl-P')
    expect(screen.queryByTestId('fast-ai')).not.toBeInTheDocument()
  })
})
