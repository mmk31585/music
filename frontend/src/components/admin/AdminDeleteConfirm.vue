<template>
  <Dialog
    v-model:visible="visible"
    modal
    :closable="deleting!"
    :draggable="false"
    :style="{ width: '420px' }"
    :pt="{
      root: { class: 'border-white/6! bg-[#141414]! rounded-2xl! shadow-2xl!' },
      header: { class: 'bg-transparent! border-0! pb-0!' },
      content: { class: 'bg-transparent! pt-0!' },
      footer: { class: 'bg-transparent! border-0!' },
      mask: { class: 'backdrop-blur-xs!' },
    }"
  >
    <template #header>
      <div class="flex items-center gap-3">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-red-500/10">
          <i aria-hidden="true" class="pi pi-trash text-red-400" />
        </div>
        <div>
          <h3 class="text-base font-semibold text-white">{{ title }}</h3>
        </div>
      </div>
    </template>

    <p class="mt-2 text-sm leading-relaxed text-slate-400">
      <slot>
        Are you sure you want to delete
        <span class="font-medium text-white">"{{ itemName }}"</span>? This action cannot be undone.
      </slot>
    </p>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          label="Cancel"
          text
          :disabled="deleting"
          class="text-slate-400! hover:text-white!"
          @click="visible = false"
        />
        <Button
          label="Delete"
          icon="pi pi-trash"
          :loading="deleting"
          severity="danger"
          class="rounded-xl! bg-red-500/10! text-red-400! ring-1! ring-red-500/20! hover:bg-red-500/20!"
          @click="$emit('confirm')"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup lang="ts">

const visible = defineModel<boolean>({ default: false })

defineProps<{
  title?: string
  itemName?: string
  deleting?: boolean
}>()

defineEmits<{
  confirm: []
}>()
</script>
