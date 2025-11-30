<template>
  <nav class="relative flex w-full items-center justify-between">
    <div class="flex items-center gap-3">
      <NuxtLink to="/" class="font-display font-bold text-charcoal-900 tracking-tight text-xl flex items-center gap-2">
        <div class="w-8 h-8 rounded-lg bg-accent-yellow flex items-center justify-center text-charcoal-900 text-xs font-bold shadow-sm">
          M
        </div>
        <span class="hidden sm:block">MilkyVerse</span>
      </NuxtLink>
    </div>

    <!-- Desktop links -->
    <div class="hidden items-center gap-1 md:flex">
      <NuxtLink 
        v-for="link in links" 
        :key="link.to" 
        :to="link.to" 
        class="px-4 py-2 rounded-full text-sm font-medium text-charcoal-500 hover:text-charcoal-900 hover:bg-bg-cream-200 transition-all duration-300"
        active-class="text-charcoal-900 bg-bg-cream-200 font-bold"
      >
        {{ link.label }}
      </NuxtLink>

      <div class="h-6 w-[1px] bg-charcoal-300 mx-3"></div>

      <NuxtLink to="/notifications" class="p-2 rounded-full text-charcoal-500 hover:text-charcoal-900 hover:bg-bg-cream-200 transition-all duration-300 relative group">
        <span class="sr-only">Notifications</span>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-5 h-5 group-hover:scale-110 transition-transform">
          <path fill-rule="evenodd" d="M5.25 9a6.75 6.75 0 0113.5 0v.75c0 2.123.8 4.057 2.118 5.52a.75.75 0 01-.297 1.206c-1.544.57-3.16.99-4.831 1.243a3.75 3.75 0 11-7.48 0 24.585 24.585 0 01-4.831-1.244.75.75 0 01-.298-1.205A8.217 8.217 0 005.25 9.75V9zm4.502 8.9a2.25 2.25 0 104.496 0 25.057 25.057 0 01-4.496 0z" clip-rule="evenodd" />
        </svg>
      </NuxtLink>

      <NuxtLink
        to="/me"
        class="ml-2 inline-flex items-center gap-2 rounded-full bg-charcoal-900 pl-1 pr-4 py-1 text-sm font-medium text-bg-cream-100 transition-all hover:bg-[#25262b] hover:shadow-lg group"
      >
        <span class="inline-flex h-7 w-7 items-center justify-center rounded-full bg-accent-yellow text-charcoal-900 text-xs font-bold group-hover:scale-105 transition-transform">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-4 w-4">
            <path fill-rule="evenodd" d="M10 2a4 4 0 00-4 4v1a4 4 0 004 4 4 4 0 004-4V6a4 4 0 00-4-4zm-6 14a6 6 0 1112 0v1a1 1 0 01-1 1H5a1 1 0 01-1-1v-1z" clip-rule="evenodd" />
          </svg>
        </span>
        我的主页
      </NuxtLink>
    </div>

    <!-- Mobile actions -->
    <div class="flex items-center gap-3 md:hidden">
      <NuxtLink to="/notifications" class="p-2 rounded-full text-charcoal-500 hover:text-charcoal-900 transition-colors">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-6 h-6">
          <path fill-rule="evenodd" d="M5.25 9a6.75 6.75 0 0113.5 0v.75c0 2.123.8 4.057 2.118 5.52a.75.75 0 01-.297 1.206c-1.544.57-3.16.99-4.831 1.243a3.75 3.75 0 11-7.48 0 24.585 24.585 0 01-4.831-1.244.75.75 0 01-.298-1.205A8.217 8.217 0 005.25 9.75V9zm4.502 8.9a2.25 2.25 0 104.496 0 25.057 25.057 0 01-4.496 0z" clip-rule="evenodd" />
        </svg>
      </NuxtLink>
      
      <button
        type="button"
        class="relative z-50 p-2 text-charcoal-900 transition-colors"
        @click="menuOpen = !menuOpen"
      >
        <div class="w-6 h-5 flex flex-col justify-between">
          <span class="w-full h-0.5 bg-current rounded-full transition-all duration-300 origin-left" :class="{ 'rotate-45 translate-x-0.5': menuOpen }"></span>
          <span class="w-full h-0.5 bg-current rounded-full transition-all duration-300" :class="{ 'opacity-0': menuOpen }"></span>
          <span class="w-full h-0.5 bg-current rounded-full transition-all duration-300 origin-left" :class="{ '-rotate-45 translate-x-0.5': menuOpen }"></span>
        </div>
      </button>
    </div>

    <!-- Mobile Full Screen Menu -->
    <Teleport to="body">
      <Transition
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="opacity-0 translate-y-4"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition duration-200 ease-in"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 translate-y-4"
      >
        <div v-if="menuOpen" class="fixed inset-0 z-40 bg-bg-cream-100/95 backdrop-blur-xl flex flex-col items-center justify-center p-6">
          <div class="w-full max-w-sm space-y-6 text-center">
            <div class="space-y-2">
              <NuxtLink
                v-for="(link, index) in links"
                :key="link.to"
                :to="link.to"
                class="block text-3xl font-display font-bold text-charcoal-500 hover:text-charcoal-900 transition-colors py-2"
                active-class="text-charcoal-900 scale-110"
                @click="closeMenu"
                :style="{ animation: `fadeIn 0.5s ease-out ${index * 0.1}s backwards` }"
              >
                {{ link.label }}
              </NuxtLink>
            </div>

            <div class="h-px w-20 bg-charcoal-300 mx-auto my-8"></div>

            <div class="space-y-4">
              <NuxtLink
                to="/me"
                class="block w-full py-4 rounded-2xl bg-charcoal-900 text-lg font-medium text-bg-cream-100 hover:scale-105 transition-all"
                @click="closeMenu"
              >
                我的主页
              </NuxtLink>

              <button
                v-if="auth.isAuthenticated || auth.hasToken"
                type="button"
                class="block w-full py-4 rounded-2xl border border-charcoal-300 text-lg font-medium text-charcoal-500 hover:text-charcoal-900 hover:bg-bg-cream-200 transition-all"
                @click="handleLogout"
              >
                退出登录
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </nav>
</template>

<script setup lang="ts">
const router = useRouter()
const route = useRoute()
const auth = useAuth()
const menuOpen = ref(false)

const links = [
  { to: '/search', label: '探索' },
  { to: '/community', label: '社区' },
  { to: '/creator', label: '创作' },
  { to: '/store', label: '商店' },
]

const closeMenu = () => {
  menuOpen.value = false
}

watch(
  () => route.fullPath,
  () => closeMenu(),
)

// Lock body scroll when menu is open
watch(menuOpen, (isOpen) => {
  if (process.client) {
    document.body.style.overflow = isOpen ? 'hidden' : ''
  }
})

const handleLogout = () => {
  auth.logout()
  router.push('/login')
  closeMenu()
}
</script>
