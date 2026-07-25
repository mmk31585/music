<template>
  <Dialog
    v-model:visible="visible"
    modal
    :closable="true"
    :draggable="false"
    :style="{ width: '520px' }"
    :pt="{
      root: { class: 'border-white/6! bg-surface-raised! rounded-2xl! shadow-2xl!' },
      header: { class: 'bg-transparent! border-0! pb-2!' },
      content: { class: 'bg-transparent! px-6! pt-0! pb-2!' },
      footer: { class: 'bg-transparent! border-0!' },
      mask: { class: 'backdrop-blur-xs!' },
    }"
  >
    <template #header>
      <div class="flex items-center gap-3" lang="fa" dir="rtl">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-500/10">
          <Headphones aria-hidden="true" class="text-emerald-400"  />
        </div>
        <div>
          <h3 class="text-base font-semibold text-white">
            موجا رو کشف کن
          </h3>
          <p class="text-xs text-slate-500">
            برای ادامه شنیدن، یه حساب رایگان بساز. همین الان، همین جا.
          </p>
        </div>
      </div>
    </template>

    <div class="mt-6 flex flex-col items-center gap-4 text-center" lang="fa" dir="rtl">
      <div class="flex flex-col gap-3 w-full">
        <Button
          label="ثبت‌نام رایگان"
          class="w-full! rounded-xl! bg-emerald-500! text-black! hover:bg-emerald-400! py-3! text-sm! font-bold!"
          @click="router.push('/auth/register')"
        />
        <Button
          label="ورود"
          text
          class="w-full! rounded-xl! text-slate-400! hover:text-white! py-3! text-sm! font-bold!"
          @click="router.push('/auth/login')"
        />
      </div>

      <div class="mt-2 flex items-center gap-2 text-xs text-white/40">
        <span>{{ remainingPlays }} از {{ GUEST_PLAY_LIMIT }} پخش رایگان امروز باقی مونده</span>
      </div>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { Headphones } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { useGuestSession } from '@/composables/useGuestSession'

const router = useRouter()
const { remainingPlays, GUEST_PLAY_LIMIT } = useGuestSession()

const visible = defineModel<boolean>('visible', { default: false })
</script>
