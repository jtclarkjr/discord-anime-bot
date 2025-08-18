/**
 * Format countdown time from seconds to a human readable string
 */
export function formatCountdown(timeUntil: number): string {
  const days = Math.floor(timeUntil / (24 * 60 * 60))
  const hours = Math.floor((timeUntil % (24 * 60 * 60)) / (60 * 60))
  const minutes = Math.floor((timeUntil % (60 * 60)) / 60)

  let timeString = ''
  if (days > 0) timeString += `${days} day${days > 1 ? 's' : ''} `
  if (hours > 0) timeString += `${hours} hour${hours > 1 ? 's' : ''} `
  if (minutes > 0) timeString += `${minutes} minute${minutes > 1 ? 's' : ''}`

  return timeString.trim()
}

/**
 * Format a date to user-friendly format (Discord timestamp)
 * Format: <t:timestamp:D> - Shows date only (e.g., "December 25, 2023")
 */
export function formatDate(date: Date): string {
  return `<t:${Math.floor(date.getTime() / 1000)}:D>`
}

/**
 * Format time to user-friendly format (Discord timestamp)
 * Format: <t:timestamp:t> - Shows time only (e.g., "3:30 PM")
 */
export function formatTime(date: Date): string {
  return `<t:${Math.floor(date.getTime() / 1000)}:t>`
}

/**
 * Format date and time for air date display (Discord timestamp)
 * Format: <t:timestamp:F> - Full date and time (e.g., "Monday, December 25, 2023 3:30 PM")
 */
export function formatAirDate(date: Date): string {
  return `<t:${Math.floor(date.getTime() / 1000)}:F>`
}

/**
 * Format compact date and time for lists (Discord timestamp)
 * Format: <t:timestamp:f> - Short date and time (e.g., "December 25, 2023 3:30 PM")
 */
export function formatCompactDateTime(date: Date): string {
  return `<t:${Math.floor(date.getTime() / 1000)}:f>`
}

/**
 * Format a time as Discord relative timestamp
 * Format: <t:timestamp:R> - Relative time (e.g., "in 2 hours", "3 days ago")
 */
export function formatRelativeTimestamp(date: Date): string {
  return `<t:${Math.floor(date.getTime() / 1000)}:R>`
}
