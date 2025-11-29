export default defineNuxtRouteMiddleware(async (to) => {
  const auth = useAuth()
  if (!auth.isAuthenticated.value && auth.hasToken.value) {
    await auth.fetchProfile()
  }
  if (auth.isAuthenticated.value) {
    return
  }
  const redirectTo =
    typeof to.fullPath === 'string' && to.fullPath.length > 1 ? to.fullPath : '/'
  return navigateTo({
    path: '/login',
    query: { redirect: redirectTo },
  })
})
