import { describe, it, expect } from 'vitest'
import { formatDuration } from '@/utils/format'

describe('formatDuration', () => {
  it('formats 0 as "0:00"', () => {
    expect(formatDuration(0)).toBe('0:00')
  })

  it('formats 65 as "1:05"', () => {
    expect(formatDuration(65)).toBe('1:05')
  })

  it('formats 3661 as "61:01"', () => {
    expect(formatDuration(3661)).toBe('61:01')
  })

  it('handles null', () => {
    expect(formatDuration(null)).toBe('0:00')
  })

  it('handles undefined', () => {
    expect(formatDuration(undefined)).toBe('0:00')
  })

  it('handles negative values', () => {
    expect(formatDuration(-1)).toBe('0:00')
  })
})
