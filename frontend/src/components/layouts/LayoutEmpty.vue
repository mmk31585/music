<template>
  <div>
    <!-- Skip link for keyboard users on error pages -->
    <a href="#main-content" class="skip-link">Skip to main content</a>
    <router-view v-slot="{ Component, route }">
      <main id="main-content" class="tw:mx-auto">
        <component :is="Component" :key="route.path" />
      </main>
    </router-view>
  </div>
</template>
