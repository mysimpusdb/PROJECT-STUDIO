export const formatDate = (dateString: string | null | undefined): string => {
  if (!dateString) return ''
  const date = new Date(dateString)
  if (isNaN(date.getTime())) return ''
  
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  
  return `${y}-${m}-${d}`
}

export const getDaysBetween = (start: string, end: string): number => {
  const d1 = new Date(start)
  const d2 = new Date(end)
  
  if (isNaN(d1.getTime()) || isNaN(d2.getTime())) return 0
  
  const diffTime = Math.abs(d2.getTime() - d1.getTime())
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
}
