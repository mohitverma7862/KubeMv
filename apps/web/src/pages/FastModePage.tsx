import { KPCommandPaletteHint } from '../design-system/KPCommandPalette'
import { useUIStore } from '../stores/uiStore'

export function FastModePage() {
  const { clusterID, namespace } = useUIStore()

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-semibold">Fast Mode</h1>
      <p className="text-sm text-[var(--color-kp-muted)]">
        Keyboard-first operator surface (K9s-inspired). Context: cluster <code>{clusterID}</code>, namespace{' '}
        <code>{namespace}</code>.
      </p>
      <KPCommandPaletteHint />
    </div>
  )
}
