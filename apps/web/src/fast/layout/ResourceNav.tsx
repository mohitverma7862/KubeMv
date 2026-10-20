import type { ResourceKindPath, ResourceNavItem } from '../types'

export function ResourceNav({
  items,
  activeKind,
  onSelectKind,
}: {
  items: ResourceNavItem[]
  activeKind?: ResourceKindPath
  onSelectKind?: (kind: ResourceKindPath) => void
}) {
  return (
    <nav className="km-fast-nav" aria-label="Resources" data-testid="fast-nav">
      {items.map((item) => (
        <button
          key={item.id}
          type="button"
          className="km-fast-tab"
          aria-current={item.kind && item.kind === activeKind ? 'page' : undefined}
          disabled={!item.available || !item.kind}
          title={item.available ? item.label : `${item.label} is not available in this API`}
          onClick={() => {
            if (item.kind) onSelectKind?.(item.kind)
          }}
        >
          [{item.shortcut}] {item.label}
        </button>
      ))}
    </nav>
  )
}
