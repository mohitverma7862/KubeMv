import type { ResourceNavItem } from '../types'

export function CommandPalette({
  open,
  items,
  onClose,
  onSelectKind,
}: {
  open: boolean
  items: ResourceNavItem[]
  onClose?: () => void
  onSelectKind?: (kind: NonNullable<ResourceNavItem['kind']>) => void
}) {
  if (!open) return null
  return (
    <div className="km-fast-overlay" role="dialog" aria-label="Command palette" data-testid="fast-palette">
      <div>Command palette</div>
      <ul>
        {items
          .filter((item) => item.available && item.kind)
          .map((item) => (
            <li key={item.id}>
              <button
                type="button"
                onClick={() => {
                  if (item.kind) onSelectKind?.(item.kind)
                  onClose?.()
                }}
              >
                [{item.shortcut}] {item.label}
              </button>
            </li>
          ))}
      </ul>
      <button type="button" onClick={onClose}>
        Close
      </button>
    </div>
  )
}
