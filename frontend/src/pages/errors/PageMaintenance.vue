<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-black px-4 text-center">
    <!-- Aurora background effect -->
    <div class="pointer-events-none fixed inset-0 overflow-hidden" aria-hidden="true">
      <div class="aurora-spot absolute -top-40 left-1/2 h-[500px] w-[500px] -translate-x-1/2 rounded-full bg-amber-500/5 blur-[120px]" />
      <div class="aurora-spot absolute -bottom-40 right-0 h-[400px] w-[400px] rounded-full bg-orange-500/5 blur-[100px]" />
    </div>

    <div class="relative z-10 flex flex-col items-center gap-6">
      <!-- Icon -->
      <div class="flex h-20 w-20 items-center justify-center rounded-full bg-amber-500/10">
        <AlertTriangle aria-hidden="true" class="text-4xl text-amber-400"  />
      </div>

      <!-- Title -->
      <h1 class="text-3xl font-bold text-white md:text-4xl">
        Under Maintenance
      </h1>

      <!-- Description -->
      <p class="max-w-md text-base text-slate-400">
        We're currently performing scheduled maintenance to improve your experience.
        This should be brief — please check back shortly.
      </p>

      <!-- Status indicator -->
      <div class="flex items-center gap-2 rounded-full bg-white/5 px-4 py-2 text-sm text-slate-400">
        <span class="relative flex h-2 w-2">
          <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-amber-400 opacity-75" />
          <span class="relative inline-flex h-2 w-2 rounded-full bg-amber-500" />
        </span>
        Maintenance in progress
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { AlertTriangle } from 'lucide-vue-next'
// Simple maintenance page — no complex dependencies needed.
// Styled to work without the main app layout (layout-empty).
</script>
