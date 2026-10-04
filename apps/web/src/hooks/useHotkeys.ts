import { useEffect } from 'react'

export function useHotkeys(bindings: Record<string, (e: KeyboardEvent) => void>) {
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const key = [
        e.ctrlKey || e.metaKey ? 'Ctrl' : '',
        e.shiftKey ? 'Shift' : '',
        e.key.length === 1 ? e.key.toUpperCase() : e.key,
      ]
        .filter(Boolean)
        .join('+')
      const fn = bindings[key]
      if (fn) {
        e.preventDefault()
        fn(e)
      }
    }
    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [bindings])
}
