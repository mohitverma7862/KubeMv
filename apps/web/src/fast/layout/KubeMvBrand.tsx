/** Product mark from docs/assets/kubemv-fast-mode-reference.svg (hex + KubeMv, #4cc9ff). */
export function KubeMvBrand({ className }: { className?: string }) {
  return (
    <span className={className ?? 'km-fast-brand'} data-testid="fast-brand">
      <svg className="km-fast-logo" width="20" height="20" viewBox="0 0 24 24" aria-hidden>
        <path
          d="M12 2 20.66 7v10L12 22 3.34 17V7L12 2z"
          fill="none"
          stroke="currentColor"
          strokeWidth="1.6"
          strokeLinejoin="round"
        />
      </svg>
      <span>KubeMv</span>
    </span>
  )
}
