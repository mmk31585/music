/**
 * Triggers a download of arbitrary data.
 * @param filename Desired file name (including extension).
 * @param data     Blob‑compatible data (string, ArrayBuffer, etc.).
 * @param contentType MIME type, e.g. "application/json".
 */
export function downloadDataAsFile(filename: string, data: BlobPart, contentType: string): void {
  const blob = new Blob([data], { type: contentType })

  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = filename
  link.click()

  // Firefox needs a short delay before revoking the object URL.
  setTimeout(() => URL.revokeObjectURL(link.href), 100)
}
