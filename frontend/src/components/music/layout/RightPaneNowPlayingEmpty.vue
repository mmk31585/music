<template>
  <div class="flex flex-col items-center justify-center gap-4 px-6 py-16 text-center">
    <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/5">
      <Music :size="28" class="text-slate-500" />
    </div>
    <div>
      <p class="text-sm font-bold text-slate-400">No track playing</p>
      <p class="mt-1 text-xs text-slate-500">Select a track to start listening</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Music } from 'lucide-vue-next'
</script>
