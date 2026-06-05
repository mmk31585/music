<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section class="rounded-[2rem] border border-white/10 bg-white/[0.05] p-8 text-white">
      <p class="text-sm font-bold tracking-[0.35em] text-[#1db954] uppercase">Your music</p>

      <h1 class="mt-3 text-4xl font-black md:text-6xl">Library</h1>

      <p class="mt-4 max-w-2xl text-slate-300">
        Saved songs, favorite artists, albums, and playlists will appear here.
      </p>
    </section>

    <section class="mt-10 grid gap-5 md:grid-cols-3">
      <div
        v-for="item in items"
        :key="item.title"
        class="rounded-3xl border border-white/10 bg-black/20 p-6"
      >
        <div
          class="flex h-14 w-14 items-center justify-center rounded-2xl text-2xl"
          :class="item.iconClass"
        >
          <i :class="item.icon" />
        </div>

        <h2 class="mt-5 text-xl font-black text-white">
          {{ item.title }}
        </h2>

        <p class="mt-2 text-sm text-slate-400">
          {{ item.description }}
        </p>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const items = [
  {
    title: 'Liked Songs',
    description: 'Tracks you liked will be collected here.',
    icon: 'pi pi-heart',
    iconClass: 'bg-pink-500/20 text-pink-300',
  },
  {
    title: 'Saved Albums',
    description: 'Albums saved to your library.',
    icon: 'pi pi-images',
    iconClass: 'bg-purple-500/20 text-purple-300',
  },
  {
    title: 'Followed Artists',
    description: 'Artists you follow and listen to often.',
    icon: 'pi pi-users',
    iconClass: 'bg-[#1db954]/20 text-[#1db954]',
  },
]
</script>
