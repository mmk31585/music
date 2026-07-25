<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section
      class="rounded-2xl bg-linear-to-br from-spotify via-primary-700 to-surface-base p-8 text-white"
    >
      <p class="text-sm font-bold tracking-[0.35em] text-white/70 uppercase">Made for you</p>

      <h1 class="mt-3 text-4xl font-black md:text-6xl">Recommendations</h1>

      <p class="mt-4 max-w-2xl text-white/80">
        Discover popular, recent, and personalized tracks based on your listening experience.
      </p>
    </section>

    <section class="mt-10 grid gap-5 md:grid-cols-3">
      <RouterLink
        v-for="card in cards"
        :key="card.to"
        :to="card.to"
        class="group rounded-3xl border border-white/10 bg-white/5 p-6 transition hover:-translate-y-1 hover:bg-white/8"
      >
        <div
          class="flex h-14 w-14 items-center justify-center rounded-2xl text-2xl"
          :class="card.iconClass"
        >
          <i aria-hidden="true" :class="card.icon" />
        </div>

        <h2 class="mt-5 text-xl font-black text-white">
          {{ card.title }}
        </h2>

        <p class="mt-2 text-sm leading-6 text-slate-400">
          {{ card.description }}
        </p>
      </RouterLink>
    </section>

    <section class="mt-10 rounded-3xl border border-white/10 bg-black/20 p-6">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-sm font-semibold text-spotify">Quick picks</p>
          <h2 class="mt-1 text-2xl font-black text-white">Start exploring</h2>
        </div>

        <Headphones aria-hidden="true" class="text-3xl text-white/20"  />
      </div>

      <p class="mt-4 max-w-2xl text-sm leading-6 text-slate-400">
        Recommendation feeds can be connected to backend endpoints such as popular tracks, recently
        played, user history, liked tracks, or genre-based suggestions.
      </p>
    </section>
  </div>
</template>

<script setup lang="ts">
import { Headphones } from 'lucide-vue-next'
const cards = [
  {
    title: 'Popular Tracks',
    description: 'The most played tracks across the catalog.',
    to: '/recommendations/popular',
    icon: 'pi pi-chart-line',
    iconClass: 'bg-spotify/20 text-spotify',
  },
  {
    title: 'Best Tracks',
    description: 'Curated high-quality tracks for the best listening session.',
    to: '/recommendations/best',
    icon: 'pi pi-star',
    iconClass: 'bg-yellow-400/20 text-yellow-300',
  },
  {
    title: 'Recently Added',
    description: 'Fresh tracks recently added to the platform.',
    to: '/recommendations/recent',
    icon: 'pi pi-clock',
    iconClass: 'bg-sky-400/20 text-sky-300',
  },
]
</script>
