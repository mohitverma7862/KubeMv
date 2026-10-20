import type { ResourceNavItem } from './types'

/** Primary resource row from the Fast Mode reference. Unavailable items stay visible and disabled. */
export const referenceNav: ResourceNavItem[] = [
  { id: 'pods', label: 'Pods', shortcut: 'K', kind: 'pods', available: true },
  { id: 'deployments', label: 'Deployments', shortcut: 'D', kind: 'deployments', available: true },
  { id: 'services', label: 'Services', shortcut: 'S', kind: 'services', available: true },
  { id: 'nodes', label: 'Nodes', shortcut: 'N', kind: 'nodes', available: true },
  { id: 'ingresses', label: 'Ingresses', shortcut: 'I', available: false },
  { id: 'events', label: 'Events', shortcut: 'E', kind: 'events', available: true },
  { id: 'jobs', label: 'Jobs', shortcut: 'J', kind: 'jobs', available: true },
  { id: 'cronjobs', label: 'CronJobs', shortcut: 'C', kind: 'cronjobs', available: true },
  { id: 'configmaps', label: 'ConfigMaps', shortcut: 'M', kind: 'configmaps', available: true },
  { id: 'secrets', label: 'Secrets', shortcut: 'R', kind: 'secrets', available: true },
]

/** Kinds reachable from More or `:`. They stay off the primary row. */
export const moreNav: ResourceNavItem[] = [
  { id: 'namespaces', label: 'Namespaces', shortcut: ':', kind: 'namespaces', available: true },
  { id: 'statefulsets', label: 'StatefulSets', shortcut: ':', kind: 'statefulsets', available: true },
  { id: 'daemonsets', label: 'DaemonSets', shortcut: ':', kind: 'daemonsets', available: true },
  { id: 'hpa', label: 'HPA', shortcut: ':', kind: 'hpa', available: true },
  { id: 'pdb', label: 'PDB', shortcut: ':', kind: 'pdb', available: true },
  { id: 'crds', label: 'CRDs', shortcut: ':', kind: 'crds', available: true },
]
