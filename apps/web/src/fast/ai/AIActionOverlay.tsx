import type { AIFinding } from '../types'

export function AIActionOverlay({
  open,
  finding,
  onClose,
}: {
  open: boolean
  finding?: AIFinding | null
  onClose?: () => void
}) {
  if (!open) return null
  return (
    <div className="km-fast-overlay" role="dialog" aria-label="AI actions" data-testid="fast-ai" data-fast-overlay="">
      <div style={{ color: 'var(--km-ai)' }}>AI</div>
      {finding ? (
        <>
          <h2>Root Cause</h2>
          <p>{finding.rootCause}</p>
          <h2>Evidence</h2>
          <ul>
            {finding.evidence.map((item) => (
              <li key={item}>{item}</li>
            ))}
          </ul>
          <h2>Recommended Fix</h2>
          <p>{finding.recommendedFix}</p>
          <h2>Commands</h2>
          <ul>
            {finding.commands.map((item) => (
              <li key={item}>
                <code>{item}</code>
              </li>
            ))}
          </ul>
        </>
      ) : (
        <p className="km-fast-muted">No investigation is loaded. Commands are never run from this overlay.</p>
      )}
      <button type="button" onClick={onClose}>
        Close
      </button>
    </div>
  )
}
