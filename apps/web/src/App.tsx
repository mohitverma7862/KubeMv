import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { FastModePage } from './pages/FastModePage'
import { LoginPage } from './pages/LoginPage'
import { useSessionStore } from './stores/sessionStore'

const queryClient = new QueryClient()

function Console() {
  const token = useSessionStore((s) => s.token)
  if (!token) return <Navigate to="/login" replace />
  return <FastModePage />
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/*" element={<Console />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
