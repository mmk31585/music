<template>
  <div class="flex min-h-screen items-center justify-center bg-surface-base px-5 py-8">
    <div
      class="auth-card w-full max-w-105 animate-reveal rounded-2xl border border-white/6 bg-black/40 p-8 shadow-2xl backdrop-blur-xl"
    >
      <div class="mb-8 text-center">
        <div class="mx-auto mb-5 flex h-14 w-14 items-center justify-center rounded-xl bg-spotify/10">
          <i class="pi pi-lock-open text-2xl text-spotify" />
        </div>
        <h1 class="font-display text-2xl font-bold text-white tracking-tight">Forgot password</h1>
        <p class="mt-1.5 text-sm text-white/40">Enter your email and we'll send you a reset link</p>
      </div>

      <Transition name="fade-slide">
        <div
          v-if="sent"
          class="mb-6 flex items-start gap-2.5 rounded-xl border border-spotify/20 bg-spotify/8 px-4 py-3 text-sm text-spotify-300"
          role="alert"
        >
          <i class="pi pi-check-circle mt-0.5 shrink-0 text-spotify" />
          <span>If an account with that email exists, you'll receive a password reset link shortly.</span>
        </div>
      </Transition>

      <form v-if="!sent" class="space-y-5" @submit.prevent="onSubmit">
        <!-- Email -->
        <div class="group">
          <label for="reset-email" class="mb-1.5 block text-sm font-medium text-white/60 group-focus-within:text-spotify transition-colors duration-200">
            Email
          </label>
          <span class="relative block">
            <i class="pi pi-envelope absolute top-1/2 left-3 -translate-y-1/2 text-sm text-white/30" />
            <InputText
              id="reset-email"
              v-model="email"
              type="email"
              autocomplete="email"
              placeholder="you@example.com"
              class="auth-input w-full pl-10"
              :class="{ 'ring-1 ring-red-500/50': error }"
              aria-required="true"
            />
          </span>
          <Transition name="fade-slide">
            <small
              v-if="error"
              class="mt-1 block text-xs text-red-400"
              role="alert"
            >
              {{ error }}
            </small>
          </Transition>
        </div>

        <Button
          type="submit"
          label="Send reset link"
          icon="pi pi-send"
          icon-pos="right"
          :loading="loading"
          class="auth-btn w-full border-0 bg-spotify text-black font-semibold hover:bg-spotify-hover transition-all duration-200"
        />
      </form>

      <div class="mt-6 text-center text-sm text-white/40">
        <RouterLink to="/auth/login" class="font-medium text-spotify transition-colors duration-200 hover:text-spotify-hover">
          <i class="pi pi-arrow-left mr-1" />
          Back to login
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { axiosClient } from '@/plugins'

const email = ref('')
const error = ref('')
const loading = ref(false)
const sent = ref(false)

async function onSubmit() {
  error.value = ''

  if (!email.value.trim()) {
    error.value = 'Email is required'
    return
  }

  loading.value = true
  try {
    await axiosClient.post('/auth/forgot-password', { email: email.value.trim() })
    sent.value = true
  } catch (err: any) {
    // Always show success to prevent email enumeration
    sent.value = true
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-card {
  animation: card-enter 0.5s var(--ease-out-expo, cubic-bezier(0.19, 1, 0.22, 1)) both;
}

@keyframes card-enter {
  from {
    opacity: 0;
    transform: translateY(16px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

:deep(.auth-input) {
  background: rgba(255, 255, 255, 0.04) !important;
  border: 1px solid rgba(255, 255, 255, 0.08) !important;
  color: white !important;
  border-radius: 12px !important;
  padding-top: 10px !important;
  padding-bottom: 10px !important;
  font-size: 0.875rem !important;
  transition: all 0.2s ease !important;
}

:deep(.auth-input:focus) {
  background: rgba(255, 255, 255, 0.06) !important;
  border-color: #1db954 !important;
  box-shadow: 0 0 0 3px rgba(29, 185, 84, 0.15) !important;
  outline: none !important;
}

:deep(.auth-input::placeholder) {
  color: rgba(255, 255, 255, 0.4) !important;
}

:deep(.auth-btn) {
  border-radius: 12px !important;
  padding: 10px 0 !important;
  font-size: 0.9rem !important;
  box-shadow: 0 4px 16px rgba(29, 185, 84, 0.25) !important;
}

:deep(.auth-btn:hover) {
  box-shadow: 0 6px 24px rgba(29, 185, 84, 0.35) !important;
  transform: translateY(-1px);
}

.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.2s ease;
}

.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
