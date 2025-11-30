<template>
  <div class="flex min-h-[80vh] w-full overflow-hidden card-soft">
    <!-- Left Side: Visual Brand -->
    <div class="relative hidden w-1/2 overflow-hidden lg:block bg-bg-cream-200">
      <div class="absolute inset-0 bg-gradient-to-br from-bg-cream-100 to-bg-pink-100 opacity-60"></div>
      
      <!-- Abstract Shapes -->
      <div class="absolute -left-20 -top-20 h-96 w-96 rounded-full bg-accent-yellow/20 blur-[100px] animate-float"></div>
      <div class="absolute bottom-0 right-0 h-80 w-80 rounded-full bg-accent-pink-soft/30 blur-[80px] animate-float-delayed"></div>

      <div class="relative z-10 flex h-full flex-col justify-between p-12">
        <div class="text-2xl font-display font-bold text-charcoal-900 tracking-tight flex items-center gap-2">
          <div class="w-8 h-8 rounded-lg bg-accent-yellow flex items-center justify-center text-charcoal-900 text-xs font-bold shadow-sm">M</div>
          MilkyVerse
        </div>
        <div class="space-y-4">
          <h2 class="text-4xl font-display font-bold text-charcoal-900 leading-tight">
            探索无限<br />AI 可能性
          </h2>
          <p class="text-charcoal-500 max-w-xs text-lg">
            连接、创造、分享。您的私人 AI 角色创作平台。
          </p>
        </div>
        <div class="text-xs text-charcoal-500">© 2024 MilkyVerse Inc.</div>
      </div>
    </div>

    <!-- Right Side: Login Form -->
    <div class="flex w-full flex-col justify-center p-8 lg:w-1/2 lg:p-16 relative bg-white">
      <div class="mx-auto w-full max-w-sm space-y-8">
        <div class="space-y-2 text-center lg:text-left">
          <h1 class="text-3xl font-display font-bold text-charcoal-900">欢迎回来</h1>
          <p class="text-charcoal-500 text-sm">请输入您的账号信息以继续</p>
        </div>

        <form class="space-y-6" @submit.prevent="handleSubmit">
          <div class="space-y-4">
            <div class="space-y-2">
              <label class="text-xs font-medium text-charcoal-500 uppercase tracking-wider">账号 / 邮箱</label>
              <input
                v-model.trim="form.identifier"
                class="glass-input"
                placeholder="demo 或邮箱"
                autocomplete="username"
                required
              />
            </div>
            <div class="space-y-2">
              <label class="text-xs font-medium text-charcoal-500 uppercase tracking-wider">密码</label>
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

          <div v-if="errorMessage" class="p-3 rounded-lg bg-status-error/10 border border-status-error/20 text-status-error text-sm text-center">
            {{ errorMessage }}
          </div>

          <button
            type="submit"
            class="w-full btn-primary py-3 flex items-center justify-center gap-2 group shadow-lg shadow-charcoal-900/10"
            :disabled="isAuthenticating"
          >
            <span v-if="isAuthenticating" class="w-4 h-4 border-2 border-bg-cream-100/30 border-t-bg-cream-100 rounded-full animate-spin"></span>
            <span>{{ isAuthenticating ? '登录中...' : '立即登录' }}</span>
            <svg v-if="!isAuthenticating" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4 group-hover:translate-x-1 transition-transform">
              <path fill-rule="evenodd" d="M3 10a.75.75 0 01.75-.75h10.638L10.23 5.29a.75.75 0 111.04-1.08l5.5 5.25a.75.75 0 010 1.08l-5.5 5.25a.75.75 0 11-1.04-1.08l4.158-3.96H3.75A.75.75 0 013 10z" clip-rule="evenodd" />
            </svg>
          </button>

          <p class="text-center text-sm text-charcoal-500 space-y-1">
            <span class="block">
              还没有账号？
              <NuxtLink to="/register" class="text-charcoal-900 hover:text-accent-yellow transition-colors font-bold decoration-accent-yellow underline decoration-2 underline-offset-2">立即注册</NuxtLink>
            </span>
            <span class="block">
              忘记密码？
              <NuxtLink to="/forgot" class="text-charcoal-500 hover:text-charcoal-900 transition-colors font-medium">找回密码</NuxtLink>
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
