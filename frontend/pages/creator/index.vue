<template>
  <section class="space-y-8">
    <header class="card-soft p-8 animate-fade-up">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-2xl font-bold text-charcoal-900">创作者控制台</h1>
          <p class="text-sm text-charcoal-500 mt-1">实时监控您的角色表现与收益数据。</p>
        </div>
        <div class="p-2 rounded-full bg-accent-pink/10 border border-accent-pink/20">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-6 h-6 text-accent-pink">
            <path d="M15.5 2A1.5 1.5 0 0014 3.5v13a1.5 1.5 0 001.5 1.5h1a1.5 1.5 0 001.5-1.5v-13A1.5 1.5 0 0016.5 2h-1zM9.5 6A1.5 1.5 0 008 7.5v9A1.5 1.5 0 009.5 18h1a1.5 1.5 0 001.5-1.5v-9A1.5 1.5 0 0010.5 6h-1zM3.5 10A1.5 1.5 0 002 11.5v5A1.5 1.5 0 003.5 18h1A1.5 1.5 0 006 16.5v-5A1.5 1.5 0 004.5 10h-1z" />
          </svg>
        </div>
      </div>
      
      <div class="grid gap-4 sm:grid-cols-3">
        <div class="bg-bg-cream-100 border border-charcoal-100 rounded-2xl p-4">
          <div class="text-sm text-charcoal-500 mb-1">角色总数</div>
          <div class="text-2xl font-bold text-charcoal-900">{{ stats?.total_roles || 0 }}</div>
        </div>
        <div class="bg-bg-cream-100 border border-charcoal-100 rounded-2xl p-4">
          <div class="text-sm text-charcoal-500 mb-1">草稿箱</div>
          <div class="text-2xl font-bold text-charcoal-900">{{ stats?.draft_roles || 0 }}</div>
        </div>
        <div class="bg-accent-pink/10 border border-accent-pink/20 rounded-2xl p-4">
          <div class="text-sm text-accent-pink mb-1">已发布</div>
          <div class="text-2xl font-bold text-accent-pink">{{ stats?.published || 0 }}</div>
        </div>
      </div>
    </header>

    <section class="card-soft p-8 animate-fade-up delay-100">
      <div class="flex items-center justify-between gap-3 mb-4">
        <div class="flex items-center gap-2">
          <span class="w-1.5 h-6 rounded-full bg-accent-yellow"></span>
          <div>
            <h2 class="text-lg font-bold text-charcoal-900">虚拟钱包</h2>
            <p class="text-xs text-charcoal-500">只展示钱包余额与流水，收入事件不单列。</p>
          </div>
        </div>
        <div class="flex gap-2 text-right">
          <div class="px-3 py-2 rounded-xl bg-bg-cream-100 border border-charcoal-100">
            <div class="text-[12px] text-charcoal-400">可用余额</div>
            <div class="text-lg font-semibold text-charcoal-900">{{ stats?.wallet?.available_balance ?? 0 }}</div>
          </div>
          <div class="px-3 py-2 rounded-xl bg-bg-cream-100 border border-charcoal-100">
            <div class="text-[12px] text-charcoal-400">冻结</div>
            <div class="text-lg font-semibold text-charcoal-900">{{ stats?.wallet?.frozen_balance ?? 0 }}</div>
          </div>
          <div class="px-3 py-2 rounded-xl bg-bg-cream-100 border border-charcoal-100">
            <div class="text-[12px] text-charcoal-400">累计收入</div>
            <div class="text-lg font-semibold text-charcoal-900">{{ stats?.wallet?.total_earned ?? 0 }}</div>
          </div>
        </div>
      </div>

      <div class="rounded-2xl border border-charcoal-100 bg-bg-cream-100/50">
        <div class="flex items-center justify-between px-4 py-3 border-b border-charcoal-100 text-sm text-charcoal-600">
          <span>收入流水</span>
          <span class="text-xs text-charcoal-400">最多 {{ pageSize }} 条/页，可滚动</span>
        </div>
        <div class="max-h-80 overflow-y-auto custom-scroll">
          <div v-if="!walletEvents.length" class="text-center py-10 text-charcoal-400 text-sm">
            暂无流水，等待第一笔收入。
          </div>
          <ul v-else class="divide-y divide-charcoal-100">
            <li
              v-for="event in pagedEvents"
              :key="event.id"
              class="flex items-center justify-between px-4 py-3 hover:bg-bg-cream-200 transition-colors"
            >
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-full bg-accent-yellow/15 border border-accent-yellow/30 flex items-center justify-center text-accent-yellow-dark">
                  ¥
                </div>
                <div>
                  <div class="text-charcoal-900 text-sm">+{{ event.amount }}</div>
                  <div class="text-[12px] text-charcoal-500">{{ event.event_type || '收入' }}</div>
                </div>
              </div>
              <div class="text-right">
                <div class="text-[12px] text-charcoal-400">{{ new Date(event.created_at).toLocaleString() }}</div>
                <div class="text-[11px] text-status-success">{{ event.status || 'confirmed' }}</div>
              </div>
            </li>
          </ul>
        </div>
        <div v-if="totalPages > 1" class="flex items-center justify-between px-4 py-3 border-t border-charcoal-100 text-sm text-charcoal-600">
          <button class="px-3 py-1 rounded-lg bg-white border border-charcoal-200 disabled:opacity-40 hover:bg-charcoal-50" :disabled="page === 1" @click="page--">
            上一页
          </button>
          <div class="text-xs">第 {{ page }} / {{ totalPages }} 页</div>
          <button class="px-3 py-1 rounded-lg bg-white border border-charcoal-200 disabled:opacity-40 hover:bg-charcoal-50" :disabled="page === totalPages" @click="page++">
            下一页
          </button>
        </div>
      </div>
    </section>

    <section class="card-soft p-8 animate-fade-up delay-150">
      <div class="flex items-center justify-between flex-wrap gap-4 mb-4">
        <div>
          <h2 class="text-lg font-bold text-charcoal-900">角色卡 / 世界书 / 预设 构建器</h2>
          <p class="text-sm text-charcoal-500 mt-1">在可视化工具里编排 Persona、世界书、示例对话，并导出 nebula/ST 友好格式。</p>
        </div>
        <div class="flex gap-3">
          <NuxtLink to="/creator/editor/new" class="btn-primary">立即构建</NuxtLink>
          <NuxtLink to="/creator/roles" class="btn-secondary">我的角色</NuxtLink>
          <NuxtLink to="/creator/presets" class="btn-secondary">预设工坊</NuxtLink>
        </div>
      </div>
      <div class="grid gap-4 md:grid-cols-2">
        <div class="rounded-2xl border border-charcoal-100 bg-bg-cream-100 p-4">
          <h3 class="text-charcoal-900 font-semibold">可视化编辑</h3>
          <p class="text-sm text-charcoal-500 mt-1">拆分 Persona、世界书、示例对话、模型提示、安全片段，所见即所得。</p>
        </div>
        <div class="rounded-2xl border border-charcoal-100 bg-bg-cream-100 p-4">
          <h3 class="text-charcoal-900 font-semibold">导出/兼容</h3>
          <p class="text-sm text-charcoal-500 mt-1">一键复制 nebula-1 JSON，或生成 ST 友好版，便于跨平台创作。</p>
        </div>
      </div>
    </section>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'

const { dashboard: stats, loadDashboard } = useCreatorStats()
const page = ref(1)
const pageSize = 5
const walletEvents = computed(() => stats.value?.recent_events || [])
const totalPages = computed(() => Math.max(1, Math.ceil(walletEvents.value.length / pageSize)))
const pagedEvents = computed(() => {
  const start = (page.value - 1) * pageSize
  return walletEvents.value.slice(start, start + pageSize)
})
watch(walletEvents, () => { page.value = 1 })
onMounted(loadDashboard)
</script>
