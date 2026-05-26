<template>
  <aside class="flex h-screen flex-col border-r border-white/10 bg-[#0a0a0a] p-4">
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

    <nav class="space-y-2">
      <RouterLink
        to="/admin"
        class="flex items-center gap-3 rounded-xl px-3 py-3 text-slate-300 transition hover:bg-white/5 hover:text-white"
      >
        <i class="pi pi-th-large" />
        <span>Dashboard</span>
      </RouterLink>

      <RouterLink
        to="/admin/media"
        class="flex items-center gap-3 rounded-xl px-3 py-3 text-slate-300 transition hover:bg-white/5 hover:text-white"
      >
        <i class="pi pi-upload" />
        <span>Media Upload</span>
      </RouterLink>

      <RouterLink
        to="/admin/catalog"
        class="flex items-center gap-3 rounded-xl px-3 py-3 text-slate-300 transition hover:bg-white/5 hover:text-white"
      >
        <i class="pi pi-database" />
        <span>Catalog</span>
      </RouterLink>
    </nav>

    <div class="mt-8 rounded-2xl bg-white/5 p-4">
      <p class="text-sm font-semibold text-white">Admin tools</p>
      <p class="mt-2 text-xs text-slate-400">Add media, curate tracks, and manage metadata.</p>
    </div>
    <RouterLink
      to="/admin/catalog/tracks"
      class="flex items-center gap-3 rounded-xl px-3 py-3 text-slate-300 transition hover:bg-white/5 hover:text-white"
    >
      <i class="pi pi-play-circle" />
      <span>Tracks</span>
    </RouterLink>

    <div class="mt-auto space-y-2">
      <RouterLink
        to="/"
        class="flex items-center gap-3 rounded-xl px-3 py-3 text-slate-300 transition hover:bg-white/5 hover:text-white"
      >
        <i class="pi pi-home" />
        <span>Back to app</span>
      </RouterLink>

      <button
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
import { useAuth } from '@/composables/auth/useAuth'

const { logout } = useAuth()

function handleLogout() {
  void logout()
}
</script>
