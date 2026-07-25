<template>
  <div
    class="flex min-h-screen items-center justify-center bg-linear-to-br from-surface-base via-surface-overlay to-surface-base px-4"
  >
    <div
      class="w-full max-w-xl rounded-3xl border border-white/10 bg-white/5 p-8 shadow-2xl backdrop-blur-xl"
    >
      <div class="mb-8 text-center">
        <div
          class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-spotify/20"
        >
          <Music aria-hidden="true" class="text-3xl text-spotify"  />
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
          class="animate-pulse rounded-xl border border-white/10 bg-white/4 px-5 py-2.5 text-sm"
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
              ? 'border-spotify bg-spotify/15 text-spotify shadow-[0_0_12px] shadow-spotify/15'
              : 'border-white/10 bg-white/4 text-slate-300 hover:border-white/20 hover:bg-white/8 hover:text-white'
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
          class="w-full max-w-xs border-0 bg-spotify px-6 text-black hover:bg-spotify-hover"
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
import { Music } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useGenresApi } from '@/services/api/catalog/genres'
import { client } from '@/composables'

interface GenreItem {
  id: string
  name: string
}

const router = useRouter()
const toast = useToast()
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
    toast.add({
      severity: 'success',
      summary: 'Saved',
      detail: 'Your genre preferences have been saved.',
      life: 3000,
    })
    await router.replace({ name: 'app.home' })
  } catch {
    apiError.value = 'ذخیره‌سازی با مشکل مواجه شد. دوباره تلاش کن.'
  } finally {
    saving.value = false
  }
}

async function skip() {
  // Auto-save with a few popular defaults so recommendations work
  if (genres.value.length > 0) {
    const autoSelect = genres.value.slice(0, 3).map((g) => g.id)
    try {
      await client.post('/onboarding/genres', { genre_ids: autoSelect })
    } catch {
      // Silently ignore — user chose to skip, no need to block navigation
    }
  }
  router.replace({ name: 'app.home' })
}
</script>

<style scoped>
.spring {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
