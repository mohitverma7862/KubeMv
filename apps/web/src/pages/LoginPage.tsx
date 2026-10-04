import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { KPButton } from '../design-system/KPButton'
import { api } from '../services/api'
import { useSessionStore } from '../stores/sessionStore'

export function LoginPage() {
  const navigate = useNavigate()
  const setSession = useSessionStore((s) => s.setSession)
  const [username, setUsername] = useState('platform-admin')
  const [password, setPassword] = useState('dev')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    setError(null)
    try {
      const result = await api.login(username, password)
      setSession(result.token, result.principal)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-full items-center justify-center p-6">
      <form
        onSubmit={onSubmit}
        className="w-full max-w-md rounded-xl border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] p-6 shadow-xl"
      >
        <h1 className="text-2xl font-semibold">KubeMv</h1>
        <p className="mt-1 text-sm text-[var(--color-kp-muted)]">
          Phase 0 dev login — credentials are not stored in source control.
        </p>
        <label className="mt-6 block text-sm">
          Username
          <input
            className="mt-1 w-full rounded-md border border-[var(--color-kp-border)] bg-[var(--color-kp-surface-2)] px-3 py-2"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </label>
        <label className="mt-4 block text-sm">
          Password
          <input
            type="password"
            className="mt-1 w-full rounded-md border border-[var(--color-kp-border)] bg-[var(--color-kp-surface-2)] px-3 py-2"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </label>
        {error ? <p className="mt-3 text-sm text-[var(--color-kp-critical)]">{error}</p> : null}
        <KPButton className="mt-6 w-full" type="submit" disabled={loading}>
          {loading ? 'Signing in…' : 'Sign in'}
        </KPButton>
      </form>
    </div>
  )
}
