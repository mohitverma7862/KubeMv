import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type UIMode = 'fast' | 'visual'

interface UIState {
  mode: UIMode
  clusterID: string
  namespace: string
  setMode: (mode: UIMode) => void
  setClusterID: (id: string) => void
  setNamespace: (ns: string) => void
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      mode: 'visual',
      clusterID: 'local',
      namespace: 'default',
      setMode: (mode) => set({ mode }),
      setClusterID: (clusterID) => set({ clusterID }),
      setNamespace: (namespace) => set({ namespace }),
    }),
    { name: 'kubemv-ui' },
  ),
)
