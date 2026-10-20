export function EventViewer({ text }: { text?: string }) {
  return <pre data-testid="fast-events">{text || 'No events loaded.'}</pre>
}
