import { SecurityExplorer } from '../features/security/SecurityExplorer'

export function SecurityPage() {
  return (
    <div className="flex h-full flex-col gap-4">
      <header>
        <h1 className="text-lg font-semibold">Security</h1>
        <p className="text-sm text-[var(--color-kp-muted)]">
          RBAC explorer and policy hints — application roles never bypass Kubernetes authorization.
        </p>
      </header>
      <SecurityExplorer />
    </div>
  )
}
