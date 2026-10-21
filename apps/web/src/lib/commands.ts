import type { ResourceKindPath } from '../services/types'

export interface CommandDefinition {
  id: string
  title: string
  aliases: string[]
  kind?: ResourceKindPath
  shortcut?: string
  risk?: 'READ' | 'LOW' | 'MEDIUM' | 'HIGH'
  category?: 'resource' | 'system'
}

/** Resource views aligned with k9s aliases (see k9s internal/view aliases). */
export const resourceCommands: CommandDefinition[] = [
  { id: 'pods', title: 'Pods', aliases: [':pods', ':po', ':pod', 'pods'], kind: 'pods', category: 'resource' },
  {
    id: 'deployments',
    title: 'Deployments',
    aliases: [':deploy', ':deployments', ':dp', ':deploys'],
    kind: 'deployments',
    category: 'resource',
  },
  { id: 'services', title: 'Services', aliases: [':svc', ':services', ':service'], kind: 'services', category: 'resource' },
  { id: 'nodes', title: 'Nodes', aliases: [':nodes', ':no', ':node'], kind: 'nodes', category: 'resource' },
  { id: 'events', title: 'Events', aliases: [':events', ':ev', ':event'], kind: 'events', category: 'resource' },
  { id: 'namespaces', title: 'Namespaces', aliases: [':ns', ':namespaces', ':namespace'], kind: 'namespaces', category: 'resource' },
  { id: 'configmaps', title: 'ConfigMaps', aliases: [':configmap', ':configmaps', ':cm'], kind: 'configmaps', category: 'resource' },
  { id: 'secrets', title: 'Secrets (metadata)', aliases: [':secret', ':secrets', ':sec'], kind: 'secrets', category: 'resource' },
  { id: 'crds', title: 'CRDs', aliases: [':crd', ':crds'], kind: 'crds', category: 'resource' },
  {
    id: 'statefulsets',
    title: 'StatefulSets',
    aliases: [':sts', ':statefulsets', ':statefulset'],
    kind: 'statefulsets',
    category: 'resource',
  },
  { id: 'daemonsets', title: 'DaemonSets', aliases: [':ds', ':daemonsets', ':daemonset'], kind: 'daemonsets', category: 'resource' },
  { id: 'jobs', title: 'Jobs', aliases: [':job', ':jobs'], kind: 'jobs', category: 'resource' },
  { id: 'cronjobs', title: 'CronJobs', aliases: [':cronjob', ':cronjobs', ':cj'], kind: 'cronjobs', category: 'resource' },
  { id: 'hpa', title: 'HPA', aliases: [':hpa', ':hpas'], kind: 'hpa', category: 'resource' },
  { id: 'pdb', title: 'PDB', aliases: [':pdb', ':pdbs'], kind: 'pdb', category: 'resource' },
]

export const systemCommands: CommandDefinition[] = [
  { id: 'quit', title: 'Quit', aliases: [':q', ':quit', ':exit'], category: 'system', risk: 'LOW' },
  { id: 'help', title: 'Help', aliases: [':h', ':help'], category: 'system' },
  { id: 'ctx', title: 'Contexts', aliases: [':ctx', ':context', ':contexts'], category: 'system' },
  { id: 'clear', title: 'Clear filter', aliases: [':clear', ':clr'], category: 'system' },
]

export const allCommands: CommandDefinition[] = [...resourceCommands, ...systemCommands]

export type ColonParseResult =
  | { type: 'resource'; kind: ResourceKindPath; filter?: string; namespace?: string }
  | { type: 'quit' }
  | { type: 'help' }
  | { type: 'ctx' }
  | { type: 'clear-filter' }
  | { type: 'unknown' }

export function parseColonCommand(input: string): ColonParseResult {
  const trimmed = input.trim()
  if (!trimmed.startsWith(':')) return { type: 'unknown' }
  const parts = trimmed.slice(1).trim().split(/\s+/)
  const cmdPart = parts[0]?.toLowerCase() ?? ''
  const rest = parts.slice(1).join(' ')

  for (const cmd of systemCommands) {
    if (cmd.aliases.some((a) => a === `:${cmdPart}`)) {
      if (cmd.id === 'quit') return { type: 'quit' }
      if (cmd.id === 'help') return { type: 'help' }
      if (cmd.id === 'ctx') return { type: 'ctx' }
      if (cmd.id === 'clear') return { type: 'clear-filter' }
    }
  }

  const match = resourceCommands.find((c) => c.aliases.some((a) => a === `:${cmdPart}`))
  if (!match?.kind) return { type: 'unknown' }

  let filter = ''
  let namespace: string | undefined
  if (rest) {
    const nsMatch = rest.match(/(?:-n|--namespace)\s+(\S+)/)
    if (nsMatch) namespace = nsMatch[1]
    const slash = rest.match(/\/(.+)/)
    if (slash) filter = slash[1]
    else if (/\w+=\S+/.test(rest)) filter = rest
    else if (!namespace) filter = rest
  }
  return { type: 'resource', kind: match.kind, filter, namespace }
}

export function searchCommands(query: string): CommandDefinition[] {
  const q = query.toLowerCase().replace(/^:/, '')
  return allCommands.filter(
    (c) =>
      c.title.toLowerCase().includes(q) ||
      c.id.includes(q) ||
      c.aliases.some((a) => a.replace(':', '').includes(q)),
  )
}
