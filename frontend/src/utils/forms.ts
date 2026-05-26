export function buildFormData(formData: FormData, data: any, prefix: string = '') {
  if (data === null || data === undefined) return

  const isFile = (v: any) => v instanceof File || v instanceof Blob

  if (isFile(data)) {
    if (prefix.trim() !== '') formData.append(prefix, data)

    return
  }

  if (typeof data === 'object' && !isFile(data)) {
    if (Array.isArray(data)) {
      data.forEach((value, i) => {
        const key = `${prefix}[${i}]`
        buildFormData(formData, value, key)
      })
    } else {
      Object.keys(data).forEach((key) => {
        const value = data[key]
        const formKey = prefix ? `${prefix}[${key}]` : key
        buildFormData(formData, value, formKey)
      })
    }

    return
  }

  formData.append(prefix, String(data))
}
