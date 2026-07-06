<template>
  <div
    class="auth-card w-full animate-reveal rounded-2xl border border-white/6 bg-black/40 p-8 shadow-2xl backdrop-blur-xl"
  >
    <div class="mb-8 text-center">
      <div class="mx-auto mb-5 flex h-14 w-14 items-center justify-center rounded-xl bg-spotify/10">
        <i class="pi pi-headphones text-2xl text-spotify" />
      </div>
      <h1 class="font-display text-2xl font-bold text-white tracking-tight">Welcome back</h1>
      <p class="mt-1.5 text-sm text-white/40">Log in to continue your music journey</p>
    </div>

    <form class="space-y-5" @submit.prevent="onSubmit">
      <!-- Email -->
      <div class="group">
        <label for="login-email" class="mb-1.5 block text-sm font-medium text-white/60 group-focus-within:text-spotify transition-colors duration-200">
          Email
        </label>
        <span class="relative block">
          <i class="pi pi-envelope absolute top-1/2 left-3 -translate-y-1/2 text-sm text-white/30" />
          <InputText
            id="login-email"
            v-model="form.email"
            type="email"
            autocomplete="email"
            placeholder="you@example.com"
            class="auth-input w-full pl-10"
            :class="{ 'ring-1 ring-red-500/50': errors.email }"
            :aria-describedby="errors.email ? 'login-email-error' : undefined"
            aria-required="true"
            @blur="validateField('email')"
          />
        </span>
        <Transition name="fade-slide">
          <small
            v-if="errors.email"
            id="login-email-error"
            class="mt-1 block text-xs text-red-400"
            role="alert"
          >
            {{ errors.email }}
          </small>
        </Transition>
      </div>

      <!-- Password -->
      <div class="group">
        <label for="login-password" class="mb-1.5 block text-sm font-medium text-white/60 group-focus-within:text-spotify transition-colors duration-200">
          Password
        </label>
        <span class="relative block">
          <i class="pi pi-lock absolute top-1/2 left-3 -translate-y-1/2 text-sm text-white/30 z-10" />
          <Password
            id="login-password"
            v-model="form.password"
            autocomplete="current-password"
            placeholder="Enter password"
            class="w-full"
            input-class="auth-input w-full pl-10"
            :feedback="false"
            toggle-mask
            :class="{ 'ring-1 ring-red-500/50': errors.password }"
            :aria-describedby="errors.password ? 'login-password-error' : undefined"
            aria-required="true"
            @blur="validateField('password')"
          />
        </span>
        <Transition name="fade-slide">
          <small
            v-if="errors.password"
            id="login-password-error"
            class="mt-1 block text-xs text-red-400"
            role="alert"
          >
            {{ errors.password }}
          </small>
        </Transition>
      </div>

      <!-- Forgot password -->
      <div class="flex justify-end -mt-3">
        <RouterLink
          to="/auth/forgot-password"
          class="text-xs text-white/40 transition hover:text-spotify"
        >
          Forgot password?
        </RouterLink>
      </div>

      <!-- Remember me -->
      <div class="flex items-center gap-2">
        <Checkbox
          id="login-remember"
          v-model="rememberMe"
          :binary="true"
          input-id="login-remember-input"
        />
        <label for="login-remember-input" class="cursor-pointer text-sm text-white/50 select-none hover:text-white/70 transition-colors">
          Remember me for 30 days
        </label>
      </div>

      <!-- API error -->
      <Transition name="fade-slide">
        <div
          v-if="apiError"
          class="flex items-start gap-2.5 rounded-xl border border-red-500/20 bg-red-500/8 px-4 py-3 text-sm text-red-300"
          role="alert"
        >
          <i class="pi pi-exclamation-circle mt-0.5 shrink-0 text-red-400" />
          <span>{{ apiError }}</span>
        </div>
      </Transition>

      <!-- Submit -->
      <Button
        type="submit"
        label="Log in"
        icon="pi pi-arrow-right"
        icon-pos="right"
        :loading="loading"
        class="auth-btn w-full border-0 bg-spotify text-black font-semibold hover:bg-spotify-hover transition-all duration-200"
      />

      <!-- Divider -->
      <div class="relative my-6">
        <div class="absolute inset-0 flex items-center">
          <div class="w-full border-t border-white/6" />
        </div>
        <div class="relative flex justify-center text-xs">
          <span class="bg-black/40 px-3 text-white/30">or</span>
        </div>
      </div>

      <!-- Guest link -->
      <RouterLink
        to="/"
        class="flex items-center justify-center gap-2 rounded-xl border border-white/8 px-4 py-2.5 text-sm text-white/50 transition-all duration-200 hover:border-white/20 hover:text-white/80 hover:bg-white/3"
      >
        <i class="pi pi-user text-xs" />
        Continue as guest
      </RouterLink>

      <!-- Switch to register -->
      <p class="text-center text-sm text-white/40">
        Don't have an account?
        <RouterLink to="/auth/register" class="font-medium text-spotify transition-colors duration-200 hover:text-spotify-hover">
          Create one
        </RouterLink>
      </p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { useLoginForm } from '@/composables/auth/useLoginForm.ts'

const { form, errors, apiError, loading, rememberMe, validateField, onSubmit } = useLoginForm()
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

:deep(.p-password-input) {
  /* Password component wraps inputtext, ensure the input gets the same styling */
}

:deep(.p-password .p-inputtext) {
  background: rgba(255, 255, 255, 0.04) !important;
  border: 1px solid rgba(255, 255, 255, 0.08) !important;
  color: white !important;
  border-radius: 12px !important;
  padding-top: 10px !important;
  padding-bottom: 10px !important;
  font-size: 0.875rem !important;
}

:deep(.p-password .p-inputtext:focus) {
  background: rgba(255, 255, 255, 0.06) !important;
  border-color: #1db954 !important;
  box-shadow: 0 0 0 3px rgba(29, 185, 84, 0.15) !important;
}

:deep(.p-password .p-input-icon) {
  color: rgba(255, 255, 255, 0.3) !important;
}

:deep(.p-password-panel) {
  background: #1a1a1a !important;
  border: 1px solid rgba(255, 255, 255, 0.08) !important;
  border-radius: 12px !important;
}

/* PrimeVue password toggle icon */
:deep(.p-password .p-password-toggle-icon) {
  color: rgba(255, 255, 255, 0.3) !important;
  right: 12px !important;
}

/* PrimeVue password input wrapper */
:deep(.p-password) {
  display: flex !important;
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
