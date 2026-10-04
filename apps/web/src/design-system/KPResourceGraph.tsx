import {
  Background,
  Controls,
  MiniMap,
  ReactFlow,
  type Edge,
  type Node,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { useMemo } from 'react'
import type { GraphEdge, GraphNode } from '../services/types'

const kindColumn: Record<string, number> = {
  Internet: 0,
  Ingress: 1,
  Service: 2,
  Deployment: 3,
  ReplicaSet: 4,
  Pod: 5,
  ConfigMap: 6,
  Secret: 6,
  ServiceAccount: 6,
}

function healthColor(score: number) {
  if (score >= 85) return '#3ecf8e'
  if (score >= 60) return '#f5b942'
  return '#ff5d5d'
}

function toFlow(nodes: GraphNode[], edges: GraphEdge[]): { nodes: Node[]; edges: Edge[] } {
  const columnCounts: Record<number, number> = {}
  const flowNodes: Node[] = nodes.map((n) => {
    const col = kindColumn[n.kind] ?? 3
    const row = columnCounts[col] ?? 0
    columnCounts[col] = row + 1
    return {
      id: n.id,
      position: { x: col * 220, y: row * 90 },
      data: { label: `${n.kind}\n${n.name}\n${n.healthScore}/100` },
      style: {
        border: `2px solid ${healthColor(n.healthScore)}`,
        background: '#121820',
        color: '#e8edf4',
        fontSize: 11,
        padding: 8,
        width: 180,
        whiteSpace: 'pre-wrap',
      },
    }
  })
  const flowEdges: Edge[] = edges.map((e) => ({
    id: `${e.source}-${e.target}`,
    source: e.source,
    target: e.target,
    label: e.relation,
    animated: e.relation === 'owns',
    style: { stroke: '#4f8cff' },
  }))
  return { nodes: flowNodes, edges: flowEdges }
}

export function KPResourceGraph({ graph }: { graph: { nodes: GraphNode[]; edges: GraphEdge[] } }) {
  const { nodes, edges } = useMemo(() => toFlow(graph.nodes, graph.edges), [graph])

  return (
    <div className="h-[520px] w-full rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)]">
      <ReactFlow nodes={nodes} edges={edges} fitView proOptions={{ hideAttribution: true }}>
        <Background color="#2a3441" gap={16} />
        <MiniMap pannable zoomable />
        <Controls />
      </ReactFlow>
    </div>
  )
}
