import type { ResourceKindPath } from '../../services/types'
import type { ResourceNavItem } from '../types'

/** Digit keys 1–9 jump to primary nav resources (k9s-style view switching). */
export function kindForDigit(items: ResourceNavItem[], digit: string): ResourceKindPath | undefined {
  const index = Number(digit) - 1
  if (!Number.isFinite(index) || index < 0) return undefined
  const available = items.filter((item) => item.available && item.kind)
  return available[index]?.kind
}

/** Letter mnemonics from the reference nav row (e.g. K → pods). */
export function kindForNavLetter(items: ResourceNavItem[], letter: string): ResourceKindPath | undefined {
  const upper = letter.toUpperCase()
  const hit = items.find((item) => item.available && item.kind && item.shortcut.toUpperCase() === upper)
  return hit?.kind
}
