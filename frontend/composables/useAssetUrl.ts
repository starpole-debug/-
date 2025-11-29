export const useAssetUrl = () => {
  const config = useRuntimeConfig()

  const resolveAssetUrl = (url?: string | null) => {
    if (!url) return ''
    const trimmed = url.trim()
    if (/^https?:\/\//i.test(trimmed)) return trimmed
    const apiBase = config.public.apiBase || ''
    const assetBase = config.public.assetBase || ''
    const pickOrigin = () => {
      // assetBase 优先（可在环境变量里显式配置，如 http://localhost:8080）
      const candidate = assetBase || apiBase
      if (candidate) {
        try {
          const parsed = new URL(candidate, typeof window !== 'undefined' ? window.location.origin : 'http://localhost')
          return parsed.origin
        } catch {
          // ignore
        }
      }
      if (typeof window !== 'undefined') return window.location.origin
      return ''
    }
    const origin = pickOrigin()
    if (!origin) return trimmed
    return origin.replace(/\/+$/, '') + (trimmed.startsWith('/') ? trimmed : `/${trimmed}`)
  }

  return { resolveAssetUrl }
}
