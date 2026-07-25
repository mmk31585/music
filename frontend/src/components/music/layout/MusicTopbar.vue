<template>
  <header
    class="sticky top-0 z-20 border-b border-border-default bg-surface-raised/70 px-4 py-4 backdrop-blur-xs md:px-6"
  >
    <div class="flex items-center justify-between gap-4">
      <div class="flex items-center gap-3">
        <button type="button" aria-label="Menu" class="rounded-full bg-surface-active p-2 text-primary lg:hidden">
          <Menu aria-hidden="true" class=""  />
        </button>
        <div>
          <p class="text-xs tracking-[0.3em] text-muted uppercase">Music App</p>
          <h2 class="text-lg font-semibold text-primary">Browse</h2>
        </div>
      </div>

      <InputText v-model="query" placeholder="Search tracks, artists, albums..." aria-label="Search" class="hidden rounded-full bg-surface-overlay px-4 py-2 text-sm text-secondary transition hover:bg-surface-active md:block border-0! outline-hidden! focus-visible:ring-2 focus-visible:ring-accent/50" />
    </div>
  </header>
</template>

<script setup lang="ts">
import { Menu } from 'lucide-vue-next'
import { ref } from 'vue'

const query = ref('')
</script>
