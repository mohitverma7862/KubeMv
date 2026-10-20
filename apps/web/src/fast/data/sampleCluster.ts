import type { ResourceDetail, ResourceKindPath, ResourceRow } from '../../services/types'
import { rowKey } from '../resource/rowKey'
import type { OwnerNode } from '../types'
import type { FastDataSource } from './source'

function row(
  kind: string,
  name: string,
  status: string,
  age: string,
  extra: Record<string, string> = {},
  labels: Record<string, string> = {},
  namespace = 'production',
): ResourceRow {
  return { kind, namespace, name, status, age, extra, labels }
}

const apiPod = row(
  'Pod',
  'api-7c9f8b6d4-k2m9x',
  'Running',
  '4h',
  { ready: '1/1', restarts: '0', ip: '10.8.1.14', node: 'ip-10-0-1-12', cpu: '42m', mem: '128Mi' },
  { app: 'api', tier: 'backend' },
)

const workerPod = row(
  'Pod',
  'payments-worker-6b8d4c7f9-q1w2e',
  'CrashLoopBackOff',
  '27m',
  { ready: '0/1', restarts: '14', ip: '10.8.1.22', node: 'ip-10-0-1-18', cpu: '8m', mem: '64Mi' },
  { app: 'payments-worker', tier: 'backend' },
)

const checkoutPod = row(
  'Pod',
  'checkout-5f4d3c2b1-p9n4c',
  'Pending',
  '2m',
  { ready: '0/1', restarts: '0', ip: '—', node: '—', cpu: '—', mem: '—' },
  { app: 'checkout' },
)

const webPod = row(
  'Pod',
  'web-6d5c4b3a2-h7t1a',
  'Running',
  '1d',
  { ready: '1/1', restarts: '1', ip: '10.8.2.4', node: 'ip-10-0-1-12', cpu: '18m', mem: '96Mi' },
  { app: 'web', tier: 'frontend' },
)

const pods = [apiPod, workerPod, checkoutPod, webPod]

const deployments = [
  row('Deployment', 'api', 'Available', '12d', { ready: '3/3', restarts: '0' }, { app: 'api' }),
  row('Deployment', 'payments-worker', 'Progressing', '12d', { ready: '0/2', restarts: '14' }, { app: 'payments-worker' }),
  row('Deployment', 'checkout', 'Available', '6d', { ready: '1/1', restarts: '0' }, { app: 'checkout' }),
  row('Deployment', 'web', 'Available', '20d', { ready: '2/2', restarts: '1' }, { app: 'web' }),
]

const services = [
  row('Service', 'api', 'ClusterIP', '12d', { ip: '10.96.0.21', ports: '80/TCP' }, { app: 'api' }),
  row('Service', 'web', 'LoadBalancer', '20d', { ip: 'pending', ports: '443/TCP' }, { app: 'web' }),
]

const nodes = [
  row('Node', 'ip-10-0-1-12', 'Ready', '40d', { cpu: '42%', mem: '61%' }, {}, ''),
  row('Node', 'ip-10-0-1-18', 'Ready', '40d', { cpu: '18%', mem: '44%' }, {}, ''),
  row('Node', 'ip-10-0-2-4', 'NotReady', '3h', { cpu: '—', mem: '—' }, {}, ''),
]

const events = [
  row('Event', 'payments-worker.18a3', 'Warning', '1m', { reason: 'BackOff', object: 'payments-worker-6b8d4c7f9-q1w2e' }),
  row('Event', 'api.9c21', 'Normal', '4h', { reason: 'Started', object: 'api-7c9f8b6d4-k2m9x' }),
  row('Event', 'checkout.11b0', 'Warning', '2m', { reason: 'FailedScheduling', object: 'checkout-5f4d3c2b1-p9n4c' }),
]

const jobs = [
  row('Job', 'migrate-db', 'Complete', '6h', { ready: '1/1' }),
  row('Job', 'nightly-report', 'Running', '8m', { ready: '0/1' }),
]

const cronjobs = [row('CronJob', 'nightly-report', 'Active', '20d', { schedule: '0 2 * * *' })]

const configmaps = [row('ConfigMap', 'api-config', '2 keys', '12d', { keys: 'LOG_LEVEL, UPSTREAM' })]

const secrets = [row('Secret', 'api-tls', 'kubernetes.io/tls', '12d', { keys: '2' })]

const namespaces = [
  row('Namespace', 'production', 'Active', '40d', {}, {}, ''),
  row('Namespace', 'payments', 'Active', '30d', {}, {}, ''),
]

const statefulsets = [row('StatefulSet', 'redis', 'Available', '20d', { ready: '1/1' })]
const daemonsets = [row('DaemonSet', 'node-exporter', 'Available', '40d', { ready: '3/3' })]
const hpa = [row('HorizontalPodAutoscaler', 'api', 'Active', '12d', { ready: '3/10' })]
const pdb = [row('PodDisruptionBudget', 'api', 'Sufficient', '12d', { ready: '2/3' })]
const crds = [row('CustomResourceDefinition', 'prometheusrules.monitoring.coreos.com', 'Established', '40d', {}, {}, '')]

const rows: Partial<Record<ResourceKindPath, ResourceRow[]>> = {
  pods,
  deployments,
  services,
  nodes,
  events,
  jobs,
  cronjobs,
  configmaps,
  secrets,
  namespaces,
  statefulsets,
  daemonsets,
  hpa,
  pdb,
  crds,
}

function manifest(resource: ResourceRow): string {
  return [
    'apiVersion: v1',
    `kind: ${resource.kind}`,
    'metadata:',
    `  name: ${resource.name}`,
    resource.namespace ? `  namespace: ${resource.namespace}` : '  namespace: ""',
    '  # sample manifest; secret data is omitted',
  ].join('\n')
}

function detail(resource: ResourceRow, message: string): ResourceDetail {
  return {
    row: resource,
    yaml: manifest(resource),
    events: [
      {
        type: resource.status === 'Running' || resource.status === 'Available' ? 'Normal' : 'Warning',
        reason: resource.status,
        message,
        object: resource.name,
        age: resource.age,
        namespace: resource.namespace,
      },
    ],
  }
}

const details: Record<string, ResourceDetail> = {}
const trees: Record<string, OwnerNode[]> = {}
const logs: Record<string, string> = {}
const describe: Record<string, string> = {}

for (const list of Object.values(rows)) {
  for (const resource of list ?? []) {
    const key = rowKey(resource)
    details[key] = detail(resource, `${resource.kind} ${resource.name} is ${resource.status}.`)
    describe[key] = [
      `Name: ${resource.name}`,
      `Namespace: ${resource.namespace || '—'}`,
      `Status: ${resource.status}`,
      `Node: ${resource.extra?.node ?? '—'}`,
      `IP: ${resource.extra?.ip ?? '—'}`,
      `Age: ${resource.age}`,
    ].join('\n')
    logs[key] = `${resource.name} sample log\nlevel=info msg="ready" status=${resource.status}`
  }
}

const apiKey = rowKey(apiPod)
trees[apiKey] = [
  {
    id: 'deploy-api',
    kind: 'Deployment',
    name: 'api',
    children: [
      {
        id: 'rs-api',
        kind: 'ReplicaSet',
        name: 'api-7c9f8b6d4',
        children: [
          {
            id: 'pod-api',
            kind: 'Pod',
            name: apiPod.name,
            children: [{ id: 'ctr-api', kind: 'Container', name: 'api' }],
          },
        ],
      },
    ],
  },
]
logs[apiKey] = 'api sample log\nlevel=info msg="listening" addr=:8080\nlevel=info msg="upstream ok"'
logs[rowKey(workerPod)] = [
  'payments-worker sample log',
  'level=error msg="upstream refused" addr=payments:5432',
  'level=error msg="backoff" restart=14',
].join('\n')

export const sampleCluster: FastDataSource = {
  id: 'sample',
  header: {
    clusterName: 'prod-cluster',
    userLabel: 'admin',
    contextName: 'prod',
    namespace: 'production',
    production: true,
    readOnly: true,
    connection: 'live',
  },
  rows,
  details,
  trees,
  logs,
  describe,
}
