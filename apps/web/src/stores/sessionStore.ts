import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { Principal } from '../services/types'

interface SessionState {
  token: string | null
  principal: Principal | null
  setSession: (token: string, principal: Principal) => void
  clear: () => void
}

export const useSessionStore = create<SessionState>()(
  persist(
    (set) => ({
      token: null,
      principal: null,
      setSession: (token, principal) => set({ token, principal }),
      clear: () => set({ token: null, principal: null }),
    }),
    { name: 'kubemv-session' },
  ),
)
