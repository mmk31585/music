<template>
  <header
    class="sticky top-0 z-20 border-b border-white/10 bg-black/60 px-4 py-4 backdrop-blur md:px-6"
  >
    <div class="flex items-center justify-between gap-4">
      <div>
        <p class="text-xs tracking-[0.3em] text-slate-500 uppercase">Administration</p>
        <h2 class="text-lg font-semibold text-white">Manage your music platform</h2>
      </div>

      <div class="flex items-center gap-3">
        <RouterLink
          to="/admin/media"
          class="rounded-full bg-[#1db954] px-4 py-2 text-sm font-medium text-black transition hover:opacity-90"
        >
          <i class="pi pi-upload mr-2" />
          Upload
        </RouterLink>
      </div>
    </div>
  </header>
</template>
