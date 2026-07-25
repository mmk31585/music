<template>
  <header
    class="sticky top-0 z-30 border-b border-border-default bg-surface-raised/70 px-4 py-3 backdrop-blur-xs md:px-6"
  >
    <div class="flex items-center justify-between gap-4">
      <div class="flex items-center gap-3">
        <button
          class="flex h-9 w-9 items-center justify-center rounded-lg text-secondary transition hover:bg-surface-active hover:text-primary lg:hidden"
          @click="$emit('toggleMobile')"
        >
          <Menu aria-hidden="true" class="text-base"  />
        </button>

        <button
          class="hidden h-9 w-9 items-center justify-center rounded-lg text-tertiary transition-all hover:bg-surface-active hover:text-primary active:scale-95 lg:flex"
          :title="collapsed ? 'Expand sidebar' : 'Collapse sidebar'"
          @click="$emit('toggleCollapse')"
        >
          <i
            aria-hidden="true"
            class="pi text-sm transition-transform duration-200"
            :class="collapsed ? 'pi-angle-right' : 'pi-angle-left'"
          />
        </button>

        <div class="hidden h-5 w-px bg-border-default lg:block" />

        <div>
          <p class="text-[10px] tracking-[0.2em] text-muted uppercase">Administration</p>
          <h2 class="text-base font-semibold text-primary">{{ pageTitle }}</h2>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <RouterLink
          to="/"
          class="inline-flex items-center gap-2 rounded-full border border-border-default px-4 py-1.5 text-sm font-medium text-secondary transition hover:bg-surface-active hover:text-primary"
        >
          <Home aria-hidden="true" class="text-xs"  />
          <span class="hidden sm:inline">Back to app</span>
        </RouterLink>
        <RouterLink
          to="/admin/media"
          class="inline-flex items-center gap-2 rounded-full bg-accent px-4 py-1.5 text-sm font-semibold text-accent-text transition hover:opacity-90"
        >
          <Upload aria-hidden="true" class="text-xs"  />
          <span class="hidden sm:inline">Upload</span>
        </RouterLink>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { Home, Menu, Upload } from 'lucide-vue-next'
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
