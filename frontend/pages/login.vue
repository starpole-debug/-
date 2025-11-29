<template>
  <div class="flex min-h-[80vh] w-full overflow-hidden rounded-3xl border border-white/5 bg-[#151621] shadow-2xl">
    <!-- Left Side: Visual Brand -->
    <div class="relative hidden w-1/2 overflow-hidden lg:block">
      <div class="absolute inset-0 bg-gradient-to-br from-indigo-600 to-purple-900 opacity-40"></div>
      <div class="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-20 mix-blend-overlay"></div>
      
      <!-- Abstract Shapes -->
      <div class="absolute -left-20 -top-20 h-96 w-96 rounded-full bg-indigo-500/30 blur-[100px] animate-float"></div>
      <div class="absolute bottom-0 right-0 h-80 w-80 rounded-full bg-purple-500/30 blur-[80px] animate-float-delayed"></div>

      <div class="relative z-10 flex h-full flex-col justify-between p-12">
        <div class="text-2xl font-bold text-white tracking-tight">Nebula Studio</div>
        <div class="space-y-4">
          <h2 class="text-4xl font-bold text-white leading-tight">
            探索无限<br />AI 可能性
          </h2>
          <p class="text-indigo-200/80 max-w-xs">
            连接、创造、分享。您的私人 AI 角色创作平台。
          </p>
        </div>
        <div class="text-xs text-white/20">© 2024 Nebula Inc.</div>
      </div>
    </div>

    <!-- Right Side: Login Form -->
    <div class="flex w-full flex-col justify-center p-8 lg:w-1/2 lg:p-16 relative">
      <div class="mx-auto w-full max-w-sm space-y-8">
        <div class="space-y-2 text-center lg:text-left">
          <h1 class="text-3xl font-bold text-white">欢迎回来</h1>
          <p class="text-slate-400 text-sm">请输入您的账号信息以继续</p>
        </div>

        <form class="space-y-6" @submit.prevent="handleSubmit">
          <div class="space-y-4">
            <div class="space-y-2">
              <label class="text-xs font-medium text-slate-300 uppercase tracking-wider">账号 / 邮箱</label>
              <input
                v-model.trim="form.identifier"
                class="glass-input"
                placeholder="demo 或邮箱"
                autocomplete="username"
                required
              />
            </div>
            <div class="space-y-2">
              <label class="text-xs font-medium text-slate-300 uppercase tracking-wider">密码</label>
              <input
                v-model="form.password"
                class="glass-input"
                placeholder="••••••"
                type="password"
                autocomplete="current-password"
                required
              />
            </div>
          </div>

          <div v-if="errorMessage" class="p-3 rounded-lg bg-red-500/10 border border-red-500/20 text-red-400 text-sm text-center">
            {{ errorMessage }}
          </div>

          <button
            type="submit"
            class="w-full btn-primary py-3 flex items-center justify-center gap-2 group"
            :disabled="isAuthenticating"
          >
            <span v-if="isAuthenticating" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <span>{{ isAuthenticating ? '登录中...' : '立即登录' }}</span>
            <svg v-if="!isAuthenticating" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4 group-hover:translate-x-1 transition-transform">
              <path fill-rule="evenodd" d="M3 10a.75.75 0 01.75-.75h10.638L10.23 5.29a.75.75 0 111.04-1.08l5.5 5.25a.75.75 0 010 1.08l-5.5 5.25a.75.75 0 11-1.04-1.08l4.158-3.96H3.75A.75.75 0 013 10z" clip-rule="evenodd" />
            </svg>
          </button>

          <p class="text-center text-sm text-slate-500 space-y-1">
            <span class="block">
              还没有账号？
              <NuxtLink to="/register" class="text-indigo-400 hover:text-indigo-300 transition-colors font-medium">立即注册</NuxtLink>
            </span>
            <span class="block">
              忘记密码？
              <NuxtLink to="/forgot" class="text-indigo-400 hover:text-indigo-300 transition-colors font-medium">找回密码</NuxtLink>
            </span>
          </p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['guest'], // 登录页使用 guest 守卫，已登录用户直接跳回受保护页面。
})

const router = useRouter()
const route = useRoute()
const { login, authError, isAuthenticating } = useAuth()

const form = reactive({
  identifier: 'demo',
  password: '123456',
})

const errorMessage = computed(() => authError.value)

const handleSubmit = async () => {
  try {
    await login({ identifier: form.identifier, password: form.password })
    const redirectPath = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    await router.push(redirectPath)
  } catch (error) {
    console.error('Login failed', error)
  }
}
</script>
