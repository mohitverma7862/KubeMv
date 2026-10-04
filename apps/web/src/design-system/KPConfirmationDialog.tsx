import { KPButton } from './KPButton'

export function KPConfirmationDialog({
  open,
  title,
  message,
  risk,
  onConfirm,
  onCancel,
}: {
  open: boolean
  title: string
  message: string
  risk: string
  onConfirm: () => void
  onCancel: () => void
}) {
  if (!open) return null
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
      <div className="w-full max-w-md rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] p-4 shadow-xl">
        <h3 className="text-lg font-semibold">{title}</h3>
        <p className="mt-2 text-sm text-[var(--color-kp-muted)]">{message}</p>
        <p className="mt-2 text-xs uppercase tracking-wide text-[var(--color-kp-warning)]">Risk: {risk}</p>
        <div className="mt-4 flex justify-end gap-2">
          <KPButton variant="secondary" onClick={onCancel}>Cancel</KPButton>
          <KPButton variant="danger" onClick={onConfirm}>Confirm</KPButton>
        </div>
      </div>
    </div>
  )
}
