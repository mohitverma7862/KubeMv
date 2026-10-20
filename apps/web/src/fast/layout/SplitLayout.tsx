import type { ReactNode } from 'react'

export function SplitLayout({
  main,
  side,
  testId,
}: {
  main: ReactNode
  side: ReactNode
  testId: string
}) {
  return (
    <div className="km-fast-split" data-testid={testId}>
      <div className="km-fast-panel">{main}</div>
      <div className="km-fast-panel km-fast-panel-alt">{side}</div>
    </div>
  )
}
