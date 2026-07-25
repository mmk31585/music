export function buildFormData(
  formData: FormData,
  data: any,
  prefix: string = ''
) {
  if (data === null || data === undefined) return

  const isFile = (v: any) => v instanceof File || v instanceof Blob

  if (isFile(data)) {
    if (prefix.trim() !== '') formData.append(prefix, data as File | Blob)
    return
  }

  if (typeof data === 'object' && !isFile(data)) {
    if (Array.isArray(data)) {
      for (const [i, value] of data.entries()) {
        const key = `${prefix}[${i}]`
        buildFormData(formData, value, key)
      }
    } else if (data !== null) {
      for (const key of Object.keys(data)) {
        const value = (data as Record<string, any>)[key]
        const formKey = prefix ? `${prefix}[${key}]` : key
        buildFormData(formData, value, formKey)
      }
    }
    return
  }

  formData.append(prefix, String(data))
}
