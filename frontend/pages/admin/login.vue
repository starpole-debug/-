<template>
  <section class="mx-auto mt-24 max-w-md glass-card p-10 animate-fade-up">
    <div class="text-center mb-8">
      <h1 class="text-2xl font-bold text-white">管理员入口</h1>
      <p class="text-sm text-slate-400 mt-2">需要先验证访问口令，才可继续登录</p>
    </div>

    <!-- 访问口令校验 -->
    <form v-if="!hasAccess" class="space-y-5" @submit.prevent="unlock">
      <div class="space-y-1">
        <label class="text-xs font-medium text-slate-400 ml-1">访问口令（X-Admin-Access）</label>
        <input v-model="accessKeyInput" type="password" placeholder="仅内部人员知晓的口令" class="glass-input" />
      </div>
      <p v-if="accessError" class="text-sm text-rose-300 text-center">{{ accessError }}</p>
      <button class="w-full btn-primary py-3 mt-4 shadow-lg shadow-indigo-500/20">验证并继续</button>
    </form>

    <!-- 管理员账号登录 -->
    <form v-else class="space-y-5" @submit.prevent="submit">
      <div class="space-y-1">
        <label class="text-xs font-medium text-slate-400 ml-1">用户名</label>
        <input v-model="username" placeholder="admin" class="glass-input" />
      </div>
      <div class="space-y-1">
        <label class="text-xs font-medium text-slate-400 ml-1">密码</label>
        <input v-model="password" type="password" placeholder="••••••••" class="glass-input" />
      </div>
      <div class="space-y-1">
        <label class="text-xs font-medium text-slate-400 ml-1">Admin Secret</label>
        <input v-model="secret" type="password" placeholder="••••••••" class="glass-input" />
      </div>
      <p v-if="errorMessage" class="text-sm text-rose-300 text-center">{{ errorMessage }}</p>
      <button
        class="w-full btn-primary py-3 mt-4 shadow-lg shadow-indigo-500/20 disabled:opacity-60 disabled:cursor-not-allowed"
        :disabled="adminAuth.isAuthenticating.value"
      >
        {{ adminAuth.isAuthenticating.value ? '登录中...' : '进入后台' }}
      </button>
      <p class="text-center text-xs text-slate-500">
        已验证访问口令，可直接登录。若口令变更，请<a class="text-indigo-300 cursor-pointer" @click="resetAccess">重新验证</a>。
      </p>
    </form>
  </section>
</template>

<script setup lang="ts">
const router = useRouter()
const adminAuth = useAdminAuth()
const username = ref('')
const password = ref('')
const secret = ref('')
const accessKeyCookie = useCookie<string | null>('admin-access', {
  sameSite: 'lax',
  secure: process.env.NODE_ENV === 'production',
  default: () => null,
})
const accessKeyInput = ref('')
const accessError = ref('')
const hasAccess = computed(() => Boolean(accessKeyCookie.value))
const errorMessage = computed(() => adminAuth.authError.value)

if (process.client) {
  watchEffect(() => {
    if (adminAuth.isAuthenticated.value) {
      router.push('/admin/models')
    }
  })
}

const submit = async () => {
  try {
    await adminAuth.login({
      username: username.value,
      password: password.value,
      adminSecret: secret.value,
      accessKey: accessKeyCookie.value || accessKeyInput.value,
    })
  } catch (error) {
    console.error('Admin login failed', error)
    if ((error as Error)?.message?.toLowerCase().includes('access denied')) {
      resetAccess()
      accessError.value = '访问口令无效，请重新输入'
    }
  }
}

const unlock = () => {
  accessError.value = ''
  if (!accessKeyInput.value) {
    accessError.value = '请输入访问口令'
    return
  }
  try {
    accessKeyCookie.value = encodeURIComponent(accessKeyInput.value)
  } catch {
    accessKeyCookie.value = accessKeyInput.value
  }
}

const resetAccess = () => {
  accessKeyCookie.value = null
  accessKeyInput.value = ''
  adminAuth.clearAccessKey()
}
</script>
