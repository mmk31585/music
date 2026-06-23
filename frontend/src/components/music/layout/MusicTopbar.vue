<template>
  <header
    class="sticky top-0 z-20 border-b border-white/10 bg-black/60 px-4 py-4 backdrop-blur-xs md:px-6"
  >
    <div class="flex items-center justify-between gap-4">
      <div class="flex items-center gap-3">
        <button type="button" aria-label="Menu" class="rounded-full bg-white/10 p-2 text-white lg:hidden">
          <i aria-hidden="true" class="pi pi-bars" />
        </button>
        <div>
          <p class="text-xs tracking-[0.3em] text-slate-500 uppercase">Music App</p>
          <h2 class="text-lg font-semibold">Browse</h2>
        </div>
      </div>

      <InputText v-model="query" placeholder="Search tracks, artists, albums..." aria-label="Search" class="hidden rounded-full bg-white/10 px-4 py-2 text-sm text-slate-300 transition hover:bg-white/15 md:block border-0! outline-hidden!" />
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const query = ref('')
</script>
