export function timeSince(timestamp: number): string {
  const intervals = [
    { label: 'year', seconds: 31536000 },
    { label: 'month', seconds: 2592000 },
    { label: 'day', seconds: 86400 },
    { label: 'hour', seconds: 3600 },
    { label: 'minute', seconds: 60 },
    { label: 'second', seconds: 1 }
  ]

  const date = new Date(timestamp)
  const secondsSince = Math.floor((Date.now() - date.getTime()) / 1000)
  if (secondsSince === 0) {
    return 'Just now'
  }
  if (secondsSince < 0) {
    return 'In the future 😱'
  }

  const interval = intervals.find(i => i.seconds < secondsSince) ?? intervals[0]!
  const count = Math.floor(secondsSince / interval.seconds)
  return `${count} ${interval.label}${count !== 1 ? 's' : ''} ago`
}

export function timeBetween(startTime: number, endTime: number): string {
  const intervals = [
    { label: 'year', seconds: 31536000 },
    { label: 'month', seconds: 2592000 },
    { label: 'day', seconds: 86400 },
    { label: 'hour', seconds: 3600 },
    { label: 'minute', seconds: 60 },
    { label: 'second', seconds: 1 }
  ]

  const seconds = Math.floor((endTime - startTime) / 1000)
  const interval = intervals.find(i => i.seconds < seconds) ?? intervals[0]!
  const count = Math.floor(seconds / interval.seconds)
  return `${count} ${interval.label}${count !== 1 ? 's' : ''}`
}

export function timeDurationToReadable(durationInSeconds: number): string {
  if (durationInSeconds <= 0) {
    return durationInSeconds.toFixed(0) + ' seconds'
  }

  const intervals = [
    { label: 'year', seconds: 31536000 },
    { label: 'month', seconds: 2592000 },
    { label: 'day', seconds: 86400 },
    { label: 'hour', seconds: 3600 },
    { label: 'minute', seconds: 60 },
    { label: 'second', seconds: 1 }
  ]

  const interval = intervals.find(i => i.seconds < durationInSeconds) ?? intervals[0]!
  const count = Math.floor(durationInSeconds / interval.seconds)
  return `${count} ${interval.label}${count !== 1 ? 's' : ''}`
}

export function bytesToReadable(bytes: number): string {
  let displayVal = bytes
  let unitIndex = 0
  const units = ['B', 'KB', 'MB', 'GB', 'TB']

  while (displayVal > 1000 && unitIndex - 1 < units.length) {
    displayVal /= 1000
    unitIndex++
  }

  return `${displayVal.toFixed(0)}${units[unitIndex]}`
}
