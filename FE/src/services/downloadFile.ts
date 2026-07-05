import { apiClient } from './api.client'

export async function downloadReportExport(
  path: string,
  params: Record<string, string | undefined>,
  filename: string
): Promise<void> {
  const response = await apiClient.get(path, { params, responseType: 'blob' })
  const url = URL.createObjectURL(response.data as Blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}
