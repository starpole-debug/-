<template>
  <section class="space-y-12">
    <!-- Community Hero -->
    <header class="relative rounded-3xl overflow-hidden bg-[#151621] border border-white/5 p-8 md:p-12 animate-fade-up">
      <div class="pointer-events-none absolute top-0 right-0 w-64 h-64 bg-indigo-500/20 rounded-full blur-[80px] animate-float"></div>
      <div class="pointer-events-none absolute bottom-0 left-0 w-64 h-64 bg-purple-500/20 rounded-full blur-[80px] animate-float-delayed"></div>
      
      <div class="relative z-10 flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
        <div class="space-y-4 max-w-2xl">
          <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-indigo-500/10 border border-indigo-500/20 text-indigo-300 text-xs font-medium">
            <span class="w-1.5 h-1.5 rounded-full bg-indigo-400 animate-pulse"></span>
            Community Hub
          </div>
          <h1 class="text-4xl md:text-5xl font-bold text-white tracking-tight">
            探索 <span class="text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 to-purple-400">无限灵感</span>
          </h1>
          <p class="text-slate-400 text-lg">
            这里汇聚了创作者们的奇思妙想。发现、分享、连接，让每一个 AI 角色都鲜活起来。
          </p>
        </div>
        
        <div class="flex items-center gap-2">
          <NuxtLink to="/community/presets" class="px-3 py-2 rounded-xl bg-white/5 border border-white/10 text-slate-200 hover:bg-white/10 transition-colors text-sm">
            预设市场
          </NuxtLink>
          <NuxtLink to="/community/new" class="btn-primary flex items-center gap-2 shadow-lg shadow-indigo-500/20 group">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5 group-hover:rotate-90 transition-transform">
              <path d="M10.75 4.75a.75.75 0 00-1.5 0v4.5h-4.5a.75.75 0 000 1.5h4.5v4.5a.75.75 0 001.5 0v-4.5h4.5a.75.75 0 000-1.5h-4.5v-4.5z" />
            </svg>
            发布动态
          </NuxtLink>
        </div>
      </div>

      <!-- Filter Tabs -->
      <div class="mt-10 flex flex-col gap-3 border-b border-white/5 pb-4 md:flex-row md:items-center md:justify-between relative z-10">
        <div class="flex items-center gap-2">
        <button 
          type="button"
          @click="setTab('latest')"
          :class="[
            'px-4 py-2 text-sm font-medium transition-colors border-b-2',
            currentTab === 'latest' ? 'text-white border-indigo-500' : 'text-slate-400 border-transparent hover:text-white'
          ]"
        >
          最新动态
        </button>
        <button 
          type="button"
          @click="setTab('hot')"
          :class="[
            'px-4 py-2 text-sm font-medium transition-colors border-b-2',
            currentTab === 'hot' ? 'text-white border-indigo-500' : 'text-slate-400 border-transparent hover:text-white'
          ]"
        >
          热门推荐
        </button>
        <button 
          type="button"
          @click="setTab('following')"
          :disabled="!auth.isAuthenticated.value"
          :class="[
            'px-4 py-2 text-sm font-medium transition-colors border-b-2',
            currentTab === 'following'
              ? 'text-white border-indigo-500'
              : !auth.isAuthenticated.value
                ? 'text-slate-600 border-transparent cursor-not-allowed'
                : 'text-slate-400 border-transparent hover:text-white'
          ]"
        >
          关注的人
        </button>
        </div>
        <form class="relative w-full md:w-80" @submit.prevent="onSearch">
          <input
            v-model.trim="searchQuery"
            type="search"
            placeholder="搜索标题或内容关键词"
            class="w-full rounded-full bg-white/5 border border-white/10 px-4 py-2 text-sm text-white placeholder:text-slate-500 focus:border-indigo-400 focus:ring-2 focus:ring-indigo-500/30"
          />
          <button
            type="submit"
            class="absolute right-2 top-1/2 -translate-y-1/2 rounded-full bg-indigo-500 px-3 py-1 text-xs font-medium text-white hover:bg-indigo-400"
          >
            搜索
          </button>
        </form>
      </div>
    </header>

    <!-- Masonry Grid -->
    <div class="columns-1 md:columns-2 lg:columns-3 gap-6 space-y-6">
      <CommunityPostCard 
        v-for="(post, index) in community.posts.value" 
        :key="post.id" 
        :post="post"
        class="break-inside-avoid animate-fade-up mb-6"
        :style="{ animationDelay: `${index * 100}ms` }"
      >
        <template #footer>
          <NuxtLink :to="`/community/post/${post.id}`" class="text-xs font-medium text-indigo-400 hover:text-indigo-300 transition-colors flex items-center gap-1">
            阅读详情
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-3 h-3">
              <path fill-rule="evenodd" d="M3 10a.75.75 0 01.75-.75h10.638L10.23 5.29a.75.75 0 111.04-1.08l5.5 5.25a.75.75 0 010 1.08l-5.5 5.25a.75.75 0 11-1.04-1.08l4.158-3.96H3.75A.75.75 0 013 10z" clip-rule="evenodd" />
            </svg>
          </NuxtLink>
        </template>
      </CommunityPostCard>
    </div>
  </section>
</template>

<script setup lang="ts">
const community = useCommunity()
const auth = useAuth()
const currentTab = ref('latest')
const searchQuery = ref('')
const loadingTab = ref(false)

const setTab = async (tab: string) => {
  if (loadingTab.value) return
  const isFollowingTab = tab === 'following'
  // 未登录时点击“关注的人”回退到“最新”
  const effectiveTab = isFollowingTab && !auth.isAuthenticated.value ? 'latest' : tab
  currentTab.value = effectiveTab
  const mapping: Record<string, { sort: string; filter: string }> = {
    latest: { sort: 'latest', filter: 'all' },
    hot: { sort: 'hot', filter: 'all' },
    following: { sort: 'latest', filter: 'following' },
  }
  const target = mapping[effectiveTab] || mapping.latest
  loadingTab.value = true
  try {
    await community.fetchFeed({ sort: target.sort, filter: target.filter, search: searchQuery.value })
  } finally {
    loadingTab.value = false
  }
}

const onSearch = () => setTab(currentTab.value)

onMounted(() => {
  setTab('latest')
})
</script>
