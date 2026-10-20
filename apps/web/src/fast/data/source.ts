import type { ResourceDetail, ResourceKindPath, ResourceRow } from '../../services/types'
import type { FastHeaderModel, OwnerNode } from '../types'

/** Cluster snapshot the shell can render. Phase 5 replaces the sample source with the API. */
export type FastDataSource = {
  id: 'sample' | 'cluster'
  header: Omit<FastHeaderModel, 'clock'>
  rows: Partial<Record<ResourceKindPath, ResourceRow[]>>
  details: Record<string, ResourceDetail>
  trees: Record<string, OwnerNode[]>
  logs: Record<string, string>
  describe: Record<string, string>
}
