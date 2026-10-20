import type { ResourceDetail } from '../../services/types'

export function ResourceInspector({ detail }: { detail?: ResourceDetail | null }) {
  return (
    <aside data-testid="fast-inspector" style={{ padding: 10 }}>
      <div className="km-fast-muted">Resource Inspector</div>
      {detail ? (
        <dl>
          <dt>Name</dt>
          <dd>{detail.row.name}</dd>
          <dt>Status</dt>
          <dd>{detail.row.status}</dd>
          <dt>Namespace</dt>
          <dd>{detail.row.namespace || '—'}</dd>
          <dt>Age</dt>
          <dd>{detail.row.age}</dd>
          {detail.row.extra?.ready ? (
            <>
              <dt>Ready</dt>
              <dd>{detail.row.extra.ready}</dd>
            </>
          ) : null}
          {detail.row.extra?.node ? (
            <>
              <dt>Node</dt>
              <dd>{detail.row.extra.node}</dd>
            </>
          ) : null}
          {detail.row.extra?.ip ? (
            <>
              <dt>IP</dt>
              <dd>{detail.row.extra.ip}</dd>
            </>
          ) : null}
        </dl>
      ) : (
        <p className="km-fast-muted">No resource selected.</p>
      )}
    </aside>
  )
}
