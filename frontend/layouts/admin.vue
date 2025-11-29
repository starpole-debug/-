<template>
  <div class="min-h-screen nebula-bg text-white">
    <div class="flex min-h-screen">
      <aside class="w-64 space-y-2 border-r border-white/5 bg-white/5 p-6 backdrop-blur-xl">
        <div class="mb-8 px-2">
          <h1 class="text-xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 to-cyan-300">
            Admin Console
          </h1>
        </div>
        <NuxtLink
          v-for="item in menu"
          :key="item.to"
          :to="item.to"
          class="block rounded-xl px-4 py-2.5 text-sm font-medium text-slate-400 hover:bg-white/5 hover:text-white transition-all"
          active-class="bg-indigo-500/10 text-indigo-300"
        >
          {{ item.label }}
        </NuxtLink>
      </aside>
      <main class="flex-1 p-8 overflow-y-auto min-h-screen">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
const menu = [
  { to: '/admin/models', label: '模型管理' },
  { to: '/admin/images', label: '绘图管理' },
  { to: '/admin/dictionary', label: '词典' },
  { to: '/admin/roles', label: '角色审核' },
  { to: '/admin/community', label: '社区守望' },
  { to: '/admin/users', label: '用户管理' },
  { to: '/admin/comments', label: '评论管理' },
  { to: '/admin/payouts', label: '提现申请' },
  { to: '/admin/messages', label: '站内信' },
]
</script>
