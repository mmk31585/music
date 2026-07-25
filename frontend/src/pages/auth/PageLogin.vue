<script setup lang="ts">
import { onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import AuthHeroPanel from '@/components/auth/AuthHeroPanel.vue'
import LoginForm from '@/components/auth/LoginForm.vue'

const toast = useToast()

onMounted(() => {
  const reason = sessionStorage.getItem('auth_redirect_reason')
  if (reason) {
    sessionStorage.removeItem('auth_redirect_reason')
    toast.add({
      severity: 'info',
      summary: reason,
      life: 5000,
    })
  }
})
</script>

<template>
  <div class="flex min-h-screen bg-surface-base">
    <div class="hidden flex-1 lg:block">
      <AuthHeroPanel />
    </div>
    <div class="flex flex-1 items-center justify-center px-5 py-8">
      <div class="w-full max-w-105">
        <LoginForm />
      </div>
    </div>
  </div>
</template>
