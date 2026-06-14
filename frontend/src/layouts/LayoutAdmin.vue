<template>
  <div class="min-h-screen bg-black text-white">
    <div class="flex min-h-screen">
      <!-- Desktop sidebar -->
      <AdminSidebar
        :collapsed="sidebarCollapsed"
        class="hidden lg:flex"
        @toggle="sidebarCollapsed = !sidebarCollapsed"
      />

      <!-- Mobile overlay backdrop -->
      <Transition name="fade">
        <div
          v-if="mobileOpen"
          class="fixed inset-0 z-40 bg-black/60 backdrop-blur-sm lg:hidden"
          @click="mobileOpen = false"
        />
      </Transition>

      <!-- Mobile sidebar drawer -->
      <Transition name="slide">
        <AdminSidebar
          v-if="mobileOpen"
          :collapsed="false"
          class="lg:!hidden"
          @close="mobileOpen = false"
        />
      </Transition>

      <!-- Main content area -->
      <div class="flex min-h-screen flex-1 flex-col">
        <AdminTopbar
          :collapsed="sidebarCollapsed"
          @toggle-mobile="mobileOpen = !mobileOpen"
          @toggle-collapse="sidebarCollapsed = !sidebarCollapsed"
        />
        <main class="flex-1 overflow-y-auto bg-linear-to-b from-[#151515] to-black">
          <router-view />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { AdminSidebar, AdminTopbar } from '@/components/admin'

const sidebarCollapsed = ref(false)
const mobileOpen = ref(false)
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.25s ease;
}
.slide-enter-from,
.slide-leave-to {
  transform: translateX(-100%);
}
</style>
