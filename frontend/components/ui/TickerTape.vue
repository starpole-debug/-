<template>
  <div 
    class="fixed left-0 right-0 pointer-events-none z-10 flex justify-center overflow-hidden"
    :style="{ 
      top: position === 'top' ? '15%' : 'auto',
      bottom: position === 'bottom' ? '15%' : 'auto',
      transform: `rotate(${angle}deg)`
    }"
  >
    <!-- Container for the stretch animation -->
    <div 
      class="bg-charcoal-900 text-bg-cream-100 h-10 sm:h-12 flex items-center whitespace-nowrap overflow-hidden shadow-xl origin-center animate-stretch-in w-[120vw]"
    >
      <!-- Marquee Content -->
      <div class="animate-marquee flex items-center gap-4 px-4 font-display font-medium tracking-wider text-sm sm:text-base uppercase">
        <span v-for="n in 20" :key="n" class="flex items-center gap-4">
          {{ text }} <span class="text-accent-yellow">•</span>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  text?: string
  angle?: number
  position?: 'top' | 'bottom'
}

withDefaults(defineProps<Props>(), {
  text: 'CREATE • SHARE • INSPIRE • VIBE • CONNECT • MILKYVERSE • FUTURE • COMMUNITY',
  angle: -2,
  position: 'top'
})
</script>
