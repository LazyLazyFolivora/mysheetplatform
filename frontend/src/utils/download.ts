/**
 * 通过鉴权接口拉取文件字节流并触发本地下载。
 * 不走 axios JSON 拦截器，避免把二进制当 JSON 处理。
 */
export async function downloadAuthBlob(apiPath: string, fallbackName: string): Promise<void> {
  const baseUrl = (import.meta.env.VITE_API_BASE_URL || '/api').replace(/\/$/, '')
  const path = apiPath.startsWith('/') ? apiPath : `/${apiPath}`
  const url = `${baseUrl}${path}`
  const token = localStorage.getItem('token') || ''

  const res = await fetch(url, {
    method: 'GET',
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    credentials: 'include',
  })

  const contentType = res.headers.get('Content-Type') || ''
  if (!res.ok) {
    let message = `下载失败（${res.status}）`
    if (contentType.includes('application/json')) {
      try {
        const data = await res.json()
        if (data?.message) message = data.message
      } catch {
        // ignore
      }
    }
    throw new Error(message)
  }

  const blob = await res.blob()
  if (!blob || blob.size === 0) {
    throw new Error('下载失败：文件为空')
  }

  const filename = parseFilename(res.headers.get('Content-Disposition')) || fallbackName
  const objectUrl = URL.createObjectURL(blob)
  try {
    const a = document.createElement('a')
    a.href = objectUrl
    a.download = filename
    a.rel = 'noopener'
    document.body.appendChild(a)
    a.click()
    a.remove()
  } finally {
    URL.revokeObjectURL(objectUrl)
  }
}

function parseFilename(contentDisposition: string | null): string | null {
  if (!contentDisposition) return null
  // filename*=UTF-8''encoded
  const star = /filename\*\s*=\s*UTF-8''([^;]+)/i.exec(contentDisposition)
  if (star?.[1]) {
    try {
      return decodeURIComponent(star[1].trim().replace(/^"|"$/g, ''))
    } catch {
      return star[1].trim().replace(/^"|"$/g, '')
    }
  }
  const plain = /filename\s*=\s*("?)([^";]+)\1/i.exec(contentDisposition)
  if (plain?.[2]) return plain[2].trim()
  return null
}
