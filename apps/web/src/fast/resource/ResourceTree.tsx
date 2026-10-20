import type { OwnerNode } from '../types'

function TreeList({ nodes, depth }: { nodes: OwnerNode[]; depth: number }) {
  return (
    <ul style={{ listStyle: 'none', margin: 0, paddingLeft: depth === 0 ? 0 : 14 }}>
      {nodes.map((node) => (
        <li key={node.id}>
          {depth > 0 ? '└ ' : ''}
          {node.kind} {node.name}
          {node.children?.length ? <TreeList nodes={node.children} depth={depth + 1} /> : null}
        </li>
      ))}
    </ul>
  )
}

export function ResourceTree({ nodes }: { nodes: OwnerNode[] }) {
  return (
    <section data-testid="fast-tree" style={{ padding: 10 }}>
      <div className="km-fast-muted">Resource Relationship</div>
      {nodes.length === 0 ? <p className="km-fast-muted">No owner chain.</p> : <TreeList nodes={nodes} depth={0} />}
    </section>
  )
}
