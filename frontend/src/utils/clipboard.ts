export async function copyText(text: string): Promise<boolean> {
  if (!text) return false
  if (typeof document === 'undefined') return false

  try {
    if (navigator?.clipboard?.writeText) {
      await navigator.clipboard.writeText(text)
      return true
    }
  } catch {
    // fallback below
  }

  const ta = document.createElement('textarea')
  ta.value = text
  ta.style.position = 'fixed'
  ta.style.top = '-9999px'
  ta.setAttribute('readonly', 'true')
  document.body.appendChild(ta)

  const activeElement = document.activeElement as HTMLElement | null
  ta.focus({ preventScroll: true })
  ta.select()

  try {
    return document.execCommand('copy')
  } finally {
    document.body.removeChild(ta)
    activeElement?.focus?.({ preventScroll: true })
  }
}
