import { useEffect, useMemo, useState } from 'react'
import { allCommands, searchCommands, type CommandDefinition } from '../../lib/commands'
import type { ResourceKindPath, ResourceNavItem } from '../types'

export function CommandPalette({
  open,
  items,
  onClose,
  onSelectKind,
  onRunCommand,
}: {
  open: boolean
  items: ResourceNavItem[]
  onClose?: () => void
  onSelectKind?: (kind: NonNullable<ResourceNavItem['kind']>) => void
  onRunCommand?: (command: CommandDefinition) => void
}) {
  const [query, setQuery] = useState('')
  useEffect(() => {
    if (open) setQuery('')
  }, [open])
  const matches = useMemo(() => {
    const q = query.trim()
    if (!q) return allCommands
    return searchCommands(q)
  }, [query])

  if (!open) return null

  function pick(command: CommandDefinition) {
    if (command.category === 'resource' && command.kind) {
      onSelectKind?.(command.kind)
    } else {
      onRunCommand?.(command)
    }
    onClose?.()
    setQuery('')
  }

  return (
    <div className="km-fast-overlay" role="dialog" aria-label="Command palette" data-testid="fast-palette" data-fast-overlay="">
      <div>Command palette</div>
      <input
        aria-label="Search commands"
        data-testid="fast-palette-search"
        value={query}
        onChange={(event) => setQuery(event.target.value)}
        placeholder="Search :commands and views…"
        autoFocus
        style={{
          width: '100%',
          marginBottom: 8,
          background: 'transparent',
          border: '1px solid var(--km-border)',
          borderRadius: 'var(--km-radius)',
        }}
      />
      <ul>
        {matches.map((command) => (
          <li key={command.id}>
            <button type="button" onClick={() => pick(command)}>
              {command.aliases[0]} — {command.title}
            </button>
          </li>
        ))}
      </ul>
      <details>
        <summary>Resource views ({items.filter((item) => item.available && item.kind).length})</summary>
        <ul>
          {items
            .filter((item) => item.available && item.kind)
            .map((item) => (
              <li key={item.id}>
                <button
                  type="button"
                  onClick={() => {
                    if (item.kind) onSelectKind?.(item.kind as ResourceKindPath)
                    onClose?.()
                    setQuery('')
                  }}
                >
                  [{item.shortcut}] {item.label}
                </button>
              </li>
            ))}
        </ul>
      </details>
      <button type="button" onClick={onClose}>
        Close
      </button>
    </div>
  )
}
