<template>
  <section class="space-y-10">
    <!-- Hero -->
    <header class="hero-card relative overflow-hidden rounded-3xl border border-white/10 p-6 md:p-8">
      <div class="pointer-events-none absolute inset-y-0 right-0 w-1/2 opacity-60 blur-3xl" :style="heroStyle"></div>
      <div class="relative flex flex-col gap-6 md:flex-row md:items-center md:justify-between">
        <div class="flex items-center gap-5">
          <div class="relative">
            <img :src="home?.user?.avatar_url || defaultAvatar" class="h-20 w-20 rounded-2xl border-2 border-white/40 object-cover" alt="avatar" />
            <span class="pulse-dot"></span>
          </div>
          <div>
            <p class="text-sm text-slate-200/70">欢迎回来</p>
            <h1 class="text-3xl font-semibold text-white tracking-tight">
              {{ home?.user?.nickname || auth.user.value?.nickname || '旅者' }}
            </h1>
            <p class="text-xs text-slate-200/60 mt-1">ID: {{ home?.user?.id || auth.user.value?.id }}</p>
            <p class="text-xs text-slate-200/60">今日心情指数：{{ engagementScore }} / 100</p>
          </div>
        </div>
        <div class="flex flex-wrap gap-3">
          <NuxtLink class="btn-secondary" to="/me/settings">管理个人资料</NuxtLink>
          <NuxtLink class="btn-secondary" to="/me/wallet">虚拟钱包</NuxtLink>
          <NuxtLink class="btn-primary" to="/creator">
            {{ home?.creator?.has_creator_access ? '进入创作者中心' : '成为创作者' }}
          </NuxtLink>
        </div>
      </div>
      <div class="mt-6 flex flex-wrap gap-3">
        <button
          v-for="action in quickActions"
          :key="action.label"
          class="action-chip"
          @click="go(action.to)"
        >
          <span class="text-base">{{ action.icon }}</span>
          <span>{{ action.label }}</span>
        </button>
      </div>
      <p v-if="errorMessage" class="mt-4 text-sm text-rose-300">{{ errorMessage }}</p>
    </header>

    <!-- Stats -->
    <section class="grid gap-4 md:grid-cols-3">
      <div v-for="card in statCards" :key="card.title" class="glass-card relative overflow-hidden">
        <div class="absolute inset-0 opacity-50" :style="card.bg"></div>
        <div class="relative">
          <p class="text-sm text-slate-200/70">{{ card.title }}</p>
          <p class="mt-2 text-3xl font-semibold text-white">{{ card.value }}</p>
          <p class="text-xs text-slate-200/60">{{ card.hint }}</p>
        </div>
      </div>
    </section>

    <!-- Interactive widgets -->
    <section class="grid gap-6 lg:grid-cols-3">
      <div class="glass-card col-span-2">
        <header class="flex items-center justify-between text-white">
          <div>
            <h2 class="text-lg font-semibold">收藏内容</h2>
            <p class="text-sm text-slate-200/70">快速回到已收藏的灵感。</p>
          </div>
          <NuxtLink class="link text-sm" to="/me/favorites">查看全部</NuxtLink>
        </header>
        <div v-if="loading" class="mt-5 text-sm text-slate-200/70">加载中...</div>
        <div v-else-if="!home?.favorites?.length" class="mt-5 text-sm text-slate-200/70">还没有收藏，去社区探索吧。</div>
        <div v-else class="mt-5 space-y-4">
          <CommunityPostCard v-for="post in home?.favorites || []" :key="post.id" :post="post" />
        </div>
      </div>

      <div class="glass-card space-y-6">
        <div>
          <h3 class="text-sm text-slate-200/70">创作者钱包快照</h3>
          <p class="mt-2 text-3xl font-semibold text-white">{{ home?.creator?.wallet_balance ?? 0 }}</p>
          <p class="text-xs text-slate-200/60">可提取余额</p>
        </div>
        <div class="relative flex items-center justify-center">
          <div class="ring-progress">
            <svg viewBox="0 0 120 120">
              <defs>
                <linearGradient id="gradientRing" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#a855f7" />
                  <stop offset="100%" stop-color="#f472b6" />
                </linearGradient>
              </defs>
              <circle cx="60" cy="60" r="52" />
              <circle
                cx="60"
                cy="60"
                r="52"
                :style="{ strokeDashoffset: 326 - 326 * creatorProgress }"
              />
            </svg>
            <div class="ring-label">
              <p class="text-xs text-slate-200/70">角色发布率</p>
              <p class="text-xl font-semibold text-white">{{ Math.round(creatorProgress * 100) }}%</p>
            </div>
          </div>
        </div>
        <NuxtLink class="btn-primary w-full justify-center" to="/creator">激活创作者能力</NuxtLink>
      </div>
    </section>

    <!-- Feed sections -->
    <section class="grid gap-6 lg:grid-cols-2">
      <div class="glass-card">
        <header class="flex items-center justify-between text-white">
          <div>
            <h2 class="text-lg font-semibold">近期浏览</h2>
            <p class="text-sm text-slate-200/70">记录你最近看过的帖子，便于继续阅读。</p>
          </div>
          <NuxtLink class="link text-sm" to="/community">继续逛逛</NuxtLink>
        </header>
        <div v-if="loading" class="mt-5 text-sm text-slate-200/70">加载中...</div>
        <div v-else-if="!home?.recent_views?.length" class="mt-5 text-sm text-slate-200/70">暂无浏览记录。</div>
        <div v-else class="mt-5 space-y-4">
          <CommunityPostCard v-for="post in home?.recent_views || []" :key="post.id" :post="post" />
        </div>
      </div>

      <div class="glass-card">
        <header class="flex flex-wrap items-center justify-between gap-3 text-white">
          <div>
            <h2 class="text-lg font-semibold">我的帖子</h2>
            <p class="text-sm text-slate-200/70">展示你在社区发表的内容。</p>
          </div>
          <NuxtLink class="btn-primary" to="/community/new">发布新帖</NuxtLink>
        </header>
        <div v-if="loading" class="mt-5 text-sm text-slate-200/70">加载中...</div>
        <div v-else-if="!home?.my_posts?.length" class="mt-5 text-sm text-slate-200/70">还没有发表内容，快去分享吧！</div>
        <div v-else class="mt-5 space-y-4">
          <CommunityPostCard v-for="post in home?.my_posts || []" :key="post.id" :post="post" />
        </div>
      </div>
    </section>

    <!-- Creator CTA -->
    <section class="rounded-3xl border border-white/10 bg-gradient-to-r from-indigo-500/30 via-purple-500/30 to-rose-500/30 p-6 shadow-xl shadow-indigo-900/30">
      <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between text-white">
        <div>
          <p class="text-sm uppercase tracking-[0.3em] text-white/80">Creator Journey</p>
          <h2 class="mt-2 text-2xl font-semibold">创作者成长引导</h2>
          <p class="text-sm text-white/80">
            已创建 {{ home?.creator?.roles_total || 0 }} 个角色，发布 {{ home?.creator?.published_roles || 0 }} 个，创作者钱包可用余额
            {{ home?.creator?.wallet_balance || 0 }}。
          </p>
        </div>
        <NuxtLink class="btn-primary bg-white/90 text-slate-800 hover:bg-white" to="/creator">前往创作者中心</NuxtLink>
      </div>
    </section>
  </section>
</template>

<script setup lang="ts">
import type { UserHomePayload, PaymentOrder } from '@/types'

definePageMeta({
  middleware: ['auth'],
})

const api = useApi()
const router = useRouter()
const auth = useAuth()
const home = ref<UserHomePayload | null>(null)
const loading = ref(false)
const errorMessage = ref('')
const defaultAvatar = 'https://placehold.co/80x80?text=Me'
const accentPool = [
  ['#a855f7', '#6366f1'],
  ['#4ade80', '#22d3ee'],
  ['#f97316', '#f43f5e'],
]
const heroStyle = computed(() => {
  const palette = accentPool[new Date().getHours() % accentPool.length]
  return `background: radial-gradient(circle at 30% 20%, ${palette[0]}55, transparent), radial-gradient(circle at 80% 0%, ${palette[1]}88, transparent)`
})

const quickActions = [
  { label: '收藏夹', icon: '📁', to: '/me/favorites' },
  { label: '我的钱包', icon: '💰', to: '/me/wallet' },
  { label: '充值记录', icon: '🧾', to: '/me/payments' },
  { label: '创作者中心', icon: '🚀', to: '/creator' },
  { label: '通知中心', icon: '🔔', to: '/notifications' },
]

const statCards = computed(() => [
  {
    title: '虚拟货币余额',
    value: home.value?.assets?.balance ?? 0,
    hint: '平台内可用代币',
    bg: 'background: linear-gradient(135deg, rgba(99,102,241,0.2), rgba(56,189,248,0.2))',
  },
  {
    title: '本月权益',
    value: home.value?.assets?.monthly_tickets ?? 0,
    hint: '剩余月票',
    bg: 'background: linear-gradient(135deg, rgba(34,197,94,0.2), rgba(252,211,77,0.2))',
  },
  {
    title: '聊天会话',
    value: home.value?.stats?.session_total ?? 0,
    hint: '最近活跃记录',
    bg: 'background: linear-gradient(135deg, rgba(248,113,113,0.2), rgba(251,191,36,0.2))',
  },
  {
    title: '收藏帖子',
    value: home.value?.stats?.favorite_total ?? 0,
    hint: '来自个人收藏夹',
    bg: 'background: linear-gradient(135deg, rgba(192,38,211,0.2), rgba(59,130,246,0.2))',
  },
  {
    title: '近期浏览',
    value: home.value?.stats?.recent_view_total ?? 0,
    hint: '记录最近的足迹',
    bg: 'background: linear-gradient(135deg, rgba(14,165,233,0.2), rgba(236,72,153,0.2))',
  },
  {
    title: '已发布帖子',
    value: home.value?.stats?.post_total ?? 0,
    hint: '你的创作数量',
    bg: 'background: linear-gradient(135deg, rgba(248,250,252,0.2), rgba(110,231,183,0.2))',
  },
])

const engagementScore = computed(() => {
  const stats = home.value?.stats
  if (!stats) return 40
  const base = (stats.favorite_total ?? 0) * 8 + (stats.session_total ?? 0) * 5 + (stats.recent_view_total ?? 0) * 4
  return Math.min(100, Math.max(20, base))
})

const creatorProgress = computed(() => {
  const total = home.value?.creator?.roles_total ?? 0
  const published = home.value?.creator?.published_roles ?? 0
  if (!total) return 0
  return Math.min(1, published / total)
})

const go = (path: string) => router.push(path)

const load = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await api.get<{ data: UserHomePayload }>('/me/home')
    home.value = res.data
  } catch (error: any) {
    errorMessage.value = error?.data?.error || error?.message || '无法加载个人主页数据，请稍后再试'
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.btn-primary {
  @apply inline-flex items-center gap-2 rounded-full bg-primary px-5 py-2 text-sm font-semibold text-white shadow hover:opacity-90 transition-all;
}
.btn-secondary {
  @apply inline-flex items-center rounded-full border border-white/30 px-5 py-2 text-sm text-white/80 hover:border-white/60 transition;
}
.hero-card {
  background: rgba(15, 23, 42, 0.6);
  box-shadow: 0 25px 80px rgba(15, 23, 42, 0.45);
}
.glass-card {
  @apply relative rounded-3xl border border-white/10 p-6;
  background: rgba(15, 23, 42, 0.55);
  box-shadow: 0 10px 40px rgba(15, 23, 42, 0.3);
}
.action-chip {
  @apply inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/10 px-4 py-2 text-xs text-white/80 backdrop-blur transition hover:bg-white/20;
}
.pulse-dot {
  position: absolute;
  right: -4px;
  bottom: -4px;
  width: 14px;
  height: 14px;
  border-radius: 999px;
  background: #34d399;
  box-shadow: 0 0 10px #34d399;
  animation: pulse 2s infinite;
}
.ring-progress {
  position: relative;
  width: 140px;
  height: 140px;
}
.ring-progress svg {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}
.ring-progress circle {
  fill: none;
  stroke-width: 6px;
  stroke: rgba(148, 163, 184, 0.2);
  stroke-linecap: round;
}
.ring-progress circle:last-child {
  stroke: url(#gradientRing);
  stroke-dasharray: 326;
  stroke-dashoffset: 326;
  transition: stroke-dashoffset 0.8s ease;
}
.ring-label {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.link {
  @apply text-primary hover:text-white transition;
}
@keyframes pulse {
  0% {
    transform: scale(0.8);
    opacity: 0.6;
  }
  50% {
    transform: scale(1);
    opacity: 1;
  }
  100% {
    transform: scale(0.8);
    opacity: 0.6;
  }
}
</style>
