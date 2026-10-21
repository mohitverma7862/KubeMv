import { helpSections } from './registry'

export function HelpOverlay({ open, onClose }: { open: boolean; onClose?: () => void }) {
  if (!open) return null
  const sections = helpSections()
  return (
    <div className="km-fast-overlay" role="dialog" aria-label="Help" data-testid="fast-help" data-fast-overlay="">
      <div>Quick Help (k9s-style)</div>
      {Object.entries(sections).map(([title, items]) =>
        items.length === 0 ? null : (
          <section key={title}>
            <h3>{title}</h3>
            <ul>
              {items.map((item) => (
                <li key={item.id}>
                  <kbd>{item.keys}</kbd> {item.label}
                </li>
              ))}
            </ul>
          </section>
        ),
      )}
      <button type="button" onClick={onClose}>
        Close
      </button>
    </div>
  )
}
