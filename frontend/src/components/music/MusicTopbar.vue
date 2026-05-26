<template>
  <header
    class="sticky top-0 z-20 border-b border-white/10 bg-black/60 px-4 py-4 backdrop-blur md:px-6"
  >
    <div class="flex items-center justify-between gap-4">
      <div class="flex items-center gap-3">
        <button class="rounded-full bg-white/10 p-2 text-white lg:hidden">
          <i class="pi pi-bars" />
        </button>
        <div>
          <p class="text-xs tracking-[0.3em] text-slate-500 uppercase">Music App</p>
          <h2 class="text-lg font-semibold">Browse</h2>
        </div>
      </div>

      <RouterLink
        to="/search"
        class="hidden rounded-full bg-white/10 px-4 py-2 text-sm text-slate-300 transition hover:bg-white/15 md:block"
      >
        <i class="pi pi-search mr-2" />
        Search tracks, artists, albums...
      </RouterLink>
    </div>
  </header>
</template>
