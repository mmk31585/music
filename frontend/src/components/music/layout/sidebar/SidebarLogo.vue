<template>
  <RouterLink to="/" class="flex items-center gap-3 px-1 py-2">
    <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-accent text-accent-text">
      <Music :size="20" />
    </div>
    <div class="flex flex-col">
      <div class="text-base font-extrabold text-primary tracking-tight leading-tight">Muse</div>
      <div class="text-[10px] font-medium tracking-wider text-muted uppercase">Music</div>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import { Music } from 'lucide-vue-next'
</script>
