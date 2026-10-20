export function YamlViewer({ text }: { text?: string }) {
  return <pre data-testid="fast-yaml">{text || 'No YAML loaded.'}</pre>
}
