<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <section class="relative z-10 min-h-screen bg-[#050505] py-30">
    <!-- Aurora background -->
    <div class="aurora-bg pointer-events-none fixed inset-0" aria-hidden="true">
      <div class="aurora-spot-1 -top-40 -left-40 bg-spotify/10" />
      <div class="aurora-spot-2 -top-60 -right-40 bg-aurora-blue/8" />
      <div class="aurora-spot-1 top-20 left-1/3 bg-electric-purple/8" style="animation-delay: -8s; width: 350px; height: 350px;" />
    </div>

    <div class="relative z-10 mx-auto w-full max-w-100 px-4 text-center">
      <div class="mb-6 flex justify-center">
        <div class="flex h-24 w-24 items-center justify-center rounded-full bg-white/4 ring-1 ring-white/6">
          <i aria-hidden="true" class="pi pi-exclamation-circle text-4xl text-slate-500" />
        </div>
      </div>

      <h2 class="mb-2 text-[80px] font-black leading-none text-white md:text-[100px]">
        404
      </h2>

      <h4 class="mb-3 text-[22px] font-bold text-white">
        Page not found
      </h4>

      <p class="mb-8 text-sm text-slate-400">
        The page you're looking for doesn't exist or has been moved.
      </p>

      <div class="flex items-center justify-center gap-3">
        <router-link
          to="/"
          class="inline-flex items-center gap-2 rounded-full bg-spotify px-8 py-3 text-sm font-bold text-black transition hover:bg-spotify-hover hover:scale-105"
        >
          <i aria-hidden="true" class="pi pi-home" />
          Go Home
        </router-link>

        <button
          class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/4 px-8 py-3 text-sm font-bold text-white/80 transition hover:border-white/30 hover:bg-white/8 hover:text-white"
          @click="router.go(-1)"
        >
          <i aria-hidden="true" class="pi pi-arrow-left" />
          Go Back
        </button>
      </div>
    </div>
  </section>
</template>
