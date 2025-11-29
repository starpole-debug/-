<template>
  <section class="min-h-screen pb-20">
    <!-- Header -->
    <header class="sticky top-0 z-30 bg-gray-900/80 backdrop-blur-md border-b border-white/5 -mx-4 px-4 py-4 mb-8">
      <div class="max-w-7xl mx-auto flex items-center justify-between">
        <div>
          <h1 class="text-xl font-bold text-white flex items-center gap-2">
            <span v-if="isSearchMode" class="text-indigo-400">搜索</span>
            <span v-else>发现</span>
          </h1>
        </div>
        
        <div class="flex items-center gap-3">
          <transition name="fade" mode="out-in">
            <button 
              v-if="!isSearchMode"
              @click="isSearchMode = true"
              class="group flex items-center gap-2 px-4 py-2 rounded-full bg-white/5 hover:bg-white/10 border border-white/5 hover:border-white/10 transition-all"
            >
              <div class="i-ph-magnifying-glass text-slate-400 group-hover:text-white" />
              <span class="text-sm text-slate-500 group-hover:text-slate-300">搜索角色...</span>
            </button>
            
            <button 
              v-else
              @click="isSearchMode = false"
              class="p-2 rounded-full hover:bg-white/10 text-slate-400 hover:text-white transition-colors"
            >
              <div class="i-ph-x text-xl" />
            </button>
          </transition>
        </div>
      </div>
    </header>

    <!-- Search Mode -->
    <div v-if="isSearchMode" class="animate-fade-in max-w-7xl mx-auto">
      <div class="relative mb-8">
        <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
          <div class="i-ph-magnifying-glass text-xl text-indigo-500" />
        </div>
        <input
          v-model.trim="keyword"
          placeholder="输入角色名称、描述或标签..."
          class="w-full bg-slate-800/50 border border-slate-700 rounded-2xl pl-12 pr-4 py-4 text-base focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all shadow-lg shadow-black/20"
          auto-focus
        />
      </div>

      <div v-if="isLoading" class="py-20 text-center text-slate-400">
        <div class="i-svg-spinners-90-ring-with-bg text-3xl mx-auto mb-4 text-indigo-500" />
        <p>正在搜索全宇宙...</p>
      </div>

      <div v-else-if="filtered.length === 0" class="py-32 text-center text-slate-500">
        <div class="i-ph-ghost text-6xl mx-auto mb-6 opacity-30" />
        <p class="text-lg">没有找到相关角色</p>
        <p class="text-sm mt-2">换个关键词试试？</p>
      </div>

      <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        <RoleRoleCard v-for="role in filtered" :key="role.id" :role="role" />
      </div>
    </div>

    <!-- Dashboard Mode -->
    <div v-else class="space-y-12 max-w-7xl mx-auto animate-fade-in">
      <!-- Hot Roles -->
      <section>
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-2">
            <div class="p-1.5 rounded-lg bg-orange-500/10 text-orange-500">
              <div class="i-ph-fire-fill text-xl" />
            </div>
            <h2 class="text-xl font-bold text-white">热门推荐</h2>
          </div>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <RoleRoleCard 
            v-for="(role, index) in featured.slice(0, 6)" 
            :key="role.id" 
            :role="role"
            class="animate-fade-up hover:-translate-y-1 transition-transform duration-300"
            :style="{ animationDelay: `${index * 50}ms` }"
          />
        </div>
      </section>

      <!-- Latest Roles -->
      <section>
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-2">
            <div class="p-1.5 rounded-lg bg-indigo-500/10 text-indigo-500">
              <div class="i-ph-sparkle-fill text-xl" />
            </div>
            <h2 class="text-xl font-bold text-white">最新发布</h2>
          </div>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          <RoleRoleCard 
            v-for="(role, index) in roles" 
            :key="role.id" 
            :role="role" 
            class="animate-fade-up hover:-translate-y-1 transition-transform duration-300"
            :style="{ animationDelay: `${index * 50}ms` }"
          />
        </div>
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import RoleRoleCard from '@/components/Role/RoleCard.vue'

definePageMeta({
  middleware: ['auth'],
})

const { roles, featured, fetchRoles, fetchFeatured, isLoading, hasLoaded } = useRole()
const isSearchMode = ref(false)
const keyword = ref('')

// Ensure data is loaded
onMounted(async () => {
  await Promise.all([fetchRoles(true), fetchFeatured()])
})

const filtered = computed(() => {
  const k = keyword.value.trim().toLowerCase()
  if (!k) return roles.value
  return roles.value.filter(r => 
    r.name.toLowerCase().includes(k) || 
    r.description.toLowerCase().includes(k) ||
    r.tags?.some(t => t.toLowerCase().includes(k))
  )
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
