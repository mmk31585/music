<template>
  <div class="space-y-10">
    <div class="h-56 animate-pulse rounded-2xl bg-white/[3%]" />

    <div class="space-y-4">
      <div class="h-4 w-32 animate-pulse rounded bg-white/[4%]" />
      <div class="h-6 w-48 animate-pulse rounded bg-white/[4%]" />
      <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        <div v-for="i in 6" :key="i" class="h-16 animate-pulse rounded-xl bg-white/[3%]" />
      </div>
    </div>

    <div class="space-y-4">
      <div class="h-4 w-28 animate-pulse rounded bg-white/[4%]" />
      <div class="h-6 w-40 animate-pulse rounded bg-white/[4%]" />
      <div class="grid grid-cols-5 gap-3">
        <div v-for="i in 5" :key="i" class="aspect-square animate-pulse rounded-2xl bg-white/[3%]" />
      </div>
    </div>

    <div class="space-y-4">
      <div class="h-4 w-24 animate-pulse rounded bg-white/[4%]" />
      <div class="h-6 w-36 animate-pulse rounded bg-white/[4%]" />
      <div class="space-y-2">
        <div v-for="i in 5" :key="i" class="h-12 animate-pulse rounded-xl bg-white/[3%]" />
      </div>
    </div>
  </div>
</template>
