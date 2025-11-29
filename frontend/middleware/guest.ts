export default defineNuxtRouteMiddleware(async (to) => {
  const auth = useAuth()
  if (!auth.isAuthenticated.value && auth.hasToken.value) {
    await auth.fetchProfile()
  }
  if (!auth.isAuthenticated.value) {
    return
  }
  const fallbackRedirect =
    typeof to.query?.redirect === 'string' ? (to.query.redirect as string) : '/'
  return navigateTo(fallbackRedirect)
})
