<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      title="Settings"
      eyebrow="Admin"
      description="Platform configuration and preferences"
    />

    <AdminEmptyState
      icon="pi pi-cog"
      title="Settings"
      description="Platform settings are coming soon. Configure site-wide preferences here."
    />
  </div>
</template>

<script setup lang="ts">
import { AdminSectionHeader, AdminEmptyState } from '@/components/admin'
</script>
