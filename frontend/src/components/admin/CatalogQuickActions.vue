<template>
  <div class="rounded-2xl border border-white/[0.06] bg-white/[0.03] p-6">
    <h3 class="text-sm font-semibold uppercase tracking-wider text-slate-400">Quick actions</h3>

    <div class="mt-4 grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
      <RouterLink
        v-for="action in actions"
        :key="action.to"
        :to="action.to"
        class="group flex items-center gap-3 rounded-xl border border-white/[0.06] bg-white/[0.02] px-4 py-3.5 transition-all duration-200 hover:border-emerald-500/20 hover:bg-emerald-500/5"
      >
        <div
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-white/5 transition-colors group-hover:bg-emerald-500/10"
        >
          <i :class="[action.icon, 'text-sm text-slate-400 group-hover:text-emerald-400 transition-colors']" />
        </div>
        <div class="min-w-0">
          <p class="text-sm font-medium text-white">{{ action.label }}</p>
          <p class="truncate text-xs text-slate-500">{{ action.description }}</p>
        </div>
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
const actions = [
  {
    label: 'Manage Tracks',
    description: 'View & edit tracks',
    icon: 'pi pi-play-circle',
    to: '/admin/tracks',
  },
  {
    label: 'Manage Artists',
    description: 'View & edit artists',
    icon: 'pi pi-users',
    to: '/admin/artists',
  },
  {
    label: 'Manage Albums',
    description: 'View & edit albums',
    icon: 'pi pi-book',
    to: '/admin/albums',
  },
  {
    label: 'Manage Genres',
    description: 'View & edit genres',
    icon: 'pi pi-tags',
    to: '/admin/genres',
  },
]
</script>
