// middleware/admin.ts
export default defineNuxtRouteMiddleware((to) => {
  if (!to.path.startsWith('/admin')) return
  if (to.path === '/admin/login') {
    // 允许在登录页录入访问口令
    return
  }

  // 访问口令缺失，强制回登录页
  const accessKey = useCookie<string | null>('admin-access')
  if (!accessKey.value) {
    return navigateTo('/admin/login')
  }

  const adminAuth = useAdminAuth()
  if (adminAuth.isAuthenticated.value) {
    return
  }

  const auth = useAuth()
  if (auth.isAuthenticated.value && auth.user.value && !auth.user.value.isAdmin) {
    return navigateTo('/')
  }

  return navigateTo('/admin/login')
})
