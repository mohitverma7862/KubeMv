import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../services/api'
import { useSessionStore } from '../stores/sessionStore'
import '../fast/theme/kubemv-terminal.css'

export function LoginPage() {
  const navigate = useNavigate()
  const setSession = useSessionStore((s) => s.setSession)
  const [username, setUsername] = useState('platform-admin')
  const [password, setPassword] = useState('dev')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function onSubmit(event: React.FormEvent) {
    event.preventDefault()
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
    <div className="km-fast" style={{ height: '100%' }}>
      <header className="km-fast-header">
        <span className="km-fast-brand">KubeMv</span>
        <span>Fast Mode (K9s-Parity)</span>
      </header>
      <form onSubmit={onSubmit} style={{ padding: 16 }} aria-label="Sign in">
        <div>sign in</div>
        <label className="km-fast-muted" style={{ display: 'block', marginTop: 12 }}>
          user
          <input
            aria-label="Username"
            value={username}
            onChange={(event) => setUsername(event.target.value)}
            style={{ marginLeft: 8 }}
          />
        </label>
        <label className="km-fast-muted" style={{ display: 'block', marginTop: 8 }}>
          password
          <input
            aria-label="Password"
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            style={{ marginLeft: 8 }}
          />
        </label>
        {error ? <p style={{ color: 'var(--km-error)' }}>{error}</p> : null}
        <button type="submit" className="km-fast-tab" disabled={loading} style={{ marginTop: 12 }}>
          {loading ? 'entering…' : 'enter'}
        </button>
      </form>
    </div>
  )
}
