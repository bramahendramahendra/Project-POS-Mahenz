export function formatTime(time: string): string {
  const [hour, minute] = time.split(':')
  return `${hour.padStart(2, '0')}:${minute.padStart(2, '0')}`
}

export function formatShiftTime(start_time: string, end_time: string): string {
  return `${formatTime(start_time)} – ${formatTime(end_time)}`
}
