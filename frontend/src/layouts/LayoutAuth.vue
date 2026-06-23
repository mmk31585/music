<template>
  <div class="min-h-screen bg-surface-raised text-white" :dir="dir">
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { useRTL } from '@/composables'
const { dir } = useRTL()
</script>
