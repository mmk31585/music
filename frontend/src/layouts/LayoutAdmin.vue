<template>
  <div class="min-h-screen bg-black text-white">
    <div class="grid min-h-screen grid-cols-1 lg:grid-cols-[280px_1fr]">
      <AdminSidebar class="hidden lg:flex" />

      <div class="flex min-h-screen flex-col">
        <AdminTopbar />
        <main class="flex-1 overflow-y-auto bg-linear-to-b from-[#151515] to-black">
          <router-view />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { AdminSidebar, AdminTopbar } from '@/components/admin'
</script>
