export function LogViewer({ text }: { text?: string }) {
  return <pre data-testid="fast-logs">{text || 'No logs loaded.'}</pre>
}
