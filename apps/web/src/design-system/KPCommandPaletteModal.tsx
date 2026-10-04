import { useEffect, useMemo, useState } from 'react'
import { searchCommands, type CommandDefinition } from '../lib/commands'
import { cn } from '../lib/cn'

export function KPCommandPaletteModal({
  open,
  onClose,
  onSelect,
}: {
  open: boolean
  onClose: () => void
  onSelect: (command: CommandDefinition) => void
}) {
  const [query, setQuery] = useState('')
  const results = useMemo(() => searchCommands(query), [query])

  useEffect(() => {
    if (!open) setQuery('')
  }, [open])

  useEffect(() => {
    if (!open) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [open, onClose])

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center bg-black/60 p-8" onClick={onClose}>
      <div
        className="w-full max-w-xl rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] shadow-2xl"
        onClick={(e) => e.stopPropagation()}
      >
        <input
          autoFocus
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search commands…"
          className="w-full border-b border-[var(--color-kp-border)] bg-transparent px-4 py-3 text-sm outline-none"
        />
        <ul className="max-h-80 overflow-auto py-2">
          {results.map((cmd) => (
            <li key={cmd.id}>
              <button
                type="button"
                className={cn(
                  'flex w-full items-center justify-between px-4 py-2 text-left text-sm hover:bg-[var(--color-kp-surface-2)]',
                )}
                onClick={() => {
                  onSelect(cmd)
                  onClose()
                }}
              >
                <span>{cmd.title}</span>
                <span className="text-xs text-[var(--color-kp-muted)]">{cmd.aliases[0]}</span>
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  )
}
