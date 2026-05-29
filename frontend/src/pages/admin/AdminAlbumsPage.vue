<template>
  <div class="mx-auto w-full max-w-6xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Catalog"
      title="Albums"
      description="Create, edit, and manage albums in your catalog."
    />

    <AdminAlbumsCard />
  </div>
</template>

<script setup lang="ts">
import { AdminSectionHeader, AdminAlbumsCard } from '@/components/admin'
</script>
