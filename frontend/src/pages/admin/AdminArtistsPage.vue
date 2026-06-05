<template>
  <div class="mx-auto w-full max-w-6xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Catalog"
      title="Artists"
      description="Create, edit, and manage artists in your catalog."
    />

    <AdminArtistsCard />
  </div>
</template>

<script setup lang="ts">
  import AdminSectionHeader from '@/components/admin/AdminSectionHeader.vue'
  import AdminArtistsCard from '@/components/admin/AdminArtistsCard.vue'
</script>
