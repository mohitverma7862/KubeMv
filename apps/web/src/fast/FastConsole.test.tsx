import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { sampleCluster } from './data/sampleCluster'
import { FastConsole } from './FastConsole'

describe('FastConsole', () => {
  it('renders the sample operator console', () => {
    render(<FastConsole clock="12:00:00" />)
    expect(screen.getByTestId('fast-header')).toHaveTextContent('prod-cluster')
    expect(screen.getByTestId('fast-header')).toHaveTextContent('PRODUCTION')
    expect(screen.getByTestId('fast-header')).toHaveTextContent('READ ONLY')
    expect(screen.getByTestId('fast-header')).toHaveTextContent('12:00:00')
    expect(screen.getByRole('button', { name: '[K] Pods' })).toHaveAttribute('aria-current', 'page')
    expect(screen.getByRole('button', { name: '[I] Ingresses' })).toBeDisabled()
    expect(screen.getByText('payments-worker-6b8d4c7f9-q1w2e')).toBeInTheDocument()
    expect(screen.getByTestId('fast-table-header')).toHaveTextContent('sample')
    expect(screen.getByTestId('fast-inspector')).toHaveTextContent('api-7c9f8b6d4-k2m9x')
    expect(screen.getByTestId('fast-tree')).toHaveTextContent('ReplicaSet')
    expect(screen.getByTestId('fast-logs')).toHaveTextContent('listening')
    expect(screen.getByTestId('fast-status')).toHaveTextContent('ctx prod')
    expect(screen.getByTestId('fast-status')).toHaveTextContent('ns production')
    expect(screen.queryByTestId('fast-ai')).not.toBeInTheDocument()
    expect(JSON.stringify(sampleCluster)).not.toMatch(/BEGIN |password|token/i)
  })

  it('routes shortcuts by focus', () => {
    render(<FastConsole clock="12:00:00" />)

    fireEvent.keyDown(window, { key: '?' })
    expect(screen.getByTestId('fast-help')).toBeInTheDocument()
    fireEvent.keyDown(window, { key: 'Escape' })
    expect(screen.queryByTestId('fast-help')).not.toBeInTheDocument()

    fireEvent.keyDown(window, { key: '/' })
    const filter = screen.getByRole('textbox', { name: 'Filter' })
    expect(filter).toHaveFocus()
    fireEvent.keyDown(filter, { key: '?' })
    expect(screen.queryByTestId('fast-help')).not.toBeInTheDocument()
    fireEvent.change(filter, { target: { value: 'status=CrashLoopBackOff' } })
    expect(screen.getAllByText('payments-worker-6b8d4c7f9-q1w2e').length).toBeGreaterThan(0)
    expect(screen.queryByText('web-6d5c4b3a2-h7t1a')).not.toBeInTheDocument()

    fireEvent.keyDown(filter, { key: 'Escape' })
    fireEvent.keyDown(window, { key: ':' })
    const command = screen.getByRole('textbox', { name: 'Command' })
    fireEvent.change(command, { target: { value: 'deploy' } })
    fireEvent.keyDown(command, { key: 'Enter' })
    expect(screen.getByRole('button', { name: '[D] Deployments' })).toHaveAttribute('aria-current', 'page')

    fireEvent.keyDown(window, { key: 'p', ctrlKey: true })
    expect(screen.getByTestId('fast-palette')).toBeInTheDocument()
    fireEvent.keyDown(window, { key: 'a' })
    expect(screen.queryByTestId('fast-ai')).not.toBeInTheDocument()
    fireEvent.keyDown(window, { key: 'Escape' })
    fireEvent.keyDown(window, { key: 'a' })
    expect(screen.getByTestId('fast-ai')).toHaveTextContent('never run')
  })

  it('moves the selection and restores the previous view', () => {
    render(<FastConsole clock="12:00:00" />)
    expect(screen.getByTestId('fast-count')).toHaveTextContent('4')
    expect(screen.getByTestId('fast-columns-head')).toHaveTextContent('READY')
    expect(screen.getByTestId('fast-columns-head')).toHaveTextContent('MEM')

    fireEvent.keyDown(window, { key: 'j' })
    expect(screen.getByTestId('fast-inspector')).toHaveTextContent('payments-worker-6b8d4c7f9-q1w2e')
    fireEvent.keyDown(window, { key: 'k' })
    expect(screen.getByTestId('fast-inspector')).toHaveTextContent('api-7c9f8b6d4-k2m9x')
    fireEvent.keyDown(window, { key: 'G' })
    expect(screen.getByTestId('fast-inspector')).toHaveTextContent('web-6d5c4b3a2-h7t1a')
    fireEvent.keyDown(window, { key: 'g' })
    expect(screen.getByTestId('fast-inspector')).toHaveTextContent('api-7c9f8b6d4-k2m9x')

    fireEvent.change(screen.getByRole('textbox', { name: 'Filter' }), { target: { value: 'web' } })
    fireEvent.click(screen.getByRole('button', { name: '[D] Deployments' }))
    expect(screen.getByRole('button', { name: '[D] Deployments' })).toHaveAttribute('aria-current', 'page')
    fireEvent.keyDown(window, { key: 'Escape' })
    expect(screen.getByRole('button', { name: '[K] Pods' })).toHaveAttribute('aria-current', 'page')
    expect(screen.getByRole('textbox', { name: 'Filter' })).toHaveValue('web')

    fireEvent.keyDown(window, { key: 'Enter' })
    expect(screen.getByTestId('fast-describe')).toBeInTheDocument()
    fireEvent.keyDown(window, { key: 'Escape' })
    expect(screen.queryByTestId('fast-describe')).not.toBeInTheDocument()
    expect(screen.getByTestId('fast-logs')).toBeInTheDocument()

    fireEvent.keyDown(window, { key: 'r', ctrlKey: true })
    expect(screen.getByTestId('fast-table-header')).toHaveTextContent('◷ 0s')
  })

  it('switches resource views from nav mnemonics', () => {
    render(<FastConsole clock="12:00:00" />)
    const nav = [
      { key: 'E', label: '[E] Events', kind: 'events' },
      { key: 'N', label: '[N] Nodes', kind: 'nodes' },
      { key: 'M', label: '[M] ConfigMaps', kind: 'configmaps' },
      { key: 'D', label: '[D] Deployments', kind: 'deployments' },
      { key: 'S', label: '[S] Services', kind: 'services' },
      { key: 'J', label: '[J] Jobs', kind: 'jobs' },
      { key: 'C', label: '[C] CronJobs', kind: 'cronjobs' },
      { key: 'R', label: '[R] Secrets', kind: 'secrets' },
      { key: 'K', label: '[K] Pods', kind: 'pods' },
    ] as const
    for (const item of nav) {
      fireEvent.keyDown(window, { key: item.key })
      expect(screen.getByRole('button', { name: item.label })).toHaveAttribute('aria-current', 'page')
    }
    fireEvent.keyDown(window, { key: 'n' })
    expect(screen.getByRole('button', { name: '[N] Nodes' })).toHaveAttribute('aria-current', 'page')
    fireEvent.keyDown(window, { key: '1' })
    expect(screen.getByRole('button', { name: '[K] Pods' })).toHaveAttribute('aria-current', 'page')
  })
})