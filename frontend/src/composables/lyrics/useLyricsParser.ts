export interface WordPart {
  text: string
  timeSeconds: number
}

export interface ParsedLine {
  timeSeconds: number
  text: string
  words: WordPart[]
  isActive: boolean
}

export function parseLRCLines(content: string): ParsedLine[] {
  const lines: ParsedLine[] = []
  const lineRegex = /\[(\d{2}):(\d{2})[\.:](\d{2,3})\](.*)/
  for (const raw of content.split('\n')) {
    const match = raw.trim().match(lineRegex)
    if (match) {
      const min = parseInt(match[1]!, 10)
      const sec = parseInt(match[2]!, 10)
      const ms = parseInt(match[3]!, 10)
      const timeSeconds = min * 60 + sec + ms / (match[3]!.length === 3 ? 1000 : 100)
      const textPart = match[4]!.trim()
      const words = parseWordTimings(textPart)
      lines.push({ timeSeconds, text: textPart, words, isActive: false })
    }
  }
  return lines.sort((a, b) => a.timeSeconds - b.timeSeconds)
}

export function parsePlainLines(content: string): ParsedLine[] {
  return content
    .split('\n')
    .filter(Boolean)
    .map((text) => ({
      timeSeconds: 0,
      text,
      words: [],
      isActive: false,
    }))
}

function parseWordTimings(text: string): WordPart[] {
  const wordRegex = /<(\d{2}:\d{2}[\.:]\d{2,3})>([^<]+)/g
  const words: WordPart[] = []
  let match
  while ((match = wordRegex.exec(text)) !== null) {
    const parts = match[1]!.split(/[:.]/)
    const min = parseInt(parts[0]!, 10)
    const sec = parseInt(parts[1]!, 10)
    const ms = parseInt(parts[2]!, 10)
    const timeSeconds = min * 60 + sec + ms / (parts[2]!.length === 3 ? 1000 : 100)
    words.push({ text: match[2]!.trim(), timeSeconds })
  }
  if (words.length === 0) {
    return text
      .split(/\s+/)
      .filter(Boolean)
      .map((word) => ({
        text: word,
        timeSeconds: -1,
      }))
  }
  return words
}
