<template>
  <section class="space-y-20 sm:space-y-32">
    <!-- Hero Section -->
    <header class="text-center space-y-8 py-10 sm:py-20 relative z-10">
      <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white border border-charcoal-300 text-sm font-medium text-charcoal-500 mb-4 animate-fade-up shadow-sm">
        <span class="relative flex h-2 w-2">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-accent-yellow opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-accent-yellow"></span>
        </span>
        The softest corner of the internet
      </div>
      
      <h1 class="text-5xl sm:text-7xl md:text-8xl font-display font-bold tracking-tight animate-fade-up delay-100 leading-[1.1] text-charcoal-900">
        <span class="block">随心创造</span>
        <span class="block text-charcoal-500">连接 AI 伙伴</span>
      </h1>
      
      <p class="max-w-2xl mx-auto text-lg sm:text-xl text-charcoal-500 animate-fade-up delay-200 leading-relaxed px-4">
        所有角色、对话、设定均由创作者自定义维护，通过统一的 Nuxt3 + Go 平台确保安全可控、可持续扩展。
      </p>
      
      <div class="flex flex-col sm:flex-row items-center justify-center gap-4 pt-8 animate-fade-up delay-300 px-4">
        <NuxtLink to="/search" class="btn-primary w-full sm:w-auto text-lg py-4 shadow-xl shadow-charcoal-900/10">
          <span class="i-ph-compass text-xl"></span>
          开始探索
        </NuxtLink>
        <NuxtLink to="/creator" class="btn-secondary w-full sm:w-auto text-lg py-4">
          <span class="i-ph-magic-wand text-xl"></span>
          成为创作者
        </NuxtLink>
      </div>
    </header>

    <!-- Featured Roles Section -->
    <section class="relative z-10">
      <div class="flex items-center justify-between mb-8 px-2">
        <h2 class="text-3xl font-display font-bold text-charcoal-900 flex items-center gap-3">
          <span class="w-1.5 h-8 rounded-full bg-accent-yellow"></span>
          精选角色
        </h2>
        <button class="group flex items-center gap-2 text-sm font-medium text-charcoal-500 hover:text-charcoal-900 transition-colors" @click="fetchFeatured">
          <span class="group-hover:rotate-180 transition-transform duration-500">
            <span class="i-ph-arrows-clockwise text-lg"></span>
          </span>
          换一批
        </button>
      </div>
      
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8">
        <RoleRoleCard 
          v-for="(role, index) in featured" 
          :key="role.id" 
          :role="role"
          class="animate-fade-up"
          :style="{ animationDelay: `${400 + index * 100}ms` }"
        />
      </div>
    </section>
  </section>
</template>

<script setup lang="ts">
const { featured, fetchFeatured } = useRole()
onMounted(fetchFeatured)
</script>
