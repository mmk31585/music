<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import PageProgressBar from '@/components/widgets/page-progressbar/PageProgressBar.vue'
import { registerRouter, registerToast } from '@/composables'
import { useToast } from 'primevue/usetoast'

const route = useRoute()

const layout = computed(() => {
  return route?.meta?.layout ?? 'layout-empty'
})

// Injection for request wrapper working
registerToast(useToast())
registerRouter(useRouter())

</script>

<template>
  <PageProgressBar />
  <Toast />

  <component :is="layout">
    <router-view />
  </component>
</template>
