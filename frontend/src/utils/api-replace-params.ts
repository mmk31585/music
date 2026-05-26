type ParamValue = string | number | boolean | null | undefined

export function apiReplaceParams(url: string, params: Record<string, ParamValue> = {}): string {
  // Replace placeholders: :name, :name?, {name}, {name?}
  const segmentPattern = /(\/?)([:{])([a-zA-Z0-9_]+)(\?)?(}?)/g

  let result = url.replace(segmentPattern, (match, leadingSlash, opener, key, optionalMark) => {
    const isOptional = !!optionalMark
    const val = params[key]

    if (val === undefined || val === null || val === '') {
      if (isOptional) {
        // remove the whole leading slash + placeholder to avoid leftover slashes
        return leadingSlash ? '' : ''
      }
      throw new Error(`Missing required URL parameter: ${key}`)
    }

    const encoded = encodeURIComponent(String(val))
    // preserve a leading slash if it existed in the template
    return (leadingSlash ?? '') + encoded
  })

  // Collapse multiple slashes introduced by removals
  result = result.replace(/\/{2,}/g, '/')

  // Remove trailing slash if original didn't intend it (optional)
  // If template had trailing slash, keep it. Otherwise, remove if now trailing.
  if (!url.endsWith('/') && result.endsWith('/')) result = result.slice(0, -1)

  // Append leftover params as query string
  // First, collect all placeholder keys from the template
  const placeholderKeys = Array.from(url.matchAll(/[:{]([a-zA-Z0-9_]+)\??}?/g)).map((m) => m[1])
  const queryEntries = Object.entries(params).filter(([k, _]) => !placeholderKeys.includes(k))

  if (queryEntries.length) {
    const qs = queryEntries
      .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(String(v ?? ''))}`)
      .join('&')
    result += (result.includes('?') ? '&' : '?') + qs
  }

  return result
}
