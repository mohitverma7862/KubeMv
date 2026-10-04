import { GitOpsExplorer } from '../features/gitops/GitOpsExplorer'

export function GitOpsPage() {
  return (
    <div className="flex h-full flex-col gap-4">
      <header>
        <h1 className="text-lg font-semibold">GitOps</h1>
        <p className="text-sm text-[var(--color-kp-muted)]">
          Application sync status, manifest drift vs Git, and deployment pipeline activity.
        </p>
      </header>
      <GitOpsExplorer />
    </div>
  )
}
