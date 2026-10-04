import type { ResourceKindPath } from '../services/types'

export interface CommandDefinition {
  id: string
  title: string
  aliases: string[]
  kind?: ResourceKindPath
  shortcut?: string
  risk?: 'READ' | 'LOW' | 'MEDIUM' | 'HIGH'
}

export const resourceCommands: CommandDefinition[] = [
  { id: 'pods', title: 'Pods', aliases: [':pods', ':po', 'pods'], kind: 'pods', shortcut: 'Ctrl+P' },
  { id: 'deployments', title: 'Deployments', aliases: [':deploy', ':deployments', ':dp'], kind: 'deployments' },
  { id: 'services', title: 'Services', aliases: [':svc', ':services'], kind: 'services' },
  { id: 'nodes', title: 'Nodes', aliases: [':nodes', ':no'], kind: 'nodes' },
  { id: 'events', title: 'Events', aliases: [':events', ':ev'], kind: 'events' },
  { id: 'namespaces', title: 'Namespaces', aliases: [':ns', ':namespaces'], kind: 'namespaces' },
  { id: 'configmaps', title: 'ConfigMaps', aliases: [':configmap', ':cm'], kind: 'configmaps' },
  { id: 'secrets', title: 'Secrets (metadata)', aliases: [':secret', ':secrets'], kind: 'secrets' },
  { id: 'crds', title: 'CRDs', aliases: [':crd', ':crds'], kind: 'crds' },
  { id: 'statefulsets', title: 'StatefulSets', aliases: [':sts', ':statefulsets'], kind: 'statefulsets' },
  { id: 'daemonsets', title: 'DaemonSets', aliases: [':ds', ':daemonsets'], kind: 'daemonsets' },
  { id: 'jobs', title: 'Jobs', aliases: [':job', ':jobs'], kind: 'jobs' },
  { id: 'cronjobs', title: 'CronJobs', aliases: [':cronjob', ':cronjobs'], kind: 'cronjobs' },
  { id: 'hpa', title: 'HPA', aliases: [':hpa'], kind: 'hpa' },
  { id: 'pdb', title: 'PDB', aliases: [':pdb'], kind: 'pdb' },
]

export function parseColonCommand(input: string): { kind?: ResourceKindPath; filter?: string; namespace?: string } {
  const trimmed = input.trim()
  if (!trimmed.startsWith(':')) return {}
  const parts = trimmed.slice(1).trim().split(/\s+/)
  const cmdPart = parts[0]?.toLowerCase() ?? ''
  const rest = parts.slice(1).join(' ')
  const match = resourceCommands.find((c) => c.aliases.some((a) => a === `:${cmdPart}`))
  if (!match?.kind) return {}
  let filter = ''
  let namespace: string | undefined
  if (rest) {
    const nsMatch = rest.match(/-n\s+(\S+)/)
    if (nsMatch) namespace = nsMatch[1]
    const slash = rest.match(/\/(.+)/)
    if (slash) filter = slash[1]
    else if (/\w+=\S+/.test(rest)) filter = rest
  }
  return { kind: match.kind, filter, namespace }
}

export function searchCommands(query: string): CommandDefinition[] {
  const q = query.toLowerCase()
  return resourceCommands.filter(
    (c) => c.title.toLowerCase().includes(q) || c.aliases.some((a) => a.includes(q)),
  )
}
