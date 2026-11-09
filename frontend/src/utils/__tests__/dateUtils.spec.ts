import { describe, it, expect } from 'vitest'
import { formatDate, getDaysBetween } from '../dateUtils'

describe('dateUtils', () => {
  it('formatDate formats date correctly', () => {
    expect(formatDate('2026-09-21T07:56:42+07:00')).toBe('2026-09-21')
    expect(formatDate(null)).toBe('')
    expect(formatDate('')).toBe('')
    expect(formatDate('invalid-date')).toBe('')
  })

  it('getDaysBetween calculates difference correctly', () => {
    expect(getDaysBetween('2026-09-01', '2026-09-10')).toBe(9)
    expect(getDaysBetween('2026-09-10', '2026-09-01')).toBe(9)
    expect(getDaysBetween('invalid', '2026-09-01')).toBe(0)
  })
})
