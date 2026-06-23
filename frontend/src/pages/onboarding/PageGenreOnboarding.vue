<template>
  <div
    class="flex min-h-screen items-center justify-center bg-gradient-to-br from-[#0a0a0a] via-[#0f0f1a] to-[#0a0a0a] px-4"
  >
    <div
      class="w-full max-w-xl rounded-3xl border border-white/10 bg-white/5 p-8 shadow-2xl backdrop-blur-xl"
    >
      <div class="mb-8 text-center">
        <div
          class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-[#1db954]/20"
        >
          <i aria-hidden="true" class="pi pi-music text-3xl text-[#1db954]" />
        </div>
        <h1 class="text-2xl font-bold text-white md:text-3xl">چه سبک موسیقی رو دوست داری؟</h1>
        <p class="mt-2 text-sm text-slate-400">
          حداقل ۳ تا انتخاب کن تا بتونیم پیشنهادهای بهتری برات داشته باشیم
        </p>
      </div>

      <!-- Loading skeleton -->
      <div v-if="loading" class="mb-6 flex flex-wrap justify-center gap-3">
        <div
          v-for="n in 12"
          :key="n"
          class="animate-pulse rounded-xl border border-white/10 bg-white/[0.04] px-5 py-2.5 text-sm"
        >
          <span class="text-transparent">Loading</span>
        </div>
      </div>

      <!-- Genre chips -->
      <div
        v-else-if="genres.length > 0"
        class="mb-6 flex flex-wrap justify-center gap-3"
      >
        <button
          v-for="genre in genres"
          :key="genre.id"
          type="button"
          class="spring rounded-xl border px-5 py-2.5 text-sm font-medium transition-all"
          :class="
            selectedIds.includes(genre.id)
              ? 'border-[#1db954] bg-[#1db954]/15 text-[#1db954] shadow-[0_0_12px_rgba(29,185,84,0.15)]'
              : 'border-white/10 bg-white/[0.04] text-slate-300 hover:border-white/20 hover:bg-white/[0.08] hover:text-white'
          "
          @click="toggleGenre(genre.id)"
        >
          {{ genre.name }}
        </button>
      </div>

      <!-- Empty state -->
      <div v-else-if="!loading" class="mb-6 text-center text-sm text-slate-500">
        <p>در حال حاضر سبکی برای نمایش وجود ندارد</p>
      </div>

      <div v-if="apiError" class="mb-4 text-center text-sm text-red-400" role="alert">
        {{ apiError }}
      </div>

      <div class="flex flex-col items-center gap-3">
        <Button
          label="ادامه"
          icon="pi pi-arrow-left"
          icon-class="ml-2"
          :loading="saving"
          :disabled="selectedIds.length < 3 || saving"
          class="w-full max-w-xs border-0 bg-[#1db954] px-6 text-black hover:bg-[#1ed760]"
          @click="saveGenres"
        />
        <button
          type="button"
          class="text-xs text-slate-500 transition-colors hover:text-slate-300"
          @click="skip"
        >
          بعداً انجام می‌دم
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import Button from 'primevue/button'
import { useGenresApi } from '@/services/api/catalog/genres'
import { client } from '@/composables'

interface GenreItem {
  id: string
  name: string
}

const router = useRouter()
const genresApi = useGenresApi()

const genres = ref<GenreItem[]>([])
const selectedIds = ref<string[]>([])
const saving = ref(false)
const loading = ref(true)
const apiError = ref('')

onMounted(async () => {
  try {
    const data = await genresApi.getGenres()
    genres.value = (data as unknown as GenreItem[]) || []
  } catch {
    apiError.value = 'بارگیری سبک‌ها با مشکل مواجه شد. بعداً تلاش کن.'
  } finally {
    loading.value = false
  }
})

function toggleGenre(id: string) {
  const idx = selectedIds.value.indexOf(id)
  if (idx >= 0) {
    selectedIds.value.splice(idx, 1)
  } else {
    selectedIds.value.push(id)
  }
}

async function saveGenres() {
  if (selectedIds.value.length < 3) return
  saving.value = true
  apiError.value = ''

  try {
    await client.post('/onboarding/genres', { genre_ids: selectedIds.value })
    await router.replace({ name: 'app.home' })
  } catch {
    apiError.value = 'ذخیره‌سازی با مشکل مواجه شد. دوباره تلاش کن.'
  } finally {
    saving.value = false
  }
}

function skip() {
  router.replace({ name: 'app.home' })
}
</script>

<style scoped>
.spring {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
