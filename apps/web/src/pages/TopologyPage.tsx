import { useSearchParams } from 'react-router-dom'
import { TopologyExplorer } from '../features/topology/TopologyExplorer'

export function TopologyPage() {
  const [params] = useSearchParams()
  const root = params.get('root') ?? undefined
  const mode = (params.get('mode') as 'workload' | 'network' | 'dependency') ?? 'workload'
  return <TopologyExplorer initialRoot={root} initialMode={mode} />
}
