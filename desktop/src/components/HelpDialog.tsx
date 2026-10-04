import { useUiStore } from "../stores/ui";

const shortcuts = [
  ["Ctrl K", "Open the command palette"],
  ["?", "Show this shortcut list"],
  ["Esc", "Close the palette, dialog, or help"],
  ["g o", "Go to Overview"],
  ["g c", "Go to Clusters"],
  ["g p", "Go to Plugins"],
  ["g s", "Go to Settings"],
];

export function HelpDialog() {
  const open = useUiStore((state) => state.helpOpen);
  const setOpen = useUiStore((state) => state.setHelpOpen);
  if (!open) {
    return null;
  }
  return (
    <div className="km-modal-back" onMouseDown={() => setOpen(false)}>
      <div
        className="km-modal"
        role="dialog"
        aria-modal="true"
        aria-labelledby="help-title"
        data-testid="help-dialog"
        onMouseDown={(event) => event.stopPropagation()}
      >
        <header id="help-title">Keyboard</header>
        <div className="km-dialog">
          <dl className="km-shortcut-list">
            {shortcuts.map(([keys, description]) => (
              <span key={keys} style={{ display: "contents" }}>
                <dt>{keys}</dt>
                <dd>{description}</dd>
              </span>
            ))}
          </dl>
        </div>
        <footer className="km-actions">
          <button type="button" className="km-btn" onClick={() => setOpen(false)}>
            Close
          </button>
        </footer>
      </div>
    </div>
  );
}
