<template>
  <header
    class="sticky top-0 z-30 border-b border-white/10 bg-black/60 px-4 py-3 backdrop-blur md:px-6"
  >
    <div class="flex items-center justify-between gap-4">
      <div class="flex items-center gap-3">
        <!-- Mobile hamburger -->
        <button
          class="flex h-9 w-9 items-center justify-center rounded-lg text-slate-400 transition hover:bg-white/5 hover:text-white lg:hidden"
          @click="$emit('toggleMobile')"
        >
          <i class="pi pi-bars text-base" />
        </button>

        <!-- Desktop collapse toggle -->
        <button
          class="hidden h-9 w-9 items-center justify-center rounded-lg text-slate-500 transition hover:bg-white/5 hover:text-white lg:flex"
          @click="$emit('toggleCollapse')"
        >
          <i class="pi text-sm" :class="collapsed ? 'pi-angle-right' : 'pi-angle-left'" />
        </button>

        <div>
          <p class="text-[10px] tracking-[0.2em] text-slate-600 uppercase">Administration</p>
          <h2 class="text-base font-semibold text-white">{{ pageTitle }}</h2>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <RouterLink
          to="/admin/media"
          class="inline-flex items-center gap-2 rounded-full bg-[#1db954] px-4 py-1.5 text-sm font-semibold text-black transition hover:opacity-90"
        >
          <i class="pi pi-upload text-xs" />
          <span class="hidden sm:inline">Upload</span>
        </RouterLink>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

defineProps<{
  collapsed: boolean
}>()

defineEmits<{
  toggleMobile: []
  toggleCollapse: []
}>()

const route = useRoute()

const pageTitle = computed(() => {
  return (route.meta?.title as string) || 'Admin'
})
</script>
