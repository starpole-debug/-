<template>
  <nav class="flex items-center gap-1">
    <NuxtLink to="/" class="mr-4 font-bold text-white tracking-tight">Nebula</NuxtLink>
    
    <div class="h-4 w-[1px] bg-white/20 mx-2"></div>

    <NuxtLink to="/search" class="px-3 py-1.5 rounded-lg text-sm text-slate-300 hover:text-white hover:bg-white/10 transition-all">探索</NuxtLink>
    <NuxtLink to="/community" class="px-3 py-1.5 rounded-lg text-sm text-slate-300 hover:text-white hover:bg-white/10 transition-all">社区</NuxtLink>
    <NuxtLink to="/creator" class="px-3 py-1.5 rounded-lg text-sm text-slate-300 hover:text-white hover:bg-white/10 transition-all">创作</NuxtLink>
    <NuxtLink to="/store" class="px-3 py-1.5 rounded-lg text-sm text-slate-300 hover:text-white hover:bg-white/10 transition-all">商店</NuxtLink>
    
    <div class="h-4 w-[1px] bg-white/20 mx-2"></div>

    <NuxtLink to="/notifications" class="p-1.5 rounded-full text-slate-300 hover:text-white hover:bg-white/10 transition-all">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
        <path d="M10 2a6 6 0 00-6 6v3.586l-.707.707A1 1 0 004 14h12a1 1 0 00.707-1.707L16 11.586V8a6 6 0 00-6-6zM10 18a3 3 0 01-3-3h6a3 3 0 01-3 3z" />
      </svg>
    </NuxtLink>

    <NuxtLink
      to="/me"
      class="ml-3 inline-flex items-center gap-2 rounded-full bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 px-4 py-1.5 text-sm font-semibold text-white shadow-lg shadow-indigo-900/40 transition-all hover:translate-y-0.5 hover:shadow-xl"
    >
      <span class="inline-flex h-6 w-6 items-center justify-center rounded-full bg-white/20 backdrop-blur">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-4 w-4">
          <path
            fill-rule="evenodd"
            d="M10 2a4 4 0 00-4 4v1a4 4 0 004 4 4 4 0 004-4V6a4 4 0 00-4-4zm-6 14a6 6 0 1112 0v1a1 1 0 01-1 1H5a1 1 0 01-1-1v-1z"
            clip-rule="evenodd"
          />
        </svg>
      </span>
      我的主页
    </NuxtLink>

    <button
      v-if="auth.isAuthenticated || auth.hasToken"
      type="button"
      class="ml-2 inline-flex items-center gap-2 rounded-full border border-white/30 px-3 py-1.5 text-sm text-white/90 transition-all hover:border-white hover:text-white hover:bg-white/10"
      @click="handleLogout"
    >
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="h-4 w-4">
        <path
          fill-rule="evenodd"
          d="M7.5 4.75A.75.75 0 0 1 8.25 4h6.5a.75.75 0 0 1 .75.75v10.5a.75.75 0 0 1-.75.75h-6.5a.75.75 0 0 1-.75-.75v-2a.75.75 0 0 1 1.5 0V15h5v-10h-5v1.75a.75.75 0 0 1-1.5 0z"
          clip-rule="evenodd"
        />
        <path d="M3.72 10.53a.75.75 0 0 1 0-1.06l2-2a.75.75 0 1 1 1.06 1.06L5.81 9.5H12a.75.75 0 0 1 0 1.5H5.81l.97.97a.75.75 0 1 1-1.06 1.06z" />
      </svg>
      退出登录
    </button>
  </nav>
</template>

<script setup lang="ts">
const router = useRouter()
const auth = useAuth()

const handleLogout = () => {
  auth.logout()
  router.push('/login')
}
</script>
