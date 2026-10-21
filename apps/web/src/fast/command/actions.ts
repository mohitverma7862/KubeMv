import type { WorkspaceTab } from '../types'

/** k9s-style single-key resource actions (pod/workload context). */
export const resourceKeyToTab: Partial<Record<string, WorkspaceTab>> = {
  l: 'logs',
  d: 'describe',
  y: 'yaml',
  e: 'yaml',
  s: 'terminal',
  a: 'terminal',
  f: 'history',
}

export const workspaceTabs: WorkspaceTab[] = [
  'logs',
  'events',
  'yaml',
  'describe',
  'terminal',
  'rollout',
  'history',
]

export function nextWorkspaceTab(current: WorkspaceTab, direction: 1 | -1): WorkspaceTab {
  const index = workspaceTabs.indexOf(current)
  if (index < 0) return workspaceTabs[0]
  const next = (index + direction + workspaceTabs.length) % workspaceTabs.length
  return workspaceTabs[next]
}
