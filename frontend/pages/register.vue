<template>
  <section class="mx-auto flex min-h-[70vh] max-w-md flex-col justify-center space-y-8">
    <header class="space-y-2 text-center">
      <h1 class="text-3xl font-display font-bold text-charcoal-900">创建新账号</h1>
      <p class="text-sm text-charcoal-500">填写用户名与密码后即可开始体验 AI 角色聊天。</p>
    </header>
    <form class="space-y-5 card-soft p-8" @submit.prevent="handleSubmit">
      <div class="space-y-2">
        <label class="text-sm font-medium text-charcoal-700">用户名</label>
        <input
          v-model.trim="form.username"
          class="glass-input"
          placeholder="your-name"
          autocomplete="username"
          required
        />
      </div>
      <div class="space-y-2">
        <label class="text-sm font-medium text-charcoal-700">邮箱</label>
        <input
          v-model.trim="form.email"
          class="glass-input"
          placeholder="name@example.com"
          type="email"
          autocomplete="email"
          required
        />
      </div>
      <div class="space-y-2">
        <label class="text-sm font-medium text-charcoal-700 flex items-center justify-between">
          验证码
          <button
            type="button"
            class="text-xs text-charcoal-900 hover:text-accent-yellow font-bold hover:underline disabled:opacity-60 disabled:cursor-not-allowed transition-colors"
            :disabled="sendingCode || countdown > 0"
            @click="handleSendCode"
          >
            {{ countdown > 0 ? `重新获取 (${countdown}s)` : sendingCode ? '发送中...' : '获取验证码' }}
          </button>
        </label>
        <input
          v-model.trim="form.code"
          class="glass-input"
          placeholder="请输入邮箱验证码"
          autocomplete="one-time-code"
          required
        />
      </div>
      <div class="space-y-2">
        <label class="text-sm font-medium text-charcoal-700 flex items-center justify-between">
          昵称 <span class="text-xs text-charcoal-500">可选，默认与用户名一致</span>
        </label>
        <input
          v-model.trim="form.nickname"
          class="glass-input"
          placeholder="展示名称"
          autocomplete="nickname"
        />
      </div>
      <div class="space-y-2">
        <label class="text-sm font-medium text-charcoal-700">密码</label>
        <input
          v-model="form.password"
          class="glass-input"
          placeholder="请输入不少于 6 位"
          type="password"
          autocomplete="new-password"
          required
          minlength="6"
        />
      </div>
      <div class="space-y-2">
        <label class="text-sm font-medium text-charcoal-700">确认密码</label>
        <input
          v-model="form.confirm"
          class="glass-input"
          placeholder="再次输入密码"
          type="password"
          autocomplete="new-password"
          required
          minlength="6"
        />
      </div>
      <p v-if="errorMessage" class="text-sm text-status-error text-center">
        {{ errorMessage }}
      </p>
      <button
        type="submit"
        class="w-full btn-primary py-3 shadow-lg shadow-charcoal-900/10"
        :disabled="isAuthenticating"
      >
        {{ isAuthenticating ? '注册中...' : '创建账号并登录' }}
      </button>
      <p class="text-center text-sm text-charcoal-500">
        已有账号？
        <NuxtLink to="/login" class="text-charcoal-900 hover:text-accent-yellow font-bold decoration-accent-yellow underline decoration-2 underline-offset-2 transition-colors">返回登录</NuxtLink>
      </p>
    </form>
  </section>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['guest'],
})

const router = useRouter()
const route = useRoute()
const { register, sendCode, authError, isAuthenticating } = useAuth()

const form = reactive({
  username: '',
  email: '',
  code: '',
  nickname: '',
  password: '',
  confirm: '',
})

const localError = ref('')
const errorMessage = computed(() => localError.value || authError.value)
const sendingCode = ref(false)
const countdown = ref(0)
let timer: NodeJS.Timeout | null = null

onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
})

const handleSubmit = async () => {
  localError.value = ''
  if (!form.email.includes('@')) {
    localError.value = '请填写有效的邮箱'
    return
  }
  if (!form.code) {
    localError.value = '请先获取邮箱验证码'
    return
  }
  if (form.password !== form.confirm) {
    localError.value = '两次输入的密码不一致'
    return
  }
  try {
    await register({
      username: form.username,
      email: form.email,
      code: form.code,
      password: form.password,
      nickname: form.nickname,
    })
    const redirectPath = typeof route.query.redirect === 'string' ? route.query.redirect : '/search'
    await router.push(redirectPath)
  } catch (error) {
    console.error('Register failed', error)
  }
}

const handleSendCode = async () => {
  localError.value = ''
  if (!form.email.includes('@')) {
    localError.value = '请填写有效的邮箱'
    return
  }
  try {
    sendingCode.value = true
    await sendCode(form.email, 'signup')
    countdown.value = 60
    timer && clearInterval(timer)
    timer = setInterval(() => {
      countdown.value -= 1
      if (countdown.value <= 0 && timer) {
        clearInterval(timer)
        timer = null
      }
    }, 1000)
  } catch (error: any) {
    localError.value = error?.message || '验证码发送失败，请稍后重试'
  } finally {
    sendingCode.value = false
  }
}
</script>
