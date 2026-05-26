<template>
  <div class="rounded-3xl border border-white/10 bg-white/5 p-8 shadow-2xl backdrop-blur-xl">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-white">Log in</h1>
      <p class="mt-2 text-sm text-slate-400">Continue to your music app</p>
    </div>

    <form class="space-y-5" @submit.prevent="onSubmit">
      <div>
        <label class="mb-2 block text-sm font-medium text-slate-300">Email</label>
        <InputText v-model="form.email" type="email" placeholder="you@example.com" class="w-full" />
        <small v-if="errors.email" class="mt-1 block text-red-400">
          {{ errors.email }}
        </small>
      </div>

      <div>
        <label class="mb-2 block text-sm font-medium text-slate-300">Password</label>
        <Password
          v-model="form.password"
          placeholder="Enter password"
          class="w-full"
          input-class="w-full"
          :feedback="false"
          toggle-mask
        />
        <small v-if="errors.password" class="mt-1 block text-red-400">
          {{ errors.password }}
        </small>
      </div>

      <Button
        type="submit"
        label="Log in"
        icon="pi pi-sign-in"
        :loading="loading"
        class="w-full border-0 bg-[#1db954] text-black"
      />

      <p class="text-center text-sm text-slate-400">
        Don’t have an account?
        <RouterLink to="/auth/register" class="font-medium text-[#1db954] hover:underline">
          Create one
        </RouterLink>
      </p>
    </form>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import { useLoginForm } from '@/composables/auth/useLoginForm.ts'

const { form, errors, loading, onSubmit } = useLoginForm()
</script>
