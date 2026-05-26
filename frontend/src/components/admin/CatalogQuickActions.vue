<template>
  <div class="rounded-3xl border border-white/10 bg-white/5 p-6 shadow-lg">
    <div class="mb-6">
      <h2 class="text-xl font-bold text-white">Catalog quick actions</h2>
      <p class="mt-2 text-sm text-slate-400">
        Shortcuts for managing tracks, artists, albums, and genres.
      </p>
    </div>

    <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <RouterLink
        to="/admin/media"
        class="rounded-2xl border border-white/10 bg-black/20 p-5 transition hover:border-[#1db954]/40 hover:bg-white/5"
      >
        <div
          class="mb-3 flex h-12 w-12 items-center justify-center rounded-2xl bg-[#1db954]/15 text-[#1db954]"
        >
          <i class="pi pi-upload" />
        </div>
        <h3 class="font-semibold text-white">Upload media</h3>
        <p class="mt-2 text-sm text-slate-400">Upload audio and media assets.</p>
      </RouterLink>

      <button
        class="rounded-2xl border border-white/10 bg-black/20 p-5 text-left transition hover:border-[#1db954]/40 hover:bg-white/5"
      >
        <div
          class="mb-3 flex h-12 w-12 items-center justify-center rounded-2xl bg-[#1db954]/15 text-[#1db954]"
        >
          <i class="pi pi-plus-circle" />
        </div>
        <h3 class="font-semibold text-white">Create track</h3>
        <p class="mt-2 text-sm text-slate-400">Add new tracks to the catalog.</p>
      </button>

      <button
        class="rounded-2xl border border-white/10 bg-black/20 p-5 text-left transition hover:border-[#1db954]/40 hover:bg-white/5"
      >
        <div
          class="mb-3 flex h-12 w-12 items-center justify-center rounded-2xl bg-[#1db954]/15 text-[#1db954]"
        >
          <i class="pi pi-users" />
        </div>
        <h3 class="font-semibold text-white">Manage artists</h3>
        <p class="mt-2 text-sm text-slate-400">Create or edit artist records.</p>
      </button>

      <button
        class="rounded-2xl border border-white/10 bg-black/20 p-5 text-left transition hover:border-[#1db954]/40 hover:bg-white/5"
      >
        <div
          class="mb-3 flex h-12 w-12 items-center justify-center rounded-2xl bg-[#1db954]/15 text-[#1db954]"
        >
          <i class="pi pi-book" />
        </div>
        <h3 class="font-semibold text-white">Manage albums</h3>
        <p class="mt-2 text-sm text-slate-400">Update album metadata and covers.</p>
      </button>
    </div>
  </div>
</template>
