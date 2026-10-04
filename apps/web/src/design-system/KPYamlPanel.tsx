export function KPYamlPanel({ yaml }: { yaml: string }) {
  return (
    <pre className="max-h-80 overflow-auto rounded-md border border-[var(--color-kp-border)] bg-[#0a0d12] p-3 text-xs leading-relaxed text-[var(--color-kp-text)]">
      {yaml}
    </pre>
  )
}
