<template>
  <section
    class="relative mb-10 overflow-hidden rounded-2xl"
    role="region"
    aria-label="Home hero"
  >
    <div class="absolute inset-0 bg-linear-to-br from-aurora-purple/20 via-black/60 to-aurora-blue/10" />
    <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_top_right,_var(--tw-gradient-stops))] from-spotify/10 via-transparent to-transparent" />

    <div class="relative z-10 flex flex-col gap-6 px-6 py-10 md:flex-row md:items-center md:justify-between md:px-10 md:py-14">
      <div class="max-w-xl space-y-4">
        <p class="text-[10px] font-bold tracking-[0.3em] text-white/40 uppercase">
          {{ periodLabel }}
        </p>
        <h1 class="text-3xl font-black text-white md:text-5xl" style="letter-spacing: -0.03em">
          {{ greeting }}
        </h1>
        <p class="text-base text-white/50 md:text-lg">
          {{ subtitle }}
        </p>
        <div class="flex flex-wrap gap-3 pt-2">
          <button
            type="button"
            class="inline-flex h-11 cursor-pointer items-center gap-2 rounded-full bg-spotify px-7 text-sm font-bold text-black transition hover:bg-spotify-hover hover:scale-105 active:scale-95 focus-visible:ring-2 focus-visible:ring-spotify focus-visible:ring-offset-2 focus-visible:outline-hidden"
            @click="$emit('explore')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4"><polygon points="5 3 19 12 5 21 5 3"/></svg>
            Explore Music
          </button>
          <button
            type="button"
            class="inline-flex h-11 cursor-pointer items-center gap-2 rounded-full border border-white/15 bg-white/4 px-6 text-sm font-medium text-white/80 backdrop-blur-xs transition hover:bg-white/10 hover:text-white active:scale-95 focus-visible:ring-2 focus-visible:ring-spotify focus-visible:ring-offset-2 focus-visible:outline-hidden"
            @click="$emit('discover')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4"><circle cx="12" cy="12" r="10"/><polygon points="16.24 7.76 14.12 14.12 7.76 16.24 9.88 9.88 16.24 7.76" fill="currentColor" stroke="none"/></svg>
            Discover
          </button>
        </div>
      </div>

      <div v-if="hasArt" class="hidden shrink-0 md:block">
        <div class="relative">
          <img
            :src="artSrc"
            alt=""
            class="h-40 w-40 rounded-2xl object-cover shadow-2xl ring-1 ring-white/10 md:h-48 md:w-48"
            style="transform: rotate(3deg);"
          />
          <div class="absolute -bottom-3 -right-3 flex h-10 w-10 items-center justify-center rounded-full bg-spotify text-black shadow-lg">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5"><polygon points="5 3 19 12 5 21 5 3"/></svg>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  greeting: string
  subtitle: string
  periodLabel: string
  artSrc?: string
  showExplore?: boolean
}>()

defineEmits<{
  explore: []
  discover: []
}>()

const hasArt = computed(() => !!props.artSrc)
</script>
