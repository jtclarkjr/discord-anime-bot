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
 * Format a date to user-friendly format
 */
export function formatDate(date: Date): string {
  const dateOptions: Intl.DateTimeFormatOptions = {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  }
  return date.toLocaleDateString('en-US', dateOptions)
}

/**
 * Format time to user-friendly format (12-hour)
 */
export function formatTime(date: Date): string {
  const timeOptions: Intl.DateTimeFormatOptions = {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true
  }
  return date.toLocaleTimeString('en-US', timeOptions)
}

/**
 * Format date and time for air date display
 */
export function formatAirDate(date: Date): string {
  return `${formatDate(date)} at ${formatTime(date)}`
}

/**
 * Format compact date and time for lists
 */
export function formatCompactDateTime(date: Date): string {
  const timeOptions: Intl.DateTimeFormatOptions = {
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
    hour12: true
  }
  return date.toLocaleDateString('en-US', timeOptions)
}
