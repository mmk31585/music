<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="fixed inset-0 z-50 flex items-center justify-center bg-bg-overlay/60 backdrop-blur-xs"
      @click.self="emit('close')"
    >
      <div class="glass-strong mx-4 w-full max-w-lg rounded-2xl p-8">
        <h2 class="mb-6 text-xl font-bold text-primary">بحث جدید</h2>

        <div class="space-y-4">
          <input
            v-model="title"
            type="text"
            placeholder="عنوان بحث"
            aria-label="عنوان بحث"
            autofocus
            class="w-full rounded-xl border border-border-default bg-surface-overlay px-4 py-3 text-sm text-primary placeholder:text-muted outline-hidden transition focus:border-border-strong"
            dir="rtl"
          />
          <textarea
            v-model="body"
            placeholder="متن بحث..."
            rows="5"
            aria-label="متن بحث"
            class="w-full rounded-xl border border-border-default bg-surface-overlay px-4 py-3 text-sm text-primary placeholder:text-muted outline-hidden transition focus:border-border-strong"
            dir="rtl"
          />
          <div
            v-if="errorMessage"
            class="rounded-lg bg-danger-subtle px-4 py-2.5 text-xs text-danger"
          >
            {{ errorMessage }}
          </div>
        </div>

        <div class="mt-6 flex gap-3">
          <button
            class="flex-1 rounded-xl bg-surface-overlay py-3 text-sm font-medium text-primary/50 transition hover:bg-surface-active"
            @click="emit('close')"
          >
            انصراف
          </button>
          <button
            class="flex-1 rounded-xl bg-accent py-3 text-sm font-bold text-black transition hover:bg-accent/90 disabled:opacity-40"
            :disabled="!title.trim() || !body.trim() || creating"
            @click="handleCreate"
          >
            {{ creating ? '...' : 'انتشار بحث' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useSocialApi } from '@/services/api/social'
import { useAppToast } from '@/composables/useAppToast'

const props = defineProps<{
  visible: boolean
  clubId: string
}>()
const emit = defineEmits<{
  close: []
  created: [discussionId: string]
}>()

const socialApi = useSocialApi()
const toast = useAppToast()

const title = ref('')
const body = ref('')
const creating = ref(false)
const errorMessage = ref('')

async function handleCreate() {
  if (!title.value.trim() || !body.value.trim() || creating.value) return
  creating.value = true
  errorMessage.value = ''
  try {
    const d = await socialApi.createClubDiscussion(props.clubId, {
      title: title.value.trim(),
      body: body.value.trim(),
    })
    if (d?.id) {
      toast.success('Discussion created!')
      emit('created', d.id)
    }
  } catch (err: any) {
    const msg = err?.response?.data?.message || err.message || 'Failed to create discussion'
    errorMessage.value = msg
    toast.apiError(err, 'Failed to create discussion')
  } finally {
    creating.value = false
  }
}
</script>
