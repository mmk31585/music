import { describe, it, expect } from 'vitest'

function formatDate(dateStr: string): string {
  try {
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return dateStr
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    })
  } catch {
    return dateStr
  }
}

describe('formatDate', () => {
  it('formats ISO date string correctly', () => {
    expect(formatDate('2024-06-15T10:00:00Z')).toBe('Jun 15, 2024')
  })

  it('returns original string on invalid input', () => {
    expect(formatDate('not-a-date')).toBe('not-a-date')
  })

  it('handles empty string', () => {
    expect(formatDate('')).toBe('')
  })
})
