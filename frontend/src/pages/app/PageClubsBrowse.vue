<template>
  <div class="mx-auto max-w-5xl space-y-8 px-4 pt-20 pb-24 md:px-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-white">کلاب‌ها</h1>
        <p class="mt-1 text-sm text-white/40">انجمن‌های موسیقی رو کشف کن</p>
      </div>
      <button
        class="inline-flex items-center gap-1.5 rounded-xl bg-spotify px-5 py-2.5 text-sm font-bold text-black transition hover:bg-spotify-hover"
        @click="showCreateDialog = true"
      >
        ساخت کلاب جدید
      </button>
    </div>

    <!-- Genre Filter Chips -->
    <div class="flex flex-wrap gap-2">
      <button
        v-for="g in genres"
        :key="g.value"
        class="rounded-full px-4 py-1.5 text-xs font-medium transition"
        :class="selectedGenre === g.value
          ? 'bg-white/15 text-white'
          : 'bg-white/4 text-white/40 hover:bg-white/8 hover:text-white/60'"
        @click="selectedGenre = g.value"
      >
        {{ g.label }}
      </button>
    </div>

    <SkeletonLoader v-if="loading" variant="card" class="h-48" />

    <div v-else-if="!clubs.length" class="flex flex-col items-center gap-4 rounded-2xl border border-dashed border-white/6 py-20 text-center" role="status">
      <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/4">
        <i aria-hidden="true" class="pi pi-building text-2xl text-slate-500" />
      </div>
      <p class="text-sm font-medium text-white/40">هنوز کلابی با این فیلتر وجود نداره</p>
    </div>

    <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <MusicClubCard
        v-for="club in clubs"
        :key="club.id"
        :club="club"
        @click="router.push({ name: 'social.club', params: { id: club.id } })"
      />
    </div>

    <CreateClubDialog
      :visible="showCreateDialog"
      @close="showCreateDialog = false"
      @created="onClubCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import MusicClubCard from '@/components/social/MusicClubCard.vue'
import CreateClubDialog from '@/components/social/CreateClubDialog.vue'
import { useSocialApi } from '@/services/api/social'
import type { MusicClub } from '@/services/api/social'

const router = useRouter()
const api = useSocialApi()

const loading = ref(true)
const clubs = ref<MusicClub[]>([])
const selectedGenre = ref('')
const showCreateDialog = ref(false)

const genres = [
  { value: '', label: 'همه' },
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

async function loadClubs() {
  loading.value = true
  try {
    const params: { limit?: number; genre?: string } = { limit: 50 }
    if (selectedGenre.value) params.genre = selectedGenre.value
    const res = await api.listClubsWithGenre(params)
    clubs.value = Array.isArray(res) ? res : []
  } catch (err) {
    console.error('Failed to load clubs:', err)
    clubs.value = []
  } finally {
    loading.value = false
  }
}

function onClubCreated(clubId: string) {
  showCreateDialog.value = false
  router.push({ name: 'social.club', params: { id: clubId } })
}

watch(selectedGenre, () => loadClubs())

onMounted(loadClubs)
</script>
