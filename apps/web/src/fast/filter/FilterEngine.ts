import type { ResourceRow } from '../../services/types'

export type ParsedFilter = {
  text: string
  labels: Record<string, string>
  status?: string
  regex?: RegExp
  invalidRegex?: string
}

export function parseFilter(input: string): ParsedFilter {
  const labels: Record<string, string> = {}
  const words: string[] = []
  let status: string | undefined
  let regex: RegExp | undefined
  let invalidRegex: string | undefined

  for (const token of input.trim().split(/\s+/).filter(Boolean)) {
    if (token.startsWith('status=')) {
      status = token.slice('status='.length)
      continue
    }
    if (token.startsWith('/') && token.endsWith('/') && token.length > 2) {
      const source = token.slice(1, -1)
      try {
        regex = new RegExp(source, 'i')
      } catch {
        invalidRegex = source
      }
      continue
    }
    const label = token.match(/^([A-Za-z0-9_.-]+)=(.+)$/)
    if (label) {
      labels[label[1]] = label[2]
      continue
    }
    words.push(token.toLowerCase())
  }

  return { text: words.join(' '), labels, status, regex, invalidRegex }
}

export function rowMatches(row: ResourceRow, filter: ParsedFilter): boolean {
  if (filter.invalidRegex) return false
  if (filter.status && row.status.toLowerCase() !== filter.status.toLowerCase()) return false
  for (const [key, value] of Object.entries(filter.labels)) {
    if (row.labels?.[key] !== value) return false
  }
  const haystack = `${row.name} ${row.namespace} ${row.status} ${row.kind}`.toLowerCase()
  if (filter.text && !haystack.includes(filter.text)) return false
  if (filter.regex && !filter.regex.test(haystack)) return false
  return true
}

export function applyFilter(rows: ResourceRow[], input: string): ResourceRow[] {
  const parsed = parseFilter(input)
  return rows.filter((row) => rowMatches(row, parsed))
}
