<script setup lang="ts">
import AuthHeroPanel from '@/components/auth/AuthHeroPanel.vue'
import LoginForm from '@/components/auth/LoginForm.vue'
</script>
<template>
  <div class="grid min-h-screen grid-cols-1 lg:grid-cols-2">
    <AuthHeroPanel />
    <div class="flex items-center justify-center px-6 py-10">
      <div class="w-full max-w-md">
        <LoginForm />
      </div>
    </div>
  </div>
</template>


