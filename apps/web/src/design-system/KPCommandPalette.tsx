import { KPButton } from './KPButton'

export function KPCommandPalette({ onOpen }: { onOpen?: () => void }) {
  return (
    <button
      type="button"
      onClick={onOpen}
      className="flex w-full max-w-xl items-center justify-between rounded-md border border-[var(--color-kp-border)] bg-[var(--color-kp-surface-2)] px-3 py-2 text-left text-sm text-[var(--color-kp-muted)] hover:border-[var(--color-kp-accent)]/50"
    >
      <span>Command palette — press Ctrl+K (Phase 1)</span>
      <kbd className="rounded border border-[var(--color-kp-border)] px-1.5 py-0.5 text-[10px]">Ctrl K</kbd>
    </button>
  )
}

export function KPCommandPaletteHint() {
  return (
    <div className="rounded-md border border-dashed border-[var(--color-kp-border)] p-4 text-sm text-[var(--color-kp-muted)]">
      K9s-style <code className="text-[var(--color-kp-text)]">:</code> commands and searchable actions land in Phase 1.
      <div className="mt-3">
        <KPButton variant="secondary" size="sm" disabled>Coming in Phase 1</KPButton>
      </div>
    </div>
  )
}
