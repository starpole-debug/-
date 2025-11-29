<template>
  <section class="space-y-16">
    <!-- Hero Section -->
    <header class="text-center space-y-6 py-10 relative">
      <div class="inline-block px-4 py-1.5 rounded-full bg-white/5 border border-white/10 text-xs font-medium text-indigo-300 mb-4 animate-fade-up">
        AI Persona Playground
      </div>
      <h1 class="text-5xl md:text-7xl font-bold tracking-tight animate-fade-up delay-100">
        <span class="hero-text">随心创造</span><br />
        <span class="text-white">连接 AI 伙伴</span>
      </h1>
      <p class="max-w-2xl mx-auto text-lg text-slate-400 animate-fade-up delay-200 leading-relaxed">
        所有角色、对话、设定均由创作者自定义维护，通过统一的 Nuxt3 + Go 平台确保安全可控、可持续扩展。
      </p>
      <div class="flex items-center justify-center gap-4 pt-4 animate-fade-up delay-300">
        <NuxtLink to="/search" class="btn-primary">开始探索</NuxtLink>
        <NuxtLink to="/creator" class="btn-secondary">成为创作者</NuxtLink>
      </div>
    </header>

    <!-- Bento Grid Section -->
    <section>
      <div class="flex items-center justify-between mb-8 px-2">
        <h2 class="text-2xl font-bold text-white flex items-center gap-2">
          <span class="w-2 h-8 rounded-full bg-indigo-500"></span>
          精选角色
        </h2>
        <button class="text-sm font-medium text-indigo-400 hover:text-indigo-300 transition-colors" @click="fetchFeatured">
          换一批
        </button>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
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
