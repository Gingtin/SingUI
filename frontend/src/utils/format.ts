import dayjs from 'dayjs'

export function formatBytes(bytes: number, decimals = 2): string {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

export function formatTime(timestamp: number | string | Date): string {
  if (!timestamp || timestamp === 0) return '永久有效'
  // If in milliseconds
  if (typeof timestamp === 'number' && timestamp > 100000000000) {
    return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
  }
  // If in seconds
  if (typeof timestamp === 'number') {
    return dayjs(timestamp * 1000).format('YYYY-MM-DD HH:mm:ss')
  }
  return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
}

export function formatDuration(seconds: number): string {
  if (!seconds || seconds <= 0) return '0秒'
  const d = Math.floor(seconds / (3600 * 24))
  const h = Math.floor((seconds % (3600 * 24)) / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = Math.floor(seconds % 60)

  const parts = []
  if (d > 0) parts.push(`${d}天`)
  if (h > 0) parts.push(`${h}小时`)
  if (m > 0) parts.push(`${m}分`)
  if (s > 0 && d === 0) parts.push(`${s}秒`)
  return parts.join(' ') || '0秒'
}
