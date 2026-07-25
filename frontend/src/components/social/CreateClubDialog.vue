<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="fixed inset-0 z-50 flex items-center justify-center bg-bg-overlay/60 backdrop-blur-xs"
      @click.self="emit('close')"
    >
      <div class="glass-strong mx-4 w-full max-w-md rounded-2xl p-8">
        <h2 class="mb-6 text-xl font-bold text-primary">ساخت کلاب جدید</h2>

        <div class="space-y-4">
          <input
            v-model="name"
            type="text"
            placeholder="اسم کلاب"
            aria-label="اسم کلاب"
            autofocus
            class="w-full rounded-xl border border-border-default bg-surface-overlay px-4 py-3 text-sm text-primary placeholder:text-muted outline-hidden transition focus:border-border-strong"
            dir="rtl"
          />

          <textarea
            v-model="description"
            placeholder="توضیحات (اختیاری)"
            rows="3"
            aria-label="توضیحات"
            class="w-full rounded-xl border border-border-default bg-surface-overlay px-4 py-3 text-sm text-primary placeholder:text-muted outline-hidden transition focus:border-border-strong"
            dir="rtl"
          />

          <select
            v-model="genre"
            class="w-full rounded-xl border border-border-default bg-surface-overlay px-4 py-3 text-sm text-primary outline-hidden transition focus:border-border-strong"
            aria-label="دسته‌بندی"
            dir="rtl"
          >
            <option value="" disabled selected>دسته‌بندی</option>
            <option v-for="g in genres" :key="g.value" :value="g.value">
              {{ g.label }}
            </option>
          </select>

          <input
            v-model="coverUrl"
            type="text"
            placeholder="لینک تصویر (اختیاری)"
            aria-label="لینک تصویر"
            class="w-full rounded-xl border border-border-default bg-surface-overlay px-4 py-3 text-sm text-primary placeholder:text-muted outline-hidden transition focus:border-border-strong"
            dir="rtl"
          />
        </div>

        <div class="mt-6 flex gap-3">
          <button
            class="flex-1 rounded-xl bg-surface-overlay py-3 text-sm font-medium text-secondary transition hover:bg-surface-active"
            @click="emit('close')"
          >
            انصراف
          </button>
          <button
            class="flex-1 rounded-xl bg-accent py-3 text-sm font-bold text-black transition hover:bg-accent/90 disabled:opacity-40"
            :disabled="!name.trim() || creating"
            @click="handleCreate"
          >
            {{ creating ? '...' : 'ساختن' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useSocialApi } from '@/services/api/social'

defineProps<{ visible: boolean }>()
const emit = defineEmits<{
  close: []
  created: [clubId: string]
}>()

const api = useSocialApi()

const name = ref('')
const description = ref('')
const genre = ref('')
const coverUrl = ref('')
const creating = ref(false)

const genres = [
  { value: 'pop', label: 'پاپ' },
  { value: 'rap', label: 'رپ' },
  { value: 'traditional', label: 'سنتی' },
  { value: 'rock', label: 'راک' },
  { value: 'electronic', label: 'الکترونیک' },
  { value: 'jazz', label: 'جز' },
  { value: 'classical', label: 'کلاسیک' },
  { value: 'folk', label: 'فولک' },
  { value: 'metal', label: 'متال' },
  { value: 'rnb', label: 'R&B' },
  { value: 'other', label: 'سایر' },
]

async function handleCreate() {
  if (!name.value.trim() || creating.value) return
  creating.value = true
  try {
    const slug = name.value
      .trim()
      .toLowerCase()
      .replace(/[^a-z0-9\u0600-\u06FF\s-]/g, '')
      .replace(/\s+/g, '-')
      .replace(/-+/g, '-')
      .replace(/^-|-$/g, '')

    const club = await api.createClub({
      name: name.value.trim(),
      slug,
      description: description.value.trim() || undefined,
      genre: genre.value || undefined,
      cover_url: coverUrl.value.trim() || undefined,
      is_public: true,
    })

    if (club?.id) {
      emit('created', club.id)
    }
  } catch (err) {
    console.error('Failed to create club:', err)
  } finally {
    creating.value = false
  }
}
</script>
