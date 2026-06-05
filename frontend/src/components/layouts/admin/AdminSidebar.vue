<template>
  <aside class="flex h-screen w-72 flex-col border-r border-white/10 bg-[#0a0a0a] p-4">
    <!-- Brand -->
    <div class="mb-8 flex items-center gap-3 px-2 pt-2">
      <div
        class="flex h-10 w-10 items-center justify-center rounded-full bg-[#1db954] font-bold text-black"
      >
        A
      </div>
      <div>
        <h1 class="text-lg font-bold text-white">Admin Panel</h1>
        <p class="text-xs text-slate-400">Catalog & media control</p>
      </div>
    </div>

    <!-- Main nav -->
    <nav class="space-y-2">
      <RouterLink
        to="/admin"
        exact-active-class="bg-white/10 text-white"
        class="flex items-center gap-3 rounded-xl px-3 py-3 text-slate-300 transition hover:bg-white/5 hover:text-white"
      >
        <i class="pi pi-th-large" />
        <span>Dashboard</span>
      </RouterLink>



      <RouterLink
        to="/admin/catalog"
        exact-active-class="bg-white/10 text-white"
        class="flex items-center gap-3 rounded-xl px-3 py-3 text-slate-300 transition hover:bg-white/5 hover:text-white"
      >
        <i class="pi pi-database" />
        <span>Catalog</span>
      </RouterLink>

      <div class="pt-4">
        <p class="mb-2 px-3 text-xs font-semibold tracking-wide text-slate-500 uppercase">
          Manage catalog
        </p>

        <div class="space-y-2">
          <RouterLink
            to="/admin/tracks"
            active-class="bg-[#1db954]/15 text-[#1db954]"
            class="flex items-center gap-3 rounded-xl px-3 py-3 text-slate-300 transition hover:bg-white/5 hover:text-white"
          >
            <i class="pi pi-play-circle" />
            <span>Tracks</span>
          </RouterLink>

          <RouterLink
            to="/admin/artists"
            active-class="bg-[#1db954]/15 text-[#1db954]"
            class="flex items-center gap-3 rounded-xl px-3 py-3 text-slate-300 transition hover:bg-white/5 hover:text-white"
          >
            <i class="pi pi-users" />
            <span>Artists</span>
          </RouterLink>

          <RouterLink
            to="/admin/albums"
            active-class="bg-[#1db954]/15 text-[#1db954]"
            class="flex items-center gap-3 rounded-xl px-3 py-3 text-slate-300 transition hover:bg-white/5 hover:text-white"
          >
            <i class="pi pi-images" />
            <span>Albums</span>
          </RouterLink>

          <RouterLink
            to="/admin/genres"
            active-class="bg-[#1db954]/15 text-[#1db954]"
            class="flex items-center gap-3 rounded-xl px-3 py-3 text-slate-300 transition hover:bg-white/5 hover:text-white"
          >
            <i class="pi pi-tags" />
            <span>Genres</span>
          </RouterLink>
        </div>
      </div>
    </nav>

    <!-- Info card -->
    <div class="mt-8 rounded-2xl border border-white/10 bg-white/5 p-4">
      <p class="text-sm font-semibold text-white">Admin tools</p>
      <p class="mt-2 text-xs leading-5 text-slate-400">
        Add media, create tracks, manage artists, albums, genres, and catalog metadata.
      </p>
    </div>

    <!-- Bottom actions -->
    <div class="mt-auto space-y-2">
      <RouterLink
        to="/"
        class="flex items-center gap-3 rounded-xl px-3 py-3 text-slate-300 transition hover:bg-white/5 hover:text-white"
      >
        <i class="pi pi-home" />
        <span>Back to app</span>
      </RouterLink>

      <button
        type="button"
        class="flex w-full items-center gap-3 rounded-xl px-3 py-3 text-left text-slate-300 transition hover:bg-white/5 hover:text-white"
        @click="handleLogout"
      >
        <i class="pi pi-sign-out" />
        <span>Logout</span>
      </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useAuth } from '@/composables/auth/useAuth.ts'

const { logout } = useAuth()

function handleLogout() {
  void logout()
}
</script>
