<template>
  <Dialog
    v-model:visible="visible"
    modal
    :closable="loading!"
    :draggable="false"
    :style="{ width, maxWidth }"
    :pt="{
      root: { class: 'border-white/6! bg-[#141414]! rounded-2xl! shadow-2xl!' },
      header: { class: 'bg-transparent! border-0! pb-2!' },
      content: { class: 'bg-transparent! px-6! pt-0! pb-2!' },
      footer: { class: 'bg-transparent! border-0!' },
      mask: { class: 'backdrop-blur-xs!' },
    }"
  >
    <template #header>
      <slot name="header" />
    </template>
    <slot />
    <template #footer>
      <slot name="footer" />
    </template>
  </Dialog>
</template>

<script setup lang="ts">

const visible = defineModel<boolean>('visible', { default: false })

withDefaults(defineProps<{
  loading?: boolean
  width?: string
  maxWidth?: string
}>(), {
  loading: false,
  width: '520px',
  maxWidth: '',
})
</script>
