import { create } from 'zustand'
import type { ResourceDetail, ResourceKindPath, ResourceRow } from '../services/types'

interface ResourceState {
  kind: ResourceKindPath
  filter: string
  labelFilter: string
  selected: ResourceRow | null
  detail: ResourceDetail | null
  setKind: (kind: ResourceKindPath) => void
  setFilter: (filter: string) => void
  setLabelFilter: (label: string) => void
  select: (row: ResourceRow | null) => void
  setDetail: (detail: ResourceDetail | null) => void
}

export const useResourceStore = create<ResourceState>((set) => ({
  kind: 'pods',
  filter: '',
  labelFilter: '',
  selected: null,
  detail: null,
  setKind: (kind) => set({ kind, selected: null, detail: null }),
  setFilter: (filter) => set({ filter }),
  setLabelFilter: (labelFilter) => set({ labelFilter }),
  select: (selected) => set({ selected, detail: null }),
  setDetail: (detail) => set({ detail }),
}))
