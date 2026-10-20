import { useState } from 'react'
import type { ResourceKindPath, ResourceNavItem } from '../types'

export function ResourceNav({
  items,
  moreItems = [],
  activeKind,
  onSelectKind,
}: {
  items: ResourceNavItem[]
  moreItems?: ResourceNavItem[]
  activeKind?: ResourceKindPath
  onSelectKind?: (kind: ResourceKindPath) => void
}) {
  const [moreOpen, setMoreOpen] = useState(false)
  const moreActive = moreItems.some((item) => item.kind === activeKind)
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
      {moreItems.length > 0 ? (
        <div style={{ position: 'relative' }}>
          <button
            type="button"
            className="km-fast-tab"
            aria-expanded={moreOpen}
            aria-current={moreActive ? 'page' : undefined}
            onClick={() => setMoreOpen((open) => !open)}
          >
            [&gt;] More
          </button>
          {moreOpen ? (
            <ul className="km-fast-more" data-testid="fast-more">
              {moreItems.map((item) => (
                <li key={item.id}>
                  <button
                    type="button"
                    className="km-fast-tab"
                    disabled={!item.available || !item.kind}
                    onClick={() => {
                      if (item.kind) onSelectKind?.(item.kind)
                      setMoreOpen(false)
                    }}
                  >
                    {item.label}
                  </button>
                </li>
              ))}
            </ul>
          ) : null}
        </div>
      ) : null}
    </nav>
  )
}
