<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section
      class="relative overflow-hidden rounded-[2rem] border border-white/10 bg-[#121212] p-8 text-white"
    >
      <div class="absolute -top-20 -left-20 h-72 w-72 rounded-full bg-purple-500/20 blur-3xl" />

      <div class="relative z-10">
        <p class="text-sm font-bold tracking-[0.35em] text-purple-300 uppercase">Collections</p>

        <h1 class="mt-3 text-4xl font-black md:text-6xl">Playlists</h1>

        <p class="mt-4 max-w-2xl text-slate-300">
          Create, organize, and play custom music collections.
        </p>
      </div>
    </section>

    <section class="mt-10">
      <button
        type="button"
        class="rounded-full bg-[#1db954] px-6 py-3 text-sm font-bold text-black transition hover:scale-105 hover:bg-[#1ed760]"
      >
        <i class="pi pi-plus mr-2" />
        Create Playlist
      </button>

      <div class="mt-6 rounded-3xl border border-white/10 bg-black/20 px-6 py-16 text-center">
        <div
          class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-white/10 text-2xl text-white"
        >
          <i class="pi pi-list" />
        </div>

        <h2 class="mt-5 text-xl font-black text-white">No playlists yet</h2>

        <p class="mt-2 text-sm text-slate-400">Start by creating your first playlist.</p>
      </div>
    </section>
  </div>
</template>
