<template>
  <div class="min-h-screen text-charcoal-900 relative overflow-x-hidden selection:bg-accent-yellow-soft/50 font-sans">
    <!-- Ticker Tapes Background -->
    <UiTickerTape text="CREATE • SHARE • INSPIRE • VIBE • CONNECT • MILKYVERSE • FUTURE • COMMUNITY" :angle="-2" position="top" />
    <UiTickerTape text="DESIGN • FUTURE • MILKYVERSE • COMMUNITY • SOFT • COZY • PLAYFUL • WARM" :angle="1.5" position="bottom" />

    <!-- Floating Navigation Container -->
    <div class="floating-nav">
      <NavMainNav />
    </div>

    <!-- Main Content -->
    <main class="pt-24 sm:pt-32 pb-20 sm:pb-24 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto relative z-20 w-full min-h-[calc(100vh-6rem)] flex flex-col">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
const auth = useAuth()
onMounted(() => auth.fetchProfile())
</script>
