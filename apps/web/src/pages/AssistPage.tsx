import { AssistExplorer } from '../features/assist/AssistExplorer'

export function AssistPage() {
  return (
    <div className="flex h-full flex-col gap-4">
      <header>
        <h1 className="text-lg font-semibold">AI Assist</h1>
        <p className="text-sm text-[var(--color-kp-muted)]">
          Rule-based triage, suggested runbooks, and dry-run automation hooks — no cluster mutations without approval.
        </p>
      </header>
      <AssistExplorer />
    </div>
  )
}
