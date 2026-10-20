import { describe, expect, it } from 'vitest'
import { freshView, goBack, visitKind } from './navigation'

describe('navigation history', () => {
  it('restores the previous view on back and keeps a remembered visit', () => {
    const pods = { ...freshView('pods'), selectedKey: 'production/Pod/api', filter: 'web', scrollTop: 64 }
    const forward = visitKind([], pods, 'deployments')
    expect(forward.current.kind).toBe('deployments')
    expect(forward.current.filter).toBe('')

    const back = goBack(forward.stack)
    expect(back?.current.selectedKey).toBe('production/Pod/api')
    expect(back?.current.filter).toBe('web')
    expect(back?.current.scrollTop).toBe(64)

    const again = visitKind(back?.stack ?? [], back?.current ?? pods, 'deployments', forward.current)
    expect(again.current.kind).toBe('deployments')
  })
})
