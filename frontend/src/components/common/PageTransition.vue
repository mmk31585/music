<template>
  <RouterView v-slot="{ Component: RouteComponent }">
    <Transition name="page" mode="out-in">
      <KeepAlive :max="keepAliveMax" :include="keepAliveInclude">
        <component :is="RouteComponent" />
      </KeepAlive>
    </Transition>
  </RouterView>
</template>

<script setup lang="ts">
defineProps<{
  keepAliveMax?: number
  keepAliveInclude?: string[]
}>()
</script>
