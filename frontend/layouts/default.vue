<template>
  <div class="min-h-screen nebula-bg text-slate-200 relative overflow-x-hidden selection:bg-indigo-500/30">
    <!-- Floating Navigation Container -->
    <div class="floating-nav">
      <NavMainNav />
    </div>

    <!-- Main Content -->
    <main class="pt-32 pb-20 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto relative z-10">
      <slot />
    </main>

    <!-- Ambient Background Elements -->
    <div class="fixed top-0 left-0 w-full h-full overflow-hidden pointer-events-none z-0">
      <div class="absolute top-[-10%] left-[-10%] w-[50vw] h-[50vw] bg-indigo-600/10 rounded-full blur-[120px] animate-float"></div>
      <div class="absolute bottom-[-10%] right-[-10%] w-[50vw] h-[50vw] bg-blue-600/10 rounded-full blur-[120px] animate-float-delayed"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
const auth = useAuth()
onMounted(() => auth.fetchProfile())
</script>
