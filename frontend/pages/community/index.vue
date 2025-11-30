<template>
  <section class="community-shell">
    <div class="ribbon ribbon-top">
      <div class="ribbon-track">
        <span v-for="n in 10" :key="`top-${n}`">CREATE • SHARE • INSPIRE • VIBE • CONNECT • MILKYVERSE</span>
      </div>
    </div>
    <div class="ribbon ribbon-bottom">
      <div class="ribbon-track">
        <span v-for="n in 10" :key="`bottom-${n}`">DESIGN • FUTURE • COMMUNITY • SOFT • COZY • PLAYFUL</span>
      </div>
    </div>

    <div class="max-w-6xl mx-auto relative z-10 space-y-8">
      <header class="hero-card animate-fade-up">
        <div class="hero-main">
          <div class="hero-copy">
            <p class="eyebrow">MilkyVerse • Community Hub</p>
            <h1>
              探索 <span>无限灵感</span>
            </h1>
            <p class="lede">
              分享你的灵感、故事与角色。让温柔柔软的社区氛围，帮助每一个 AI 角色被看见、被连接。
            </p>
            <div class="cta-row">
              <NuxtLink to="/community/new" class="cta-btn">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M10.75 4.75a.75.75 0 00-1.5 0v4.5h-4.5a.75.75 0 000 1.5h4.5v4.5a.75.75 0 001.5 0v-4.5h4.5a.75.75 0 000-1.5h-4.5v-4.5z" />
                </svg>
                发布动态
              </NuxtLink>
              <NuxtLink to="/community/presets" class="ghost-btn">预设市场</NuxtLink>
            </div>
          </div>

          <div class="hero-side">
            <div class="stat">
              <p>社区动态</p>
              <h3>{{ postCount }}</h3>
              <span>实时更新的创作灵感</span>
            </div>
            <div class="stat">
              <p>活跃话题</p>
              <h3>{{ topicCount }}</h3>
              <span>大家正在讨论的主题</span>
            </div>
          </div>
        </div>

        <div class="hero-bottom">
          <div class="left-stack">
            <form class="search-box" @submit.prevent="onSearch">
              <input
                v-model.trim="searchQuery"
                type="search"
                placeholder="搜索 vibes、角色、关键词..."
              />
              <button type="submit">搜索</button>
            </form>

            <div class="toy-card" @mousemove="onPlayMove" @mouseleave="onPlayLeave">
              <div class="toy-header">
                <div class="toy-icon">✨</div>
                <div>
                  <p class="toy-title">情绪软糖球</p>
                  <p class="toy-sub">拖动光斑，看看它的色彩变化</p>
                </div>
                <button class="toy-btn" type="button" @click.stop="randomizeOrb">换个颜色</button>
              </div>
              <div class="toy-stage">
                <div class="toy-ring"></div>
                <div
                  class="toy-orb"
                  :style="{
                    background: orbColor,
                    transform: `translate(${playTilt.x}px, ${playTilt.y}px) scale(${playTilt.scale})`
                  }"
                ></div>
              </div>
            </div>
          </div>

          <div class="side-stack">
            <div class="trending-card">
              <div class="trending-title">
                <span class="icon">↗</span>
                Trending Now
              </div>
              <ul>
                <li v-for="(item, idx) in trendingTopics" :key="idx">
                  <div class="label"># {{ item.label }}</div>
                  <div class="meta">{{ item.hint }}</div>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </header>

      <div class="feed-controls">
        <div class="tabs">
          <button
            type="button"
            class="tab-btn"
            :class="{ active: currentTab === 'latest' }"
            @click="setTab('latest')"
          >
            最新动态
          </button>
          <button
            type="button"
            class="tab-btn"
            :class="{ active: currentTab === 'hot' }"
            @click="setTab('hot')"
          >
            热门推荐
          </button>
        </div>
        <p class="hint">柔和、温暖、好玩的社区氛围，欢迎加入 ✨</p>
      </div>

      <div class="columns-1 md:columns-2 lg:columns-3 gap-6 space-y-6">
        <CommunityPostCard 
          v-for="(post, index) in community.posts.value" 
          :key="post.id" 
          :post="post"
          class="break-inside-avoid animate-fade-up mb-6"
          :style="{ animationDelay: `${index * 100}ms` }"
        >
          <template #footer>
            <NuxtLink :to="`/community/post/${post.id}`" class="text-xs font-medium text-charcoal-500 hover:text-accent-yellow transition-colors flex items-center gap-1">
              阅读详情
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-3 h-3">
                <path fill-rule="evenodd" d="M3 10a.75.75 0 01.75-.75h10.638L10.23 5.29a.75.75 0 111.04-1.08l5.5 5.25a.75.75 0 010 1.08l-5.5 5.25a.75.75 0 11-1.04-1.08l4.158-3.96H3.75A.75.75 0 013 10z" clip-rule="evenodd" />
              </svg>
            </NuxtLink>
          </template>
        </CommunityPostCard>
      </div>
    </div>

  </section>
</template>

<script setup lang="ts">
const community = useCommunity()
const currentTab = ref('latest')
const searchQuery = ref('')
const loadingTab = ref(false)
const postCount = computed(() => community.posts.value?.length || 0)
const topicCount = computed(() => {
  const set = new Set<string>()
  community.posts.value.forEach((post) => {
    ;(post.topic_ids || []).forEach((t) => set.add(String(t)))
  })
  return set.size || Math.max(3, Math.min(12, postCount.value))
})
const trendingTopics = computed(() => {
  const seen = new Set<string>()
  const derived = community.posts.value
    .map((p, idx) => ({
      label: p.title || `热帖 #${idx + 1}`,
      hint: p.content?.slice(0, 42) || '社区热议话题',
    }))
    .filter((item) => {
      const key = item.label.trim()
      if (!key || seen.has(key)) return false
      seen.add(key)
      return true
    })
    .slice(0, 4)
  const fallback = [
    { label: 'MilkyVerse', hint: '柔软又有设计感的社区' },
    { label: 'AI 角色', hint: '分享你的角色设定与故事' },
    { label: '创意灵感', hint: '灵感、草稿、随想都欢迎' },
    { label: '社区活动', hint: '一起参与讨论与连接' },
  ]
  return derived.length ? derived : fallback
})

const orbPalettes = [
  'linear-gradient(135deg, #FFD93D, #FF9A62)',
  'linear-gradient(135deg, #A78BFA, #F472B6)',
  'linear-gradient(135deg, #34D399, #60A5FA)',
  'linear-gradient(135deg, #F59E0B, #EF4444)',
  'linear-gradient(135deg, #22D3EE, #3B82F6)',
]
const orbColor = ref(orbPalettes[0])
const playTilt = reactive({ x: 0, y: 0, scale: 1 })
const onPlayMove = (event: MouseEvent) => {
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
  const offsetX = event.clientX - rect.left - rect.width / 2
  const offsetY = event.clientY - rect.top - rect.height / 2
  playTilt.x = offsetX * 0.04
  playTilt.y = offsetY * 0.04
  playTilt.scale = 1.04
}
const onPlayLeave = () => {
  playTilt.x = 0
  playTilt.y = 0
  playTilt.scale = 1
}
const randomizeOrb = () => {
  const pool = orbPalettes.filter((p) => p !== orbColor.value)
  orbColor.value = pool[Math.floor(Math.random() * pool.length)] || orbColor.value
}

const setTab = async (tab: string) => {
  if (loadingTab.value) return
  const mapping: Record<string, { sort: string; filter: string }> = {
    latest: { sort: 'latest', filter: 'all' },
    hot: { sort: 'hot', filter: 'all' },
  }
  const target = mapping[tab] || mapping.latest
  currentTab.value = mapping[tab] ? tab : 'latest'
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

<style scoped>
.community-shell {
  position: relative;
  min-height: 100vh;
  padding: 2.5rem 1rem 5rem;
  overflow: hidden;
  background: radial-gradient(circle at 10% 10%, #fff7d1, transparent 32%),
              radial-gradient(circle at 82% 20%, #ffe5f0, transparent 30%),
              linear-gradient(180deg, #fffdf5 0%, #fff6db 55%, #ffeef6 100%);
}

.ribbon {
  position: absolute;
  left: -12%;
  right: -12%;
  height: 42px;
  background: #0f1115;
  color: #ffd93d;
  box-shadow: 0 18px 40px -28px rgba(0, 0, 0, 0.6);
  overflow: hidden;
  z-index: 0;
}
.ribbon-top {
  top: 86px;
  transform: rotate(-6deg);
}
.ribbon-bottom {
  top: 152px;
  transform: rotate(-3deg);
}
.ribbon-track {
  display: flex;
  gap: 32px;
  white-space: nowrap;
  font-weight: 700;
  letter-spacing: 0.08em;
  animation: ribbon-scroll 14s linear infinite;
  font-size: 12px;
  padding: 0 18px;
}
.ribbon-track span {
  opacity: 0.85;
}

.hero-card {
  position: relative;
  border-radius: 28px;
  padding: 2.5rem;
  background: rgba(255, 255, 255, 0.82);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.6);
  box-shadow: 0 22px 60px -28px rgba(26, 27, 30, 0.35);
  overflow: hidden;
}
.hero-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at 18% 18%, rgba(255, 217, 61, 0.18), transparent 42%),
              radial-gradient(circle at 80% 18%, rgba(255, 229, 240, 0.32), transparent 42%),
              linear-gradient(120deg, rgba(255, 255, 255, 0.4), rgba(255, 250, 245, 0.5));
  z-index: 0;
}
.hero-main {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}
@media (min-width: 768px) {
  .hero-main {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }
}
.hero-copy h1 {
  font-family: 'Outfit', 'Quicksand', system-ui, sans-serif;
  font-size: clamp(2.6rem, 3vw, 3.6rem);
  line-height: 1.1;
  font-weight: 800;
  color: #1a1b1e;
  margin-bottom: 0.5rem;
}
.hero-copy h1 span {
  background: linear-gradient(120deg, #ffd93d 0%, #ffb347 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
.eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 999px;
  background: rgba(255, 217, 61, 0.16);
  color: #3a3b40;
  font-size: 12px;
  font-weight: 700;
  border: 1px solid rgba(255, 217, 61, 0.4);
}
.lede {
  color: #4d4e54;
  font-size: 16px;
  max-width: 40rem;
}
.cta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 12px;
}
.cta-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  border-radius: 999px;
  padding: 12px 18px;
  background: #1a1b1e;
  color: #fffdf5;
  font-weight: 700;
  box-shadow: 0 14px 32px -16px rgba(26, 27, 30, 0.65);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.cta-btn svg {
  width: 20px;
  height: 20px;
}
.cta-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 16px 36px -18px rgba(26, 27, 30, 0.7);
}
.ghost-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  border-radius: 999px;
  padding: 12px 16px;
  background: #ffffff;
  color: #3a3b40;
  font-weight: 700;
  border: 1px solid #e6e7eb;
  box-shadow: 0 10px 18px -16px rgba(26, 27, 30, 0.4);
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}
.ghost-btn:hover {
  transform: translateY(-1px);
  border-color: #ffd93d;
  box-shadow: 0 16px 28px -18px rgba(26, 27, 30, 0.4);
}

.hero-side {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  min-width: 240px;
}
.stat {
  padding: 12px 14px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.82);
  border: 1px solid rgba(255, 255, 255, 0.7);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.7), 0 10px 24px -18px rgba(26, 27, 30, 0.35);
}
.stat p {
  margin: 0;
  color: #6e7077;
  font-size: 13px;
}
.stat h3 {
  margin: 4px 0;
  font-size: 28px;
  font-weight: 800;
  color: #1a1b1e;
}
.stat span {
  color: #9a9ca2;
  font-size: 12px;
}

.hero-bottom {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
  margin-top: 22px;
  align-items: start;
}
@media (min-width: 960px) {
  .hero-bottom {
    grid-template-columns: 2fr 1fr;
  }
}
.left-stack {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.search-box {
  display: flex;
  align-items: center;
  gap: 12px;
  border-radius: 999px;
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.75);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.75), 0 12px 24px -22px rgba(26, 27, 30, 0.35);
  align-self: start;
  max-width: 620px;
}
.search-box input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  color: #3a3b40;
  font-size: 15px;
  padding: 8px 4px;
}
.search-box input::placeholder {
  color: #c4c5ca;
}
.search-box button {
  border: none;
  border-radius: 999px;
  padding: 10px 16px;
  background: #1a1b1e;
  color: #fffdf5;
  font-weight: 700;
  font-size: 14px;
  cursor: pointer;
  transition: transform 0.15s ease, background 0.15s ease;
}
.search-box button:hover {
  background: #25262b;
  transform: translateY(-1px);
}

.trending-card {
  border-radius: 18px;
  padding: 14px 16px;
  background: rgba(255, 255, 255, 0.88);
  border: 1px solid rgba(255, 255, 255, 0.75);
  box-shadow: 0 18px 45px -32px rgba(17, 17, 17, 0.6);
}
.trending-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 800;
  color: #1a1b1e;
  margin-bottom: 10px;
}
.trending-title .icon {
  color: #ffd93d;
}
.trending-card ul {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 10px;
}
.trending-card li {
  padding: 10px 12px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.7);
}
.trending-card .label {
  font-weight: 700;
  color: #1a1b1e;
}
.trending-card .meta {
  font-size: 12px;
  color: #6e7077;
  margin-top: 2px;
}

.side-stack {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.toy-card {
  position: relative;
  border-radius: 18px;
  padding: 12px 14px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.75);
  box-shadow: 0 16px 40px -30px rgba(17, 17, 17, 0.45);
  overflow: hidden;
}
.toy-header {
  display: flex;
  align-items: center;
  gap: 10px;
  justify-content: space-between;
}
.toy-icon {
  width: 32px;
  height: 32px;
  border-radius: 10px;
  background: #fff7d1;
  display: grid;
  place-items: center;
  font-size: 14px;
}
.toy-title {
  margin: 0;
  font-weight: 800;
  color: #1a1b1e;
}
.toy-sub {
  margin: 0;
  font-size: 12px;
  color: #6e7077;
}
.toy-btn {
  border: none;
  border-radius: 999px;
  padding: 8px 12px;
  background: #1a1b1e;
  color: #fffdf5;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: transform 0.15s ease, background 0.15s ease;
}
.toy-btn:hover {
  background: #25262b;
  transform: translateY(-1px);
}
.toy-stage {
  position: relative;
  margin-top: 12px;
  height: 130px;
  border-radius: 14px;
  background: linear-gradient(135deg, rgba(255, 217, 61, 0.18), rgba(255, 229, 240, 0.22));
  overflow: hidden;
}
.toy-ring {
  position: absolute;
  inset: 16px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.7), transparent 60%);
  filter: blur(18px);
}
.toy-orb {
  position: absolute;
  top: 40%;
  left: 42%;
  width: 78px;
  height: 78px;
  border-radius: 50%;
  filter: drop-shadow(0 12px 25px rgba(0, 0, 0, 0.18));
  transition: transform 0.18s ease, background 0.3s ease;
}

.feed-controls {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 14px 16px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.75);
  box-shadow: 0 12px 28px -22px rgba(26, 27, 30, 0.45);
}
@media (min-width: 768px) {
  .feed-controls {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }
}
.tabs {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
.tab-btn {
  position: relative;
  padding: 10px 14px;
  border-radius: 999px;
  border: 1px solid transparent;
  background: rgba(255, 255, 255, 0.85);
  color: #6e7077;
  font-weight: 700;
  transition: all 0.2s ease;
}
.tab-btn::after {
  content: '';
  position: absolute;
  left: 18px;
  right: 18px;
  bottom: 8px;
  height: 3px;
  border-radius: 999px;
  background: linear-gradient(90deg, #ffd93d, #ffb347);
  opacity: 0;
  transform: translateY(6px);
  transition: all 0.2s ease;
}
.tab-btn.active {
  color: #1a1b1e;
  border-color: rgba(255, 217, 61, 0.8);
  box-shadow: 0 10px 22px -18px rgba(26, 27, 30, 0.4);
}
.tab-btn.active::after {
  opacity: 1;
  transform: translateY(0);
}
.tab-btn.disabled {
  color: #c4c5ca;
  cursor: not-allowed;
  border-color: rgba(0, 0, 0, 0.02);
}
.hint {
  color: #6e7077;
  font-size: 13px;
}

@keyframes ribbon-scroll {
  0% {
    transform: translateX(0);
  }
  100% {
    transform: translateX(-50%);
  }
}
</style>
