<template>
  <router-view v-slot="{ Component, route }">
    <div :key="route.path" class="tw:mx-auto">
      <component :is="Component" :key="route.path" />
    </div>
  </router-view>
</template>
