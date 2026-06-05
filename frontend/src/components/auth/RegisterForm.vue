<template>
  <div class="rounded-3xl border border-white/10 bg-white/5 p-8 shadow-2xl backdrop-blur-xl">
    <div class="mb-8 text-center">
      <h1 class="text-3xl font-bold text-white">Create account</h1>
      <p class="mt-2 text-sm text-slate-400">Start your music journey</p>
    </div>

    <form class="space-y-5" @submit.prevent="onSubmit">
      <div>
        <label class="mb-2 block text-sm font-medium text-slate-300">Display name</label>
        <InputText
          v-model.trim="form.displayName"
          placeholder="Your name"
          class="w-full"
          :invalid="!!errors.displayName"
        />
        <small v-if="errors.displayName" class="mt-1 block text-red-400">
          {{ errors.displayName }}
        </small>
      </div>

      <div>
        <label class="mb-2 block text-sm font-medium text-slate-300">Username</label>
        <InputText
          v-model.trim="form.username"
          placeholder="your_username"
          class="w-full"
          :invalid="!!errors.username"
        />
        <small v-if="errors.username" class="mt-1 block text-red-400">
          {{ errors.username }}
        </small>
      </div>

      <div>
        <label class="mb-2 block text-sm font-medium text-slate-300">Email</label>
        <InputText
          v-model.trim="form.email"
          type="email"
          placeholder="you@example.com"
          class="w-full"
          :invalid="!!errors.email"
        />
        <small v-if="errors.email" class="mt-1 block text-red-400">
          {{ errors.email }}
        </small>
      </div>

      <div>
        <label class="mb-2 block text-sm font-medium text-slate-300">Password</label>
        <Password
          v-model="form.password"
          placeholder="Create password"
          class="w-full"
          input-class="w-full"
          :feedback="true"
          toggle-mask
          :invalid="!!errors.password"
        />
        <small v-if="errors.password" class="mt-1 block text-red-400">
          {{ errors.password }}
        </small>
      </div>

      <div
        v-if="apiError"
        class="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300"
      >
        {{ apiError }}
      </div>

      <Button
        type="submit"
        label="Create account"
        icon="pi pi-user-plus"
        :loading="loading"
        :disabled="loading"
        class="w-full border-0 bg-[#1db954] text-black hover:bg-[#1ed760]"
      />

      <p class="text-center text-sm text-slate-400">
        Already have an account?
        <RouterLink to="/auth/login" class="font-medium text-[#1db954] hover:underline">
          Log in
        </RouterLink>
      </p>
    </form>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import { useRegisterForm } from '@/composables/auth/useRegisterForm'

const { form, errors, apiError, loading, onSubmit } = useRegisterForm()
</script>
