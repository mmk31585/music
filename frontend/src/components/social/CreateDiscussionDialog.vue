<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
      @click.self="emit('close')"
    >
      <div class="glass-strong mx-4 w-full max-w-lg rounded-2xl p-8">
        <h2 class="mb-6 text-xl font-bold text-white">بحث جدید</h2>

        <div class="space-y-4">
          <input
            v-model="title"
            type="text"
            placeholder="عنوان بحث"
            class="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder-white/20 outline-none transition focus:border-white/20"
            dir="rtl"
          />
          <textarea
            v-model="body"
            placeholder="متن بحث..."
            rows="5"
            class="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder-white/20 outline-none transition focus:border-white/20"
            dir="rtl"
          />
        </div>

        <div class="mt-6 flex gap-3">
          <button
            class="flex-1 rounded-xl bg-white/5 py-3 text-sm font-medium text-white/50 transition hover:bg-white/10"
            @click="emit('close')"
          >
            انصراف
          </button>
          <button
            class="flex-1 rounded-xl bg-[#1db954] py-3 text-sm font-bold text-black transition hover:bg-[#1db954]/90 disabled:opacity-40"
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

const props = defineProps<{
  visible: boolean
  clubId: string
}>()
const emit = defineEmits<{
  close: []
  created: [discussionId: string]
}>()

const socialApi = useSocialApi()

const title = ref('')
const body = ref('')
const creating = ref(false)

async function handleCreate() {
  if (!title.value.trim() || !body.value.trim() || creating.value) return
  creating.value = true
  try {
    const d = await socialApi.createClubDiscussion(props.clubId, {
      title: title.value.trim(),
      body: body.value.trim(),
    })
    if (d?.id) {
      emit('created', d.id)
    }
  } catch {
    // ignore
  } finally {
    creating.value = false
  }
}
</script>
