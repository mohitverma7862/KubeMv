import { useMutation, useQuery } from '@tanstack/react-query'
import { useState } from 'react'
import { api } from '../../services/api'
import type { HookDryRunResult } from '../../services/types'
import { useSessionStore } from '../../stores/sessionStore'
import { useUIStore } from '../../stores/uiStore'
import { KPButton } from '../../design-system/KPButton'

export function AssistExplorer() {
  const token = useSessionStore((s) => s.token)!
  const { clusterID, namespace } = useUIStore()
  const [kind, setKind] = useState('Deployment')
  const [name, setName] = useState('payment-api')
  const [dryRun, setDryRun] = useState<HookDryRunResult | null>(null)

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['assist-bundle', clusterID, namespace, kind, name],
    queryFn: () => api.assistBundle(token, clusterID, { namespace, kind, name }),
  })

  const dryRunMutation = useMutation({
    mutationFn: (hookId: string) =>
      api.assistHookDryRun(token, clusterID, hookId, { namespace, kind, name }),
    onSuccess: setDryRun,
  })

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-end gap-3">
        <label className="text-sm">
          Kind
          <select
            className="ml-2 rounded border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] px-2 py-1"
            value={kind}
            onChange={(e) => setKind(e.target.value)}
          >
            <option>Deployment</option>
            <option>Pod</option>
          </select>
        </label>
        <label className="text-sm">
          Name
          <input
            className="ml-2 rounded border border-[var(--color-kp-border)] bg-[var(--color-kp-surface)] px-2 py-1 font-mono"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </label>
        <KPButton size="sm" variant="secondary" onClick={() => refetch()}>
          Analyze
        </KPButton>
      </div>

      {isLoading && <p className="text-sm text-[var(--color-kp-muted)]">Building assist bundle…</p>}

      {data && (
        <>
          <section className="rounded-lg border border-[var(--color-kp-border)] bg-[var(--color-kp-surface-elevated)] p-4">
            <h2 className="text-sm font-semibold uppercase text-[var(--color-kp-muted)]">Triage</h2>
            <p className="mt-2 text-sm">{data.triage.summary}</p>
            <p className="mt-1 text-xs text-[var(--color-kp-muted)]">
              Confidence: {(data.triage.confidence * 100).toFixed(0)}% — {data.triage.disclaimer}
            </p>
            <ul className="mt-3 space-y-2">
              {data.triage.hypotheses.map((h) => (
                <li key={h.title} className="text-sm">
                  <span className="font-medium">{h.title}</span>
                  <span className="text-[var(--color-kp-muted)]"> ({h.likelihood})</span> — {h.evidence}
                </li>
              ))}
            </ul>
            <div className="mt-2 flex flex-wrap gap-2 text-xs font-mono text-[var(--color-kp-muted)]">
              {data.triage.signals.map((s) => (
                <span key={s} className="rounded bg-[var(--color-kp-surface)] px-2 py-0.5">{s}</span>
              ))}
            </div>
          </section>

          <section>
            <h2 className="mb-2 text-sm font-semibold uppercase text-[var(--color-kp-muted)]">
              Runbook — {data.runbook.title}
            </h2>
            <ol className="list-decimal space-y-2 pl-5 text-sm">
              {data.runbook.steps.map((step) => (
                <li key={step.order}>
                  <div className="font-medium">{step.title}</div>
                  {step.command && <pre className="mt-1 overflow-x-auto rounded bg-[var(--color-kp-surface)] p-2 text-xs">{step.command}</pre>}
                  {step.caution && <div className="text-xs text-[var(--color-kp-muted)]">{step.caution}</div>}
                </li>
              ))}
            </ol>
          </section>

          <section>
            <h2 className="mb-2 text-sm font-semibold uppercase text-[var(--color-kp-muted)]">Guarded automation</h2>
            <div className="space-y-2">
              {data.hooks.map((hook) => (
                <div
                  key={hook.id}
                  className="flex flex-wrap items-center justify-between gap-2 rounded border border-[var(--color-kp-border)] px-3 py-2"
                >
                  <div>
                    <div className="font-medium">{hook.title}</div>
                    <div className="text-xs text-[var(--color-kp-muted)]">{hook.description}</div>
                    <div className="text-xs">Risk: {hook.risk}</div>
                  </div>
                  <KPButton
                    size="sm"
                    variant="secondary"
                    disabled={!hook.dryRunSupported || dryRunMutation.isPending}
                    onClick={() => dryRunMutation.mutate(hook.id)}
                  >
                    Dry-run
                  </KPButton>
                </div>
              ))}
            </div>
          </section>

          {dryRun && (
            <section className="rounded border border-emerald-500/30 bg-emerald-500/10 p-3 text-sm">
              <div className="font-medium">Dry-run: {dryRun.hookId}</div>
              <p>{dryRun.message}</p>
              <ul className="mt-2 list-disc pl-5 font-mono text-xs">
                {dryRun.plannedActions.map((a) => (
                  <li key={a}>{a}</li>
                ))}
              </ul>
              <p className="mt-1 text-xs text-[var(--color-kp-muted)]">Audit: {dryRun.auditId}</p>
            </section>
          )}
        </>
      )}
    </div>
  )
}
