<template>
  <div class="min-h-screen pb-20 space-y-8">
    <!-- Profile Header & Assets (Bento Row 1) -->
    <section class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Identity Card -->
      <div class="lg:col-span-2 card-soft p-0 overflow-hidden relative group">
        <!-- Abstract Cover -->
        <div class="h-32 bg-gradient-to-r from-accent-pink/20 via-bg-cream-200 to-accent-yellow/20 relative">
          <div class="absolute inset-0 bg-[url('https://www.transparenttextures.com/patterns/cubes.png')] opacity-30"></div>
          <div class="absolute top-4 right-4 flex gap-2">
            <NuxtLink to="/me/settings" class="p-2 rounded-full bg-white/50 hover:bg-white text-charcoal-600 transition-all backdrop-blur-sm" title="设置">
              <div class="i-ph-gear-six-fill text-lg" />
            </NuxtLink>
            <NuxtLink to="/notifications" class="p-2 rounded-full bg-white/50 hover:bg-white text-charcoal-600 transition-all backdrop-blur-sm relative" title="通知">
              <div class="i-ph-bell-fill text-lg" />
              <span class="absolute top-2 right-2 w-2 h-2 bg-status-error rounded-full border border-white" v-if="false"></span>
            </NuxtLink>
          </div>
        </div>
        
        <div class="px-8 pb-8 relative">
          <!-- Avatar -->
          <div class="absolute -top-12 left-8">
            <div class="relative">
              <img 
                :src="home?.user?.avatar_url || defaultAvatar" 
                class="h-24 w-24 rounded-2xl border-4 border-white shadow-lg object-cover bg-white" 
                alt="avatar" 
              />
              <div class="absolute bottom-1 right-1 w-5 h-5 bg-status-success rounded-full border-2 border-white" title="在线"></div>
            </div>
          </div>

          <!-- Info -->
          <div class="mt-14 flex flex-col md:flex-row md:items-end justify-between gap-4">
            <div>
              <h1 class="text-3xl font-display font-bold text-charcoal-900 tracking-tight flex items-center gap-2">
                {{ home?.user?.nickname || auth.user.value?.nickname || '旅者' }}
                <span class="px-2 py-0.5 rounded-md bg-charcoal-100 text-charcoal-500 text-xs font-bold uppercase tracking-wider">LV.{{ userLevel }}</span>
              </h1>
              <div class="flex items-center gap-3 mt-2 text-sm text-charcoal-500">
                <span class="font-mono bg-bg-cream-200 px-2 py-0.5 rounded text-charcoal-600">ID: {{ home?.user?.id || auth.user.value?.id }}</span>
                <span class="flex items-center gap-1">
                  <div class="i-ph-smiley text-accent-pink" />
                  心情指数 {{ engagementScore }}
                </span>
              </div>
              <p class="mt-3 text-charcoal-600 max-w-md text-sm leading-relaxed">
                {{ home?.user?.bio || '这个家伙很懒，什么都没有写...' }}
              </p>
            </div>
            
            <!-- Quick Stats -->
            <div class="flex gap-6 border-t md:border-t-0 md:border-l border-charcoal-100 pt-4 md:pt-0 md:pl-6">
              <div class="text-center">
                <div class="text-xl font-bold text-charcoal-900">{{ home?.stats?.post_total || 0 }}</div>
                <div class="text-xs text-charcoal-400 font-medium uppercase tracking-wide">帖子</div>
              </div>
              <div class="text-center">
                <div class="text-xl font-bold text-charcoal-900">{{ home?.stats?.favorite_total || 0 }}</div>
                <div class="text-xs text-charcoal-400 font-medium uppercase tracking-wide">收藏</div>
              </div>
              <div class="text-center">
                <div class="text-xl font-bold text-charcoal-900">{{ home?.stats?.session_total || 0 }}</div>
                <div class="text-xs text-charcoal-400 font-medium uppercase tracking-wide">会话</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Wallet & Assets Card -->
      <div class="card-soft p-6 flex flex-col justify-between relative overflow-hidden">
        <div class="absolute top-0 right-0 w-40 h-40 bg-accent-yellow/10 rounded-full blur-3xl -mr-10 -mt-10 pointer-events-none"></div>
        
        <div>
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-lg font-bold text-charcoal-900 flex items-center gap-2">
              <div class="i-ph-wallet-fill text-accent-yellow-dark" />
              我的资产
            </h2>
            <NuxtLink to="/me/wallet" class="text-xs font-medium text-charcoal-500 hover:text-charcoal-900 flex items-center gap-1">
              明细 <div class="i-ph-caret-right" />
            </NuxtLink>
          </div>

          <div class="space-y-4">
            <div class="bg-bg-cream-100 rounded-xl p-4 border border-charcoal-100">
              <div class="text-xs text-charcoal-500 mb-1">平台币余额</div>
              <div class="text-3xl font-bold text-charcoal-900 font-display flex items-baseline gap-1">
                {{ home?.assets?.balance?.toLocaleString() ?? 0 }}
                <span class="text-sm font-normal text-charcoal-400">Coins</span>
              </div>
            </div>

            <div class="flex gap-3">
              <div class="flex-1 bg-white border border-charcoal-100 rounded-xl p-3 flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-accent-pink/10 flex items-center justify-center text-accent-pink">
                  <div class="i-ph-ticket-fill text-xl" />
                </div>
                <div>
                  <div class="text-lg font-bold text-charcoal-900">{{ home?.assets?.monthly_tickets ?? 0 }}</div>
                  <div class="text-[10px] text-charcoal-400 uppercase">月票</div>
                </div>
              </div>
              <div class="flex-1 bg-white border border-charcoal-100 rounded-xl p-3 flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-accent-yellow/10 flex items-center justify-center text-accent-yellow-dark">
                  <div class="i-ph-gift-fill text-xl" />
                </div>
                <div>
                  <div class="text-lg font-bold text-charcoal-900">0</div>
                  <div class="text-[10px] text-charcoal-400 uppercase">道具</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="mt-6 pt-6 border-t border-charcoal-100">
          <NuxtLink to="/store" class="w-full btn-primary justify-center py-3 shadow-xl shadow-accent-yellow/20">
            <div class="i-ph-plus-circle-fill" />
            立即充值
          </NuxtLink>
        </div>
      </div>
    </section>

    <!-- Creator & Tools (Bento Row 2) -->
    <section class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- Creator Journey Banner -->
      <div class="lg:col-span-3 rounded-3xl p-1 bg-gradient-to-br from-charcoal-900 via-charcoal-800 to-charcoal-900 shadow-xl shadow-charcoal-900/10 group relative overflow-hidden">
        <div class="absolute inset-0 bg-[url('https://www.transparenttextures.com/patterns/stardust.png')] opacity-20"></div>
        <div class="absolute right-0 top-0 bottom-0 w-1/2 bg-gradient-to-l from-accent-pink/20 to-transparent"></div>
        
        <div class="bg-charcoal-900/50 backdrop-blur-sm h-full rounded-[20px] p-6 md:p-8 flex flex-col md:flex-row items-center justify-between gap-6 relative z-10">
          <div class="flex-1">
            <div class="flex items-center gap-2 mb-2">
              <span class="px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-widest bg-accent-yellow text-charcoal-900">Creator Studio</span>
              <span class="text-charcoal-400 text-xs">Level {{ creatorLevel }}</span>
            </div>
            <h3 class="text-2xl font-bold text-white mb-2">释放你的创造力</h3>
            <p class="text-charcoal-300 text-sm max-w-lg">
              已创建 <span class="text-white font-bold">{{ home?.creator?.roles_total || 0 }}</span> 个角色，
              发布 <span class="text-white font-bold">{{ home?.creator?.published_roles || 0 }}</span> 个。
              <span v-if="(home?.creator?.wallet_balance || 0) > 0">
                累计收益 <span class="text-accent-yellow font-bold">{{ home?.creator?.wallet_balance }}</span>。
              </span>
              <span v-else>
                开启创作之旅，赚取收益。
              </span>
            </p>
          </div>
          <div class="flex items-center gap-4">
             <div class="hidden md:block text-right mr-4">
                <div class="text-xs text-charcoal-400 mb-1">发布进度</div>
                <div class="w-32 h-2 bg-charcoal-700 rounded-full overflow-hidden">
                  <div class="h-full bg-gradient-to-r from-accent-yellow to-accent-pink" :style="{ width: `${creatorProgress * 100}%` }"></div>
                </div>
             </div>
             <NuxtLink to="/creator" class="px-6 py-3 rounded-xl bg-white text-charcoal-900 font-bold hover:bg-bg-cream-100 transition-colors flex items-center gap-2 shadow-lg">
               进入控制台
               <div class="i-ph-arrow-right-bold" />
             </NuxtLink>
          </div>
        </div>
      </div>

      <!-- Quick Actions Grid -->
      <div class="grid grid-cols-2 gap-3">
        <NuxtLink to="/me/favorites" class="card-soft p-4 flex flex-col items-center justify-center gap-2 hover:-translate-y-1 transition-transform cursor-pointer group">
          <div class="w-10 h-10 rounded-full bg-bg-cream-200 group-hover:bg-accent-yellow/20 flex items-center justify-center text-charcoal-600 group-hover:text-charcoal-900 transition-colors">
            <div class="i-ph-folder-star-fill text-xl" />
          </div>
          <span class="text-xs font-medium text-charcoal-600">收藏夹</span>
        </NuxtLink>
        <NuxtLink to="/me/history" class="card-soft p-4 flex flex-col items-center justify-center gap-2 hover:-translate-y-1 transition-transform cursor-pointer group">
          <div class="w-10 h-10 rounded-full bg-bg-cream-200 group-hover:bg-accent-pink/20 flex items-center justify-center text-charcoal-600 group-hover:text-charcoal-900 transition-colors">
            <div class="i-ph-clock-counter-clockwise-fill text-xl" />
          </div>
          <span class="text-xs font-medium text-charcoal-600">浏览历史</span>
        </NuxtLink>
        <NuxtLink to="/me/settings" class="card-soft p-4 flex flex-col items-center justify-center gap-2 hover:-translate-y-1 transition-transform cursor-pointer group">
          <div class="w-10 h-10 rounded-full bg-bg-cream-200 group-hover:bg-charcoal-100 flex items-center justify-center text-charcoal-600 group-hover:text-charcoal-900 transition-colors">
            <div class="i-ph-sliders-horizontal-fill text-xl" />
          </div>
          <span class="text-xs font-medium text-charcoal-600">偏好设置</span>
        </NuxtLink>
        <button class="card-soft p-4 flex flex-col items-center justify-center gap-2 hover:-translate-y-1 transition-transform cursor-pointer group" @click="logout">
          <div class="w-10 h-10 rounded-full bg-bg-cream-200 group-hover:bg-status-error/10 flex items-center justify-center text-charcoal-600 group-hover:text-status-error transition-colors">
            <div class="i-ph-sign-out-fill text-xl" />
          </div>
          <span class="text-xs font-medium text-charcoal-600 group-hover:text-status-error">退出登录</span>
        </button>
      </div>
    </section>

    <!-- Content Tabs (Bento Row 3) -->
    <section class="card-soft min-h-[500px]">
      <div class="border-b border-charcoal-100 px-6 pt-6">
        <div class="flex items-center gap-8">
          <button 
            v-for="tab in tabs" 
            :key="tab.id"
            @click="activeTab = tab.id"
            class="pb-4 text-sm font-bold transition-all relative"
            :class="activeTab === tab.id ? 'text-charcoal-900' : 'text-charcoal-400 hover:text-charcoal-600'"
          >
            {{ tab.label }}
            <span 
              v-if="activeTab === tab.id" 
              class="absolute bottom-0 left-0 right-0 h-0.5 bg-accent-yellow rounded-t-full"
            ></span>
          </button>
        </div>
      </div>

      <div class="p-6">
        <div v-if="loading" class="py-20 text-center text-charcoal-400 flex flex-col items-center gap-3">
          <div class="i-svg-spinners-90-ring-with-bg text-2xl" />
          <p class="text-sm">加载数据中...</p>
        </div>

        <div v-else>
          <!-- Recent Views -->
          <div v-if="activeTab === 'recent'" class="space-y-6">
            <div v-if="!home?.recent_views?.length" class="py-20 text-center text-charcoal-400">
              <div class="i-ph-footprints-fill text-4xl mx-auto mb-3 opacity-20" />
              <p class="text-sm">暂无浏览记录，去<NuxtLink to="/community" class="text-accent-pink hover:underline">社区</NuxtLink>逛逛吧。</p>
            </div>
            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <CommunityPostCard v-for="post in home?.recent_views" :key="post.id" :post="post" />
            </div>
          </div>

          <!-- Favorites -->
          <div v-if="activeTab === 'favorites'" class="space-y-6">
             <div v-if="!home?.favorites?.length" class="py-20 text-center text-charcoal-400">
              <div class="i-ph-star-fill text-4xl mx-auto mb-3 opacity-20" />
              <p class="text-sm">还没有收藏任何内容。</p>
            </div>
            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <CommunityPostCard v-for="post in home?.favorites" :key="post.id" :post="post" />
            </div>
          </div>

          <!-- My Posts -->
          <div v-if="activeTab === 'posts'" class="space-y-6">
             <div v-if="!home?.my_posts?.length" class="py-20 text-center text-charcoal-400">
              <div class="i-ph-pencil-simple-slash-fill text-4xl mx-auto mb-3 opacity-20" />
              <p class="text-sm">你还没有发布过任何帖子。</p>
              <NuxtLink to="/community/new" class="btn-primary mt-4 text-xs">发布第一条动态</NuxtLink>
            </div>
            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <CommunityPostCard v-for="post in home?.my_posts" :key="post.id" :post="post" />
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { UserHomePayload } from '@/types'

definePageMeta({
  middleware: ['auth'],
})

const api = useApi()
const auth = useAuth()
const router = useRouter()
const home = ref<UserHomePayload | null>(null)
const loading = ref(false)
const defaultAvatar = 'https://placehold.co/150x150?text=Me'

const activeTab = ref('recent')
const tabs = [
  { id: 'recent', label: '近期浏览' },
  { id: 'favorites', label: '我的收藏' },
  { id: 'posts', label: '我的发布' },
]

const engagementScore = computed(() => {
  const stats = home.value?.stats
  if (!stats) return 40
  const base = (stats.favorite_total ?? 0) * 8 + (stats.session_total ?? 0) * 5 + (stats.recent_view_total ?? 0) * 4
  return Math.min(100, Math.max(20, base))
})

const userLevel = computed(() => {
  return Math.floor(engagementScore.value / 10) + 1
})

const creatorProgress = computed(() => {
  const total = home.value?.creator?.roles_total ?? 0
  const published = home.value?.creator?.published_roles ?? 0
  if (!total) return 0
  return Math.min(1, published / total)
})

const creatorLevel = computed(() => {
  const published = home.value?.creator?.published_roles ?? 0
  if (published > 10) return 3
  if (published > 3) return 2
  return 1
})

const load = async () => {
  loading.value = true
  try {
    const res = await api.get<{ data: UserHomePayload }>('/me/home')
    home.value = res.data
  } catch (error) {
    console.error('Failed to load profile', error)
  } finally {
    loading.value = false
  }
}

const logout = async () => {
  if (confirm('确定要退出登录吗？')) {
    auth.logout()
    router.push('/login')
  }
}

onMounted(load)
</script>

<style scoped>
/* Custom Scrollbar for Bento Grid if needed */
.i-ph-gear-six-fill, .i-ph-bell-fill, .i-ph-smiley, .i-ph-wallet-fill, .i-ph-caret-right, .i-ph-ticket-fill, .i-ph-gift-fill, .i-ph-plus-circle-fill, .i-ph-arrow-right-bold, .i-ph-folder-star-fill, .i-ph-clock-counter-clockwise-fill, .i-ph-sliders-horizontal-fill, .i-ph-sign-out-fill, .i-ph-footprints-fill, .i-ph-star-fill, .i-ph-pencil-simple-slash-fill {
  display: inline-block;
  vertical-align: middle;
}
</style>
