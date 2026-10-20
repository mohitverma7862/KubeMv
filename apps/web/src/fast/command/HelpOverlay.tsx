import { shortcuts } from './registry'

export function HelpOverlay({ open, onClose }: { open: boolean; onClose?: () => void }) {
  if (!open) return null
  return (
    <div className="km-fast-overlay" role="dialog" aria-label="Help" data-testid="fast-help">
      <div>Quick Help</div>
      <ul>
        {shortcuts.map((item) => (
          <li key={item.id}>
            <kbd>{item.keys}</kbd> {item.label}
          </li>
        ))}
      </ul>
      <button type="button" onClick={onClose}>
        Close
      </button>
    </div>
  )
}
