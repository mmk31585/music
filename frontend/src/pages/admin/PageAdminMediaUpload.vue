<template>
  <div class="mx-auto w-full max-w-5xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Media"
      title="Upload media"
      description="Send audio files and other media assets to the admin upload endpoint."
    />

    <UploadMediaCard />
  </div>
</template>

<script setup lang="ts">
import { AdminSectionHeader, UploadMediaCard } from '@/components/admin'
</script>
