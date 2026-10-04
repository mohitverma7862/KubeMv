import { QueryClient, QueryClientProvider, useQuery } from '@tanstack/react-query'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { AppShell } from './layouts/AppShell'
import { FastModePage } from './pages/FastModePage'
// workloads route reuses Fast Mode explorer in Phase 1
import { LoginPage } from './pages/LoginPage'
import { OverviewPage } from './pages/OverviewPage'
import { SecurityPage } from './pages/SecurityPage'
import { GitOpsPage } from './pages/GitOpsPage'
import { TopologyPage } from './pages/TopologyPage'
import { ObservabilityPage } from './pages/ObservabilityPage'
import { VisualModePage } from './pages/VisualModePage'
import { api } from './services/api'
import { useSessionStore } from './stores/sessionStore'

const queryClient = new QueryClient()

function ProtectedApp() {
  const token = useSessionStore((s) => s.token)
  const { data: clusters = [] } = useQuery({
    queryKey: ['clusters'],
    enabled: Boolean(token),
    queryFn: () => api.clusters(token!),
  })

  if (!token) {
    return <Navigate to="/login" replace />
  }

  return (
    <Routes>
      <Route element={<AppShell clusters={clusters} />}>
        <Route index element={<OverviewPage />} />
        <Route path="fast" element={<FastModePage />} />
        <Route path="visual" element={<VisualModePage />} />
        <Route path="workloads" element={<FastModePage />} />
        <Route path="topology" element={<TopologyPage />} />
        <Route path="network" element={<TopologyPage />} />
        <Route path="observability" element={<ObservabilityPage />} />
        <Route path="security" element={<SecurityPage />} />
        <Route path="gitops" element={<GitOpsPage />} />
      </Route>
    </Routes>
  )
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/*" element={<ProtectedApp />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
