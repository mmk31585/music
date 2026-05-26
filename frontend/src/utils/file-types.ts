const IMAGE_EXT = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg', 'avif']
const VIDEO_EXT = ['mp4', 'webm', 'ogg', 'mov', 'mkv']
const SOUND_EXT = ['mp3', 'wav', 'aac', 'flac', 'ogg', 'm4a']
const PDF_EXT = ['pdf']
const WORD_EXT = ['doc', 'docx', 'rtf']
const EXCEL_EXT = ['xls', 'xlsx', 'csv', 'ods']
const POWER_POINT_EXT = ['ppt', 'pptx', 'odp']
const TEXT_EXT = ['txt', 'md', 'log', 'json', 'xml', 'yml', 'yaml']
const ZIP_EXT = ['zip', 'rar', '7z', 'tar', 'gz', 'bz2']

export type ResolvedMediaType =
  | 'image'
  | 'video'
  | 'sound'
  | 'pdf'
  | 'excel'
  | 'power-point'
  | 'word'
  | 'text'
  | 'zip'
  | 'other'

export function resolveMediaType(ext?: string): ResolvedMediaType {
  if (!ext) return 'other'

  const e = ext.toLowerCase().replace('.', '')

  if (IMAGE_EXT.includes(e)) return 'image'
  if (VIDEO_EXT.includes(e)) return 'video'
  if (SOUND_EXT.includes(e)) return 'sound'
  if (PDF_EXT.includes(e)) return 'pdf'
  if (WORD_EXT.includes(e)) return 'word'
  if (EXCEL_EXT.includes(e)) return 'excel'
  if (POWER_POINT_EXT.includes(e)) return 'power-point'
  if (TEXT_EXT.includes(e)) return 'text'
  if (ZIP_EXT.includes(e)) return 'zip'

  return 'other'
}
