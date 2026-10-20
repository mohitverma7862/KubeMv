export function DescribeViewer({ text }: { text?: string }) {
  return <pre data-testid="fast-describe">{text || 'No describe loaded.'}</pre>
}
