import { parseColonCommand, resourceCommands } from '../../lib/commands'

export function CommandAutocomplete({ value }: { value: string }) {
  if (!value.startsWith(':')) return null
  const parsed = parseColonCommand(value)
  const prefix = value.slice(1).split(/\s+/)[0]?.toLowerCase() ?? ''
  const matches = resourceCommands.filter((command) =>
    command.aliases.some((alias) => alias.slice(1).startsWith(prefix)),
  )
  if (matches.length === 0) return null
  return (
    <ul data-testid="fast-command-suggest" className="km-fast-muted">
      {matches.slice(0, 6).map((command) => (
        <li key={command.id}>
          {command.aliases[0]}{' '}
          {parsed.type === 'resource' && parsed.kind === command.kind ? '(match)' : ''}
        </li>
      ))}
    </ul>
  )
}
