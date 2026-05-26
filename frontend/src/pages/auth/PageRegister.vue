<template>
  <div class="grid min-h-screen grid-cols-1 lg:grid-cols-2">
    <AuthHeroPanel />
    <div class="flex items-center justify-center px-6 py-10">
      <div class="w-full max-w-md">
        <RegisterForm />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { AuthHeroPanel, RegisterForm } from '@/components/auth'
</script>
