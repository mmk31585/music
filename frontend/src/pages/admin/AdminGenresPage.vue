<template>
  <div class="mx-auto w-full max-w-6xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Catalog"
      title="Genres"
      description="Create, edit, and manage music genres."
    />

    <AdminGenresCard />
  </div>
</template>

<script setup lang="ts">
import { AdminSectionHeader, AdminGenresCard } from '@/components/admin'
</script>
