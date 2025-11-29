export const useUpload = () => {
  const api = useApi()
  const { resolveAssetUrl } = useAssetUrl()

  const uploadAvatar = async (file: File) => {
    const form = new FormData()
    form.append('file', file)
    const res = await api.post<{ data: { url: string } }>('/uploads/avatar', form)
    return resolveAssetUrl(res.data?.url)
  }

  const uploadPostImage = async (file: File) => {
    const form = new FormData()
    form.append('file', file)
    const res = await api.post<{ data: { url: string } }>('/uploads/posts', form)
    // 返回后端相对路径，展示时再用 resolveAssetUrl 统一拼接正确域名
    return res.data?.url
  }

  const uploadRoleAvatar = async (file: File) => {
    const form = new FormData()
    form.append('file', file)
    const tryUpload = async (path: string) => api.post<{ data: { url: string } }>(path, form)
    try {
      const res = await tryUpload('/uploads/roles')
      return res.data?.url
    } catch (error: any) {
      const status = error?.statusCode ?? error?.response?.status ?? error?.status
      // 旧版本后端没有 /uploads/roles 时回退到 avatar 上传接口
      if (status === 404 || status === 405) {
        const res = await tryUpload('/uploads/avatar')
        return res.data?.url
      }
      throw error
    }
  }

  return {
    uploadAvatar,
    uploadPostImage,
    uploadRoleAvatar,
    resolveAssetUrl,
  }
}
